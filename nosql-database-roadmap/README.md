# NoSQL Master Curriculum
### Zero → Production-Ready → Senior/Principal-Level Architecture & Interviews

I’d structure this as a **course/book rather than a MongoDB tutorial**. The goal is to understand the underlying distributed-database ideas first, then see how MongoDB, Redis, DynamoDB, Cassandra/ScyllaDB, and Neo4j make different trade-offs.

Each major stage includes **theory + implementation + architecture + interview preparation**.

---

## MODULE 0 — Prerequisites & Mental Model

### Chapter 0.1 What Is a Database?
    - Storage vs database
    - Database engine
    - Storage engine
    - Memory/cache
    - Disk/SSD
    - Query processing
    - Transactions/concurrency
    - Durability/availability
    - Latency/throughput

### Chapter 0.2 What Happens When You Execute a Query?
    - Client → DB
    - Parsing
    - Planning
    - Index lookup
    - Data retrieval
    - Concurrency
    - Result construction
    - Network response
    - Sources of latency

### Chapter 0.3 Relational Databases as the Baseline
    - Tables
    - Rows/columns
    - PK/FK
    - JOINs
    - Normalization
    - Relational transactions

### Chapter 0.4 SQL vs NoSQL Mental Model
    - Why NoSQL emerged
    - Distributed data
    - Denormalization
    - Partitioning
    - Replication
    - Flexible schemas
    - Trade-offs

### Chapter 0.5 CPU, I/O, Memory & Database Latency

### Chapter 0.6 Transactions & ACID

### Chapter 0.7 B-Trees, Indexes & Read Replicas

### Chapter 0.8 Relational Scaling & Horizontal Scaling

### Chapter 0.9 Flexible Schemas & Relational vs NoSQL

### Chapter 0.10 Module 0 Interview Checkpoint

---

# MODULE 1 — Why NoSQL Exists

## Chapter 1.1 — The Problems NoSQL Tries to Solve

- Massive datasets
- Massive request volumes
- Globally distributed applications
- High availability
- Low-latency access
- Flexible schemas
- Rapid schema evolution
- Distributed workloads

## Chapter 1.2 — SQL vs. NoSQL: Under the Hood

### Relational approach

- Strong schema
- Relations
- JOINs
- ACID transactions
- Centralized coordination

### Typical NoSQL approach

- Flexible data structures
- Denormalization
- Application-controlled relationships
- Partitioning
- Replication
- Distributed coordination
- Tunable consistency

## Chapter 1.3 — The NoSQL Landscape

Four major families:

```text
                 NoSQL
                   |
       +-----------+-----------+
       |           |           |
   Document    Key-Value   Wide-Column
       |
     Graph
```

Then understand that real products increasingly overlap.

Examples:

- MongoDB → Document
- Redis → Key-value/data structures
- DynamoDB → Key-value/document
- Cassandra → Wide-column
- ScyllaDB → Wide-column
- Neo4j → Graph

## Chapter 1.4 — NoSQL Myths

- "NoSQL means no schema"
- "NoSQL doesn't support transactions"
- "NoSQL is always eventually consistent"
- "NoSQL is always faster"
- "NoSQL replaces SQL"
- "MongoDB is just JSON storage"
- "Horizontal scaling is free"

### Interview Topics

- Why would a company choose MongoDB instead of PostgreSQL?
- Why might DynamoDB be a terrible choice for a particular application?
- When is SQL clearly the better choice?
- What does "schema flexibility" actually mean?

---

# MODULE 2 — Distributed Systems Foundations

This is one of the **most important modules for interviews**.

## Chapter 2.1 — Distributed Systems Fundamentals

- Nodes
- Clients
- Networks
- Network latency
- Network partitions
- Node failures
- Partial failure
- Timeouts
- Retries
- Failure detection
- Split brain
- Clock differences
- Distributed coordination

## Chapter 2.2 — CAP Theorem

Understand it from first principles.

- Consistency
- Availability
- Partition tolerance
- What a network partition actually means
- Why partition tolerance isn't really optional
- CP systems
- AP systems

### Important distinction

