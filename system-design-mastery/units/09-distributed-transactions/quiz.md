# Unit 9 Quiz: Distributed Transactions & Data Consistency Patterns

Answers in `quiz-answers.md`.

**Q1 (MCQ, easy).** What is the main weakness of two-phase commit (2PC)?
A. It can't guarantee atomicity.
B. It can block participants holding locks if the coordinator fails after they've voted to commit, and it couples availability of all participants.
C. It only works with NoSQL.
D. It requires no network communication.

**Q2 (MCQ, easy).** In a saga, what is a compensating transaction?
A. A database rollback command.
B. A new action that semantically undoes the effect of a previously committed local step.
C. A retry of the failed step.
D. A transaction with a longer timeout.

**Q3 (MCQ, medium).** Which property do sagas *lack* compared with ACID transactions?
A. Atomicity
B. Durability
C. Isolation
D. Consistency

**Q4 (short, medium).** Give two strategies to avoid needing a distributed transaction in the first place.

**Q5 (application, medium).** Order the steps of a saga for booking a trip: `sendConfirmationEmail`, `reserveHotelRoom (hold)`, `chargeCard`, `reserveFlightSeat (hold)`. Identify the pivot and the compensable and retriable steps.

**Q6 (MCQ, medium).** The payment call times out. What should the orchestrator do?
A. Treat it as failed and compensate immediately.
B. Treat it as unknown: query the payment status or retry with the same idempotency key before deciding.
C. Assume success and continue.
D. Drop the saga.

**Q7 (short, hard).** Explain "semantic countermeasures" against the lack of isolation in sagas with two examples.

**Q8 (MCQ, hard).** Which invariant can a CRDT-based design maintain *without* coordination?
A. Usernames must be unique.
B. Account balance must never go below zero.
C. A like counter that only ever increases and merges by summing per-replica counts.
D. Exactly one leader at a time.

**Q9 (application, hard).** Design the compensation for `chargeCard` and explain two ways it can be problematic compared with "undo".

**Q10 (short, hard).** Why is reconciliation necessary even in a well-designed saga system?
