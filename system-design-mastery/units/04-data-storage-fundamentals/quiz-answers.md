# Unit 4 Quiz: Answer Key

**Q1: B (LSM-tree).**
- Writes hit an in-memory memtable and WAL, flush sequentially as SSTables, and merge in the background.
- A: B-trees update pages in place (random I/O). C: hash indexes are lookup structures, not a storage engine strategy for write throughput. D: a heap with no index is fast to append but the option isn't about compaction and makes reads a full scan.

**Q2: B.**
- Equality on the leftmost column plus ordering by the second: seek to user 42, read the first 10 in order.
- A: skips the leftmost column, so the index can't seek by created_at. C: neither column is a usable prefix (`created_at` isn't first; `email` not indexed). D: no email index.

**Q3: C (lost update).**
- Both read the same value, then one write overwrites the other's effect. 
- A: a dirty read means reading uncommitted data. B: phantom is about new rows appearing in a range. D: write skew involves writes to *different* rows violating an invariant; here it's the same row.

**Q4: Sample answer.**
1. Each write must update every index: more I/O, lower write throughput, longer transactions.
2. Indexes consume memory/disk (and cache space that evicts useful data); many are never used by the planner, and low-selectivity indexes give no benefit.
- Also acceptable: more vacuum/maintenance work, slower schema changes.

**Q5: B.**
- Relational DBs provide ACID transactions, constraints (foreign keys, checks), and strong consistency.
- A: caches aren't durable systems of record. C: search indexes are derived, eventually consistent. D: object stores lack multi-row transactions and constraints.

**Q6: B.**
- Read Committed takes a fresh snapshot per statement, so the same row can differ between two statements.
- A is wrong (that's a Repeatable Read guarantee). C is wrong: lost updates are possible with app-level read-modify-write. D is wrong: Serializable is far stricter.

**Q7: `INDEX(customer_id, created_at DESC) INCLUDE (status, total)`.**
- The leading `customer_id` seeks the customer's entries, `created_at DESC` supplies the ordering so no sort is needed and `LIMIT 20` stops after 20 entries, and the included columns make it covering so no table lookups are needed.
- Partial credit: `(customer_id, created_at)` without INCLUDE (still avoids the sort).

**Q8. Sample answer.**
- Two on-call doctors each check that at least two are on call, then each removes *themselves*; both commit under snapshot isolation because they modified different rows; now zero are on call, violating the invariant.
- Fixes: (1) run at Serializable isolation and retry aborts; (2) lock the rows read (`SELECT … FOR UPDATE`) so the second transaction blocks; (3) create a single guard row/counter with a CHECK constraint that both transactions update (turning skew into a write-write conflict).

**Q9: B.**
- Append-heavy, per-device range reads: LSM strengths, partitioned by device for parallelism and locality. 
- A: a single Postgres primary at 200k writes/s with 8 indexes would be overwhelmed. C: search engines aren't systems of record and are expensive on ingest. D: no durability.

**Q10. Sample answer.**
- Questions: (1) What exactly is slow (queries, writes, connections)? Show `EXPLAIN ANALYZE`, slow query log. (2) Are indexes, query patterns, or connection pooling the problem? (3) What are our access patterns and data size (is it a scale problem or a tuning problem)? (4) What do we lose (transactions, joins) and who will run it?
- Justified when: measured limits that tuning/replicas can't fix (e.g., write rate beyond a single node's capability or data too large to shard manually) *and* the access patterns are simple key/document lookups that fit the new model.
