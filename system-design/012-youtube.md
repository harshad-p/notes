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