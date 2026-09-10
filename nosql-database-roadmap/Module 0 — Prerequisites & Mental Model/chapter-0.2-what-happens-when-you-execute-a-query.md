# Lesson 0.2 — What Happens When You Execute a Query?

Now we're going one level deeper.

Suppose your .NET application executes:

```csharp
var user = await db.Users
    .FirstOrDefaultAsync(x => x.Id == 42);
```

At the application level, that looks like one operation.

Underneath, there is a chain of events.

```text
.NET Application
       |
       | request
       ↓
   Database
       |
       ↓
  Query processing
       |
       ↓
  Query plan
       |
       ↓
  Index / data lookup
       |
       ↓
   Storage engine
       |
       ↓
  RAM / SSD
```

Let's take this apart.

---

## 1. The Client Sends a Request

Your application doesn't directly manipulate the database's files.

It communicates with the database through a protocol.

For example:

```text
.NET Application
       |
       | TCP/network connection
       ↓
PostgreSQL
```

or:

```text
.NET Application
       |
       ↓
MongoDB
```

Even if the database happens to be running on your own Mac:

```text
.NET Application
       |
    localhost
       |
       ↓
   Database
```

there is still communication between two processes.

This distinction becomes **very important once the database is remote**.

---

# 2. The Database Receives the Query

Imagine:

```sql
SELECT *
FROM users
WHERE id = 42;
```

The database has to understand what you've asked.

There are broadly several stages:

```text
Query
  ↓
Parsing
  ↓
Validation
  ↓
Planning
  ↓
Execution
  ↓
Result
```

---

# 3. Parsing

The database first determines whether your query is syntactically valid.

For example:

```sql
SELECT * FROM users WHERE id = 42;
```

is valid.

Something like:

```sql
SELCT * FROM users
```

isn't.

The database builds an internal representation of the query.

You don't need to know the implementation details yet. Just understand:

> The database has to turn your query text into something its engine can execute.

---

# 4. Query Planning

Now comes one of the most important database concepts.

Suppose there are **10 million users**.

The database needs to find:

```text
id = 42
```

It could do:

```text
User 1
User 2
User 3
...
User 10,000,000
```

That's a full scan.

Or perhaps there is an index:

```text
              Index
                |
                ↓
             id = 42
                |
                ↓
           User record
```

The database's **query planner/optimizer** determines an efficient execution strategy.

This is why an index isn't simply:

> "A thing that makes queries faster."

More accurately:

> An index gives the database additional structures from which it can construct a more efficient execution plan.

That distinction becomes important later.

---

# 5. The Query Planner Can Choose Different Strategies

Imagine:

```sql
SELECT *
FROM users
WHERE country = 'Germany';
```

Suppose 40% of your 10 million users are in Germany.

An index on `country` might not necessarily be the best strategy.

The database could decide:

```text
Index scan
```

or:

```text
Sequential/full scan
```

depending on the situation.

Why?

Because an index has a cost too.

For example, if the database has to retrieve huge numbers of records anyway, scanning the underlying data may be cheaper than jumping through an index and then fetching millions of records.

**This is why "I created an index, therefore my query will be faster" is an oversimplification.**

---

# 6. The Storage Engine Gets Involved

Once the database has an execution plan, the storage engine actually retrieves the required data.

Conceptually:

```text
Query Planner
     ↓
"Find record 42"
     ↓
Storage Engine
     ↓
RAM?
  /   \
yes    no
 |      |
 ↓      ↓
return  SSD
        ↓
      return
```

And this brings us to an important concept:

## The database wants data in memory whenever possible.

RAM is much faster than SSD storage.

So databases maintain various caches/buffers.

---

# 7. Buffer/Cache

Suppose you repeatedly request:

```text
User 42
```

The database may already have the relevant data in memory.

Then:

```text
Query
 ↓
Memory
 ↓
Result
```

rather than:

```text
Query
 ↓
Memory
 ↓
SSD
 ↓
Memory
 ↓
Result
```

That difference can be substantial.

