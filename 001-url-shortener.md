# System Design #1 — Design a URL Shortener

**Interview question:**

> "Design a URL shortening service like Bitly."

Why start here?

Because it teaches the basic building blocks without overwhelming you:

- API design
- database choice
- scalability
- caching
- unique IDs
- reliability
- observability
- trade-offs

And later, we'll reuse the same thinking for **the much more relevant question: designing an internal developer platform**.

---

# Step 1 — Don't start designing

This is important for your interview.

If they say:

> "Design a URL shortening service."

Don't immediately start talking about databases and AWS.

First clarify the requirements.

You could say:

> "Sure. Before I design it, I'd like to clarify a few requirements."

Then ask:

### Functional requirements

**1. What should the system do?**

> "Users should be able to submit a long URL and get a short URL."

> "When someone opens the short URL, they should be redirected to the original URL."

That's enough for now.

### Non-functional requirements

Then ask:

> "Do we have any scale requirements?"

Suppose the interviewer says:

> 10 million users and 100 million redirects per day.

Now you have something to design around.

---

# Step 2 — Identify the two main operations

We have only two important operations:

### Create short URL

```text
POST /urls

longUrl → shortUrl
```

Example:

```text
https://www.example.com/products/12345
                 ↓
https://short.ly/a8K2x
```

### Redirect

```text
GET /a8K2x

a8K2x → original URL
```

That's it.

Notice how simple the system is.

---

# Step 3 — Think about the data

We need to remember:

```text
shortCode → longUrl
```

So our database could contain:

| short_code | long_url |
|---|---|
| a8K2x | https://example.com/products/12345 |
| B72pq | https://google.com/... |

We could also store:

```text
created_at
expires_at
user_id
```

if those requirements exist.

---

# Step 4 — Choose a database

For this problem, I'd start with a **relational database**.

You could say:

> "I'd start with a relational database because the data model is simple and we need reliable persistence. We can scale it later if needed."

Don't immediately say:

> "We'll use DynamoDB because AWS."

That's over-engineering.

Start simple.

---

# Step 5 — Design the basic architecture

Now we can draw:

```text
                ┌─────────────┐
                │   Client    │
                └──────┬──────┘
                       │
                       ▼
                ┌─────────────┐
                │ Load Balancer│
                └──────┬──────┘
                       │
                       ▼
                ┌─────────────┐
                │  URL API     │
                └──────┬──────┘
                       │
                       ▼
                ┌─────────────┐
                │  Database    │
                └─────────────┘
```

That's your **first version**.

Don't add Redis, Kafka, Kubernetes, 17 microservices, etc.

---

# Step 6 — Now think about scale

Here's where the interview gets interesting.

The interviewer says:

> "We have 100 million redirects per day."

Ask yourself:

**Are reads or writes more frequent?**

Clearly:

```text
Create URL → relatively rare

Redirect → extremely frequent
```

So we have a **read-heavy system**.

That immediately suggests:

> **Caching.**

---

# Step 7 — Add a cache

Now:

```text
                  ┌─────────────┐
                  │   Client    │
                  └──────┬──────┘
                         │
                         ▼
                  ┌─────────────┐
                  │ Load Balancer│
                  └──────┬──────┘
                         │
                         ▼
                  ┌─────────────┐
                  │  URL API     │
                  └──────┬──────┘
                         │
                    ┌────▼────┐
                    │  Cache  │
                    └────┬────┘
                         │
                         │ cache miss
                         ▼
                  ┌─────────────┐
                  │  Database    │
                  └─────────────┘
```

The redirect flow becomes:

```text
GET /a8K2x
      ↓
Check cache
      ↓
Found?
 ┌────┴────┐
Yes        No
 ↓          ↓
URL       Database
 ↓          ↓
Redirect   Cache
            ↓
         Redirect
```

This is a very common system-design pattern.

---

# Step 8 — Talk about reliability

Now bring in another keyword from your JD:

**reliability.**

You could say:

