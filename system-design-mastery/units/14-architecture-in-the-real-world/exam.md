# Unit 14 Exam: Architecture in the Real World

Answers in `exam-answers.md`.

---

### Scenario (Q1-Q6): Parcelo

Parcelo, a logistics SaaS (8 years old, 60 engineers), runs a Rails-style monolith + 6 services, one shared Postgres for most data. Facts: deploys take 3 hours and are done weekly; change failure rate ~25%; MTTR ~3 hours; the "shipments" table (4 TB) is queried by 5 services directly; a nightly pipeline (cron + scripts, 5 years, reliable) feeds the finance team; incidents last quarter: 3 (all after deploys), one caused by a Kafka misconfiguration nobody owns; on-call is burning out; the CTO wants "a full move to microservices on Kubernetes with a new event-driven platform in 6 months."

**Q1 (analyze, hard).** Critique the CTO's plan and Parcelo's current architecture: list the evidence you'd gather first, three structural causes of the deploy pain, and which parts you would *not* change (with reasons).

**Q2 (apply, hard).** Write a one-page design doc outline (all headings with 1-2 lines each) for "Give the shipments data a single owner service," including goals/non-goals, alternatives, rollout, and operational plan.

**Q3 (apply, hard).** Write an ADR for "Adopt canary deployments with automated rollback" including Context (use numbers from above), Decision, Consequences, Alternatives, and revisit conditions.

**Q4 (evaluate, hard).** Draft the blameless postmortem summary (timeline in 5-6 lines, root causes vs contributing factors, 5 action items with categories: prevent/detect/mitigate and owners) for the Kafka misconfiguration incident: a producer's `acks=1` and an unreplicated topic lost 40 minutes of shipment events during a broker restart; the reconciliation job caught it 6 hours later.

**Q5 (evaluate, hard).** Prioritize a 12-month roadmap for Parcelo (quarterly) with rationale, expected metric improvements (DORA and reliability), and how you'd communicate it to leadership skeptical of "not doing microservices."

**Q6 (create, hard).** Plan and run a design review for the shipments-ownership proposal: attendees/lenses, pre-read, agenda with times, decision rule, and what you'd do if two senior engineers strongly disagree.

---

### Independent items

**Q7 (MCQ, medium).** Which is the best way to handle "technical debt" in planning?
A. Ignore it until a rewrite.
B. Track it visibly with owner and business impact, and allocate ongoing capacity to pay down items that reduce measurable pain.
C. Pay all of it before any new features.
D. Never document it.

**Q8 (MCQ, hard).** In a design review, the most senior person immediately states their preferred solution and the room goes quiet. What is the best facilitation move?
A. Adopt it to save time.
B. Ask each reviewer to write their assessment/risks independently first (or start with the most junior), then discuss alternatives against the stated criteria.
C. Postpone the review indefinitely.
D. Vote by seniority.

**Q9 (short, hard).** Explain "disagree and commit" and when it's appropriate vs harmful.

**Q10 (create, hard).** For the nightly finance pipeline: write a short argument for *not* rewriting it, the monitoring you'd add, and the two conditions under which you'd revisit.
