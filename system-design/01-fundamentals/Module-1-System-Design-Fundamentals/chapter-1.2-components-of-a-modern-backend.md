# 1.2 Components of a Modern Backend

A backend is the server-side part of an application: it receives requests, applies business rules, reads or changes data, and returns results. A production backend is rarely one program and one database. It is usually a set of components, each responsible for a specific kind of work.

The goal is not to memorize a standard diagram. Start with the problem a component solves, then decide whether the problem exists in the system being designed. A small internal tool might need only an API server and a database; a globally used photo service may need most of the components below.

## A simple example: a photo-sharing application

Assume users can upload photos, view profiles, and search public captions. The authoritative record of a user and a photo's metadata—owner, caption, upload time, and visibility—is stored in a database. The large photo file itself is stored separately. This distinction matters: metadata is source-of-truth transactional data, while a thumbnail or a search index is derived data that can be recreated from the authoritative data.

```text
Mobile app / browser
        |
       DNS
        |
   Load balancer
        |
    API servers ------------------> Database
        |                               |
        |                               +--> cache (temporary copies)
        |
        +--> Object storage (original photo files)
        |
        +--> Message queue --> image-processing workers
                                      |
                                      +--> thumbnails in object storage

Browser <---- CDN <---- public images from object storage
Search request --> Search engine (derived index)
```

This is one possible topology, not a required architecture. We will now establish what each component does.

## Client

A **client** is software that uses the backend: a browser, mobile application, desktop application, another backend service, or a command-line tool. It initiates requests, renders results, and often holds user-facing state such as the current screen or unsaved form input.

The client should not be trusted merely because it is an official application. A user can modify browser requests or write their own client. The backend therefore validates requests and decides what data or operations the caller may access. We will treat API design and authorization in depth later; for now, the key point is that the client is outside the backend's trust boundary.

## DNS

**Domain Name System (DNS)** translates a human-readable name such as `api.photos.example` into a network destination. Clients use the name; the network needs an address. DNS lets an operator change where the name points without releasing an updated mobile app or changing every client.

For example, DNS can direct `api.photos.example` to a load balancer rather than directly to any one API server. DNS is a naming and initial routing mechanism. It is not itself the application server that processes the photo request. The precise lookup and connection sequence comes in the next lesson.

## Load balancer

A **load balancer** receives incoming network traffic and distributes it across multiple healthy backend servers. If three identical API servers can process requests, sending all traffic to one wastes two servers and makes that one server a single point of failure. A load balancer routes requests so the group behaves as one API endpoint.

Health checking is the load balancer's essential companion responsibility: it must avoid sending new requests to a server that is unavailable. This improves availability of the API tier, but it does not make a database or another dependency healthy.

For the photo service, the load balancer might send one profile request to API server A and the next to API server B. That works well only if either server can handle the request. In the next lessons, we will connect this requirement to state and horizontal scaling.

## API server

An **API server** is the application component that implements the service's operations. It receives an HTTP request such as `POST /photos`, validates it, applies business rules, coordinates calls to storage systems, and returns an HTTP response.

When a user creates a photo record, the API server might verify the user's identity, confirm that the uploaded object exists, and write the metadata row to the database. The database remains the source of truth for that metadata. The API server contains the decision-making logic; it should not treat a cache or a search result as the definitive answer for a transaction.

Multiple API servers are usually replicas of the same application code. They provide more request capacity and allow service to continue if one instance fails. They do not automatically solve database capacity or data consistency problems.

## Database

A **database** stores structured data that the application needs to query and change reliably. For the photo service, examples are users, photo metadata, comments, and access permissions. It provides durable storage and query capabilities that application memory cannot provide.

The database is commonly the source of truth for transactional application data: the canonical record the system relies on to decide whether a photo exists and who may see it. Other systems may hold copies optimized for speed or a particular access pattern, but they must be understood as derived from this authoritative data unless the design explicitly defines another source of truth.

A database is not automatically the right place for every byte. Storing multi-megabyte photo files inside relational rows makes database storage, backup, and serving paths carry a workload they are not optimized for. Object storage addresses that different problem.

## Cache

A **cache** stores temporary copies of data so a later request can be answered faster or without repeating expensive work. A profile read might fetch user information from the database the first time, then reuse a cached copy for a short period.

The primary benefit is lower latency and less load on the database. The trade-off is that a cache can be missing an entry or holding an older copy. Therefore, a cache is normally not the source of truth. On a cache miss, the application obtains the authoritative data from the database; after a write, the cached copy must eventually be updated, removed, or allowed to expire. The exact strategies come in the caching module.

A cache helps only when the same data or computation is requested again. It is not useful merely because an application has data.

