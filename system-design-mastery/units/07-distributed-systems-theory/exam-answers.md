# Unit 7 Exam: Answer Key

**Q1.**
- **Booking (req. 1):** linearizable writes on the seat record (single leader or consensus/CAS); **CP/PC-EC** posture: during a partition, the side without a quorum must refuse bookings (accept unavailability rather than double-book); pays latency for coordination.
- **Confirmation view (req. 2):** read-your-writes/session consistency at minimum; simplest: after booking, read from the same authoritative store (or return the committed result and route that user's reads to the leader/up-to-date replica for a period). Not necessarily global linearizability.
- **Browsing (reqs. 3-4):** eventual consistency from local replicas/caches; **PA/EL**: stays available during partitions, low latency, seconds stale.
- Key point: different operations get different models: strong where invariants live, weak elsewhere.

**Q2.**
- Sequence: Leader A acquires the lease (15 s), starts a booking, pauses 25 s. Lease expires at 15 s; B acquires it and books seat 12 for customer 2. A resumes, still believing it's leader (its local clock check hasn't run), and books seat 12 for customer 1: double booking.
- Fix: the lock service returns a **fencing token** (monotonic, e.g., etcd revision or ZooKeeper zxid): A has token 41, B token 42. Every write to the DB includes the token: `UPDATE seats SET holder=$c, fence=$tok WHERE seat=$id AND fence <= $tok AND holder IS NULL` (or a stored procedure that rejects `fence < current`). A's write with 41 is rejected after B wrote with 42. Additionally use a conditional write/unique constraint on the seat so even without tokens there's at most one holder.

**Q3.**
- 5 nodes: 2/2/1. Isolate **us (2 nodes)** from the other 3: the us side has 2 (< 3 majority) → cannot commit; the eu+ap side has 3 → majority, elects a leader and continues. us bookings fail (CP), us browsing continues from local replicas.
- Layout 3/1/1 with the 3-node region isolated: that region holds a majority (3 of 5) → it continues; the other regions (2 nodes) can't commit → bookings unavailable in eu/ap, browsing still works.
- Insight: placement decides which region survives; the 2/2/1 layout survives any single-region loss (the remaining 3 form a majority), but 3/1/1 does not survive losing the 3-node region (only 2 nodes remain).
- So 2/2/1 tolerates loss of any one region; 3/1/1 does not tolerate loss of the big region.

**Q4.**
- Latency ≈ ap→us leader 200 ms RTT (client request + response) + majority commit ~80 ms (leader replicates to a nearby node in eu and needs 3 acks; with 2 in us, quorum is us-us-eu ≈ 80 ms) → **~280 ms** plus app time (~300 ms total).
- Reduce perceived latency: (1) **optimistic UI/pre-reservation hold**: show "reserving..." instantly with an async confirmation (booking not final until confirmed), (2) **home-region sharding**: assign each event/venue a leader in the region closest to the buyers (eu venues led from eu), so most bookings avoid the long haul, (3) **batching or coalescing** and connection reuse to trim overhead. Requirement 1 remains intact because the single-leader-per-seat-shard commit path is unchanged.

**Q5.**
- Failure: two regions set the same seat's status at nearly the same real time; region A's clock is 200 ms ahead, so its (earlier) "booked by Ann" carries a later timestamp than eu's "booked by Bob"; LWW keeps Ann and silently discards Bob (or both believe they booked it before convergence): violates req. 1 (both were told success).
- LWW can't enforce uniqueness at all: both writes were accepted before merging. Need either coordination (consensus/CAS/single owner per seat) or bounded clock uncertainty with commit-wait (TrueTime-like) plus still a conditional write. LWW is acceptable only for data where lost updates don't matter (e.g., a user's last-seen timestamp).

**Q6: C.**
- X's causal chain (A→B→C→X) has no communication with D's write Y, so neither happened-before the other: concurrent (vector clocks are incomparable). A/B claim causality that doesn't exist; D confuses concurrency with identity.

**Q7: B.**
- CP describes behavior under partition: minority sides sacrifice availability. It doesn't mean never down, and many "CP" stores are only linearizable for certain operations/configurations (e.g., default read modes). A/C/D are wrong.

**Q8. Sample answer.**
- Every commit needs acks from a majority (2 of 3), so latency is at least the RTT to the nearest other replica: cross-region 50-200 ms vs <1 ms in-region.
- Trade-off options: keep the consensus group in one region (or one region + close neighbors, e.g., two nearby AZs/regions plus one distant) and accept regional-failure risk; use many small consensus groups placed near their users (Spanner-style leaders per shard); use async replication to other regions for DR with a defined RPO; reserve cross-region consensus for small critical metadata.

**Q9. Sample answer.**
- Use etcd: a lease with TTL (e.g., 10 s) and key `/scheduler/leader` created with `if not exists` transaction (`Txn` compare create_revision == 0), attached to the lease. The create revision (or a counter) serves as the **fencing token**.
- Keepalive/renewal at TTL/3 (~3 s) to tolerate missed heartbeats; the scheduler stops work if it can't renew within (TTL − safety margin) using a monotonic clock.
- On lease loss (renewal fails, or the key is deleted/changed): immediately stop scheduling, abort in-flight assignments (or ensure they're idempotent), drop leader-only state, and rejoin as a follower/candidate, watching the key.
- Every scheduler action against workers/DB includes the token; downstream rejects stale tokens.

**Q10. Sample answer.**
- Typing (collaborative edits): **eventual with convergence via CRDT/OT** (concurrent edits merge; low latency, offline-capable).
- Saving a snapshot: **read-your-writes/monotonic** for the author; durable write acknowledged by a quorum.
- Permission changes (revoke access): **linearizable/strongly consistent** so a revoked user can't read after the revocation is acknowledged.
- View counts: **eventual**; approximate is fine.
- Comment threads: **causal consistency** so replies never precede their parents.
