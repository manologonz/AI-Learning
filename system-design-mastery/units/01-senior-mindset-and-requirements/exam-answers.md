# Unit 1 Exam: Answer Key

**Q1.**
- (a) **FR**: describes a capability of the system.
- (b) **NFR** (correctness/consistency/idempotency guarantee): it constrains *how* the transfer behaves under retries. Many would call it a "reliability" NFR; it drives design (idempotency keys, unique constraints).
- (c) **NFR** (latency with a percentile and number).
- (d) **Out of scope (for now)**: a feature explicitly deferred; adds regulatory and technical complexity that doesn't serve the current core promise.
- Common error: calling (b) an FR because it "sounds like a feature." Test: does it describe a capability, or a quality of an existing capability? Quality → NFR.

**Q2. Sample answer.**
- **SLI:** good = `POST /transfers` requests that return a non-5xx status *and* complete within 1s; valid = all `POST /transfers` requests excluding client errors (4xx from invalid input) and synthetic/health traffic. Measured at the API edge on real traffic.
- **SLO:** 99.95% over a rolling 30 days.
- **Justification:** money movement is the core promise and failures cause support load and press; users tolerate occasional retry but not sustained failure. 99.95% is achievable on a well-run monolith + managed Postgres with replicas, whereas 99.99% would force multi-region and heavy investment the company can't fund yet. (Any well-argued number between 99.9% and 99.99% is fine if justified.)
- **Budget:** (1 − 0.9995) × 43,200 = **21.6 minutes** per 30 days.
- Grading: the SLI must be user-facing and include valid vs good events; the number must be justified by cost/benefit, not "because more is better."

**Q3. Sample arguments.**
1. **The budget is overspent.** 99.95% allows 0.05% bad; actual bad is 0.07%, meaning 140% of the budget is consumed (0.07/0.05). Error-budget policy says to slow feature launches and prioritize reliability. "Basically" ignores that the budget is the mechanism for deciding.
2. **The SLO exists to give a safety margin above the SLA and to drive behavior.** If the SLA is 99.9%, we're still above it, but we've lost margin; a small further degradation breaches the contract. Treating misses as noise erodes the credibility of SLOs.
3. Also acceptable: the shortfall is concentrated at peak times (paydays), when it hurts most; averaged windows hide the pain of that bad hour.
- Full credit: two distinct arguments, at least one quantitative use of the error budget.

**Q4. Sample answer.**
"A 20-service rewrite trades a working monolith for large upfront cost: months of feature freeze, distributed-systems failure modes (partial failure, network latency, data consistency across services), operational load on a 15-person team, and higher risk to a payments system where correctness is paramount. Nothing here shows the monolith can't handle the current load or that the growth targets require it. Instead, stage it: (1) load-test and scale the monolith vertically and with read replicas, caching, and a queue for async work; (2) enforce module boundaries inside the monolith (modular monolith); (3) extract a service only where a boundary is clear and independent scaling or team autonomy justifies it. Triggers for sharding: primary write throughput sustained above ~60% of tested capacity after tuning, the DB exceeding what a single node's storage/IOPS can serve (e.g., >70% disk or hot working set no longer fitting memory), or replication lag/maintenance windows breaching the SLO."
- Full credit: names concrete costs, offers a staged path, gives ≥2 *measurable* triggers.

**Q5. Sample answer.**
- (i) *Home-screen balance:* "We choose a cached balance (≤5s stale) over reading the primary every time because latency and DB load matter more than sub-5s freshness for a display value. Holds while the number is informational; revisit if users complain about mismatch after a transfer, in which case read-your-writes via cache invalidation on the user's own transfers."
- (ii) *Authorization check:* "We choose strongly consistent reads (primary, within the transfer's transaction) over cached because correctness (no overdraft, no double spend) matters more than the ~40ms extra latency. This holds always for money movement; there's no reasonable reversal condition short of a redesign of the ledger."
- **Different choices** in each context, and that's the point: consistency requirements attach to *use of data*, not to the data itself.
- Common error: choosing one option globally. Senior answers differentiate by context.

