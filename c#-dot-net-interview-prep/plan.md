## Part 0 — Interview Foundations

- How to approach technical questions
- "Why X over Y?" questions
- Explaining trade-offs
- Diagnosing instead of guessing
- How to reason through unfamiliar code
- Coding-interview approach
- System-design approach

This should be short. It's about **how to think**, not C# syntax.

---

# Part 1 — C# Core

### 1.1 Types and Type System
- Value types vs reference types
- Stack vs heap — what this actually means in .NET
- `struct` vs `class`
- `record` vs `record class` vs `record struct`
- Boxing/unboxing
- Nullable value types
- Nullable reference types

### 1.2 Classes and Objects
- Constructors
- Properties
- Fields
- Access modifiers
- `static`
- `readonly`
- `const`
- `init`

### 1.3 Object-Oriented Programming
- Encapsulation
- Inheritance
- Polymorphism
- Abstraction
- Composition vs inheritance
- Virtual/override/new
- Abstract classes vs interfaces

### 1.4 Equality
- Reference equality
- Value equality
- `Equals`
- `GetHashCode`
- `==`
- Records and generated equality

---

# Part 2 — Collections

- `Array`
- `List<T>`
- `Dictionary<TKey,TValue>`
- `HashSet<T>`
- `Queue<T>`
- `Stack<T>`
- `LinkedList<T>`
- `IEnumerable<T>`
- `ICollection<T>`
- `IList<T>`
- `IReadOnlyCollection<T>`
- `IReadOnlyList<T>`
- Choosing the right collection
- Big-O characteristics
- Hash tables and buckets
- Dictionary collision handling

**Interview questions:**

- `Dictionary` vs `HashSet`
- `List` vs `LinkedList`
- Why does dictionary lookup usually take O(1)?
- What happens when a dictionary grows?
- When would you use `ConcurrentDictionary`?

---

# Part 3 — Generics

- Why generics exist
- Generic methods
- Generic classes
- Type safety
- Generic constraints
- `where T : class`
- `where T : struct`
- `where T : new()`
- Covariance
- Contravariance
- Invariance

---

# Part 4 — Delegates, Lambdas and Events

- Delegates
- `Action`
- `Func`
- `Predicate`
- Lambda expressions
- Closures
- Events
- Event handlers
- Delegates vs interfaces
- Common interview traps

---

# Part 5 — LINQ

Not just syntax.

### Fundamentals

- Deferred execution
- Immediate execution
- `IEnumerable`
- `IQueryable`
- `Select`
- `Where`
- `SelectMany`
- `OrderBy`
- `GroupBy`
- `Join`
- `Any`
- `All`
- `First`
- `Single`
- `FirstOrDefault`
- `ToList`
- `ToDictionary`

### Deeper

- LINQ execution pipeline
- Multiple enumeration
- LINQ performance
- LINQ-to-Objects vs LINQ-to-SQL
- Expression trees
- Why `IQueryable` behaves differently

This directly connects to the question you were asked:

> **IEnumerable vs IQueryable**

---

# Part 6 — Memory & Performance

- Garbage collection
- Generations
- Managed vs unmanaged memory
- Large Object Heap
- `IDisposable`
- `using`
- `IAsyncDisposable`
- Finalizers
- Object allocation
- String immutability
- `StringBuilder`
- `Span<T>`
- `Memory<T>`
- `ArrayPool<T>`
- Performance profiling
- Allocation vs CPU bottlenecks

---

# Part 7 — Exceptions

- Exception hierarchy
- `try/catch/finally`
- `throw`
- `throw ex`
- Inner exceptions
- Exception wrapping
- Custom exceptions
- When to catch
- When NOT to catch
- Exceptions vs result types
- Exceptions and performance
- Global exception handling

Especially:

```csharp
catch (Exception ex)
{
    throw ex;
}
```

vs

```csharp
catch (Exception ex)
{
    throw;
}
```

and why the difference matters.

---

# Part 8 — Async/Await

This gets a **deep treatment** because it came up repeatedly.

- What `async` actually does
- What `await` actually does
- State machines
- Continuations
- SynchronizationContext
- ThreadPool
- I/O-bound work
- CPU-bound work
- Why async can free a thread
- Why async doesn't create a thread
- `Task`
- `ValueTask`
- `Task<T>`
- `Task.WhenAll`
- `Task.WhenAny`
- Cancellation
- `CancellationToken`
- Async exception propagation
- Async deadlocks
- `ConfigureAwait`

