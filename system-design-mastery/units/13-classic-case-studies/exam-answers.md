# Unit 13 Exam: Answer Key (model solutions and rubrics)

**Rubric for design questions (each ~10 pts):** requirements & scope (1), estimates with stated assumptions (2), data model/API (1), architecture (2), deep-dive on the hard part (2), failure modes & trade-offs (2). Reward correct reasoning even if the components differ.

---

**Q1. Trending topics: model solution.**
- **Requirements:** functional: top-10 hashtags per country; NFR: freshness ≤ 2 min, approximate OK, spam-resistant; out of scope: personalization, historical analytics.
- **Estimates:** 500M × 30% = 150M hashtag posts/day ≈ 1.7k/s avg, ~5k/s peak (assuming ~1.5 hashtags per post → ~8k hashtag events/s peak). Trivially small stream! The challenge is aggregation and heavy hitters, not raw QPS.
- **Design:** post events → Kafka (partition by hashtag) → stream processors doing **windowed counting** (sliding window, e.g., 1-min tumbling buckets aggregated over 60 min with exponential decay for recency) per (country, hashtag) → write top-K per country to a serving store (Redis sorted set / in-memory service) every minute → homepage reads from cache/CDN (TTL ~30-60 s). Use **count-min sketch + heap** or **Space-Saving** for approximate heavy hitters when cardinality is huge; exact counts in a keyed state store are fine at this scale.
- **Hard part:** distinguishing *trending* (velocity relative to baseline) from *popular*: score = (current window count) / (historical baseline expected count) or z-score; plus spam resistance: count **unique users** (HyperLogLog) not raw posts, dedupe per user per window, downweight new/low-reputation accounts, cap per-user contributions.
- **Failures:** processor crash → state restored from changelog/checkpoint, at-least-once → dedupe by event id or accept approximate; skewed partitions (a viral hashtag) → two-stage aggregation (local pre-aggregate then global); serving store failure → serve last good top-K (stale is better than empty). Trade-off: approximate counting for memory/latency; freshness vs cost of window recomputation.

**Q2. Job scheduler: model solution.**
- **Requirements:** register/update/delete jobs; trigger at time T; retries; missed-run policy (skip/catch-up-once/catch-up-all, configurable); NFR: durable, ≤5 s late p99, fair across tenants.
- **Estimates:** 50M jobs; 20k executions/s at minute boundaries (crons cluster on :00). Job record ~500 B → 25 GB (fits in a sharded DB easily). The challenge is the *thundering herd at minute boundaries* and timing accuracy.
- **Data model:** `jobs(job_id, tenant, schedule, next_run_at, payload, target, status, version)` sharded by `job_id` hash; index on `(shard, next_run_at)`. `executions(job_id, run_at, attempt, status)`.
- **Timing:** each shard owned by one scheduler (via lease/consensus, with fencing tokens); it loads jobs due in the next window (e.g., next 30 s) into an in-memory **timing wheel/priority queue**; at fire time it **enqueues an execution message** (with deterministic idempotency key `job_id:run_at`) to a work queue, updates `next_run_at` using the cron expression in the same transaction (or outbox). Workers pull from the queue (per-tenant fair queues/weights), call the target, ack, record results; retries with exponential backoff + jitter, DLQ after N.
- **Spreading the herd:** pre-enqueue with small random jitter within a tolerance (e.g., ±1 s if allowed), rate-limit per tenant, autoscale workers, split hot shards.
- **Failure handling:** scheduler crash → another instance takes the shard lease and rescans `next_run_at <= now` (jobs picked up late, but not lost; policy applies for missed runs); double-fire prevented by fencing tokens and the idempotency key; queue/at-least-once → target must dedupe on the key (documented); clock skew: use DB/queue time as reference and NTP; long downtime: apply the missed-run policy (catch up at most N).
- **Trade-offs:** DB-polling schedulers are simple and durable but need index/partition care; in-memory wheels are fast but rebuildable from the DB.

