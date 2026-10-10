# Unit 2 Quiz: Estimation & Capacity Planning

Formative and low stakes. Show your work. Answers in `quiz-answers.md`.

**Q1 (MCQ, easy).** Roughly how many QPS does 1 billion requests per day represent (average)?
A. ~120
B. ~1,200
C. ~12,000
D. ~120,000

**Q2 (MCQ, easy).** Which operation is roughly 1,000x slower than a main-memory reference?
A. An L1 cache reference
B. An SSD random read
C. A round trip within the same data center
D. An HDD seek

**Q3 (short, easy).** A service has 50M DAU and each user makes 20 requests/day. What's the average QPS, and the peak QPS assuming a 3x peak factor?

**Q4 (application, medium).** A photo app receives 10M new photos/day, 2 MB each. How much raw storage per day, and per year with 3x replication? Give answers in TB/PB.

**Q5 (application, medium).** At peak, 20,000 requests/s each return a 500 KB payload. What egress bandwidth is that in Gbps?

**Q6 (MCQ, medium).** A team estimates capacity from average QPS and provisions servers at exactly 100% of it. What's the most important flaw?
A. Average QPS can't be computed from daily totals.
B. Traffic peaks above average and servers need headroom for deploys, failures, and spikes.
C. Servers should be sized by storage, never by QPS.
D. Peak factors only apply to databases.

**Q7 (application, medium).** A chat app has 30M DAU and 15% of them are online at peak, each holding one WebSocket connection. One tuned server holds 50k connections. How many servers are needed just for connections (before redundancy)? Why might this be a more important number than message QPS?

**Q8 (MCQ, hard).** You estimate 1 TB/year of *new* rows for a service and conclude "1 TB is nothing, one machine is fine forever." Which omission most undermines that conclusion?
A. You forgot that disks can't hold 1 TB.
B. You ignored replication, indexes, backups, retention, and growth; and didn't check IOPS/throughput or working-set-in-memory, not just capacity.
C. You should have used bits instead of bytes.
D. 1 TB/year is always too large for one machine.

**Q9 (application, hard).** 60,000 peak QPS. One app server handles 2,000 QPS at test; you target 60% utilization and want to survive losing one of 3 AZs with the remaining two still carrying peak. How many servers total?

**Q10 (short, hard).** A social feed: one post by a user with 5M followers triggers "fan-out on write" to each follower's feed. If 1,000 such celebrity posts happen per hour, how many feed writes per second on average does this create (give the order of magnitude)? Explain why this can dominate a system whose "normal" write QPS looks small.
