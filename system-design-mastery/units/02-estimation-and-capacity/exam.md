# Unit 2 Exam: Estimation & Capacity Planning

Summative. Show assumptions and arithmetic. Answers in `exam-answers.md`.

---

### Scenario (Q1-Q6): PicShare

PicShare is a photo-sharing app. Assumptions given by product: 30M DAU; each user views 100 photos/day; 5% of users upload one photo/day. Each uploaded photo is stored in 3 renditions totalling 4 MB. When viewing, the client downloads a ~200 KB rendition. Peak traffic is 3x average. Use 86,400 s/day.

**Q1 (apply, easy).** Compute average and peak read QPS (photo views) and average and peak upload QPS. What is the read:write ratio?

**Q2 (apply, medium).** Compute daily storage growth, yearly growth, and yearly growth with 3x replication. Would you use erasure coding (~1.5x overhead) instead? State the saving and one cost.

**Q3 (apply, medium).** Compute peak egress bandwidth in Gbps if every view is served from origin. Then compute the origin bandwidth if a CDN has a 90% hit rate.

**Q4 (analyze, hard).** Which resource is the dominant cost/constraint for PicShare: read QPS, write QPS, storage, or egress bandwidth? Defend your answer using your numbers, and say what architectural component this points to.

**Q5 (apply, hard).** The metadata (photo id, owner, timestamps, S3 key, etc.) is 1 KB per photo. How large is one year of metadata, and does it justify sharding on storage grounds alone? Read QPS for metadata equals your Q1 peak; explain what you'd do to serve it, given a Postgres primary handles ~10k simple QPS and Redis ~100k QPS per node.

**Q6 (evaluate, hard).** Product announces a campaign that will 4x uploads and add a viral event doubling peak read traffic for two weeks. Recompute the impact on (a) upload QPS, (b) storage/day, (c) peak egress at origin (assume same CDN hit rate). Which part of the plan needs pre-scaling, and which is already comfortable?

---

### Independent items

**Q7 (MCQ, medium).** A design review states: "Our service handles 5,000 QPS on average, and each server does 1,000 QPS, so we need 5 servers." Which critique is the most complete?
A. It's fine; average QPS is the correct sizing metric.
B. It ignores peak factor, target utilization, and redundancy across failure domains; realistic sizing could be ~2-3x more.
C. It should use bytes rather than QPS.
D. It's wrong only because 1,000 QPS per server is too low for any real server.

**Q8 (MCQ, hard).** Which statement about estimation is most accurate?
A. If your estimate is off by 30%, the design is invalid.
B. Estimates should be as precise as possible since they're the foundation of the design.
C. Being within the right order of magnitude is what matters, provided assumptions are explicit and the dominant resource is identified.
D. Estimation is only needed in interviews; production teams use load tests.

**Q9 (application, hard).** An IoT platform has 2M devices each sending a 200-byte reading every 10 seconds. Compute (a) average ingest QPS, (b) ingest bandwidth, (c) storage per year raw (before replication), and (d) name the resource that will worry you first and why.

**Q10 (create, hard).** Choose any product you know (e.g., a ride-hailing app, a video conferencing tool, an e-commerce checkout). Write a one-page estimation: state 5 assumptions, derive avg/peak QPS, storage/year, and bandwidth, and finish with a "so what": the two architectural decisions your numbers force.
