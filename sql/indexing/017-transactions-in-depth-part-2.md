# Transactions in Depth — Part 2: Concurrency & Isolation

The next important piece is **what happens when multiple requests try to read or modify the same data at the same time**. 

## 1. The classic problem: two requests update the same row

Suppose we have:

```text
Listings
Id    Price
123   500000
```

Two users simultaneously try to change the price:

- Request A wants `520000`
- Request B wants `530000`

A naïve sequence could be:

```text
A reads 500000
B reads 500000

A writes 520000
B writes 530000
```

The final value is `530000`.

But neither request knew that the other had modified the row. **A's update has effectively been lost.**

This is called a **lost update**.

There are two major approaches to solving this:

1. **Pessimistic concurrency**
2. **Optimistic concurrency**

---

# 2. Pessimistic concurrency

The idea is:

> "I'm going to assume someone else might modify this data, so I'll lock it while I'm working with it."

The database prevents another transaction from modifying the relevant row until the first transaction is finished.

Conceptually:

```text
Transaction A
    ↓
Lock row 123
    ↓
Read / modify row
    ↓
Commit
    ↓
Unlock row

Transaction B
    ↓
Tries to modify row 123
    ↓
Waits
    ↓
Gets the row after A commits
```

This gives you strong protection, but there is a cost: **concurrency decreases because other transactions may have to wait.**

And if transactions hold locks for too long, you can get significant blocking.

This is why we discussed earlier that you don't want:

```text
BEGIN TRANSACTION

update listing

call external payment API
wait 4 seconds

send email
wait 2 seconds

COMMIT
```

The database may be holding locks for all those seconds.

---

# 3. Optimistic concurrency

Optimistic concurrency takes the opposite approach:

> "I'll assume conflicts are relatively rare. Let everyone work, but detect whether someone changed the data before I commit my change."

This is extremely common in applications.

Add a version column:

```text
Listings
Id    Price     Version
123   500000    7
```

Request A reads:

```text
Price = 500000
Version = 7
```

Request B also reads:

```text
Price = 500000
Version = 7
```

A attempts:

```sql
UPDATE Listings
SET Price = 520000,
    Version = 8
WHERE Id = 123
  AND Version = 7;
```

One row is updated.

Now the database contains:

```text
Price = 520000
Version = 8
```

Then B attempts:

```sql
UPDATE Listings
SET Price = 530000,
    Version = 8
WHERE Id = 123
  AND Version = 7;
```

**Zero rows are updated.**

Why?

Because the row is now:

```text
Version = 8
```

not `7`.

So B knows:

> "Someone modified this record after I read it."

The application can then decide what to do:

- reject the update and tell the user,
- reload the latest version,
- ask the user to resolve the conflict,
- or retry automatically if that makes sense.

This is generally called an **optimistic concurrency check**.

---

# 4. Why this is different from a transaction

This distinction is important.

A transaction answers:

> **Which database operations succeed or fail together?**

Optimistic concurrency answers:

> **Did someone modify the data between the time I read it and the time I tried to modify it?**

They solve related but different problems.

You can use both.

For example:

```text
BEGIN TRANSACTION

UPDATE Listings
WHERE Id = 123
AND Version = 7

if rows affected = 0:
    rollback / conflict

COMMIT
```

The transaction gives you atomicity, while the version check detects the concurrent modification.

---

# 5. What is the `Version` column actually doing?

It doesn't have to literally be an integer.

Common approaches include:

## Integer version

```text
Version = 1
Version = 2
Version = 3
...
```

Simple and easy to understand.

## Timestamp/row-version mechanism

SQL Server has a `rowversion` data type specifically useful for optimistic concurrency.

The application doesn't need to decide the next version value itself; SQL Server generates it when the row changes.

## EF Core concurrency token

EF Core supports this concept directly.

For example, you can mark a property as a concurrency token. EF Core then effectively includes the original value in the `UPDATE` condition.

If somebody else changed the row, EF Core can detect that **zero rows were affected** and raise a concurrency exception.

This is one of those cases where understanding the underlying SQL is more valuable than memorizing the EF Core API.

---

# 6. When would you choose pessimistic vs optimistic?

A useful interview answer is:

**Optimistic concurrency** is usually attractive when:

- conflicts are relatively uncommon,
- you want high concurrency,
- users may read data and modify it later,
- you can gracefully handle conflicts.

**Pessimistic concurrency** is useful when:

- conflicts are frequent,
- conflicting updates are expensive,
- you really need exclusive access while performing a sequence of operations.

For example, imagine inventory:

```text
Stock = 1
```

Two customers simultaneously purchase the final item.

You cannot simply let both applications read `Stock = 1` and then independently write `Stock = 0`.

You need a concurrency-safe operation.

Depending on the system, you might use locking or an atomic conditional update such as:

```sql
UPDATE Products
SET Stock = Stock - 1
WHERE Id = 123
  AND Stock > 0;
```

Then check the number of affected rows.

If it is `1`, the purchase succeeded.

If it is `0`, there wasn't available stock (or the row didn't match).

This is an excellent example of why **atomic database operations** can sometimes be preferable to "read, calculate, then write."

---

# 7. One more subtle point: isolation level vs locking

Don't mix these concepts.

**Isolation level** defines the consistency guarantees between concurrent transactions.

**Locks** are one mechanism SQL Server can use to implement those guarantees.

And SQL Server also has **row-versioning** mechanisms, such as snapshot-based isolation, where readers can work from versions of rows rather than simply waiting on writers.

So if an interviewer asks:

> "How do you prevent two users from updating the same row?"

Don't immediately answer:

> "Use a transaction."

Instead, explain the actual concurrency strategy:

> "It depends on the operation. For an application where conflicts are relatively rare, I'd typically use optimistic concurrency with a row version/concurrency token. The update includes the version that was originally read, and if zero rows are affected, I know another transaction modified the record. For cases where exclusive access is required, I'd consider pessimistic locking."

That's a much stronger answer.

---

## The transaction topics we've now covered

- Transaction boundaries
- Atomicity and rollback
- Isolation levels
- Dirty / non-repeatable / phantom reads
- Savepoints
- EF Core transactions
- Pessimistic concurrency
- Optimistic concurrency
- Lost updates
- Version/concurrency tokens
- Blocking implications

**Next:** the last major transaction topic is **deadlocks and practical transaction problems**—how they happen, how SQL Server detects them, and how you'd prevent/debug them.