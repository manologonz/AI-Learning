# Unit 14: Architecture in the Real World

## Objectives
By the end of this unit you can:
- Write an Architecture Decision Record (ADR) and a design doc that others can review and act on.
- Critique an existing architecture and separate symptoms from structural causes.
- Plan a zero-downtime migration with verification, rollback, and staged rollout.
- Run a blameless postmortem and turn incidents into durable improvements.
- Lead a design review: build alignment, surface risks, and make decisions under uncertainty.
- Make and defend build-vs-buy and technical-debt decisions.

## Prerequisites
Units 1-13. This unit is about the *human and organizational* half of senior system design.

## Core content

### 1. Seniority is mostly influence and judgment
By now you can produce a technically sound design. At senior level, the design is only part of the job. You must also:
- **Frame the problem** so the team solves the right one.
- **Get alignment** among engineers, product, ops, and security.
- **Decide under uncertainty** with reversible/irreversible categories (Bezos's *one-way vs two-way doors*): move fast on two-way doors; invest analysis in one-way doors (data models, public APIs, shard keys, vendor lock-in).
- **Communicate in writing** so decisions outlive meetings and people.
- **Reduce risk incrementally**: small steps, feature flags, canaries, kill switches.
- **Teach and multiply** others: reviews, mentoring, shared patterns.

### 2. Design documents that work
A good design doc is a tool for **thinking and for review**, not paperwork. Keep it short (3-8 pages) and structured:
1. **Context & problem**: why now, who's affected, links to data.
2. **Goals and non-goals**: explicit; non-goals prevent scope creep.
3. **Requirements** (functional + NFR/SLOs; see Unit 1) and **key numbers** (estimates).
4. **Proposed design**: diagram, data model, APIs, flows; enough detail to implement and to find flaws.
5. **Alternatives considered**: with honest pros/cons and why rejected. This is the most important section for reviewers: it shows the decision wasn't obvious.
6. **Trade-offs and risks**: what could go wrong, mitigations, reversal conditions.
7. **Rollout/migration plan**: phases, flags, rollback, verification.
8. **Operational plan**: observability, SLOs/alerts, runbooks, on-call impact, capacity, cost.
9. **Security/privacy/compliance** considerations.
10. **Open questions** and timeline/staffing.
Tips: lead with a summary and the decision needed; use diagrams (C4 style: context → containers → components); keep prose plain; state assumptions and numbers; invite disagreement early (circulate a draft before the meeting).

### 3. Architecture Decision Records (ADRs)
A short, immutable-once-accepted note capturing **one decision**: **Title, Status** (proposed/accepted/superseded), **Context** (forces at play), **Decision**, **Consequences** (good and bad), **Alternatives**. Stored in the repo (`/docs/adr/0007-use-outbox-for-events.md`). They answer the future question "*why on earth did we do it this way?*" and stop teams from re-litigating decisions without new information. When circumstances change, write a *new* ADR that supersedes the old one; don't rewrite history.

```md
# ADR-0012: Adopt transactional outbox for order events
Status: Accepted (2026-03-04)
Context: Dual writes to Postgres and Kafka lost events during broker outages (INC-231). Need atomic state change + event.
Decision: Write events to an `outbox` table in the same transaction; publish via Debezium.
Consequences: + no lost events, at-least-once (consumers must be idempotent); − extra table/CDC ops burden; ~200 ms added publish latency.
Alternatives: 2PC (rejected: broker unsupported, availability), publish-after-commit (rejected: loses events on crash).
Revisit if: event volume > 50k/s per DB or CDC lag exceeds SLO.
```

### 4. Critiquing an existing architecture
Approach (in order):
1. **Understand goals and constraints**: what is the system *for*, what are the SLOs, team size, history? Many "bad" designs were sensible when written.
2. **Measure before judging**: traffic, latency percentiles, error budgets, cost, deploy frequency, incident history, change failure rate, onboarding time. Use the **DORA metrics** (deployment frequency, lead time, change failure rate, MTTR) as team-level health.
3. **Find the actual bottleneck/pain**: constraint theory: improving non-bottlenecks is waste.
4. **Look for structural smells**: shared databases, distributed monoliths, unbounded queues, missing timeouts, single points of failure, no idempotency, manual deploys, cache-dependent databases, unclear ownership, unowned "god" services, config drift.
5. **Separate symptoms from causes** (5 whys): "checkout is slow" → DB CPU → N+1 queries → ORM lazy loading in a shared library → no query budget in review.
6. **Propose the smallest change that removes the constraint**, with a measured target; avoid rewrites.
7. **Rank by risk × impact × effort**; present as a roadmap.

### 5. Evolving systems: migrations and technical debt
- **Zero-downtime migration checklist**: (1) dual-run/shadow to validate; (2) expand → migrate → contract for schemas and contracts; (3) backfill in throttled batches with progress tracking and idempotence; (4) verify with checksums/row counts and sampled comparisons; (5) cut over gradually (1% → 100%) behind flags; (6) keep rollback available until confidence; (7) decommission the old path (delete!); (8) communicate timelines to dependents.
- **Data migrations** are the riskiest: prefer reversible steps; keep old and new in sync via CDC/dual writes while verifying; plan the *reverse* sync if rollback after cutover must preserve new writes.
- **Technical debt**: deliberate/prudent debt (ship now, pay later with a plan) is a tool; reckless/inadvertent debt is risk. Track it visibly with an owner and *impact* (e.g., "adds 2 days to every release", "causes 1 incident/quarter"), and pay it down alongside feature work, tied to outcomes. Rewrites are a last resort: the strangler pattern (Unit 11) beats rewrite.
- **Build vs buy vs adopt open source**: build only where it's a *differentiator* or no fit exists; buy/managed services for commodities (auth, email, payments, monitoring). Evaluate **total cost of ownership** (including on-call, upgrades, expertise, integration), vendor lock-in and exit cost, SLAs vs your SLOs, data-residency/compliance, and maturity. A managed service's *outage* becomes yours: design for its failure (Unit 10).

### 6. Incidents and postmortems
- **Incident response**: detect → triage → mitigate first (rollback, failover, shed load), diagnose second; clear roles (incident commander, comms, ops); one channel; timeline notes; customer comms.
- **Blameless postmortem**: assume people acted reasonably with the info they had; focus on *systems, incentives, and missing safeguards*. Structure: summary and impact (SLO/budget), timeline, root causes and contributing factors (not "human error": ask why the system allowed it), what went well, where we got lucky, **action items with owners and dates** (prioritized: prevent, detect, mitigate), and follow-through tracking. Share widely; treat as learning.
- Metrics: MTTD, MTTR, recurrence rate, % action items closed.

### 7. Leading a design review
- **Before:** circulate the doc 2-3 days ahead; assign reviewers with specific lenses (security, SRE, data, product); state the decision you need.
- **During:** timebox; start with the problem and constraints, not the solution; encourage "what would make this fail?" questions; separate *blocking* concerns from suggestions; record decisions and open questions live.
- **After:** update the doc/ADR, follow up on action items, and communicate the outcome.
- **Style:** critique the design, not the person; seek disconfirming evidence; use data; "disagree and commit" once decided; change your mind visibly when evidence changes.
- **Common failure modes:** bikeshedding on details, HiPPO (highest-paid person's opinion) decisions, analysis paralysis on two-way doors, groupthink, designs presented as fait accompli.

### 8. Putting the senior toolkit together
A senior system designer's habits:
1. **Requirements and numbers first**, explicit SLOs.
2. **Simplest design that meets the numbers**; each component justified.
3. **Name the hard part; deep-dive there.**
4. **Assume failure**; design for degradation, recovery, and observability.
5. **State trade-offs with reversal conditions**; record them (ADRs).
6. **Evolve incrementally**; every big change is a series of reversible steps.
7. **Make the team better**: write, review, mentor, standardize.

## Worked examples

### Straightforward: writing an ADR
*Decision: use Postgres row-level security for tenant isolation.* Context: two cross-tenant leak bugs in a year. Decision: RLS with `tenant_id` set per transaction; app roles lack BYPASSRLS. Consequences: + defense in depth; − per-query overhead ~2-5%, careful connection-pool session handling (PgBouncer transaction mode needs `SET LOCAL`). Alternatives: schema-per-tenant (ops burden), app-only checks (proven insufficient). Revisit if: > 5k tenants or noisy-neighbor issues require cell isolation.

### Typical: zero-downtime migration from a monolith DB table to a service
*Move `notifications` tables into the new Notification service.* (1) Service reads/writes its own DB; monolith continues writing to the old tables. (2) CDC streams old-table changes to the service (backfill + ongoing). (3) Shadow reads compare results. (4) Flip *reads* to the service for 1% → 100% of users. (5) Flip *writes* (monolith calls the service API; old tables receive reverse-sync via CDC as a rollback safety net). (6) After 2 weeks with no diffs, stop the reverse sync and drop the old tables. Each step has a rollback and a success metric (diff rate, latency, error rate).

### Edge case: the "obviously wrong" architecture that shouldn't be changed
*A reviewer wants to replace a 6-year-old cron + shell-script batch pipeline with a streaming platform. The pipeline runs nightly, has 99.99% success over 3 years, costs $200/month, and the only consumer is a morning report.* The elegant redesign has high cost (months of work), new operational burden, and no user-visible benefit. Right call: leave it, add monitoring and documentation, and revisit only if requirements change (e.g., intra-day reports, data volume ×100). **Senior judgment includes knowing what not to build.**

## Common misconceptions
- **"Design docs are bureaucracy."** They're cheap thinking; the cost of a bad design found in code is 10-100x more.
- **"ADRs must be long."** One page, one decision.
- **"Blameless means no accountability."** It means accountability for fixing systems, not scapegoating; owners still own action items.
- **"A rewrite will fix our debt."** Rewrites reproduce old requirements with new bugs; migrate incrementally.
- **"Buy is always cheaper than build."** Integration, lock-in, and operational realities can dominate; compute TCO.
- **"Consensus in reviews means quality."** Groupthink is real; explicitly assign a devil's-advocate role.
- **"Seniors have the right answers."** They ask better questions, know the failure modes, and communicate trade-offs.

## Check your understanding
1. What is the most valuable section of a design doc for reviewers? *(Alternatives considered/trade-offs: it exposes the reasoning and lets reviewers challenge it.)*
2. One-way vs two-way door: classify a shard-key choice and a feature flag default. *(Shard key = one-way (costly to change); flag default = two-way.)*
3. Why should a postmortem action item have an owner and date? *(Otherwise it never happens; learning requires follow-through.)*
4. When should you *not* replace an ugly-but-working system? *(When cost/risk exceeds user-visible benefit and requirements haven't changed.)*

## Summary
- Senior work = judgment + communication: design docs, ADRs, reviews, and postmortems make decisions durable and improve the team.
- Critique with measurements and constraints; find the true bottleneck; make the smallest change that removes it.
- Migrate incrementally with verification and rollback; never big-bang.
- Build only differentiators; evaluate TCO and failure modes of buy decisions.
- Blameless postmortems turn incidents into system improvements.
