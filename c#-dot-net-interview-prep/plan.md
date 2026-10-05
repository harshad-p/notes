# Curriculum

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

## Part 1 — C# Core

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

## Part 2 — Collections & Data Structures

This module covers .NET collections first, then the underlying data structures they are based on, their interfaces, complexity, and how to choose the right structure for a problem.

### 2.1 `Array`

- Fixed-size arrays
- Zero-based indexing
- `Length`
- Single-dimensional arrays
- Multidimensional arrays
- Jagged arrays
- Array initialization
- Array copying
- Array resizing
- Memory layout
- Performance characteristics
- Interfaces implemented by arrays
- When to use an array
- Advantages and disadvantages

---

### 2.2 `List<T>`

- What `List<T>` is
- Dynamic-array implementation
- Indexing
- `Count`
- `Capacity`
- Initial capacity
- `EnsureCapacity`
- Capacity growth
- Resizing and copying
- `Add`
- `AddRange`
- `Insert`
- `InsertRange`
- `Remove`
- `RemoveAt`
- `RemoveRange`
- `Contains`
- `IndexOf`
- `Clear`
- `ToArray`
- Interfaces implemented by `List<T>`
- Internal implementation
- Time complexity
- Memory considerations
- When to use `List<T>`
- Advantages and disadvantages

---

### 2.3 `Dictionary<TKey,TValue>`

- What `Dictionary<TKey,TValue>` is
- Key/value relationship
- Key uniqueness
- `Count`
- `Capacity`
- Initial capacity
- `EnsureCapacity`
- Resizing
- `Add`
- Indexer
- `TryAdd`
- `TryGetValue`
- `ContainsKey`
- `ContainsValue`
- `Remove`
- `Clear`
- Hashing
- Hash codes
- Buckets
- Entries
- Collisions
- Collision handling
- Equality and `GetHashCode`
- Comparers
- Internal implementation
- Time complexity
- Memory considerations
- Interfaces implemented by `Dictionary<TKey,TValue>`
- When to use `Dictionary<TKey,TValue>`
- Advantages and disadvantages

---

### 2.4 `HashSet<T>`

- What `HashSet<T>` is
- Set semantics
- Uniqueness
- `Count`
- `Capacity`
- Initial capacity
- `EnsureCapacity`
- Resizing
- `Add`
- `Remove`
- `Contains`
- `Clear`
- `UnionWith`
- `IntersectWith`
- `ExceptWith`
- `SymmetricExceptWith`
- `IsSubsetOf`
- `IsSupersetOf`
- `SetEquals`
- Hashing
- Equality
- `GetHashCode`
- Comparers
- Internal implementation
- Buckets and entries
- Time complexity
- Memory considerations
- Interfaces implemented by `HashSet<T>`
- When to use `HashSet<T>`
- Advantages and disadvantages

---

### 2.5 `Queue<T>`

- FIFO
- Creating and initializing a queue
- `Count`
- `Capacity`
- Initial capacity
- `EnsureCapacity`
- Capacity growth
- `Enqueue`
- `Dequeue`
- `TryDequeue`
- `Peek`
- `TryPeek`
- `Contains`
- `Clear`
- Internal array
- Circular-buffer implementation
- Resizing
- Amortized complexity
- Time complexity
- Memory considerations
- Interfaces implemented by `Queue<T>`
- `Queue<T>` vs `List<T>`
- `Queue<T>` vs `Stack<T>`
- Practical use cases
- Advantages and disadvantages

---

### 2.6 `Stack<T>`

- LIFO
- Creating and initializing a stack
- `Count`
- `Capacity`
- Initial capacity
- `EnsureCapacity`
- Capacity growth
- `Push`
- `Pop`
- `TryPop`
- `Peek`
- `TryPeek`
- `Contains`
- `Clear`
- Internal array
- Resizing
- Amortized complexity
- Time complexity
- Memory considerations
- Interfaces implemented by `Stack<T>`
- `Stack<T>` vs `Queue<T>`
- Stack and recursion
- Practical use cases
- Advantages and disadvantages

---

### 2.7 `PriorityQueue<TElement,TPriority>`

