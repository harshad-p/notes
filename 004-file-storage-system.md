Absolutely. We'll add a **hospitality-specific system design problem** later, ideally one that resembles the kinds of backend/platform problems a company like Apaleo could have. I won't force it into every exercise.

For now:

# System Design #4 — Design a File Storage System

Think **Google Drive / Dropbox**, but we'll keep the initial requirements modest.

### The problem

> **Design a system where users can upload files, download them later, and organize them into folders.**

Let's assume:

- Users can upload files up to **5 GB**.
- They can download their files.
- Files can be organized into folders.
- Files should not be lost.
- Multiple users may use the system simultaneously.
- We should be able to scale to millions of users.

---

## Step 1 — Identify the two kinds of data

The first thing I'd separate is:

### File contents

The actual:

```text
photo.jpg
video.mp4
document.pdf
```

These can be **very large**.

### Metadata

Information *about* the file:

```text
File ID
Name
Size
Owner
Folder
CreatedAt
Storage location
```

The metadata is relatively small.

So I wouldn't put the actual 5 GB file into a normal relational database.

Instead:

```text
                 File Storage System
                        │
              ┌─────────┴─────────┐
              ▼                   ▼
        Metadata DB          Object Storage
                             
                           actual files
```

For something like AWS, object storage would be **S3**.

---

## Step 2 — Upload flow

A naïve design would be:

```text
User
 ↓
Our API
 ↓
Our server
 ↓
S3
```

But there's a problem.

If 10,000 users upload 5 GB files, our servers become a bottleneck because all those large files are flowing through them.

Instead, we can let the client upload **directly to object storage**.

```text
User
  │
  │ 1. "I want to upload photo.jpg"
  ▼
Our API
  │
  │ 2. Generate upload URL
  ▼
User
  │
  │ 3. Upload directly
  ▼
S3
```

This is a very important pattern.

The API handles **metadata and authorization**, while object storage handles the large file transfer.

---

## Step 3 — Why a pre-signed URL?

Our API can generate a temporary URL that gives the user permission to upload a specific object.

Conceptually:

```text
POST /files
        ↓
API authenticates user
        ↓
Creates file metadata
        ↓
Generates temporary upload URL
        ↓
Returns URL to client
```

The client then uploads directly to S3.

The URL can be:

- temporary
- restricted to a particular object
- restricted to an operation such as upload

So we don't have to give the client permanent AWS credentials.

---

### Our initial architecture

```text
                         User
                           │
                           ▼
                      File API
                       /     \
                      /       \
                     ▼         ▼
              Metadata DB    S3
                              ▲
                              │
                         Direct upload
                              ▲
                              │
                            User
```

Already we have an important system-design principle:

> **Don't send large data through your application servers when a specialized storage service can handle it directly.**

Next we'll tackle something that becomes important with **5 GB files**:

> **What happens if the upload fails halfway through?**

That's where **multipart/resumable uploads** come in.

## Step 4 — Large and interrupted uploads

Suppose a user is uploading a **5 GB video**.

If we treat it as one request:

```text
5 GB
──────────────→ S3
              ↓
           Network fails
              ↓
           Upload lost
```

The user may have to start all over again.

That's not a good experience.

### Multipart upload

Instead, we split the file into parts:

```text
5 GB file

┌────┬────┬────┬────┬────┬────┐
│ P1 │ P2 │ P3 │ P4 │ P5 │ ...│
└────┴────┴────┴────┴────┴────┘
```

Each part is uploaded separately:

```text
P1 ✓
P2 ✓
P3 ✓
P4 ✗
P5 ✓
```

If P4 fails, we only retry P4.

```text
P4 ✗
 ↓
retry
 ↓
P4 ✓
```

S3 (and similar object-storage systems) can then assemble the parts into the final object.

---

## Why this is better

It gives us:

- **Resumability** — interrupted uploads can continue.
- **Parallelism** — multiple parts can upload simultaneously.
- **Better failure handling** — retry only failed parts.
- **Better performance** — particularly for large files.

For our 5 GB requirement, I'd definitely use multipart uploads.

---

## Step 5 — What does our database know?

We need to track the file independently from its physical storage.

For example:

```text id="c6x0pz"
File
-------------------------
id          847291
name        vacation.mp4
ownerId     123
size        5 GB
status      UPLOADING
storageKey  users/123/847291
createdAt   ...
```

