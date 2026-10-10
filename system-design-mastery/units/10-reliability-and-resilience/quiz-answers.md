# Unit 10 Quiz: Answer Key

**Q1: B.**
- Random sleeps in [0, cap] spread retries and avoid thundering herds. A is wrong (jitter adds waiting), C false, D false (timeouts are still needed).

**Q2: B.**
- Half-open allows a few probes to test recovery. A describes open; C describes closed; D is nonsense.

**Q3: C (27).**
- Attempts multiply across layers: 3 (client) × 3 (gateway) × 3 (service) = 27 DB attempts, arriving exactly when the DB is struggling.
- A ignores the layers. B (9) counts only two layers or adds wrongly. D (6) adds (3+3) instead of multiplying. Fix: retry at one layer, with a budget.

**Q4: 50 succeed immediately; ~6 seconds for the rest.**
- Bucket capacity is 50, so 50 succeed; 30 remain. Refill 5/s → 30 tokens take 6 s.
- Insight: capacity = max burst; refill = sustained rate.

**Q5: B.**
- Waiting callers hold threads/connections; resource exhaustion cascades. A/C/D are irrelevant.

**Q6. Sample answer.**
1. **Trigger/purpose:** rate limiting enforces per-client quotas/fairness regardless of load; load shedding reacts to the *system's* overload, regardless of who's sending.
2. **Selectivity:** rate limits are per key/tenant with fixed thresholds; shedding is dynamic and priority/cost-based (drop low-priority or expensive requests first).
- Also acceptable: limits are proactive/policy; shedding is reactive/protective.

**Q7: B.**
- Goodput collapsed: the server is busy producing responses nobody will read. Use bounded queues, shed early, propagate deadlines and drop requests whose deadline has passed. A is wrong (throughput without goodput). C makes it worse. D hides the problem and increases latency and resource holding.

**Q8. Sample answer.**
- **Pilot light:** minimal core (data replicated, servers off/small) in region B; on disaster, scale up and switch traffic: RTO tens of minutes to hours, RPO = replication lag (seconds-minutes); low cost, moderate complexity (runbooks, automation, regular tests).
- **Active-active:** both regions serve traffic; RTO near zero (traffic shifts), RPO near zero-seconds; highest cost (full capacity ×2 or headroom) and complexity (conflict resolution/data partitioning, global routing, consistency trade-offs).

**Q9. Sample answer.**
- Algorithm: **token bucket** per API key (capacity 20, refill 100/60 ≈ 1.67 tokens/s). It allows bursts up to 20, sustained 100/min.
- State: Redis (key = `rl:{apiKey}`), updated atomically via a Lua script (tokens + timestamp), TTL to expire idle keys; or local limiters with a distributed share for scale.
- If Redis fails: **fail-open** for availability (allow traffic, use local approximate limits per node as a fallback, alert), unless the API is security-sensitive (login), in which case fail-closed or use stricter local limits.

**Q10. Sample answer.**
- Backups may be corrupt, incomplete, or unrestorable, and recovery may be too slow for the RTO. Verify: (1) **restore tests**: you can restore fully and the app works with it; (2) **RTO/RPO measured**: restore time and data currency meet objectives; (3) **isolation/integrity**: backups are immutable, stored in a separate account/region, and encryption keys are available; also dependencies like DNS, secrets, and infrastructure-as-code needed to rebuild.
