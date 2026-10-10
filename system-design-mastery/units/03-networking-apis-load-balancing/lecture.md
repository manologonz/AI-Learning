# Unit 3: Networking, APIs & Load Balancing

## Objectives
By the end of this unit you can:
- Trace a request end to end (DNS → CDN → load balancer → app → data) and say where latency and failure occur at each hop.
- Choose between REST, gRPC, GraphQL, and real-time transports (WebSocket, SSE, long polling) for a given need.
- Choose an L4 or L7 load balancer and an algorithm (round robin, least connections, hashing, power-of-two-choices) and explain the trade-offs.
- Design APIs that are idempotent, paginated with cursors, and versionable.

## Prerequisites
Units 1-2 (requirements, estimation).

## Core content

### 1. The life of a request
`Client → DNS → (CDN edge) → Load balancer → App server → Cache/DB → response`

Every hop has cost and a failure mode:

| Hop | Cost | What can go wrong |
|---|---|---|
| **DNS** | 1-100 ms first lookup, cached by TTL | Stale records after failover (TTL), resolver outages |
| **TCP handshake** | 1 RTT | SYN floods, packet loss → retransmit delays |
| **TLS handshake** | 1 RTT (TLS 1.3), 2 for 1.2; 0-RTT resumption | Cert expiry (a classic outage), CPU cost |
| **CDN edge** | Near the user | Stale cache, cache miss goes to origin |
| **Load balancer** | Sub-ms | Single point of failure if not redundant; bad health checks |
| **App → DB** | ~0.5 ms/hop in-DC | N+1 queries: 100 sequential queries = 50 ms of pure network |

**Why it matters:** a cold HTTPS request from another continent costs DNS + TCP + TLS + request = ~4 RTTs ≈ 500 ms before a byte of your response. Reuse connections (keep-alive, HTTP/2 multiplexing), put a CDN/edge nearby, and terminate TLS at the edge.

### 2. Protocols in one page
- **TCP**: reliable, ordered, congestion controlled; head-of-line blocking (one lost packet stalls the stream). **UDP**: no guarantees, low overhead (DNS, video, QUIC is built on it).
- **HTTP/1.1**: one request at a time per connection (browsers open ~6). **HTTP/2**: multiplexed streams over one TCP connection, header compression; still suffers TCP head-of-line blocking. **HTTP/3 (QUIC over UDP)**: independent streams, faster handshakes, better on lossy mobile networks.

### 3. API styles

| Style | Good for | Watch out |
|---|---|---|
| **REST/JSON over HTTP** | Public APIs, CRUD, cacheability (GET + HTTP caches), broad tooling | Over/under-fetching, chatty for complex screens |
| **gRPC (HTTP/2 + protobuf)** | Internal service-to-service: strict schemas, streaming, small payloads, codegen | Poor browser support without a proxy, harder to debug by eye |
| **GraphQL** | Many client types needing flexible shapes (mobile vs web); aggregating backends | N+1 resolver problems, hard to cache/rate-limit by cost, query-complexity attacks |

**Real-time options:**
- **Polling**: simple, wasteful; fine for low-frequency updates.
- **Long polling**: request held open until data or timeout; works everywhere; connection churn.
- **SSE (Server-Sent Events)**: one-way server→client over HTTP; auto-reconnect; great for feeds/notifications.
- **WebSocket**: full-duplex, low overhead; stateful long-lived connections (Unit 2's connection-count problem), needs sticky routing/graceful drain on deploys.

### 4. Designing good APIs (the senior details)

**Idempotency.** Networks retry. A `POST /payments` retried after a timeout must not charge twice. Pattern: client sends an `Idempotency-Key`; server stores the key with the response and returns the stored response on replay. `GET`, `PUT`, `DELETE` are idempotent by definition; `POST` is not unless you make it so.

```ts
// Idempotency-Key middleware sketch (Express-like)
app.post("/payments", async (req, res) => {
  const key = req.header("Idempotency-Key");
  if (!key) return res.status(400).json({ error: "Idempotency-Key required" });

  // Atomic claim: only one request wins the insert
  const claimed = await db.query(
    "INSERT INTO idem (key, status) VALUES ($1,'in_progress') ON CONFLICT DO NOTHING RETURNING key",
    [key]
  );
  if (claimed.rowCount === 0) {
    const prior = await db.query("SELECT status, response FROM idem WHERE key=$1", [key]);
    if (prior.rows[0].status === "in_progress") return res.status(409).json({ error: "in flight" });
    return res.status(200).json(prior.rows[0].response); // replay stored result
  }
  const result = await chargeCard(req.body);   // ideally in same txn / with provider's own key
  await db.query("UPDATE idem SET status='done', response=$2 WHERE key=$1", [key, result]);
  res.status(201).json(result);
});
```
(Production also needs key expiry, request-hash checking so the same key with different bodies is rejected, and crash recovery for stuck `in_progress` rows.)

**Pagination.**
- **Offset** (`?page=500&limit=20` → `OFFSET 10000`): simple, but the DB scans and discards 10,000 rows (slow deep pages), and results shift when rows are inserted/deleted (duplicates/skips).
- **Cursor/keyset** (`?after=<opaque cursor>` → `WHERE (created_at,id) < ($1,$2) ORDER BY created_at DESC, id DESC LIMIT 20`): stable and O(log n) with an index. Preferred for feeds and large datasets. Cost: no random access to "page 37".

**Versioning & compatibility.** Additive changes (new optional fields) are non-breaking; removing/renaming fields or changing semantics is breaking. Strategies: URL (`/v2/`), header, or evolve additively with deprecation. Clients you don't control (mobile apps in the wild) force long compatibility windows.

**Errors & retries.** Return meaningful status codes: 4xx = client's fault (don't retry unless 429/408), 5xx = server's fault (retryable, with backoff+jitter, Unit 10). Use `429` with `Retry-After` for rate limiting.

