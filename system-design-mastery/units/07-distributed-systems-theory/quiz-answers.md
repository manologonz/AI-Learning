# Unit 7 Quiz: Answer Key

**Q1: B.**
- 2f+1 = 5 → f = 2; majority is 3. A undercounts; C/D would leave fewer than a majority.

**Q2: B.**
- CP preserves consistency by making the minority side unavailable. A describes an AP system. C: only the minority side is affected. D: nonsense.

**Q3: B.**
- Skewed clocks reorder writes: a truly later write can carry an earlier timestamp and lose silently. A/C/D are wrong: timestamps are small, readable, and unrelated to deadlocks.

**Q4. Sample answer.**
- A monotonically increasing number issued by the lock/lease service with each grant. The protected resource remembers the highest token seen and rejects requests with lower tokens. A paused ex-leader that wakes up holds an old token, so its writes are rejected even though it still believes it holds the lock.

**Q5: B.**
- CAP's C is linearizability. A is ACID's C (different). C is eventual consistency (the opposite). D is causal consistency (weaker).

**Q6.**
- (a) Uniqueness: **linearizable** (needs a single agreed order/compare-and-set; or a unique constraint on a strongly consistent store).
- (b) View counts: **eventual** consistency is fine.
- (c) Comment thread: **causal consistency** (or consistent prefix): replies after parents, no need for a global order.

**Q7: B.**
- Majority of 4 is 3, so it still tolerates only 1 failure (same as 3 nodes: majority 2), but costs more and every commit needs 3 acks. A/C/D are wrong arithmetic.

**Q8. Sample answer.**
- Possible states: (1) request never arrived; (2) arrived and is still processing (slow); (3) processed and committed but the response was lost; (4) crashed midway, with a partial/no effect, possibly to be recovered.
- Principle: make operations **idempotent** (idempotency keys) so retries are safe, and provide a way to query the outcome (status endpoint/reconciliation).

**Q9: B.**
- Lamport clocks guarantee: if A → B then L(A) < L(B), but not the converse. So L(e1) < L(e2) doesn't prove causality; they could be concurrent. What you can conclude is e2 did not happen before e1 (else L(e2) < L(e1)). A over-claims; C is false; D over-claims (can't tell).

**Q10. Sample answer.**
- PACELC: *if there's a Partition, choose Availability or Consistency; Else (normal operation), choose Latency or Consistency.* Even without failures, coordinating replicas for strong consistency costs latency.
- PA/EL example: Cassandra/Dynamo default: stays available under partitions and favors low latency by reading/writing few replicas (eventual consistency).
- PC/EC example: Spanner/etcd-like systems: refuse minority operations under partitions and pay coordination latency (consensus round trips, commit-wait) for strong consistency in normal operation.