CAP is **not**:

> "Pick two of three whenever you want."

Understand what happens **during a partition**.

## Chapter 2.3 — BASE

- Basically Available
- Soft state
- Eventual consistency
- Why BASE emerged
- ACID vs. BASE
- Why BASE isn't simply "bad consistency"

## Chapter 2.4 — PACELC

- Partition
- Availability
- Else
- Latency
- Consistency

Understand why a system still has a trade-off **when there is no partition**.

## Chapter 2.5 — Consistency Models

From strongest to weaker models:

- Linearizability
- Sequential consistency
- Causal consistency
- Read-your-writes
- Monotonic reads
- Eventual consistency

## Chapter 2.6 — Quorums

- N
- R
- W
- R + W > N
- Read quorum
- Write quorum
- Strong vs. eventual reads
- Failure scenarios
- Quorum mathematics

## Chapter 2.7 — Replication

- Primary/replica
- Leader/follower
- Multi-leader
- Leaderless
- Synchronous replication
- Asynchronous replication
- Replication lag
- Failover
- Election

### Interview Topics

- Explain CAP without using textbook terminology.
- Why is CAP misunderstood?
- What happens when R + W ≤ N?
- Can an AP system provide strong consistency?
- Is eventual consistency always unacceptable?
- Why does synchronous replication hurt latency?
- What happens if the leader dies during a write?

### Hands-On Project

**Build a tiny distributed key-value store simulation**

Start with:

```text
Client
  |
  +---- Node A
  +---- Node B
  +---- Node C
```

Simulate:

- Replication
- Node failure
- Network partition
- Quorum reads/writes
- Eventual convergence

This will make CAP/PACELC much more concrete.

---

# MODULE 3 — Scaling: Partitions, Shards & Replicas

## Chapter 3.1 — Vertical Scaling

- Bigger CPU
- More RAM
- Faster storage
- Limitations
- Cost curves
- Single-node bottlenecks

## Chapter 3.2 — Horizontal Scaling

- Adding nodes
- Distributed workloads
- Stateless application servers
- Database distribution

## Chapter 3.3 — Partitioning

- What partitioning means
- Range partitioning
- Hash partitioning
- Composite partitioning

## Chapter 3.4 — Sharding

- Shard key
- Hash-based sharding
- Range-based sharding
- Consistent hashing
- Virtual nodes
- Shard distribution

## Chapter 3.5 — Hot Partitions

- Poor shard keys
- Sequential IDs
- Celebrity/user hotspots
- Time-based hotspots
- Write hotspots
- Read hotspots

## Chapter 3.6 — Rebalancing

- Adding nodes
- Moving partitions
- Consistent hashing
- Data movement
- Rebalancing cost

## Chapter 3.7 — Replication vs. Sharding

This distinction needs to become instinctive.

```text
Replication:
"How many copies?"

Sharding:
"Which node owns the data?"
```

## Chapter 3.8 — Replica Sets

- Primary
- Secondary
- Election
- Failover
- Majority
- Replication lag

### Interview Topics

- Sharding vs. partitioning
- Sharding vs. replication
- How do you choose a shard key?
- What makes a bad shard key?
- How do you handle hot partitions?
- What happens when a shard becomes unavailable?

### Hands-On Project

Take a large synthetic dataset and design:

- 3-node cluster
- Partition key
- Replica strategy
- Failure scenario
- Rebalancing strategy

---

# MODULE 4 — NoSQL Data Modeling

This is where NoSQL starts becoming substantially different from relational modeling.

## Chapter 4.1 — Entity Modeling vs. Access-Pattern Modeling

Traditional thinking:

> What entities exist?

NoSQL thinking:

> What queries must the application perform?

## Chapter 4.2 — Denormalization

- Why duplicate data?
- Read optimization
- Write amplification
- Maintaining duplicated data
- Stale copies
- Materialized views

## Chapter 4.3 — Embedding

- One-to-one
- One-to-many
- Nested objects
- Read locality
- Atomic updates

## Chapter 4.4 — Referencing

- When data should remain separate
- Large collections
- Independent lifecycles
- Many-to-many relationships

