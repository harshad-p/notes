# Lesson 33 — Middleware in Depth: `ResponseWriter` Wrapping, Status Codes & Request Logging

In Lesson 32, we saw that middleware can wrap a handler:

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // before

        next.ServeHTTP(w, r)

        // after
    })
}
```

We also identified a limitation:

> If middleware wants to know which HTTP status the handler returned, `http.ResponseWriter` doesn't provide a `StatusCode()` method.

This lesson solves that problem.

---

# 1. Why would middleware need the response status?

Consider a request log:

```text
GET /products/42 200 3ms
```

versus:

```text
GET /products/42 404 2ms
```

versus:

```text
POST /products 400 1ms
```

For production logging, knowing only:

```text
GET /products/42 3ms
```

isn't enough.

We often want at least:

```text
method
path
status
duration
```

The problem is that the endpoint controls the status:

```go
w.WriteHeader(http.StatusNotFound)
```

but the middleware owns the outer `w` value.

So we need a way for the middleware to **observe what the handler does with the response writer**.

---

# 2. `http.ResponseWriter` is an interface

We've already encountered:

```go
w http.ResponseWriter
```

`http.ResponseWriter` is an interface.

The relevant methods are:

```go
type ResponseWriter interface {
    Header() Header
    Write([]byte) (int, error)
    WriteHeader(statusCode int)
}
```

We've used:

```go
w.Header()
```

and:

```go
w.WriteHeader(...)
```

and JSON encoding ultimately writes the response body through:

```go
w.Write(...)
```

The important point is that `w` isn't necessarily a particular concrete struct that we can inspect for internal state.

We're interacting with it through its interface.

And the interface doesn't expose:

```go
StatusCode()
```

So we can't simply ask:

```go
status := w.StatusCode()
```

---

# 3. The wrapping technique

The solution is similar to what middleware already does with handlers.

We wrap the `ResponseWriter`.

We create our own type:

```go
type responseWriter struct {
    http.ResponseWriter
    statusCode int
}
```

There are two things to notice.

### Embedded `http.ResponseWriter`

This:

```go
http.ResponseWriter
```

is an **embedded field**.

We learned struct embedding earlier.

It means our `responseWriter` contains an underlying `http.ResponseWriter` and promotes its methods.

So if we haven't overridden a method, calls can still be delegated to the embedded writer.

### `statusCode`

We add our own field:

```go
statusCode int
```

This is where we'll record the status selected by the handler.

So conceptually:

```text
responseWriter
├── underlying ResponseWriter
└── statusCode
```

---

# 4. Intercepting `WriteHeader`

Now we want our wrapper to notice when the handler calls:

```go
w.WriteHeader(http.StatusCreated)
```

We can define our own `WriteHeader` method:

```go
func (w *responseWriter) WriteHeader(statusCode int) {
    w.statusCode = statusCode
    w.ResponseWriter.WriteHeader(statusCode)
}
```

This is the important part.

When the handler calls:

```go
w.WriteHeader(201)
```

it's actually calling our wrapper's method.

Our method does two things:

```text
1. Remember the status
2. Forward the status to the real ResponseWriter
```

So:

```text
handler
   |
   | WriteHeader(201)
   ↓
our responseWriter
   |
   ├── statusCode = 201
   |
   └── underlying.WriteHeader(201)
             |
             ↓
       actual HTTP response
```

The middleware can now inspect:

```go
w.statusCode
```

after the handler finishes.

---

# 5. Why do we forward the call?

This part is crucial.

We could write:

```go
func (w *responseWriter) WriteHeader(statusCode int) {
    w.statusCode = statusCode
}
```

But that would only record the status.

It would **not actually send the status to the client**.

The wrapper is an observer/decorator, not a replacement for the underlying response writer.

Therefore:

```go
w.ResponseWriter.WriteHeader(statusCode)
```

must still happen.

This pattern is extremely common when wrapping interfaces:

```text
intercept
  ↓
record/modify something
  ↓
delegate to underlying implementation
```

---

# 6. The default status problem

There's an important subtlety.

A handler doesn't have to explicitly call:

```go
w.WriteHeader(http.StatusOK)
```

For example:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello")
}
```

The handler writes the body, and Go implicitly sends:

```text
200 OK
```

But our wrapper's `WriteHeader` isn't necessarily called explicitly by the handler.

Therefore, if we initialize:

```go
statusCode int
```

to zero, our middleware could incorrectly log:

```text
0
```

instead of:

```text
200
```

We need to represent the default behavior.

A simple solution is:

```go
rw := &responseWriter{
    ResponseWriter: w,
    statusCode:     http.StatusOK,
}
```

Now the wrapper assumes `200` unless the handler explicitly selects another status.

This works well for our current purpose.

---

# 7. But what happens when the handler calls `Write`?

There's another subtlety.

Suppose the handler does:

```go
fmt.Fprintln(w, "Hello")
```

or:

```go
json.NewEncoder(w).Encode(product)
```

Those eventually write the response body.

The HTTP response is committed at that point if no status has already been written, meaning the effective status becomes `200`.

Our wrapper should therefore be aware of `Write` as well if we want its status tracking to accurately mirror HTTP behavior.

We can override it:

```go
func (w *responseWriter) Write(data []byte) (int, error) {
    return w.ResponseWriter.Write(data)
}
```

But this alone doesn't update anything.

We could make it explicitly establish the default status:

```go
func (w *responseWriter) Write(data []byte) (int, error) {
    if !w.wroteHeader {
        w.WriteHeader(http.StatusOK)
    }

    return w.ResponseWriter.Write(data)
}
```

This introduces another piece of state:

```go
wroteHeader bool
```

Now the wrapper can distinguish:

```text
header hasn't been written
```

from:

```text
header has already been written
```

Our type becomes:

```go
type responseWriter struct {
    http.ResponseWriter
    statusCode  int
    wroteHeader bool
}
```

and:

```go
func (w *responseWriter) WriteHeader(statusCode int) {
    if w.wroteHeader {
        return
    }

    w.statusCode = statusCode
    w.wroteHeader = true

    w.ResponseWriter.WriteHeader(statusCode)
}
```

Then:

```go
func (w *responseWriter) Write(data []byte) (int, error) {
    if !w.wroteHeader {
        w.WriteHeader(http.StatusOK)
    }

    return w.ResponseWriter.Write(data)
}
```

Now our wrapper tracks the important behavior more accurately.

---

# 8. Why ignore a second `WriteHeader`?

Notice this:

```go
if w.wroteHeader {
    return
}
```

HTTP response headers can only effectively be committed once.

For example:

```go
w.WriteHeader(http.StatusCreated)
w.WriteHeader(http.StatusBadRequest)
```

doesn't mean the client receives both statuses.

Once the response has been committed as `201`, you can't change it to `400`.

Our wrapper therefore needs to preserve that behavior.

The first status wins.

This is also why this ordering is important:

```go
w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(product)
```

rather than trying to change the status after the body has started.

---

# 9. The complete wrapper

Putting the pieces together:

```go
type responseWriter struct {
    http.ResponseWriter
    statusCode  int
    wroteHeader bool
}

func (w *responseWriter) WriteHeader(statusCode int) {
    if w.wroteHeader {
        return
    }

    w.statusCode = statusCode
    w.wroteHeader = true

    w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriter) Write(data []byte) (int, error) {
    if !w.wroteHeader {
        w.WriteHeader(http.StatusOK)
    }

    return w.ResponseWriter.Write(data)
}
```

This is our first genuinely useful example of **decorating an interface implementation**.

We're not replacing the HTTP server's response writer.

We're putting another layer around it.

---

# 10. Using the wrapper in middleware

Now our logging middleware can create the wrapper:

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        rw := &responseWriter{
            ResponseWriter: w,
            statusCode:     http.StatusOK,
        }

        next.ServeHTTP(rw, r)

        duration := time.Since(start)

        log.Printf(
            "%s %s %d %v",
            r.Method,
            r.URL.Path,
            rw.statusCode,
            duration,
        )
    })
}
```

The important change is:

```go
next.ServeHTTP(rw, r)
```

rather than:

```go
next.ServeHTTP(w, r)
```

We're deliberately giving the downstream handler **our wrapper**.

The handler still sees:

```go
http.ResponseWriter
```

because `responseWriter` satisfies that interface.

But internally, our wrapper can observe the status.

---

# 11. What the handler actually sees

This is a particularly nice property of interfaces.

Our endpoint doesn't need to know that middleware has wrapped the writer.

It still receives:

```go
func handler(w http.ResponseWriter, r *http.Request)
```

From its perspective, `w` is simply an `http.ResponseWriter`.

It can do:

```go
w.WriteHeader(http.StatusNotFound)
```

or:

```go
json.NewEncoder(w).Encode(product)
```

normally.

The wrapper transparently forwards those operations to the underlying writer.

So the dependency looks like:

```text
Handler
   |
   | http.ResponseWriter
   ↓
responseWriter
   |
   | http.ResponseWriter
   ↓
