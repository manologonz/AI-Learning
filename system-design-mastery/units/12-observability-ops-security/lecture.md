# Unit 12: Observability, Operations & Security

## Objectives
By the end of this unit you can:
- Instrument a service with metrics, logs, and traces and explain what each answers.
- Design SLO-based alerts using error budgets and burn rates.
- Choose a deployment strategy (rolling, blue-green, canary) and safe database migration approach.
- Design authentication/authorization boundaries (OAuth2/OIDC, JWT, service-to-service auth) and manage secrets.
- Reason about multi-tenancy isolation and the cost of an architecture.

## Prerequisites
Units 1, 10, 11.

## Core content

### 1. Observability: understanding a system from the outside
**Monitoring** tells you *whether* something is wrong (known failure modes); **observability** lets you ask *why* about unknown ones. Three pillars (plus events/profiles):

- **Metrics**: numeric time series, cheap to store and aggregate; great for dashboards, alerting, and trends. Use the **RED** method for request-driven services (**R**ate, **E**rrors, **D**uration) and **USE** for resources (**U**tilization, **S**aturation, **E**rrors). Histograms (not averages) for latency. **Beware high cardinality labels** (user_id as a metric label explodes storage/cost).
- **Logs**: discrete events with context; **structured** (JSON) with fields (`request_id`, `user_id`, `service`, `level`, `latency_ms`), consistent levels, and sampling for high volume. Great for debugging specifics; expensive at scale. Never log secrets/PII.
- **Traces**: follow one request across services as a tree of **spans** (start, duration, attributes) sharing a **trace ID** propagated in headers (W3C `traceparent`). Reveal where latency lives and which dependency failed. **OpenTelemetry** is the vendor-neutral standard for instrumentation. Sampling: head-based (decide at start, cheap) vs tail-based (keep slow/error traces, better but needs buffering).
- **Correlation IDs** link logs, traces, and metrics exemplars; attach the trace ID to log lines.

```ts
// OpenTelemetry-style span + structured log (conceptual)
import { trace } from "@opentelemetry/api";
const tracer = trace.getTracer("orders");

async function placeOrder(req) {
  return tracer.startActiveSpan("placeOrder", async (span) => {
    span.setAttributes({ "order.items": req.items.length, "tenant.id": req.tenant });
    try {
      const res = await charge(req);                       // child spans created by instrumented HTTP client
      logger.info({ msg: "order placed", orderId: res.id, traceId: span.spanContext().traceId });
      return res;
    } catch (e) { span.recordException(e); span.setStatus({ code: 2 }); throw e; }
    finally { span.end(); }
  });
}
```

### 2. Alerting on symptoms with SLOs
Alert on **user-visible symptoms**, not causes (CPU 90% may be fine; checkout failing is not). Every alert should be **actionable and urgent**; otherwise it's noise (alert fatigue → missed real incidents).

