# Further resources: Unit 5

- **Amazon Builders' Library: "Caching challenges and strategies"** (article, aws.amazon.com/builders-library/caching-challenges-and-strategies) — a senior-engineer-level treatment of cold caches, bimodal behavior, and why the backend must survive without the cache.
- **Redis docs: "Client-side caching" and "Eviction policies"** (docs, redis.io) — the primary source for LRU/LFU configuration and server-assisted invalidation.
- **"Scaling Memcache at Facebook"** (paper, NSDI 2013) — real-world lessons on stampedes (leases), invalidation across regions, and hot keys; long but the section on leases is worth the read alone.
- **"Optimal Probabilistic Cache Stampede Prevention"** (Vattani et al., paper/summary) — the XFetch algorithm behind probabilistic early recomputation; a compact, advanced follow-up to single-flight.
- **MDN: HTTP caching** (docs, developer.mozilla.org/en-US/docs/Web/HTTP/Caching) — the reference for `Cache-Control`, ETag, and `stale-while-revalidate` when you design the CDN/browser layers.
