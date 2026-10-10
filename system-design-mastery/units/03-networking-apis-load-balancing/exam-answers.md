# Unit 3 Exam: Answer Key

**Q1.**
- Cause: a cold HTTPS request costs ~4 RTTs (DNS, TCP, TLS, request) = ~4 × 220 = **~880 ms** before the HTML's first byte; then the page needs many more requests (JS, CSS, API calls, images) across the same long link, often serialized by dependencies, easily reaching ~3 s. All bytes travel to us-east.
- Mitigations by impact: (1) **CDN/edge** for static assets and cacheable pages, with TLS terminated at the edge (cuts most RTT-bound work); (2) **connection reuse** (HTTP/2/3, keep-alive, edge-to-origin persistent connections) and fewer round trips (bundle/preload, fewer sequential API calls, BFF endpoint); (3) **regional deployment** (read replicas/app fleet in Asia) for dynamic content, with GeoDNS/anycast routing. Also acceptable: TLS 1.3 + 0-RTT/session resumption, DNS prefetch.

**Q2.**
- Root cause: a non-idempotent `POST /charges` was retried after a timeout; the server had processed the first request but the response was lost. The client can't distinguish "not received" from "not processed," so the retry executed the charge again.
- Fix: client generates a unique **Idempotency-Key** per logical payment intent and reuses it on retries. Server stores (key, request hash, status, response) with a unique constraint and atomically claims the key before charging.
  - (a) Replay after success: return the stored response (same status/body), no new charge.
  - (b) Replay while in progress: return `409`/`425`-style "in flight" (or block briefly) so the client retries later.
  - (c) Same key, different body: reject with `422`/`400` (key reuse with a mismatched request hash indicates a client bug).
- Also pass a derived key to the payment provider so the provider dedupes too; expire keys after ~24h.

**Q3.**
- Likely flaw: **deep health checks** (checking the DB) on every instance, so when the DB primary failed over, all instances failed health checks and the LB removed the whole fleet, taking down even DB-independent pages. Recovery was worse because instances returned all at once, causing a thundering herd against a cold DB/connection pool, and possibly re-tripping health checks; plus DNS/connection caches on the DB endpoint could slow reconnection.
- Fixes: (1) shallow liveness checks at the LB; (2) LB "panic threshold"/min-healthy so it never ejects everything; (3) graceful degradation: only checkout depends on the DB; other endpoints serve cached/read-only data; (4) circuit breakers + timeouts to fail fast on DB calls; (5) connection pool warm-up and jittered reconnect/slow-start on the LB to avoid herds; (6) automated DB failover with a stable endpoint (proxy/virtual IP) and low DNS TTL.

**Q4.**
- Offset scans and discards n rows: 50k-order users hit `OFFSET 49,980`, which reads ~50k rows for one page (O(n) per page), so deep pages are slow, plus the total for a loop through all pages is O(n²).
- Keyset SQL: `SELECT * FROM orders WHERE user_id=$1 AND (created_at, id) < ($2,$3) ORDER BY created_at DESC, id DESC LIMIT 20;` with an index on `(user_id, created_at DESC, id DESC)`.
- API evolution: add `?cursor=` and return `next_cursor` (opaque base64 of the last row's `(created_at,id)`) in the response; keep `?page=` working for old clients but cap it (e.g., max page 100, else 400 with a message), mark it deprecated, and use metrics to see when old clients age out (mobile apps need a long window). Additive change; no breaking removal.

**Q5.**
- **Polling** (every 10 s): 2M users → 200k QPS of mostly "no change" responses: wasteful. **WebSocket:** bidirectional is unnecessary; 2M long-lived connections stateful, deploys require draining. **SSE:** one-way, plain HTTP, auto-reconnect: best fit. 
- Estimate: 2M concurrent connections ÷ ~50k per server = ~40 servers (plus headroom); updates are rare so bandwidth is trivial. Alternatively hybrid: SSE for the active-order screen, push notifications when the app is backgrounded.
- Operational concern (any): deploys must drain connections gradually (avoid reconnect storms: client jittered reconnect); proxy idle timeouts need periodic heartbeats/comments; use `Last-Event-ID` for resume; scale by subscribing connection servers to a pub/sub of order events.

**Q6: B.**
- Power of two choices uses local counts to pick the less-loaded of two random backends: near-optimal with no shared state.
- A ignores request cost (slow requests pile up). C creates affinity and skew, and ignores cost. D is like round robin: balances counts, not work.

**Q7: B.**
- Sticky sessions couple users to servers: hot spots, session loss on crash/deploy, and awkward autoscaling. Keep sessions in Redis/DB or signed tokens; servers stateless.
- A: no relation. C: cookie leaking is a separate issue. D: it doesn't disable HTTP/2.

**Q8: B.**
- (i) public/partner: REST/JSON for ubiquity, cacheability, tooling. (ii) internal high-QPS: gRPC for efficiency and typed contracts. (iii) aggregation: BFF or GraphQL to avoid six round trips from the mobile client.
- A puts gRPC on a public API and WebSocket for aggregation (wrong tool). C, D misuse SSE/GraphQL/gRPC for the wrong problems.

**Q9. Sample answer.**
- L7 terminating TLS can inspect headers/paths, route, add WAF rules, apply retries and rate limiting, centralize certificates, and compress/offload crypto from app servers. Without termination the LB sees only encrypted bytes and can only route by IP/port/SNI.
- Passthrough when end-to-end encryption is required by policy/compliance (the LB must never see plaintext), or when clients use mutual TLS the backend must validate itself.

**Q10. Sample answer** (budget ≈ 600 ms):
- DNS: anycast/GeoDNS resolver ~20 ms (nearby). TCP+TLS to **Tokyo CDN edge**: 2 RTTs × 10 ms = 20 ms. Edge cache hit for static HTML/JS: 10 ms → first byte ~50 ms. For a dynamic API call on a miss: edge→regional origin (Tokyo region) over a warm persistent connection: 1 RTT ~10 ms + LB 1 ms + service 20 ms + Redis 1 ms + DB read replica 5 ms ≈ 40 ms; if the request needs the primary in us-east (writes): +220 ms.
- Total for first byte: ~50-100 ms for static; ~150 ms for cached dynamic; ~400 ms for a write to us-east; all under 600 ms.
- Eliminated: long-haul DNS/TCP/TLS handshakes (now to nearby edge), origin round trips for static content, cold TLS to us-east. Full credit: budgets add up, and the hops that remain cross-ocean are explicitly identified as writes.
