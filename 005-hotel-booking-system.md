# System Design #5 — Hotel Booking System

Let's start with the requirements and keep the first version simple.

> **Design a hotel booking system where guests can search availability and book rooms without double bookings.**

Assume:

- Hotels have multiple room types.
- A room can be booked for a date range.
- Users can search availability.
- Users can make a reservation.
- We must prevent two users from booking the same room.
- Payment can be handled by an external provider.
- We need to support many hotels and many concurrent users.

---

## Step 1 — The core data

At minimum:

```text id="j0v6m2"
Hotel
 └── Room
      ├── Room Type
      └── Availability

Reservation
 ├── Guest
 ├── Room
 ├── Check-in
 ├── Check-out
 └── Status
```

A simplified database could look like:

```text id="8s8r7q"
Hotel
---------
id
name

Room
---------
id
hotel_id
type

Reservation
---------
id
room_id
guest_id
check_in
check_out
status
```

---

## Step 2 — Searching availability

The user searches:

> Hotel X, August 15–18, 2 guests

Our API receives:

```http id="f7k9de"
GET /hotels/123/availability
    ?checkIn=2026-08-15
    &checkOut=2026-08-18
```

The system needs to determine which rooms aren't already reserved during that period.

Conceptually:

```text id="7fdm3x"
Requested:
Aug 15 ───────── Aug 18

Existing reservation:
      Aug 16 ─── Aug 17
```

That room isn't available.

The important overlap condition is essentially:

```text id="w8r4gp"
existing.checkIn < requested.checkOut
AND
existing.checkOut > requested.checkIn
```

If both are true, the reservations overlap.

---

# Step 3 — Here's the real problem

Searching is relatively easy.

**Booking is where system design gets interesting.**

Imagine only one room remains:

```text id="4z6w0c"
Room 101
Available
```

Two users search at almost exactly the same time:

```text id="g0b5j8"
User A ──→ "Available ✓"
User B ──→ "Available ✓"
```

Both click **Book**.

Without protection:

```text id="t3h4k1"
User A ──→ Book Room 101 ✓
User B ──→ Book Room 101 ✓
```

We have a **double booking**.

That's unacceptable.

So the central system-design question becomes:

> **How do we guarantee that only one booking succeeds?**

And this is where we'll spend our time, because it introduces **concurrency control and consistency**—two concepts that are particularly important in booking systems.

Next we'll design the actual booking transaction and look at **database constraints vs locking vs distributed locking**.

## Step 4 — Preventing double bookings

Let's start with the simplest solution: **let the database protect us.**

Suppose two users try to book Room 101 at the same time.

We want:

```text
User A ──┐
         ├──→ Database ──→ only ONE succeeds
User B ──┘
```

### Option 1 — Database transaction + locking

When booking, we run a transaction:

```text
BEGIN TRANSACTION

1. Check room availability
2. Lock the relevant room/availability record
3. Create reservation
4. COMMIT
```

If User A gets the lock first:

```text
User A → 🔒 Room 101
User B → waits
```

User A creates the reservation:

```text
Room 101 → booked
```

Then commits:

```text
User A → COMMIT
       → 🔓 lock released
```

User B now checks availability and sees:

```text
Room 101 → already booked
```

So B gets:

```text
409 Conflict
Room no longer available
```

That's already a perfectly reasonable solution for a first version.

---

# But there's a problem with simply locking the room

Suppose we're booking:

```text
Room 101
August 15 → August 18
```

Then another guest wants:

```text
Room 101
August 20 → August 22
```

These reservations **don't overlap**.

We don't want the second request unnecessarily waiting just because someone else is booking the same physical room.

More importantly, in a real hotel we usually don't even need to assign a specific room during the initial booking.

The guest might book:

> **"Deluxe Room"**

rather than:

> **"Room 101."**

So a more realistic model is to manage **room inventory**.

---

# Step 5 — Room inventory

Suppose the hotel has:

```text
Deluxe Rooms: 10
```

For August 15–18:

```text
Available: 10
```

A booking consumes one:

```text
Available: 9
```

Another booking:

```text
Available: 8
```

And so on.

The critical operation becomes:

```text
UPDATE room_inventory
SET available = available - 1
WHERE hotel_id = 123
  AND room_type = 'DELUXE'
  AND date = '2026-08-15'
  AND available > 0;
```

Then we check how many rows were updated.

```text
1 row updated → booking can continue
0 rows updated → no availability
```

The database performs this atomically.

This is a powerful pattern because we don't have:

