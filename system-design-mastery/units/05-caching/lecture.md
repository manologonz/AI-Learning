# Unit 5: Caching

## Objectives
By the end of this unit you can:
- Choose where to cache (client, CDN, in-process, distributed, DB buffer) and which pattern to use (cache-aside, read-through, write-through, write-back, write-around).
- Compute the effect of hit rate on latency and backend load.
- Predict and prevent the failure modes: stampede, penetration, hot keys, stale data races.
- Design invalidation and TTL strategies that match consistency requirements.

## Prerequisites
Units 3-4 (request path, storage and consistency).

## Core content

### 1. Why caching works and what it costs
Caches exploit **locality**: a small hot set gets most of the traffic (often Zipf-like; ~20% of items serve ~80% of requests). A cache trades **consistency and complexity for latency and backend load**. "There are only two hard things in computer science: cache invalidation and naming things."

**The core math**
- Backend load = `QPS × (1 − hit_rate)`. Going from 90% to 99% hit rate cuts backend load **10x** (10% → 1% misses), far more than the "9 points" suggests.
- Average latency = `h × t_cache + (1 − h) × (t_cache + t_backend)`. With t_cache = 1 ms, t_backend = 50 ms, h = 0.95 → 0.95×1 + 0.05×51 = 3.5 ms; at h = 0.5 → 26 ms.
- **The database must survive a cold cache.** If your DB only works because the cache absorbs 95% of reads, then a cache restart is an outage (Unit 10).

### 2. Where to cache (layers)
1. **Client/browser** (`Cache-Control`, ETag): free latency, but you can't invalidate it once served; use versioned URLs (`app.3f2a.js`) for immutable assets.
2. **CDN/edge**: static assets and cacheable API responses near users; purge APIs, `stale-while-revalidate`.
3. **API gateway / reverse proxy** (Nginx, Varnish): shared response cache.
4. **In-process** (a `Map`/LRU inside each app instance): nanosecond access, no network, but duplicated per instance, and instances disagree until TTL expires.
5. **Distributed cache** (Redis, Memcached): shared, sub-ms network hop, large capacity; a new failure domain.
6. **Database** (buffer pool, query cache): automatic; keep the working set in memory.

Use several layers, but each adds an invalidation problem. A common strong combo: CDN + small in-process LRU for the hottest keys + Redis + DB.

### 3. Caching patterns

| Pattern | Read path | Write path | Good for | Risks |
|---|---|---|---|---|
| **Cache-aside (lazy loading)** | App checks cache; on miss loads DB and populates cache | App writes DB, then invalidates (deletes) cache key | General purpose default | Stale reads in races; cold start; stampedes |
| **Read-through** | Cache library loads from DB on miss | (same as above) | Simplifies app code | Coupled to cache library |
| **Write-through** | Read from cache | Write to cache *and* DB synchronously | Read-heavy, need fresh cache | Slower writes; caches data never read |
| **Write-back (write-behind)** | Read from cache | Write to cache; flushed to DB asynchronously | Write-heavy, tolerable loss window | **Data loss if cache dies** before flush; complexity |
| **Write-around** | On miss from DB | Write to DB only; cache fills on read | Data written once, rarely read | First read after write is slow |

**Prefer invalidate (delete) over update on write** for cache-aside: updating the cache from app code races with other writers and can persist stale values indefinitely; deleting lets the next reader repopulate from the source of truth.

```ts
// Cache-aside with single-flight to avoid stampedes
const inflight = new Map<string, Promise<any>>();

async function getUser(id: string) {
  const key = `user:${id}`;
  const hit = await redis.get(key);
  if (hit) return JSON.parse(hit);

  if (!inflight.has(key)) {
    inflight.set(key, (async () => {
      const row = await db.query("SELECT * FROM users WHERE id=$1", [id]);
      // TTL + jitter so keys don't expire in sync
      const ttl = 300 + Math.floor(Math.random() * 60);
      await redis.set(key, JSON.stringify(row), "EX", ttl);
      return row;
    })().finally(() => inflight.delete(key)));
  }
  return inflight.get(key);
}

async function updateUser(id: string, patch: object) {
  await db.query("UPDATE users SET ... WHERE id=$1", [id]); // source of truth first
  await redis.del(`user:${id}`);                            // then invalidate
}
```

