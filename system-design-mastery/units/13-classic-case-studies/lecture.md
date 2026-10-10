# Unit 13: Classic Case Studies

## Objectives
By the end of this unit you can:
- Apply the full design process (clarify → estimate → interface/data → high-level → deep-dive → failure/evolution) to unfamiliar prompts.
- Identify the "hard part" of each classic system and design specifically around it.
- Compare alternative designs using numbers and explicit trade-offs.
- Recognize which earlier building blocks (caching, partitioning, queues, sagas, rate limiting) each case exercises.

## Prerequisites
Units 1-12.

## How to read this unit
Each case follows the same template: **Requirements → Estimates → API/Data → Design → Hard part → Trade-offs & failure modes → Evolution.** Don't memorize the answers; memorize the *reasoning moves*. In the exam you'll get prompts you haven't seen.

---

## Case 1: URL Shortener (the warm-up)
**Requirements:** create short link for a long URL; redirect; custom aliases (optional); expiry; analytics (click counts). NFRs: redirect p99 < 50 ms, 99.99% availability for redirects, links are never wrong (a redirect to the wrong URL is worse than an outage).
**Estimates:** 100M new links/month ≈ 40 writes/s; read:write 100:1 → ~4k redirects/s avg, ~12k peak. Storage: 6B links over 5 years × ~500 B ≈ 3 TB. Tiny by our standards.
**API/Data:** `POST /links {url, alias?, ttl?} → {short}`; `GET /{code} → 301/302`. Table `links(code PK, url, created_at, expires_at, owner)`. KV/wide-column or Postgres.
**Design:** load balancer → stateless app → cache (Redis) → DB; CDN in front for hot redirects.
**Hard part: generating unique short codes.** Base62 with 7 chars = 62⁷ ≈ 3.5×10¹² ≫ 6B needed.
- *Hash the URL* (MD5/SHA → first 7 chars): collisions need handling (check-and-retry); same URL → same code (dedupe) but that leaks/joins users' links.
- *Counter/ID generator → base62*: no collisions, but predictable/enumerable; a central counter is a bottleneck/SPOF. Fix: **range allocation** (each app node reserves a block of 1,000 IDs from a coordination service such as ZooKeeper or a DB sequence) or **Snowflake-style IDs** (timestamp + machine ID + sequence).
- *Random 7 chars + unique constraint*: simple, non-guessable; retry on collision (rare).
```ts
const ALPHABET = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";
const toBase62 = (n: bigint) => { let s = ""; do { s = ALPHABET[Number(n % 62n)] + s; n /= 62n; } while (n > 0n); return s; };
```
**Trade-offs:** 301 (permanent, cacheable by browsers/CDN, but breaks analytics and expiry) vs 302 (temporary, every click hits you: needed for counting/expiry). Analytics via async events (Kafka) so redirects never wait on counting. Abuse: rate limits, malware scanning. **Reads dominate → cache-aside with high hit rate**; partition by `code` hash when data outgrows one node.

---