## Chapter 4.5 — Modeling Trade-offs

- Read-heavy vs. write-heavy
- Small vs. large documents
- Frequently updated data
- Data duplication
- Consistency requirements

## Chapter 4.6 — Cardinality

- Low cardinality
- High cardinality
- One-to-few
- One-to-many
- Many-to-many

## Chapter 4.7 — Write Amplification & Read Amplification

- Duplicate writes
- Multiple reads
- Fan-out
- Fan-in
- Query-driven schema design

### Interview Topics

Given:

> "Users can follow other users and view their feed."

Ask:

- Embed or reference?
- What is the access pattern?
- How would you model the feed?
- What happens with a user who has 50 million followers?

### Hands-On Project

Design a **social-media backend data model**.

Requirements:

- Users
- Posts
- Comments
- Likes
- Followers
- Feed
- Notifications

First model it relationally.

Then redesign it for:

- MongoDB
- DynamoDB
- Cassandra

The differences will teach you more than reading three tutorials.

---

# MODULE 5 — Document Databases

## Chapter 5.1 — Document Model

- BSON/JSON
- Documents
- Collections
- Nested objects
- Arrays
- Dynamic schemas
- Document boundaries

## Chapter 5.2 — MongoDB Architecture

- `mongod`
- Replica sets
- Primary
- Secondaries
- Elections
- Journaling
- WiredTiger
- Storage engine basics

## Chapter 5.3 — MongoDB CRUD

- Insert
- Find
- Update
- Delete
- Operators
- Arrays
- Nested fields
- Projection

## Chapter 5.4 — Querying

- Filtering
- Sorting
- Pagination
- Aggregation
- `$match`
- `$group`
- `$lookup`
- `$unwind`
- `$project`

## Chapter 5.5 — MongoDB Indexes

- Single-field
- Compound
- Multikey
- Unique
- Sparse
- Partial
- TTL
- Geospatial
- Text/search
- Covered queries

## Chapter 5.6 — Query Planning

- Explain plans
- Collection scans
- Index scans
- Selectivity
- Index intersection
- Sort stages
- Query optimization

## Chapter 5.7 — Transactions

- Single-document atomicity
- Multi-document transactions
- Sessions
- Read concern
- Write concern
- Causal consistency

## Chapter 5.8 — MongoDB Sharding

- Sharded clusters
- Shard keys
- `mongos`
- Config servers
- Chunks
- Balancer
- Zone sharding

## Chapter 5.9 — Production MongoDB

- Connection pooling
- Timeouts
- Monitoring
- Backups
- Restore testing
- Security
- Authentication
- Authorization
- TLS
- Encryption

### Hands-On Project

**Production-style hotel reservation API**

Use:

- MongoDB
- .NET Web API
- Redis later
- Indexes
- Transactions
- Pagination
- Search
- Reservation expiration

This will also fit naturally with your .NET background.

---

# MODULE 6 — Key-Value Databases

## Chapter 6.1 — The Key-Value Model

```text
key → value
```

Understand why this simplicity is powerful.

## Chapter 6.2 — Redis

- Strings
- Hashes
- Lists
- Sets
- Sorted sets
- Streams
- Bitmaps
- HyperLogLog

## Chapter 6.3 — Redis Expiration

- TTL
- EXPIRE
- Key eviction
- Eviction policies
- Cache-aside

## Chapter 6.4 — Redis as a Cache

- Cache-aside
- Write-through
- Write-behind
- Read-through
- Cache invalidation
- Stampede
- Penetration
- Avalanche

## Chapter 6.5 — Redis Distributed Patterns

- Distributed locks
- Counters
- Rate limiting
- Leaderboards
- Sessions
- Pub/Sub
- Streams

## Chapter 6.6 — DynamoDB

- Partition key
- Sort key
- Item
- Table
- Query
- Scan
- GSIs
- LSIs
- Capacity units
- On-demand vs. provisioned

## Chapter 6.7 — DynamoDB Modeling

- Single-table design
- Composite keys
- Adjacency patterns
- Access patterns
- Sparse indexes
- Write sharding

