# Unit 8 Exam: Async Messaging & Event-Driven Architecture

Answers in `exam-answers.md`.

---

### Scenario (Q1-Q6): FoodRush

FoodRush processes ~5k orders/s at peak. On `POST /orders` the API currently (in one request) writes the order to Postgres, charges the card via an external API, calls the restaurant service, sends an SMS, and updates a search index. Incidents: the SMS provider went down and orders failed; sometimes the order exists but the restaurant never got notified; customers occasionally get two SMSs; during a promo the restaurant service slowed and API latency hit 30 s.

**Q1 (analyze, medium).** List the coupling problems in the current design and classify which steps must be synchronous vs can be async, and why.

**Q2 (apply, hard).** Redesign the flow using the outbox pattern and a Kafka topic. Draw the sequence in text: what happens in the API transaction, who publishes, and which consumers exist. Include partition key choice and why.

**Q3 (apply, hard).** Customers get two SMSs: explain two different causes in the new design and write the consumer logic that makes SMS sending effectively-once, including the limits of what you can guarantee with an external provider.

**Q4 (analyze, hard).** The restaurant-notification consumer starts failing for one restaurant (its endpoint returns 500 for 20 minutes). Describe the failure mode for the partition, and design retry/DLQ/backoff so other restaurants aren't affected.

**Q5 (evaluate, hard).** Traffic doubles; consumer lag on `orders.placed` rises. Diagnose systematically (what you'd measure, in order), and list scaling options with their constraints given 24 partitions and 30 consumer instances.

**Q6 (evaluate, hard).** Product wants an "order timeline" (placed → accepted → preparing → picked up → delivered) with a full audit trail and the ability to rebuild analytics. Would you use event sourcing for the order aggregate, or an outbox with a status column and history table? Recommend and justify with trade-offs and a reversal condition.

---

### Independent items

**Q7 (MCQ, medium).** Which statement about Kafka's exactly-once semantics (EOS) is most accurate?
A. It guarantees any external side effect (email, HTTP call) happens exactly once.
B. It provides exactly-once processing for consume-transform-produce within Kafka using idempotent producers and transactions; external side effects still need idempotency.
C. It removes the need for consumer offsets.
D. It makes ordering global across partitions.

**Q8 (MCQ, hard).** SQS visibility timeout is 30 s, but an image-processing task takes up to 5 minutes. What happens?
A. Nothing; SQS extends automatically.
B. The message becomes visible again after 30 s and may be processed by another worker concurrently: duplicates; fix by increasing the timeout/extending it with heartbeats and making processing idempotent.
C. The message is lost.
D. The message goes to the DLQ immediately.

**Q9 (short, hard).** Compare choreography and orchestration for the order flow (payment, restaurant, delivery). Which would you pick for 5 steps with compensations, and why?

**Q10 (create, hard).** Define the event schema (fields, versioning approach) for `OrderPlaced` and describe two backward-compatible and two breaking changes, and how you'd roll out a breaking change safely.
