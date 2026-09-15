# Module 1 — System Design Fundamentals

First, establish the basic vocabulary and way of thinking.

### 1.1 What is a distributed system?

- What makes a system distributed
- Single-machine vs distributed systems
- Why we distribute systems
- Scalability
- Availability
- Reliability
- Performance
- Fault tolerance
- Durability
- Why distributed systems are difficult

### 1.2 Components of a modern backend

Understand the role of:

- Client
- DNS
- Load balancer
- API gateway (cross-cutting edge concerns; detailed after load balancing)
- API server
- Database
- Cache
- Object storage
- Message queue
- Event stream
- Search engine
- CDN

Not designing anything yet. Just understand **what each component does and why it exists**.

### 1.3 Request lifecycle

Understand what actually happens when:

```text
User → example.com → API → database → response
```

Including:

- DNS lookup
- TCP connection
- TLS
- HTTP request
- load balancing
- application processing
- database request
- response

We don't need to turn this into a networking course. Just enough to reason about system design.

### 1.4 Stateful vs stateless services

- What state means
- Stateful servers
- Stateless servers
- Why stateless APIs are easier to scale
- Where state should live instead
- Sessions
- Session storage

### 1.5 Horizontal vs vertical scaling

- Vertical scaling
- Horizontal scaling
- Scaling up vs scaling out
- When each works
- Why horizontal scaling eventually becomes important

### 1.6 Requirements and capacity estimation

This is a very practical missing piece before design work begins.

- Functional vs non-functional requirements
- Read-heavy vs write-heavy workloads
- Peak traffic vs average traffic
- Requests per second
- Data volume growth over time
- Storage estimates
- Bandwidth and egress estimates
- Latency targets (average vs tail — measured precisely in observability)
- Availability targets
- Cost constraints
- What to clarify before designing
- Little's Law (relationship of concurrency, latency, and throughput — sanity-check in-flight work from RPS and latency targets)

The purpose is to turn vague requirements like "millions of users" into concrete constraints that guide architecture decisions.

---

# Module 2 — Scalability

This is where we learn how systems handle increasing traffic.

### 2.1 Capacity and bottlenecks

- Requests per second
- Throughput
- Latency
- Concurrent requests
- Little's Law: average concurrency ≈ throughput × average latency (ties to Module 1.6)
- CPU-bound workloads
- Memory-bound workloads
- I/O-bound workloads
- Identifying bottlenecks

### 2.2 Load balancing

- Why load balancers exist
- Layer 4 vs Layer 7
- Round robin
- Least connections
- Health checks
- Connection draining
- Load-balancer failure
- Multiple load balancers
- Sticky sessions / session affinity
- When affinity helps vs when it fights horizontal scaling

### 2.2.1 API gateway vs load balancer

Introduced here immediately after load-balancer basics; expanded again in Module 9 when service boundaries matter.

- Load balancer: distribute connections and requests across healthy backends
- API gateway: edge entry point for routing, aggregation, and cross-cutting API concerns
- Termination of TLS and protocol handling at the edge (conceptual)
- Request routing to the correct backend service
- How gateways relate to load balancers in a typical topology
- Backend for frontend (BFF) as a special case of edge/API shaping

### 2.3 Horizontal application scaling

- Multiple API instances
- Stateless applications
- Session handling (and when sticky sessions are a workaround)
- Shared state
- Service discovery
- Auto-scaling

### 2.4 Database scaling

- Vertical database scaling
- Read replicas
- Primary/replica architecture
- Read/write splitting
- Replication lag
- Database connection limits

### 2.5 Database sharding

This deserves substantial attention.

- What sharding means
- Why sharding is needed
- Shard key
- Hash-based sharding
- Range-based sharding
- Consistent hashing
- Hot partitions
- Cross-shard queries
- Resharding

And importantly:

**Replication vs sharding**

This distinction needs to become automatic for you.

### 2.6 Resource capacity and exhaustion

- Database connection pools
- Connection limits
- HTTP keep-alive and connection reuse
- Thread-pool and event-loop exhaustion
- File descriptors and other finite resources
- Queue and worker capacity
- Concurrency limits
- Why a fast dependency can still exhaust callers
- Backpressure across service boundaries
- Load shedding and admission control (reject or delay work early to protect core capacity — expanded in Module 7.4; distinct from rate limiting in Module 9.2)
- Multi-tenant noisy-neighbor effects at shared pools (expanded with tenancy in Module 9.7)

