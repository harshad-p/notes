# Chapter 1.4 — NoSQL Myths

The term *NoSQL* tends to produce a collection of oversimplified statements. Most of them contain a small amount of truth, but become wrong when treated as universal rules.

The goal here is to replace those statements with more accurate mental models.

---

## 1. “NoSQL means no schema”

**False.**

NoSQL generally means that the database does not require the same kind of rigid relational schema—not that the data has no structure.

Consider a document database:

```json
{
  "id": 42,
  "name": "Alice",
  "address": {
    "city": "Berlin"
  }
}
```

The document clearly has a structure.

The distinction is that a document database may allow another document to have a different structure:

```json
{
  "id": 43,
  "name": "Bob",
  "phone": "+49..."
}
```

Whether this is actually *allowed* and whether it is *desirable* depends on the database and application.

### Schema can exist at multiple levels

Even if the database doesn't enforce a rigid schema, you can still have:

- application-level schemas
- validation rules
- API contracts
- serialization models
- domain models
- indexes that assume particular fields exist

So:

> **Schema flexibility ≠ absence of schema.**

A better description is:

> NoSQL databases often move some schema enforcement away from a rigid database-wide relational schema and allow greater structural flexibility.

And importantly, **schema flexibility is not automatically a benefit**. Uncontrolled schema variation can make querying, validation, migrations, and maintenance harder.

---

# 2. “NoSQL doesn't support transactions”

**False.**

Many NoSQL databases support transactions.

The important question is:

> **What transactional guarantees does this particular database provide, and at what scope?**

There are several dimensions:

- Can a transaction modify one item?
- Multiple documents?
- Multiple partitions?
- Multiple collections/tables?
- Multiple nodes?
- What isolation guarantees exist?
- What happens during failures?

Different NoSQL databases make different choices.

For example, modern MongoDB supports multi-document transactions. DynamoDB also provides transactional operations.

However, this does **not** mean every NoSQL database behaves like a traditional relational database transactionally.

A distributed database may deliberately limit transaction scope because coordinating a transaction across partitions or nodes is expensive.

### The important distinction

Don't think:

```text
SQL      → transactions
NoSQL    → no transactions
```

Think:

```text
SQL/NoSQL
    ↓
Specific database
    ↓
Specific transaction model
    ↓
Specific guarantees + scope + performance characteristics
```

This distinction will become important when we reach our dedicated **Transactions, Consistency & Concurrency** module.

---

# 3. “NoSQL is always eventually consistent”

**False.**

This myth comes from the fact that **eventual consistency is common in distributed NoSQL systems**, but it is not universal.

First, the concept:

### Eventual consistency

Suppose data is replicated:

```text
        Write
          ↓
       Node A
       /    \
      ↓      ↓
   Node B   Node C
```

The write reaches Node A first.

For some period, B and C might still contain the old value.

Eventually, assuming the system continues operating normally, the replicas converge.

That's eventual consistency.

The key property is:

> **The system does not require every relevant replica to agree immediately before an operation can succeed.**

This can reduce coordination and improve availability/latency.

But some NoSQL databases or operations can provide **strong consistency**.

Even within the same database, consistency may sometimes be configurable.

So:

> **Consistency is a property of a specific system and operation—not a defining property of "NoSQL."**

Also, don't confuse:

**eventual consistency**

with:

**inconsistent data forever.**

Eventual consistency means the system has a defined convergence model; it does not mean correctness is irrelevant.

We'll examine consistency models much more deeply later.

---

# 4. “NoSQL is always faster”

**False.**

This is probably the most common misconception.

Performance depends on:

- workload
- data model
- query pattern
- indexes
- dataset size
- partitioning
- hardware
- caching
- network latency
- concurrency
- consistency requirements
- database implementation
- workload distribution

A key-value lookup can be extremely fast because the operation is simple:

```text
key → value
```

But that doesn't mean a key-value database will outperform a relational database for every workload.

Likewise, a relational database can be dramatically faster when its query model matches the workload.

### The deeper point

NoSQL databases often achieve excellent performance by **constraining or specializing the problem**.

If the database knows that requests will primarily look like:

```text
Get item by key
```

it doesn't need to provide the same general-purpose query machinery required by a relational system supporting arbitrary JOINs, aggregations, predicates, and transactional operations.