During upload:

```text id="v5bx6w"
UPLOADING
```

After successful completion:

```text id="x3tqj7"
AVAILABLE
```

If the upload fails permanently:

```text id="8h4j0w"
FAILED
```

This prevents the application from treating an incomplete object as a valid file.

---

## One important design detail

We shouldn't rely only on the client telling us:

> "Upload finished."

We should verify that the upload actually completed successfully.

Depending on the object-storage system, we can use an upload-completion operation, metadata checks, or storage events.

Conceptually:

```text id="c6a7ov"
S3
 ↓
"Object uploaded"
 ↓
Our system
 ↓
File status = AVAILABLE
```

This keeps our metadata consistent with what's actually in storage.

---

### Our design now

```text id="5ik0d8"
                         User
                           │
                           ▼
                       File API
                      /        \
                     ▼          ▼
              Metadata DB      S3
                                  ▲
                                  │
                       Multipart upload
                       ┌────┬────┬────┐
                       │ P1 │ P2 │ P3 │ ...
                       └────┴────┴────┘
```

The next interesting problem is **downloads and access control**:

> If I own a file, how do you let me download it directly from S3 without exposing the file publicly or giving me permanent access?

That brings us back to **pre-signed URLs, authorization, and security**, but in a more concrete way.

ok, do tell me about the next part. But also, I am noticing we are not going into the "how" or the technology aspect of it. For example, I understand if I don't talk about Redis, or DataDog, and just say "distributed cache", or observability. But here, things are getting a little too abstract. Maybe because of lack of my knowledge. You said we can let the user upload to S3 directly. I can see so much complexity in it. The user hits our API, our API generates an S3 URL for user, but then how do I connect the incoming stream to the generated S3 URL, etc. etc. Also, you said splitting the file. How does that happen? Maybe these patterns exist, but right now they sound very complex to me)

Yes — this is a good observation. In the previous problems, I was deliberately keeping technology abstract because the **architecture** was the important part. Here, that abstraction is hiding the actual mechanism.

For a system design interview, you should know both:

1. **What** the architecture does.
2. **How** the major pieces actually work.

You don't need to know the internals of S3, but you should understand the flow well enough to explain it.

---

# Step 5 — How does direct S3 upload actually work?

Your question is exactly right:

> "Our API gives the client a URL. How does the file actually get there?"

The important thing is: **our API does not receive or forward the file.**

Let's use a concrete example.

The user selects:

```text
vacation.mp4
5 GB
```

### 1. Client asks our API to start an upload

```http
POST /files/upload
```

The request contains metadata such as:

```json
{
  "fileName": "vacation.mp4",
  "size": 5368709120
}
```

Our API authenticates the user and decides:

> "Yes, this user is allowed to upload this file."

---

### 2. Our backend asks S3 to create an upload

Our backend uses the **AWS SDK**.

Conceptually:

```text
Our API
   │
   │ AWS SDK
   ▼
S3
```

S3 creates an **upload session** and gives our backend an `uploadId`.

Our API doesn't download or upload the file.

---

### 3. Our API gives information back to the client

Something like:

```json
{
  "fileId": "847291",
  "uploadId": "abc123",
  "storageKey": "users/123/847291"
}
```

Depending on the implementation, the backend can also generate **pre-signed URLs** for the individual upload parts.

---

### 4. The client talks directly to S3

This is the part that initially sounds strange.

The browser/mobile app already has the file.

It uses the URL provided by our backend to send the bytes **directly to S3**:

```text
                 5 GB file
                    │
                    │ HTTPS
                    ▼
User's Browser ─────────────→ S3
                             
                 NOT:
Browser → Our API → S3
```

Our API isn't involved in the data transfer.

The URL essentially says:

> "S3, temporarily allow this client to upload this particular object."

That's what the **pre-signed URL** is doing.

---

# But how does splitting the file work?

This happens **on the client**.

We don't have to write some complicated file-splitting system on our backend.

Modern browsers expose APIs that allow JavaScript to work with a file as a collection of **byte ranges**.

For example:

```text
5 GB file

0 MB ───── 10 MB
        Part 1

10 MB ─── 20 MB
        Part 2

20 MB ─── 30 MB
        Part 3
...
```

The browser can read a particular portion of the file and upload that portion.

Conceptually:

```javascript
const part = file.slice(start, end);
```

