## 10. Memory grants, spills, and `tempdb`

This is one of those SQL Server topics that makes a lot more sense once you understand **what SQL Server is actually trying to do with memory**.

The key idea is:

> SQL Server sometimes needs temporary working memory to execute a query, in addition to the memory required simply to hold the data being returned.

That working memory is called a **memory grant**.

---

### 1. Why does a query need extra memory?

Consider:

```sql
SELECT *
FROM Orders
ORDER BY CreatedAt;
```

SQL Server needs to produce the rows in `CreatedAt` order.

If the rows aren't already available in that order through an appropriate index, SQL Server may need to perform a **Sort**.

Sorting requires working memory.

For example, conceptually, SQL Server may have to hold a substantial number of rows while it determines their correct order.

Similarly, a **Hash Join** needs memory for its hash table.

A **Hash Aggregate** needs memory to maintain its grouping information.

So these operators can request memory before execution.

---

# 2. What is a memory grant?

Before executing certain operators, SQL Server's optimizer estimates how much memory the query will need and requests that amount from the SQL Server memory manager.

For example, the optimizer might determine:

> "I expect this sort to involve 100,000 rows, and I'll need approximately 50 MB."

SQL Server grants memory to the query.

The query then uses that memory for its execution.

So there are really two different things:

**Buffer pool memory**

Used extensively for caching database pages.

**Query execution memory**

Temporarily granted to operators such as sorts and hash operations.

You don't want to think of SQL Server memory as one giant bucket where every operation uses memory in exactly the same way.

---

# 3. Why does the optimizer have to estimate memory?

Because the query hasn't run yet.

Suppose:

```sql
SELECT *
FROM Listings
WHERE CategoryId = 10
ORDER BY Price;
```

The optimizer needs to estimate how many rows will reach the Sort.

Suppose it estimates:

```text
Estimated rows: 10,000
```

It can calculate a memory requirement based on that estimate.

But what if the actual result is:

```text
Actual rows: 5,000,000
```

Now the query needs substantially more working memory than SQL Server expected.

And this leads to an important problem.

---

# 4. What is a spill?

If an operator needs more memory than it was granted, it can sometimes use **`tempdb`** to temporarily store intermediate data.

That's called a **spill**.

For example, a Sort might not fit entirely in memory.

Instead of keeping everything in RAM, SQL Server can write portions of the intermediate data to `tempdb`, process them in pieces, and eventually produce the required result.

The query can still succeed.

But it can become significantly slower because you're now doing additional I/O and processing.

This is a key point:

> **A spill isn't necessarily a query failure. It's usually a sign that an operator couldn't perform all of its work in the memory it had available.**

---

# 5. A concrete Sort example

Imagine:

```sql
SELECT *
FROM Listings
ORDER BY Price;
```

There are 10 million rows.

SQL Server estimates:

```text
100,000 rows
```

and grants memory appropriate for roughly that amount.

During execution it discovers:

```text
10,000,000 rows
```

need to be sorted.

It can't keep the entire working set in its granted memory.

So parts of the sort can spill to `tempdb`.

Instead of:

```text
RAM → sort → result
```

you effectively end up doing something more like:

```text
RAM + temporary storage → multiple stages of sorting → final result
```

The important part isn't the diagram; it's the mechanism:

**the operator breaks the work into manageable pieces and uses temporary storage because the available working memory isn't sufficient.**

---

# 6. Hash joins can spill too

Remember the Hash Join from the previous topic.

SQL Server builds a hash structure from one side of the join.

Suppose it estimates:

```text
100,000 rows
```

but there are actually:

```text
10,000,000 rows
```

The hash structure may exceed the available memory.

SQL Server can partition the work and use `tempdb` for portions that don't fit in memory.

That's a **hash spill**.

So when you see a Hash Join with spills, don't immediately conclude:

> "Hash Join is bad."

Instead ask:

> "Why did SQL Server underestimate how much data this operation would process?"

And that brings us directly back to **cardinality estimation and statistics**.