### 5. Load balancing

**L4 (transport)**: routes TCP/UDP connections by IP/port; very fast, no payload inspection; can't route by path or header. **L7 (application)**: understands HTTP; routes by host/path/header/cookie, terminates TLS, does retries, A/B and canary routing, WAF; more CPU, more features.

**Algorithms:**
- **Round robin / weighted**: simple; ignores actual load (a slow request pins a server).
- **Least connections / least outstanding requests**: adapts to uneven request cost.
- **Power of two random choices**: pick 2 servers at random, send to the less loaded; near-optimal balance with no global state. Widely used.
- **Hash-based (IP/session/key)**: gives affinity (same client → same server; caches stay warm); breaks distribution if keys are skewed; naive `hash % N` reshuffles almost everything when N changes (fixed by consistent hashing, Unit 6).
- **Sticky sessions** keep state on a server but hurt balance and failover; prefer stateless servers with state in a shared store.

**Health checks & redundancy.** Active checks (probe an endpoint) and passive checks (eject on errors/timeouts). A *shallow* check (process is up) vs a *deep* check (dependencies OK): deep checks can cause cascading ejection when a shared dependency blips (every server fails its check, LB removes all of them → total outage). Prefer shallow LB checks + circuit breakers/metrics for dependency health. LBs themselves need redundancy (multiple nodes, anycast/floating IP, or managed LB).

**Global routing.** DNS-based (GeoDNS, latency-based) sends users to a region, but DNS TTL delays failover; anycast IP routing with BGP steers to the nearest healthy PoP faster. CDNs terminate at the edge.

## Worked examples

### Straightforward: choosing an API style
*A mobile app needs a home screen combining user profile, 10 recent orders, and 5 recommendations.* With REST that's 3 calls (or one custom endpoint, "backend-for-frontend"); with GraphQL one query shaped by the client. For a single screen, a BFF endpoint is enough; GraphQL is justified when many clients need many different shapes. Between internal services (orders ↔ inventory), use gRPC for typed contracts and low overhead.

### Typical: idempotent payment retry
*Client posts a $50 charge, times out after 10s, retries.* Without an idempotency key the server might have processed the first request (response lost), so retry double-charges. With the key: first attempt inserts the key, charges, stores the response. Retry hits `ON CONFLICT`, finds `done`, returns the *stored* 201 without re-charging. If the first attempt is still `in_progress`, the retry gets 409 and the client backs off. The client must reuse the same key on retries and generate a new key for a *new* intent.

### Edge case: the deep health check that took down everything
*An LB health check hits `/health`, which queries the DB. The DB has a 20-second hiccup; every app instance fails its check; the LB ejects all 40 instances; even after the DB recovers, instances come back one by one and get overwhelmed by the returning traffic (thundering herd) and are ejected again.* A shallow liveness check plus graceful degradation (serve cached/partial data, fail specific endpoints) would have kept the fleet in rotation. Rule: an LB health check should answer "can this instance serve traffic?", not "is every dependency perfect?". Also configure a minimum healthy percentage ("panic mode": if <X% healthy, route to all).

## Common misconceptions

- **"HTTP/2 or 3 makes everything faster automatically."** They reduce connection overhead; they don't fix slow queries, and HTTP/2 still has TCP head-of-line blocking.
- **"GraphQL replaces REST."** It shifts complexity (caching, rate limiting, N+1 in resolvers) to the server; it's a tool for flexible clients, not a universal upgrade.
- **"Round robin balances load."** It balances *request counts*, not *work*; heterogeneous requests need least-connections or power-of-two-choices.
- **"Sticky sessions are fine at scale."** They create hot servers, complicate deploys and failover; make servers stateless where possible.
- **"Offset pagination is fine."** It degrades with depth and produces duplicates/gaps under concurrent writes.
- **"Idempotent means the same response."** It means the same *effect*; replays typically return the original result but the guarantee is about side effects.
- **"Health check = deep dependency check."** Deep checks can convert a partial failure into a total outage.
- **"DNS failover is instant."** Resolvers and clients cache beyond TTL; plan for minutes.

## Check your understanding

1. Cross-continent RTT is 150 ms. How long before the first response byte of a fresh HTTPS (TLS 1.3) request, ignoring server time? *(DNS ~1 RTT if not cached + TCP 1 RTT + TLS 1 RTT + request 1 RTT ≈ 4 × 150 = ~600 ms; with cached DNS ~450 ms.)*
2. Why is `hash(client_ip) % N` a bad way to pick a server when N changes? *(Almost all keys remap to different servers, wiping caches/affinity; consistent hashing fixes this.)*
3. Which LB type would you need to route `/api/*` and `/static/*` to different backends? *(L7.)*
4. When would you choose SSE over WebSocket? *(One-way server→client updates over plain HTTP, simple proxy support, auto-reconnect: feeds, notifications, progress bars.)*

## Summary
- A request crosses DNS, TCP, TLS, CDN, LB, app, and data; reuse connections and move work to the edge to cut RTTs.
- REST for public/CRUD, gRPC for internal, GraphQL for flexible multi-client shapes; SSE/WebSocket for realtime.
- Retries are inevitable: make writes idempotent with idempotency keys; use cursor pagination; evolve APIs additively.
- L4 = fast, dumb; L7 = smart, feature-rich. Prefer least-request or power-of-two-choices over round robin for uneven work.
- Keep health checks shallow; keep servers stateless; plan for DNS caching delays.