- What a priority queue is
- Priority vs FIFO ordering
- Creating and initializing a priority queue
- `Count`
- Capacity
- Initial capacity
- `Enqueue`
- `Dequeue`
- `TryDequeue`
- `Peek`
- `TryPeek`
- `EnqueueDequeue`
- `DequeueEnqueue`
- `Clear`
- Priority semantics
- Min-priority behavior
- What happens when multiple elements have the same priority
- Priority comparers
- Custom priority types
- Internal implementation
- Binary heap
- Heap representation using an array
- Heap insertion
- Heap removal
- Heapify
- Resizing
- Time complexity
- Memory considerations
- Interfaces implemented by `PriorityQueue<TElement,TPriority>`
- `PriorityQueue` vs `Queue`
- `PriorityQueue` vs sorting
- Practical use cases
- Scheduling
- Shortest-path algorithms
- Best-first search
- Advantages and disadvantages

---

### 2.8 `LinkedList<T>`

- What a linked list is
- Singly linked vs doubly linked lists
- `LinkedList<T>` implementation
- `LinkedListNode<T>`
- Head
- Tail
- `First`
- `Last`
- `Count`
- `AddFirst`
- `AddLast`
- `AddBefore`
- `AddAfter`
- `Remove`
- `RemoveFirst`
- `RemoveLast`
- `Find`
- `FindLast`
- `Contains`
- `Clear`
- Node references
- Why lookup is O(n)
- Why insertion/removal can be O(1) when the node is known
- Memory overhead
- Interfaces implemented by `LinkedList<T>`
- `LinkedList<T>` vs `List<T>`
- When a linked list is actually useful
- Advantages and disadvantages

---

### 2.9 `SortedSet<T>`

- What `SortedSet<T>` is
- Sorted unique values
- Ordering
- `Count`
- `Add`
- `Remove`
- `Contains`
- Set operations
- Range operations
- Minimum and maximum values
- Comparers
- Equality vs ordering
- Internal tree-based implementation
- Balanced tree concept
- Time complexity
- Memory considerations
- Interfaces implemented by `SortedSet<T>`
- `SortedSet<T>` vs `HashSet<T>`
- When to use `SortedSet<T>`
- Advantages and disadvantages

---

### 2.10 `SortedDictionary<TKey,TValue>`

- What `SortedDictionary<TKey,TValue>` is
- Sorted keys
- Key/value relationship
- `Count`
- `Add`
- Indexer
- `TryAdd`
- `TryGetValue`
- `ContainsKey`
- `Remove`
- `Clear`
- Key ordering
- Comparers
- Internal tree-based implementation
- Balanced tree concept
- Time complexity
- Memory considerations
- Interfaces implemented by `SortedDictionary<TKey,TValue>`
- `SortedDictionary<TKey,TValue>` vs `Dictionary<TKey,TValue>`
- When to use `SortedDictionary<TKey,TValue>`
- Advantages and disadvantages

---

### 2.11 `SortedList<TKey,TValue>`

- What `SortedList<TKey,TValue>` is
- Sorted keys
- `Count`
- `Capacity`
- Initial capacity
- Capacity growth
- Index-based access
- `Add`
- Indexer
- `TryGetValue`
- `ContainsKey`
- `Remove`
- `RemoveAt`
- `Clear`
- Binary search
- Internal array-based implementation
- Insertion and removal costs
- Time complexity
- Memory considerations
- Interfaces implemented by `SortedList<TKey,TValue>`
- `SortedList` vs `SortedDictionary`
- `SortedList` vs `Dictionary`
- When to use `SortedList`
- Advantages and disadvantages

---

### 2.12 Deque — Double-Ended Queue

This is primarily a data-structure concept because .NET does not provide a general-purpose `Deque<T>` collection equivalent to `Queue<T>`.

- What a deque is
- Double-ended insertion
- Double-ended removal
- Front
- Back
- FIFO and LIFO behavior as special cases
- Array-based implementations
- Circular-buffer implementations
- Linked implementations
- Time complexity
- `Deque` vs `Queue`
- `Deque` vs `Stack`
- Sliding-window problems
- BFS variants
- Monotonic queues
- Scheduling use cases
- Advantages and disadvantages

---

Underlying Data Structures

### 2.13 Heap

- What a heap is
- Complete binary tree
- Min-heap
- Max-heap
- Heap property
- Parent/child relationships
- Array representation
- Parent index calculation
- Child index calculation
- Insert
- Extract minimum
- Extract maximum
- Sift up
- Sift down
- Heapify
- Building a heap
- Heap sort
- Time complexity
- Space complexity
- Heap vs sorted collection
- Heap vs `PriorityQueue<TElement,TPriority>`
- How `PriorityQueue<TElement,TPriority>` uses a heap

---

### 2.14 Binary Tree

