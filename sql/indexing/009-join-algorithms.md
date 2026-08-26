# 9. How SQL Server actually performs joins

When you write:

```sql
SELECT *
FROM Customers c
JOIN Orders o
    ON c.Id = o.CustomerId;
```

you have specified **what relationship you want**.

You haven't specified **how SQL Server should find the matching rows**.

The optimizer can choose several algorithms. The three major ones are:

1. **Nested Loops**
2. **Hash Join**
3. **Merge Join**

These are fundamentally different algorithms. Understanding them is much more useful than memorizing when each one is "good."

---

## 1. Nested Loops

Let's start with the simplest.

Suppose SQL Server has:

```text
Customers

Id
---
10
20
30
```

and:

```text
Orders

CustomerId
----------
10
10
20
30
30
30
```

A nested-loops join essentially takes a row from one input and searches the other input for matching rows.

For customer `10`, it finds the orders belonging to `10`.

Then it does the same for customer `20`, then `30`.

The important question is:

> **How does it search the other input?**

If `Orders.CustomerId` has an index, SQL Server can perform an index seek for each customer.

For example, conceptually:

```text
Customer 10 → seek Orders for 10
Customer 20 → seek Orders for 20
Customer 30 → seek Orders for 30
```

If there are only 3 customers, that's very cheap.

Now imagine 5 million customers.

Doing millions of searches into the other table may be much more expensive.

So nested loops is particularly attractive when **one side has relatively few rows and the other side can be searched efficiently**, often through an index.

---

## 2. Why the index matters so much

Consider:

```sql
SELECT *
FROM Customers c
JOIN Orders o
    ON c.Id = o.CustomerId
WHERE c.Id = 10;
```

Suppose only one customer survives the filter.

A nested-loop plan might do:

1. Find customer `10`.
2. Search `Orders` for `CustomerId = 10`.
3. Return the matching orders.

If `Orders.CustomerId` is indexed, that second operation can be very efficient.

But without the index, SQL Server could potentially have to scan the Orders table looking for matching rows.

Now the algorithm becomes much less attractive.

So when you see **Nested Loops + Index Seek**, that's a combination worth understanding rather than simply labeling "good."

---

# 3. Nested loops can also be terrible

Suppose:

```text
Customers → 1,000,000 rows
Orders    → 10,000,000 rows
```

and the optimizer chooses Customers as the outer input.

It could potentially perform a lookup into Orders for every customer.

That's a million searches.

If the optimizer estimated only 100 customers but there are actually a million, you can get exactly the kind of performance problem we discussed with cardinality estimation.

This is why the following combination is important:

> **Nested Loops + unexpectedly large outer input = something worth investigating.**

The algorithm itself isn't bad.

The **number of iterations** is what matters.

---

# 4. Hash Join

Hash joins solve the problem differently.

Instead of repeatedly searching the second table, SQL Server can build an in-memory hash structure from one input.

Suppose we're joining:

```text
Customers
10
20
30
```

with:

```text
Orders
10
10
20
30
30
```

SQL Server can take the customer IDs and use a **hash function** to determine where each ID belongs in a hash table.

For example, conceptually, a hash function might transform:

```text
10 → bucket 4
20 → bucket 7
30 → bucket 2
```

The exact hashing implementation is an internal detail; the important idea is that SQL Server can quickly determine where to look for a particular key.

It then reads the other input and hashes its `CustomerId` using the same scheme.

When it sees:

```text
Order.CustomerId = 20
```

it can calculate the corresponding hash bucket and look there for matching customer rows.

So instead of doing:

> "Search the entire Customers table for every order."

it can build a structure that makes matching much cheaper.

---

# 5. Why is it called a hash join?

Because the matching process uses a **hash table**.

The hash table is an in-memory data structure that maps a key to a location, called a bucket.

The important property is that hashing lets SQL Server locate the relevant candidate rows without scanning the entire build input for every row from the other input.

