# Further resources: Unit 6

- **Designing Data-Intensive Applications, Ch. 5 (Replication) and Ch. 6 (Partitioning)** (Kleppmann, book) — the canonical treatment of lag anomalies, quorums, and partitioning; read alongside this unit.
- **"Dynamo: Amazon's Highly Available Key-value Store"** (paper, SOSP 2007) — the origin of leaderless replication, consistent hashing with virtual nodes, hinted handoff, and read repair; very readable.
- **Cassandra docs: "Dynamo" and "Data Modeling"** (docs, cassandra.apache.org) — shows quorum settings and partition/clustering key design in a production system.
- **Instagram/Discord engineering posts on sharding** (e.g., Discord's "How Discord Stores Trillions of Messages") — real-world shard-key decisions, hot partitions, and migrations, a good preview of Unit 13.
- **Patroni documentation** (docs, patroni.readthedocs.io) — a concrete, open-source implementation of leader election and failover for Postgres, useful for grounding the failover discussion.
