# Unit 2: Estimation & Capacity Planning

## Objectives
By the end of this unit you can:
- Convert user counts into average and peak QPS, storage growth, and bandwidth.
- Use latency orders of magnitude to judge whether a design is plausible.
- Size a server fleet with headroom and failure-domain (AZ) redundancy.
- Sanity-check a design against numbers and identify which resource is the bottleneck.

## Prerequisites
Unit 1 (requirements and NFRs). Estimation turns "10M users" into engineering constraints.

## Core content

### 1. Why estimate at all?
Estimates decide *whether you need a distributed system*. 200 writes/s fits one Postgres node; 200k writes/s does not. Seniors estimate to:
- Kill unnecessary complexity ("that's 30 QPS, we don't need Kafka").
- Find the dominant cost or bottleneck (storage? bandwidth? QPS? connections?).
- Justify a design in a review with numbers instead of taste.

The goal is the **right order of magnitude**, not precision. 10k vs 12k QPS never changes a design; 10k vs 1M does.

### 2. The toolkit

**Powers and time**
- 1 day = 86,400 s ≈ **10^5 s** (round up for easy math; error ~15%).
- 1M requests/day ≈ **~12 QPS**. 1B/day ≈ **~12k QPS**.
- 1 year ≈ 3×10^7 s. 1 month ≈ 2.6×10^6 s.