```text
SELECT available
      ↓
"There's 1 left!"
      ↓
UPDATE available
```

as two independent operations.

Instead, the **check and decrement happen together**.

---

## Why this is useful

Imagine the last room:

```text
Available = 1
```

Two requests arrive:

```text
User A ─┐
        ├──→ UPDATE ... WHERE available > 0
User B ─┘
```

The database serializes the updates.

Result:

```text
User A → 1 row updated → success
User B → 0 rows updated → failure
```

No double booking.

---

### One important hospitality detail

Hotels often have **inventory by room type**, not necessarily a specific physical room:

```text
Hotel
 ├── Deluxe: 10 rooms
 ├── Suite:   3 rooms
 └── Family:  5 rooms
```

So our reservation might initially say:

```text
Reservation
----------------
hotel_id
room_type_id
check_in
check_out
guest_id
status
```

The actual physical room can be assigned later.

That's actually closer to how a hotel booking platform would often work.

---

### Next

Now we have another important question:

> **What happens when a guest selects a room, goes to payment, and takes 5 minutes to complete payment?**

We can't simply hold a room forever, but we also don't want someone else to book it while the first guest is paying.

That leads us into **reservation holds, expiration, and payment consistency**—probably the most interesting part of this system.

## Step 6 — Holding inventory during payment

Suppose there's **1 Deluxe room left**.

User A clicks **Book**, but payment takes a few minutes.

We don't want User B to take that room while A is paying.

But we also don't want A to hold it forever.

### Temporary reservation hold

We introduce a temporary state:

```text id="8u1t6m"
AVAILABLE
    ↓
HELD
    ↓
CONFIRMED
```

For example:

```text id="0r9m0n"
Room type: Deluxe
Status: HELD
Held by: User A
Expires: 14:35
```

The hold might last 10 minutes.

During that period:

```text id="zxskl6"
User A → Payment
User B → "No availability"
```

---

## What if User A doesn't pay?

The hold expires:

```text id="ahf1bb"
HELD
  ↓
10 minutes
  ↓
EXPIRED
  ↓
AVAILABLE
```

We can implement this with a background job that periodically finds expired holds.

But there's an important detail:

**We shouldn't rely only on the background job to make the room available.**

Suppose the hold expires at 14:35, but our cleanup job doesn't run until 14:36.

A booking request arrives at 14:35:30.

When checking availability, we should consider:

```text
hold.expires_at > NOW()
```

rather than blindly trusting the status.

The cleanup job can then remove expired records asynchronously.

---

# Step 7 — Payment consistency

Now imagine:

```text id="m5bh8n"
Room → HELD
        ↓
Payment Provider
        ↓
Payment succeeds ✓
```

We then want:

```text id="gkq4kj"
HELD → CONFIRMED
```

But what if the payment succeeds and our system crashes before updating the reservation?

We could end up with:

```text id="7fsq8g"
Payment → SUCCESS ✓
Reservation → HELD
```

This is why we need to design the payment workflow carefully.

We give the reservation a unique ID:

```text id="x0z1q5"
Reservation: R12345
```

and use that as the reference/idempotency key with the payment provider where supported.

Then if we retry the payment operation, we don't accidentally charge the customer twice.

---

## A simplified flow

```text id="3h8w2a"
              User
               │
               ▼
        Create reservation
               │
               ▼
             HELD
               │
               ▼
           Payment
          /       \
       Success    Failure
          │          │
          ▼          ▼
      CONFIRMED    RELEASE
```

If payment fails:

```text id="x8az3m"
HELD → AVAILABLE
```

If payment succeeds:

```text id="1g3z7a"
HELD → CONFIRMED
```

---

# One more real-world issue

What if payment succeeds **after the hold has expired**?

For example:

```text id="5b6x5j"
14:30 → Room held
14:40 → Hold expires
14:41 → Payment succeeds
```

Now we have a business problem.

We shouldn't blindly confirm the reservation.

We need a rule such as:

> If the hold has expired, don't confirm automatically. Verify availability and either restore the reservation or refund the payment.

This is a good example of why **business rules are part of system design**.

The technical system can't decide everything for us.

---

### The core state machine

I'd actually draw this in an interview:

```text id="zx9h3g"
             ┌─────────────┐
             │  AVAILABLE  │
             └──────┬──────┘
                    Hold
                     ↓
             ┌─────────────┐
             │    HELD     │
             └─────┬───────┘
                /       \
        Payment OK     Expired/Failed
           ↓                ↓
    ┌─────────────┐   ┌─────────────┐
    │  CONFIRMED  │   │  AVAILABLE  │
    └─────────────┘   └─────────────┘
```

