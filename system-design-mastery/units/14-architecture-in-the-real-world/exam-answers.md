# Unit 14 Exam: Answer Key

**Q1. Model answer.**
- **Evidence to gather:** deploy pipeline stage timings (why 3 hours: test suite, manual steps, environment provisioning?), change failure causes (which changes fail: schema, config, cross-service?), incident postmortems, DORA metrics trend, DB dependency map (who reads/writes shipments, query patterns), on-call load, cost, ownership map (who owns Kafka?).
- **Structural causes of deploy pain (examples):** (1) shared database and tables read by 5 services → schema changes require lockstep coordination and risk; (2) slow, undifferentiated test/deploy pipeline (no per-service pipelines, weekly batch releases → big-bang changes, 25% failure); (3) no progressive delivery/rollback automation (canary/flags), and unclear ownership (Kafka).
- **CTO plan critique:** a 6-month full rewrite/migration to microservices + K8s + new event platform is a big-bang with high risk, doesn't address the identified root causes directly, adds operational load to already burning-out on-call, and risks feature stagnation. The *causes* (weekly big deploys, shared DB, no ownership/observability) can be addressed incrementally.
- **Don't change:** the nightly finance pipeline (reliable, cheap, low change need); the monolith's stable modules that aren't pain points; existing services that work (leave unless they're a bottleneck). Invest in the parts causing incidents/cost.