- What a tree is
- Root
- Node
- Parent
- Child
- Sibling
- Leaf
- Edge
- Depth
- Height
- Subtree
- Binary tree
- Complete binary tree
- Full binary tree
- Perfect binary tree
- Balanced vs unbalanced trees
- Tree representation
- Recursive representation
- Array representation
- Tree traversal
- Preorder traversal
- Inorder traversal
- Postorder traversal
- Level-order traversal
- Breadth-first traversal
- Depth-first traversal
- Time and space complexity
- Practical applications

Example:

```text
        A
       / \
      B   C
     / \
    D   E
```

---

### 2.15 Binary Search Tree

- What a BST is
- BST ordering property
- Search
- Insert
- Delete
- Minimum value
- Maximum value
- Successor
- Predecessor
- Inorder traversal
- Average-case complexity
- Worst-case complexity
- Degenerate BST
- Balanced vs unbalanced BST
- Why a BST can degrade to O(n)
- Practical use cases
- BST vs hash table
- BST vs sorted collections

---

### 2.16 Balanced Trees

- Why tree balancing is necessary
- Balanced tree concept
- AVL trees
- AVL rotations
- Left rotation
- Right rotation
- Red-black trees
- Red-black tree properties
- Why red-black trees are useful
- Complexity
- Relationship to .NET sorted collections
- `SortedSet<T>` and balanced trees
- `SortedDictionary<TKey,TValue>` and balanced trees
- When balanced trees are preferable to hash tables

We don't need to implement a production-quality red-black tree from scratch, but we should understand the concepts well enough to explain the trade-offs in an interview.

---

Other Core Data Structures

### 2.17 Graph

- What a graph is
- Vertices/nodes
- Edges
- Directed graphs
- Undirected graphs
- Weighted graphs
- Unweighted graphs
- Cyclic graphs
- Acyclic graphs
- Connected graphs
- Adjacency list
- Adjacency matrix
- Edge lists
- BFS
- DFS
- Visited tracking
- Connected components
- Cycle detection
- Graph traversal complexity
- Practical applications
- Graph vs tree

Example:

```text
A ----- B
|       |
|       |
C ----- D
```

---

### 2.18 Trie

- What a trie is
- Prefix tree
- Nodes
- Root
- Character paths
- Word termination
- Insert
- Search
- Prefix search
- Delete
- Time complexity
- Space complexity
- Autocomplete
- Dictionary lookup
- Prefix matching
- Trie vs `Dictionary`
- Practical applications
- Advantages and disadvantages

Example:

```text
        root
       /    \
      c      d
      |
      a
      |
      t
```

---

### 2.19 Disjoint Set / Union-Find

- What Disjoint Set is
- Sets and components
- `Find`
- `Union`
- Parent representation
- Root
- Path compression
- Union by rank
- Union by size
- Complexity
- Connected components
- Cycle detection
- Kruskal's algorithm
- Practical applications
- Advantages and disadvantages

---

### 2.20 LRU Cache

- What an LRU cache is
- Least Recently Used semantics
- `Get`
- `Put`
- Eviction
- Cache capacity
- Why a dictionary alone isn't sufficient
- Dictionary + doubly linked list
- Maintaining recency
- Moving recently accessed items
- Evicting the least recently used item
- Achieving O(1) `Get`
- Achieving O(1) `Put`
- Memory considerations
- Practical backend use cases
- LRU vs other cache eviction policies
- Implementation exercise

Conceptually:

```text
Dictionary
    |
    +---- key -> node
                 |
                 v
        Doubly Linked List

Most Recent <-> ... <-> Least Recent
```

---

Collection Interfaces

### 2.21 `IEnumerable<T>`

- What `IEnumerable<T>` represents
- Generic vs non-generic `IEnumerable`
- `GetEnumerator`
- `IEnumerator<T>`
- `foreach`
- Enumeration
- Deferred execution concept
- What operations the interface exposes
- Why collections implement it
- `IEnumerable<T>` vs concrete collections

---

### 2.22 `ICollection<T>`

- What `ICollection<T>` represents
- `Count`
- `IsReadOnly`
- `Add`
- `Remove`
- `Clear`
- `Contains`
- `CopyTo`
- Relationship with `IEnumerable<T>`
- Why different collections expose different capabilities through the interface

---

### 2.23 `IList<T>`

- What `IList<T>` represents
- Index-based access
- `Count`
- `IsReadOnly`
- `Add`
- `Insert`
- `Remove`
- `RemoveAt`
- Indexer
- Relationship with `ICollection<T>`
- `IList<T>` vs `ICollection<T>`
- Why `List<T>` implements it
- Limitations of the abstraction

