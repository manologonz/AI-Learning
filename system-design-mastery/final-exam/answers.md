# Final Exam: Answer Key and Rubric

Each entry gives the correct answer, the reasoning, why distractors fail, and the **weak-unit map** (which unit to review if you missed it).

---

## Part A (3 points each; 36 total)

**A1: C (~4.3 minutes).** *Unit 1.*
- 0.0001 × 43,200 min = 4.32 min. A (52 min) is the *yearly* 99.99% figure (window confusion). B (43 min) is the 99.9% budget (one nine off). D (26 s) is 99.999%.

**A2: avg ≈ 57.9k QPS; peak ≈ 231k QPS; 231 servers.** *Unit 2.*
- Requests/day = 200M × 25 = 5B; ÷ 86,400 ≈ 57.9k; ×4 ≈ 231k peak.
- Per-server usable = 2,500 × 0.6 = 1,500 → 231k / 1,500 = 154 servers to carry peak.
- Survive loss of 1 of 3 AZs: 2 surviving AZs must hold 154 → 77 each → **231 total** (1.5×). Give credit for 154 (no AZ math) as partial (1.5 pts).

**A3: Idempotency key.** *Unit 3.*
- Client generates one key per *intent* and reuses it on retries. Server: (1) replay after success → return the stored response, no new transfer; (2) replay while in-flight → 409/425 "in progress" (or wait); (3) same key with a different body → 422 reject (key misuse). Store key + request hash + status + response atomically with a unique constraint; expire after ~24 h. (1 pt fix, 2 pts for the three behaviors.)

**A4: B.** *Unit 4.*
- The atomic guarded update removes the read-modify-write race (and prevents overdraft). A (Read Uncommitted) is weaker, C/D don't address concurrency. (Also valid in practice: `SELECT ... FOR UPDATE`, optimistic version check, Serializable with retry.)

**A5:** *Unit 5.*
- (a) Max miss rate = 8k/120k = 6.7% → hit rate ≥ **93.3%** (aim 97%+ for headroom).
- (b) **Cache stampede (dogpile)**. Defenses: single-flight/request coalescing (one loader); TTL jitter; stale-while-revalidate; distributed lock; proactive refresh for hot keys; in-process cache for the hot key. (1 pt each part; 1 pt for two defenses.)

**A6: B (W = 3, R = 3).** *Unit 6.*
- Overlap requires W + R > N: 3 + 3 = 6 > 5. A: 2+3=5 (not >5); C: 1+4=5; D: 4+1=5, not > 5, so a single-replica read can miss the write. Quorums still aren't linearizable under sloppy quorums/concurrent writes.

**A7:** *Unit 7.*
- (a) 2f+1 = 9 → tolerates **f = 4**. (b) Majority = **5**.
- (c) A process pause (GC/VM) can let the lease expire while the old leader still believes it holds it, so it may act after a new leader took over. Fix: **fencing tokens** (monotonic token per grant; the resource rejects lower tokens) or compare-and-set on the resource. (1 pt each.)

**A8:** *Unit 8.*
```
BEGIN;
  ins = INSERT INTO processed(event_id) VALUES ($id) ON CONFLICT DO NOTHING;
  if ins.rowCount == 1:
     INSERT INTO ledger(order_id, amount, event_id) VALUES (...);   -- side effect in same txn
COMMIT;
ACK(message);    -- after commit
```
- Ack after commit: a crash before ack leads to redelivery, which is deduped by `processed`; acking first would risk losing the message (at-most-once). Dedupe + effect must be atomic. (2 pts code, 1 pt ack reasoning.)

**A9: B.** *Unit 9.*
- Timeout = unknown outcome. A can cancel an order for a payment that actually succeeded; C/D are unsafe. Use status query/same idempotency key, and reconciliation as a safety net.

**A10: 81 attempts (3⁴); use retries at one layer with a budget.** *Unit 10.*
- Attempts multiply per layer. Recommended: retry only at the edge/client (1 retry), only idempotent/retryable errors (503/timeouts/429), full-jitter exponential backoff, a retry budget (≤10% extra traffic), deadline propagation, and no retries in inner layers (fail fast); add circuit breakers and load shedding.

**A11: B.** *Unit 11.*
- Shared schema + lockstep deploys + deep sync chains = all costs, no independence. A, C are normal; D: count alone doesn't define it.

**A12:** *Unit 12.*
- Allowed error rate = 0.1%. Observed 2.16% → **burn rate 21.6x**. Both 1 h and 5 min windows exceed 14.4 → **page immediately** (also exceeds 6x/6 h if sustained).
- One hour at 21.6x consumes 21.6 × (1 h / 720 h) = **3% of the 30-day budget**. (1 pt burn rate, 1 pt action, 1 pt budget fraction.)