### 2.7 Distributed identifiers

After sharding and capacity reasoning — IDs are part of keys, ordering, and deduplication everywhere later.

- Why auto-increment and single-node sequences stop scaling
- Universally unique identifiers (UUIDs / ULIDs)
- Time-ordered IDs (Snowflake-style) and rough ordering trade-offs
- Clock skew risks for time-based IDs (ties to Module 5.1.1)
- Choosing IDs as shard keys vs separate shard routing
- IDs for idempotency keys, events, and messages (used in Module 5.4 and Module 6)

---

# Module 3 — Caching

Caching comes up constantly in system design, so we should understand it properly rather than treating "use Redis" as a magic answer.

### 3.1 What is caching?

- Why caches work
- Cache hit
- Cache miss
- Hit ratio
- Cacheable vs non-cacheable data

### 3.2 Cache patterns

- Cache-aside
- Read-through
- Write-through
- Write-back
- Write-around

### 3.3 Cache invalidation

- TTL
- Explicit invalidation
- Stale data
- Cache consistency
- Cache stampede
- Cache penetration
- Bloom filters and negative caching (when they help vs operational cost)
- Cache avalanche

### 3.4 Distributed caching

- Why local memory caches aren't enough
- Redis/Memcached
- Multiple application servers
- Shared cache
- Cache partitioning
- Replication

### 3.5 CDN caching

- Edge caching
- Origin server
- Cache-Control
- Immutable content
- Cache invalidation
- Geographic distribution

This will make the CDN discussion from YouTube much easier to understand.

---

# Module 4 — Data Storage

Here we'll build the database knowledge required for choosing storage correctly.

## 4.1 Relational databases

- Tables
- Primary keys
- Foreign keys
- Indexes
- Transactions
- ACID
- Isolation
- Read phenomena and isolation anomalies
- Deadlocks
- Lock contention
- Transaction duration and throughput
- Constraints
- Normalization
- Denormalization

You already know SQL, so this won't be a beginner SQL course.

The focus will be:

> **What properties of relational databases matter when designing a distributed system?**

### 4.1.1 Data modeling for system design

- Access-pattern-driven schema design
- Choosing keys and indexes
- Index write and storage costs
- Hot rows and hot keys
- Time-based data
- Retention and archival
- Denormalization boundaries
- Separating transactional and analytical workloads

### 4.1.2 Database performance in production

- Query plans
- Slow queries
- Index selection and index maintenance cost
- Lock contention
- Hot rows
- Connection pool pressure
- Read/write workload separation

---

## 4.2 NoSQL databases

- Key-value stores
- Document databases
- Wide-column databases
- Graph databases and time-series stores (when access patterns justify them — not a default alongside key-value, document, or wide-column)
- When NoSQL makes sense
- Access-pattern-driven design
- Partitioning
- Replication
- Consistency models

We'll connect this with the NoSQL curriculum you're already doing rather than duplicating everything.

---

## 4.3 Object storage

This is particularly important because it appeared in YouTube.

- Object vs row/document
- Bucket
- Object key
- Metadata
- Object lifecycle
- Large file storage
- Multipart/resumable uploads
- Pre-signed URLs
- Object storage vs database
- Object storage vs filesystem

This should make things like:

> `videos/123/1080p/segment-001`

immediately intuitive.

---

## 4.4 Database consistency

- Strong consistency
- Eventual consistency
- Read-after-write consistency
- Replication lag
- Consistency vs latency
- When eventual consistency is acceptable
- Serializability

At this stage, focus on the practical distinction between stale reads and strongly consistent reads. Advanced consistency models, quorum, and conflict resolution are introduced after replication and coordination in Module 5.

This is where the vocabulary around distributed consistency becomes precise enough for interview reasoning.

---

# Module 5 — Distributed Systems Fundamentals

This is probably the **most important module** for you.

And this is where terms like **split brain** will be introduced properly rather than appearing unexpectedly in a system-design answer.

## 5.1 Nodes and communication

- Nodes
- Clients
- Networks
- Network latency
- Network partitions
- Node failures
- Partial failure

### 5.1.1 Time and ordering

