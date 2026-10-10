# Unit 6 Quiz: Answer Key

**Q1: B.**
- Replicas serve reads and provide failover. Writes still go to the leader (A), each replica stores the full dataset (C), and joins across shards aren't affected (D).

**Q2: (3,3) and (1,5).**
- (2,2)=4 no. (3,3)=6 yes. (1,5)=6 yes. (4,1)=5 not greater than 5: no.
- Correct set: **(3,3) and (1,5)**. (4,1) sums to exactly 5, so a single read replica might be the one that missed the write.

**Q3: B.**
- The read hit a follower that hadn't received the write. Read recent writes from the leader (for a time window or until replica LSN catches up).
- A: monotonic reads is about going backwards between reads, and "pin to leader forever" removes read scaling. C: split brain is dual leaders. D: phantoms are a transaction isolation issue.

**Q4: C (`timestamp` alone).**
- With range partitioning, all current writes land in the newest range: one hot shard.
- A, D spread writes by hash. B partitions by device first, then sorts by time: writes spread across devices and range scans per device work.

**Q5: ~90% of keys move; ~1/11 (~9%) stay.**
- Changing the modulus changes almost every key's remainder; only keys where `h mod 10 == h mod 11` (about 1/11) stay. The result is a near-total cache flush and a DB stampede. Consistent hashing moves only ~1/N of keys.

**Q6: B.**
- Vnodes give each physical node many positions, smoothing load and spreading the load of a failed node over many others. A/C/D are unrelated.

**Q7. Sample answer.**
- **Local index:** read = scatter-gather to all shards, so latency is the max of all shards (tail latency ↑, load ∝ shard count); write = cheap and transactional within the shard; fresh.
- **Global index (partitioned by email):** read = one shard lookup then a fetch by user_id (2 hops); write = the base row and index entry live on different shards, needing a distributed write (often async, so the index can be stale or briefly inconsistent).
- Pick local when reads on that attribute are rare or shard count is small; global when the query is on a hot path and slightly stale index is acceptable.

**Q8: B.**
- Lost writes (never replicated) and split brain (old leader believes it's still leader). Requires fencing tokens/STONITH, demotion, and reconciling divergent writes.
- A understates it; C irrelevant; D false.

**Q9. Sample answer.**
- **Last-write-wins:** simple, but silently loses one edit and depends on clocks. 
- **Merge via CRDT/OT:** preserves both edits, but complex to implement and limited to data types with merge semantics.
- **Application-level resolution/ask the user (or keep siblings):** preserves information, but pushes complexity to the app/UX.
- (Also: route each document's writes to a single "home" region to avoid conflicts; downside: cross-region latency for the non-home user.)

**Q10. Sample answer.**
- Cache the celebrity's posts aggressively (CDN + in-process + Redis) since posts are read-heavy and largely immutable.
- Replicate the hot data: read replicas dedicated to that shard, or key replication with random suffixes.
- Split the hot partition: sub-shard the celebrity's data (`(user_id, hash(post_id))`), and rate-limit/degrade (serve stale) to protect the shard.
