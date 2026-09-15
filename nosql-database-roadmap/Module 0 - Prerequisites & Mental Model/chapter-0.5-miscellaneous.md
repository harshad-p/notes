# Chapter 0.5 — CPU, I/O, Memory and Database Latency

We mentioned these, but they deserve an explicit mental model.

When a database executes something, several resources can become bottlenecks.

```text
                 Database
                    │
        ┌───────────┼───────────┐
        ▼           ▼           ▼
       CPU         RAM         Storage
        │           │           │
   computation    cache       SSD/HDD
   planning       buffers     reads/writes
   execution
```

### CPU

CPU is used for things such as:

- query execution
- filtering
- sorting
- JOIN processing
- aggregation
- parsing
- compression/decompression
- encryption
- index processing

If CPU is saturated, adding an SSD won't necessarily help.

---

### Memory

RAM is dramatically faster than persistent storage.

A database can keep frequently accessed data in memory:

```text
Query
  ↓
RAM/cache
  ↓
found?
 ├── yes → return
 └── no
       ↓
     SSD
       ↓
     RAM/cache
       ↓
     return
```

This is one reason the same query can behave very differently between:

```text
cold cache
```

and:

```text
warm cache
```

---

### I/O

I/O means input/output operations, commonly involving storage or network.

For a database:

```text
SSD read
SSD write
network read
network write
```

can all contribute to latency.

A query might therefore be:

```text
CPU-bound
I/O-bound
memory-bound
network-bound
```

This is an important production concept.

---

## Why database latency isn't simply "disk speed"

Suppose an operation takes 20 ms.

That 20 ms might consist of:

```text
Network       2 ms
Query parsing 0.2 ms
Planning      1 ms
Lock waiting  5 ms
CPU           3 ms
Storage I/O   4 ms
Result        1 ms
Network       3.8 ms
---------------------
Total        20 ms
```

The numbers are illustrative, but the principle matters.

**A slow query does not automatically mean slow disk.**

And sometimes the database doesn't touch disk at all because the required pages are already cached.

---

# Chapter 0.6 — Transactions and ACID

A **transaction** is a logical unit of work that should satisfy certain guarantees.

Consider transferring €100:

```text
Account A: -100
Account B: +100
```

We don't want:

```text
A: -100
B: unchanged
```

because the system crashed between the two operations.

A transaction lets the database treat the operations as one logical unit.

---

## ACID

### A — Atomicity

All operations in the transaction happen, or none of them do.

```text
Transfer
 ├── debit A
 └── credit B
```

Either:

```text
both succeed
```

or:

```text
both are rolled back
```

---

### C — Consistency

The transaction should take the database from one valid state to another valid state, respecting defined constraints and rules.

For example:

```text
Account balance >= 0
```

if that is an enforced business/database constraint.

Important: **database consistency here is not the same thing as "all replicas always have the latest value."**

That distinction becomes extremely important in distributed NoSQL systems.

---

### I — Isolation

Concurrent transactions shouldn't improperly interfere with each other.

Imagine:

```text
Transaction A
      │
      ├── reads balance
      │
Transaction B
      │
      ├── modifies balance
      │
Transaction A
      └── continues...
```

The database needs rules governing what A can see from B.

This leads to isolation levels such as:

```text
Read Uncommitted
Read Committed
Repeatable Read
Serializable
```

We'll go much deeper into this later.

---

### D — Durability

Once the database says:

> COMMIT succeeded

the committed data should survive a crash.

This can involve mechanisms such as:

- WAL
- journaling
- flushing
- replication
- recovery procedures

We already introduced durability earlier.

---

### ACID mental model

```text
Transaction
    │
    ├── Atomicity
    ├── Consistency
    ├── Isolation
    └── Durability
```

Don't memorize the acronym alone. Understand the problem each property solves.

---

# Chapter 0.7 — B-Trees, Indexes and Read Replicas

We discussed indexes, but not **what an important traditional database index actually looks like**.

## B-tree

A B-tree is a tree-like data structure designed to efficiently locate sorted data, particularly when working with storage systems.

Conceptually:

```text
                 [50]
                /    \
             [20]    [80]
            /   \    /   \
          [10] [30] [70] [90]
```

Instead of scanning:

```text
1
2
3
4
5
...
50 million
```

the database can navigate the tree toward the desired value.

That's why an index can turn something conceptually like:

```text
O(n)
```

into approximately:

```text
O(log n)
```

for locating values, though real database performance involves many more factors.

Modern databases use variations and other index structures too. Don't equate **"database index" = "B-tree"** universally.

We'll later cover:

- B-tree/B+tree concepts
- hash indexes
- compound indexes
- covering indexes
- selectivity
- index-only scans
- clustered/non-clustered indexes
- database-specific index structures

---

## Read replicas

Suppose we have:

```text
                Application
                    │
              ┌─────┴─────┐
              ▼           ▼
           Primary      Replica
```

The primary handles writes:

```text
INSERT
UPDATE
DELETE
```

and changes are replicated to the replica.

The application can potentially send reads to replicas:

```text
Writes ──────→ Primary
Reads  ───────→ Replica
```

This can increase **read capacity**.

But there is a critical catch:

### Replication isn't necessarily instantaneous.

For example:

```text
Primary
  balance = 500
      │
      │ replication
      ▼
Replica
  balance = 400
```

For a short period, a read from the replica may return stale data.

This introduces a fundamental distributed-systems trade-off:

> **Do I want more read scalability, or do I require the freshest possible data?**

That question becomes much more important in NoSQL.

---

# Chapter 0.8 — Relational Scaling and the "Why Can't We Just Shard SQL?"

This deserves a proper answer because it is one of your planned interview questions.

First:

> **Relational databases absolutely can be horizontally scaled.**

So the interview question is slightly misleading.

The real question is:

> Why isn't horizontal scaling as straightforward for a traditional relational workload as simply adding another independent server?

Because relational databases often have **cross-entity relationships and operations**.

Imagine:

```text
Customers
Orders
OrderItems
Products
Payments
```

Now imagine distributing them:

```text
Node A → Customers
Node B → Orders
Node C → Products
Node D → Payments
```

A single business operation might require:

```text
Customer
   ↓
Order
   ↓
OrderItems
   ↓
Products
   ↓
Payment
```

Now you potentially have:

- distributed JOINs
- distributed transactions
- network communication
- coordination
- distributed locking
- consistency concerns
- failure handling

None of these are impossible.

They are simply **expensive and complicated**.

This is why many NoSQL systems make a deliberate trade:

> Design the data around the queries so that common operations can be performed within a partition/node, minimizing cross-node coordination.

That's a much more precise explanation than:

> "SQL can't scale horizontally."

---

# Chapter 0.9 — Flexible Schemas, Relational vs. NoSQL

One more missing topic: **schema flexibility**.

Consider a relational table:

```text
Users
-------------------------
Id | Name | Email | Age
```

The structure is explicitly defined.

A document database might have:

```json
{
  "id": 1,
  "name": "Alice",
  "email": "alice@example.com"
}
```

while another document could contain:

```json
{
  "id": 2,
  "name": "Bob",
  "email": "bob@example.com",
  "age": 37,
  "preferences": {
    "language": "en"
  }
}
```

This is commonly called **schema flexibility** or **schema-on-read**.

But:

> Flexible schema ≠ no schema.

Your application still has assumptions about what data looks like.

If half your documents contain:

```text
email
```

and the other half don't, your application needs to handle that.

You can also enforce schemas at the application or database level.

So don't fall into:

> SQL = schema  
> NoSQL = no schema

A better model is:

```text
Relational
→ schema is generally explicit and strongly structured

NoSQL
→ schema can often be more flexible
→ structure is frequently driven by application access patterns
```

---

# Chapter 0.10 — Module 0 Interview Checkpoint

Now let's explicitly close the interview topics that were part of the original plan.

### Q1. Why can't relational databases simply scale horizontally?

**Good senior answer:**

> They can, but horizontal scaling becomes more complex when workloads require cross-partition JOINs, transactions, and coordination. Distributed relational systems exist, but distributing relational operations introduces network communication, coordination, and consistency costs.

---

### Q2. Why are JOINs expensive in distributed systems?

Because the participating data may reside on different nodes.

```text
Node A ──network── Node B
```

The system may need to move data or coordinate operations across the network.

---

### Q3. Is NoSQL actually "non-relational"?

Not necessarily in the simplistic sense.

"NoSQL" is an umbrella term for several database models that generally don't use the traditional relational model as their primary data model.

Some NoSQL databases may support relationships or SQL-like querying.

So:

```text
NoSQL ≠ one specific database model
```

---

### Q4. Is NoSQL always faster than SQL?

**Absolutely not.**

Performance depends on:

- workload
- data model
- query pattern
- indexes
- data distribution
- consistency requirements
- hardware
- implementation
- network
- workload volume

A badly modeled MongoDB database can be much slower than a well-designed PostgreSQL database.

---

### Q5. Why did NoSQL databases become necessary?

Better answer than "SQL couldn't scale":

> New workloads increasingly required massive horizontal scaling, high write/read throughput, flexible data models, globally distributed systems, and different consistency/availability trade-offs. Existing relational databases could solve many of these problems, but specialized NoSQL systems were designed around particular distributed workloads and could make different trade-offs.

---

### One remaining tiny gap

In **0.2**, I said "result construction" but didn't actually explain it. We can consider that covered at a basic level by the query-processing pipeline, but since you specifically want a syllabus where we don't quietly skip things, I'll explicitly add this concept to our notes:

```text
Execution
   ↓
Rows/documents found
   ↓
Filtering/projection/aggregation
   ↓
Result set constructed
   ↓
Serialization
   ↓
Network response
   ↓
Application
```

For a large result, **constructing, serializing, transferring, and deserializing the result can itself become expensive**. This is one reason `SELECT *` and returning thousands/millions of records can be problematic even when the database finds them efficiently.

**With that, I consider Module 0 complete against the original plan.**