- Clock skew
- Why wall-clock timestamps cannot always establish event order
- Logical clocks
- Timestamp ordering
- Causality and happens-before relationships
- Practical uses and limitations of timestamps in distributed systems

## 5.2 Failure models

- Crash failure
- Network failure
- Timeout
- Message loss
- Duplicate messages
- Delayed messages
- Byzantine failure (briefly; mainly so you know the term)

### 5.3 Timeouts and retries

- Why distributed calls need timeouts
- Retryable vs non-retryable failures
- Retry storms
- Exponential backoff
- Jitter
- Maximum retries

### 5.4 Idempotency

- What idempotency means
- Why retries create duplicate operations
- Idempotency keys
- Idempotent APIs
- Idempotent message processing
- Ambiguous outcomes after timeouts (especially payments)
- Reconciliation: matching internal state with external provider or ledger truth
- Reconciliation jobs vs inline retries

This is extremely important for interviews.

---

## 5.5 Replication

- Why replication exists
- Primary/replica
- Leader/follower terminology
- Synchronous replication
- Asynchronous replication
- Replication lag
- Failover
- Promotion

### 5.6 Failover

- Health detection
- Choosing a new primary
- Promoting a replica
- Routing traffic
- Recovery

### 5.7 Split brain

Finally, the term you encountered.

- What split brain means
- How it happens
- Why it is dangerous
- Two nodes believing they are primary
- Fencing
- Leader election
- Quorum
- Preventing two writers

This should be a full lesson because it appears in many distributed systems.

### 5.8 Concurrency and coordination

This is essential for real systems, especially booking, payments, inventory, and ordering.

- Atomicity and atomic operations
- Optimistic vs pessimistic concurrency control
- Lost updates
- Compare-and-set / version checks
- Distributed locks and leases
- Lease expiry and fencing tokens
- Reservation races
- Idempotent writes under retries
- Ordering guarantees
- Why distributed transactions are hard
- Outbox pattern
- Sagas / workflow compensation
- Consensus at a conceptual level
- Coordination services such as etcd or ZooKeeper as examples
- Why distributed locks and leader election require coordination

The goal is to understand how multiple actors safely update state without corrupting it.

### 5.9 Advanced consistency and coordination

Only after understanding replication, failures, and coordination:

- Linearizability
- Monotonic reads
- Causal consistency
- Session guarantees
- Quorum reads and writes
- Read and write quorums
- Conflict resolution
- Consistency versus availability during partitions

---

# Module 6 — Messaging and Event-Driven Systems

This is another huge system-design topic.

Teach **queues vs event streams** before diving into either implementation detail.

## 6.1 Queues vs event streams

- Task queue vs log / event stream (do not conflate them)
- Point-to-point work distribution vs durable ordered log
- Consumer pulls work vs multiple independent consumer groups
- Delete-on-ack vs retention and replay
- When a queue is the right abstraction
- When an event stream is the right abstraction
- Fan-out with queues vs topics/partitions
- How this choice affects delivery guarantees and debugging

## 6.2 Message queues

- Producer
- Consumer
- Queue
- Message
- Acknowledgement
- Visibility timeout
- Consumer scaling
- Backpressure

## 6.3 Queue patterns

- Work queues
- Fan-out
- Competing consumers
- Delayed messages
- Priority queues

## 6.4 Event streams

- Topic
- Partition
- Producer
- Consumer
- Consumer group
- Offset
- Ordering within a partition
- Retention
- Replay
- Event stream as durable log vs stream processing (continuous aggregation — view counts, metrics; architectural meaning, not cluster administration)

Kafka is the common example; the goal is architectural meaning, not cluster administration.

---

## 6.5 Delivery guarantees

- At-most-once
- At-least-once
- Exactly-once
- Why exactly-once is difficult
- Deduplication
- Idempotent consumers
- Stream processing for counters and aggregates (builds on event streams in 6.4 — deduplication and at-least-once delivery still apply)

This will directly connect to the YouTube view-count discussion.

---

## 6.6 Dead-letter queues and poison messages

- Why messages fail
- Poison / toxic messages that never succeed
- Retry policy
- Maximum retries
- DLQ
- Reprocessing and manual intervention
- Idempotency when reprocessing from a DLQ

---

# Module 7 — Reliability and Availability

Now we learn how to design systems that continue operating when things fail.

