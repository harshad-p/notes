# SARGability

SARGability is one of those SQL terms that sounds much more complicated than the underlying idea.

**SARG** comes from **Search ARGument**.

A predicate is **SARGable** when SQL Server can use it efficiently to search an index rather than having to calculate something for every row first.

The important thing is not simply:

> "Does this query have an index?"

It's:

> **"Can SQL Server use the index to narrow down the rows efficiently?"**

---

## A simple example

Suppose we have:

```sql
CREATE INDEX IX_Users_Email
ON Users(Email);
```

And we query:

```sql
SELECT *
FROM Users
WHERE Email = 'harshad@example.com';
```

This is SARGable.

SQL Server can essentially use the index to locate the relevant part of the index directly.

Now consider:

```sql
SELECT *
FROM Users
WHERE LOWER(Email) = 'harshad@example.com';
```

The problem is that we've put a function around the indexed column.

SQL Server can't simply ask the index:

> "Find me the value `harshad@example.com`."

It needs to determine the `LOWER()` value for the rows before it can compare them.

That can prevent an efficient index seek and result in much more work.

---

# Why this matters

Imagine there are **10 million users**.

With:

```sql
WHERE Email = 'harshad@example.com'
```

the index can potentially take SQL Server directly to the relevant location.

With:

```sql
WHERE LOWER(Email) = 'harshad@example.com'
```

SQL Server may have to examine a huge number of rows and evaluate `LOWER(Email)`.

That's a massive difference.

So SARGability is essentially about preserving SQL Server's ability to **navigate an index using the predicate**.

---

# Another common example: dates

Suppose:

```sql
CREATE INDEX IX_Orders_OrderDate
ON Orders(OrderDate);
```

And we want all orders from 2025:

```sql
SELECT *
FROM Orders
WHERE YEAR(OrderDate) = 2025;
```

This looks perfectly reasonable.

But again, we're applying a function to the indexed column.

SQL Server can't simply seek to:

> "the beginning of 2025"

based on that expression.

A better formulation is:

```sql
SELECT *
FROM Orders
WHERE OrderDate >= '2025-01-01'
  AND OrderDate <  '2026-01-01';
```

Now the predicate describes a **range of values in the index**.

That is much more naturally searchable.

---

# Another example: arithmetic on a column

Suppose:

```sql
WHERE Price * 1.2 > 1000
```

You're asking SQL Server to calculate `Price * 1.2` before determining whether the row qualifies.

You can often rearrange the condition:

```sql
WHERE Price > 833.33
```

Now SQL Server can search the `Price` index directly.

The exact mathematical transformation obviously depends on the expression, but the principle is important:

> **Try to put the transformation on the value you're comparing against rather than on the indexed column.**

---

# What about LIKE?

This is a particularly useful one.

Suppose:

```sql
WHERE Name LIKE 'Harsh%'
```

This can be SARGable because SQL Server knows that you're looking for values beginning with `"Harsh"`.

It can potentially seek into the appropriate part of the index.

But:

```sql
WHERE Name LIKE '%shad'
```

is fundamentally different.

SQL Server doesn't know where in the index the matching values begin because the beginning is unknown.

It may need to examine many values.

And:

```sql
WHERE Name LIKE '%arsh%'
```

has the same fundamental problem.

So:

```text
'Harsh%'
```

can generally use a normal B-tree index effectively.

```text
'%Harsh'
```

generally cannot use it for an efficient seek.

(This is one reason dedicated full-text/search systems exist for more sophisticated text searching.)

---

# Implicit conversions

Here's another one that can be surprisingly important.

Suppose your column is:

```sql
UserId INT
```

and you query:

```sql
WHERE UserId = '123'
```

SQL Server may implicitly convert the string to an integer, which is generally fine.

But if you have a situation where SQL Server has to convert the **column itself** for comparison, you can end up with poor index usage.

This is why **matching data types matters**.

For example, don't casually compare:

```text
VARCHAR column
```

against something that causes SQL Server to convert the entire column to another type.

---

# SARGability doesn't mean "every function is bad"

This is important.

Don't memorize:

> "Functions in WHERE clauses are bad."

That's too simplistic.

The actual question is:

> **"Does this expression prevent SQL Server from efficiently navigating the available index?"**

For example, there are ways to make some computed expressions indexable, such as indexed computed columns, depending on the situation.

And SQL Server can sometimes transform predicates internally.

So SARGability is about **how the optimizer can access the data**, not about following a syntax rule.

---

# How do I recognize a SARGability problem?

Suppose you have:

```sql
CREATE INDEX IX_Listings_CreatedAt
ON Listings(CreatedAt);
```

Then someone reports:

> "This query is slow."

```sql
SELECT *
FROM Listings
WHERE CAST(CreatedAt AS DATE) = '2026-09-01';
```

I'd immediately be suspicious.

I'd inspect the execution plan.

If I see something like an **Index Scan** rather than an efficient seek, and there are many rows, I'd investigate whether the expression is preventing efficient index access.

I'd potentially rewrite it as a range:

```sql
WHERE CreatedAt >= '2026-09-01'
  AND CreatedAt <  '2026-09-02'
```

Then compare the actual execution plans and runtime.

That's the important workflow:

**Don't assume. Test.**

---

# SARGability and our previous indexing discussion

This connects directly to what we've already covered.

Suppose you carefully designed this index:

```sql
CREATE INDEX IX_Listings_CategoryId
ON Listings(CategoryId);
```

But then write:

```sql
WHERE SomeFunction(CategoryId) = ...
```

Having the index doesn't guarantee SQL Server can exploit it effectively.

So there are really two separate questions:

> **1. Do I have an appropriate index?**

and:

> **2. Can my query predicate actually use that index efficiently?**

You need both.

---

# One subtle point: Scan ≠ bad

Don't fall into the trap of thinking:

> Seek = good  
> Scan = bad

That's not true.

Suppose your table contains 100 rows and you're asking for 90 of them.

An index seek followed by thousands of lookups might be considerably worse than simply scanning the table.

Likewise, if your query legitimately needs most of an index/table, a scan may be exactly what SQL Server should do.

So when you see a scan, ask:

> **"How many rows does the query need compared with how many rows exist?"**

and:

> **"What is the actual cost of the alternative?"**

Again, diagnosis rather than memorization.

---

# The interview-level answer

If someone asks:

> **"What is SARGability?"**

I'd answer:

> "SARGability describes whether a predicate can be efficiently used by the database engine to search an index. For example, if I have an index on `OrderDate`, `WHERE OrderDate >= '2025-01-01' AND OrderDate < '2026-01-01'` is much more index-friendly than `WHERE YEAR(OrderDate) = 2025`, because applying the function to the indexed column can prevent an efficient index seek. I'd look at the execution plan to verify the actual behavior rather than assuming that every function necessarily causes a problem."

That's enough for an interview.

But the deeper thing to remember is:

> **SARGability is about preserving the database engine's ability to turn your predicate into an efficient search through the index.**

---

### Where we are in the plan

We've now covered:

1. ✅ Execution plans in practice
2. ✅ SARGability
3. **Next: Parameter sniffing / parameter-sensitive plans**
4. Identifying the actual bottleneck

Then we'll move on to **transactions in depth**.