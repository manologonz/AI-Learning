# Unit 5 Exam: Answer Key

**Q1.**
- 70% of 15k = 10.5k QPS allowed at DB. Miss rate ≤ 10.5k / 80k = 13.1% → hit rate ≥ **~87%**. Aim for 95%+ (4k QPS misses).
- A 10-minute TTL can leave corrections stale up to 10 minutes, violating the 30 s requirement. Fix: **invalidate on write** (delete key on edit) plus a shorter safety TTL (e.g., 60 s with jitter, or 30 s for headline pages) and an editor "preview" path that bypasses the cache (reads from the primary).

**Q2.**
- On expiry, all concurrent requests for the article (~32k QPS → within a 100 ms DB fetch, ~3,200 requests) miss simultaneously and each queries the DB: a spike of ~3k+ concurrent queries against a DB capable of ~15k QPS; if the query takes longer (say 500 ms) the pile-up is 16k requests, causing timeouts and retries: a stampede.
- Defenses (trust order): (1) **single-flight/request coalescing** (per-instance) plus a distributed lock or "leader refresh" so one request refills; (2) **stale-while-revalidate**: keep serving the expired value while one background request refreshes; (3) **proactive refresh** of known-hot keys before expiry / TTL jitter. Also: local in-process cache with a few seconds TTL to cut Redis/DB pressure.

**Q3.**
- Ops: 32k ops/s < 100k: not saturated on ops.
- Bandwidth: 32k × 200 KB = 6.4 GB/s ≈ **51 Gbps** from one node, which exceeds a typical 10-25 Gbps NIC. **Yes, network bandwidth is the bottleneck.**
- Fixes: **in-process cache** with 1-5 s TTL on each app instance (serving from local memory; Redis sees ~1 request/instance/TTL); serve via **CDN** (articles are public and cacheable: the real answer); compress payloads; replicate the key across nodes (`article:123#k`).

**Q4.**
- Interleaving: two editors save nearly simultaneously. W1: DB update (body A), W2: DB update (body B), W2: `SET cache B`, W1: `SET cache A` (delayed). Final: DB = B, cache = A: stale/wrong until TTL, and nobody's operation failed. (Also possible: a reader repopulating with an older value.)
- Recommended: **delete on write** (cache-aside), TTL as a safety net, and for stricter needs write a monotonic version and only set if newer (Lua compare-and-set), or route writes through a single ordered writer.

**Q5. Sample plan.**
- **CDN:** public articles, `Cache-Control: s-maxage=30, stale-while-revalidate=60`; **explicit purge** on edit (via API) so a correction shows within seconds; protects origin from 80k QPS (hit rate 99%+). Failure: falls back to origin (must be able to survive; use origin shielding).
- **In-process LRU:** top ~1,000 articles, TTL 2-5 s; protects Redis from hot key; failure: none critical (just misses).
- **Redis:** cache-aside, key `article:{id}:{version}` or delete-on-write; TTL 60 s + jitter; protects DB; failure: circuit breaker, single-flight, DB protected by rate limiting.
- **DB:** source of truth; replicas handle misses.
- Correction propagation: editor saves → DB write → delete Redis key → CDN purge → local caches expire within ≤5 s → total ≤ ~30 s worst case, usually < 5 s. Editors' preview reads primary directly.

**Q6: B.**
- A cold cache is the failure; you must protect the DB (load shedding, coalescing), and reduce cold-start risk (warming, replicas/persistence, gradual traffic ramp).
- A: longer TTL doesn't protect against a flush/restart. C: eviction policy is unrelated to a cold start. D: removes the protection instead of fixing it.

**Q7: B.**
- Expensive to compute, moderately dynamic: cache with a TTL aligned to meaningful change, precompute in the background, and serve stale while refreshing.
- A ignores staleness; C throws away a big win; D loses server-side sharing and control.

**Q8: ≈ 21 GB; hot cache ≈ 3.2 GB.**
- 10M × (2 KB + 100 B) = 10M × 2.1 KB = **21 GB**.
- Hot set: 15% of keys = 1.5M keys × 2.1 KB ≈ **3.15 GB**, serving 90% of traffic → 10% of traffic misses to the DB. Trade-off: 6.7x less memory for a hit rate around 90% versus a much higher hit rate with the full set; choose based on DB capacity.

**Q9. Sample answer.**
- Read-your-writes means a user who has just written data always sees their own write, even if others may see stale data. Implementation: after a profile edit, return the updated object in the response and have the client use it; server-side, set a short-lived per-user flag (or session marker with timestamp) so reads for that user bypass or refresh the cache for N seconds; or read from the primary for that user after a write.

**Q10. Sample postmortem.**
- Root cause: `FLUSHALL` in the deploy emptied the entire cache during peak traffic, so all reads went to the DB (80k QPS vs 15k capacity), saturating CPU and creating timeouts; retries amplified load and caches couldn't refill efficiently due to the stampede.
- Contributing: no protection for the DB (no load shedding/coalescing), no warming, deploy at peak (Monday 09:00), no guardrail against destructive commands, DB capacity depends on the cache.
- Corrective actions: (1) remove/ban `FLUSHALL` from deploy tooling (rename/disable command; use key-versioning for cache busting); (2) add single-flight and per-DB concurrency limits with load shedding so degraded mode returns partial results rather than collapse; (3) cache warming + gradual ramp after cold starts and run deploys off-peak; (4) test cold-cache behavior in load tests, and alert on hit-rate drops.
