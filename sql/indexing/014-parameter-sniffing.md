# Parameter Sniffing

This is one of the more interesting SQL Server performance problems because the query itself can be perfectly reasonable, the indexes can be perfectly reasonable, and **SQL Server can still choose a terrible execution plan for some parameter values.**

The key is understanding that SQL Server usually **doesn't create a completely new execution plan every time you execute a parameterized query**.

---

## Start with a parameterized query

Imagine:

```sql
SELECT *
FROM Listings
WHERE CityId = @CityId;
```

Suppose we have:

- 10 million listings
- `CityId` is indexed

But the data isn't evenly distributed.

Maybe:

```text
Berlin      → 3,000,000 listings
Munich      →   500,000 listings
Bremen      →    20,000 listings
Small town  →       50 listings
```

The best way to execute the query can therefore depend heavily on `@CityId`.

For a city with **50 listings**, an index seek followed by lookups might be excellent.

For Berlin with **3 million listings**, that same strategy might be terrible. Scanning a large portion of the table could be cheaper.

So there isn't necessarily one universally optimal plan.

---

# What does SQL Server do?

When SQL Server first encounters a parameterized query, it has to compile an execution plan.

Suppose the first execution is:

```sql
@CityId = 123
```

and City 123 happens to have only 50 listings.

SQL Server examines the statistics and thinks:

> "Only a small number of rows will come back. An index seek looks like a good strategy."

It creates a plan based on that situation.

Then SQL Server can **cache that execution plan**.

Now someone executes the exact same query with:

```sql
@CityId = 42
```

where City 42 happens to be Berlin with 3 million listings.

SQL Server may reuse the previously compiled plan.

So you can end up with:

> Plan optimized for a tiny result set → reused for a huge result set.

That is **parameter sniffing**.

---

# Why is it called "sniffing"?

When SQL Server compiles the query, it can look at the parameter value that was supplied during compilation.

It effectively says:

> "I see that `@CityId` is 123. Based on the statistics, I estimate 50 rows."

It has **sniffed** the parameter value.

The problem isn't that SQL Server looked at the parameter.

That's actually useful.

The problem occurs when the plan created for that particular value isn't appropriate for subsequent values.

---

# A concrete example

Imagine this:

```sql
CREATE PROCEDURE GetListings
    @CityId INT
AS
BEGIN
    SELECT *
    FROM Listings
    WHERE CityId = @CityId;
END
```

First execution:

```sql
EXEC GetListings @CityId = 999;
```

City 999 has 20 listings.

SQL Server might choose:

> Index Seek → Key Lookups

Excellent.

Then:

```sql
EXEC GetListings @CityId = 1;
```

City 1 has 3 million listings.

The same plan might now perform millions of lookups.

Suddenly the query becomes very slow.

Then someone executes:

```sql
EXEC GetListings @CityId = 999;
```

again.

It may be fast.

So you can get the strange situation where:

> **"The query is sometimes fast and sometimes slow depending on the parameter."**

That's a classic clue.

---

# Why doesn't SQL Server simply compile a new plan every time?

Because compilation isn't free.

Imagine a high-traffic application executing a query thousands of times per second.

If SQL Server had to:

1. parse the query
2. optimize it
3. create a plan

for every execution, that would itself consume substantial CPU.

Plan caching allows SQL Server to reuse a compiled plan.

That's normally a **good thing**.

Parameter sniffing is essentially a consequence of that optimization when different parameter values have dramatically different data distributions.

---

# How do statistics fit into this?

This is important.

SQL Server uses **statistics** to estimate how many rows a predicate will return.

For example, statistics might tell it that:

```text
CityId = 999 → approximately 20 rows
```

and:

```text
CityId = 1 → approximately 3,000,000 rows
```

Those estimates influence the execution plan.

So parameter sniffing isn't really:

> "SQL Server randomly chose a bad plan."

It's more:

> **"SQL Server created a plan based on one parameter value, and that plan isn't appropriate for another parameter value."**

---

# How would you diagnose it?

Suppose an application reports:

> "This stored procedure is sometimes extremely fast and sometimes extremely slow."

That's a clue.

I'd investigate:

### 1. Look at the actual execution plans

Compare the plan for different parameter values.

