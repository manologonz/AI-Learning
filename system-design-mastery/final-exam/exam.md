# System Design Mastery: Final Exam

**Cumulative across Units 1-14. Weighted toward your goal (real-world architecture at work): reliability & resilience, data scaling & consistency, and case-study application.**

Suggested time: 4-5 hours (Parts A + B ≈ 2 h; Capstone ≈ 2-3 h). Open-notes is fine; do not open `answers.md` until you've finished. Show assumptions and arithmetic.

| Part | Content | Points |
|---|---|---|
| A | 12 applied questions (recall/apply, cross-unit) | 36 |
| B | 4 analysis/evaluation scenarios | 36 |
| C | Capstone design doc | 28 |
| | **Total** | **100** |

Passing bar for "senior-ready": ≥ 80 overall **and** ≥ 20/28 on the capstone. Below 70: use the weak-unit map in `answers.md` for targeted review.

---

## Part A: Applied questions (3 points each)

**A1 (Unit 1, MCQ).** A service has a 99.99% availability SLO over a 30-day window. What is its error budget?
A. ~52 minutes
B. ~43 minutes
C. ~4.3 minutes
D. ~26 seconds

**A2 (Unit 2, calculation).** 200M DAU each make 25 requests/day; peak is 4x the average. One app server handles 2,500 QPS at test; target 60% utilization; you must survive losing one of three AZs with the remaining two carrying peak. Compute average QPS, peak QPS, and total servers.

**A3 (Unit 3, short).** A client retries `POST /transfers` after a timeout and the money moves twice. Give the fix and the three server-side behaviors it must define (replay after success, replay while in-flight, same key with different body).

**A4 (Unit 4, MCQ).** Two concurrent transactions each read `balance = 100` and each write `balance = 70`, and the final balance is 70 (should be 40). Which is the best fix?
A. Increase isolation to Read Uncommitted.
B. Use an atomic guarded update: `UPDATE acct SET balance = balance - 30 WHERE id = $1 AND balance >= 30`.
C. Add an index on `balance`.
D. Cache the balance in Redis.

**A5 (Unit 5, calculation + short).** 120k read QPS; the DB can serve 8k QPS. (a) What minimum cache hit rate is needed? (b) A hot key expires and 4,000 requests miss simultaneously. Name the failure and two defenses.

**A6 (Unit 6, MCQ).** Which leaderless configuration guarantees that read and write quorums overlap for N = 5?
A. W = 2, R = 3
B. W = 3, R = 3
C. W = 1, R = 4
D. W = 4, R = 1

**A7 (Unit 7, short).** A 9-node Raft cluster: (a) how many failures does it tolerate, (b) what is the majority, and (c) why does a paused ex-leader still endanger correctness even with leases, and what fixes it?

**A8 (Unit 8, short).** Write the consumer transaction (in pseudocode) that gives effectively-once processing for a `PaymentCaptured` event, and state where the ack goes and why.

**A9 (Unit 9, MCQ).** In a saga, the payment call times out. What should the orchestrator do?
A. Assume failure and compensate immediately.
B. Treat the outcome as unknown: query status or retry with the same idempotency key, and compensate only on definitive failure (plus reconcile later).
C. Assume success and continue.
D. Abort the whole saga silently.

**A10 (Unit 10, calculation).** Four layers (client, gateway, service, worker) each make up to 3 attempts against the layer below. What's the worst-case attempts against the database, and what retry policy do you recommend instead?

**A11 (Unit 11, MCQ).** Which is the strongest sign of a distributed monolith?
A. Services communicate over HTTP.
B. Services share a database schema, must be deployed together, and call each other in long synchronous chains.
C. Each service has its own CI pipeline.
D. There are more than 10 services.

**A12 (Unit 12, calculation).** SLO 99.9%. Errors run at 2.16% for the last hour and for the last 5 minutes. What is the burn rate, and what should happen (using 14.4x/1 h and 6x/6 h rules)? What fraction of a 30-day budget would be consumed if this persisted for one hour?

---

## Part B: Analysis and evaluation (9 points each)

**B1 (Reliability incident, Units 5, 8, 10, 12).** *Scenario:* At 09:00 Monday a deploy runs `FLUSHALL` on Redis. Read QPS is 80k; the DB handles 15k. Clients and the gateway retry 3 times each. The DB CPU is 100% for 40 minutes; after the deploy is "fixed" the outage continues 30 more minutes. A queue consumer for `orders.placed` also falls behind by 5M messages and starts double-processing orders on rebalances.
(a) Diagnose the failure chain including why it outlasted the trigger. (b) Give the immediate recovery procedure. (c) Give five prioritized design changes with the unit-level technique each uses. (d) Define SLO/alert changes that would have caught this in minutes.

