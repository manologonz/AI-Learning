# Unit 11 Quiz: Service Architecture

Answers in `quiz-answers.md`.

**Q1 (MCQ, easy).** Which best describes a "distributed monolith"?
A. A monolith deployed on several servers behind a load balancer.
B. Multiple services that must be deployed together, share a database, and call each other in long synchronous chains: all the costs of distribution with few benefits.
C. A monolith with many modules.
D. A serverless architecture.

**Q2 (MCQ, easy).** Which rule best describes data ownership between microservices?
A. Any service may read any table if it's read-only.
B. Each service owns its data and others access it only via its API or events.
C. All services share one schema for consistency.
D. Data ownership is decided by the DBA team.

**Q3 (MCQ, medium).** A 10-engineer team building one product with unclear domain boundaries is considering 12 microservices. What's the strongest recommendation?
A. Go ahead; microservices are the industry standard.
B. Start with a modular monolith with enforced module boundaries and split only when specific pain (deploy contention, divergent scaling/compliance needs) appears.
C. Use a single giant service without modules.
D. Use serverless functions for everything.

**Q4 (short, medium).** Compute the availability of a request that synchronously calls 4 services in series, each 99.9% available. What does this suggest about long synchronous chains?

**Q5 (MCQ, medium).** What is the strangler fig pattern?
A. Rewriting the system entirely and switching over on a set date.
B. Incrementally routing functionality to new services behind a facade until the old system can be retired.
C. Running two copies of the monolith.
D. Freezing feature development.

**Q6 (short, medium).** Name two signals that a service boundary has been drawn in the wrong place.

**Q7 (MCQ, hard).** A service needs customer names from the Customer service on every request (50k QPS), and the Customer service occasionally slows down. Which approach best improves resilience and latency?
A. Increase the timeout to 30 s.
B. Maintain a local read-only copy of the needed customer fields, updated via Customer events (accepting eventual consistency), or cache aggressively.
C. Read directly from the Customer database.
D. Retry more aggressively.

**Q8 (short, hard).** Describe how you'd add a required field to an event consumed by 6 teams without breaking anyone (expand → migrate → contract).

**Q9 (application, hard).** Sketch the strangler-fig plan to extract a `search` module from a monolith: the routing layer, data strategy, verification, and rollback.

**Q10 (short, hard).** When is a service mesh justified and when is it overkill? Name two benefits and two costs.