### Interview Topics

- Redis vs. MongoDB
- Redis vs. DynamoDB
- Cache vs. database
- Why shouldn't Redis automatically become your primary database?
- What happens when Redis loses data?
- How do you prevent cache stampedes?
- Why is DynamoDB `Scan` usually dangerous?

### Hands-On Project

Build a **high-traffic API caching system**:

```text
Client
  ↓
.NET API
  ↓
Redis
  ↓ cache miss
MongoDB
```

Add:

- TTL
- Cache invalidation
- Rate limiting
- Distributed lock
- Metrics

---

# MODULE 7 — Wide-Column / Column-Family Databases

## Chapter 7.1 — Why Wide-Column Databases Exist

- Massive write throughput
- Distributed storage
- Large datasets
- Availability
- Predictable access patterns

## Chapter 7.2 — Cassandra Architecture

- Cluster
- Node
- Token
- Partition
- Replication factor
- Coordinator
- Gossip
- Snitch

## Chapter 7.3 — Data Model

- Keyspace
- Table
- Partition key
- Clustering columns
- Rows
- SSTables

## Chapter 7.4 — Cassandra Query Model

- CQL
- Partition-oriented queries
- Why arbitrary queries are problematic
- ALLOW FILTERING
- Query-first design

## Chapter 7.5 — Consistency

- ONE
- QUORUM
- LOCAL_QUORUM
- ALL
- Consistency trade-offs

## Chapter 7.6 — Cassandra Internals

- Commit log
- Memtable
- SSTable
- Compaction
- Bloom filters
- Tombstones
- Read path
- Write path

## Chapter 7.7 — ScyllaDB

- Architecture
- Shard-per-core
- Performance model
- Cassandra compatibility
- Operational differences

## Chapter 7.8 — Cassandra Failure Handling

- Hinted handoff
- Read repair
- Repair
- Anti-entropy
- Node replacement
- Bootstrap

### Interview Topics

This module contains **excellent senior-level interview material**.

- Why are Cassandra writes so fast?
- Why are reads sometimes expensive?
- Why can't Cassandra perform arbitrary JOINs?
- Why is the partition key so important?
- What is a tombstone?
- What happens when a node dies?
- What is compaction?
- Why can a bad partition key destroy performance?

### Hands-On Project

Build a **high-volume airport event ingestion system**:

```text
Flights
Aircraft
Gate events
Runway events
Passenger events
```

Optimize for:

- Millions/billions of events
- Time-based queries
- High write throughput
- Multi-node deployment

---

# MODULE 8 — Graph Databases

## Chapter 8.1 — Graph Fundamentals

- Nodes
- Relationships
- Properties
- Labels
- Directed graphs
- Weighted graphs

## Chapter 8.2 — Why Graph Databases?

Compare:

```text
SQL JOINs
      vs.
Graph traversal
```

## Chapter 8.3 — Neo4j

- Property graph model
- Nodes
- Relationships
- Properties
- Labels

## Chapter 8.4 — Cypher

- MATCH
- WHERE
- CREATE
- MERGE
- RETURN
- OPTIONAL MATCH
- Aggregations
- Paths

## Chapter 8.5 — Graph Algorithms

- Shortest path
- Centrality
- Community detection
- Similarity
- PageRank

## Chapter 8.6 — Graph Modeling

- Social networks
- Fraud
- Recommendations
- Knowledge graphs
- Dependency graphs

### Interview Topics

- When should you use a graph database?
- Why not simply use SQL?
- Why not represent a graph in MongoDB?
- What causes graph traversal performance problems?
- When does a graph database become unnecessary complexity?

### Hands-On Project

Build a **fraud detection graph**:

```text
Customer
   ↓
Account
   ↓
Transaction
   ↓
Device
   ↓
IP Address
```

Find suspicious relationships and transaction paths.

---

# MODULE 9 — Indexing Deep Dive

This deserves its own module because indexing knowledge transfers across almost every database.

## Chapter 9.1 — Why Indexes Exist

- Full scan
- Lookup
- Selectivity
- Cardinality