This state model is probably more important to understand than any particular AWS technology.

Next, we'll tackle **search vs booking at scale**—because availability searches can be vastly more frequent than actual bookings, and we don't want millions of searches hitting the transactional database.

I did not understand this: hold.expires_at > NOW()
Also, regarding this line: ...and use that as the reference/idempotency key with the payment provider where supported.
Then if we retry the payment operation, we don't accidentally charge the customer twice.
Is it similar to the previous system design where we "try" to avoid a notification being sent twice using the notification ID. But we can 100% prevent that until the receive request. But once we send it to the provider, we cannot guarantee no duplicates unless the provider explicitely supports it 
Anyway, you can also continue with the next thing you wanted to tackle.

Yes — your understanding of the payment/idempotency issue is **exactly right**. And there's a useful distinction between the two cases.

### First: `hold.expires_at > NOW()`

This is simply a time comparison.

Suppose the database contains:

```text
hold.expires_at = 15:00
```

And the current time is:

```text
NOW() = 14:55
```

Then:

```text
15:00 > 14:55  → TRUE
```

So the hold is **still active**.

At 15:05:

```text
15:00 > 15:05  → FALSE
```

So the hold has **expired**.

We can therefore check availability with something like:

```text
WHERE status = 'HELD'
AND expires_at > NOW()
```

Meaning:

> "Treat this hold as occupying inventory only if its expiration time is still in the future."

This protects us from depending entirely on a cleanup job to change `HELD → EXPIRED`.

---

# And yes — payment idempotency is the same fundamental issue

Your comparison with the notification system is correct.

### Our side

We can guarantee that **we don't intentionally send the same payment request twice** by tracking:

```text
Reservation ID = R12345
```

and making our own operation idempotent.

But imagine:

```text
Our system
    ↓
Payment Provider
    ↓
Payment succeeds ✓
    ↓
Network failure
    ↓
Our system doesn't receive response
```

Our system doesn't know whether the payment succeeded.

If we simply retry:

```text
Our system
    ↓
Payment Provider
    ↓
Charge again ❌
```

we could charge twice.

### If the payment provider supports idempotency keys

We send:

```text
Idempotency-Key: R12345
```

The first request:

```text
R12345 → €200 → SUCCESS
```

The retry:

```text
R12345 → €200
          ↓
Provider: "I've already processed R12345."
          ↓
Same result returned
```

Now we have much stronger protection.

So the principle is:

> **We can guarantee exactly what happens inside our system. Once we cross the boundary into an external system, exactly-once behavior depends on that system supporting the necessary semantics.**

That's a very good point to bring up in an interview.

---

# Step 8 — Search vs. Booking

Now let's look at the scalability problem.

Imagine a hotel with:

```text
100 rooms
```

But:

```text
10,000 users
```

are searching for availability at the same time.

Usually:

```text
Searches ≫ Bookings
```

We don't want every search hitting our transactional database.

### We can cache availability

For example:

```text
User
 ↓
Availability API
 ↓
Redis
 ↓
Hotel DB
```

A frequently searched hotel/date combination can be cached:

```text
hotel=123
roomType=DELUXE
dates=Aug15-18

available = 7
```

The cache makes reads very fast.

---

## But here's the important part

**We should not trust the cache to make the final booking decision.**

Imagine Redis says:

```text
7 rooms available
```

Meanwhile, seven people have just booked those rooms.

The cache might still say:

```text
7 available
```

So when the user clicks **Book**, we go to the authoritative transactional system:

```text
Search
 ↓
Redis → fast, possibly stale
```

but:

```text
Book
 ↓
Database → authoritative
```

The booking transaction performs the final availability check and decrements inventory atomically.

---

### This gives us a classic system-design pattern

> **Use a fast, potentially stale read path for search, but an authoritative strongly consistent path for writes.**

That's particularly relevant for hotel booking.

---

## One more problem: keeping Redis reasonably fresh

When a booking succeeds:

```text
Database
   ↓
inventory: 7 → 6
```

we should invalidate or update the corresponding Redis entry.

For example:

```text
DB updated
   ↓
Invalidate cache
```

Or use an event:

```text
Booking confirmed
       ↓
Event
       ↓
Cache invalidation/update
```

There can still be a short period where the cache is stale.

That's okay because **the booking database is the source of truth**.

---

