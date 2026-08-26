Exactly. Breaking things down is useful; **repeating the same concept from slightly different angles isn't**. I'll keep the depth, but make each step earn its place.

And since the role is platform-heavy, I'd actually make the next problem a little more relevant than a generic REST API:

## System Design #3 — Design a scalable notification system

This is a good exercise because it introduces several classic distributed-system concepts without being tied to REST:

> **Design a system that sends notifications to users through email, SMS, and push notifications.**

We'll encounter:

- queues and asynchronous processing
- high traffic / burst handling
- retries and dead-letter queues
- idempotency
- rate limiting
- multiple workers
- reliability
- observability
- prioritization

And importantly, we'll **start with a simple design and scale it only when the requirements demand it**.

That is much closer to the way you should approach the Karat interview than trying to demonstrate every technology you know.

### Step 1 — Requirements

Let's assume:

- Other applications send us notification requests.
- We support **email, SMS, and push**.
- Notifications don't need to be delivered immediately; a delay of a few seconds is acceptable.
- We need to handle traffic spikes.
- We should not accidentally send the same notification twice.
- We need to know whether delivery succeeded or failed.

That's enough to start.

The first question I'd ask the interviewer is:

> **"Are notifications expected to be delivered synchronously, or is asynchronous delivery acceptable?"**

Given the requirement above, we'd choose **asynchronous**.

So our first architecture is simply:

```text
Application
     ↓
Notification Service
     ↓
   Queue
     ↓
  Workers
   / | \
  ↓  ↓  ↓
Email SMS Push
```

That's our starting point.

**Next, we'll look at why the queue is the key component and what happens when 1 million notifications arrive during a traffic spike.**

Sure. This distinction is important because it affects the whole architecture.

## Synchronous vs. asynchronous notification delivery

Suppose an application wants to send an email:

```text
Application → Notification Service → Email Provider
```

### Synchronous

The application **waits** for the notification service to finish.

```text
Application
    │
    │ "Send this email"
    ▼
Notification Service
    │
    │ send email
    ▼
Email Provider
    │
    │ success
    ▼
Notification Service
    │
    │ success
    ▼
Application
```

The application might receive:

```text
200 OK
```

only after the email provider has responded.

**Problem:** If the email provider takes 3 seconds or is unavailable, the application's request is also delayed or fails.

---

### Asynchronous

The application only asks the notification service to **accept the request**.

```text
Application
    │
    │ "Send this email"
    ▼
Notification Service
    │
    │ put message in queue
    ▼
   Queue
    │
    ▼
   202 Accepted
    │
    ▼
Application continues
```

Later:

```text
Queue
  ↓
Worker
  ↓
Email Provider
```

The application doesn't wait for the email to actually be delivered.

This is usually the better choice for notifications because **sending the notification isn't normally part of the application's critical request path**.

For example, after a user places an order:

```text
Place Order
    ↓
Order successful
    ↓
Return response to user
```

We don't want the user waiting because:

> "We're still trying to send your confirmation email."

Instead:

```text
Order Service
     ↓
Notification Queue
     ↓
Email Worker
     ↓
Email Provider
```

The email can be delivered a few seconds later.

---

# Step 2 — Why the queue matters

Now imagine we suddenly receive:

**1 million notifications in one minute.**

If we synchronously call the email/SMS/push providers:

```text
1,000,000 requests
        ↓
Notification Service
        ↓
External providers
```

we could overwhelm our service or the providers.

With a queue:

```text
1,000,000 requests
        ↓
Notification Service
        ↓
      Queue
        ↓
  Workers process
  at controlled rate
```

The queue acts as a **buffer**.

We can scale workers horizontally:

```text
             Queue
          /    |    \
         ↓     ↓     ↓
      Worker Worker Worker
         ↓     ↓     ↓
       Providers
```

If the incoming rate is higher than the processing rate, messages wait in the queue rather than overwhelming the workers.

---

### There's an important trade-off

Asynchronous processing gives us **better scalability and resilience**, but we lose immediate confirmation of delivery.

The application can know:

> "The notification was accepted."

It cannot immediately know:

> "The user received the email."

Those are two different states.

That's an important distinction to mention in a system-design interview.

---

### Our design so far

```text
                    ┌──────────────┐
                    │ Applications │
                    └──────┬───────┘
                           ↓
                  ┌──────────────────┐
                  │ Notification API │
                  └────────┬─────────┘
                           ↓
                        Queue
                           ↓
                    ┌──────┴──────┐
                    ↓             ↓
                 Workers        Workers
                    ↓             ↓
                 Email          SMS/Push
                 Provider        Providers
```