## Chapter 9.2 — B-Trees

- Structure
- Search
- Insert
- Delete
- Range queries

## Chapter 9.3 — Hash Indexes

- Equality lookup
- Hash collisions
- Why range queries don't work well

## Chapter 9.4 — Compound Indexes

- Field ordering
- Prefix rules
- Equality + range
- Sort optimization

## Chapter 9.5 — Specialized Indexes

- Geospatial
- TTL
- Text/search
- Partial
- Sparse
- Multikey

## Chapter 9.6 — Vector Indexes

- Embeddings
- Similarity search
- ANN
- HNSW
- Approximate vs. exact search
- Vector + metadata filtering

## Chapter 9.7 — Index Costs

Indexes aren't free.

- Memory
- Disk
- Write amplification
- Maintenance
- Build time
- Too many indexes

## Chapter 9.8 — Query Optimization

- Explain plans
- Selectivity
- Covered queries
- Sorts
- Scans
- Slow query analysis

### Hands-On Project

Take one dataset and deliberately create:

1. No index
2. Wrong index
3. Single-field index
4. Compound index
5. Covering index

Measure the differences.

---

# MODULE 10 — Consistency, Transactions & Concurrency

## Chapter 10.1 — ACID in Distributed Databases

- Atomicity
- Consistency
- Isolation
- Durability

## Chapter 10.2 — Isolation

- Read uncommitted
- Read committed
- Repeatable read
- Snapshot isolation
- Serializable

## Chapter 10.3 — Distributed Transactions

- Two-phase commit
- Coordinator
- Prepare
- Commit
- Failure scenarios

## Chapter 10.4 — Why Distributed Transactions Are Expensive

- Network round trips
- Coordination
- Locking
- Failure recovery

## Chapter 10.5 — Alternatives

- Saga pattern
- Compensating transactions
- Outbox pattern
- Idempotency
- Event-driven workflows

## Chapter 10.6 — Eventual Consistency in Practice

- Stale reads
- Read-your-writes
- Conflict windows
- User-visible inconsistency

### Interview Topics

- When should you avoid distributed transactions?
- Saga vs. 2PC
- What is idempotency?
- How do you guarantee a payment isn't processed twice?
- How do you handle partial failure?

---

# MODULE 11 — Multi-Region & Global Databases

## Chapter 11.1 — Why Multi-Region?

- Disaster recovery
- Latency
- Availability
- Data sovereignty
- Global users

## Chapter 11.2 — Replication Topologies

- Single-region
- Active-passive
- Active-active
- Multi-leader
- Leaderless

## Chapter 11.3 — Conflict Resolution

- Last-write-wins
- Version numbers
- Timestamps
- Vector clocks
- CRDTs
- Application-level conflict resolution

## Chapter 11.4 — Global Consistency

- Global strong consistency
- Regional consistency
- Eventual consistency
- Causal consistency

## Chapter 11.5 — Disaster Scenarios

- Region failure
- Network partition
- Replication lag
- Split brain
- Failover
- Failback

### Hands-On Project

Design a **global e-commerce system** with:

- Europe
- US
- Asia

Decide:

- Where data lives
- Who owns writes
- How replication works
- What happens during partition
- How conflicts resolve

---

# MODULE 12 — Reliability & Failure Engineering

## Chapter 12.1 — Failure Is Normal

- Hardware failures
- Process crashes
- Network failures
- Disk failures
- Human error
- Bad deployments

## Chapter 12.2 — Timeouts & Retries

- Retry storms
- Exponential backoff
- Jitter
- Retry budgets

## Chapter 12.3 — Idempotency

- Idempotency keys
- Duplicate requests
- At-least-once delivery

## Chapter 12.4 — Queues

- Kafka
- RabbitMQ
- SQS
- Consumer groups
- Ordering
- Delivery semantics

## Chapter 12.5 — Exactly-Once Semantics

- Why "exactly once" is difficult
- At-most-once
- At-least-once
- Effectively-once processing

---

# MODULE 13 — Backup, Recovery & Disaster Recovery

## Chapter 13.1 — Backups