**Burn-rate alerting** (from Unit 1's error budget): *burn rate* = observed error rate ÷ allowed error rate (1 − SLO). Burn rate 1 = you exactly spend the whole budget over the SLO window. Multi-window, multi-burn-rate alerts (Google SRE Workbook):
- **Page** when burn rate ≥ 14.4 over 1 h (and 5 min confirm): at that rate, 2% of a 30-day budget burns in an hour.
- **Page** at ≥ 6 over 6 h (5% of budget).
- **Ticket** at ≥ 1 over 3 days.
The short confirm window prevents alerts that continue after the problem is fixed.

Example: SLO 99.9% (budget 0.1%). A 1.44% error rate for one hour is a 14.4x burn → page. A 0.05% error rate is a 0.5x burn → no page.

Also: dashboards per service (golden signals: latency, traffic, errors, saturation), runbooks linked in alerts, and **on-call practices** (rotation, escalation, blameless postmortems).

### 3. Deployment and change safety
Most incidents follow a change. Strategies:
- **Rolling**: replace instances gradually; simple, but mixed versions run simultaneously (requires backward compatibility) and rollback is another rollout.
- **Blue-green**: two full environments; switch traffic atomically; instant rollback; costs double capacity and data/state must be compatible.
- **Canary**: send 1-5% traffic to the new version, compare SLIs against baseline (automated analysis), then progressively increase; automatic rollback on regressions. Best for large systems.
- **Feature flags**: decouple *deploy* from *release*; roll out by user/tenant %, kill-switch instantly. Debt: clean up old flags.
- **Progressive delivery across regions/cells**: never deploy to all regions at once.

**Database migrations (expand → migrate → contract)**: never make a change the old code can't tolerate. Add nullable column (expand) → deploy code writing both → backfill in small batches → switch reads → deploy code that stops using old column → drop it (contract). Avoid long locks (`ALTER TABLE` on big tables, adding an index without `CONCURRENTLY`), and make migrations reversible or forward-fixable.

**Config changes are deploys**: validate, roll out gradually, and version them.

### 4. Security architecture essentials
**Principles**: defense in depth, least privilege, secure by default, zero trust (never trust a request because it's "inside" the network), and minimize the attack surface.

**Authentication (who are you?) and authorization (what may you do?)**
- **OAuth 2.0** = delegated *authorization* (access tokens with scopes). **OpenID Connect (OIDC)** adds *authentication* (ID token) on top. Use an identity provider (Auth0, Cognito, Keycloak, Okta) rather than building your own. Use **Authorization Code + PKCE** for browsers/mobile; **client credentials** for machine-to-machine.
- **Sessions vs JWTs**: server-side sessions (opaque ID in cookie → store lookup): easy revocation, requires shared store. **JWT** (signed, self-contained claims): stateless verification (no lookup), but **hard to revoke before expiry**: use short lifetimes (5-15 min) + refresh tokens (rotated, revocable). Never put secrets in a JWT (it's signed, not encrypted); validate signature, `exp`, `aud`, and `iss`; pin algorithms.
- **Authorization models**: RBAC (roles), ABAC (attributes/policies), ReBAC (relationships, e.g., Google Zanzibar-style: "can Alice view document X because she's in group G that owns folder F"). Centralize policy decisions (policy engine, e.g., OPA/Cedar) and **enforce at every service**, not only the gateway.
- **Service-to-service**: **mTLS** (mutual TLS with short-lived certs from a mesh/SPIFFE) or signed service tokens; propagate the *end-user identity* (token exchange) so downstream services authorize on behalf of the user (avoid the "confused deputy").

**Data protection**: TLS in transit everywhere; encryption at rest (KMS-managed keys, envelope encryption); field-level encryption for highly sensitive data; hashing passwords with bcrypt/argon2 (never reversible encryption); **secrets** in a manager (Vault, AWS Secrets Manager), injected at runtime, rotated, never in code/images/logs.

**Common attack surface**: injection (parameterized queries), SSRF, broken access control (**IDOR**: verify object ownership on every access; the #1 OWASP risk), insecure deserialization, dependency/supply-chain vulnerabilities (pin, scan, SBOM), DDoS (edge/CDN/WAF, rate limits), abuse (bots, credential stuffing: rate limit + MFA).

**Compliance & privacy**: data classification, residency (GDPR, data localization → regional partitioning), retention and deletion (right to erasure vs immutable logs/backups; design deletion pipelines), audit logs (append-only), PCI-DSS scope reduction (tokenize card data with the PSP so your systems never touch raw PANs).

### 5. Multi-tenancy
SaaS choices along an **isolation vs cost** spectrum:
- **Shared everything** (tenant_id column, row-level security): cheapest, highest density; risk of cross-tenant leaks (enforce with RLS + tests) and noisy neighbors (per-tenant quotas/rate limits).
- **Schema-per-tenant**: better logical isolation and per-tenant migrations/backup; operationally heavier at thousands of tenants.
- **Database-per-tenant / cell-based**: strongest isolation and blast-radius control; higher cost/ops; suits enterprise/regulated customers.
- **Cell-based architecture**: partition the entire stack into independent **cells** (each serving a subset of tenants) so failures and bad deploys are contained (blast radius ↓), with a thin routing layer above.
Hybrid: pooled for small tenants, dedicated for large/regulated ones (the "whale" solution from Unit 6).

### 6. Cost as an architectural dimension
Cost is an NFR (Unit 1). Know the big levers: **egress bandwidth** (CDN, cross-AZ/region traffic charges), storage tiering (hot/warm/cold, lifecycle policies), compute right-sizing/spot/reserved capacity, managed service premiums vs ops labor, over-provisioning for peaks (autoscaling), log/metric retention and cardinality, and cross-AZ chatty microservices. Practice **FinOps**: tag resources by team/service, show costs per unit of business value (cost per order), set budgets/alerts, and include cost in design docs ("this design costs about $X/month at Y load").

## Worked examples

### Straightforward: from average latency alert to SLO alert
*Alert today: "CPU > 80% for 5 min", fires nightly during batch jobs and is ignored.* Replace with an SLO: 99.9% of playback starts succeed within 2 s. Alerts: burn rate ≥ 14.4 (1 h + 5 min) pages; ≥ 6 (6 h) pages; ≥ 1 (3 d) tickets. CPU becomes a dashboard/saturation signal, not a page.

### Typical: safe column rename
*Rename `orders.cust_id` to `customer_id` with zero downtime.* (1) Add `customer_id` (nullable). (2) Deploy code that writes both columns and reads `cust_id`. (3) Backfill in batches of 10k rows off-peak. (4) Deploy code reading `customer_id` (canary first). (5) Deploy code that stops writing `cust_id`. (6) Drop `cust_id` after a soak period. Each step is independently rollbackable; a single `ALTER TABLE RENAME` would break old instances during the rollout.

### Edge case: IDOR behind a "secure" gateway
*The gateway validates the JWT, and `GET /invoices/{id}` in the billing service trusts that any authenticated request is allowed.* A user increments `id` and reads other tenants' invoices: broken object-level authorization. The gateway authenticated but couldn't know ownership. Fix: every service authorizes each object access using the *end-user identity + tenant* (`WHERE id=$1 AND tenant_id=$2`, or a policy check), Postgres row-level security as a backstop, non-guessable IDs (UUIDs; not a defense by themselves), and automated authorization tests per endpoint. Also: internal services must not blindly trust network location (zero trust).

## Common misconceptions

- **"Observability = dashboards."** Dashboards answer known questions; observability supports exploring unknown ones with high-cardinality, correlated data.
- **"Alert on everything to be safe."** Excess alerts train people to ignore them; page only on user-impacting, urgent, actionable conditions.
- **"Average latency is fine."** Use histograms/percentiles (Unit 1).
- **"JWTs are more secure than sessions."** They're stateless, not more secure; revocation is the trade-off.
- **"The gateway handles security."** Authorization must be enforced where the data lives.
- **"Encryption at rest protects against everything."** It protects lost disks/snapshots; an authorized-but-compromised application still reads plaintext.
- **"Blue-green means zero risk."** Shared databases and stateful changes still bite; canaries catch issues with real traffic.
- **"Logs are free."** Volume and retention drive real cost; sample and structure.

## Check your understanding

1. SLO 99.95%: what error rate is a burn rate of 10? *(Allowed error 0.05%; 10x → 0.5%.)*
2. Why is `user_id` a bad metric label? *(High cardinality explodes time series count; put it in logs/traces.)*
3. Name two ways to safely roll out a risky change. *(Canary with automatic analysis; feature flag with gradual ramp/kill switch; region/cell-by-cell.)*
4. Why can't a gateway alone prevent IDOR? *(It can't know object ownership; each service must authorize per object.)*

## Summary
- Metrics (what/how much), logs (details), traces (where) work together via correlation IDs; use OpenTelemetry and watch cardinality/cost.
- Alert on SLO burn rate for user-visible symptoms; make alerts actionable with runbooks.
- Deploy progressively (canary/flags/cells) and use expand→migrate→contract for schemas.
- Authenticate at the edge, authorize at every service, propagate user identity, use mTLS between services, manage secrets centrally.
- Choose tenant isolation by risk and cost; treat cost as a first-class design dimension.
