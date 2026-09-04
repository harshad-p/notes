# Transactions in Depth

Let's move to the next topic from our agreed plan.

The important thing here is to get beyond **"a transaction makes several SQL statements atomic."** That's true, but there's considerably more going on.

---

## 1. What is a transaction actually for?

Imagine a bank transfer:

```sql
UPDATE Accounts
SET Balance = Balance - 100
WHERE Id = 1;

UPDATE Accounts
SET Balance = Balance + 100
WHERE Id = 2;
```

There are two separate database operations.

We don't want this situation:

- Account 1 loses $100.
- Application crashes.
- Account 2 never receives the $100.

The two operations represent **one logical business operation**.

A transaction allows us to tell the database:

> "Treat these operations as one unit. Either the whole unit succeeds, or none of it does."

For example:

```sql
BEGIN TRANSACTION;

UPDATE Accounts
SET Balance = Balance - 100
WHERE Id = 1;

UPDATE Accounts
SET Balance = Balance + 100
WHERE Id = 2;

COMMIT;
```

If something goes wrong before `COMMIT`:

```sql
ROLLBACK;
```

The changes made inside that transaction are undone.

---

# 2. Atomicity

This is the **A** in ACID.

Atomicity means that the transaction is treated as one indivisible unit from the database's perspective.

For example:

```text
Transaction:
    deduct $100
    credit $100
```

Either both changes become part of the committed database state, or neither does.

It isn't simply a programming-language concept.

The database maintains internal information that allows it to recover from failures and undo uncommitted changes.

---

# 3. `COMMIT` is the important boundary

Consider:

```sql
BEGIN TRANSACTION;

UPDATE Orders
SET Status = 'Paid'
WHERE Id = 123;

-- lots of application work here

COMMIT;
```

Before `COMMIT`, the transaction has not been successfully completed.

After `COMMIT`, the database has accepted the transaction as committed.

This gives us an important design principle:

> **The transaction boundary should normally correspond to the smallest logical unit of work that must succeed together.**

You don't want to keep a transaction open unnecessarily.

---

# 4. Why transaction boundaries matter

Imagine your API does this:

```text
BEGIN TRANSACTION

Update order

Call payment provider

Send email

Generate PDF

COMMIT
```

That's generally a terrible transaction boundary.

Why?

Because you're keeping database resources involved while waiting for external systems.

The payment provider might take 3 seconds.

The email service might take 2 seconds.

Now your database transaction is sitting open for several seconds.

That can mean:

- locks are held longer
- other transactions may have to wait
- blocking increases
- deadlocks become more likely
- database resources are occupied unnecessarily

Instead, you generally want the database transaction to contain **only the database work that must be atomic**.

External operations need a different reliability strategy.

---

# 5. Transactions don't magically make everything consistent

This is a subtle but important distinction.

Suppose:

```text
Database transaction
    ↓
Update order
    ↓
COMMIT
    ↓
Send email
```

The database transaction guarantees the database changes.

It does **not** guarantee:

> "The email will definitely be sent."

If the application crashes after the commit but before sending the email, the order is paid, but the email wasn't sent.

This is one reason patterns such as **transactional outbox** exist.

We'll come back to that when we discuss distributed systems/reliability, because it sits at the boundary between database transactions and messaging.

---

# 6. Isolation

Atomicity is only one part of transactions.

The other major concept we need to understand is **isolation**.

Imagine two transactions operating at the same time.

Transaction A:

```sql
UPDATE Accounts
SET Balance = Balance - 100
WHERE Id = 1;
```

Transaction B simultaneously tries to read that account.

What should B see?

Should it see:

- the old balance?
- the new uncommitted balance?
- wait until A commits?
- some consistent snapshot?

That's what **transaction isolation** deals with.

It determines how much one transaction is allowed to observe the intermediate effects of another transaction.

---

# 7. The classic isolation problems

There are several phenomena you need to know.

### Dirty read

Transaction A changes something but hasn't committed.

Transaction B reads the changed value.

Then A rolls back.

B has read data that never actually existed in the committed database state.

Example:

```text
A: balance = $1,000 → $500

B: reads $500

A: ROLLBACK

Actual committed balance = $1,000
```

B read a **dirty** value.

---

### Non-repeatable read

Transaction A reads a row:

```text
Price = $500
```

Transaction B changes it and commits:

```text
Price = $600
```

Transaction A reads the same row again within its transaction.

Now it sees:

```text
Price = $600
```

The same query produced a different value during the same transaction.

---

### Phantom read

This involves **rows appearing or disappearing** between reads.

Suppose transaction A runs:

```sql
SELECT *
FROM Listings
WHERE Price < 500000;
```

It gets 100 listings.

Transaction B inserts another listing satisfying that condition and commits.

Transaction A runs the query again.

Now it gets 101 listings.

The new row is a **phantom**.

---

# 8. SQL Server isolation levels

SQL Server provides different isolation levels that determine these behaviors.

The commonly important ones are:

### READ UNCOMMITTED

Very weak isolation.

A transaction can read changes that haven't been committed.

This permits dirty reads.

You may encounter:

```sql
SELECT *
FROM Listings WITH (NOLOCK);
```

`NOLOCK` is associated with this kind of behavior.

It's important not to treat `NOLOCK` as a magical performance optimization. You're explicitly accepting weaker consistency and potentially problematic reads.

---

### READ COMMITTED

This is the normal default isolation level in SQL Server.

Broadly, you don't read another transaction's uncommitted changes.

Depending on the database configuration and whether row-versioning is enabled, the exact behavior differs in terms of how reads interact with writers.

---

