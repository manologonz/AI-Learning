# Unit 11 Exam: Answer Key

**Q1. Sample recommendation.**
- **First: `catalog` (read side) and `search`**: high read volume competing with checkout on Black Friday; read-mostly, derived data (search index is naturally separate, needs different tech), low transactional coupling. Also `notifications` (async by nature, bursty, low coupling) as an easy early win.
- **Second: `promotions`/`pricing`** as they need independent scaling/availability (checkout must survive promotions being down): design with a fallback (default price/no promo).
- **Keep in the monolith for now:** `checkout`, `orders`, `inventory`, `cart` (tightly coupled with transactional invariants, shared data; splitting requires sagas and is high-risk); `users` (central entity coupling: extract later as identity/auth, not CRUD).
- Rationale: extract along *pain + clean boundaries + low transactional coupling*; protect checkout by isolating the things that compete with it. Also invest in a modular monolith for remaining modules to speed deploys (module-level CI, ownership).

**Q2. Sample plan.**
- API: read endpoints (`GET /products/{id}`, `GET /categories/{id}/products`) plus events `ProductCreated/Updated/PriceChanged` on a topic (outbox from the monolith while it's still the writer).
- Data ownership: phase 1: monolith remains the system of record for product *writes* (admin tooling); catalog service holds a read model (own DB/cache) built via CDC/outbox events + initial backfill. Phase 2: writes move to the catalog service, and the monolith reads via API/events (tables owned by catalog; monolith access removed).
- Transition: routing layer sends `GET` product traffic to catalog by percentage; the monolith code path stays behind a flag.
- Verify: shadow reads compare responses (fields, price, availability) and latency; monitor divergence (event lag, freshness).
- Rollback: flip the route back to the monolith (still consistent since the monolith is still authoritative in phase 1); after stability and write migration, delete old code.

**Q3.**
- Serial availability: 0.999^5 ≈ 0.995 (99.5%): ~3.6 hours downtime/month; plus latency adds up.
- Changes: (1) make **promotions** optional: 100-200 ms timeout, circuit breaker, and fallback (no promo or cached promos: checkout proceeds at list price with an explicit business rule); (2) **pricing/users**: cache or local replicas/claims in the session token so they're not synchronous hops; (3) **inventory**: reserve via a fast local check or a reservation (TCC) with a short timeout; (4) call independent services **in parallel** rather than serially to fit the 1.5 s budget (latency = max not sum); (5) payments remains the critical, guarded dependency with idempotency and timeout budget; (6) use bulkheads per dependency.
- Result: critical path = checkout + payments (+ inventory), availability ~99.7-99.8%; optional dependencies degrade gracefully.

**Q4. Critique and alternative.**
- Shared domain-model library couples every service's release and model to the shared library version (lockstep upgrades, hidden coupling); a shared master database is a shared database (schema coupling, contention, no ownership, single point of failure): it recreates the monolith's problems while adding network hops.
- Alternative: each service owns its data and its own model (bounded contexts; "customer" means different things in different contexts), share *contracts* not models (schemas in a registry, generated clients), replicate reference data via events into local read models, define an authoritative owner per data type (e.g., Product master = catalog), and use small stable libraries only for cross-cutting concerns (logging, auth).

**Q5. Sample answer.**
- Prerequisites before service #3: CI/CD templates per service with automated tests and canary deploys; centralized logging with correlation IDs; distributed tracing; metrics/SLO dashboards and burn-rate alerts per service; on-call rotation and runbooks; service catalog/ownership; standardized health checks, config/secrets management; API gateway with auth/rate limits; schema registry and contract tests; per-service IaC.
- Roadmap: **Q1**: platform foundations (CI/CD, observability, gateway), modular-monolith cleanup; extract notifications (low risk). **Q2**: extract catalog+search behind strangler routing with shadow traffic; introduce outbox/CDC. **Q3 (months 5-6)**: promotions extraction with degradation design, load test for Black Friday (chaos: kill promotions), decommission old code.
- Success metrics: deploy frequency per team (e.g., from 2/week to daily), lead time for change, change failure rate, checkout p99 and availability during peak, DB CPU during peak, MTTR, number of cross-team release coordination events.

**Q6: B.**
- Conway's Law: system design mirrors communication structures; use the "inverse Conway maneuver" (shape teams to get desired architecture). A/C/D are incorrect.

**Q7: B.**
- Entity/CRUD services become central coupling points and single points of failure; prefer capabilities, replicated read models, or identity tokens/claims. A false; C only adds capacity, not resolving coupling; D irrelevant.

**Q8. Sample answer.**
- Duplicate a sample of production requests to the new service (shadow), discard its responses to the user, and compare with the old path's responses/latency/error rates offline.
- Avoid side effects: shadow instances run read-only or against sandbox/no-op adapters for writes (no emails, no charges), use idempotency/test flags, isolate data stores.
- Pitfalls: nondeterminism/timestamps causing false diffs, load doubling on downstream systems, privacy/PII copying, hidden side effects, drift in reference data, and false confidence if shadowed traffic doesn't cover rare paths.

**Q9. Sample answer.**
- Inventory owns stock counts; Orders never writes them. On order creation, Orders sends a **reserve** command (sync request with short timeout and an idempotency key, since Orders needs an immediate answer) or publishes `ReserveStock` and awaits `StockReserved`/`Rejected` (async, with an order in `PENDING`). Inventory creates a **reservation with TTL** (e.g., 15 min).
- Confirmation: after payment, Orders sends `ConfirmReservation` (decrement stock). If payment fails or the TTL expires, Inventory releases the reservation automatically and emits `ReservationExpired`; Orders handles it by cancelling or re-trying. Idempotent handlers everywhere; periodic reconciliation between reservations and orders.

**Q10. Sample answer.**
- Example: `pricing` and `promotions` services that always change together, share the same data, and require synchronous calls on every request with a lockstep release process.
- Evidence: deploy history showing coupled releases, incident postmortems for cross-service latency/failures, high call volume between them, frequent cross-team PRs, no independent scaling need.
- Process: define the merged bounded context/owner, merge codebases (start with a module boundary inside a single deployable), migrate consumers to the unified API, consolidate data stores carefully (one-way migration with verification), then decommission the extra service.
