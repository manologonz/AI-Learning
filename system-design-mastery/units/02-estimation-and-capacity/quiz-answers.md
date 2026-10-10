# Unit 2 Quiz: Answer Key

**Q1: C (~12,000).**
- 10^9 / 86,400 ≈ 11.6k. Memorize: 1B/day ≈ 12k QPS.
- A: off by 100x (would be 10M/day). B: off by 10x (would be 100M/day). D: off by 10x the other way (~10B/day). These come from mis-remembering 10^5 seconds per day.

**Q2: B (SSD random read, ~100 µs vs ~100 ns).**
- A is wrong: L1 is ~100x *faster* than main memory. C is wrong: a same-DC round trip (~0.5 ms) is ~5,000x memory. D is wrong: an HDD seek (~10 ms) is ~100,000x memory.

**Q3: ≈ 11.6k avg, ≈ 35k peak.**
- 50M × 20 = 1B/day → ~11.6k QPS; × 3 → ~35k.

**Q4: 20 TB/day raw; ~7.3 PB/yr raw; ~22 PB/yr with 3x.**
- 10M × 2 MB = 20 TB/day. × 365 = 7,300 TB = 7.3 PB. × 3 = 21.9 PB.
- Common error: forgetting replication or mixing TB and PB.

**Q5: 80 Gbps.**
- 20,000 × 500 KB = 10,000,000 KB/s = 10 GB/s. × 8 = **80 Gbps**.
- Common error: forgetting ×8 (bytes vs bits) and stopping at "10".

**Q6: B.**
- Systems fail at peaks; no headroom leaves no room for deploys, GC pauses, node failure, or traffic spikes.
- A is wrong: averages are trivially computed. C is wrong: sizing depends on the binding resource (could be QPS, storage, connections, bandwidth). D is wrong: peak factors apply to any component that sees traffic.

**Q7: 90 servers.**
- 30M × 15% = 4.5M concurrent connections. 4.5M / 50k = **90 servers**.
- It matters because connection count (memory per connection, file descriptors) can be the dominant constraint even when message QPS looks modest; a fleet sized for QPS alone would be far too small, and connection servers need special handling for failover (reconnect storms).

**Q8: B.**
- Capacity is one axis. Real footprint = data × replication × indexes × backups; retention and growth over years; and performance depends on IOPS/throughput and whether the working set fits in RAM.
- A is wrong: modern disks hold multiple TB. C is a non-issue. D is wrong: 1 TB is small for storage but the reasoning, not the number, was flawed.

**Q9: 75 servers.**
- Servers to carry peak at 60% utilization: 60,000 / (2,000 × 0.6) = 50.
- Surviving the loss of 1 of 3 AZs means the remaining 2 AZs must together hold those 50 → 25 per AZ → 3 × 25 = **75** (i.e., 1.5x).
- Common error: answering 50 (ignores failure domains) or 100 (doubling for "redundancy" without reasoning about AZ count). With 3 AZs the overhead is 50%, with 2 AZs it would be 100%; more AZs means cheaper redundancy.

**Q10: ≈ 1.4 million feed writes/s on average.**
- 1,000 posts/hour × 5M followers = 5×10^9 writes/hour; ÷ 3,600 ≈ 1.4×10^6 writes/s.
- Why it dominates: the "post rate" looks tiny (~0.3 posts/s), but each post multiplies into millions of writes. Fan-out amplification, not raw QPS, is the bottleneck. That's why celebrity accounts use fan-out-on-read or hybrid models (Unit 13).
