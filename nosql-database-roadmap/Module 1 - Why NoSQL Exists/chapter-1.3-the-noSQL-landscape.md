# Chapter 1.3 — The NoSQL Landscape

NoSQL is **not one type of database**. It is an umbrella term covering several fundamentally different data models.

The four major families we will use throughout this course are:

```text
                         NoSQL
                           |
             +-------------+-------------+
             |             |             |
         Document      Key-Value    Wide-Column
             |
           Graph
```

The diagram is only a learning simplification. **Graph is a fourth peer family**, not a subtype of Document. A better representation is:

```text
                         NoSQL
                           |
       +-------------------+-------------------+
       |                   |                   |
   Document            Key-Value         Wide-Column
       |
     Graph
```

Conceptually:

```text
                         NoSQL
                           |
       +-------------------+-------------------+-------------------+
       |                   |                   |                   |
   Document            Key-Value         Wide-Column           Graph
```

These families differ primarily in **how they represent data and what kinds of access patterns they are optimized for**.

---

# 1. Document Databases

A document database stores data as **documents**, usually represented using JSON-like structures.

Conceptually:

```text
{
    "id": 42,
    "name": "Alice",
    "address": {
        "city": "Berlin",
        "country": "Germany"
    },
    "orders": [
        {
            "id": 1001,
            "total": 250
        }
    ]
}
```

The important characteristic is that the document can contain **nested structures and collections of related data**.

Unlike a traditional relational table, the structure does not necessarily have to be identical for every document.

### Strengths

Document databases are particularly useful when:

- data naturally forms hierarchical structures
- an application usually retrieves a whole object/document
- related information can reasonably be stored together
- schema evolution is frequent
- horizontal distribution is important

### Trade-offs

Because data can be denormalized:

- information may be duplicated
- updates may need to modify multiple documents
- relationships are generally less natural than in relational databases
- complex cross-document operations can become expensive or awkward

### Representative database

 → **Document**

MongoDB is the primary document database we'll study in this course.

---

# 2. Key-Value Databases

The simplest NoSQL model is:

```text
Key → Value
```

For example:

```text
"user:42" → <some value>
```

The database fundamentally knows how to associate a key with a value.

The key is normally the primary way of locating the data.

This model is extremely attractive when the application already knows **exactly what piece of data it wants**.

### Strengths

Key-value systems can provide:

- very fast lookups
- simple data access
- easy horizontal distribution
- high request throughput
- straightforward partitioning by key

The simplicity of the model is actually an architectural advantage.

The database has fewer relationships and query semantics to coordinate.

### Trade-offs

The simplicity also limits what you can naturally ask.

A pure key-value model isn't designed around queries such as:

> "Find all customers in Berlin whose balance exceeds €5,000 and sort them by balance."

Instead, the application generally knows the key it needs.

This means **data modeling becomes heavily driven by access patterns**.

### Representative databases

 → **Key-value / data structures**

Redis goes beyond a simple string value. It provides data structures such as:

- strings
- hashes
- lists
- sets
- sorted sets
- streams

So calling Redis merely a "key-value database" is technically useful but incomplete.

 → **Key-value + document**

DynamoDB supports both key-value and document-oriented data models.

We'll examine this distinction carefully later.

---

# 3. Wide-Column Databases

Wide-column databases are quite different from both document and traditional relational databases.

They organize data around concepts such as:

- rows
- partition keys
- clustering keys/columns
- column families/tables

The important idea is that the database is designed around **partitioned, distributed datasets and predictable access patterns**.

A simplified conceptual structure is:

```text
Partition
    |
    +-- Row
    +-- Row
    +-- Row
```

Rows belonging to the same partition are typically stored together.

The partition key therefore becomes an extremely important modeling decision.

### Why this model exists

Wide-column systems are designed for workloads involving:

- enormous datasets
- very high write volumes
- distributed storage
- predictable query patterns
- high availability
- horizontal scaling

