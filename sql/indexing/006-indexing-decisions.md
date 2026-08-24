## 6. How to decide what to index

This is the part that matters most in practice.

Knowing how to write:

```sql
CREATE INDEX ...
```

is easy. The harder senior-level question is:

> **"Why did you choose those columns, in that order, and why do you think this index will help?"**

You should be able to answer that from the query and the workload.

---

# 1. Start with the query, not the table

Suppose you get:

```sql
SELECT Id, Price, Title
FROM Listings
WHERE CategoryId = 42
  AND Status = 'ACTIVE'
ORDER BY Price DESC;
```

Don't immediately think:

> "I need an index on CategoryId."

Break the query down.

### Filtering

```sql
WHERE CategoryId = 42
  AND Status = 'ACTIVE'
```

Potential index keys:

```text
CategoryId
Status
```

### Ordering

```sql
ORDER BY Price DESC
```

Potentially relevant:

```text
Price
```

### Data being returned

```text
Id
Price
Title
```

Potentially relevant as included columns.

Now we're in a position to reason about an index.

---

# 2. Equality predicates are usually good index candidates

Consider:

```sql
WHERE CategoryId = 42
```

An index on `CategoryId` is an obvious candidate.

Similarly:

```sql
WHERE UserId = @userId
```

or:

```sql
WHERE Status = 'ACTIVE'
```

could potentially benefit from an index.

But there's an important caveat:

**Being used in a `WHERE` clause does not automatically mean a column should be indexed.**

We'll get to why.

---

# 3. Selectivity matters

Suppose:

```text
1,000,000 rows
```

and:

```text
Status:
ACTIVE → 950,000
INACTIVE → 50,000
```

An index on `Status` isn't particularly useful for:

```sql
WHERE Status = 'ACTIVE'
```

because you're asking for 95% of the table.

Compare:

```text
CategoryId:
1 → 100 rows
2 → 120 rows
3 → 90 rows
...
```

An index on `CategoryId` could be much more useful because it can narrow the search dramatically.

This is **selectivity**.

A highly selective predicate eliminates a large proportion of the rows.

For example:

```text
Id = 12345
```

is extremely selective if `Id` is unique.

```text
Status = 'ACTIVE'
```

might be poorly selective.

---

# 4. But selectivity isn't the whole story

Don't turn this into:

> "Always index the most selective column."

That's another oversimplification.

Suppose:

```sql
WHERE Status = 'ACTIVE'
  AND CategoryId = 42
```

Individually:

```text
Status       → low selectivity
CategoryId   → medium selectivity
```

But together:

```text
Status + CategoryId
```

might identify very few rows.

So a composite index can exploit the combination.

This is one reason we discussed column order earlier.

---

# 5. Look at JOIN conditions too

Consider:

```sql
SELECT ...
FROM Categories c
JOIN Listings l
    ON c.Id = l.CategoryId
```

`Listings.CategoryId` is a potential index candidate because it participates in the join.

A very common mistake is to look only at:

```text
WHERE
```

and completely ignore:

```text
JOIN
```

For large tables, indexes on join columns can be extremely important.

For example:

```sql
CREATE INDEX IX_Listings_CategoryId
ON Listings(CategoryId);
```

could make finding listings belonging to a particular category much more efficient.

---

# 6. Foreign keys don't automatically mean "there is an index"

This is worth remembering.

If you have:

```sql
CategoryId INT REFERENCES Categories(Id)
```

that establishes a **foreign-key constraint**.

It does not necessarily mean SQL Server automatically created an index on `Listings.CategoryId`.

You may need to create one yourself if the workload benefits from it.

This is a useful interview point.

---

# 7. Look at ORDER BY

Suppose:

```sql
SELECT *
FROM Orders
WHERE CustomerId = 123
ORDER BY CreatedAt DESC;
```

A potentially useful index is:

```sql
CREATE INDEX IX_Orders_Customer_Created
ON Orders(CustomerId, CreatedAt DESC);
```

Why?

The database can first locate:

```text
CustomerId = 123
```

and the matching entries are organized by:

```text
CreatedAt
```