**Q2. Design doc outline: "Single owner for shipments data".**
1. **Summary/Decision needed:** approve creating a Shipments service that owns the `shipments` tables and exposes APIs/events.
2. **Context & problem:** 5 services query the 4 TB table directly; schema changes blocked; 2 incidents from lock/queries; no clear owner. Data: query counts per service, slow queries.
3. **Goals/Non-goals:** Goals: single owner, schema changes without cross-team coordination, p99 read < 100 ms, no downtime; Non-goals: splitting other tables, changing UI, replacing Postgres.
4. **Requirements/numbers:** read/write QPS per consumer, freshness needs (which consumers can tolerate eventual), SLO 99.95%.
5. **Proposed design:** Shipments service (API + events via outbox), tables private to it; consumers move to API or read models via events; diagram.
6. **Alternatives:** (a) keep shared DB with a schema-versioning convention (cheap, doesn't fix coupling), (b) DB views/read replicas per consumer (short-term relief), (c) full split into multiple services (too large). Why rejected.
7. **Trade-offs & risks:** latency of API vs direct SQL, consistency (eventual read models), migration effort per consumer, availability dependency on the new service.
8. **Rollout/migration:** per consumer: API parity + shadow reads → flag flip 1% → 100% → revoke DB grants; sequence consumers by risk; rollback via flag; final step: enforce access via DB roles.
9. **Operational plan:** SLOs/alerts, dashboards, runbook, capacity (cache), on-call ownership, cost.
10. **Security/compliance:** authZ on APIs, audit logging, data classification.
11. **Open questions/timeline:** owners, staffing, dates.

**Q3. ADR-0033: Adopt canary deployments with automated rollback.**
- **Status:** Proposed.
- **Context:** weekly, 3-hour deploys with a 25% change failure rate and MTTR ~3 h; all 3 incidents last quarter followed deploys; big batched releases make root-causing hard; we have SLOs on the API (99.9%).
- **Decision:** Deploy each service via canary (1% → 10% → 50% → 100%) with automated analysis of error rate and p99 vs baseline for 10 minutes/step; auto-rollback on regression; feature flags for risky behavior changes; move toward daily deploys.
- **Consequences:** + smaller blast radius, faster detection/rollback (MTTR target < 30 min), enables higher deploy frequency; − needs traffic-splitting infrastructure (mesh/ingress), good SLI instrumentation, investment in flag hygiene; canary analysis needs enough traffic (low-traffic services need longer steps).
- **Alternatives:** blue-green (fast rollback but full-capacity duplicate and exposes all users at cutover), manual staged rollout by region (no automatic analysis), status quo + more testing (doesn't catch production-only issues).
- **Revisit if:** canary analysis false-positive rate > 10%, or services lack traffic for meaningful comparisons (switch to flag-based rollouts).

**Q4. Postmortem: Kafka misconfiguration (model).**
- **Summary/impact:** 40 minutes of shipment events lost during a broker restart on <date>; downstream tracking and notifications were incomplete until reconciliation repaired data 6 hours later; ~X customers saw missing status updates; SLO for event delivery breached (budget burn).
- **Timeline:** 10:02 broker-3 restarted for maintenance; 10:03 topic leader unavailable, `acks=1` producers acknowledged writes only to the leader; 10:03-10:43 unreplicated partitions lost on leader failure/unclean election; 10:45 broker healthy; 16:10 reconciliation flagged mismatch; 16:40 incident declared; 18:30 events replayed from the source DB outbox/backfill.
- **Root causes (system):** topic configured with replication factor 1 and producer `acks=1` (no durability guarantee); no ownership/standards for Kafka configuration; maintenance procedure didn't check under-replicated partitions.
- **Contributing factors:** no alert on under-replicated/offline partitions or producer error rates; reconciliation only ran every 6 h; no config review or policy-as-code; "nobody owns Kafka."
- **What went well:** reconciliation detected the gap; the outbox made replay possible. **Where we got lucky:** no financial-impact events were in the lost window.
- **Action items:** (1) *Prevent:* enforce RF=3, `min.insync.replicas=2`, producer `acks=all` via cluster policy/linting: Platform team, 2 weeks. (2) *Prevent:* assign Kafka ownership (platform team) with runbooks, on-call: Eng Director, 1 week. (3) *Detect:* alerts on under-replicated partitions, offline partitions, producer failed sends, consumer lag: SRE, 2 weeks. (4) *Detect/Mitigate:* run reconciliation hourly and build a replay tool from the outbox: Shipments team, 4 weeks. (5) *Prevent:* maintenance checklist (rolling restarts only when ISR healthy) + game-day test of broker loss: Platform, 6 weeks.

**Q5. Roadmap (model).**
- **Q1: Stabilize and make change safe.** Ownership of Kafka; alerts/SLO burn-rate alerting; CI speedup (parallelize tests, target < 30 min); canary + auto-rollback (ADR); on-call improvements/runbooks. Expected: MTTR 3 h → < 1 h, change failure rate 25% → < 15%.
- **Q2: Decouple the data.** Shipments ownership service (design doc); DB access via roles; first consumers migrated; modular monolith boundaries and per-module CI. Expected: deploy frequency weekly → several per week; fewer schema-related incidents.
- **Q3: Extract for pain, not fashion.** Extract 1-2 high-value seams (e.g., notifications/search) via strangler; more consumers off shared DB; progressive delivery everywhere; game days. Expected: daily deploys for extracted services, change failure < 10%.
- **Q4: Consolidate and scale.** Finish shipments migration; remove dead code paths/old tables; cost/capacity review; decide next extractions with data; retro on DORA metrics.
- **Communication to leadership:** frame in outcomes (deploy frequency, incidents, on-call health, time-to-market); show the plan *is* an incremental path to a better architecture with lower risk (strangler fig), give the cost/risk of the 6-month big-bang, define metrics and quarterly checkpoints, and reserve the right to adjust with data; offer to revisit microservices scope after Q2 evidence.

**Q6. Design review plan (model).**
- **Attendees/lenses:** Shipments team (owner), two consumer-team leads, SRE (operability), security, data/DBA, product/EM (timeline), an independent senior engineer as devil's advocate.
- **Pre-read:** design doc circulated 3 days before; reviewers submit comments/questions in advance; decision needed is stated at the top.
- **Agenda (60 min):** 5 min problem/constraints; 10 min proposal walkthrough; 25 min discussion of top-rated concerns (risks, alternatives); 10 min migration/rollout/ops plan; 10 min decision and action items.
- **Decision rule:** the design owner decides after input (single accountable decision-maker), with blocking concerns needing explicit resolution or documented acceptance; capture in ADR.
- **Strong disagreement:** surface the criteria/assumptions behind each position, look for data or a cheap experiment (spike, load test, prototype) to settle it, consider a time-boxed reversible trial, escalate to the accountable owner if unresolved, then "disagree and commit," recording the dissent and revisit conditions.

**Q7: B.**
- Visible tracking + capacity tied to outcomes. A and D ignore debt; C is unrealistic.

**Q8: B.**
- Independent written input first prevents anchoring; then discuss alternatives against criteria. A/D are HiPPO; C avoids the problem.

**Q9. Sample answer.**
- After the team has debated and a decision is made by the accountable owner, dissenters commit to executing it wholeheartedly (and can log their concerns and revisit conditions). Appropriate for two-way doors or when more debate costs more than the difference between options; harmful if used to silence valid safety/security concerns, when new evidence emerges (should reopen), or when dissent wasn't actually heard.

**Q10. Sample answer.**
- Argument: the pipeline is reliable (99.99% over 3 years), cheap (~$200/mo), meets the requirement (a morning report), and a rewrite costs months and adds an operational burden without user-visible benefit; risk of regressions; opportunity cost.
- Monitoring to add: run success/failure alerts with dead-man's-switch (alert if it *didn't* run), data-freshness SLI, row-count/quality checks, runtime trend, and documentation/runbook with an owner.
- Revisit if: requirements change (intra-day/near-real-time reports, or 100x data volume causing runtime > window), or the maintainer/platform (scripts/OS) becomes unsupported/unmaintainable.
