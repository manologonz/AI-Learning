# Unit 1: The Senior Mindset & Requirements

## Objectives
By the end of this unit you can:
- Turn a vague prompt into explicit functional requirements, non-functional requirements, and out-of-scope items.
- Define an SLI, SLO, and SLA for a service, and compute an error budget.
- State a design decision as a trade-off: what you gain, what you pay, and under what conditions you'd reverse it.
- Run a repeatable design process (clarify → estimate → API/data → high-level → deep-dive → failure/evolve).
- Identify the "hard part" of a problem, the one constraint that drives the architecture.

## Prerequisites
None. You know how to build services; this unit is about how to *decide* what to build.

## Core content

### 1. What "senior" actually means in system design

A mid-level engineer answers: *"How do I build this?"*
A senior engineer answers: *"What should we build, what will it cost us, what will break, and how will we know?"*

Three habits separate them:

1. **They pin down the problem before drawing boxes.** Most bad architectures solve the wrong problem well.
2. **They think in trade-offs, never "best" solutions.** Every choice buys something and costs something. If you can't name the cost, you don't understand the choice.
3. **They design for failure and change.** The happy path is 10% of the work. What happens when a dependency is slow, when traffic 10x's, when the team needs to change the schema?

Analogy for a JS/TS dev: choosing between `Promise.all` and sequential `await` isn't "which is correct", it's latency vs. load vs. failure isolation. Architecture is the same reasoning at system scale.

### 2. Requirements: the three buckets

**Functional requirements (FR)**: *what the system does.* "Users can post a message. Followers see it in a feed."

**Non-functional requirements (NFR)**: *how well it does it.* Latency, availability, durability, consistency, throughput, scalability, security, cost, compliance. **NFRs, not FRs, drive architecture.** Two systems with identical features but different NFRs (a hobby blog vs. a payments processor) have completely different designs.

**Out of scope**: what you're explicitly *not* building. Saying this out loud is a senior move: it prevents scope creep and shows you've made a choice.

Questions that extract NFRs (memorize the shape, not the list):
- *Scale*: how many users, requests/sec, data volume, growth rate?
- *Read/write ratio and access pattern*: read-heavy? bursty? global?
- *Latency*: what does "fast enough" mean, and at what percentile?
- *Availability*: what does an hour of downtime cost the business?
- *Consistency*: can a user ever see stale data? Which data must never be wrong (money) vs. may lag (like counts)?
- *Durability*: is losing data ever acceptable?
- *Constraints*: team size, existing stack, budget, compliance (GDPR, PCI).

### 3. Measuring quality: SLI, SLO, SLA

- **SLI (Service Level Indicator)**: a *measurement*. A ratio of good events to valid events. Example: `successful requests with latency < 300ms / total valid requests`.
- **SLO (Service Level Objective)**: an internal *target* for an SLI over a window. "99.9% of requests are good over 30 days."
- **SLA (Service Level Agreement)**: an external *contract* with consequences (refunds, penalties). SLAs are set looser than SLOs so you breach your own target before you breach the contract.

**Error budget** = `1 − SLO`. It's the amount of unreliability you're allowed. A 99.9% SLO over 30 days allows 0.1% × 30 days ≈ **43.2 minutes** of full downtime (or equivalent bad requests). The budget turns reliability into a currency: while budget remains, ship features; when it's burned, prioritize stability.

The "nines" table (per 30 days):

| Availability | Downtime allowed |
|---|---|
| 99% | ~7.2 hours |
| 99.9% | ~43 min |
| 99.99% | ~4.3 min |
| 99.999% | ~26 sec |

Each extra nine costs roughly an order of magnitude more effort/money and constrains your design (multi-region, no manual deploys, automated failover). A senior asks: *does the business actually need this nine?*

**Percentiles, not averages.** Latency SLIs use p50/p95/p99, because averages hide the tail. If p99 is 4s, 1 in 100 requests is awful, and a page that makes 50 backend calls hits that tail on most page loads (see "tail amplification": P(all 50 calls avoid the slow 1%) = 0.99^50 ≈ 60%, so ~40% of pages hit at least one slow call).

```ts
// Error budget calculator
function errorBudgetMinutes(slo: number, windowDays = 30): number {
  return (1 - slo) * windowDays * 24 * 60;
}
errorBudgetMinutes(0.999);  // 43.2
errorBudgetMinutes(0.9999); // 4.32

// Request-based SLI: good events / valid events
const sli = (good: number, valid: number) => good / valid;
sli(9_985_000, 10_000_000); // 0.9985 → misses a 99.9% SLO
```

### 4. Trade-off thinking

Every architectural decision is a point on a tension. The classics you'll meet all course:

- Latency ↔ consistency
- Availability ↔ consistency (under partition)
- Throughput ↔ latency
- Simplicity ↔ flexibility
- Cost ↔ redundancy/reliability
- Read speed ↔ write speed (indexes, denormalization, caches)
- Build speed now ↔ change speed later

**Trade-off statement template** (use it in every design you write):

> "We choose **X** over **Y** because **[requirement]** matters more than **[cost]**. This holds while **[assumption]**; if **[assumption breaks]**, we'd revisit."

Example: "We choose a single Postgres primary with read replicas over a sharded store because our write volume is ~300/s (well within one node) and we need transactions. If writes exceed ~5k/s sustained, we'd revisit partitioning."

Note the last clause: seniors attach **reversal conditions**. It's what makes a decision an engineering judgment rather than a preference.

### 5. The design process

A repeatable skeleton keeps you from jumping to "let's use Kafka":