> "The cache should improve performance, but the database remains the source of truth. If the cache is unavailable, we should still be able to serve requests from the database."

That's a good senior-level statement.

You're showing that you're not blindly depending on the cache.

---

# Step 9 — Observability

And now you can naturally bring in another JD keyword:

**observability.**

I'd monitor:

### Metrics

- request latency
- requests per second
- error rate
- cache hit rate
- database latency

### Logs

- failed redirects
- invalid URLs
- database errors

### Traces

Useful when a request goes through:

```text
Load Balancer
 → API
 → Cache
 → Database
```

You can say:

> "I'd add metrics, logs, and tracing so we can identify where latency or failures are coming from."

You don't need to say "Datadog" unless they ask what tooling you'd use.

---

# Step 10 — Your final answer

If the interviewer asks:

> **"Design a URL shortener."**

Your high-level answer could eventually become:

> "I'd start with a simple API and a relational database. We need two main operations: creating a short URL and redirecting a short URL to the original URL.
>
> Since redirects are much more frequent than URL creation, this will be a read-heavy system. I'd put a cache in front of the database to handle frequent redirects and reduce database load.
>
> I'd keep the database as the source of truth. The API would sit behind a load balancer so we can scale horizontally.
>
> For reliability and observability, I'd monitor latency, error rates, cache hit rate, and database performance, and I'd use logs and tracing to troubleshoot problems.
>
> I would start with this simple architecture and only introduce more components if the scale or requirements justify them."

**That is already a good senior-level answer.**

Notice what you *didn't* do:

❌ 15 microservices  
❌ Kafka  
❌ Kubernetes  
❌ 5 different databases  
❌ complicated AWS architecture  

You started with **simplicity**, then added complexity because the requirements justified it.

That is exactly the mindset I want you to develop for this interview.

---

I have 2 doubts. Why did you say with just 6 characters... wouldn't it depend upon the ID from the db? 
And how do you make sure the distributed databases don't generate the same ID. Actually, I have some doubts about the distributedness also. Are there duplicate copies in different locations? Which is eventually copied? So if the data resides in EU, and a server in USA doesn't find it, does the request internally go to EU?

## 1. "With 6 characters" — doesn't it depend on the DB ID?

**Yes.** You caught an important nuance.

When I said:

> 6 characters gives us 62⁶ possible codes

I meant **the address space**, not that every database ID will fit into 6 characters.

Base62 maps a number to a string:

```text
DB ID       Base62
1           1
10          A
100         1C
1,000       G8
...
```

The largest number representable with 6 Base62 characters is:

**62⁶ − 1 ≈ 56.8 billion**

So:

- ID `12345` → maybe 3 characters
- ID `1,000,000` → maybe 4 characters
- ID `56 billion` → 6 characters
- ID above ~56.8 billion → **7 characters**

So you're absolutely right: **the length depends on the ID you're encoding.**

And there's another interesting point:

### What if we have 100 billion URLs?

Then 6 characters aren't enough.

We'd need 7 characters:

```text
62⁷ ≈ 3.5 trillion
```

That's one of the trade-offs you could mention in an interview.

---

# 2. Distributed IDs — how do you prevent duplicates?

This is where things get more interesting.

Imagine we have:

```text
           USA DB
           ID = 100
           
           EU DB
           ID = 100
```

If both independently use auto-increment IDs, **we have a collision**.

So you can't simply say:

> "I'll put an auto-increment ID in every database."

You need a **globally unique ID generation strategy**.

One approach is a Snowflake-style ID.

Instead of:

```text
100
101
102
```

you generate something based on:

```text
timestamp
+
machine/region ID
+
sequence number
```

For example, conceptually:

```text
┌──────────────┬───────────┬──────────┐
│ timestamp    │ server ID │ sequence │
└──────────────┴───────────┴──────────┘
```

Suppose:

```text
USA server = 01
EU server  = 02
```

At the same moment:

```text
USA → timestamp + 01 + 001
EU  → timestamp + 02 + 001
```

