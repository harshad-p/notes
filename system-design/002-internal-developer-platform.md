Yes. I think we've covered the URL shortener enough for now. You now understand the important concepts rather than just memorizing an architecture:

- API design
- database
- caching
- read-heavy systems
- horizontal vs vertical scaling
- replication vs sharding
- regional distribution
- unique ID generation
- Base62
- collision handling
- distributed IDs
- reliability and fault tolerance
- trade-offs

### Next problem: **Design an Internal Developer Platform**

This is the one I'd choose for your **Immowelt/AVIV interview** because it maps directly to their JD.

The interview question could be:

> **"Design an internal platform that allows developers to create, deploy, and monitor a new service."**

This lets us work directly with their keywords:

**self-service workflows**  
**golden paths**  
**engineering standards**  
**CI/CD**  
**automation**  
**platform services**  
**reusable building blocks**  
**observability**  
**developer experience**

And it connects nicely with your UPF experience with internal applications and your Git workflow.

---

### We'll learn it the same way

We'll start extremely simply:

**Step 1:** What problem are we solving?

**Step 2:** What does a developer actually want to do?

**Step 3:** Define the basic workflow.

**Step 4:** Design the simplest architecture.

Then we'll progressively ask:

> What happens when 100 teams use it?

> How do we standardize projects?

> How does CI/CD work?

> How do we make it self-service?

> How do we handle deployments?

> How do we monitor everything?

> What happens when something fails?

That will give you a much more realistic **platform-engineering system design** than jumping straight into a giant architecture diagram.

### So the next question we'll tackle is:

> **Design a platform where a developer can create a new service and deploy it to production without needing to manually configure repositories, CI/CD, infrastructure, and monitoring.**

This is the one I'd focus on next.

# System Design #2 — Internal Developer Platform

Let's treat this like an actual Karat question and build the answer **step by step**.

### The question

> **"Design an internal platform that allows developers to create, deploy, and monitor a new service."**

Don't design anything yet.

---

## Step 1 — Understand the problem

Imagine Immowelt has hundreds of developers and many product teams.

A developer wants to create a new service.

Without an internal platform, they might have to manually:

1. Create a Git repository
2. Set up the project
3. Configure CI/CD
4. Create cloud infrastructure
5. Configure deployment
6. Set up monitoring
7. Configure permissions
8. Set up alerts

That's a lot of repetitive work.

The platform's goal is essentially:

> **Give developers a simple, standardized path from "I need a service" to "my service is running."**

That's what the JD means by:

**self-service workflows** and **golden paths**.

---

# Step 2 — What does the developer actually do?

Let's keep the developer experience extremely simple.

Imagine the developer goes to an internal portal and selects:

```text
Create New Service
```

They enter:

```text
Service name: PaymentService
Language: C#
Environment: Production
```

Then they click:

```text
Create
```

The platform does the rest.

Conceptually:

```text
Developer
    ↓
Internal Portal
    ↓
"Create Service"
    ↓
Platform
    ↓
Repository
    ↓
CI/CD
    ↓
Infrastructure
    ↓
Deployment
    ↓
Monitoring
```

That's our basic workflow.

---

# Step 3 — What should the platform create?

When the developer clicks **Create**, we might automatically create:

### 1. Repository

```text
GitHub
   ↓
PaymentService
```

With a standard structure:

```text
PaymentService/
├── src/
├── tests/
├── README
├── Dockerfile
└── ...
```

### 2. CI/CD pipeline

Automatically configured:

```text
Commit
  ↓
Build
  ↓
Tests
  ↓
Security checks
  ↓
Build Docker image
  ↓
Deploy
```

### 3. Infrastructure

The platform could create the required infrastructure using **Infrastructure as Code**, such as Terraform.

For example:

```text
PaymentService
      ↓
Terraform
      ↓
AWS resources
```

### 4. Monitoring

Automatically configure:

- logs
- metrics
- alerts
- dashboards

