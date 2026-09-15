# Module 0 — Prerequisites & Mental Model

We’ll do this as an actual course, **one lesson at a time**, rather than covering the entire module in one dump.

---

# Lesson 0.1 — What Is a Database?

You already know what a database is at the application level. For this course, I want to go **one layer deeper**.

Suppose your application executes:

```sql
SELECT * FROM users WHERE id = 42;
```

At the application level, it looks incredibly simple:

```text
Application
     |
     | SQL query
     ↓
  Database
     |
     ↓
   User
```

But the database has to do quite a lot of work.

Conceptually:

```text
Application
     |
     ↓
Database Server
     |
     ├── Query processing
     ├── Concurrency control
     ├── Cache / memory
     ├── Indexes
     ├── Storage engine
     └── Disk / SSD
```

Let's understand each piece.

---

## 1. Database vs. Database Server

These terms are often used interchangeably, but they're slightly different.

A **database** is the organized collection of data.

For example:

```text
Users
Orders
Products
Payments
```

A **database server/engine** is the software responsible for managing that data.

Examples:

- PostgreSQL
- SQL Server
- MongoDB
- Redis
- Cassandra
- Neo4j

So:

```text
MongoDB = database software

your application data = data managed by MongoDB
```

When you install MongoDB, you're installing software that provides the database engine/server.

---

# 2. Where Is the Data Actually Stored?

Ultimately, persistent database data needs to survive a process restart or machine reboot.

That means it has to reach persistent storage, typically an SSD.

Conceptually:

```text
Application
     ↓
Database
     ↓
Memory/RAM
     ↓
SSD/Disk
```

Why use RAM?

Because RAM is dramatically faster than persistent storage.

So databases try very hard to avoid unnecessarily going to disk.

This is one reason **caching** will become a major topic later.

---

# 3. The Database Is More Than "Files on Disk"

Imagine MongoDB storing:

```json
{
    "id": 42,
    "name": "Alice",
    "country": "Germany"
}
```

You might initially think:

> "The database just writes this JSON somewhere."

Not really.

The database has to solve problems such as:

- How do I locate Alice quickly?
- What if 100 clients request Alice simultaneously?
- What if the process crashes during a write?
- What if two clients modify Alice at the same time?
- How do I make sure the data survives?
- How do I recover after a crash?
- How do I handle millions of records?
- How do I replicate the data to another server?

This is where the **database engine** comes in.

---

# 4. Storage Engine

The **storage engine** is the component responsible for managing how data is stored and retrieved.

Different databases use different storage mechanisms.

For example, MongoDB uses **WiredTiger** as its primary storage engine.

A relational database might use B-tree indexes and various page/storage structures.

A wide-column database such as Cassandra uses a very different architecture involving things like:

- Commit logs
- Memtables
- SSTables
- Compaction

We'll eventually study these in detail.

For now, remember:

> **The data model is what you see. The storage engine is how the database actually manages it underneath.**

That's an extremely important distinction for this course.

---

# 5. Query Processing

Suppose you ask MongoDB:

```javascript
db.users.find({ age: 37 })
```

The database needs to determine:

> How do I find these users efficiently?

It might have an index:

```text
age index
   |
   +---- 20 → ...
   +---- 25 → ...
   +---- 30 → ...
   +---- 37 → User A
   +---- 37 → User B
   +---- 40 → ...
```

Or it might have no useful index.

Then it may have to examine every document:

```text
User 1  → no
User 2  → no
User 3  → no
User 4  → yes
User 5  → no
...
User 10,000,000 → no
```

That's a **collection scan**.

This distinction will become extremely important when we study indexing and performance.

---

# 6. Concurrency

Now imagine:

```text
Request A:
Change Alice's balance to €500

Request B:
Change Alice's balance to €700
```

Both arrive almost simultaneously.

The database has to deal with concurrent operations safely.

This leads to concepts such as:

- Locks
- Isolation
- Atomicity
- MVCC
- Transactions
- Optimistic concurrency
- Pessimistic concurrency

