# SQL Server Practical Troubleshooting

Now we move from **understanding SQL performance problems** to **actually finding them in a running SQL Server system**.

The three tools worth knowing for interviews are:

1. **Query Store** — which queries are slow/expensive, and how their performance changes over time.
2. **DMVs** — what SQL Server is doing right now and what it has observed.
3. **Extended Events** — detailed event-level diagnostics, especially for things like deadlocks.

We'll start with **Query Store**, because it's probably the most useful one to understand first.

---

## 1. Query Store

Imagine someone tells you:

> "Our application has become slow over the last two weeks."

You don't necessarily know:

- which query became slow,
- when it became slow,
- whether its execution plan changed,
- whether CPU increased,
- whether a particular query suddenly started doing much more work.

Query Store is designed to retain this kind of information.

Think of it as a **historical performance record of queries and their execution plans**.

### Without Query Store

You might inspect what's running **right now**.

But if the bad query ran yesterday at 2 PM and is fine now, you've missed it.

### With Query Store

You can investigate historical behavior:

```text
Query
  ↓
Execution statistics
  ↓
Execution plans
  ↓
Performance over time
```

This makes it particularly useful for production troubleshooting.

---

# 2. What does Query Store actually record?

At a high level, Query Store associates:

- queries,
- query plans,
- runtime statistics.

For example, suppose you have:

```sql
SELECT *
FROM Listings
WHERE City = @city;
```

Query Store might show that this query historically had:

```text
Average duration: 80 ms
CPU:              35 ms
Executions:       50,000
```

Then suddenly:

```text
Average duration: 2,800 ms
CPU:              1,900 ms
Executions:       50,000
```

That's a huge clue.

You can then investigate what changed.

---

# 3. Query Store is especially useful for plan changes

Remember **parameter sniffing** and execution plans from the previous section?

Suppose a query originally had:

```text
Plan A
Index Seek
Good performance
```

Then SQL Server starts using:

```text
Plan B
Large Index Scan
Poor performance
```

The query text may be exactly the same.

Query Store can help you see that **the query has multiple plans and that performance changed when a different plan was used**.

This is one of the most practical reasons to know Query Store.

---

# 4. A very realistic troubleshooting scenario

Imagine your API endpoint:

```text
GET /listings?city=Berlin
```

suddenly becomes slow.

You investigate Query Store and discover:

```text
Query:
SELECT ...
FROM Listings
WHERE City = @city

Plan 1:
Average duration = 50 ms

Plan 2:
Average duration = 3,200 ms
```

Now you have a much stronger hypothesis:

> "The query itself hasn't changed, but its execution plan changed."

You can then inspect the two plans and determine why.

Perhaps one plan uses an index seek while the other performs a large scan.

Or perhaps the query suffers from parameter-sensitive behavior because Berlin has millions of listings while another city has only a few hundred.

That's much more useful than blindly adding another index.

---

# 5. Query Store vs execution plan

They're related, but don't confuse them.

An **execution plan** tells you:

> "How is SQL Server executing this query?"

Query Store helps answer:

> "What queries and plans have been used, and how has their performance behaved over time?"

So:

**Execution plan = detailed execution strategy**

**Query Store = historical performance/plan evidence**

---

# 6. What would you actually look for?

If you're investigating slow SQL, useful things to rank by include:

### Duration

Which queries take the longest?

```text
Query A   50 ms
Query B   200 ms
Query C   5,000 ms   ← investigate
```

### CPU

Which queries consume the most CPU?

A query might be slow because it's doing computationally expensive joins, sorting, aggregation, etc.

### Logical reads

Which queries are reading huge amounts of data?

This can expose inefficient scans or poor indexing.

### Execution count

A query taking 10 ms isn't necessarily harmless.

If it runs **10 million times**, it may be a major contributor to system load.

This is an important point:

> **Don't only look for the slowest individual query. Look for queries that are expensive in aggregate.**

For example:

```text
Query A
10 seconds × 2 executions
= 20 seconds total

Query B
20 ms × 1,000,000 executions
= 20,000 seconds total
```

Query B could be much more important to the system despite being individually fast.

---

# 7. Query Store and production troubleshooting

A reasonable troubleshooting sequence is:

```text
Application is slow
       ↓
Is the DB actually the bottleneck?
       ↓
Identify expensive queries
       ↓
Query Store
       ↓
Look at duration / CPU / reads / execution count
       ↓
Inspect execution plan
       ↓
Compare good vs bad plans
       ↓
Form hypothesis
       ↓
Test change
```

Notice that this connects directly to what we learned earlier.

We aren't jumping straight to:

> "Add an index."

We're first asking:

> **"What is actually consuming the time/resources?"**

---

# 8. Query Store can also help after a deployment

Suppose you deploy a new application version.

Before deployment:

```text
Query X
Average duration = 40 ms
```

After deployment:

```text
Query X
Average duration = 1,500 ms
```

If Query Store has the relevant history, you have evidence that performance changed around that period.

Then you can compare:

- query text,
- plans,
- execution statistics.

This makes Query Store particularly useful for **regression detection**.

---

# 9. What about fixing a bad plan?

SQL Server has mechanisms for **plan forcing** through Query Store.

Suppose:

```text
Plan A → 50 ms
Plan B → 3,000 ms
```

and you've determined Plan A is reliably better.

You can potentially force SQL Server to use Plan A while you work on the underlying issue.

But this should not be treated as:

> "Always force the good-looking plan."

Because data distributions and workloads can change.

Plan forcing can be a **temporary mitigation**, while the real issue might be:

- statistics,
- indexing,
- parameter sensitivity,
- schema/query changes,
- data distribution.

---

# 10. Query Store vs DMVs

This distinction is important enough to remember.

### DMVs

Dynamic Management Views provide information about the **current state and accumulated server information**.

For example:

> What queries are consuming CPU right now?

or:

> What sessions are waiting?

or:

> What indexes are being used?

### Query Store

Better for:

> What happened historically with this query?

So if a query is slow **right now**, DMVs can be extremely useful.

If someone says:

> "This query was terrible yesterday afternoon, but it's fine now."

Query Store becomes especially valuable.

---

## 11. One interview-quality answer

If you're asked:

> **"How would you investigate a query that became slow in production?"**

I'd answer something like:

> "First I'd determine whether the database is actually the bottleneck by looking at application latency and database metrics. I'd then use Query Store to identify the expensive query and compare its historical execution statistics and plans. If the query has multiple plans, I'd check whether a plan regression or parameter-sensitive behavior is involved. I'd inspect the actual execution plan and look at CPU, logical reads, row estimates versus actual rows, scans, lookups, sorts, and waits. Then I'd make one targeted change and compare the measurements before and after."

That's a very solid senior-level answer.

---

### Next: DMVs

We'll look at the specific SQL Server DMVs you'd use to answer questions such as:

- **What queries are expensive right now?**
- **What is SQL Server waiting for?**
- **Which indexes are actually being used?**
- **Who is blocking whom?**
- **What sessions are currently running?**

Then we'll finish with **Extended Events** and a complete production troubleshooting workflow.