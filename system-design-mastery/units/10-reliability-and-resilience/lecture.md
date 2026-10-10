# Unit 10: Reliability & Resilience Engineering

## Objectives
By the end of this unit you can:
- Set timeouts, retries (with backoff and jitter), and retry budgets that don't cause retry storms.
- Apply circuit breakers, bulkheads, load shedding, and backpressure to stop cascading failures.
- Design a rate limiter (token bucket, sliding window) and choose where to enforce it.
- Define failure domains and design multi-AZ/multi-region availability with explicit RTO/RPO.
- Plan graceful degradation and test failure modes (chaos, game days).

## Prerequisites
Units 1-9. This unit is where earlier ideas (SLOs, caching, replication, idempotency) turn into production survival skills.

## Core content

### 1. How systems actually fail
Big outages are rarely a single component dying. They're **cascading failures**: one component slows, callers wait and pile up (threads/connections exhausted), retries multiply load, dependencies collapse, and the whole system goes down, often *after* the original trigger has ended (**metastable failure**: the system stays broken because retries/overload sustain themselves).

Design principles:
- **Assume every dependency will be slow or down.** Slow is worse than down: down fails fast; slow ties up resources.
- **Fail fast, fail small, degrade gracefully.**
- **Protect the scarce resource** (threads, connections, DB capacity) at every layer.
- **Availability math**: serial dependencies multiply (0.999 × 0.999 = 0.998); redundancy in parallel improves (`1 − (1−a)^n` if independent, which they rarely are: shared config, shared deploys, correlated failures).

### 2. Timeouts
Every network call needs a timeout. **No timeout = unbounded resource holding.** Set from observed latency: e.g., slightly above p99.9 of the dependency's healthy latency, and within the caller's own deadline. Use **deadline propagation**: the request carries a deadline (`grpc-timeout`, `X-Request-Deadline`); each hop subtracts elapsed time and abandons work when the caller has already given up (otherwise servers waste capacity on responses no one will read). Include connect, read, and total timeouts.

### 3. Retries (and retry storms)
Retries fix transient failures but can amplify outages:
- **Retry only retryable errors** (timeouts, 503, connection reset; not 400/404), and only **idempotent** operations (or with idempotency keys).
- **Exponential backoff with jitter**: `sleep = random(0, min(cap, base × 2^attempt))` ("full jitter"). Jitter desynchronizes clients to avoid thundering herds.
- **Retry budgets**: limit retries to a fraction of traffic (e.g., ≤10% extra), or use per-client token buckets. **Retry only at one layer**: if 3 layers each retry 3 times, one failure becomes 3³ = 27 calls (retry amplification).
- **Hedged requests**: send a duplicate to another replica after the p95 latency and take the first response; cuts tail latency for idempotent reads at ~5% extra load.

```ts
// Retry with capped exponential backoff + full jitter + deadline awareness
async function withRetry<T>(fn: () => Promise<T>, { max = 3, base = 100, cap = 2000, deadline = Date.now() + 5000 } = {}) {
  for (let attempt = 0; ; attempt++) {
    try { return await fn(); }
    catch (e: any) {
      const retryable = e.code === "ETIMEDOUT" || e.status === 503 || e.status === 429;
      if (!retryable || attempt >= max) throw e;
      const sleep = Math.random() * Math.min(cap, base * 2 ** attempt);    // full jitter
      if (Date.now() + sleep >= deadline) throw e;                           // no time left
      await new Promise(r => setTimeout(r, sleep));
    }
  }
}
```

### 4. Circuit breakers, bulkheads, and timeouts working together
- **Circuit breaker**: tracks failures to a dependency. **Closed** (normal) → **Open** (after threshold: fail immediately, no calls, giving the dependency room to recover) → **Half-open** (after a cool-down, allow a few probe calls; success closes, failure re-opens). Prevents piling requests onto a dying service and frees your threads.
- **Bulkheads**: isolate resources so one failing dependency can't consume everything: separate thread/connection pools or concurrency limits per dependency, separate queues per tenant/priority. Like ship compartments.
- **Fallbacks/graceful degradation**: return cached/stale data, a default, or omit the feature (recommendations blank, "prices may be outdated"). Decide per feature what's *critical* (checkout) vs *optional* (reviews).