- Full
- Incremental
- Differential
- Snapshot
- Continuous backup

## Chapter 13.2 — Recovery

- Restore
- Point-in-time recovery
- Recovery testing

## Chapter 13.3 — RPO & RTO

- Recovery Point Objective
- Recovery Time Objective
- Availability targets

## Chapter 13.4 — Disaster Recovery Strategies

- Backup/restore
- Warm standby
- Hot standby
- Multi-region

### Interview Topics

> "Your production database was accidentally deleted. What do you do?"

You should eventually be able to answer this as an operational procedure rather than simply saying "restore the backup."

---

# MODULE 14 — Security

## Chapter 14.1 — Authentication

- Passwords
- Certificates
- IAM
- Service identities

## Chapter 14.2 — Authorization

- RBAC
- Least privilege
- Roles
- Permissions

## Chapter 14.3 — Encryption

- At rest
- In transit
- Key management
- Rotation

## Chapter 14.4 — Database Network Security

- Private networks
- Firewalls
- Security groups
- TLS
- Network segmentation

## Chapter 14.5 — Data Security

- PII
- GDPR considerations
- Data masking
- Auditing
- Retention
- Deletion

---

# MODULE 15 — Observability & Production Operations

## Chapter 15.1 — Database Metrics

- Latency
- Throughput
- CPU
- Memory
- Disk
- IOPS
- Connections
- Cache hit rate
- Replication lag

## Chapter 15.2 — Query Monitoring

- Slow queries
- Query frequency
- Query plans
- Failed queries

## Chapter 15.3 — Distributed-System Metrics

- Partition imbalance
- Hot shards
- Node health
- Replication lag
- Compaction
- Tombstones

## Chapter 15.4 — Logging

- Structured logs
- Correlation IDs
- Audit logs

## Chapter 15.5 — Tracing

- Distributed tracing
- Database spans
- End-to-end latency

## Chapter 15.6 — Alerting

- What deserves an alert?
- Thresholds
- Symptoms vs. causes
- Alert fatigue

### Hands-On Project

Instrument your previous applications with:

- OpenTelemetry
- Metrics
- Logs
- Traces
- Database monitoring

Then intentionally create:

- Slow query
- Connection exhaustion
- Cache failure
- Replication lag

Diagnose each one.

---

# MODULE 16 — Performance Engineering

## Chapter 16.1 — Latency

- Average
- p50
- p95
- p99
- Tail latency

## Chapter 16.2 — Throughput

- Requests/sec
- Writes/sec
- Reads/sec
- Saturation

## Chapter 16.3 — Bottlenecks

- CPU
- Memory
- Disk
- Network
- Lock contention
- Connection pools

## Chapter 16.4 — Load Testing

- Workload generation
- Realistic distributions
- Warm-up
- Steady state
- Failure testing

## Chapter 16.5 — Capacity Planning

- Current traffic
- Growth
- Headroom
- Scaling thresholds
- Cost

---

# MODULE 17 — Cloud NoSQL

## Chapter 17.1 — AWS

- DynamoDB
- ElastiCache
- DocumentDB
- Key design
- Capacity
- Global tables

## Chapter 17.2 — Azure

- Cosmos DB
- Partitioning
- Consistency levels
- Global distribution

## Chapter 17.3 — GCP

- Firestore
- Bigtable
- Datastore concepts

## Chapter 17.4 — Managed vs. Self-Hosted

- Operational burden
- Cost
- Flexibility
- Vendor lock-in
- SLA
- Backup
- Upgrades

---

# MODULE 18 — Polyglot Persistence

A production application rarely has to choose **one database for everything**.

## Chapter 18.1 — Multiple Databases

Example:

```text
                 Application
                     |
       +-------------+-------------+
       |             |             |
   PostgreSQL      Redis        MongoDB
       |                           |
 transactions                 documents
```

## Chapter 18.2 — Choosing the Right Database Per Workload

- Transactional data
- Cache
- Search
- Analytics
- Relationships
- Event streams
- Time series
- Object storage

## Chapter 18.3 — Keeping Multiple Systems Consistent

- CDC
- Outbox
- Events
- Change streams
- ETL/ELT

