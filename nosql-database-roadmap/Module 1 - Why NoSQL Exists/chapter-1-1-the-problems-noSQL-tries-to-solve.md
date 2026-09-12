# MODULE 1 — Why NoSQL Exists

## Chapter 1.1 — The Problems NoSQL Tries to Solve

The first thing to understand is that NoSQL did **not** emerge because relational databases suddenly became bad.

It emerged because certain workloads created combinations of requirements that were difficult, expensive, or operationally inconvenient to satisfy with traditional relational architectures.

The important word is **combination**.

A system might simultaneously require:

- enormous amounts of data
- enormous request volume
- low latency
- high availability
- horizontal scaling
- global distribution
- flexible data structures
- rapid evolution
- distributed operation

Each requirement individually can often be handled by relational technology.

The challenge is satisfying many of them **simultaneously and at very large scale**.

---

# 1. Massive datasets

The first problem is simply **scale of stored data**.

A database may grow from:

```text
GB
→ TB
→ tens of TB
→ hundreds of TB
→ PB
```

At sufficiently large sizes, storing everything on one machine becomes increasingly problematic.

There are several dimensions to this problem.

### Storage capacity

One machine has finite:

- disk capacity
- memory
- CPU
- I/O bandwidth

If the dataset exceeds practical capacity, it has to be distributed.

---

### Data access

Even if a single machine technically has enough storage, accessing an enormous dataset efficiently is another problem.

You need to consider:

- how data is located
- how much data must be scanned
- how indexes are distributed
- how much memory can be used for caching
- how much I/O the machine can sustain

Simply having enough disk space does not mean the database can efficiently serve the workload.

---

### Maintenance and operational concerns

Large datasets also affect:

- backups
- restores
- replication
- migrations
- reindexing
- hardware replacement
- failure recovery

The larger the dataset, the more expensive many of these operations become.

This motivates **partitioning data across multiple machines**.

We'll study partitioning in detail later.

---

# 2. Massive request volumes

Dataset size and request volume are **different scaling problems**.

You could have:

```text
Very large data
+
relatively few requests
```

or:

```text
Moderate data
+
enormous request volume
```

The second can be just as challenging.

A database may need to handle:

- thousands of requests/second
- tens of thousands
- hundreds of thousands
- potentially millions of operations/second in specialized systems

The important metric isn't just the average request rate.

Production systems also care about:

- peak traffic
- latency under load
- read/write ratio
- concurrency
- burst behavior
- consistency requirements

A database that handles 50,000 operations/second under ideal conditions isn't necessarily capable of maintaining acceptable latency when those operations arrive concurrently.

---

# 3. Why one machine eventually becomes a bottleneck

Suppose a database is running on one server:

```text
                 Database
                    │
        ┌───────────┼───────────┐
        │           │           │
       CPU         RAM         I/O
```

Increasing workload eventually saturates one or more resources.

You can attempt to solve this through **vertical scaling**:

```text
more CPU
more RAM
faster storage
more I/O bandwidth
```

But vertical scaling has limits.

Eventually:

- hardware has finite capacity
- increasingly powerful hardware becomes expensive
- some hardware upgrades require downtime
- a single machine remains a failure domain
- capacity growth becomes constrained by the machine

So we reach the fundamental alternative:

> **Horizontal scaling.**

Instead of making one machine bigger:

```text
        One very large server
```

we use:

```text
       ┌────────┐
       │ Node 1 │
       ├────────┤
       │ Node 2 │
       ├────────┤
       │ Node 3 │
       ├────────┤
       │ Node 4 │
       └────────┘
```

Now the workload and/or data can be distributed.

This is one of the central ideas behind many NoSQL systems.

---

# 4. Globally distributed applications

Modern applications aren't necessarily serving users from one geographic location.

Users may be distributed across:

```text
Europe
North America
Asia
Australia
...
```

If the application's database exists only in one geographic region, users far away may experience additional network latency.

This creates a desire to place data closer to users.

Conceptually:

```text
Europe users
     ↓
Europe database region

US users
     ↓
US database region

Asia users
     ↓
Asia database region
```

But now we have introduced a much harder problem:

> **There are multiple copies or partitions of data in different geographic locations.**

The system has to deal with:

- replication
- synchronization
- network latency
- network failures
- conflicting writes
- consistency
- regional failures
- routing requests to appropriate regions

This is fundamentally a **distributed-systems problem**.

---

# 5. High availability

Availability means that the system remains capable of serving requests when failures occur.

A single database server is inherently vulnerable:

```text
Application
     │
     ▼
 Database
     X
   failure
```

If that machine fails, the database may become unavailable.

A distributed architecture can instead have multiple nodes:

```text
             Database
          ┌────┼────┐
          ▼    ▼    ▼
        Node A Node B Node C
```

If one node fails:

```text
Node A → failed

Node B
Node C
   ↓
continue serving
```

This requires **replication** or other mechanisms for maintaining multiple usable copies/partitions of data.