If one parameter produces a good plan and another produces a terrible one, that's suspicious.

### 2. Compare estimated vs actual rows

For example:

```text
Estimated: 50
Actual:    3,000,000
```

That's a massive estimation problem.

### 3. Check whether the data distribution is skewed

If some values occur vastly more frequently than others, parameter sensitivity becomes much more likely.

---

# How can you fix it?

There isn't one universal solution.

That's important.

The correct solution depends on why the plan is problematic.

---

## Option 1: `OPTION (RECOMPILE)`

You can tell SQL Server to compile the query for that execution:

```sql
SELECT *
FROM Listings
WHERE CityId = @CityId
OPTION (RECOMPILE);
```

Now SQL Server sees the actual parameter value and creates a plan specifically for it.

This can solve parameter-sniffing problems.

But there's a cost:

> **You lose the benefit of reusing the cached plan and pay compilation cost repeatedly.**

So you don't blindly put `RECOMPILE` everywhere.

It's particularly useful when:

- the query isn't executed extremely frequently
- parameter values produce dramatically different optimal plans
- getting the right plan is much more important than compilation overhead

---

# Option 2: Optimize for a particular value

SQL Server has hints such as:

```sql
OPTION (OPTIMIZE FOR ...)
```

You can essentially tell SQL Server to optimize using a particular assumption.

But this is a fairly specialized tool.

You're deliberately saying:

> "Optimize this query as though the parameter behaves like X."

That can be useful in certain workloads, but it can also become fragile if the data distribution changes.

---

# Option 3: Parameter Sensitive Plan optimization

This is particularly interesting because newer SQL Server versions have a feature specifically designed to address parameter-sensitive workloads.

Instead of insisting on one plan for every parameter value, SQL Server can maintain **multiple plans for different parameter ranges** when appropriate.

Conceptually:

```text
Small result set
       ↓
Plan A

Medium result set
       ↓
Plan B

Large result set
       ↓
Plan C
```

Then SQL Server chooses an appropriate plan based on the parameter value.

This is called **Parameter Sensitive Plan optimization (PSP)**.

This is a much more sophisticated solution than simply saying:

> "Disable parameter sniffing."

---

# One important correction to a common misconception

You may hear:

> "Parameter sniffing is bad."

That's not really correct.

**Parameter sniffing is normally beneficial.**

SQL Server looking at the parameter value and using that information to optimize the query can produce a better plan.

The problem is:

> **When one cached plan isn't suitable for the range of parameter values the query receives.**

That's the distinction I'd want you to understand.

---

# How this relates to our previous topics

We now have several things that can produce poor performance:

### Bad index

You don't have an appropriate way to access the data.

### Non-SARGable predicate

You have an index, but your predicate prevents efficient use of it.

### Bad cardinality estimate

SQL Server thinks:

> 100 rows

but reality is:

> 1,000,000 rows.

### Parameter sniffing

SQL Server created a plan based on one parameter value, then reused it for a substantially different parameter value.

These can all interact.

For example:

```text
Parameter value
      ↓
Statistics
      ↓
Cardinality estimate
      ↓
Execution plan
      ↓
Actual workload
```

If the assumptions at the beginning don't match reality, the chosen plan can be very poor.

---

# What I'd say in an interview

If asked:

> **"What is parameter sniffing?"**

I'd say:

> "When SQL Server compiles a parameterized query, it can use the parameter value from that execution to estimate the number of rows and choose an execution plan. That plan is then normally cached and reused. If the data distribution is highly skewed, the plan that is optimal for one parameter value might be very inefficient for another. That's the parameter-sniffing problem. I'd diagnose it by comparing execution plans and estimated versus actual row counts for different parameter values. Depending on the workload, possible solutions include recompilation, query hints, or SQL Server's Parameter Sensitive Plan optimization."

That's a strong answer because you're explaining **why** it happens rather than simply memorizing the definition.

---

### Where we are

We've now covered:

- ✅ Execution plans
- ✅ SARGability
- ✅ Parameter sniffing / parameter-sensitive plans
- **Next: identifying the actual performance bottleneck**

That last part is important because an execution plan alone doesn't tell you everything. We need to distinguish **CPU, I/O, locking/blocking, memory pressure, and other waits** before deciding what to change.