### Hands-On Project

Build a **production-style booking platform** using:

- PostgreSQL → transactional source of truth
- MongoDB → document-oriented operational data
- Redis → cache/session/rate limiting
- Elasticsearch/OpenSearch → search
- Kafka → events

The important part isn't simply making them work.

It's answering:

> **Why does each piece exist?**

---

# MODULE 19 — Advanced Architecture

## Chapter 19.1 — Distributed Data Architecture

- Data ownership
- Service boundaries
- Database-per-service
- Shared databases
- Data duplication

## Chapter 19.2 — CQRS

- Commands
- Queries
- Read models
- Write models
- Eventual consistency

## Chapter 19.3 — Event Sourcing

- Events as source of truth
- Event replay
- Snapshots
- Schema evolution

## Chapter 19.4 — Change Data Capture

- CDC
- Database logs
- Debezium
- Streaming changes

## Chapter 19.5 — Materialized Views

- Precomputed data
- Read optimization
- Refresh strategies

## Chapter 19.6 — Time-Series Databases

Additional NoSQL family worth knowing:

- InfluxDB
- TimescaleDB
- Metrics
- Events
- Retention
- Downsampling

## Chapter 19.7 — Search Engines

Another important adjacent category:

- Elasticsearch
- OpenSearch
- Inverted indexes
- Full-text search
- Relevance
- Faceting

## Chapter 19.8 — Vector Databases

- Embeddings
- Similarity search
- ANN
- HNSW
- RAG
- Hybrid search
- Metadata filtering

---

# MODULE 20 — SQL vs. NoSQL Decision Framework

You should eventually be able to answer this without saying:

> "It depends."

Instead, evaluate systematically.

## Chapter 20.1 — Workload

- Read/write ratio
- Query patterns
- Transaction complexity
- Data relationships
- Dataset size

## Chapter 20.2 — Consistency

- Strong consistency required?
- Eventual consistency acceptable?
- Cross-record transactions?

## Chapter 20.3 — Scaling

- Vertical scaling sufficient?
- Horizontal scaling required?
- Global distribution?

## Chapter 20.4 — Operational Requirements

- Team expertise
- Managed services
- Backup requirements
- SLA
- Cost

## Chapter 20.5 — Decision Matrix

Eventually compare:

| Requirement | PostgreSQL | MongoDB | DynamoDB | Cassandra | Redis | Neo4j |
|---|---|---|---|---|---|---|
| Complex transactions | ★★★★★ | ★★★ | ★★ | ★ | ★ | ★★ |
| Flexible documents | ★★★ | ★★★★★ | ★★★★ | ★★ | ★★ | ★★ |
| Massive distributed writes | ★★★ | ★★★★ | ★★★★★ | ★★★★★ | ★★★★★ | ★★ |
| Graph traversal | ★★ | ★★ | ★ | ★ | ★ | ★★★★★ |
| Caching | ★★ | ★★ | ★★ | ★ | ★★★★★ | ★ |
| Global distribution | ★★★ | ★★★★ | ★★★★★ | ★★★★★ | ★★★★ | ★★★ |

The important lesson will be understanding **why** each score exists rather than memorizing the table.

---

# MODULE 21 — NoSQL Interview Mastery

This will run **throughout the curriculum**, rather than being left until the end.

## Level 1 — Fundamentals

- What is NoSQL?
- SQL vs. NoSQL
- Document vs. key-value
- What is denormalization?
- What is eventual consistency?
- What is CAP?

## Level 2 — Practical Developer

- Embed vs. reference
- Choosing indexes
- Pagination
- TTL
- Transactions
- Redis caching
- MongoDB aggregation
- DynamoDB partition keys

## Level 3 — Senior Developer

- Design a scalable data model
- Choose a shard key
- Diagnose a slow query
- Handle replication lag
- Design idempotent writes
- Handle cache failure
- Design retry logic

## Level 4 — Staff/Principal

- Design a globally distributed database
- Active-active architecture
- Multi-region conflict resolution
- Hot-partition mitigation
- Disaster recovery
- Capacity planning
- Consistency guarantees
- Cost/performance trade-offs

