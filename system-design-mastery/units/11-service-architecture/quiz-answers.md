# Unit 11 Quiz: Answer Key

**Q1: B.**
- A distributed monolith has distribution costs without independence. A is just horizontal scaling; C is a modular monolith (fine); D is unrelated.

**Q2: B.**
- Shared tables create hidden coupling. A (read-only access still couples schemas), C, D are wrong.

**Q3: B.**
- Boundaries are unclear and the team is small: cost of distribution outweighs benefits; a modular monolith keeps options open. A is an appeal to fashion; C ignores modularity; D swaps one default for another.

**Q4: ≈ 99.6%.**
- 0.999^4 ≈ 0.996. Each added synchronous hop reduces availability and adds latency; deep chains are fragile: keep them shallow, make non-critical calls async/optional, or replicate data locally.

**Q5: B.**
- Route incrementally through a facade, migrate functionality piece by piece, then delete the old code. A is a big-bang rewrite (the anti-pattern); C/D are unrelated.

**Q6. Sample answer.**
- Chatty synchronous calls between two services for one user action; lockstep deployments/changes across services; shared database tables; need for distributed transactions on common operations; "entity services" exposing CRUD.

**Q7: B.**
- A local replica (event-carried state) or cache removes the synchronous dependency and reduces latency; staleness is accepted for names. A hides the problem, C breaks ownership, D amplifies load.

**Q8. Sample answer.**
- **Expand:** add the new field as optional in the schema (registry backward-compatible) and start producing it; consumers ignore unknown fields.
- **Migrate:** each of the 6 teams updates their consumer to use it (track adoption via consumer-driven contract tests/metrics).
- **Contract:** once all consumers have migrated (verified), make the field required in the schema/validation or remove the old field. Never make it required at step 1.

**Q9. Sample answer.**
- Routing: gateway/proxy (or in-monolith interface + flag) that routes `/search` to old or new implementation by percentage/tenant.
- Data: the search service builds its own index from CDC/outbox events emitted by the monolith; initial backfill from a snapshot; the index is derived data, so no ownership conflict.
- Verify: shadow traffic (send queries to both, compare results/latency), relevance metrics; dual-run before cutover.
- Rollout: 1% → 10% → 100%, watching SLOs; rollback = flip the route back (old path remains until the new one is stable); then delete the old code.

**Q10. Sample answer.**
- Benefits: uniform mTLS/identity between services; consistent retries/timeouts/circuit-breaking policy and traffic shifting (canaries) without app changes; rich telemetry.
- Costs: added latency/CPU/memory per hop (sidecar), operational complexity and upgrade risk, and a steep learning curve; misconfigured policies (retries) can cause outages.
- Justified with dozens of services/polyglot stacks and platform team capacity; overkill for a handful of services (use a library + gateway).