## Case 2: News Feed / Timeline
**Requirements:** post; follow; see a reverse-chronological (or ranked) feed. NFR: feed load p99 < 300 ms; new posts visible within seconds (eventual consistency OK); very read-heavy.
**Estimates:** 300M DAU, feed views 10/day each = 3B reads/day ≈ 35k/s avg (100k peak); posts 30M/day ≈ 350/s. Average 200 followers, but celebrities have 30M.
**Hard part: fan-out.**
- **Fan-out on write (push):** on post, insert the post ID into every follower's precomputed feed list (Redis list/sorted set per user). Reads are trivial (one lookup). Write amplification = followers: 200 avg fine; **30M for a celebrity is a disaster** (Unit 2's fan-out math).
- **Fan-out on read (pull):** at read time, fetch recent posts from everyone you follow and merge. Cheap writes; expensive reads (200 queries per feed load).
- **Hybrid (what real systems do):** push for normal users; **pull for celebrities** and merge the celebrity's recent posts into the feed at read time. Threshold by follower count (e.g., > 100k).
**Design:** Post service → DB (posts by user, partitioned by user_id) + event to a fan-out queue → fan-out workers write `feed:{user}` = sorted set of (post_id, ts) capped at ~800 entries (memory: 800 × 8 B ≈ 6.4 KB × 300M users ≈ **1.9 TB** in Redis, feasible sharded, and only for active users; inactive users' feeds are built lazily). Feed read: get IDs from the feed cache, **hydrate** post/user data via multi-get (from a post cache), merge celebrity posts, apply ranking, return a page with a cursor.
**Trade-offs/failure modes:** ranking (ML) adds latency and needs precomputed features; feed cache loss → rebuild from follow graph + posts (pull path as fallback); hot posts → cache; deleted/blocked content filtered at read time; fan-out lag under bursts → queue depth alerting; consistency (a user's own post should appear immediately: read-your-writes by inserting it into the author's own feed synchronously).

---

## Case 3: Chat / Messaging
**Requirements:** 1:1 and group messaging, delivery in order per conversation, offline delivery, read receipts, presence, history. NFR: delivery latency < 200 ms online, no message loss, per-conversation ordering.
**Estimates:** 50M DAU × 40 msgs = 2B/day ≈ 23k/s avg, ~70k peak; **10M concurrent WebSocket connections** at peak (Unit 2) → ~200 gateway servers.
**Design:** Clients hold WebSockets to **connection gateways** (stateful, holding user→connection maps in a **session registry**, Redis). Sending: client → gateway → **chat service** → persist message (partitioned by `conversation_id`, wide-column store, Unit 4 Q10) with a per-conversation **sequence number** assigned by the conversation's shard → publish to a broker → delivery workers look up each recipient's gateway in the registry and push; **offline users**: message stays in their inbox (durable) and is fetched on reconnect via `sync(since=seq)`; also push notification (APNs/FCM).
**Hard part:** ordering + reliable delivery + connection management.
- Ordering: single writer per conversation (partition by conversation), sequence numbers; clients render by sequence and gap-fill.
- Reliability: at-least-once delivery + client acks + dedupe by `message_id` (client-generated idempotency key); the server persists *before* acking to the sender.
- Groups: fan-out on write to member inboxes for small groups; for huge channels use read-time fetch (same hybrid as feeds).
- Presence: heartbeats with TTL in Redis; don't broadcast every change to everyone (subscribe to visible contacts only; coalesce).
**Failure modes:** gateway crash → 50k clients reconnect simultaneously (**reconnect storm**: jittered backoff, connection draining on deploys); registry stale → deliver via inbox on next sync; end-to-end encryption complicates server-side search/multi-device (key management).

---

## Case 4: Typeahead / Search Autocomplete
**Requirements:** suggest top-k completions for a prefix as the user types. NFR: p99 < 100 ms including network, suggestions reflect trends within ~an hour, high availability (degrade to no suggestions).
**Estimates:** 5B searches/day ≈ 58k/s avg; each keystroke fires a request (~5 per search) → ~290k/s avg, ~1M/s peak. Aggressive caching required.
**Design:** **Offline/nearline pipeline**: aggregate query logs (stream → counts per query per window) → build a **prefix → top-k** index (trie or a hash map from each prefix to its precomputed top-10 list). Serve from memory (sharded by prefix range or hash; replicate for throughput); **CDN/browser caching** for popular prefixes (short TTL, 1-5 min); client-side debounce (~100-150 ms) and cancel stale requests.
**Hard part:** latency at 1M QPS and freshness. Precompute top-k at each trie node (space vs compute), update by periodic rebuild + swap (hourly) with a small real-time layer for breaking trends; personalization/ranking as a re-ranking step over the candidates. Filter offensive/unsafe terms. **Trade-off:** precomputed lists are fast but stale; on-the-fly is fresh but slow; use a blend.

---

## Case 5: Ride-Hailing / Proximity Search
**Requirements:** riders request rides; match to nearby available drivers; live location; ETAs; trip lifecycle. NFR: match in seconds; location freshness ~4 s; high availability.
**Estimates:** 1M active drivers sending location every 4 s ≈ **250k writes/s**; each ~100 B → 25 MB/s (trivial bandwidth; the *rate* matters). Location data is ephemeral (we only need the latest), so it lives in memory, not a durable DB per update.
**Hard part:** geospatial indexing at high write rate. Options: **geohash** (string prefix = cell; neighbors need adjacent-cell queries at boundaries), **quadtree** (adaptive to density), **S2/H3** (hierarchical cells, good neighbor handling). Store `cell → set of drivers` in **in-memory sharded stores** (Redis GEO or a custom service), sharded by cell/region; query = the rider's cell + neighbor cells, filter by distance/availability, rank by ETA.
**Design:** driver app → location gateway → update in-memory index (+ Kafka stream for analytics/trip history sampling) → **matching service** (query nearby, offer to the top candidates, timeouts/reoffers) → **trip service** (state machine: requested → matched → en-route → in-trip → completed) with a durable DB and the **saga** pattern for payment. Dense cities → hot cells: sub-divide dynamically.
**Failure modes:** double-dispatch of a driver to two riders (use an atomic claim/lock on the driver in the state store), stale locations (TTL), surge traffic at events (pricing + admission control), split-brain across regions (keep a city in one region).

---

## Case 6: File Storage & Sync (Dropbox-style)
**Requirements:** upload/download files, sync across devices, sharing, version history. NFR: durability 11 nines, efficient bandwidth use, eventual consistency on sync with conflict handling.
**Design:** clients split files into **chunks** (~4 MB), compute a **content hash** per chunk; **metadata service** (relational, sharded by user/namespace) stores the file → chunk list mapping; **block store** (object storage) stores chunks by hash → **deduplication** across files/users, and **delta sync** (only changed chunks upload). Sync: clients get change notifications via long-poll/WebSocket from a **notification service** and pull the metadata change log (`changes since cursor`).
**Hard part:** consistency and conflicts, and efficiency. Concurrent edits on two devices → **conflict copy** ("file (Alice's conflicted copy)") using version vectors/base-revision checks (optimistic concurrency: upload succeeds only if `base_version == current_version`). Large uploads: resumable, parallel chunk upload, checksums; dedupe leaks existence of a file via hash (mitigate with per-user salts for sensitive cases).
**Estimates:** 500M users × 2 GB avg = 1 EB (1,000 PB) of logical data before replication. Real systems store far less physical data thanks to dedupe and most users being nearly empty, but it shows why storage cost dominates the design: tiering, erasure coding, and dedupe are core, not optimizations.

---

## Case 7: Payments & Ledger
**Requirements:** charge customers via a payment provider (PSP), record every money movement, refunds, reconciliation. NFR: **correctness above all**, exactly-once *effects*, auditability, high availability for the payment path, PCI scope minimization.
**Design:** `POST /payments` with an **idempotency key**; a payments service with a **state machine** (created → authorized → captured → settled; failed/voided/refunded), calls the PSP with the same key, persists the state and an **immutable double-entry ledger** (debit/credit rows summing to zero) in one DB transaction; emits events via **outbox**. Webhooks from the PSP are the source of truth for asynchronous outcomes (idempotent handlers, signature verification). **Reconciliation**: daily job comparing our ledger to PSP settlement reports.
**Hard part:** unknown outcomes (timeouts) and duplicates. On timeout, do not assume failure: query the PSP by idempotency key/reference or wait for the webhook; a payment in `unknown` state is resolved by polling/reconciliation. Never store card numbers (tokenize via the PSP). Ledger append-only; corrections via reversing entries. Strong consistency (single-leader SQL) per account/merchant; shard by merchant/account; multi-currency and rounding rules (store integer minor units).

---

## Case 8: Notification System & Rate Limiter Service (brief)
- **Notifications** (Unit 1's example): API → priority queues (transactional vs marketing) → per-channel workers (push/email/SMS) → provider adapters with retries, backoff, and rate limits per provider; user preferences and quiet hours; templating; dedupe by `(user, event_id)`; DLQ + reconciliation. Hard part: reliable delivery through flaky third parties without duplicates or losing critical messages.
- **Distributed rate limiter** (Unit 10): token bucket in Redis (Lua) keyed by API key; local pre-limits to reduce Redis load; fail-open vs closed policy; hot keys for large customers → shard counters (`key#0..k`) and sum approximately.

---

## The meta-process (use it every time)
1. **Clarify** functional/non-functional requirements and scope (write down what's out of scope).
2. **Estimate** QPS, storage, bandwidth, connections → identify the dominant resource.
3. **Define APIs and the data model** around the queries; choose stores by access pattern.
4. **Draw the simplest design** that meets the numbers.
5. **Deep-dive the hard part** (fan-out, ID generation, ordering, geo index, correctness).
6. **Stress it:** failures, hot keys, growth 10x, consistency anomalies, abuse.
7. **State trade-offs and reversal conditions.**

## Worked examples (reasoning moves)

### Straightforward: pick the ID scheme for the shortener
*Requirement: non-guessable codes, no coordination on the hot path, ≤ 7 chars.* → Random 7-char base62 with a unique constraint and retry (collision probability at 6B/3.5T ≈ 0.17% per attempt at full scale; trivial with retries). If enumeration doesn't matter and analytics needs sorted creation, use range-allocated counters.

### Typical: decide push vs pull for a follower graph
*Average 150 followers; 1% of accounts have > 100k; feed p99 300 ms.* Hybrid: push for accounts under 100k followers (fan-out ≈ 150 writes/post), pull for the big ones merged at read time (each user follows maybe ~5 celebrities → 5 extra lookups, cached). Numbers show pure push would cost 30M writes per celebrity post; pure pull would cost ~150 queries per read × 100k reads/s.

### Edge case: the duplicate charge after a timeout
*Payment API times out at 3 s; the client retries with a new idempotency key generated per HTTP attempt.* Double charge. Fix: the key belongs to the **user intent** (generated once when the user taps "Pay" and reused on retries); server dedupes on it; timeouts trigger status queries rather than fresh attempts.

## Common misconceptions
- **"Memorize the standard answer."** Interviewers/reviewers probe deviations; understand the constraints that produce each design.
- **"Start with the diagram."** Start with requirements and numbers; a diagram without them is decoration.
- **"More components = better design."** Every component must be justified by a requirement or number.
- **"Fan-out on write is always better."** It fails for high-follower accounts; hybrid is the answer.
- **"Geohash gives exact neighbors."** Neighbors across cell boundaries need extra cell queries.
- **"301 redirects are strictly better for a shortener."** They defeat analytics and expiry.
- **"Ledger = balance column."** A ledger is append-only entries; the balance is derived.

## Check your understanding
1. Why can't a URL shortener rely on the hash of the URL alone? *(Collisions and same-URL-same-code semantics; need collision handling.)*
2. What's the memory estimate for feed caches and why only for active users? *(~2 TB for 300M × 6.4 KB; inactive users are built lazily.)*
3. How does a chat system keep per-conversation order? *(Single writer/partition per conversation with sequence numbers; clients order by sequence.)*
4. Why is location data kept in memory rather than a durable DB write per ping? *(250k writes/s of ephemeral data; only the latest matters; durable sampling for history.)*

## Summary
- Same process, different hard parts: ID generation (shortener), fan-out (feed), ordering + connections (chat), latency at scale with precomputation (typeahead), geo index at high write rate (ride-hailing), chunking/dedupe/conflicts (file sync), correctness and unknown outcomes (payments).
- Use numbers to justify each component and to pick between push/pull, cache/compute, sync/async.
- Always finish by stress-testing the design (failures, hot keys, growth) and stating trade-offs.
