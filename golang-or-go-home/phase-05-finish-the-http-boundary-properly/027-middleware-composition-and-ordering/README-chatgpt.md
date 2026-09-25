# Lesson 34 — Middleware Composition, Ordering, and Production Patterns

## 1. Middleware composition: moving beyond nested calls

Previously we built middleware like:

```go
handler := loggerMiddleware(
    versionMiddleware(
        userHandler,
    ),
)
```

This works because middleware is fundamentally function composition.

The outer middleware receives the inner handler:

```text
logger
   |
   v
version
   |
   v
handler
```

The problem appears when an application has many middleware:

```text
- panic recovery
- request ID
- logging
- metrics
- tracing
- authentication
- authorization
- rate limiting
- CORS
```

Manually nesting them becomes difficult to maintain.

The common solution is a **middleware chain helper**.

---

# 2. Defining a Middleware type

Instead of repeatedly writing:

```go
func(http.Handler) http.Handler
```

we can define a type:

```go
type Middleware func(http.Handler) http.Handler
```

This does not create new behavior.

It gives a name to an existing pattern.

Now we can write:

```go
func loggingMiddleware(next http.Handler) http.Handler
```

as:

```go
func loggingMiddleware(next http.Handler) http.Handler
```

and store it as:

```go
[]Middleware
```

The benefit is readability and API design.

A developer reading:

```go
func Chain(middlewares ...Middleware)
```

immediately understands:

"This function accepts middleware functions."

---

# 3. Building a middleware chain

A chain helper typically looks like:

```go
func Chain(middlewares ...Middleware) Middleware {
    return func(handler http.Handler) http.Handler {
        for i := len(middlewares) - 1; i >= 0; i-- {
            handler = middlewares[i](handler)
        }

        return handler
    }
}
```

The important part is not the loop.

The important part is the ordering.

Suppose we register:

```go
handler := Chain(
    loggingMiddleware,
    authMiddleware,
    recoveryMiddleware,
)(router)
```

A developer usually expects:

```text
logging
auth
recovery
router
```

But middleware wraps from the outside inward.

The resulting structure is:

```text
logging(
    auth(
        recovery(
            router,
        ),
    ),
)
```

Execution:

```text
logging before
    auth before
        recovery before
            router
        recovery after
    auth after
logging after
```

This is why the chain helper applies middleware in reverse order.

---

# 4. Middleware ordering changes behavior

Middleware ordering is not cosmetic.

Changing:

```text
logging
authentication
handler
```

to:

```text
authentication
logging
handler
```

changes what gets logged.

Example:

A request arrives without authentication.

With:

```text
logging
    authentication
        handler
```

the log middleware sees:

```text
GET /users/42 -> 401
```

With:

```text
authentication
    logging
        handler
```

the logging middleware may never execute because authentication stops the request.

---

Another example:

## Recovery placement

Consider:

```text
recovery
    logging
        handler
```

If the handler panics:

```text
recovery catches panic
```

because it surrounds the handler.

But:

```text
logging
    recovery
        handler
```

means logging is outside recovery.

The ordering affects:

- what gets executed
- what gets logged
- whether failures are captured
- whether metrics are recorded

This is why middleware ordering is an interview topic.

---

# 5. Build middleware once during startup

Middleware registration belongs during application initialization.

Example:

```go
func main() {
    router := createRouter()

    handler := Chain(
        recoveryMiddleware,
        requestIDMiddleware,
        loggingMiddleware,
    )(router)

    http.ListenAndServe(":8080", handler)
}
```

The important lifecycle:

```text
Application starts
        |
        v
Create router
        |
        v
Build middleware chain
        |
        v
Start HTTP server
        |
        v
Serve requests
```

Do not build middleware during every request.