Now the developer doesn't have to manually configure all of this.

---

# Step 4 — Our first architecture

Don't overcomplicate it.

Start with:

```text
                    Developer
                        │
                        ▼
                ┌──────────────┐
                │ Internal      │
                │ Developer     │
                │ Portal        │
                └──────┬───────┘
                       │
                       ▼
                ┌──────────────┐
                │ Platform API  │
                └──────┬───────┘
                       │
          ┌────────────┼────────────┐
          ▼            ▼            ▼
       GitHub        CI/CD      Terraform
                                    │
                                    ▼
                                   AWS
```

And after deployment:

```text
AWS
 │
 ▼
Monitoring / Observability
```

That's already a reasonable first architecture.

**Notice that we're not creating 20 microservices.**

---

# Step 5 — Where does the "golden path" come in?

Suppose every developer can create their own project however they want.

One developer uses:

```text
.NET 8
GitHub Actions
Terraform
Datadog
```

Another uses:

```text
.NET 9
CircleCI
some custom deployment
different logging
```

Another does something completely different.

Now the platform has failed at **standardization**.

Instead, we provide templates.

For example:

```text
Create Service
     │
     ├── .NET API
     ├── Node.js API
     └── Python API
```

Developer chooses:

```text
.NET API
```

The platform gives them the company's approved setup.

That's a **golden path**.

The important idea is:

> **The platform makes the recommended way easy, without necessarily preventing teams from doing something different when they have a good reason.**

That's an important distinction.

---

# Step 6 — Why is this called Developer Experience?

Think about the original JD:

> "Make building software fast, reliable, and enjoyable."

Without the platform:

```text
Developer
   ↓
2 days of setup
   ↓
Finally start coding
```

With the platform:

```text
Developer
   ↓
5-minute setup
   ↓
Start coding
```

That's **developer productivity / developer experience**.

And this is why the JD says:

> "Developer Experience is treated as a first class product."

The platform isn't just infrastructure.

The **developer is the customer**.

That's a very important concept for this interview.

---

## Pause here

Don't worry about CI/CD, AWS, Terraform, Kubernetes, observability, scaling, etc. yet.

We've established the **core idea**:

> **A developer asks the platform for a service, and the platform automatically creates a standardized, deployable, observable service.**

Next, we'll tackle the interesting question:

> **"What exactly happens inside the Platform API when the developer clicks Create?"**

That's where we'll introduce **self-service workflows, asynchronous processing, queues, retries, and failure handling** — and the architecture will start becoming much more realistic.

## Step 7 — What happens when the developer clicks "Create"?

We currently have:

```text
Developer
    ↓
Internal Portal
    ↓
Platform API
```

The developer clicks:

> **Create .NET Service**

What should the Platform API do?

At first, you might think:

```text
Platform API
    ↓
Create GitHub repo
    ↓
Create AWS resources
    ↓
Create CI/CD
    ↓
Configure monitoring
    ↓
Return "Done"
```

But there's a problem.

### This could take a long time.

Creating infrastructure, repositories, pipelines, permissions, etc. might involve several external systems.

What happens if GitHub succeeds but AWS fails?

You could end up with:

```text
GitHub       ✓
CI/CD        ✓
AWS          ❌
Monitoring   ❌
```

Now the developer has a half-created service.

So we need to think about **workflow management**.

---

# Step 8 — Make the workflow asynchronous

Instead of making the developer wait for everything:

```text
Developer
   ↓
Platform API
   ↓
"Create service"
   ↓
Do everything
   ↓
Response
```

we can do:

```text
Developer
   ↓
Platform API
   ↓
Create request
   ↓
Queue
   ↓
Worker
```

The API can immediately respond:

> "Your service is being created."

The worker then performs the actual work.

```text
                  Queue
                    ↓
                  Worker
               /    |    \
              ↓     ↓     ↓
           GitHub  AWS   CI/CD
```

