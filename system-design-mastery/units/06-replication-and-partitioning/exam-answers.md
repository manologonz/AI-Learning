# Unit 6 Exam: Answer Key

**Q1.**
- Replicas copy every write; they offload *reads* but every replica must also apply all 2M writes/s (replicas can even fall behind). The bottleneck is write throughput and storage on one primary (30 TB, +1 TB/week). Replication provides read scaling, failover, and durability, not write or storage scaling; you need partitioning.

**Q2.**
- Store: shard by `channel_id` (hash) so a channel's messages live together. Primary key `((channel_id, bucket), message_id DESC)` where `bucket` = time window (e.g., week) and `message_id` is time-sortable (ULID/Snowflake).
- Fit: "latest 50" reads the newest bucket of one partition in clustering order; scroll-back walks earlier buckets; all single-partition.
- Expensive queries: anything cross-channel: "all messages by user", global search, per-user unread counts across channels (need separate denormalized tables/indexes or a search system).

**Q3.**
- Hot spot 1: **mega-channels** (5M members) concentrate reads/writes on one partition. Mitigation: sub-shard hot channels `(channel_id, shard = hash(message_id) % k)` with merged reads, cache the latest messages, and serve fan-out reads from replicas/CDN-like caches.
- Hot spot 2: **time bucket**: within one channel the newest bucket takes all writes; if a channel is huge, one partition grows/hots. Mitigation: shorter buckets (hourly) for high-volume channels, and the sub-shard above. (Time isn't a hot spot *across* channels because channel_id is hashed.)

**Q4.**
- (a) `hash % N`: going 20→30, a key stays only if `h mod 20 == h mod 30`, which holds for 1/3 of keys, so **~67% of data moves**; drawback: massive data movement/cache flush.
- (b) Consistent hashing with vnodes: about **1/3 of the data** moves (ideal: 10/30, only to the new nodes, none between existing ones); drawback: vnode bookkeeping, imbalance risk if vnodes are too few.
- (c) Fixed 1,024 partitions: move ~1/3 of *partitions* (~341) to new nodes; drawback: partition count is a hard ceiling on parallelism/granularity and must be chosen up front (too few = can't scale beyond; too many = overhead).
- Note: the key insight is (a) moves ~2/3 vs (b),(c) ~1/3 and only to new nodes.

**Q5.**
- **Consistent prefix violation**: the reply and original are in different partitions (or read from different replicas) that replicate independently, so a reader observes the effect before the cause.
- Prevention: keep a conversation's messages in the same partition (shard by `channel_id`, which we did), assign monotonically ordered IDs per channel, and have clients render by sequence number and buffer/gap-fill (hold a reply until its parent arrives), or read from a single replica per channel.

**Q6.**
- **Best: region-pinned leaders per shard (single-leader per channel, home region).** Each channel has a home region; writes to that channel go to its leader; other regions read replicas (async) and forward writes cross-region for the minority of channels with remote members. It avoids conflicts entirely and preserves per-channel ordering.
- Multi-leader would need conflict resolution and complicates ordering (chat requires a total order per channel); leaderless (quorum) gives no strong ordering and is awkward for ordered logs.
- Trade-off: users in a channel whose home is remote pay cross-region write latency (~100-200 ms); mitigate by homing a channel in the region of most of its active members and migrating homes based on traffic.

**Q7: B.**
- W+R=3 = N, not > N, so a read may hit the one replica that missed the write. A false; C false (W=2 tolerates one failure for writes); D false (quorum isn't linearizability).

**Q8: B.**
- Semi-sync: an acknowledged write exists on the leader and at least one follower; promoting the synchronous follower loses no acknowledged writes; promoting an async follower can lose recent writes. A overclaims; C exaggerates; D nonsense.

**Q9.**
- The query waits for the *slowest* shard. P(at least one of 100 exceeds its p99) = 1 − 0.99^100 ≈ **63%**. So the query's median experiences a shard's p99.
- Mitigations: hedged requests (send a duplicate after a short delay and use the first response), partial results with timeouts, reducing fan-out (better shard key/global index), caching, or precomputation.

**Q10. Sample answer.**
- **Detection:** health from a consensus-backed manager (Patroni + etcd, or managed service): several missed heartbeats over ~10-30 s (avoid flapping); require quorum agreement, not a single observer.
- **Candidate:** most advanced replica by LSN (least data loss); prefer the sync replica if configured.
- **Fencing:** ensure the old primary can't accept writes: revoke its lease/leader key in the consensus store, STONITH/power off, or block at the network/proxy layer; on return it must rejoin as a replica (pg_rewind).
- **Redirect:** clients connect via a proxy/VIP/DNS with a low TTL that follows the leader key; connection pools reconnect with jitter.
- **Guarantees:** *Async:* RPO > 0 (typically seconds of writes could be lost; report the replication lag), RTO ~30-60 s. *Sync (one follower):* RPO = 0 for acknowledged writes, but a follower outage can stall writes unless you fall back to async; RTO similar. Communicate the trade-off with business numbers (cost of losing X seconds of orders vs write latency +1-2 ms).
