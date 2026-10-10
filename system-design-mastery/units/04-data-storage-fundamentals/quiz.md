# Unit 4 Quiz: Data Storage Fundamentals

Answers in `quiz-answers.md`.

**Q1 (MCQ, easy).** Which storage engine design is optimized for very high write throughput through sequential writes and background compaction?
A. B-tree
B. LSM-tree
C. Hash index
D. Heap file with no index

**Q2 (MCQ, easy).** For an index on `(user_id, created_at)`, which query can use it efficiently?
A. `WHERE created_at > '2025-01-01'`
B. `WHERE user_id = 42 ORDER BY created_at DESC LIMIT 10`
C. `WHERE created_at = '2025-01-01' AND email = 'x'`
D. `WHERE email = 'x'`

**Q3 (MCQ, medium).** Two transactions each read balance = 100, then each writes balance = 100 - 30 using values computed in application code. Final balance is 70 instead of 40. What anomaly is this?
A. Dirty read
B. Phantom read
C. Lost update
D. Write skew

**Q4 (short, medium).** Give two reasons to *not* create an index on every column of a heavily written table.

**Q5 (MCQ, medium).** Which is the best fit as a primary system of record for financial ledger entries requiring multi-row transactions and strong constraints?
A. A key-value cache
B. A relational database
C. A full-text search engine
D. An object store

**Q6 (MCQ, medium).** Which statement about Postgres's default isolation level (Read Committed) is correct?
A. It prevents non-repeatable reads.
B. Each statement sees a snapshot of data committed before that statement began, so two reads of the same row in one transaction can differ.
C. It prevents lost updates for all read-modify-write patterns.
D. It is equivalent to Serializable.

**Q7 (application, medium).** A query fetches `SELECT status, total FROM orders WHERE customer_id = ? ORDER BY created_at DESC LIMIT 20`. Design the best index and explain in two sentences why.

**Q8 (short, hard).** Explain write skew with a concrete example and name two different fixes.

**Q9 (MCQ, hard).** You have an append-heavy event workload at 200k writes/s with key-based range reads per device. Which choice is best justified?
A. A single Postgres primary with 8 indexes
B. An LSM-based wide-column or time-series store partitioned by device
C. A search engine as system of record
D. An in-memory cache with no persistence

**Q10 (short, hard).** A teammate says "we should move to MongoDB because our Postgres is slow." Give three questions you'd ask before agreeing, and one situation in which the move *would* be justified.
