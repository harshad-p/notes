# Advanced SQL Performance Diagnosis

The key change from what we've done before is this:

Previously, we talked about **how to make queries faster**.

Now we're going to focus on:

> **How do I figure out *why* this particular query is slow?**

That's a different skill.

A query being slow doesn't immediately tell you what is wrong.

For example, this:

```sql
SELECT *
FROM Listings
WHERE CategoryId = 42;
```

could be slow because:

- there's no useful index
- SQL Server chose the wrong index
- SQL Server estimated the number of rows incorrectly
- the query is doing an expensive sort
- the query is waiting for another transaction
- the disk/storage is slow
- the server is CPU-bound
- the query is suffering from parameter sniffing
- you're returning millions of rows
- the application is slow rather than SQL Server
- something else is consuming the database resources

So the first senior-level habit is:

> **Don't immediately change the query or add an index. Diagnose first.**

---

# 1. Execution plans

This is where we start.

When SQL Server receives a query, it doesn't simply execute the SQL text literally from top to bottom.

It determines **how** to retrieve the requested data.

For example:

```sql
SELECT Name, Price
FROM Listings
WHERE CategoryId = 42;
```

SQL Server has choices.

It could scan the entire table:

> "Read every listing and check CategoryId."

Or it could use an index:

> "Use the CategoryId index to find the relevant rows."

Or it could use some combination of indexes and other operations.

The **execution plan** describes the strategy SQL Server chose.

---

## Estimated vs actual execution plan

This distinction is important.

### Estimated execution plan

SQL Server says:

> "Based on what I know, I expect this plan to work well."

It hasn't actually executed the query.

### Actual execution plan

SQL Server executes the query and then gives you information about what actually happened.

This is much more useful for diagnosing problems.

For example, SQL Server might estimate:

> 10 rows

but actually process:

> 500,000 rows

That's a huge clue.

---

# 2. Why estimates matter

Imagine SQL Server is deciding between two strategies.

It estimates:

```text
Plan A → process 20 rows
Plan B → process 500,000 rows
```

Naturally, it may choose Plan A.

But suppose reality is:

```text
Plan A → actually processes 500,000 rows
Plan B → actually processes 500,000 rows
```

Now the optimizer made a bad decision because its estimate was wrong.

And that can happen because of things like:

- outdated statistics
- data distribution
- parameter sniffing
- expressions that make predicates difficult to estimate
- correlations between columns

This is why **estimated rows vs actual rows** is one of the first things I look at when investigating an execution plan.

---

# 3. Reading an execution plan

You don't need to memorize every operator.

Start with the important ones.

### Table Scan

SQL Server reads the table.

This isn't automatically bad.

If the table has 20 rows, scanning it is perfectly reasonable.

The question is:

> **Was scanning the table appropriate for this query?**

---

### Index Seek

SQL Server uses an index to directly locate relevant rows.

Generally desirable when you're looking for a relatively small portion of a large table.

---

### Index Scan

SQL Server reads a substantial portion of an index.

Again, **not automatically bad**.

Suppose you're returning 70% of a table.

Scanning an index may be perfectly sensible.

---

### Key Lookup

This one is particularly useful to understand.

Suppose you have:

```sql
CREATE INDEX IX_Listings_CategoryId
ON Listings(CategoryId);
```

And your query is:

```sql
SELECT Name, Price
FROM Listings
WHERE CategoryId = 42;
```

The index can find the relevant rows through `CategoryId`.

But the index doesn't contain:

- `Name`
- `Price`

So SQL Server may need to go back to the actual table for each matching row.

That's a **Key Lookup**.

If only 5 rows match, no big deal.

If 500,000 rows match, you could have a serious problem.

That's one situation where a covering index could help—which connects directly to the indexing work we already did.

---

# 4. Sort

Suppose you run:

```sql
SELECT *
FROM Listings
WHERE CategoryId = 42
ORDER BY Price;
```

SQL Server may need to retrieve the rows and then sort them.

A large sort can consume:

- CPU
- memory

And if insufficient memory is available, parts of the operation can spill to tempdb.

So when looking at a plan, an expensive **Sort** can be a clue.

---

# 5. Hash Match

You'll frequently encounter Hash Match in joins and aggregations.

For example:

```sql
SELECT c.Name, AVG(l.Price)
FROM Categories c
JOIN Listings l
    ON c.Id = l.CategoryId
GROUP BY c.Name;
```

SQL Server may choose a hash-based strategy.

Again, seeing a Hash Match doesn't mean:

> "Hash Match = bad."

You need to understand:

> Why did SQL Server choose it, and how much work is it doing?

---

# 6. The most important thing: cost

Execution plans often show percentages such as:

> 65%  
> 25%  
> 10%

These are **estimated relative costs within that execution plan**.

Don't interpret:

> "65% means this operation uses 65% of the server's CPU."

It doesn't.

It's telling you that SQL Server's optimizer estimated this operation to account for a particular proportion of the **query's estimated cost**.

This is a common interview trap.

---

# 7. A practical example

Imagine your query is:

```sql
SELECT Name, Price
FROM Listings
WHERE CategoryId = 42
ORDER BY Price;
```

You inspect the actual execution plan and discover:

```text
Estimated rows: 100
Actual rows:    400,000
```

That's immediately interesting.

Then you see a large Key Lookup.

Now you have a hypothesis:

> SQL Server thought only ~100 rows would be returned, so the lookup strategy looked cheap. In reality, 400,000 rows were returned, causing huge lookup work.

Now you're not blindly saying:

> "Add an index."

You're reasoning from evidence.

You might then consider whether a covering index is appropriate.

That's what I mean by **performance diagnosis**.

---

# 8. Don't stop at the execution plan

This is extremely important.

Suppose someone tells you:

> "This query takes 10 seconds."

You shouldn't automatically assume the query itself is doing 10 seconds of work.

It could be:

```text
Query execution:      500 ms
Waiting for lock:    9,000 ms
```

The SQL statement isn't necessarily computationally expensive.

It's **waiting**.

Or:

```text
Query execution: 8 seconds
CPU:              200 ms
I/O:              huge
```

Now your investigation goes in a completely different direction.

Or:

```text
Query execution: 8 seconds
CPU:              7.8 seconds
```

Now you're looking at CPU-intensive processing.

So we're going to learn to separate:

**"The query is expensive"**

from

**"The query is waiting."**

That's the next major piece of this topic.

---

## The mental model I want you to develop

When someone hands you a slow SQL query, don't immediately think:

> "Which index should I add?"

Think:

> **"Where is the time actually going?"**

Then investigate:

1. **What plan did SQL Server choose?**
2. **What did SQL Server estimate?**
3. **What actually happened?**
4. **Is it doing too much CPU work?**
5. **Is it doing too much I/O?**
6. **Is it waiting on another resource, such as a lock?**
7. **Is the plan inappropriate for the parameters/data distribution?**
8. **Only then: what change would address the actual problem?**

That's the foundation for the rest of this section.

Next, we'll go into **SARGability**, because it explains a surprisingly large number of cases where you have an index but SQL Server still can't use it effectively.