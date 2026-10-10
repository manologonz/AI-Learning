# Unit 8 Quiz: Async Messaging & Event-Driven Architecture

Answers in `quiz-answers.md`.

**Q1 (MCQ, easy).** Which best describes at-least-once delivery?
A. Messages may be lost but are never duplicated.
B. Messages are never lost but may be delivered more than once.
C. Each message is delivered exactly once by the network.
D. Messages are delivered in random order but always once.

**Q2 (MCQ, easy).** You need multiple independent systems (search, analytics, fraud) to each consume every order event, and to replay history for a new consumer. Which fits best?
A. A single work queue with competing consumers
B. A log-based broker like Kafka with a consumer group per system
C. Direct HTTP calls to each system
D. A cron job

**Q3 (MCQ, medium).** A consumer acks the message first and then processes it; the process crashes mid-way. Result?
A. The message is redelivered and processed.
B. The message is lost (at-most-once behavior for that processing).
C. The broker retries automatically because acks are ignored.
D. The message is duplicated.

**Q4 (short, medium).** Explain the dual-write problem in one example and how the transactional outbox resolves it.

**Q5 (application, medium).** A topic has 8 partitions and a consumer group of 12 instances. How many instances process messages and what should you do if lag grows?

**Q6 (MCQ, medium).** Which operation is naturally idempotent?
A. `balance = balance + 10`
B. `INSERT` a new row with an auto-increment id
C. `UPDATE orders SET status = 'shipped' WHERE id = 42`
D. Send an email

**Q7 (application, hard).** Design an idempotent consumer for a `PaymentCaptured` event (id, order_id, amount) that updates a ledger. Provide the transaction steps and explain where the ack goes.

**Q8 (MCQ, hard).** Your team doubles partitions from 12 to 24 on a keyed topic to increase throughput. What is the main risk?
A. Consumers will stop working.
B. Key-to-partition mapping changes, so events for the same key may land on different partitions during/after the change, breaking per-key ordering for in-flight data.
C. Kafka deletes old messages.
D. Nothing; ordering is global.

**Q9 (short, hard).** A message that always fails blocks a partition. Give two designs to handle it and one trade-off for each.

**Q10 (short, hard).** When would you choose event sourcing, and name two costs you'd accept?
