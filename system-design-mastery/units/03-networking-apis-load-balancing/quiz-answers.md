# Unit 3 Quiz: Answer Key

**Q1: D (POST).**
- POST creates/acts and repeated calls can repeat side effects. A is safe+idempotent. B (PUT replaces a resource with the same state each time) and C (DELETE: deleting twice leaves the same end state) are idempotent by definition, though the *response* may differ (404 on second delete).

**Q2: B (L7).**
- L7 understands HTTP: host, path, headers, cookies. L4 sees only IPs/ports. C is wrong: L4 can't inspect paths. D is wrong: DNS resolves names, not paths.

**Q3: B.**
- The server dedupes on the key and returns the stored result, making the retry safe. A is wrong: GET must not create orders. C is wrong: retries are necessary for reliability; removing them trades duplicates for lost requests. D is wrong: a longer timeout doesn't remove the ambiguity (client can't tell if the server processed it).

**Q4: Sample answer.**
1. **Performance:** offset scans and discards all skipped rows (`OFFSET 1,000,000` reads a million rows); keyset uses the index to jump to the position: O(log n).
2. **Consistency:** with concurrent inserts, offset pages shift, so users see duplicates or miss items; a cursor anchors to a specific row position.
- Also acceptable: works well with infinite scroll; cost: no random page jumps.

**Q5: B.**
- With modulo, going from N to N+1 changes the target for nearly every key (only ~1/(N+1) stay). Result: near-total cache miss storm hitting the DB.
- A is the property *consistent hashing* gives, not modulo. C is wrong. D is not how modulo works.

**Q6: B (SSE).**
- One-way, plain HTTP (proxy-friendly), automatic reconnect built into EventSource.
- A works but is overkill for one-way and more proxy/infrastructure friction. C isn't browser-native. D wastes resources and adds latency (10 requests/s per client).

**Q7: ~800 ms.**
- DNS lookup: 1 RTT = 200 ms (uncached, recursive resolver may add more). TCP handshake: 1 RTT = 200 ms. TLS 1.3: 1 RTT = 200 ms. HTTP request → first byte: 1 RTT = 200 ms. Total ≈ **4 × 200 = 800 ms**. (Lowers with connection reuse, DNS cache, TLS resumption/0-RTT, or a nearby edge.)
- Common error: forgetting DNS and TLS.

**Q8: B.**
- Deep check ties every instance's health to a shared dependency; when it stalls, all fail together and the LB has nothing to route to, even for endpoints that don't need the DB. Recovery can also stampede.
- A ignores the check coupling. C isn't automatic behavior. D: checks are not cached in a way that avoids this.
- Fix: shallow health checks, graceful degradation, a min-healthy/panic threshold.

**Q9. Sample answer.**
- `created_at` alone is non-unique: many rows can share a timestamp (bulk inserts, same-ms writes), so `WHERE created_at < $t` skips or repeats rows at the boundary.
- Use a compound cursor (created_at, id) with a total order: `WHERE (created_at, id) < ($1, $2) ORDER BY created_at DESC, id DESC LIMIT 20`, supported by an index on `(created_at DESC, id DESC)`. The cursor is an opaque, base64-encoded pair of the last row's values; the server returns `next_cursor`.

**Q10. Sample answer.**
- **Benefit:** protobuf is compact and fast to (de)serialize; schemas (`.proto`) enforce contracts and generate typed clients, eliminating drift.
- **Cost:** tooling/migration effort, harder debugging (binary), needs an HTTP/2-aware LB and a gateway for browser/public clients; teams must adopt schema governance.
- **Don't migrate if:** the services are mostly I/O-bound so serialization isn't the real bottleneck (verify the 30% claim), the API is public/browser-facing, or the team can't own the tooling; consider fixing drift with OpenAPI + contract tests first.
