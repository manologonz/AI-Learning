# Unit 4 Exam: Data Storage Fundamentals

Answers in `exam-answers.md`.

---

### Scenario (Q1-Q5): TicketNow

TicketNow sells event tickets. Tables: `events(id, name, starts_at, venue_id, capacity)`, `seats(id, event_id, section, row, number, status, held_until, version)`, `orders(id, user_id, event_id, total, created_at)`, `order_items(order_id, seat_id)`. Traffic: 20M users, browse 50:1 vs buy; hot on-sale moments (50k users chase 20k seats within a minute). Requirements: a seat is never sold twice; browsing is fast; users see "your orders" newest-first.

**Q1 (apply, medium).** The purchase code does: `SELECT status FROM seats WHERE id=$1` then, if 'available', `UPDATE seats SET status='sold'`. Show the interleaving that oversells a seat under Read Committed, and give two distinct correct fixes with SQL.

**Q2 (analyze, hard).** Customers buying *multiple* seats together must get all or none. A customer has 4 seats in one order; another customer wants 2 of the same seats. Describe the transaction you'd write, what isolation level or locking you'd use, and how you avoid deadlocks.

**Q3 (apply, medium).** "My orders" is slow for heavy buyers. Propose an index (write it) and explain the read benefit and the write cost. Estimate index size for 500M orders assuming ~40 bytes per entry.

**Q4 (evaluate, hard).** A PM asks for full-text event search with typo tolerance and filtering by city/date. Where does it live, how does it stay in sync with the source of truth, and what consistency does the user see?

**Q5 (evaluate, hard).** During on-sale, seat-hold writes spike to 30k/s and Postgres CPU is pegged. Give an analysis of options (at least three), including what you'd measure first, and pick one with a stated trade-off and reversal condition.

---

### Independent items

**Q6 (MCQ, medium).** Which is *true* about MVCC?
A. Readers block writers.
B. Each transaction reads from a snapshot of row versions, so readers don't block writers, but old versions must be cleaned up.
C. It removes the need for isolation levels.
D. It stores only the latest version of each row.

**Q7 (MCQ, hard).** A wide-column store (Cassandra-style) table `messages_by_channel(channel_id, ts, msg)` is queried also by user ("all messages by user X"). What is the idiomatic approach?
A. Add a secondary index on user and scan across all partitions at query time.
B. Create a second table `messages_by_user(user_id, ts, msg)` and write to both, accepting duplicated data and the responsibility to keep them consistent.
C. Use joins in the query.
D. Store everything in one giant row.

**Q8 (MCQ, hard).** You run `EXPLAIN` and see a sequential scan on a 200M row table for `WHERE status = 'active'` where 95% of rows are 'active'. What is the most likely conclusion?
A. You must create an index on status immediately.
B. An index would not help much (low selectivity); the planner is right to scan; consider whether the query needs 95% of rows at all or use a partial index on the rare values.
C. The statistics are wrong; disable the planner.
D. Convert the column to JSON.

**Q9 (short, hard).** Explain why LSM compaction can cause p99 latency spikes and how a team can mitigate it (two techniques).

**Q10 (create, hard).** Design the data model for a chat app's message history (10B messages, reads: last 50 messages in a channel, scroll back; writes: 50k/s). Choose the store, the primary/partition key with clustering order, and identify one hot-partition risk with a mitigation.
