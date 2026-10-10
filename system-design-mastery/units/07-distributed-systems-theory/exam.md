# Unit 7 Exam: Distributed Systems Theory

Answers in `exam-answers.md`.

---

### Scenario (Q1-Q5): SeatLock

SeatLock runs a reservation service across three regions (us, eu, ap). Each region has an app fleet and a local database replica. Requirements: (1) a seat must never be booked by two customers; (2) booking confirmation must be visible to that customer immediately on any device; (3) browse pages can be seconds stale; (4) service should keep accepting *browsing* even if inter-region links fail; (5) the reservation coordinator uses a leader elected via a lock service with 15-second leases.

**Q1 (analyze, medium).** For requirements 1-4, assign a consistency model and CAP/PACELC posture to each of booking, confirmation view, and browsing. Justify.

**Q2 (analyze, hard).** The leader worker pauses 25 s due to a long GC. Describe the sequence that leads to double booking with the lease design, and give the corrected design with fencing tokens. Show what the database check looks like.

**Q3 (apply, hard).** The lock service is a 5-node consensus cluster spread as 2 (us), 2 (eu), 1 (ap). A partition isolates us from the others. Who can commit? What if instead the layout were 3/1/1 and the 3-node region is isolated from the rest? Discuss availability of booking in each region.

**Q4 (evaluate, hard).** Booking requests from ap must reach the leader in us (RTT ~200 ms) plus a majority commit from nodes in us and eu (say 80 ms). Estimate booking latency and propose two ways to reduce perceived latency without breaking requirement 1.

**Q5 (evaluate, hard).** A teammate proposes using timestamp-based LWW across regions for seat status to avoid coordination. Give a concrete failure, and state precisely what guarantee you would need to make LWW acceptable (or an alternative).

---

### Independent items

**Q6 (MCQ, medium).** A system's node A sends to B, then B sends to C, then C writes X. Another node D independently writes Y at about the same time with no communication. In vector clocks, what is the relationship between X and Y?
A. X happened-before Y.
B. Y happened-before X.
C. They are concurrent (neither happened-before the other).
D. They are identical.

**Q7 (MCQ, hard).** An engineer says: "Our database is CP, so it's always consistent and never down." Which critique is most accurate?
A. CP databases never experience partitions.
B. CP means that during a partition some nodes become unavailable to preserve consistency; it doesn't guarantee freedom from downtime, and "consistent" here means linearizable only under specific configuration.
C. CP means the DB is eventually consistent.
D. CP systems ignore quorum.

**Q8 (short, hard).** Why is a Raft cluster with 3 nodes across three regions typically slower to commit than 3 nodes in one region, and how would you architect the trade-off?

**Q9 (application, hard).** Design a leader election with lease + fencing using etcd or ZooKeeper for a job scheduler: describe the keys/nodes used, the lease TTL and renewal cadence relative to timeouts, and what the scheduler does when it loses the lease.

**Q10 (create, hard).** Choose a consistency model per operation for a collaborative document editor (typing, saving, sharing permission changes, view counts, comment threads). For each, state the model and one sentence of justification.
