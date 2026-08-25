# 7. Execution plans — how SQL Server actually executes your query

We've been talking about indexes as though SQL Server simply looks at an index and uses it.

That's not really what happens.

When you submit a query, SQL Server has to decide **how to execute it**.

For example:

```sql
SELECT c.Name, AVG(l.Price)
FROM Categories c
JOIN Listings l
    ON c.Id = l.CategoryId
WHERE l.Status = 'ACTIVE'
GROUP BY c.Id, c.Name
ORDER BY AVG(l.Price) DESC;
```

There are many possible ways to execute this.

SQL Server's **query optimizer** examines the query, indexes, statistics, estimated row counts, available operators, and other information and chooses an execution strategy.

The result is an **execution plan**.

---

# 1. What is an execution plan?

An execution plan is essentially SQL Server's description of:

> **"Here's how I intend to obtain and process the data needed to produce this result."**

It might contain operations such as:

```text
Index Seek
   ↓
Nested Loops
   ↓
Key Lookup
   ↓
Sort
   ↓
Aggregate
```

Each of those is an **operator**.

The plan isn't necessarily executed top-to-bottom in the visual order you see.

The arrows represent the **flow of rows between operators**.

That's important.

---

# 2. The three things you should distinguish

When investigating performance, you'll encounter:

### Query

What you wrote:

```sql
SELECT ...
```

### Execution plan

What SQL Server decided to do:

```text
Index Seek
Hash Match
Sort
...
```

### Runtime statistics

What actually happened when it ran:

```text
Estimated rows: 100
Actual rows: 850,000
```

That third one is extremely important.

Because SQL Server can make a bad decision **even when its plan looks reasonable on paper**, if its estimates are wrong.

We'll get to that.

---

# 3. Index Seek

This is usually what people mean when they say:

> "The query is using the index."

Suppose:

```sql
SELECT *
FROM Listings
WHERE Id = 12345;
```

and `Id` is indexed.

SQL Server can navigate directly toward the relevant part of the index.

Conceptually:

```text
Index
  ↓
find 12345
  ↓
matching row
```

That's an **Index Seek**.

A seek doesn't necessarily mean "one row."

You can seek into a range:

```sql
WHERE CreatedAt >= '2026-01-01'
  AND CreatedAt < '2026-02-01'
```

and retrieve thousands of rows.

So:

> **Seek means SQL Server can navigate into the index to locate the relevant portion rather than processing the entire index.**

---

# 4. Index Scan

A scan means SQL Server is reading through a substantial portion of an index.

For example:

```text
Index
 ↓
read many/all index entries
 ↓
filter rows
```

This isn't automatically bad.

That's important.

Suppose your query needs 90% of the rows.

An index seek followed by thousands or millions of lookups could be worse than simply scanning an index.

So:

> **Index Scan ≠ bad.**

You need to ask:

> Why did SQL Server choose it?

and:

> How much data is it scanning?

---

# 5. Table Scan

If the table is a **heap**—meaning it has no clustered index—you can get a:

**Table Scan**

That means SQL Server is reading the heap's pages looking for qualifying rows.

Conceptually:

```text
Table
 ↓
page 1
 ↓
page 2
 ↓
page 3
 ↓
...
```

If the table has a clustered index, you won't normally see a "Table Scan" for that table.

Instead, scanning the entire clustered structure is called a:

**Clustered Index Scan**

And this is another useful distinction:

```text
Heap
    → Table Scan

Clustered table
    → Clustered Index Scan
```

Both can mean "read a lot of the table," but the underlying structures differ.

---

# 6. Why would SQL Server choose a scan when an index exists?

This is a very important question.

Imagine:

```text
Table: 10,000,000 rows
```

Query:

```sql
WHERE Status = 'ACTIVE'
```

And:

```text
9,000,000 rows are ACTIVE
```

Suppose you have:

```sql
CREATE INDEX IX_Listings_Status
ON Listings(Status);
```

SQL Server could do:

```text
Index Seek
   ↓
9,000,000 matching entries
   ↓
millions of lookups
```

That could be awful.

Instead:

```text
Clustered Index Scan
   ↓
read table once
```

could be cheaper.

So SQL Server isn't thinking:

> "An index exists, therefore I must use it."

It's thinking:

> **"Which available strategy is estimated to be cheapest?"**

That distinction is fundamental.

---

# 7. Key Lookup

We've already touched this, but now you can see it in the context of an execution plan.

Suppose:

```sql
SELECT Price, Title
FROM Listings
WHERE CategoryId = 42;
```

You have:

```sql
CREATE INDEX IX_Listings_Category
ON Listings(CategoryId);
```

The plan might look conceptually like:

```text
Index Seek
    ↓
Nested Loops
    ↓
Key Lookup
```

The index finds the rows, but doesn't contain `Price` and `Title`.

So SQL Server goes back to the clustered index for each qualifying row.

If there are:

```text
Actual rows = 5
```

that's probably fine.

If:

```text
Actual rows = 500,000
```

that can be very expensive.

