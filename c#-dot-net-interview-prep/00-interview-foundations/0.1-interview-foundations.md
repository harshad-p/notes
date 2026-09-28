# Part 0 — Interview Foundations

This part is intentionally short. The goal is to build a **repeatable way of thinking** that you can apply to C#, .NET, SQL, backend, and system-design questions.

You don't need to memorize scripts. You want to develop a habit of **understanding the problem first, then reasoning toward an answer**.

---

## 0.1 How to Approach Technical Questions

A technical interview question often looks like it is testing whether you know a fact:

> “What's the difference between `IEnumerable` and `IQueryable`?”

But a senior-level interviewer is often testing something deeper:

- Do you understand what each abstraction represents?
- Do you know when the difference matters?
- Can you explain the trade-off?
- Can you apply it to a real situation?
- Can you recognize when your answer depends on context?

A useful mental process is:

### 1. Understand what is actually being asked

Don't immediately start explaining everything you know about the topic.

For example:

> “What's the difference between a thread and a task?”

The interviewer may not be asking for definitions. They may really be asking:

> **Why would I choose one over the other?**

So first identify the type of question.

Common categories are:

| Question type | What you're being asked to demonstrate |
|---|---|
| **What is X?** | Understanding |
| **How does X work?** | Mechanism |
| **Why use X?** | Purpose |
| **X vs Y?** | Differences + trade-offs |
| **Why X over Y?** | Decision-making |
| **What happens if...?** | Reasoning |
| **How would you improve...?** | Diagnosis + design |
| **How would you design...?** | Systematic problem solving |

---

### 2. Start from the simplest correct explanation

Don't start with the most advanced thing you know.

For example, if asked:

> “Why use a cache?”

Start with:

> “A cache stores frequently accessed data closer to the application so we can avoid repeatedly going to the database.”

Then go deeper if necessary:

> “The trade-off is that cached data can become stale, so we need to think about expiration and invalidation.”

Then deeper again if the interviewer asks:

> “What happens if multiple application servers are using the cache?”

Now you're discussing distributed caching.

The principle is:

**Start simple → establish correctness → add depth when needed.**

---

### 3. Explain the mechanism, not just the result

Compare these two answers:

> “`IQueryable` is better for database queries.”

versus:

> “`IQueryable` represents a query that can be translated by the provider, such as EF Core, and executed by the database. That allows filtering to happen in the database rather than retrieving all the data and filtering it in memory.”

The second answer demonstrates understanding because you can explain **why** the behavior occurs.

This is particularly important for senior interviews.

---

### 4. Connect the concept to consequences

After explaining how something works, ask yourself:

> **“So what?”**

For example:

`LOWER(Name) = 'harshad'`

The important thing isn't merely that this is “non-SARGable.”

The reasoning is:

1. SQL Server has an index on `Name`.
2. The query applies a function to the column.
3. SQL Server may not be able to use the index efficiently for the search.
4. It may need to examine many rows.
5. That can become expensive as the table grows.

Now you've demonstrated understanding rather than terminology.

---

## 0.2 "Why X Over Y?" Questions

This is one of the most important patterns for senior interviews.

Examples:

- Why `Task` instead of `Thread`?
- Why `IQueryable` instead of `IEnumerable`?
- Why SQL instead of NoSQL?
- Why `lock` instead of `Interlocked`?
- Why a queue instead of calling the service directly?
- Why Redis instead of querying the database?
- Why `async/await` instead of `Task.Run()`?
- Why a distributed cache instead of an in-memory cache?

A weak answer usually sounds like:

> “X is better than Y because X is faster.”

A strong answer follows this reasoning:

### 1. Establish the purpose of both

Don't assume Y is bad.

> “Both approaches can solve the problem, but they are designed for different situations.”

### 2. Identify the important difference

Ask:

> **What property actually matters for this decision?**

For example:

`Task` vs `Thread` → abstraction and resource management.

SQL vs NoSQL → data relationships, consistency, query requirements, scaling model.

`lock` vs `Interlocked` → complexity of the operation being protected.

### 3. Explain the trade-off

Every serious engineering decision has a cost.

For example:

> “A distributed cache can reduce database load and latency, but now we have another distributed component and have to deal with stale data, cache failures and invalidation.”

### 4. State when you would choose each

This is the part that makes the answer senior-level.

Don't say:

> “X is better.”

Say:

> “I'd choose X when ..., while Y makes more sense when ...”

That shows that you understand **context**, rather than memorizing a preferred technology.

---

## 0.3 Explaining Trade-offs

A trade-off means:

> **You gain something, but you give something else up.**

For example, adding a cache:

**Gain:**
- lower latency
- fewer database queries
- potentially higher throughput

**Give up:**
- data can become stale
- cache invalidation becomes a problem
- another component can fail
- memory is required
- application behavior becomes more complex

A useful way to reason about any technical choice is:

### What do I gain?

Performance?  
Scalability?  
Reliability?  
Simplicity?  
Consistency?  
Development speed?

### What do I give up?

Complexity?  
Cost?  
Consistency?  
Operational effort?  
Maintainability?

