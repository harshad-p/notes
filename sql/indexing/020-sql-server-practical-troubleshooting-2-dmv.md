# SQL Server Practical Troubleshooting — Part 2: DMVs

**DMVs (Dynamic Management Views)** are one of the main ways you inspect what SQL Server is doing.

The key distinction from Query Store is:

- **Query Store:** historical performance — *"What has been happening?"*
- **DMVs:** current/accumulated server state — *"What is happening, and what has SQL Server observed?"*

You don't need to memorize dozens of DMV names. For interviews, understand the important categories and what questions they answer.

---

## 1. Finding expensive queries

One useful DMV is:

```sql
sys.dm_exec_query_stats
```

It contains execution statistics for cached query plans.

You can combine it with:

```sql
sys.dm_exec_sql_text
```

to retrieve the actual SQL text.

For example:

```sql
SELECT TOP 10
    qs.execution_count,
    qs.total_worker_time,
    qs.total_logical_reads,
    qs.total_elapsed_time,
    st.text
FROM sys.dm_exec_query_stats AS qs
CROSS APPLY sys.dm_exec_sql_text(qs.sql_handle) AS st
ORDER BY qs.total_worker_time DESC;
```

Don't worry about memorizing the exact syntax yet. Understand what we're doing.

We're essentially asking:

> "Among the queries whose statistics are currently available, which ones have consumed the most CPU?"

---

## 2. CPU vs duration vs reads

This is where your earlier performance-diagnosis knowledge becomes useful.

### `total_worker_time`

Roughly represents **CPU consumed** by the query executions.

If this is very high, investigate CPU-heavy work.

For example:

- expensive joins,
- sorting,
- aggregation,
- expressions,
- poor execution plans.

### `total_elapsed_time`

How much elapsed time the executions consumed.

High elapsed time doesn't necessarily mean high CPU.

A query could spend most of its time **waiting**.

### `total_logical_reads`

How much data SQL Server had to read from memory/cache.

Very high logical reads can point toward:

- scans,
- inefficient predicates,
- poor indexing,
- excessive lookups,
- queries processing far more rows than necessary.

---

# 3. Average vs total matters

This is a subtle but important point.

Suppose:

```text
Query A
Average duration: 5 seconds
Executions:       2

Query B
Average duration: 20 ms
Executions:       2,000,000
```

If you sort only by average duration, Query A looks terrible.

But Query B may be consuming far more total resources.

So when investigating production performance, you should consider both:

**Per-execution cost**

and

**Aggregate cost**

The DMV statistics can help you examine both.

---

# 4. Finding currently running queries

Another useful DMV is:

```sql
sys.dm_exec_requests
```

This is particularly useful for:

> "What is SQL Server doing **right now**?"

You can find information such as:

- currently executing requests,
- elapsed time,
- CPU time,
- wait information,
- blocking session,
- database,
- SQL text.

This is different from `sys.dm_exec_query_stats`.

`dm_exec_query_stats` is primarily about **accumulated statistics for cached plans**.

`dm_exec_requests` is about **currently executing requests**.

---

# 5. Finding blocking

This is one of the practical questions you might get in an interview:

> "The database is slow. How would you find out if blocking is the cause?"

You can inspect currently running requests and their blocking relationships.

A simplified query might look like:

```sql
SELECT
    session_id,
    blocking_session_id,
    wait_type,
    wait_time,
    status
FROM sys.dm_exec_requests
WHERE blocking_session_id <> 0;
```

Suppose you see:

```text
Session 52
blocking_session_id = 41
```

That means session 52 is waiting because session 41 is blocking it.

Now you investigate session 41.

Perhaps it is running:

```sql
BEGIN TRANSACTION;

UPDATE Orders
SET Status = 'PROCESSING'
WHERE Id = 123;

-- application hasn't committed yet
```

And that's why another request is waiting.

---

# 6. A common production situation

Imagine users report:

> "Everything is slow."

You check the application and find database requests taking 10 seconds.

You inspect the query itself and discover:

```text
CPU:          20 ms
Elapsed time: 10,000 ms
```

That's an enormous clue.