**Q3. Critique: model solution.**
- **Problems:** (1) 250k updates/s to a single Postgres table with indexes is far beyond a primary's capacity (WAL, index maintenance, vacuum bloat from constant updates); (2) durable writes for ephemeral data; (3) `ORDER BY distance` over all available drivers doesn't scale, and Euclidean/geodetic distance ≠ road ETA; (4) race conditions: two riders can match the same driver (read-then-write without atomic claim); (5) single point of failure/availability: DB down = no matching; (6) no timeouts/reoffer logic if the driver doesn't accept; (7) hot cells in dense areas; (8) status updates and location updates contend (row locks); (9) stale drivers (no TTL) remain "available".
- **Redesign:** in-memory geospatial index (S2/H3 cells → driver sets) sharded by region; location updates go to the index (TTL per driver) and a stream for durable sampling; matching service queries the rider's cell + neighbors, ranks by ETA (routing service or heuristic), **atomically claims** a driver (compare-and-set on driver state in Redis/Lua) with an offer timeout and reoffer; trip state in a durable DB via a state machine; per-city sharding, replicas for failover, and degradation (widen radius, queue riders) when supply is low.

**Q4. Stories: model solution.**
- **Estimates:** views = 200M × 30 = 6B/day ≈ 70k/s avg (~200k peak); uploads 20M/day ≈ 230/s; storage/day: 20M × 2 MB = 40 TB/day, but expires after 24 h → steady-state ~40-80 TB live (plus transcoded variants ×3 → ~150 TB); egress: 200k/s × ~500 KB delivered variant ≈ 100 GB/s ≈ 800 Gbps peak: **CDN is essential** (origin ≈ 5-10% of that). Metadata: 20M/day × 1 KB = 20 GB/day.
- **Storage/expiry:** media in object storage with **lifecycle rules** (TTL 24 h + grace) and CDN TTL ≤ 24 h; metadata in a partitioned KV store with TTL (`expires_at`); "exactly 24 h" enforced at read time (`now < expires_at`) so expiry is exact even if physical deletion lags (deletion is asynchronous).
- **Feed:** stories ring: per user, list of followed authors having active stories, ordered by affinity/recency; maintain `active_stories:{author}` (small, TTL); the tray for a viewer = pull for followed authors (hybrid push for small authors) with cached results; affinity scoring done by a ranking service with precomputed features.
- **Hard part:** read-heavy egress and the **view receipts**: 6B views/day → writes 70k/s; batch/aggregate view events (Kafka), store `seen(story_id, viewer)` as a probabilistic/partial structure for large audiences (count via HyperLogLog; exact list only for authors who open the viewers list, or paginated from an append log). Also the "tray" freshness and blocked-user filtering at read time.
- **Failures/trade-offs:** CDN outage → degrade quality/latency; view receipts are best-effort (eventual); deletion requests (privacy) must propagate to CDN via purge.

**Q5. ID generator comparison: model answer.**
| | (A) DB sequence with range allocation | (B) Snowflake-style |
|---|---|---|
| Availability | DB/coordination is a dependency, but nodes can serve from their local block for a while; need HA for allocation | Fully local generation; no runtime dependency (only machine-ID assignment at startup) |
| Ordering | Roughly ordered per block; globally unordered across nodes | Roughly time-ordered (k-sortable) |
| Capacity | Local block gives very high per-node rate | 4,096 IDs/ms/node = 4M/s per node (12-bit seq) → 500k/s across the fleet is trivial |
| Coordination | Needs the allocator; blocks are lost on crash (gaps) | Needs unique machine IDs (via ZooKeeper/config); clock discipline |
| Failure scenario | Allocator down → nodes exhaust blocks and stall (mitigate with large blocks, HA) | Clock jumps backward (NTP) → duplicates/collisions unless the generator refuses to issue (waits until clock catches up) or uses the last timestamp |
- **Recommendation:** (B) Snowflake for high-rate, decentralized, time-sortable IDs, with guardrails against clock regression; use (A) if strictly compact, sequential, gap-tolerant numbers are needed and you can protect the allocator. (Either is defensible with justification.)

**Q6. Synthesis: sample answer.**
- Chat and news feed share: (1) **hybrid fan-out** (push for small audiences, pull for huge ones), (2) **cursor-based pagination and per-entity ordering via sequence/time-sortable IDs**, (3) **cache + durable store separation** (Redis inboxes/feeds over a durable log), plus idempotent delivery.
- They *should differ*: chat requires strict per-conversation ordering and reliable no-loss delivery with acks (a message must never be dropped), whereas a feed tolerates approximate ordering, staleness, and dropped/regenerated items; so chat persists to a durable log first, while the feed cache can be rebuilt from source at any time.
- Any two cases with three well-justified shared patterns and one justified difference earns full credit.