Those IDs are different.

The server IDs are assigned so that **two generators don't use the same identity**.

That's how you can generate IDs independently without collisions.

---

# 3. But do we even need distributed ID generation?

**Not necessarily.**

This is where I want you to develop a good interview instinct.

We started with:

> "Let's use a database auto-increment ID."

Then we asked:

> "What if we have multiple databases?"

That's when distributed ID generation becomes relevant.

But if our system has:

```text
              API
               │
        ┌──────▼──────┐
        │ One primary │
        │  database   │
        └─────────────┘
```

then the database can safely generate the IDs.

You don't need Snowflake.

**Don't introduce distributed IDs just because they're impressive.**

---

# 4. Your EU/USA example

You asked:

> "If the data resides in EU, and a server in USA doesn't find it, does the request internally go to EU?"

**It can. But that's not necessarily how you'd design it.**

Let's say we have replicated databases:

```text
              URL API
             /       \
            /         \
         EU DB       US DB
          │            │
       replica       replica
```

If the user in the US requests:

```text
short.ly/abc123
```

the US server would normally check the **US replica**.

If the data has been replicated there:

```text
US → finds abc123 → redirect
```

No trip to Europe.

---

## What if it isn't there?

There are several possibilities.

### Option A — Route to the correct region

The system knows where the data belongs:

```text
abc123 → EU
```

So:

```text
US request
     ↓
US service
     ↓
EU service/database
     ↓
URL
```

This creates cross-region latency.

---

### Option B — Replicate the data

You could instead replicate the URL mappings:

```text
                Primary
                   EU
                  /  \
                 /    \
               US      Asia
```

Then the US request is served locally.

This gives better latency but introduces:

- replication cost
- replication delay
- consistency concerns

---

# 5. There's another possibility: cache

And for our URL shortener, this is actually very interesting.

Remember our architecture:

```text
Client
  ↓
API
  ↓
Cache
  ↓
Database
```

Suppose `abc123` is extremely popular.

The US server may have:

```text
Cache:
abc123 → https://example.com
```

Then it doesn't care where the database is.

```text
US user
   ↓
US API
   ↓
US cache
   ↓
redirect
```

**No database request at all.**

That's one reason caching is so useful for read-heavy systems.

---

# 6. And this brings us back to our URL shortener

We're starting to see why system design is about **trade-offs** rather than drawing boxes.

We started with:

```text
API → Database
```

Then:

**Problem:** too many reads.

→ Add cache.

```text
API → Cache → Database
```

Then:

**Problem:** system needs to serve users globally.

→ Consider regional deployment/replication.

Then:

**Problem:** multiple regions generate IDs.

→ Consider globally unique ID generation.

And now you have something like:

```text
                    Global users
                         │
              ┌──────────┴──────────┐
              ▼                     ▼
           EU Region             US Region
              │                     │
          API + Cache           API + Cache
              │                     │
              └──────────┬──────────┘
                         ▼
                  Data storage
```

But **we don't automatically need all of this**.

If the interviewer says:

> "This service only needs to operate in Europe."

Then suddenly:

- no US region
- no cross-region replication
- no distributed ID generation

might be necessary.

And that's exactly what you should say in an interview:

> **"I'd start with a single region because that's sufficient for the current requirements. If we later need global availability, we can introduce regional replicas and a distributed ID generation strategy."**

That's a **much stronger senior answer** than immediately drawing a globally distributed architecture.

**distributed systems can help with scalability, fault tolerance, reliability, and availability.** They aren't synonymous with any one of them.

For example, Immowelt might have:

```text
                    Users
                      ↓
               Load Balancer
                /          \
               ↓            ↓
            Server        Server
               \            /
                \          /
                 Cache
                   ↓
              DB Primary
              /        \
             ↓          ↓
        Replica EU   Replica EU
```

Here:

- multiple servers → **horizontal scaling**
- multiple DB copies → **replication**
- multiple instances → **fault tolerance**
- cache → **performance/scalability**
- database replicas → **read scalability**