### So our hotel system now has two paths

```text
                  ┌──────────────┐
                  │     User     │
                  └──────┬───────┘
                         │
                 ┌───────┴────────┐
                 ↓                ↓
              SEARCH            BOOK
                 ↓                ↓
               Redis          Database
                 ↓                ↓
             Fast read      Transaction
                                  ↓
                            Inventory
                                  ↓
                              Payment
```

This is a much more realistic architecture for a booking platform.

And this is where I think **Apaleo-style hospitality systems get particularly interesting**: availability isn't just a CRUD lookup. You have inventory, reservations, pricing, payment, cancellation, and external channels all interacting with each other.

**Next I'd tackle the harder hospitality problem: what happens when the same hotel's inventory is being sold through multiple channels at once—your own booking engine, Booking.com, Expedia, etc.?** That's where distributed systems and real-world hotel platforms get especially interesting.

## Step 9 — Multiple booking channels

Now let's make the hotel system more realistic.

Suppose a hotel has **10 Deluxe rooms**.

Those rooms are being sold through:

```text
                 Hotel
                   │
        ┌──────────┼──────────┐
        ↓          ↓          ↓
   Hotel website  Booking.com  Expedia
```

A guest books a room on Booking.com.

Our system needs to know about it, otherwise our website might still show:

> **1 room available**

when actually there are zero.

### The core problem

Inventory can change from multiple places:

```text
Hotel website ────┐
Booking.com ──────┼──→ Hotel inventory
Expedia ──────────┘
```

We need a **single source of truth** for inventory.

---

# The typical solution: a central inventory/reservation system

Instead of allowing every channel to independently maintain inventory:

```text
Booking.com → its own inventory
Expedia     → its own inventory
Website     → its own inventory
```

we have:

```text
                    ┌───────────────┐
Hotel website ────→ │               │
Booking.com ──────→ │ Central       │
Expedia ──────────→ │ Reservation & │
                    │ Inventory     │
                    │ System        │
                    └───────────────┘
```

This is conceptually close to what a hospitality platform such as **Apaleo** is dealing with.

---

## What happens when Booking.com makes a reservation?

Suppose:

```text
Deluxe rooms = 10
Available = 1
```

Booking.com sends:

```text
"Reserve 1 Deluxe room"
```

Our system:

```text
1. Check inventory
2. Atomically decrement availability
3. Create reservation
4. Return confirmation
```

Now:

```text
Available = 0
```

We can then propagate that change to the other channels.

```text
Central system
      │
      ├──→ Hotel website: 0 available
      ├──→ Booking.com: 0 available
      └──→ Expedia: 0 available
```

---

# But there's a problem: propagation isn't instantaneous

Suppose:

```text
10:00:00 — Last room booked on Booking.com
10:00:01 — Website still thinks 1 room exists
10:00:02 — Website receives update
```

For those two seconds, the website has **stale inventory**.

Someone could try to book it.

So again:

> **Cached/channel availability is not the final authority.**

The central reservation system must perform the final atomic inventory check.

If the website tries to book the room while another channel has already taken it:

```text
Website
   ↓
Central system
   ↓
available > 0 ?
   ↓
NO
   ↓
Booking rejected
```

The website then tells the guest:

> "Unfortunately, this room is no longer available."

Not ideal, but **correct**.

---

# Step 10 — What if two channels book simultaneously?

This is where our previous concurrency discussion becomes useful.

Suppose there is exactly **one room left**:

```text
Available = 1
```

At almost the same moment:

```text
Booking.com ──┐
              ├──→ Central inventory
Website ──────┘
```

Both try:

```text
available = available - 1
WHERE available > 0
```

The database guarantees that only one succeeds.

```text
Booking.com → 1 row updated → SUCCESS
Website     → 0 rows updated → SOLD OUT
```

No double booking.

---

# But now we have another distributed-system problem

What if Booking.com successfully creates the reservation, but our response back to Booking.com gets lost?

```text
Booking.com
    ↓
Our system
    ↓
Reservation created ✓
    ↓
Network failure
    ↓
Booking.com: "Did it work?"
```

Booking.com may retry:

```text
Booking.com
    ↓
Same reservation request
```

We don't want to create another reservation.

So once again, **idempotency** matters.

The external booking channel should provide a unique reservation/reference ID.

For example:

```text
channel = BOOKING_COM
externalReservationId = ABC123
```

Our database can enforce uniqueness:

```text
UNIQUE(channel, externalReservationId)
```

