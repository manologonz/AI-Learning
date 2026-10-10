# Unit 6 Exam: Replication & Partitioning

Answers in `exam-answers.md`.

---

### Scenario (Q1-Q6): ChatLoop

ChatLoop is a group-messaging app: 100M users, 2M messages/s at peak (writes), messages read mostly right after being written and again when scrolling history. Data: messages (channel_id, message_id, sender_id, body, ts). The current system is one Postgres primary with 5 read replicas; the primary is at 95% CPU and storage is 30 TB and growing 1 TB/week. There are 20M channels; a few "mega-channels" (announcement channels) have 5M members.

**Q1 (analyze, medium).** Why won't adding more read replicas fix the primary's problems? What does replication do and not do here?

**Q2 (apply, hard).** Propose a partitioning scheme and shard key for messages. State the primary key/clustering order, why it fits the two main queries (latest 50 in a channel; scroll back), and what queries it makes expensive.

**Q3 (analyze, hard).** Identify two hot-spot risks in your design (one from mega-channels, one from time) and give a mitigation for each.

**Q4 (apply, hard).** You have 20 nodes today and expect to grow to 30 next quarter. Compare (a) `hash % N`, (b) consistent hashing with vnodes, and (c) fixed 1,024 partitions. For each, give the fraction of data that moves when going from 20 to 30 nodes (or "approximately") and one drawback.

**Q5 (evaluate, hard).** Users complain that in group chats they sometimes see a reply before the message it replies to. Explain which replication anomaly this is and how partitioning choices caused or could prevent it.

**Q6 (evaluate, hard).** The team wants multi-region (US/EU/APAC) with users writing to their local region. Choose between single-leader-per-shard with region-pinned leaders, multi-leader, and leaderless. Justify the choice and state how conflicts (or lack thereof) are handled for chat messages.

---

### Independent items

**Q7 (MCQ, medium).** A Cassandra-style cluster has N=3, W=2, R=1. Which statement is true?
A. Reads always see the latest acknowledged write.
B. A read may return stale data because W+R = 3 is not greater than N.
C. Writes can tolerate two node failures.
D. The cluster is linearizable.

**Q8 (MCQ, hard).** With synchronous replication to one follower and asynchronous to two others (semi-sync), what does the leader crash lose?
A. Nothing ever
B. No acknowledged writes if failover promotes the synchronous follower; writes acknowledged only to async followers could be lost otherwise
C. All writes from the last hour
D. Only reads

**Q9 (short, hard).** Why can a scatter-gather query across 100 shards have terrible p99 latency even if each shard's p99 is 20 ms? Give the approximate probability that at least one shard exceeds its p99, and a mitigation.

**Q10 (create, hard).** Design the failover procedure for a Postgres primary with 2 replicas: detection, choosing the candidate, fencing the old primary, client redirect, and what you'd tell the business about data loss guarantees (RPO) and downtime (RTO) for async vs sync replication.
