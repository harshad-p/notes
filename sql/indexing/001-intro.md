Absolutely. We'll go through the indexing plan **in order**, with enough depth to understand how SQL Server actually behaves, but without re-explaining basic SQL concepts you already know.

I'll also keep the interview context in mind: the goal isn't to memorize index syntax; it's to be able to look at a query and reason about **whether an index helps, what kind, and why**.

# 1. What a database index actually is

Let's start with the physical idea, because everything else follows from it.

Suppose you have a `Listings` table with 10 million rows:

```text
Id | CategoryId | Status | Price
---------------------------------
1  | 17         | Active | 120
2  | 4          | Active | 80
3  | 17         | Sold   | 200
...
10 million rows
```

Without an index, SQL Server may have to examine a large portion of the table to find:

```sql
SELECT *
FROM Listings
WHERE CategoryId = 17;
```

The simplest mental model is:

```text
Table
┌──────────────────────────────┐
│ row │ row │ row │ row │ ... │
└──────────────────────────────┘
             ↓
       examine rows
```

That's essentially a **table scan**.

An index creates another data structure organized around one or more columns.

For example, an index on `CategoryId` gives SQL Server a structure roughly like:

```text
CategoryId
   ↓
  1 → locations of matching rows
  2 → locations of matching rows
  3 → locations of matching rows
 ...
 17 → locations of matching rows
 ...
```

So instead of examining millions of rows, SQL Server can navigate the index to the relevant entries.

The important point is:

> **An index is a separate data structure maintained by the database that makes certain access patterns faster.**

It is **not** simply a copy of the column.

---

# 2. The structure SQL Server normally uses: B-tree

For ordinary SQL Server indexes, the important structure to understand is a **B+ tree**.

You don't need to implement one, but you should understand why it makes lookups efficient.

Very simplified:

```text
                 [50]
               /      \
          [10, 30]    [70, 90]
          /  |  \      /  |  \
         ... ... ...   ... ... ...
```

The tree has multiple levels.

If you're looking for:

```text
CategoryId = 70
```

SQL Server doesn't start from the first row and inspect everything.

It navigates through the tree:

```text
Root
 ↓
appropriate branch
 ↓
appropriate branch
 ↓
leaf
```

This is why an index can reduce the amount of data SQL Server needs to inspect.

The exact implementation details are more complicated, but the interview-level concept is:

> **The tree keeps values ordered and allows SQL Server to navigate to the relevant range rather than scanning every row.**

---

# 3. Why is this called an index?

Think of a book.

Without an index:

> "Find every occurrence of the word `database`."

You might read every page.

With an index at the back:

```text
database → pages 12, 48, 103, 217...
```

You jump directly to the relevant locations.

A database index is conceptually similar, except SQL Server maintains the structure automatically and uses it to execute queries.

---

# 4. An index doesn't necessarily contain the whole row

This is an important distinction.

Suppose:

```sql
CREATE INDEX IX_Listings_CategoryId
ON Listings(CategoryId);
```

The index is primarily organized around `CategoryId`.

It also contains enough information for SQL Server to locate the corresponding table rows.

So if you execute:

```sql
SELECT *
FROM Listings
WHERE CategoryId = 17;
```

SQL Server might:

```text
Index
  ↓
find CategoryId = 17
  ↓
find corresponding rows
  ↓
go back to Listings table
  ↓
retrieve remaining columns
```

That last step is important.

It's commonly called a **key lookup** when SQL Server uses a nonclustered index to locate rows and then fetches additional columns from the underlying table.

And this leads directly to a major indexing concept we'll cover later:

> **Included columns can sometimes allow SQL Server to satisfy the query entirely from the index.**

We'll get there.

---

# 5. Why not create an index on every column?

This is where indexes become interesting.

Indexes make **reads** potentially faster, but they have costs.

Suppose you have:

```text
Listings
 ├── Id
 ├── CategoryId
 ├── Status
 ├── Price
 ├── CreatedAt
 ├── Description
 └── ...
```

and you create indexes on everything:

```text
Index on Id
Index on CategoryId
Index on Status
Index on Price
Index on CreatedAt
Index on Description
...
```

Now whenever you insert:

```sql
INSERT INTO Listings ...
```

SQL Server doesn't just have to insert the row into the table.

It also has to update all the relevant indexes.

Conceptually:

```text
INSERT
  ↓
Table
  ↓
Index 1 updated
  ↓
Index 2 updated
  ↓
Index 3 updated
  ↓
Index 4 updated
```

So indexes introduce:

- additional storage
- additional maintenance
- additional work on inserts
- additional work on updates
- additional work on deletes

Therefore:

> **An index is a trade-off: faster access for certain queries in exchange for storage and write/maintenance overhead.**

That's one of the most important principles to understand.

---

# 6. An index is useful only if it matches an access pattern

Suppose your application frequently executes:

```sql
SELECT *
FROM Listings
WHERE CategoryId = @categoryId;
```

An index on:

```text
CategoryId
```

is an obvious candidate.

But suppose your application instead almost always does:

```sql
SELECT *
FROM Listings
WHERE Status = 'Active'
  AND CategoryId = @categoryId;
```

Now we have more information.

The query is filtering using **two columns**.

We might consider a composite index:

```sql
CREATE INDEX IX_Listings_Status_CategoryId
ON Listings(Status, CategoryId);
```

But don't take that as a universal recommendation yet.

**Column order matters.**

And that's where indexing gets considerably more interesting.

---

# 7. Single-column vs composite indexes

A **single-column index**:

```sql
CREATE INDEX IX_Listings_CategoryId
ON Listings(CategoryId);
```

is organized around one column.

A **composite index**:

```sql
CREATE INDEX IX_Listings_Status_CategoryId
ON Listings(Status, CategoryId);
```

is organized around multiple columns.

The ordering is:

```text
Status
   ↓
CategoryId
```

not:

```text
CategoryId
   ↓
Status
```

Those two indexes are **not interchangeable**.

For example:

```sql
(Status, CategoryId)
```

and:

```sql
(CategoryId, Status)
```

can behave quite differently depending on the query.

This is one of the things we'll spend significant time on later because it's probably the most important practical indexing concept after understanding indexes themselves.

---

# 8. How do you determine whether a column needs an index?

Don't start with:

> "This column appears in `WHERE`, therefore index it."

That's too simplistic.

Instead, start with the **queries your application actually executes**.

Suppose your application frequently executes:

```sql
SELECT Id, Name, Price
FROM Listings
WHERE CategoryId = @categoryId
ORDER BY Price DESC;
```

Now you have several things happening:

- filtering by `CategoryId`
- ordering by `Price`
- selecting `Id`, `Name`, `Price`

That query gives you considerably more information about what an appropriate index might look like.

You could potentially consider something involving:

```text
CategoryId
Price
```

and potentially included columns:

```text
Id
Name
```

But whether that's actually beneficial depends on:

- how many rows exist
- how selective `CategoryId` is
- how often this query runs
- how frequently the table is modified
- existing indexes
- the query plan
- the size of the columns
- the workload overall

So **index design is query/workload-driven**, not simply column-driven.

---

# 9. Selectivity

This is one concept you should understand well.

Suppose you have 10 million listings.

### `Status`

Maybe:

```text
Active   → 7,000,000
Sold     → 2,500,000
Expired →   500,000
```

Only three possible values.

`Status` has relatively **low selectivity**.

If you ask:

```sql
WHERE Status = 'Active'
```

you're still getting 70% of the table.

An index may not provide a dramatic benefit, depending on the rest of the query and data distribution.

Now consider:

```text
Id
```

There may be 10 million distinct values.

That's extremely selective.

A query:

```sql
WHERE Id = 7348291
```

can identify one row.

An index is extremely useful here.

So:

> **Selectivity describes how effectively a column's value narrows down the candidate rows.**

High selectivity generally makes an index more attractive for filtering.

But—and this is important—**low selectivity does not automatically mean "never index it."**

A low-cardinality column can still be useful in the right index or query pattern.

For example, a filtered index can sometimes be very effective for:

```sql
WHERE Status = 'Active'
```

especially when the filtered subset is small.

We'll cover that later.

---

# 10. Cardinality vs selectivity

You'll hear both terms.

**Cardinality** generally refers to the number of distinct values.

For example:

```text
Status:
Active
Sold
Expired
```

Cardinality ≈ 3 distinct values.

Whereas:

```text
CustomerId:
1
2
3
...
1,000,000
```

has very high cardinality.

**Selectivity** is more about how much a predicate narrows the result set.

They're related, but not exactly synonymous.

For example:

```sql
WHERE CustomerId = 123
```

is highly selective if CustomerId is unique.

Whereas:

```sql
WHERE Status = 'Active'
```

may have poor selectivity if 90% of rows are active.

This distinction becomes useful when reasoning about why the optimizer chooses an index—or chooses **not** to use one.

---

# 11. Why might SQL Server ignore your index?

This is another misconception worth eliminating early.

You create:

```sql
CREATE INDEX IX_Listings_Status
ON Listings(Status);
```

and then run:

```sql
SELECT *
FROM Listings
WHERE Status = 'Active';
```

You might expect:

> "SQL Server has an index, therefore it will use it."

Not necessarily.

Suppose 95% of the table is active.

SQL Server could decide:

> "Using the index would mean finding 9.5 million matching entries and then doing enormous amounts of additional work to retrieve the rows. It may simply be cheaper to scan the table."

The optimizer estimates the cost of different strategies.

So:

> **Having an index doesn't guarantee an index seek.**

This is extremely important when we get to execution plans.

---

# 12. Index Seek vs Index Scan

You'll hear these constantly when looking at SQL Server execution plans.

### Index Seek

SQL Server navigates the index to find the relevant portion.

Conceptually:

```text
Index
  ↓
find matching range
  ↓
read relevant entries
```

Generally desirable for selective lookups.

### Index Scan

SQL Server reads a substantial portion—or all—of the index.

```text
Index
  ↓
read many/all entries
```

A scan isn't automatically bad.

If you genuinely need a large percentage of the data, scanning can be the correct and efficient choice.

This is another interview trap:

> **"Index Scan = bad."**

No.

The question is whether the chosen access method is appropriate for the query and data.

---

# 13. The first practical rule I'd give you

When someone asks:

> **"Which column should I index?"**

don't answer immediately.

Think:

### Step 1
**What query are we trying to make faster?**

### Step 2
**How does that query access the table?**

Look at:

```sql
WHERE
JOIN
ORDER BY
GROUP BY
```

### Step 3
**How selective are those predicates?**

### Step 4
**What indexes already exist?**

### Step 5
**What does the execution plan show?**

### Step 6
**What is the write cost of adding the index?**

That's the reasoning process I want you to develop.

---

## Where we're going next

We've covered the foundation:

- what an index is
- why B-trees matter
- index/table relationship
- read vs write trade-offs
- single vs composite indexes
- selectivity/cardinality
- why SQL Server may or may not use an index
- seek vs scan
- how to reason about candidate columns

**Next: clustered vs nonclustered indexes.**

That's where we'll get into an especially important SQL Server concept: **what the table itself looks like when you have a clustered index, how the clustered key relates to the physical organization of the data, and why a table can have only one clustered index but many nonclustered indexes.**