## 5. Important SQL Server index types

Now that we understand clustered/nonclustered indexes and `INCLUDE`, let's look at the other index types you're likely to encounter.

One important distinction first:

**"Index type" can mean different things.**

For example:

- clustered vs nonclustered describes **how the index relates to the table**
- unique describes **a constraint/property of the index**
- filtered describes **which rows the index contains**
- columnstore describes **how data is physically organized/compressed**

So these aren't always mutually exclusive categories.

---

# 1. Unique indexes

A unique index guarantees that the indexed key values are unique.

For example:

```sql
CREATE UNIQUE INDEX IX_Users_Email
ON Users(Email);
```

Now the database won't allow:

```text
alice@example.com
alice@example.com
```

as two separate values in that indexed key.

This is useful when uniqueness is a **business/data requirement**, not merely a performance optimization.

For example:

```text
User.Email
Username
ExternalSystemId
Product.SKU
```

might need to be unique.

### Important distinction: unique index vs primary key

A primary key also enforces uniqueness, but it has additional semantic meaning:

> **The primary key identifies the row as the table's primary identifier.**

You can have:

```text
one primary key
many unique constraints/indexes
```

For example:

```text
User
 ├── Id          → PRIMARY KEY
 ├── Email       → UNIQUE
 └── Username    → UNIQUE
```

You could therefore have several unique indexes/constraints on a table.

---

# 2. Primary key index

When you create:

```sql
CREATE TABLE Users
(
    Id INT PRIMARY KEY,
    Name VARCHAR(100)
);
```

SQL Server needs an index to enforce the uniqueness of the primary key.

By default, SQL Server commonly creates a **clustered** primary key unless you specify otherwise.

You can explicitly say:

```sql
PRIMARY KEY CLUSTERED
```

or:

```sql
PRIMARY KEY NONCLUSTERED
```

So remember:

> **Primary key and clustered index are different concepts that are often combined.**

We've already discussed why that distinction matters.

---

# 3. Filtered indexes

This one is particularly interesting for the query you were asked.

A **filtered index contains only rows satisfying a condition**.

For example:

```sql
CREATE INDEX IX_Listings_Active
ON Listings(CategoryId)
WHERE Status = 'Active';
```

Conceptually:

```text
Listings
────────────────────────
Active      ← included
Active      ← included
Sold        ← NOT included
Expired     ← NOT included
Active      ← included
```

Instead of indexing every listing, the index contains only active listings.

---

## Why could this be useful?

Suppose you have:

```text
10 million listings

Active   → 500,000
Sold     → 8 million
Expired  → 1.5 million
```

And your application constantly asks:

```sql
SELECT ...
FROM Listings
WHERE Status = 'Active'
  AND CategoryId = 17;
```

A filtered index containing only active listings could be much smaller than an index covering the entire table.

For example:

```sql
CREATE INDEX IX_Listings_Active_Category
ON Listings(CategoryId)
WHERE Status = 'Active';
```

Now SQL Server has a structure specifically for the active subset.

---

# 4. Why not always use a filtered index?

Because it is specialized.

Suppose you later query:

```sql
WHERE Status = 'Sold'
```

Your filtered active index isn't useful for that.

Also, filtered indexes have restrictions around what predicates are allowed and how the optimizer can match queries to them.

So they are excellent when:

> **A relatively stable subset of rows is queried frequently.**

They're not a replacement for ordinary indexes.

---

# 5. Composite filtered index

You can combine ideas.

For example:

```sql
CREATE INDEX IX_Listings_Active_Category
ON Listings(CategoryId)
INCLUDE (Price)
WHERE Status = 'Active';
```

Now think about what this index contains:

```text
Rows:
    Active listings only

Key:
    CategoryId

Included:
    Price
```

That's a very targeted index for queries such as:

```sql
SELECT CategoryId, AVG(Price)
FROM Listings
WHERE Status = 'Active'
GROUP BY CategoryId;
```

Whether it's actually the best index would still need to be validated with the execution plan and workload.

---

# 6. Columnstore indexes

Now we're moving into a different category.

Most of the indexes we've discussed are **row-oriented**.

Columnstore indexes store data in a **column-oriented format**.

This is particularly useful for:

- analytics
- reporting
- aggregations
- large datasets
- data warehouses

Imagine:

```text
Row storage:

Row 1 → Id, Category, Price, Status
Row 2 → Id, Category, Price, Status
Row 3 → Id, Category, Price, Status
```

Column-oriented storage conceptually looks more like:

```text
Id:
1
2
3
...

Category:
17
5
17
...

Price:
100
200
50
...

Status:
Active
Sold
Active
...
```

If a query only needs:

```sql
SELECT AVG(Price)
FROM Listings
WHERE Status = 'Active';
```

column-oriented storage can be very efficient because the database can work primarily with the relevant columns rather than repeatedly processing entire rows.

Columnstore also uses compression heavily.

---

# 7. Clustered columnstore vs nonclustered columnstore

SQL Server supports both:

```text
Clustered columnstore index
Nonclustered columnstore index
```

A clustered columnstore fundamentally changes how the table's data is stored.

A nonclustered columnstore provides a column-oriented structure alongside the normal rowstore table.

> **Rowstore indexes are generally the normal choice for transactional workloads; columnstore is particularly valuable for large analytical workloads.**

