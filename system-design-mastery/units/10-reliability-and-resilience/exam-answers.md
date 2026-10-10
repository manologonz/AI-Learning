# Unit 10 Exam: Answer Key

**Q1. Flaws → techniques (need ≥5).**
1. 30 s default timeout to an optional dependency → **tight timeouts** (~200-300 ms) + deadline propagation.
2. Recommendations is optional but coupled on the critical path → **fallback/graceful degradation** (skip or serve cached recommendations).
3. Shared thread pool → one slow dependency exhausts all threads → **bulkheads** (per-dependency concurrency limits/pools).
4. No circuit breaker → continued calls to a failing service → **circuit breaker** with half-open probing.
5. Multi-layer aggressive retries (client + gateway) → amplification → **single-layer retries, retry budget, jittered backoff**.
6. Retries reached the DB and no admission control → **load shedding** and protection of the DB (rate limits, connection pooling limits).
7. Bad deploy hit everything at once → **canary/gradual rollout with automated rollback on SLO burn**.
8. Slow recovery/manual traffic control → **runbooks, slow start, automatic recovery ramp**.

**Q2. Sample design.**
- Timeout: recommendations normally 50 ms; p99.9 maybe ~150 ms → **timeout 250 ms** (total request budget for playback start is 2 s; recs is optional so keep small).
- Bulkhead: Little's Law `L = λ × W`: 500 QPS × 0.05 s = **25 concurrent** in steady state. With headroom (2-3x for bursts and slower moments) → **concurrency limit ~60-75** for recommendations per instance (e.g., 60), separate from the main pool; beyond that fail fast to the fallback. During degradation (all calls hitting 250 ms timeout): 500 × 0.25 = 125 concurrent would be needed, but the limit of 60 rejects the excess immediately and protects threads.
- Circuit breaker: open when failure/timeout rate ≥ 50% over a 10 s window with min 20 requests; stay open 15-30 s; half-open with 5 probe requests; close after ≥4/5 succeed.
- Fallback: return cached last-known-good recommendations (per user or global "popular now" list, TTL ~1 h) or an empty list with the UI hiding the row. Playback start never waits on recommendations at all (call it asynchronously/lazy).

**Q3.**
- Sustaining loop: the rollback fixed the trigger, but the system was in overload: DB CPU saturated, latency high, requests timing out → clients/gateway retried (multiplying load ~9x at the worst) → more timeouts → more retries. Queues held old, already-expired requests (goodput ≈ 0). The cache had also emptied/expired during the outage, sending more reads to the DB. The retry load itself exceeded capacity even without the original fault.
- Recovery procedure: shed load hard (block or cap traffic at the edge at a fraction of capacity), let queues drain and the DB recover, warm caches, then ramp traffic up gradually (e.g., 10% steps) while watching latency/error rates.
- Design changes: retry budgets/one-layer retries, adaptive load shedding (admit at DB capacity), circuit breaker on the DB path, bounded queues + deadline propagation (drop expired work), cache warming/serve-stale, automated slow-start ramp on recovery, and chaos/load tests for "dependency slow + retry" scenarios.

**Q4.**
- Amplification: 3 client attempts × 3 gateway attempts = **9 requests** to the backend per client request; if `playback` retries to services/DB as well, more. At 50k QPS, 9x = 450k QPS.
- Policy: **retry only at the client (or only at the gateway)**: 2 total attempts (1 retry), only for idempotent requests and retryable errors (503/timeouts/429 with `Retry-After`), **full-jitter exponential backoff** `sleep = rand(0, min(2 s, 100 ms × 2^n))`, a **retry budget** of ≤10% of requests (token bucket for retries: each success adds 0.1 tokens; a retry costs 1), and a deadline so retries stop when the overall budget (e.g., 2 s) is used. Internal services don't retry (they fail fast and propagate errors).

**Q5.**
- Detection: queue wait time/in-flight concurrency vs. the tested limit, CPU > 85%, DB pool wait time, p99 latency > 2x normal, or error rate; use latency-based (CoDel-like) or adaptive concurrency limits rather than static CPU.
- Shed order (lowest value first): (1) background/prefetch calls, low-priority endpoints (recommendation refresh, analytics), (2) browse page enrichments, (3) new-session browse traffic, (4) *never* shed playback for already-playing or paying users first; health checks and playback starts have the highest priority given the SLO is on playback.
- Response: `503` (or `429`) with `Retry-After` and a jittered hint; clients back off; cheap rejection at the gateway before touching services/DB. Justification: the SLO is playback start success; sacrificing optional traffic protects the 99.95% objective and its error budget.

**Q6. Sample options.**
- **Backup & restore:** RTO hours, RPO up to backup interval (hours); cheapest; doesn't meet 30 min/1 min.
- **Pilot light / warm standby:** replicate DB async cross-region (RPO seconds-minute), core services deployed scaled-down; on failure scale up and shift DNS/anycast: RTO ~15-30 min; moderate cost (~20-40% extra).
- **Active-active:** RTO ~0, RPO ~0-seconds; highest cost/complexity (data conflicts, global routing).
- Recommendation: **warm standby**: continuous cross-region replication for billing with RPO ~<1 min (use sync or semi-sync for billing writes if 1 minute is too risky), degraded browse from replicated/cached content, automated failover runbook tested quarterly; active-active for the CDN/video path is already effectively provided by the CDN. Meets 30 min/1 min at moderate cost.

**Q7: B.**
- 0.999 × 0.9995 × 0.9999 ≈ 0.99840 → **99.84%**. A/D ignore the multiplication; C is the best component, not the system.

**Q8: B.**
- Deadlines shrink hop by hop so downstream work stops when it can't complete in time, saving capacity. A is wrong (that's static per-hop timeouts, losing the budget), C wrong, D wrong (complementary).

**Q9. Sample answer.**
- Static stability: a system continues operating correctly when a dependency (typically a control plane) is unavailable, by not requiring changes during the failure. Example: keep enough pre-provisioned capacity across AZs so that losing an AZ doesn't require launching new instances (the launch API might also be impaired); use cached config/credentials/DNS so data-plane traffic doesn't depend on control-plane calls.

**Q10. Sample game-day plan.**
- **Title:** Loss of AZ-b for StreamBox production (Tuesday 11:00, off-peak).
- **Hypothesis:** With N+AZ capacity headroom (≥150% at peak), losing AZ-b causes no SLO violation: playback start success stays ≥99.95%, and failover of Postgres replicas/leader completes within 60 s.
- **Injected failure:** Block all traffic to/from AZ-b instances (network ACL/instance stop via chaos tool), including the DB replica or primary if there; terminate 100% of AZ-b service instances.
- **Blast radius:** Limit to one AZ, 30 minutes; alert on-call; feature flags ready; only during business hours with engineers present.
- **Metrics:** playback success rate and p99 latency, error rate per service, DB failover time, replication lag, cache hit rate, autoscaling events, retries per request, queue depths.
- **Abort criteria:** SLO burn rate > 10x for 5 minutes, or error rate >2% for playback, or any data-integrity anomaly → restore the AZ immediately.
- **Success:** stays within SLO (or brief, bounded degradation within the error budget), failover < 60 s, no manual scaling required; findings documented with action items (e.g., raise capacity headroom, fix hard-coded AZ endpoints).
