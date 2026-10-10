# Further resources: Unit 13

- **Alex Xu, "System Design Interview" Vol. 1 and 2** (books) — the most widely used collection of worked case studies (URL shortener, news feed, chat, typeahead, proximity service, ledger); use as a source of extra practice prompts after you've derived your own designs.
- **High Scalability: "Architecture" case studies** (highscalability.com) — real production architectures summarized with numbers; good for seeing how actual systems deviate from textbook answers.
- **Discord Engineering: "How Discord Stores Trillions of Messages"** (blog) — a real chat-storage story (Cassandra hot partitions → ScyllaDB), directly reinforcing Case 3 and Unit 6.
- **Twitter/X Engineering: "Timelines at Scale" (InfoQ talk, Raffi Krikorian)** (talk) — the canonical description of the fan-out-on-write vs read hybrid used by a real feed.
- **Stripe Engineering: "Designing robust and predictable APIs with idempotency"** (blog, stripe.com/blog/idempotency) — the payments idempotency design in Case 7, from a company that runs it.
