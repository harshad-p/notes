# Chapter 1.2 — SQL vs. NoSQL: Under the Hood

This chapter will focus specifically on **what changes inside the architecture when we move from a traditional relational approach toward NoSQL**, and why those changes address the problems from Chapter 1.1.

## 1. Relational approach

### Strong schema

A relational database generally defines a structured schema before data is stored:

- Tables
- Columns
- Data types
- Constraints
- Relationships

The database can therefore enforce rules such as:

- a column must contain an integer
- a value cannot be `NULL`
- a foreign key must reference an existing row
- a value must be unique

This gives you strong structural guarantees.

The important point is that the schema isn't merely documentation. **The database itself understands and enforces much of the structure.**

---

### Relations

Relational databases represent data as relations—practically, tables containing rows and columns.

Relationships between pieces of data are represented using keys.

For example:

```text
Customers
    Id
    Name

Orders
    Id
    CustomerId
    Total
```

`Orders.CustomerId` can reference `Customers.Id`.

The data is therefore kept relatively independent rather than repeatedly embedded into other records.

This is closely connected to **normalization**.

---

### JOINs

Because related information is stored separately, retrieving a complete piece of information can require a JOIN.

Conceptually:

```text
Customers
    ↓
    JOIN
    ↓
Orders
```

On one database server, the database engine can perform this efficiently using indexes, query planning, buffering, etc.

But as data becomes distributed across machines, JOINs become substantially more complicated.

The database may need to:

1. locate the relevant data on multiple nodes
2. send requests across the network
3. retrieve intermediate results
4. combine those results
5. coordinate the operation

The **network** therefore becomes part of query execution.

This is one of the fundamental difficulties NoSQL architectures try to avoid.

---

## 2. ACID transactions

Relational databases traditionally provide strong transactional semantics.

A transaction can modify multiple pieces of related data and make them appear as one atomic operation.

For example, conceptually:

```text
Update A
Update B
Update C
    ↓
Commit everything
```

or:

```text
something fails
    ↓
rollback
```

This is extremely valuable when correctness depends on several changes succeeding together.

But distributed transactions are considerably harder.

If A is on Node 1 and B is on Node 2, committing the transaction may require coordination between those machines.

Now you have:

- network communication
- partial failures
- timeouts
- coordination protocols
- distributed locking or equivalent mechanisms
- additional latency

So the problem isn't that **ACID is bad**.

The problem is:

> **Strong transactional guarantees become more expensive as the amount of distributed coordination increases.**

---

# 3. Centralized coordination

Traditional relational databases are often designed around a relatively centralized architecture.

That doesn't mean modern relational databases are necessarily running on exactly one machine. They can have:

- read replicas
- clustering
- partitioning
- distributed implementations
- sharding

But the traditional relational model naturally works well when a database can coordinate operations within a relatively centralized system.

This gives the database a lot of control.

It can centrally reason about:

- transactions
- constraints
- indexes
- relationships
- query execution
- locking
- consistency

That makes complex operations easier to reason about.

The downside is that when you want to distribute everything across many machines and regions, **that coordination itself becomes a scalability challenge**.

---

# 4. What changes with NoSQL?

NoSQL databases don't simply remove SQL syntax.

They often change the **data model and architectural assumptions** so that distribution becomes easier.

The major shift is:

> Instead of starting with highly normalized relationships and asking the distributed system to coordinate them, many NoSQL systems start with the application's access patterns and organize data so that requests can be handled more independently.

This leads to several important differences.

---

## 5. Denormalization

Instead of always storing related information separately, NoSQL systems—particularly document databases—often duplicate related data.

Conceptually:

```text
One document
    ├── customer information
    ├── order information
    └── relevant details
```

Now retrieving that information may require only one lookup.

The benefit:

**less cross-record coordination and fewer distributed JOINs.**

The cost:

**duplicated data and more complicated updates.**

If the duplicated information changes, you may need to update multiple copies.

So NoSQL isn't eliminating complexity.

It's often **moving complexity from reads/relationships into data duplication and consistency management.**

---

# 6. Partitioning

A distributed NoSQL database can divide its dataset across multiple machines.

