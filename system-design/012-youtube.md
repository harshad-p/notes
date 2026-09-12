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