---

# 7. This is why inaccurate statistics can have multiple consequences

Earlier we talked about:

```text
Estimated rows: 100
Actual rows: 10,000,000
```

Now you can see that this can affect more than just the choice of join algorithm.

It can also affect:

### Join choice

SQL Server may choose Nested Loops when Hash Join would have been better.

### Memory grant

SQL Server may allocate too little memory.

### Sort

The sort may spill to `tempdb`.

### Hash operation

The hash table may spill to `tempdb`.

So one bad cardinality estimate can have a cascade of consequences.

---

# 8. Can SQL Server grant too much memory?

Yes.

This is the other side of the problem.

Suppose SQL Server estimates:

```text
10,000,000 rows
```

but the actual query produces:

```text
10,000 rows
```

SQL Server might request a very large memory grant that isn't actually needed.

That memory is then unavailable to other queries while the grant is being held.

Imagine multiple concurrent queries doing this.

You can end up with:

> Queries waiting for memory grants even though the server still has substantial memory activity elsewhere.

This is called **memory grant contention**.

So both extremes can be problematic:

**Too little memory**

→ spills and slower execution.

**Too much memory**

→ other queries may have to wait for memory.

---

# 9. This is another reason estimates matter

You can now connect the entire chain:

```text
Statistics
    ↓
Cardinality estimate
    ↓
Execution plan
    ↓
Memory requirement estimate
    ↓
Memory grant
    ↓
Actual execution
```

If the cardinality estimate is badly wrong, the downstream decisions can also be wrong.

This is why I wanted to establish statistics before teaching memory grants.

---

# 10. What exactly is `tempdb`?

`tempdb` is a special SQL Server system database.

Unlike your application's database, it is intended for **temporary and intermediate work**.

SQL Server uses it for many things, including:

- temporary tables
- table variables in certain circumstances
- temporary internal structures
- sorting that spills to disk
- hash spills
- row versioning in relevant isolation-level scenarios
- various other internal operations

And `tempdb` is shared by the SQL Server instance.

That last point matters.

If one query is causing huge amounts of temporary work, it can affect other workloads using `tempdb`.

---

# 11. `tempdb` isn't necessarily "disk = bad"

Another oversimplification to avoid:

> "If tempdb is involved, the query is bad."

No.

SQL Server legitimately uses `tempdb` for many normal operations.

For example, temporary tables are explicitly designed to use it.

The problem is generally **excessive or unexpected temporary work**.

A query that spills a few small amounts of data isn't necessarily something you need to panic about.

A high-volume query repeatedly spilling gigabytes to `tempdb` is a different story.

---

# 12. Why does `tempdb` sometimes become a bottleneck?

Imagine dozens of concurrent queries performing:

- large sorts
- hash joins
- temporary table operations
- row-versioning operations

All of them can place pressure on `tempdb`.

Now the problem isn't necessarily:

> "This individual query is slow."

It can become:

> **"Multiple workloads are competing for the same temporary storage resources."**

You can therefore have a perfectly reasonable query that becomes slow because the server is under `tempdb` pressure.

---

# 13. Memory grants and concurrency

Here's an important scenario.

Suppose you have 100 simultaneous requests.

Each query gets a large memory grant.

Even if the server has enough memory for one query, it may not have enough to satisfy all 100 grants simultaneously.

Some queries can be forced to wait.

This is why a query that performs well in isolation can behave badly under production load.

Performance testing only one request at a time doesn't necessarily reveal this.

**Hence it is advisable to perform performance testing under laod**

---

# 14. Why `SELECT *` can make this worse

Suppose you're sorting:

```sql
SELECT *
FROM Listings
ORDER BY Price;
```

and `Listings` contains:

```text
Id
CategoryId
Price
Title
Description
CreatedAt
...
```

The rows being manipulated may be considerably wider than if you only needed:

```sql
SELECT Id, Price
FROM Listings
ORDER BY Price;
```

Wider rows can increase memory requirements for operations that need to hold or manipulate those rows.

