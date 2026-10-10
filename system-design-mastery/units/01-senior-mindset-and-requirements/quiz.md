# Unit 1 Quiz: The Senior Mindset & Requirements

Formative, low stakes. Answer in the chat or write your answers down, then I'll grade with feedback. (Answers are in `quiz-answers.md`; no peeking.)

**Q1 (MCQ, easy).** Which best describes a Service Level Objective (SLO)?
A. A legally binding promise to customers with financial penalties.
B. A measurement of some aspect of service behavior, such as request latency.
C. An internal target value for an SLI over a time window.
D. The maximum load a service can handle before failing.

**Q2 (short answer, easy).** A service has a 99.9% availability SLO over a 30-day window. How many minutes of downtime does the error budget allow? Show the calculation.

**Q3 (MCQ, medium).** A team defines its availability SLI as "the percentage of times our load balancer's `/healthz` probe returns 200." Users report frequent checkout failures, yet the SLI shows 99.99%. What's the most fundamental problem?
A. The SLO target is too high.
B. The SLI doesn't reflect the user-facing experience of real requests.
C. Health checks should run more frequently.
D. They should have used an SLA instead of an SLO.

**Q4 (MCQ, medium).** You're asked to design a ride-sharing app. Which item is a **non-functional** requirement?
A. Riders can request a ride to a destination.
B. Drivers can accept or decline a request.
C. Matching a rider to a nearby driver should complete within 2 seconds at p99.
D. Riders can see their trip history.

**Q5 (short answer, medium).** Explain in 2-3 sentences why an SLA is normally set looser than the internal SLO.

**Q6 (MCQ, medium).** A dashboard shows average latency of 120ms and the team declares the service healthy. Which is the strongest objection?
A. Averages can't be computed for latency.
B. Averages hide tail latency, so a meaningful slice of users may be having a terrible experience.
C. Latency should be measured in seconds, not milliseconds.
D. Only the median matters, so the average is redundant.

**Q7 (application, medium).** A page load fans out to 30 backend calls in parallel. Each call independently has a 1% chance of exceeding the latency target (i.e., p99 = target). Approximately what fraction of page loads will contain at least one slow call? Show the formula.

**Q8 (MCQ, hard).** A team says: "We'll use Kafka, Kubernetes, and 12 microservices from day one so we can scale." They currently have 200 users and a team of three. Which critique best reflects senior-level thinking?
A. Kafka is not a scalable technology.
B. Every component adds operational and cognitive cost and should be justified by a named requirement or number; start with the simplest design that meets current needs and define reversal conditions.
C. Microservices are always worse than monoliths.
D. They should have used serverless instead, since it scales automatically.

**Q9 (short answer, hard).** Write a trade-off statement (using "we choose X over Y because... holds while... revisit if...") for this decision: *choosing eventual consistency for a "like count" on posts instead of strongly consistent counts.* Include a concrete revisit condition.

**Q10 (application, hard).** For the prompt "Design a food-delivery order-tracking feature (customer sees the courier's live location on a map)", list: (a) two functional requirements, (b) three non-functional requirements with concrete-ish numbers or qualities, (c) one out-of-scope item, and (d) what you'd identify as the "hard part."
