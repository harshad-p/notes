# Chapter 2.3 — BASE

CAP gave us the **distributed-systems constraint**:

> During a network partition, you cannot simultaneously guarantee strong consistency and availability.

BASE gives us a way of thinking about systems that **favor availability and allow state to converge over time**.

The three letters stand for:

- **BA — Basically Available**
- **S — Soft State**
- **E — Eventual Consistency**

But unlike ACID, BASE isn't a precise transactional specification. It's more of an **architectural philosophy** describing systems that relax some traditional consistency guarantees in exchange for availability, scalability, and reduced coordination.

---

## 2.3.1 — Why BASE Emerged

Traditional relational systems were often designed around a strong model:

```text
Write
  ↓
Maintain a consistent database state
  ↓
Commit
  ↓
Read the committed state
```

That model is extremely useful.

But imagine a globally distributed system with replicas in:

```text
Berlin
New York
Tokyo
```

A write in Berlin may need to be propagated to the other replicas.

If we require all replicas to agree before confirming the write, we introduce:

- network latency
- coordination
- dependency on remote nodes
- reduced availability when communication fails

For some applications, waiting for global agreement isn't worth the cost.

Instead, the system might say:

> “Accept the write locally and let the replicas catch up.”

That gives us a different set of trade-offs:

```text
       Strong coordination
              ↑
              │
      consistency
              │
              │
              ↓
       availability
       + scalability
       + lower coordination
```

BASE emerged from this general approach.

---

# 2.3.2 — Basically Available

**Basically Available** means that the system is designed to remain available even when some parts of the system experience failures or degraded conditions.

It does **not** mean:

> “The system is guaranteed to always respond.”

That's too strong.

Instead, the system attempts to provide a useful response despite failures.

For example, imagine a shopping application where the recommendation service is temporarily unavailable.

The application might still return:

```text
Product page
+ price
+ description
+ inventory
```

while omitting:

```text
"Customers who bought this also bought..."
```

The system remains **basically available** even though some functionality is degraded.

This idea is closely related to the broader distributed-systems principle of **graceful degradation**.

### Important distinction

Availability in CAP and “Basically Available” in BASE are related but **not identical terms**.

CAP availability is a formal property within the theorem.

BASE's “Basically Available” is an architectural philosophy: keep the system serving useful requests despite failures or temporary inconsistencies.

---

# 2.3.3 — Soft State

This is probably the least intuitive part of BASE.

In a traditional strongly consistent system, you might think:

> “If I write some state, that state remains the authoritative state until another transaction changes it.”

With **soft state**, the state of a distributed system may change over time **even without a new user action**, because replicas are synchronizing or converging.

For example:

```text
t0:
Node A → X = 10
Node B → X = 10

t1:
A receives update → X = 20

t2:
A → X = 20
B → X = 10

t3:
replication occurs

t4:
A → X = 20
B → X = 20
```

At `t2`, the system's state isn't yet globally settled.

Node B can change from:

```text
X = 10
```

to:

```text
X = 20
```

without a new user request directly modifying B.

That's the intuition behind **soft state**.

The state can evolve as the system processes replication, synchronization, expiration, conflict resolution, or other background activity.

---

# 2.3.4 — Eventual Consistency

This is the most recognizable part of BASE.

**Eventual consistency** means that if no new updates occur and the system can communicate normally, replicas will **eventually converge to the same value**.

For example:

```text
Initial:

A → X = 10
B → X = 10


Write X = 20 to A:

A → X = 20
B → X = 10


Replication:

A → X = 20
B → X = 20
```

During the convergence period, different clients might observe different values.

After convergence, assuming no further updates:

```text
A → X = 20
B → X = 20
```

### What “eventual” means

It means:

> **Not necessarily immediately, but eventually under the required conditions.**

It does **not** mean:

> “The database is randomly inconsistent.”

There is an intended path toward convergence.

---

# 2.3.5 — Eventual Consistency Doesn't Mean No Guarantees

This is an important nuance.

An eventually consistent database can still provide many guarantees.

For example, it may guarantee:

- writes aren't silently lost
- replicas eventually converge
- a particular key has deterministic conflict resolution
- writes from the same client are observed in order
- reads from a particular replica are monotonic

The exact guarantees depend on the database.

So:

> **Eventual consistency is a consistency model, not an absence of consistency.**

This is why saying:

> “BASE means bad consistency”

is incorrect.

The system has deliberately chosen **weaker timing guarantees about when all replicas reflect the same state**.

---

# 2.3.6 — A Practical Example

Consider a social-media application.

You update your profile:

```text
name = "Harshad"
```

Suppose your profile has replicas in several regions.

Immediately after the write:

```text
Europe replica → Harshad
US replica     → old value
Asia replica   → old value
```

For a short period, different users might see different versions.

Eventually:

```text
Europe → Harshad
US     → Harshad
Asia   → Harshad
```

Would this be acceptable?

For many profile or social-feed scenarios, yes.

Would it be acceptable for every workload?

No.

Consider:

```text
Bank balance → €10
```

If one replica says:

```text
€10
```

and another says:

```text
€100
```

during a transaction, that could be unacceptable.

So consistency isn't inherently “good” or “bad.”

It is a **business requirement**.

---

# 2.3.7 — ACID vs BASE

It's tempting to treat these as opposing database technologies.

They aren't.

They're different ways of thinking about **system behavior and guarantees**.

