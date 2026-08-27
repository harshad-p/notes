Prompt:  
Every new or edited listing (\~500k/day, peaks of 100k/hour during agency bulk imports) must pass automated checks: banned phrases, contact info leakage, suspected duplicates, suspicious price. Trusted agencies auto-publish; private sellers always go through automated screening. SLA 5 min; flagged items go to a human moderator queue with a 30-minute SLA for 95% of items. Design the moderation pipeline.

## 1. Clarify requirements

> **"Before I design it, I'd like to clarify a few requirements."**

### Functional requirements

> "My understanding is that every new or edited listing must go through automated moderation checks before it can be published, with different rules for trusted agencies and private sellers."

I'd clarify:

- What exactly does **trusted agency** mean? Is that a property of the agency account that we can retrieve synchronously?
- Does an edited listing go through the **same checks as a new listing**?
- What happens when an automated check flags a listing? Is it always sent to a human moderator?
- Can a moderator approve, reject, or request changes?
- Can a moderator override an automated decision?
- Do we need to notify the seller when a listing is rejected or needs changes?
- Do we need an audit trail of automated decisions and moderator decisions?
- Do we need to retain the original and edited versions of listings?

I'd assume yes to the audit trail, and that moderators can approve/reject/override.

### Non-functional requirements

> "The important requirements I see are the five-minute SLA for automated screening and the 30-minute SLA for 95% of human-review items."

We have:

- ~500k listings/day
- Peak: **100k/hour**
- Automated screening: **≤5 minutes**
- Human moderation: **95% processed within 30 minutes**
- High availability
- Listings shouldn't be lost
- Duplicate processing shouldn't result in inconsistent decisions

100k/hour is about **28 listings/sec**, so the average throughput isn't particularly extreme. The important issue is the **burst** from agency imports.

(I'd explicitly point this out. It demonstrates that I'm not blindly designing for millions of requests/sec.)

---

# 2. Data to store

> "I'd keep the listing itself separate from the moderation state and moderation history."

I'd have something conceptually like:

```text
Listing
- ListingId
- SellerId
- SellerType
- Title
- Description
- Price
- Location
- ...
```

Then moderation state:

```text
Moderation
- ListingId
- Version
- Status
- CreatedAt
- StartedAt
- CompletedAt
- Decision
- RiskScore
```

And an audit/history record:

```text
ModerationResult
- ListingId
- Version
- CheckType
- Result
- Reason
- CreatedAt
```

For example:

> `BannedPhrase → Failed → "..."`

> `PriceCheck → Passed`

> `DuplicateCheck → Flagged`

And for human moderation:

```text
ModerationReview
- ListingId
- Version
- ModeratorId
- Decision
- Reason
- CreatedAt
```

> "I'd include the listing version because the same listing can be edited while an earlier version is still being processed."

That is an important concurrency consideration.

---

# 3. Database

> "For the authoritative listing and moderation state, I'd start with a relational database."

I have clear relationships between:

- listings
- sellers/agencies
- moderation jobs
- moderation results
- moderators
- audit history

I also want strong consistency around the **current moderation state**.

> "I don't need a NoSQL database just because we're processing a lot of listings. 500,000 per day is not inherently a NoSQL requirement."

(I'd rather demonstrate that I'm choosing based on requirements rather than volume alone.)

---

# 4. First version

> "I'll start simple."

```text
Client
   |
   v
API
   |
   v
Database
```

When a listing is created or edited:

> "I would persist the listing and mark it as requiring moderation."

For a private seller:

> `PendingModeration`

For a trusted agency:

> "The requirement says trusted agencies auto-publish, but I still need to clarify whether that means they bypass all automated checks or whether they are published after automated checks."

I'd ask that explicitly.

Assuming **trusted agencies bypass the screening**, their flow is different.

---

# 5. Introduce asynchronous processing

The five-minute SLA tells me that automated screening **doesn't need to happen inside the HTTP request**.

> "I'd make moderation asynchronous."

So:

```text
Client
   |
   v
API
   |
   +----> Database
   |
   +----> Moderation Queue
              |
              v
       Moderation Workers
```

> "The API persists the listing and publishes a moderation job. The worker performs the checks asynchronously."