**Q6. Sample answer.**
- **Hard part: correctness of money movement** (atomic, exactly-once-effect transfers with an auditable ledger under retries, failures, and concurrency), with availability as a close second. If we get correctness wrong, the core promise fails and trust is destroyed; scale is achievable, corrupt balances aren't recoverable by scaling.
- **NFRs to quantify first:** (1) *consistency/correctness requirements* (e.g., zero tolerance for double-apply; what reconciliation guarantees; audit requirements) and (2) *availability target and peak load* (payday peaks, concurrent transfers per second), since they determine the SLO, capacity plan, and cost. (Compliance/PCI is a defensible alternative.)

**Q7: B.**
- **Why:** Serial dependencies multiply: 0.999 × (0.9999)^n ≈ ~0.998-0.999 at best. You cannot exceed the availability of a hard serial dependency without redundancy (multi-vendor), caching, queuing, or graceful degradation.
- A is wrong: weakest link, not strongest, sets the ceiling for serial dependencies.
- C is wrong: an SLO is a target, not a fact; math constrains what's achievable.
- D is wrong: credits are a financial remedy, not availability; users are still down.

**Q8: B.**
- **Why:** P(at least one slow) = 1 − 0.99^40 ≈ 1 − 0.669 ≈ **33%**. Tail amplification means per-call SLOs don't compose to page-level SLOs.
- A is wrong: SLOs don't compose that way.
- C is wrong: page latency for parallel calls is governed by the *slowest* call (max), not the average.
- D is wrong: SLOs apply to whatever the user experiences; the SLI should be measured at the page/journey level.

**Q9. Sample answer.**
- **Availability:** the fraction of time/requests the system responds successfully. **Durability:** the probability that data, once acknowledged as stored, is not lost.
- *Durability over availability:* a banking ledger or medical records store; better to reject a write than to acknowledge it and lose it.
- *Availability over per-write durability:* live location pings or metrics/telemetry ingestion; dropping some points is fine, but the ingest endpoint must stay up.
- Common error: treating them as the same thing ("uptime"). A system can be always up but lose data (high availability, low durability), or never lose data but be frequently unreachable.

**Q10. Sample answer** (grade structure and specificity).
- **FRs:** upload a photo (JPEG/PNG up to 10MB); crop/replace/remove it; photo appears wherever the user's avatar is shown, at appropriate sizes.
- **NFRs:** upload completes within 3s at p95 for 5MB on a typical connection; avatar fetch p95 < 100ms (served via CDN, ≥95% cache hit); durability: uploaded photos are not lost (11 nines-class object storage); availability of avatar reads 99.95% (fallback to a default avatar when unavailable); moderation/safety scanning within 1 minute; scale: 10M users, ~1M uploads/day (~12/s average), reads ~500x more than writes.
- **Out of scope:** video avatars; face recognition/tagging.
- **SLI/SLO:** SLI = successful avatar GETs served under 200ms / valid avatar GETs; SLO = 99.9% over 30 days.
- **Hard part:** read-heavy delivery at low latency and cost (CDN + resizing/variant strategy). Upload volume is modest; the read-side egress and image variants dominate.
- **Trade-off statement:** "We choose object storage (e.g., S3) fronted by a CDN over storing blobs in the relational DB because photos are large, immutable, read-heavy, and durability/cost per GB matters more than transactional coupling with user rows. This holds while we only need whole-object reads/writes; we'd revisit if we needed transactional guarantees between metadata and blob (keep a DB row for metadata plus an upload-complete state machine instead) or on-the-fly transformations at a scale that makes a dedicated image service cheaper."
- Full credit: measurable NFRs, an FR/NFR split with no leakage, a hard part that follows from the numbers, and a trade-off with a reversal condition.
