# Unit 14 Quiz: Architecture in the Real World

Answers in `quiz-answers.md`.

**Q1 (MCQ, easy).** Which best describes an Architecture Decision Record (ADR)?
A. A long specification of every component.
B. A short, versioned note capturing one significant decision: context, decision, consequences, and alternatives.
C. A meeting agenda.
D. A monitoring dashboard.

**Q2 (MCQ, easy).** Which section of a design doc most helps reviewers evaluate the proposal?
A. The acknowledgments
B. Alternatives considered, with trade-offs
C. The company history
D. The list of technologies used

**Q3 (MCQ, medium).** "One-way door" decisions deserve more up-front analysis. Which is most clearly a one-way door?
A. Renaming an internal variable
B. Choosing a shard key for a 50 TB dataset
C. Toggling a feature flag default
D. Changing a log level

**Q4 (short, medium).** List four elements a blameless postmortem should contain.

**Q5 (MCQ, medium).** During an incident, what should be the first priority?
A. Finding the root cause
B. Mitigating customer impact (rollback, failover, load shedding)
C. Writing the postmortem
D. Deciding who is responsible

**Q6 (short, medium).** Give three signals you'd measure before criticizing an architecture.

**Q7 (application, hard).** Write an ADR (Title, Status, Context, Decision, Consequences, Alternatives, Revisit if) for: "Use Kafka for inter-service events instead of direct HTTP calls" at a company with 12 services and recurring cascading failures.

**Q8 (short, hard).** Outline a zero-downtime plan to migrate a large table to a new database, including verification and rollback.

**Q9 (MCQ, hard).** Evaluating build vs buy for an identity/authentication system in a 30-person company. Which analysis is most complete?
A. Build it: we have engineers, so it's free.
B. Buy/adopt a managed identity provider unless authentication is the product differentiator; compare TCO (integration, on-call, security expertise), vendor lock-in and exit cost, SLAs vs SLOs, compliance, and plan for provider outages.
C. Buy the most expensive option.
D. Buy nothing; use a shared password file.

**Q10 (short, hard).** A senior engineer proposes a full rewrite of a legacy system "to fix technical debt." List three questions you'd ask and an alternative approach.