And:

> Why `Task.Run(() => database.GetUser(id))` is generally not the solution when the database API doesn't provide async I/O.

---

# Part 9 — Threads, Tasks & Parallelism

- `Thread`
- `Task`
- ThreadPool
- `Task.Run`
- Long-running work
- CPU-bound parallelism
- `Parallel`
- `Parallel.ForEach`
- `Parallel.ForEachAsync`
- Thread creation cost
- Thread starvation
- Synchronization

Key comparisons:

- `Thread` vs `Task`
- `Task.Run` vs `await`
- Parallelism vs concurrency
- Async vs parallelism

---

# Part 10 — Concurrent Programming

This deserves its own part.

- Race conditions
- Atomicity
- `lock`
- `Monitor`
- `Interlocked`
- `volatile`
- `SemaphoreSlim`
- `Mutex`
- `ReaderWriterLockSlim`
- Concurrent collections
- `ConcurrentDictionary`
- `ConcurrentQueue`
- `ConcurrentBag`
- Thread-safe caching
- Immutable data
- Read/write patterns
- Deadlocks
- Lock contention

And the exact reasoning problems we practiced:

> "Two related values must be updated atomically. What do you use?"

> "Multiple threads access this cache. How do you make it safe?"

---

# Part 11 — Dependency Injection

- Why DI exists
- Dependency inversion
- Constructor injection
- Method injection
- Property injection
- Service lifetimes
- Transient
- Scoped
- Singleton
- Service locator anti-pattern
- Captive dependencies
- Factory registration
- `IServiceProvider`
- DI scopes

Also distinguish:

> **Dependency injection types**

from:

> **DI service lifetimes**

This was something you specifically asked about earlier.

---

# Part 12 — .NET Runtime

- CLR
- .NET runtime
- JIT
- IL
- Assemblies
- Metadata
- Reflection
- Runtime type information
- Garbage collector
- ThreadPool
- Native interop
- Application startup

Useful question:

> What actually happens when I execute a C# program?

---

# Part 13 — ASP.NET Core

- Request lifecycle
- Controllers
- Minimal APIs
- Routing
- Model binding
- Validation
- Dependency injection
- Configuration
- Logging
- Middleware
- Filters
- Authentication
- Authorization
- Error handling

---

# Part 14 — HTTP & APIs

- HTTP methods
- Status codes
- Headers
- Cookies
- Authentication
- Authorization
- JWT
- Access tokens
- Refresh tokens
- OAuth 2.0
- OpenID Connect
- Idempotency
- Pagination
- Versioning
- Rate limiting
- Retries
- Timeouts
- REST principles

Even though the Immowelt JD didn't explicitly mention REST, this is still useful backend knowledge.

---

# Part 15 — Middleware & Pipelines

- Middleware pipeline
- Request/response pipeline
- Middleware ordering
- Short-circuiting
- Exception middleware
- Authentication middleware
- Authorization middleware
- Logging middleware
- Custom middleware
- Filters vs middleware

Also:

> What happens to a request as it travels through ASP.NET Core?

---

# Part 16 — EF Core & Data Access

- `DbContext`
- DbSet
- Tracking
- No-tracking
- Change tracking
- Relationships
- Loading strategies
- Eager loading
- Explicit loading
- Lazy loading
- N+1
- Transactions
- Concurrency
- Optimistic concurrency
- `SaveChanges`
- `SaveChangesAsync`
- Connection pooling
- Compiled queries
- EF Core performance

---

# Part 17 — SQL

This should be substantial.

### SQL fundamentals

- SELECT
- JOIN
- GROUP BY
- HAVING
- Subqueries
- CTEs
- Window functions

### Data modelling

- Normalization
- Denormalization
- Constraints
- Primary keys
- Foreign keys

### Indexing

- Clustered indexes
- Non-clustered indexes
- Composite indexes
- Included columns
- Covering indexes
- Filtered indexes
- Selectivity
- Index trade-offs

### Performance

- Execution plans
- SARGability
- Parameter sniffing
- Statistics
- Query Store
- DMVs
- Locking
- Blocking
- Deadlocks
- I/O vs CPU
- Query optimization

### Transactions

- ACID
- Isolation levels
- Atomicity
- Rollback
- Savepoints
- Concurrency
- EF Core transactions

---

# Part 18 — Messaging

