# Further resources: Unit 3

- **High Performance Browser Networking, Ch. 1-4 (Ilya Grigorik, free at hpbn.co)** (book) — the definitive, readable account of latency, TCP/TLS handshake cost, and HTTP/2; ties directly to the RTT budgeting in this unit.
- **Stripe API docs: "Idempotent requests"** (docs, stripe.com/docs/api/idempotent_requests) — a real production example of the idempotency-key contract, including key reuse rules.
- **"The Power of Two Random Choices"** (Mitzenmacher, paper/overview) — short explanation of why sampling two servers nearly matches global least-load balancing; useful for justifying LB algorithms.
- **Use The Index, Luke: "Paging Through Results"** (use-the-index-luke.com/no-offset) — the clearest explanation of keyset vs offset pagination with SQL and index details.
- **Envoy docs: "Load balancing" and "Health checking"** (docs, envoyproxy.io) — a real L7 proxy's algorithms, panic thresholds, and outlier detection; good for seeing the health-check ideas in configuration.
