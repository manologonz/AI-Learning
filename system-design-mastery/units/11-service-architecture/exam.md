# Unit 11 Exam: Service Architecture

Answers in `exam-answers.md`.

---

### Scenario (Q1-Q5): RetailCo

RetailCo has a 9-year-old Rails/Java-style monolith (one 2 TB Postgres) with modules: catalog, pricing, cart, checkout, orders, inventory, promotions, users, notifications, search. 90 engineers in 12 teams. Pain: deploys twice a week are blocked by test suites and cross-team merges; Black Friday overloads the whole app because catalog reads and checkout compete for DB and app capacity; checkout must remain up even if promotions are down; search relevance work needs different tech (Elasticsearch). Leadership wants "microservices."

**Q1 (evaluate, hard).** Decide which modules to extract first, second, and which to keep in the monolith for now. Justify with the pain points and boundary quality (data coupling, transactional needs).

**Q2 (apply, hard).** Design the extraction of `catalog` (read-heavy). Define the API/events, data strategy (who owns what tables), and how the monolith continues to work during the transition. Include the routing/verification/rollback plan.

**Q3 (analyze, hard).** Checkout calls promotions, pricing, inventory, users, and payments synchronously. Draw the dependency, compute the availability if each is 99.9% and all are serial, and propose changes so checkout survives promotions being down and stays under a 1.5 s p99 target.

**Q4 (evaluate, hard).** After extraction, teams propose a shared `common-lib` with domain models and a shared "master data" database for customer/product info "for consistency." Critique this and propose an alternative.

**Q5 (create, hard).** Define the prerequisites/platform capabilities RetailCo must have before extracting the *third* service, and a 6-month roadmap (quarter granularity) with success metrics for the migration.

---

### Independent items

**Q6 (MCQ, medium).** Conway's Law says:
A. Systems should be designed independently of organizations.
B. Organizations design systems that mirror their communication structures, so team boundaries and service boundaries interact.
C. Microservices always improve communication.
D. Larger teams produce better architecture.

**Q7 (MCQ, hard).** A team's "UserService" exposes CRUD on the users table and is called by 20 services on every request. What's the most likely problem?
A. It's a well-designed bounded context.
B. It's an entity service creating a central coupling point and availability bottleneck; better to expose business capabilities and distribute data via events/local replicas or tokens carrying needed claims.
C. It needs more replicas only.
D. It should be rewritten in another language.

**Q8 (short, hard).** Explain shadow traffic (dark launching) in a migration: what you compare, how you avoid side effects, and what can go wrong.

**Q9 (application, hard).** Two services (Orders, Inventory) need to stay in sync on stock levels. Design the interaction: sync call vs events, ownership, and how to handle a reservation that expires.

**Q10 (short, hard).** Give an example where you'd merge two microservices back together, the evidence you'd cite, and the process.