### REPEATABLE READ

Once your transaction reads a row, SQL Server prevents another transaction from modifying that row in a way that would make your repeated read return a different value.

It's stronger than READ COMMITTED.

But it can result in more locking/blocking.

---

### SERIALIZABLE

This provides the strongest traditional isolation.

It prevents other transactions from inserting rows that would satisfy certain range predicates while your transaction is operating on that range.

This addresses phantom reads.

But the price is potentially significant locking and reduced concurrency.

---

### SNAPSHOT

SQL Server can also provide snapshot-based isolation using row versions.

Instead of necessarily waiting for another transaction's changes, a reader can see a consistent version of the data from an appropriate point in time.

This can substantially reduce reader/writer blocking.

But it isn't free: SQL Server has to maintain row versions, which creates additional storage/tempdb-related overhead.

---

# 9. The trade-off

This is the central idea behind isolation levels:

> **Stronger isolation gives you stronger consistency guarantees, but can reduce concurrency and increase resource usage.**

You therefore don't automatically say:

> "Use SERIALIZABLE because it's safest."

You choose the level appropriate for the business operation.

For many ordinary application queries, `READ COMMITTED` is sufficient.

For some workloads, snapshot-based isolation can be useful.

For particularly sensitive operations, stronger isolation may be justified.

---

# 10. Transactions and concurrent updates

This connects directly to a question you were asked in your interview:

> **"How would you ensure data consistency if the same row in the database is being updated?"**

There isn't one answer.

Suppose two requests simultaneously modify:

```text
Listing 123
Price = $500,000
```

Request A wants:

```text
Price = $510,000
```

Request B wants:

```text
Price = $520,000
```

You need to determine what semantics the application wants.

One option is **pessimistic concurrency**: use database locking so conflicting operations can't proceed simultaneously.

Another is **optimistic concurrency**: allow concurrent work but detect whether someone changed the row before committing your update.

For example, the table might have a version/concurrency token.

Conceptually:

```sql
UPDATE Listings
SET Price = 520000
WHERE Id = 123
  AND Version = 7;
```

If another request has already changed the row and its version is now 8:

```text
Rows affected = 0
```

The application knows:

> "The row changed since I read it."

It can then reject the update or ask the caller to retry/reload.

This is a very important concept for EF Core, so we'll cover it properly when we get to EF Core transactions/concurrency.

---

# 11. Savepoints

Now let's go one level deeper.

A transaction doesn't necessarily have to be:

```text
BEGIN
    everything
COMMIT
```

You can establish a **savepoint** inside a transaction.

Conceptually:

```sql
BEGIN TRANSACTION;

-- operation A

SAVE TRANSACTION MySavepoint;

-- operation B

ROLLBACK TRANSACTION MySavepoint;

-- operation C

COMMIT;
```

The rollback doesn't necessarily undo the entire transaction.

It rolls back to the savepoint.

So:

- operation A remains
- operation B is undone
- operation C can still happen
- the transaction can still be committed

This is useful when part of a larger transaction is optional or recoverable.

---

# 12. EF Core

In EF Core, `SaveChanges()` itself is transactional for the changes it sends as part of that call.

For example:

```csharp
context.Listings.Add(listing);
context.ModerationResults.Add(result);

await context.SaveChangesAsync();
```

EF Core will ensure those changes are persisted atomically as part of the save operation.

But suppose you have:

```csharp
await context.SaveChangesAsync();

await DoSomethingElseAsync();

await context.SaveChangesAsync();
```

Those aren't automatically one giant transaction.

If both database operations need to succeed or fail together, you may explicitly create a transaction.

For example:

```csharp
await using var transaction =
    await context.Database.BeginTransactionAsync();

try
{
    // database operation 1

    await context.SaveChangesAsync();

    // database operation 2

    await context.SaveChangesAsync();

    await transaction.CommitAsync();
}
catch
{
    await transaction.RollbackAsync();
    throw;
}
```

The exact structure can vary, but the important idea is:

> **Explicitly define the transaction when multiple database operations must share one atomic boundary.**

---

# 13. One important EF Core warning

Remember your earlier discussion about `DbContext`.

`DbContext` isn't designed for concurrent operations.

So you shouldn't do:

```csharp
await Task.WhenAll(
    context.Users.ToListAsync(),
    context.Listings.ToListAsync()
);
```

on the same `DbContext`.

Transactions don't change that.

A transaction answers:

> "Which database operations belong to one atomic unit?"

`DbContext` concurrency answers:

> "Can this particular EF Core context execute operations concurrently?"

They're related concepts, but they're **not the same thing**.

---

# What I want you to retain

Don't memorize the SQL syntax first.

The important hierarchy is:

**Transaction**

A group of database operations that must be treated as one logical unit.

**Atomicity**

Either the transaction commits as a whole or its changes are rolled back.

**Transaction boundary**

Defines exactly which operations belong to that unit.

**Isolation**

Controls what concurrent transactions can observe and how they interact.

**Savepoint**

Allows partial rollback inside a larger transaction.

**EF Core**

`SaveChanges()` provides transactional behavior for its own save operation; explicit transactions are used when you need multiple operations to share one transaction.

And perhaps the most important senior-level principle:

> **Keep transactions as small as the business requirement allows. Don't hold a database transaction open while waiting on external services or doing unrelated work.**

### Next

We've covered the first major pieces of **Transactions in Depth**. The next part should be **concurrency + isolation in considerably more detail**, particularly **optimistic vs pessimistic concurrency, row/version locking, and how SQL Server actually behaves when two requests update the same row**. That will connect directly to the interview question you were asked.