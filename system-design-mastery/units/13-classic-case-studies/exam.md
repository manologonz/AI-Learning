# Unit 13 Exam: Classic Case Studies

Answers in `exam-answers.md`. Each design question expects: requirements/scope, estimates, API/data, high-level design, deep-dive on the hard part, failure modes, and trade-offs. Budget ~25 minutes per design.

---

**Q1 (design, hard): Design a "trending topics" system.** Users post short messages containing hashtags. The homepage shows the top 10 trending hashtags per country, updated every minute, over the last hour with recency weighting. Scale: 500M posts/day, 30% with a hashtag, peak 3x. Requirements: freshness ≤ 2 minutes, approximate counts acceptable, spam/bot resistance. Deliver estimates, pipeline design, data structures, and failure modes.

**Q2 (design, hard): Design a distributed job scheduler (cron-as-a-service).** Customers register jobs (cron expression or one-off, HTTP callback or queue message). Scale: 50M active jobs, ~20k executions/s at peak minute boundaries. Requirements: no job lost, at-least-once with idempotency support, ≤ 5 s late at p99, per-tenant fairness, retries with backoff, missed-run policy after downtime. Deliver the data model, timing wheel/queue approach, leader/sharding strategy, and failure handling.

**Q3 (analyze, hard): Critique this ride-hailing matching design.** "A single Postgres table `drivers(id, lat, lng, status)` updated on every ping (250k/s); match queries run `SELECT ... WHERE status='available' ORDER BY distance LIMIT 10` using a PostGIS index; the matching service picks the nearest driver and sets `status='busy'`." List at least five problems (performance, correctness, availability) and give the redesign.

**Q4 (design, hard): Design a photo-sharing "Stories" feature** (24-hour expiring content). 200M DAU, each views 30 stories/day; 20M new stories/day (avg 2 MB media). Requirements: ordered by recency and affinity, view receipts ("seen by"), expiry at exactly 24 h, global. Deliver estimates (storage, egress, QPS), storage/expiry design, feed generation approach, and the hard part.

**Q5 (evaluate, hard): Choose between two designs for a global unique-ID generator** used by all services at 500k IDs/s: (A) a central DB sequence with range allocation; (B) Snowflake-style 64-bit IDs (41-bit timestamp, 10-bit machine ID, 12-bit sequence). Compare on availability, ordering, clock issues, capacity per node, coordination needs, and give the recommendation and one failure scenario each.

**Q6 (create, hard): Cross-cutting synthesis.** Pick any two of the case studies in this unit (or your own) and describe three reusable patterns that appear in both, plus one place where they *should* differ despite looking similar (with the reason).
