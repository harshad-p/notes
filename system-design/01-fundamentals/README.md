Let's build the vocabulary and fundamentals first.

The goal of this curriculum is:

> **When we eventually return to system-design questions, I should be able to say "use a read replica here because..." rather than having to teach you what a read replica is while designing the system.**

And importantly, this is **not a generic computer-science curriculum**. It is specifically the knowledge you need to reason about backend system design interviews.

# System Design Fundamentals Curriculum

The curriculum is **12 modules**. Refer to [plan.md](plan.md) for the full outline. That's the level we want before returning to full system designs.

# The order matters

I would **not** teach these in random order.

The dependency chain should be:

```text
Module 1 — Fundamentals (components, request lifecycle, capacity)
    ↓
Module 2 — Scalability
    ↓
Module 3 — Caching
    ↓
Module 4 — Data storage
    ↓
Module 5 — Distributed systems fundamentals
    ↓
Module 6 — Messaging and event-driven systems
    ↓
Module 7 — Reliability and availability
    ↓
Module 8 — Global systems
    ↓
Module 9 — System-design building blocks
    ↓
Module 10 — Observability
    ↓
Module 11 — Putting the fundamentals together (exercises)
    ↓
Module 12 — Full system-design interview practice
```

For example, I don't want to teach you **split brain** before you understand replication, failover, leaders, and distributed failure.

Likewise, I don't want to throw **Kafka partitions** at you before you understand why a queue/event stream exists.

And I don't want to teach **CAP** as a memorized acronym before you've experienced the actual consistency/availability problem that CAP describes.

---

# What I want you to be able to do afterward

When we eventually return to system-design problems, I want you to be able to hear:

> "We have millions of users, the database is overloaded, and users are worldwide."

and immediately start reasoning about:

**load → bottleneck → replication → caching → sharding → regions → consistency → failure modes**

rather than thinking:

> "Hmm, maybe we need Redis... and Kafka... and Kubernetes..."

That distinction is important.

You should be able to **justify every component you introduce**.

And then our system-design practice can finally follow the interview format you originally wanted:

> **"Before I design it, I'd like to clarify a few requirements..."**

From that point onward, **you explain the design to the interviewer**. I won't be simultaneously teaching you what every component means. I'll only intervene when you make a mistake, miss an important issue, or need to go deeper.

### Recommended starting point

I'd start with **Module 1 — System Design Fundamentals**, specifically **1.1 What is a distributed system?**

Once we have worked through the curriculum, we'll return to system-design problems and the conversation should feel *very* different from the YouTube exercise.
