# Unit 12 Quiz: Observability, Operations & Security

Answers in `quiz-answers.md`.

**Q1 (MCQ, easy).** Which telemetry type best shows where latency accumulates for a single request across five services?
A. Metrics
B. Distributed traces
C. Unstructured logs
D. Health checks

**Q2 (MCQ, easy).** What is the main risk of using `user_id` as a metric label?
A. Metrics can't have labels.
B. High cardinality creates an enormous number of time series, hurting cost and performance.
C. It leaks passwords.
D. It breaks tracing.

**Q3 (application, medium).** An SLO is 99.9% over 30 days. Errors are running at 0.72% for the last hour. What's the burn rate, and should this page someone under the 14.4x rule?
Show your calculation.

**Q4 (MCQ, medium).** Which is the best deployment strategy for detecting a subtle performance regression using real production traffic with minimal blast radius?
A. Big-bang release to all users at 3 a.m.
B. Canary release with automated comparison of SLIs to baseline and automatic rollback.
C. Blue-green with an instant 100% switch.
D. Rolling update without health checks.

**Q5 (short, medium).** List the steps to safely rename a database column with zero downtime.

**Q6 (MCQ, medium).** Which statement about JWTs is correct?
A. JWT payloads are encrypted by default.
B. JWTs allow stateless verification but are hard to revoke before expiry, so use short lifetimes plus refresh tokens.
C. JWTs eliminate the need to validate signatures.
D. A JWT's `alg` header can be trusted from the client without pinning.

**Q7 (MCQ, hard).** A user changes `/orders/123` to `/orders/124` and sees someone else's order. The API gateway validates tokens. What's the vulnerability and fix?
A. SQL injection; use parameterized queries.
B. Broken object-level authorization (IDOR); check ownership/tenant for every object access in the service (plus RLS as a backstop).
C. CSRF; add tokens.
D. XSS; escape output.

**Q8 (short, hard).** Compare shared-schema (tenant_id column) multi-tenancy with database-per-tenant on isolation, cost, noisy-neighbor risk, and operational load.

**Q9 (application, hard).** Design three alerts (page/ticket) for an API with a 99.95% SLO using multi-window burn rates. Give thresholds and the rationale for the short confirmation window.

**Q10 (short, hard).** Give three architectural levers that reduce cloud cost significantly for a video/image-heavy product, and one trade-off for each.