---

## Part B (9 points each; 36 total)

**Scoring guide:** (a)-(e) items are worth ~2 pts each with partial credit; reward specific numbers, named techniques, and stated trade-offs.

### B1: Reliability incident
- **(a) Failure chain:** `FLUSHALL` → 0% hit rate → 80k QPS hits a 15k DB (5.3x overload) → latency ↑ → timeouts → **client + gateway retries (up to 9x)** amplify load → goodput collapses; queues hold expired requests. The system is **metastable**: retries + cold cache sustain overload after the trigger is gone. Consumer lag of 5M and rebalances cause **redelivery**; consumers weren't idempotent → double-processed orders. (Also: no DB protection, deploy at peak, no guardrail on destructive commands.)
- **(b) Recovery:** shed/cap traffic at the edge (e.g., admit ~50% of DB capacity), disable retries (or lower budgets), stop nonessential jobs, scale/recover the DB, warm the cache for hot keys (or serve stale/CDN), **ramp traffic gradually** (10% steps); for the consumer: pause/scale carefully, ensure dedupe (processed table) before resuming, then drain the backlog at a controlled rate; verify order duplicates and reconcile/refund.
- **(c) Design changes (prioritized):**
  1. Remove/ban `FLUSHALL`; use versioned keys; deploy guardrails/off-peak (process).
  2. Protect the DB: load shedding, concurrency limits, request coalescing (single-flight), circuit breaker (Unit 10, 5).
  3. Retry policy: one layer, budget, jitter, deadline propagation (Unit 10).
  4. Idempotent consumers with dedupe + ack-after-commit; bounded consumer lag alerts; DLQ (Unit 8).
  5. Cache resilience: replicated/persistent cache tier, warming, stale-while-revalidate, TTL jitter (Unit 5).
  6. Canary deploys with automatic rollback (Unit 12).
- **(d) SLO/alerts:** SLO on core journey (e.g., 99.9% success < 500 ms); burn-rate alerts (14.4x/1 h + 5 min pages; 6x/6 h); alerts on cache hit-rate drops, DB saturation (CPU/connections), retry rate, consumer lag/age, and duplicate-processing signals. Would have paged within ~5 minutes.

### B2: Data scaling (model)
- **(a) Numbers:** 40k writes/s vs a primary's practical ~5-15k complex writes/s (even with tuning); 6 TB/+100 GB/wk (5 TB/yr → ~11 TB in a year) strains a single node's storage/backup/restore/vacuum; 200k reads/s needs replicas/caching regardless. **First try:** index/query tuning, batching writes, connection pooling, read replicas + Redis cache for `order_id` reads, partitioning tables by time for maintenance, bigger hardware. Sharding is justified when *write throughput* (40k/s, 35% from 20 tenants) can't be met.
- **(b) Shard key:** primary key `hash(order_id)`-based routing for 80% of reads (`order_id` lookups hit one shard) but "my orders" needs `customer_id` locality. Two options: (1) shard by `customer_id` (hash) so both order placement and "my orders" are single-shard; `order_id` embeds the shard (e.g., ID contains shard/customer prefix) so lookups route without a directory; skew: give the 20 enterprise customers **dedicated shards or sub-shard** them by `(customer_id, hash(order_id) % k)` with merged reads. (2) shard by `order_id` + a global secondary index for `customer_id` (extra async index writes, staleness). **Recommended: (1)** with whale sub-sharding. Analytics: 5% by date → CDC into a warehouse/columnar store.
- **(c) Cross-shard:** "my orders" stays single-shard (except whales: fan-out to k sub-shards, cached); analytics off the OLTP path via CDC/stream to a warehouse; global uniqueness via order_id generator (Snowflake with shard bits).
- **(d) Online migration:** stand up new shards; **dual-write/CDC** from the old primary keyed by the new routing; backfill in throttled batches; verify with checksums/sampled comparisons; **shadow reads**; flip reads gradually by customer cohort, then writes (per-customer cutover with brief write pause/fencing token), keep reverse CDC for rollback; decommission after soak. Use an expand→migrate→contract approach for routing code (a routing layer/proxy reads a shard map).
- **(e) Failure modes:** hot shard (whale) → sub-shard/dedicated; resharding pain (fixed logical partitions, e.g., 1,024 virtual shards mapped to nodes so nodes can be added by moving partitions); single-shard outage affects a slice of customers → replicas + automated failover; cross-shard transaction needs (avoid; use sagas); router/shard-map SPOF → replicate/cached map.

