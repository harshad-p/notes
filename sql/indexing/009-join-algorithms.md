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

# 9. Join algorithms: Nested Loops, Hash Join, and Merge Join

When you write:

```sql
SELECT *
FROM Customers c
JOIN Orders o
    ON c.Id = o.CustomerId;
```

you've told SQL Server **what result you want**.

You haven't told it **how to find the matching rows**.

SQL Server's optimizer can choose among several algorithms. The three major ones you should understand are:

- **Nested Loops**
- **Hash Join**
- **Merge Join**

The important question isn't "which one is best?"

There isn't one.

The optimizer chooses based largely on things like:

- how many rows it expects from each side
- whether useful indexes exist
- whether the inputs are already ordered
- available memory
- the cost of sorting or scanning
- the join predicate

---

# 1. Nested Loops

Let's start with the easiest one.

Imagine:

```text
Customers:

Id
---
10
20
30
```

and:

```text
Orders:

CustomerId
----------
10
10
20
30
30
30
```

A nested-loops join can take one customer at a time and look for matching orders.

It effectively does:

> Take customer 10. Find all orders for customer 10.

Then:

> Take customer 20. Find all orders for customer 20.

Then:

> Take customer 30. Find all orders for customer 30.

The critical question is **how it finds those orders**.

If `Orders.CustomerId` has an index, SQL Server can perform an index seek for each customer.

So if there are only 3 customers:

```text
Customer 10 → index seek → matching orders
Customer 20 → index seek → matching orders
Customer 30 → index seek → matching orders
```

That's potentially extremely efficient.

---

# 2. Why Nested Loops can become terrible

Now imagine:

```text
Customers = 5,000,000 rows
Orders    = 10,000,000 rows
```

and SQL Server decides to use Nested Loops.

It could potentially perform millions of searches against the other input.

Even if each individual search is reasonably fast, doing it millions of times can become extremely expensive.

So Nested Loops tends to be attractive when **the outer input is relatively small** and the inner input can be searched efficiently.

This is a key sentence to remember:

> **Nested Loops is good when you have relatively few outer rows and an efficient way to find matching inner rows.**

---

# 3. What actually happens internally?

Conceptually, the algorithm is:

```text
for each row in outer input
{
    find matching rows in inner input;

    output the matching combinations;
}
```

It doesn't necessarily literally execute a C# `for` loop, of course. That's the algorithmic idea.

If the inner side has an index on the join key, the "find matching rows" operation can be an index seek.

This is why an execution plan might show:

**Nested Loops → Index Seek**

That combination can be perfectly healthy.

---

# 4. Hash Join

Now imagine both tables are large.

Suppose:

```text
Customers = 10 million
Orders    = 100 million
```

There may be no useful index for the join.

Doing a Nested Loops join would potentially require an enormous number of searches.

A Hash Join can solve the problem differently.

Suppose we're joining:

```sql
Customers.Id = Orders.CustomerId
```

SQL Server can take one side and build an in-memory **hash table** keyed by the join column.

For example, conceptually, it processes:

```text
Customer 101
Customer 205
Customer 900
...
```

and calculates a hash value from each customer's ID.

The hash determines which internal bucket the row belongs to.

The important purpose of that hash isn't that SQL Server "sorts" the customers.

It doesn't.

Instead, hashing gives SQL Server a way to quickly determine:

> "Where should I look for a customer with this particular ID?"

Once that hash structure has been built, SQL Server processes the other input—the orders.

For each order, it takes its `CustomerId`, calculates the corresponding hash, and checks the appropriate bucket for matching customer rows.

So rather than repeatedly searching the entire Customers table, SQL Server can quickly narrow down where the matching rows should be.

---

# 5. Why is it called a Hash Join?

Because the join uses a **hash table** to organize one side of the join.

The important thing is that hashing provides fast lookup based on the join key.

For example, suppose:

```text
Customer ID = 12345
```

A hash function transforms that value into a hash value.

SQL Server uses that value to determine the relevant bucket.

It then compares the actual join keys of the candidates in that bucket to make sure they really match.

The hash itself doesn't prove that two values are equal; it narrows down the candidates.

That's an important implementation detail.

Two different values can produce the same hash value—a **hash collision**—so SQL Server still has to compare the actual keys.

---

# 6. Why not always use a Hash Join?

Because building the hash table has a cost.

If you have:

```text
Customers = 5 rows
Orders = 10 million rows
```

building a hash structure just to match five customers may be unnecessary work.

Nested Loops could simply find the five customers' orders through an index.

So Hash Join becomes attractive when you're dealing with **larger inputs** and especially when useful indexes aren't available.

---

# 7. The memory problem with Hash Join

Here's where things get interesting.

The hash table needs memory.

Suppose SQL Server estimates:

```text
100,000 rows
```

and allocates an appropriate amount of memory.

But the actual number is:

```text
20,000,000 rows
```

The hash table may not fit comfortably in memory.

SQL Server can then use temporary storage, commonly involving `tempdb`, to handle the excess data.

That is called a **hash spill**.

A hash spill can significantly hurt performance.

And now you can see the connection with what we just learned:

**Bad cardinality estimate → inappropriate memory grant / execution strategy → potentially expensive execution.**

This is why I wanted you to understand cardinality estimation before diving deeper into join algorithms.

---

# 8. Merge Join

The third major algorithm is Merge Join.

Merge Join has a very different requirement:

