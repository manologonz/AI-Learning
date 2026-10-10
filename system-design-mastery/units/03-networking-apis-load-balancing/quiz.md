# Unit 3 Quiz: Networking, APIs & Load Balancing

Answers in `quiz-answers.md`.

**Q1 (MCQ, easy).** Which HTTP method is *not* idempotent by definition?
A. GET
B. PUT
C. DELETE
D. POST

**Q2 (MCQ, easy).** Which load balancer type can route requests based on the URL path?
A. L4
B. L7
C. Both, equally
D. Neither; only DNS can do that

**Q3 (MCQ, medium).** A client retries `POST /orders` after a network timeout and the customer is charged twice. What's the standard fix?
A. Change the method to GET.
B. Have the client send an idempotency key that the server stores and uses to return the original result on replay.
C. Disable retries entirely on the client.
D. Increase the server timeout to 60 seconds.

**Q4 (short, medium).** Give two concrete reasons cursor (keyset) pagination beats offset pagination for a feed with millions of rows and constant inserts.

**Q5 (MCQ, medium).** You use `hash(user_id) % N` to choose a cache server. You add one server (N → N+1). What happens?
A. Only 1/(N+1) of the keys move.
B. Roughly N/(N+1) of keys remap, causing a mass cache miss.
C. Nothing; hash modulo is stable.
D. Keys owned by the new server are duplicated on all others.

**Q6 (MCQ, medium).** Which real-time transport is best for a live sports-score ticker (server → client only) that must work through typical corporate proxies with automatic reconnect?
A. WebSocket
B. Server-Sent Events
C. gRPC bidirectional streaming
D. Short polling every 100 ms

**Q7 (application, medium).** A user 200 ms RTT away opens a fresh HTTPS connection (TLS 1.3) with no cached DNS. Roughly how long until the request has been sent and the server can begin responding, ignoring server processing? Break it into steps.

**Q8 (MCQ, hard).** An LB health check calls `/health`, which runs `SELECT 1` on the shared database. The DB stalls for 30 seconds. What is the most likely outcome?
A. Only the DB is affected; app servers stay in rotation.
B. Every app instance fails its check, the LB ejects the whole fleet, and users see a total outage even for endpoints that don't need the DB.
C. The LB routes traffic to a backup DB automatically.
D. The health checks are cached so nothing happens.

**Q9 (application, hard).** Design the cursor for `GET /posts?after=` on a table sorted by `created_at DESC`. Explain why `created_at` alone is insufficient and what the SQL predicate should look like.

**Q10 (short, hard).** Your internal services use REST/JSON today and profiling shows 30% of CPU on JSON (de)serialization plus frequent schema drift bugs between teams. Would you move to gRPC? Give one benefit, one cost, and one condition under which you'd *not* migrate.
