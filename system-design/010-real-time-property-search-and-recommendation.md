# Real-Time Property Search & Recommendation System

**Interviewer:** Design a real-time property search and recommendation system.

**Me:**

> **Before I design it, I'd like to clarify a few requirements.**

#### 1. Functional requirements

I'd like the system to support:

- Searching active property listings.
- Filtering by location, price range, property type, number of rooms, amenities, etc.
- Geospatial search, such as properties within a radius or map bounding box.
- Sorting by relevance, price, or newest.
- Pagination through potentially millions of listings.
- Showing newly created or edited listings in search relatively quickly.
- Personalized recommendations based on things like a user's searches, favorites, and viewed properties.

For recommendations, I'll initially assume we can provide reasonable content-based or preference-based recommendations rather than requiring a sophisticated ML model from day one.

#### 2. Non-functional requirements

I'd expect this to be a **very read-heavy system**, since many users will search much more frequently than agencies create or modify listings.

I'd target:

- Low search latency, say p95 below 200–300 ms.
- High availability.
- Millions of listings.
- High concurrent search traffic.
- Search freshness of perhaps less than a minute for new or updated listings.
- Eventual consistency is acceptable for search results, but the primary property data must remain strongly consistent.

(I'd specifically clarify that search does not need to be strongly consistent because that gives us considerably more flexibility for scaling the read path.)

### 3. What data do we need?

The source-of-truth data would contain things such as:

**Listing**

- `id`
- `property_id`
- `agency_id`
- `price`
- `status`
- `property_type`
- `rooms`
- `description`
- `created_at`
- `updated_at`

**Property**

- address
- latitude/longitude
- size
- amenities

And for recommendations I'd eventually want user behavior such as:

- searches
- viewed listings
- favorites
- contacted agencies

I'd keep timestamps such as `created_at` and `updated_at` because freshness and ordering matter here.

### 4. Database choice

For the source of truth, I'd choose a **relational database**.

There are relationships between users, agencies, properties, and listings, and we need reliable updates to listing state.

However, I wouldn't make the relational database handle the entire search workload.

The search requirements are quite different: full-text search, many filters, geospatial queries, relevance ranking, and potentially very high read volume.

So I'd introduce a **dedicated search engine/index**, such as Elasticsearch or OpenSearch.

The important distinction is:

> **SQL remains the source of truth. The search index is a read-optimized representation of that data.**

### 5. APIs

I'd expose something like:

```text
GET /listings/search
    ?location=Berlin
    &minPrice=500
    &maxPrice=2000
    &rooms=2
    &propertyType=APARTMENT
    &sort=relevance
    &page=...

GET /listings/{listingId}

GET /users/me/recommendations
```

For behavioral events, we could have something like:

```text
POST /users/me/events
```

where the event might represent a listing view, favorite, search, etc.

### 6. Start with the simplest architecture

I'd initially keep the architecture straightforward:

```text
Client
   |
   v
  API
   |
   +----> SQL Database
   |
   +----> Search Index
```

The API handles the request, while the search index handles search-specific queries.

(I don't want to introduce Redis, Kafka, Kubernetes, microservices, etc. just because this is a large system. I'll introduce additional components only when a requirement justifies them.)

---

### 7. How does data get into the search index?

This is where the interesting part starts.

Suppose an agency creates a listing.

I want the SQL transaction to be the authoritative operation. After the listing is successfully stored, the change needs to reach the search index.

I would **not** simply do:

```text
write SQL
write search index
```

because if SQL succeeds but the application crashes before updating the search index, the two systems become inconsistent.

I'd use a **transactional outbox**.

The transaction writes:

1. The listing.
2. An outbox event describing the change.

A background worker then processes those events and updates the search index.

```text
                 +----------------+
                 | SQL transaction|
                 +-------+--------+
                         |
                 listing + outbox
                         |
                         v
                      Queue
                         |
                         v
                  Indexing Worker
                         |
                         v
                    Search Index
```

The index may therefore be a few seconds behind SQL, but that's acceptable under our freshness requirement.

I'd also include the listing version or `updated_at` in the event so that an older event cannot overwrite a newer version in the search index.

### 8. Search performance

The search engine is now handling:

- text search
- price filters
- property-type filters
- amenities
- geo-radius/bounding-box queries
- relevance ranking

For very large result sets, I'd avoid deep offset pagination where possible and use a cursor/search-after style approach.

(I'd expect this to matter once someone asks, "What happens when the user requests page 10,000?")

### 9. Caching

Search is read-heavy, so caching is useful.

I'd cache particularly popular searches and potentially listing details.

For example, a cache key could incorporate the normalized search criteria, sorting and pagination cursor.

I wouldn't blindly cache every possible query because the number of possible combinations is enormous.

Once we have multiple API servers, I'd use a **distributed cache** rather than an in-memory cache on each server.

The cache can tolerate some staleness because the search system itself is already eventually consistent.

### 10. Scaling

At this point, if traffic increases, I can scale the API horizontally:

```text
                 Load Balancer
                /      |      \
              API     API     API
               |       |       |
               +-------+-------+
                       |
               Search / SQL
```

The search infrastructure can scale independently from the API.

The search index can use:

- multiple nodes
- shards for distributing data
- replicas for read capacity and availability

If SQL becomes a bottleneck for listing details, I could introduce read replicas where appropriate.

(The important architectural point is that increasing search traffic shouldn't force me to scale the primary SQL database proportionally, because search traffic is going to the search infrastructure.)

---

### 11. Recommendations

I'd initially keep recommendations relatively simple.

For a user who frequently searches for:

> 2-bedroom apartments in Berlin under $2,000

I can combine the user's preferences with listing attributes and produce relevant candidates.

Later, when we have enough behavioral data, I'd introduce an ML-based recommendation system.

User interactions such as views, searches and favorites can be processed asynchronously. A recommendation service can periodically generate recommendations, which we can cache and serve quickly.

So recommendations don't need to sit synchronously in the critical search path.

---

### 12. Reliability

If the search index is temporarily unavailable, I wouldn't immediately start sending millions of arbitrary search queries to SQL as a fallback—that could take down the primary database.

I'd prefer controlled degradation:

- serve cached popular results where possible
- retry indexing operations
- monitor index health
- potentially provide a limited fallback search capability

For indexing failures, I'd use retries with exponential backoff and a dead-letter queue for events that repeatedly fail.

The original listing remains safe because SQL is still the source of truth.

---

### 13. Observability

I'd monitor:

- search p50/p95/p99 latency
- search error rate
- zero-result searches
- search-index freshness/lag
- indexing failures
- queue depth
- cache hit rate
- SQL latency
- recommendation latency
- recommendation click-through rate

I'd also use structured logs and distributed tracing so that I can follow a request across the API, cache, search service and database.

---

### Final architecture

So my final design would evolve into:

```text
                         Client
                           |
                     Load Balancer
                           |
                    Multiple API Servers
                     /        |        \
                    /         |         \
               Cache       Search     SQL DB
                             Index        |
                                         |
                                  Transactional Outbox
                                         |
                                       Queue
                                         |
                                  Indexing Workers
                                         |
                                         v
                                    Search Index
                                         
User Events
     |
     v
   Queue
     |
     v
Recommendation Processing
     |
     v
Recommendation Cache
```

(I'd emphasize that the diagram is showing logical components rather than saying each one must be a separate microservice.)

**The key design decisions are:**

1. **SQL is the source of truth.**
2. **A dedicated search index handles high-volume search, filtering and geospatial queries.**
3. **Transactional outbox + asynchronous indexing keeps SQL and search eventually consistent without unsafe dual writes.**
4. **Caching handles popular read-heavy workloads.**
5. **API servers and search infrastructure scale independently.**
6. **Recommendations are initially simple and can evolve toward ML without changing the core search architecture.**
7. **Queues, retries and DLQs make asynchronous processing resilient.**

---

### Interview follow-ups

**Interviewer:** *What happens if a listing is edited five times very quickly?*

I'd include a version or `updated_at` with every indexing event. The indexing worker only applies an event if it represents a newer version than what's currently indexed. That prevents an older asynchronous event from overwriting newer data.

**Interviewer:** *What if the search index says a property is active but SQL says it's already sold?*

Search results are allowed to be eventually consistent. When the user opens the property, the authoritative SQL-backed service can verify its current status. We can also prioritize rapid propagation of status changes.

**Interviewer:** *What if search traffic suddenly increases 10x?*

I'd first scale the search cluster and API servers horizontally, increase cache utilization for hot queries, and monitor whether the search cluster or database is the bottleneck. I wouldn't immediately scale everything.

**Interviewer:** *Would you use Kafka?*

Potentially, yes, if event volume becomes large enough to justify a durable distributed event stream—for example, millions of listing updates and user interaction events. For a smaller system, a simpler queue may be sufficient.

**Interviewer:** *Would you use SQL for the search?*

For simple filtering, potentially yes. But once we require full-text relevance, complex filtering, geospatial search and very high read volume, I'd move that workload to a specialized search engine.

**Interviewer:** *How would you rank results?*

I'd start with search relevance plus business signals such as freshness and property quality. Later, personalized signals from user behavior could be incorporated into the ranking model.

**Interviewer:** *How would you make geospatial search efficient?*

I'd store the property's latitude/longitude as a geo-point in the search index and let the search engine execute radius or bounding-box queries using its spatial indexing capabilities.

**Interviewer:** *How do you distinguish availability from scalability here?*

Multiple API servers and search replicas help both, but they're solving different problems. **Scalability** is the ability to handle increasing traffic by adding resources. **Availability** is keeping the system operational despite failures, such as losing an API server or search node.