Instead of designing the database primarily around arbitrary relational queries, you design the data around **how the application will retrieve it**.

### Trade-offs

The price is reduced flexibility for arbitrary querying.

A poorly chosen partition key can cause:

- hot partitions
- uneven data distribution
- overloaded nodes
- poor performance

So wide-column modeling requires strong understanding of workload and access patterns.

### Representative databases

 → **Wide-column**

 → **Wide-column**

We'll study both because they are particularly valuable for understanding distributed database architecture.

---

# 4. Graph Databases

Graph databases model information as a graph:

```text
Nodes ── Relationships ── Nodes
```

More formally:

```text
Entities → Relationships → Entities
```

A graph database therefore treats **relationships as first-class data**.

For example, the important information might not simply be:

```text
Person A
Person B
```

but:

```text
Person A
   |
   | WORKS_WITH
   ↓
Person B
```

and:

```text
Person A
   |
   | KNOWS
   ↓
Person C
   |
   | WORKS_AT
   ↓
Company X
```

This makes graph databases particularly suitable for workloads where the **connections between entities** are central to the problem.

Common areas include:

- recommendation systems
- fraud detection
- dependency analysis
- network analysis
- knowledge graphs
- social relationships

### Representative database

 → **Graph**

---

# 5. The Most Important Difference Between the Families

Don't memorize the four families merely as four product categories.

Understand the **primary data-access mental model**:

| Family | Think primarily in terms of |
|---|---|
| Document | Documents / objects |
| Key-value | Key → value |
| Wide-column | Partitions + rows + access patterns |
| Graph | Nodes + relationships |

This distinction becomes extremely important when we start **data modeling**.

---

# 6. Why Do We Need Different NoSQL Families?

Because "NoSQL" doesn't describe a single workload.

Consider four fundamentally different requirements:

### Workload A

> Retrieve an entire customer profile by ID.

A document model is natural.

### Workload B

> Retrieve a value immediately using a known key.

A key-value model is natural.

### Workload C

> Continuously ingest enormous amounts of distributed event data and retrieve it according to known partition-oriented access patterns.

A wide-column model can be appropriate.

### Workload D

> Traverse relationships between millions of entities.

A graph model is natural.

The important lesson isn't the examples themselves.

It's this:

> **The data model should make the important operations natural and efficient.**

---

# 7. Real Products Don't Fit Perfectly Into Four Boxes

This is an important part of the original syllabus.

The four families are **conceptual categories**, not rigid product boundaries.

Modern databases increasingly overlap.

For example:

| Database | Primary classification | Important additional characteristics |
|---|---|---|
| MongoDB | Document | Rich querying, indexing, aggregation |
| Redis | Key-value | Multiple data structures, streams, scripting |
| DynamoDB | Key-value / Document | Secondary indexes, rich attribute types |
| Cassandra | Wide-column | Distributed, partition-oriented architecture |
| ScyllaDB | Wide-column | Cassandra-compatible distributed architecture |
| Neo4j | Graph | Graph traversal/query capabilities |

Therefore:

> **"What type of database is X?" usually has a useful primary answer, but that answer doesn't describe everything the database can do.**

---

# 8. Database Category vs. Product Capability

This distinction will become increasingly important.

A database's **category** describes its fundamental data model.

Its **capabilities** describe what else the product provides.

For example, a document database may support:

- secondary indexes
- aggregation
- transactions
- replication
- sharding
- caching
- change streams

That doesn't make it a relational database.

Similarly, a key-value database may support sophisticated data structures without ceasing to be fundamentally key-oriented.

So avoid thinking:

```text
One feature = one database category
```

Instead think:

```text
Data model
     +
Query model
     +
Distribution architecture
     +
Consistency model
     +
Other capabilities
```

These collectively define what a database actually is good at.

---

# 9. NoSQL vs. "No SQL"

Another important clarification:

**NoSQL does not mean "a database that doesn't use SQL."**

Some NoSQL databases have their own query languages.

Some use APIs.

Some use SQL-like query languages.

