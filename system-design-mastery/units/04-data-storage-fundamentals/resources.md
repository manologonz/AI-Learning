# Further resources: Unit 4

- **Designing Data-Intensive Applications, Ch. 3 (Storage and Retrieval) and Ch. 7 (Transactions)** (Kleppmann, book) — the best explanation of B-tree vs LSM and of isolation anomalies including write skew; our lecture is a compressed version.
- **Use The Index, Luke** (free online book, use-the-index-luke.com) — practical SQL indexing (composite order, covering indexes, why the planner ignores an index), directly applicable to production tuning.
- **Jepsen: "Consistency Models" map** (jepsen.io/consistency) — an authoritative visual map of isolation/consistency models; use it as a reference when you meet a database's marketing claims.
- **PostgreSQL docs: Chapter 13 "Concurrency Control"** (docs, postgresql.org) — the primary source on MVCC and what Read Committed / Repeatable Read / Serializable actually do in Postgres.
- **Hands-on sandbox:** run Postgres locally (Docker) and reproduce the lost update and write skew examples with two `psql` sessions; watching them happen is far more durable than reading.