There is an important caveat: hash collisions can occur. Two different values can produce the same hash bucket. SQL Server therefore still has to check the actual join keys after locating the bucket.

So don't think:

> "The hash tells SQL Server exactly which row it is."

It's more accurate to think:

> **"The hash quickly narrows down where potential matches are, and SQL Server then verifies the actual keys."**

---

# 6. Build input and probe input

A hash join has two conceptual sides.

### Build side

SQL Server reads one input and builds the hash table from it.

### Probe side

SQL Server reads the other input and uses its join key to look for matches in the hash table.

Choosing which side to build is important because the hash table consumes memory.

Generally, SQL Server wants the build side to be the smaller input.

For example, if you have:

```text
Categories = 50 rows
Listings    = 10,000,000 rows
```

building a hash table from Categories makes much more sense than building one from all 10 million listings.

This is one reason cardinality estimates matter.

If SQL Server incorrectly believes Categories has 10 million rows when it actually has 50, it can make a poor decision.

---

# 7. What happens if the hash table doesn't fit in memory?

This is an important deeper detail.

A hash join prefers to operate in memory because that's fast.

But suppose SQL Server needs to build a huge hash table and doesn't have enough memory available.

It may have to **spill to tempdb**.

That means part of the intermediate data gets written to disk rather than remaining entirely in memory.

This can dramatically hurt performance.

In an execution plan, you can often see warnings associated with spills.

So if you see a Hash Match and it is slow, don't immediately conclude:

> "Hash joins are slow."

Investigate things such as:

- how many rows were processed
- whether the estimates were wrong
- how much memory was granted
- whether the hash operation spilled to `tempdb`

That's much more useful.

---

# 8. Merge Join

Merge join works differently again.

It is particularly effective when **both inputs are already sorted by the join key**.

Imagine:

```text
Customers

10
20
30
40
```

and:

```text
Orders

10
10
20
30
30
40
```

Both inputs are ordered by the join key.

SQL Server can walk through them simultaneously.

It doesn't need to repeatedly search one table, and it doesn't need to construct a hash table.

It compares the current keys and advances through the appropriate input.

For example:

- Customer = 10, Order = 10 → match.
- Another Order = 10 → match with the same customer.
- Order advances to 20.
- Customer advances to 20.
- Continue.

This can be extremely efficient because each input can essentially be consumed in order.

---

# 9. Why would the inputs already be sorted?

Indexes are one major reason.

Suppose you have an index whose key is:

```text
Orders(CustomerId)
```

The leaf level of that index is ordered by `CustomerId`.

So SQL Server may already have data available in the required order.

This is another reason indexes aren't merely:

> "Fast lookup structures."

They can also provide **ordering** that other operators can exploit.

---

# 10. Comparing the three

Here's the conceptual difference:

### Nested Loops

> Take a row from one input and find matching rows in the other input.

Best suited to situations such as:

```text
small outer input
+
efficient lookup into inner input
```

### Hash Join

> Build a hash table from one input, then use the hash table to find matches from the other input.

Often useful when:

```text
large inputs
+
no particularly useful ordering/index
```

### Merge Join

> Walk through two inputs that are already ordered by the join key and match them as you advance.

Particularly useful when:

```text
both inputs are already sorted
```

---

# 11. There is no universal "best" join

This is something I want you to really internalize.

Suppose someone asks:

> "Which is faster, Nested Loops or Hash Join?"

The correct answer is:

> **It depends on the size and characteristics of the inputs, available indexes, ordering, and the estimated number of rows.**

For example:

```text
10 customers
10 million orders
```

with an index on:

```text
Orders.CustomerId
```

Nested Loops could be excellent.

But:

```text
5 million customers
10 million orders
```

might make a hash join considerably more attractive.

And if both inputs are already ordered appropriately, a merge join could be excellent.

---

# 12. Cardinality estimation directly affects this decision

Now connect this to what we just learned.

Suppose SQL Server estimates:

```text
Customers → 100 rows
Orders    → 10,000,000 rows
```

