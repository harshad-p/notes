# Lesson 32 — Middleware: Intercepting and Wrapping HTTP Requests

We've now built handlers that deal directly with requests:

```text
HTTP request
    ↓
router
    ↓
handler
    ↓
response
```

But there are things we often want to do for **many or all endpoints**:

- log every request
- measure request duration
- authenticate a request
- add common response headers
- recover from unexpected panics
- enforce request-level rules

We could put that code into every handler, but that would duplicate logic.

**Middleware** gives us a way to place code around handlers.

---

## 1. What middleware actually is

A middleware is essentially a function that:

1. receives a handler
2. does something before that handler runs
3. calls the handler
4. optionally does something after the handler finishes

Conceptually:

```text
             middleware
                 │
HTTP request ────┤
                 ↓
              handler
                 │
                 ↓
HTTP response ───┘
```

A useful example is request logging.

We might want:

```text
GET /products/42
```

to produce a log such as:

```text
GET /products/42 12ms
```

We don't want every individual handler responsible for measuring its own execution time.

Middleware is a natural place for that cross-cutting behavior.

---

# 2. The middleware function shape

Before writing middleware, we need one new Go type from `net/http`:

```go
http.Handler
```

`http.Handler` is an interface representing something capable of handling an HTTP request.

Its definition is essentially:

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

So anything that provides:

```go
ServeHTTP(w, r)
```

satisfies `http.Handler`.

This is the interface underlying the HTTP server and the router we've already been using.

Our `http.HandlerFunc` functions are compatible with this system because Go provides an adapter that turns a function with the appropriate signature into an `http.Handler`.

The important distinction is:

```text
http.Handler
    ↓
interface

http.HandlerFunc
    ↓
function type implementing that interface
```

We'll use this distinction directly when writing middleware.

---

# 3. The basic middleware pattern

A middleware commonly has this shape:

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // before

        next.ServeHTTP(w, r)

        // after
    })
}
```

There are several things here that are new, so let's focus on them.

### `next http.Handler`

`next` represents the **next handler in the chain**.

Middleware doesn't necessarily know whether `next` is:

- an individual endpoint handler
- another middleware
- a router

It only needs to know that it satisfies `http.Handler`.

### `next.ServeHTTP(w, r)`

This invokes that handler.

So this is the crucial operation:

```go
next.ServeHTTP(w, r)
```

Without it, the request stops at the middleware.

The middleware has effectively intercepted the request but never passed it onward.

---

# 4. Middleware is a wrapper

Think of:

```go
loggingMiddleware(handler)
```

as producing a new handler:

```text
loggingMiddleware
       │
       ↓
   original handler
```

The returned handler performs the middleware's work and then delegates to the original handler.

This is essentially the **decorator/wrapper** pattern.

The important part is that middleware doesn't need to modify the original handler.

It wraps it.

---

# 5. A real logging middleware

Let's make the example useful.

We want to log:

- HTTP method
- URL path
- how long the handler took

For timing, we'll use Go's `time` package.

`time.Now()` gives us a `time.Time` representing the current time.

`time.Since(t)` calculates the duration elapsed since a particular `time.Time`.

So:

```go
start := time.Now()
```

captures the beginning.

Then later:

```go
duration := time.Since(start)
```

gives us the elapsed duration.

Our middleware becomes:

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        next.ServeHTTP(w, r)

        duration := time.Since(start)

        log.Printf("%s %s %v",
            r.Method,
            r.URL.Path,
            duration,
        )
    })
}
```

The important sequence is:

```text
request arrives
      ↓
start timer
      ↓
call next handler
      ↓
handler executes
      ↓
handler finishes
      ↓
calculate duration
      ↓
log result
```

This gives us logging without putting timing code into every endpoint.

---

# 6. Why the code after `ServeHTTP` executes

This is a useful property of middleware.

Consider:

```go
next.ServeHTTP(w, r)

duration := time.Since(start)
```

`next.ServeHTTP(...)` is a normal synchronous function call.

The middleware waits for it to return.

Therefore, code after it runs **after the downstream handler has completed**.

This allows middleware to perform both:

### Before

```go
start := time.Now()
```

### During

```go
next.ServeHTTP(w, r)
```

### After

```go
duration := time.Since(start)
```

That's why middleware is useful for things such as request timing.

---

# 7. Attaching middleware to the router

Suppose we have:

```go
mux := http.NewServeMux()

mux.HandleFunc("GET /products", getProductsHandler)
mux.HandleFunc("GET /products/{id}", getProductHandler)
mux.HandleFunc("POST /products", createProductHandler)
```

Previously, we passed the router directly to the server:

```go
http.ListenAndServe(":8080", mux)
```

Now we can wrap the router:

```go
handler := loggingMiddleware(mux)

http.ListenAndServe(":8080", handler)
```

The resulting structure is:

```text
HTTP request
     ↓
loggingMiddleware
     ↓
    mux
     ↓
matched endpoint
```

Notice that we wrapped the **router**, not each individual endpoint.

Therefore every request handled by that router passes through the middleware.

---

# 8. Multiple middleware

This becomes particularly useful when we have several pieces of cross-cutting behavior.

Suppose we have:

```text
logging
authentication
recovery
```

We could conceptually have:

```text
HTTP request
     ↓
logging middleware
     ↓
authentication middleware
     ↓
recovery middleware
     ↓
router
     ↓
endpoint
```

Each middleware calls the next handler.

For example:

```go
func authenticationMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // authentication logic

        next.ServeHTTP(w, r)
    })
}
```

Then:

```go
handler := loggingMiddleware(
    authenticationMiddleware(
        mux,
    ),
)
```

This works, but nested calls become difficult to read as the number of middleware grows.

We'll eventually look at cleaner ways to compose middleware.

For now, the important thing is understanding the underlying mechanism.

---

# 9. Middleware order matters

Suppose we have:

```text
A
 ↓
B
 ↓
handler
```

The request flows:

```text
A before
  ↓
B before
  ↓
handler
  ↓
B after
  ↓
A after
```

So middleware behaves somewhat like nested function calls.

If:

```go
A(B(handler))
```

then `A` is the outer wrapper and `B` is inside it.

This matters.

For example, suppose authentication happens inside logging:

```text
logging
   ↓
authentication
   ↓
handler
```

Then even rejected authentication attempts can be logged.

That's potentially desirable.

If the order were reversed, the behavior would be different.

Middleware ordering isn't merely stylistic; it can affect application behavior.

---

# 10. Middleware can stop the request

Calling the next handler is not mandatory.

For example, authentication middleware might determine that a request has no valid credentials.

It could send a response:

```go
http.Error(w, "Unauthorized", http.StatusUnauthorized)
```

and **not call**:

```go
next.ServeHTTP(w, r)
```

The chain stops there.

Conceptually:

```text
request
   ↓
authentication
   ↓
invalid
   ↓
401 response
```

The endpoint never executes.

This is one of the most important uses of middleware:

> Middleware can enforce a condition before allowing a request to proceed.

---

# 11. Middleware vs handler

It's worth establishing a clean distinction.

### Handler

A handler generally represents the behavior of a particular endpoint.

For example:

```text
GET /products/{id}
```

might retrieve one product.

### Middleware

Middleware represents behavior that surrounds or applies to other handlers.

For example:

```text
log request
authenticate request
measure duration
recover from panic
```

So:

```text
Handler:
"What should this endpoint do?"

Middleware:
"What should happen around requests before/after endpoints?"
```

This isn't an absolute separation—there are edge cases—but it's a useful architectural boundary.

---

# 12. Middleware doesn't automatically mean "global"

One subtle but useful point: middleware doesn't have to apply to every endpoint.

Because middleware is just a handler wrapper, you can choose what you wrap.

For example:

```text
router
├── public endpoints
└── protected endpoints
```

You might have:

```text
public handler
```

directly available, while protected handlers are wrapped in authentication middleware.

Conceptually:

```text
              router
             /      \
            /        \
       public       auth middleware
       handler          ↓
                    protected
                     handler
```

This is useful because not every endpoint necessarily requires authentication.

---

# 13. A subtle issue: middleware and response status

Suppose middleware wants to log the HTTP status code returned by the handler.

You might initially think:

```go
next.ServeHTTP(w, r)

log.Printf("status: ???")
```

But `http.ResponseWriter` doesn't directly provide a method like:

```go
w.StatusCode()
```

So the middleware needs to **wrap the `ResponseWriter` itself** and observe calls such as `WriteHeader`.

This is a more advanced middleware technique.

We won't implement it yet because it introduces another layer of HTTP internals that deserves its own explanation.

The important thing for now is:

> Middleware can wrap both the handler and, when necessary, the response writer.

We'll return to this when we build more production-oriented middleware.

---

# 14. Middleware and dependency injection

There is also a useful connection to the interfaces we've already learned.

Middleware accepts:

```go
next http.Handler
```

rather than a concrete handler type.

That means middleware doesn't care what the downstream handler actually is.

It could be:

```text
router
another middleware
specific endpoint
custom handler
```

As long as it implements `http.Handler`.

This is another practical example of Go's interface-based design: **depend on the behavior you need, not the concrete implementation.**

---

# 15. The complete flow

With our current API, we can now have:

```text
                 HTTP request
                      │
                      ▼
             Logging Middleware
                      │
                      ▼
             Authentication
                Middleware
                      │
                      ▼
                   ServeMux
                      │
              ┌───────┼────────┐
              ▼       ▼        ▼
           GET /   GET /{id}  POST /
              │       │        │
              ▼       ▼        ▼
           Handler Handler  Handler
```

And the response travels back through the chain as the nested handlers return.

This is the foundation for a lot of functionality you'll encounter in Go HTTP applications.

---

# Exercise

Take the product API you've been building.

Create a middleware that measures how long each request takes and logs:

```text
METHOD PATH DURATION
```

For example, the output should conceptually look like:

```text
GET /products 1.2ms
```

Apply it to the router so that **all product endpoints** pass through it.

Then add a second middleware that adds this response header:

```text
X-API-Version: 1
```

Think about:

- Where should the header be set?
- Does the middleware need to call `next`?
- What happens if you put the timing code before vs. after `next.ServeHTTP()`?
- In what order do the two middleware functions execute if you wrap one around the other?

Don't implement status-code capture yet—we'll deliberately introduce that mechanism later.

---

### Next: Lesson 33 — Middleware in Depth: ResponseWriter Wrapping, Status Codes & Request Logging