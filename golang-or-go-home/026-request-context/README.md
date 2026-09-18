# Lesson 34 — Request Context: Cancellation, Deadlines & Request-Scoped Data

We’ll cover three things:

1. **What `context.Context` is**
2. **Cancellation and deadlines**
3. **Request-scoped values**

The important idea is that a request is not just data coming into your handler. It also has a **lifecycle**: it can be cancelled, it can have a deadline, and it can carry a small amount of request-specific metadata.

---

## 1. What is `context.Context`?

Go's `context` package provides a standard way to carry **request lifetime and request-scoped information** through a call chain.

You can think of a context as an object representing:

> "This operation is part of this request, and these are the conditions under which the operation should continue."

For an HTTP request, you get its context with:

```go
ctx := r.Context()
```

`r.Context()` returns a `context.Context`.

You don't normally create a new context for every HTTP request yourself. The HTTP server creates one and associates it with the request.

So your call chain might look like:

```text
HTTP request
    │
    ▼
handler
    │
    ▼
service
    │
    ▼
repository
    │
    ▼
database
```

The same context can travel through that entire chain:

```text
request context
      │
      ├── handler
      ├── service
      ├── repository
      └── database
```

This becomes particularly important when something causes the request to stop.

---

# 2. Cancellation

Suppose your endpoint does something expensive:

```text
GET /reports/123
       │
       ▼
    Handler
       │
       ▼
   Database query
       │
       ▼
   External API
```

Now imagine the client closes the browser before the operation finishes.

There is little value in continuing expensive work for a response nobody will receive.

The request's context can become **cancelled**.

The cancellation can propagate down the call chain.

That's one of the major purposes of `context.Context`.

---

## 3. `ctx.Done()`

A context exposes:

```go
ctx.Done()
```

This gives you a notification mechanism that tells you:

> "The context has been cancelled."

`Done()` returns a **channel**.

We haven't covered Go channels yet, so don't worry about their mechanics right now. For this lesson, the important thing is:

```text
ctx.Done()
    ↓
notification that cancellation occurred
```

You can wait for that notification using:

```go
<-ctx.Done()
```

The `<-` here means "receive from the channel."

For example:

```go
select {
case <-ctx.Done():
    return ctx.Err()
}
```

We'll cover `select` and channels properly when we get to Go concurrency.

For now, understand the API-level behavior rather than the channel implementation.

---

# 4. `ctx.Err()`

When a context has been cancelled or its deadline has expired, you can ask why:

```go
ctx.Err()
```

It returns an error.

Common possibilities are:

```go
context.Canceled
```

or:

```go
context.DeadlineExceeded
```

The distinction matters:

```text
context.Canceled
        ↓
operation was cancelled

context.DeadlineExceeded
        ↓
operation exceeded its allowed time
```

---

# 5. Deadlines and timeouts

Cancellation is useful, but sometimes **we want to impose our own limit**.

For example:

> "This database operation must finish within 2 seconds."

The context package lets us derive a context with a timeout:

```go
ctx, cancel := context.WithTimeout(parent, 2*time.Second)
defer cancel()
```

This creates a **child context** from `parent`.

Conceptually:

```text
parent context
      │
      ▼
WithTimeout(..., 2 seconds)
      │
      ▼
child context
```

The child context automatically becomes cancelled when the two seconds expire.

### Why `defer cancel()`?

`WithTimeout` creates resources associated with the timer.

Calling:

```go
cancel()
```

releases them when you're finished early.

So the standard pattern is:

```go
ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
defer cancel()
```

Notice something important here:

The variable on the right-hand side is itself a context:

```go
ctx, cancel := context.WithTimeout(ctx, ...)
```

The new `ctx` replaces the old context for the operation we're about to perform.

---

# 6. `WithTimeout` vs `WithDeadline`

There are two closely related APIs.

### `WithTimeout`

You specify a duration:

```go
context.WithTimeout(ctx, 2*time.Second)
```

Meaning:

> Give this operation at most 2 seconds from now.

### `WithDeadline`

You specify an absolute point in time:

```go
context.WithDeadline(ctx, deadline)
```

Meaning:

> This operation must finish by this particular time.

So:

```text
WithTimeout
"2 seconds from now"

WithDeadline
"by 11:30:00"
```

For most request-level time limits, `WithTimeout` is convenient.

---

# 7. Why pass context into repositories?

This is where context becomes really useful in backend code.

Suppose your repository has:

```go
func (r *UserRepository) GetByID(ctx context.Context, id int) (User, error)
```

The handler obtains:

```go
ctx := r.Context()
```

and passes it down:

```text
HTTP request
     │
     ▼
   ctx
     │
     ├──────────────► service
     │                    │
     │                    ▼
     │               repository
     │                    │
     │                    ▼
     └──────────────► database
```

The database operation can then respect cancellation/deadlines.

For example, database libraries commonly provide APIs such as:

```go
db.QueryContext(ctx, query)
```

The important part isn't the method name itself.

It's the design:

> **The operation receives the request's context, so it knows when it should stop.**

Without context propagation, you could have:

```text
Client disconnected
       │
       ▼
HTTP handler no longer useful
       │
       X
       │
       ▼
Database query continues anyway
```

With proper context propagation:

```text
Client disconnected
       │
       ▼
request context cancelled
       │
       ▼
database operation notices
       │
       ▼
operation stops
```

This is one reason you'll see `context.Context` almost everywhere in serious Go backend APIs.

---

# 8. Context should generally be passed explicitly

You'll often see:

```go
func (s *UserService) GetUser(
    ctx context.Context,
    id int,
) (User, error)
```

rather than storing the context inside the service:

```go
type UserService struct {
    ctx context.Context // generally a bad design
}
```

Context represents the **lifetime of a particular operation/request**, not the lifetime of your service object.

A service might handle thousands of requests:

```text
Service
 ├── request A → context A
 ├── request B → context B
 ├── request C → context C
 └── request D → context D
```

Therefore the context belongs to the operation:

```go
GetUser(ctx, id)
```

not to the long-lived service.

---

# 9. Request-scoped values

Context has another capability: carrying values associated with a request.

For example, middleware might determine a request ID:

```text
HTTP request
     │
     ▼
logging middleware
     │
     │ request ID = abc-123
     ▼
handler
     │
     ▼
service
```

The request ID can be stored in the context.

The API is:

```go
context.WithValue(...)
```

Conceptually:

```text
parent context
      │
      ▼
WithValue(...)
      │
      ▼
child context
      │
      └── request ID = abc-123
```

Then downstream code can retrieve it.

---

# 10. But don't turn Context into a parameter bag

This is an important Go convention.

Don't do this:

```go
ctx = context.WithValue(ctx, "userName", "Harshad")
ctx = context.WithValue(ctx, "age", 37)
ctx = context.WithValue(ctx, "country", "Germany")
ctx = context.WithValue(ctx, "customer", customer)
```

and then use context as a replacement for ordinary function parameters.

That's difficult to understand because the function signature no longer tells you what the function actually needs.

Prefer:

```go
func ProcessOrder(ctx context.Context, customer Customer) error
```

over hiding `customer` inside the context.

### Good candidates

Context values are most appropriate for **request-scoped metadata** that needs to cross many layers, such as:

- request/correlation ID
- trace information
- authentication-related metadata in appropriate designs
- logging metadata

The general rule is:

> **Context carries request lifetime and request-scoped metadata, not normal application data.**

---

# 11. One practical middleware example

This connects directly to the middleware lesson.

Suppose we want every request to have a timeout of 5 seconds.

We can create middleware:

```go
func timeoutMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx, cancel := context.WithTimeout(
            r.Context(),
            5*time.Second,
        )
        defer cancel()

        r = r.WithContext(ctx)

        next.ServeHTTP(w, r)
    })
}
```

There are two new things here.

### `r.WithContext(ctx)`

An HTTP request contains its context.

`WithContext` returns a **new request associated with the supplied context**.

So:

```go
r = r.WithContext(ctx)
```

means:

> From this point onward, downstream handlers should see this new context.

The flow becomes:

```text
incoming request
       │
       ▼
r.Context()
       │
       ▼
WithTimeout(..., 5 sec)
       │
       ▼
new context
       │
       ▼
r.WithContext(ctx)
       │
       ▼
next handler
       │
       ▼
service
       │
       ▼
repository
```

Now every downstream component that uses `r.Context()` receives the timeout-aware context.

---

# 12. One subtle but important point

A timeout on the context **does not magically terminate your Go function**.

For example, if you write:

```go
ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
defer cancel()

doSomethingExpensive()
```

`doSomethingExpensive()` will not automatically stop merely because two seconds passed.

The operation must **cooperate with the context**.

For example, a database API that accepts `ctx` can monitor it:

```go
db.QueryContext(ctx, ...)
```

Or your own long-running code can periodically check the context.

So context is fundamentally a **cooperative cancellation mechanism**.

That's an important distinction.

---

# 13. The mental model

Keep this model:

```text
                  HTTP request
                       │
                       ▼
                 request.Context()
                       │
          ┌────────────┴────────────┐
          │                         │
      cancellation               timeout
          │                         │
          └────────────┬────────────┘
                       ▼
                 service calls
                       │
                       ▼
                 DB / external API
```

And separately:

```text
Context
  │
  ├── lifetime / cancellation
  ├── deadline
  └── request-scoped metadata
```

That's the core of `context.Context`.

---

## Exercise

Add a **5-second timeout middleware** to the API you've been building.

Then create an endpoint whose handler performs some deliberately slow operation.

Your goal is to observe the difference between:

- an operation that **ignores** the context
- an operation that **checks/uses** the context

Also inspect `ctx.Err()` after cancellation and determine whether you get `context.Canceled` or `context.DeadlineExceeded`.

Don't worry about channels or `select` yet—we'll properly cover those when we reach Go concurrency.