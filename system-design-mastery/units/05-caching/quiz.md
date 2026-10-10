# Unit 5 Quiz: Caching

Answers in `quiz-answers.md`.

**Q1 (MCQ, easy).** In cache-aside, what should the application do after writing to the database?
A. Update the cache entry with the new value in the background.
B. Invalidate (delete) the cache entry so the next read repopulates it.
C. Nothing; the cache detects changes.
D. Flush the entire cache.

**Q2 (application, easy).** A service gets 60k QPS. The cache hit rate is 95%. How many QPS reach the database? What if the hit rate rises to 99%?

**Q3 (MCQ, medium).** A popular key's TTL expires and 5,000 concurrent requests all query the database simultaneously. What is this called and what's the standard mitigation?
A. Cache penetration; use a Bloom filter.
B. Cache stampede; use single-flight/request coalescing (plus TTL jitter or stale-while-revalidate).
C. Cache avalanche; use a larger cache.
D. Hot key; use a bigger Redis node.

**Q4 (MCQ, medium).** Requests keep asking for user IDs that don't exist, bypassing the cache and hammering the DB. Which fix targets the root cause?
A. Increase TTL.
B. Negative caching and/or a Bloom filter of valid IDs, plus input validation.
C. Switch to write-through.
D. Use LFU eviction.

**Q5 (short, medium).** Give one advantage and one risk of the write-back (write-behind) pattern.

**Q6 (application, medium).** Average cache latency is 1 ms, DB latency 40 ms, hit rate 90%. What is the average latency (assume a miss costs cache lookup + DB)? What is it at 99%?

**Q7 (MCQ, hard).** In cache-aside, a slow reader loads an old value from the DB, then a writer updates the DB and deletes the key, and finally the slow reader stores its old value in the cache. What's the best description?
A. Cache penetration
B. A stale-data race; it leaves stale data until TTL and can be mitigated with versioned writes, TTLs, or delayed double-delete
C. Cache avalanche
D. A Redis bug; upgrade Redis

**Q8 (short, hard).** One celebrity's profile key receives 300k reads/s and saturates its Redis shard. List two mitigations and one downside of each.

**Q9 (MCQ, hard).** Which statement is most accurate?
A. If the cache hit rate is 98%, the database no longer needs to handle production traffic.
B. The database must be able to survive a cold or failed cache (or you need protective measures like load shedding and warmed replicas), because cache failure sends full traffic to it.
C. Caches are always consistent with the database in write-through mode.
D. TTLs are only needed when you can't invalidate.

**Q10 (short, hard).** For each data type, choose a caching approach and TTL/consistency stance in one line each: (a) product catalog descriptions, (b) a user's shopping cart, (c) bank account balance shown at authorization time.