which can potentially reduce or eliminate an additional sort.

This is why indexes aren't just about filtering.

---

# 8. Look at GROUP BY

Suppose:

```sql
SELECT CategoryId, AVG(Price)
FROM Listings
WHERE Status = 'ACTIVE'
GROUP BY CategoryId;
```

Potentially useful information for the index includes:

```text
Status
CategoryId
Price
```

One possible index:

```sql
CREATE INDEX IX_Listings_Status_Category
ON Listings(Status, CategoryId)
INCLUDE (Price);
```

Again, this is a **candidate**, not automatically the correct answer.

The actual execution plan and data distribution determine whether it helps.

---

# 9. Don't forget the SELECT list

Suppose:

```sql
SELECT Name, Price, CreatedAt
FROM Listings
WHERE CategoryId = 42;
```

An index:

```sql
CREATE INDEX IX_Listings_Category
ON Listings(CategoryId);
```

can find the matching rows.

But then SQL Server might need key lookups to retrieve:

```text
Name
Price
CreatedAt
```

So you might consider:

```sql
CREATE INDEX IX_Listings_Category
ON Listings(CategoryId)
INCLUDE (Name, Price, CreatedAt);
```

Now the index may cover the query.

This is exactly why we covered included columns before moving on.

---

# 10. Think about the workload, not one query

This is a very important senior-level distinction.

Suppose you have:

```text
1 query that runs 10 times/day
```

and:

```text
another query that runs 50,000 times/minute
```

You shouldn't necessarily optimize the first one just because it's slow.

You need to consider:

- how often the query runs
- how expensive it is
- how much data it processes
- how important it is
- whether it is latency-sensitive
- how frequently the underlying table is modified

An index that saves 100 ms from a query executed once per day might not justify itself.

An index that saves 20 ms from a query executed 50,000 times per minute absolutely might.

---

# 11. Reads vs writes

This is the fundamental tradeoff.

Indexes generally make **reads faster**.

But indexes make **writes more expensive**.

Imagine:

```text
Listings
```

with:

```text
1 clustered index
5 nonclustered indexes
```

When you insert a row:

```text
INSERT
  ↓
clustered index must be updated
  ↓
nonclustered index #1
  ↓
nonclustered index #2
  ↓
nonclustered index #3
  ↓
...
```

Similarly, an update can require maintaining indexes whose keys or included values are affected.

So you don't want:

> "An index for every column."

You want:

> **"Indexes that provide enough read benefit to justify their storage and write-maintenance cost."**

---

# 12. What about a column with only two possible values?

For example:

```text
IsActive BIT
```

You might think:

> "Only two values—surely an index is useless."

Not necessarily.

Suppose:

```text
10 million rows

IsActive = 1 → 20,000
IsActive = 0 → 9,980,000
```

An index on `IsActive = 1` could be useful.

A **filtered index** might be particularly attractive:

```sql
CREATE INDEX IX_Listings_Active
ON Listings(CategoryId)
WHERE IsActive = 1;
```

So cardinality and data distribution matter more than simply asking:

> "How many distinct values does this column have?"

---

# 13. What about functions on indexed columns?

This is a common performance trap.

Suppose you have:

```sql
WHERE LOWER(Name) = 'harshad'
```

and an index on:

```text
Name
```

The database may not be able to use that ordinary index as efficiently as it could for:

```sql
WHERE Name = 'Harshad'
```

because you're transforming the indexed value.

Similarly:

```sql
WHERE YEAR(CreatedAt) = 2026
```

is often less index-friendly than a range:

```sql
WHERE CreatedAt >= '2026-01-01'
  AND CreatedAt <  '2027-01-01'
```

The second form gives SQL Server a range it can seek into.

This is often called making the predicate **SARGable**.

You don't need to memorize the acronym, but you absolutely should understand the concept.

> **Write predicates in a way that allows the database to use the index efficiently.**

---

# 14. Another classic example: leading wildcard

Suppose:

```sql
WHERE Name LIKE '%phone%'
```

A normal B-tree index on `Name` generally can't efficiently seek to "phone" because the database doesn't know where the matching strings begin.