### 7.1 Availability

- What availability means
- Uptime percentages
- 99%
- 99.9%
- 99.99%
- Downtime calculations
- Single points of failure

### 7.2 Fault tolerance

- Redundancy
- Failover
- Replication
- Multiple instances
- Multiple availability zones
- Multiple regions

### 7.3 Disaster recovery

- Backup
- Restore
- Replication
- RPO
- RTO
- Disaster recovery strategies

### 7.4 Graceful degradation and overload protection

- Partial functionality
- Fallbacks
- Cached responses
- Feature degradation
- Dependency failure
- Load shedding under overload (fail fast, drop non-critical work)
- Admission control (limit in-flight work entering the system)
- How shedding differs from rate limiting (Module 9.2) and from backpressure (Module 2.6)

### 7.5 Circuit breakers

- Why cascading failures happen
- Open/closed/half-open
- When circuit breakers help

### 7.6 Bulkheads

- Isolation
- Resource pools
- Preventing one workload from consuming everything

---

# Module 8 — Global Systems

Only after everything above should we tackle global architecture.

### 8.1 Availability zones and regions

- Data center
- Availability Zone
- Region
- Regional failure

### 8.2 Global traffic routing

- DNS-based routing
- Geo routing
- Latency-based routing
- Health checks
- Failover routing

### 8.3 Multi-region databases

This is where the questions from YouTube become much easier:

- Primary in one region
- Regional read replicas
- Multi-primary
- Local writes
- Replication
- Conflict resolution
- Consistency trade-offs

### 8.4 Global data

- What should be globally replicated?
- What can remain regional?
- User-local data
- Global data
- CDN
- Regional caches

### 8.5 CAP theorem

We'll learn this **after** understanding failures, replication, partitions and consistency.

Not before.

- Consistency
- Availability
- Partition tolerance
- Why partition tolerance isn't really optional in distributed systems
- What CAP actually says
- Common CAP misconceptions

### 8.6 PACELC

- CAP isn't the whole story
- Latency vs consistency during normal operation
- Why distributed databases make these trade-offs

---

# Module 9 — System Design Building Blocks

This is the bridge between fundamentals and actual interview design.

Here we'll learn the standard components and when they are appropriate.

### 9.1 API design

- REST
- RPC
- Request and response serialization
- JSON, binary formats, and serialization trade-offs
- Payload size, latency, and CPU costs
- Compression and its trade-offs
- Schema compatibility
- Schema registries
- Pagination
- Cursor pagination
- Authentication vs authorization (who is the caller vs what they may do)
- Sessions vs bearer tokens (conceptual trade-offs; not a specific OAuth/JWT course)
- Resource ownership checks (e.g. user may only access their own objects — ties to file and listing designs)
- Secrets and credentials (never in logs or source control; rotation at a conceptual level)
- Encryption in transit (builds on TLS in Module 1.3) and encryption at rest for sensitive data
- PII, privacy, and data minimization (what to store; cross-link Module 8.4 global/regional data and Module 9.9 retention)
- Audit and access logs when requirements demand accountability
- Inbound webhooks and outbound callbacks: verify sender, signatures, replay protection, idempotent handling
- Abuse and credential-stuffing protection at the edge (distinct from overload load shedding in Module 7.4; complements rate limiting in Module 9.2)
- Pre-signed and time-limited access (ties to Module 4.3; leakage risk if URLs are shared)
- Rate limiting (detailed in Module 9.2)
- Idempotency (ties to Module 5.4, especially payments)
- API versioning

### 9.2 Rate limiting

- Why it exists
- Fair use and abuse throttling (bots, credential stuffing — builds on security concerns in Module 9.1)
- Token bucket
- Leaky bucket
- Fixed window
- Sliding window
- Distributed rate limiting

### 9.3 Search systems

- Why SQL isn't a search engine
- Inverted indexes
- Elasticsearch/OpenSearch-style systems
- Indexing pipelines
- Search consistency

### 9.4 File processing

- Upload
- Object storage
- Async processing
- Queues
- Workers
- Progress/status tracking

### 9.5 Notifications

- Email
- Push
- SMS
- Webhooks and provider callbacks (delivery status, bounces; signature verification and idempotency — ties to Module 9.1)
- Queue-based delivery
- Retry
- Deduplication
- Provider failures

