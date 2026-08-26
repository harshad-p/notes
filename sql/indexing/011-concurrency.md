Yes. I'll keep the remaining roadmap explicit so we can stop and resume without restarting topics.

### Remaining topics

1. **Locking, blocking, deadlocks & isolation levels** ← **next**
2. **Advanced SQL performance diagnosis**
   - execution plans in practice
   - parameter sniffing / parameter-sensitive plans
   - SARGability
   - identifying the actual bottleneck
3. **Transactions in depth**
   - transaction boundaries
   - atomicity/rollback
   - savepoints
   - EF Core transactions
4. **Advanced indexing**
   - filtered indexes
   - included columns
   - covering indexes
   - index design trade-offs
   - when *not* to index
5. **SQL Server practical troubleshooting**
   - Query Store
   - DMVs
   - Extended Events / Profiler
   - finding expensive queries

That's **5 major topics** remaining. Some have several subtopics, but I won't turn every subtopic into a separate lesson unless it's genuinely necessary.

And yes, if you pause after this one and come back later, we'll use this roadmap rather than starting over.

---

# 11. Locking, blocking, deadlocks, and isolation levels

You've already encountered optimistic and pessimistic concurrency, so let's go one level deeper into what SQL Server is actually doing.

## 1. Why does SQL Server need locks?

Imagine two transactions simultaneously modifying the same row.

The row contains:

```text
Balance = 100
```

Transaction A wants to withdraw 80.

Transaction B wants to withdraw 50.

If both transactions read `100` before either writes its result, you could end up with:

```text
A calculates: 100 - 80 = 20
B calculates: 100 - 50 = 50
```

and one update overwrites the other.

You have effectively lost one operation.

SQL Server needs mechanisms to coordinate concurrent access.

**Locks are one of those mechanisms.**

---

# 2. What does a lock actually protect?

A lock isn't simply:

> "This row belongs to Transaction A."

It's more precise than that.

SQL Server can acquire locks on different resources, including:

- rows
- keys
- pages
- tables
- databases
- other internal resources

And locks have different modes.

The two most important ones to understand initially are:

### Shared (`S`)

Used when reading data under locking-based isolation.

Multiple transactions can generally hold compatible shared locks on the same resource.

### Exclusive (`X`)

Used when modifying data.

An exclusive lock prevents other transactions from performing operations that conflict with that modification.

So if Transaction A has an exclusive lock on a row, Transaction B can't simply modify that same row at the same time.

---

# 3. A simple example

Transaction A:

```sql
BEGIN TRANSACTION;

UPDATE Accounts
SET Balance = Balance - 100
WHERE Id = 1;
```

Transaction A hasn't committed yet.

Transaction B tries:

```sql
UPDATE Accounts
SET Balance = Balance - 50
WHERE Id = 1;
```

Transaction B may have to **wait** because Transaction A currently holds a lock that conflicts with B's requested lock.

If A eventually executes:

```sql
COMMIT;
```

the lock is released and B can continue.

This waiting is **blocking**.

---

# 4. Blocking is not automatically a problem

This distinction is important.

People sometimes say:

> "Blocking is bad."

Not exactly.

If Transaction A modifies a row and Transaction B needs that same row, **some coordination is necessary**.

If A holds the lock for 2 milliseconds and B waits 1 millisecond, that's completely normal.

The problem is **excessive or prolonged blocking**.

For example:

```text
Transaction A
    |
    | holds lock for 30 seconds
    |
    ↓
Transaction B waits
Transaction C waits
Transaction D waits
Transaction E waits
```

Now a single transaction is potentially causing a chain of requests to stall.

---

# 5. Why would a transaction hold a lock for 30 seconds?

This is where application design matters.

Consider:

```csharp
using var transaction = await db.Database.BeginTransactionAsync();

var user = await GetUserAsync();

await CallExternalServiceAsync();

await UpdateUserAsync();

await transaction.CommitAsync();
```

Imagine `CallExternalServiceAsync()` takes 10 seconds.

Depending on the operations and isolation level, locks acquired by the transaction may remain held while you're waiting for that external service.

That's dangerous.

You've effectively done this:

> Start database transaction → acquire locks → leave database → wait 10 seconds → come back → finish transaction.

A transaction should generally be kept **as short as practical**.

---

# 6. Lock escalation

SQL Server doesn't necessarily maintain millions of individual row locks forever.

If a transaction acquires a large number of locks, SQL Server can perform **lock escalation**.