Compare:

```sql
WHERE Name LIKE 'phone%'
```

Now the beginning of the value is known, so an ordinary index can potentially be useful.

This is one reason specialized full-text search exists for more complex text-search requirements.

---

# 15. Don't blindly index every JOIN column either

Consider:

```sql
JOIN SmallLookupTable
```

with only:

```text
50 rows
```

An index may provide almost no meaningful benefit.

SQL Server could simply scan 50 rows extremely cheaply.

This is where **table size and cardinality** matter.

An index isn't automatically useful just because a column appears in:

```sql
JOIN
```

---

# 16. The execution plan is the final judge

This is the most important thing I want you to take away.

You don't ultimately determine:

> "This column needs an index."

by looking at the SQL alone.

You make a hypothesis.

Then:

```text
Query
 ↓
Execution plan
 ↓
Observe actual behavior
 ↓
Create/modify index
 ↓
Run again
 ↓
Compare
```

For example, you might think:

> "The query is filtering on Status, so I'll index Status."

But SQL Server may tell you:

```text
Table Scan
```

because:

```text
Status = 'ACTIVE'
```

matches 90% of the table.

Or you might create:

```text
(Status, CategoryId)
```

and discover:

```text
Index Seek
    ↓
Key Lookup × 400,000
```

Then you investigate a covering index.

That's much more mature than simply saying:

> "I indexed the WHERE column."

---

# 17. How I'd answer "How would you improve this query?"

This is very close to what your interviewer asked.

Don't immediately answer:

> "I'd add an index."

Instead:

> "I'd first look at the execution plan and actual execution statistics to determine where the query is spending time. I'd check whether we're scanning a large table, whether the join is causing expensive lookups, whether sorting or aggregation is expensive, and whether existing indexes are being used effectively. Based on that I'd consider an appropriate index, potentially a composite or covering index, and then compare the execution plan and runtime before and after."

**That's a senior answer.**

Then, if they ask:

> "Okay, what index would you try?"

Now you can say:

> "For this particular query, I'd consider an index beginning with `Status` and `CategoryId`, potentially including `Price`, because those columns are used for filtering, joining/grouping, and aggregation. But I'd validate that against the actual data distribution and execution plan rather than assuming it's optimal."

That's much stronger than your original:

> "I'd index Status because it's in the WHERE clause."

Your original answer wasn't wrong. **It was just incomplete.**

---

# 18. The indexing decision framework

When you're handed a query, mentally go through:

### Step 1 — What filters rows?

```text
WHERE
```

### Step 2 — What connects tables?

```text
JOIN
```

### Step 3 — What determines ordering?

```text
ORDER BY
```

### Step 4 — What determines grouping?

```text
GROUP BY
```

### Step 5 — What data does the query actually need?

```text
SELECT
```

### Step 6 — How selective are the predicates?

```text
How many rows survive?
```

### Step 7 — How often does this query run?

```text
10/day?
100/sec?
```

### Step 8 — How often is the table written to?

```text
Read-heavy?
Write-heavy?
```

### Step 9 — What does the execution plan actually show?

```text
Scan?
Seek?
Lookup?
Sort?
Hash?
```

### Step 10 — Test the change

```text
Before
   ↓
change
   ↓
After
   ↓
compare
```

That is the thought process I want you to develop.

---

## Where we are in the indexing topic

We've now covered:

1. **Clustered vs nonclustered**
2. **Composite indexes and column order**
3. **Included columns / covering indexes / key lookups**
4. **Unique, filtered, columnstore and other index types**
5. **How to decide what to index**

There are two remaining pieces I'd consider **important enough for your interview prep**:

### Next: execution plans

We'll go through **Index Seek vs Index Scan vs Table Scan, Key Lookup, Sort, Hash Match, Nested Loops, estimated vs actual rows, and how to read a plan without getting lost in it.**

After that, we'll do **index maintenance and statistics** briefly—fragmentation, updating statistics, rebuilding vs reorganizing, and why SQL Server sometimes makes a bad indexing decision because its statistics are stale.

Then I'd consider the indexing topic complete.