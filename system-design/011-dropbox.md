# Distributed File Storage & Sharing System

**Interviewer:** Design a Dropbox-like file storage and sharing system.

**Me:**

> **Before I design it, I'd like to clarify a few requirements.**

#### 1. Functional requirements

I'd like the system to support:

- Uploading files.
- Downloading files.
- Creating folders and organizing files.
- Renaming and deleting files.
- Sharing files or folders with other users.
- Permissions such as viewer and editor.
- Listing the files in a folder.
- Handling large files efficiently.

(I'll assume version history and real-time collaborative editing are out of scope for the initial design.)

#### 2. Non-functional requirements

I'd expect:

- Millions of users.
- Potentially billions of files.
- Large files, potentially several GB.
- High availability.
- High durability—we must not lose uploaded files.
- Reasonable upload/download performance.
- Metadata operations should have low latency.
- File contents don't need strong consistency in the same way that permissions and metadata do.

### 3. What data do we need?

There are really two different kinds of data here: **file metadata** and **file contents**.

For metadata I'd have something like:

**File**

- `id`
- `owner_id`
- `parent_folder_id`
- `name`
- `size`
- `content_type`
- `storage_key`
- `status`
- `created_at`
- `updated_at`

And for sharing:

**FilePermission**

- `file_id`
- `user_id`
- `permission`
- `created_at`

The actual file contents could be huge, so I would **not store them directly in the relational database**.

I'd store them in object storage.

### 4. Database choice

For metadata, I'd choose a **relational database**.

We have relationships between users, files, folders and permissions, and operations such as sharing and permission changes benefit from transactional consistency.

For the actual file bytes, I'd use **object storage** designed for large objects and high durability.

So conceptually:

- Relational DB → metadata and permissions
- Object storage → file contents

### 5. APIs

I'd have APIs such as:

```text
POST   /files/upload
GET    /files/{fileId}
GET    /folders/{folderId}/files
PATCH  /files/{fileId}
DELETE /files/{fileId}

POST   /files/{fileId}/share
DELETE /files/{fileId}/share/{userId}

GET    /files/{fileId}/download
```

For large files, though, I wouldn't send the entire file through my API servers.

That's the first important architectural issue.

### 6. Simplest architecture

I'd start with:

```text
Client
   |
   v
API
   |
   +----> Metadata DB
   |
   +----> Object Storage
```

But there's an immediate problem.

If a user uploads a 5 GB file, sending 5 GB through my API server means the API server becomes a bottleneck.

So I'd change the upload flow.

### 7. Large file uploads

I'd have the API authenticate the user and create an upload session.

The API can then provide the client with a **pre-signed upload URL** or equivalent temporary authorization to upload directly to object storage.

The flow becomes:

```text
Client
   |
   | create upload
   v
API -----> Metadata DB
   |
   | temporary upload authorization
   v
Client ---------> Object Storage
```

The API therefore handles **control-plane operations**, while object storage handles the **data plane**.

This is an important distinction because our API servers don't need to carry the actual file bytes.

### 8. Chunked uploads

For very large files, I'd support multipart/chunked uploads.

For example, a 5 GB file could be split into chunks.

The client can upload several chunks independently, potentially retrying only the failed chunk rather than restarting the entire upload.

Object storage can then assemble the completed object.

I'd keep an upload session containing things like:

- upload ID
- file ID
- expected size
- uploaded chunks
- status
- expiration time

If the client disappears halfway through, the incomplete upload can eventually be cleaned up.

### 9. Completing an upload

I don't want the metadata DB to say:

> `status = COMPLETE`

before the object actually exists.

So I'd have a completion step.

The client tells the API that the upload has finished. The backend verifies the object exists and has the expected properties, then marks the metadata as complete.

For example:

```text
UPLOADING
    |
    v
UPLOADED
    |
    v
AVAILABLE
```

(I'd keep this state explicit because partial uploads and failed uploads are important edge cases.)

### 10. Sharing and permissions

Sharing is primarily a metadata operation.

Suppose Alice shares a file with Bob.

I'd create or update:

```text
FilePermission
-------------------------
file_id
user_id
permission
```

The authorization check happens before returning the download authorization.

Importantly, I wouldn't make the object-storage URL itself permanently public.

Instead:

1. User requests download.
2. API authenticates user.
3. API checks permission.
4. API generates a short-lived download authorization.
5. Client downloads directly from object storage.

That way, removing Bob's permission prevents him from obtaining new download URLs.

### 11. Caching

Folder listings and file metadata are much more cacheable than the actual file contents.

I'd consider caching:

- folder listings
- file metadata
- permission information where appropriate

I wouldn't try to put multi-GB files into Redis or another application cache.

Object storage already handles the large-content delivery problem.

If the application is running on multiple API servers, I'd use a distributed cache for shared metadata caching.

### 12. Scaling

API servers can be horizontally scaled behind a load balancer.

The object storage system scales independently and handles the heavy bandwidth workload.

The metadata DB may eventually become the bottleneck.

I'd first optimize queries and indexes, then consider:

- read replicas for read-heavy metadata operations
- partitioning if the dataset becomes very large
- database sharding if we eventually reach the point where a single database cannot handle the workload

For example, folder listing is likely to be a very common operation, so I'd make sure the database has an efficient index around the parent-folder relationship.

### 13. Asynchronous processing

There are several things that don't need to block the upload response.

For example:

- virus/malware scanning
- generating thumbnails/previews
- extracting metadata
- generating search indexes
- sending sharing notifications

After an upload completes, I'd publish an event to a queue.

Workers can process these operations independently.

```text
Object Storage
      |
      v
 Upload Completed Event
      |
      v
    Queue
   /  |  \
  /   |   \
Scan Preview Metadata
```

I'd use retries with exponential backoff for transient failures and a dead-letter queue for messages that repeatedly fail.

(I wouldn't make virus scanning part of the synchronous upload path if the requirement allows asynchronous processing.)

### 14. Reliability

The most important reliability requirement is **file durability**.

I'd rely on object storage's replication/durability guarantees rather than implementing file replication myself at the application layer.

For metadata, I'd use database backups and replication.

If an API server dies during an upload, the client can retry the upload using the upload session.

If a particular chunk fails, only that chunk needs to be retried.

If a queue worker fails, the message can be retried.

If a notification fails, it shouldn't cause the file upload itself to roll back.

### 15. Observability

I'd monitor:

- upload/download latency
- upload failure rate
- chunk retry rate
- object-storage errors
- metadata DB latency
- queue depth
- processing failures
- incomplete upload count
- cache hit rate
- storage capacity
- permission-denied rate

I'd use metrics, structured logs and distributed traces to investigate failures across the API, database, queue and workers.

### Final architecture

I'd summarize the design as:

```text
                         Client
                           |
                    Load Balancer
                           |
                    API Servers
                    /         \
                   /           \
           Metadata DB       Cache
                   |
                   |
             Upload Session
                   |
                   v
              Object Storage
                   |
                   v
             Upload Completed
                   |
                 Queue
              /    |     \
             /     |      \
          Scan   Preview  Metadata
```

The important architectural decision is that **the API manages metadata, authorization and upload sessions, while the client transfers large file contents directly to object storage**.

That keeps the API tier lightweight and lets the storage layer scale independently.