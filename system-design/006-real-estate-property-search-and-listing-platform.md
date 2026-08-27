Perfect. We'll do **Real Estate Property Search & Listing Platform** and follow your exact framework, step by step.

We’ll stay in **interview mode**: I’ll phrase things as you would actually say them to the interviewer, rather than turning it into a lesson. We won't jump ahead to Redis/Kafka/microservices until there's a requirement that justifies them.

## Step 1 — Clarify requirements

I'd start the interview like this:

> **"Before I start designing the system, I'd like to clarify a few requirements."**

### Functional requirements

> **"What are the main operations users should be able to perform?"**

I'd expect the interviewer to establish something along these lines:

- Users can **create and manage property listings**.
- Users can **search for properties**.
- They can filter by things such as:
  - location
  - price range
  - property type
  - number of rooms
- Users can **view the details of a listing**.
- Users can **save/favorite listings**.

For the first version, I would **not** add things like recommendations, messaging, payments, image processing, maps, etc. unless the interviewer introduces them.

Then I'd clarify:

> **"Who can create listings? Is it any registered user, or are there separate roles such as landlords, agents, and regular buyers?"**

And:

> **"Should users be able to edit and delete their own listings?"**

Those questions establish authorization and ownership without prematurely designing an authentication service.

---

### Non-functional requirements

Then I'd ask:

> **"What are the expected traffic and scale? For example, approximately how many users, listings, and searches per second should we support?"**

Then:

> **"Are there any latency requirements for search and listing details?"**

For example, if the interviewer says:

> "Search should normally respond within 200 ms."

Now we have something concrete that can influence later decisions.

I'd also ask:

> **"How important is availability? Can the search functionality tolerate a short outage, or do we need it to be highly available?"**

And:

> **"How fresh does the search data need to be? If someone creates a listing, does it need to appear in search immediately, or is a delay of a few seconds acceptable?"**

That last question is particularly useful because it can eventually influence whether we can use asynchronous processing, caching, search indexes, etc.

Finally:

> **"Are there any specific consistency or durability requirements for listings?"**

For example, we probably don't want a successfully created property listing to simply disappear.

---

### What we have established so far

For our exercise, let's assume the interviewer gives us:

- Registered users can create, edit, and delete their own listings.
- Users can search and filter listings.
- Users can view listing details.
- Users can favorite listings.
- Search is **very read-heavy** compared with listing creation.
- Search should normally respond within roughly **200 ms**.
- Listing changes should eventually appear in search; a small delay is acceptable.
- Listing data must be durable.
- We expect **millions of listings** and significantly more search requests than listing writes.
- The system should be highly available.

Got it. I'll stay **in interviewee mode** from here. Anything I want you to notice about my reasoning will be in `(parentheses)`. I won't stop after each step; I'll take the system all the way through, including likely interviewer follow-ups and how I'd respond.

# Real Estate Property Search & Listing Platform

### 1. Clarify requirements

> **"Before I design it, I'd like to clarify a few requirements."**

> "From a functional perspective, I understand that users should be able to create, update, and delete their own property listings. Other users should be able to search and filter listings, view listing details, and save listings as favorites."

> "For the search functionality, what filters should we support? I'm assuming things like location, price range, property type, number of rooms, and possibly availability."

> "Who can create listings? Is it any registered user, or do we have different roles such as landlords and real-estate agents?"

> "Should a listing become immediately visible in search after it's created, or is eventual consistency acceptable?"

> "What happens when a listing is updated or deleted? Should those changes also be reflected in search immediately?"

> "For non-functional requirements, what kind of traffic are we expecting? In particular, I'm interested in the number of listings, concurrent users, and search requests per second."

> "Do we have a latency target for search?"

> "How important is availability? Should search continue working if some other part of the system is unavailable?"

> "And are there any specific security or authorization requirements?"

(I'm deliberately asking questions that can actually affect the architecture. I'm not asking for every possible requirement just to demonstrate that I know the questions.)

Let's assume the interviewer gives me:

- Millions of listings.
- Much more searching than listing creation/update.
- Search should normally respond within ~200 ms.
- A few seconds of delay before a new/updated listing appears in search is acceptable.
- Listings must be durable.
- The system should be highly available.
- Users can create and manage their own listings.
- Users can search, view, and favorite listings.

