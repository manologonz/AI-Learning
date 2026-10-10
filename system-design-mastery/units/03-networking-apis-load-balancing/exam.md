# Unit 3 Exam: Networking, APIs & Load Balancing

Answers in `exam-answers.md`.

---

### Scenario (Q1-Q5): ShopFast

ShopFast is a global e-commerce site: web + mobile apps, ~40k peak QPS, users in North America, Europe, Asia. Stack: a fleet of stateless Node services behind one region's (us-east) load balancer, one Postgres primary. Complaints: (1) Asia users report 3-second first loads; (2) some customers were charged twice during a network incident; (3) during a DB failover last month the entire site went down for 4 minutes even though only checkout depends on the DB; (4) the "orders" list API times out for power users with 50k orders.

**Q1 (analyze, medium).** Explain the likely cause of the 3-second first loads for Asia users, using RTT reasoning (assume ~220 ms RTT to us-east). List three mitigations, ordered by impact.

**Q2 (apply, medium).** Explain the double-charge root cause in terms of request semantics, and design the fix: what does the client send, what does the server store, and what happens on (a) a replay after success, (b) a replay while the first attempt is still running, (c) the same key with a *different* body?

**Q3 (analyze, hard).** Diagnose the 4-minute site-wide outage during the DB failover. Name the likely design flaw in the LB configuration, what recovery-time behavior made it worse, and the fix (at least three concrete changes).

**Q4 (apply, hard).** The orders list uses `LIMIT 20 OFFSET n`. Explain why power users time out, write the keyset SQL replacing it, and describe how you'd evolve the API response and keep old mobile clients working.

**Q5 (evaluate, hard).** Product wants a live "order status" screen (packed → shipped → out for delivery). Compare polling, SSE, and WebSocket for this, pick one, and justify with an estimate: 2M concurrent users watching, updates rare (a few per order). Include one operational concern of your choice (deploys, scaling, proxies).

---

### Independent items

**Q6 (MCQ, medium).** Which load balancing algorithm is best when requests vary widely in cost (some take 5 ms, others 2 s) and you have no shared load state across many LB nodes?
A. Round robin
B. Power of two random choices using each backend's outstanding request count
C. Hash on client IP
D. Random with equal weights

**Q7 (MCQ, hard).** A team uses sticky sessions (cookie affinity) to keep user sessions in server memory. What is the most serious consequence at scale?
A. Higher TLS cost
B. Uneven load, lost sessions on server failure or deploy, and harder autoscaling; the fix is to externalize session state and make servers stateless
C. The cookie leaks user IDs
D. It disables HTTP/2

**Q8 (MCQ, hard).** Choose the best API style/transport pairing for: (i) public partner API for order lookups; (ii) internal high-QPS payment-to-ledger calls; (iii) a mobile app screen assembling data from six services.
A. (i) gRPC, (ii) REST, (iii) WebSocket
B. (i) REST/JSON, (ii) gRPC, (iii) a backend-for-frontend or GraphQL layer aggregating the services
C. (i) GraphQL, (ii) SSE, (iii) REST
D. (i) SSE, (ii) GraphQL, (iii) gRPC

**Q9 (short, hard).** Explain why "L7 load balancer terminating TLS" is often preferred over TLS passthrough at L4, and name one situation where you'd choose passthrough.

**Q10 (create, hard).** Draw (in text) the request path for a first-time user in Tokyo hitting ShopFast after you've applied your fixes: DNS, CDN, LB, service, cache, DB. For each hop, write the expected latency budget in ms so the total first-byte time is under 600 ms, and state which hops you eliminated versus the original.