So a more accurate statement is:

> **Some NoSQL systems can achieve excellent performance for particular workloads because their data model and architecture are optimized for those workloads.**

---

# 5. “NoSQL replaces SQL”

**False.**

NoSQL did not make relational databases obsolete.

Both approaches continue to exist because they solve different problems.

Relational databases remain extremely strong when you need things such as:

- complex relationships
- JOINs
- sophisticated querying
- strong relational integrity
- mature transactional semantics
- structured data
- ad-hoc analytical queries

NoSQL can be preferable when requirements strongly favor things such as:

- massive horizontal distribution
- very high throughput
- predictable access patterns
- flexible document structures
- globally distributed workloads
- specialized data models

And sometimes **both are used in the same system**.

That's the idea behind polyglot persistence, which we'll study later.

The right question isn't:

> "SQL or NoSQL?"

It's:

> **"What data model and operational characteristics does this workload require?"**

---

# 6. “MongoDB is just JSON storage”

**Very false.**

This description captures one superficial characteristic while ignoring most of the database.

MongoDB stores BSON documents—not simply raw JSON files.

BSON is a binary representation with additional data types and is used internally by MongoDB.

But more importantly, MongoDB provides database functionality around those documents, including:

- querying
- indexes
- aggregation
- updates
- transactions
- replication
- sharding
- concurrency control
- durability mechanisms
- change streams
- access control

So this:

```text
MongoDB = JSON files
```

is roughly analogous to saying:

```text
PostgreSQL = files containing rows
```

It describes where information ultimately becomes persisted, but completely misses the **database engine and distributed-system behavior** surrounding it.

The useful mental model is:

> **MongoDB is a document database whose primary data model is BSON documents, with a full database engine built around storing, querying, indexing, updating, replicating, and distributing those documents.**

We'll study MongoDB properly when we reach the Document Databases module.

---

# 7. “Horizontal scaling is free”

**Definitely false.**

Adding machines does not magically produce proportional performance gains.

Suppose one machine handles:

```text
100,000 requests/sec
```

You cannot automatically assume:

```text
2 machines → 200,000
10 machines → 1,000,000
```

Several things can prevent linear scaling.

### Data distribution

The workload must actually be distributable.

If most requests target one partition, adding machines doesn't help much.

```text
Node A → 95% of traffic
Node B → 2.5%
Node C → 2.5%
```

You've added capacity, but the hot partition remains the bottleneck.

### Coordination

If nodes constantly need to coordinate, communication between them introduces:

- network latency
- synchronization overhead
- contention
- additional failure modes

### Rebalancing

Adding or removing nodes can require moving data between machines.

That consumes:

- network bandwidth
- CPU
- storage I/O

and can affect normal workload performance.

### Operational complexity

More machines mean more things that can fail:

```text
1 node
    ↓
1 failure domain

100 nodes
    ↓
100 possible failure points
    +
network interactions
    +
coordination
    +
monitoring
    +
recovery
```

Distributed systems therefore exchange some forms of **hardware limitation** for **distributed-systems complexity**.

This is one of the central themes of the entire course.

---

# The Seven Myths — Correct Mental Models

| Myth | Better mental model |
|---|---|
| NoSQL means no schema | NoSQL often provides **schema flexibility**, not absence of schema |
| NoSQL doesn't support transactions | Transaction support and guarantees **vary by database and scope** |
| NoSQL is always eventually consistent | Consistency models **vary by database and operation** |
| NoSQL is always faster | Performance depends on **workload + data model + architecture** |
| NoSQL replaces SQL | SQL and NoSQL solve **different classes of problems** |
| MongoDB is just JSON storage | MongoDB is a **full document database engine** |
| Horizontal scaling is free | Horizontal scaling introduces **distribution and operational complexity** |

---

## The bigger lesson

The dangerous thing about NoSQL myths is that they usually start with something that is **partially true**:

- flexible schemas → becomes "no schema"
- different transaction models → becomes "no transactions"
- eventual consistency exists → becomes "always eventually consistent"
- excellent performance for certain workloads → becomes "always faster"
- horizontal scaling is a strength → becomes "scaling is free"

As you progress through this course, you'll repeatedly encounter this principle:

> **Database architecture is about trade-offs, not absolutes.**

