# Unit 13 Quiz: Answer Key

**Q1: D.**
- 62⁶ ≈ 5.7×10¹⁰ already exceeds 6×10⁹ by ~10x, and 62⁷ ≈ 3.5×10¹² gives ~600x headroom (important for random-code collision rates). A (62⁵ ≈ 9.2×10⁸) is too small. Both B and C suffice; 7 is the common choice for headroom and low collision probability with random generation.

**Q2: B.**
- Write amplification: one post → 30M inserts. A/D describe pull-side costs, C is nonsense.

**Q3. Sample answer.**
- Push (fan-out on write) for most authors: on post, insert into followers' precomputed feeds → cheap reads. Pull (fan-out on read) for high-follower authors ("celebrities"): don't fan out; at read time fetch their recent posts and merge into the precomputed feed.
- Criterion: follower count threshold (e.g., > 100k) or an estimated fan-out cost vs read-time cost; may also consider posting rate and follower activity.

**Q4: B.**
- A single writer per conversation assigns sequence numbers. A: client clocks are unreliable. C: random UUIDs have no order. D: TCP orders within one connection only; multi-device/multi-server paths break it.

**Q5: 250,000 writes/s.**
- 1,000,000 / 4 = 250k/s. Each update supersedes the previous one; durable relational writes at this rate are wasteful (and expensive: WAL, indexes) for ephemeral data. Keep the latest location in an in-memory geo-index (Redis/custom), sharded by cell; durably sample points (e.g., every 30 s) via a stream for history/analytics.

**Q6: B.**
- 301 = permanent; clients/CDNs cache it and skip your service. For analytics/expiry/edits you need each click to hit you (302, or a 301 with short cache headers, though behavior varies). A is false, C nonsense, D false.

**Q7. Sample answer.**
- Treat the outcome as **unknown**. Do: mark the payment `pending/unknown`, query the PSP by our idempotency key/reference (status API) or wait for the webhook; retry the *same* request with the *same* idempotency key if needed; resolve via reconciliation.
- Must not: assume failure and create a new charge with a new key, or tell the customer "failed" then charge anyway; must not mark it failed and release goods before knowing.

**Q8. Sample answer.**
- `chunks(hash PK, size, storage_key)` in object storage (content-addressed); `files(file_id, namespace_id, path, current_version)`; `file_versions(file_id, version, chunk_list[hashes], size, mtime, author_device)`; `changes(namespace_id, seq, file_id, version)` change log with a monotonically increasing sequence per namespace.
- Dedupe: client hashes chunks; asks the server which hashes are missing; uploads only missing chunks (across users/files).
- Conflict detection: an upload of a new version carries `base_version`; the server accepts if `base_version == current_version` (compare-and-set); otherwise it stores it as a conflict copy and notifies the client.

**Q9. Sample answer.**
- Precomputed prefix → top-k lists served from memory; aggressive caching at CDN/browser for popular prefixes (short TTL); client debounce + cancel in-flight requests; sharded in-memory replicas close to users; limit response size (top 5-10).
- Freshness trade-off: precomputed lists are rebuilt hourly (stale for trending queries); add a small real-time layer/boost for breaking trends at extra complexity.

**Q10. Sample answer.**
- Auditability/immutability (every movement has a permanent record; corrections are reversing entries); easier reconciliation and debugging (balance is derived and can be recomputed/verified; entries trace to events); also concurrency friendliness (appends vs contended updates) and built-in invariants (debits = credits).