Bad:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    middlewareChain := Chain(logging, auth)(router)

    middlewareChain.ServeHTTP(w, r)
}
```

Why?

Because middleware configuration is application setup.

Requests should execute an already-created handler chain.

---

# 6. Panic recovery middleware

A production HTTP server should usually prevent a panic from crashing the entire server process.

A recovery middleware looks conceptually like:

```go
func recoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        defer func() {
            if err := recover(); err != nil {
                // log panic
                // return error response
            }
        }()

        next.ServeHTTP(w, r)
    })
}
```

The important part is not `recover()` itself.

The important question is:

> What can recovery actually do when the response has already started?

---

## Response commitment

An HTTP response has two important stages:

### Before headers are written

The server has not committed the response.

You can still:

```go
w.WriteHeader(500)
w.Write(...)
```

### After headers are written

The response is already committed.

Example:

```go
w.WriteHeader(200)

panic("database failed")
```

At this point:

- the client already received `200`
- changing it to `500` is impossible

The recovery middleware cannot undo bytes already sent.

---

This means recovery middleware can:

✅ log the panic

✅ prevent process crash

✅ close the request safely

Possibly:

✅ send a 500 response

only if the response has not already started.

This is why production recovery middleware often works together with a custom `ResponseWriter` wrapper that tracks whether headers were written.

## Middleware ownership and responsibility

Middleware exists for **cross-cutting concerns**.

A cross-cutting concern is behavior that applies to many or most requests but is not the actual business operation being performed.

Examples:

```text
Request
   |
   v
Middleware
   |
   +-- authentication
   +-- authorization
   +-- logging
   +-- metrics
   +-- tracing
   +-- rate limiting
   +-- request IDs
   +-- panic recovery
   |
   v
Handler
   |
   +-- application behavior
```

Middleware answers questions like:

- "Is this request allowed to enter the system?"
- "How long did this request take?"
- "Did this request fail?"
- "Who made this request?"
- "Should this request be limited?"

It should generally **not answer business questions**.

---

## Example: authentication vs business logic

Authentication middleware:

```text
Does this request have a valid identity?
```

Good:

```go
func authMiddleware(next http.Handler) http.Handler {
    // validate token
    // attach user identity
    // continue request
}
```

The middleware does not know what the user is trying to do.

---

A handler:

```go
func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
    // create order
}
```

owns the business operation.

It decides:

- Is this product available?
- Is this order valid?
- Should inventory be reduced?
- Should payment be created?

---

## Example: authorization boundary

Authorization can be tricky.

Consider:

> "Can this user access order 123?"

This is not always middleware responsibility.

A middleware can check:

```text
Is the user authenticated?
Does the user have the "customer" role?
```

But:

```text
Does this specific user own order 123?
```

usually requires business data.

That belongs closer to the application/service layer:

```text
Handler
   |
   v
OrderService
   |
   v
Repository
```

because the rule depends on domain data.

---

## Avoid turning middleware into a "god layer"

A common mistake is:

```go
func middleware(next http.Handler) http.Handler {
    authenticate()
    authorize()
    validateInput()
    checkBusinessRules()
    updateDatabase()
    sendEmail()
}
```

This creates a hidden application layer.

Problems:

- difficult to test
- difficult to reuse
- unclear ownership
- middleware ordering becomes fragile

---

## A useful rule

A good question to ask:

> "Would this behavior apply to almost every request regardless of the business operation?"

If yes, middleware is a candidate.

Examples:

| Concern | Middleware? |
|-|-|
| Request logging | Yes |
| Request ID | Yes |
| Authentication token validation | Yes |
| Rate limiting | Yes |
| CORS | Yes |
| Metrics | Yes |
| Create invoice | No |
| Calculate shipping cost | No |
| Check order ownership | Usually no |
| Apply discount rules | No |

---

This addition would make Lesson 34 much more complete.

It also sets up later lessons nicely:
- **Context** → how middleware passes request-scoped information downstream
- **Authentication/Authorization** → where those responsibilities actually live
- **Application architecture** → handler/service/repository boundaries

So I would not consider Lesson 34 complete without this section.

---

# 7. Structured logging with `log/slog`

Previously we used:

```go
log.Printf(
    "%s %s %v",
    r.Method,
    r.URL.Path,
    elapsed,
)
```

This creates human-readable text.

Example:

```text
GET /users 15ms
```

The problem:

Machines do not understand this easily.

Production systems usually need logs that can be searched and analyzed.

For example:

```json
{
  "method": "GET",
  "path": "/users",
  "duration_ms": 15
}
```

This is structured logging.

---

Go provides:

```go
log/slog
```

The package provides structured logging using key-value attributes.

Example:

```go
logger.Info(
    "request completed",
    "method", r.Method,
    "path", r.URL.Path,
    "duration", elapsed,
)
```

Instead of embedding data into a string:

```go
"GET /users took 15ms"
```

we attach fields:

```text
message:
request completed