The query isn't necessarily computationally expensive.

It's spending almost all its time **waiting**.

You then investigate:

- wait type,
- blocking session,
- locks,
- transaction duration.

And discover another transaction has held a lock for 9.9 seconds.

So the problem wasn't:

> "The query needs an index."

It was:

> **"The query is blocked."**

This is exactly why blindly optimizing the query is dangerous.

---

# 7. Waits

SQL Server has many types of waits.

The general concept is:

> **A wait tells you what a worker is waiting for.**

For example, a query might be waiting for:

- another transaction's lock,
- disk/I/O,
- memory,
- CPU scheduling,
- parallelism-related resources,
- network/client activity.

You don't need to memorize every wait type.

The useful interview-level skill is recognizing:

> **A slow query isn't necessarily doing work; it may be waiting for something.**

That's a major part of database troubleshooting.

---

# 8. Finding missing/unused indexes

There are also DMVs related to index usage.

For example:

```sql
sys.dm_db_index_usage_stats
```

This can give you information about things such as:

- seeks,
- scans,
- lookups,
- updates.

This can help answer:

> "Is this index actually being used?"

For example, suppose a table has:

```text
Index A
Index B
Index C
Index D
Index E
Index F
```

and you discover some indexes receive virtually no reads but are constantly maintained because the table is frequently updated.

That can be a sign that the indexes deserve investigation.

But again, **don't automatically delete an apparently unused index**.

DMV statistics have limitations:

- they reset under certain circumstances,
- usage may be workload-dependent,
- an index might be used rarely but be critical when it is used.

So this is evidence, not an automatic verdict.

---

# 9. DMVs are not permanent historical records

This is an important difference from Query Store.

Imagine you have:

```text
Monday:
Query A terrible

Tuesday:
Query A terrible

Wednesday:
SQL Server restarted
```

Some DMV statistics may be reset/lost after events such as a restart or plan cache changes.

So if you ask:

> "How was this query performing three weeks ago?"

DMVs may not be the right source.

That's where Query Store is much more useful.

---

# 10. The mental model

You can think of the tools like this:

```text
                    SQL Server troubleshooting

                         What happened?
                              │
                         Query Store
                              │
                 historical queries/plans
                              │
                              ↓
                    What is happening now?
                              │
                             DMVs
                              │
          ┌───────────────────┼──────────────────┐
          ↓                   ↓                  ↓
      Requests             Waits             Usage
          │                   │                  │
     running work        what it's waiting   indexes/resources
```

The important part isn't memorizing the diagram.

It's knowing **which question each tool answers**.

---

# 11. A realistic troubleshooting example

Suppose `/search` suddenly becomes slow.

You could investigate in this order:

### Step 1 — Query Store

Find:

```text
Search query
Average duration increased
```

### Step 2 — Execution plan

You discover:

```text
Large Index Scan
Millions of rows
```

### Step 3 — DMVs

You check whether the problem is currently happening.

Perhaps you discover that the current executions aren't actually CPU-heavy—they're waiting on locks.

Now your hypothesis changes.

### Step 4 — Investigate blocking

You find:

```text
Session 72
    ↓ blocks
Session 91
    ↓ blocks
Session 94
    ↓ blocks
...
```

Now you've discovered the real problem.

The query plan may be perfectly reasonable.

The system is slow because another transaction is holding locks for too long.

That distinction is **exactly** what separates database troubleshooting from merely knowing SQL syntax.

---

## The three DMV categories I'd remember

You don't need to memorize every DMV. Remember these:

| Question | Useful DMV |
|---|---|
| What queries are currently running? | `sys.dm_exec_requests` |
| Which cached queries have consumed CPU/reads/time? | `sys.dm_exec_query_stats` |
| How are indexes being used? | `sys.dm_db_index_usage_stats` |

And remember:

> **DMVs give you operational evidence; Query Store gives you historical evidence.**

---

### Next: Extended Events

That's the last major tool in this SQL Server troubleshooting section.

We'll cover **how you capture things that you specifically want to observe**, including deadlocks, long-running queries, errors, and other events—and why **Extended Events is preferred over the old SQL Server Profiler approach**.