Don't reach for columnstore simply because a query contains `AVG()`.

---

# 8. Full-text indexes

SQL Server also has **full-text search** capabilities.

This is different from a normal index.

Suppose you have:

```text
Product.Description
```

and want searches such as:

```text
"wireless headphones with noise cancellation"
```

A normal B-tree index isn't designed to understand natural-language text searches like that.

Full-text indexing is designed for text-search scenarios.

It's useful for things like:

```text
searching words
phrases
prefixes
linguistic matching
```

---

# 9. Spatial indexes

SQL Server also has **spatial indexes** for spatial/geographic data.

For example:

```text
location
latitude/longitude
geography
geometry
```

If you have:

```text
Find restaurants within 5 km of this point
```

a normal B-tree index isn't the right structure for answering arbitrary spatial relationships efficiently.

Spatial indexing addresses that class of problem.

Again, specialized—but worth recognizing.

---

# 10. XML indexes

SQL Server supports indexes specifically for XML data.

They're useful when you store XML and need to query the XML structure efficiently.

> **Different data access patterns sometimes require different index structures.**

---

# 11. Hash indexes — an important qualification

You may encounter the term **hash index** when reading about SQL Server.

These are associated primarily with **memory-optimized tables / In-Memory OLTP**, rather than being the normal index you create on an ordinary disk-based table.

Hash indexes are designed around equality lookups such as:

```sql
WHERE Id = 123
```

rather than general range queries such as:

```sql
WHERE Price > 100
```

So don't confuse them with the ordinary B-tree indexes we've been discussing.

---

# 12. Let's organize all of this

Rather than memorizing a giant list, categorize them:

### Rowstore indexes

The normal indexes you'll use most often:

```text
Clustered
Nonclustered
Unique
Filtered
```

And some of these can overlap.

For example:

```text
Unique + Nonclustered
Filtered + Nonclustered
```

### Specialized indexes

```text
Columnstore
Full-text
Spatial
XML
Hash (memory-optimized scenarios)
```

Different problems, different structures.

---

# 13. How do you actually create them?

You've already seen the basic syntax.

### Nonclustered

```sql
CREATE INDEX IX_Listings_Status
ON Listings(Status);
```

### Composite

```sql
CREATE INDEX IX_Listings_Status_Category
ON Listings(Status, CategoryId);
```

### Included columns

```sql
CREATE INDEX IX_Listings_Category
ON Listings(CategoryId)
INCLUDE (Price, Name);
```

### Unique

```sql
CREATE UNIQUE INDEX IX_Users_Email
ON Users(Email);
```

### Filtered

```sql
CREATE INDEX IX_Listings_Active
ON Listings(CategoryId)
WHERE Status = 'Active';
```

### Clustered

```sql
CREATE CLUSTERED INDEX IX_Listings_Id
ON Listings(Id);
```

And you can remove an index with:

```sql
DROP INDEX IX_Listings_Status
ON Listings;
```

---

# 14. A practical question: "Should I index Status?"

Let's revisit your interview.

You said:

> "I'd index `Status` because it's used in the WHERE clause."

That's a **reasonable first thought**, but now we can make it more sophisticated.

Imagine:

```text
Listings: 10 million

Active: 9.5 million
Sold: 400,000
Expired: 100,000
```

An index on:

```text
Status
```

isn't necessarily very useful for:

```sql
WHERE Status = 'Active'
```

because you're asking for almost the entire table.

Now imagine:

```text
Active: 100,000
Sold: 9 million
Expired: 900,000
```

The same index becomes potentially much more useful.

Or perhaps the workload is overwhelmingly:

```sql
WHERE Status = 'Active'
  AND CategoryId = ?
```

Now a composite or filtered index could be much more appropriate.

**The column doesn't determine the index. The workload does.**

---

# 15. This is why execution plans matter

We're getting closer to the part I think will be most useful for you.

You can have:

```text
Query
+
Index
```

and still not know whether the index actually helped.

You need to inspect what SQL Server actually did.

For example:

```text
Table Scan
```

might tell you:

> "SQL Server decided scanning the table was cheaper."

Or:

```text
Index Seek
   ↓
Key Lookup × 500,000
```

might tell you:

> "The index is helping find rows, but retrieving the rest of the data is expensive."

Or:

```text
Index Seek
   ↓
Stream Aggregate
```

might be exactly what you wanted.

That's why **execution plans are the next major topic**.

---

## One thing to retain from this section

Don't memorize:

> "There are 10 types of indexes."

Instead remember the purpose:

| Type | Main purpose |
|---|---|
| Clustered | Organize table data around a key |
| Nonclustered | Separate searchable structure |
| Unique | Enforce uniqueness |
| Filtered | Index only a subset of rows |
| Columnstore | Large analytical/aggregation workloads |
| Full-text | Text search |
| Spatial | Geographic/spatial queries |
| XML | XML-specific querying |
| Hash | Equality lookups in memory-optimized scenarios |

And several properties can coexist. For example, an index can be **unique and nonclustered**, or **filtered and nonclustered**.

### Next: determining *which* columns to index in a real query

We'll take actual queries and reason from `WHERE`, `JOIN`, `ORDER BY`, `GROUP BY`, selectivity, cardinality, existing indexes, and write frequency to arrive at an index. That's the part that will directly help with questions like **"How would you improve the performance of this query?"**