### B3: Consistency & transactions (model)
- **(a)** Booking: **linearizable** (single authoritative writer/CAS per seat), **CP/PC-EC** posture: minority side refuses bookings during partitions. Seat browsing: eventual (cached, seconds stale), PA/EL. Order history: read-your-writes/session consistency (causal fine).
- **(b)** Saga: `createOrder(PENDING)` → `holdSeat` (TTL 10 min, atomic CAS on seat) → `authorizePayment` (voidable) → **pivot: `capturePayment` + `confirmSeat`** → retriable: `issueTicket`, `sendEmail`. Compensations: `releaseHold`, `voidAuthorization`, `cancelOrder` (idempotent, retried until success). Idempotency keys per step; timeouts = unknown outcome → status query; reconciliation between orders, payments, seat states; the orchestrator (Temporal-style) persists progress.
- **(c)** Partition seats by event; each event has a **home region** (leader in the region closest to most buyers; consensus/CAS for seat state within a region group); other regions read replicas (async) for browsing and forward booking requests to the home region (extra RTT only for remote buyers) or use optimistic UI (pending hold). Avoid global consensus on the hot path.
- **(d)** During a partition, regions not connected to the event's home region can't book that event's seats (CP) but can browse from replicas; the home region continues; holds expire safely on TTL; users see "try again" or a queue. Optionally allow *pre-allocated seat blocks* per region to keep selling within a block during partitions (relaxing global choice for availability, at a cost of stranded inventory).
- **(e)** Two regions both accept a booking for seat 12A, each writing `status=SOLD, buyer=X` with local timestamps; LWW merges after the partition and keeps one buyer, silently discarding the other, though both were told "confirmed" (violating uniqueness). Uniqueness needs coordination (CAS/single owner), not LWW.

### B4: Architecture critique (model)
- **(a)** Don't add a mesh or second region yet. Gather evidence: incident history and causes, DORA metrics, latency breakdown per call chain (needs tracing), dependency map, cost, on-call load, business SLO/availability needs (does the business need a second region? RTO/RPO?). A mesh adds complexity to a team of 25 with 3 on-call; a second region multiplies operational and data complexity. Prioritize: tracing + SLOs, deploy safety, decouple shared DBs and shallow the call chains, then revisit DR (a warm standby before active-active) with an RTO/RPO-driven case.
- **(b)** Problems → smallest fixes: (1) 2 shared databases → assign data ownership, revoke cross-service access gradually via APIs/events; (2) depth-6 sync chains → parallelize, collapse chatty calls, cache/replicate reference data, make non-critical calls async/optional with timeouts/breakers; (3) no tracing → OpenTelemetry + correlation IDs; (4) weekly manual deploys → CI/CD automation, canary + flags; (5) 18 services for 25 engineers → consolidate services that always change together into a modular monolith or fewer services; (6) on-call of 3 → ownership per service, runbooks, SLO-based alerting to cut noise.
- **(c)** "We choose to **consolidate** the five checkout-path services into one deployable with internal modules over keeping them separate, because the services change together, share data, and their synchronous chain costs latency/availability more than independent deploys benefit us. This holds while the team is one group owning checkout; we'd revisit if the team splits into independent groups with different scaling or compliance needs."
- **(d)** Metrics: deploy frequency (weekly → daily), change failure rate (and MTTR), checkout p99/availability vs SLO (and depth of call chain), on-call pages/week, plus cost or lead time for change.

---

## Part C: Capstone rubric (28 points)

Use the section weights in `exam.md`. Below is a **model outline** and the criteria for full credit. Any coherent alternative with justified trade-offs is acceptable.

**1. Requirements (3):** FR: create/edit/share docs, real-time co-editing, comments/suggestions, version history, search, permissions, media. NFR/SLOs: edit propagation p99 < 300 ms to collaborators in-region; document open p99 < 1.5 s; availability 99.95% (SLI: successful edit/open operations); durability of acknowledged edits (no loss, 11 nines for stored docs); permission revocation effective within seconds; residency for EU tenants; audit logs. Out of scope (v1): spreadsheets, offline-first mobile editing (or partial), advanced ML features.

**2. Estimates (3):** e.g., concurrent editors: 20M docs/day × ~3 editors × ~30 min ≈ peak concurrent editors ~ (20M×3×0.5h/24h≈1.25M average, peak ~3-4M) → **~4M concurrent WebSocket connections** (dominant resource: connections/memory + fan-out); ops/s: ~4M × 1 op/s (only while typing, maybe 30% active) ≈ 1.2M ops/s peak; storage: 1B docs × 200 KB = 200 TB text (+ history ×3-5 = ~1 PB) plus media: 300M × 2 MB = 600 TB; bandwidth: ops small (~100 B) → ~120 MB/s ingest; fan-out to avg 3 editors; media via CDN. (Accept different assumptions if consistent.) Identify: connection count and per-document ordering/fan-out as dominant.

