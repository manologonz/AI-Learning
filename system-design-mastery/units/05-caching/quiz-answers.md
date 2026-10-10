# Unit 5 Quiz: Answer Key

**Q1: B.**
- Deleting avoids concurrent-writer races; the next read reloads the source of truth.
- A: updating from app code can persist stale values after races. C: standard caches don't watch the DB. D: flushing everything causes a stampede.

**Q2: 3,000 QPS at 95%; 600 QPS at 99%.**
- 60k × 0.05 = 3,000; 60k × 0.01 = 600. A 4-point hit-rate gain cuts DB load 5x.

**Q3: B.**
- Stampede/dogpile on expiry of a hot key. Single-flight (one request loads, others wait on it), jittered TTLs, locking, or serving stale while revalidating.
- A: penetration is about non-existent keys. C: avalanche is many keys expiring at once/whole-tier failure. D: hot key is about a single node's throughput, not simultaneous misses.

**Q4: B.**
- These are penetration requests; the key never exists, so it never gets cached. Cache the "not found" for a short TTL, use a Bloom filter to reject impossible keys, and validate input.
- A: TTL doesn't help for missing keys. C: write-through doesn't address reads for nonexistent keys. D: eviction policy is unrelated.

**Q5: Sample answer.**
- Advantage: very low write latency and high write throughput (writes acknowledged from memory; batching/coalescing to the DB reduces DB writes).
- Risk: data loss if the cache node fails before flushing; also ordering/consistency complexity and the DB lagging behind.

**Q6: 5 ms at 90%; 1.4 ms at 99%.**
- Avg = h×1 + (1−h)×(1+40). At 90%: 0.9 + 0.1×41 = 0.9 + 4.1 = **5.0 ms**. At 99%: 0.99 + 0.01×41 = **1.4 ms**.

**Q7: B.**
- Classic race. TTL bounds the damage; version/timestamp checks, delayed double delete, or write-through with ordering reduce it.
- A, C are different failure modes; D is not a bug, it's inherent to the pattern.

**Q8: Sample answer.**
1. **In-process local cache** with a 1-2 s TTL in each app instance. Downside: instances may briefly serve different/stale values; memory duplicated.
2. **Key replication/sharding** (`profile:123#0..9`, reads pick a random suffix). Downside: writes/invalidation must touch all copies; more memory.
3. (Also valid) CDN/edge caching. Downside: invalidation harder, and personalization limited.

**Q9: B.**
- The DB is the fallback; if the cache fails or restarts cold, the full load lands on it. Design for it (replicas, warming, load shedding).
- A is wrong: 2% of 500k QPS is still 10k QPS. C: write-through narrows the window but failure between steps or multiple writers can still cause inconsistency. D: TTL is a safety net even with invalidation because invalidation can fail.

**Q10: Sample answers.**
- (a) Catalog descriptions: cache-aside + CDN, TTL minutes-hours with purge on edit; stale by minutes is fine.
- (b) Shopping cart: store in Redis/DB keyed by user with write-through (cart is user-owned, read-your-writes required); TTL for abandoned carts; no cross-user staleness issue.
- (c) Balance at authorization: **do not cache**; read the source of truth in the transaction. (A cached balance is fine for display only.)