### 9.6 Background jobs

- Workers
- Schedulers
- Queues
- Retries
- Idempotency
- Job state

### 9.7 Service boundaries, gateways, and synchronous dependencies

- Modular monolith vs microservices
- Identifying service boundaries
- Synchronous call chains
- Fan-out
- Dependency graphs
- Cascading latency and failure
- API contracts
- Backward compatibility
- When not to split a service
- API gateway in multi-service designs (builds on Module 2.2.1)
- Edge routing, request aggregation, and protocol translation
- Backend for frontend (BFF) patterns
- Multi-tenancy: tenant isolation, noisy neighbors, fair sharing of shared pools
- When tenancy affects sharding, caching, and rate limits

Microservices should be introduced as one possible response to a problem, not as a default architecture.

### 9.8 Schema evolution and zero-downtime changes

This is a big real-world topic that often gets skipped too early.

- Schema changes in production
- Backward-compatible changes
- Forward-compatible changes
- Expand and contract patterns
- Backfill jobs
- Dual writes and read paths
- Versioning events and APIs
- Serialization compatibility across producers and consumers
- Schema registry use and compatibility modes
- Avoiding breaking clients
- Rolling deployments and migration safety

Production systems do not remain static. We need to understand how to evolve schemas and contracts safely.

### 9.8.1 Deployment strategies

- Rolling deployments
- Blue-green deployments
- Canary deployments
- Progressive traffic shifting
- Health checks and readiness during deployment
- Rollback strategies
- Backward-compatible application and database changes
- Coordinating deployments with schema migrations

### 9.9 Data lifecycle and workload separation

- Retention policies
- Archival and cold storage
- TTL-based cleanup
- Online versus offline processing
- Transactional, analytical, and search workloads
- Read models and materialized views
- Change data capture (CDC): propagating writes to search indexes, caches, and read models without dual-write races (ties to Module 9.3 and Module 9.9.1)

### 9.9.1 Event sourcing and CQRS

Only after messaging, delivery guarantees, outbox (Module 5.8), and read models above.

- Command Query Responsibility Segregation (CQRS): separate write model from read-optimized projections
- Event sourcing: append-only event log as the system of record
- Projections and rebuild from replay (ties to Module 6.4)
- Comparison with outbox + asynchronous consumers
- Consistency and lag on read models
- Operational cost: replay, schema evolution on events, debugging
- When these patterns are worth it vs overkill

### 9.10 Real-time delivery to clients

Before chat and live features in interview practice — builds on stateful vs stateless (Module 1.4) and horizontal scaling (Module 2.3).

- Request/response limits for live updates
- Long polling
- Server-Sent Events (SSE)
- WebSockets and persistent connections
- Scaling connection-heavy tiers (affinity, shared pub/sub backplane — introduced here, not before)
- Where connection state should live
- Push notifications vs in-app real-time (ties to Module 9.5)

---

# Module 10 — Observability

This is relatively short but important.

### 10.1 Metrics

- Counters
- Gauges
- Histograms
- Percentiles (p50, p95, p99)
- Tail latency and the long tail
- Why averages mislead for user-facing SLOs
- Throughput
- Error rate
- Latency

### 10.2 Logging

- Structured logs
- Correlation IDs
- Log levels
- Centralized logging

### 10.3 Distributed tracing

- Trace
- Span
- Trace propagation
- Finding latency across services

### 10.4 What should we monitor?

- Golden signals
- Saturation
- Dependency failures
- Queue depth
- Database health
- Replication lag

### 10.5 SLOs, operational trade-offs, and cost

These are part of how real system designs are evaluated.

- SLIs and SLOs
- Error budgets
- Operational readiness
- Runbooks and incident response
- Cost of redundancy
- Cost of consistency and replication
- Regional vs single-region trade-offs
- Simpler systems often win when requirements allow
- Storage, compute, and egress costs
- Replication and redundancy costs
- Cost per request or per active user
- Operational complexity as a cost
- Cost-versus-latency trade-offs

Architecture is not only about technical elegance. We need to reason about what is worth operating and what is worth paying for.

### 10.6 Operational validation

- Load testing
- Capacity testing
- Performance testing
- Failure injection
- Testing timeouts, retries, and degraded dependencies
- Backup and restore testing
- Failover testing
- Verifying recovery objectives in practice
- Observability needed to interpret test results