---

### 2.24 `ISet<T>`

- What `ISet<T>` represents
- Set semantics
- Uniqueness
- `Add`
- Set operations
- Subset/superset relationships
- Set equality
- Relationship with `ICollection<T>`
- `HashSet<T>` implementation
- `SortedSet<T>` implementation
- `ISet<T>` vs concrete set implementations

---

### 2.25 `IReadOnlyCollection<T>`

- Read-only collection abstraction
- `Count`
- Relationship with `IEnumerable<T>`
- Difference between read-only API and immutable data
- Why returning `IReadOnlyCollection<T>` can be useful
- Implementations
- When to use it

---

### 2.26 `IReadOnlyList<T>`

- Read-only list abstraction
- Index-based access
- `Count`
- Indexer
- Relationship with `IReadOnlyCollection<T>`
- `IReadOnlyList<T>` vs `IList<T>`
- Read-only view vs immutable collection
- When to use it

---

Cross-Cutting Collection Concepts

### 2.27 Choosing the Right Collection

Decision-making based on requirements:

- Fast indexed access
- Fast key lookup
- Uniqueness
- FIFO processing
- LIFO processing
- Priority-based processing
- Sorted unique values
- Sorted key/value pairs
- Frequent insertion/removal
- Sequential iteration
- Read-only APIs
- Memory constraints
- Known collection size
- Unknown collection size
- Ordering requirements
- Equality requirements

Comparison of:

```text
Array
List<T>
Dictionary<TKey,TValue>
HashSet<T>
Queue<T>
Stack<T>
PriorityQueue<TElement,TPriority>
LinkedList<T>
SortedSet<T>
SortedDictionary<TKey,TValue>
SortedList<TKey,TValue>
```

---

### 2.28 Big-O Characteristics

- What Big-O means
- Time complexity
- Space complexity
- Best case
- Average case
- Worst case
- Amortized complexity
- O(1)
- O(log n)
- O(n)
- O(n log n)
- O(n²)

Create a consolidated complexity table for the major .NET collections.

---

### 2.29 Capacity, Resizing & Memory

- `Count` vs `Capacity`
- Initial capacity
- Why capacity exists
- Capacity growth
- Resizing
- Array allocation
- Copying elements
- Amortized operations
- Memory overhead
- Preallocating capacity
- `EnsureCapacity`
- When specifying capacity helps
- When capacity does not matter
- Which collections have meaningful capacity concepts

---

### 2.30 Hash Tables & Buckets

- Hash tables
- Hash functions
- Hash codes
- Buckets
- Entries
- Mapping keys to buckets
- Load factor concept
- Collisions
- Resizing
- Rehashing
- Average O(1) lookup
- Worst-case behavior
- Memory trade-offs
- Relationship to `Dictionary<TKey,TValue>`
- Relationship to `HashSet<T>`

---

### 2.31 Dictionary Collision Handling

- Why collisions happen
- Multiple keys mapping to the same bucket
- Collision resolution
- Chaining concept
- .NET's implementation approach
- Equality checks after hash matching
- `Equals`
- `GetHashCode`
- Comparers
- What happens when a dictionary grows
- Why good hash distribution matters
- Performance implications

---

### 2.32 Enumerators

- `IEnumerable`
- `IEnumerable<T>`
- `IEnumerator`
- `IEnumerator<T>`
- `GetEnumerator`
- `MoveNext`
- `Current`
- `Reset`
- `foreach`
- How `foreach` works internally
- Enumeration vs indexing
- Enumeration of different collection types
- Enumerator invalidation
- Modifying a collection while enumerating
- Deferred enumeration concepts

---

### 2.33 Sorting & Comparers

- `IComparable<T>`
- `IComparable`
- `IComparer<T>`
- `Comparer<T>`
- `Comparison<T>`
- Natural ordering
- Custom ordering
- `Sort`
- Sorting collections
- Comparers for dictionaries and sets
- Ordering vs equality
- `SortedSet<T>`
- `SortedDictionary<TKey,TValue>`
- `SortedList<TKey,TValue>`
- Stable vs unstable sorting
- Practical sorting scenarios

---

### 2.34 Equality & Collections

- Reference equality
- Value equality
- `Equals`
- `GetHashCode`
- `IEquatable<T>`
- `==`
- Equality comparers
- `IEqualityComparer<T>`
- Hash-based collections
- `Dictionary<TKey,TValue>`
- `HashSet<T>`
- Sorted collections
- Equality vs ordering
- Why overriding `Equals` and `GetHashCode` together matters
- Common collection equality bugs
- Custom comparers

