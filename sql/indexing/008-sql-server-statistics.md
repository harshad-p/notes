## 8. Statistics, cardinality estimation, and why the optimizer can get things wrong

This is actually one of the most important pieces of understanding SQL Server performance, because it explains **why SQL Server chooses a particular execution plan in the first place**.

### 1. What are statistics?

Suppose you have a `Listings` table with 10 million rows and a `Status` column.

The values might look roughly like:

- `Active`: 8,500,000 rows
- `Sold`: 1,000,000 rows
- `Expired`: 500,000 rows

If SQL Server receives:

```sql
SELECT *
FROM Listings
WHERE Status = 'Active';
```

it needs to decide how to retrieve those rows.

It could use an index on `Status`, or it could scan the table.

But to make that decision, it needs to estimate:

> "How many rows do I expect `Status = 'Active'` to return?"

That's where **statistics** come in.

Statistics contain information about the **distribution of values in a column**, allowing the optimizer to estimate how many rows a predicate is likely to match.

They aren't a copy of your data.

They're a summarized description of the data distribution.

---

# 2. Why does the optimizer need an estimate?

Because different numbers of matching rows can make completely different execution strategies optimal.

Suppose there are only 10 matching rows out of 10 million.

An index could be extremely useful:

```text
Find the 10 relevant rows directly.
```

But suppose there are 9.5 million matching rows.

Using an index might mean:

> Find millions of index entries and potentially perform millions of additional lookups.

At that point, scanning the table may be cheaper.

So the optimizer needs to make a prediction **before executing the query**.

Statistics are one of the major sources of information it uses to make that prediction.

---

# 3. What is actually inside a statistic?

For a normal single-column statistic, one important component is a **histogram**.

A histogram divides the values in a column into ranges and records information about how frequently values occur within those ranges.

Imagine a `Price` column containing:

```text
10
15
20
25
...
100
...
1000
...
50000
```

SQL Server doesn't need to remember every row in the statistic.

Instead, it can build a summarized representation of the distribution.

Conceptually, it can know things such as:

> "Most prices are below €500."

> "There are relatively few listings above €10,000."

> "The value €100 occurs very frequently."

That information helps estimate predicates such as:

```sql
WHERE Price > 10000
```

or:

```sql
WHERE Price = 100
```

---

# 4. Histogram isn't just "a list of values"

This distinction matters.

A histogram isn't simply:

```text
10 → 100 rows
20 → 200 rows
30 → 50 rows
...
```

for every distinct value.

SQL Server has a limited number of histogram steps, so it summarizes the distribution into ranges.

That's important because a table can contain millions of distinct values.

The statistic therefore represents the distribution **approximately**.

That immediately gives you an important insight:

> **Cardinality estimation is inherently an estimation process.**

SQL Server doesn't necessarily know the exact number of rows that will come out of every operation before executing it.

---

# 5. What is cardinality?

**Cardinality** simply means the number of rows in a set/result.

For example:

```sql
SELECT *
FROM Listings
WHERE Status = 'Active';
```

If 8,500,000 rows satisfy the predicate, the actual cardinality is:

**8,500,000 rows.**

When SQL Server hasn't executed the query yet, it has an **estimated cardinality**.

For example:

```text
Estimated: 8,300,000
Actual:    8,500,000
```

That's a reasonably good estimate.

But you might see:

```text
Estimated: 100
Actual:    2,000,000
```

That's a catastrophic estimation error.

And now we can understand something from the previous section much more deeply.

---

# 6. How an estimation error changes the execution plan

Imagine this query:

```sql
SELECT *
FROM Listings
WHERE CategoryId = 42;
```

Suppose SQL Server estimates:

```text
100 rows
```

It might think:

> "Only 100 rows. I'll use an index and retrieve those rows."

That could result in an index seek followed by lookups.

But suppose the actual data contains:

```text
2,000,000 rows with CategoryId = 42
```

The chosen strategy may now be terrible.