fields:
method=GET
path=/users
duration=15ms
```

---

Why is this useful?

A log aggregation system can query:

```text
show all requests where:
duration > 1000ms
```

or:

```text
show all POST requests to /payments
```

without parsing text.

This is why production systems use structured logs.

---

# 8. Request logging middleware

A typical logging middleware measures:

- HTTP method
- path
- status code
- duration
- request ID
- user identity (if appropriate)

Example:

```go
func loggingMiddleware(logger *slog.Logger) Middleware {
    return func(next http.Handler) http.Handler {

        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

            start := time.Now()

            next.ServeHTTP(w, r)

            logger.Info(
                "request completed",
                "method", r.Method,
                "path", r.URL.Path,
                "duration", time.Since(start),
            )
        })
    }
}
```

Later, when we improve our `ResponseWriter` wrapper, we can also capture:

```text
status=404
status=500
```

---

# 9. Request IDs

A request ID is a unique identifier assigned to each incoming request.

Example:

```text
Request:
GET /orders/123

Request-ID:
abc-123-def
```

Why?

Because one user request may create many operations:

```text
HTTP request
      |
      +-- database query
      |
      +-- external API call
      |
      +-- message publish
```

A request ID allows us to connect related logs.

---

At this stage, we only generate and log the ID.

Example:

```go
requestID := uuid.New().String()

logger.Info(
    "request started",
    "request_id", requestID,
)
```

Later, we will put this ID into:

```go
context.Context
```

so downstream handlers and services can access it.

That belongs to the context lesson.

---

# 10. What should NOT be logged

Logs are operational data.

They often end up in:

- log servers
- monitoring systems
- backups
- developer dashboards

Therefore, avoid logging sensitive data.

Do not log:

## Authentication data

Bad:

```text
Authorization: Bearer eyJ...
```

## Passwords

Bad:

```json
{
 "password":"secret123"
}
```

## Personal data unnecessarily

Examples:

- full addresses
- private identifiers
- unnecessary customer details

## Entire request bodies

Especially:

- payment requests
- authentication requests
- uploaded documents

---

A useful rule:

> Log enough information to diagnose a problem, but not enough information to become a security problem.

---

# Interview Questions

## 1. Why does middleware order matter?

Because middleware wraps handlers.

Changing the order changes:

- execution order
- whether middleware runs
- what information it can observe
- how failures are handled

---

## 2. Can recovery middleware always return HTTP 500 after a panic?

No.

Only if the response has not already been committed.

Once headers/body have been sent, the client may already have received the response.

Recovery can still:

- log the panic
- prevent server crash
- perform cleanup

but cannot rewrite the already-sent response.

---

# Exercise

Implement:

1. A `Middleware` type.
2. A `Chain` helper.
3. Three middleware functions:
   - logging
   - authentication check
   - recovery

Register them using the chain helper.

Verify the execution order by logging:

```text
before
after
```

from each middleware.

Then answer:

1. What changes if authentication middleware is placed before logging?
2. What happens if recovery is the innermost middleware instead of the outermost?

---

This completes the first production-level middleware composition lesson. The next logical step is **Lesson 35 — Request context, cancellation, deadlines, and propagating request-scoped data**, where the request ID moves from "generated and logged" to something every downstream component can use.