## Level 5 — "Production Incident" Questions

These are particularly valuable because documentation rarely teaches them.

Examples:

> Production latency suddenly went from 50 ms to 3 seconds. How do you investigate?

> One partition is receiving 70% of the traffic. What happened?

> Replicas are falling increasingly behind. What do you check?

> Your cache is down. Does your system survive?

> A user submits the same payment request three times. What happens?

> Your database is healthy, but the application is timing out. What do you investigate?

> Your primary database region goes down. Walk me through recovery.

> Your MongoDB cluster has 20 indexes and write performance collapsed. What do you do?

> Cassandra has millions of tombstones. Why is this happening?

> A DynamoDB table suddenly gets throttled. What could cause it?

These will become part of the curriculum rather than being treated as a separate interview-cramming section.

---

# MODULE 22 — Capstone Architecture Challenges

Finally, stop learning databases individually and start designing systems.

## Capstone 1 — URL Shortener

Design:

- MongoDB/DynamoDB
- Cache
- IDs
- Partitioning
- Analytics

Focus:

- Key design
- Read-heavy workload
- Hot keys

---

## Capstone 2 — Social Network

Design:

- Users
- Followers
- Posts
- Feed
- Notifications

Compare:

- MongoDB
- Cassandra
- DynamoDB
- Neo4j

---

## Capstone 3 — Hotel Booking Platform

Design:

```text
Search
Availability
Pricing
Reservations
Payments
Notifications
```

Address:

- Transactions
- Concurrency
- Double booking
- Idempotency
- Cache
- Search

---

## Capstone 4 — Global E-Commerce Platform

Address:

- Multi-region
- Replication
- Inventory
- Orders
- Payments
- Eventual consistency
- Disaster recovery

---

## Capstone 5 — High-Scale Event Platform

Design for:

> 100 million events/day.

Address:

- Cassandra/ScyllaDB
- Kafka
- Partitioning
- Retention
- Compaction
- Monitoring
- Backups

---

# The Learning Sequence

I recommend we **do not simply march through the database products one after another**.

The dependency graph should be:

```text
                    NoSQL
                      │
             ┌────────┴────────┐
             │                 │
       Distributed Systems   Data Models
             │                 │
       CAP / BASE / PACELC     │
             │                 │
       Replication             │
       Partitioning            │
       Consistency             │
             └────────┬────────┘
                      │
                Data Modeling
                      │
        ┌─────────────┼─────────────┐
        │             │             │
     Document      Key-Value    Wide-Column
        │             │             │
     MongoDB       Redis       Cassandra
        │             │             │
        └─────────────┼─────────────┘
                      │
                    Graph
                      │
                    Neo4j
                      │
              Indexing & Performance
                      │
           Transactions & Consistency
                      │
            Multi-region / DR / Security
                      │
                Observability
                      │
             Advanced Architecture
                      │
               Capstone Systems
```

### Recommended major stages

**Stage 1:** Modules 0–3 — Foundations & distributed systems  
**Stage 2:** Module 4 — NoSQL data modeling  
**Stage 3:** Module 5 — MongoDB  
**Stage 4:** Module 6 — Redis + DynamoDB  
**Stage 5:** Module 7 — Cassandra + ScyllaDB  
**Stage 6:** Module 8 — Neo4j  
**Stage 7:** Module 9 — Indexing & query performance  
**Stage 8:** Modules 10–13 — consistency, transactions, reliability, DR  
**Stage 9:** Modules 14–17 — security, observability, performance, cloud  
**Stage 10:** Modules 18–19 — polyglot persistence & advanced architecture  
**Stage 11:** Module 20 — SQL vs. NoSQL decisions  
**Stage 12:** Modules 21–22 — senior/principal interviews + architecture

The key principle I'll follow if we work through this as lessons is: **we won't treat "knowing MongoDB commands" as knowing NoSQL.** You should be able to look at a workload, derive its access patterns, choose a data model, choose partition/replication/consistency strategies, predict failure modes, and defend the trade-offs in a senior-level interview.