This is a very common system-design pattern.

---

# Step 9 — Why use a queue?

There are several reasons.

### 1. The API doesn't have to wait

The user gets a quick response.

### 2. We can retry failures

Suppose:

```text
GitHub ✓
AWS ❌
```

The worker can retry the AWS operation.

### 3. We can handle many requests

Suppose 50 developers create services at the same time.

Instead of 50 API requests doing expensive work simultaneously:

```text
50 requests
     ↓
   Queue
     ↓
Workers process them
```

We can control how much work happens at once.

That's useful for **scalability**.

---

# Step 10 — But retries introduce another problem

Suppose our worker sends:

```text
Create repository
```

GitHub creates the repository successfully.

But the response gets lost.

The worker thinks:

> "GitHub failed."

So it retries.

Now we have:

```text
Create repository
     ↓
Repository already exists
```

We don't want retries creating duplicate resources.

This introduces an important concept:

## Idempotency

An operation is **idempotent** if performing it multiple times produces the same final result.

For example:

```text
Create service PaymentService
```

The worker should be able to safely retry it.

It might check:

```text
Does PaymentService already exist?
       ↓
      Yes
       ↓
Skip creation
```

rather than blindly creating another resource.

---

# Step 11 — Track the workflow

We also need to tell the developer what's happening.

So we can store a workflow record:

```text
Service: PaymentService

Status:
✓ Repository created
✓ CI/CD configured
⏳ Infrastructure creating
○ Monitoring
```

The portal could show:

```text
Creating PaymentService...

✓ Repository
✓ CI/CD
⏳ AWS infrastructure
○ Monitoring
```

This is part of good **Developer Experience**.

The developer doesn't have to wonder:

> "Did it work?"

---

# Step 12 — Our architecture is getting better

Now we have:

```text
                         Developer
                             │
                             ▼
                    ┌────────────────┐
                    │ Internal Portal│
                    └───────┬────────┘
                            │
                            ▼
                    ┌────────────────┐
                    │  Platform API  │
                    └───────┬────────┘
                            │
                            ▼
                         Queue
                            │
                            ▼
                         Worker
                      /      |      \
                     ▼       ▼       ▼
                  GitHub    AWS     CI/CD
                            │
                            ▼
                      Monitoring
```

And somewhere we have:

```text
              Workflow Database
                    ↑
                    │
                  Worker
```

which stores the status.

---

# This is a very important system-design pattern

You should recognize this structure:

```text
API
 ↓
Queue
 ↓
Worker
 ↓
External systems
```

It's useful when:

- work takes a long time
- external services are involved
- you need retries
- you need to process many requests
- you don't want your API request blocked

And it's directly relevant to the Immowelt JD because they specifically mention:

> **"self-service workflows"**

and

> **"automation to boost developer productivity and reduce manual effort."**

We're essentially designing the machinery behind those workflows.

---

### One question before we continue

Imagine the worker is halfway through:

```text
✓ GitHub
✓ CI/CD
✓ AWS
❌ Monitoring
```

and retries fail three times.

**What should happen?**

Should we:

1. Roll everything back?
2. Leave the partially created service?
3. Mark the workflow as failed and let someone fix it?
4. Something else?

That's our next system-design problem, and it introduces **failure handling, rollback, compensation, and recovery**.

I think 3, but maybe it's dependent upon the user. WHat they want. But I would choose 3. Anyway, go ahead with the next step.

Yes — **3 is the right default** for this kind of platform. And you're right that rollback can depend on the resource and the user's requirements. We shouldn't automatically delete a working service just because one later step failed.

# Step 13 — Failure handling and recovery

Our workflow currently looks like:

```text
Create Service
      ↓
   Queue
      ↓
   Worker
      ↓
 ┌────┼────┬────┐
 ↓    ↓    ↓    ↓
GitHub CI/CD AWS Monitoring
  ✓     ✓    ✓      ❌
```

Suppose monitoring setup fails three times.