We **started introducing distributed concepts**, but we didn't actually finish the distributed design of the URL shortener. Let's do that now, step by step.

And this is useful because it brings together everything you just learned.

# URL Shortener — Making It Distributed

Let's start with our simple design:

```text
User
  ↓
API Server
  ↓
Cache
  ↓
Database
```

This works. But now suppose the service becomes very popular.

We have two different problems to solve:

1. **The API servers need to handle more traffic.**
2. **We don't want the service to go down if one server fails.**

---

## Step 1 — Multiple API servers

We can add more API servers:

```text
                    Users
                      ↓
                Load Balancer
                 /     |     \
                ↓      ↓      ↓
             API 1   API 2   API 3
                \      |      /
                 \     |     /
                    Cache
                      ↓
                   Database
```

Now we've distributed the **application layer**.

This gives us:

### Scalability

If traffic increases:

```text
3 servers → 6 servers → 10 servers
```

### Fault tolerance

If API 2 dies:

```text
API 1 ✓
API 2 💥
API 3 ✓
```

The load balancer sends traffic to API 1 and API 3.

---

# Step 2 — What about the cache?

If each API server has its own memory cache:

```text
API 1 → Cache 1
API 2 → Cache 2
API 3 → Cache 3
```

we have a problem.

Suppose:

```text
abc123 → https://example.com
```

is cached in API 1.

A request goes to API 2.

API 2 doesn't know about it.

So we could use a **shared distributed cache**, such as Redis:

```text
             API 1 ──┐
             API 2 ──┼──→ Redis
             API 3 ──┘
```

Now all API servers can access the same cache.

That's another form of distribution.

---

# Step 3 — What about the database?

This is where things get more interesting.

Our database might initially be:

```text
API servers
     ↓
   Database
```

But the database eventually becomes a bottleneck.

One option is **replication**:

```text
              Primary DB
              /        \
             ↓          ↓
        Read Replica  Read Replica
```

Writes go to the primary:

```text
POST /urls
     ↓
Primary DB
```

Reads can potentially go to replicas:

```text
GET /abc123
     ↓
Read Replica
```

This is called **read replication**.

---

# Step 4 — What if the primary database fails?

Now we're thinking about **fault tolerance**.

We could have:

```text
              Primary DB
                  │
            replication
                  ↓
             Replica DB
```

If the primary fails, we can promote the replica.

Now we have a highly available database.

There are complexities around:

- replication lag
- failover
- consistency

But conceptually that's the idea.

---

# Step 5 — What about different geographical regions?

Now let's say or our fictional URL service needs to serve users globally.

We could have:

```text
                 Global Users
                      │
             Global Load Balancer
                /            \
               ↓              ↓
           EU Region       US Region
              │               │
        ┌─────┴─────┐   ┌─────┴─────┐
        │API API API│   │API API API│
        └─────┬─────┘   └─────┬─────┘
              │               │
           Cache            Cache
              │               │
             DB              DB
```

Now we have **geographical distribution**.

Users can be routed to a nearby region.

This can reduce latency.

---

# Step 6 — But what happens to the data?

Now we have to make a major design decision.

### Option A — One primary region

For example:

```text
             EU Database
                  ↑
            source of truth

US API ───────────┘
```

US requests may need to communicate with Europe.

Advantages:

- simpler
- easier consistency

Disadvantages:

- higher latency for US users
- EU becomes a dependency
- potentially a single-region failure

---

### Option B — Replicate data across regions

```text
              EU DB
             ↕
          replication
             ↕
              US DB
```

Now both regions have copies.

Advantages:

- low latency
- better availability
- regional failure doesn't necessarily take down the service

Disadvantages:

- replication complexity
- replication delay
- consistency problems
- more infrastructure

---

# Step 7 — Now we have the ID problem

Remember your earlier question.

If both regions create URLs:

```text
EU → ID 100
US → ID 100
```