---

# Module 11 — Putting the fundamentals together

This is the transition from isolated concepts to architecture reasoning. It is still not full interview practice.

Instead, we'll do small architecture exercises.

For example:

> "Your API has 100 instances and uses a database. One database instance is overloaded. What options do you have?"

You should be able to reason:

- vertical scaling
- read replicas
- caching
- sharding
- query optimization
- workload separation

Or:

> "A payment request times out, but you don't know whether the payment actually succeeded."

You should immediately think:

- retry
- duplicate operation
- idempotency key
- transaction state
- reconciliation

Or:

> "The primary database dies."

You should immediately know the vocabulary:

- replica
- failover
- promotion
- leader
- fencing
- split brain
- RPO
- RTO

### 11.1 Capacity and bottleneck exercises

- Turn requirements into traffic, storage, bandwidth, and latency estimates
- Identify the first likely bottleneck
- Compare vertical scaling, caching, replication, sharding, and workload separation
- Include connection, queue, and worker capacity in estimates

### 11.2 Data and consistency exercises

- Choose storage based on access patterns
- Identify hot keys and hot partitions
- Decide where strong consistency is required
- Decide where eventual consistency is acceptable
- Explain the consequences of replication lag

### 11.3 Reliability and failure exercises

- Trace failures through a dependency graph
- Choose timeout, retry, backoff, and circuit-breaker behavior
- Reason about idempotency, duplicate work, and reconciliation
- Choose a queue vs an event stream for a stated workload
- Distinguish availability, durability, RPO, and RTO
- Explain graceful degradation, load shedding, and recovery

### 11.4 Evolution and operations exercises

- Perform an expand-and-contract schema migration
- Evolve an API or event without breaking consumers
- Define useful metrics, logs, traces, and SLOs
- Estimate operational and infrastructure cost
- Design a load, failover, or restore test

---

# Module 12 — Full System-Design Interview Practice

This is the dedicated practice phase after the fundamentals. The goal is to apply the concepts under realistic interview constraints, not to learn new terminology for the first time.

## 12.1 A repeatable design process

- Clarify functional requirements
- Clarify non-functional requirements
- Identify scale, traffic shape, and data size
- Define core APIs and important data models
- Draw the simplest viable high-level architecture
- Identify bottlenecks and single points of failure
- Add scaling, reliability, and consistency mechanisms only when justified
- Discuss observability, operations, and cost
- State assumptions and trade-offs clearly

## 12.2 Interview communication

- Narrating decisions and assumptions
- Keeping the design at the right level of detail
- Explaining alternatives without losing focus
- Handling changing requirements
- Responding to interviewer follow-up questions
- Recognizing when a simpler design is sufficient
- Summarizing the final design and its limitations

## 12.3 Practice designs

Use progressively more complex designs:

- URL shortener
- Rate-limited API
- File storage and processing
- Notification system
- Search and indexing system
- News feed
- Chat or messaging system
- Booking or inventory system
- Payment workflow
- Video upload and delivery
- Globally distributed user-facing service
- Webhook-driven third-party integration (payment or email provider status callbacks)
- Metrics or analytics pipeline (counter aggregation at scale; hot keys, idempotency, optional stream processing)
- Leaderboard or global counter at scale

For every design, practice:

- Requirements and capacity estimates
- API and data model
- Read and write paths
- Storage choice
- Caching strategy
- Asynchronous work
- Consistency and concurrency
- Failure handling
- Observability
- Cost and operational complexity

## 12.4 Follow-up and failure drills

- Database overload
- Cache outage or stampede
- Queue backlog
- Poison messages and DLQ buildup
- Duplicate message processing
- Stale replica reads
- Primary failure and failover
- Regional outage
- Hot partition
- Dependency timeout
- Shedding overload while preserving core functionality
- WebSocket or live-connection tier failure
- Backward-incompatible schema or API change

## 12.5 Evaluation checklist

- Did the design satisfy the stated requirements?
- Were capacity assumptions explicit and plausible?
- Were bottlenecks identified?
- Were replication and sharding used for the correct purposes?
- Were consistency and failure trade-offs explained?
- Were retries made safe with idempotency?
- Were operational costs and complexity considered?
- Could the design be evolved without downtime?
