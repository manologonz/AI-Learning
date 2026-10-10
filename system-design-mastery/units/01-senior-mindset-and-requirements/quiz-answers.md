# Unit 1 Quiz: Answer Key

**Q1: C.**
- **Why:** An SLO is an internal target for an SLI over a window (e.g., 99.9% good requests over 30 days).
- A is wrong: that's an SLA (external, with penalties). Confusing SLO/SLA is the most common mix-up.
- B is wrong: that's an SLI, the measurement itself, not the target.
- D is wrong: that's capacity/throughput limit, unrelated to service level objectives.

**Q2: 43.2 minutes.**
- Window = 30 × 24 × 60 = 43,200 min. Budget = (1 − 0.999) × 43,200 = 0.001 × 43,200 = **43.2 min**.
- Common error: using 365 days or forgetting the window; the budget is always relative to the SLO window.

**Q3: B.**
- **Why:** The SLI measures a proxy (probe status) rather than user-experienced success. The health endpoint can return 200 while checkout fails, the DB is down, or responses are slow. Good SLIs are measured on real user requests and include correctness and latency.
- A is wrong: a "too high" target doesn't explain a reported 99.99% alongside real failures; the measurement is misleading, not the target.
- C is wrong: more frequent probes still measure the wrong thing.
- D is wrong: SLA vs SLO doesn't fix a bad SLI; an SLA built on this SLI would be just as misleading.

**Q4: C.**
- **Why:** It describes *how well* (latency at a percentile), not *what*. NFR = quality attribute with a measurable bar.
- A, B, D are all functional: they describe capabilities/features.

**Q5: Sample answer.**
The SLO is our own target; the SLA is a contract with financial or legal consequences. If the SLA equals the SLO, breaching our internal target would already mean paying penalties, leaving no time to react. A stricter SLO creates a buffer: alerts fire and the team can freeze launches and fix problems before the contractual line is crossed.
- Full credit needs: (1) SLA has penalties, (2) SLO is stricter to provide a safety margin/early warning.

**Q6: B.**
- **Why:** Averages mask the tail: 95% of requests at 50ms and 5% at 2s can still average low. Users at the tail have a bad experience, and in fan-out systems the tail dominates.
- A is wrong: averages can be computed; they're just poor for this purpose.
- C is wrong: units are irrelevant.
- D is wrong: the median also hides the tail; you need p95/p99 too.

**Q7: ≈ 26%.**
- P(no slow call) = 0.99^30 ≈ 0.740. P(at least one slow) = 1 − 0.740 ≈ **0.26**.
- Insight: fan-out amplifies the tail. A "1% problem" per call is a ~26% problem per page.
- Common error: 30 × 1% = 30% (an overcount that ignores overlap; close here but wrong in method and diverges as n grows).

**Q8: B.**
- **Why:** Seniors justify each component with a requirement or number; complexity has real cost (operations, debugging, hiring, cognitive load). With 200 users and 3 engineers, a modular monolith on a managed database meets needs, with stated conditions to revisit.
- A is wrong: Kafka scales very well; it's just unnecessary here.
- C is wrong: over-generalization; microservices are sometimes right. The critique is about fit, not dogma.
- D is wrong: swaps one default for another without analysis; serverless has its own trade-offs.

**Q9: Sample answer.**
"We choose eventual consistency for like counts over strongly consistent counts because low write latency and high availability on a hot path matter more than showing the exact count instantly; a count that lags by a few seconds harms no one. This holds while likes are purely informational. We'd revisit if likes gained financial or ranking-critical meaning (e.g., paid promotion billed per like) or if lag exceeded ~1 minute."
- Full credit: names gain (latency/availability/throughput), cost (temporarily stale), assumption (informational data), and a concrete reversal condition.

**Q10: Sample answer** (yours may differ; grade on structure and specificity).
- (a) FRs: customer can see the courier's location updating on a map for an active order; customer sees ETA; tracking starts at pickup and ends at delivery.
- (b) NFRs: location freshness ≤ 5s end-to-end (latency); handles e.g. 100k concurrent active deliveries → ~20k location writes/s if each courier reports every 5s; availability high (99.9%) but brief staleness tolerable (best-effort, not durable per point); battery/bandwidth efficiency of the courier app; privacy (only the order's customer can see the courier's location, and only during the delivery).
- (c) Out of scope: route optimization/dispatch, historical route analytics, in-app chat.
- (d) Hard part: high-frequency write ingestion and fan-out of live updates to viewers (long-lived connections), while tolerating lossy, out-of-order data. Durability of every point is *not* required; freshness is.
- Full credit: at least one measurable NFR, one that's clearly non-functional (not a feature), and a hard part that follows from the NFRs.