**Data sizes**: 1 KB = 10^3 B, 1 MB = 10^6, 1 GB = 10^9, 1 TB = 10^12, 1 PB = 10^15. (Use decimal; the ~2.4% binary difference doesn't matter.)

**Typical object sizes**: a tweet/text row ≈ 200-500 B; a DB row of metadata ≈ 1 KB; a compressed web image ≈ 200 KB-2 MB; a phone photo ≈ 2-5 MB; a minute of 1080p video ≈ 50-100 MB; a log line ≈ 200 B.

**Latency numbers (orders of magnitude)**

| Operation | Time |
|---|---|
| L1 cache / CPU register ref | ~1 ns |
| Main memory reference | ~100 ns |
| SSD random read | ~100 µs (1,000x memory) |
| Round trip within same data center | ~0.5 ms |
| Read 1 MB sequentially from SSD | ~1 ms |
| HDD seek | ~10 ms |
| Round trip cross-continent | ~100-150 ms |

Rules of thumb: memory is ~1000x faster than SSD; SSD ~100x faster than HDD seeks; a network hop in a DC costs about a thousand memory reads; **crossing an ocean costs a hundred million CPU cycles**, so keep chatty calls off long links.

**Rough server capacity (CRUD app server with DB calls)**: ~500-2,000 QPS per commodity instance; a well-tuned Postgres primary handles roughly 5-20k simple QPS (far fewer for heavy writes/joins); Redis ~100k ops/s per node; a Kafka broker roughly tens-hundreds MB/s. Treat these as anchors; say "assume 1k QPS per app server" out loud and move on.

### 3. The formulas

1. **Average QPS** = requests per day ÷ 86,400 (≈ ÷10^5).
2. **Peak QPS** = average × peak factor (2-3x for daily cycles; 5-10x for events like a sale or a goal in a World Cup game). State the factor.
3. **Read/write QPS**: derive both separately; the ratio decides caching and replication.
4. **Storage/day** = writes/day × object size; **yearly** = ×365; **with replication** multiply by copies (3x typical, 1.5x with erasure coding). Add a 20-30% cushion for indexes/overhead.
5. **Bandwidth** = QPS × payload size (bytes/s → ×8 for bits/s). Ingress (writes) and egress (reads) separately.
6. **Fleet size** = peak QPS ÷ (per-server QPS × target utilization), then add redundancy.
7. **Memory for a cache** ≈ hot working set: if 20% of items serve 80% of traffic, cache ≈ 20% of daily-accessed data (see Unit 5).
8. **Concurrent connections** = concurrent users × connections per user (matters for chat, live location, websockets more than QPS).

**Target utilization**: run at 50-70% of tested capacity so you can absorb spikes, deploys, and failures. **Failure domains**: if you want to survive losing one of 3 AZs with the remaining 2 carrying full peak, provision 150% of the required capacity: each of the surviving 2 AZs must hold half the peak fleet, so the total across 3 AZs is 1.5 × the single-fleet requirement.

```ts
// Reusable estimation helpers
const SECONDS_PER_DAY = 86_400;
const avgQps = (perDay: number) => perDay / SECONDS_PER_DAY;
const peakQps = (perDay: number, factor = 3) => avgQps(perDay) * factor;
const storagePerYearTB = (perDay: number, bytes: number, replication = 3) =>
  (perDay * bytes * 365 * replication) / 1e12;
const serversNeeded = (peak: number, perServer: number, util = 0.6) =>
  Math.ceil(peak / (perServer * util));
// N+AZ redundancy: survive losing 1 of `azs` zones
const withAzRedundancy = (servers: number, azs = 3) =>
  Math.ceil((servers / (azs - 1))) * azs;

serversNeeded(30_000, 1_500);      // 34
withAzRedundancy(34);              // 51
```

### 4. Method: always do these in order
1. State assumptions explicitly (DAU, actions/user/day, sizes). Interviewers/reviewers care more about assumptions than arithmetic.
2. Derive average, then peak.
3. Separate reads and writes.
4. Compute storage, bandwidth.
5. Identify the **dominant resource** and check it against single-machine limits. That tells you what must be distributed.
6. Sanity check: does the result feel plausible? (The whole internet's traffic is Tbps-scale; if your startup needs Pbps, you erred.)

## Worked examples

### Straightforward: QPS
*100M DAU, each makes 10 API reads per day. Average and peak QPS?*
- 100M × 10 = 1B requests/day.
- ÷ 86,400 ≈ **11.6k QPS** (≈ 12k). Peak at 3x ≈ **35k QPS**.
- One app server at ~1k QPS → ~35 servers at 100% utilization, ~58 at 60%. Not a "one box" problem, but not exotic either.

### Typical: storage and bandwidth for a text+image service
*300M DAU, 10% post once per day (30M posts/day). Each post: 300 B text; 20% include a 1 MB image. Reads: each user views 50 posts/day, each ~1 KB payload (text + metadata), images served separately.*

1. **Write QPS**: 30M/86,400 ≈ 350 avg; peak 3x ≈ 1,000. Trivial for one primary DB.
2. **Read QPS**: 300M × 50 = 15B/day ÷ 86,400 ≈ 174k avg; peak ≈ **~500k QPS**. Read:write ≈ 500:1. **This** is the hard part: reads dominate → caching and replicas essential.
3. **Storage/day**: text: 30M × 300 B = 9 GB/day. Images: 30M × 20% × 1 MB = 6M MB = **6 TB/day**. Image storage dominates text by ~700x. Yearly images: 6 TB × 365 ≈ 2.2 PB; ×3 replication ≈ 6.6 PB. Conclusion: images live in object storage + CDN; the DB holds only metadata and stays small (9 GB/day ≈ 3.3 TB/yr, one shard for years).
4. **Egress**: if 30% of views load an image at 1 MB: 500k × 30% ≈ 150k image fetches/s at peak × 1 MB = 150 GB/s = **1.2 Tbps**. Impossible to serve from origin; CDN with ~95% hit rate leaves ~60 Gbps at origin. Bandwidth (cost!) is the true bottleneck.

### Edge case: the number that hides the real constraint
*A chat app: 50M DAU, each sends 40 messages/day. Average QPS?* 2B/day → ~23k QPS. "Easy," you say. But each *online* user holds a long-lived WebSocket. If peak concurrency is 20% of DAU = **10M concurrent connections**. At ~50k connections per server (tuned) that's 200 servers just to hold connections, before any message work, and memory per connection (~10-50 KB) becomes ~100-500 GB total. **The dominant resource is connections/memory, not QPS.**

Lesson: QPS isn't always the bottleneck. Check connections, memory, disk IOPS, network, and per-item fan-out (one post × 10M followers = 10M writes).

## Common misconceptions

- **"Precision matters."** It doesn't; orders of magnitude do. Spending five minutes to get 11.6k instead of 12k is wasted effort.
- **"Use average QPS to size the system."** Systems fail at peak. Always apply a peak factor and utilization headroom.
- **"Storage growth is the same as data volume produced today."** Forget replication (3x), indexes, backups, and retention and you under-estimate by 3-5x.
- **"QPS is the only load metric."** Connections, bandwidth, storage IOPS, memory, and fan-out amplification are often the binding constraint.
- **"Bits and bytes are interchangeable."** Bandwidth is quoted in bits/s; multiply bytes/s by 8.
- **"Read all latency numbers from memory verbatim."** Memorize the *ratios* (memory : SSD : network : disk) and rounded values; the exact numbers vary by hardware generation.
- **"Provision exactly to capacity."** No headroom → deploys, GC pauses, and a single AZ failure become outages.

## Check your understanding

1. 5M requests/day is roughly how many QPS? *(5M/86,400 ≈ 58 QPS ≈ 60.)*
2. A service stores 2M records/day at 5 KB each with 3x replication. How many TB per year? *(2M × 5 KB = 10 GB/day; × 365 = 3.65 TB; × 3 = ~11 TB/yr.)*
3. 40k peak QPS, 2k QPS per server, 50% target utilization. Servers needed? *(40,000 / (2,000×0.5) = 40.)*
4. Why can a system with only 500 average write QPS still need sharding? *(Data volume, e.g. hundreds of TB, fan-out amplification, or hot keys; QPS isn't the only resource.)*

## Summary
- Estimate to choose architecture and expose the dominant resource; aim for order of magnitude.
- 1 day ≈ 10^5 seconds; 1M req/day ≈ 12 QPS; 1B/day ≈ 12k QPS.
- Derive average → peak → read/write split → storage (with replication) → bandwidth → fleet with utilization and AZ headroom.
- Memory ≈ 1000x faster than SSD, SSD ≈ 100x faster than HDD seek, a DC hop ≈ 0.5 ms, cross-continent ≈ 100+ ms.
- Check non-QPS resources: connections, bandwidth, storage, fan-out.