## Object storage

**Object storage** stores large, unstructured blobs such as images, video, documents, and backups. Data is addressed by an object key, for example `photos/42/original.jpg`, and usually kept in a logical container called a bucket.

For the photo service, the original image and generated thumbnails belong in object storage; the database stores metadata and the object key. This separation lets the database efficiently answer questions such as “which public photos belong to user 42?” while object storage efficiently holds and delivers large immutable files.

Object storage is not a replacement for a database: it does not offer the same transactional querying model for application records. Conversely, a database is often a poor primary store for large media files. We will examine the model and trade-offs in the storage module.

## Message queue

A **message queue** holds units of work until a worker can process them. A producer places a message on the queue; a consumer, often called a worker, receives and performs the task. This separates the API request from work that is slow or does not need to finish before the user receives a response.

After an image upload, the API server can record the photo and add a “generate thumbnails for photo 42” message. A worker processes it later and writes thumbnail objects. The user can receive an upload-success response without waiting for CPU-intensive image conversion.

The queue smooths bursts and lets workers scale separately from API servers. It does not make the work disappear: if messages arrive faster than workers can complete them, waiting time grows. It also means thumbnail availability becomes asynchronous; the original photo metadata can exist before a thumbnail has been generated. The message queue is a mechanism for distributing tasks, not an authoritative database of the photo record.

## Event stream

An **event** is a record that something happened, such as `PhotoUploaded` or `UserRegistered`. An **event stream** is a durable, ordered record of events that can be retained and read by multiple independent consumers. Each consumer can use the same event history for its own purpose.

For example, one consumer may update a search index after `PhotoCaptionUpdated`; another may calculate daily upload metrics. Neither becomes the source of truth for the photo record: both construct derived data from events originating in the authoritative write path.

An event stream differs from a work queue in intent. A queue assigns a task to one worker so it gets done once. An event stream publishes a fact that several consumers may independently observe, including later consumers that need to replay retained events. We will explore this distinction before the detailed messaging lessons.

## Search engine

A **search engine** is a storage and query system optimized for finding relevant text or documents. A relational database can filter rows, but it is not generally optimized for features users expect from search: tokenization, relevance ranking, partial text matching, and language-aware analysis.

The photo service might place public photo captions and tags into a search index. When a caption changes, an asynchronous process updates that index. The index is derived data: a fresh database record might not appear in search immediately, and a stale search result still needs an API-level permission check before revealing private content. The search engine is chosen when search behavior justifies the additional data pipeline and operational cost.

## CDN

A **content delivery network (CDN)** is a geographically distributed set of servers that caches and delivers content close to users. It is especially useful for large, frequently requested, mostly immutable public files such as photo thumbnails, JavaScript bundles, and videos.

Without a CDN, every viewer may fetch a popular image from the photo service's origin storage in one region. With a CDN, a nearby edge server can serve a cached copy after the first request in that area. This reduces latency for distant users and reduces traffic reaching the origin.

A CDN is not the database and is not a general solution for private, constantly changing responses. It needs explicit rules for what may be cached and how long a cached copy remains valid. Later, we will study those cache-control decisions.

## API gateway

An **API gateway** is an edge component that sits in front of backend APIs and applies API-wide concerns before routing a request to the appropriate service. Examples include authentication checks, request logging, rate limits, protocol translation, and combining responses from several backend services.

An API gateway is distinct from a load balancer. A load balancer's core job is distributing traffic among healthy instances of one backend tier. A gateway's core job is shaping and governing API traffic. A real deployment can use both: the gateway decides that `/photos` belongs to the photo API, then a load balancer selects a healthy instance of that API. Some products combine capabilities, but the conceptual responsibilities remain different.

Do not add a gateway just because an architecture diagram has one. It becomes valuable when several APIs need shared edge behavior or a single stable API surface. For one small service, the complexity may not be justified.

## Recap

- The client uses DNS to reach the backend; DNS maps the service name to a network destination.
- A load balancer distributes requests across healthy API servers. API servers implement business rules and coordinate storage.
- A database commonly stores authoritative transactional data. Caches, search indexes, and often CDN copies are derived data optimized for a particular read path.
- Object storage holds large blobs; a message queue distributes tasks; an event stream lets several independent consumers observe retained facts.
- A CDN accelerates suitable content near users. An API gateway applies shared API-edge behavior and routes to backend services; it is not the same thing as a load balancer.

## Check your understanding

The photo service has an API server, a database, and object storage. Uploading a photo now takes eight seconds because thumbnail generation runs inside the upload request. Which component would you add first to separate thumbnail generation from the request, and what user-visible state must the system represent while that work has not finished?
