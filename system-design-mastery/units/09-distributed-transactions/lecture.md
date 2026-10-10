# Unit 9: Distributed Transactions & Data Consistency Patterns

## Objectives
By the end of this unit you can:
- Explain two-phase commit (2PC), its guarantees, and why it's avoided across service boundaries.
- Design a saga (orchestrated or choreographed) with compensating actions for a multi-service workflow.
- Handle partial failure, duplicates, and out-of-order events with idempotency, timeouts, and reconciliation.
- Recognize when CRDTs or other convergent designs beat coordination.
- Choose between "avoid the distributed transaction" (co-locate data), 2PC, and sagas for a given invariant.

## Prerequisites
Units 7-8 (consistency, consensus, messaging, outbox).

## Core content

### 1. The problem
A business operation spans multiple systems: *place order → reserve inventory → charge card → schedule delivery*. Each has its own datastore. You want **atomicity**: all succeed or none. But there's no shared transaction, and any call can time out with an ambiguous outcome (Unit 7).

**First question: can you avoid it?** The best distributed transaction is the one you don't need:
- **Co-locate** data that must change atomically in one database/shard (shard by tenant or order so a single ACID transaction covers it).
- **Redraw service boundaries** so the invariant lives inside one service (a common sign that two "microservices" should be one).
- **Weaken the invariant** with business tolerance (overbooking with apology and compensation, as airlines do).
Only when the invariant truly spans independent stores do you need one of the patterns below.

### 2. Two-phase commit (2PC)
A **coordinator** drives participants:
1. **Prepare (voting)**: each participant durably records that it *can* commit and locks the resources, then votes yes/no.
2. **Commit/abort**: if all voted yes, the coordinator logs the decision and tells everyone to commit; else abort.

