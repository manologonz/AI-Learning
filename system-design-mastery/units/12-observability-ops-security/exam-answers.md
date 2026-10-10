# Unit 12 Exam: Answer Key

**Q1. Sample mapping.**
1. Slow PDFs unnoticed for 2 days → missing **SLO-based, symptom alerting** (latency/success SLI for PDF generation with burn-rate alerts) and per-journey dashboards; also synthetic checks.
2. Doubled p99 invisible → dashboards used **average** latency; need **histograms/percentiles**, canary analysis comparing p99/error SLIs, and automatic rollback.
3. Cross-tenant invoice → **broken object-level authorization**; fix per-service ownership checks (tenant + object), RLS backstop, authorization tests.
4. Bulk import noisy neighbor → **per-tenant quotas/rate limits, bulkheads, queueing/priority lanes**, isolated worker pools; possibly dedicated cell for whales.
5. Committed DB password → **secrets management** (Vault/Secrets Manager), pre-commit/CI secret scanning, rotation, least-privilege DB credentials, short-lived credentials.
6. Cost spike from verbose logs → **log volume governance**: levels, sampling, retention tiers, structured/dedupe, cost dashboards per service.

**Q2.**
- SLI: good = `POST /invoices/{id}/pdf` requests completing with success and latency ≤ 5 s (or async job completed ≤ 30 s), valid = all requests excluding client errors (4xx from invalid input) and synthetic traffic. Measured at the service edge/queue.
- SLO: **99.5% over 30 days.**
- Budget: 0.5% × 20M = **100,000 bad requests** per month.
- Alerts: page at burn ≥ 14.4 over 1 h + 5 min (error ≥ 7.2%); page at burn ≥ 6 over 6 h + 30 min (≥ 3%); ticket at burn ≥ 1 over 3 days (≥ 0.5%).
- Also monitor saturation (queue depth/age) as a dashboard signal.

**Q3.**
- **Checks at every service**, per object: every query includes `tenant_id` (from the verified token) and object ownership: `WHERE id=$1 AND tenant_id=$2` (return 404 not 403 to avoid leaking existence).
- **Identity propagation:** the gateway validates the JWT (issuer, audience, expiry, signature), and forwards the end-user identity + tenant to services via signed internal token/token exchange (or service mTLS + user context in a signed header); services never accept `tenant_id` from request bodies/URLs without checking against the token.
- **Defense in depth (DB):** Postgres **row-level security** with `tenant_id = current_setting('app.tenant_id')`, set per transaction; least-privilege DB roles (no BYPASSRLS for app roles).
- **Testing:** automated authorization tests per endpoint (user A cannot read B's object; matrix tests), fuzzing IDs, integration tests in CI, periodic pen-tests; monitoring for suspicious access patterns (many 404/403s).

**Q4.**
- Options: (1) shared everything with quotas: cheap but risky; (2) schema-per-tenant: better isolation, migration overhead; (3) database/cell per tenant: strongest, costly.
- Hybrid: **pooled multi-tenant cells** for small/medium tenants with **per-tenant rate limits, concurrency caps, fair queueing** for imports (separate worker pools/priority queues, bulk imports run as throttled background jobs), and **dedicated cells** (own DB + compute) for the top ~1% of tenants (10M invoices) or regulated ones. Trigger to move: sustained usage > X% of a cell's capacity, latency impact on others, or contractual isolation requirements. Routing layer maps tenant → cell; migrations by tenant with dual-write/CDC and cutover.

**Q5.**
- **First hour:** revoke/rotate the credential immediately (assume compromised, even if the repo is private); check DB audit logs for unauthorized use; remove from git history (BFG/filter-repo) and force-rotate anything derived; notify security; assess blast radius.
- **Systemic:** secret manager with runtime injection (no secrets in repos/env files), pre-commit hooks and CI secret scanning (gitleaks, GitHub push protection), short-lived DB credentials (IAM auth/Vault dynamic secrets), least-privilege roles, regular rotation, and developer training/runbook.

**Q6. Sample strategy.**
- Keep at 100%: errors/warnings, audit/security events, traces for errors and slow requests (tail-based sampling), metrics (aggregated; drop high-cardinality labels).
- Sample: successful request logs at 1-10% (sample by trace ID so entire traces are consistent); debug logs off by default, enabled per service/tenant temporarily via dynamic log levels.
- Tiering/retention: hot 7 days (searchable), warm 30 days (cheaper/slower), cold 1 year in object storage (audit/compliance); drop verbose logs after 3 days.
- Structure: JSON with fixed fields, no large payload dumps, redact PII/secrets, dedupe repeated errors (rate-limit identical messages).
- Governance: per-service log volume budgets/dashboards, cost by team. Expected ~50% savings while preserving debuggability through traces + error logs.

**Q7: B.**
- Authorization Code + PKCE is the recommended flow for public clients; A (implicit flow) is deprecated; C: secrets in mobile apps are extractable; D is insecure.

**Q8: B.**
- With ~300 requests, a 0.01 percentage-point difference isn't statistically meaningful, and a +8% p99 could be noise or real. Extend the canary/raise the share until enough traffic, use statistical tests and predefined thresholds. A: premature; C: overreacts; D: latency regressions matter (SLO).

**Q9. Sample answer.**
- Confused deputy: a service with broad privileges is tricked into performing an action for a caller who lacks permission, because the service authorizes based on *its own* identity/authority rather than the caller's (e.g., the reporting service can read all invoices; a user asks it for another tenant's report).
- Mitigation: propagate the **end-user identity/tenant** (token exchange yielding a downstream token scoped to the user and the audience) so downstream services authorize on behalf of the user (least privilege), and services authorize per request.

**Q10. Sample answer.**
```
        Clients → Edge/CDN/WAF → Global router (tenant → cell map, control plane)
                                   │
        ┌──────────────┬───────────┴───────────┬──────────────┐
      Cell A          Cell B                 Cell C         Cell D (dedicated: tenant X)
   (25 services,   (25 services,          (…)             (own DB, own quotas)
   own Postgres,   own Postgres,
   own queues)     own queues)
   Shared: identity provider, billing, observability pipeline, CI/CD, secrets (control plane)
```
- Routing: thin, stateless router maps tenant_id → cell (lookup cached, highly available).
- Deployment: progressive across cells (canary cell first, then 1 cell at a time, automatic rollback), never all cells at once.
- Failure containment: a bad deploy, poison message, or noisy tenant affects only one cell (~1/N of tenants); cells are independent (no cross-cell synchronous calls); the control plane is static-stable (cell keeps serving if control plane is down).
- Tenant moves between cells via data migration tooling.
