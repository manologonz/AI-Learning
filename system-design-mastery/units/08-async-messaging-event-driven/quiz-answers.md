# Unit 8 Quiz: Answer Key

**Q1: B.**
- At-least-once retries until ack: no loss, possible duplicates. A describes at-most-once. C is not achievable end-to-end by the network. D is unrelated.

**Q2: B.**
- Kafka-style log: independent consumer groups, retention, replay. A: competing consumers split messages (each system would see only some). C couples producers and can't replay. D isn't event driven.

**Q3: B.**
- Acking before processing means the broker considers it done; a crash loses it. A describes ack-after processing. C/D are wrong.

**Q4. Sample answer.**
- Example: an order service commits the order to the DB and then publishes `OrderPlaced` to Kafka; if it crashes between the two, the DB has the order but no event was sent (or, reversed, an event exists for an order that failed to commit).
- Outbox: write the event to an `outbox` table in the same DB transaction as the order, then a relay publishes outbox rows (at-least-once) and marks them sent; atomicity comes from the local transaction, and consumers dedupe.

**Q5: 8 instances process; 4 idle.**
- Max parallelism in a consumer group = partition count. If lag grows: optimize consumer processing, batch, or **increase partitions** (with ordering/rebalance implications) *before* adding instances beyond 8. Also check downstream bottlenecks.

**Q6: C.**
- Setting an absolute status is idempotent. A (relative increment) double-applies; B creates duplicates; D repeats a side effect.

**Q7. Sample answer.**
```
BEGIN;
  INSERT INTO processed_events(event_id) VALUES ($id) ON CONFLICT DO NOTHING;
  -- if 0 rows inserted: COMMIT (nothing to do) and ack
  INSERT INTO ledger(order_id, amount, event_id) VALUES ($order, $amount, $id);  -- or upsert by event_id
COMMIT;
ACK message  -- after commit
```
- Dedupe row and ledger entry commit atomically; if the process crashes before ack, redelivery finds the event_id already processed and skips. Ack goes after the commit.

**Q8: B.**
- `hash(key) % partitions` changes with the partition count, so a key can move; events for one key can be split across old/new partitions with no cross-partition ordering. A/C/D are wrong. (Mitigate: plan partition counts up front, or drain/cutover.)

**Q9. Sample answer.**
1. **Bounded retries → DLQ** (park the message and continue). Trade-off: subsequent messages for that key may be processed out of order; need alerting and replay tooling.
2. **Per-key blocking with skip/park of the whole key**, or a retry topic with backoff. Trade-off: delays that entity's processing (and increases complexity), but preserves per-key ordering.
- Also: schema validation at the producer to prevent poison messages.

**Q10. Sample answer.**
- Choose when the history is itself the product/requirement: strict audit trail (finance), need to rebuild many read models, time-travel debugging, or complex domain workflows.
- Costs: event schema versioning forever, projection/eventual consistency complexity, replay/snapshot management, harder ad-hoc queries, GDPR/deletion challenges, steeper team learning curve.
