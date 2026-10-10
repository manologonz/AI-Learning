# Unit 13 Quiz: Classic Case Studies

Answers in `quiz-answers.md`.

**Q1 (MCQ, easy).** For a URL shortener needing ~6 billion codes over 5 years, how many base62 characters are sufficient?
A. 5 (≈ 9×10⁸)
B. 6 (≈ 5.7×10¹⁰)
C. 7 (≈ 3.5×10¹²)
D. Both B and C are sufficient; 7 gives more headroom
Pick the best answer.

**Q2 (MCQ, easy).** What is the main problem with pure fan-out-on-write for a celebrity with 30M followers?
A. Reads become slow.
B. One post triggers ~30M feed writes, overwhelming the write path.
C. The celebrity's posts can't be stored.
D. Feed reads need multiple queries.

**Q3 (short, medium).** Explain the hybrid feed approach and give the criterion for choosing push vs pull per author.

**Q4 (MCQ, medium).** In a chat system, what is the best way to ensure per-conversation message ordering?
A. Order by client timestamps.
B. Route all messages of a conversation to a single partition/writer that assigns monotonically increasing sequence numbers.
C. Use random UUIDs and sort them.
D. Rely on TCP ordering across connections.

**Q5 (application, medium).** 1M drivers send location every 4 seconds. Compute writes/s and explain why you wouldn't write each update to a durable disk-based relational DB.

**Q6 (MCQ, medium).** Why choose a 302 over a 301 redirect for a shortener with analytics and link expiry?
A. 302 is faster.
B. 301 is cached by browsers/CDNs, so later clicks may never reach your servers (no counting, expiry not enforced); 302 keeps requests flowing through you.
C. 301 breaks HTTPS.
D. 302 is required by the HTTP spec for redirects.

**Q7 (short, hard).** In a payments system, a call to the PSP times out. Describe precisely what the service should do next and what it must not do.

**Q8 (application, hard).** Design the data model for file sync (chunks, file metadata, versions). Show how dedupe and conflict detection work.

**Q9 (short, hard).** For typeahead at ~1M requests/s peak, list three techniques to meet a 100 ms budget and one freshness trade-off.

**Q10 (short, hard).** Why is a "ledger" better than updating a balance column for money? Name two benefits.