### 4. Eviction and TTL
- **LRU** (least recently used): default; good for recency-driven access. **LFU** (least frequently used): better when popularity is stable; resists scan pollution. **FIFO/Random**: simple. **TTL**: bounds staleness even if invalidation is missed; **always set a TTL** as a safety net.
- Redis `maxmemory-policy`: `allkeys-lru` for pure caches; `noeviction` for data stores (don't mix roles).
- Size the cache from the hot set (Unit 2): if 20% of a 500 GB daily-accessed set is hot → ~100 GB; consider item overhead (Redis ~50-100 B per key).

### 5. Failure modes (the real content of this unit)

1. **Cache stampede (dogpile / thundering herd).** A hot key expires; 10,000 concurrent requests all miss and hit the DB at once. Fixes: **single-flight/request coalescing** (one loader, others wait), **locks** (`SET lock NX`), **TTL jitter**, **stale-while-revalidate** (serve stale while one request refreshes), **probabilistic early refresh**, and **cache warming**.
2. **Cache penetration.** Requests for keys that *don't exist* (attackers or bugs) always miss and hit the DB. Fixes: **negative caching** (cache "not found" with short TTL), **Bloom filter** of valid keys, input validation.
3. **Cache avalanche.** Many keys expire simultaneously or the whole cache tier fails; the DB collapses. Fixes: TTL jitter, multi-tier caches, replicas/clustering, load shedding and DB-protective rate limits.
4. **Hot key.** One key gets so much traffic it saturates a single cache node (e.g., celebrity profile). Fixes: **local in-process cache with short TTL** in front of Redis, **key replication** (`key#1..k` with random suffix on reads), CDN, or splitting the value.
5. **Stale data races (cache-aside).** T1 reads DB (old value); T2 updates DB and deletes cache; T1 then writes its old value into the cache → stale until TTL. Mitigations: short TTLs, versioned values (write only if newer), delayed double-delete, or write-through with a single writer. Accept and document bounded staleness where possible.
6. **Inconsistency across layers/regions.** A user updates a profile and sees the old one from a stale layer. Fix: **read-your-writes** by bypassing/refreshing caches for the writer or pinning them briefly to the primary path.
7. **Cache as a dependency.** If the cache is down, requests fall through to the DB and may kill it. Design fallback behavior (circuit breaker, degrade features) intentionally.

## Worked examples

### Straightforward: hit-rate math
*Product pages: 40k QPS; DB can handle 5k QPS. What hit rate is required, and what happens at 85%?*
- Max miss rate = 5k/40k = 12.5% → need ≥ **87.5%** hit rate, ideally 95%+ for headroom.
- At 85%: DB sees 40k × 0.15 = 6k QPS > 5k → overloaded; queueing, timeouts, cascading failure. At 95% → 2k QPS (comfortable).

### Typical: choosing a strategy
*User profile service, reads 100:1 over writes, must reflect a user's own edits immediately, others can be a few seconds stale.* → **Cache-aside with delete-on-write** + 5-minute TTL with jitter, single-flight loading, and read-your-writes: after an edit, the API response returns the fresh object and the client uses it (or the writer's session bypasses cache for a few seconds). Others see the change after the delete/TTL. Add a tiny in-process LRU (1-2 s TTL) for the very hottest profiles.

### Edge case: the race that leaves stale data
*T1 (reader): cache miss → reads DB: `name = "Ann"`. T2 (writer): updates DB to "Anna", deletes the cache key. T1 (delayed): writes "Ann" into cache.* Cache now says "Ann" until TTL expires; the DB says "Anna". No error was logged. Mitigations: (1) rely on a modest TTL to bound the damage; (2) store a version/updated_at and use `SET` only if version is newer (Lua script/compare); (3) on write, delete, wait briefly (~500 ms), delete again ("delayed double delete"); (4) if strict, avoid cache for that data or use write-through with per-key ordering. **Lesson:** cache-aside gives bounded staleness, not linearizability.

## Common misconceptions

- **"A cache makes the system faster, period."** It moves work and adds staleness, invalidation logic, and a new failure mode. Measure hit rate and the p99 of misses.
- **"Update the cache on write."** Racy under concurrent writers; delete/invalidate and let readers repopulate.
- **"TTL alone solves invalidation."** TTL bounds staleness, but choosing it is a business trade-off (how stale can it be?).
- **"95% hit rate is great."** It depends: at 100k QPS, 5% misses = 5k QPS to the DB. Always translate hit rate into backend load.
- **"Redis is durable, so I can use it as the primary store."** By default it's a cache with optional persistence; know the loss window (RDB/AOF) before making it a system of record.
- **"Distributed caches remove the need for local caches."** Network hops to Redis (~0.5 ms) are 1000x memory; for hot keys a local cache matters.
- **"Caching fixes a slow query."** It hides it until a miss or cold start; fix the query/index too.

## Check your understanding

1. 50k QPS, cache hit rate 96%. DB load? *(50k × 4% = 2k QPS.)*
2. Why add jitter to TTLs? *(Prevents synchronized expiry → avalanche/stampede.)*
3. A key is requested 200k times/s and is stuck on one Redis node at 100k ops/s capacity. Fix? *(Local in-process cache with short TTL, or replicate the key across nodes/suffix shards.)*
4. Why prefer deleting a cache entry to updating it on writes? *(Avoids races between concurrent writers and stale overwrites; the next read loads the truth.)*

## Summary
- Caches trade consistency/complexity for latency/load; translate hit rate into backend QPS.
- Cache-aside with delete-on-write + TTL (with jitter) is the default; write-through/back and write-around fit specific patterns.
- Know the failure modes: stampede, penetration, avalanche, hot keys, stale races; design defenses up front.
- Layers multiply invalidation problems; be explicit about acceptable staleness per data type.
- The database must survive cache failure gracefully.