**Next, we should tackle a more interesting problem: what happens when an email provider is temporarily down?**

That's where retries, backoff, dead-letter queues, and idempotency become useful—and we'll keep those concepts tied to the actual problem rather than expanding them unnecessarily.

## Step 3 — What if the email provider is down?

Suppose our worker picks up a notification:

```text id="q8zj3a"
Queue
  ↓
Worker
  ↓
Email Provider ❌
```

We don't want to immediately mark the notification as permanently failed. The provider might be down for only a few seconds.

### Retry with backoff

We can retry:

```text id="8q4j5x"
Attempt 1 → fails
     ↓
wait 1 sec
     ↓
Attempt 2 → fails
     ↓
wait 5 sec
     ↓
Attempt 3 → fails
     ↓
wait 30 sec
     ↓
Attempt 4 → succeeds
```

The increasing delay is **exponential backoff**.

We can also add some randomness (**jitter**) so thousands of workers don't all retry at exactly the same time.

---

## What if it keeps failing?

We shouldn't retry forever.

After, say, 5 attempts:

```text id="qj4m7s"
Worker
  ↓
5 failed attempts
  ↓
Dead Letter Queue (DLQ)
```

The notification is now considered failed and can be investigated or retried later.

This is better than silently losing it.

---

## But there's another problem

Suppose the provider actually **received the email**, but our worker didn't receive the response because of a network failure.

The worker thinks:

```text id="x9a2fr"
"Email failed."
```

and retries.

The provider might then send the same email twice.

So we need **idempotency**.

For each notification, we give it a unique ID:

```text id="v2j3qk"
Notification ID: 847291

"Your order has shipped."
```

The provider or our system can use that ID to recognize:

> "I've already processed notification 847291."

So a retry doesn't create a duplicate.

---

### Our system now looks like:

```text id="z6y8q0"
                    Application
                         ↓
                Notification API
                         ↓
                       Queue
                         ↓
                      Worker
                         ↓
                  Email Provider
                    ↙          ↘
                Success       Failure
                               ↓
                            Retry
                               ↓
                         Max attempts?
                          /          \
                        No            Yes
                        ↓              ↓
                     Retry           DLQ
```

The important concepts here are:

**queue → retry → exponential backoff → DLQ → idempotency**

These are worth knowing for your interview because they're common patterns in distributed systems.

### Next step

Let's make the system more realistic:

> **What happens when we have 10 million notifications and need to support email, SMS, and push, each with different throughput limits?**

That's where we'll deal with **multiple queues, worker scaling, and rate limiting**.

I get your idea of idempotency, but I don't see how your approach of a notification ID solves it. For it to work, the provider needs to be setup to handle duplicates. But that is beyond our control, right? Or the providers these days can be setup for such a situation? 
ANyway, you can explain this along with the next step.

You're right. **A notification ID by itself does not make the operation idempotent.** I oversimplified that part.

If we're calling an external email provider, we can't assume it will understand our `notificationId`.

There are actually two levels of idempotency here.

### 1. Our system can be idempotent

Suppose our worker receives the same notification twice:

```text
Queue
 ↓
Worker
 ↓
Notification DB
```

We can have:

```text
notificationId = 847291
status = SENT
```

If the worker sees `847291` again, it can avoid calling the provider again.

That protects **our own system** from duplicate processing.

### 2. But there's a nasty failure case

Consider:

```text
Worker → Email Provider
             ↓
         Email sent ✓
             ↓
      Network failure
             ↓
Worker doesn't know
```

Our worker thinks the operation failed.

It retries:

```text
Worker → Email Provider
             ↓
         Email sent again ❌
```

Now the user gets two emails.

**We cannot completely solve this from our side if the provider doesn't support idempotency.**

That's an important distributed-systems limitation.

---

## What can we do?

Many modern APIs **do support idempotency keys** for operations where duplicate execution matters. For example, an API might let us send:

```text
Idempotency-Key: 847291
```

and guarantee that repeating the request with that key won't create another operation.

But **we must check the specific provider's API**. We shouldn't design our system assuming every email/SMS provider supports this.

If the provider doesn't support it, we can reduce duplicates using our own state tracking, but we can't guarantee exactly-once delivery.

And this leads to an important interview concept:

> **Exactly-once processing is difficult in distributed systems. We usually design for at-least-once processing and make the operation idempotent where possible.**

That's the more accurate way to think about it.

---

# Next step — Multiple channels and rate limits

Now suppose we have:

```text
10 million notifications
```

and three channels:

```text
              Notification Queue
                 /      |      \
                ↓       ↓       ↓
             Email     SMS     Push
             Queue     Queue    Queue
                ↓       ↓       ↓
             Workers  Workers  Workers
```

Why separate queues?

Because the providers have different characteristics.

For example:

```text
Email → can process 10,000/sec
SMS   → limited to 1,000/sec
Push  → can process 20,000/sec
```

If everything shared one queue, slow SMS processing could affect email and push.

With separate queues, we can scale independently:

```text
Email Queue → 20 workers
SMS Queue   → 5 workers
Push Queue  → 30 workers
```

---

## Rate limiting

There's another issue.

Suppose the SMS provider allows:

> **1,000 requests per second**

but our queue contains 100,000 SMS notifications.

We can't simply start 100 workers and let them send as fast as possible.

We need a **rate limiter**:

```text
SMS Queue
    ↓
Rate Limiter
    ↓
≤ 1,000 requests/sec
    ↓
SMS Workers
    ↓
Provider
```

The queue absorbs the spike, while the rate limiter protects the provider and keeps us within its limits.

---

### So the architecture is now:

```text
                       Applications
                            ↓
                    Notification API
                            ↓
                      ┌───────────┐
                      │ Dispatcher│
                      └─────┬─────┘
                            │
              ┌─────────────┼─────────────┐
              ↓             ↓             ↓
         Email Queue    SMS Queue     Push Queue
              ↓             ↓             ↓
         Rate Limiter  Rate Limiter  Rate Limiter
              ↓             ↓             ↓
           Workers       Workers       Workers
              ↓             ↓             ↓
           Provider      Provider      Provider
```

The important new idea is **isolation**: one slow or overloaded notification channel shouldn't bring down the others.

Next, I'd move to **observability and delivery status**. That's the last major piece I'd add before we call this system complete.

## Step 4 — Observability and delivery status

Now imagine the system is running at scale. We need to answer questions like:

- How many notifications are being sent?
- How many are failing?
- Is SMS slower than usual?
- Is a provider down?
- How long are messages sitting in the queue?

This is where **observability** comes in.

### Three things to monitor

**Metrics**

```text
notifications_sent
notifications_failed
queue_depth
delivery_latency
provider_error_rate
```

These give us the numbers.

**Logs**

For an individual notification:

```text
Notification 847291
Channel: SMS
Status: FAILED
Provider response: timeout
Attempt: 3
```

Useful for troubleshooting.

**Traces**

If a notification passes through several services:

```text
Application
   ↓
Notification API
   ↓
Queue
   ↓
Worker
   ↓
SMS Provider
```

A distributed trace lets us follow that request across the system.

---

## Delivery status

I'd also maintain a notification record:

```text
ID: 847291
Channel: SMS
Status: SENT
Created: 10:01:02
Sent: 10:01:04
Attempts: 1
```

Possible states could be:

```text
QUEUED
PROCESSING
SENT
FAILED
```

This lets the application or support team answer:

> "What happened to this notification?"

One subtle point: **`SENT` doesn't necessarily mean the user saw the notification.** It generally means our provider accepted/sent it. If the provider gives us delivery callbacks/webhooks, we can track a more precise status such as `DELIVERED`.

---

# Final architecture

At this point, I'd be comfortable stopping:

```text
                         Applications
                              │
                              ▼
                      Notification API
                              │
                         Dispatcher
                    ┌─────────┼─────────┐
                    ▼         ▼         ▼
                Email Q     SMS Q     Push Q
                    │         │         │
                Rate Limit Rate Limit Rate Limit
                    │         │         │
                 Workers   Workers    Workers
                    │         │         │
                    ▼         ▼         ▼
                Providers  Providers  Providers
                    │         │         │
                    └─────────┼─────────┘
                              ▼
                        Delivery Status
                              │
                        Observability
```

### If the interviewer asks how you'd scale it

You can say:

> "The API and workers can scale horizontally. Queues absorb traffic spikes, and I'd scale each channel independently because the providers have different throughput limits. I'd also use rate limiting to protect external providers."

### If they ask about reliability

> "I'd use durable queues, retries with exponential backoff, dead-letter queues for persistent failures, and idempotency where the provider supports it. I wouldn't claim exactly-once delivery because that's difficult to guarantee with external systems."

That's a **solid senior-level system design answer** without needing to over-engineer it.

For your Karat preparation, I'd consider this problem **done**. The next one I'd recommend is **a file storage system like Google Drive/Dropbox**, because it introduces a completely different set of problems: large files, object storage, metadata, uploads, downloads, consistency, and access control.