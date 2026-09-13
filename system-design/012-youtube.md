We will follow:

1. Requirements
2. Data and data-storage decisions
3. Database choice and why
4. APIs
5. Simplest architecture
6. Identify bottlenecks / external or slow components
7. Queues and asynchronous processing
8. Caching
9. Global distribution and replication
10. Scaling
11. Consistency and failure scenarios
12. Observability
13. Final architecture + trade-offs
14. Interview follow-ups

And this time, **the global SQL question is going to be a major part of the design**, not something we bolt on at the end.

---

# System Design: YouTube

**Interviewer:** Design YouTube.

**Me:**

> **Before I design it, I'd like to clarify a few requirements.**

### 1. Functional requirements

I'd like the system to support:

- Users uploading videos.
- Users watching videos.
- Searching for videos.
- Viewing video metadata such as title, description, channel and statistics.
- Likes and comments.
- Subscribing to channels.
- A home/recommendation feed.

For this design, I'll focus primarily on **video upload, video playback and globally distributed video access**. I'll keep things like live streaming and real-time collaborative features out of scope.

### 2. Non-functional requirements

Since we're designing YouTube, I'll assume:

- Users are distributed globally.
- There are billions of video views.
- Video files can be very large.
- Read traffic is vastly larger than write traffic.
- Playback latency should be low.
- The system should be highly available.
- Videos should be durable and not lost.
- We should continue serving popular content even if an individual server or region fails.
- Some data can be eventually consistent.

