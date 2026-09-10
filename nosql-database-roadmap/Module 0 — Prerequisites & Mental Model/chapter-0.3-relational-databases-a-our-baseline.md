# Module 0 — Chapter 0.3
## Data Storage vs. Database + Relational Database Baseline

Before we start talking about NoSQL, I want to make sure the distinction between **storage**, **database**, and **database engine** is completely solid. This becomes important later when we discuss MongoDB, Cassandra, Redis, etc.

---

## 1. Data storage is not a database

Suppose you have:

```text
users.json
```

containing:

```json
[
  { "id": 1, "name": "Alice" },
  { "id": 2, "name": "Bob" }
]
```

This is **data stored somewhere**.

You could put it on:

- an SSD
- a hard disk
- an S3 object
- a filesystem
- a USB drive
- etc.

But that doesn't make it a database.

A database provides a system for **managing data**, not merely storing bytes.

For example, it needs to answer questions like:

> Find user with ID 2.

> Give me all orders from customer 123.

> Two requests updated the same record simultaneously. What happens?

> The server crashed halfway through a write. What should the data look like after recovery?

> How can I find this record without examining 50 million records?

That's where database software comes in.

---

# 2. Three layers you should distinguish

Think of this:

```text
                    DATABASE SYSTEM
                         │
              ┌──────────┴──────────┐
              │                     │
        Database Engine        Data Storage
              │                     │
      query processing       SSD / disk / etc.
      indexes
      transactions
      concurrency
      recovery
      caching
```

More concretely:

### Storage

Actually holds bytes.

```text
SSD
 └── files / blocks / pages
```

### Database engine

Software that manages those bytes intelligently.

Examples:

- PostgreSQL
- SQL Server
- MongoDB
- Cassandra
- Redis

### Database

The logical collection of data managed by the engine.

For example:

```text
PostgreSQL server
    ├── database: Hotel
    │      ├── customers
    │      ├── reservations
    │      └── rooms
    │
    └── database: Reporting
           ├── ...
```

The terminology varies somewhat between products, which is why you shouldn't become overly attached to the exact word "database."

---

# 3. Now establish the relational baseline

You already know SQL, so we're going to use relational databases as our reference point.

Suppose we have:

```text
Customers
--------------------------------
Id       Name
1        Alice
2        Bob
3        Charlie
```

and:

```text
Orders
--------------------------------
Id       CustomerId       Amount
101      1                200
102      1                150
103      2                500
```

Conceptually:

```text
Customers
    │
    │  CustomerId
    ▼
Orders
```

`Customers.Id` is typically the **primary key**.

`Orders.CustomerId` is a **foreign key** referring to it.

---

# 4. Why relational databases use tables

The relational model organizes data into **relations**, which we normally visualize as tables.

A table has:

```text
columns → structure
rows    → individual records
```

For example:

```text
Orders
+-----+------------+--------+
| Id  | CustomerId | Amount |
+-----+------------+--------+
| 101 | 1          | 200    |
| 102 | 1          | 150    |
| 103 | 2          | 500    |
+-----+------------+--------+
```

This gives us a strong structure.

Every row follows the same basic schema:

```text
Id
CustomerId
Amount
```

---

# 5. Normalization

Relational databases traditionally encourage **normalization**.

For example, instead of:

```text
Order
--------------------------------
OrderId
CustomerName
CustomerEmail
ProductName
ProductPrice
```

you might have:

```text
Customer
---------
Id
Name
Email

Order
---------
Id
CustomerId

Product
---------
Id
Name
Price

OrderItem
---------
OrderId
ProductId
Quantity
```

Why?

To avoid unnecessary duplication and maintain consistency.

If Alice changes her email address, you ideally update it **once** in `Customer`.

---

# 6. But then we need JOINs

Now imagine you want:

> Give me all orders along with the customer's name.

The database needs to combine information:

```sql
SELECT
    o.Id,
    c.Name,
    o.Amount
FROM Orders o
JOIN Customers c
    ON c.Id = o.CustomerId;
```

Conceptually:

```text
Orders
   │
   │ CustomerId
   ▼
Customers
```

JOINs are extremely powerful.

But they also become interesting when data is distributed across multiple machines.

Suppose:

```text
Server A
Customers

Server B
Orders
```

Now a JOIN potentially becomes a **distributed operation**.

That can be substantially more expensive than joining data on one machine.

This is one of the important ideas that eventually leads us toward NoSQL modeling.

---

# 7. Relational databases aren't "bad at scaling"

This is important.

Don't fall into the simplistic:

> SQL doesn't scale, NoSQL scales.