If Booking.com sends `ABC123` again:

```text
Already exists → return existing reservation
```

This is something we **can guarantee on our side**, unlike exactly-once behavior at the external provider boundary.

---

# This is getting close to real hospitality architecture

We now have:

```text
                   ┌─────────────────┐
                   │ Central Booking │
                   │ & Inventory     │
                   └────────┬────────┘
                            │
               ┌────────────┼────────────┐
               ↓            ↓            ↓
          Hotel website  Booking.com   Expedia
```

And internally:

```text
                     Booking API
                         │
             ┌───────────┼───────────┐
             ↓           ↓           ↓
         Inventory   Reservation   Payment
             │
             ↓
        PostgreSQL
             │
             ↓
          Redis
       (fast search)
```

The key principle is:

> **Channels can have stale views of availability, but the central reservation system is authoritative.**

---

## One important real-world complication

Hotels don't always want to expose **all** their inventory to every channel.

For example:

```text
Deluxe rooms = 10

Website      → 5 rooms
Booking.com  → 3 rooms
Expedia      → 2 rooms
```

This is called an **allocation/allotment** model.

Alternatively, channels may draw from a shared pool:

```text
Total inventory = 10
All channels compete for those 10 rooms
```

Which model we use is a **business decision**, not purely a technical one.

---

### Why I like this problem for your preparation

We've now touched several senior-level concepts naturally:

- Transactions
- Concurrency
- Inventory management
- Idempotency
- Caching
- Eventual consistency
- External integrations
- Rate/availability propagation
- Multiple consumers/channels
- Source of truth

And importantly, these aren't random concepts we're adding to make the design complicated. **They're emerging from the actual hospitality problem.**

Next, I'd tackle **cancellations and modifications**. That's where inventory gets returned, payments may need refunds, and the system has to deal with another set of consistency problems.

## Step 11 — Cancellation and modification

Now a guest wants to cancel:

```text
Reservation: R123
Room type: Deluxe
Aug 15–18
Status: CONFIRMED
```

The simplest flow is:

```text
Cancel request
     ↓
Validate cancellation policy
     ↓
Cancel reservation
     ↓
Return inventory
     ↓
Process refund if applicable
```

So:

```text
CONFIRMED
    ↓
CANCELLED
    ↓
Inventory +1
```

### But again, payment creates a consistency problem

Suppose we do:

```text
1. Cancel reservation ✓
2. Return inventory ✓
3. Refund payment ✗
```

Now the room is available again, but the customer hasn't received their refund.

So we shouldn't think of cancellation as one database transaction covering everything. The payment provider is external.

A practical approach is to track the refund separately:

```text
Reservation
  status: CANCELLED

Refund
  status: PENDING
```

Then a worker handles the refund:

```text
PENDING
   ↓
Payment Provider
   ↓
SUCCESS → REFUNDED
```

If the provider fails temporarily, we retry.

---

## Modification is more interesting

Suppose the guest changes:

```text
Aug 15–18
```

to:

```text
Aug 20–23
```

We shouldn't simply modify the existing dates without checking availability.

Conceptually:

```text
Existing reservation
       ↓
Check new dates
       ↓
Available?
   /          \
 Yes           No
 ↓              ↓
Move booking   Reject change
 ↓
Release old inventory
```

And this needs to be done carefully so we don't temporarily lose or double-count inventory.

For example, if the new dates have only one room left, we should **reserve the new inventory before releasing the old inventory**.

Otherwise:

```text
Release old room
     ↓
Another booking takes new room
     ↓
New reservation fails
```

The guest could end up with neither reservation.

So for a modification:

> **Secure the new inventory first, then release the old inventory.**

---

## One more important concept: cancellation policy

The technical system also needs to know the hotel's business rules.

For example:

```text
Cancel > 48 hours before check-in → full refund
Cancel < 48 hours → 50% refund
No-show → no refund
```

These rules affect both:

- whether the reservation can be cancelled
- how much money should be refunded

So I'd keep the **reservation state** separate from the **payment/refund state**.

---

At this point, I think we've covered enough of the booking system. The important mental model is:

```text
Search
  ↓
Cached availability

Booking
  ↓
Authoritative inventory
  ↓
Temporary hold
  ↓
Payment
  ↓
Confirmed

Cancellation
  ↓
Release inventory
  ↓
Refund asynchronously
```

I'd move on from this problem now rather than adding more edge cases. The next system-design problem should be from a **different domain**, so you're practicing the underlying patterns rather than memorizing hotel booking.