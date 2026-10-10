# Unit 9 Exam: Distributed Transactions & Data Consistency Patterns

Answers in `exam-answers.md`.

---

### Scenario (Q1-Q5): TripBook

TripBook sells bundled trips: flight + hotel + car. Each is a separate internal service with its own database and external supplier APIs (some slow, some without idempotency support). Rules: never charge a customer for a bundle they don't fully get; holds on suppliers expire after 10-15 min; the payment gateway supports authorize (voidable) and capture. Traffic: 300 bookings/s at peak. Past incident: the hotel supplier timed out, the orchestrator assumed failure and released the flight; the hotel actually confirmed; the customer had a hotel but no flight and was charged.

**Q1 (analyze, medium).** Diagnose the incident's root cause in terms of unknown outcomes, and describe the fix (three concrete changes).

**Q2 (apply, hard).** Design the booking saga: list steps in order, mark the pivot, give each step's compensation, and state which are retriable after the pivot. Justify the order (consider hold expiry, failure likelihood, and cost of undo).

**Q3 (apply, hard).** The car supplier's API has no idempotency key and sometimes duplicates a booking on retry. How do you design the car step so retries can't double-book? What residual risk remains?

**Q4 (analyze, hard).** While the saga runs, the "my bookings" page shows partial state and customer-support agents can cancel a booking mid-saga. Enumerate two anomalies caused by lack of isolation and give a countermeasure for each.

**Q5 (evaluate, hard).** Should TripBook use an orchestrator (Temporal/Step Functions) or choreography? Give a recommendation, including how the orchestrator survives its own crash, and one failure it introduces.

---

### Independent items

**Q6 (MCQ, medium).** Which statement is true about 2PC and Spanner?
A. Spanner avoids 2PC entirely.
B. Spanner uses 2PC across Paxos groups (each group replicated), which mitigates coordinator failure blocking because the participants and coordinator state are highly available replicated state machines.
C. Spanner uses sagas.
D. Spanner is eventually consistent.

**Q7 (MCQ, hard).** Which step order is best for a saga: A (charge card, hard to undo), B (create shipping label, cheap to cancel), C (reserve inventory, cheap to cancel)?
A. A, B, C
B. B, C, then A last only if all else succeeded, with A as the pivot
C. A, C, B
D. Any order is fine.

**Q8 (short, hard).** Write the state machine (states and allowed transitions) for an order in a saga, and explain how it rejects out-of-order events.

**Q9 (application, hard).** A double-entry ledger records a purchase and a later refund. Show the ledger entries for a $50 purchase with a $50 refund, and explain why it aids reconciliation compared with updating a balance column.

**Q10 (create, hard).** Choose the pattern (single-DB transaction, 2PC, saga, CRDT, or "relax the invariant") for each and justify in one line: (a) decrementing stock and creating an order in the same service/DB; (b) transferring loyalty points between two users in different regions where brief inconsistency is acceptable; (c) booking a flight across two separate airline systems; (d) a collaborative shopping list.