The API can return quickly instead of making the user wait for every moderation check.

---

# 6. Important detail: database + queue consistency

There's a potential problem.

Suppose:

1. API saves listing.
2. API tries to publish queue message.
3. Queue operation fails.

Now:

```text
Database → listing exists
Queue    → no moderation job
```

The listing could sit there forever.

I wouldn't want that.

> "I'd use an outbox pattern."

The API transaction writes both:

```text
Listing
OutboxEvent
```

into the same database transaction.

Then a background publisher reads the outbox and publishes the moderation event to the queue.

That gives us much stronger reliability between the database and messaging system.

---

# 7. Moderation pipeline

Now I have several checks:

- Banned phrases
- Contact information leakage
- Suspected duplicate
- Suspicious price

I could run these sequentially, but I'd prefer to run independent checks concurrently.

> "Once the worker has the listing, the independent checks can execute in parallel."

Conceptually:

```text
                   → Banned phrase check
                  /
Moderation Job → → Contact info check
                  \
                   → Duplicate check
                    \
                     → Price check
```

> "This reduces the overall processing time because we're not waiting for one check to finish before starting the next."

(I'd only do this if the checks are genuinely independent and the downstream systems can handle the concurrency.)

---

# 8. What does the worker actually do?

> "The worker retrieves the listing version and executes the checks."

Suppose we get:

```text
Banned phrases     → PASS
Contact leakage    → PASS
Duplicate          → FLAG
Suspicious price   → PASS
```

Then:

> "The overall automated decision becomes `FLAGGED`, because at least one check requires human review."

If everything passes:

> `APPROVED`

If a rule is definitive, such as a clearly prohibited phrase:

> `REJECTED`

I'd define these states explicitly rather than letting every check invent its own interpretation.

---

# 9. Idempotency

Now I need to consider duplicate queue messages.

Suppose the same moderation job is delivered twice.

Without protection:

```text
Job 123
 ↓
Worker 1 → moderation
Job 123
 ↓
Worker 2 → moderation again
```

> "I'd make the moderation operation idempotent using the listing ID plus version, or a unique moderation-job ID."

So:

```text
ListingId = 123
Version = 7
```

identifies one particular moderation operation.

If a worker receives it twice, the second worker can detect that the work has already been completed or is currently being processed.

---

# 10. What happens if a worker crashes?

Suppose:

```text
Worker
 ↓
runs checks
 ↓
crashes
```

> "The queue should not permanently lose the message. The message becomes available for redelivery according to the queue's acknowledgement semantics."

I'd also make the processing idempotent, because redelivery is expected in distributed systems.

---

# 11. Retry policy

For temporary failures:

> "I'd retry with exponential backoff and a maximum retry count."

For example, if the duplicate-detection service is temporarily unavailable:

```text
Attempt 1 → failure
Attempt 2 → failure
Attempt 3 → failure
...
Maximum retries
       ↓
Dead-letter queue
```

But I wouldn't retry a deterministic business result.

For example:

> "Banned phrase detected"

isn't a transient error.

There's no point retrying it.

---

# 12. What happens after automated moderation?

There are three possible outcomes.

### Approved

> "I mark the listing as approved and make it publishable."

### Rejected

> "I mark it as rejected and notify the seller if that's part of the product requirements."

### Flagged

> "I create a human moderation task."

So:

```text
Automated Moderation
        |
        +---- Approved → Publish
        |
        +---- Rejected → Reject
        |
        +---- Flagged → Moderator Queue
```

---

# 13. Human moderator queue

I'd create a separate queue for human moderation.

> "I wouldn't use the same queue as the automated workers because the processing characteristics are completely different."

Automated moderation can scale horizontally.

Human moderators can't simply be added automatically whenever traffic increases.

So:

```text
Automated Moderation
        |
        v
Human Moderation Queue
        |
        v
Moderator UI
```

---

# 14. How do we meet the 30-minute SLA?

I'd prioritize the moderation queue.

Every moderation task gets:

```text
CreatedAt
DueAt
Priority
```

The moderation UI can show:

> Oldest/highest-priority items first.

I'd monitor:

- queue depth
- oldest item age
- time waiting for moderator
- percentage completed within 30 minutes

The key metric is:

> **95th percentile moderation completion time ≤ 30 minutes.**

---

# 15. What if there are too many flagged listings?

This is where the SLA becomes interesting.

Suppose a bulk import suddenly generates:

```text
100,000 listings/hour
```

and 20% get flagged.

That's:

```text
20,000 human-review items/hour
```

We may not have enough moderators.

> "I'd monitor the queue's arrival rate versus the moderator processing rate."

If the queue is growing:

> "We have a capacity problem, not a database problem."

I'd consider:

- prioritization
- increasing moderator capacity
- improving automated classification
- automatically resolving low-risk flags
- different SLAs by severity

But I wouldn't silently violate the 30-minute SLA.

---

# 16. Duplicate detection

This is one of the more interesting checks.

> "I wouldn't necessarily query the primary SQL database for every possible duplicate."

Depending on how sophisticated duplicate detection needs to be, I'd consider a **search/indexing system**.

For example, we could index characteristics such as:

- location
- price
- property type
- rooms
- normalized title
- normalized description
- potentially image fingerprints

Then:

> "The duplicate detector searches for listings that look sufficiently similar and returns candidates."

The application can then determine whether the similarity is sufficient to flag the listing.

(I'd introduce this only because duplicate detection is specifically a search/similarity problem. I wouldn't introduce Elasticsearch merely because we're building a moderation system.)

---

# 17. Suspicious price

This could start simple.

For example:

> "Compare the listing price against historical prices for similar properties."

But as the system evolves, this could become a machine-learning model.

I'd keep the interface abstract:

> `PriceRiskService → risk score`

For example:

```text
Risk score = 0.02 → normal
Risk score = 0.91 → suspicious
```

Then the moderation policy determines what to do with the score.

---

# 18. Contact information leakage

This can be a deterministic rule-based service initially.

For example:

- phone numbers
- email addresses
- URLs
- prohibited contact patterns

> "I'd normalize the text first, because users may try to bypass simple pattern matching by inserting spaces or punctuation."

Again, the important system-design point is that this is an independent moderation check and can run in parallel with the others.

---

# 19. Scaling

Our peak is around **28 listings/sec**, but we need to tolerate bursts.

The queue gives us buffering.

> "I'd horizontally scale the moderation workers based on queue depth and processing latency."

For example:

```text
Moderation Queue
      |
      +---- Worker 1
      +---- Worker 2
      +---- Worker 3
      +---- Worker 4
```

If imports cause the queue to grow:

> "Add more workers."

Once the burst finishes:

> "Scale them back down."

---

# 20. What about the five-minute SLA?

I'd monitor the age of the oldest moderation job.

For example:

> `Moderation queue oldest message age`

If it approaches the SLA threshold, we need to increase capacity.

I might also reserve some worker capacity specifically for moderation so that unrelated background workloads can't consume all the workers.

---

# 21. What if the moderation service itself becomes overloaded?

> "The queue provides backpressure."

The listing API doesn't have to process all moderation immediately.

The queue absorbs the burst.

But the queue isn't magic:

> "If the arrival rate remains higher than our processing rate for long enough, the queue will keep growing."

So I'd scale workers and monitor the backlog.

---

# 22. Trusted agencies

For trusted agencies, I'd have:

```text
Agency
   |
Trusted?
   |
   +---- Yes → Publish
   |
   +---- No  → Moderation Pipeline
```

But I would clarify whether trusted agencies completely bypass moderation or simply have a different moderation policy.

> "If they truly auto-publish, I would still record that the listing was published under the trusted-agency policy."

That gives us an audit trail.

---

# 23. Edited listings

This is important.

Suppose listing version 5 was approved.

The seller edits the description.

Now we have:

```text
Version 5 → APPROVED
Version 6 → PENDING
```

> "I would never let the edit inherit the previous approval automatically."

The edited version needs to go through the required checks.

And the published state should remain associated with the correct version.

Otherwise, you could get:

> approved version 5 → seller edits prohibited content → version 6 accidentally remains published.

---

# 24. Concurrency problem

Suppose version 6 is being moderated while the seller edits again:

```text
Version 6 → moderation running
Version 7 → submitted
```

Now version 6 could finish after version 7.

> "I'd include the listing version in every moderation job and verify that the result still applies to the current version before changing the listing's publish state."

So version 6 cannot accidentally overwrite the state of version 7.

This is one of those details I'd expect a senior engineer to catch.

---

# 25. Reliability

I'd ask:

> "What happens if the database is temporarily unavailable?"

The moderation job remains in the queue and can be retried.

> "What happens if the duplicate-detection service is unavailable?"

Retry the individual check rather than losing the entire moderation job.

> "What happens if the email notification service is unavailable?"

The moderation decision should not be lost just because notification delivery failed. Notification delivery should be handled asynchronously with its own retry policy.

---

# 26. Observability

I'd track:

### Automated moderation

- listings processed/sec
- processing latency
- percentage completed within 5 minutes
- queue depth
- oldest queue item
- failure rate
- retry count
- dead-letter count
- percentage approved/flagged/rejected

### Human moderation

- queue depth
- oldest pending review
- median review time
- 95th percentile review time
- percentage meeting 30-minute SLA
- moderator throughput

### Individual checks

- banned phrase latency
- duplicate detection latency
- price check latency
- contact-information check latency
- failure rates for each

> "I'd particularly monitor each check separately because one slow or failing check shouldn't be hidden inside the overall moderation latency."

---

# 27. Final architecture

At this point I'd draw:

```text
                         ┌──────────────────────┐
                         │   Listing Database   │
                         └──────────▲───────────┘
                                    │
Client → Load Balancer → Listing API
                                    │
                                    ↓
                              Outbox Event
                                    │
                                    ↓
                           Moderation Queue
                                    │
                                    ↓
                         Moderation Workers
                           /      |      \
                          /       |       \
                         ↓        ↓        ↓
                  Phrase Check  Contact  Duplicate
                                      \
                                       → Price Check
                                    │
                                    ↓
                           Moderation Decision
                            /       |        \
                           /        |         \
                          ↓         ↓          ↓
                      Approved   Rejected    Flagged
                          |         |           |
                          ↓         ↓           ↓
                       Publish   Notify    Human Queue
                                               |
                                               ↓
                                        Moderator UI
                                               |
                                      Approve / Reject
```

With retries and dead-letter queues around the asynchronous processing paths.

---

# 28. Final answer I'd give the interviewer

> **"I'd start with the listing service and relational database as the source of truth. Every new or edited listing gets a version and a moderation state. For private sellers, the listing enters the moderation pipeline; trusted agencies follow the trusted-agency policy."**
>
> **"I wouldn't run the moderation synchronously inside the API request. I'd persist the listing and create an asynchronous moderation job, using an outbox pattern so that we don't end up with a listing in the database but no corresponding moderation event."**
>
> **"Moderation workers consume those jobs and run the independent checks in parallel: banned phrases, contact information leakage, duplicate detection, and suspicious pricing. The results are combined into an overall decision. Approved listings can be published, definitive failures can be rejected, and suspicious cases go into a separate human-moderation queue."**
>
> **"I'd make the processing idempotent using the listing ID and version because queue messages can be delivered more than once. Retries would use exponential backoff, and messages that repeatedly fail would go to a dead-letter queue."**
>
> **"For the five-minute automated SLA, I'd monitor queue age and processing latency and horizontally scale the moderation workers during bulk-import spikes. For the human SLA, I'd prioritize the moderator queue based on age and severity and monitor the 95th percentile review time."**
>
> **"I'd also be careful with edits. If an approved listing is edited, the new version must be screened independently, and an older moderation result must never overwrite the state of a newer version."**
>
> **"For duplicate detection, I could eventually introduce a search index or similarity service because that problem is fundamentally about finding similar listings. I wouldn't add it to the initial architecture unless the requirements justified it."**
>
> **"Finally, I'd monitor the entire pipeline: queue depth, processing latency, SLA compliance, failure and retry rates, dead-letter messages, and the performance of each individual moderation check."**

(The main design decisions here are driven by **asynchronous processing, burst handling, independent checks, reliability, idempotency, and the two SLAs**. I wouldn't add Kafka, Redis, Kubernetes, or a search engine simply to make the diagram look more distributed.)