1. **Clarify**: FRs, NFRs, out of scope (this unit).
2. **Estimate**: QPS, storage, bandwidth (Unit 2).
3. **Interface & data model**: APIs and core entities; access patterns drive storage choice.
4. **High-level design**: the simplest architecture that meets the requirements. Start with the boring one.
5. **Deep-dive the hard part**: the bottleneck, the risky component, the place where NFRs bite.
6. **Failure & evolution**: what breaks, how you detect it, how it grows, what you'd change at 10x.

Principle: **start simple, add complexity only against a named requirement.** Every component (cache, queue, shard) must be justified by a number or an NFR. "Because Netflix does it" is not a justification.

### 6. Finding the hard part

Every problem has a dominant constraint. Find it and you know where to spend depth:

| System | The hard part |
|---|---|
| URL shortener | Read-heavy scale, unique ID generation |
| News feed | Fan-out: write amplification for high-follower users |
| Chat | Persistent connections, ordering, delivery guarantees |
| Payments | Correctness: exactly-once effects, consistency, audit |
| Video streaming | Bandwidth/egress cost, CDN |

Heuristic: ask "if I got this one thing wrong, would the system fail its core promise?" That's your hard part.

## Worked examples

### Straightforward: SLO and error budget
*A service has an SLO of 99.95% availability over 30 days. How much downtime is allowed? Last month it was down 25 minutes. Did it meet the SLO?*

1. Window = 30 × 24 × 60 = 43,200 minutes.
2. Budget = (1 − 0.9995) × 43,200 = 0.0005 × 43,200 = **21.6 minutes**.
3. Actual downtime 25 min > 21.6 min → **SLO missed** by 3.4 minutes. Budget is 116% consumed → per error-budget policy, freeze risky launches and prioritize reliability.

### Typical: clarifying a vague prompt
*Prompt: "Design a notification system."*

Bad response: start drawing queues.
Good response, in order:
- **Clarify FRs**: channels (push, email, SMS)? Triggered by events or scheduled? User preferences/opt-out? Templating? → Assume: push + email + SMS, event-triggered, opt-out supported.
- **NFRs**: scale (say 10M users, 50M notifications/day ≈ 580/s average, peaks maybe 10x); latency (transactional like "password reset" within seconds; marketing can take minutes); delivery guarantee (at-least-once acceptable, duplicates tolerable? for OTP, yes duplicates are annoying but not harmful; for "you were charged $500", dedupe matters); availability (a failure to send must not block the triggering action → async).
- **Out of scope**: in-app notification center, rich analytics, A/B testing.
- **Hard part**: reliably delivering through flaky third-party providers with retries, rate limits, and priorities, without duplicates or losing critical messages.

Now you know the design centers on a durable queue, per-channel workers, provider retry with backoff, and priority lanes, all derived from requirements instead of guessed.

### Edge case: the SLO that lies
*A team reports "99.99% availability" based on the load balancer's health checks returning 200. Users complain the app is broken. Why?*

The SLI is measuring the wrong thing. Health checks hit `/healthz`, which returns 200 even when the database is down, checkout returns 500s, or responses take 12 seconds. A good SLI is measured **from the user's perspective**, on real traffic, with both correctness and latency in the definition:

`good = status < 500 AND latency < 500ms` on `POST /checkout`, `valid = all non-4xx-client-error requests`.

Lesson: an SLO is only as honest as its SLI. Also, averaging across all endpoints hides a broken critical one: track SLIs per critical user journey, not one global number.

## Common misconceptions

- **"More nines is always better."** Each nine costs ~10x more and slows delivery. Over-reliability is a real cost; the right SLO is what users actually notice and the business needs.
- **"SLO and SLA are the same thing."** SLO = internal target; SLA = external contract with penalties. SLA should be *looser* than SLO.
- **"Average latency is a good SLI."** Averages hide the tail. Use percentiles; users experience the tail, especially with fan-out.
- **"NFRs are polish we add later."** NFRs choose the architecture. Retrofitting consistency or availability into a design not built for it often means a rewrite.
- **"There's a best architecture."** There's only fit-to-requirements. If you can't state what your choice costs, you haven't analyzed it.
- **"Start with microservices/Kafka/Kubernetes to be 'scalable'."** Complexity is itself a cost (operational, cognitive). Start with the simplest thing that meets the numbers.
- **"Availability = uptime of servers."** Availability is whether *users* can complete requests successfully and fast enough.

## Check your understanding

1. A 99.99% SLO over 30 days: how many minutes of budget? *(Answer: 0.0001 × 43,200 = 4.32 min.)*
2. Your SLA promises 99.9%. Why should the internal SLO be stricter, say 99.95%? *(So the team is alerted and can act before contractual penalties trigger; the gap is a safety margin.)*
3. Two products both let users upload photos. One is a family photo album, the other stores medical scans. Name two NFRs that differ and how they change the design. *(Durability/compliance: medical needs stronger durability guarantees, audit logs, encryption, retention rules; family album can tolerate simpler, cheaper storage. Also availability/latency expectations may differ.)*
4. A page calls 20 services in parallel, each with p99 = 1s. Roughly what fraction of page loads see at least one ≥1s call? *(1 − 0.99^20 ≈ 18%.)*

## Summary
- Seniors clarify requirements first; **NFRs drive architecture**, FRs mostly don't.
- SLI = measurement, SLO = internal target, SLA = external contract; **error budget = 1 − SLO** and turns reliability into a trade-off currency.
- Use percentiles; measure SLIs from the user's perspective on real traffic.
- Every decision is a trade-off: state what you gain, what you pay, and the condition that would reverse it.
- Follow a process: clarify → estimate → interface/data → high-level → deep-dive → failure/evolution. Start simple.
- Find the **hard part**, the constraint that would break the core promise, and spend your depth there.
