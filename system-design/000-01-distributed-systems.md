
### 1. Vertical scaling ≠ distributed

You can have one powerful server:

```text
             Users
               ↓
       ┌───────────────┐
       │    Server     │
       │  32 CPU       │
       │  128 GB RAM   │
       └───────────────┘
```

If it becomes overloaded, you can make it more powerful:

> 32 CPU → 64 CPU → 128 CPU

That's **vertical scaling (scale up)**.

It's still a single machine.

---

### 2. Horizontal scaling = multiple servers

Instead:

```text
                Users
                  ↓
            Load Balancer
             /    |    \
            ↓     ↓     ↓
         Server Server Server
```

Now we're distributing the workload across multiple machines.

That's **horizontal scaling (scale out)**.

And yes, this is one of the major reasons we distribute systems: **scalability**.

---

### 3. Fault tolerance / availability

Here's another reason.

With one server:

```text
Server 💥
   ↓
Everything is down
```

With multiple servers:

```text
             Load Balancer
              /    |    \
             ↓     ↓     ↓
           Server Server Server
              💥
```

One fails → the others continue serving requests.

That's **fault tolerance / high availability**.

So your instinct is correct: distribution can improve reliability.

---

### 4. Replication is another concept

Suppose you have:

```text
        Primary DB
        /        \
       ↓          ↓
   Replica 1   Replica 2
```

Now you have **copies of the data**.

That's **replication**.

It can help with:

- availability
- disaster recovery
- read scalability

But replication itself isn't the same thing as "distributed."

---

### 5. Sharding is another form of distribution

Instead of copying the same data:

```text
DB 1 → Users A–M
DB 2 → Users N–Z
```

Now the data is **partitioned** across machines.

That's **sharding**.

It can help with:

- scaling storage
- scaling throughput

---

## The easiest way to remember this

Think of these as separate concepts:

| Concept | Main idea |
|---|---|
| **Vertical scaling** | Make one machine more powerful |
| **Horizontal scaling** | Add more machines |
| **Distribution** | Spread computation/data across machines |
| **Replication** | Keep copies of data |
| **Sharding** | Split data between machines |
| **Fault tolerance** | Continue working when something fails |
| **High availability** | Keep the service available |
| **Scalability** | Handle increasing load |

And they often work **together**.

So when you're asked in a system-design interview:

> **"Why do we need multiple servers?"**

Don't automatically answer **"for scalability."**

A stronger answer is:

> "There can be several reasons. We may need more capacity to scale horizontally, but multiple instances also improve availability because we can continue serving traffic if one instance fails."

# what does "distributed database" actually mean?

This is a very important distinction.

There are two concepts people often mix up:

### Replication

You have **copies of the same data**.

```text
             Primary
                │
        ┌───────┴────────┐
        ▼                ▼
      EU DB             US DB
      copy              copy
```

### Sharding

You have **different pieces of the data**.

```text
EU DB → URLs 1–50 billion
US DB → URLs 50–100 billion
```

Those are completely different architectures.
