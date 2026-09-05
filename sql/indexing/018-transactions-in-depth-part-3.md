# Transactions in Depth — Part 3: Deadlocks

A **deadlock** is one of the most important practical transaction problems to understand because it can happen even when every individual transaction is logically correct.

## 1. What is a deadlock?

A deadlock occurs when **two or more transactions are waiting for each other to release resources**, so none of them can continue.

Consider two rows:

```text
A
B
```

Transaction 1:

```text
locks A
   ↓
tries to lock B
```

Transaction 2:

```text
locks B
   ↓
tries to lock A
```

Now:

- Transaction 1 cannot continue until Transaction 2 releases B.
- Transaction 2 cannot continue until Transaction 1 releases A.

Neither can proceed.

That's a deadlock.

---

# 2. Concrete SQL example

Transaction 1:

```sql
BEGIN TRANSACTION;

UPDATE Accounts
SET Balance = Balance - 100
WHERE Id = 1;

-- later...
UPDATE Accounts
SET Balance = Balance + 100
WHERE Id = 2;

COMMIT;
```

Transaction 2 starts at roughly the same time:

```sql
BEGIN TRANSACTION;

UPDATE Accounts
SET Balance = Balance - 50
WHERE Id = 2;

-- later...
UPDATE Accounts
SET Balance = Balance + 50
WHERE Id = 1;

COMMIT;
```

The sequence could become:

```text
Transaction 1                  Transaction 2

locks Account 1
                                locks Account 2

tries Account 2
                                tries Account 1

WAITING                         WAITING
```

That's the circular dependency that creates the deadlock.

---

# 3. What does SQL Server do?

SQL Server doesn't simply leave both transactions waiting forever.

It has **deadlock detection**.

It detects the cycle and chooses one transaction as the **deadlock victim**.

Conceptually:

```text
Transaction 1 ──waiting──> Transaction 2
     ↑                         │
     └────────waiting──────────┘

              ↓

SQL Server detects cycle

              ↓

Kills/rolls back one transaction

              ↓

Other transaction continues
```

The application receives a deadlock error for the victim transaction.

This is important:

**A deadlock is not necessarily a database failure.**

It is SQL Server resolving an impossible waiting situation by sacrificing one transaction.

---

# 4. Why do deadlocks happen in real applications?

Usually because transactions acquire locks in **different orders**.

For example:

Transaction A:

```text
Customer → Order
```

Transaction B:

```text
Order → Customer
```

If both execute concurrently, they can end up waiting for each other.

That's why one of the simplest and most effective deadlock-prevention techniques is:

> **Acquire resources in a consistent order.**

For example, always update:

```text
Customer first
Order second
```

instead of allowing some code paths to do the reverse.

---

# 5. Transaction duration matters

Remember our earlier rule:

> Keep transactions as short as the business requirement allows.

Long transactions hold locks longer.

Imagine:

```sql
BEGIN TRANSACTION;

UPDATE Orders
SET Status = 'PROCESSING'
WHERE Id = 123;

-- 10 seconds of application work

-- external API call
-- payment processing
-- HTTP request
-- PDF generation

COMMIT;
```

You've potentially held database resources while doing work that has nothing to do with the database.

That increases the opportunity for:

- blocking,
- lock contention,
- deadlocks,
- transaction log growth,
- poor throughput.

So instead, do as much non-DB work as possible **outside** the transaction.

---

# 6. Blocking vs deadlock

These are often confused in interviews.

## Blocking

Transaction A holds a lock.

Transaction B wants that resource.

So B waits.

```text
A ──holds lock──> Row 1

B ──waiting─────> Row 1
```

Eventually A commits, and B continues.

That's **blocking**.

## Deadlock

A waits for B **while B waits for A**.

```text
A waits for B
B waits for A
```

Neither can ever proceed without intervention.

That's **deadlock**.

So:

> **Every deadlock involves waiting, but not every wait is a deadlock.**

---

# 7. Should you simply retry deadlocks?

Often, yes—but **carefully**.

A deadlock victim can retry the transaction because the failure is generally transient.

For example:

```text
Attempt 1
    ↓
deadlock
    ↓
short delay
    ↓
Attempt 2
    ↓
success
```

Typically you'd use a limited retry policy with **exponential backoff**.

For example:

```text
retry 1 → small delay
retry 2 → longer delay
retry 3 → longer again
then fail
```

But don't blindly retry everything.

If the transaction fails because of:

- invalid data,
- constraint violation,
- permission problem,
- bad SQL,

retrying won't fix it.

A deadlock is different because another concurrent transaction may have caused a temporary conflict.

---

# 8. How would you reduce deadlocks?

In an interview, I'd give several concrete strategies:

## 1. Consistent lock/resource ordering

Make different transactions acquire resources in the same order.

## 2. Keep transactions short

Don't perform unnecessary work inside them.

## 3. Access only the data you need

Unnecessarily touching many rows can increase locking and contention.

## 4. Use appropriate indexes

This one is subtle.

Suppose an update has to find a row by `CustomerId`, but there is no useful index. SQL Server might need to examine many rows, potentially increasing the amount of locking and the duration of the operation.

A good index can therefore help **indirectly** with concurrency.

But don't say:

> "Indexes prevent deadlocks."

They don't. They can reduce the likelihood by making operations faster and more targeted.

## 5. Use appropriate isolation levels

Higher isolation can require more locking or versioning resources, so don't automatically choose `SERIALIZABLE` everywhere.

## 6. Retry transient deadlocks

Have the application retry when appropriate, with a bounded backoff.

---

# 9. How would you investigate a deadlock?

This is where your earlier **SQL Server troubleshooting** topic will eventually connect.

You don't want to just say:

> "I'd add retries."

First you want to know **why the deadlock happened**.

SQL Server can capture deadlock information, including:

- which transactions were involved,
- which resources they held,
- which resources they were waiting for,
- which transaction was chosen as the victim.

The classic representation is called a **deadlock graph**.

It lets you see something like:

```text
Transaction A
  holds: Row A
  waits for: Row B

Transaction B
  holds: Row B
  waits for: Row A
```

Then you fix the underlying transaction behavior rather than merely hiding the symptom with retries.

We'll cover the SQL Server tooling for this later when we get to **Extended Events and DMVs**.

---

# 10. One interview-quality answer

If you're asked:

> **"What happens if two transactions update the same data and cause a deadlock?"**

A strong answer would be:

> "A deadlock occurs when transactions hold resources that the other transaction needs, creating a circular wait. SQL Server detects the deadlock and chooses one transaction as the victim, rolling it back so the other transaction can continue. I'd investigate the deadlock graph to identify the conflicting operations, then try to prevent it by keeping transactions short, acquiring resources in a consistent order, and making sure queries use appropriate indexes. For transient deadlock errors, the application can also retry the transaction with a bounded exponential backoff."

That's the level of answer I'd expect from someone with your experience.

---

## Transactions are now essentially complete.

The remaining SQL track we had planned is:

**SQL Server practical troubleshooting**
1. Query Store
2. DMVs
3. Extended Events
4. Finding expensive queries
5. Putting all of this together into a real troubleshooting workflow

That's the next topic.