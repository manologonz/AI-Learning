# Further resources: Unit 10

- **AWS Builders' Library: "Timeouts, retries, and backoff with jitter"** (article, aws.amazon.com/builders-library) — the canonical, production-grounded guide to retry storms, jitter, and timeout selection; read it in full.
- **AWS Builders' Library: "Using load shedding to avoid overload"** and **"Static stability using Availability Zones"** (articles) — the sources for load shedding priorities and static stability used in this unit.
- **Google SRE Book, Ch. 21-22: "Handling Overload" and "Addressing Cascading Failures"** (book chapters, free at sre.google) — the definitive treatment of overload, queue behavior, and cascading failure with real Google examples.
- **"Metastable Failures in Distributed Systems" (Bronson et al., HotOS 2021)** (paper) — names and explains the sustained-failure phenomenon seen in the StreamBox incident, with real cases.
- **Netflix Tech Blog: "Performance Under Load" (adaptive concurrency limits)** and Resilience4j docs (blog + library docs) — practical implementations of adaptive limits, circuit breakers, and bulkheads you can adapt in a TypeScript/Node setting.
