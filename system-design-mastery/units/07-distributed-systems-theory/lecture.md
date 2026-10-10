# Unit 7: Distributed Systems Theory

## Objectives
By the end of this unit you can:
- Describe distributed-system failure modes (partial failure, network partitions, slow-vs-dead ambiguity) and design around them.
- Explain why wall-clock time is unreliable and use logical clocks and fencing tokens correctly.
- Apply CAP and PACELC without the common misreadings.
- Distinguish consistency models: linearizable, sequential, causal, read-your-writes, eventual.
- Explain how consensus (Raft) elects leaders and replicates logs, and compute quorum tolerance (2f+1).

## Prerequisites
Unit 6 (replication, quorums, failover).

## Core content

### 1. What makes distributed systems different
On one machine, a function call either returns or the process dies. In a distributed system you get **partial failure**: some components fail while others continue, and you often **can't tell** which. Sending a request and hearing nothing could mean: the request was lost, the server is slow, the server crashed before or after processing, or the response was lost. **The network gives you no way to distinguish these.** Everything in this unit follows from that.

**Fallacies of distributed computing**: the network is reliable; latency is zero; bandwidth is infinite; the network is secure; topology doesn't change; there is one administrator; transport cost is zero; the network is homogeneous.

**Failure models** (weakest assumptions = most robust and most costly):
- **Crash-stop**: nodes fail and never return.
- **Crash-recovery**: nodes crash and may return with durable state (typical assumption).
- **Omission**: messages are lost.
- **Byzantine**: nodes behave arbitrarily/maliciously (blockchains, aircraft; not typical in a single trusted org).
- **Network partition**: the network splits into groups that can't communicate.

**Timeouts are the only failure detector**, and they're inherently imperfect: too short → false positives (healthy node declared dead, causing failover storms); too long → slow recovery. Adaptive detectors (phi-accrual) tune to observed variance.

### 2. Time and ordering
- **Time-of-day clocks** (NTP-synced): can jump backward/forward, differ across machines by milliseconds to seconds. **Never use them to order events across nodes or as a lock's correctness basis.** Last-write-wins on wall clocks silently loses writes when clocks skew.
- **Monotonic clocks**: only for measuring durations on one machine.
- **Logical clocks**:
  - **Lamport clocks**: a counter incremented on each event; on receiving a message, `clock = max(local, received) + 1`. If A happened-before B then L(A) < L(B) (not the converse); gives a consistent total order (with node-id tiebreaks) but can't detect concurrency.
  - **Vector clocks**: a vector of counters (one per node); can tell whether two events are causally ordered or **concurrent** (conflict); size grows with nodes. Used in Dynamo-style siblings.
  - **Hybrid logical clocks (HLC)**: physical time + logical counter; used by CockroachDB.
- **TrueTime (Google Spanner)**: exposes clock uncertainty `[earliest, latest]`; commits *wait out* the uncertainty (~7 ms) to get external consistency: using atomic clocks/GPS.
- **Process pauses**: GC, VM migration, or `SIGSTOP` can freeze a process for seconds while others declare it dead, then it resumes believing it's still the leader. This is why leases alone are dangerous.

**Fencing tokens.** A lock/lease service issues an increasing token with each grant. The protected resource (DB, storage) rejects any request carrying a token lower than one it has already seen. Then a paused old leader waking up can't corrupt state.

```ts
// Resource side: enforce fencing tokens
let highestToken = 0;
function write(token: number, data: unknown) {
  if (token < highestToken) throw new Error("stale leader (fenced)");
  highestToken = token;
  storage.apply(data);
}
// Client A holds token 33, pauses (GC). Lease expires. Client B gets token 34 and writes.
// A wakes up and writes with token 33 → rejected. Lease timeouts alone would NOT have saved you.
```

### 3. CAP and PACELC (and what they don't say)
**CAP**: in the presence of a network **P**artition, a system must choose between **C**onsistency (linearizability: every read sees the latest write) and **A**vailability (every non-failing node answers). Partitions *will* happen, so the real choice is CP vs AP **during a partition**.
- **CP** (e.g., ZooKeeper/etcd, HBase, Spanner): minority side refuses requests to avoid inconsistency.
- **AP** (e.g., Cassandra/Dynamo defaults, DNS): all sides keep serving, possibly with stale/conflicting data reconciled later.

Common misreadings: (1) "pick 2 of 3": you can't give up P; (2) CAP's "C" means *linearizability*, not ACID consistency; (3) CAP says nothing about the *normal* (no-partition) case. That's what **PACELC** adds: **if Partition → A or C; Else → Latency or Consistency.** Even without failures, stronger consistency costs latency (coordination across replicas/regions). Examples: Dynamo/Cassandra = PA/EL; Spanner/MongoDB-default = PC/EC-ish; PNUTS = PC/EL.

### 4. Consistency models (strongest to weakest)
1. **Linearizability (strong/atomic consistency)**: operations appear to take effect atomically at a single point between invocation and response; behaves like a single copy. Needed for: leader election/locks, uniqueness constraints, account-balance checks. Costs latency and availability.
2. **Sequential consistency**: all processes see the same order of operations, consistent with each process's program order, but not necessarily real-time order.
3. **Causal consistency**: causally related operations are seen in order everywhere (reply after question); concurrent ones may vary. Achievable while staying available under partitions (with caveats), a sweet spot for many apps.
4. **Session guarantees**: read-your-writes, monotonic reads, monotonic writes: per-client guarantees.
5. **Eventual consistency**: if writes stop, replicas converge. No ordering promises in the meantime.