### 5. Load shedding and backpressure
- **Load shedding**: when overloaded, reject excess work early (HTTP 503/429) to keep serving what you can, rather than accepting all and failing everything after long queues. Shed by priority (health checks and paying customers first), by cost, or by queue latency. Reject *before* expensive work; cheap rejection preserves goodput. Little's Law: `L = λW`: as latency W grows under overload, in-flight L explodes unless bounded.
- **Bounded queues and concurrency limits**: an unbounded queue turns overload into latency (long delays, timeouts, then useless work). Prefer small queues + rejection. **Adaptive concurrency limits** (AIMD/gradient like Netflix's) find the safe concurrency automatically.
- **Backpressure**: propagate "slow down" upstream (blocking producers, TCP flow control, reactive streams, HTTP 429 with `Retry-After`) so overload doesn't just move the buffer somewhere else.
- **Goodput vs throughput**: at overload, throughput may stay high while *goodput* (useful completed work) drops to zero because responses arrive after clients timed out.

### 6. Rate limiting
Protects capacity, enforces fairness, and prevents abuse.
- **Token bucket**: bucket holds up to B tokens, refilled at rate r; each request takes one; allows bursts up to B, average r. (Most common.)
- **Leaky bucket**: smooths output to a constant rate.
- **Fixed window counter**: simple but allows 2x bursts at window boundaries.
- **Sliding window log/counter**: accurate, more memory/computation; the counter approximation is cheap and good.
- **Where**: edge/API gateway (per API key/IP), per service, per tenant, and *client-side* (adaptive throttling). **Distributed limiters**: central Redis counter (Lua script for atomicity), or local limits with a share of the global rate (approximate but no hot dependency). Decide *fail-open vs fail-closed* when the limiter store is down (usually fail-open for availability, closed for security-critical).
- Respond with **429** + `Retry-After`; make limits observable and documented.

```ts
// Token bucket in Redis (atomic via Lua) — conceptual
const lua = `
local tokens = tonumber(redis.call('HGET', KEYS[1], 't') or ARGV[1])
local ts = tonumber(redis.call('HGET', KEYS[1], 'ts') or ARGV[3])
local delta = math.max(0, ARGV[3] - ts) * ARGV[2]        -- refill = elapsed * rate
tokens = math.min(ARGV[1], tokens + delta)
local ok = tokens >= 1
if ok then tokens = tokens - 1 end
redis.call('HSET', KEYS[1], 't', tokens, 'ts', ARGV[3]); redis.call('EXPIRE', KEYS[1], 3600)
return ok and 1 or 0`;
// args: capacity B, refill rate (tokens/ms), now(ms)
```

### 7. Failure domains, redundancy, and disaster recovery
- **Failure domains**: process → host → rack → **availability zone (AZ)** → region → provider/global control plane. Spread replicas across domains that fail independently; **static stability**: the system keeps working when a control plane or dependency is down, by not needing to make changes (pre-provisioned capacity, cached config).
- **Multi-AZ** (default for production): survive a data center failure with sync-ish replication and automatic failover; N+1 capacity across AZs (Unit 2).
- **Multi-region** patterns: **backup & restore** (cheap, RTO hours), **pilot light** (minimal core running, scale up on disaster), **warm standby** (scaled-down copy), **active-active** (serves everywhere; needs conflict handling and data partitioning). Cost and complexity rise steeply; choose from the *business's* RTO/RPO.
- **RTO** (Recovery Time Objective): how long until service is restored. **RPO** (Recovery Point Objective): how much data loss is acceptable (time). Async cross-region replication → RPO = replication lag; sync → RPO 0 at latency cost.
- **Backups are not disaster recovery until you've restored from them.** Test restores; keep backups isolated (immutability, separate account) from ransomware/operator error.
- **Deploy safety**: most outages follow changes. Canary, gradual rollout, feature flags, automatic rollback on SLO burn, config validation, and avoiding synchronized changes across all regions.

### 8. Testing and operating for failure
- **Chaos engineering**: inject failures (kill instances, add latency, drop dependencies) in a controlled way to verify assumptions; start in staging, then production with a small blast radius.
- **Game days** and **load tests**, including *cold-cache*, *dependency-slow*, and *AZ-loss* scenarios.
- **Blameless postmortems**: root causes are system and process conditions; produce concrete action items with owners.
- **Runbooks and alerts on SLO burn rate** (Unit 12), not on causes.

## Worked examples

### Straightforward: retry amplification
*Client → API → Service A → DB. Each of the 3 hops retries 3 times on failure. The DB has a brief blip.* One client request can produce up to 3 × 3 × 3 = **27** DB attempts, right when the DB is struggling, prolonging the outage. Fix: retry at one layer only (the edge/client), use a retry budget (≤10%), and propagate a deadline so inner layers don't retry beyond the caller's patience.

### Typical: protect a service from a slow dependency
*The checkout service calls Recommendations (optional) and Payments (critical), sharing a pool of 200 threads. Recommendations becomes slow (10 s).* Threads pile up waiting; within seconds all 200 are blocked and checkout can't call Payments: total outage from a non-critical dependency. Fix: **timeouts** (e.g., 200 ms on recommendations), a **bulkhead** (max 20 concurrent recommendation calls, separate pool), a **circuit breaker** (open after 50% failures in 10 s; half-open probes), and a **fallback** (skip recommendations). Result: checkout unaffected.

### Edge case: metastable failure from retries and a cold cache
*A brief DB failover causes timeouts; clients retry aggressively, doubling load; the cache empties from expirations; the DB can't recover under retry + cold-cache load even after the failover ends; it stays down for an hour until traffic is manually cut.* The trigger was transient; the *sustaining effect* was retries and overload. Fixes: retry budgets and jittered backoff, load shedding that reduces admitted traffic to the DB's capacity, a **circuit breaker** that pauses calls to the recovering DB, gradual ramp-up (slow start) after recovery, and cache warming. Recovery procedure: shed load (or block traffic) until queues drain, then ramp up.

## Common misconceptions

- **"Retry on every failure to be robust."** Unbudgeted retries amplify outages; retry only idempotent, retryable errors, at one layer, with jitter.
- **"Timeouts are just a nice-to-have."** Missing timeouts are a leading cause of cascading failures.
- **"Multi-AZ = disaster recovery."** It protects against a data-center failure, not a region-wide outage, data corruption, or a bad deploy replicated everywhere.
- **"Redundant components make availability 1−(1−a)^n."** Only if failures are independent; shared dependencies, config pushes, and correlated bugs break this.
- **"Autoscaling will handle overload."** It reacts in minutes; overload kills you in seconds; and scaling a stressed dependency (DB) may not be possible. Shed load.
- **"Rate limiting = DoS protection."** It's a start; volumetric DDoS needs edge/CDN/anycast absorption.
- **"Circuit breakers fix failures."** They limit blast radius and give dependencies time to recover.
- **"Backups exist, so we're safe."** Untested restores are hopes.
- **"Graceful degradation is a UX detail."** It's an architectural decision: it must be planned per dependency.

## Check your understanding

1. Three layers each retry 2 times (3 attempts total each). Worst-case attempts at the bottom for one request? *(3 × 3 × 3 = 27.)*
2. Why is "slow" often worse than "down"? *(Slow ties up threads/connections; down fails fast.)*
3. A token bucket has capacity 100 and refill 10 tokens/s. Maximum sustained rate and max burst? *(10/s sustained; burst of 100.)*
4. RPO of 5 minutes means what? *(You can lose up to 5 minutes of data on recovery.)*

## Summary
- Failures cascade; slowness kills more than crashes. Protect scarce resources everywhere.
- Timeouts + deadline propagation + limited, jittered retries at one layer; hedge tail latency for idempotent reads.
- Circuit breakers, bulkheads, fallbacks isolate failures; load shedding and bounded queues preserve goodput; backpressure moves signals upstream.
- Token bucket is the workhorse rate limiter; decide fail-open vs fail-closed.
- Design failure domains, RTO/RPO, and DR tiers from business requirements; test restores, failover, and failure modes; most outages follow changes.