we can't simply Base62-encode those IDs.

We need a **globally unique ID**.

We could use:

- distributed ID generator
- UUID
- Snowflake-style IDs
- region-specific ID ranges

For example:

```text
EU IDs: 1–999,999,999
US IDs: 1,000,000,000–1,999,999,999
```

Then:

```text
EU → 123
US → 1,000,000,123
```

No collision.

Or use a Snowflake-style generator where the ID contains information such as:

```text
timestamp + region/server ID + sequence
```

Then Base62-encode the result.

---

# So what does our final distributed URL shortener look like?

Something like:

```text
                         GLOBAL USERS
                              │
                    Global Load Balancer
                       /             \
                      ↓               ↓
                 EU REGION        US REGION
                 /   |   \        /   |   \
                API API API      API API API
                  \   |   /        \   |   /
                   Redis             Redis
                      │               │
                     DB  ←──replication──→ DB
```

And our IDs are globally unique.

---

# But here's the most important interview lesson

**We wouldn't necessarily build this entire thing.**

If the interviewer says:

> "The service only needs to operate in Germany."

Then your design could simply be:

```text
Load Balancer
      ↓
API × 3
      ↓
Redis
      ↓
Primary DB
      ↓
Replica
```

No US region.

No cross-region replication.

No complicated distributed ID generator.

If they then say:

> "Now we need to support the US."

**That's when you evolve the design.**

You can say:

> "For the initial design, I'd keep a single region because it keeps the architecture simple. If we need global availability and lower latency later, I'd introduce regional deployments, replicated data, and a globally unique ID-generation strategy."

That is **exactly the style of thinking I want you to practice for the system-design section**:

**Start simple → identify the bottleneck → introduce distribution where it solves a real problem → explain the trade-off.**

I still have doubts about the base62 unique id generation. You said we could use 6 characters and it would give us so many billion url's. Let's say if ID = 1, Then to base62 it stays /1 or /000001 ?

Yes — **both are possible**, but normally you'd use **`/1`**, not `/000001`.

The key point is that **Base62 conversion doesn't inherently produce a fixed-length string**.

### Example

Suppose our database gives us sequential IDs:

| DB ID | Base62 | Short URL |
|---:|---|---|
| 1 | `1` | `short.ly/1` |
| 10 | `A` | `short.ly/A` |
| 61 | `z` | `short.ly/z` |
| 62 | `10` | `short.ly/10` |
| 63 | `11` | `short.ly/11` |
| 100 | `1C` | `short.ly/1C` |
| 1000 | `G8` | `short.ly/G8` |

So:

```text
ID = 1
       ↓
Base62
       ↓
"1"
```

The resulting URL is:

```text
https://short.ly/1
```

Not:

```text
https://short.ly/000001
```

---

### So where did the "6 characters = 56.8 billion" come from?

That's the **maximum capacity if we allow up to 6 characters**.

Think of it like this:

```text
1 character → 62 possibilities
2 characters → 62²
3 characters → 62³
4 characters → 62⁴
5 characters → 62⁵
6 characters → 62⁶
```

So a 6-character code can represent:

**62⁶ = 56,800,235,584 different values**

But you don't have to use all six characters.

You start with:

```text
1
2
3
...
z
10
11
...
```

and eventually:

```text
zzzzz
zzzzzz
```

So the URLs **gradually get longer as the ID increases**.

---

### One subtle issue

If you literally use the database's auto-increment ID, then your first URLs are very predictable:

```text
short.ly/1
short.ly/2
short.ly/3
short.ly/4
...
```

Someone could easily guess other URLs.

That's one reason a real production URL shortener might choose a different ID-generation strategy, such as a random or non-sequential ID, even though **sequential ID + Base62 is a very simple and good starting design**.

And this is exactly the kind of follow-up an interviewer might ask:

> **"What are the disadvantages of using sequential IDs?"**

We can tackle that next, because it leads naturally into **collision avoidance, security, and why you'd choose Base62 at all.**

