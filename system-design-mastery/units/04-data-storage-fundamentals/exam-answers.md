# Unit 4 Exam: Answer Key

**Q1.**
- Interleaving: T1 reads seat 17 = 'available'; T2 reads seat 17 = 'available'; T1 updates to 'sold' and commits; T2 updates to 'sold' (its earlier check is stale) and commits. Both customers get the seat. Read Committed doesn't prevent it because the read and write are separate statements.
- Fix A (atomic guarded update): `UPDATE seats SET status='sold' WHERE id=$1 AND status='available'` and check `rowCount = 1`; the DB serializes the row update, the second sees 'sold' and updates 0 rows.
- Fix B (pessimistic lock): `BEGIN; SELECT status FROM seats WHERE id=$1 FOR UPDATE;` then check/update and commit. The second transaction blocks until the first commits, then sees 'sold'.
- Also valid: optimistic `version` column check; Serializable with retry.

**Q2.**
- One transaction: lock all requested seats in a **consistent order** (sort by seat id) using `SELECT … WHERE id = ANY($1) ORDER BY id FOR UPDATE`, verify all 'available' (or hold expired), update them all to 'sold', insert the order and order_items, commit. If any seat isn't available, roll back (all or none).
- Deadlock avoidance: acquire locks in the same global order (ascending id), keep transactions short (no external calls like payment inside the lock; use a *hold* step with `held_until` first), and retry on deadlock detection.
- Alternative: single guarded statement `UPDATE seats SET status='held' WHERE id = ANY($1) AND status='available'` and verify `rowCount = array length`, otherwise roll back.
- Isolation: Read Committed with explicit locking suffices; Serializable also works with retries.

**Q3.**
- `CREATE INDEX ON orders (user_id, created_at DESC) INCLUDE (event_id, total);`
- Read benefit: seek to the user's block, read the first N in order, no sort, no table fetches (covering).
- Write cost: every order insert also writes an index entry; adds storage and slightly slower writes (orders are low volume relative to browse, so acceptable).
- Size: 500M × 40 B = 20 GB (plus overhead → ~25-30 GB); fits on one node but should be mostly cached in memory for hot users.

**Q4.**
- Lives in a **search engine** (Elasticsearch/OpenSearch) as a derived read model; Postgres remains the source of truth.
- Sync: change data capture (CDC, e.g., Debezium reading the WAL) or the transactional **outbox** pattern publishing events to a queue that indexers consume; a periodic reindex/backfill job to repair drift.
- Consistency: eventually consistent (seconds of lag): a new event may not appear immediately; availability/prices/seat status must *not* come from search results; the detail page reads from the DB (or cache). Never use the search index for purchase decisions.

**Q5.** Sample answer (grade on breadth + measurement + a justified pick).
- Measure first: is it write CPU from index maintenance, lock contention on hot rows, connection saturation, WAL/fsync, or bad queries? (`pg_stat_statements`, lock waits, `EXPLAIN`.)
- Options: (1) **Move seat holds off Postgres** into Redis with atomic ops (`SET seat:17 user NX EX 300`), persisting to Postgres only on purchase; (2) **batch/queue** writes and virtual waiting room to smooth the spike (admit N users at a time); (3) reduce index/write amplification (fewer indexes on `seats`, fillfactor/HOT updates), tune connection pooling (PgBouncer) to cut contention; (4) **shard by event_id** so hot events land on separate nodes; (5) scale up hardware as the immediate lever.
- Pick example: waiting room + Redis holds. Trade-off: added infrastructure and a second system where holds may diverge from Postgres; acceptable because holds are ephemeral and losing them just re-opens a seat. Reversal: if Postgres handles the peak after tuning (CPU < 60%), drop Redis to reduce complexity.

**Q6: B.**
- A is wrong: MVCC's point is that readers don't block writers. C: isolation levels still define visibility rules. D: it keeps multiple versions (that's why vacuum is needed).

**Q7: B.**
- NoSQL wide-column stores model per query; you denormalize into a second table and write to both (via batch/idempotent writes), owning consistency.
- A: cross-partition scans (scatter-gather) don't scale. C: no joins in this model. D: giant partitions are hot spots and have size limits.

**Q8: B.**
- Low selectivity: the planner correctly prefers a sequential scan since an index lookup on 95% of rows costs more than scanning. Consider a **partial index** for the rare values (`WHERE status <> 'active'`) if those are the queries that matter.
- A: an index wouldn't help. C/D are nonsense.

**Q9. Sample answer.**
- Compaction consumes disk and CPU bandwidth in bursts and competes with foreground reads/writes; it can also leave many SSTables to check when it falls behind (read amplification), and write stalls occur if memtable flushes back up.
- Mitigations: throttle/rate-limit compaction, choose a strategy suited to the workload (e.g., time-window compaction for time-series, leveled for read-heavy), provision I/O headroom (NVMe), spread compaction across nodes, and monitor pending compactions. Bloom filters reduce read amplification.

**Q10. Sample answer.**
- Store: wide-column LSM (Cassandra/ScyllaDB/DynamoDB): 50k writes/s, 10B messages, range reads by channel.
- Table: `PRIMARY KEY ((channel_id, bucket), message_id DESC)` where `bucket` is a time window (e.g., month or day) and `message_id` is a time-sortable ID (ULID/Snowflake). Partition = channel+bucket; clustering order = newest first; "last 50" reads one partition; scroll-back walks to older buckets.
- Hot-partition risk: a huge/very active channel concentrating writes on one partition → partitions grow unbounded/hot. Mitigation: smaller buckets (hourly) for high-traffic channels, or add a shard suffix `(channel_id, bucket, shard)` and merge on read; cap partition size (~100MB).