That's why **estimated/actual row counts** matter.

---

# 8. Nested Loops

Now we're getting into operators that aren't themselves indexes.

A **Nested Loops** join essentially does:

> For each row from one input, find matching rows in the other input.

Imagine:

```text
Customers
---------
Alice
Bob
Carol
```

and:

```text
Orders
------
Alice → Order 1
Alice → Order 2
Bob   → Order 3
```

Conceptually:

```text
Customer Alice
    ↓
find Alice's orders

Customer Bob
    ↓
find Bob's orders

Customer Carol
    ↓
find Carol's orders
```

Nested Loops can be **extremely efficient** when the outer input is small and the inner side has an appropriate index.

For example:

```text
10 customers
    ↓
indexed lookup into
10 million orders
```

can be perfectly reasonable.

But:

```text
5 million customers
    ↓
repeatedly search
10 million orders
```

could be disastrous.

Again, context matters.

---

# 9. Hash Match

A **Hash Match** is another join/aggregation strategy.

For a join, conceptually SQL Server can:

```text
Build a hash table from one input
             ↓
Scan the other input
             ↓
Look for matching keys in the hash table
```

For example:

```text
Categories
-----------
1 → Clothes
2 → Electronics
3 → Books
```

SQL Server can build an in-memory hash structure around the category IDs and then process listings looking for matching category IDs.

Hash joins are particularly useful when:

- inputs are relatively large
- suitable indexes aren't available
- SQL Server expects to process many rows

They're not inherently bad either.

Seeing:

```text
Hash Match
```

doesn't mean:

> "My query is broken."

You need to understand why it was chosen.

---

# 10. Merge Join

There is another important join strategy:

**Merge Join**

It works particularly well when both inputs are already ordered by the join key.

For example:

```text
Categories
1
2
3
4

Listings.CategoryId
1
1
2
2
3
4
```

Because both sides are ordered, SQL Server can walk through them together.

Conceptually:

```text
Categories     Listings
   1   ←→       1
   2   ←→       2
   3   ←→       3
   4   ←→       4
```

This can be highly efficient when the required ordering already exists.

Again:

**Nested Loops, Hash Match, and Merge Join are not "good/bad" rankings.**

They're different algorithms suited to different situations.

---

# 11. Sort

If you write:

```sql
ORDER BY Price DESC
```

SQL Server may need a **Sort** operator.

Conceptually:

```text
Rows
 ↓
Sort by Price
 ↓
Return rows
```

Sorting can become expensive when there are many rows.

This is one reason indexes can help with `ORDER BY`.

If the relevant index already provides the required ordering, SQL Server may be able to avoid an explicit sort.

For example:

```text
Index ordered by:

CustomerId
CreatedAt
```

can potentially help:

```sql
WHERE CustomerId = 123
ORDER BY CreatedAt;
```

---

# 12. Aggregate operators

Your interview query had:

```sql
AVG(l.Price)
GROUP BY c.Id, c.Name
```

SQL Server needs to aggregate rows.

You may see operators such as:

### Stream Aggregate

Works well when the input is already appropriately ordered.

Conceptually:

```text
Category 1
Category 1
Category 1
Category 2
Category 2
Category 3
     ↓
Stream Aggregate
     ↓
Category 1 → average
Category 2 → average
Category 3 → average
```

### Hash Aggregate

SQL Server can instead build hash structures to group rows.

Conceptually:

```text
Listing
   ↓
CategoryId
   ↓
hash bucket
   ↓
accumulate Price/count
```

Again, neither is inherently bad.

---

# 13. The really important part: estimated vs actual rows

This is one of the most useful things to understand about execution plans.

Suppose SQL Server estimates:

```text
Estimated rows: 10
```

and the query actually produces:

```text
Actual rows: 500,000
```

That's a **massive estimation error**.

Why does this matter?

Because the optimizer makes decisions based on those estimates.

It might think:

```text
Only 10 rows?
Nested Loops + Key Lookup sounds cheap.
```

But reality is:

```text
500,000 rows
×
500,000 Key Lookups
```

Now the chosen plan is terrible.

This is one of the ways **statistics** become important.

---

# 14. Where do the estimates come from?

SQL Server maintains **statistics** describing data distribution.

For example, it might know that:

```text
Status:
Active → 90%
Sold   → 9%
Expired → 1%
```

or approximately understand the distribution of values in an indexed column.

The optimizer uses this information to estimate:

> "How many rows will this predicate probably return?"

Those estimates influence:

- seek vs scan
- join algorithm
- join order
- memory allocation
- aggregation strategy
- sorting strategy

This is why an index isn't just a data structure.

**Indexes and statistics influence the optimizer's decisions.**

---

# 15. A practical example

Imagine:

```text
Listings: 10 million rows
```

Query:

```sql
SELECT *
FROM Listings
WHERE CategoryId = 123;
```

SQL Server's statistics say:

```text
Estimated rows: 100
```

So it chooses:

```text
Index Seek
   ↓
Nested Loops
   ↓
Key Lookup
```