> **Both inputs need to be ordered by the join key.**

Suppose we have:

```text
Customers:
10
20
30
40
```

and:

```text
Orders.CustomerId:
10
10
20
30
30
40
```

Because both inputs are ordered, SQL Server can walk through them together.

It doesn't need to repeatedly search the orders table.

It also doesn't need to build a hash table.

It essentially advances through the two ordered inputs and compares the current keys.

If the keys match, it produces the joined rows.

If one key is smaller than the other, it advances the side with the smaller key.

For example:

```text
Customer: 20
Order:    10
```

The order's key is smaller, so SQL Server advances the orders input.

Then:

```text
Customer: 20
Order:    20
```

They match.

It produces the result and continues.

---

# 9. Why would the inputs already be ordered?

Indexes are one reason.

Suppose you have an index that provides rows ordered by:

```text
Orders.CustomerId
```

SQL Server may be able to read that index in the required order.

Similarly, a clustered index can provide ordering based on its key.

If both sides are already suitably ordered, Merge Join can be very efficient.

If they aren't ordered, SQL Server might need to sort them first.

And that can completely change whether Merge Join is attractive.

---

# 10. Comparing the three

Now the differences should make more sense:

| Algorithm | Basic idea | Often attractive when |
|---|---|---|
| **Nested Loops** | Find matches for each outer row | Small outer input + useful index |
| **Hash Join** | Build a hash structure and probe it | Large inputs, especially without useful indexes |
| **Merge Join** | Walk two already-sorted inputs together | Both inputs already ordered by join key |

Don't memorize those as absolute rules.

They're **tendencies**, not laws.

---

# 11. A really important example

Suppose SQL Server estimates:

```text
Customers: 10 rows
Orders: 10 million rows
```

and `Orders.CustomerId` has a good index.

Nested Loops is very attractive:

```text
10 customers
×
indexed search
```

Now suppose SQL Server estimates:

```text
Customers: 5 million rows
Orders: 10 million rows
```

Nested Loops suddenly looks much less attractive.

A Hash Join might be better.

Or if both sides are already ordered appropriately, Merge Join could be attractive.

**The estimated row count can therefore change the join algorithm.**

This is exactly why inaccurate statistics can cause performance problems.

---

# 12. What happens if SQL Server chooses the wrong join?

Suppose the actual situation is:

```text
Outer table = 5 million rows
```

but SQL Server estimates:

```text
Outer table = 20 rows
```

It might choose Nested Loops because 20 rows sounds cheap.

But then execution begins and there are actually millions of rows.

Now SQL Server may have to perform a huge number of inner operations.

The join algorithm itself isn't necessarily defective.

The optimizer made a decision based on a bad estimate.

This is one of the most useful connections between the topics we've covered:

> **Statistics influence cardinality estimates. Cardinality estimates influence the optimizer's choice of join algorithm. The join algorithm can have a huge impact on runtime.**

---

# 13. Join order matters too

This is another detail that's easy to overlook.

Suppose you have:

```sql
A
JOIN B
JOIN C
JOIN D
```

SQL Server doesn't necessarily join them in the order you wrote them.

It can determine that a different order is cheaper.

For example, suppose:

```text
A → 10 million rows
B → 5 million rows
C → 20 rows
D → 2 million rows
```

If joining A with C immediately reduces the result dramatically, SQL Server may prefer that strategy.

The optimizer is essentially searching through possible execution strategies and trying to find a low-cost plan.

With more tables, the number of possible combinations grows rapidly, which is one reason query optimization is a complex problem.

---

# 14. A subtle point: joins don't necessarily mean "load both tables"

When you write:

```sql
SELECT ...
FROM Customers c
JOIN Orders o
    ON c.Id = o.CustomerId
```

SQL Server doesn't necessarily read every row from both tables.

Depending on the predicates and chosen plan, it may be able to eliminate huge portions of the data early.

For example:

```sql
WHERE c.Id = 123
```

could make a Nested Loops strategy extremely efficient:

1. Find customer 123.
2. Use `Orders.CustomerId` index to find that customer's orders.
3. Return the results.

It doesn't need to scan all 10 million orders.

This is another reason indexes on foreign-key/join columns can matter.

---

# 15. One more important concept: physical vs logical join

Your SQL expresses a **logical join**:

```sql
A JOIN B ON A.Id = B.AId
```

That's the relational operation you want.

The physical implementation might be:

- Nested Loops
- Hash Join
- Merge Join

The optimizer is free to choose the physical implementation as long as it produces the correct result.

This is analogous to what we discussed earlier:

> Your SQL describes **what** you want; the execution plan describes **how SQL Server will obtain it**.

---

## The mental model I want you to have

When you see a join in an execution plan, don't immediately think:

> "Hash Join = bad"  
> "Nested Loops = good"

Instead ask:

**What are the estimated row counts?**

**Are the inputs indexed?**

**Are they already ordered?**

**How many times will the inner side be accessed?**

**Does the hash table fit in memory?**

**Are the actual row counts close to the estimates?**

That lets you reason about *why* SQL Server chose the algorithm instead of memorizing rules.

### One particularly useful thing to learn next

The next level after this is **memory grants and `tempdb`**—because that's where Hash Joins, Sorts, spills, large aggregations, and some seemingly mysterious SQL Server performance problems start to make sense.

It's a deeper topic, but it ties directly into what we've just learned, so I recommend doing that next.