We'll study these later.

The important thing for now is:

> A database isn't simply reading and writing data. It is coordinating potentially thousands or millions of concurrent operations.

---

# 7. Durability

Consider:

```text
Application
    ↓
Database
    ↓
"Payment recorded"
```

The database tells the application:

```text
SUCCESS
```

Then immediately afterward:

```text
💥 Server crashes
```

What should happen?

If the database had promised that the write was durable, the payment shouldn't simply disappear.

This is the **D in ACID: Durability**.

Later we'll look at how databases achieve this using mechanisms such as:

- Write-ahead/commit logs
- Journaling
- Replication
- Flushes
- Snapshots
- Persistent storage

---

# 8. Availability

Now imagine the database consists of one server:

```text
        Database
           |
        Server A
```

Server A dies.

Your application can't access the database.

Now imagine:

```text
             Database
          /     |     \
      Server A Server B Server C
```

If Server A dies, perhaps B can continue serving requests.

This is where **replication** becomes important.

And this is one of the fundamental reasons NoSQL and distributed databases become interesting.

---

# 9. Latency vs. Throughput

These two are frequently confused.

### Latency

How long does **one operation** take?

Example:

```text
GET user 42 → 8 ms
```

### Throughput

How many operations can the system handle over time?

Example:

```text
100,000 requests/second
```

A system could have:

```text
Very low latency
but
low throughput
```

or:

```text
High throughput
but
higher individual latency
```

Production architecture often involves balancing both.

---

# 10. The Mental Model I Want You to Build

By the end of this course, when you hear:

> "MongoDB database"

I don't want your mental model to simply be:

```text
MongoDB
   ↓
JSON documents
```

I want something closer to:

```text
                     Application
                          |
                       Network
                          |
                    MongoDB Node
                          |
              +-----------+-----------+
              |                       |
         Query Engine             Memory/Cache
              |                       |
              +-----------+-----------+
                          |
                    Storage Engine
                          |
                    Persistent Data
                          |
                       Disk
```

And eventually:

```text
                       Application
                            |
                       Load Balancer
                            |
             +--------------+--------------+
             |              |              |
          Node A          Node B          Node C
             |              |              |
             +--------------+--------------+
                            |
                     Replication
                            |
                    Distributed Data
```

That second diagram is where our NoSQL journey is ultimately heading.

---

# Quick Check

Before moving on, make sure these distinctions are clear:

**Database**

→ The data being managed.

**Database engine/server**

→ Software that manages the data.

**Storage engine**

→ Component responsible for how data is stored/retrieved.

**Index**

→ Additional data structure allowing faster lookup.

**Replication**

→ Keeping multiple copies of data.

**Partitioning/sharding**

→ Distributing data across different nodes.

**Latency**

→ Time for an operation.

**Throughput**

→ Amount of work handled per unit of time.

These eight concepts will keep coming back throughout the course.

---

### Small exercise

Don't look this up. Just answer from your current understanding:

**1.** If a database has 10 million documents and you search for one document without an appropriate index, what is the database likely to do?

**2.** If the database has three replicas and one server dies, what problem is replication trying to solve?

**3.** What's the difference between *replication* and *sharding*?

You can answer in your own words—even if you're unsure. 

---

**1.** The database will do a complete scan.&#x20;

**2.** Availability

**3.** Replication vs. Sharding
- **Replication** = multiple nodes have copies of the same logical data.
- **Sharding/partitioning** = different nodes are responsible for different portions of the data.

So with 300 GB:

```text
Replication:

Node A → 300 GB
Node B → 300 GB
Node C → 300 GB
```

Whereas sharding might be:

```text
Shard A → 100 GB
Shard B → 100 GB
Shard C → 100 GB
```

And production systems commonly combine them:

```text
              300 GB logical dataset
                       |
          +------------+------------+
          |            |            |
       Shard 1      Shard 2      Shard 3
       100 GB       100 GB       100 GB
        /  \          /  \          /  \
       A    B        C    D        E    F
      replica       replica       replica
```