The next problem is **how to generate the ID safely**, especially once we have multiple API servers or regions.

We have established:

```text
Database ID → Base62 → short code
```

For example:

```text
12345 → 3D7 → short.ly/3D7
```

Now imagine **three API servers receive requests at almost exactly the same time**:

```text
API 1 ──┐
API 2 ──┼──→ ??? ID generation
API 3 ──┘
```

We need to guarantee that two requests **never receive the same ID**.

So the next question is:

> **How do we generate unique IDs safely when many requests arrive concurrently?**

We'll look at three approaches:

1. **Database auto-increment** — simplest
2. **Random IDs** — simple but introduces collision handling
3. **Distributed ID generation** — useful when the system grows across multiple servers/regions

And we'll build them up **one at a time**, starting with the database auto-increment approach.

### Next step

Before we move to a more complicated design, I suggest we do **Step 11: unique short-code generation**. That's where this seemingly simple system starts introducing interesting design decisions—**database IDs vs hashing vs random strings**, collisions, concurrency, and scalability.

There are several approaches.

### Option 1 — Auto-incrementing ID + Base62

This is probably the simplest practical approach.

Suppose the database gives us:

```text
ID = 125
```

We convert `125` to Base62:

```text
125 → cb
```

So:

```text
https://example.com/very/long/url
              ↓
ID = 125
              ↓
Base62(125) = cb
              ↓
https://short.ly/cb
```

Base62 uses:

```text
a-z
A-Z
0-9
```

So there are **62 possible characters**.

With just 6 characters, you get:

```text
62^6 ≈ 56.8 billion
```

possible codes.

That's plenty for many systems.

---

### Option 2 — Random string

Generate something like:

```text
a8K2xP
```

and check whether it already exists.

If it does:

```text
generate again
```

This is simple, but you have to deal with **collisions**.

---

### Option 3 — Hash the original URL

For example:

```text
SHA-256(longUrl)
```

and take part of the hash.

The problem is that:

- hashes are much longer than we need
- truncating them introduces collision risk
- the same URL always produces the same code unless you add additional information

So I wouldn't choose this as my first solution.

---

### Option 4 — Distributed ID generator

At very large scale, you could use something like:

```text
Snowflake-style ID
        ↓
Base62
        ↓
short code
```

This avoids depending on a single database's auto-increment sequence.

But again, **don't start here**.

---

## What would I say in the interview?

I'd start with:

> "For generating the short code, I'd probably start with a unique numeric ID and encode it using Base62. This gives us a short, URL-safe identifier and avoids collision checks."

Then the interviewer might ask:

> "But doesn't that expose how many URLs you've created?"

Yes.

That's a trade-off.

You could then say:

> "Yes. If sequential IDs are a security or privacy concern, I'd use a random or distributed ID instead. I'd choose based on the requirements."

**That's system design thinking.**

You're not trying to find *the* correct answer. You're showing that you understand the **trade-offs**.

# Step 11 — Unique ID generation: Database auto-increment

We currently have:

```text
POST /urls
      ↓
API Server
      ↓
Database
      ↓
ID = ?
      ↓
Base62
      ↓
Short URL
```

The easiest solution is to let the database generate the ID.

### Example

Our table:

```text
Url
-------------------------
Id          BIGINT
ShortCode   VARCHAR(...)
LongUrl     VARCHAR(...)
```

And `Id` is an auto-incrementing primary key.

So requests might produce:

```text
Request 1 → ID 1
Request 2 → ID 2
Request 3 → ID 3
Request 4 → ID 4
```

Then:

```text
1 → Base62 → "1"
2 → Base62 → "2"
3 → Base62 → "3"
```

### But what if requests arrive simultaneously?

Suppose:

```text
Request A ──┐
Request B ──┼──→ Database
Request C ──┘
```

You don't have to manually synchronize them in your application.

The database's auto-increment mechanism guarantees that each inserted row gets a different ID.

So you might get:

