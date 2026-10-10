# Unit 12 Exam: Observability, Operations & Security

Answers in `exam-answers.md`.

---

### Scenario (Q1-Q6): LedgerLane

LedgerLane is a B2B invoicing SaaS: 4,000 tenants (a few with 10M invoices), 25 services on Kubernetes, one shared Postgres cluster (tenant_id column), a gateway that validates JWTs, and a mobile app. Recent issues: (1) customers complained about slow invoice PDFs for two days before engineering noticed; (2) a deploy changed a query and p99 doubled, but dashboards showed only average latency, which barely moved; (3) a tenant reported seeing another tenant's invoice via a manipulated URL; (4) a large tenant's bulk import slowed everyone; (5) an engineer committed a database password to the repo; (6) cloud bill up 40% after adding verbose logging.

**Q1 (analyze, medium).** For each numbered incident, name the missing capability and the fix (be specific: telemetry type, alert, control).

**Q2 (apply, hard).** Define an SLO and SLIs for "invoice PDF generation" (choose good/valid events, target, window), compute the error budget in requests for 20M requests/month at 99.5%, and design page/ticket alerts with burn rates.

**Q3 (apply, hard).** Design the authorization model to eliminate incident 3: where checks happen, what identity is propagated between services, and how you'd test it. Include defense in depth at the database layer.

**Q4 (evaluate, hard).** Propose tenant isolation options for LedgerLane's noisy-neighbor problem (incident 4) and pick a hybrid. Include per-tenant rate limits/quotas and when a tenant moves to a dedicated cell.

**Q5 (apply, hard).** Deliver a remediation plan for the leaked secret: immediate steps in the first hour, and systemic controls so it can't recur.

**Q6 (evaluate, hard).** Propose a logging/telemetry strategy that cuts cost ~50% without losing debuggability: sampling, retention tiers, structure, cardinality, and what you'd keep at 100%.

---

### Independent items

**Q7 (MCQ, medium).** Which is the best practice for a mobile app's OAuth flow?
A. Implicit flow storing tokens in local storage.
B. Authorization Code flow with PKCE, short-lived access tokens, and refresh token rotation.
C. Hardcode a client secret in the app.
D. Send the user's password to every API call.

**Q8 (MCQ, hard).** A canary at 5% shows error rate 0.1% vs 0.09% baseline, latency p99 +8%. Traffic is low, and the sample has ~300 requests. What's the best action?
A. Promote to 100% since the difference is small.
B. Treat the result as statistically inconclusive; extend the canary or increase its share until there's enough traffic for a confident comparison (or use stricter guardrails), then decide.
C. Roll back immediately without analysis.
D. Ignore latency; only errors matter.

**Q9 (short, hard).** Explain "confused deputy" in a microservice context and how token exchange/propagation of the end-user identity mitigates it.

**Q10 (create, hard).** Draw (in text) a cell-based architecture for LedgerLane: cells, routing layer, shared components, deployment strategy across cells, and failure containment story.
