# Further resources: Unit 7

- **Designing Data-Intensive Applications, Ch. 8 (The Trouble with Distributed Systems) and Ch. 9 (Consistency and Consensus)** (Kleppmann, book) — the clearest treatment of unreliable clocks, fencing, linearizability, and consensus; it's the backbone of this unit.
- **The Raft paper and visualization** (paper "In Search of an Understandable Consensus Algorithm" plus raft.github.io) — the visualization lets you kill leaders and partition the cluster, which builds real intuition.
- **Martin Kleppmann: "How to do distributed locking"** (blog post, martin.kleppmann.com) — the definitive argument for fencing tokens over lease-only locks.
- **Jepsen analyses** (jepsen.io/analyses) — pick one database you use and read how it was tested; shows how real systems violate their consistency claims.
- **"Time, Clocks, and the Ordering of Events in a Distributed System"** (Lamport, 1978 paper) — short and foundational; worth reading once to see where logical clocks and "happened-before" come from.