**3. API/data model (4):** `documents(doc_id, tenant_id, owner, region, current_snapshot_ref, acl_ref, created_at)`; **operation log** per document (append-only, `(doc_id, seq, client_id, op)`) plus periodic **snapshots** (compaction); `permissions` as a relationship graph (ReBAC: user/group/org → doc/folder with roles); `comments(doc_id, thread_id, anchor, ...)`; media in object storage with CDN; search index (Elasticsearch) as derived. Stores: log/snapshots in a partitioned KV/wide-column or object store + metadata in Postgres/Spanner-like (sharded by tenant/doc); permissions in a strongly consistent store (Zanzibar-like) with caching. Justify by access pattern.

**4. Architecture (4):** clients → CDN/edge → gateway (auth, rate limits) → **real-time collab servers** (WebSocket; consistent-hash/route by `doc_id` so all editors of a doc land on the same "document room" owner) → operation log + snapshot store; async pipeline (Kafka) → search indexer, notification, audit, version-history compaction, analytics; permission service; media service; regional cells; control plane for routing (tenant → region/cell).

**5. Hard part (5):** choose **CRDT (e.g., Yjs/Automerge-like)** or **OT** with a central sequencer per document. Model answer: single authoritative *room owner* per document (lease + fencing token via etcd/consensus) sequences operations (OT/CRDT merge) and appends them durably (ack after quorum write) before broadcasting; clients apply optimistic local edits and reconcile; presence/cursors are ephemeral (not durable, coalesced); snapshots every N ops to bound replay; ownership failover: new owner loads latest snapshot + log tail; clients resend unacked ops idempotently (client_id, op_seq); **200 editors:** cap fan-out via batching/coalescing (e.g., 50-100 ms broadcast windows), send presence as sampled updates, split viewers vs editors (read-only viewers via pub/sub tiers). Offline edits: CRDTs merge naturally (OT needs rebasing). Hot documents: dedicate resources, backpressure/rate-limit per client. Trade-off: CRDT metadata growth (tombstones → garbage collection) vs OT central server dependency.

**6. Consistency/permissions/multi-tenancy (3):** edits: eventual convergence (CRDT) with per-doc ordering by the owner; **permission checks/revocation: strongly consistent** (revoked users disconnected and rejected on the next op; propagate via push + short-TTL caches, with fencing on the room); read-your-writes for the author; whales: dedicated cells/DBs, per-tenant quotas; residency: region-pinned tenants (EU cell), no cross-region replication of content (backups stay in-region), keep control-plane metadata minimal and non-sensitive; audit logs append-only in-region.

**7. Reliability (3):** collab server crash → ownership lease expires (e.g., 10 s), new owner takes over, clients reconnect with jittered backoff (avoid reconnect storms), unacked ops retried idempotently; region loss → warm standby/pilot light per residency zone with RPO ≈ seconds (async replication within the geography, e.g., two EU regions) and RTO ≈ 15-30 min; degrade: read-only mode if the log store is impaired, disable presence/comments first (load shedding priorities), search stale; retries: budget + jitter; deploys: cell-by-cell canaries, connection draining; SLO alerts on edit success/latency; test with game days (kill owners, partition regions).

**8. ADRs (2):** e.g., ADR-1: "Use CRDT (Yjs-style) with a per-document sequencer/persistence owner instead of OT" (alternatives: OT, last-writer-wins on paragraphs; revisit if metadata growth > X or if rich-text semantics needs exceed CRDT support). ADR-2: "Region-pinned tenant cells for residency and blast-radius control" (alternatives: global multi-region active-active DB; revisit if cross-region collaboration between tenants becomes core).

**9. Rollout/cost (1):** v1: single-region cells, no offline editing, basic comments; add suggestions/search relevance/regions later; cost drivers: connection servers (memory), storage of op history (compaction/TTL, tiering old versions to cold storage), media egress (CDN), search index size; levers: compaction, cold tiering, CDN hit rate, right-sizing, per-tenant quotas.

---

## Score interpretation and weak-unit map
- **≥ 90:** senior-level design fluency; practice by reviewing real designs at work and writing ADRs.
- **80-89:** strong; review the units tied to your lowest-scoring items.
- **70-79:** solid mid-to-senior; targeted re-teach on weak units, then re-attempt Part B.
- **< 70:** revisit units: A1 → U1; A2 → U2; A3 → U3; A4 → U4; A5 → U5; A6 → U6; A7 → U7; A8 → U8; A9 → U9; A10, B1 → U10; A11, B4 → U11; A12 → U12; B2 → U4/U6; B3 → U7/U9; Capstone → U13/U14. Then redo the unit exams before retaking.
