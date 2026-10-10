# Backend Development (Node.js / TypeScript) — Course Syllabus

**Learner starting level:** Intermediate programmer, aiming for senior-level backend competence
**Goal:** Get a (senior-track) backend job: interview-ready fundamentals, real production skills, system design fluency
**Stack for examples:** Node.js + TypeScript, PostgreSQL, Redis
**Estimated total time:** ~60 hours

## Units

1. **How the Backend Actually Works** (~3h): Networking, HTTP, and the Node runtime
   - Objectives: trace a request from browser to DB and back (DNS, TCP, TLS, HTTP); explain the event loop, libuv, and why blocking it is fatal; distinguish HTTP/1.1, /2, /3 trade-offs; use curl and dev tools to debug
   - Prerequisites: none
2. **Designing and Building APIs** (~5h): REST, HTTP semantics, and API design
   - Objectives: use methods/status codes/idempotency correctly; design resources, pagination, filtering, versioning; build a typed server (Fastify/Express); compare REST, GraphQL, gRPC
   - Prerequisites: 1
3. **Relational Databases and SQL** (~6h)
   - Objectives: model schemas and normalize; write joins, aggregates, CTEs, window functions; read EXPLAIN plans; design indexes
   - Prerequisites: 1
4. **Transactions, Isolation, and Data Access Patterns** (~5h)
   - Objectives: explain ACID and isolation levels and anomalies; avoid race conditions (locking, optimistic concurrency); use ORMs/query builders, avoid N+1; run migrations safely
   - Prerequisites: 3
5. **Authentication and Authorization** (~5h)
   - Objectives: hash passwords properly; sessions vs JWT trade-offs; OAuth2/OIDC flows; RBAC/ABAC; CSRF/CORS
   - Prerequisites: 2
6. **Application Security** (~4h)
   - Objectives: OWASP Top 10; injection, SSRF, XSS, broken access control; secrets management; input validation; rate limiting
   - Prerequisites: 2, 5
7. **Architecture, Structure, and Testing** (~5h)
   - Objectives: layered/hexagonal design, dependency injection; error handling; config; unit/integration/contract tests; test DBs
   - Prerequisites: 2, 4
8. **Caching and Redis** (~4h)
   - Objectives: cache-aside/write-through/etc.; invalidation and stampede; HTTP caching, CDNs; Redis data structures
   - Prerequisites: 3
9. **Async Work: Queues, Events, and Background Jobs** (~5h)
   - Objectives: queues vs pub/sub vs streams; at-least-once and idempotent consumers; retries, dead-letter, outbox pattern; Kafka/RabbitMQ/BullMQ basics
   - Prerequisites: 4, 8
10. **Concurrency and Performance in Node** (~4h)
    - Objectives: worker threads, clustering, streams and backpressure; profiling; connection pools; load testing
    - Prerequisites: 1, 3
11. **Distributed Systems Fundamentals** (~5h)
    - Objectives: CAP/PACELC, consistency models; replication and sharding; timeouts, retries, circuit breakers; sagas; consensus intuition
    - Prerequisites: 4, 9
12. **Deployment, Containers, and CI/CD** (~4h)
    - Objectives: Dockerize properly; environments, 12-factor; CI pipelines; zero-downtime deploys, migrations; Kubernetes concepts
    - Prerequisites: 7
13. **Observability and Reliability** (~3h)
    - Objectives: structured logging, metrics, tracing; SLIs/SLOs; alerting; incident debugging; graceful shutdown, health checks
    - Prerequisites: 12
14. **System Design and Interview Mastery** (~5h)
    - Objectives: run a system design interview framework; capacity estimation; design URL shortener, feed, chat, rate limiter; NoSQL/search/blob storage choices; discuss trade-offs like a senior
    - Prerequisites: all previous

## How this course works
Each unit: **lecture** (with 3-tier worked examples and misconceptions) → **quiz** (low stakes) → **exam** (synthesis) → **resources**. A cumulative **final exam** with a capstone design/build task closes the course. Scores and weak objectives are tracked in `progress.md`, and gaps get targeted re-teaching. I'll build ahead of you, one unit at a time.
