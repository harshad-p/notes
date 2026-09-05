# Property Viewing Appointment System

**Interviewer:**

> Design a system that allows prospective tenants/buyers to book viewing appointments for real-estate properties.

**Candidate:**

Before I design it, I'd like to clarify a few requirements.

### Functional requirements

I'd clarify:

1. Can a property have multiple available viewing slots?
2. Can multiple people book the same slot, or is each slot limited to one appointment?
3. Can users cancel or reschedule?
4. Can the property owner/agency create, modify, and cancel viewing slots?
5. Do we need reminders?
6. Should users be able to see their upcoming and past appointments?
7. What happens if two users try to book the same slot at approximately the same time?

(That last question is particularly important. I want to establish whether we're dealing with a concurrency/consistency problem before choosing the database and booking mechanism.)

I'd also clarify the non-functional requirements:

- What is the expected number of properties?
- How many viewing slots might exist?
- What's the expected booking volume?
- What availability and latency requirements do we have?
- How important is strong consistency for booking?
- How much downtime is acceptable?
- Are duplicate bookings unacceptable?

For this design, I'll assume:

- Millions of properties
- Each property can have many viewing slots
- A slot can be booked by **one user**
- Users can cancel
- Agencies can create and modify slots
- Users should receive confirmation and reminder notifications
- Booking must prevent double-booking
- The system should remain responsive during periods of high demand

(The most important requirement here is that **we must never successfully book the same slot for two users**.)

---

# Data

I'd start by identifying the core entities.

We need something like:

### User

```text
User
-----
Id
Name
Email
CreatedAt
```

### Property

```text
Property
--------
Id
OwnerId / AgencyId
Address
...
CreatedAt
```

### ViewingSlot

```text
ViewingSlot
-----------
Id
PropertyId
StartTime
EndTime
Status
CreatedAt
```

`Status` might be something like:

```text
Available
Booked
Cancelled
```

### Appointment

```text
Appointment
-----------
Id
ViewingSlotId
UserId
Status
CreatedAt
CancelledAt
```

I'd probably also keep timestamps such as `CreatedAt`, `UpdatedAt`, and potentially `ExpiresAt` where they make sense.

(I'm deliberately separating **ViewingSlot** from **Appointment**. The slot represents availability created by the agency; the appointment represents a user's booking of that slot.)

---

# Database

I'd choose a **relational database** here.

The reason is primarily the consistency requirement.

We have relationships between:

- users
- properties
- viewing slots
- appointments

More importantly, we need a strong guarantee around:

> **One slot can have at most one active appointment.**

This is a case where transactional guarantees and constraints provided by a relational database are very useful.

I don't need NoSQL here simply because the system needs to scale. I'd first determine whether the relational database can handle the expected workload.

---

# API

I'd expose APIs along these lines.

For users:

```text
GET  /properties/{propertyId}/viewing-slots
POST /viewing-slots/{slotId}/appointments
GET  /users/me/appointments
POST /appointments/{appointmentId}/cancel
```

For agencies:

```text
POST   /properties/{propertyId}/viewing-slots
PUT    /viewing-slots/{slotId}
DELETE /viewing-slots/{slotId}
```

I'd use `/users/me/appointments` rather than requiring a user ID from the client because the authenticated identity already tells the API which user's appointments should be returned.

(It also avoids trusting a client-supplied user ID for an operation that should belong to the authenticated user.)

---

# First version of the architecture

I'd start very simply.

```text
Client
   ↓
API Server
   ↓
SQL Database
```

The API handles:

- authentication/authorization
- retrieving slots
- creating appointments
- cancellation
- agency slot management

I wouldn't introduce Redis, Kafka, Kubernetes, or multiple services yet.

(The first version should solve the actual business problem. I'll introduce infrastructure only when a requirement gives me a reason to.)

---

# The important part: booking

The most important operation is:

```text
POST /viewing-slots/{slotId}/appointments
```

Suppose two users do this at almost exactly the same time:

```text
User A → Slot 123
User B → Slot 123
```

I cannot simply do:

```text
1. Check whether slot is available
2. If available, create appointment
```

because both requests could do the check before either one creates the appointment.

You could end up with:

```text
User A: "Is slot available?" → Yes
User B: "Is slot available?" → Yes

User A: Create appointment
User B: Create appointment
```

Now we have a double booking.

So the availability check and booking operation need to be **atomic**.

---

# How I'd enforce it

I'd use the database to enforce the invariant.

For example, I could have a uniqueness constraint ensuring that a viewing slot can have only one active appointment.

The booking operation would happen inside a transaction.

Conceptually:

```text
Begin transaction

Create appointment for Slot 123

If another appointment already exists:
    transaction fails

Commit

```

The exact implementation would depend on the database schema and cancellation semantics.

For example, if cancelled appointments remain in the database, I might use a constraint/index that enforces uniqueness only for active appointments.

(The important point isn't the particular SQL syntax. The important point is that **the database, rather than application-level timing, enforces the invariant**.)

---

# What happens when two requests arrive?

Imagine:

```text
User A ──────┐
             ├── API ── SQL
User B ──────┘
```

Both attempt to book slot 123.

The database coordinates the concurrent writes.

One transaction succeeds.

The other encounters the uniqueness/locking constraint and fails.

The API then returns something like:

```text
201 Created
```

to the successful request and something like:

```text
409 Conflict
```

to the request that lost the race.

I'd return a message indicating that the slot is no longer available.

This is preferable to trying to coordinate the race entirely inside the API servers.

---

# Idempotency

There's another problem.

Suppose the user clicks **Book** and the request succeeds, but the response gets lost.

The client might retry:

```text
POST /viewing-slots/123/appointments
```

Without protection, we could interpret the retry as another booking attempt.

I'd therefore consider an **idempotency key** for the booking operation.

For example:

```text
Idempotency-Key: abc123
```

The server stores the result associated with that key.

If the same request is retried with the same key, the API can return the previous result rather than performing another booking operation.

(Concurrency protection and idempotency solve **different problems**. The database constraint protects against two competing bookings; idempotency protects against the same client retrying the same operation.)

---

# Notifications

After a successful booking, we probably need:

- confirmation email
- perhaps push notification
- reminder before the viewing

I wouldn't send those synchronously as part of the database transaction.

Instead:

```text
API
 ↓
SQL transaction
 ↓
Appointment successfully created
 ↓
Queue
 ↓
Notification workers
 ↓
Email / Push provider
```

The booking itself remains fast and doesn't depend on an external email provider being available.

I'd make the notification processing idempotent as well, because queues and workers can retry messages.

---

# What if the notification provider is down?

The appointment should still exist.

The notification worker can retry with exponential backoff.

For example:

```text
Attempt 1
   ↓
failure
   ↓
wait
   ↓
Attempt 2
   ↓
failure
   ↓
wait longer
   ↓
Attempt 3
```

After the maximum number of retries, I'd move the message to a **dead-letter queue** for investigation/reprocessing.

(The important reliability principle is that an email failure must not cause the booking itself to disappear.)

---

# Now let's think about scale

Suppose the API becomes heavily loaded.

I'd add multiple API servers:

```text
                       ┌── API Server 1 ──┐
Client → Load Balancer ├── API Server 2 ──┤
                       └── API Server 3 ──┘
                                ↓
                           SQL Database
```

The load balancer distributes requests across the API servers.

This gives us:

- horizontal scalability
- better availability
- ability to handle more concurrent requests

Because the API servers are stateless, any server can handle a request.

---

# Caching

Viewing-slot availability is potentially read-heavy.

For example, users constantly ask:

> "What viewing slots are available for this property?"

I could introduce a distributed cache.

```text
Client
  ↓
Load Balancer
  ↓
API Servers
  ↓
Distributed Cache
  ↓
SQL Database
```

For a cache hit, the API can return the cached slots without querying SQL.

For a cache miss, it retrieves the data from SQL and populates the cache.

However, **I would be very careful about caching booking availability**.

The cache can be slightly stale.

That's acceptable for displaying information such as:

> "These are the currently available slots."

But I would **never rely on the cache to guarantee that a slot can be booked**.

The database remains the source of truth.

(The user might see a stale "available" slot, attempt to book it, and receive a conflict. That's acceptable. Double-booking is not.)

---

# Cache invalidation

Whenever an appointment is successfully created or cancelled, I'd invalidate/update the relevant cached availability.

For example:

```text
Booking succeeds
       ↓
Database updated
       ↓
Invalidate slot availability cache
```

There is still a small window where another request could see stale data, which is another reason the database must ultimately enforce the booking invariant.

---

# Reliability

I'd design the system so that failure of individual components doesn't destroy the core data.

For example:

**Cache fails**

The API falls back to SQL.

**Notification service fails**

The appointment remains stored and the notification is retried.

**One API server fails**

The load balancer sends traffic to the remaining servers.

**A booking loses its API response**

The client can retry using its idempotency key.

**Two users book simultaneously**

The database allows only one transaction to establish the booking.

---

# Observability

I'd monitor:

### Metrics

- booking requests per second
- booking latency
- booking failure/conflict rate
- database latency
- queue depth
- notification processing latency
- cache hit rate
- API error rate

### Logs

I'd particularly want logs for:

- failed booking attempts
- concurrency conflicts
- database errors
- failed notifications
- retry attempts
- dead-lettered messages

### Traces

For example, a booking request might be traced through:

```text
Client
 ↓
Load Balancer
 ↓
API
 ↓
SQL
 ↓
Queue
 ↓
Notification Worker
 ↓
Email Provider
```

This lets us determine where latency or failures are occurring.

---

# If the system gets much larger

If SQL becomes the bottleneck, I'd investigate **why** before immediately replacing it.

I'd look at:

- slow queries
- execution plans
- indexes
- connection pool usage
- locking/blocking
- CPU
- I/O

If the workload is predominantly reads, I'd consider **read replicas** for appropriate read-only operations.

However, I'd keep the booking write path going to the authoritative database because we need strong consistency there.

(We don't want a stale replica deciding whether a viewing slot is available.)

If traffic becomes enormous, we could further separate the read/search path from the transactional booking path.

That could eventually lead to a more distributed architecture, but I wouldn't introduce that complexity until the scale and requirements justify it.

---

# Final summary to the interviewer

> "I'd start with a relational database because the core of this system is transactional booking with a strong invariant that a viewing slot can only have one active appointment. I'd begin with a stateless API and SQL database, and use database constraints and transactions to handle concurrent booking attempts rather than relying on application-level checks. I'd introduce idempotency keys to make booking retries safe. Notifications would be asynchronous through a queue so that external notification providers don't affect the booking transaction. As traffic grows, I'd add multiple API servers behind a load balancer and use a distributed cache for read-heavy availability queries, while keeping the database as the source of truth for bookings. I'd then monitor the system and introduce read replicas or further separation of the read path only if measurements showed the database had become the bottleneck."

(That is the architecture I'd defend. The important thing isn't how many boxes are on the diagram; it's being able to explain why every box exists and what requirement caused me to add it.)