But high availability introduces another trade-off:

> The more independent nodes you have, the more coordination you may need.

And coordination itself can increase:

- latency
- network traffic
- complexity
- failure modes

This tension is central to distributed database design.

---

# 6. Low-latency access

High throughput and low latency are **not the same thing**.

A system might process a huge number of requests overall while individual requests still take too long.

For example, an application might care about:

```text
p50 latency
p95 latency
p99 latency
```
(more on them at the end)

rather than merely:

```text
requests/second
```

Why?

Because users experience individual requests.

A database designed for a particular access pattern can potentially make the critical operation very efficient by:

- locating data directly
- minimizing scanning
- minimizing JOINs
- minimizing network hops
- keeping frequently accessed data close to compute
- choosing an appropriate partition key
- using appropriate indexes

This is closely connected to the NoSQL principle:

> **Design the data model around the application's access patterns.**

That principle will become extremely important when we reach NoSQL data modeling.

---

# 7. Flexible schemas

Traditional relational systems generally expect a clearly defined schema.

Conceptually:

```text
Users
------------------------
Id
Name
Email
Age
```

Every row follows that structural model.

NoSQL systems—particularly document databases—can allow records to have different structures.

For example, the logical model can evolve from:

```text
User
{
    id,
    name
}
```

to:

```text
User
{
    id,
    name,
    preferences,
    addresses,
    ...
}
```

without necessarily requiring a traditional table-altering migration before the new structure can be stored.

This is called **schema flexibility**.

But remember what we established in Module 0:

> **Flexible schema does not mean no schema.**

The application still needs to understand what fields mean and which structures are valid.

A schema can exist at:

- application level
- validation layer
- database level
- serialization layer
- API contract level

The difference is often **where and how strictly the schema is enforced**.

---

# 8. Rapid schema evolution

Flexible schemas become particularly valuable when the data model changes frequently.

Suppose a product evolves rapidly.

Its data might go through:

```text
Version 1
A B C

Version 2
A B C D

Version 3
A B C D E F
```

In a highly structured relational environment, schema evolution can involve:

- ALTER TABLE operations
- migrations
- constraint changes
- index changes
- application/database version coordination
- deployment sequencing
- compatibility with old application versions

None of these are inherently bad.

In fact, controlled relational migrations are extremely valuable.

But when the system changes rapidly and different versions of the application may coexist, schema flexibility can reduce some of that friction.

This is especially useful when:

- product requirements change frequently
- data structures vary significantly between records
- new attributes are introduced regularly
- different application versions coexist
- development velocity is important

---

# 9. Distributed workloads

This is broader than simply having a large database.

A **distributed workload** means the processing itself may need to happen across multiple machines.

Instead of:

```text
              One machine
                  │
          ┌───────┴───────┐
          │               │
        data            queries
```

we might have:

```text
                 Workload
              ┌────┼────┐
              ▼    ▼    ▼
            Node A Node B Node C
```

Now the database must decide:

- where data belongs
- where requests should go
- how data is replicated
- how operations are coordinated
- how failures are detected
- how nodes recover
- how nodes are added/removed
- how workload is balanced

This is where database design begins to overlap heavily with distributed-systems engineering.

---

# 10. The important connection between these problems

These requirements aren't independent.

They interact.

For example:

```text
Massive dataset
      ↓
Need multiple machines
      ↓
Partition data
      ↓
Need replication
      ↓
Need failure handling
      ↓
Distributed system
      ↓
Network communication
      ↓
Consistency/latency trade-offs
```

Or:

```text
Massive request volume
      ↓
One server becomes bottleneck
      ↓
Horizontal scaling
      ↓
Multiple nodes
      ↓
Requests must be distributed
      ↓
Data must be distributed
```

And:

```text
Global users
      ↓
Need data closer to users
      ↓
Multiple geographic regions
      ↓
Replication / distributed writes
      ↓
Consistency becomes harder
```

This is why NoSQL isn't merely about replacing:

```text
table → document
```

The deeper story is:

> **NoSQL systems were designed to make particular trade-offs around scale, distribution, availability, data modeling, and performance.**

---

# 11. The central trade-off

At small scale, you can often optimize for simplicity:

```text
One database
One location
Strong consistency
Complex queries
JOINs
Transactions
```

As requirements increase, you may need:

```text
Multiple machines
Multiple regions
Partitioning
Replication
High availability
Very high throughput
Low latency
```

But distributing the system introduces costs:

```text
Network
Coordination
Consistency problems
Failure handling
Operational complexity
Data duplication
More complicated modeling
```

So NoSQL isn't:

> "A way to get all the benefits for free."

It is largely about choosing **different points in the design trade-off space**.

---

# 12. What NoSQL actually buys you

Depending on the specific NoSQL technology, you may get architectural characteristics such as:

- easier horizontal scaling
- partition-oriented architecture
- high throughput for particular workloads
- low-latency access for predictable access patterns
- flexible document structures
- built-in replication
- geographically distributed deployment
- reduced need for cross-entity JOINs
- data models optimized for particular access patterns

But these advantages generally come with corresponding costs:

- denormalization
- duplicated data
- more application responsibility
- weaker or different consistency guarantees
- limited query flexibility in some systems
- more complicated distributed operations
- operational complexity

And different NoSQL databases make **different** trade-offs.

That's why we'll study them separately rather than treating "NoSQL" as one technology.

---

# 13. A crucial distinction: "necessary" doesn't mean "mandatory"

I want to correct one phrase from the original syllabus:

> "Why NoSQL databases became necessary"

Don't interpret that literally.

NoSQL was not universally necessary.

For many applications today, a relational database remains the best choice.

The better interpretation is:

> **Why did workloads emerge for which specialized non-relational database architectures became useful?**

That is the question we're actually answering.

---

# 14. The seven problems in one mental model

For your notes, I'd summarize Chapter 1.1 like this:

```text
                 WHY NOSQL?
                     │
       ┌─────────────┼─────────────┐
       │             │             │
 Massive data   Massive traffic   Global users
       │             │             │
       └─────────────┼─────────────┘
                     ▼
             Horizontal scaling
                     │
                     ▼
             Distributed system
                     │
          ┌──────────┼──────────┐
          ▼          ▼          ▼
     Availability   Low       Flexible
                   latency     schemas
                     │
                     ▼
             Different trade-offs
```

The key progression is:

**Scale → distribution → new trade-offs → specialized database architectures.**

---

# Interview-level understanding

You should now be able to answer:

### "Why did NoSQL emerge?"

A strong answer would be:

> NoSQL emerged to address workloads where very large data volumes, high request rates, horizontal scaling, global distribution, low latency, high availability, and flexible data models became important simultaneously. Relational databases can address many of these requirements, but distributed relational architectures can involve significant coordination around transactions, JOINs, and consistency. NoSQL systems were designed to make different trade-offs for specific workload patterns.

That's much better than:

> "SQL doesn't scale, so NoSQL was invented."

---

### One distinction to keep in mind

We're deliberately **not yet diving into CAP, partitioning algorithms, replication strategies, consistency models, or sharding mechanics**.

Those are the *solutions and distributed-system mechanisms* we'll study next.

This chapter is answering only:

> **What problems are we trying to solve?**

The next chapters will answer:

> **How do distributed databases actually solve them?**

## What are p50, p95, etc.
They are **percentiles of observed latency**. They tell you how fast a request is for a typical user and, importantly, how bad things get for slower requests.

Suppose you measure the response time of **1,000 database requests** and sort them from fastest to slowest:

```text
fastest ─────────────────────────────────────────── slowest
  1       2       3       ...       500   ...   950 ... 1000
```

### P50 — 50th percentile

The latency at which **50% of requests are faster or equal** and 50% are slower.

If:

```text
P50 = 20 ms
```

then roughly half of your requests complete in **20 ms or less**.

This is essentially the **median**.

---

### P95 — 95th percentile

If:

```text
P95 = 80 ms
```

then **95% of requests complete in 80 ms or less**, while the slowest 5% take longer.

```text
95%                    5%
───────────────────────┬────────────
       ≤ 80 ms         │ > 80 ms
```

---

### P99 — 99th percentile

If:

```text
P99 = 300 ms
```

then **99% of requests complete in 300 ms or less**, while the slowest 1% take longer.

```text
99%                                      1%
─────────────────────────────────────────┬──────
              ≤ 300 ms                  │ > 300 ms
```

---

## Why not just use the average?

Suppose you have these latencies:

```text
10, 10, 10, 10, 10, 10, 10, 10, 10, 1000 ms
```

The average is:

```text
109 ms
```

But that's misleading.

**90% of requests took only 10 ms.**

The 1000 ms request is an outlier that heavily affects the average.

Percentiles expose this distribution much better.

---

## Why P95/P99 matter in production

Imagine an API with:

```text
P50 = 20 ms
P95 = 40 ms
P99 = 500 ms
```

At first glance, the API looks excellent:

> "Typical requests take only 20 ms."

But **1% of requests take 500 ms or more**.

At:

```text
100 requests/sec
```

that's roughly:

```text
1 request/sec
```

experiencing that high latency.

At:

```text
100,000 requests/sec
```

that's potentially:

```text
1,000 requests/sec
```

in the slow tail.

So production engineers care a lot about the **tail latency** represented by P95, P99, and sometimes P99.9.

---

### The mental model

```text
P50   → typical request
P95   → slower tail
P99   → very slow tail
P99.9 → extreme tail
```

And when we talk about **low-latency NoSQL systems**, saying "the database is fast" isn't enough. We want to know something like:

```text
Read latency:
P50  = 3 ms
P95  = 7 ms
P99  = 15 ms
```

because that tells us much more about the actual behavior under load.
