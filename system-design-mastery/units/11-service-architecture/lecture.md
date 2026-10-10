# Unit 11: Service Architecture: Monolith to Microservices

## Objectives
By the end of this unit you can:
- Decide when to keep a (modular) monolith and when to split into services, with explicit criteria.
- Draw service boundaries using domain-driven design (bounded contexts) and data ownership.
- Choose sync vs async communication and design service contracts that evolve safely.
- Plan an incremental migration (strangler fig) with data migration and rollback.
- Explain the roles of an API gateway, BFF, and service mesh, and their costs.

## Prerequisites
Units 3, 8, 9 (APIs, messaging, distributed transactions).

## Core content

### 1. The real trade-off
Microservices trade **in-process simplicity for organizational scalability and independent deployability**, and you pay with distributed-systems complexity: network failures, latency, eventual consistency, versioned contracts, observability, and operational overhead (deployments, service discovery, security between services, on-call).

What you gain (when it works):
- **Independent deploy and scale**: teams ship without coordinating; scale hot components separately.
- **Team autonomy** (Conway's Law: system structure mirrors org structure).
- **Fault isolation** (if bulkheads/timeouts exist; not automatic).
- **Technology flexibility** (rarely as valuable as expected).

What you pay: every function call becomes a network call (latency ≈ 1000x, can fail, needs timeouts/retries/idempotency); no cross-service transactions (Unit 9); data duplication and eventual consistency; harder debugging (distributed tracing needed); contract versioning; platform and DevOps investment.

**Decision heuristic:**
- Team ≤ ~20-30 engineers, single product, unclear domain → **modular monolith**. Modules with enforced boundaries (separate packages, no cross-module DB access, interface-only calls) give 80% of the benefit and are easy to split later.
- Split when **you can name the pain**: deploy contention (teams blocking each other), a component with radically different scaling/availability/security/compliance needs, a stable and well-understood boundary, and org maturity (CI/CD, observability, on-call).
- **"Distributed monolith"** is the worst outcome: services that must deploy together, share a database, and call each other synchronously in long chains: all the costs, none of the benefits.

### 2. Finding boundaries (domain-driven design lite)
- **Bounded context**: a boundary within which a model and its language are consistent. "Customer" in Billing (account, payment method) differs from "Customer" in Support (tickets, tone). Each context owns its model.
- **Aggregate**: a cluster of entities treated as one consistency unit (`Order` with `OrderLines`); transactions don't cross aggregates. **Service ≈ one or a few bounded contexts.**
- **Signals of a good boundary**: high cohesion inside, low coupling across; changes tend to be localized; a natural team owner; clear data ownership; the interface is a small set of business-level operations rather than CRUD on tables.
- **Signals of a bad boundary**: chatty synchronous calls for a single user action, frequent lockstep changes across services, shared tables, "entity services" (`UserService` exposing CRUD used by everyone), needing distributed transactions for common operations.
- **Data ownership rule**: each service owns its data; others access it *only* via its API or events, never by reading its tables. Shared databases are hidden coupling.

### 3. Communication styles
- **Synchronous (REST/gRPC)**: simple request/response; needed when the caller needs the result now. Risks: temporal coupling, cascading failures, latency chains (availability multiplies: 5 services at 99.9% ≈ 99.5%). Mitigate with timeouts, circuit breakers, bulkheads (Unit 10), and avoid deep call chains (depth ≥ 3 is a smell).
- **Asynchronous (events/commands via broker)**: decoupled and resilient; eventual consistency; needs idempotency and schema discipline (Unit 8).
- **Data replication via events**: a service keeps a local read-only copy of data it needs from another service (built from its events) to avoid synchronous lookups: faster and more available, but stale and duplicated ("data on the outside").
- **API composition/BFF**: an aggregator calls services and shapes a response for a client type.
- **Contracts**: treat APIs and events as public interfaces. **Consumer-driven contract tests** (Pact) catch breaking changes in CI. Follow *expand → migrate → contract*: add new field/version, migrate consumers, remove old later. Postel's law: be liberal in what you accept (ignore unknown fields).

### 4. Infrastructure components
- **API gateway**: single entry point for external clients: auth, rate limiting, routing, TLS, request aggregation, protocol translation. Risk: becomes a monolithic bottleneck or a place where business logic accumulates. Keep it thin.
- **BFF (Backend for Frontend)**: one gateway/aggregator per client type (mobile, web) that shapes responses.
- **Service discovery**: how services find each other (DNS, Kubernetes services, Consul, client-side load balancing).
- **Service mesh** (Istio, Linkerd): sidecar proxies handle mTLS, retries, timeouts, traffic shifting, and telemetry uniformly. Benefits: consistent policy without app changes; costs: extra latency (~1-3 ms/hop), operational complexity, resource overhead. Worth it at dozens of services; overkill at five.
- **Configuration & secrets**, **feature flags**, **shared libraries vs sidecars**: shared libraries create coupling across languages and versions; keep them small.

### 5. Migrating: the strangler fig pattern
Never big-bang rewrite. Incrementally route functionality from the monolith to new services behind a **facade/proxy**:
1. Identify a seam (a module with a clean boundary and business value to extract).
2. Put a routing layer (gateway/proxy) in front.
3. Build the new service; **dual-run/shadow traffic** to compare outputs.
4. Migrate **data**: options: (a) shared DB temporarily (transitional, with strict table ownership), (b) sync via CDC/dual writes (with outbox) then cut over reads, then writes, (c) one-time migration for cold data. The database is the hardest part.
5. Shift traffic gradually (1% → 10% → 100%) with rollback.
6. Delete the old code (this step is what people forget).

**Branch by abstraction** (introduce an interface in the monolith, implement it twice, flip via flag) is the in-process version.

### 6. Operability prerequisites
Before splitting, have: automated CI/CD per service, centralized logs with **correlation IDs**, **distributed tracing**, metrics and SLOs per service (Unit 12), standardized health checks and deploy templates, on-call ownership, and a documented service catalog. Without these, microservices multiply pain.

## Worked examples

### Straightforward: monolith or services?
*A 12-engineer startup with one product, one Postgres database, deploying twice a week. Marketing wants "modern architecture."* Keep a **modular monolith**: modules (`billing`, `catalog`, `orders`) with enforced boundaries and no cross-module table access. Pain today: none from the monolith. Revisit when: deploy conflicts block teams weekly, or one module needs a different scaling/availability profile.

### Typical: extract a service by strangler fig
*The monolith's `notifications` module (email/SMS/push) is bursty, uses third-party APIs with rate limits, and causes noisy-neighbor CPU spikes.* Good candidate: clear boundary, async by nature, minimal shared data. Steps: define an event/command contract (`SendNotification`), publish from the monolith via outbox, build the notification service consuming it, shadow-send to a test sink and compare, cut over channel by channel (email first), keep the old code path behind a flag for rollback, then delete it. The service owns its own templates/preferences tables (copied via a one-time migration and kept in sync by events until the monolith stops writing them).

### Edge case: the distributed monolith
*A team split a monolith into 15 services sharing one DB schema. A "place order" request synchronously calls 9 services in a chain; releases require coordinating 6 teams; one service's slow query takes down checkout.* Diagnosis: shared database (no data ownership), deep synchronous chains, no bulkheads, wrong boundaries (entity services). Remediation: consolidate services that always change together into one (yes, merge back!), give each remaining service its own data, replace read lookups with event-carried state or local replicas, orchestrate the order flow with a saga, and add timeouts/breakers. Lesson: **splitting along technical layers or entities creates coupling; splitting along business capabilities and data ownership reduces it.** Merging services is a legitimate architectural move.

## Common misconceptions

- **"Microservices are more scalable."** Monoliths scale horizontally too; microservices mostly scale *teams* and let you scale hotspots independently.
- **"Each service must have its own language/stack."** Standardization lowers the cost; polyglot is a cost not a goal.
- **"Services can share a database as long as they use different tables."** Shared DBs couple schema changes and deployments; ownership must be strict.
- **"Smaller services are better."** Nano-services multiply network calls and operational load; size services by cohesion and team ownership, not lines of code.
- **"A service mesh makes services resilient by itself."** It provides plumbing (retries, timeouts, mTLS); you still need sensible policies and idempotency.
- **"We can rewrite the monolith in the background."** Big-bang rewrites frequently fail; strangle incrementally.
- **"An API gateway solves inter-service communication."** It's for edge traffic; east-west traffic needs its own strategy.

## Check your understanding

1. Name three signals that a boundary is wrong. *(Chatty sync calls for one user action, lockstep deploys, shared tables/needs cross-service transactions.)*
2. 5 serial synchronous services at 99.9% each: overall availability? *(0.999^5 ≈ 99.5%.)*
3. Why is the database the hardest part of a migration? *(Data ownership/transactions, dual writes, consistency, and cutover risk.)*
4. When is merging services the right move? *(When they always change/deploy together and share data: a distributed monolith.)*

## Summary
- Microservices trade simplicity for team autonomy and independent deploy/scale; start with a modular monolith.
- Draw boundaries around business capabilities/bounded contexts and data ownership; avoid entity services and shared DBs.
- Prefer async/events or local replicas over long sync chains; guard sync calls with timeouts and breakers.
- Contracts are public APIs: version additively; test with consumer-driven contracts.
- Migrate with the strangler fig: route, shadow, migrate data, shift traffic, delete old code.