We can mark the workflow:

```text
Status: FAILED

✓ Repository created
✓ CI/CD configured
✓ Infrastructure created
✗ Monitoring failed
```

The important thing is:

> **Don't hide the failure.**

The developer should know what happened and what they can do about it.

---

## Step 14 — Why not automatically roll everything back?

Imagine we created:

```text
GitHub repository       ✓
AWS infrastructure      ✓
CI/CD                   ✓
Monitoring              ✗
```

We could try to delete everything:

```text
Delete AWS
Delete CI/CD
Delete GitHub
```

But what if the developer has already started using the repository?

Or what if the AWS infrastructure was successfully created but the monitoring service had a temporary problem?

Automatically destroying everything could actually make things worse.

So for this platform, I'd generally prefer:

> **Keep successfully created resources, mark the workflow as failed, and provide a way to retry or fix the failed step.**

That's what you suggested.

---

# Step 15 — Retry only the failed step

Suppose the developer sees:

```text
PaymentService

✓ Repository
✓ CI/CD
✓ Infrastructure
✗ Monitoring

[Retry monitoring]
```

They click **Retry**.

We don't want to start from the beginning.

We want:

```text
Monitoring
   ↓
Retry
```

not:

```text
GitHub
   ↓
CI/CD
   ↓
AWS
   ↓
Monitoring
```

Again, this is where **idempotency** matters.

The worker should know:

> "The repository already exists, so don't create it again."

---

# Step 16 — What if the workflow keeps failing?

We don't want an infinite retry loop.

For example:

```text
Attempt 1 → failure
Attempt 2 → failure
Attempt 3 → failure
```

After some limit:

```text
Status: FAILED
```

Then we can:

- notify the developer
- expose the error
- allow manual retry
- log the failure
- alert the platform team if necessary

This gives us a clear recovery path.

---

# Step 17 — What should we log?

This brings us back to **observability**.

For every workflow, we want to know:

```text
Service: PaymentService
Workflow ID: 12345

10:01 Repository creation started
10:01 Repository created
10:02 CI/CD creation started
10:02 CI/CD created
10:03 AWS provisioning started
10:04 AWS provisioning completed
10:04 Monitoring setup started
10:04 Monitoring failed
10:05 Retry
10:05 Monitoring failed
```

This is much better than:

> "Something went wrong."

---

# Step 18 — What if the Platform API itself goes down?

Now let's think one level higher.

Suppose:

```text
Developer
    ↓
Platform API 💥
```

A developer tries to create a service.

The request might never reach the queue.

We can make the API **stateless** and run multiple instances:

```text
                 Load Balancer
                 /     |     \
                ↓      ↓      ↓
             API 1   API 2   API 3
```

If API 2 fails:

```text
API 1 ✓
API 2 ✗
API 3 ✓
```

The load balancer sends requests to the healthy instances.

That's our earlier concept of **horizontal scaling + fault tolerance**.

---

# Step 19 — What if the worker goes down?

Suppose:

```text
Queue
  ↓
Worker 💥
```

We don't want the workflow to disappear.

That's one of the major benefits of the queue.

The message remains in the queue until a worker successfully processes it, depending on the queue technology and acknowledgement model.

Another worker can pick it up:

```text
             Queue
            /     \
           ↓       ↓
       Worker 1  Worker 2
          💥        ✓
                    ↓
                Continue
```

Again:

**Distribution isn't just about scaling.**

Here we're using multiple workers primarily for **reliability and availability**, while they can also give us more processing capacity.

---

# Step 20 — Now we have a pretty solid platform

Our architecture is now:

```text
                         Developer
                             │
                             ▼
                    ┌────────────────┐
                    │ Internal Portal│
                    └───────┬────────┘
                            │
                            ▼
                    ┌────────────────┐
                    │ Load Balancer  │
                    └───────┬────────┘
                            │
                     ┌──────┴──────┐
                     ↓             ↓
                   API 1         API 2
                     │             │
                     └──────┬──────┘
                            ↓
                          Queue
                            │
                     ┌──────┴──────┐
                     ↓             ↓
                  Worker 1      Worker 2
                     │             │
             ┌───────┼─────────────┤
             ↓       ↓             ↓
          GitHub    CI/CD         AWS
                                    │
                                    ↓
                              Observability

                     Workflow Database
                            ↑
                            │
                         Workers
```

We're now handling:

- **self-service**
- **automation**
- **standardization**
- **scalability**
- **fault tolerance**
- **retries**
- **idempotency**
- **observability**
- **developer experience**

And importantly, we arrived at this architecture **because we had problems to solve**. We didn't add components just because they're common system-design buzzwords.

---

## Next step: Security

There's a major piece we haven't addressed yet.

The developer clicked:

> **Create Service**

Our platform is now capable of:

- creating GitHub repositories
- creating AWS resources
- configuring CI/CD
- deploying software
- configuring monitoring

That's a **very powerful system**.

So the next question an interviewer could ask is:

> **"How would you handle authentication, authorization, secrets, and permissions?"**

That's our next step, and it's particularly relevant to an internal platform used by **hundreds of engineers**.

Yes, I get what you mean. I'll keep the teaching **focused and progressive**: introduce only the new concept, explain why it matters, and avoid repeating the architecture we've already established unless something changes.

## Step 15 — Security

Our platform can create AWS resources, GitHub repositories, CI/CD pipelines, and deployments. So we need to control **who can do what**.

### 1. Authentication — Who are you?

The developer logs into the internal portal using the company's identity provider.

For example:

```text
Developer
   ↓
Company SSO
   ↓
Internal Platform
```

The platform knows who the developer is.

---

### 2. Authorization — What are you allowed to do?

Being authenticated doesn't mean the developer can do everything.

For example:

| Action | Developer | Platform Admin |
|---|---|---|
| Create service | ✅ | ✅ |
| Deploy service | ✅ | ✅ |
| Delete production infrastructure | ❌ | ✅ |
| Change platform templates | ❌ | ✅ |

This is usually handled with **roles and permissions**.

A good interview phrase:

> "I'd separate authentication from authorization. SSO would authenticate the user, and role-based access control would determine what actions they can perform."

---

### 3. Secrets

Our platform will need credentials for things like:

- AWS
- GitHub
- Datadog
- databases

We **shouldn't store these directly in our application code or configuration**.

Instead, we'd use a secrets manager, such as AWS Secrets Manager or another approved company solution.

The platform retrieves secrets when needed.

---

### 4. One important principle: least privilege

The platform shouldn't have unrestricted access to everything.

For example, if a worker only needs permission to create a repository, it shouldn't have permission to delete every repository in the organization.

So we'd give each component **only the permissions it actually needs**.

That's the **principle of least privilege**.

---

### The key answer I'd give in the interview

> "I'd use company SSO for authentication and role-based access control for authorization. I'd keep credentials in a secrets manager rather than in the application or repositories, and I'd follow least privilege so each service only has the permissions it needs."

That's enough for this part.

### Next problem

Now let's tackle something very relevant to the JD:

> **How do we make sure every service created through the platform follows the company's engineering standards?**

That's where **templates, golden paths, reusable libraries, and CI/CD standards** come together.

## Step 16 — Engineering standards and golden paths

The problem is:

> If 100 teams use our platform, how do we stop every team from setting things up differently?

We don't want:

```text
Team A → different CI/CD
Team B → different logging
Team C → different project structure
Team D → different security checks
```

Instead, the platform provides **standard templates**.

### Service template

A developer chooses:

```text
Create Service
   ↓
.NET API
```

The platform generates:

```text
.NET API template
├── src/
├── tests/
├── Dockerfile
├── CI/CD configuration
├── logging
├── health checks
└── documentation
```