Ask: *what anomaly can the user notice and does it hurt?* Use the weakest model that protects the invariant.

### 5. Consensus, in one page
**Consensus**: multiple nodes agree on a value (or a sequence of values, i.e., a log) despite failures. It underpins leader election, atomic broadcast, distributed locks, and replicated state machines. **FLP impossibility**: in a fully asynchronous system, no deterministic consensus algorithm can guarantee termination if even one node may crash. In practice, algorithms (Paxos, Raft, Zab, Viewstamped Replication) stay **safe** always and are **live** when the network is reasonably behaved (using timeouts).

**Raft in outline**
- **Roles**: leader, followers, candidates. Time is divided into **terms**.
- **Leader election**: a follower whose election timer fires (randomized to avoid ties) becomes a candidate, increments the term, and requests votes; **a majority** wins. Each node votes once per term, and only for a candidate whose log is at least as up-to-date as its own.
- **Log replication**: clients send commands to the leader; it appends to its log and replicates to followers; an entry is **committed** once stored on a **majority**; then it's applied to the state machine and acknowledged. A leader change never loses committed entries (elected leader must hold them).
- **Safety**: at most one leader per term; committed entries persist. **Liveness** needs a stable leader and majority connectivity.

**Quorum arithmetic**: a cluster of **2f+1** nodes tolerates **f** failures (majority = f+1). 3 nodes → 1 failure; 5 → 2; 7 → 3. Even sizes add cost without adding fault tolerance (4 nodes still tolerate only 1). Two disjoint majorities can't exist, which is what prevents split brain: the minority side of a partition can't elect a leader or commit.

**Cost**: consensus requires a round trip to a majority per commit: latency = the RTT to the median-nearest replica; cross-region consensus costs tens to hundreds of ms. **Use consensus for small, critical coordination data** (metadata, leadership, config, locks, small transactional state) via etcd/ZooKeeper/Consul, *not* as your high-volume data path, unless you're building a system like Spanner/CockroachDB that partitions data into many small consensus groups.

## Worked examples

### Straightforward: quorum and partition
*A 5-node etcd cluster splits into groups of 3 and 2.* The 3-group has a majority (≥ 3), continues to elect a leader and commit. The 2-group cannot: reads/writes fail (or reads may be stale if serviced locally). This is CP behavior; no split brain. If it splits 2/2/1, nobody has a majority → cluster unavailable until healed.

### Typical: a distributed lock done wrong and right
*Job workers use a Redis lock with a 30 s expiry to ensure one worker updates a file in object storage.* Worker A takes the lock, hits a 40 s GC pause; the lock expires; Worker B takes it and writes; A resumes and overwrites B's work: corruption, even though the "lock" logic was correct. Fix: a lock service that returns **fencing tokens**, and storage that rejects lower tokens (or use conditional writes: `If-Match` on ETag/version). Lesson: a lease is a *liveness* optimization; correctness needs fencing or compare-and-set at the resource.

### Edge case: LWW with clock skew
*Two data centers use last-write-wins on timestamps. DC-1's clock runs 300 ms ahead. A user changes their name in DC-2 at real time t; DC-1's earlier write (real time t−100 ms) carries timestamp t+200 ms and "wins".* The later real update is silently discarded. No error; no alert. LWW is only safe if lost updates are acceptable or clocks are tightly bounded (TrueTime-like), or if you use per-key version vectors/CRDTs instead.

## Common misconceptions

- **"CAP: choose any two."** Partitions aren't optional; the choice is C vs A *during* a partition, and CAP's "C" is linearizability.
- **"Eventual consistency means data might be wrong forever."** It converges when writes stop; the question is how long anomalies are visible and whether they hurt.
- **"NTP-synced clocks are good enough to order events."** Skew of tens of ms is normal; seconds happen. Use logical clocks or fencing.
- **"A lock with a timeout is safe."** Pauses and delays let two holders coexist; you need fencing/CAS.
- **"Raft/Paxos make everything strongly consistent and fast."** They cost a majority round trip per commit and throughput is bounded by the leader.
- **"More nodes in a consensus group is better."** More nodes = higher commit latency and no better failure tolerance at even sizes; 3 or 5 is typical.
- **"Timeouts detect failures."** They detect *silence*, not death; a slow node is indistinguishable from a dead one.
- **"Strong consistency isn't needed."** Uniqueness, money, and leader election need it; sprinkle it where invariants demand.

## Check your understanding

1. 7-node cluster: how many failures tolerated, and what majority is needed? *(3 failures; majority of 4.)*
2. Why can't you rely on the local clock to decide "my lease is still valid"? *(Clock drift/jumps and process pauses; use fencing tokens.)*
3. Classify: a bank account debit check; a "like" counter; a chat conversation. Which consistency model? *(Debit: linearizable/serializable; likes: eventual; chat: causal/consistent prefix.)*
4. What does PACELC add to CAP? *(Trade-off between latency and consistency when there is no partition.)*

## Summary
- Partial failure and unreliable networks mean you can't distinguish slow from dead; design for ambiguity (idempotency, timeouts, retries).
- Wall clocks can't order events; use logical clocks, and fencing tokens for correctness.
- CAP = C vs A during partitions; PACELC adds latency vs consistency otherwise; C in CAP means linearizability.
- Pick the weakest consistency model that protects your invariant; use linearizability for locks, uniqueness, and money.
- Consensus (Raft) uses majority quorums (2f+1 tolerates f) for elections and log commits; use it for small critical state.
