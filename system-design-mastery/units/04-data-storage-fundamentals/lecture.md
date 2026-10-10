# Unit 4: Data Storage Fundamentals

## Objectives
By the end of this unit you can:
- Choose a datastore family (relational, key-value, document, wide-column, search, object, time-series) from an access pattern and NFRs.
- Explain how B-tree and LSM-tree storage engines differ and what each optimizes.
- Design indexes (composite, covering) and predict their read benefit and write cost.
- Explain ACID and the isolation levels, and identify anomalies (lost update, write skew) and their fixes.
- Model data around query patterns rather than around entities alone.

## Prerequisites
Unit 2 (estimation: you'll size data and IOPS).

## Core content

### 1. Start from access patterns, not technology
Ask: *What are the queries? How often? What consistency? How big?* Then pick the store. Write down the top 5 queries before choosing anything.

| Family | Model | Strengths | Weaknesses | Examples |
|---|---|---|---|---|
| **Relational (RDBMS)** | Tables, joins, SQL, transactions | Flexible ad-hoc queries, strong consistency, constraints | Vertical scaling limits; sharding is manual | Postgres, MySQL |
| **Key-value** | `key → blob` | Very fast, simple, scales horizontally | Only lookups by key | Redis, DynamoDB (KV mode) |
| **Document** | JSON documents by key + secondary indexes | Flexible schema, data locality for aggregates | Joins weak; consistency varies | MongoDB, Firestore |
| **Wide-column** | Partition key + sorted clustering columns | Massive write throughput, predictable scale | Query patterns fixed at design time | Cassandra, HBase, Bigtable |
| **Search** | Inverted index | Full-text, faceting, ranking | Not a system of record; eventual consistency | Elasticsearch, OpenSearch |
| **Time-series** | Timestamped points, downsampling | Compression, range queries | Narrow use | TimescaleDB, InfluxDB, Prometheus |
| **Object store** | Blobs in buckets | Cheap, durable (11 nines), huge | High latency, no partial updates | S3, GCS |
| **Graph** | Nodes/edges | Multi-hop traversals | Scaling, niche | Neo4j |

**Default for most products: Postgres.** It handles a surprising range (JSON, full-text, geo, queues) and gives transactions. Choose something else when a *number* forces you (write volume, data size, latency, access pattern), not for fashion.

### 2. Storage engines: B-tree vs LSM-tree
Both structures turn "find by key" into fewer disk reads; they differ in write/read trade-offs.

- **B-tree** (Postgres, MySQL/InnoDB): data in fixed-size pages, updated **in place**. Reads: O(log n) page reads (typically 3-4 levels, top cached). Writes: random page updates + write-ahead log (WAL) for crash safety. Predictable reads, good for range scans, but **write amplification** and random I/O.
- **LSM-tree** (Cassandra, RocksDB, LevelDB, HBase): writes go to an in-memory **memtable** (+ WAL), flushed as immutable sorted files (**SSTables**); background **compaction** merges files. Writes are sequential and very fast; reads may check memtable + several SSTables (mitigated by **Bloom filters** and indexes); compaction burns I/O and can cause latency spikes.

Rule of thumb: **write-heavy, append-like workloads → LSM; read-heavy, transactional, range/lookups → B-tree.**

### 3. Indexes
An index is a separate sorted structure that trades **write cost + space for read speed**.
- **Composite index** `(user_id, created_at)` serves `WHERE user_id=? ORDER BY created_at` and `WHERE user_id=?` (**leftmost-prefix rule**), but not `WHERE created_at=?` alone.
- **Covering index** includes all columns the query needs, so the DB answers from the index without touching the table.
- **Selectivity**: an index on a boolean/low-cardinality column rarely helps; the planner scans anyway.
- **Cost**: each write updates every index; 6 indexes ≈ 7 writes per insert. Too many indexes hurt write throughput and memory. Check with `EXPLAIN ANALYZE`.
- Range queries and sorting benefit from B-tree order; hash indexes only serve equality.

### 4. Transactions: ACID
- **Atomicity**: all or nothing. **Consistency**: constraints hold (application-defined invariants). **Isolation**: concurrent transactions don't see each other's partial work (to a degree). **Durability**: committed = survives crash (WAL fsync, replication).

**Isolation levels & anomalies** (SQL standard; real engines differ):

| Level | Dirty read | Non-repeatable read | Phantom | Notes |
|---|---|---|---|---|
| Read Uncommitted | possible | possible | possible | rarely used |
| **Read Committed** (Postgres default) | no | possible | possible | each statement sees committed data as of its start |
| **Repeatable Read / Snapshot** (MySQL default, Postgres RR) | no | no | mostly prevented (Postgres RR uses snapshot) | MVCC snapshot at txn start; **write skew still possible** |
| **Serializable** | no | no | no | Behaves as if transactions ran one at a time; costs aborts/retries or locks |

Anomalies you must recognize:
- **Lost update**: two transactions read-modify-write the same row; one overwrites the other. Fix: atomic `UPDATE t SET n = n + 1`, `SELECT … FOR UPDATE`, optimistic version column (`WHERE version = $v`), or Serializable.
- **Write skew**: two transactions read overlapping data, then write *different* rows, jointly violating an invariant. Classic: two doctors both check "at least 2 on call" then both go off-call. Snapshot isolation allows this. Fix: Serializable, or lock the checked rows (`FOR UPDATE`) / materialize a conflict row / use a constraint.
- **Phantom**: a re-run range query returns new rows inserted by another txn.

**MVCC** (multi-version concurrency control): readers don't block writers; each transaction sees a snapshot of row versions. Cost: old versions must be vacuumed (Postgres `VACUUM`).

```ts
// Lost update: read-modify-write in app code (WRONG under concurrency)
const { balance } = (await db.query("SELECT balance FROM acct WHERE id=$1", [id])).rows[0];
await db.query("UPDATE acct SET balance=$1 WHERE id=$2", [balance - 50, id]); // overwrites concurrent change

// Fix 1: atomic update with guard
await db.query("UPDATE acct SET balance = balance - 50 WHERE id=$1 AND balance >= 50", [id]);

// Fix 2: optimistic concurrency
const r = await db.query(
  "UPDATE doc SET body=$1, version=version+1 WHERE id=$2 AND version=$3", [body, id, seenVersion]);
if (r.rowCount === 0) throw new Error("conflict, reload and retry");
```

### 5. Data modeling
- **Relational**: normalize to avoid anomalies (3NF), **denormalize deliberately** for read performance (copy the author name into each comment row) and accept update complexity.
- **NoSQL**: model **per query**. In Cassandra/DynamoDB you design tables around access patterns (`messages_by_channel(channel_id, ts)`), duplicating data for each query shape. Consistency across copies is your problem.
- **Choose the partition/primary key deliberately**: it determines locality, hot spots, and query possibilities (Unit 6).
- **Store blobs elsewhere**: metadata in the DB, files in object storage referenced by key.
- **Schema evolution**: adding nullable columns is cheap; rewriting large tables locks or runs long; plan expand → migrate → contract steps.

## Worked examples

### Straightforward: pick a datastore
*A user-profile service: 20M users, lookups by user_id at 30k QPS, occasional updates, need transactions when changing email + credentials.* → **Postgres** (system of record, relational integrity, transactions) with a Redis cache for hot reads (Unit 5). A document/KV store could serve lookups but the transactional credentials update and unique email constraint favor relational.

### Typical: index design
*Query: "the last 20 orders for a customer, showing id, total, status", `orders` has 500M rows.*
- Without an index: full scan (~minutes).
- `INDEX(customer_id)`: finds the customer's rows, but a sort by `created_at` is still needed and each row requires a table fetch.
- Best: **`INDEX(customer_id, created_at DESC) INCLUDE (total, status)`**: seek to the customer, read the first 20 entries in order, answer from the index alone (covering). Cost: extra ~40 bytes per row × 500M ≈ 20 GB of index, and slower inserts. Justified because this query is on the hot path.

### Edge case: write skew on-call rota
*Rule: at least one doctor must stay on call. Alice and Bob are on call. Both run in parallel at snapshot isolation: `SELECT count(*) WHERE on_call` → 2; each sets themselves off-call.* Each transaction updates a different row, so no write-write conflict is detected; both commit; zero doctors on call. Fixes: run at **Serializable** (one txn aborts and retries); or `SELECT … FROM shifts WHERE on_call FOR UPDATE` to lock the rows read; or a single `oncall_count` row updated with a `CHECK (count >= 1)`. Lesson: "repeatable read" is not "safe"; think about invariants that span rows.

## Common misconceptions

- **"NoSQL is faster/more scalable than SQL."** Different trade-offs. A well-tuned Postgres handles enormous workloads; NoSQL wins when access patterns are simple key-based and you need horizontal write scale.
- **"More indexes = faster."** Every index taxes writes and memory; unused indexes are pure cost.
- **"ACID means the DB is always consistent."** The C in ACID is your invariants; isolation is configurable, and the default level allows anomalies.
- **"Repeatable Read prevents all anomalies."** Snapshot-based levels still allow write skew.
- **"Transactions make read-modify-write in app code safe."** Not without locks, atomic statements, or serializable isolation.
- **"LSM-trees are always faster."** Faster writes; reads and compaction can cost more; space amplification and latency spikes.
- **"JSON columns mean I don't need a schema."** You move schema enforcement into application code, and every reader has to handle every historical shape.

## Check your understanding

1. Index on `(a, b, c)`. Which of these use it efficiently: `WHERE a=1`, `WHERE b=2`, `WHERE a=1 AND b=2`, `WHERE b=2 AND c=3`? *(First and third; b and c without a can't use the leftmost prefix.)*
2. Why are LSM writes fast? *(Sequential appends to a memtable/WAL and immutable files; compaction is deferred to background.)*
3. Two users both buy the last concert ticket at snapshot isolation using "SELECT count; if >0 then UPDATE another row". What anomaly? *(Write skew/lost invariant; fix with a row lock or serializable.)*
4. Why does a covering index help? *(The query is answered from the index alone, avoiding random table page fetches.)*

## Summary
- Choose stores from access patterns and numbers; Postgres is a strong default.
- B-tree: in-place, read-optimized, transactional; LSM: append + compaction, write-optimized.
- Indexes speed reads and tax writes; use composite/covering indexes and verify with EXPLAIN.
- Know isolation levels: Read Committed and Snapshot still permit lost updates/write skew; use atomic updates, locks, optimistic versioning, or Serializable.
- Model data around queries; keep blobs in object storage; plan schema migrations.
