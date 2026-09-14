# Module 2 — Distributed Systems Foundations

## Chapter 2.1 — Distributed Systems Fundamentals

This is an important transition from **“what NoSQL databases are”** to **“why distributed databases behave differently.”**

You already know the basic idea of multiple machines from Module 0/1—partitioning, replication, horizontal scaling, and the fact that network communication introduces latency and failure. So I won't re-teach those basics. Here we're going deeper into what makes a distributed system fundamentally difficult.

### 2.1.1 — What Is a Distributed System?

A **distributed system** is a system in which multiple independent computers cooperate to provide some overall functionality.

The important word is **independent**.

Each machine has its own:

- CPU
- memory
- storage
- operating-system state
- network connection
- clock
- failure state

Yet the application may need to make them behave as if they are part of one coherent system.

For example, a database cluster might have:

```text
Application
     |
     v
  Node A
   /   \
  v     v
Node B Node C
```

The application may think:

> “I am writing to the database.”

But internally, the system may need to:

1. choose a node,
2. write the data,
3. replicate it,
4. wait for acknowledgements,
5. handle a node that doesn't respond,
6. determine whether a timeout means failure,
7. potentially retry,
8. maintain consistency.

That is where distributed-systems complexity begins.

---

# 2.1.2 — Nodes

A **node** is an individual machine or process participating in the distributed system.

Depending on the system, a node might be:

- a physical server
- a virtual machine
- a cloud instance
- a container
- sometimes a database process

Don't assume that **node = physical computer**.

For example, three database containers could technically be three nodes even if they run on fewer physical machines. Conversely, one physical machine can host multiple independent processes.

The important concept is:

> A node is an independently operating participant in the distributed system.

Each node can fail independently.

That leads to one of the most important ideas in distributed systems:

### Partial failure

In a single-machine application, you might think:

> “The system is either working or it's down.”

In a distributed system, that's no longer true.

Node A might be working.

Node B might be down.

Node C might be working but unreachable from A.

The system is therefore **partially functioning**.

This is called **partial failure**.

---

# 2.1.3 — Clients

A **client** is something that communicates with the distributed system.

For a database, the client is usually your application:

```text
.NET API
   |
   | database request
   v
Database cluster
```

But clients can also be:

- another service
- a CLI
- a background worker
- a web application
- an ETL process

The important distributed-systems issue is that the client and server communicate **over a network**.

That means communication itself can fail.

---

# 2.1.4 — Networks Are Not Reliable Function Calls

This is one of the most important mental shifts.

When your application calls:

```csharp
result = CalculatePrice(order);
```

you normally expect:

- the function executes,
- it returns,
- you know whether it succeeded.

A network call is fundamentally different:

```text
Application
    |
    | request
    v
   Node
```

The request can:

- arrive
- arrive late
- never arrive
- arrive but processing takes too long
- be processed successfully but the response is lost

That last case is particularly important.

Suppose your application sends:

```text
"Charge €100"
```

The payment service processes the payment successfully.

But the response gets lost.

Your application sees:

```text
TIMEOUT
```

What does that mean?

It **doesn't know** whether:

- the request never arrived,
- the server rejected it,
- the server is still processing it,
- or the server completed it successfully and only the response was lost.

This uncertainty is fundamental to distributed systems.

---

# 2.1.5 — Network Latency

**Network latency** is the time required for communication between participants.

Even inside a data center, communication isn't instantaneous.

For example:

```text
Application → DB node
```

might involve:

- network stack
- switches
- routing
- serialization
- deserialization
- processing

And if the DB node needs to communicate with another node:

```text
Application
     |
     v
 Node A
     |
     v
 Node B
```

you have another network interaction.

This matters because operations that were cheap on one machine can become expensive when they require coordination between machines.

For example:

```text
Local operation
CPU → Memory
```

versus:

```text
Distributed operation
Node A → Network → Node B → Network → Node A
```

This is one reason distributed databases try to minimize unnecessary cross-node coordination.

### Important distinction

**Latency ≠ throughput.**

- **Latency:** how long one operation takes.
- **Throughput:** how many operations the system can process per unit of time.

You already encountered this distinction in Module 0. Here it becomes particularly important because network communication adds latency to distributed operations.

---

# 2.1.6 — Network Partitions

A **network partition** occurs when parts of a distributed system cannot communicate with each other.

Imagine:

```text
Node A  ←──X──→  Node B
```

Both nodes are alive.

Their CPUs work.

Their memory works.

Their disks work.

But the network connection between them has failed.

Now each node has a different view of the system.

With three nodes:

```text
        Network
       /       \
    Node A     Node B
       \
        X
         \
        Node C
```

Perhaps A and B can communicate, but C cannot reach them.

This is different from a node failure.

### Node failure

```text
Node A → DEAD
```

### Network partition

```text
Node A → ALIVE
Node B → ALIVE

A ←X→ B
```

Both nodes are alive, but cannot communicate.

This distinction becomes extremely important when we later study **consistency, replication, consensus, and CAP**.

---

# 2.1.7 — Timeouts

Because network operations can potentially wait forever, distributed systems use **timeouts**.

For example:

```text
Send request
     |
     v
Wait 2 seconds
     |
     +---- response → success
     |
     +---- no response → timeout
```

A timeout means:

> “I waited long enough that I am no longer willing to wait.”

It does **not necessarily mean**:

> “The operation failed.”

That's a crucial distinction.

If a database request times out, the operation might have:

- failed before execution,
- never reached the server,
- executed but response was delayed,
- executed successfully but response was lost.

Therefore:

> **Timeout is evidence of uncertainty, not proof of failure.**

This becomes extremely important when we discuss retries.

---

# 2.1.8 — Retries

A retry means sending the operation again after a failure or timeout.

For example:

```text
Request
   ↓
Timeout
   ↓
Retry
   ↓
Success
```

Retries can improve reliability when failures are temporary.

But retries introduce a dangerous problem:

## What if the first request actually succeeded?

Consider:

```text
Client
  |
  | "Create reservation"
  v
Server
  |
  | reservation created
  X
  response lost
  |
Client receives timeout
```

The client retries:

```text
"Create reservation"
```

Now the operation might execute twice.

For operations that aren't naturally safe to repeat, this can cause serious problems.

This is why distributed systems care about **idempotency**.

An operation is **idempotent** if performing it multiple times has the same intended effect as performing it once.

For example, conceptually:

```text
Set status = "PAID"
```

can be designed to be idempotent.

Whereas:

```text
Charge customer €100
```

is not inherently idempotent.

This is why production systems often use:

- idempotency keys
- request IDs
- deduplication
- retry policies
- exponential backoff

We'll go deeper into these later rather than prematurely expanding this chapter.

---

# 2.1.9 — Failure Detection

Here's another subtle problem.

Suppose Node A sends a request to Node B:

```text
A → B
```

B doesn't respond.

Is B dead?

You **cannot know with certainty merely from the missing response**.

Possibilities include:

1. B crashed.
2. B is overloaded.
3. The network is partitioned.
4. The request was lost.
5. The response was lost.
6. B is processing the request very slowly.

This is why distributed systems use mechanisms such as:

- heartbeats
- health checks
- timeouts
- failure detectors
- leader-election mechanisms

But these mechanisms don't magically provide perfect knowledge.

A failure detector is effectively making an inference:

> “Based on the evidence available to me, I believe this node is unavailable.”

That distinction becomes very important later when we study distributed consensus.

---

# 2.1.10 — Split Brain

**Split brain** occurs when different parts of a distributed system independently believe they should be in control.

For example:

```text
        Network partition
             X
            / \
           /   \
        Node A Node B
          ↓     ↓
       "I'm     "I'm
        leader" leader"
```

Both sides believe they are the leader.

That can be dangerous if both accept writes independently.

For example:

```text
Client 1 → A → update X = 10

Client 2 → B → update X = 20
```

Now the system has conflicting histories.

Preventing split brain is one reason distributed databases and distributed systems need mechanisms for **coordination and leader election**.

We'll study those mechanisms in more depth later.

---

# 2.1.11 — Clock Differences

Every machine has its own clock.

They are **not perfectly synchronized**.

For example:

```text
Node A: 12:00:01.100
Node B: 12:00:01.250
```

Node B's clock is 150 ms ahead.

This becomes surprisingly important.

Suppose:

```text
Node A:
"Order created at 12:00:01.100"

Node B:
"Payment created at 12:00:01.050"
```

Can we conclude the payment happened before the order?

**Not necessarily.**

The clocks themselves may be different.

This is why distributed systems cannot blindly use local wall-clock timestamps to determine the exact global ordering of events.

There are techniques for dealing with this, including:

- clock synchronization
- logical clocks
- hybrid logical clocks
- database-specific timestamp mechanisms

We'll encounter these later when they become relevant.

---

# 2.1.12 — Distributed Coordination

This is the underlying problem connecting many of the concepts above.

Suppose multiple nodes need to agree on something:

> Who is the leader?

> Which write happened first?

> Has this data been committed?

> Which replica should accept writes?

> Is another node still alive?

> Can we safely proceed after a network failure?

The nodes need **coordination**.

But coordination itself requires communication.

And communication is:

- delayed,
- potentially lost,
- potentially duplicated,
- potentially partitioned.

So we get a fundamental distributed-systems tension:

> **To coordinate, nodes must communicate. But communication itself can fail.**

This is why distributed systems are fundamentally harder than single-node systems.

---

# The Mental Model

The important concepts in this chapter connect like this:

| Concept | Core problem |
|---|---|
| **Nodes** | Independent participants |
| **Clients** | Consumers communicating with nodes |
| **Network latency** | Communication takes time |
| **Network partition** | Nodes may be unable to communicate |
| **Node failure** | A participant may disappear |
| **Partial failure** | Some components work while others don't |
| **Timeout** | We stop waiting, but may not know what happened |
| **Retry** | Improves transient-failure handling but can duplicate operations |
| **Failure detection** | Inferring whether another node is unavailable |
| **Split brain** | Multiple parts believe they are in control |
| **Clock differences** | Local time cannot reliably establish global ordering |
| **Coordination** | Nodes must agree despite unreliable communication |

### The central lesson

A distributed system isn't difficult merely because there are **more machines**.

It's difficult because:

> **Independent machines must cooperate even though communication is delayed and can fail, machines can fail independently, and no single node necessarily has complete knowledge of the system.**

That is the foundation we'll use for the next chapters.