```text
A → 101
B → 102
C → 103
```

even though they arrived at roughly the same time.

That's one of the nice things about letting the database handle uniqueness.

---

## But there's a limitation

Now imagine we have **multiple database servers**:

```text
             API
           /     \
          ↓       ↓
       EU DB    US DB
```

Both databases have their own auto-increment sequence.

They could produce:

```text
EU DB → 101
US DB → 101
```

Now we have:

```text
Base62(101) → "1D"
Base62(101) → "1D"
```

**Collision.**

That's why a single database with auto-increment is simple, but it doesn't automatically solve globally distributed ID generation.

---

### So our design currently looks like this:

**Single region:**

```text
          API servers
         /    |    \
        ↓     ↓     ↓
             DB
              │
        Auto-increment ID
              │
            Base62
```

This is perfectly reasonable for a first version.

If the interviewer asks:

> "How would you scale this globally?"

**Then** we need to replace the simple ID generation approach.

And that's where **distributed IDs** come in.

Next we'll look at **random IDs**, because it's a surprisingly simple alternative and it introduces the important concept of **collision probability**.

## Step 12 — Random IDs

Instead of asking the database for:

```text
1 → 2 → 3 → 4 → 5...
```

we could generate the short code ourselves.

For example:

```text
a8K2xP
```

Then store:

```text
short_code = a8K2xP
long_url   = https://example.com/...
```

### The problem

What if we randomly generate:

```text
a8K2xP
```

and that code already exists?

We have a **collision**.

So the process becomes:

```text
Generate random code
       ↓
Does it exist?
   /       \
 Yes        No
 ↓           ↓
Generate    Save it
again
```

For example:

```text
Generate → a8K2xP
             ↓
          Already exists
             ↓
Generate → 7Hd91Q
             ↓
          Doesn't exist
             ↓
           Save
```

---

## Why can random IDs work?

Because we have a huge number of possible combinations.

If we use **6 Base62 characters**:

```text
62⁶ ≈ 56.8 billion
```

possible codes.

So if our system has only a few million URLs, the chance of a collision for any *single* generated code is fairly small.

But there's an important catch.

### The birthday problem

As the number of existing URLs increases, the probability of **some collision somewhere** becomes much higher than you might intuitively expect.

You don't need to understand the mathematics for the interview unless they ask.

Just understand:

> More URLs → more chance of collision.

That's why we **always check the database** before accepting a randomly generated code.

---

# Why might we prefer random IDs?

Compared with sequential IDs:

### Sequential

```text
1
2
3
4
5
```

Advantages:

- Simple
- No collision handling
- Easy to generate with a DB

Disadvantages:

- Predictable
- Can reveal how many URLs have been created
- Requires a central ID generator/database

### Random

```text
a8K2xP
7Hd91Q
k92LmA
```

Advantages:

- Harder to guess
- Doesn't require sequential IDs
- Multiple servers can generate them independently

Disadvantages:

- Collisions are possible
- Must check for collisions
- Need enough random space

---

## One important distinction

Random generation doesn't mean:

> "We don't need a database."

We **still need the database** to store the mapping:

```text
a8K2xP → https://example.com/something
```

And the database is what ultimately tells us:

> "Yes, `a8K2xP` already exists."

---

## What would I choose?

For our interview design, I'd say:

> "For a simple implementation, I could use a database-generated ID and Base62 encoding. If we need non-sequential IDs or multiple independent ID generators, I could use random codes. In that case, I'd check for collisions before storing the mapping."

That's a perfectly reasonable answer.

But now we reach the interesting question:

> **If we have thousands of servers generating IDs simultaneously, can we generate unique IDs without constantly checking for collisions?**

That's where **distributed ID generation** becomes useful.

And that's the next step.

## Step 13 — Distributed ID generation

Now let's make the problem slightly bigger.

Imagine our URL shortener has:

```text
                  Users
                    ↓
              Load Balancer
             /      |      \
            ↓       ↓       ↓
          API 1   API 2   API 3
```