This is why database performance isn't simply:

> "How fast is the SSD?"

The database's memory usage, cache behavior, indexes, workload, CPU, network, and storage all matter.

---

# 8. What About a Write?

Reads are relatively easy to visualize.

Writes introduce another important problem.

Suppose:

```sql
UPDATE users
SET name = 'Bob'
WHERE id = 42;
```

The database needs to ensure that the update is handled correctly.

Among other things, it has to consider:

- Where is the existing data?
- Does another transaction modify it?
- What indexes need updating?
- When is the change considered committed?
- What happens if the process crashes?
- How does the change reach persistent storage?
- If replicas exist, when do they receive it?

This is where concepts such as:

- Transactions
- Concurrency control
- WAL/journaling
- Durability
- Replication

become important.

We'll go much deeper into those later.

---

# 9. Now Add a Distributed Database

So far we've assumed:

```text
Application
     ↓
One database server
```

But imagine:

```text
                  Application
                       |
                 Database cluster
                 /       |       \
              Node A   Node B   Node C
```

Now the query might involve additional decisions.

For example:

> Which node has the data?

Or:

> Should I read from the primary or a replica?

Or:

> Does this write need to be replicated before I tell the client it succeeded?

Or:

> What happens if Node B doesn't respond?

This is where NoSQL becomes particularly interesting.

---

# 10. The Network Is Now Part of Your Database

This is a crucial mental shift.

With a single local database:

```text
Application → Database
```

you can often pretend the database is a single thing.

With a distributed database:

```text
              Network
          /      |      \
       Node A  Node B  Node C
```

the network itself becomes part of the system's behavior.

Networks can:

- be slow
- drop packets
- partition
- temporarily disconnect nodes
- deliver responses at different times

And nodes can:

- crash
- restart
- become overloaded
- fall behind

This is the foundation for why we'll later need:

**CAP → replication → quorum → consistency → partitioning → failure handling.**

---

# The Complete Mental Model

For now, think of a database query like this:

```text
Application
     |
     | request
     ↓
Database Server
     |
     ↓
Parse / Validate
     |
     ↓
Query Planner
     |
     ↓
Execution Plan
     |
     ↓
Index / Data Access
     |
     ↓
Storage Engine
     |
     +------→ RAM/Cache
     |
     +------→ SSD/Disk
     |
     ↓
Result
     |
     ↓
Application
```

And in a distributed NoSQL system:

```text
                    Application
                         |
                    Query/Write
                         |
                  Database Cluster
                 /       |       \
              Node A   Node B   Node C
                 \       |       /
                    Replication
                         |
                    Persistent
                       data
```

We'll spend a **lot** of time understanding that second diagram.

---

## One important interview distinction

If an interviewer asks:

> **"Why is my database query slow?"**

A weak answer is:

> "Maybe you need an index."

A stronger answer starts investigating the whole path:

```text
Application
    ↓
Network latency?
    ↓
Connection pool?
    ↓
Database
    ↓
Query plan?
    ↓
Index?
    ↓
CPU?
    ↓
Memory/cache?
    ↓
Disk I/O?
    ↓
Lock/contention?
    ↓
Replication?
    ↓
Database node overloaded?
```

That's the level of thinking we're building toward.

---

### Quick check

Consider this query:

```sql
SELECT *
FROM orders
WHERE customer_id = 123;
```

There are **50 million orders**.

Answer these three:

1. What could happen if there is **no index** on `customer_id`?
2. What changes if there **is** an appropriate index?
3. If the query is still slow despite the index, name **two possible reasons**.

---

1. **No index:** generally a full/large scan of the relevant table/collection. Whether that is “slow” depends on data size, storage, caching, selectivity, etc.
2. **With an index:** yes. The optimizer considers the index, but may reject it if it estimates a scan is cheaper—for example, if `customer_id = 123` matches a huge percentage of the rows.
3. **Still slow:** CPU, I/O, blocking/locks, memory pressure, poor query plan, low selectivity, returning too many rows, network latency, etc. Good production thinking.
