# Unit 14 Quiz: Answer Key

**Q1: B.**
- ADRs are short, single-decision records. A is a spec; C/D unrelated.

**Q2: B.**
- Alternatives and trade-offs expose reasoning and let reviewers challenge it. A/C/D add little to evaluation.

**Q3: B.**
- A shard key is costly/painful to change (data migration, downtime risk). A, C, D are cheap to reverse (two-way doors).

**Q4. Sample answer.**
- Summary and impact (customer/SLO); timeline; root causes and contributing factors; what went well/where we got lucky; action items with owners and due dates (prevent/detect/mitigate); lessons learned; follow-up tracking. Blameless language throughout.

**Q5: B.**
- Mitigate first (restore service), diagnose after. A can wait; C after resolution; D contradicts the blameless approach.

**Q6. Sample answer.**
- Latency percentiles/SLO attainment, error rates and incident history, cost, deploy frequency/lead time/change failure rate/MTTR (DORA), saturation/utilization of key resources, on-call load, onboarding time, dependency graph/ownership.

**Q7. Sample ADR.**
- **Title:** ADR-0021: Use Kafka for inter-service domain events.
- **Status:** Proposed.
- **Context:** 12 services communicate via synchronous HTTP chains; recurring cascading failures (INC-x, INC-y); consumers need asynchronous, replayable events; team has some Kafka experience.
- **Decision:** Publish domain events (via transactional outbox) to Kafka topics keyed by aggregate ID; consumers are idempotent; keep synchronous calls only where the caller needs an immediate response.
- **Consequences:** + decoupling and resilience, replay, multiple consumers; − eventual consistency, operational burden (cluster, schema registry), need for idempotency/DLQs, debugging complexity.
- **Alternatives:** keep HTTP with circuit breakers/timeouts (rejected as the sole fix: temporal coupling remains), RabbitMQ/SQS (no replay/ordering per key), managed streaming (Kinesis/MSK) as the hosting option.
- **Revisit if:** operating cost or complexity outweighs benefit, e.g., event volume stays < 100/s and a managed queue suffices.

**Q8. Sample answer.**
- Stand up the new DB; backfill from a snapshot in throttled, idempotent batches; keep in sync via CDC/dual-write (with outbox) so changes during backfill land; **shadow reads** comparing results (sampled) and verify with row counts/checksums; migrate reads gradually (1% → 100%) behind a flag; then migrate writes (writer switches with reverse sync from new to old for rollback); monitor SLOs; after a soak period stop reverse sync and decommission the old table. Rollback at each step = flip the flag back (old system still current until the write cutover, then via reverse sync).

**Q9: B.**
- Identity is a commodity and security-critical; buy/adopt unless it's the differentiator; evaluate TCO, lock-in/exit, SLA vs your SLO, compliance, and failure planning. A ignores TCO/opportunity cost; C/D are unreasoned or unsafe.

**Q10. Sample answer.**
- Questions: (1) What specific pain (measurable: incidents, lead time, cost) will the rewrite fix, and by how much? (2) What's the plan for parity/hidden requirements and data migration; how long until value is delivered; what's the rollback plan? (3) Why can't we fix it incrementally?
- Alternative: **strangler fig** (Unit 11): identify the worst-pain seams, extract/refactor incrementally with shadow traffic and flags, delivering value continuously; pay debt with tracked, outcome-linked work.