**B2 (Data scaling, Units 2, 4, 6).** *Scenario:* An `orders` table on one Postgres primary: 6 TB, growing 100 GB/week, 40k writes/s at peak, 200k reads/s (80% by `order_id`, 15% by `customer_id` "my orders", 5% analytics by date). Customers and their orders are heavily skewed: 20 enterprise customers generate 35% of writes. Requirements: no downtime, strong consistency for order placement, "my orders" p99 < 100 ms.
(a) Show the numbers that force sharding (and what you'd try first). (b) Choose the shard key(s) and justify against the three access patterns and the skew. (c) Handle cross-shard access (`customer_id` queries, analytics). (d) Describe the online migration plan with rollback. (e) State two failure modes of your sharded design and mitigations.

**B3 (Consistency & transactions, Units 7, 8, 9).** *Scenario:* A global ticketing platform (us/eu/ap) sells seats; a seat must never be sold twice; browsing is read-heavy; payments use an external PSP with authorize/capture; holds last 10 minutes; regions may lose connectivity to each other.
(a) Choose a consistency model and CAP/PACELC posture for seat booking, seat browsing, and order history. (b) Design the booking workflow (saga steps, pivot, compensations, idempotency). (c) Where does the seat's authoritative state live and how do you avoid cross-region latency on the hot path? (d) What happens to bookings during an inter-region partition? (e) Give one scenario where LWW would have caused a double sale.

**B4 (Architecture critique, Units 1, 11, 12, 14).** *Scenario:* A 25-engineer company has 18 microservices sharing 2 databases, synchronous call chains of depth 6 on checkout, no tracing, weekly manual deploys, and an on-call rotation of 3 people. Leadership wants to add a service mesh and a second region.
(a) Evaluate the proposal, ordering what you'd do first with evidence to gather. (b) Identify five structural problems and the smallest fix for each. (c) Write the key trade-off statement (with reversal condition) for consolidating services. (d) Define 4 metrics of success.

---

## Part C: Capstone (28 points)

**Design and defend a real-time collaborative document platform ("DocLoom").**

*Product:* users create documents (rich text, up to 50 MB with embedded media), edit them **simultaneously** with others (live cursors, comments, suggestions), share with individuals/groups/organizations with fine-grained permissions, search across accessible documents, and see version history. Enterprise tenants demand data residency (EU data stays in the EU), audit logs, SSO, and 99.95% availability.

*Scale (assume, and state if you change them):* 50M DAU; 20M documents edited per day; typical active document has 1-5 concurrent editors, but a few reach 200 concurrent editors; average edit rate 1 op/s/editor while typing; 1B documents stored, average 200 KB text (media stored separately, avg 2 MB per media item, 30% of documents have media); 100k tenants (10 whales with 10M+ documents each).

**Deliver a design document (≤ 6 pages equivalent) containing:**
1. **Requirements**: functional, non-functional with SLOs (define SLIs), out-of-scope. (3 pts)
2. **Estimates**: QPS/ops rate, connections, storage, bandwidth; identify the dominant resource. (3 pts)
3. **API and data model**: documents, operations/versions, permissions, comments; choice of stores per data type with justification. (4 pts)
4. **Core architecture** with a diagram (text OK): edge, real-time collaboration servers, storage, search, async pipelines. (4 pts)
5. **Deep dive on the hard part**: concurrent editing (choose CRDT vs OT, session/doc ownership and routing, ordering, persistence/compaction, presence, offline edits, scaling to 200 editors). (5 pts)
6. **Consistency, permissions, and multi-tenancy**: what needs strong vs eventual consistency (e.g., permission revocation), isolation strategy for whales, data residency design. (3 pts)
7. **Reliability**: failure modes (collab server crash, region loss, hot documents), degradation, DR with RTO/RPO, retry/backpressure strategy, deploy strategy. (3 pts)
8. **Two ADRs** (short) for your two most consequential decisions, with alternatives and reversal conditions. (2 pts)
9. **Rollout/evolution & cost**: phased delivery plan, what you'd *not* build in v1, the top cost drivers and levers. (1 pt)
