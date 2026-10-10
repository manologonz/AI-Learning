# Unit 10 Exam: Reliability & Resilience Engineering

Answers in `exam-answers.md`.

---

### Scenario (Q1-Q6): StreamBox

StreamBox is a video streaming service in one region (3 AZs). Stack: API gateway → `browse`, `playback`, `recommendations`, `billing` services → Postgres (primary + replicas) + Redis + object storage/CDN. Peak 50k API QPS. The SLO for playback start is 99.95% success within 2 s.

Incident timeline (Friday 20:00 peak): (1) `recommendations` starts responding in ~8 s instead of ~50 ms due to a bad deploy; (2) `browse` and `playback` (which call recommendations synchronously, default 30 s HTTP timeout, shared pool of 300 threads per instance) start timing out; (3) clients and gateway retry aggressively (3 attempts each at gateway and client); (4) DB CPU spikes because retries reach the DB; (5) after rolling back the deploy at 20:25, the outage lasts until 21:10 (45 more minutes) until traffic is manually blocked and slowly re-admitted.

**Q1 (analyze, medium).** Identify at least five design flaws visible in the incident and map each to a resilience technique from this unit.

**Q2 (apply, hard).** Redesign the call from `playback` to `recommendations`: give timeout values, bulkhead sizing (show your Little's Law math), circuit breaker parameters, and the fallback. Assume playback handles 500 QPS per instance and recommendations normally takes 50 ms.

**Q3 (analyze, hard).** Explain why the outage persisted 45 minutes after the rollback (metastable failure). Describe the sustaining loop and give the recovery procedure and the design changes preventing it.

**Q4 (apply, hard).** Compute the worst-case retry amplification from the client + gateway (3 attempts each) and propose a specific retry policy (where, how many, budget, backoff formula).

**Q5 (evaluate, hard).** Design load shedding for the playback API: what do you shed first, how do you detect overload, and what do you return? Justify with the SLO.

**Q6 (evaluate, hard).** The business asks: "What if the whole region goes down?" Present three DR options with RTO/RPO/cost and recommend one, assuming the business tolerates 30 minutes of downtime and 1 minute of lost data for billing, but only degraded browse for the rest.

---

### Independent items

**Q7 (MCQ, medium).** A system has three serial components with availabilities 99.9%, 99.95%, and 99.99%. Approximately what is the overall availability?
A. 99.9%
B. 99.84%
C. 99.99%
D. 99.95%

**Q8 (MCQ, hard).** Which statement about deadline propagation is correct?
A. It sets the same timeout at every hop regardless of elapsed time.
B. Each hop receives the remaining time budget, so downstream services can abandon work that can no longer meet the caller's deadline.
C. It's only useful for batch jobs.
D. It replaces circuit breakers.

**Q9 (short, hard).** Explain static stability with one example (e.g., what should an autoscaled service do if the control plane API is unavailable?).

**Q10 (create, hard).** Write a game-day plan (title, hypothesis, injected failure, blast radius, metrics to watch, abort criteria, success criteria) to test what happens when one AZ is lost for StreamBox.
