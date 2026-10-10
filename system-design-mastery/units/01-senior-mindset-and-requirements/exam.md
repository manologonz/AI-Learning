# Unit 1 Exam: The Senior Mindset & Requirements

Summative and harder than the quiz: these questions combine ideas from across the unit. Take your time and show your reasoning. Answers in `exam-answers.md`.

---

### Scenario (for Q1-Q6)

**PayLite** is a startup building a peer-to-peer payments app. Today: 50k users, 15 engineers, one Node/TypeScript monolith on Postgres. The CEO wants to reach 5M users within 18 months. Users can add a bank card, send money to a friend by phone number, and see transaction history. An outage during peak hours (evenings, paydays) generates a flood of support tickets and press attention.

---

**Q1 (analyze, medium).** Classify each as FR, NFR, or out of scope (for now), and briefly justify:
(a) "Send money to a friend by phone number."
(b) "A transfer must never be applied twice, even if the user's client retries."
(c) "Transaction history loads within 400ms at p95."
(d) "Cryptocurrency support."

**Q2 (apply, medium).** PayLite decides its critical journey is "send money." Propose an SLI (define both the good and valid events) and an SLO with a window for this journey, and justify your target number in 2-3 sentences. Then compute the resulting error budget in minutes for a 30-day window.

**Q3 (analyze, hard).** After launch, the SLI for "send money" is 99.93% against a 99.95% SLO, but the engineering lead says, "We're fine, that's basically 99.95%." Give two distinct arguments why she is wrong (or why the situation needs action), referencing the error budget concept.

**Q4 (evaluate, hard).** An architect proposes: "To hit 5M users, let's immediately rewrite into 20 microservices with a sharded database." Write a critique in 4-6 sentences that (i) names the risks/costs, (ii) proposes a staged alternative, and (iii) specifies at least two measurable triggers that would justify sharding later.

**Q5 (evaluate, hard).** PayLite must choose between two designs for the balance shown on the home screen:
- *Design A:* read the balance from the primary DB every time (strongly consistent, ~40ms, DB load grows with users).
- *Design B:* read a cached balance (may lag by up to 5s).
Write a trade-off statement for each of the two contexts: (i) the balance shown on the home screen; (ii) the balance check used to authorize a transfer. Do you make the same choice in both? Why or why not?

**Q6 (analyze, hard).** Identify the "hard part" of PayLite's system, defend it in 2-3 sentences, and say which *two* non-functional requirements would you push the product team to quantify first, and why.

---

### Independent items

**Q7 (MCQ, medium).** A vendor's SLA promises 99.9% uptime, with service credits if breached. Your product's SLO is 99.95%, and it depends serially on this vendor (your service is down whenever the vendor is down), plus your own components each with 99.99% availability. Roughly what is the best-case availability of your product, and what does that imply?
A. 99.99%, the strongest component dominates.
B. About 99.9%, since serial dependencies multiply; you can't promise 99.95% relying solely on this vendor without redundancy or graceful degradation.
C. 99.95%, because the SLO is what you set, not what dependencies deliver.
D. 100%, since the vendor pays credits when they fail.

**Q8 (MCQ, hard).** A team's p99 latency SLO is 500ms, and it's currently met. A new feature makes each page load fan out to 40 backend calls, each with an independent p99 of 500ms. What is the best description of the likely result?
A. Page-level p99 stays 500ms because each call meets its SLO.
B. Roughly a third of page loads will contain at least one call over 500ms, so page-level tail latency is much worse than per-call SLOs suggest; mitigations include fewer calls, hedged requests, or timeouts with degraded results.
C. Page latency will be the average of the 40 calls, so it improves.
D. The SLO no longer applies because fan-out is an architectural change.

**Q9 (short answer, hard).** "Availability" and "durability" are often confused. Define each in one sentence and give one example system where you'd prioritize durability over availability, and one where you'd prioritize availability over durability of every single write.

**Q10 (create, hard).** Take the prompt: *"Design a system that lets users upload profile photos and shows them across the app."* Produce a mini requirements document containing: 3 FRs, 4 NFRs (with numbers or measurable qualities), 2 out-of-scope items, 1 SLI/SLO pair, and the hard part. Then write one trade-off statement (with a reversal condition) for the choice of where to store the images.