Conceptually:

```text
Node A → portion of data
Node B → another portion
Node C → another portion
```

Instead of one machine having to handle the entire dataset and workload, different nodes can handle different portions.

This is fundamental to horizontal scaling.

But partitioning introduces new problems:

- How do we decide where data belongs?
- How do we find it?
- What happens when a node fails?
- What happens when a node becomes overloaded?
- How do we redistribute data?
- What happens when one partition receives far more traffic than others?

We'll study partitioning and sharding deeply later.

---

# 7. Replication

NoSQL systems commonly replicate data across multiple nodes.

Conceptually:

```text
        Data
       /    \
   Node A  Node B
```

This provides redundancy.

If one machine fails, another copy may still be available.

Replication therefore helps address the **availability and failure** problems from Chapter 1.1.

But replication introduces another problem:

> **How quickly must all copies agree?**

If Node A has the newest value but Node B hasn't received it yet, they temporarily contain different states.

That leads directly into consistency models.

---

# 8. Distributed coordination

This is one of the most important differences to understand.

A centralized database can coordinate operations relatively easily because much of the relevant state exists within one system.

A distributed database has to coordinate **machines**.

Machines communicate over a network, and networks can:

- be slow
- lose messages
- partition
- time out
- become overloaded

Machines can also fail independently.

Therefore, distributed databases have to explicitly deal with problems that are much less significant in a single-node architecture.

This is why distributed databases aren't simply:

> "A normal database running on five computers."

The architecture itself changes.

---

# 9. Consistency becomes a design choice

Traditional relational systems generally emphasize strong consistency and transactional guarantees.

NoSQL systems have historically been willing to offer different consistency guarantees depending on the workload.

For example, a system may allow:

- strong consistency
- eventual consistency
- configurable consistency
- consistency scoped to particular operations or partitions

The advantage is that an application can sometimes prioritize:

- availability
- latency
- geographic distribution
- scalability

instead of requiring every operation to wait for every relevant replica to agree.

But weaker consistency introduces application-level consequences.

The application has to understand that:

> **The value it reads may not always represent the most recent write.**

This trade-off is central to distributed NoSQL architecture.

We'll go much deeper into consistency later.

---

# 10. The fundamental architectural trade-off

This is the key idea of this chapter.

### Relational approach

Generally favors:

```text
Structured schema
      +
Relationships
      +
JOINs
      +
Strong transactions
      +
Centralized coordination
```

This makes complex relationships and transactional correctness easier to express.

### NoSQL approach

Many NoSQL systems instead emphasize:

```text
Data organized for access patterns
      +
Denormalization
      +
Partitioning
      +
Replication
      +
Distributed execution
      +
Flexible consistency
```

This can make massive distributed workloads easier to scale.

But you pay for it through things such as:

- duplicated data
- more application-aware modeling
- weaker or different transactional guarantees
- consistency considerations
- partitioning problems
- operational complexity
- distributed failure handling

---

# 11. Why this connects directly to Chapter 1.1

The problems we identified previously now have corresponding architectural responses:

| Problem | Typical distributed/NoSQL response |
|---|---|
| Massive datasets | Partitioning |
| Massive request volume | Partitioning + replication + horizontal scaling |
| Global applications | Geographic distribution + replication |
| High availability | Replication |
| Low latency | Locality + partitioning + caching/replication |
| Flexible schemas | Document/flexible data models |
| Rapid schema evolution | Less rigid storage schemas |
| Distributed workloads | Partitioning + distributed execution |

Notice something important:

**NoSQL isn't one solution.**

Different NoSQL databases make different choices about these mechanisms.

That's why MongoDB, Redis, Cassandra, DynamoDB, and Neo4j can all be called "NoSQL" while having radically different architectures.

---

## The mental model to keep

Don't remember:

> SQL = old, NoSQL = new.

Remember:

> **Relational systems optimize around relationships, structured schemas, and strong transactional semantics.**

> **Many NoSQL systems optimize around distributing data and workload across machines while tailoring the data model to access patterns.**

Neither is universally superior.

The engineering question is:

> **Which set of trade-offs matches the workload?**

That is the foundation for everything that follows in this course.