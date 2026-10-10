# Unit 10 Quiz: Reliability & Resilience Engineering

Answers in `quiz-answers.md`.

**Q1 (MCQ, easy).** What does "full jitter" exponential backoff accomplish?
A. It makes retries happen faster.
B. It randomizes the sleep between 0 and the exponential cap so clients don't retry in synchronized waves.
C. It guarantees that a retry succeeds.
D. It removes the need for timeouts.

**Q2 (MCQ, easy).** In the circuit breaker's half-open state, what happens?
A. All requests are blocked forever.
B. A limited number of probe requests are allowed through; success closes the circuit, failure re-opens it.
C. All requests pass without checking.
D. The breaker is disabled.

**Q3 (MCQ, medium).** A client calls a gateway, which calls a service, which calls the DB. The client, the gateway, and the service each make up to 3 attempts (1 try + 2 retries) when their downstream call fails. In the worst case, how many DB attempts does one client request cause?
A. 3
B. 9
C. 27
D. 6

**Q4 (application, medium).** A token bucket has capacity 50 and refills 5 tokens/s. A client is idle for a long time, then sends 80 requests instantly. How many succeed immediately? How long until the remaining ones (if retried steadily) would all be admitted?

**Q5 (MCQ, medium).** Which is the best explanation of why "slow" dependencies are dangerous?
A. Slow responses use more bandwidth.
B. Callers hold threads/connections while waiting, exhausting resources and cascading failure to otherwise healthy parts of the system.
C. Slow calls corrupt data.
D. Slow calls make DNS fail.

**Q6 (short, medium).** Give two differences between load shedding and rate limiting.

**Q7 (MCQ, hard).** A service is overloaded: requests queue up for 30 s, and clients time out at 5 s. The server keeps processing timed-out requests. What is the best description and fix?
A. High throughput, so everything is fine.
B. Goodput has collapsed (work done after the client gave up is wasted); fix with bounded queues, load shedding, and deadline propagation so expired requests are dropped.
C. Add more retries.
D. Increase the timeout to 60 s.

**Q8 (short, hard).** Compare pilot light and active-active for multi-region DR: RTO/RPO, cost, and complexity.

**Q9 (application, hard).** Design a rate limiter for an API (100 req/min per API key, bursts OK to 20). Choose the algorithm and where to store state, and describe what happens if the state store fails.

**Q10 (short, hard).** Why is "we have backups" not equivalent to "we can recover"? Name three things to verify.
