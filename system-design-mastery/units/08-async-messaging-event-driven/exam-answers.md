# Unit 8 Exam: Answer Key

**Q1.**
- Coupling: temporal (all steps in one request), availability (SMS provider down → orders fail), latency (slowest dependency dictates API latency), atomicity (DB write + external calls aren't atomic → order exists but restaurant not notified), retry semantics (retrying the whole request duplicates SMS/charges).
- Synchronous: validating the order and persisting it (customer needs an answer), and typically **payment authorization** (the customer needs to know it succeeded; some flows use async with a pending state). Async: restaurant notification (needs reliability, not immediate response), SMS/notifications, search indexing, analytics.
- Principle: keep on the critical path only what the user must know now; everything else via events.

**Q2.**
- API transaction: `BEGIN; INSERT orders(... status='placed'); INSERT outbox(event_id, topic='orders.placed', key=order_id, payload); COMMIT;` → return 201 (pending confirmation) or after a sync payment authorization.
- Relay: Debezium (CDC on outbox) or a poller publishes to Kafka `orders.placed`, key = `order_id` (per-order ordering for subsequent events on that order; all events of one order in one partition).
- Consumers (separate groups): restaurant-notifier, SMS-notifier, search-indexer, analytics, fraud. Each idempotent (dedupe by `event_id`).
- Partition key: `order_id` (high cardinality → even spread; ordering per order). Not `restaurant_id` (a big restaurant → hot partition).

**Q3.**
- Causes: (1) at-least-once redelivery (consumer crashed after sending SMS but before acking; or rebalance), (2) outbox relay republished the event; also (3) producer retries.
- Consumer logic: `INSERT INTO sms_sent(event_id, order_id) ON CONFLICT DO NOTHING`, claim first (status `sending`), then call the provider with an **idempotency key = event_id** (if the provider supports it), then mark `sent`; on crash between "claim" and "sent", a redelivery finds `sending` and either re-sends with the same provider idempotency key or checks the provider's status API.
- Limits: if the provider has no idempotency support and the crash happens after the SMS was sent but before recording, a duplicate is possible. Choose the failure you prefer: (a) mark-then-send → possible **lost** SMS on crash (at-most-once); (b) send-then-mark → possible **duplicate** (at-least-once). For notifications, a rare duplicate is usually preferable to a lost message.

**Q4.**
- Failure mode: if the consumer retries in place forever (or blocks on the 500), that partition's other restaurants' events are stuck behind the failing message (head-of-line blocking): lag grows for restaurants sharing the partition.
- Design: on failure, publish the message to a **retry topic with delay tiers** (e.g., 1 min, 5 min, 30 min, exponential backoff + jitter) and commit the offset so the main partition proceeds; after N attempts → **DLQ** with alert; add a **circuit breaker** per restaurant endpoint (skip attempts while open, route to the retry tier); provide tooling to replay the DLQ once fixed. Ordering trade-off: later messages for that restaurant may overtake; use event versions/timestamps or a status check to ignore outdated notifications.

**Q5.**
- Measure in order: (1) consumer lag per partition (skew = hot partition/key vs uniform = capacity), (2) consumer CPU/processing time per message, downstream latency (DB/HTTP calls), (3) rebalancing frequency, (4) producer rate vs consumption rate, (5) batch sizes/fetch configs, (6) errors/retries inflating work.
- Options: optimize the handler (batch DB writes, async I/O, remove synchronous downstream calls), scale instances **up to 24** (only 24 partitions; the extra 6 of 30 instances sit idle), increase partitions (ordering remapping risk, plan a migration), fix hot keys, shed/skip low-priority events, add capacity to the bottleneck downstream. Note: at 30 instances you're already beyond the useful count; further scaling requires partitions or per-consumer concurrency (multi-threaded processing preserving per-key order).

**Q6. Sample recommendation.**
- **Outbox + status column + append-only `order_events` history table** (a lightweight audit log), rather than full event sourcing. It gives an audit trail and analytics feed (the log table/CDC stream) with ordinary CRUD reads, without the schema-evolution and projection burden of full event sourcing.
- Trade-offs: you lose "state is a pure fold of events" and free time travel; acceptable if audit + analytics rebuild are met by the history table/CDC into the warehouse.
- Reversal: adopt full event sourcing if we need to reconstruct arbitrary past states for disputes/legal reasons, or many read models with complex domain rules that must be rebuilt regularly; at that point the cost is justified.

**Q7: B.**
- EOS covers Kafka read-process-write via transactions. A is wrong: external effects need idempotency. C is wrong: offsets are part of the mechanism. D is wrong: ordering remains per partition.

**Q8: B.**
- The message reappears after the visibility timeout while the first worker is still working: concurrent duplicate processing. Fix: set the timeout above max processing time or extend via `ChangeMessageVisibility` heartbeats, and keep processing idempotent. A is false; C is false; D happens only after max receive count.

**Q9. Sample answer.**
- Choreography: services publish/react to events (payment → `PaymentCaptured` → restaurant → ...); loosely coupled, no central component, but the flow is implicit, hard to visualize, and compensation logic is scattered.
- Orchestration: a coordinator (workflow engine like Temporal/Step Functions) calls steps and handles retries/timeouts/compensations explicitly; visible state, easier monitoring, and central point to change; coupling to orchestrator and its availability.
- Pick **orchestration** for 5 steps with compensations: flow clarity, timeouts, and rollback logic are the hard parts; use choreography for simple fan-out notifications.

**Q10. Sample answer.**
- Schema: `event_id (uuid), event_type ("OrderPlaced"), schema_version (int), occurred_at (ISO-8601 UTC), order_id, customer_id, items[{sku, qty, price_cents}], total_cents, currency, restaurant_id, metadata{trace_id}`. Use Avro/Protobuf with a schema registry (compatibility mode BACKWARD/FULL).
- Backward-compatible: add an optional field (`promo_code`), add a new enum value that consumers treat as unknown/default (if tolerated), add a new event type.
- Breaking: rename or remove a field, change a field's type (e.g., `total` dollars → cents), change semantics (`total` now excludes tax).
- Safe rollout of a breaking change: publish the new version as a **new event type/topic** (`OrderPlaced.v2`) or dual-publish both versions; migrate consumers one by one; monitor v1 consumption; retire v1 after all consumers move (expand → migrate → contract).