The optimizer didn't necessarily make a nonsensical decision.

It made a decision based on **bad information**.

That's a crucial distinction.

---

# 7. Why can statistics become inaccurate?

Because your database changes.

Imagine SQL Server creates statistics when your table looks like:

```text
Active: 100,000
Sold:   9,900,000
```

The optimizer learns that `Active` is relatively rare.

Then over the next year your application changes the data:

```text
Active: 8,000,000
Sold:   2,000,000
```

If the statistics haven't been updated sufficiently, the optimizer may still have an outdated understanding of the distribution.

The database has changed.

The statistic hasn't caught up.

Now SQL Server can make poor cardinality estimates.

---

# 8. Does SQL Server automatically update statistics?

Yes.

SQL Server has **automatic statistics management**.

When `AUTO_UPDATE_STATISTICS` is enabled—which is the normal/default configuration for modern SQL Server databases—SQL Server can determine that statistics have become sufficiently stale and update them.

But "automatic" doesn't mean:

> "Statistics are continuously and perfectly synchronized with the table."

There is overhead associated with calculating statistics, so SQL Server uses thresholds to decide when an update is warranted.

That means statistics can be temporarily stale.

---

# 9. What does `UPDATE STATISTICS` do?

You can explicitly tell SQL Server to update statistics:

```sql
UPDATE STATISTICS Listings;
```

Or for a particular statistic:

```sql
UPDATE STATISTICS Listings IX_Listings_Status;
```

You can also use:

```sql
UPDATE STATISTICS Listings WITH FULLSCAN;
```

`FULLSCAN` tells SQL Server to examine all rows when building the statistic rather than using a sample.

That's potentially much more expensive on a huge table.

So you shouldn't casually do:

```sql
UPDATE STATISTICS ... WITH FULLSCAN
```

on everything.

You're trading more accurate statistics for more work.

---

# 10. Why would sampling be sufficient?

Imagine a table with:

```text
1 billion rows
```

Reading every row just to build statistics could itself be expensive.

If SQL Server can inspect a representative sample and determine:

> "The distribution appears to look like this."

it can create useful statistics at a fraction of the cost.

For many datasets, that's perfectly adequate.

But sampling can have problems when the data distribution is unusual—for example, if important values are extremely rare.

That's one reason `FULLSCAN` exists.

---

# 11. Statistics aren't only created because you explicitly create an index

This is an important detail.

SQL Server can create **auto-created statistics** for columns used by queries when it determines that statistics would help optimization.

So don't think:

> "Statistics = index metadata."

They're related, but they're not the same thing.

You can have statistics associated with an index, but SQL Server can also maintain statistics that aren't simply the index itself.

---

# 12. Composite indexes and statistics

Now let's connect this to something we already covered.

Suppose you create:

```sql
CREATE INDEX IX_Listings_Status_Category
ON Listings(Status, CategoryId);
```

SQL Server has information associated with the leading index key and can use that information during optimization.

The **order of columns matters** because the optimizer can reason much more directly about the leading key.

For example, the index is ordered conceptually by:

```text
Status first
CategoryId second
```

So it can efficiently reason about:

```sql
WHERE Status = 'Active'
```

and:

```sql
WHERE Status = 'Active'
  AND CategoryId = 42
```

But a query filtering only on:

```sql
WHERE CategoryId = 42
```

doesn't get the same benefit from that index structure, because `CategoryId` isn't the leading key.

This is one reason the leftmost/leading-column concept matters.

---

# 13. A deeper problem: correlated columns

Here's where cardinality estimation becomes genuinely interesting.

Suppose you have:

```text
Status
CategoryId
```

and your data has a relationship between them.

For example, perhaps:

```text
CategoryId = 10
```

contains almost exclusively:

```text
Status = 'Active'
```

while:

```text
CategoryId = 20
```

contains mostly:

```text
Status = 'Sold'
```

Now consider:

```sql
WHERE Status = 'Active'
AND CategoryId = 10
```

The optimizer needs to estimate how many rows satisfy **both** conditions.

A simplistic approach might treat the columns as independent.

But they're not independent.

That's called **correlation**.

This is one of the reasons multi-column statistics and the optimizer's cardinality-estimation model become important for more complex workloads.

You don't need to memorize the mathematical model, but understand the underlying problem:

> **Knowing the distribution of each column independently doesn't necessarily tell you the distribution of combinations of columns.**

That's a real source of estimation error.

---

# 14. Why this matters for your indexing decisions

Suppose you see:

```sql
WHERE Status = 'Active'
AND CategoryId = 10
```

and someone says:

> "Status is low-selectivity, so don't index it."

That's too simplistic.

The combination:

```text
Status + CategoryId
```

might be highly selective.

And SQL Server's statistics need to understand enough about that combination to estimate the result correctly.

This is why **query optimization is more complicated than simply counting distinct values in individual columns.**

---

# 15. Now let's move to fragmentation

This is a completely different issue from statistics.

People often mix these together:

> "The index is fragmented, so the optimizer has bad statistics."

Those are not the same thing.

**Statistics describe data distribution.**

**Fragmentation describes how index pages are physically arranged.**

You need to keep those concepts separate.

---

# 16. What is an index page?

SQL Server stores data in fixed-size pages.

For normal data pages, we're talking about **8 KB pages**.

An index is therefore not one giant continuous block of memory.

It's a collection of pages containing index records.

For a B-tree index, those pages are connected into a structure that allows SQL Server to navigate from the root toward the relevant leaf pages.

---

# 17. What is fragmentation?

Imagine an index logically contains keys in this order:

```text
10, 20, 30, 40, 50, 60, 70
```

The logical order of those keys matters for the index.

But the physical pages containing those values may not be stored in perfectly sequential physical order on disk.

As data is inserted, updated, and deleted, pages can become less optimally arranged.

That's **fragmentation**.

There are actually different notions involved, but the basic idea is:

> **The logical ordering of index pages and their physical arrangement can become less efficient as the index changes.**

---

# 18. Page splits

This is one of the mechanisms that causes fragmentation.

Suppose an index page is nearly full.

Now you insert a value that belongs in that page.

There isn't enough room.

SQL Server may have to split the page.

Conceptually:

```text
Before:

Page A
[10 20 30 40 50]

Insert 25

After:

Page A
[10 20 25]

Page B
[30 40 50]
```

The actual behavior is more nuanced, but that's the basic mechanism.

Page splitting means SQL Server has to redistribute rows between pages and maintain the index structure.

Frequent page splits can create overhead and contribute to fragmentation.

---

# 19. Why would SQL Server leave empty space?

This is where **fill factor** comes in.

Suppose you create an index with a fill factor of 80%.

SQL Server doesn't initially try to pack every page completely full.

It leaves some free space.

Why?

Because if you're going to insert values into the index later, that free space can accommodate some of those inserts without immediately requiring page splits.

There's a tradeoff.

### Higher fill factor

More data per page initially.

Advantages:

- fewer pages
- potentially better storage efficiency
- fewer pages to read

Disadvantage:

- less room for future inserts
- potentially more page splits

### Lower fill factor

More free space.

Advantages:

- more room for inserts
- potentially fewer page splits

Disadvantages:

- more pages
- more storage
- potentially more I/O

So:

> **Fill factor is a write-pattern optimization, not a magical "set it low" performance setting.**

---

# 20. Should you rebuild fragmented indexes?

This is where old SQL Server advice can get dangerous.

You might hear:

> "If an index is fragmented, rebuild it."

That's not a sufficiently good rule.

You need to consider:

- how fragmented it actually is
- how large the index is
- how frequently it's accessed
- what type of workload you have
- whether fragmentation is actually affecting performance
- the cost of rebuilding it

A tiny index being 40% fragmented may not matter at all.

