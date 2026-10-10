# Unit 6: Scaling Data: Replication & Partitioning

## Objectives
By the end of this unit you can:
- Design a replication topology (single-leader, multi-leader, leaderless) and state its availability and consistency consequences.
- Diagnose replication-lag anomalies and apply fixes (read-your-writes, monotonic reads).
- Compute quorum settings (N, W, R) and reason about which guarantees they give.
- Choose a partitioning scheme and shard key; explain consistent hashing; anticipate hot partitions and resharding.
- Handle secondary indexes and cross-shard queries.

## Prerequisites
Unit 4 (storage), Unit 5 (caching).

## Core content

### 1. Two tools, two purposes
- **Replication** = copies of the *same* data on multiple nodes → **availability, read scalability, durability, lower latency by locality**. Doesn't increase write capacity (usually every write goes to every copy).
- **Partitioning (sharding)** = splitting *different* data across nodes → **write and storage scalability**.
Large systems use both: each shard is replicated.

Scale in order: optimize queries/indexes → cache → vertical scaling → read replicas → partition. Each step is more complex than the last; take the next only when numbers force it.

### 2. Replication topologies

**Single-leader (primary/replica)**: all writes go to the leader, which ships a log (WAL/binlog) to followers. Reads can go to followers. Simple, no write conflicts. Postgres, MySQL, MongoDB replica sets.
- **Synchronous**: leader waits for follower ack: no data loss on leader crash, but a slow/down follower blocks writes. **Asynchronous**: leader doesn't wait: fast, but recent writes can be lost on failover, and followers lag. **Semi-sync**: wait for one follower; a common compromise.
- **Failover** (promote a follower): detect failure (timeouts are ambiguous: slow vs dead), choose the most up-to-date follower, redirect clients, demote the old leader when it returns. Hazards: **lost writes** (async lag), **split brain** (two nodes think they're leader → divergent data; use fencing/quorum-based leader election, Unit 7), and **premature failover** from a temporary blip.

**Replication lag anomalies** (with async followers):
- **Read-your-writes violation**: user posts, refreshes, reads from a lagging replica, and their post is missing. Fix: read own recent writes from the leader (for N seconds or per-user "last write timestamp/LSN"), or route to replicas that have caught up to the user's LSN.
- **Monotonic reads violation**: a user sees a newer value, then an older one from a different replica. Fix: pin a user to one replica (hash of user id).
- **Consistent prefix violation**: an answer appears before its question (writes to different partitions replicate independently). Fix: keep causally related writes in the same partition or track causality.

**Multi-leader**: several nodes accept writes (one per data center, or per offline device). Lower write latency and datacenter-failure tolerance, but **write conflicts** are inevitable: two leaders update the same record. Resolution: last-write-wins (LWW; silently drops data, depends on clocks), merge (CRDTs, operational transforms), or app-level conflict handling. Avoid unless you need multi-region writes or offline editing.

**Leaderless (Dynamo-style: Cassandra, Riak, DynamoDB internals)**: the client (or coordinator) writes to N replicas and reads from several.
- **Quorum**: with **N** replicas, write to **W**, read from **R**. If **W + R > N**, every read set overlaps the most recent write set → reads see the latest acknowledged write *in the absence of failures/edge cases*. Typical: N=3, W=2, R=2. Tunable: W=1 (fast writes, weak), R=1 (fast reads), W=N (slow writes, tolerate many read failures).
- Repair mechanisms: **read repair** (fix stale replicas when a read notices divergence), **hinted handoff** (store writes for a down node), **anti-entropy** (Merkle-tree background sync).
- Quorums aren't magic: concurrent writes, sloppy quorums, clock-based LWW, partial failures, and failed writes can still yield stale or lost values. They provide *probabilistic* freshness, not linearizability.

### 3. Partitioning (sharding)

Goals: even data distribution, even *load* distribution, minimal cross-shard operations.

**Strategies**
- **Range partitioning** (key ranges: A-F, G-M...): efficient range scans; risk of **hot spots** (time-based keys write to the latest range only). Used by Bigtable/HBase, Spanner.
- **Hash partitioning** (`hash(key) → partition`): spreads load evenly; **destroys range queries** (adjacent keys land everywhere). Compromise: **compound keys** (hash the first part, sort by second: Cassandra `(user_id, timestamp)`).
- **Directory/lookup-based**: a service maps key → shard; flexible (move tenants freely), adds a lookup hop and a SPOF to manage.
- **Geo/tenant partitioning**: by region or customer; aligns with data residency and isolation; watch for one giant tenant.

**Naive `hash % N`** remaps almost all keys when N changes. Better:
- **Fixed number of partitions** (e.g., 1,000 partitions over 10 nodes; add a node by stealing whole partitions). Used by Elasticsearch, Riak, Couchbase.
- **Consistent hashing**: nodes and keys are placed on a ring; each key belongs to the next node clockwise. Adding/removing a node moves only ~1/N of keys. **Virtual nodes** (many ring positions per physical node) even out load and let heterogeneous nodes take weighted shares.
- **Dynamic splitting**: split a partition when it grows too large/hot (Bigtable, MongoDB, CockroachDB).

```ts
// Minimal consistent hash ring with virtual nodes
import { createHash } from "crypto";
const h = (s: string) => parseInt(createHash("md5").update(s).digest("hex").slice(0, 8), 16);

class Ring {
  private ring: { pos: number; node: string }[] = [];
  constructor(nodes: string[], private vnodes = 100) { nodes.forEach(n => this.add(n)); }
  add(node: string) {
    for (let i = 0; i < this.vnodes; i++) this.ring.push({ pos: h(`${node}#${i}`), node });
    this.ring.sort((a, b) => a.pos - b.pos);
  }
  remove(node: string) { this.ring = this.ring.filter(e => e.node !== node); }
  get(key: string): string {
    const p = h(key);
    // binary search for first entry with pos >= p (wraps around)
    let lo = 0, hi = this.ring.length;
    while (lo < hi) { const mid = (lo + hi) >> 1; this.ring[mid].pos < p ? (lo = mid + 1) : (hi = mid); }
    return this.ring[lo % this.ring.length].node;
  }
}
```

**Choosing the shard key (the most consequential decision)**
- High cardinality, even distribution, matches the dominant query (so most queries hit one shard), and doesn't grow monotonically.
- Bad keys: `country` (few values, skewed), `created_at` (all writes hit the latest shard), `status`.
- Good: `user_id`, `tenant_id`, `hash(order_id)`.
- Changing a shard key later is a painful migration.

**Hot partitions/keys**: a celebrity, a viral item, or a huge tenant. Mitigate: split the hot key by appending a random suffix and merging reads (`key#0..9`), dedicated shards for big tenants, caching, and rate limits.