For example, it may move from many finer-grained locks toward a table-level lock.

Why?

Because maintaining huge numbers of individual locks has memory and management overhead.

But escalation can increase blocking.

Instead of:

> "I'm locking these particular rows."

you can end up with something much closer to:

> "I'm locking a much larger portion of this table."

So large transactions can have significant concurrency consequences.

---

# 7. Now: what is a deadlock?

Blocking is:

> A waits for B.

A deadlock is:

> A waits for B **while B waits for A**.

For example:

Transaction A:

```sql
UPDATE Accounts
SET Balance = Balance - 100
WHERE Id = 1;
```

A now holds a lock on Account 1.

Then A tries to update Account 2.

Meanwhile Transaction B did:

```sql
UPDATE Accounts
SET Balance = Balance - 50
WHERE Id = 2;
```

B now holds a lock on Account 2.

Then B tries to update Account 1.

You now have:

```text
Transaction A
holds Account 1
    ↓
wants Account 2
    ↑
Transaction B
holds Account 2
    ↓
wants Account 1
```

Neither can proceed.

That's a **deadlock**.

---

# 8. SQL Server has to break the deadlock

Neither transaction can magically continue.

SQL Server detects the deadlock and chooses one transaction as the **deadlock victim**.

That transaction is rolled back, releasing its locks.

The other transaction can then continue.

Your application will receive an error for the victim transaction.

This is why deadlocks are different from ordinary blocking.

**Blocking can resolve naturally.**

**A deadlock cannot resolve unless one participant is aborted.**

---

# 9. How do you prevent deadlocks?

One of the most effective techniques is **consistent access ordering**.

Suppose your application always updates accounts in ascending ID order.

Transaction A:

```text
Account 1
Account 2
```

Transaction B also does:

```text
Account 1
Account 2
```

Then both transactions attempt to acquire locks in the same order.

You greatly reduce the possibility of:

```text
A owns 1 → wants 2
B owns 2 → wants 1
```

because B would try to acquire Account 1 first.

This is a very practical deadlock-prevention technique.

Other causes can include:

- unnecessarily long transactions
- missing indexes causing more rows to be touched
- inconsistent access patterns
- inappropriate transaction scope
- high concurrency

---

# 10. Isolation levels

Now we get to the more interesting part.

An isolation level determines **how much one transaction is isolated from the effects of other concurrent transactions**.

SQL Server provides several isolation levels:

1. Read Uncommitted
2. Read Committed
3. Repeatable Read
4. Serializable
5. Snapshot

And SQL Server also supports **Read Committed Snapshot Isolation (RCSI)** as a database-level option.

These aren't just interview vocabulary.

They change the behavior of reads and writes.

---

# 11. Read Uncommitted

This is the least restrictive standard isolation level.

A transaction can read data that another transaction has modified but **hasn't committed yet**.

That's called a **dirty read**.

Example:

Transaction A:

```text
Balance = 100
```

A changes it to:

```text
Balance = 0
```

but hasn't committed.

Transaction B reads the balance and sees:

```text
0
```

Then A rolls back.

The real balance returns to:

```text
100
```

B just read a value that was never actually committed.

That's a dirty read.

You can explicitly request this behavior with:

```sql
SELECT *
FROM Accounts WITH (NOLOCK);
```

But `NOLOCK` is **not** simply a magical performance optimization.

It effectively requests read-uncommitted behavior for that table reference, and comes with correctness trade-offs.

---

# 12. Read Committed

This is the default isolation level in SQL Server unless configured otherwise.

The basic guarantee is:

> You don't read uncommitted changes from another transaction.

If another transaction is modifying a row and hasn't committed, a locking-based Read Committed query may have to wait.

So:

```text
Transaction A
modifies row
     |
     | lock
     ↓
Transaction B
tries to read row
     |
     ↓
waits
```

Once A commits or rolls back, B can proceed.

However, Read Committed does **not** guarantee that if you read the same row twice within one transaction, you'll necessarily get the same value both times.

That leads us to Repeatable Read.

---

# 13. Repeatable Read

Suppose Transaction A reads:

```text
Price = 100
```

and keeps the transaction open.

Under Repeatable Read, SQL Server prevents another transaction from modifying that row in a way that would make A's repeated read return a different value.

So A reads:

```text
100
```

then B tries to change it to:

```text
200
```

B may have to wait until A finishes.

That's the "repeatable" part:

> If I read a row and continue my transaction, another transaction can't modify that row underneath me in a way that changes my subsequent read.

---

# 14. But Repeatable Read doesn't prevent new rows

This is a subtle but important distinction.

Suppose you execute:

```sql
SELECT *
FROM Orders
WHERE CustomerId = 10;
```

You get:

```text
Order 1
Order 2
```

Another transaction could potentially insert:

```text
Order 3
CustomerId = 10
```

Your second execution could then see:

```text
Order 1
Order 2
Order 3
```

That's called a **phantom read**.

Repeatable Read protects rows you've already read, but it doesn't necessarily prevent another transaction from inserting new rows matching your predicate.

---

# 15. Serializable

Serializable is the strongest traditional locking-based isolation level.

It prevents the kind of phantom behavior we just described.

Conceptually, if you execute:

```sql
SELECT *
FROM Orders
WHERE CustomerId = 10;
```

under Serializable, SQL Server needs to protect not only the rows currently matching the predicate but also the **range** where another matching row could be inserted.

That's why Serializable can cause considerably more blocking.

You're essentially asking:

> "Make this transaction behave as if concurrent transactions aren't allowed to interfere with this range of data."

Very strong consistency.

Potentially much lower concurrency.

---

# 16. Snapshot

Snapshot takes a fundamentally different approach.

Instead of making readers wait for writers in many situations, SQL Server can give a transaction a **consistent version of the data**.

Suppose:

Transaction A updates:

```text
Price = 200
```

while Transaction B is reading.

Under snapshot-based behavior, B can continue reading the version of the row appropriate to its transaction's snapshot rather than necessarily waiting for A's write lock.

SQL Server maintains row versions to make this possible.

Those versions are stored using `tempdb`.

And now you can see another connection to the previous topic.

We just talked about `tempdb` being used for temporary/internal work.

**Row versioning is another reason `tempdb` matters.**

---

# 17. RCSI vs Snapshot

These two are easy to confuse.

### Read Committed Snapshot Isolation

Changes the behavior of the **Read Committed** isolation level so that reads can use row versions instead of taking the traditional shared locks in many cases.

It is enabled at the database level.

### Snapshot isolation

A transaction explicitly uses Snapshot isolation, and reads are based on a transaction-level consistent snapshot.

The exact semantics differ, so don't treat them as interchangeable names.

---

# 18. Why would you use row versioning?

One major benefit is reducing reader/writer blocking.

Imagine an application where:

- lots of requests read data
- occasional requests update data
- you don't want ordinary reads sitting behind short-lived writes

Row versioning can allow readers to continue using an appropriate committed version rather than waiting for the writer.

But there's a trade-off:

**More row-version storage and `tempdb` activity.**

So again, there isn't a universally "best" isolation level.

---

# 19. The practical hierarchy

A useful way to think about the isolation levels is:

**Read Uncommitted**

"I accept dirty/inconsistent reads."

**Read Committed**

"Don't show me uncommitted data."

**Repeatable Read**

"Don't let rows I've read change underneath me."

**Serializable**

"Also protect the ranges so matching rows can't appear/disappear underneath me."

**Snapshot**

"Give me a consistent version of the data instead of relying primarily on readers blocking writers."

The exact implementation details are more nuanced, especially once row-versioning options are involved, but that's the conceptual progression.

---

# 20. What would you say in a senior interview?

If asked:

> "How would you handle two users updating the same row?"

Don't immediately say:

> "Use a lock."

I'd say something like:

> "It depends on the consistency requirement. I'd first consider optimistic concurrency, for example using a row-version/concurrency token so that the second update detects that the row changed since it was read. If conflicts are expected to be frequent and I need to prevent concurrent modification up front, pessimistic locking may be more appropriate. I'd also consider the transaction scope and isolation level because those affect blocking and concurrency."

That's much stronger than simply naming a locking mechanism.

---

### One final connection

You now have three separate concepts that often get mixed together:

**Concurrency control** — how simultaneous operations are coordinated.

**Isolation** — what one transaction is allowed to observe from other transactions.

**Locking / row versioning** — mechanisms SQL Server uses to implement those guarantees.

They're related, but they're not synonyms.

And that's the level of distinction I want you to have going forward.

**Roadmap remains: 4 major topics after this one:**

1. Advanced SQL performance diagnosis
2. Transactions in depth
3. Advanced indexing
4. SQL Server practical troubleshooting

We can stop there whenever you want and resume from exactly this roadmap later.