Guarantees **atomic commit** across participants. Problems:
- **Blocking**: if the coordinator crashes after participants voted yes but before sending the decision, participants hold locks *in doubt*, unable to commit or abort unilaterally, until the coordinator recovers. (3PC and consensus-based commit, like Spanner's Paxos-replicated 2PC, mitigate this.)
- **Latency and lock holding** across network round trips, hurting throughput; contention amplifies.
- **Availability**: all participants must be up (availability multiplies down).
- **Coupling**: participants must support XA-like protocols; most SaaS APIs and NoSQL stores don't.
Where it's used: within a single DB engine across shards (Spanner, CockroachDB, Vitess with care), XA between a DB and a broker in old enterprise stacks. **Across independently owned services, avoid it.**

### 3. Sagas
A **saga** is a sequence of *local* transactions, each with a **compensating transaction** that semantically undoes it. If step k fails, run compensations for steps k−1…1 (in reverse).

Example: `CreateOrder → ReserveInventory → ChargePayment → ScheduleDelivery`. If ChargePayment fails: `ReleaseInventory` then `CancelOrder`.

Key properties:
- **No isolation**: intermediate states are visible (inventory reserved, payment not yet done); saga = ACD without I. Other transactions might see or act on partial results ("dirty reads"), so you need **semantic countermeasures**: status flags (`PENDING`), reservations with expiry, commutative updates, re-reading before acting, and pessimistic ordering (do the hardest-to-undo step last).
- **Compensations must be idempotent, retryable, and always succeed eventually** (they can't "fail"; keep retrying). Some actions are **non-compensable** (email sent, physical goods shipped): put them at the end, after the **pivot transaction** (the point of no return), or use *forward recovery* (retry until success) after the pivot.
- **Ordering of steps**: (1) compensable steps, (2) pivot (the step that decides success), (3) retriable steps that must eventually succeed.

**Orchestration vs choreography**
- **Orchestrated** (a saga coordinator/workflow engine, e.g., Temporal, AWS Step Functions, Camunda): explicit state machine, timeouts, retries, and compensation; easy to observe. Recommended for anything beyond 2-3 steps.
- **Choreographed**: each service reacts to events and emits new ones; loose coupling but implicit flow, cyclic dependencies risk, hard to debug/evolve.

```ts
// Orchestrated saga sketch with persisted state and compensations
type Step = { name: string; run: (ctx: Ctx) => Promise<void>; compensate?: (ctx: Ctx) => Promise<void> };

const steps: Step[] = [
  { name: "reserveInventory", run: c => inventory.reserve(c.orderId, c.items, { idem: c.orderId + ":inv" }),
                              compensate: c => inventory.release(c.orderId, { idem: c.orderId + ":inv-rel" }) },
  { name: "chargePayment",    run: c => payments.charge(c.orderId, c.amount, { idem: c.orderId + ":pay" }),   // pivot
                              compensate: c => payments.refund(c.orderId, { idem: c.orderId + ":pay-rf" }) },
  { name: "scheduleDelivery", run: c => delivery.schedule(c.orderId, { idem: c.orderId + ":del" }) }, // retriable after pivot
];

async function runSaga(ctx: Ctx) {
  const done: Step[] = [];
  for (const s of steps) {
    try { await withRetry(() => s.run(ctx)); await saveProgress(ctx, s.name); done.push(s); }
    catch (e) {
      for (const d of done.reverse()) if (d.compensate) await retryForever(() => d.compensate!(ctx));
      throw e; // mark saga failed
    }
  }
}
// State (which steps are done) is persisted so a crashed orchestrator resumes; every call carries an idempotency key.
```

### 4. Supporting patterns
- **Idempotency keys everywhere** (Unit 3/8): every step and compensation can be safely retried.
- **Timeouts and unknown outcomes**: on timeout, don't assume failure; **query the status** or retry idempotently; for money, use a status API or reconciliation.
- **Reservations / TCC (Try-Confirm-Cancel)**: *Try* places a hold (reserve inventory, authorize card) with an expiry; *Confirm* finalizes; *Cancel* releases. Reduces the visibility of unfinished state and avoids hard-to-undo actions.
- **Reconciliation**: periodic jobs compare systems of record (orders vs payments vs shipments) and repair or flag discrepancies. A safety net you need regardless: distributed systems drift.
- **Ledger/double-entry**: model money as append-only entries with balancing debits and credits; corrections are new entries (compensations) not updates; makes audits and reconciliation tractable.
- **Outbox** (Unit 8) for reliably emitting saga events/commands.
- **Deduplication of out-of-order events**: version numbers/state machines that reject invalid transitions (`shipped` can't precede `paid`).

### 5. Convergent designs (avoiding coordination)
- **CRDTs (Conflict-free Replicated Data Types)**: data structures whose merge operation is commutative, associative, and idempotent, so replicas converge without coordination. Examples: G-Counter/PN-Counter (likes, inventory deltas), OR-Set (shopping cart adds/removes), LWW-Register (with the clock caveat), sequence CRDTs (collaborative text, e.g., Yjs/Automerge). Trade-off: only certain invariants fit; **can't enforce global constraints** like "balance ≥ 0" or "unique username".
- **Invariant confluence**: if an invariant can't be violated by concurrent operations (e.g., counters that only increase, unique IDs from a generator), you don't need coordination; if it can (uniqueness, non-negative balance), you do. **Coordinate only where invariants demand it.**

## Worked examples

### Straightforward: choosing the pattern
*Transfer money between two accounts in the same Postgres database.* One ACID transaction. No saga, no 2PC. *Transfer between two accounts in different banks.* No shared transaction possible → saga/ledger with pending state, or a clearing process with reconciliation (this is how real banking works: debit → in-flight → credit or reverse).

### Typical: order saga with failure at the pivot
*Steps: create order (PENDING) → reserve inventory (hold 15 min) → charge card → confirm order → schedule delivery.* Payment declines: orchestrator runs `releaseInventory` (idempotent, retried until success) and `cancelOrder`; customer sees "payment failed". Since inventory is a *hold with expiry*, even if the orchestrator dies, the hold expires and the inventory frees itself. The charge is the pivot: once it succeeds, remaining steps (confirm, schedule delivery) are retried forward until success; if scheduling can't ever succeed (no couriers), the compensation is a refund, decided by a human/policy rather than automatic rollback.

### Edge case: the timeout that lies
*The orchestrator calls `payments.charge`; the call times out after 5 s. It decides "failed" and compensates (releases inventory, cancels the order). In fact the payment service did process the charge 6 s later.* Customer is charged for a cancelled order. Fixes: (1) treat a timeout as **unknown**, not failure; (2) retry the charge with the *same idempotency key* to learn the true outcome, or query the payment status; (3) only compensate on a definitive failure; (4) a reconciliation job that detects "charged but order cancelled" and refunds. The saga's correctness depends on making outcomes *knowable*.

## Common misconceptions

- **"Microservices need distributed transactions everywhere."** Good boundaries make most operations local; excessive cross-service transactions signal wrong boundaries.
- **"Sagas give you ACID."** They give atomicity via compensation, not isolation; intermediate states are visible.
- **"Compensation is rollback."** It's a *new business action* (refund, release) that may leave traces (a charged-then-refunded transaction appears on the statement).
- **"2PC is obsolete."** It's used inside distributed databases; it's the *cross-service* use that's problematic.
- **"A timeout means failure."** It means unknown.
- **"CRDTs solve consistency."** They solve convergence for specific data types, not global invariants.
- **"Choreography is simpler."** Simple for 2 services; hard to reason about for long flows.
- **"Compensations can fail."** They must be designed to eventually succeed (retry, alert a human); an unrecoverable compensation is an incident.

## Check your understanding

1. What blocks in 2PC if the coordinator crashes after all participants vote yes? *(Participants remain in-doubt holding locks until the coordinator recovers.)*
2. In a saga, which step should go last: sending a confirmation email or reserving inventory? *(Email; do non-compensable steps after the pivot; reservations are compensable.)*
3. Can a CRDT enforce "no negative inventory"? *(Not by itself; that's a global invariant requiring coordination or reservation.)*
4. Why must saga compensations be idempotent? *(They're retried on failures/crashes; running twice must not double-refund.)*

## Summary
- Avoid distributed transactions by co-locating data, choosing boundaries wisely, or relaxing invariants.
- 2PC gives atomic commit but blocks and couples; use it inside a database, not across services.
- Sagas = local transactions + compensations; no isolation; orchestrate with persisted state and idempotency.
- Treat timeouts as unknown outcomes; add status checks, reservations (TCC), reconciliation, and ledgers.
- CRDTs/invariant confluence let you skip coordination when invariants allow it.
