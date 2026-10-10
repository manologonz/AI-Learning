# Further resources: Unit 9

- **Chris Richardson: microservices.io, "Saga" pattern** (pattern write-up, microservices.io/patterns/data/saga.html) — the clearest concise description of orchestration vs choreography and their trade-offs.
- **Hector Garcia-Molina & Kenneth Salem, "Sagas"** (1987 paper) — the original idea, short and still relevant; worth reading to see the motivation (long-lived transactions).
- **Temporal documentation: "Workflows" and "Saga pattern"** (docs, docs.temporal.io) — a real durable-execution engine; shows how crash-resilient orchestrators replay history and why activities must be idempotent.
- **Designing Data-Intensive Applications, Ch. 7 & 9 (Distributed Transactions and Consensus)** (Kleppmann, book) — a rigorous treatment of 2PC's failure modes and why it's problematic in practice.
- **"Coordination Avoidance in Database Systems" (Bailis et al.)** (paper, VLDB 2015) — introduces invariant confluence, the theoretical basis for "coordinate only where invariants demand it".
