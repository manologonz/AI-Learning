# Unit 5 Exam: Caching

Answers in `exam-answers.md`.

---

### Scenario (Q1-Q5): NewsFlash

NewsFlash serves articles. Peak: 80k QPS for article reads, 20 writes/s (edits, view counters excluded). The article DB (Postgres primary + 2 replicas) tops out at ~15k read QPS total. Redis cluster (6 nodes, ~100k ops/s per node) sits in front with cache-aside and 10-minute TTL. Breaking news causes a single article to receive ~40% of traffic. Editors need corrections to appear within 30 seconds; the editors' preview must be instant.

**Q1 (apply, medium).** What minimum hit rate is needed to keep the DB under 70% of its capacity at peak? Is a 10-minute TTL a problem for the editors' requirement? Propose a fix.

**Q2 (analyze, hard).** Describe what happens when the breaking-news article's key expires at peak, quantify the DB load spike, and give three defenses ordered by how much you'd trust them.

**Q3 (analyze, hard).** The hot article is ~32k QPS on one key. One Redis node is at 100k ops/s. Is Redis the bottleneck for this key? What about the network bandwidth if the article payload is 200 KB? Compute and propose a fix.

**Q4 (evaluate, hard).** A junior proposes updating the cache from the write path (`SET article:123 newBody`) instead of deleting. Explain a concrete interleaving where this leaves stale data even though nothing "failed," and recommend a design.

**Q5 (create, hard).** Design the multi-layer caching plan for articles: CDN, in-process, Redis, DB. For each layer give TTL/invalidation, what it protects, and its failure behavior. Include how a correction propagates within 30 seconds.

---

### Independent items

**Q6 (MCQ, medium).** After a Redis outage and restart (empty cache), your DB (15k QPS capacity) receives 80k QPS and collapses. What preventive design would *most directly* address this?
A. Increase the TTL to 24 hours.
B. Load shedding/rate limiting to the DB, request coalescing, and cache warming/replicated cache so a cold cache doesn't send full traffic to the DB.
C. Switch from LRU to LFU.
D. Remove the cache; the DB will be simpler.

**Q7 (MCQ, hard).** You're caching per-user recommendation lists that take 800 ms to compute and are viewed ~3 times per session. Which is the best judgment?
A. Cache aggressively with a long TTL regardless of freshness.
B. Cache with a TTL matching how quickly recommendations meaningfully change (minutes), precompute asynchronously if possible, and use stale-while-revalidate so users never wait 800 ms.
C. Never cache personalized data.
D. Cache in the browser only.

**Q8 (application, hard).** A cache holds 10M keys of ~2 KB values with Redis overhead of ~100 B/key. How much memory does it need at minimum (before replication)? If the hot set is 15% of keys and it serves 90% of traffic, what's the smaller cache size and the tradeoff?

**Q9 (short, hard).** Explain "read-your-writes" for cached data and give one implementation approach for a profile-edit scenario.

**Q10 (evaluate, hard).** Write a one-paragraph postmortem-style analysis of this incident: *"A deploy flushed the cache (`FLUSHALL`) at 09:00 Monday; p99 latency rose from 80 ms to 12 s; the DB CPU was 100% for 20 minutes."* Include root cause, contributing factors, and three corrective actions.
