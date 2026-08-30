# 1. SQL / Relational Database

**Choose when:**
- Data has clear relationships.
- You need joins.
- You need strong transactional guarantees.
- Data integrity and constraints are important.
- The schema is reasonably well-defined.
- You need complex queries across related entities.

**Examples:** SQL Server, PostgreSQL, MySQL, Oracle.

**Trade-offs:**
- Excellent consistency and relational querying.
- Scaling vertically is straightforward, but horizontal scaling can become more complicated.
- Schema changes can require more coordination.
- Very large, highly distributed workloads may require additional architecture.

**Typical example:**

Users, properties, bookings, payments:

```text
User
  ↓
Property
  ↓
Booking
  ↓
Payment
```

Those relationships make SQL a natural starting point.

---

# 2. NoSQL — Document Database

**Choose when:**
- Data naturally fits into documents.
- Schema needs to evolve frequently.
- You don't need many relational joins.
- You have very high scale or distributed workloads.
- Your access patterns are relatively well understood.

**Examples:** MongoDB, Cosmos DB.

**Trade-offs:**
- Flexible schema.
- Can scale horizontally very well.
- Often requires denormalization.
- You need to think carefully about partition keys and access patterns.
- Complex relationships/joins are generally less natural than in SQL.

**Example:**

A saved search:

```text
{
    userId,
    location,
    priceRange,
    propertyType,
    rooms,
    notificationPreferences,
    ...
}
```

This is relatively self-contained and can evolve as new search criteria are introduced.

---

# 3. NoSQL — Key-Value

**Choose when:**
- You primarily retrieve something using a key.
- You need extremely fast lookups.
- You need massive horizontal scale.
- The data doesn't require complex relationships.

**Examples:** DynamoDB, Redis.

Think:

> "Given this key, give me the value."

**Trade-offs:**
- Extremely efficient for the right access pattern.
- Excellent scalability.
- Poor fit for arbitrary relational queries.

---

# 4. NoSQL — Wide Column

**Choose when:**
- You have enormous datasets.
- Very high write/read throughput is required.
- Data can be distributed across many machines.
- Access patterns are known beforehand.

**Examples:** Cassandra, ScyllaDB.

**Trade-off:**

You get excellent distributed scalability, but you generally design the data model around your queries rather than expecting arbitrary querying like SQL.

---

# 5. NoSQL — Graph

**Choose when relationships are the actual thing you're querying.**

Examples:

- Social networks
- Recommendation systems
- Fraud detection
- Dependency graphs

**Example:**

> "Find people connected to Alice through two or three relationships."

**Example:** Neo4j.

---

# 6. Database Replication

**Choose when:**
- You need higher availability.
- You need additional read capacity.
- You need redundancy if a database server fails.

A common model:

```text
              Primary
             /       \
            ↓         ↓
       Replica 1   Replica 2
```

Writes → primary.

Reads → potentially replicas.

**Trade-offs:**
- More availability and read capacity.
- Additional infrastructure.
- Replication lag can mean stale reads.

**Important:** Replication doesn't magically solve a database performance problem. You first need to understand whether the bottleneck is reads, writes, storage, CPU, locking, or queries.

---

# 7. Database Sharding / Partitioning

**Choose when:**
- One database instance can no longer handle the data or workload.
- You need to distribute data across multiple machines.

For example:

```text
Users 1–10M       → DB 1
Users 10M–20M     → DB 2
Users 20M–30M     → DB 3
```

**Trade-offs:**
- Massive horizontal scalability.
- More operational complexity.
- Cross-shard queries and transactions become difficult.
- Choosing a poor partition key can create hot partitions.

---

# 8. Cache

**Choose when:**
- The same data is read repeatedly.
- The underlying operation is expensive.
- The data doesn't change constantly.
- Low latency is important.

Example:

> Frequently requested property details.

Instead of hitting the DB every time:

```text
Request
  ↓
Cache
  ↓
DB only when cache misses
```

**Trade-offs:**
- Much faster reads.
- Reduces database load.
- Introduces stale-data/invalidation problems.
- Additional infrastructure.

---

# 9. Distributed Cache

**Choose when:**
- You have multiple application servers.
- Cached data needs to be shared between them.

For example:

```text
             Redis
            /     \
           ↓       ↓
       API 1     API 2
```

Without distributed caching:

```text
API 1 → its own cache
API 2 → different cache
```

You could have the same data cached multiple times.

**Trade-offs:**
- Shared cache across servers.
- More infrastructure.
- Network dependency.
- Cache failure needs to be handled.

**Typical choice:** Redis.

---

# 10. Load Balancer

**Choose when:**
- You have multiple application instances.
- You want to distribute traffic.
- You want availability if one instance fails.

```text
              Load Balancer
             /      |      \
            ↓       ↓       ↓
          API 1   API 2   API 3
```

**Trade-offs:**
- Better scalability.
- Better availability.
- Additional infrastructure.
- Your application generally needs to be stateless or otherwise handle distributed state correctly.

---

# 11. Message Queue

**Choose when:**
- Work doesn't need to happen during the request.
- You want to decouple services.
- You need buffering during traffic spikes.
- You need retry capabilities.

Examples:

- RabbitMQ
- Amazon SQS
- Azure Service Bus

Example:

> User uploads a property → API stores it → queue contains "process property images" → worker processes them.

**Trade-offs:**
- Better resilience and decoupling.
- Introduces asynchronous behavior.
- Eventual consistency.
- Requires handling duplicates, retries, ordering, and failed messages.

---

# 12. Kafka / Event Streaming

**Choose when:**
- You have high-volume streams of events.
- Multiple consumers need the events.
- You want to retain events.
- Consumers may need to replay events.
- Ordering/partitioning matters.

Example:

