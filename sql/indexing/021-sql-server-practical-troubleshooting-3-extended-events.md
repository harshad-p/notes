# SQL Server Practical Troubleshooting — Part 3: Extended Events

**Extended Events (XE)** is SQL Server's event-monitoring system.

Query Store tells you what happened historically, and DMVs tell you what's happening/currently accumulated. Extended Events lets you say:

> **"I want SQL Server to record these specific events when they happen."**

This is particularly useful when the problem is intermittent and difficult to reproduce.

---

## 1. Why do we need Extended Events?

Imagine someone reports:

> "Every few hours, one of our queries takes 30 seconds."

You check the query now:

```text
Average duration: 40 ms
Current execution: fine
```

Query Store may show that it happened, but you might want much more detail about **the specific occurrence**.

You can create an Extended Events session that watches for things such as:

- deadlocks,
- long-running queries,
- errors,
- excessive waits,
- specific SQL statements,
- specific databases or applications.

When the event occurs, SQL Server captures the information you've configured.

---

# 2. Think of it as an event listener

Suppose you care about deadlocks.

You configure:

```text id="n2w4yx"
Extended Events session

Watch for:
    deadlock

When it happens:
    capture details
```

Then:

```text id="z0o3k4"
Transaction A ──┐
                ├── deadlock
Transaction B ──┘
       ↓
Extended Events captures event
```

You can then inspect the captured event and determine what the transactions were doing.

This is why Extended Events is particularly useful for **deadlock investigation**.

---

# 3. What information can it capture?

Depending on the event, you can capture things such as:

- SQL statement,
- database,
- session ID,
- application,
- login,
- duration,
- CPU time,
- logical reads,
- error information,
- wait information,
- deadlock graph.

The important idea is that **you choose what you want to observe**.

You don't generally want to capture absolutely everything on a busy production server.

---

# 4. Extended Events vs Profiler

You may hear:

> "SQL Server Profiler."

Profiler is the older SQL Server tracing tool.

For modern SQL Server troubleshooting, **Extended Events is generally preferred**.

Why?

Because Extended Events was designed to provide detailed event monitoring with much lower overhead and more flexibility than the old tracing approach.

So if an interviewer asks:

> "How would you trace production SQL Server problems?"

A good answer is:

> "I'd use Extended Events rather than relying on the older SQL Server Profiler, and I'd configure a targeted session for the events I'm interested in."

You don't need to claim that Profiler is completely useless. The important point is that **XE is the modern mechanism**.

---

# 5. Don't capture everything

This is an important operational principle.

Imagine:

```text
Production server
    ↓
10,000 queries/sec
```

And you configure:

> "Capture every event and every piece of information."

You're generating a huge amount of diagnostic data and potentially adding unnecessary overhead.

Instead, make the session targeted.

For example:

> Capture deadlock events for the production database.

Or:

> Capture queries exceeding a particular duration.

Or:

> Capture errors of interest from a specific application.

The goal is:

**collect enough information to answer the question without turning diagnostics into another performance problem.**

---

# 6. Example: investigating deadlocks

Suppose your application occasionally reports:

```text
Transaction failed.
```

You suspect deadlocks.

You create an Extended Events session watching for deadlock events.

Later:

```text
Request A
    ↓
updates Orders
    ↓
waits for Customer

Request B
    ↓
updates Customer
    ↓
waits for Order

        DEADLOCK
```

XE captures the deadlock information.

You inspect it and discover:

```text
Transaction A
holds: Orders
waits: Customer

Transaction B
holds: Customer
waits: Orders
```

Now you know the exact circular dependency.

You can fix the application by making both code paths acquire resources in the same order.

That's much better than randomly changing indexes or increasing command timeouts.

---

# 7. Extended Events and slow queries

You can also configure XE to capture long-running statements.

For example, conceptually:

```text id="7czg9r"
Application
    ↓
SQL query
    ↓
takes > threshold
    ↓
Extended Events captures it
```

You could then investigate:

- which query ran,
- how long it took,
- which application/session executed it,
- what database it targeted,
- other event information relevant to the problem.

This is particularly useful for **intermittent problems**.

---

# 8. Query Store vs DMVs vs Extended Events

Now we can put the three together.

| Tool | Main question |
|---|---|
| **Query Store** | What happened historically? |
| **DMVs** | What is happening / what has SQL Server accumulated? |
| **Extended Events** | What specific events do I want to capture when they happen? |

A useful mental model:

```text
Query Store
    "Show me the history."

DMVs
    "Show me the current state / accumulated statistics."

Extended Events
    "Watch for this particular thing and record it."
```

They aren't competing tools. They complement each other.

---

# 9. Putting everything together

Suppose your production API has suddenly become slow.

You could approach it like this:

### First: establish whether SQL Server is actually the problem

Don't assume.

Compare:

```text
API latency
DB query latency
downstream service latency
```

If the API takes 5 seconds but SQL takes 50 ms, you have a different problem.

---

### Second: use Query Store

Ask:

> "Which queries have become expensive?"

Look at:

- duration,
- CPU,
- logical reads,
- execution count,
- plan changes.

---

### Third: inspect the execution plan

Look for evidence such as:

- unexpected scans,
- excessive lookups,
- bad cardinality estimates,
- expensive sorts/hashes,
- spills,
- parameter-sensitive behavior.

---

### Fourth: use DMVs

If the problem is happening now:

- what requests are running?
- what are they waiting for?
- is something blocking them?
- which session is responsible?
- what indexes are being used?

---

### Fifth: use Extended Events when necessary

If the problem is intermittent or requires detailed event-level evidence:

- capture deadlocks,
- capture long-running statements,
- capture relevant errors/waits,
- inspect the event details.

---

# 10. The complete troubleshooting mindset

You've now covered enough SQL Server performance material that I want you to shift your mindset slightly.

When somebody says:

> **"The database is slow."**

Don't immediately think:

> "I need an index."

Think:

```text
What is actually slow?
        ↓
Where is the time going?
        ↓
CPU?
I/O?
Blocking?
Deadlock?
Memory?
Bad plan?
Too many executions?
Application?
External dependency?
        ↓
What evidence can prove the cause?
        ↓
Make one targeted change
        ↓
Measure again
```

That's the real skill.

---

# SQL track completed

We've now covered the SQL curriculum we set out to do:

### Core SQL
- SQL querying
- joins
- grouping/aggregation
- subqueries
- filtering
- ordering
- relational concepts

### Indexing
- clustered/nonclustered indexes
- composite indexes
- included columns
- covering indexes
- filtered indexes
- index trade-offs
- when indexes hurt

### Performance
- execution plans
- scans vs seeks
- key lookups
- cardinality estimates
- SARGability
- parameter sniffing
- identifying the actual bottleneck

### Transactions
- atomicity
- transaction boundaries
- rollback
- savepoints
- isolation levels
- dirty/non-repeatable/phantom reads
- optimistic concurrency
- pessimistic concurrency
- lost updates
- locking
- blocking
- deadlocks
- retries

### Production troubleshooting
- Query Store
- DMVs
- Extended Events
- Profiler vs Extended Events
- diagnosing expensive queries
- diagnosing blocking/deadlocks
- performance investigation workflow

**SQL is now at a pretty solid senior-interview level.**

The next time SQL comes up in an interview, you should be able to go beyond *"add an index"* and actually reason about **plans, concurrency, waits, and evidence**.

Next we can move to the **next major backend topic** rather than doing another round of SQL repetition.