| | ACID | BASE |
|---|---|---|
| Primary emphasis | Strong transactional guarantees | Availability and distributed operation |
| State | Transactionally controlled | Can temporarily change/converge |
| Consistency | Stronger transactional consistency | May allow temporary inconsistency |
| Coordination | Often more coordination | Attempts to reduce coordination |
| Replicas | May require stronger synchronization | Can converge asynchronously |
| Failure behavior | May reject/block operations | Often continues with degraded/temporary inconsistency |
| Typical use | Financial transactions, inventory, relational workflows | Large distributed systems, feeds, caches, some globally distributed workloads |

But don't conclude:

> SQL = ACID  
> NoSQL = BASE

That's false.

Modern databases often provide **both ACID transactions and distributed/replicated behavior**.

Likewise, NoSQL databases can provide strong consistency.

BASE describes a **design philosophy**, not a mandatory property of every NoSQL database.

---

# 2.3.8 — Why Would Anyone Accept Temporary Inconsistency?

Because strong consistency has a **cost**.

Suppose replicas are geographically distributed:

```text
Berlin ←→ New York ←→ Tokyo
```

If every write requires agreement across all regions, a write may have to wait for network communication.

That means:

- higher latency
- more coordination
- greater sensitivity to network failures
- potentially reduced availability

But if we allow:

```text
Write locally
      ↓
Return success
      ↓
Replicate asynchronously
      ↓
Eventually converge
```

we can potentially achieve:

- lower write latency
- higher availability
- better geographic scalability
- less cross-region coordination

The trade-off is that some reads may temporarily observe older state.

That's the fundamental reason BASE-style designs exist.

---

# 2.3.9 — BASE Isn't Simply "Bad Consistency"

This deserves explicit emphasis.

Imagine two systems.

### System A

Guarantees:

> Every read immediately reflects the latest globally committed value.

But during a network partition it must reject many requests.

### System B

Guarantees:

> Reads may temporarily return an older value, but the system continues serving requests and replicas converge.

Neither is objectively better.

It depends on the business requirement.

For:

**Payment authorization**

you probably care more about correctness than continuing with uncertain state.

For:

**Social-media likes**

you may prefer:

> “Let people continue liking posts and reconcile the count later.”

The architecture should reflect the **cost of inconsistency**.

---

# 2.3.10 — BASE and CAP

The connection to the previous chapter is important.

CAP says:

> During a partition, strong consistency and availability cannot both be guaranteed.

BASE represents one common architectural response:

> Favor availability and allow temporary inconsistency/convergence.

Conceptually:

```text
CAP:

Partition occurs
      ↓
Consistency ←→ Availability


BASE-style approach:

Partition/degraded communication
      ↓
Continue serving requests
      ↓
Allow temporary divergence
      ↓
Replication/reconciliation
      ↓
Convergence
```

So BASE is particularly associated with systems that make the **AP-oriented trade-off**, although BASE itself is not a formal synonym for AP.

---

# 2.3.11 — A Subtle Point: Eventual Consistency Has Preconditions

The word **eventual** matters.

The guarantee generally assumes conditions such as:

- updates stop, and
- the system can eventually communicate, and
- the replication/reconciliation mechanism continues functioning.

If the system remains permanently partitioned, you cannot simply say:

> “Eventually everything will converge.”

There may never be an opportunity to synchronize.

So eventual consistency isn't a promise that replicas magically become equal regardless of failures.

It's a convergence guarantee **under appropriate conditions**.

---

# 2.3.12 — Interview Mental Model

If asked:

> **“What is BASE?”**

A good answer is:

> BASE is an architectural approach for distributed systems that emphasizes availability and scalability by relaxing some immediate consistency guarantees. Basically Available means the system tries to continue serving useful requests despite failures or degraded conditions. Soft State means distributed state can change as replicas synchronize, even without a new user operation. Eventual Consistency means that, assuming updates stop and communication and replication recover, replicas eventually converge.

If asked:

> **“Is BASE just bad consistency?”**

Say:

> No. It's a deliberate trade-off. The system accepts temporary inconsistency in exchange for properties such as availability, lower coordination, and potentially lower latency. Whether that's appropriate depends on the business consequences of stale or conflicting data.

---

## The Mental Model

The progression from the last two chapters is:

**CAP**

> A partition occurs. What guarantee do we preserve?

**BASE**

> If we favor availability, we can allow replicas to temporarily disagree and converge later.

So:

```text
Network partition / distributed operation
                ↓
      coordination is expensive
                ↓
       don't require everything
       to agree immediately
                ↓
        continue serving requests
                ↓
       replicas may diverge
                ↓
       synchronize/reconcile
                ↓
          converge eventually
```

### One final distinction to remember

**ACID** asks largely:

> “Can this transaction be executed with strong transactional guarantees?”

**BASE** asks more broadly:

> “Can this distributed system remain useful and available while allowing its state to converge over time?”

Neither philosophy is universally correct. **The business semantics of the data determine which trade-off is acceptable.**

### Topics I added

I added four small areas that weren't explicitly in your original plan because they prevent common misunderstandings:

- CAP vs BASE: how the two concepts relate
- Eventual consistency's preconditions
- BASE guarantees vs “no guarantees”
- Graceful degradation as an intuition for “Basically Available”

I deliberately **did not go into specific consistency models, conflict resolution algorithms, or CRDTs** yet. Those belong later and would dilute this chapter.