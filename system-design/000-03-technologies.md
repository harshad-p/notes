# 1. Caching

Used to avoid repeatedly doing expensive work or repeatedly retrieving the same data.

Examples:
- In-memory cache — `.NET MemoryCache`
- Distributed cache — **Redis**, Memcached
- CDN / edge caching — Cloudflare, Amazon CloudFront, Azure Front Door

---

# 2. Databases

## Relational / SQL

Data is organized into tables with relationships, schemas, constraints, transactions, joins, etc.

Examples:
- **SQL Server**
- PostgreSQL
- MySQL
- Oracle
- Azure SQL Database
- Amazon Aurora

## NoSQL

Several different types exist.

**Document databases**
- **MongoDB**
- **Azure Cosmos DB** (supports multiple APIs/data models; commonly used as a NoSQL database)

**Key-value databases**
- **DynamoDB**
- Redis

**Wide-column databases**
- Cassandra
- ScyllaDB

**Graph databases**
- Neo4j
- Amazon Neptune

---

# 3. Message Queues

Used when one component needs to send work to another component **asynchronously**.

Examples:
- **RabbitMQ**
- Amazon SQS
- Azure Service Bus
- ActiveMQ

Typical use:

> API receives request → puts work into queue → worker processes it later.

This gives you decoupling, buffering, retries, etc.

---

# 4. Event Streaming / Distributed Messaging

This is where **Kafka** belongs.

Examples:
- **Apache Kafka**
- Apache Pulsar
- Amazon Kinesis
- Azure Event Hubs

Kafka is particularly useful when you have a stream of events that multiple consumers may independently consume and when you care about things like **retention, replay, partitions, and high throughput**.

For example:

> Listing created → Kafka → Search indexing, notification service, analytics service, etc.

---

# 5. Search

This is where **Elasticsearch** belongs.

Examples:
- **Elasticsearch**
- OpenSearch
- Apache Solr

Used for things that go beyond ordinary database querying, such as:

- Full-text search
- Fuzzy search
- Relevance ranking
- Autocomplete
- Faceted search
- Geographic search

For our Immowelt example:

> "Show me 2-bedroom apartments within 5 km of Berlin Mitte under €1,500."

A dedicated search engine can be extremely useful for this kind of workload.

---

# 6. API / Communication

Ways for services to communicate.

**HTTP/REST**
- ASP.NET Core Web API
- FastAPI
- Spring Boot

**gRPC**
- gRPC
- ASP.NET Core gRPC

**GraphQL**
- Apollo
- Hot Chocolate (.NET)

**WebSockets**
- ASP.NET Core SignalR
- native WebSockets

---

# 7. Containers / Deployment

Used to package and run applications consistently.

Examples:
- Docker
- Kubernetes
- Azure Kubernetes Service (AKS)
- Amazon EKS
- Amazon ECS

---

# 8. Load Balancing

Distributes requests across multiple application instances.

Examples:
- NGINX
- HAProxy
- AWS Elastic Load Balancer
- Azure Load Balancer
- Azure Application Gateway

---

# 9. Object / File Storage

For large files rather than structured database records.

Examples:
- Amazon S3
- Azure Blob Storage
- Google Cloud Storage

Typical use:

> Property images → Object storage

rather than storing the actual image bytes inside SQL Server.

---

# 10. Observability

Used to understand what the system is doing in production.

**Metrics**
- Prometheus
- Azure Monitor
- Amazon CloudWatch

**Logs**
- Elasticsearch / OpenSearch
- Splunk
- Serilog + centralized log storage

**Tracing**
- OpenTelemetry
- Jaeger
- Zipkin

Commercial platforms:
- Datadog
- New Relic
- Grafana

---

# 11. Identity / Authentication

Used for authentication and authorization.

Examples:
- Microsoft Entra ID (Azure AD)
- Auth0
- Keycloak
- Amazon Cognito
- Okta

Protocols you'll encounter:
- OAuth 2.0
- OpenID Connect
- SAML

---

# 12. Configuration / Secrets

Used to keep configuration and sensitive credentials out of application code.

Examples:
- Azure Key Vault
- AWS Secrets Manager
- HashiCorp Vault
- Kubernetes Secrets

---

# 13. Background Job Processing

Similar to queues, but specifically focused on scheduled/background application work.

Examples:
- Hangfire
- Quartz.NET
- Celery
- Azure Functions
- AWS Lambda

For example:

> Every night at 2 AM → process expired saved searches.

---

# The compact version I'd actually memorize

| Category | Examples |
|---|---|
| **Cache** | Redis, MemoryCache, Memcached |
| **SQL DB** | SQL Server, PostgreSQL, MySQL, Oracle |
| **NoSQL** | MongoDB, Cosmos DB, DynamoDB, Cassandra |
| **Message Queue** | RabbitMQ, SQS, Azure Service Bus |
| **Event Streaming** | Kafka, Kinesis, Event Hubs |
| **Search** | Elasticsearch, OpenSearch, Solr |
| **API communication** | REST, gRPC, GraphQL |
| **Containers** | Docker, Kubernetes |
| **Load Balancer** | NGINX, AWS ELB, Azure Application Gateway |
| **Object Storage** | S3, Azure Blob Storage |
| **Observability** | OpenTelemetry, Prometheus, Grafana, Datadog |
| **Authentication** | Entra ID, Auth0, Keycloak |
| **Secrets** | Key Vault, Secrets Manager, Vault |
| **Background Jobs** | Hangfire, Quartz, Lambda, Functions |

And for **your interview specifically**, I would make sure you can explain the differences between these particularly well:

**SQL vs NoSQL → Redis → RabbitMQ vs Kafka → Elasticsearch → Load Balancer → Docker/Kubernetes → Object Storage → Observability.**

Those are much more valuable than memorizing dozens of product names.