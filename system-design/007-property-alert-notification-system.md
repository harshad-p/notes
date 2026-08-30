# Property Alert / Notification System

## 1. Clarify requirements

> **"Before I design the system, I'd like to clarify a few requirements."**

### Functional requirements

> "My understanding is that users should be able to create saved property searches. For example, a user could specify a location, price range, property type, and number of rooms."

> "When a new property listing is created, we should determine which saved searches match that listing and notify the corresponding users."

> "What notification channels do we need to support? I'm assuming email, but should we also support SMS or push notifications?"

> "Should users be able to have multiple saved searches?"

> "Can users modify and delete their saved searches?"

> "Should users receive a notification for every matching listing, or should we aggregate multiple matches into a periodic notification?"

> "Do we need to notify users about updates to existing listings as well, or only newly created listings?"

> "Should a user receive the same notification more than once if the system processes the same listing event multiple times?"

(I've asked that last question because it will become important when we discuss retries and duplicate messages.)

### Non-functional requirements

> "What kind of scale are we expecting? Approximately how many users, saved searches, and new listings do we expect per day?"

Let's assume:

- 10 million users
- 50 million saved searches
- 1 million new/updated listings per day
- potentially millions of notifications per day
- notification processing can be asynchronous
- a notification doesn't have to arrive immediately; a delay of a few seconds is acceptable
- high availability is required
- eventual consistency is acceptable

> "How quickly should a new listing trigger a notification?"

Assume:

> "Within a few seconds under normal load."

> "Do we need strong consistency when users create or modify saved searches?"

Assume:

> "No. Eventual consistency is acceptable."

> "Do we expect the saved-search criteria to evolve over time? For example, could we add new types of filters?"

Assume:

> "Yes."

(That last requirement is going to be one of the reasons I consider a document-oriented NoSQL database.)

---

# 2. Determine what data needs to be stored

> "I'd identify three main types of data: users, saved searches, and notification state."

A saved search might contain:

```text
SavedSearch
- Id
- UserId
- Location
- MinPrice
- MaxPrice
- PropertyType
- MinRooms
- MaxRooms
- CreatedAt
- UpdatedAt
- ExpiresAt
- Enabled
```

I'd also have user notification preferences:

```text
NotificationPreference
- UserId
- EmailEnabled
- SmsEnabled
- PushEnabled
```

And I'd need some state around notifications:

```text
Notification
- Id
- UserId
- SavedSearchId
- ListingId
- Channel
- Status
- CreatedAt
- SentAt
- RetryCount
```

> "I'd also want timestamps such as `CreatedAt` and `UpdatedAt`, and potentially `ExpiresAt` if saved searches can expire."

---

## Interviewer: "Would you put all of this in one database?"

> "Not necessarily. I'd first look at the access patterns."

> "Users and notification preferences are relatively straightforward. The saved searches are more interesting because we have potentially tens of millions of them, they're queried frequently when listings arrive, and their structure may evolve as we add new search criteria."

---

# 3. Choose the database

> "For saved searches, I'd consider a document-oriented NoSQL database."

> "The saved search is naturally represented as a document, and different versions of the application may introduce additional search criteria over time. We don't necessarily need to model every possible criterion as a separate relational table."

For example, one document could conceptually look like:

```text
{
    id: 123,
    userId: 456,
    location: "Berlin",
    minPrice: 500,
    maxPrice: 1500,
    propertyType: "Apartment",
    minRooms: 2,
    notificationChannels: ["email", "push"],
    enabled: true
}
```

> "The exact NoSQL technology would depend on the access patterns and scale, but I'm choosing the NoSQL category because the workload is large, the data is naturally document-shaped, and we expect the structure of saved-search criteria to evolve."

(I am deliberately not saying "NoSQL because it's faster." NoSQL isn't inherently faster; the data model and access pattern are what matter.)

---

## Interviewer: "Why not SQL?"

> "SQL could absolutely handle this. I'm not saying it couldn't."

> "I'd choose SQL if we needed strong relational constraints, complex relationships, or transactional operations across these entities."

> "In this case, the saved search is relatively self-contained, the structure may evolve, and we're primarily interested in high-volume access based on specific lookup patterns. A document-oriented model can fit that very well."

---

# 4. Think about the access pattern

Now I'd ask:

> "The most important access pattern seems to be: given a newly created listing, find the saved searches that could match it."

This is important because NoSQL database design is often **access-pattern driven**.

I don't want to design the document structure first and only later discover that my most important query is extremely expensive.

Suppose a listing arrives:

```text
Berlin
Apartment
€1,200
3 rooms
```

I need to efficiently identify potentially matching saved searches.

> "I'd therefore design the keys and indexes around the queries we actually need to perform."

---

# 5. Define the APIs

For saved searches:

```text
POST   /saved-searches
GET    /saved-searches
GET    /saved-searches/{id}
PUT    /saved-searches/{id}
DELETE /saved-searches/{id}
```

For notification preferences:

```text
GET /users/me/notification-preferences
PUT /users/me/notification-preferences
```

> "The API is primarily responsible for managing saved searches and preferences. Notification generation itself doesn't need to happen synchronously inside these HTTP requests."

---

# 6. First architecture

> "I'll start with the simplest architecture."

```text
Client
   |
   v
API
   |
   v
NoSQL Database
```

> "The API handles saved-search creation and modification, while the database stores the saved searches and notification preferences."

I wouldn't add queues yet.

---

# 7. Now consider what happens when a listing is created

The interviewer tells me:

> "A separate listing system creates the property listings."

I'd ask:

> "Does our system receive an event when a listing is created?"

Assume yes.

> "Then I wouldn't have the listing system synchronously call our notification API for every user. I'd consume the listing-created event asynchronously."

I'd introduce a queue:

```text
Listing Service
      |
      v
   Queue
      |
      v
Alert Processor
      |
      v
NoSQL Database
```

> "The alert processor consumes listing events and determines which saved searches match the new listing."

---

# 8. Why a queue?

> "The queue gives us buffering and decouples listing creation from notification processing."

Suppose there is suddenly a large number of listings.

Without a queue, notification processing could slow down the listing creation workflow.

With a queue:

> "The listing service can continue publishing events, while our workers process them at their own rate."

If we suddenly have 100,000 events waiting:

> "The queue absorbs the spike instead of forcing the listing service to wait."

---

# 9. Matching saved searches

Now comes one of the harder parts.

> "For every listing, I need to find the saved searches whose criteria match it."

I don't want to simply retrieve all 50 million saved searches and test each one.

That would be extremely expensive.

> "I'd create indexes based on the most important filtering dimensions, such as location and potentially property type."

For example, if location is one of our major access patterns:

```text
Location = Berlin
```

we can narrow the candidates substantially before applying the remaining criteria.

> "I'd distinguish between candidate retrieval and exact matching. The database can narrow the candidate set, and the application can then evaluate the remaining criteria."

(I don't need every possible condition to be represented perfectly in a database index. I need the overall operation to be efficient.)

---

# 10. Interviewer: "What if Berlin has 10 million saved searches?"

> "Then location alone isn't selective enough."

> "I'd look at the actual distribution and access patterns and potentially use a compound key or additional partitioning dimensions."

For example, I might partition based on:

- location
- property type
- other high-selectivity attributes

But:

> "I would be careful not to create a partitioning scheme that produces hot partitions."

---

# 11. What is a hot partition?

Suppose we partition saved searches by city.

If our system has:

```text
Berlin     20 million searches
Munich      2 million
Hamburg     1 million
Other      27 million
```

Berlin receives disproportionately more traffic.

> "The Berlin partition could become a hot partition and become the bottleneck even though the overall database has plenty of capacity."

So I'd potentially introduce a mechanism to distribute heavily accessed data more evenly.

---

# 12. Scaling the alert processors

Suppose we now receive millions of listing events.

> "The alert processor is stateless, so we can horizontally scale it."

```text
                    Worker 1
                   /
Queue →           Worker 2
                   \
                    Worker 3
```

> "Multiple workers can consume messages concurrently."

But I need to make sure the queue's consumer semantics support safe concurrent processing.

---

# 13. Idempotency

This is critical.

Suppose:

```text
ListingCreated event
        ↓
Worker processes it
        ↓
Notification sent
        ↓
Worker crashes before acknowledging message
```

The queue may deliver the same event again.

Now we could send:

```text
Notification
Notification
```

> "I'd therefore make notification processing idempotent."

I'd give the event or notification a unique identifier.

Before sending:

> "The worker checks whether this notification has already been successfully processed. If it has, it skips it."

There is another subtle issue here.

If I check:

```text
Has notification been sent?
```

and then separately:

```text
Send notification
```

those aren't necessarily one atomic operation.

The process could crash between them.

So I'd design the operation carefully around an idempotency key and the capabilities of the notification provider.

---

# 14. External notification providers

Now we have:

```text
Alert Processor
      |
      v
Notification Queue
      |
      +------> Email Worker → Email Provider
      |
      +------> SMS Worker   → SMS Provider
      |
      +------> Push Worker  → Push Provider
```

> "I'd separate the notification workloads because email, SMS, and push can have different throughput limits, costs, and failure characteristics."

(It's the same Bulkhead reasoning we discussed earlier.)

---

# 15. Rate limiting

Suppose the SMS provider allows:

> 100 requests per second.

> "I'd rate-limit the SMS worker so that we don't exceed the provider's limits."

The same principle applies to email and push providers if they impose their own limits.

---

# 16. Retry mechanism

Suppose the email provider returns a temporary failure.

> "I'd retry the operation rather than immediately failing permanently."

I'd use:

> "Exponential backoff with a maximum retry count."

For example:

```text
Attempt 1 → fail
       ↓
wait
Attempt 2 → fail
       ↓
wait longer
Attempt 3 → fail
       ↓
...
```

> "If we reach the maximum number of retries, I'd move the message to a dead-letter queue."

Then the message can be inspected or reprocessed later.

---

# 17. Separate queues and failure isolation

Suppose the SMS provider is completely down.

> "I don't want SMS failures to prevent email notifications from being processed."

That's why:

```text
Email Queue → Email Workers
SMS Queue   → SMS Workers
Push Queue  → Push Workers
```

> "Now a failure in one notification channel doesn't consume all the capacity of the others."

---

# 18. Caching

Now I'd ask:

> "Are there any objects we're repeatedly reading that change relatively infrequently?"

For example, user notification preferences.

If they're accessed extremely frequently:

> "I could consider caching them."

But I wouldn't automatically cache all saved searches.

Why?

Because saved searches are being actively queried as part of the matching process, and caching tens of millions of them may be expensive and difficult to invalidate.

> "I'd first measure the access pattern and only cache data where the benefit justifies the additional complexity."

---

# 19. Scale

Now I'd look at the major scaling dimensions:

### API

> "The API is stateless, so I can horizontally scale it behind a load balancer."

### Workers

> "The notification workers can also scale horizontally based on queue depth."

### Database

> "The NoSQL database can be partitioned horizontally based on the access patterns and partition-key design."

### Queue

> "The queue itself needs to support the required throughput and partitioning."

---

# 20. Distributed architecture

At this point, my architecture is roughly:

```text
                         ┌── Email Queue ── Email Workers ── Email Provider
                         │
Listing Service → Queue ─┼── SMS Queue ─── SMS Workers ─── SMS Provider
                         │
                         └── Push Queue ── Push Workers ── Push Provider
                             
                                   ^
                                   |
                              Alert Processor
                                   |
                                   v
                              NoSQL Database
```

And separately:

```text
Client
   |
   v
Load Balancer
   |
   +-------- API Server 1
   |
   +-------- API Server 2
   |
   +-------- API Server 3
              |
              v
         NoSQL Database
```

---

# 21. Interviewer: "Why not have the alert processor directly send the email?"

> "I could for a simple system, but I'd rather separate matching from notification delivery."

> "The matching process answers the question 'who should be notified?' while the notification workers handle 'how do I deliver the notification?'"

> "That allows each side to scale independently and prevents a slow email provider from slowing down matching."

---

# 22. Interviewer: "What if one listing matches one million saved searches?"

This is an important scaling problem.

> "I wouldn't necessarily try to synchronously create one million notifications in a single operation."

> "I'd break the work into manageable batches and publish notification tasks asynchronously."

For example:

```text
Listing
   |
   v
Find matching searches
   |
   v
Create notification jobs in batches
   |
   v
Notification queues
```

> "I'd also put limits around the amount of work one listing can generate so that a single pathological event doesn't consume the entire worker pool."

---

# 23. Interviewer: "What if one user has 500 saved searches and the same listing matches 50 of them?"

> "I'd probably avoid sending 50 separate notifications to the same user for the same listing."

> "I'd deduplicate at the user/listing level, so the user receives one notification for that listing even if multiple saved searches matched it."

That gives me another useful idempotency/deduplication key:

```text
UserId + ListingId + NotificationType
```

> "I'd enforce that combination according to the notification semantics we choose."

---

# 24. What if a user creates a saved search while events are being processed?

> "Because we're accepting eventual consistency, I don't necessarily need every component to observe the change immediately."

But I'd still define the behavior.

For example:

> "If the saved search exists before the listing event is processed, it should be eligible for that event. If the saved search is created afterward, I wouldn't necessarily go back and generate notifications for historical listings."

That makes the semantics clear.

---

# 25. What if the NoSQL database goes down?

> "The database needs to provide the required availability and durability guarantees, including replication across failure domains."

If it becomes temporarily unavailable:

> "The alert processor shouldn't simply lose the listing event. The queue retains the event, so processing can resume once the database becomes available."

That's one of the major benefits of having the queue between the event source and processing.

---

# 26. What if the notification provider is down for an hour?

> "The notification job remains in the queue while we retry according to our retry policy."

> "If the provider remains unavailable long enough to exceed the retry limit, the messages go to the dead-letter queue."

> "I'd also monitor the queue depth and provider error rate so we'd know that delivery is falling behind."

---

# 27. Observability

I'd monitor:

### Metrics

- API latency
- API requests/sec
- API error rate
- queue depth
- queue processing latency
- notification processing rate
- notification failure rate
- retry count
- dead-letter queue size
- database latency
- database throttling
- worker utilization
- notifications sent per channel
- provider response rates

### Logs

> "I'd log failed notification attempts, invalid events, database failures, provider errors, and dead-lettered messages."

### Tracing

> "I'd use distributed tracing where possible to follow an operation across the API, queue, worker, and external provider."

For example:

```text
API
 ↓
Queue
 ↓
Alert Processor
 ↓
Notification Queue
 ↓
Email Worker
 ↓
Email Provider
```

---

# 28. Security

> "I'd authenticate users and authorize access to their saved searches."

> "I wouldn't trust a user-provided user ID when accessing their own searches if the identity is already available from the authentication context."

So:

```text
GET /saved-searches
```

is preferable to exposing another user's searches through something like:

```text
GET /users/{userId}/saved-searches
```

unless there's a legitimate administrative use case.

I'd also protect provider credentials and avoid putting sensitive information into logs.

---

# 29. Interviewer: "Would you use Kafka?"

> "Potentially, but I wouldn't decide that simply because we're building an event-driven system."

> "I'd choose a messaging technology based on the required throughput, ordering requirements, delivery semantics, retention, replay requirements, and operational constraints."

> "If we need high-throughput event streaming and replayable event history, a system like Kafka could be appropriate. If we primarily need a work queue where workers consume tasks, a traditional message queue could be a better fit."

---

# 30. Interviewer: "Would you use microservices?"

> "I'd probably separate the notification processing boundary because it has very different scaling and failure characteristics from the API."

> "But I wouldn't automatically create a microservice for every component."

> "I'd start with clear logical boundaries and extract independently deployable services where there's a concrete benefit."

---

# 31. Final architecture

At this point I'd present the design as:

```text
                              ┌───────────────┐
                              │ Email Queue   │
                              └───────┬───────┘
                                      ↓
                                Email Worker
                                      ↓
                                Email Provider


Listing Service
      |
      v
 Listing Event Queue
      |
      v
 Alert Processor
      |
      v
   NoSQL DB
      ^
      |
      |
Client → Load Balancer → API Servers
                           |
                           |
                           v
                    Saved Searches
                    & Preferences


Alert Processor
      |
      v
Notification Jobs
      |
      +──────────────→ Email Queue
      |
      +──────────────→ SMS Queue
      |
      +──────────────→ Push Queue
```

With retries and dead-letter queues associated with the notification queues.

---

# 32. Final summary to the interviewer

> **"I'd start with a stateless API and a NoSQL database for the saved searches and notification preferences. The document model fits the relatively self-contained saved-search objects, and the schema can evolve as we add new search criteria."**
>
> **"Listing creation is asynchronous from the notification system. When a listing is created, we receive an event through a queue. An alert processor finds saved searches that could match the listing and creates notification jobs."**
>
> **"I'd separate email, SMS, and push processing because they have different throughput limits and failure characteristics. Each has its own workers, rate limiting, retry policy with exponential backoff, and dead-letter queue."**
>
> **"I'd make processing idempotent because queue delivery can result in duplicate processing. I'd also deduplicate notifications at the user/listing level so that one listing doesn't generate multiple notifications for the same user simply because it matched several saved searches."**
>
> **"As the system grows, I'd horizontally scale the API and workers, and partition the NoSQL database according to the actual access patterns while watching for hot partitions."**
>
> **"I'd only introduce caching where measurement shows repeated reads that benefit from it, rather than caching everything."**
>
> **"Finally, I'd monitor API latency, queue depth, processing latency, database performance, notification failures, retries, and provider failures, with logs and distributed tracing for diagnosing individual operations."**
>
> **"The main trade-off is that this design favors scalability and availability and accepts eventual consistency. The additional queues and asynchronous processing make the system more resilient to spikes and external-provider failures, but they also introduce more operational complexity."**

## 1. What is an access pattern?

An **access pattern is simply a way the application needs to retrieve or modify data**.

For our saved searches, examples might be:

1. "Give me all saved searches belonging to user 123."
2. "Get saved search 456."
3. "Find saved searches that could match this new listing."
4. "Disable saved search 456."
5. "Delete all saved searches belonging to user 123."

Those are our access patterns.

The important one for this system is #3:

> **A new property was created. Which saved searches might match it?**

That's very different from:

> "Give me saved search #456."

And the database needs to be designed so that the important operations can be performed efficiently.

---

# 2. Let's start with a simple NoSQL document

Suppose we store saved searches as documents:

```text
SavedSearch
{
    id: 1001,
    userId: 25,
    city: "Berlin",
    propertyType: "Apartment",
    minPrice: 500,
    maxPrice: 1500,
    minRooms: 2,
    enabled: true
}
```

And we have millions of these:

```text
SavedSearch 1001
SavedSearch 1002
SavedSearch 1003
...
SavedSearch 50,000,000
```

Now imagine a new listing arrives:

```text
Listing
{
    id: 90001,
    city: "Berlin",
    propertyType: "Apartment",
    price: 1200,
    rooms: 3
}
```

We need to find the saved searches that this listing satisfies.

---

# 3. The naïve approach

We could theoretically ask the database:

> "Give me all 50 million saved searches."

Then our application could examine every one:

```text
Is city Berlin?
Is apartment?
Is price between 500 and 1500?
Are rooms >= 2?
Is it enabled?
```

But that's obviously terrible.

We're doing 50 million checks to find perhaps 10,000 relevant searches.

So we need a way to **narrow down the candidates before doing all the detailed matching**.

That's where keys and indexes come in.

---

# 4. In SQL, you're accustomed to starting with the data model

Suppose we have:

```sql
SavedSearch
----------------
Id
UserId
City
PropertyType
MinPrice
MaxPrice
MinRooms
Enabled
```

You might create indexes such as:

```sql
CREATE INDEX IX_SavedSearch_City
ON SavedSearch(City);
```

Then SQL Server can use that index to efficiently find rows where:

```sql
WHERE City = 'Berlin'
```

You generally have quite a lot of flexibility in asking different queries later.

---

# 5. NoSQL often makes you think differently

With many NoSQL databases, you don't want to start with:

> "What fields does my object have?"

You start with:

> **"What operations will my application need to perform?"**

For our system, we know:

> "When a listing arrives, I need to find saved searches associated with this listing's characteristics."

So we design our partition keys and indexes around that operation.

That's what I meant by:

> "Design the keys and indexes around the queries we actually need to perform."

---

# 6. Let's make this concrete

Suppose our most important question is:

> "Find enabled saved searches for Berlin apartments."

We could create an index that allows the database to locate documents based on something like:

```text
city + propertyType
```

Conceptually:

```text
Berlin + Apartment
        ↓
saved searches for that combination
```

Now when this listing arrives:

```text
Berlin
Apartment
€1200
3 rooms
```

we can first retrieve the saved searches that are relevant to:

```text
Berlin + Apartment
```

Instead of examining all 50 million.

Suppose that reduces the candidate set from:

**50 million → 2 million**

That's already a huge improvement.

Then we apply the remaining criteria:

```text
price >= minPrice
price <= maxPrice
rooms >= minRooms
enabled = true
```

Maybe:

**2 million → 15,000 matches**

Those 15,000 users are the ones we need to notify.

---

# 7. But here's the important part: why not index everything?

You might think:

> "Fine. Let's create indexes on city, property type, min price, max price, rooms, enabled, etc."

You *can* have multiple indexes in many NoSQL systems, but indexes aren't free.

An index itself requires storage.

And when you insert or update a saved search, the database may have to update the corresponding indexes.

So now:

```text
Write SavedSearch
       ↓
Store document
       ↓
Update index A
       ↓
Update index B
       ↓
Update index C
```

The more indexes you maintain, the more expensive writes can become.

So we're making a trade-off:

> **Optimize the data model and indexes for the operations that matter most, rather than indexing every field just in case.**

---

# 8. There's an even bigger issue: partitioning

This is where NoSQL gets particularly interesting.

Suppose we have 50 million saved searches.

A distributed NoSQL database might spread those documents across many machines.

We need to decide **how to distribute them**.

That's where a **partition key** comes in.

Imagine we choose:

```text
partition key = city
```

Then conceptually:

```text
Berlin    → Partition A
Munich    → Partition B
Hamburg   → Partition C
...
```

The database can therefore determine where the relevant data lives.

---

# 9. Why does that matter?

Imagine our listing is:

```text
Berlin
Apartment
€1200
3 rooms
```

If `city` is the partition key, the database knows:

> "I only need to look at the Berlin partition."

It doesn't have to search all partitions.

That's one of the major reasons **access patterns matter when choosing partition keys**.

---

# 10. But now we discover a problem

Suppose:

```text
Berlin → 20 million saved searches
Munich → 2 million
Hamburg → 1 million
```

Berlin is disproportionately popular.

Now the Berlin partition gets hammered.

This is what we discussed earlier as a **hot partition**.

So simply saying:

> "Use city as the partition key"

isn't enough.

I'd have to examine the actual distribution and workload.

Maybe I'd need something more granular, such as combining attributes or introducing another distribution mechanism.

For example, conceptually:

```text
Berlin + Apartment + bucket 1
Berlin + Apartment + bucket 2
Berlin + Apartment + bucket 3
...
```

The exact solution depends heavily on the particular NoSQL database and its partitioning model.

---

# 11. Now consider a completely different access pattern

Suppose the user opens their account page.

The application asks:

> "Give me all saved searches belonging to user 25."

That's a completely different query.

Our earlier design might have been optimized around:

> listing → matching saved searches

But now we also need:

> user → their saved searches

So we might need another key/index that supports that access pattern.

Conceptually:

```text
UserId = 25
    ↓
SavedSearch 1001
SavedSearch 1025
SavedSearch 1088
...
```

This is why you can't simply say:

> "I'll make `Id` the key."

and assume everything will be efficient.

You need to think about **how the application actually accesses the data**.

---

# 12. This is where NoSQL can feel backwards compared with SQL

With SQL, you might think:

> "Here's my normalized schema. Now I'll write queries against it and add indexes to improve them."

With many NoSQL designs, you often think:

> "Here are the queries my application needs. How should I structure the data so those queries are efficient?"

That's the essence of **access-pattern-driven data modeling**.

---

# 13. What about our price range?

There's an interesting complication here.

Our listing has:

```text
price = €1200
```

But the saved searches have:

```text
minPrice = €500
maxPrice = €1500
```

We aren't simply looking for:

```text
price = 1200
```

We're asking:

```text
minPrice <= 1200
AND
maxPrice >= 1200
```

That's a **range-overlap query**.

And that can be much harder to optimize than an equality lookup such as:

```text
city = Berlin
```

This is exactly why I'd want to understand the actual query patterns before deciding the NoSQL schema.

I might decide that:

> "Location and property type are good candidate filters for the initial lookup, and the application performs the remaining range checks."

Or, depending on scale and the database technology, I might design a more sophisticated indexing strategy.

---

# 14. This is also why our Search Index discussion matters

Remember the dedicated search index we discussed for the previous system?

The same problem starts becoming relevant here.

We're essentially asking:

> "Given this property, find all saved searches whose criteria match it."

That's a potentially complicated search operation.

At relatively small scale, the NoSQL database might handle it.

At very large scale, I might eventually introduce a specialized matching/search architecture.

But I wouldn't jump there immediately.

---

# 15. The key thing I want you to take away

When you hear **"access pattern" in a NoSQL interview**, think:

> **"What does the application need to ask the database to do?"**

For our system, examples are:

> "Get all saved searches for this user."

> "Get this saved search by ID."

> "Find saved searches that could match this listing."

> "Find enabled searches for this location and property type."

Then the NoSQL design asks:

> **"What partition keys, sort keys, indexes, and possibly denormalized data structures will make those operations efficient at our expected scale?"**

That's what I meant by **designing keys and indexes around access patterns**.

And one important correction to our earlier interview answer: **we shouldn't confidently claim that a document database is the right choice merely because the saved search has a flexible structure.** The *workload and access patterns* are the stronger justification. Schema flexibility is a supporting reason.