This is one of the practical reasons that:

> **Retrieving unnecessary columns isn't just about network traffic. It can affect query processing as well.**

It also connects to the projection discussion we had earlier.

---

# 15. How do you identify a spill?

When inspecting an actual execution plan in SQL Server Management Studio, certain operators can indicate that they spilled.

You may see warnings associated with:

- Sort
- Hash Match
- other operators capable of spilling

The actual execution plan can provide information about the spill and the number of passes involved.

You can also investigate runtime/query-performance tooling for memory-grant information.

The important diagnostic process is:

> **Find the spill → determine why the operator needed more memory than expected → investigate the cardinality estimate → determine whether statistics, predicates, indexing, or the query itself can be improved.**

Don't simply throw more RAM at it.

---

# 16. Can adding an index eliminate a spill?

Sometimes.

Suppose:

```sql
SELECT *
FROM Listings
ORDER BY CreatedAt;
```

requires a huge Sort.

If you create an appropriate index that provides the rows in `CreatedAt` order, SQL Server may no longer need to perform that expensive Sort.

That can eliminate the memory requirement associated with the Sort.

But that's not the same as:

> "Indexes reduce memory."

The more precise statement is:

> **An appropriate access path can allow SQL Server to avoid an operation that otherwise requires substantial working memory.**

That's the kind of distinction I want you to make.

---

# 17. Can rewriting the query help?

Yes.

For example, if you're sorting millions of rows but only need the first 20:

```sql
SELECT TOP 20 *
FROM Listings
ORDER BY CreatedAt DESC;
```

SQL Server can potentially use a suitable index to avoid processing everything.

Likewise, filtering data earlier can reduce the number of rows reaching a sort, join, or aggregate.

The general principle is:

> **Reduce unnecessary data as early as possible, while allowing the optimizer to choose an efficient access path.**

---

# 18. Why `tempdb` comes up in interviews

If someone asks:

> "Why would a query suddenly become very slow?"

A strong answer isn't just:

> "Maybe there's an index missing."

You should also be able to think about:

- bad cardinality estimates
- excessive memory requirements
- sort spills
- hash spills
- memory-grant contention
- `tempdb` contention
- blocking
- parameter-sensitive plans
- stale statistics

You don't need to claim that any one of these is definitely the problem without evidence.

You say:

> **"I'd inspect the execution plan and runtime statistics to determine which resource or operator is actually responsible."**

---

# 19. One subtle point: memory grant ≠ all memory used by the query

This distinction is important.

A query can use memory in different ways.

The **memory grant** specifically concerns memory SQL Server grants for certain query execution operators.

It doesn't mean:

> "This is the total memory consumed by the query."

SQL Server also uses memory for things such as:

- buffer pool pages
- cached execution plans
- internal structures
- other engine components

So don't look at:

```text
Memory grant = 500 MB
```

and conclude:

> "This query consumes exactly 500 MB."

That's not what the number means.

---

# 20. The deeper picture

At this point, several things we've covered start fitting together.

Suppose your query has a Hash Join.

SQL Server needs to decide:

> How many rows will participate?

Statistics help answer that.

Then:

> Which join algorithm should I use?

The optimizer considers cardinality and cost.

Then:

> How much memory will the hash operation need?

The optimizer estimates that too.

Then:

> Can the operation fit in its memory grant?

If not:

> It may spill into `tempdb`.

So a performance problem that initially looks like:

> "The Hash Join is slow."

might actually originate much earlier:

> **The statistics caused a bad cardinality estimate, which led to a poor memory estimate and an execution that spilled.**

That's the kind of causal chain you want to be able to reason through.

---

### Next topic

I'd go one level deeper into **query execution and concurrency: locking, blocking, deadlocks, and isolation levels**.

You've already learned optimistic/pessimistic concurrency and some locking concepts, but we haven't tied them together at the **SQL Server engine level**—what locks actually protect, why one query blocks another, how deadlocks arise, and how isolation levels change the behavior.