---

# 2. Determine the data

> "I'll first identify the main entities we need to persist."

I'd expect something like:

### User

```text
User
- Id
- Name
- Email
- CreatedAt
```

### Listing

```text
Listing
- Id
- UserId
- Title
- Description
- Price
- PropertyType
- NumberOfRooms
- Area
- Address
- City
- PostalCode
- Status
- CreatedAt
- UpdatedAt
- ExpiresAt
```

(I'd also think about fields such as `CreatedAt`, `UpdatedAt`, `ExpiresAt`, ownership, and status because these often become important for filtering, lifecycle management, auditing, and cleanup.)

### Favorite

```text
Favorite
- UserId
- ListingId
- CreatedAt
```

(I'd probably make `(UserId, ListingId)` unique because a user shouldn't be able to favorite the same listing multiple times.)

There could also be images:

```text
ListingImage
- Id
- ListingId
- Url
- SortOrder
```

I wouldn't necessarily store the image itself in the database. I'd store metadata/reference information and use object storage for the actual image.

---

### Interviewer follow-up: "Would you store the address as one field?"

> "I could, but because we're going to search by geographic attributes, I'd probably separate at least the city, postal code, and potentially latitude and longitude. That gives us more flexibility for filtering and eventually supporting location-based searches."

### Interviewer: "What if we want users to search within 5 km of a location?"

> "Then I'd want latitude and longitude available, and I'd choose a database or indexing mechanism that supports efficient geospatial queries."

(I wouldn't immediately introduce another technology. I'd first establish that this requirement exists.)

---

# 3. Choose the database

> "Based on the requirements, I'd initially choose a relational database."

> "The core data has clear relationships between users, listings, and favorites. We also need transactional consistency when creating or modifying those records, and the data has a reasonably structured schema."

> "The workload is heavily read-oriented, so I'd focus on making the read path efficient rather than choosing a database simply because the system has high traffic."

(I'd deliberately avoid naming SQL Server/PostgreSQL/etc. unless the interviewer asks. The important architectural decision here is relational versus another data model.)

### Interviewer: "Why not NoSQL?"

> "It's possible, especially if we had extremely high scale or a data model that benefited from denormalization. But I don't see a strong reason to give up relational modeling and transactional guarantees for the core entities yet."

> "I'd start with the relational database and introduce additional storage or indexing mechanisms if the search workload justified it."

---

# 4. Define the APIs

> "Now I'll define the API surface based on the functional requirements."

For listings:

```text
POST   /listings
GET    /listings/{id}
PUT    /listings/{id}
DELETE /listings/{id}
```

For searching:

```text
GET /listings?city=Berlin&minPrice=500&maxPrice=1500&rooms=2
```

For favorites:

```text
POST   /listings/{id}/favorite
DELETE /listings/{id}/favorite
GET    /users/me/favorites
```

> "I'd also make sure the listing modification endpoints verify that the authenticated user owns the listing before allowing the operation."

### Interviewer: "Would you use POST or PUT for updating?"

> "I'd generally use PUT when the operation represents replacing the resource representation, and PATCH if we're explicitly supporting partial updates. I'd choose based on the API semantics we're defining."

### Interviewer: "What would the search response look like?"

> "I wouldn't return everything in the listing entity. I'd return a search-specific representation containing the fields needed for the results page, such as ID, title, price, location, rooms, area, thumbnail, and perhaps some summary information."

(I want to avoid accidentally coupling the search response to the database entity.)

I'd also add pagination:

```text
GET /listings?...&page=2&pageSize=20
```

But at larger scale I'd prefer **cursor-based pagination** if the requirements call for deep pagination or very large result sets.

---

# 5. First architecture

> "I'll start with the simplest architecture that satisfies the requirements."

I'd draw:

```text
Client
   |
   v
API
   |
   v
Database
```

> "The API handles authentication, authorization, validation, listing operations, search, and favorites. The database stores the persistent state."

(I intentionally haven't drawn Redis, Kafka, Kubernetes, microservices, or anything else. I want every component to have a reason to exist.)

---

# 6. Think about the search workload

Now I'd focus on the biggest obvious characteristic:

> "The search workload is significantly heavier than writes, and we have a 200 ms latency target."

> "Before introducing another technology, I'd make sure the database queries are properly indexed and that we're only retrieving the columns we actually need."

For example, searches might frequently filter on:

```text
City
Price
PropertyType
NumberOfRooms
Status
```

and potentially sort by:

```text
Price
CreatedAt
```

> "I'd inspect the actual query plans and workload before deciding exactly which indexes to create. I don't want to create indexes on every column because indexes also increase storage and make inserts and updates more expensive."

---

# 7. Interviewer: "The search query is becoming slow. What do you do?"

> "I'd first measure rather than immediately add infrastructure."

> "I'd look at the actual execution plan, estimated versus actual row counts, index usage, expensive operators, scans versus seeks, sort operations, and whether we're retrieving significantly more data than necessary."

> "I'd also look at the query's selectivity and whether the existing indexes actually match the filtering and ordering patterns."

> "If the database itself is still the bottleneck after query and index optimization, then I'd consider scaling the read path."

---

# 8. Caching

Now I have a concrete reason to consider caching.

> "Because listing searches and listing details are read-heavy, I'd consider caching frequently requested data if we see enough repeated reads."

> "For example, listing details are often requested repeatedly, and some popular listings may receive a large number of reads."

I'd update the architecture:

```text
Client
   |
   v
API
   |
   v
Cache
   |
   v
Database
```

For a listing detail request:

> "The API first checks the cache. If the listing is there and hasn't expired, we return it without querying the database."

> "If it's not there, we query the database, return the result, and populate the cache."

So:

```text
Cache hit → return cached listing
Cache miss → database → return result → populate cache
```

### Interviewer: "What about stale data?"

> "That's a trade-off I'd explicitly define. Since the requirements allow a small delay, I'd give cached entries a reasonable TTL and also invalidate or update the cache when a listing changes."

> "The exact strategy depends on how frequently listings change and how stale we're willing to let search results become."

---

# 9. Interviewer: "What if the cache goes down?"

> "The database remains the source of truth. If the cache is unavailable, the API can fall back to the database."

> "The trade-off is that database load and latency will increase temporarily, so I'd want monitoring and possibly protection against a cache failure causing a sudden database overload."

(This is an important reliability point: caching should improve the system, not become the only place where important data exists.)

---

# 10. Asynchronous processing

Now I'd look for operations that don't need to block the user's request.

Suppose the interviewer says:

> "When a listing is created, we want to send email notifications to users who have saved searches matching that listing."

I'd say:

> "I wouldn't make the listing creation request wait for all those notifications to be sent."

I'd introduce a queue:

```text
Client
   |
   v
API
   |
   +----> Database
   |
   +----> Queue
             |
             v
        Notification Worker
             |
             v
       Email Provider
```

> "The API persists the listing and publishes a notification event. A background worker consumes the event and handles the notification."

> "That keeps the user-facing request fast and prevents a slow external email provider from making listing creation slow."

---

# 11. Idempotency

> "I'd also make the asynchronous processing idempotent."

Suppose the worker receives the same event twice.

I don't want:

```text
same event
    ↓
email
email
email
```

I'd give the event a unique identifier and keep track of successfully processed events.

> "If the worker receives an event it has already processed, it can safely ignore it."

(Exactly-once processing is generally difficult to guarantee across distributed systems. Designing the operation to be safely repeatable is usually more practical.)

---

# 12. Retries

Suppose the email provider is temporarily unavailable.

> "I wouldn't immediately mark the notification as permanently failed."

I'd retry with exponential backoff.

For example:

```text
Attempt 1 → fail
wait
Attempt 2 → fail
wait longer
Attempt 3 → fail
wait even longer
...
```

> "I'd also put a maximum number of retries so that a permanently failing message doesn't consume resources forever."

After the retry limit:

> "I'd move the message to a dead-letter queue so it can be inspected and potentially reprocessed later."

---

# 13. Rate limiting

Suppose the external email provider allows only 100 requests per second.

> "I'd rate-limit the worker so we don't overwhelm the provider."

And if the API itself is public:

> "I'd also consider rate limiting incoming API requests to protect the application from abusive or unexpectedly high traffic."

---

# 14. Separate queues

Suppose we now have:

- email notifications
- SMS notifications
- image processing

I wouldn't necessarily put everything onto one queue.

> "I'd consider separate queues because these workloads have different processing characteristics and failure modes."

For example, if image processing suddenly becomes extremely busy, I don't want it consuming all the worker capacity and delaying time-sensitive notifications.

This is essentially the same reasoning behind the **Bulkhead pattern** we discussed earlier.

---

# 15. Search indexing

Now suppose the interviewer says:

> "We have 50 million listings and complex searches are now too slow for our 200 ms requirement."

This is where I'd consider introducing a dedicated search system.

> "At this point, I'd consider separating the transactional storage of listings from the read-optimized search representation."

The architecture becomes conceptually:

```text
                  ┌──→ Database
                  |
Client → API ─────┤
                  |
                  └──→ Search Index
```

But I wouldn't necessarily write search data synchronously into both places.

> "The database remains the source of truth. When a listing changes, we publish an event and asynchronously update the search index."

For example:

```text
Create listing
     ↓
Database
     ↓
ListingCreated event
     ↓
Queue
     ↓
Search-index worker
     ↓
Search index
```

That fits our earlier requirement that search can tolerate a few seconds of eventual consistency.

---

# 16. Interviewer: "Why not just use the database?"

> "I would prefer the database initially because it's simpler. But if search becomes the dominant workload and we're doing complex filtering, text search, geographic searches, relevance ranking, or faceted search at very high volume, a dedicated search index can be much better suited to that workload."

> "The important trade-off is that now we have two representations of the data, so we need to deal with synchronization and eventual consistency."

---

# 17. Scale the API

Suppose traffic increases significantly.

> "The API is stateless, so we can horizontally scale it."

I'd draw:

```text
                    API Server 1
                  /
Client → Load Balancer → API Server 2
                  \
                    API Server 3
```

> "The load balancer distributes requests across the API instances."

> "This helps scalability because we can add instances as traffic increases, and it also improves availability because the system doesn't depend on a single API server."

---

# 18. What does "stateless" mean here?

Interviewer:

> "What do you mean by stateless?"

> "I mean the API servers don't keep important session or application state only in their local memory. Any API instance should be able to handle a request."

> "Persistent state belongs in shared infrastructure such as the database, and shared caching state belongs in a distributed cache if we need it."

This becomes particularly important now.

---

# 19. Distributed cache

> "Once we have multiple API instances, I wouldn't want each instance to maintain its own independent cache."

For example, if:

```text
Request 1 → API Server A
Request 2 → API Server B
```

and A has a cached listing but B doesn't, we can get inconsistent cache behavior and duplicate cache population.

So I'd introduce a distributed cache:

```text
                    API Server 1
                  /
Client → Load Balancer
                  \
                    API Server 2
                         |
                         v
                  Distributed Cache
                         |
                         v
                     Database
```

> "Now all API instances can access the same cache."

---

# 20. Database scaling

Now I'd ask:

> "At this point I'd also want to determine whether the database has become the bottleneck."

If reads are the problem:

> "We could consider read replicas, depending on our consistency requirements."

Writes would continue going to the primary, while appropriate read traffic could go to replicas.

But:

> "I'd need to account for replication lag. If a user creates a listing and immediately searches for it, they might not see it on a replica yet."

That's acceptable in our assumed requirements because we already established that search can be eventually consistent.

---

# 21. Database becomes too large

Interviewer:

> "What if we have hundreds of millions or billions of listings?"

> "I'd first look at whether indexing, partitioning, archiving expired listings, read replicas, and query optimization are sufficient."

> "If the scale eventually exceeds what a single database architecture can handle, I'd consider partitioning or sharding based on the access patterns."

I wouldn't immediately shard.

> "Sharding introduces significant application and operational complexity, so I'd only introduce it when the scale actually requires it."

---

# 22. Availability

Now I'd explicitly consider failures.

> "I'd identify the major single points of failure."

For example:

- API server
- load balancer
- cache
- database
- queue
- search index
- external providers

For API servers:

> "Multiple instances behind the load balancer prevent one failed instance from taking down the API."

For the database:

> "I'd use the database's high-availability capabilities, such as replication/failover, depending on the chosen technology."

For the cache:

> "The database remains the source of truth, so cache failure degrades performance rather than losing listing data."

For the search index:

> "The database remains the source of truth. If the search index becomes temporarily unavailable, we need to decide whether to temporarily disable search, fall back to the database for limited functionality, or queue index updates until it recovers."

---

# 23. What happens if the queue is down?

> "I'd want the database transaction and event publication to be reliable. If we simply write to the database and then publish to the queue as two independent operations, we can have a failure between them where the listing exists but the event was never published."

This is an important distributed-systems problem.

I'd consider the **Outbox pattern**.

Conceptually:

```text
Database transaction
    |
    +-- Listing
    |
    +-- Outbox event
```

Both are committed atomically.

A background process then reads the outbox and publishes the event to the queue.

> "That prevents us from losing the event simply because the application crashed after committing the listing but before publishing the message."

---

# 24. Observability

At this point I'd add observability.

> "I'd want metrics, logs, and distributed tracing."

### Metrics

I'd monitor:

- API request latency
- requests per second
- error rate
- database latency
- cache hit rate
- cache latency
- queue depth
- queue processing latency
- failed/retried messages
- search latency
- search-index update lag

### Logs

I'd log things such as:

- failed listing operations
- authorization failures
- database errors
- failed external-provider calls
- queue processing failures

### Traces

> "Tracing would be particularly useful for following a request across multiple components."

For example:

```text
Client
  ↓
Load Balancer
  ↓
API
  ↓
Cache
  ↓
Database
```

Or for asynchronous work:

```text
API
  ↓
Queue
  ↓
Worker
  ↓
External Provider
```

> "I'd want a correlation or trace ID so that I can connect the work across these components."

---

# 25. Security

I'd also explicitly address security before wrapping up.

> "I'd authenticate users and authorize listing operations so that users can only modify listings they own."

> "I'd validate all incoming data and enforce reasonable request-size limits, especially because listings may contain images or large descriptions."

> "I'd also rate-limit public endpoints and protect sensitive operations from abuse."

And for data:

> "I'd use encrypted connections between services and ensure sensitive information isn't written into application logs."

---

# 26. Interviewer: "Would you make this microservices?"

This is a very likely follow-up.

> "Not initially."

> "I'd start with a modular monolith or a small number of services because the initial requirements don't justify the operational complexity of many independently deployed services."

> "The architecture can still have clear boundaries—for example, listing management, search, and notifications—without immediately turning each boundary into a separate service."

> "If a particular component needs to scale independently or has significantly different availability or deployment requirements, that's when I'd consider extracting it."

---

# 27. Interviewer: "Where would you introduce Kubernetes?"

> "I wouldn't introduce Kubernetes simply because we're scaling the API."

> "If we have enough services and containers that orchestration, automated deployment, scaling, health management, and service scheduling become useful, Kubernetes could be appropriate."

> "But it's an implementation choice rather than a fundamental requirement of the architecture."

---

# 28. Final architecture

At this point, after requirements have justified the additional components, I'd have something roughly like:

```text
                         ┌───────────────┐
                         │ Search Index  │
                         └───────▲───────┘
                                 │
                                 │ async update
                                 │
Client
   |
   v
Load Balancer
   |
   +-------------------+
   |                   |
   v                   v
API Server 1        API Server 2
   |                   |
   +---------+---------+
             |
             v
      Distributed Cache
             |
             v
          Database
             |
             |
             v
          Outbox
             |
             v
           Queue
          /     \
         /       \
 Notification   Search
   Worker       Worker
      |            |
      v            v
Email/SMS      Search Index
Provider
```

(I wouldn't necessarily draw this entire thing at the beginning. The important part of the interview is showing *why each piece appeared*.)

---

# 29. Final summary to the interviewer

> **"I'd start with a relatively simple architecture: a stateless API backed by a relational database. The database would remain the source of truth for users, listings, and favorites."**
>
> **"Because the workload is heavily read-oriented, I'd optimize the database queries and indexes first, and then introduce caching for frequently accessed data where the consistency requirements allow it."**
>
> **"As traffic grows, I'd horizontally scale the stateless API behind a load balancer and use a distributed cache so all API instances can share cached data."**
>
> **"For operations that don't need to complete synchronously, such as notifications and search-index updates, I'd introduce queues and workers. Those workers would use idempotency, retries with exponential backoff, rate limiting, and dead-letter queues."**
>
> **"If search becomes a major bottleneck because of the number and complexity of queries, I'd separate search from transactional storage and maintain a dedicated search index asynchronously. That gives us better search performance while accepting some eventual consistency."**
>
> **"For reliability, I'd make the database the source of truth, use appropriate database high-availability mechanisms, and use an outbox pattern where reliable event publication is required."**
>
> **"Finally, I'd add metrics, logs, and distributed tracing so we can identify latency, errors, database problems, cache behavior, and asynchronous processing failures."**
>
> **"The main principle is that I wouldn't start with all of these components. I'd start simple and introduce each one when a specific requirement or bottleneck justifies the additional complexity."**

That is the version I'd be comfortable giving in an actual system-design interview.

Yes — these are exactly the things we should clarify before continuing the design. And you're right about the search index: **I introduced a concept we hadn't covered yet.** We should cover it properly rather than casually throwing it into the architecture.

### 1. "Transactional guarantees" / "transactional consistency"

Yes, I'm referring to **database transactions**, but the idea is broader than simply "SQL has transactions."

A transaction lets you treat multiple database operations as **one logical unit of work**.

For example, when creating a listing, imagine we need to:

1. Insert the listing.
2. Insert its images.
3. Insert some related record.

Suppose step 1 and 2 succeed, but step 3 fails.

Without a transaction, you could end up with partially completed data:

```text
Listing       ✓
Images        ✓
Related data  ✗
```

With a transaction:

```text
All operations succeed → COMMIT
Any operation fails    → ROLLBACK
```

So the database can guarantee that the transaction's changes are applied **atomically**.

That's one aspect of transactional consistency.

There are actually several properties commonly associated with transactions, summarized as **ACID**:

- **Atomicity** — all changes happen, or none do.
- **Consistency** — the transaction takes the database from one valid state to another, respecting its constraints/rules.
- **Isolation** — concurrent transactions don't improperly interfere with each other.
- **Durability** — once committed, the data survives failures.

We've already discussed isolation in detail.

#### Does NoSQL not provide transactions?

This is where saying "SQL has transactions, NoSQL doesn't" would be wrong.

Many NoSQL databases **do support transactions**, sometimes even across multiple records/documents.

The more accurate statement is:

> **Relational databases traditionally provide very strong transactional and relational capabilities, while NoSQL databases vary considerably in the transactional guarantees and consistency models they provide.**

Some NoSQL systems deliberately make trade-offs around transactions, consistency, joins, schema flexibility, etc., in exchange for other characteristics such as horizontal scalability or a data model optimized for particular access patterns.

So in our system-design interview, I wouldn't say:

> "I chose SQL because NoSQL doesn't support transactions."

I'd say:

> "The core entities have clear relationships and require transactional operations, so a relational database is a natural starting point."

---

# 2. Why `/users/me/favorites` instead of `/users/{userId}/favorites`?

Good catch.

This is an **API design choice**, not a technical requirement.

If the user is authenticated, the server already knows who they are from their authentication credentials/token.

So:

```text
GET /users/me/favorites
```

means:

> "Give me the favorites belonging to the currently authenticated user."

The server obtains the user's identity from the authentication context rather than trusting the client to provide a user ID.

For example, the request might contain an access token identifying user `123`.

The API effectively knows:

```text
Authenticated user = 123
```

and retrieves that user's favorites.

This is useful because you don't want a normal user to simply do:

```text
GET /users/456/favorites
```

and potentially see another user's private data.

You **can** use:

```text
GET /users/{userId}/favorites
```

if the API requires that resource-oriented structure and then performs authorization:

> "Does the authenticated user have permission to access user 456's favorites?"

Both designs are valid.

For this particular system, I prefer `/users/me/favorites` because favorites are inherently **the current user's personal resource**.

---

# 3. "Scale the read path"

Yes — **read replicas are one thing I could mean**, but I deliberately said it more generally.

Suppose the database is struggling because we have:

```text
10,000 reads/sec
500 writes/sec
```

and most of the load is reads.

Before adding infrastructure, I'd optimize:

- queries
- indexes
- projections
- pagination
- unnecessary database calls

If the database is **still** the bottleneck, one option is read replicas:

```text
                 ┌──→ Read Replica 1
                 │
API → Database ──┼──→ Read Replica 2
                 │
                 └──→ Read Replica 3
```

Writes go to the primary.

Appropriate reads go to replicas.

That allows us to distribute read load.

### But there's an important problem

Replication isn't necessarily instantaneous.

Suppose:

```text
User creates listing
       ↓
Primary DB
       ↓
Immediately searches
       ↓
Read Replica
```

The replica might not have received the new listing yet.

That's **replication lag**.

And this is why our earlier requirement was useful:

> "A few seconds of delay before a new listing appears in search is acceptable."

That makes read replicas much easier to consider.

So yes, **read replicas are one example of scaling the read path**, but I could also consider:

- caching
- a dedicated search system
- database partitioning
- query/index optimization

depending on what the actual bottleneck is.

---

# 4. What is a Search Index?

Yes. **We haven't covered this yet.** Let's establish it now because it's important to the system we're designing.

A search index is a data structure/system specifically optimized for **finding records based on search criteria**.

It's different from the indexes we've been discussing inside a relational database.

Imagine we have 50 million property listings.

A user searches:

> Apartments in Berlin, between €500 and €1,500, 2+ rooms, sorted by relevance.

A normal relational database *can* handle this, especially with good indexes.

But real-estate search can become much more complicated:

- multiple filters
- full-text search
- geographic distance
- relevance ranking
- fuzzy matching
- autocomplete
- facets/aggregations
- sorting
- very large result sets

A dedicated search system is designed specifically around those kinds of operations.

---

## What does it actually contain?

It doesn't necessarily contain the exact same representation as our database.

Our database might have:

```text
Listing
Id
Title
Description
Price
City
PostalCode
NumberOfRooms
PropertyType
...
```

The search index might contain a **search-optimized representation**:

```text
Listing ID
Title
Description
Price
Location
Coordinates
Rooms
Property type
Status
...
```

The search system creates specialized structures from this data that make searching fast.

The database remains our **source of truth**.

---

## How does the data get there?

This is where our earlier queue discussion becomes relevant.

When somebody creates or updates a listing:

```text
API
 ↓
Database
 ↓
ListingChanged event
 ↓
Queue
 ↓
Search indexing worker
 ↓
Search index
```

The worker takes the listing information and updates the search index.

So if the listing is changed at 10:00:00, the database might have the new value immediately, while the search index gets it at 10:00:02.

That's **eventual consistency**.

And our requirements explicitly allow that.

---

## Why not just query the database?

We absolutely can **at the beginning**.

That's actually how I'd design the first version.

```text
Client
   ↓
API
   ↓
Database
```

Then, if we discover that property search is becoming the bottleneck because of scale and complexity:

> "I'd consider introducing a dedicated search index."

That's an architectural evolution rather than something I'd blindly add from the start.

---

## One important distinction

Don't confuse these two:

### Database index

An index inside a database that helps the database efficiently locate rows.

For example:

```text
INDEX on (City, Price)
```

It can make a SQL query much faster.

### Search index

A separate search-oriented data structure/system designed for sophisticated search workloads.

It might support things like:

> "Find apartments within 5 km of this location, under €1,500, containing the word 'balcony', with at least two rooms, ranked by relevance."

That's a substantially different workload.

---

So in **our system design**, I'd currently keep the search index **out of the initial architecture**.

We'd start with the relational database, optimize its queries/indexes, measure the actual workload, and **only introduce a dedicated search system if the requirements/scale justify it**.

That also fits perfectly with the principle you established:

> **Start simple. Add infrastructure because you have a reason for it.**