```text
Listing Created
       ↓
     Kafka
   /    |     \
  ↓     ↓      ↓
Search  Alerts Analytics
```

**Trade-offs:**
- Extremely powerful for event-driven architectures.
- High throughput.
- Replayable events.
- More complex than a simple queue.

### Kafka vs RabbitMQ

Think:

**RabbitMQ:**  
> "Here is some work. One worker should process it."

**Kafka:**  
> "Here is an event. Multiple consumers may independently consume it, and the event can remain available for later replay."

That's a simplification, but a useful starting point.

---

# 13. Search Engine

**Choose when:**
- Search is a major feature.
- Full-text search is required.
- You need relevance ranking.
- Fuzzy matching/autocomplete is important.
- Complex filtering or geographic search is required.

**Examples:** Elasticsearch, OpenSearch.

**Trade-offs:**
- Excellent search capabilities.
- Can dramatically reduce search workload on your primary database.
- Another system whose data must be synchronized with the source database.
- Usually eventually consistent with the primary database.

**Important:** It usually isn't your primary source of truth.

---

# 14. Object Storage

**Choose when:**
- You're storing large files/blobs.
- Images, videos, PDFs, etc.

Examples:

- Amazon S3
- Azure Blob Storage

Instead of:

```text
SQL Database
   ↓
10 MB property image
```

you'd typically store:

```text
Object Storage → image
Database       → URL/reference to image
```

**Trade-offs:**
- Cheap and scalable storage.
- Excellent for large files.
- Not designed for relational querying.

---

# 15. Rate Limiting

**Choose when:**
- You need to protect an API.
- An external provider has request limits.
- You want to prevent one client from consuming disproportionate resources.

Example:

> 100 requests/minute per user.

**Trade-offs:**
- Protects your system.
- Restricts legitimate traffic too if configured badly.
- Distributed rate limiting requires shared state such as Redis or another distributed mechanism.

---

# 16. Retry + Exponential Backoff

**Choose when:**
- Failure is potentially temporary.
- External services occasionally fail.
- Network failures occur.

Instead of:

> "Request failed → immediately try again 10 times."

You progressively increase the delay.

**Trade-offs:**
- Helps recover from transient failures.
- Too many retries can actually make an outage worse.
- Requires maximum retries and usually a dead-letter/failure strategy.

---

# 17. Dead-Letter Queue

**Choose when:**
- A message has repeatedly failed processing.

Instead of endlessly retrying:

```text
Queue
 ↓
Worker
 ↓
fail
 ↓
retry
 ↓
retry
 ↓
retry
 ↓
Dead-letter queue
```

**Trade-off:**

You don't lose the message, but someone/system eventually needs to inspect or reprocess it.

---

# 18. Horizontal vs Vertical Scaling

### Vertical scaling

Make one machine bigger.

> More CPU, RAM, storage.

**Choose when:**
- The workload fits comfortably on one machine.
- You want simplicity.

**Trade-off:**
- Eventually you hit hardware limits.
- Doesn't inherently give you redundancy.

### Horizontal scaling

Add more machines.

> 1 server → 5 servers → 50 servers.

**Choose when:**
- You need significant scalability.
- You need redundancy/availability.

**Trade-off:**
- Distributed-system complexity.

---

# 19. CDN

**Choose when:**
- You serve lots of static content.
- Users are geographically distributed.
- You want content closer to users.

Examples:

- CloudFront
- Cloudflare
- Azure Front Door

Excellent for:

- Images
- CSS/JS
- Videos
- Static assets

**Trade-off:**
- Cached content can become stale.
- Cache invalidation becomes a consideration.

---

# 20. Observability

Not a scaling mechanism, but essential once systems become complicated.

### Metrics

> "How is the system performing?"

Latency, throughput, error rate, CPU, queue depth, cache hit rate.

### Logs

> "What happened?"

Exceptions, failed requests, provider failures.

### Traces

> "Where did this particular request spend its time?"

Useful when a request travels through several services.

---

# The decision framework I'd actually use in an interview

When presented with a system-design problem, don't think:

> "Which technologies can I mention?"

Think:

### **Database**

Ask:

> Do I have relationships and strong transactional requirements?

**Yes → SQL is a strong candidate.**

> Is the data naturally document/key-value shaped, highly distributed, and driven by predictable access patterns?

**Yes → consider NoSQL.**

Then ask:

> Is the problem actually reads?

**Yes → optimize queries/indexes → consider cache/read replicas.**

> Is one database becoming too large?

**Consider partitioning/sharding.**

> Is search itself becoming complex?

**Consider a search engine.**

### **Traffic**

> More requests than one server can handle?

**Horizontal scaling + load balancer.**

### **Slow/independent work**

> Doesn't need to finish during the HTTP request?

**Queue/background processing.**

### **Huge event stream / multiple consumers / replay**

**Kafka/event streaming.**

### **Repeated expensive reads**

**Cache.**

### **Multiple application servers need the same cache**

**Distributed cache.**

### **Large files**

**Object storage.**

### **External dependency**

**Timeout + retry + exponential backoff + rate limiting + circuit breaker where appropriate.**

### **Failures**

> What happens if this component goes down?

That's the question that leads you toward **replication, redundancy, queues, fallbacks, and graceful degradation**.

---

## The most important SQL vs NoSQL distinction

If you remember only one thing, I'd make it this:

**SQL:**

> "I have data with relationships and need flexible querying, strong transactions, and integrity."

**NoSQL:**

> "I have a very large/distributed workload and I know the important ways the application will access the data, so I can optimize the data model around those access patterns."

And importantly:

**Don't choose NoSQL simply because "it's faster" or SQL because "it's more reliable."** Both categories have multiple products with very different characteristics. The requirements and workload determine the appropriate choice.