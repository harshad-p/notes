# 3. Composite indexes and column order

This is probably the most important practical indexing topic for you, especially because of the SQL question you were asked.

A **composite index** is an index on multiple columns:

```sql
CREATE INDEX IX_Listings_Status_CategoryId
ON Listings(Status, CategoryId);
```

The crucial thing is that this is **not equivalent** to having two separate indexes:

```sql
CREATE INDEX IX_Listings_Status
ON Listings(Status);

CREATE INDEX IX_Listings_CategoryId
ON Listings(CategoryId);
```

And the order in the composite index matters.

---

## 1. Think of the index as being sorted hierarchically

Suppose we create:

```sql
CREATE INDEX IX_Listings_Status_CategoryId
ON Listings(Status, CategoryId);
```

Conceptually, the index is ordered like:

```text
Status
  ↓
CategoryId
  ↓
row locator
```

Imagine the data:

```text
Status    CategoryId
--------------------
Active       1
Active       1
Active       2
Active       2
Active       5
Active       17
Active       17
Sold         1
Sold         2
Sold         17
```

The database first organizes by `Status`.

Within each `Status` value, it organizes by `CategoryId`.

So:

```text
Active
 ├── CategoryId 1
 ├── CategoryId 2
 ├── CategoryId 5
 └── CategoryId 17

Sold
 ├── CategoryId 1
 ├── CategoryId 2
 └── CategoryId 17
```

That hierarchy is the key to understanding column order.

---

# 2. Why column order matters

Suppose we have:

```sql
CREATE INDEX IX_Listings_Status_CategoryId
ON Listings(Status, CategoryId);
```

Now consider:

```sql
WHERE Status = 'Active'
  AND CategoryId = 17
```

Excellent match.

SQL Server can effectively navigate:

```text
Status = Active
        ↓
CategoryId = 17
        ↓
matching rows
```

Both columns participate naturally in locating the rows.

---

Now:

```sql
WHERE Status = 'Active'
```

Also a good match.

The index starts with `Status`, so SQL Server can find the `Active` portion of the index.

---

But what about:

```sql
WHERE CategoryId = 17
```

This is where things become interesting.

The index starts with `Status`.

The database doesn't have the index globally ordered by `CategoryId`.

Instead, `CategoryId = 17` is scattered across the different `Status` groups:

```text
Active → Category 17
Sold   → Category 17
Expired → Category 17
...
```

So this index isn't nearly as useful for that predicate alone.

This gives us the important concept of the **leading column** or **leftmost prefix**.

---

# 3. The leftmost-prefix idea

For:

```sql
(Status, CategoryId, CreatedAt)
```

the index naturally supports access patterns beginning with:

```text
Status
```

or:

```text
Status + CategoryId
```

or:

```text
Status + CategoryId + CreatedAt
```

But it isn't equivalent to an index beginning with:

```text
CategoryId
```

This is why:

```sql
(Status, CategoryId)
```

and:

```sql
(CategoryId, Status)
```

are two different indexes with different strengths.

---

# 4. Let's apply this to your interview query

Your query was:

```sql
SELECT c.Name, AVG(l.Price) AS avg
FROM Categories c
JOIN Listings l
    ON c.Id = l.CategoryId
WHERE l.Status = 'ACTIVE'
GROUP BY c.Id, c.Name
ORDER BY avg DESC;
```

Suppose we're considering:

```sql
CREATE INDEX IX_Listings_Status_CategoryId
ON Listings(Status, CategoryId);
```

Why might this be useful?

Because the query says:

```text
Filter listings by Status
          ↓
Join listings using CategoryId
          ↓
Group by category
          ↓
Calculate average Price
```

Our index starts with exactly:

```text
Status
CategoryId
```

So it potentially helps SQL Server narrow the listings to the relevant active rows and organize those rows by category.

But notice something:

**Price isn't in the index.**

SQL Server still needs `Price` to calculate:

```sql
AVG(l.Price)
```

This leads directly to **included columns**, which we'll cover next.

---

# 5. Why not just create three separate indexes?

You might wonder:

```sql
CREATE INDEX ... ON Listings(Status);
CREATE INDEX ... ON Listings(CategoryId);
CREATE INDEX ... ON Listings(Price);
```

Why not?

Sometimes that is appropriate.

But separate indexes don't give the same ordering as:

```sql
(Status, CategoryId, Price)
```

A composite index can allow SQL Server to navigate through a particular multi-column ordering.

For example:

```text
(Status, CategoryId)
```

gives you:

```text
Active
   → Category 1
   → Category 2
   → Category 17
Sold
   → Category 1
   → Category 2
   → Category 17
```

Whereas separate indexes are independent structures:

```text
Status index
----------------
Active → rows...
Sold   → rows...


CategoryId index
----------------
1  → rows...
2  → rows...
17 → rows...
```

The optimizer can sometimes combine multiple indexes, but you shouldn't assume it will—or that doing so will be as efficient as having the right composite index.

---

# 6. How do you decide column order?

This is where simplistic rules start causing trouble.

You'll often hear:

> "Put the most selective column first."

That's **not a universal rule**.

A better approach is:

> **Design the index around the actual query patterns and how the predicates are used.**

For example:

```sql
WHERE Status = @status
  AND CategoryId = @categoryId
```

could potentially support:

```text
(Status, CategoryId)
```

or:

```text
(CategoryId, Status)
```

Both columns are equality predicates.

In such a case, both orders may potentially work well for the filtering portion.

But other parts of the query can make one order more useful.

For example:

```sql
WHERE CategoryId = @categoryId
ORDER BY CreatedAt
```

might make:

```text
(CategoryId, CreatedAt)
```

particularly interesting.

The index can organize rows like:

```text
Category 1
   CreatedAt 1
   CreatedAt 2
   CreatedAt 3
Category 2
   CreatedAt 1
   CreatedAt 2
...
```

That can help both filtering and ordering.

---

# 7. Equality vs range predicates

This is another important consideration.

Suppose:

```sql
WHERE Status = 'Active'
  AND CreatedAt >= '2026-01-01'
```

You have:

- equality on `Status`
- range on `CreatedAt`

An index like:

```sql
(Status, CreatedAt)
```

is a natural candidate.

Conceptually:

```text
Status = Active
       ↓
CreatedAt >= Jan 1
       ↓
relevant range
```

Now imagine:

```sql
WHERE CreatedAt >= '2026-01-01'
  AND Status = 'Active'
```

The textual order of predicates in the SQL doesn't matter.

SQL Server doesn't care that you wrote `CreatedAt` first.

What matters is the **query semantics and the index definition**.

---

# 8. Composite indexes can help `ORDER BY`

This is an area people often overlook.

Suppose:

```sql
SELECT *
FROM Orders
WHERE CustomerId = 123
ORDER BY OrderDate DESC;
```

And we have:

```sql
CREATE INDEX IX_Orders_CustomerId_OrderDate
ON Orders(CustomerId, OrderDate);
```

The index is conceptually:

```text
CustomerId
    ↓
OrderDate
    ↓
row
```

So once SQL Server finds:

```text
CustomerId = 123
```

the matching rows are already organized by `OrderDate`.

That can potentially eliminate or reduce the need for a separate sorting operation.

This is why index design isn't just about `WHERE`.

You have to consider:

```text
WHERE
JOIN
ORDER BY
GROUP BY
SELECT
```

together.

---

# 9. But don't blindly put every query column into the key

Suppose:

```sql
SELECT Name, Price, Description
FROM Listings
WHERE CategoryId = 17;
```

You might think:

```sql
CREATE INDEX ...
ON Listings(CategoryId, Name, Price, Description);
```

That's generally not what you want.

Why?

Because now the index itself becomes much larger.

And large indexes:

- consume more storage
- require more I/O
- require more maintenance
- make writes more expensive
- can reduce the number of index pages that fit in memory

Instead, SQL Server provides **included columns**.

---

# 10. Included columns

You can create:

```sql
CREATE INDEX IX_Listings_CategoryId
ON Listings(CategoryId)
INCLUDE (Name, Price, Description);
```

Now:

```text
Index key:
CategoryId

Included data:
Name
Price
Description
```

The important distinction is:

> **Key columns determine the ordering/search structure of the index. Included columns are stored with the index to help satisfy the query, but aren't part of the index's search/order key.**

So if you have:

```sql
WHERE CategoryId = 17
```

and need:

```text
Name
Price
Description
```

the index may be able to provide all of those without going back to the clustered table.

That's called a **covering index** when the index contains everything needed to satisfy the query.

---

# 11. Connecting this to your interview

Your query needed:

```text
Filter:
Status

Join:
CategoryId

Aggregate:
Price
```

One possible index worth investigating could therefore be:

```sql
CREATE INDEX IX_Listings_Status_CategoryId
ON Listings(Status, CategoryId)
INCLUDE (Price);
```

Why?

```text
Key:
Status
CategoryId

Included:
Price
```

The database can potentially:

```text
Find Active listings
        ↓
Navigate by CategoryId
        ↓
Have Price available for AVG()
```

And potentially avoid looking up the clustered rows for `Price`.

**But notice the wording I'm using: "potentially."**

I would not tell an interviewer:

> "This is definitely the correct index."

I'd say:

> "I'd consider a composite index on Status and CategoryId, potentially including Price, but I'd validate it against the actual execution plan, existing indexes, data distribution, and workload."

That's a much better senior-level answer.

---

# 12. One important correction to a common misconception

A composite index does **not** mean:

> "SQL Server can only use it if every column is specified."

That's too simplistic.

For:

```sql
(Status, CategoryId)
```

a query filtering only:

```sql
Status = 'Active'
```

can still use the index.

The problem is primarily when you're trying to use the index through a column that comes **after an unconstrained leading column**.

So:

```text
(Status, CategoryId)

WHERE Status = ...
```

Good candidate.

```text
(Status, CategoryId)

WHERE Status = ...
  AND CategoryId = ...
```

Even better match.

```text
(Status, CategoryId)

WHERE CategoryId = ...
```

Much less useful for seeking directly into the index.

That's the mental model you should retain.

---

# 13. How I'd reason through a composite-index question in an interview

If they give you:

```sql
SELECT ...
FROM Listings
WHERE Status = 'Active'
  AND CategoryId = 42
ORDER BY CreatedAt DESC;
```

Don't immediately start writing:

```sql
CREATE INDEX ...
```

Walk through it:

**First:** What columns filter the data?

```text
Status
CategoryId
```

**Second:** What determines ordering?

```text
CreatedAt
```

**Third:** Which columns are actually needed to produce the result?

Maybe:

```text
Price
Name
```

**Fourth:** Which columns should be index keys versus included columns?

Now you're designing an index based on the query rather than applying a memorized recipe.

---

## One thing I particularly want you to remember

When you see:

```sql
INDEX(A, B, C)
```

read it mentally as:

> **"The index is primarily organized by A, then within A by B, then within A+B by C."**

That immediately explains why:

```text
INDEX(Status, CategoryId)
```

isn't the same thing as:

```text
INDEX(CategoryId, Status)
```

and why the **leading column matters**.

---

### Next: included columns, covering indexes, and key lookups

We touched those here, but I want to give them their own treatment rather than glossing over them. We'll look at **exactly what a key lookup is, why it can become a performance problem, how `INCLUDE` addresses it, and when adding included columns actually makes an index worse rather than better.**