Nested Loops might look attractive.

But suppose the actual number of customers is:

```text
2,000,000
```

Suddenly doing repeated lookups into Orders may be very expensive.

The optimizer's choice was based on the estimate.

That's why this:

```text
Estimated rows ≠ Actual rows
```

is often more interesting than simply seeing:

```text
Nested Loops
```

in a plan.

---

# 13. Join order matters too

Suppose you have:

```sql
SELECT ...
FROM A
JOIN B ON ...
JOIN C ON ...
JOIN D ON ...
```

SQL Server doesn't necessarily have to join them in the order you wrote.

It might determine that this is cheaper:

```text
A + C
    ↓
small intermediate result
    ↓
B
    ↓
D
```

rather than:

```text
A + B
    ↓
huge intermediate result
    ↓
C
    ↓
D
```

Why does this matter?

Because intermediate results can become enormous.

If you can reduce the number of rows early, subsequent operations have less work to do.

This is one of the reasons predicates and cardinality estimates matter so much to the optimizer.

---

# 14. A subtle point: SQL Server can change strategies based on runtime information

SQL Server has increasingly sophisticated execution mechanisms that can adapt to runtime conditions in some situations.

For example, newer SQL Server versions have features such as **adaptive joins**, where SQL Server can defer the final choice between certain join strategies until it has observed enough rows at runtime.

The underlying reason for such features is exactly the problem we've been discussing:

> **The optimizer has to make decisions based on estimates, but estimates aren't always perfect.**

You don't need to memorize adaptive-join internals for an interview, but it's a good example of why database engines are much more sophisticated than "pick an index and run the query."

---

# 15. How this appears in your execution plan

Suppose you inspect your interview query and see:

```text
Clustered Index Scan
        ↓
     Hash Match
        ↓
  Hash Aggregate
        ↓
       Sort
```

Don't just think:

> "There are four boxes."

Ask:

1. **Why is it scanning?**
2. **How many rows is it scanning?**
3. **Why did it choose a hash join?**
4. **How many rows enter the hash join?**
5. **How many come out?**
6. **Why is it using a hash aggregate?**
7. **Does the final `ORDER BY` require a sort?**
8. **Are the estimated and actual row counts close?**
9. **Did any operator spill to `tempdb`?**

That's how you read a plan.

---

# 16. One more thing: joins aren't necessarily the expensive part

This is an important debugging lesson.

Suppose a query has:

```text
Hash Match
```

taking 60% of the estimated cost.

You shouldn't automatically conclude:

> "The hash join is the problem."

The hash join might simply be receiving **10 million rows because an earlier filter failed to reduce the dataset**.

The real problem could be upstream.

For example:

```text
Table Scan
    ↓
10,000,000 rows
    ↓
Hash Join
```

Maybe the missing index or non-SARGable predicate caused the initial scan.

Fixing the access path could reduce the rows entering the join from:

```text
10,000,000
```

to:

```text
2,000
```

and the hash join might then become trivial.

Always think about **data flowing through the plan**, not just the operator that looks scary.

---

## The mental model I want you to have

When SQL Server joins two datasets, it has fundamentally different ways of finding matches:

**Nested Loops** repeatedly searches one input for each row of another input.

**Hash Join** builds an in-memory hash structure from one input and uses it to locate potential matches from the other input.

**Merge Join** takes advantage of both inputs already being ordered by the join key and walks through them together.

The optimizer chooses between them based largely on **estimated cardinalities, available indexes, ordering, and costs**.

And if those cardinality estimates are wrong, the optimizer can choose an algorithm that looks sensible according to its assumptions but performs terribly with the actual data.

That's the connection between **statistics → cardinality estimates → execution plans → join algorithms**.

Next I'd move to **SQL Server's memory grants and `tempdb`**, because that explains the deeper mechanics behind things we just mentioned—especially **hash spills, sort spills, why a query can suddenly become extremely slow, and how SQL Server decides how much memory an operation gets**.