real HTTP response writer
```

This is exactly the sort of interface-based composition Go is good at.

---

# 12. Logging the response status

Now we can produce useful logs.

For example:

```text
GET /products 200 1.3ms
GET /products/42 200 900µs
GET /products/999 404 1.1ms
POST /products 400 700µs
```

The middleware doesn't need to understand what each endpoint does.

It only observes:

- request metadata
- response status
- execution duration

That's a strong example of a **cross-cutting concern**.

---

# 13. One important production caveat

There is a subtle problem with the wrapper we've created.

`http.ResponseWriter` has historically been a fairly small interface, but concrete implementations can also expose additional optional capabilities.

For example, some HTTP handlers may need functionality associated with:

```go
http.Flusher
```

or:

```go
http.Hijacker
```

or:

```go
io.ReaderFrom
```

If we blindly wrap `ResponseWriter`, our wrapper may hide those additional interfaces from downstream code.

This matters for things such as:

- streaming responses
- WebSockets
- optimized data transfers

So production-grade response-writer middleware needs to be more careful about preserving optional interfaces.

We're **not** going to implement all of that now. It would distract from the core concept.

But I want you to know that:

> "Just wrap `http.ResponseWriter`" is conceptually correct for learning the pattern, but a production middleware implementation has additional compatibility considerations.

That's the kind of subtlety worth knowing before blindly copying this pattern into a serious application.

---

# 14. Another important caveat: status is not necessarily the whole story

Our middleware currently logs:

```text
method
path
status
duration
```

That's useful, but production observability often needs more:

```text
method
path
status
duration
request ID
user identity
response size
error information
```

We could keep adding things to this middleware.

But that doesn't necessarily mean we should.

For example, authentication and request IDs are conceptually different concerns and can be separate middleware.

A good middleware chain might eventually look like:

```text
request
   ↓
request ID middleware
   ↓
logging middleware
   ↓
authentication middleware
   ↓
router
   ↓
handler
```

Each middleware has a focused responsibility.

---

# 15. Middleware is composition, not magic

At this point, the whole mechanism can be understood without treating middleware as some special HTTP feature.

It's essentially composition:

```text
Handler
   ↑
wrapped by middleware
   ↑
wrapped by another middleware
```

And the same idea applies to the response writer:

```text
real ResponseWriter
        ↑
responseWriter wrapper
        ↑
handler
```

The handler only knows about the interface.

The middleware can add behavior around that interface.

This is one of the most useful patterns to internalize in Go.

---

# Exercise

Extend your logging middleware so that it logs:

```text
METHOD PATH STATUS DURATION
```

using your own `responseWriter` wrapper.

Test it against endpoints that produce:

- an implicit `200`
- an explicit `201`
- a `400`
- a `404`

Pay particular attention to this case:

```go
w.WriteHeader(http.StatusCreated)
w.WriteHeader(http.StatusBadRequest)
```

What should your middleware record, and why?

Also test an endpoint that **only writes a body** without explicitly calling `WriteHeader`.

You should be able to explain why your wrapper records `200` in that case.

## Why func (w *responseWriter) WriteHeader(statusCode int) has a pointer? 

There are two separate concepts here.

### 1. Implementing an interface does not require `*`

For example:

```go
type Speaker interface {
    Speak()
}

type Dog struct{}

func (d Dog) Speak() {
    fmt.Println("Woof")
}
```

`Dog` implements `Speaker`. No `*` is required.

You could also write:

```go
func (d *Dog) Speak() {
    fmt.Println("Woof")
}
```

In that case, **`*Dog` implements `Speaker`**, not `Dog`.

So the choice of value vs pointer receiver is independent of the fact that an interface is involved.

---

### 2. Why `responseWriter` needs a pointer receiver

Our type was:

```go
type responseWriter struct {
    http.ResponseWriter
    statusCode  int
    wroteHeader bool
}
```

And:

```go
func (w *responseWriter) WriteHeader(statusCode int) {
    w.statusCode = statusCode
    w.wroteHeader = true
    w.ResponseWriter.WriteHeader(statusCode)
}
```

We want `WriteHeader` to modify the **actual `responseWriter` instance**:

```go
w.statusCode = statusCode
w.wroteHeader = true
```

A value receiver:

```go
func (w responseWriter) WriteHeader(statusCode int)
```

would give the method a **copy** of the struct.

So modifications would happen to the copy, not the original wrapper.

That's why we use:

```go
*responseWriter
```

It means:

> `w` is a pointer to the actual `responseWriter`, so modifications affect the original object.

This is the same pointer-receiver concept we already covered earlier with structs and methods.

---

### The important distinction

Think of these as two independent questions:

**Question 1: What does the receiver need to be?**

```go
func (w responseWriter) ...
```

vs.

```go
func (w *responseWriter) ...
```

This is about **value vs pointer semantics**.

**Question 2: Does the type implement an interface?**

That depends on the method set of the type.

With:

```go
func (w *responseWriter) WriteHeader(...)
```

the pointer type `*responseWriter` has that method.

And because we're passing:

```go
rw := &responseWriter{...}
next.ServeHTTP(rw, r)
```

`rw` is a `*responseWriter`, so it is the type being used as the `http.ResponseWriter`.

### One correction to my Lesson 33 teaching

I should have explicitly said:

> "We're using a pointer receiver because the wrapper's methods need to mutate its state. The `*` has nothing specifically to do with implementing `http.ResponseWriter`."

You already learned pointer receivers earlier, so I should have connected this to that concept instead of silently introducing the `*` in the interface example.

### Next: Lesson 34 — Request Context: Cancellation, Deadlines & Passing Request-Scoped Data