For example, Cassandra has CQL (Cassandra Query Language), which deliberately resembles SQL syntactically.

The term **NoSQL** is therefore better understood historically as referring to **non-relational database approaches**, although even that definition has become somewhat blurry as modern systems gain more capabilities.

---

# 10. The Families Also Differ in How They Distribute Data

This is a particularly important connection to Chapter 1.2.

Different models make different distribution strategies natural.

### Document

A document can often be distributed according to a shard/partition key.

```text
Document A → Node 1
Document B → Node 2
Document C → Node 3
```

### Key-value

The key can naturally determine where the value belongs.

```text
hash(key) → partition → node
```

This is one reason key-value systems are so naturally distributable.

### Wide-column

The partition key is fundamental to determining data placement.

```text
partition key
      ↓
partition
      ↓
node(s)
```

### Graph

Distribution is more difficult when relationships frequently cross machines.

A graph traversal may need to follow relationships from one node to another and potentially cross partition boundaries.

This is one reason graph workloads present **different distributed-system challenges** from key-value workloads.

---

# 11. A Crucial Principle: Data Model Influences Distribution

This is worth remembering:

> **Your data model isn't merely about how data looks. It affects how data can be partitioned, queried, replicated, and scaled.**

That is why NoSQL data modeling is so important.

In relational databases, you can often begin with entities and relationships and let the query planner determine how to retrieve the data.

In many NoSQL systems, you need to think much more explicitly about:

**"How will this data be accessed, and how can those accesses be distributed?"**

We'll eventually make this a major part of the course.

---

# 12. Polyglot Persistence

Because the families solve different problems, a large system doesn't necessarily need to choose **one database for everything**.

A system might use:

```text
                 Application
                /     |      \
               /      |       \
        Document   Key-Value   Graph
```

Each database handles the workload for which its model is appropriate.

This is called **polyglot persistence**.

However, this introduces significant complexity:

- multiple technologies
- multiple operational systems
- data synchronization
- consistency across databases
- backup/recovery for multiple systems
- monitoring
- deployment
- developer expertise

So:

> **Using multiple databases is not automatically better architecture.**

It should be justified by genuinely different workload requirements.

We'll study this properly much later in the course.

---

# 13. One More Important Category: Search and Specialized Databases

The four families are the **core NoSQL families for this course**, but the real database ecosystem contains additional specialized systems.

For example:

- search engines
- time-series databases
- vector databases
- streaming/event stores

Some of these are sometimes grouped loosely into the broader NoSQL ecosystem, while others are treated as separate database categories.

We'll cover these later under **Advanced Architecture**, rather than confusing the basic four-family model now.

---

# 14. The Landscape You Should Have in Your Head

At this stage, your mental map should be:

```text
                              DATABASES
                                  |
                    +-------------+-------------+
                    |                           |
               RELATIONAL                    NoSQL
                    |                           |
             Tables / rows             +---------+---------+---------+
             Relationships             |         |         |         |
             JOINs                 Document  Key-Value  Wide-Column  Graph
             Strong schema             |         |         |         |
             ACID                   MongoDB    Redis   Cassandra   Neo4j
                                      |      DynamoDB   ScyllaDB
```

This is **not a ranking**.

It is a classification based primarily on **data model and access patterns**.

---

# 15. What You Should Be Able to Explain Now

If someone asks:

**"What are the major types of NoSQL databases?"**

You should be able to say:

> The four major families are document, key-value, wide-column, and graph databases. They differ primarily in their data models and the access patterns they're designed to support. MongoDB is a document database, Redis is primarily key-value with rich data structures, DynamoDB supports key-value and document models, Cassandra and ScyllaDB are wide-column databases, and Neo4j is a graph database. These categories aren't completely rigid because modern products increasingly overlap in their capabilities.

And the deeper understanding is:

> **NoSQL isn't a single architecture. It's a collection of different approaches to modeling and distributing data, each making different trade-offs.**

That distinction is foundational for the rest of the course.