# Further resources: Unit 8

- **Designing Data-Intensive Applications, Ch. 11 (Stream Processing)** (Kleppmann, book) — logs vs queues, change data capture, and event sourcing, explained with the trade-offs we used here.
- **Debezium documentation: "Outbox Event Router"** (docs, debezium.io) — a production-grade implementation of the transactional outbox with CDC; reading it shows the details our sketch skipped.
- **Confluent: "Kafka Exactly-Once Semantics" blog series** (articles, confluent.io/blog) — the authoritative explanation of what EOS does and doesn't guarantee; essential to avoid the exactly-once myth.
- **Chris Richardson: microservices.io "Transactional outbox" and "Saga" patterns** (pattern catalog) — concise problem/solution/consequences writeups that map directly to Units 8-9.
- **AWS Builders' Library: "Avoiding insurmountable queue backlogs"** (article, aws.amazon.com/builders-library) — how queue backlogs turn into outages and the design levers (partitioning by priority, shedding, backlog-aware consumers) to prevent them.
