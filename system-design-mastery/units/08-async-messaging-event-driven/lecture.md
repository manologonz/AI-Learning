# Unit 8: Async Messaging & Event-Driven Architecture

## Objectives
By the end of this unit you can:
- Choose between a message queue and a log-based broker (SQS/RabbitMQ vs Kafka) for a workload.
- Explain delivery semantics (at-most-once, at-least-once, "exactly-once") and build idempotent consumers.
- Implement the transactional outbox pattern to avoid dual-write bugs.
- Design retries, dead-letter queues, ordering (partition keys), and backpressure handling.
- Explain event sourcing and CQRS and when they are (and aren't) worth it.

## Prerequisites
Units 6-7 (partitioning, consistency, failure ambiguity).

## Core content

### 1. Why go asynchronous
Synchronous calls couple caller and callee in **time, availability, and load**: the callee must be up, fast, and able to absorb your peak. Async messaging gives:
- **Decoupling and availability**: the producer succeeds even if consumers are down; messages wait.
- **Load leveling (buffering)**: absorb spikes and let consumers drain at their own pace.
- **Fan-out**: one event, many independent consumers (email, analytics, search index).
- **Retries and durability** built in.
Costs: eventual consistency, harder debugging, duplicates and reordering, operational complexity, and no immediate result for the caller. **Use sync when the caller needs the answer now; async when it needs the work done eventually.**

### 2. Queues vs logs

| | **Message queue** (SQS, RabbitMQ, ActiveMQ) | **Log/stream** (Kafka, Kinesis, Pulsar, Redpanda) |
|---|---|---|
| Model | Broker holds messages until a consumer acks; message deleted after ack | Append-only, partitioned, retained log; consumers track their own **offset** |
| Consumption | Competing consumers share a queue; each message goes to one worker | Multiple independent **consumer groups** each read the whole stream; within a group, one consumer per partition |
| Replay | No (once acked, gone) | Yes: rewind the offset, reprocess history |
| Ordering | Weak / best-effort (FIFO queues limited in throughput) | **Per partition** strict order |
| Throughput | Moderate-high | Very high (sequential I/O, batching) |
| Best for | Task/job distribution, work queues, per-message retry/delay | Event streams, analytics, event sourcing, multi-consumer pipelines, CDC |

Rule of thumb: *"Do this task once"* → queue. *"This happened; anyone who cares can read it, now or later"* → log.

**Kafka essentials**: topic → partitions (unit of parallelism and ordering); producers choose partition by key (`hash(key) % partitions`, so all events for one entity are ordered); replication factor 3 with `acks=all` + `min.insync.replicas=2` for durability; consumer group rebalancing; retention by time/size or log compaction (keep the latest value per key). **Max parallelism per group = number of partitions.**

### 3. Delivery semantics
- **At-most-once**: send and don't retry → messages may be lost, never duplicated. (Metrics you can afford to lose.)
- **At-least-once**: retry until acknowledged → **no loss, but duplicates**. **The default in practice.**
- **Exactly-once**: the network can't provide it end to end (Two Generals problem). What systems offer is **effectively-once**: at-least-once delivery + **idempotent processing**. Kafka's "exactly-once semantics" works only within Kafka (idempotent producer + transactions across topics); side effects outside (an email, a payment API) still need idempotency.

**Idempotent consumer patterns**
- **Idempotency key / dedupe table**: record processed message IDs in the same transaction as the side effect: `INSERT INTO processed(msg_id) ... ON CONFLICT DO NOTHING`, and skip if it already existed.
- **Naturally idempotent operations**: `SET status='shipped'` (absolute), not `balance += 10` (relative). Upserts by key.
- **Conditional updates with versions**: apply only if `event.version > row.version`.

```ts
// At-least-once consumer with idempotent processing
async function handle(msg: { id: string; orderId: string; amount: number }) {
  await db.tx(async (tx) => {
    const fresh = await tx.query(
      "INSERT INTO processed_messages(id) VALUES ($1) ON CONFLICT DO NOTHING RETURNING id", [msg.id]);
    if (fresh.rowCount === 0) return;               // duplicate → skip, but still ack
    await tx.query("UPDATE invoices SET paid = paid + $1 WHERE order_id = $2", [msg.amount, msg.orderId]);
  });                                                // dedupe + effect commit atomically
  await broker.ack(msg);                             // ack AFTER commit; crash before ack → redelivery → deduped
}
```
Ack **after** processing (at-least-once). If you ack first and crash, the message is lost (at-most-once).

### 4. The dual-write problem and the transactional outbox
You must update your DB **and** publish an event. Two separate writes can't be atomic:
- DB commit succeeds, publish fails → downstream never learns (lost event).
- Publish succeeds, DB commit fails → downstream sees a phantom event.

**Transactional outbox**: in the *same DB transaction* as the business change, insert the event into an `outbox` table. A separate **relay** (poller or CDC tool like Debezium reading the WAL) publishes outbox rows to the broker and marks them sent. Delivery is at-least-once → consumers must be idempotent. This converts a distributed atomicity problem into a local transaction plus an idempotent republish.

```sql
BEGIN;
  INSERT INTO orders(id, user_id, total, status) VALUES ($1,$2,$3,'placed');
  INSERT INTO outbox(id, topic, key, payload) VALUES (gen_random_uuid(),'orders.placed',$1,$4);
COMMIT;   -- relay reads outbox → publishes → deletes/marks rows
```

### 5. Reliability mechanics
- **Retries**: transient failures → retry with exponential backoff + jitter (Unit 10). Use **delay queues / retry topics** so a poison message doesn't block others.
- **Dead-letter queue (DLQ)**: after N failures, move the message aside for inspection/replay; **alert on DLQ depth**. Unmonitored DLQs are data graveyards.
- **Poison messages**: malformed messages that always fail; without a DLQ/skip policy they block a partition forever (head-of-line blocking).
- **Ordering**: only guaranteed per partition/queue-group; choose the partition key = the entity whose order matters (`order_id`). Retries can reorder; design consumers to tolerate out-of-order (versions/timestamps) or use per-key sequential processing.
- **Backpressure**: consumer lag (Kafka: `max offset − committed offset`) is your key health metric; scale consumers up to the partition count, shed/skip low-priority work, cap producer rates, or size partitions for peak. Unbounded queues just delay failure; set retention/TTLs and alert.
- **Schema evolution**: events are contracts. Use a schema registry (Avro/Protobuf), additive changes, versioned event types; consumers ignore unknown fields.
- **Visibility timeout (SQS)**: a received message hides for T seconds; if not deleted in time it reappears (duplicate processing risk for long tasks: extend or set T > processing time).

### 6. Event sourcing and CQRS
- **Event sourcing**: store the *sequence of events* (`OrderPlaced`, `ItemAdded`, `OrderShipped`) as the source of truth; current state = fold(events). Pros: full audit trail, time travel, rebuild any projection, natural integration events. Cons: complexity, event schema evolution forever, replays are slow (use snapshots), queries need projections, and GDPR deletion is hard.
- **CQRS** (Command Query Responsibility Segregation): separate the **write model** (commands, validation, invariants) from **read models** (denormalized views optimized for queries, fed by events). Pros: independent scaling and shapes for reads; cons: eventual consistency between write and read sides, more moving parts.
- Use them when audit/history or radically different read/write shapes are core needs (finance, collaboration). **Don't** adopt for CRUD apps; a table with a `history` table or an outbox often gives 80% of the benefit.

### 7. Choreography vs orchestration (preview of Unit 9)
- **Choreography**: services react to each other's events; decentralized, loosely coupled, but the overall flow is implicit and hard to follow.
- **Orchestration**: a coordinator (workflow engine such as Temporal/Step Functions) explicitly commands steps; clearer, easier to monitor/compensate, but a central component.

## Worked examples

### Straightforward: queue vs log
*Send welcome emails after signup; each email must be sent once (best effort), retries on provider errors, no need to replay.* → **Queue** (SQS) with worker pool, retries, DLQ. *Every order event must feed search indexing, analytics, and fraud detection independently, and new consumers should be able to backfill from history.* → **Log** (Kafka) with a consumer group per system and retention of days/weeks.

### Typical: outbox + idempotent consumers
*Order service creates an order and must notify inventory and billing.* Order service: one DB transaction inserts the order row and an `OrderPlaced` outbox row. Debezium tails the WAL and publishes to Kafka topic `orders.placed`, keyed by `order_id`. Inventory and billing consumer groups each dedupe on `event_id` inside their own DB transaction. If Debezium restarts and republishes, duplicates are dropped; if the consumer crashes before ack, redelivery is deduped. No lost events, no phantom events, effectively-once processing.

### Edge case: poison message and reordering
*A Kafka consumer hits a message that throws every time (bad JSON). With `retry forever` semantics it blocks the partition; lag grows to 10M; downstream is stale for hours. Meanwhile someone "fixes" it by sending failures to a retry topic, and now `OrderCancelled` (retried later) is processed after `OrderShipped` for the same order.* Fixes: bounded retries → DLQ with alerting; for ordering-sensitive entities either process failures **in place** with limited retries before parking the *entire key* (per-key blocking), or make consumers version-aware (`if event.version <= state.version: ignore/park`). Order and retry are in tension: you must choose which invariant wins per use case.

## Common misconceptions

- **"Kafka gives exactly-once delivery."** It gives exactly-once *within Kafka* under conditions; effects outside need idempotency.
- **"Ack early to avoid duplicates."** That just trades duplicates for loss.
- **"Adding partitions preserves ordering."** Changing the partition count remaps keys (per-key order breaks across the change); partitions can't easily be reduced.
- **"A queue means ordering."** Standard queues don't guarantee order; ordering needs FIFO queues or partition keys, at throughput cost.
- **"Async is always better."** It shifts complexity to consistency, debugging, and operations; keep synchronous paths for interactions needing immediate results.
- **"Publish after commit is safe enough."** A crash between commit and publish loses the event; use the outbox.
- **"DLQ solves failures."** Only if you monitor, triage, and replay it.
- **"Event sourcing = using Kafka."** Event sourcing is a storage pattern for an aggregate's state; Kafka is a transport.

## Check your understanding

1. Consumer processed a message, then crashed before acking. What happens next and what must the consumer be? *(Redelivery; idempotent.)*
2. 12 partitions, consumer group of 20 consumers: how many are active? *(12; the extra 8 idle.)*
3. Why does the outbox pattern need consumers to be idempotent? *(The relay may publish a row more than once: at-least-once.)*
4. Which key preserves per-customer ordering in Kafka? *(Partition by `customer_id`.)*

## Summary
- Async decouples time, availability, and load, at the cost of eventual consistency and duplicates.
- Queues for tasks; logs for streams/replay/multi-consumer; ordering is per partition.
- Effectively-once = at-least-once + idempotent consumers; ack after processing.
- Outbox solves the DB+broker dual write; use DLQs with alerts; monitor consumer lag.
- Event sourcing and CQRS are powerful but costly; adopt them for audit/read-shape needs, not by default.
