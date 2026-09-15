# Chapter 2.2 — CAP Theorem

This is one of the most commonly **misunderstood** topics in distributed-systems interviews. The important thing is not memorizing “pick two.” It's understanding **what a distributed system can and cannot guarantee when communication between nodes breaks**.

You already know the basic ideas of replication, network partitions, consistency, availability, and partial failure from earlier chapters. We'll build on those rather than re-teaching them.

---

## 2.2.1 — Why CAP Exists

Consider a replicated database:

```text
          Network
        /         \
     Node A       Node B
        \         /
         Application
```

Suppose Node A and Node B both contain copies of some data.

Normally:

```text
A ←→ B
```

They can communicate and keep their replicas synchronized.

Now imagine a network partition:

```text
        X
     A ──── B
```

Both nodes are still alive, but **A cannot communicate with B**.

This creates the fundamental CAP problem:

> What should the system do if it receives a request while the nodes cannot communicate?

If the system accepts the request on one side, how can it guarantee that the other side agrees?

If it refuses the request, the system remains consistent but isn't available for that operation.

CAP formalizes this trade-off.

---

# 2.2.2 — C: Consistency

In CAP, **Consistency** has a specific meaning.

A system is consistent if every read receives the **most recent successful write** (or an error), according to the system's consistency guarantee.

A simple example:

```text
Write:
balance = €100
```

After the write succeeds:

```text
Read → €100
```

A stale replica returning:

```text
Read → €50
```

would violate this strong consistency guarantee.

### Important: CAP consistency ≠ ACID consistency

This is a very important distinction.

You already learned **ACID Consistency** in Module 0.

ACID consistency means:

> A transaction takes the database from one valid state to another valid state according to defined constraints/rules.

CAP consistency is about:

> Whether reads observe the appropriate latest value across the distributed system.

They are **different concepts despite having the same name**.

This distinction is a common interview trap.

---

# 2.2.3 — A: Availability

In CAP, **Availability** means that every request to a non-failing node receives a response.

That response doesn't necessarily have to contain the newest data.

For example:

```text
Client → Node B
        ↓
      "Give me X"
        ↓
      response
```

An available system doesn't simply stop responding because another node cannot be reached.

This is different from saying:

> “The entire system is always operational.”

CAP availability is a specific distributed-systems guarantee about requests receiving responses.

---

# 2.2.4 — P: Partition Tolerance

A **network partition** occurs when nodes that normally communicate can no longer communicate.

For example:

```text
Node A  ←──── X ────→  Node B
```

A and B are both functioning.

The network path between them isn't.

**Partition tolerance** means the system continues operating despite this communication failure.

And this is where the popular:

> “Pick two of three”

explanation becomes misleading.

---

# 2.2.5 — Why Partition Tolerance Isn't Really Optional

In a distributed system, you generally cannot simply choose:

> “We'll sacrifice partition tolerance.”

Why?

Because network partitions are a real failure mode.

Networks can experience:

- broken links
- router failures
- switch failures
- dropped packets
- connectivity loss
- firewall problems
- routing failures
- cloud/network infrastructure failures

You don't get to tell the network:

> “Please don't partition because our architecture selected C + A.”

Therefore, for a genuinely distributed system, **P is effectively a requirement**.

The meaningful question becomes:

> **When a partition occurs, do we prioritize Consistency or Availability?**

That's the real CAP trade-off.

---

# 2.2.6 — The CAP Scenario

Let's make this concrete.

Suppose we have two replicas:

```text
        Application
        /         \
       v           v
    Node A       Node B
      X=10         X=10
```

Everything is synchronized.

Now a network partition occurs:

```text
        Application
        /         \
       v           v
    Node A   X   Node B
      X=10         X=10
```

The application sends a write:

```text
X = 20
```

Suppose the request reaches Node A.

Node A now has:

```text
X = 20
```

But Node B still has:

```text
X = 10
```

Node A cannot communicate with B.

Now the system has a choice.

---

# 2.2.7 — CP: Consistency + Partition Tolerance

A **CP system** chooses to preserve consistency when a partition occurs, potentially sacrificing availability.

Suppose Node B receives:

```text
Read X
```

But B cannot determine whether Node A has a newer value.

A CP-oriented system may respond:

```text
Error / unavailable
```

rather than returning potentially stale data.

The system effectively says:

> “I'd rather reject the operation than give you data that violates our consistency guarantee.”

So during a partition:

**Consistency wins over availability.**

### Why this can be desirable

For some operations, stale data is worse than an error.

Examples include:

- financial balances
- inventory allocation
- unique ownership
- leader election
- some reservation systems

If two sides of a partition independently make decisions, you could end up with:

```text
Inventory = 1

Customer A → successfully buys it
Customer B → also successfully buys it
```

The system remained available, but consistency was compromised.

---

# 2.2.8 — AP: Availability + Partition Tolerance

An **AP system** prioritizes availability during a partition.

Suppose Node A and Node B become disconnected.

Both may continue accepting operations.

For example:

```text
Node A:
X = 20

Node B:
X = 30
```

Both sides remain operational.

The system may later reconcile the conflicting state when connectivity returns.

The philosophy is:

> “Continue serving requests even if replicas temporarily disagree.”

This is useful when temporary inconsistency is acceptable.

Examples can include:

- social-media activity
- recommendation systems
- some analytics
- distributed counters
- certain shopping/catalog workloads

The exact consistency guarantees depend heavily on the specific database; **AP does not mean “eventually consistent” by definition**.

---

# 2.2.9 — What Actually Happens During a Partition?

This is the key interview concept.

Imagine:

```text
        PARTITION
      XXXXXXXXXXXX

      A            B
```

You cannot simultaneously guarantee both:

1. **Every successful read sees the latest globally agreed value**, and
2. **Every healthy side continues accepting requests independently**

because the two sides cannot communicate.

Suppose A accepts:

```text
X = 20
```

B accepts:

```text
X = 30
```

Without communication, neither side can know about the other's write.

If both continue successfully accepting writes, they cannot guarantee a single globally consistent ordering/state.

If one side stops accepting operations until coordination is restored, availability is sacrificed.

That's the CAP trade-off.

---

# 2.2.10 — CAP Is Not “Pick Any Two”

The popular diagram:

```text
        C
       / \
      /   \
     A --- P
```

often leads to:

> “Choose any two.”

That's **not** the useful interpretation.

The important scenario is:

> **A partition has occurred.**

When there is no partition, a system can potentially provide both strong consistency and high availability.

CAP tells us about the guarantees that can be maintained **in the presence of a partition**.

So the practical framing is:

```text
Normal operation:
C + A can potentially coexist.

During partition:
C ←→ A
     ↑
   trade-off
```

P isn't something you casually switch off in a distributed system.

---

# 2.2.11 — CP vs AP

| | CP | AP |
|---|---|---|
| Partition tolerance | Yes | Yes |
| Consistency during partition | Prioritized | May be weakened |
| Availability during partition | May be sacrificed | Prioritized |
| Typical behavior | Reject/block some operations | Continue serving requests |
| Best when | Incorrect/stale state is dangerous | Availability is more important than immediate consistency |

Don't interpret this as:

> CP = good, AP = bad.

Neither is inherently better.

The correct choice depends on the **business semantics of the operation**.

---

# 2.2.12 — A Crucial Nuance: CAP Is About Guarantees, Not Database Labels

It's tempting to say:

> “MongoDB is CP.”

or:

> “DynamoDB is AP.”

Those statements can be useful as simplified interview shorthand, but they're often too simplistic.

Real databases have:

- different operations
- different consistency levels
- different replication mechanisms
- quorum configurations
- leader/follower architectures
- tunable guarantees
- different behavior during failures

Therefore, when discussing a real database, ask:

> **What consistency and availability guarantees does this specific operation/configuration provide during a partition?**

CAP describes a fundamental constraint. It does **not** completely describe a database's behavior.

We'll study the actual consistency models and quorum mechanisms later.

---

# 2.2.13 — CAP vs Latency

CAP is specifically about **partition failures**, but distributed systems also have ordinary network latency.

These aren't the same thing.

### Normal latency

```text
A ───────→ B
   10 ms
```

Communication works; it just takes time.

### Partition

```text
A ── X ── B
```

Communication cannot reliably occur.

A system might experience increased latency without being partitioned.

This distinction matters because a system can sacrifice availability not only by returning errors, but sometimes by **waiting for coordination** until a timeout occurs.

That brings us back to the Chapter 2.1 lesson:

> A timeout doesn't necessarily mean the remote operation failed.

---

# 2.2.14 — CAP and Quorum Intuition

One useful concept to preview is **quorum**.

Suppose we have three replicas:

```text
A
B
C
```

A system might require agreement from a sufficient number of replicas before considering a write successful.

For three replicas, a majority is:

**2 of 3**

If one node fails:

```text
A ✓
B ✓
C ✗
```

two nodes can still form a majority.

But during a partition:

```text
     X
A ─────── B
          C
```

the side with a majority may be able to continue while the minority side cannot safely make the same decisions.

This is one of the mechanisms through which real systems preserve consistency.

**We will study quorums properly later**, when we cover replication and consistency. For now, the important point is that CAP trade-offs are implemented through concrete mechanisms such as coordination, acknowledgements, leaders, and quorums.

---

# 2.2.15 — What CAP Does *Not* Say

CAP does **not** say:

- NoSQL databases are inherently AP.
- SQL databases are inherently CP.
- Every distributed database must sacrifice availability all the time.
- Consistency and availability can never coexist.
- Eventual consistency is the same thing as AP.
- You literally select two checkboxes when designing a database.
- A partition is the same thing as a node failure.

The theorem is specifically about **what guarantees can simultaneously be maintained when a network partition occurs**.

---

# Interview Mental Model

If an interviewer asks:

> **“Explain CAP theorem.”**

A strong answer is:

> CAP says that a distributed system cannot simultaneously guarantee strong consistency and availability when a network partition occurs. Partition tolerance is effectively unavoidable in a real distributed system, so the practical trade-off during a partition is between consistency and availability. A CP system may reject or block operations to preserve consistency, while an AP system continues serving requests and may temporarily allow inconsistent or stale state.

And if they ask:

> **“Isn't CAP just pick two?”**

Say:

> No. That's an oversimplification. CAP is specifically about behavior during a network partition. When there is no partition, consistency and availability can coexist. The important question is what guarantee the system chooses to preserve when nodes cannot communicate.

---

## What you should retain from 2.2

The entire chapter boils down to one scenario:

**Two parts of a distributed system cannot communicate.**

Now ask:

**Do we continue accepting operations independently?**

→ We preserve **availability**, potentially sacrificing strong consistency.

**Do we refuse/block operations until we can safely coordinate?**

→ We preserve **consistency**, potentially sacrificing availability.

And that's the heart of CAP.

### Topics added to the original plan

I added only the pieces needed to make the CAP model technically complete:

- CAP consistency vs ACID consistency
- precise meaning of CAP availability
- why P is effectively unavoidable
- concrete partition scenario
- why “pick two” is misleading
- CAP vs ordinary latency
- quorum intuition
- why real database labels like “CP” and “AP” are simplifications
- what CAP does **not** say

We should **not** go deeper into quorum algorithms or consistency models yet; those belong later in the syllabus.