That's wrong.

Relational databases can scale extremely well.

They can use:

- indexes
- caching
- read replicas
- partitioning
- sharding
- connection pooling
- query optimization
- vertical scaling
- distributed architectures

Databases such as PostgreSQL, SQL Server, and others can handle very large workloads.

The question is:

> **What kind of workload and architecture are you dealing with?**

---

# 8. The fundamental trade-off we're approaching

Consider an application where the most common operation is:

```text
Get reservation by ReservationId
```

and the reservation needs:

```text
Reservation
 ├── guest
 ├── hotel
 ├── room
 ├── dates
 ├── payment
 └── preferences
```

One possible relational design is:

```text
Reservations
Guests
Hotels
Rooms
Payments
Preferences
```

and retrieving the complete reservation may involve multiple JOINs.

A document database might instead store something closer to:

```json
{
  "reservationId": 123,
  "guest": {
    "name": "Alice",
    "email": "alice@example.com"
  },
  "hotel": {
    "id": 42,
    "name": "Hotel Berlin"
  },
  "room": {
    "number": 301,
    "type": "Deluxe"
  },
  "payment": {
    "status": "Paid"
  }
}
```

Now:

```text
Get reservation 123
```

can potentially retrieve most of what the application needs as **one document**.

That doesn't automatically make the document model better.

We've simply changed the trade-off.

And understanding **why** we might make that trade-off is the heart of NoSQL data modeling.

---

# 9. The mental model to keep

For the rest of this course, keep these four concepts separate:

```text
Storage
   ↓
Where bytes ultimately live

Database engine
   ↓
Software managing those bytes and providing database behavior

Data model
   ↓
How we logically represent data
   ├── relational
   ├── document
   ├── key-value
   ├── wide-column
   └── graph

Database system
   ↓
The complete system providing storage,
querying, concurrency, durability,
replication, etc.
```

This distinction will become **very important** once we get to MongoDB's document model versus WiredTiger's storage engine, Cassandra's data model versus its SSTables, and Redis's in-memory architecture.

---

## Quick check

Before we move to Chapter 0.4, answer these in your own words:

**1.** Is an SSD a database? Why or why not?

**2.** What's the difference between a **database engine** and a **storage engine**?

**3.** Why can JOINs become more problematic when data is distributed across multiple machines?

**4.** Why might a NoSQL document intentionally contain duplicated data that would normally be separated into tables in a relational database?

---

### 1. SSD vs. database

The **database is not necessarily just the collection of data**. In practical engineering usage, "database" can refer to both the logical data and the system managing it, depending on context.

```text
SSD
 └── stores bytes

Database
 ├── organizes data
 ├── retrieves data
 ├── indexes
 ├── handles concurrency
 ├── provides durability
 ├── handles recovery
 └── etc.
```

---

### 2. Database engine vs. storage engine

The **storage engine doesn't necessarily manage RAM as a persistent storage medium**. It manages how data is persisted and retrieved, while using RAM extensively for caching/buffering.

Conceptually:

```text
Database engine
 ├── query parser
 ├── query planner/optimizer
 ├── execution engine
 ├── transaction/concurrency mechanisms
 ├── indexes
 └── storage engine
       ├── pages/files
       ├── reads/writes
       ├── caching/buffering
       ├── journaling/WAL
       └── recovery
```

The exact architecture varies by database.

---

### 3. JOINs in distributed systems

There's an important distinction:

A JOIN doesn't necessarily mean **different databases**. It could be different **nodes/shards of the same distributed database**.

For example:

```text
                 JOIN
                  │
        ┌─────────┴─────────┐
        ▼                   ▼
     Node A              Node B
   Customers              Orders
```

Now the system may have to:

1. retrieve data from Node A
2. retrieve data from Node B
3. move data across the network
4. perform the JOIN
5. deal with failures/timeouts/etc.

That can be considerably more expensive than:

```text
One machine
    │
Customers + Orders
    │
    └── JOIN
```

This idea will come back repeatedly when we discuss **partitioning and distributed data modeling**.

---

### 4. Why duplicate data in NoSQL?

The key phrase:

> **This is a trade-off.**

That's one of the most important sentences in this entire course.

Relational modeling often asks:

> "How do I avoid unnecessary duplication and preserve consistency?"

NoSQL modeling often asks:

> **"How should I organize my data so that my important access patterns are efficient?"**

Those are related, but different optimization goals.

For example, if 95% of requests are:

```text
Get reservation by ID
```

then having the reservation information together may be worth duplicating some data.

But then you inherit a new problem:

> **What happens when the duplicated data changes?**

---