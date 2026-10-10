# Unit 12 Quiz: Answer Key

**Q1: B.**
- Traces show per-span durations along the request path. Metrics aggregate (no per-request path), logs are unstructured and hard to correlate, health checks are binary.

**Q2: B.**
- Each distinct label value creates a new time series; millions of users → millions of series. Put per-user data in logs/traces. A/C/D are wrong.

**Q3: burn rate = 7.2; do not page under the 14.4x rule.**
- Allowed error rate = 0.1%. Observed 0.72% → 0.72 / 0.1 = **7.2x**. Below 14.4 (the 1 h page threshold) but above 6: it would trigger the 6 h window page *if sustained* for 6 hours (≥6x over 6 h). So it doesn't page immediately but it's on course to; investigate now.

**Q4: B.**
- Canary exposes a small slice of real traffic, compares SLIs, and rolls back automatically. A maximizes blast radius; C switches everything at once (rollback is fast but detection happens after 100% exposure); D lacks health checks.

**Q5. Steps.**
1. Add the new column (nullable). 2. Deploy code that writes both columns. 3. Backfill old rows in small batches. 4. Deploy code that reads the new column. 5. Deploy code that stops writing the old column. 6. Drop the old column after a soak period. (Expand → migrate → contract.)

**Q6: B.**
- JWTs are signed (not encrypted) and stateless, so revocation is hard: use short-lived access tokens + refresh tokens. A is false; C false (always validate signature, exp, aud, iss); D is the classic `alg=none` vulnerability, so pin algorithms.

**Q7: B.**
- IDOR/BOLA: the service didn't check that the requester owns the object. Gateway authentication doesn't imply authorization. Fix: enforce per-object authorization in the service (`WHERE id=$1 AND tenant_id=$2`, policy check), RLS as a backstop, tests. A/C/D are different vulnerabilities.

**Q8. Sample answer.**
- **Shared schema:** lowest cost/highest density, simplest ops (one migration), but weakest isolation (a bug leaks data; mitigate with RLS), noisy neighbors (need quotas), and single blast radius.
- **DB-per-tenant:** strongest isolation (security/compliance), per-tenant backup/restore and tuning, contained blast radius; but higher cost, many databases to migrate/monitor (ops load ∝ tenants), and connection management complexity.
- Hybrid: pooled for small tenants, dedicated for large/regulated ones.

**Q9. Sample answer** (SLO 99.95% → budget 0.05%).
- **Page:** burn ≥ 14.4 over 1 h *and* over 5 min (error rate ≥ 0.72%): 2% of the budget per hour.
- **Page:** burn ≥ 6 over 6 h and 30 min (error ≥ 0.3%): 5% of the budget.
- **Ticket:** burn ≥ 1 over 3 days and 6 h: slow leak.
- Short window: ensures the problem is *still happening now*, so the alert resets quickly once fixed and doesn't page for a resolved incident.

**Q10. Sample answer.**
1. **CDN + high cache hit rates** to cut origin egress: trade-off: cache invalidation complexity and CDN costs/lock-in.
2. **Storage tiering/lifecycle** (move cold media to infrequent-access/archive; erasure coding): trade-off: slower/more expensive retrieval for cold data.
3. **Right-size media (adaptive bitrates/modern codecs, image resizing/WebP/AVIF)**: trade-off: encoding compute cost and pipeline complexity.
- Also: keep traffic within an AZ/region to avoid cross-AZ charges (trade-off: less cross-AZ balancing/resilience), reserved/spot compute (commitment or interruption risk).