All three servers need to create IDs.

We don't want:

```text
API 1 → 123
API 2 → 123  ❌
API 3 → 123  ❌
```

And we don't want every server constantly asking a central database:

> "Give me the next ID."

That database could become a bottleneck.

So we need a way for multiple machines to generate IDs **independently while guaranteeing uniqueness**.

---

# A simple idea: give each server a unique range

Suppose we have three servers.

We could assign:

```text
API 1 → IDs 1–1,000,000
API 2 → IDs 1,000,001–2,000,000
API 3 → IDs 2,000,001–3,000,000
```

Then:

```text
API 1 → 1
API 1 → 2
API 1 → 3

API 2 → 1,000,001
API 2 → 1,000,002

API 3 → 2,000,001
```

No collision is possible.

Each server can generate IDs locally.

### But there's a problem

What happens when API 1 uses all of its IDs?

It needs another range.

So we still need some central mechanism to assign ranges.

This can work, but it's not ideal for a huge distributed system.

---

# A more common approach: Snowflake-style IDs

A Snowflake-style ID combines several pieces of information.

Conceptually:

```text
┌────────────┬───────────┬──────────┐
│ Timestamp  │ Server ID │ Sequence │
└────────────┴───────────┴──────────┘
```

For example:

```text
Timestamp     = 2026-08-09 15:20:31
Server ID     = 42
Sequence      = 001
```

Together, these produce one large numeric ID.

Another server might generate:

```text
Timestamp     = 2026-08-09 15:20:31
Server ID     = 43
Sequence      = 001
```

Even though the timestamp and sequence are the same, the **server ID is different**, so the resulting IDs are different.

---

# Why the sequence number?

Imagine API server 42 receives 3 requests at exactly the same time:

```text
API 42

Request A → sequence 001
Request B → sequence 002
Request C → sequence 003
```

So they still get different IDs.

The combination of:

**timestamp + server ID + sequence**

gives us a way to generate unique IDs locally.

---

# Then Base62 comes back

The resulting Snowflake-style ID is usually a fairly large number.

For example, imagine:

```text
ID = 7283749283749283
```

We don't want that in our URL:

```text
short.ly/7283749283749283
```

So we convert it to Base62:

```text
7283749283749283
          ↓
      Base62
          ↓
short.ly/8Hd92Kx
```

Now we have a short code.

---

# One important correction to our earlier discussion

You might now be wondering:

> "Didn't we say 6 Base62 characters can represent 56 billion values? What happens with a huge Snowflake ID?"

Exactly.

**Six characters may not be enough.**

A Snowflake ID might require:

```text
7 characters
8 characters
9 characters
...
```

depending on the numeric range.

So you don't say:

> "I'll always have a 6-character URL."

You say:

> "I'll use Base62 to make the numeric ID more compact."

The length depends on the ID value.

---

# What have we learned?

We now have three approaches:

### 1. Database auto-increment

```text
DB → 1, 2, 3, 4...
      ↓
    Base62
```

**Simple.**

Good for a relatively centralized system.

---

### 2. Random ID

```text
API → random code
          ↓
      check DB
          ↓
       collision?
```

**Simple and unpredictable**, but requires collision checking.

---

### 3. Distributed ID

```text
API 1 → unique ID
API 2 → unique ID
API 3 → unique ID
```

using something like:

```text
timestamp + server ID + sequence
```

Then:

```text
unique ID → Base62 → short code
```

**More scalable, but more complex.**

---

## And this is the system-design principle I want you to remember

If the interviewer asks:

> **"How would you generate unique IDs?"**

Don't immediately say:

> "Snowflake."

Start with:

> "For a simple single-region system, I'd use a database-generated ID. If we need multiple independent ID generators, I'd consider a distributed ID-generation approach such as a Snowflake-style ID."

That answer demonstrates **judgment**, not just knowledge of distributed systems.

And that's much more important for a **Senior Software Engineer** interview.