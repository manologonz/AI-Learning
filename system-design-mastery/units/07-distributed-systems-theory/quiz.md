# Unit 7 Quiz: Distributed Systems Theory

Answers in `quiz-answers.md`.

**Q1 (MCQ, easy).** A 5-node Raft cluster can tolerate how many node failures while still committing writes?
A. 1
B. 2
C. 3
D. 4

**Q2 (MCQ, easy).** During a network partition, a CP system typically:
A. Serves reads and writes on both sides and reconciles later.
B. Refuses (or degrades) operations on the side that can't reach a quorum to preserve consistency.
C. Shuts down entirely.
D. Ignores the partition.

**Q3 (MCQ, medium).** Why is using wall-clock timestamps for last-write-wins conflict resolution dangerous?
A. Timestamps take too much space.
B. Clock skew can make an earlier write appear later, silently discarding the true latest update.
C. Wall clocks can't be read in distributed systems.
D. It causes deadlocks.

**Q4 (short, medium).** Explain what a fencing token is and how it protects against a paused ex-leader.

**Q5 (MCQ, medium).** Which correctly states what CAP's "C" means?
A. The ACID "consistency" that enforces application constraints.
B. Linearizability: every read reflects the latest completed write as if there were a single copy.
C. Eventual convergence.
D. Causal ordering.

**Q6 (application, medium).** Match each to the weakest adequate consistency model: (a) uniqueness of usernames, (b) view counts on a video, (c) a comment thread where replies must never appear before parents.

**Q7 (MCQ, hard).** A 4-node cluster vs a 3-node cluster running majority-quorum consensus. Which statement is correct?
A. The 4-node cluster tolerates 2 failures.
B. Both tolerate exactly 1 failure; the 4-node cluster adds cost and commit latency without extra fault tolerance.
C. The 3-node cluster tolerates 0 failures.
D. Both tolerate 2 failures.

**Q8 (short, hard).** A request to the payment service times out. List four different real-world states the system might be in, and state the client-side design principle that makes retrying safe.

**Q9 (MCQ, hard).** Two events e1 and e2 have Lamport timestamps 5 and 9 respectively. What can you conclude?
A. e1 happened-before e2.
B. Nothing definitive: e1 could have happened-before e2, or they could be concurrent; but e2 did not happen-before e1.
C. e2 happened-before e1.
D. They are concurrent.

**Q10 (short, hard).** Explain PACELC in your own words with one example system for each of PA/EL and PC/EC and the reason for its position.