### What assumption am I making?

For example:

> “I'm assuming the data doesn't need to be immediately consistent.”

That last question is particularly useful.

Many architectural decisions are really about **which assumption you're willing to make**.

---

## 0.4 Diagnose Instead of Guessing

This is particularly important for backend and performance questions.

Suppose an interviewer says:

> “This API became slow. How would you fix it?”

A weak approach is:

> “I'd add Redis.”

Or:

> “I'd add an index.”

Or:

> “I'd add more servers.”

Those are solutions looking for a problem.

A better approach is:

> **First determine where the time is actually being spent.**

For example:

1. Is the application CPU-bound?
2. Is it waiting on the database?
3. Is the database doing excessive I/O?
4. Is the query using a bad execution plan?
5. Is the request blocked by another transaction?
6. Is a downstream service slow?
7. Is there a network problem?
8. Is the application thread pool saturated?
9. Did traffic increase?
10. Did the amount of data increase?

Only after identifying the bottleneck should you choose the solution.

This principle applies everywhere.

### Don't do this:

**Symptom → favorite technology → solution**

### Do this:

**Symptom → evidence → bottleneck → cause → solution → verify**

For example:

> API latency increased  
> → check metrics  
> → database accounts for 80% of request time  
> → inspect query  
> → execution plan shows expensive scan  
> → investigate predicate/index  
> → change query/index  
> → measure again

That is **engineering reasoning**.

---

## 0.5 How to Reason Through Unfamiliar Code

You will sometimes get code you've never seen before.

Don't try to understand every line immediately.

Start from the outside.

### Step 1 — What is the code supposed to do?

Look at:

- method/class name
- parameters
- return value
- surrounding context

Ask:

> “What problem is this code solving?”

### Step 2 — Identify the main flow

Follow the data:

> Input → processing → external calls/database → output

Ignore implementation details initially.

### Step 3 — Find important state

Look for:

- mutable variables
- shared state
- collections
- database entities
- cached values
- locks
- tasks
- configuration

These often determine the behavior.

### Step 4 — Look for boundaries

Pay particular attention when code crosses boundaries:

- application → database
- application → HTTP service
- thread → shared memory
- request → background job
- process → process
- transaction → external system

These are where many bugs and performance problems appear.

### Step 5 — Ask what could go wrong

Examples:

- null values
- exceptions
- concurrency
- duplicate processing
- timeouts
- retries
- partial failure
- large datasets
- resource leaks
- race conditions

You don't need to find a problem immediately.

First understand the **normal flow**.

---

## 0.6 Coding-Interview Approach

For coding problems, don't immediately start typing.

Use this sequence:

### 1. Clarify the requirement

Understand:

- input
- output
- constraints
- edge cases
- expected behavior

If something is ambiguous, ask.

### 2. Give a simple approach

Explain what you're going to do before writing code.

For example:

> “I'll iterate through the listings, keep the active ones, calculate the aggregate per category, and then sort the result.”

### 3. Implement the straightforward solution

Get the logic correct first.

Don't prematurely optimize.

### 4. Test mentally with edge cases

Consider:

- empty input
- one element
- duplicates
- missing relationships
- nulls
- no matching records
- very large input

### 5. Analyze complexity

Ask:

> “How much work does this do as the input grows?”

For example:

`O(n)`, `O(n log n)`, `O(n²)`.

But don't just state Big-O. Understand **why**.

### 6. Then optimize

If the interviewer asks:

> “How could you improve this?”

Now look for:

- unnecessary loops
- repeated work
- inappropriate data structures
- database round trips
- unnecessary allocations
- sorting
- memory usage

And explain the trade-off of your optimization.

---

## 0.7 System-Design Approach

System design is essentially the same reasoning process at a larger scale.

Start with:

> **What does the system need to do?**

Then:

1. Clarify functional requirements.
2. Clarify non-functional requirements.
3. Determine what data needs to be stored.
4. Decide how that data needs to be accessed.
5. Choose an appropriate storage model.
6. Define APIs.
7. Draw the simplest architecture that satisfies the requirements.
8. Identify bottlenecks and failure points.
9. Add components only when there's a reason for them.
10. Explain the trade-offs.

The important principle is:

> **Don't design for hypothetical scale before understanding the actual requirements.**

If a simple API → database architecture satisfies the requirements, start there.

Then the interviewer might say:

> “What if traffic becomes 100× larger?”

Now you have a reason to introduce things such as:

- load balancing
- additional application instances
- caching
- read replicas
- queues
- partitioning
- distributed systems

Each component should answer a specific problem.

---

# The Core Mental Model

If you remember only one thing from Part 0, make it this:

> **Understand → reason → choose → explain → verify.**

When faced with a technical question:

**Understand the problem**  
↓  
**Identify what actually matters**  
↓  
**Consider the available options**  
↓  
**Compare their trade-offs**  
↓  
**Choose based on the requirements**  
↓  
**Verify that the solution actually solves the problem**

That's the mindset we'll carry into the rest of the course.

**Part 0 is complete.** The next part is **Part 1 — C# Core**, starting with the C# type system.