The developer gets a working starting point without manually configuring everything.

---

### CI/CD golden path

The template could also provide a standard pipeline:

```text
Pull Request
     ↓
Build
     ↓
Tests
     ↓
Security checks
     ↓
Build image
     ↓
Deploy
```

Teams can use this by default rather than creating their own pipeline from scratch.

---

### Reusable libraries

We can also provide internal libraries/SDKs for common functionality:

- logging
- authentication
- telemetry
- API clients
- error handling

That gives teams consistency without forcing everyone to reinvent the same code.

---

### One important detail

A **golden path shouldn't mean "you are forced to do everything our way."**

I'd say:

> "I'd make the recommended path the easiest path. Teams can deviate when they have a valid reason, but they shouldn't have to rebuild the standard setup themselves."

That fits the JD's focus on **simplicity, self-service, and developer experience**.

### Next

Now we have the service creation workflow. The next major question is:

> **How does a service actually get deployed from the developer's repository to AWS?**

That's where we'll design the **CI/CD and deployment part**.

Yes. I understand. We're spending too much time decomposing every small concern. For interview preparation, you need to see the **system-level picture and the important trade-offs**, not every implementation detail.

Let's speed this up.

# Step 17 — CI/CD and Deployment

At this point, assume we have:

- Authentication/authorization
- Service templates / golden paths
- Repository creation
- Workflow orchestration
- Infrastructure provisioning

Now the core question is:

> **How does code go from GitHub to a running service?**

A standard pipeline:

```text
Developer
   ↓
GitHub
   ↓
CI/CD
   ↓
Build + Test
   ↓
Docker Image
   ↓
Container Registry
   ↓
Deployment
   ↓
AWS
```

### Why a container registry?

We don't want every deployment to rebuild the application from source.

CI builds an immutable artifact:

```text
PaymentService:v42
```

and stores it in a registry.

The deployment system then takes that exact image and deploys it.

This gives us **repeatable deployments**.

---

## Deployment strategy

For production, we don't necessarily want:

```text
Old version → immediately replace everything
```

We could use:

### Rolling deployment

```text
v41:  v41  v41  v41
          ↓
v42:  v42  v41  v41
          ↓
v42:  v42  v42  v41
          ↓
v42:  v42  v42  v42
```

If something goes wrong, we can stop or roll back.

Other options include **blue/green** and **canary deployments**. You don't need to memorize all three yet; know what problem they solve.

---

## What should the platform standardize?

The platform can provide a default deployment pipeline:

- automated tests
- security checks
- build
- containerization
- deployment
- health checks
- monitoring
- rollback

So the developer experience becomes:

> **Push code → pipeline handles the rest.**

That is the **golden path** we're building.

---

### The complete system is now roughly:

```text
                     Developer
                         │
                         ▼
                Internal Developer Portal
                         │
                         ▼
                  Platform API
                         │
                  Workflow Engine
                         │
          ┌──────────────┼──────────────┐
          ▼              ▼              ▼
       GitHub          CI/CD       Infrastructure
          │              │              │
          └──────────────┼──────────────┘
                         ▼
                    AWS / Runtime
                         │
                         ▼
                   Observability
```

**That's the system I would want you to be able to draw in a Karat interview.**

Then, depending on the interviewer, you drill into the relevant area:

- **Scale** → multiple platform/API/worker instances
- **Reliability** → queues, retries, health checks, rollback
- **Security** → SSO, RBAC, secrets, least privilege
- **Developer experience** → templates, self-service, golden paths
- **Deployment** → CI/CD, containers, canary/rolling deployment
- **Observability** → logs, metrics, traces

We don't need another 15 steps on this particular system.

### I'd move to the next system-design problem now.

For your Immowelt JD, I recommend **Design a scalable REST API** next. It will be much closer to a traditional system-design question and will let us practice the things that are most likely to come up: **load balancing, caching, database scaling, replication, sharding, rate limiting, and observability.**