Then it sends that part to the corresponding S3 upload URL.

So:

```text
Browser
   │
   ├── Part 1 ─────────→ S3
   ├── Part 2 ─────────→ S3
   ├── Part 3 ─────────→ S3
   ├── Part 4 ─────────→ S3
   └── Part 5 ─────────→ S3
```

They can even be uploaded **in parallel**.

---

# What happens if Part 4 fails?

The browser knows:

```text
Part 1 ✓
Part 2 ✓
Part 3 ✓
Part 4 ✗
Part 5 ✓
```

So it retries only Part 4.

```text
Part 4
   ↓
retry
   ↓
S3 ✓
```

Once all parts are uploaded, the client tells our backend:

```http
POST /files/847291/complete
```

Our backend tells S3:

> "All the parts are uploaded. Combine them into the final object."

S3 does that.

Our database can then change:

```text
UPLOADING → AVAILABLE
```

---

# So the complete upload is actually:

```text
                 1. Start upload
User ───────────────────────────→ Our API
                                   │
                                   │ AWS SDK
                                   ▼
                                  S3
                                   │
                              uploadId
                                   │
                                   ▼
User ←─────────────────────────────┘


                 2. Upload bytes

User Browser
     │
     ├──── Part 1 ─────────────→ S3
     ├──── Part 2 ─────────────→ S3
     ├──── Part 3 ─────────────→ S3
     ├──── Part 4 ─────────────→ S3
     └──── Part 5 ─────────────→ S3


                 3. Complete

User ─────→ Our API ─────→ S3
             │
             ▼
        Metadata DB
        AVAILABLE
```

**That's the actual mechanism you should understand.**

You don't need to explain every HTTP request in a Karat interview unless the interviewer asks. But now you have the mental model behind the architecture.

---

# Step 6 — Downloads

Now let's reverse the process.

The user asks:

```http
GET /files/847291/download
```

Our API:

1. Authenticates the user.
2. Checks whether they have access to file `847291`.
3. Looks up the S3 object key.
4. Generates a temporary **pre-signed GET URL**.
5. Returns it.

```text
User
 │
 │ GET /files/847291/download
 ▼
Our API
 │
 │ authorization check
 ▼
Metadata DB
 │
 │ allowed
 ▼
Pre-signed S3 URL
 │
 ▼
User
 │
 │ directly downloads
 ▼
S3
```

Again, the large file doesn't travel through our API.

---

## Why not just make the S3 object public?

Because then someone who knows the URL could potentially access the file.

Instead, the pre-signed URL might only be valid for:

```text
15 minutes
```

After that:

```text
URL → expired ❌
```

So we get:

**private storage + controlled temporary access + scalable file transfer.**

---

## Now the technology should feel less abstract

For an AWS implementation, you could actually say:

> "I'd use S3 for object storage, a relational database such as PostgreSQL for metadata, and the AWS SDK from our backend to initiate multipart uploads and generate pre-signed URLs. The client would upload the parts directly to S3 rather than sending large files through our API."

That's a **much stronger system-design answer** than simply saying "object storage."

And I think this is the right balance for our remaining problems: **I'll introduce the actual technology and mechanism whenever the abstraction would otherwise hide something important.**

ok, so supporting direct S3 upload, and creating chunks needs a combination of support from the browser, AWS, and our API. 
Also, isn't there a risk of snooping where a url to allow temporary uploads and downloads is generated?

Exactly. **The whole flow is a collaboration between the client, our backend, and S3.** The browser handles reading/splitting the file, our API handles authentication/authorization and coordinates the upload, and S3 handles the actual storage and multipart upload.

And yes, there **is a security risk if a pre-signed URL is leaked**, but it's deliberately designed to limit the damage.

### What happens if someone gets the upload URL?

Suppose our API gives the browser:

```text
https://s3.../users/123/file-847291
       ?signature=...
       &expires=...
```

If an attacker gets that URL before it expires, they may be able to use it.

So a pre-signed URL is effectively a **temporary bearer credential**.

That's why we should:

- **Keep the expiry short** — e.g. a few minutes for an upload URL.
- Restrict it to **one specific object/key**.
- Restrict the operation — upload URL should only permit upload, download URL only download.
- Don't put sensitive information in the URL itself.
- Use HTTPS, so the URL isn't exposed through network sniffing.
- Have our backend perform the authorization **before generating the URL**.