- Queues
- Pub/sub
- RabbitMQ
- Kafka
- Azure Service Bus
- Delivery semantics
- At-most-once
- At-least-once
- Exactly-once considerations
- Message ordering
- Retries
- Dead-letter queues
- Idempotency
- Duplicate messages
- Backpressure

Key interview comparison:

> **RabbitMQ vs Kafka**

---

# Part 19 — Caching

- Why caching exists
- Cache-aside
- Read-through
- Write-through
- Write-behind
- Cache invalidation
- TTL
- Cache stampede
- Cache penetration
- Cache eviction
- Local vs distributed cache
- Redis
- Distributed locking

---

# Part 20 — Testing & Debugging

We already started this area, but I'd make the course more systematic.

### Testing

- Unit tests
- Integration tests
- End-to-end tests
- Test doubles
- Mocking
- Stubbing
- Dependency injection for testing
- Test isolation
- Async tests
- Database integration tests

### Debugging

- Reading stack traces
- Breakpoints
- Conditional breakpoints
- Watches
- Call stack
- Exception settings
- Debugging async code
- Debugging concurrency
- Debugging production issues
- Performance profiling

---

# Part 21 — Design Patterns

Don't turn this into "memorize 23 GoF patterns."

Focus on patterns you're actually likely to encounter.

### Creational

- Factory
- Factory Method
- Builder

### Structural

- Adapter
- Decorator
- Facade

### Behavioral

- Strategy
- Observer
- Chain of Responsibility

### .NET-specific/common patterns

- Repository
- Unit of Work
- Options pattern
- Middleware
- Dependency Injection

For each:

**Problem → pattern → implementation → trade-offs → when NOT to use it.**

---

# Part 22 — Distributed Systems

This is where the backend knowledge starts coming together.

- Horizontal scaling
- Load balancing
- Stateless services
- Replication
- Partitioning
- Sharding
- Consistency
- Availability
- CAP theorem
- Eventual consistency
- Distributed transactions
- Idempotency
- Retries
- Timeouts
- Circuit breakers
- Bulkheads
- Backpressure
- Rate limiting

---

# Part 23 — System Design

Your existing system-design framework fits here.

Problems such as:

1. Property Search & Listing Platform
2. Property Alert / Notification System
3. Listing Moderation Pipeline
4. URL Shortener
5. Payment Processing System
6. Notification System
7. File Upload/Processing System
8. Job Processing Platform

For each:

- Requirements
- Data model
- Database choice
- API
- Initial architecture
- Scaling
- Caching
- Queues
- Reliability
- Observability
- Trade-offs
- Follow-up questions

---

# Part 24 — Interview Practice

This is where everything comes together.

### "Why X over Y?"

Examples:

- `IEnumerable` vs `IQueryable`
- `Task` vs `Thread`
- `Task.Run` vs async I/O
- `lock` vs `Interlocked`
- `Dictionary` vs `ConcurrentDictionary`
- SQL vs NoSQL
- RabbitMQ vs Kafka
- Redis vs database
- SQL transaction vs distributed transaction
- interface vs abstract class
- class vs record
- `List` vs `HashSet`

### Coding

- Collections
- Grouping
- Joins
- Algorithms
- String manipulation
- Data transformation
- Concurrency
- Async
- Debugging existing code
- SQL coding

### Backend scenarios

- Slow API
- Database bottleneck
- Queue backlog
- Duplicate messages
- Race condition
- Memory leak
- Thread starvation
- Production outage

### Behavioral technical questions

- Tell me about a performance improvement
- Tell me about a failure
- Tell me about a difficult production issue
- Tell me about a technical disagreement
- Tell me about a system you owned end-to-end

---

## One thing I'd deliberately add to this course

A **"Senior Question Bank"** inside `24-interview-practice/`.

Not just questions with answers, but three levels:

```text
Question
    ↓
Your answer
    ↓
Critique
    ↓
Strong answer
    ↓
Deeper follow-up
```

That's exactly the format we've been using when I've asked you questions like:

> "Why would you use `lock` instead of `Interlocked`?"

and you answer first.

That will be much more valuable for you than simply reading prepared answers.

### Overall progression

```text
C# fundamentals
      ↓
Collections / LINQ
      ↓
Memory / exceptions
      ↓
Async / concurrency
      ↓
.NET / ASP.NET
      ↓
EF Core / SQL
      ↓
Messaging / caching
      ↓
Testing / debugging
      ↓
Design patterns
      ↓
Distributed systems
      ↓
System design
      ↓
Senior interview practice
```