That principle is more useful than memorizing a list of "SQL vs. NoSQL" differences.

## Interview Topics

### 1. Why would a company choose MongoDB instead of PostgreSQL?

A company might choose **MongoDB** when the application's data is naturally document-oriented and the workload benefits from storing and retrieving related data as a single document.

Typical reasons include:

- **Document-shaped data** — application objects map naturally to documents.
- **Flexible schema** — different documents can evolve without requiring every change to be a relational schema migration.
- **Denormalization** — related data can be embedded, reducing the need for JOINs.
- **Horizontal scaling** — MongoDB provides built-in sharding for distributing data across nodes.
- **High availability** — replica sets provide redundancy and automatic failover.
- **Rapidly evolving applications** — useful when the data model changes frequently.

PostgreSQL may still be better when the application has complex relationships, strong relational integrity requirements, extensive JOINs, or sophisticated transactional requirements.

**Interview answer:**

> I'd choose MongoDB when the data is naturally document-oriented, access patterns favor retrieving related data together, and flexible schema evolution and horizontal scaling are important. I wouldn't choose it simply because it's "NoSQL" or supposedly faster.

---

### 2. Why might DynamoDB be a terrible choice for a particular application?

DynamoDB is highly optimized for **predictable, known access patterns at large scale**. It can be a poor choice when an application requires flexible, relational-style querying.

For example, imagine an analytics application where users can arbitrarily filter, sort, aggregate, and JOIN data across many entities.

DynamoDB can become awkward because:

- Data modeling is heavily driven by **known access patterns**.
- Arbitrary queries are not its strength.
- JOINs are not a native relational operation.
- Complex relationships can require denormalization or multiple application-level queries.
- Poorly designed partition keys can create **hot partitions**.
- Ad-hoc analytical workloads are generally better suited to other database systems.

**Interview answer:**

> DynamoDB can be a terrible choice when the workload is unpredictable or query-heavy, especially when users need arbitrary filtering, complex relationships, JOINs, or analytical queries. DynamoDB works best when access patterns are well understood and the data can be modeled around those patterns.

---

### 3. When is SQL clearly the better choice?

SQL is usually the better choice when the application fundamentally depends on **relationships, constraints, transactions, and flexible querying**.

Examples include:

- Financial/accounting systems
- Order and payment processing
- Systems with many related entities
- Applications requiring strong referential integrity
- Complex reporting and analytical queries
- Workloads requiring many JOINs
- Applications where requirements cannot easily predict future query patterns

For example, an e-commerce system might have:

`Customers → Orders → OrderItems → Products → Payments`

If the application frequently needs to query and combine these relationships while maintaining strong transactional guarantees, a relational database such as PostgreSQL can be a much more natural fit.

**Interview answer:**

> SQL is clearly preferable when relationships and data integrity are central to the application, especially when we need complex JOINs, constraints, transactions, and flexible ad-hoc queries. I wouldn't introduce NoSQL just because the application needs to scale.

---

### 4. What does "schema flexibility" actually mean?

**Schema flexibility does not mean having no schema.**

It means the database does not necessarily require every record to conform to one rigid structure enforced in the same way as a traditional relational table.

For example, documents could evolve from:

```json
{
  "name": "Alice",
  "email": "alice@example.com"
}
```

to:

```json
{
  "name": "Alice",
  "email": "alice@example.com",
  "phone": "+49..."
}
```

without requiring an `ALTER TABLE` operation to add a column.

However, the application may still have an expected schema, validation rules, API contracts, serialization models, or database-level validation.

The important distinction is:

> **Schema flexibility means the stored data structure can evolve more easily; it does not mean the application has no schema or structure.**

The trade-off is that flexibility can move some responsibility from the database to the application. If poorly managed, documents can become inconsistent and difficult to maintain.

**Interview answer:**

> Schema flexibility means the database allows records to evolve without requiring every structural change to be applied uniformly through a rigid schema migration. It doesn't mean there is no schema. The application still needs a consistent data model and validation strategy.

### Module 1 status

We've now completed:

- **1.1 — The Problems NoSQL Tries to Solve** ✅
- **1.2 — SQL vs. NoSQL: Under the Hood** ✅
- **1.3 — The NoSQL Landscape** ✅
- **1.4 — NoSQL Myths** ✅

**Next: Chapter 1.5**