### But there's an important limitation

Once we've generated:

```text
User A → pre-signed URL
```

S3 generally doesn't know *who* User A is. The URL itself is the authorization.

So if User A copies the URL and sends it to User B, **User B may be able to use it until it expires**.

That's why we don't treat a pre-signed URL as a replacement for authentication. Authentication and authorization happen **before** the URL is generated.

For highly sensitive systems, we can also use additional controls such as very short-lived URLs, encryption, network restrictions, or have the application proxy the download instead—but that sacrifices some of the scalability benefit.

---

### One subtle point about the chunks

You're also right that this isn't simply:

> "The browser creates chunks and S3 understands them."

There is a defined **S3 multipart-upload protocol**:

```text
Our API → S3: initiate multipart upload
                 ↓
              uploadId
                 ↓
Browser uploads Part 1, Part 2, Part 3...
                 ↓
S3 returns identifiers for uploaded parts
                 ↓
Our API/client → S3: complete multipart upload
                 ↓
S3 assembles the object
```

The AWS SDKs and browser/client libraries handle much of this machinery for us.

**For the interview, you don't need to implement it from scratch.** You need to know that multipart upload is an S3 capability, the client can upload parts directly using authorized URLs, and our backend controls who is allowed to initiate/complete the upload.

That distinction—**understanding the mechanism without pretending you need to reinvent the AWS SDK**—is exactly where I'd aim your system-design knowledge.

Let's finish the file-storage problem with one important area, then move on.

## Step 7 — Metadata, consistency, and deletion

We have:

```text
Files → S3
Metadata → PostgreSQL
```

Now consider deletion.

A user deletes `vacation.mp4`.

We need to delete both:

```text
PostgreSQL record
S3 object
```

But what if one succeeds and the other fails?

### Don't try to make them one transaction

PostgreSQL and S3 aren't part of the same database transaction.

So I'd use a simple workflow:

```text
DELETE /files/847291
        ↓
Mark metadata as DELETING
        ↓
Delete S3 object
        ↓
Delete/mark metadata as DELETED
```

If S3 temporarily fails:

```text
DELETING
   ↓
retry
```

This is another example where **eventual consistency is acceptable**.

The user doesn't necessarily need the S3 object and database row to disappear at precisely the same millisecond.

---

## What about orphaned files?

There's another failure scenario:

```text
S3 object ✓
Database record ✗
```

For example, the application crashed after uploading the file but before updating the database.

Now we have an **orphaned object** consuming storage.

A periodic cleanup job can find objects that don't have corresponding metadata and remove them.

Similarly, we can clean up abandoned multipart uploads.

This is a good example of why **background jobs** are useful even when the main request succeeds.

---

## Step 8 — Scaling the metadata database

S3 is designed to handle enormous amounts of object storage, so the metadata database is more likely to become our bottleneck.

For example:

```text
Users
  ↓
Load Balancer
  ↓
File API instances
  ↓
PostgreSQL
```

We can initially scale the API horizontally.

If database reads become heavy, we could introduce:

```text
PostgreSQL Primary
       │
       ├── Read Replica
       └── Read Replica
```

And potentially cache frequently accessed metadata with **Redis**.

We don't need sharding immediately. We'd introduce it only when the data volume or traffic actually requires it.

---

# That's enough for the file-storage problem

The main concepts I'd want you to retain are:

1. **S3/object storage** → actual file contents.
2. **PostgreSQL** → metadata.
3. **Pre-signed URLs** → direct, temporary client ↔ S3 access.
4. **Multipart upload** → large/resumable uploads.
5. **Background jobs** → cleanup and recovery.
6. **Horizontal API scaling** → handle many users.
7. **DB replicas/cache** → scale metadata reads.

And importantly, you now understand **how the S3 flow actually works**, rather than just memorizing "use object storage."

---

### Next: Hospitality system

I think this is a good point to switch domains.

Given your **Apaleo + Immowelt** preparation, I'd like to do:

> **Design a hotel room booking/reservation system.**

We'll make it more interesting than a generic CRUD booking API:

**"Design a system that allows guests to search hotel availability and book a room, while preventing double bookings."**

That will let us explore **availability, concurrency, database transactions, distributed locking, idempotency, payment, overbooking, and consistency**—all highly relevant to hospitality and excellent system-design interview material.