---

Applying Collections & Data Structures

### 2.35 Scenarios for Using Each Collection

Work through practical scenarios such as:

- Maintain a list of users
- Fast lookup by ID
- Remove duplicates
- Process jobs in arrival order
- Undo operations
- Schedule work by priority
- Maintain sorted unique values
- Maintain sorted key/value data
- Frequently insert/remove known nodes
- Prefix search
- Autocomplete
- Graph traversal
- Shortest-path exploration
- Connected components
- LRU caching
- Sliding-window problems
- Backtracking
- Recursive traversal

For each scenario:

```text
Requirements
    ↓
Candidate data structures
    ↓
Trade-offs
    ↓
Choose one
    ↓
Explain why
```

---

### 2.36 Advantages & Disadvantages of Each Collection

For every major collection:

- Strengths
- Weaknesses
- Time complexity
- Space complexity
- Memory overhead
- Ordering behavior
- Lookup characteristics
- Insertion/removal characteristics
- Appropriate use cases
- Inappropriate use cases
- Alternatives

---

Interview Questions

### Interview Questions - Core comparisons

- `Array` vs `List<T>`
- `List<T>` vs `LinkedList<T>`
- `Dictionary<TKey,TValue>` vs `HashSet<T>`
- `Dictionary<TKey,TValue>` vs `SortedDictionary<TKey,TValue>`
- `SortedDictionary<TKey,TValue>` vs `SortedList<TKey,TValue>`
- `HashSet<T>` vs `SortedSet<T>`
- `Queue<T>` vs `Stack<T>`
- `Queue<T>` vs `PriorityQueue<TElement,TPriority>`
- `PriorityQueue<TElement,TPriority>` vs sorting
- `List<T>` vs `HashSet<T>`

### Interview Questions - Complexity and implementation

- Why is dictionary lookup usually O(1)?
- What happens when a dictionary grows?
- How are hash collisions handled?
- Why does `List<T>` have a Capacity?
- Why is `List<T>.Add()` amortized O(1)?
- Why is `Queue<T>.Enqueue()` amortized O(1)?
- Why is `Stack<T>.Push()` amortized O(1)?
- How does a circular buffer work?
- How does a binary heap work?
- How does `PriorityQueue<TElement,TPriority>` work internally?
- Why can a BST become O(n)?
- Why do balanced trees matter?
- Why isn't a `SortedDictionary` implemented as a hash table?
- Why would you use a heap instead of sorting everything?

### Interview Questions - Interfaces and abstractions

- Why does `List<T>` implement `IList<T>`?
- What is the difference between `IEnumerable<T>` and `ICollection<T>`?
- `ICollection<T>` vs `IList<T>`
- `IList<T>` vs `IReadOnlyList<T>`
- `ICollection<T>` vs `IReadOnlyCollection<T>`
- What does `ISet<T>` provide?
- Why return `IReadOnlyList<T>` instead of `List<T>`?
- What happens when a collection is exposed through an interface?

### Interview Questions - Practical scenarios

- Which collection would you use for fast lookup by ID?
- Which collection would you use to remove duplicates?
- Which collection would you use for FIFO processing?
- Which collection would you use for LIFO processing?
- Which collection would you use for priority-based processing?
- How would you implement an LRU cache?
- How would you implement autocomplete?
- How would you detect a cycle in a graph?
- How would you find connected components?
- How would you implement a shortest-path search?
- Which data structure would you use for a sliding-window problem?

---

## Part 3 — Generics

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

## Part 4 — Delegates, Lambdas and Events

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

## Part 5 — LINQ

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

## Part 6 — Memory & Performance

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

## Part 7 — Exceptions

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

## Part 8 — Async/Await

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

## Part 9 — Threads, Tasks & Parallelism

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

## Part 10 — Concurrent Programming

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

## Part 11 — Dependency Injection

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

## Part 12 — .NET Runtime

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

## Part 13 — ASP.NET Core

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

## Part 14 — HTTP & APIs

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

## Part 15 — Middleware & Pipelines

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

## Part 16 — EF Core & Data Access

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

## Part 17 — SQL

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

## Part 18 — Messaging

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

## Part 19 — Caching

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

## Part 20 — Testing & Debugging

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

## Part 21 — Design Patterns

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

## Part 22 — Distributed Systems

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

## Part 23 — System Design

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

## Part 24 — Interview Practice

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