**Secondary indexes on partitioned data**
- **Local (document-partitioned)**: each shard indexes its own rows; a query on the secondary attribute must **scatter-gather** all shards (tail latency ↑). Writes are simple.
- **Global (term-partitioned)**: index partitioned by the indexed value; reads hit one shard, but a write updates indexes on other shards (often async → stale).

**Cross-shard queries & transactions**: joins across shards are expensive; denormalize, co-locate related entities under the same shard key (orders with their user), or use search/analytics stores for broad queries. Multi-shard transactions need 2PC/sagas (Unit 9); avoid designs that need them on the hot path.

**Rebalancing and routing**: moving data while serving traffic (copy, catch up via the log, switch over). Requests find shards via a **routing tier** (proxy), a **partition-aware client**, or **gossip** among nodes; a coordination service (ZooKeeper/etcd) often stores the map.

## Worked examples

### Straightforward: quorum math
*N=5, W=3, R=3: is W+R>N? 6>5 yes. How many node failures can writes and reads tolerate?* Writes need 3 of 5 acks → tolerate 2 failures; reads likewise. If W=5, R=1: fast reads, but a single failure blocks all writes. If W=1, R=1: 2 > 5? No → stale reads are likely.

### Typical: choosing a shard key
*A multi-tenant SaaS with 10k tenants (a few with 30% of the data); queries are almost always per tenant.* Shard by `tenant_id` (hash) so tenant queries hit one shard and transactions stay local. The whales create imbalance → give the top tenants dedicated shards (directory-based override), or sub-partition a whale by `(tenant_id, hash(entity_id))`. Cross-tenant analytics run on a separate warehouse.

### Edge case: the range-shard hot spot and the consistent-hash gotcha
*A metrics table sharded by timestamp range: all current writes land on the newest shard (one node takes 100% of write load while older shards idle).* Fix: prefix the key with a hash bucket (`(bucket = hash(device_id) % 32, timestamp)`) so writes spread across 32 partitions while range scans by device remain efficient. A second gotcha: with consistent hashing *without virtual nodes*, adding a node can move all its load from one neighbor (uneven), and removing a failed node dumps its entire load on one successor → cascading overload. Vnodes spread the removed node's load across many nodes.

## Common misconceptions

- **"Replicas increase write throughput."** Every replica applies every write; reads scale, writes don't.
- **"Sync replication guarantees consistency for readers."** It protects durability of acknowledged writes; reads from other replicas can still be stale depending on the setup.
- **"Quorum reads (W+R>N) are always linearizable."** Sloppy quorums, concurrent writes, and clock-based conflict resolution break that.
- **"Just shard when things get slow."** Sharding is expensive and irreversible-ish; exhaust indexes, caching, and replicas first.
- **"Hash sharding is always better than range."** It kills range scans and ordering; it depends on queries.
- **"Consistent hashing solves hot keys."** It balances *keys*, not *load per key*; one hot key stays hot.
- **"Multi-leader is just faster single-leader."** It introduces conflict resolution, the hardest part.
- **"Failover is automatic and safe."** Timeouts can't distinguish slow from dead; naive failover causes split brain and lost writes.

## Check your understanding

1. N=3, W=1, R=1: can a read miss the latest write? *(Yes: W+R=2 ≤ 3, so read and write sets need not overlap.)*
2. Users say their just-posted comment disappears on refresh. Name the anomaly and a fix. *(Read-your-writes violation from an async replica; read own writes from the leader/up-to-date replica for a short window.)*
3. Why does consistent hashing beat `hash % N` for a cache cluster? *(Adding/removing a node remaps ~1/N of keys instead of nearly all.)*
4. Why is a monotonically increasing ID a poor range-shard key? *(All new writes hit the newest shard: a hot spot.)*

## Summary
- Replication gives availability/read scale/durability; partitioning gives write/storage scale; use both.
- Async replication brings lag anomalies and failover risk; quorums (W+R>N) give tunable, not absolute, consistency; multi-leader means conflict resolution.
- Shard by a high-cardinality key aligned to the main query; range = scans + hot spots, hash = even load + no scans.
- Consistent hashing with virtual nodes limits data movement on resize; fixed/dynamic partitions are common alternatives.
- Plan for hot keys, scatter-gather secondary indexes, cross-shard queries, and online rebalancing.
