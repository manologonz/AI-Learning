# Unit 9 Exam: Answer Key

**Q1.**
- Root cause: the orchestrator treated a **timeout as a definitive failure**, but the outcome was *unknown* (the hotel confirmed later). It compensated the flight while the hotel booking stood; then the payment logic charged despite an inconsistent bundle.
- Fixes: (1) treat timeouts as unknown → query the supplier's booking status (by our reference id) or retry with an idempotency key before deciding; (2) persist saga state and make compensations for *all* steps that may have succeeded (compensate hotel if it later confirms: "late success handling"); (3) **authorize, don't capture** until all parts are confirmed; capture at the pivot; (4) a reconciliation job comparing our bookings to supplier confirmations; alert on mismatches.

**Q2. Sample saga.**
1. `createBooking` (PENDING); compensation: mark CANCELLED.
2. `holdFlight` (compensable, hold with expiry; compensation `releaseFlightHold`).
3. `holdHotel` (compensation `releaseHotelHold`).
4. `holdCar` (compensation `releaseCarHold`).
5. `authorizePayment` (compensation: `voidAuthorization`, cheap and clean).
6. **Pivot: `confirmAllHolds` / `capturePayment`**: the point of no return: once holds are converted to confirmed bookings and the payment captured, we move forward only.
7. Retriable after pivot: `issueTickets/confirmations`, `sendEmail`, `updateBooking to CONFIRMED` (retry until success; if impossible, human/refund policy).
- Order justification: holds first (cheap to release, auto-expire, riskiest to obtain: least available/most likely to fail first); authorization before capture (voidable); the expensive/hard-to-undo capture last as the pivot; notifications after. Do the step most likely to fail early to minimize compensations; complete within the 10-min hold window (set saga timeout well below hold expiry).

**Q3.**
- Since the API isn't idempotent: (1) **check before create**: first query the supplier for an existing booking with our client reference (if the API supports a reference/lookup); (2) use a **unique client reference** in the booking request; before any retry, look it up; (3) serialize per-booking attempts (single worker per booking id via partition key) so concurrent duplicates can't occur; (4) wrap the supplier in an **adapter service** that stores request state (`pending/created`) and dedupes on our idempotency key.
- Residual risk: a crash after the supplier creates the booking but before our record is written, when the supplier can't be queried by reference → a duplicate or orphan booking; mitigate with reconciliation against supplier reports and cancelling orphans.

**Q4. Sample answer.**
- Anomaly 1: customer sees "Flight confirmed" while the hotel later fails and the flight is released (flickering/dirty state). Countermeasure: show a single `PENDING` state until the saga completes; only show confirmed bookings.
- Anomaly 2: a support agent cancels mid-saga, racing with the orchestrator continuing to hold/charge. Countermeasure: route the cancel through the saga (a cancel signal/event the orchestrator handles by compensating), and use state-machine guards so the cancel and the saga can't both proceed (reject cancel when in PIVOT-committed state → offer refund flow instead).
- Also acceptable: double-submit by the user, or two sagas for the same inventory (use holds).

**Q5.**
- **Orchestrator** (e.g., Temporal): 6-7 steps with compensations, timeouts, human-in-the-loop, and ambiguous outcomes need explicit state machines, retries with backoff, durable timers, and visibility.
- Survives crashes: the workflow engine persists an event history and replays deterministic workflow code to resume from the last completed activity; activities are idempotent (with idempotency keys).
- Failure introduced: the orchestrator becomes a critical dependency (availability, capacity, deployment/versioning of long-running workflows), and non-deterministic workflow code or incompatible version changes can break in-flight workflows.

**Q6: B.**
- Spanner runs 2PC across Paxos-replicated groups, so the coordinator and participants are replicated state machines (blocking on coordinator failure is much less likely). A is false; C/D false (Spanner is externally consistent).

**Q7: B.**
- Do cheap-to-undo steps first (B, C), and the hard-to-undo step (A) last as the pivot, so failures early only need cheap compensations. A/C put the hard-to-undo step early → refunds when later steps fail; D ignores costs.

**Q8. Sample answer.**
- States: `PENDING → HELD → PAID → CONFIRMED`; failure/exit states: `CANCELLED`, `REFUNDED`; transitions allowed: `PENDING→HELD`, `HELD→PAID`, `PAID→CONFIRMED`, `PENDING|HELD→CANCELLED`, `PAID→REFUNDED`, `CONFIRMED→REFUNDED`.
- Each event carries a target transition; the consumer applies it only if the current state permits that transition (and/or the event version is greater); otherwise it ignores it (stale/duplicate) or parks it for later (out-of-order). E.g., `PaymentCaptured` arriving while still `PENDING` (before `HELD`) is parked/retried; a duplicate `PaymentCaptured` in `PAID` is a no-op.

**Q9. Sample answer.**
- Purchase: debit `Customer receivable/Cash` $50, credit `Revenue` $50. Refund: debit `Revenue` (or `Refunds`) $50, credit `Customer/Cash` $50. Each entry is immutable and pairs of debits/credits net to zero.
- Aids reconciliation: full history is preserved (no destructive updates), the balance is derived by summing entries, each entry references the transaction/event id, so discrepancies between systems can be traced entry by entry; corrections are new entries.

**Q10. Sample answer.**
- (a) **Single-DB transaction**: same service/DB; ACID suffices.
- (b) **Relax the invariant / saga-lite**: debit-then-credit with an idempotent transfer record and reconciliation; brief inconsistency accepted (or a per-user ledger).
- (c) **Saga** with holds/TCC: two independent external systems, no shared transaction.
- (d) **CRDT** (OR-Set/sequence CRDT): concurrent edits merge without coordination; no global invariant required.
