# Unit 9 Quiz: Answer Key

**Q1: B.**
- In-doubt participants hold locks until the coordinator recovers; every participant must be available. A is wrong (it does guarantee atomic commit); C is wrong; D nonsense.

**Q2: B.**
- Compensation is a forward business action (refund, release hold) that semantically reverses a committed step. A: local transactions already committed and can't be rolled back. C is a retry; D is unrelated.

**Q3: C (Isolation).**
- Intermediate states are visible to others. A is achieved via compensation; B by local commits; D is application-defined invariants (eventually).

**Q4. Sample answer.**
1. Co-locate data that must change atomically in one database/shard (shard key aligned with the transaction).
2. Redraw service boundaries so the invariant is owned by one service.
3. (Also) Relax the invariant with business tolerance (allow overbooking + compensation).

**Q5.**
- Order: (1) `reserveFlightSeat (hold)`, (2) `reserveHotelRoom (hold)`: compensable, cancelable with expiry; (3) `chargeCard`: **pivot**; (4) `sendConfirmationEmail`: retriable after the pivot (non-compensable, do last).
- If charge fails: release both holds. After the charge, confirm bookings and email via forward recovery (retry until success).
- (Reserving the riskiest/hardest-to-obtain hold first, or which hold order minimizes failures, is a valid refinement.)

**Q6: B.**
- Timeout = unknown outcome. Compensating prematurely can leave a charge on a cancelled order. Query status or retry with the same idempotency key. A/C/D are wrong.

**Q7. Sample answer.**
- Semantic countermeasures reduce anomalies from visible intermediate state: (1) **status flags/pending state**: an order is `PENDING` until the saga completes, so other processes know not to ship/act on it; (2) **reservations with expiry** (holds) rather than immediate decrement; (3) **commutative updates** (increments/decrements) so ordering doesn't matter; (4) **re-read before act**, verifying the current state before a dependent action.

**Q8: C.**
- A grow-only/increment-only counter (G-Counter) merges by taking per-replica max/sum, so no coordination. A (uniqueness), B (non-negative balance), and D (single leader) are global invariants requiring coordination/consensus.

**Q9. Sample answer.**
- Compensation: `refund(paymentId, amount)` with an idempotency key (`orderId:refund`), retried until success, recorded as a new ledger entry.
- Problems: (1) not a true undo: the customer sees both a charge and a refund; refunds may take days to settle; (2) fees may not be returned, the refund can fail (expired card/closed account), and the intermediate state was visible (customer saw a charge). Prefer **authorize then capture** (auth is voidable) so the compensation is a "void", not a refund.

**Q10. Sample answer.**
- Because distributed systems drift: lost messages, bugs, manual interventions, partial failures beyond retry limits, and unknown outcomes leave systems inconsistent. Reconciliation compares systems of record (orders vs payments vs shipments) and repairs or flags mismatches (charged but no order; shipped but not paid). It's the safety net that catches what the saga logic missed.
