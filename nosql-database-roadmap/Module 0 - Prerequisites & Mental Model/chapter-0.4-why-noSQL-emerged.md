# 0.4 — SQL vs. NoSQL: Why NoSQL Emerged

Now we're ready for the final part of Module 0.

The goal here is **not** to conclude that NoSQL is better than SQL.

It's to understand **what problems led people to build different database models.**

---

## 1. The traditional relational model

A relational database gives you something extremely powerful:

```text
Structured data
+
Relationships
+
JOINs
+
Transactions
+
Strong consistency
+
Powerful querying
```

For many applications, this is exactly what you want.

For example:

```text
Banking
ERP
Accounting
Inventory
Order management
```

If you transfer €100:

```text
Account A: -€100
Account B: +€100
```

you generally don't want a situation where only one side succeeds.

That's where transactional guarantees become extremely valuable.

---

# 2. Then applications started changing

Modern applications began producing enormous amounts of data and traffic.

Think about:

```text
Social media
IoT
Online gaming
E-commerce
Analytics
Messaging
Streaming
Mobile applications
```

A service might have:

```text
10 million users
100 million events/day
thousands of requests/second
users distributed globally
```

And importantly, these systems often need to operate across **multiple machines**.

---

# 3. Vertical scaling has limits

One approach is:

> Buy a bigger server.

For example:

```text
Server
────────────────
8 CPU
32 GB RAM
1 TB SSD
```

becomes:

```text
Server
────────────────
64 CPU
512 GB RAM
20 TB SSD
```

This is **vertical scaling**.

It can work extremely well.

But eventually you encounter:

- hardware limits
- increasing cost
- a large failure domain
- difficult upgrades
- insufficient capacity for some workloads

So another approach is:

> Use multiple machines.

```text
        Application
             │
       ┌─────┼─────┐
       ▼     ▼     ▼
     Node1 Node2 Node3
```

That's **horizontal scaling**.

And now the database has a much harder problem.

---

# 4. Distributed databases are fundamentally harder

With one machine:

```text
Application
    │
    ▼
Database
```

The database knows where its data is.

With multiple machines:

```text
                 Application
                     │
             ┌───────┼───────┐
             ▼       ▼       ▼
           Node A  Node B  Node C
```

Now you have to think about:

- Which node owns the data?
- How do nodes communicate?
- What happens if Node B dies?
- What if the network breaks?
- How do replicas stay synchronized?
- What happens if two nodes receive conflicting updates?
- How do you maintain consistency?
- How do you distribute the workload?
- How do you add more nodes?

These are **distributed-systems problems**.

And this is where NoSQL becomes much more interesting.

---

# 5. NoSQL wasn't one invention

This is another misconception to eliminate early.

"NoSQL" doesn't mean:

> "A database that doesn't use SQL."

It is an umbrella term covering several fundamentally different models.

```text
NoSQL
 │
 ├── Document
 │      MongoDB
 │
 ├── Key-value
 │      Redis
 │      DynamoDB
 │
 ├── Wide-column
 │      Cassandra
 │      ScyllaDB
 │
 └── Graph
        Neo4j
```

They solve different problems.

There isn't one universal **NoSQL architecture**.

---

# 6. Different workloads → different models

Imagine a system storing:

```text
User
 ├── profile
 ├── preferences
 ├── addresses
 └── settings
```

A document database can represent that naturally:

```json
{
  "id": 42,
  "name": "Alice",
  "preferences": {
    "language": "en",
    "currency": "EUR"
  }
}
```

---

Now imagine:

```text
sessionId → sessionData
```

You don't necessarily need complex relationships.

A key-value database can be ideal:

```text
"session:abc123" → {...}
```

---

Now imagine billions of events:

```text
timestamp
deviceId
temperature
```

with enormous write volume.

A wide-column database may be a better fit.

---

Now imagine:

```text
Alice → knows → Bob
Bob → works_at → Company X
Alice → bought → Product Y
Product Y → belongs_to → Category Z
```

A graph database can represent those relationships directly.

---

# 7. The big conceptual shift

Relational databases traditionally encourage you to start with:

> **What data do I have, and how should I normalize it?**

Many NoSQL systems encourage you to start with:

> **What queries/access patterns does my application need?**

For example:

```text
Requirement:

"Get a user's latest 20 notifications."
```

Instead of beginning with:

```text
What tables should I create?
```

you might begin with:

```text
How will this query be executed?
What partition should contain the data?
What index/key should I use?
How many nodes will I touch?
Can I retrieve it in one operation?
```

This is a major shift in thinking.

---

# 8. And there's a price

NoSQL doesn't magically eliminate complexity.

It often **moves the complexity somewhere else**.

For example:

### Relational approach

You normalize:

```text
Customer
Order
OrderItem
Product
```

and let JOINs reconstruct information.

### Document approach

You might denormalize:

```text
Order
 ├── customer information
 ├── items
 └── product information
```

Now reading can be simpler/faster.

But updating the customer's information may become more complicated because copies might exist in multiple documents.

So:

```text
Relational
    ↓
More normalization
    ↓
More relationships / JOINs

Document
    ↓
More denormalization
    ↓
Potential duplication / consistency concerns
```

Neither is universally superior.

---

# 9. This leads directly to distributed NoSQL

The really important progression is:

```text
Traditional database
        ↓
Large workload
        ↓
Need more capacity
        ↓
Multiple machines
        ↓
Partitioning
        ↓
Replication
        ↓
Distributed systems problems
        ↓
Different consistency/availability trade-offs
        ↓
Different database models
```

And that's why the next module is so important.

---

# Module 0 complete

You now have the baseline we need:

```text
Database
    ↓
Database engine
    ↓
Storage engine
    ↓
Query processing
    ↓
Indexes
    ↓
Concurrency
    ↓
Durability
    ↓
Relational model
    ↓
JOINs / normalization
    ↓
Distributed scaling
    ↓
Why different NoSQL models exist
```

## One thing I want you to remember

Don't think:

> **SQL vs NoSQL**

Think:

> **Which data model and distributed architecture best fit the workload and its requirements?**

That mindset will serve you much better at senior/principal level.

---

## Next: Module 1 — Why NoSQL Exists

We'll start getting into the actual foundations:

1. **What problems traditional databases encounter at scale**
2. Horizontal vs vertical scaling in detail
3. Partitioning/sharding
4. Replication
5. Read/write scaling
6. Distributed database architecture
7. CAP theorem
8. BASE
9. PACELC
10. The trade-offs that actually matter in production

And importantly, **we'll start hands-on work once we reach the first concrete database technology**, rather than spending the entire course in theory.