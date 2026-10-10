# Unit 6 Quiz: Replication & Partitioning

Answers in `quiz-answers.md`.

**Q1 (MCQ, easy).** Which does adding read replicas primarily improve?
A. Write throughput
B. Read throughput and availability
C. Storage capacity per dataset
D. Cross-shard join performance

**Q2 (application, easy).** A leaderless store has N=5. Which (W, R) pair guarantees read/write quorum overlap: (2,2), (3,3), (1,5), (4,1)? List all that do.

**Q3 (MCQ, medium).** A user posts a comment, refreshes, and the comment is missing; a few seconds later it appears. Which anomaly and fix?
A. Monotonic read violation; pin sessions to the leader forever.
B. Read-your-writes violation; read the user's recent writes from the leader (or an up-to-date replica).
C. Split brain; add more replicas.
D. Phantom read; use Serializable.

**Q4 (MCQ, medium).** Which shard key is most problematic for a write-heavy events table with range partitioning?
A. `hash(device_id)`
B. `(device_id, timestamp)`
C. `timestamp` alone
D. `(tenant_id, hash(event_id))`

**Q5 (short, medium).** Why does `hash(key) % N` cause a mass cache miss when N changes from 10 to 11, and roughly what fraction of keys stay in place?

**Q6 (MCQ, medium).** In consistent hashing, what is the purpose of virtual nodes?
A. To encrypt keys.
B. To spread each physical node over many ring positions so load is balanced and a node's failure redistributes across many nodes.
C. To make the ring smaller.
D. To allow range queries.

**Q7 (application, hard).** A secondary index by `email` on a table sharded by `user_id`. Compare local and global secondary indexes for a `WHERE email = ?` query: cost on read, cost on write, and freshness.

**Q8 (MCQ, hard).** A team promotes an async follower after the leader crashes; the old leader recovers with a few seconds of writes the new leader never got, and still believes it's the leader. Which problems have occurred?
A. Only replication lag.
B. Lost acknowledged writes and split brain; requires fencing/demotion and reconciling or discarding divergent writes.
C. A cache stampede.
D. Nothing; this is normal and safe.

**Q9 (short, hard).** A multi-leader setup in two regions: two users edit the same document at once, one in each region. Name three ways to resolve the conflict and the main downside of each.

**Q10 (application, hard).** A social app shards posts by `user_id`. A celebrity with 30M followers gets 200k reads/s on their posts. List three techniques to keep that shard alive.
