# System Design Mastery — Course Syllabus

**Learner starting level:** Mid-level engineer; has used databases, caches, and queues in production but hasn't owned distributed-system design.
**Goal:** Design, defend, and evolve real-world architectures at work (senior-level promotion). Interview fluency is a side benefit.
**Code anchor:** JavaScript/TypeScript sketches (Node) where code helps; otherwise diagrams and back-of-envelope math.
**Estimated total time:** ~50 hours (14 units + capstone)

## Units

1. **The Senior Mindset & Requirements** (~3h) — How seniors frame problems: functional vs non-functional requirements, SLIs/SLOs/SLAs, trade-off thinking, and a repeatable design process.
   - Objectives: turn a vague prompt into requirements and SLOs; state trade-offs explicitly; run a structured design process; identify the "hard part" of a problem.
   - Prerequisites: none
2. **Estimation & Capacity Planning** (~3h) — Back-of-envelope math: QPS, storage, bandwidth, latency numbers.
   - Objectives: estimate QPS/storage/bandwidth from user counts; use latency orders of magnitude; size a fleet; sanity-check a design with numbers.
   - Prerequisites: Unit 1
3. **Networking, APIs & Load Balancing** (~4h) — HTTP/TCP/TLS basics, REST/gRPC/GraphQL, L4 vs L7 balancing, DNS, CDNs, API design (idempotency, pagination, versioning).
   - Objectives: choose an API style; design idempotent and paginated endpoints; pick a load-balancing strategy; explain the request path end-to-end.
   - Prerequisites: Unit 1
4. **Data Storage Fundamentals** (~4h) — Relational vs NoSQL families, indexes, B-trees vs LSM-trees, transactions, ACID, isolation levels.
   - Objectives: pick a datastore for an access pattern; explain index costs; reason about isolation anomalies; model data for query patterns.
   - Prerequisites: Unit 2
5. **Caching** (~3h) — Cache-aside, write-through/back, eviction, invalidation, stampedes, CDN and multi-layer caches.
   - Objectives: choose a caching strategy; predict failure modes (stampede, stale data, hot keys); design invalidation; compute hit-rate impact.
   - Prerequisites: Units 3, 4
6. **Scaling Data: Replication & Partitioning** (~4h) — Leader/follower, multi-leader, leaderless; sharding strategies, consistent hashing, hot partitions, rebalancing.
   - Objectives: design a replication topology; choose a shard key; explain consistent hashing; handle resharding and hot spots.
   - Prerequisites: Unit 4
7. **Distributed Systems Theory** (~4h) — Failure models, time and ordering, CAP/PACELC, consistency models, consensus (Raft/Paxos intuition), quorums.
   - Objectives: distinguish consistency models; apply CAP/PACELC correctly; explain quorum math and leader election; reason about clocks and ordering.
   - Prerequisites: Unit 6
8. **Async Messaging & Event-Driven Architecture** (~4h) — Queues vs logs (SQS/RabbitMQ vs Kafka), delivery semantics, idempotent consumers, outbox pattern, event sourcing, CQRS.
   - Objectives: choose queue vs log; achieve effective exactly-once; implement outbox/idempotency; design retries and DLQs.
   - Prerequisites: Units 6, 7
9. **Distributed Transactions & Data Consistency Patterns** (~3h) — 2PC, sagas, compensation, idempotency keys, deduplication, CRDT overview.
   - Objectives: design a saga for a multi-service workflow; explain 2PC's failure modes; handle partial failure and duplicates.
   - Prerequisites: Units 7, 8
10. **Reliability & Resilience Engineering** (~4h) — Timeouts, retries with backoff/jitter, circuit breakers, bulkheads, rate limiting, load shedding, backpressure, graceful degradation, multi-region and DR.
    - Objectives: prevent retry storms and cascading failures; design rate limiters; define failover/DR strategy with RTO/RPO; apply degradation strategies.
    - Prerequisites: Units 3, 5, 8
11. **Service Architecture: Monolith to Microservices** (~3h) — Modular monoliths, service boundaries, DDD basics, API gateways, service mesh, sync vs async communication, migration strategies.
    - Objectives: decide when (not) to split; draw service boundaries; plan a strangler-fig migration; manage cross-service contracts.
    - Prerequisites: Units 3, 8, 9
12. **Observability, Operations & Security** (~3h) — Metrics/logs/traces, alerting on SLOs, deployment strategies (canary, blue-green), authN/authZ (OAuth/JWT), secrets, multi-tenancy, cost awareness.
    - Objectives: instrument a system; design SLO-based alerts; choose a rollout strategy; design authN/authZ boundaries.
    - Prerequisites: Units 1, 11
13. **Classic Case Studies** (~6h) — Walkthroughs: URL shortener, news feed, chat, rate limiter service, notification system, video/file storage, search/typeahead, ride-hailing/geo, payments ledger.
    - Objectives: apply the full process to unfamiliar prompts; identify each design's core bottleneck; compare alternatives with numbers.
    - Prerequisites: Units 1–12
14. **Architecture in the Real World** (~3h) — ADRs, evolving systems, tech debt, build vs buy, migrations, incident postmortems, communicating designs (diagrams, design docs, review culture), leading design discussions.
    - Objectives: write an ADR and design doc; critique an existing architecture; plan a zero-downtime migration; lead a design review.
    - Prerequisites: Units 1–13
15. **Capstone** (~4h) — Design, document, and defend a complete system end to end (part of the final exam).

## How this course works
Each unit follows a lecture → quiz → exam → resources loop. Lectures have three-tier worked examples and explicit misconceptions. Quizzes are formative; unit exams are harder and synthesis-focused. A cumulative final exam (weighted toward reliability, data scaling, consistency, and case studies) includes a capstone design task. Progress and weak spots are tracked in `progress.md`, and I'll run short spaced-repetition checks on older units at the start of new sessions.

I'll build units incrementally: Unit 1 first, then keep building ahead of you.