Rebuilding it could cost more than simply leaving it alone.

---

# 21. `REORGANIZE` vs `REBUILD`

SQL Server provides two common maintenance operations.

### REORGANIZE

```sql
ALTER INDEX IX_Listings_Status
ON Listings
REORGANIZE;
```

This is a more incremental operation that reorganizes the existing index structure.

It generally requires fewer resources than a full rebuild, although it can take longer to accomplish the same degree of cleanup.

### REBUILD

```sql
ALTER INDEX IX_Listings_Status
ON Listings
REBUILD;
```

This essentially reconstructs the index.

It can produce a much more thoroughly reorganized index, but it's a heavier operation.

A rebuild can also update the associated index statistics as part of the operation.

---

# 22. Why statistics and fragmentation are often confused

Consider two separate problems.

### Problem A

Your statistics say:

```text
Estimated rows: 10
Actual rows: 2,000,000
```

This suggests a **cardinality-estimation/statistics problem**.

### Problem B

Your index has become physically fragmented.

That is an **index storage/layout problem**.

They can occur at the same time, but fixing one doesn't automatically fix the other conceptually.

This distinction is important.

---

# 23. And here's an important connection

When you rebuild an index, SQL Server can update the index's statistics.

So rebuilding can sometimes appear to fix a query that was actually suffering from stale statistics.

That doesn't mean:

> "Fragmentation caused the bad query."

It could mean:

> "The rebuild also refreshed information that the optimizer needed."

This is one reason blindly rebuilding indexes can mask the actual problem.

---

# 24. What should you actually do in production?

You don't normally wake up every morning and say:

> "Let's rebuild all indexes."

A sensible maintenance strategy considers:

- index size
- fragmentation
- workload
- statistics freshness
- query performance
- maintenance windows
- available resources

And importantly:

> **Measure whether maintenance is actually solving a problem.**

---

# 25. One more deep distinction: fragmentation doesn't always hurt

This is worth emphasizing because it gets oversimplified constantly.

On modern storage systems, especially SSD-backed systems, the traditional relationship between physical fragmentation and query performance isn't as straightforward as older SQL Server guidance might suggest.

For a query that performs a highly selective seek and reads a handful of pages, moderate fragmentation might make essentially no meaningful difference.

For large range scans, page ordering and density can matter more.

So don't memorize:

```text
30% fragmentation = rebuild
```

as some universal law.

There isn't one.

---

# 26. Putting everything together

Now you can see the chain much more clearly.

When SQL Server receives your query, it doesn't just ask:

> "Is there an index?"

It has to reason about:

**The data**

How many rows exist?

**Statistics**

What does the data distribution look like?

**Cardinality**

How many rows do I expect this predicate to return?

**Indexes**

What access paths are available?

**Cost**

Which execution strategy appears cheapest?

**Execution plan**

Which combination of scans, seeks, joins, sorts, aggregates, etc. should I use?

Then, during execution:

**Reality**

Did the number of rows actually match my estimates?

If not, the optimizer may have chosen a poor strategy.

That is the underlying reason why SQL Server performance tuning isn't simply:

> "Find slow query → add index."

It's an interaction between **data distribution, statistics, indexes, cardinality estimation, and the execution plan**.

---

### One thing I'd like you to understand before we move on

If you see this in an actual execution plan:

```text
Estimated rows: 50
Actual rows: 2,500,000
```

your first thought should **not** be:

> "I need a bigger/faster server."

It should be:

> **"Why did SQL Server believe there would only be 50 rows?"**

That question leads you toward the real cause—statistics, data distribution, predicates, correlations, parameterization, or other cardinality-estimation issues.

That's the level of reasoning I want us to build.

Next, we can go deeper into how SQL Server's optimizer actually chooses between Nested Loops, Hash Join, and Merge Join, including what each algorithm physically does, when each wins, and why the optimizer's estimated row count can completely change the chosen algorithm.