(I'd explicitly separate **video content**, **metadata**, and **user-generated data**, because they have very different storage and consistency requirements.)

---

# 3. What data do we need?

I'd start by identifying the major entities.

### Video

```text
Video
-----
id
channel_id
title
description
status
created_at
updated_at
```

The `status` could indicate things such as processing, published, or deleted.

### Channel

```text
Channel
-------
id
user_id
name
created_at
```

### User

```text
User
----
id
...
```

### Comments

```text
Comment
-------
id
video_id
user_id
text
created_at
```

### Likes

```text
VideoLike
---------
video_id
user_id
created_at
```

There is also the actual **video file**, but that's where I would make an important distinction.

I would **not store the video bytes in SQL**.

The metadata belongs in a database, while the large video objects belong in **object storage**.

---

# 4. Where does the actual video live?

I'd use object storage for the original and processed video files.

A video might go through several processing stages:

```text
Original video
      |
      v
Transcoding
      |
      +---- 360p
      +---- 720p
      +---- 1080p
      +---- 4K
```

The resulting files would be stored in object storage.

(I don't want the application servers or SQL database handling multi-GB binary data.)

So already we have two very different storage systems:

**SQL database**

→ users, channels, video metadata, comments, permissions, etc.

**Object storage**

→ actual video content.

---

# 5. Database choice

For the metadata, I'd start with a **relational database**.

There are relationships between:

- users and channels
- channels and videos
- users and comments
- users and likes

And we may need transactions and constraints for some operations.

But here's where our global requirement becomes interesting.

If I simply say:

> "We'll have one SQL database in the US and every user in Germany, India and Australia accesses it."

then we've created a latency problem.

So I need to think about **what exactly needs to be globally distributed**.

And this is where I want to distinguish three things:

### Video content

Can be distributed globally very aggressively.

### Read-heavy metadata

Can potentially be replicated to different regions.

### Authoritative transactional data

May need a primary/authoritative location depending on the consistency requirements.

This means **we don't necessarily need one database architecture for everything**.

---

# 6. APIs

I'd expose APIs such as:

```text
POST /videos
GET  /videos/{videoId}

GET  /videos/{videoId}/comments
POST /videos/{videoId}/comments

POST /videos/{videoId}/like
DELETE /videos/{videoId}/like

GET /search?q=...

GET /users/me/feed
```

For playback, I don't want:

```text
Client → API Server → Video bytes → Client
```

because that would make our API servers handle enormous amounts of bandwidth.

Instead, the API should provide the information needed to access the video from a CDN.

---

# 7. Simplest architecture

Before worrying about global distribution, I'd establish the basic architecture:

```text
                       Client
                          |
                         API
                       /     \
                      /       \
                 SQL DB     Object Storage
                                |
                                v
                             Video
```

The API handles:

- authentication
- metadata
- permissions
- comments
- likes
- search requests

Object storage handles the actual video files.

But this architecture is **not sufficient for global YouTube-scale traffic**.

Suppose a video becomes extremely popular.

Millions of people in Germany, India and the US might request the same video.

I don't want all those requests going directly to object storage.

So the next component is obvious.

---

# 8. CDN

I'd put a **CDN in front of the video content**.

```text
Germany ──┐
India ────┼──> CDN ──> Object Storage
USA ──────┘
```

A user requesting a popular video can receive it from a nearby CDN edge location.

The CDN can cache the video segments.

This is one of the most important ideas in this design:

> **Global low latency doesn't necessarily mean globally replicating the SQL database.**

For video playback, the data that matters most for latency is the **video itself**, so we distribute that data through the CDN.

---

# 9. But what about SQL?

Now we get to the question you specifically wanted to understand.

Suppose the video metadata is:

```text
Video ID: 123
Title: "My trip to Berlin"
Channel: Harshad
Views: 10,234,567
```

A user in Germany requests it.

A user in India requests the same metadata.

Do both requests go to the same SQL database?

**They could.**

But at YouTube scale, I'd probably introduce **read replicas in multiple regions**, assuming our SQL technology and replication model support the required workload.

Conceptually:

```text
                    Primary SQL
                   /     |      \
                  /      |       \
                 v       v        v
          EU Read     US Read    Asia Read
          Replica     Replica    Replica
```

Writes go to the authoritative primary.

Reads can often be served by a geographically closer replica.

This dramatically reduces the network distance for read-heavy operations.

But now we have a new problem:

### Replication lag.

Suppose a creator changes the title.

The write goes to the primary:

```text
Title = "My trip to Berlin"
        ↓
Title = "My trip to Berlin - 2026"
```

The European replica might receive that change almost immediately.

The Asian replica might receive it slightly later.

So for a short period:

```text
Europe → new title
Asia   → old title
```

That's **eventual consistency**.

For YouTube metadata, that may be perfectly acceptable.

---

# 10. Does everything need to go through the primary?

No.

This is an important distinction.

For something like:

> "Give me the title and description of video 123."

I'd be comfortable serving that from a local read replica.

But for something where correctness depends on the latest authoritative state, I may need a stronger consistency strategy.

For example, if we had:

> "Change the ownership of this channel."

I'd want the write to go to the authoritative database.

Similarly, after a user performs a write, we need to think about **read-after-write consistency**.

Suppose I upload a video in Germany.

The write reaches the primary.

Immediately afterwards I request:

```text
GET /videos/123
```

If that request gets routed to an Asian replica that hasn't received the change yet, I might get:

> Video not found.

That's undesirable.

So we'd need a strategy for operations requiring read-after-write consistency.

One simple approach is:

> After a write, route that user's relevant reads to the primary for a short period.

There are other approaches too, such as tracking replication positions/versions and only reading from a replica that has caught up.

---

And **this is exactly the part I want us to dig into deeply in this system design**.

We now have the fundamental question:

> **When should we use one primary + regional read replicas, and when do we actually need multiple writable database regions?**

We'll tackle that next rather than glossing over it, because that distinction is at the heart of globally distributed SQL systems.

### Continuing the YouTube design

**Interviewer:** *Do we really need multiple writable SQL databases around the world?*

**Me:**

> Not necessarily. I would first ask whether our workload actually requires writes to happen locally in every region.

For YouTube, the majority of traffic is **read traffic**. Video playback, video metadata, search, recommendations, comments and similar operations generate enormous numbers of reads compared with writes.

So I'd start with:

```text
                    Primary SQL
                   /     |      \
                  /      |       \
              EU Read   US Read   Asia Read
              Replica   Replica   Replica
```

Writes go to the primary, while geographically close replicas handle most reads.

This is considerably simpler than having multiple writable databases.

---

## But what if the primary is in the US?

**Interviewer:** *Wouldn't a creator in India have high latency when uploading or updating a video?*

**Me:**

> The write itself could have higher latency if it has to travel to the primary, but I wouldn't put the video upload through the SQL database anyway.

The large video upload goes directly to object storage.

For metadata writes, we're talking about relatively small requests.

If write latency becomes a significant problem, then we can consider a more distributed write architecture.

But I wouldn't introduce multi-primary writes merely because users are geographically distributed.

(I want to optimize based on the actual workload rather than assuming every piece of data must have a local writable copy.)

---

# 11. When would we actually need multiple writable regions?

Suppose YouTube has users in Europe, Asia and North America, and we discover that metadata writes are also enormous and users require very low write latency.

Now we could consider:

```text
       Europe              USA                Asia
         |                  |                   |
      SQL DB             SQL DB              SQL DB
         \                  |                  /
          \_________________|_________________/
                    Replication
```

Now all three regions can accept writes.

But we've introduced a **much harder problem**.

What happens if two regions modify the same piece of data at approximately the same time?

For example:

```text
Europe: title = "A"
Asia:   title = "B"
```

Which one wins?

We now need conflict resolution or a mechanism that coordinates the writes.

That's why **multi-primary/multi-writer databases are significantly more complicated** than primary + read replicas.

---

# 12. So what would I choose for YouTube?

I'd actually partition the problem by **type of data**.

For example:

| Data | Likely architecture |
|---|---|
| Video files | Object storage + CDN |
| Video metadata | Primary SQL + regional read replicas |
| Comments | SQL/sharded DB + replicas |
| Likes/views | Specialized high-write/event architecture |
| Recommendations | Distributed processing + caches |
| Search index | Globally distributed search infrastructure |

This is important because **"the database" isn't necessarily one database anymore.**

Different workloads can have different storage architectures.

---

# 13. But there's another scaling problem: one SQL primary

Let's say YouTube has:

- 1 billion daily users
- enormous read traffic
- millions of metadata writes per second

Even if reads are distributed to replicas, eventually the primary could become a bottleneck for writes.

At that point I'd consider **sharding**.

Instead of:

```text
             One giant SQL DB
            /      |       \
        Users    Videos   Comments
```

we could partition the data:

```text
Shard 1
Users/Videos A-H

Shard 2
Users/Videos I-P

Shard 3
Users/Videos Q-Z
```

But the exact partitioning strategy depends on the access patterns.

For example, we might shard based on `UserId`, `ChannelId`, or another suitable key.

---

# 14. Geographic partitioning vs replication

This distinction is **really important**.

Suppose we say:

> "Europe's database contains European users and Asia's database contains Asian users."

That's **partitioning**.

The data isn't duplicated. Different regions own different portions of the data.

Whereas:

```text
Primary
   |
   +---- Replica EU
   |
   +---- Replica Asia
```

means the **same data is duplicated** across regions.

That's **replication**.

And we can combine the two.

For example:

> Shard the global dataset into multiple logical partitions, and replicate each partition to multiple regions.

That's how systems can become extremely large without requiring every database node to contain the entire world's data.

---

# 15. Now let's return to YouTube's most important read path

A user in Germany watches a video.

The request might conceptually look like:

```text
User in Germany
       |
       v
Global traffic routing
       |
       v
European API
       |
       +----> EU cache / metadata replica
       |
       +----> CDN
                |
                v
          Video segments
```

The SQL database may not even participate in the actual video playback.

The metadata might come from a nearby replica or cache.

The video itself comes from a nearby CDN edge.

This is why YouTube can serve the same popular video to users across the world without sending every request to one central SQL database.

---

## 16. One subtle but important point

**Interviewer:** *If we have an EU read replica, does that mean the application automatically knows to use it?*

**Me:**

> No. We need some form of routing.

We could have regional API deployments configured with their preferred database replica.

For example:

```text
German user
    ↓
EU API
    ↓
EU SQL replica
```

while:

```text
Indian user
    ↓
Asia API
    ↓
Asia SQL replica
```

The global traffic-routing layer can direct users to an appropriate region.

But the system still needs to know whether a particular operation can safely use a replica.

A read-only operation such as:

> "Give me video 123"

can usually use the local replica.

A write operation such as:

> "Change the video title"

goes to the authoritative writer.

And an operation requiring the very latest data may need special handling.

---

# 17. What if the EU region goes down?

This is where **replication and availability** become useful.

Suppose:

```text
             Primary
                |
        +-------+-------+
        |               |
      EU DB           Asia DB
        X
```

If the EU region disappears, traffic can be routed to another healthy region.

The user's request may have slightly higher latency, but the service remains available.

For the video itself, the CDN may already have the content cached independently, so playback might continue even while an API region is unavailable.

This gives us **regional fault tolerance**.

---

# 18. But there's a dangerous assumption

We should **not** say:

> "If the primary dies, just promote any replica."

We first need to understand replication state.

Suppose the primary has processed:

```text
Video 1
Video 2
Video 3
Video 4
```

and the EU replica has only received:

```text
Video 1
Video 2
Video 3
```

If we immediately promote EU, we've potentially lost the latest write.

So failover involves questions such as:

- How far behind is the replica?
- Was the latest transaction replicated?
- Which node is allowed to become primary?
- How do other nodes discover the new primary?
- How do we prevent two primaries from accepting writes?

This is where **distributed consensus/failover mechanisms** start becoming relevant.

---

And this is the key mental model I want you to take from this design so far:

> **Global applications don't automatically mean "put a copy of the SQL database everywhere."**

Instead, we choose separately for each workload:

**Where is the authoritative data?**

**Where can copies safely exist?**

**Which copies can accept reads?**

**Which copies can accept writes?**

**How stale can a copy be?**

**What happens when the authoritative region fails?**

**Can the data be partitioned geographically or by key?**

For YouTube, a very reasonable starting architecture is:

**one authoritative SQL write path + regional read replicas + CDN/object storage for the massive video-read workload.**

Then, only where the scale or latency requirements demand it, we introduce **sharding or multi-region/multi-writer architecture**.

## Continuing: Sharding and Global Data Ownership

**Interviewer:** *Our single SQL primary is becoming a bottleneck. How would you scale the database?*

**Me:**

> I would consider sharding the database so that the workload is distributed across multiple database partitions.

The first question is **what should determine the shard**.

For YouTube, `UserId` or `ChannelId` would be strong candidates because many operations are naturally associated with a user or channel.

For example:

```text id="9d4y2m"
                 Video Data
                    |
        +-----------+-----------+
        |           |           |
      Shard A     Shard B     Shard C
      Users       Users       Users
       1-...       ...         ...
```

But I wouldn't literally use ranges like `1–1M`, `1M–2M`, etc. without considering the distribution of traffic.

A hash-based partitioning scheme can distribute users more evenly.

---

### 1. What exactly are we sharding?

I'd avoid assuming that the entire application has to use one shard key.

For example:

- User/channel data could be sharded by `UserId` or `ChannelId`.
- Video metadata could potentially be colocated with its channel.
- Comments could be partitioned by `VideoId`.
- View events could use a completely different high-throughput storage architecture.

The goal is to choose a partitioning strategy based on **access patterns**.

(I want related data that is frequently queried together to be colocated where practical.)

---

### 2. Each shard can still have replicas

Sharding and replication solve different problems.

We could have:

```text id="b1xq5w"
                    Global Dataset
                         |
             +-----------+-----------+
             |           |           |
           Shard A     Shard B     Shard C
             |           |           |
          Primary      Primary      Primary
          /    \       /    \       /    \
        EU     US     EU     US    EU     US
       replica replica replica replica replica replica
```

So:

**Sharding** distributes *different data* across databases.

**Replication** creates *copies of the same data*.

And we can use both simultaneously.

---

### 3. How does the API know which shard to query?

We need a routing layer.

For example, if we're sharding by `UserId`, the API can determine the shard from the user's ID.

Conceptually:

```text id="qz7t0u"
UserId
   |
   v
Shard routing
   |
   v
Shard B
```

The application doesn't need to ask every database:

> "Do you have this user?"

It calculates or looks up the appropriate shard.

For a hash-based scheme, something like:

`hash(UserId) → shard`

could determine the destination.

In a real system, I'd generally have a **shard map/configuration service** rather than hardcoding this logic everywhere.

---

### 4. What happens when we add more shards?

This is one of the major problems with sharding.

Suppose we have:

```text id="6c7m2k"
Shard A
Shard B
Shard C
```

and later need:

```text id="x8f2z1p"
Shard A
Shard B
Shard C
Shard D
Shard E
Shard F
```

If our hash function directly maps IDs to the number of shards, adding shards can cause a huge amount of data to move.

That's why systems often use techniques such as **consistent hashing** or a logical-partition approach.

For example, instead of saying:

> "There are exactly 6 physical shards."

we could have many logical partitions and map those partitions onto physical database nodes.

That gives us much more flexibility when scaling.

---

### 5. What about a massive YouTube channel?

Now we encounter another important problem.

Suppose a celebrity has a channel with billions of views and enormous amounts of comments.

If we shard by `ChannelId`, all that activity might end up concentrated on one shard.

That's a **hot shard**.

So the shard key isn't simply about distributing storage.

It also needs to distribute **traffic**.

For very hot data, we may need additional partitioning.

For example, comments could be partitioned using something like:

`ChannelId + VideoId`

rather than putting everything belonging to a channel onto one partition.

---

### 6. What about queries across shards?

Suppose an administrator asks:

> "Give me the 100 most-viewed videos across all of YouTube."

If videos are distributed across 1,000 shards, there isn't one database containing the answer.

We could query every shard:

```text id="wq4s8j"
Shard A ──┐
Shard B ──┤
Shard C ──┤
...       ├──> Merge results
Shard Z ──┘
```

But doing that synchronously for every user request would obviously be expensive.

So we'd generally maintain **precomputed aggregates, indexes, or specialized analytical systems** for global queries.

This is an important consequence of sharding:

> **Sharding makes individual partitions easier to scale, but cross-shard queries become harder.**

---

### 7. What about geographic ownership?

Now let's connect this back to your original question.

We could have a user whose data is primarily owned by a particular region.

For example:

```text id="q5m0k8"
German user
    |
    v
EU shard
    |
    +--> EU primary
    |
    +--> US replica
    |
    +--> Asia replica
```

And another user:

```text id="k0r8sp"
Indian user
    |
    v
Asia shard
    |
    +--> Asia primary
    |
    +--> EU replica
    |
    +--> US replica
```

Now we're combining:

**partitioning by data ownership**

with

**replication for geographic reads and availability.**

This can reduce write latency because a user's authoritative data can be owned by a geographically appropriate region.

But it introduces another question:

> What happens when users from different regions interact with the same data?

That's where things get considerably harder.

---

### 8. Do we actually need multi-region writes?

For YouTube, I'd still try to avoid them wherever possible.

Suppose a German creator owns a channel whose authoritative shard is in Europe.

An Indian user commenting on that creator's video doesn't necessarily require us to make the same SQL row writable from both Europe and India.

We could route the comment write to the appropriate owning shard.

The Indian user may experience slightly higher write latency, but the system remains much easier to reason about.

For workloads where **very low write latency in every region is mandatory**, we'd have to consider more sophisticated multi-region write architectures.

That's the point where we start discussing things like:

- multi-leader replication
- consensus
- quorum
- conflict resolution
- globally distributed SQL databases

But I wouldn't jump there for YouTube unless the requirements justify it.

---

### The important mental model

At this point, our SQL architecture isn't simply:

> **"One database with copies."**

It's closer to:

```text
                   Global SQL Data
                         |
          +--------------+--------------+
          |              |              |
       Shard A         Shard B        Shard C
          |              |              |
       Primary        Primary        Primary
       /    \         /    \         /    \
      EU    US       EU    Asia     US    Asia
    replica replica replica replica replica replica
```

Each **shard owns a different portion of the data**, while its replicas provide geographic read capacity and failover.

And importantly, **the CDN still handles the enormous video-content traffic**, so our SQL system isn't being asked to serve the actual video streams.

---

The next interesting problem is now **YouTube's views and likes**.

That's where we'll encounter a completely different scaling problem: **billions of writes**, hot videos, counters, eventual consistency, aggregation, and why we might *not* want to increment a SQL row synchronously for every view.

## YouTube: Views, Likes, and Massive Write Traffic

**Interviewer:** *YouTube gets billions of video views. How would you handle view counts and likes at this scale?*

**Me:**

> I wouldn't synchronously increment a row in our SQL database for every video view.

(That would turn a very hot video into a hot database row and create enormous write contention.)

Instead, I'd treat a view as an **event**.

A simplified flow would be:

```text
User watches video
       |
       v
API / playback service
       |
       v
Event stream / queue
       |
       +----> View aggregation
       |
       +----> Analytics
       |
       +----> Recommendation system
```

The user's video playback should **not wait for the database counter to be updated**.

---

### 1. Why not simply do this?

Suppose a video gets 100,000 views per second.

Our naive implementation would do:

```text
UPDATE Videos
SET ViewCount = ViewCount + 1
WHERE Id = 123;
```

100,000 times per second against the same row.

Even if SQL can technically process a large number of updates, we've created a **hot row**.

And there's no reason the displayed view count needs to be accurate to the exact millisecond.

So instead, we'd collect events and aggregate them.

---

### 2. Aggregating views

For example, workers could process batches:

```text
100,000 view events
        |
        v
Aggregation
        |
        v
+100,000 views
        |
        v
Persist aggregate
```

The displayed count might therefore lag by a few seconds.

For YouTube, that's generally acceptable.

This gives us a much more scalable architecture.

---

### 3. But where do the raw events go?

At YouTube scale, I'd expect a distributed event-streaming system capable of handling enormous throughput.

Something like Kafka or an equivalent distributed log would be a reasonable choice.

We could partition events by `VideoId`.

That gives us an interesting property:

```text
Video A events → Partition 1
Video B events → Partition 2
Video C events → Partition 3
```

Events for the same video can therefore be processed in an ordered stream.

But there's a problem.

### What if one video becomes extremely popular?

Suppose one video gets 10 million views per second.

If every event for that video goes to one partition, we've created a **hot partition**.

So even event partitioning needs careful thought.

We might partition using a combination of video ID and another value, allowing a very hot video's events to be spread across multiple partitions.

Then aggregation becomes a two-stage process:

```text
Raw events
    |
    v
Parallel aggregation
    |
    v
Partial counts
    |
    v
Final aggregation
```

---

### 4. Do we lose accuracy?

We need to distinguish **exact event processing** from **eventual aggregation**.

We want each legitimate view to contribute appropriately, but the displayed aggregate doesn't have to update synchronously.

We'd also need to think about duplicate events.

For example:

```text
Client sends view event
       |
       v
Server processes it
       |
       X
    response lost
       |
       v
Client retries
```

Now we might receive the same logical event twice.

So we'd want an event ID or another mechanism for deduplication where the business requirement requires it.

(Exactly-once processing across a large distributed system is expensive and complicated. Often we instead design processing to be **at-least-once** and make the aggregation idempotent.)

---

## 5. What about likes?

Likes are different.

A user can only have one active like per video.

So I care about the invariant:

> One user should not have two active likes on the same video.

That is a stronger consistency requirement than simply displaying a view count.

I could maintain something like:

```text
VideoLike
---------
video_id
user_id
created_at
```

with a uniqueness constraint on:

```text
(video_id, user_id)
```

Then a user's like operation can be made idempotent.

If the user clicks Like twice, we don't create two likes.

---

### 6. But likes can also become extremely hot

Imagine a celebrity releases a video and receives millions of likes within minutes.

We still don't necessarily want every like to synchronously update:

```text
Videos.LikeCount
```

Instead, we can separate:

**User's actual like state**

from

**displayed aggregate like count.**

The user's like relationship needs stronger consistency.

The aggregate count can be updated asynchronously.

So:

```text
UserLike
   |
   +----> authoritative state
   |
   +----> event
            |
            v
       aggregation
            |
            v
      LikeCount
```

This is a recurring distributed-systems pattern:

> **Keep the correctness-critical state separate from the derived aggregate.**

---

## 7. Where does SQL fit now?

We still absolutely can use SQL.

For example, SQL could remain the authoritative store for user/video relationships:

```text
User
Channel
Video
VideoLike
Comment
```

while the enormous stream of view events goes through a distributed event-processing system.

We might periodically write aggregated counters back into SQL or another serving store.

So we're not saying:

> "SQL can't handle YouTube."

We're saying:

> **Don't force one SQL row to synchronously handle every event generated by the entire world.**

---

## 8. What happens across regions?

Now our global architecture becomes more interesting.

A user in Europe generates a view.

A user in Asia generates another view.

We don't want both requests travelling across the world just to increment a SQL counter.

Instead:

```text
Europe                         Asia
   |                             |
Regional API                  Regional API
   |                             |
   v                             v
Regional event stream       Regional event stream
   |                             |
   +-------------+---------------+
                 |
                 v
          Global aggregation
                 |
                 v
           View counters
```

The events can be processed close to where they're generated.

The final aggregate can be globally combined.

This is much more scalable than globally synchronizing every individual view.

---

### 9. Now we have an important consistency decision

**Interviewer:** *Does every user immediately see the exact same view count?*

**Me:**

> No, I wouldn't require that.

I'd define the view count as **eventually consistent**.

If one region has processed 10,000 more views than another, users might temporarily see slightly different counts.

For a YouTube-style view counter, that's an acceptable trade-off for massive scalability and low latency.

For something like a bank balance, I would make a very different decision.

---

## 10. This gives us an important classification

We're now getting to the heart of distributed-system design.

| Data | Consistency requirement |
|---|---|
| Video file | Strong durability, CDN replication |
| Video title | Usually eventual for global reads |
| View count | Eventual |
| Like count | Eventual aggregate |
| User's like relationship | Stronger consistency |
| Channel ownership | Strong consistency |
| Payment | Strong consistency |
| Recommendations | Eventual |

So **the entire YouTube system does not have one consistency model**.

Different pieces of data have different requirements.

And that's one of the most important things to recognize when designing globally distributed systems.

---

### Where we've reached

Our architecture now has several distinct layers:

```text
                     GLOBAL USERS
                          |
                   Global Routing
                          |
              +-----------+-----------+
              |                       |
         Regional APIs            Regional APIs
              |                       |
        SQL / Cache              SQL / Cache
              |                       |
       +------+-------+       +-------+------+
       |              |       |              |
    Metadata        Events  Metadata       Events
                       \       /
                        \     /
                     Event Stream
                          |
                    Aggregation
                          |
                    Derived data


Video content
      |
      v
Object Storage
      |
      v
Global CDN
```

And now we've covered the major reason **not everything should be replicated through SQL synchronously**.

Next, I'd take us into **global SQL consistency and failover**: what happens when the primary region dies, how replicas are promoted, what replication guarantees mean, and then the harder question of **whether a globally distributed SQL database such as CockroachDB/Spanner-style architecture changes the design**.

---

## YouTube — Global Database Failure & Consistency

### Interviewer:
> You have a primary SQL database in the US and read replicas in Europe and Asia. What happens if the US region goes down?

### Candidate:

I would first clarify what kind of failure we're dealing with.

If the **US application servers** fail but the database remains healthy, that's relatively straightforward: traffic can be routed to another API region.

But if the **US database primary itself fails**, we need database failover.

Our architecture currently looks roughly like:

```text
                    Global Users
                         |
                  Global Routing
                         |
          +--------------+--------------+
          |              |              |
        EU API         US API        Asia API
          |              |              |
          v              v              v
       EU SQL          US SQL         Asia SQL
      Replica          Primary        Replica
```

Normally:

```text
                    US SQL
                    Primary
                   /       \
                  /         \
                 v           v
             EU Replica   Asia Replica
```

If US goes down, we need to promote one of the replicas.

For example:

```text
                    EU SQL
                    Primary
                   /       \
                  /         \
                 v           v
            US Replica   Asia Replica
```

Then global routing sends writes to the new EU primary.

---

## But we can't simply say "pick the closest replica"

This is an important part of the design.

Suppose the US primary had accepted:

```text
Video 123
title = "My New Video"
```

but replication to Europe was slightly behind.

The EU replica might still contain:

```text
Video 123
title = "Old Title"
```

If we immediately promote Europe, the acknowledged update might disappear.

So failover depends partly on **how much replication lag we tolerate**.

For YouTube, this gives us a trade-off.

### Asynchronous replication

```text
US Primary
    |
    | replicate later
    v
EU Replica
```

The primary doesn't wait for Europe before acknowledging the write.

That's good for write latency.

But if US suddenly disappears:

```text
US: Video title = "New Title"   ✓ acknowledged

EU: Video title = "Old Title"   ← replication hadn't arrived
```

The new title could potentially be lost during failover.

For many YouTube operations, that may be an acceptable trade-off.

A user changing a video's description isn't the same as a bank transaction.

---

## Stronger replication

We could instead require some writes to be replicated to another region before acknowledging them.

Conceptually:

```text
          Write
            |
            v
        US Primary
         /       \
        v         v
      EU DB     Asia DB
        |         |
        +----ACK--+
             |
             v
          Client
```

Now the acknowledged write is much safer if US disappears.

But there's a cost:

**The write now depends on cross-region communication.**

That increases latency.

So for YouTube, I wouldn't make every piece of data globally synchronous.

I'd classify the data according to how important the consistency is.

| Data | Consistency |
|---|---|
| Video file | Strong durability |
| Video title/description | Usually eventual |
| View count | Eventual |
| Like count | Eventual |
| User's like relationship | Stronger |
| Channel ownership | Strong |
| Comments | Moderate/strong depending on operation |
| Recommendations | Eventual |

That's a very important design principle:

> **Don't choose one consistency model for the entire system. Choose it per type of data and operation.**

---

# What actually happens during failover?

Let's make the YouTube scenario concrete.

Initially:

```text
US
└── SQL Primary
      |
      +---- EU Replica
      |
      +---- Asia Replica
```

US suddenly becomes unavailable.

### 1. Detect failure

Our database infrastructure detects that the primary isn't healthy.

We shouldn't promote a replica just because one API server couldn't connect.

We need reasonably strong evidence that the primary is actually unavailable.

---

### 2. Select a replica

Suppose:

```text
EU Replica
Replication lag: 20 ms

Asia Replica
Replication lag: 2 seconds
```

EU is the better candidate.

But we also need to know whether EU has reached a safe point in the transaction log.

---

### 3. Promote EU

EU becomes the new authoritative writer.

```text
EU SQL
Primary
```

The other regions now replicate from EU.

```text
                  EU Primary
                 /          \
                v            v
          US Replica      Asia Replica
```

---

### 4. Redirect writes

Our YouTube API services need to discover the new primary.

Previously:

```text
Write → US SQL
```

Now:

```text
Write → EU SQL
```

We shouldn't hardcode:

```text
db-server-us
```

into every API server.

Instead, the database layer/service discovery tells the application which endpoint is currently authoritative.

---

### 5. Prevent split-brain

This is extremely important.

Imagine US comes back alive.

If US thinks:

> "I'm still the primary."

while EU also thinks:

> "I'm the primary."

we could have:

```text
US Primary  <----->  EU Primary
```

and both accept writes.

Now imagine:

```text
US:
Video title = A

EU:
Video title = B
```

We've created conflicting histories.

So the old primary must be **fenced off** from accepting writes before the new primary is allowed to take over.

(At this scale, this usually involves some form of lease/consensus/fencing mechanism in the database or orchestration layer. The important interview point is that failover must guarantee **one authoritative writer**, not simply "promote a replica.")

---

# What about users watching videos during the failure?

This is where YouTube's architecture becomes interesting.

Remember that we're **not streaming video through the SQL database or even primarily through our API servers**.

We have:

```text
User
  |
  v
CDN
  |
  v
Video Object Storage
```

So suppose the US database goes down.

A user might still be able to watch a popular video because the video segments are already cached at a nearby CDN.

```text
User
 |
 v
CDN
 |
 +-- video already cached
 |
 v
Playback continues
```

The user might temporarily have trouble with:

- uploading a video
- changing metadata
- liking a video
- posting a comment
- loading personalized recommendations

But **already-cached video content can continue to be served**.

This is a major benefit of separating:

**control plane**

from

**content delivery**.

For YouTube:

```text
Control / metadata

Client
  |
  v
API
  |
  v
SQL
```

versus:

```text
Video delivery

Client
  |
  v
CDN
  |
  v
Object Storage
```

A failure in the control plane doesn't necessarily stop the entire video delivery system.

---

# What about read-after-write?

Here's another YouTube-specific problem.

Suppose I'm a creator in Germany.

I upload:

> "My Trip to Berlin"

The write goes to the current primary.

Immediately afterward I request:

```http
GET /videos/123
```

But my request might hit the European replica.

If replication hasn't caught up yet, I could see:

```text
404 Not Found
```

or old metadata.

That's obviously a bad experience immediately after an operation.

So for operations where the user expects to immediately see their own change, we can use a stronger read strategy.

For example:

```text
POST /videos
       |
       v
   Primary
       |
       v
  version = 84521
```

The client/request context can carry that version.

Then:

```text
GET /videos/123

EU replica
   |
   +-- has version 84521? → yes → serve
   |
   +-- not caught up? → primary / another replica
```

We don't need every read to go to the primary.

We only need to ensure that **reads requiring fresh state don't accidentally use an older replica**.

That's much more scalable.

---

# RPO and RTO

At this point, two useful reliability concepts naturally appear.

### RPO — Recovery Point Objective

How much acknowledged data can we potentially lose?

For example:

> RPO = 1 second

means that in the worst case, we may lose around one second of recently acknowledged data.

For YouTube:

Losing a recently updated view count isn't particularly concerning.

Losing a creator's entire uploaded video is obviously unacceptable.

So different data paths can have very different durability requirements.

---

### RTO — Recovery Time Objective

How long can the system take to recover?

For example:

> RTO = 30 seconds

means we want the service back within roughly 30 seconds after a major failure.

Again, YouTube doesn't necessarily need every subsystem to have exactly the same RTO.

Video playback may have extremely aggressive availability requirements, while some administrative metadata operations can tolerate a little more degradation.

---

# Now let's update our YouTube architecture

We're getting somewhere much more realistic:

```text
                         GLOBAL USERS
                              |
                     Global Traffic Routing
                              |
              +---------------+---------------+
              |               |               |
           EU API          US API          Asia API
              |               |               |
              +---------------+---------------+
                              |
                    YouTube Control Plane
                              |
                    +---------+---------+
                    |                   |
                 SQL DB             Event System
                    |                   |
             Metadata/state        Views/Likes/etc.
                    |
        +-----------+-----------+
        |           |           |
       EU          US          Asia
     replica      primary     replica


             VIDEO CONTENT PATH
             
User
 |
 v
CDN
 |
 +---- cached video segments
 |
 v
Object Storage
 |
 +---- original
 +---- 360p
 +---- 720p
 +---- 1080p
 +---- 4K
```

And during a US database failure:

```text
                         Global Users
                              |
                     Global Traffic Routing
                              |
                    +---------+---------+
                    |                   |
                 EU API             Asia API
                    |                   |
                    +---------+---------+
                              |
                         EU SQL
                         Primary
                        /       \
                       v         v
                 US Replica   Asia Replica


Video playback can still go:

User → CDN → cached/object-stored video
```

That's still **100% YouTube**. We're not leaving the problem.

---

# One more YouTube-specific question

Now imagine the interviewer pushes us:

> "You're saying there is one primary SQL database. YouTube has users all over the world. Won't all those writes eventually overwhelm that primary?"

Yes.

And this brings us to the next scaling problem **within YouTube**:

```text
             Millions of writes
                    |
                    v
              One Primary
                    |
              bottleneck
```

We already introduced **sharding** earlier as one possible answer.

But now we need to integrate it properly into our YouTube architecture rather than discussing sharding abstractly.

For example:

```text
                    YouTube SQL
                         |
              +----------+----------+
              |          |          |
           Shard 1    Shard 2    Shard 3
              |          |          |
           users A-F   G-M       N-Z
```

And then each shard can itself have replicas:

```text
                 Shard 1
                    |
              +-----+-----+
              |           |
           Primary      Replica
             US            EU
```

So **sharding + replication + geographic distribution** can all coexist.

The important question for YouTube becomes:

> **What exactly should we shard by?**

That's the next thing I'd tackle, because the answer isn't simply "shard by UserId." Different YouTube workloads—videos, comments, channels, views, likes, search—have different access patterns.

And that takes us deeper into the **YouTube design itself**, not into a generic database lesson.

## YouTube — Sharding the Data

### Interviewer:
> You have one SQL primary. At YouTube scale, how would you prevent it from becoming a bottleneck?

### Candidate:

I would introduce **sharding**, but I wouldn't shard the entire YouTube database using one universal key.

Different parts of YouTube have different access patterns, so I'd first separate the major workloads.

We have:

- Users / channels
- Videos
- Comments
- Likes
- View events
- Video metadata
- Recommendations
- Search

The shard key should be chosen based on **how that data is accessed**, not simply because one column is convenient.

---

### 1. Users and channels

For user-related data, `UserId` is a reasonable starting point.

```text
             User data
                 |
       +---------+---------+
       |         |         |
    Shard 1   Shard 2   Shard 3
    A-F       G-M       N-Z
```

But in practice I'd probably use **hash-based partitioning** rather than literal alphabetical ranges.

For example:

```text
hash(UserId) → logical partition → physical shard
```

This distributes users more evenly.

A user generally accesses their own:

- profile
- subscriptions
- settings
- channel information

So keeping user-related data together can be useful.

---

# 2. Videos are more interesting

A video belongs to a channel:

```text
Video
  |
  +-- ChannelId
```

We could shard videos using `ChannelId`.

That means:

```text
Channel A
   |
   +-- Video 1
   +-- Video 2
   +-- Video 3
```

would generally be colocated.

That's useful for queries such as:

```http
GET /channels/{channelId}/videos
```

because we know which shard to query.

But there's a problem.

Imagine Taylor Swift uploads thousands of videos.

Or, more importantly, imagine an enormous channel generating huge amounts of traffic.

Now one shard could become a **hot shard**.

So the shard key isn't just about distributing storage.

It must distribute **traffic and workload**.

---

# 3. Comments

Comments have a different access pattern.

The common query is:

```http
GET /videos/{videoId}/comments
```

So `VideoId` is a natural partitioning key.

```text
Video 123
   |
   +-- Comment 1
   +-- Comment 2
   +-- Comment 3
   +-- ...
```

But again, YouTube has extremely popular videos.

Suppose one video receives:

> 10 million comments.

Putting all those comments into one partition could create a hotspot.

So we might partition comments using something like:

```text
hash(VideoId, CommentBucket)
```

where the bucket can divide one extremely popular video's comments across multiple partitions.

For example:

```text
Video 123
   |
   +-- Bucket 0
   +-- Bucket 1
   +-- Bucket 2
   +-- Bucket 3
```

(The important principle is that a shard key needs to prevent **hotspots**, not merely distribute rows.)

---

# 4. Likes

Likes have a particularly useful uniqueness requirement.

We want:

```text
one user
+
one video
=
one active like
```

So we'd have something conceptually like:

```text
VideoLike
----------------
VideoId
UserId
CreatedAt
```

with a unique constraint on:

```text
(VideoId, UserId)
```

This gives us a reliable way to prevent:

```text
User 42 likes Video 123
User 42 likes Video 123 again
```

from creating two active likes.

But we **don't** want every like to synchronously update:

```text
Videos.LikeCount
```

because a hugely popular video could become a hot row.

Instead:

```text
User
 |
 v
Like API
 |
 v
Authoritative Like State
 |
 v
Event
 |
 v
Aggregation
 |
 v
LikeCount
```

So:

```text
UserLike
```

is correctness-sensitive.

Whereas:

```text
LikeCount
```

is a derived value and can be eventually consistent.

---

# 5. Views

Views are even more different.

We don't want:

```sql
UPDATE Videos
SET ViewCount = ViewCount + 1
WHERE VideoId = 123;
```

for every view.

Imagine a popular video receiving hundreds of thousands of views per second.

That would turn one database row into a massive contention point.

Instead:

```text
Users
  |
  v
Playback
  |
  v
View Events
  |
  v
Event Stream
  |
  +--------+--------+
  |        |        |
  v        v        v
Aggregation / Analytics / Recommendations
  |
  v
Aggregated View Count
```

This is one of the biggest architectural differences between:

**transactional state**

and

**high-volume events**.

A view doesn't need to synchronously modify the video's canonical row.

---

# So our YouTube data isn't all in one giant sharded SQL database

This is the important architectural realization.

Instead, we might have something more like:

```text
                         YouTube
                            |
        +-------------------+-------------------+
        |                   |                   |
        v                   v                   v
   User/Channel         Video Metadata       Comments
      SQL                  SQL                  SQL
    sharded               sharded              sharded
        |                   |                   |
        +-------------------+-------------------+
                            |
                       Event System
                            |
                +-----------+-----------+
                |           |           |
              Views       Likes    Recommendations
```

And then separately:

```text
                    Video Files
                        |
                 Object Storage
                        |
                       CDN
```

And:

```text
                    Search
                      |
               Search Index
```

(I'm intentionally separating these because **YouTube doesn't have one database that does everything**. Different workloads need different storage/access patterns.)

---

# One more important problem: cross-shard queries

Suppose the homepage needs:

> "Give me 20 recommended videos for this user."

We don't want:

```text
Query Shard 1
Query Shard 2
Query Shard 3
...
Query Shard 100
```

and then combine everything synchronously.

That would be terrible for latency.

Instead, recommendations should already be prepared using the event/analytics pipeline.

Something like:

```text
Views
Likes
Watch history
Subscriptions
Search history
       |
       v
Recommendation processing
       |
       v
User 123 → [Video 91, Video 72, Video 381, ...]
```

Then:

```http
GET /users/me/feed
```

can retrieve a relatively small precomputed result.

This is another recurring YouTube principle:

> **Don't make the online request perform expensive global computation if you can precompute the result.**

---

# Our architecture is now evolving

We started with:

```text
Client
  |
  v
API
  |
  v
SQL
```

Then:

```text
Client
  |
  v
API
  |
  +---- SQL
  |
  +---- Object Storage
```

Then:

```text
Client
  |
  +---- API → SQL
  |
  +---- CDN → Object Storage
```

Then:

```text
                         Client
                        /      \
                       /        \
                      v          v
                    API         CDN
                     |            |
                     v            v
                  SQL          Video Files
                   |
             +-----+-----+
             |     |     |
            EU     US   Asia
          replicas primary replicas
```

Then at larger scale:

```text
                         YouTube
                            |
        +-------------------+-------------------+
        |                   |                   |
        v                   v                   v
   User/Channel         Video Metadata       Comments
      Shards               Shards             Shards
        |                   |                   |
        +-------------------+-------------------+
                            |
                       Event Stream
                            |
              +-------------+-------------+
              |             |             |
            Views         Likes      Recommendations
```

And the **video bytes themselves remain outside SQL**:

```text
Upload
   |
   v
Object Storage
   |
   v
Transcoding
   |
   +-- 360p
   +-- 720p
   +-- 1080p
   +-- 4K
   |
   v
CDN
   |
   v
Global viewers
```

At this point we have a fairly serious YouTube architecture.

The next natural pressure point is **video upload itself**: uploading a multi-GB file, resumable uploads, transcoding, processing failures, multiple resolutions, and how we make sure the API servers aren't handling the actual video bytes. That's still completely within our YouTube design.

## YouTube — Video Upload Pipeline

### Interviewer:
> Let's talk about uploading a video. A user might upload a 5 GB or 20 GB video. How would you design that?

### Candidate:

I wouldn't send the entire video through our API servers.

The API should handle **metadata and coordination**, while the actual video bytes should go directly to object storage.

So instead of:

```text
Client → API Server → Object Storage
```

I'd use:

```text
Client → API
          |
          | create upload
          v
      Upload URL
          |
          v
Client → Object Storage
```

The API never has to hold the 20 GB file in memory or proxy all those bytes.

---

## 1. Start the upload

The client calls:

```http
POST /videos
```

with metadata such as:

```json
{
  "title": "My Berlin Trip",
  "description": "..."
}
```

The API creates a video record:

```text
Video
------------------------
Id
ChannelId
Title
Description
Status = UPLOADING
CreatedAt
```

Then it generates an upload authorization mechanism, typically a **pre-signed URL** or multipart/resumable-upload session for object storage.

The API returns something like:

```json
{
  "videoId": "123",
  "uploadUrl": "..."
}
```

---

## 2. Upload directly to object storage

The client now uploads the actual file directly:

```text
Client
   |
   | 5 GB video
   v
Object Storage
```

The API isn't involved in transferring those 5 GB.

That's important because API servers are designed for request processing, not massive file transfer.

---

# 3. What if the upload fails at 80%?

We don't want the user to restart 5 GB from zero.

So I'd support **resumable/multipart uploads**.

Conceptually:

```text
5 GB file

+---------+---------+---------+---------+
| Part 1  | Part 2  | Part 3  | Part 4  |
+---------+---------+---------+---------+
    ✓          ✓          ✗          ✓
```

The client can retry only Part 3.

This is particularly important for users on unstable or slow connections.

---

# 4. How do we know the upload is complete?

The object-storage upload mechanism can provide completion information.

The client can then tell our API:

```http
POST /videos/123/upload-complete
```

Or object storage can emit an event when the object is successfully created.

I prefer an **event-driven approach** at large scale:

```text
Client
  |
  v
Object Storage
  |
  | ObjectCreated event
  v
Upload Processing
```

The processing system verifies that the expected object exists and then starts the next stage.

---

# 5. Transcoding

We can't serve the original 5 GB file directly to every viewer.

A user might have:

- a slow mobile connection
- a 1080p display
- a 4K display
- limited bandwidth

So we need multiple representations.

```text
Original Video
      |
      v
 Transcoding
      |
      +---- 360p
      +---- 480p
      +---- 720p
      +---- 1080p
      +---- 1440p
      +---- 4K
```

And increasingly, we're not just producing one giant file per resolution.

We split the video into **small segments**.

For example:

```text
1080p

segment-001
segment-002
segment-003
segment-004
...
```

This allows the player to request only the portions it needs.

---

# 6. Why segments matter

Suppose you're watching:

> 45:00 of a 60-minute video.

We don't want:

```text
Download entire 60-minute video
```

Instead:

```text
Player
  |
  +-- segment 001
  +-- segment 002
  +-- segment 003
  ...
  +-- segment 045
```

The player can also switch quality dynamically.

For example:

```text
Good network

1080p → 1080p → 1080p
```

Then the connection gets worse:

```text
1080p → 720p → 720p
```

Then improves:

```text
720p → 1080p
```

This is the basic idea behind **adaptive bitrate streaming**.

---

# 7. Where does the CDN fit?

After transcoding:

```text
Object Storage
      |
      v
CDN
```

The CDN caches frequently requested video segments.

So a viewer in Germany doesn't need to retrieve every segment from some distant origin.

Instead:

```text
                    Object Storage
                         |
                         v
                        CDN
                   /     |      \
                  /      |       \
                 v       v        v
               EU      US       Asia
             viewers  viewers  viewers
```

Popular video segments can be served from edge locations close to viewers.

This is where we get the huge global-scale benefit from the CDN.

---

# 8. Processing shouldn't block the upload request

We shouldn't do:

```text
POST /videos

Upload
  ↓
Transcode
  ↓
Generate thumbnails
  ↓
Analyze video
  ↓
Return response
```

That could take minutes or longer.

Instead:

```text
POST /videos
     |
     v
Create metadata
     |
     v
Return upload information
```

Then asynchronously:

```text
Object Storage
      |
      v
Upload event
      |
      v
Processing Queue
      |
      +---- Transcoding
      |
      +---- Thumbnail generation
      |
      +---- Metadata extraction
      |
      +---- Content processing
```

Each processing step can scale independently.

---

# 9. What if transcoding fails?

This is where the queue design we discussed earlier becomes useful.

For example:

```text
Upload Queue
     |
     v
Transcoding Worker
     |
     X failure
     |
     v
Retry
     |
     X failure
     |
     v
Retry
     |
     X failure
     |
     v
Dead Letter Queue
```

We'd use:

- retries
- exponential backoff
- maximum retry count
- dead-letter queue

The video record can reflect its state:

```text
UPLOADING
     ↓
PROCESSING
     ↓
READY
```

or:

```text
PROCESSING
     ↓
FAILED
```

The client can query:

```http
GET /videos/123
```

and receive:

```json
{
  "id": "123",
  "title": "My Berlin Trip",
  "status": "PROCESSING"
}
```

Eventually:

```json
{
  "id": "123",
  "title": "My Berlin Trip",
  "status": "READY"
}
```

---

# 10. We should also think about duplicate processing

Imagine the object-storage event is delivered twice:

```text
ObjectCreated(video-123)
ObjectCreated(video-123)
```

We don't want to transcode the same 5 GB video twice unnecessarily.

So processing should be **idempotent**.

For example, the processing job can have:

```text
VideoId = 123
ProcessingVersion = 1
```

and workers can check whether that processing version has already completed.

(Exactly-once delivery is difficult in distributed systems; it's generally safer to assume events may be delivered more than once and make consumers idempotent.)

---

# 11. What about thumbnails?

We can generate them asynchronously from the uploaded video.

```text
Original
   |
   +---- Transcoding
   |
   +---- Thumbnail generation
   |
   +---- Other processing
```

The thumbnails are also stored in object storage and served through the CDN.

So when search results display:

```text
[thumbnail] My Berlin Trip
```

the thumbnail itself doesn't need to come from the SQL database.

SQL contains something like:

```text
ThumbnailUrl / ThumbnailKey
```

while the actual image lives in object storage/CDN.

---

# Our upload architecture now looks like this

```text
                         Creator
                            |
                            | POST /videos
                            v
                      YouTube API
                            |
                    Create video metadata
                            |
                            v
                    SQL: UPLOADING
                            |
                            v
                     Upload URL
                            |
                            |
Creator -------------------+
   |
   | actual video bytes
   v
Object Storage
   |
   | upload event
   v
Processing Queue
   |
   +------------------+
   |                  |
   v                  v
Transcoding       Thumbnail
 Workers           Workers
   |                  |
   +--------+---------+
            |
            v
      Processed Video
            |
            v
      Object Storage
            |
            v
           CDN
            |
            v
         Viewers
```

And importantly, **SQL remains relatively small**:

```text
SQL
 ├── VideoId
 ├── ChannelId
 ├── Title
 ├── Description
 ├── Status
 ├── timestamps
 └── references/keys to stored assets
```

It does **not** contain the 5 GB video.

---

### Interviewer:
> What happens if millions of creators upload videos at the same time?

Now we have another YouTube-specific scaling problem.

We don't want:

```text
1 million uploads
       ↓
1 million transcoding jobs
       ↓
100 workers
       ↓
everything overloaded
```

The queue gives us **backpressure**.

Uploads can complete and wait in the queue while transcoding workers process jobs at the rate the infrastructure can handle.

And because transcoding is CPU/GPU intensive, we can scale those workers independently from the API servers.

That gives us another important separation:

```text
API capacity       ≠       Transcoding capacity
```

A spike in uploads shouldn't necessarily take down:

- video metadata APIs
- search
- comments
- playback authorization

because the expensive processing work is isolated behind the queue.

---

So our YouTube design now covers **upload → object storage → asynchronous processing → transcoding → CDN playback**, including failures and scaling.

The next pressure point inside the same design is **playback itself**: how the player gets the right video segments, how CDN caching works for billions of views, what happens when a video suddenly goes viral, and how we prevent the origin/object storage from being overwhelmed.

## YouTube — Video Playback at Global Scale

### Interviewer:
> Now let's focus on playback. A user clicks a video. Walk me through what happens.

### Candidate:

I'd separate **playback authorization/metadata** from **video delivery**.

The API doesn't need to send the video itself.

A simplified flow is:

```text
User
 |
 | GET /videos/123
 v
YouTube API
 |
 +---- SQL → metadata
 |
 v
Playback information
 |
 v
CDN
 |
 v
Video segments
```

The API might return information telling the player where and how to retrieve the video streams.

---

## 1. The player needs multiple representations

For our video:

```text
Video 123
   |
   +-- 360p
   +-- 720p
   +-- 1080p
   +-- 1440p
   +-- 4K
```

Each representation is divided into segments:

```text
1080p
 ├── segment 001
 ├── segment 002
 ├── segment 003
 ├── ...
```

The player doesn't necessarily download the entire video.

It requests segments as needed.

---

# 2. Adaptive bitrate

Suppose I'm watching on a good connection.

The player might initially select:

```text
1080p
```

If bandwidth drops:

```text
1080p
  ↓
720p
```

If it becomes even worse:

```text
720p
  ↓
480p
```

The goal is to avoid:

> **buffering**

rather than blindly maximizing resolution.

So the player continuously considers things such as:

- available bandwidth
- current buffer
- device capabilities
- video resolution
- segment download time

and chooses an appropriate representation.

---

# 3. Why the CDN is critical

Imagine a video becomes extremely popular.

```text
        Video 123
           |
     10 million viewers
           |
     +-----+-----+
     |     |     |
    EU    US    Asia
```

If every viewer requested the video directly from object storage:

```text
10 million viewers
       |
       v
Object Storage
```

the origin would receive an enormous amount of traffic.

Instead:

```text
                         Object Storage
                              |
                              v
                             CDN
                       /       |       \
                      /        |        \
                    EU        US       Asia
                  viewers   viewers   viewers
```

The CDN caches the video segments.

Once a segment is cached at an edge location, subsequent viewers can receive it without going all the way back to the origin.

---

# 4. Viral video problem

Now let's make this more extreme.

A celebrity uploads a video.

Within minutes:

> 50 million people want to watch it.

This creates a **cache-miss storm** if the CDN doesn't already have the content.

For example:

```text
50 million requests
        |
        v
CDN
        |
        | cache miss
        v
Object Storage
```

We don't want 50 million origin requests.

Instead, the CDN should use mechanisms such as **request collapsing/coalescing**.

Conceptually:

```text
Viewer 1 ─┐
Viewer 2 ─┤
Viewer 3 ─┤
Viewer 4 ─┤
           ├── CDN ── one origin request
Viewer 5 ─┤
Viewer 6 ─┤
Viewer 7 ─┘
```

The first request fetches the segment from the origin.

Other requests can wait for that same fetch.

Then:

```text
                    Origin
                       |
                       v
                  CDN segment
                 / / / | \ \ \
                v v v  v  v v v
              viewers...
```

So a viral video doesn't necessarily translate into a proportional increase in origin traffic.

---

# 5. Cache hierarchy

At enormous scale, we can have multiple layers.

Conceptually:

```text
Viewer
   |
   v
Edge CDN
   |
   | miss
   v
Regional CDN/cache
   |
   | miss
   v
Object Storage
```

For a popular video:

```text
Viewer
   |
   v
Edge Cache
   |
   +-- HIT → return segment
```

For a less popular video:

```text
Viewer
   |
   v
Edge Cache
   |
   +-- MISS
        |
        v
   Regional Cache
        |
        +-- HIT → return
```

Only when necessary do we reach object storage.

(This hierarchy is an architectural option; the exact CDN topology would depend on the infrastructure chosen.)

---

# 6. What happens when a video is not popular?

We shouldn't permanently cache every video everywhere.

That would be extremely expensive.

Suppose YouTube has billions of videos, but only a small fraction are receiving heavy traffic.

So caching naturally follows popularity.

```text
Popular video
→ many CDN cache hits

Rare video
→ fewer cached copies
→ eventually fetched from origin
```

This is one of the reasons CDNs work so well for YouTube-like workloads.

---

# 7. Cache expiration and video immutability

There's a very useful property here:

**The actual video content generally doesn't change after publication.**

If:

```text
video123/1080p/segment001
```

represents a particular immutable segment, we don't need to constantly invalidate it.

We can give it a long cache lifetime.

If the creator changes the title:

```text
SQL metadata changes
```

but:

```text
video segment
```

doesn't change.

So metadata and video content have very different caching characteristics.

That's another reason we separated them.

---

# 8. What if the creator deletes the video?

Now we have an interesting consistency problem.

Suppose:

```text
SQL
Video 123 = DELETED
```

but the CDN still has:

```text
video123/1080p/segment001
```

cached.

We don't necessarily want users to continue accessing a deleted video.

So playback authorization/metadata can check the video's current state before allowing access.

The CDN can also have appropriate invalidation/expiry mechanisms.

The key distinction is:

> **Caching immutable content is easy; deciding whether the user is still allowed to access that content is a control-plane concern.**

---

# 9. What if the API is temporarily unavailable?

This is another benefit of separating metadata/control from content delivery.

Suppose:

```text
YouTube API
     X
```

but:

```text
CDN
 |
 +-- cached video segments
```

are healthy.

Already-established playback sessions may continue consuming cached segments.

But a brand-new playback request may need the API to obtain authorization or playback information.

So we don't claim:

> "The CDN makes YouTube completely independent of the API."

Rather:

> **The CDN reduces the dependency of high-volume video delivery on the API and database.**

That's a much more accurate architecture statement.

---

# 10. Regional failure

Suppose the European CDN infrastructure has a problem.

We can route the viewer to another healthy CDN location.

Conceptually:

```text
Germany
   |
   v
EU CDN
   X
   |
   v
Alternative CDN location
   |
   v
Video
```

The exact routing mechanism could use CDN-level traffic management, DNS/anycast/global routing, etc.

The important design principle is:

> **Don't make one physical region a single point of failure for global playback.**

---

# 11. Playback versus view counting

There's another important separation.

A user watches:

```text
Video 123
```

The video segments come through:

```text
CDN
```

But that doesn't mean every segment request should synchronously update:

```text
ViewCount
```

Instead, the playback system can generate an appropriate **watch/view event** after satisfying the platform's rules for what constitutes a view.

Then:

```text
Playback
   |
   v
View Event
   |
   v
Event Stream
   |
   v
Aggregation
   |
   +---- View count
   +---- Analytics
   +---- Recommendations
   +---- Creator statistics
```

This keeps the extremely high-volume playback path independent from the transactional database.

---

# Our YouTube playback architecture

We're now at:

```text
                         GLOBAL VIEWERS
                              |
                              v
                       Global Routing
                              |
                         YouTube API
                              |
                    +---------+---------+
                    |                   |
                    v                   v
                  SQL            Playback information
                                        |
                                        v
                                       CDN
                                        |
                         +--------------+--------------+
                         |              |              |
                      EU Edge        US Edge        Asia Edge
                         |              |              |
                         +--------------+--------------+
                                        |
                                  Cache miss
                                        |
                                        v
                                  Object Storage
                                        |
                                        v
                                  Video segments
```

And independently:

```text
Playback
   |
   +---- CDN → video segments
   |
   +---- View/watch events
               |
               v
          Event Stream
               |
       +-------+-------+
       |       |       |
     Views  Analytics  Recommendations
```

This is now a fairly realistic **global YouTube playback architecture**.

---

### Interviewer:
> A video suddenly goes viral and receives millions of requests per second. Are you comfortable with this architecture?

I'd say:

> Yes, because the high-volume traffic is primarily absorbed by the CDN rather than hitting our API servers, SQL database, or object-storage origin directly. CDN caching, request coalescing, and multiple edge locations prevent a viral video's traffic from translating one-for-one into origin traffic. The remaining watch events are handled asynchronously through the event pipeline rather than synchronously updating the video row.

The next issue I'd tackle is **YouTube search**. Search has a very different workload from ordinary SQL metadata lookup: we need full-text search across an enormous number of videos, relevance ranking, indexing, freshness, and eventually personalized results. That's still squarely part of our YouTube design.

## YouTube — Search

### Interviewer:
> Now design search. A user types "football highlights" and expects relevant YouTube videos. How would you handle that at YouTube scale?

### Candidate:

I wouldn't query our SQL database directly with something like:

```sql
SELECT *
FROM Videos
WHERE Title LIKE '%football highlights%';
```

At YouTube scale, that wouldn't be an appropriate search architecture.

We need a **dedicated search index**.

---

## 1. Search architecture

I'd introduce a search service:

```text id="b9d0h1"
User
 |
 | GET /search?q=football+highlights
 v
Search API
 |
 v
Search Service
 |
 v
Search Index
```

The SQL database remains the authoritative source for video metadata.

The search index is a **derived representation optimized for searching**.

For example:

```text id="6w3v7c"
SQL
Video 123
Title       = "Best Football Highlights"
Description = "..."
ChannelId   = 456
Status      = READY
```

The search index might contain:

```text id="1c2p8d"
Document 123
-----------------------------
title: "Best Football Highlights"
description: "..."
channel_id: 456
...
```

The search index is optimized for things such as:

- full-text matching
- relevance
- ranking
- filtering
- autocomplete

rather than transactional updates.

---

## 2. How does data get into the search index?

We shouldn't make video creation synchronously depend on search indexing.

Instead:

```text id="e5f0l7"
Video created/updated
        |
        v
     Event
        |
        v
   Event Stream
        |
        v
 Search Indexer
        |
        v
 Search Index
```

So when a creator uploads a video:

```text id="g9s2kj"
SQL
 |
 +-- Video metadata
 |
 +-- event
       |
       v
   Search Index
```

This means the search index can be **eventually consistent**.

For example, a video might become `READY` at 10:00:00 but appear in search at 10:00:02.

That's generally acceptable.

---

## 3. What happens when the title changes?

Suppose the creator changes:

```text id="e1w4wz"
"Berlin Travel"
```

to:

```text id="k2v7mc"
"Berlin Travel Guide 2026"
```

SQL is updated first.

Then:

```text id="4c0y4m"
VideoUpdated event
        |
        v
Search Indexer
        |
        v
Search document updated
```

Again, there's potentially a tiny delay.

We don't need a distributed transaction between SQL and the search engine.

---

## 4. Search result shouldn't contain everything

Suppose the search index returns:

```text id="i9x4q0"
Video IDs:
123
456
789
```

We then need metadata for those results.

We have two broad options.

### Option A — Store enough metadata in the search index

The index contains:

```text id="3u0c2m"
VideoId
Title
Thumbnail
ChannelName
Duration
PublishedAt
...
```

Then the search response can be generated directly from the index.

This is usually preferable for a search-heavy workload.

### Option B — Fetch metadata from SQL

```text id="qz3e5y"
Search
 |
 v
[123,456,789]
 |
 v
SQL
 |
 v
metadata
```

But this introduces additional database traffic and latency.

So I'd keep the fields required to render search results in the search index itself.

(The index is still not the authoritative source; it's a searchable projection of the underlying data.)

---

## 5. Ranking

Finding matching videos isn't enough.

Suppose there are 10 million videos matching:

> football

We need to decide which 20 appear first.

Ranking could consider signals such as:

- textual relevance
- watch time
- engagement
- freshness
- video quality
- language
- location
- user preferences
- personalization

Conceptually:

```text id="y8v3m1"
Query
  |
  v
Candidate retrieval
  |
  v
Potentially thousands of videos
  |
  v
Ranking
  |
  v
Top 20
```

At first, I'd keep ranking relatively simple.

For example:

```text id="0kz3vr"
Text relevance
+
Popularity
+
Freshness
```

Then we can introduce more sophisticated recommendation/personalization models later.

(Again: start simple and scale the part that actually becomes a bottleneck.)

---

## 6. Search index itself needs to scale

One search server obviously isn't enough.

We can partition the index into **shards**.

```text id="c7j5q8"
                    Search Index
                         |
          +--------------+--------------+
          |              |              |
       Shard 1        Shard 2        Shard 3
```

A query may need to search multiple shards:

```text id="z2v8n4"
"football highlights"
       |
       v
 +-----+-----+-----+
 |     |     |     |
 v     v     v     v
 S1    S2    S3   ...
 |     |     |
 +-----+-----+
       |
       v
 candidate results
       |
       v
     ranking
```

This is one of the places where a distributed search system naturally performs a **scatter/gather** operation.

But we'd keep the shards distributed and replicated so that one node failing doesn't make search unavailable.

---

## 7. Replication

We can have replicas of search shards:

```text id="m7x2za"
Shard 1
  |
  +-- Replica A
  +-- Replica B

Shard 2
  |
  +-- Replica A
  +-- Replica B
```

If one search node fails, another replica can serve that shard's portion.

Search doesn't need to be perfectly consistent at every instant.

If a newly uploaded video takes a few seconds to appear, that's acceptable.

---

## 8. What if the search index completely disappears?

This is another useful distinction.

The search index is **derived data**.

Our authoritative video metadata remains in SQL.

So theoretically:

```text id="w3s9x1"
SQL
 |
 | authoritative
 v
Video metadata
```

can rebuild:

```text id="x1v7q0"
SQL
 |
 v
Re-indexing pipeline
 |
 v
Search Index
```

Of course, rebuilding billions of videos would take a long time, so we'd maintain backups/snapshots and replicated indexes in practice.

But the architectural principle is important:

> **Don't make a derived search index the only copy of critical metadata.**

---

## 9. Global search

Now we're back to our global YouTube requirement.

A user in Germany shouldn't necessarily have to query a search cluster physically located in the US.

We can deploy search infrastructure regionally:

```text id="j2p5w8"
                  Global Users
                       |
              Global Traffic Routing
                       |
          +------------+------------+
          |            |            |
         EU           US          Asia
       Search       Search        Search
```

The indexes can be replicated across regions.

So:

```text id="w4k8p2"
Germany
   |
   v
EU Search
```

while:

```text id="s7m1q5"
Japan
   |
   v
Asia Search
```

gets low-latency access to search infrastructure.

---

## 10. But how do all regions get new videos?

Suppose a creator in Germany publishes:

```text id="x3n8a1"
Video 123
```

The event pipeline can distribute the update:

```text id="r7k4p9"
SQL / authoritative metadata
          |
          v
      VideoUpdated
          |
          v
     Event Stream
       /   |   \
      v    v    v
    EU    US   Asia
    Index Index Index
```

Now all regional search indexes eventually contain the video.

We accept some propagation delay.

---

## 11. Search autocomplete

There's actually another workload hidden inside search.

When the user types:

```text
foot
```

YouTube might suggest:

```text
football
football highlights
football skills
football live
```

This needs to be extremely fast.

We don't want to execute an expensive full search for every keystroke.

So I'd treat autocomplete as a separate workload:

```text id="c2y6r0"
"f"
 ↓
"fo"
 ↓
"foo"
 ↓
"foot"
```

with a specialized/autocomplete index or cached popular queries.

Because:

```text
GET /search/suggestions?q=foot
```

has very different requirements from:

```text
GET /search?q=football+highlights
```

---

## Search architecture now

Our YouTube system has another major subsystem:

```text id="a9p3l6"
                         YouTube
                            |
          +-----------------+------------------+
          |                 |                  |
          v                 v                  v
       SQL DB          Event Stream       Search System
          |                 |                  |
          |                 |             +----+----+
          |                 |             |         |
          |                 |          EU Index   US Index
          |                 |             |
          |                 |          Asia Index
          |                 |
          +-----------------+
                   |
             derived data
```

And the search flow:

```text id="r4c7v2"
User
 |
 v
Regional Search API
 |
 v
Search Index
 |
 v
Candidate videos
 |
 v
Ranking
 |
 v
Top results
```

while the authoritative data remains:

```text id="h5n8q1"
Video Metadata
      |
      v
     SQL
```

---

### Interviewer:
> Why not just use SQL with indexes for search?

I'd answer:

> SQL indexes are excellent for structured lookups and filtering, but YouTube search needs large-scale full-text retrieval, relevance scoring, ranking, autocomplete, and distributed search across an enormous corpus. A dedicated search index is optimized for that workload, while SQL remains the source of truth for transactional metadata.

---

At this point we've covered the major **content path** and **search path**.

The next major YouTube workload is **the homepage/feed and recommendations**. That's where our view, like, subscription, and watch-history events become useful: turning billions of user interactions into personalized video candidates without making the homepage request perform a massive computation.