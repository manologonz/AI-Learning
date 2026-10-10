# Unit 2 Exam: Answer Key

**Q1.**
- Reads: 30M × 100 = 3B/day ÷ 86,400 ≈ **34.7k QPS avg** (~35k); peak ×3 ≈ **~104k QPS**.
- Uploads: 30M × 5% = 1.5M/day ÷ 86,400 ≈ **17 QPS avg**; peak ≈ **52 QPS**.
- Ratio: 3B / 1.5M = **2,000:1** (extremely read-heavy).

**Q2.**
- Daily: 1.5M × 4 MB = 6,000,000 MB = **6 TB/day**. Yearly: × 365 = **2.19 PB**. With 3x replication: **~6.6 PB**.
- Erasure coding (~1.5x): ~3.3 PB, saving ~3.3 PB (~50%) of storage cost. Costs: higher CPU/latency on reads and repairs, more complex rebuilds, small objects handle poorly. Fine for cold/warm large objects (photos), less suited to hot data.
- Common error: mixing MB/TB.

**Q3.**
- Peak: 104k × 200 KB = 20.8 GB/s ≈ 20 GB/s × 8 = **~166 Gbps** (≈ 160-170).
- With 90% CDN hit: origin sees 10% → **~16.6 Gbps**.
- Insight: serving 166 Gbps from origin is very expensive; a CDN cuts origin load 10x and (more importantly) serves closer to users.

**Q4.**
- **Dominant: egress bandwidth (and storage as the secondary).** Read QPS ~100k peak is big but cacheable/CDN-able; write QPS (~52) is trivial; storage is 6.6 PB/yr (large but cheap per GB in object storage); egress at 166 Gbps peak is the recurring, highest-cost item, and the one that most constrains the architecture.
- Points to: **CDN + object storage** with aggressive caching and appropriately sized renditions. Accept "storage" if defended with cost numbers, but egress is the stronger answer since it recurs continuously and scales with reads.

**Q5.**
- 1.5M photos/day × 1 KB × 365 ≈ 0.55 TB/yr. **Storage alone doesn't justify sharding**; one Postgres node holds this for many years.
- Read QPS is the issue: ~104k peak vs a primary's ~10k. Solution: cache metadata in Redis (100k ops/s/node; e.g., 3-5 nodes for headroom) with cache-aside; add read replicas for misses (with 95% hit rate, DB sees ~5k QPS, within one primary's range; replicas add headroom). Also the CDN already absorbs image reads, but the metadata lookups are per feed page, so cache-aside is the main lever.
- Full credit: sizes storage correctly, doesn't shard needlessly, identifies read QPS as the metadata constraint and proposes caching with numbers.

**Q6.**
- (a) Uploads 4x: 52 → ~208 QPS peak (still trivial).
- (b) Storage 4x: 24 TB/day (~8.8 PB/yr raw); need to verify object-store budget and ingest bandwidth: 208 QPS × 4 MB ≈ 0.83 GB/s ≈ 6.7 Gbps ingress; fine.
- (c) Peak reads doubled: 208k QPS × 200 KB = 41.6 GB/s ≈ 333 Gbps total; origin at 10% = **~33 Gbps**.
- Needs pre-scaling: **CDN capacity/commit and origin egress** (and metadata cache tier: 208k QPS → more Redis nodes). Comfortable: upload path and metadata DB writes. Also mention: ensure the CDN hit rate holds for a viral (hot) set: the hot set is small so hit rate likely improves.

**Q7: B.**
- Peak factor (2-3x), target utilization (~60%), N+1 or AZ redundancy: 5 → ~14-17 servers, easily 3x.
- A wrong: sizing on average fails at peak. C wrong: units aren't the issue. D wrong: 1,000 QPS per server is a plausible number; the flaw is the method.

**Q8: C.**
- Order of magnitude with explicit assumptions is the goal. A and B over-value precision (inputs are themselves guesses). D is wrong: estimates *precede* load tests and decide what to build; load tests validate.

**Q9.**
- (a) 2M / 10 s = **200k QPS** average ingest.
- (b) 200k × 200 B = 40 MB/s = **320 Mbps** (small).
- (c) 40 MB/s × 86,400 = 3.46 TB/day; × 365 = **~1.26 PB/year** raw.
- (d) **Write throughput/ingest QPS (200k writes/s)**: beyond a single relational primary; needs a partitioned log/time-series store (Kafka + TSDB/columnar). Storage (1.26 PB/yr) is second; bandwidth is minor. Accept storage retention/downsampling as a strong second point. 
- Common error: thinking 320 Mbps is a problem; the *rate of small writes* is the problem.

**Q10.** Sample structure (grade on method, not the product):
- Assumptions (5): DAU, actions/user/day, object sizes, peak factor, read/write ratio.
- Derived: avg/peak QPS (separate reads and writes), storage/day → /year with replication, bandwidth (in Gbps, ingress and egress).
- "So what": two forced decisions, e.g. "reads ≫ writes → cache + read replicas"; "storage in PB → object store not DB"; "connections in millions → dedicated connection tier".
- Full credit: explicit assumptions, correct unit handling, peak vs average distinguished, dominant resource named, and decisions follow from numbers.