But reality is:

```text
Actual rows: 800,000
```

Now you've got a huge mismatch.

You might see the query performing poorly even though:

> "There is an index and SQL Server is using an Index Seek!"

The problem isn't necessarily the index itself.

The optimizer made a decision based on an incorrect estimate.

---

# 16. This is why "Index Seek = fast" is wrong

You should get rid of this mental shortcut:

```text
Seek = good
Scan = bad
```

Instead:

```text
Seek
 ↓
How many rows?
 ↓
How many lookups?
 ↓
How expensive is the rest of the plan?
```

and:

```text
Scan
 ↓
How much data?
 ↓
Is the query asking for most of the table?
 ↓
Could the scan actually be cheaper?
```

A scan over 1,000 rows might be perfectly fine.

A seek followed by 5 million lookups might be terrible.

---

# 17. Reading an execution plan without getting overwhelmed

When you first open one, you'll see a lot of boxes and numbers.

Don't try to understand everything simultaneously.

Start with these questions:

### 1. Where is the data coming from?

Look for:

```text
Table Scan
Clustered Index Scan
Index Seek
Index Scan
```

### 2. Are there expensive lookups?

Look for:

```text
Key Lookup
```

### 3. Are there expensive sorts?

Look for:

```text
Sort
```

### 4. How are joins being performed?

Look for:

```text
Nested Loops
Hash Match
Merge Join
```

### 5. Are estimates wildly different from reality?

Compare:

```text
Estimated Number of Rows
Actual Number of Rows
```

### 6. Where is the query spending its work?

The execution plan will provide cost estimates and operator information that help you identify suspicious parts of the plan.

Don't blindly trust the displayed percentage as an absolute measurement, though. It's primarily a **relative optimizer cost estimate**, not a stopwatch.

---

# 18. Actual vs estimated execution plans

SQL Server tooling can show you an **estimated execution plan** without actually running the query.

That's useful for understanding:

> "What does SQL Server think it will do?"

An **actual execution plan** includes runtime information after executing the query.

That's particularly useful for:

> "What actually happened?"

And that's where you can see things such as:

```text
Estimated rows: 100
Actual rows: 75,000
```

which immediately tells you:

> "Something about the optimizer's assumptions is wrong."

---

# 19. What I'd do with your interview query

Take:

```sql
SELECT c.Name, AVG(l.Price)
FROM Categories c
JOIN Listings l
    ON c.Id = l.CategoryId
WHERE l.Status = 'ACTIVE'
GROUP BY c.Id, c.Name
ORDER BY AVG(l.Price) DESC;
```

Before adding an index, I'd inspect the plan.

I'd want to know:

```text
Listings
   ↓
Scan or Seek?
   ↓
How many ACTIVE rows?
   ↓
How are Listings joined to Categories?
   ↓
Nested Loops / Hash Match / Merge Join?
   ↓
How is the GROUP BY performed?
   ↓
Stream Aggregate / Hash Aggregate?
   ↓
Is there a Sort for ORDER BY?
```

Then I'd look at the actual row counts.

Only after that would I decide whether an index such as:

```sql
CREATE INDEX IX_Listings_Status_CategoryId
ON Listings(Status, CategoryId)
INCLUDE (Price);
```

is actually justified.

That is the difference between:

> "I know how to create an index."

and:

> **"I know how to diagnose a query."**

---

# 20. One deeper concept: the optimizer isn't executing your SQL literally

This is worth understanding.

You write:

```sql
SELECT ...
FROM A
JOIN B ...
WHERE ...
GROUP BY ...
ORDER BY ...
```

SQL Server doesn't simply execute those clauses in the textual order you wrote them.

The optimizer transforms the query into an execution strategy.

It may:

- reorder joins
- choose different join algorithms
- choose scans instead of seeks
- push filters earlier
- use indexes
- eliminate unnecessary work
- choose different aggregation strategies

The SQL you write describes **what you want**.

The execution plan describes **how SQL Server intends to get it**.

That distinction is fundamental to understanding relational databases.

---

## Where this leaves us

You now have the major pieces:

```text
Query
 ↓
Optimizer
 ↓
Statistics + indexes + cardinality estimates
 ↓
Execution plan
 ↓
Operators
 ├── Seek / Scan
 ├── Lookup
 ├── Nested Loops
 ├── Hash Match
 ├── Merge Join
 ├── Sort
 └── Aggregate
 ↓
Actual execution
```

### Next: SQL Server statistics and index maintenance

This is the missing piece that explains **why the optimizer sometimes chooses the wrong plan in the first place**, and we'll also cover:

- what statistics actually contain
- how SQL Server creates/updates them
- stale statistics
- `UPDATE STATISTICS`
- `AUTO_UPDATE_STATISTICS`
- index fragmentation
- `REORGANIZE` vs `REBUILD`
- page splits
- fill factor
- when fragmentation actually matters
- why rebuilding every index "because it's fragmented" is not a good maintenance strategy

That will take us to the deeper internals behind the execution plans we've just discussed.