# Lesson 34 — Middleware composition and ordering

You already wrap the mux with middleware. In Lesson 32 it looked like this:

```go
handler := loggerHandler(versionEmbedHandler(mux))
```

That works for two wrappers. It gets messy when you add more. Each wrapper is another nested call. The function closest to `mux` wraps it first. The order the request actually travels is easy to misread.

This lesson does three things. First, a small helper so the chain is listed in one place. Then panic recovery. Then a request log you can actually use in production.

---

## A name for the middleware shape

You already write functions like this:

```go
func(http.Handler) http.Handler
```

Give that shape a name:

```go
type Middleware func(http.Handler) http.Handler
```

This is only a named function type. `loggerMiddleware` already matches it. Naming it lets a helper take several middleware as a list.

---

## A `chain` helper

```go
func chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
```

`...Middleware` means the function accepts any number of middleware. Inside the function they arrive as a slice.

The loop walks the slice **backwards**. That is the important part.

This:

```go
handler := chain(mux, recoverMiddleware, loggingMiddleware, versionMiddleware)
```

is the same as this:

```go
handler := recoverMiddleware(loggingMiddleware(versionMiddleware(mux)))
```

You list outermost first. The helper wraps from the inside out so that list matches how a request moves:

```text
recoverMiddleware
  loggingMiddleware
    versionMiddleware
      mux
```

If the loop walked forwards, the first name in the list would become the innermost wrapper. The list would then read backwards compared to the request. Walking backwards is a convention so the call site reads top-to-bottom. It is not a Go rule. Another helper might do the opposite. You need to know which end is outer.

**Order changes behaviour.** We will see a concrete case after panic is introduced. For now: do not treat the list as cosmetic.

Build the chain once in `main`. Pass that handler to `ListenAndServe`. Do not wrap the mux inside a request handler. That would rebuild the chain on every request.

---

## Panic

Go does not use exceptions for ordinary errors. You already handle those with `error` values.

A **panic** is different. It is for a situation the current function cannot continue: a nil pointer dereference, or an explicit `panic(...)` call.

When a panic happens, Go stops running the rest of that function. It then runs every `defer` in that function. Then it leaves the function and does the same in the caller. That walk continues up the call stack.

If nothing stops the panic, that goroutine dies. In a small `main` program, that usually kills the process.

HTTP handlers run on their own goroutines. A panic in a handler can crash more than that one request if it is not stopped.

---

## `recover`

`recover` stops a panic.

It returns the value that was passed to `panic`. If there is no panic, it returns `nil`.

It only works inside a function that was deferred. This does nothing useful:

```go
func handler(w http.ResponseWriter, r *http.Request) {
	recover() // no panic is active here
	next.ServeHTTP(w, r)
}
```

This is the pattern that works:

```go
defer func() {
	rec := recover()
	if rec == nil {
		return
	}
	// there was a panic; rec is the value
}()
```

The deferred function runs as the panic walks up the stack. That is when `recover` can catch it.

---

## Recovery middleware

`http.Server` already recovers panics in handlers. The process often stays up. That is not enough for an API.

You still want your own log of the panic. You still want your own 500 response, in the same JSON shape as other errors. And you need to handle the case where the handler already started writing the response. The server does not fix that case for you.

```go
func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			slog.Error("panic", "panic", rec)
			if !rw.wroteHeader {
				http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(rw, r)
	})
}
```

`statusWriter` is the wrapper from Lesson 33. It records the first status. The first `Write` counts as an implicit 200.

### If the response is already committed

Lesson 33: the first `WriteHeader` wins. Writing a body also commits a 200.

If the handler already committed, calling `WriteHeader(500)` does not change what the client received. Recovery can still log the panic. It must not write another body on top of a response that already started. That corrupts the stream.

So: if `wroteHeader` is true, log and stop. If it is false, send a 500.

---

## Order: recover and logging

Now the ordering point can be stated in terms you have.

Logging middleware usually does this:

```go
next.ServeHTTP(w, r)
log.Printf(...)
```

If `ServeHTTP` panics, the line after it does not run.

If recover is **outside** logging, the panic is caught in recover. Logging never reaches its `log.Printf` unless that log sits in a `defer`.

If logging is **outside** recover, logging can record a 500 after recover writes it. A panic *inside* the logger itself would not be caught.

A practical setup: recover outermost, so it covers everything inside. Put the request log in a `defer` inside logging, so a panic still records the request.

---

## `log/slog`

`log.Printf` writes one formatted string.

`log/slog` is the standard-library **structured** logger (Go 1.21+). You pass a message plus key-value fields. Tools can filter on those fields. They do not have to parse a line of text.

```go
slog.Info("request",
	"request_id", id,
	"method", r.Method,
	"path", r.URL.Path,
	"status", rw.status,
	"duration", time.Since(start),
)
```

`duration` stays a `time.Duration`. It is not baked into `"took 1.2ms"`.

For an API, JSON logs are the usual production format. Set that once at startup:

```go
slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
```

After that, `slog.Info` and `slog.Error` use this logger.

Replace `log.Printf` in the logging middleware. Use `defer` for the completion log so a panic still records the request:

```go
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		id := newRequestID()
		rw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		defer func() {
			slog.Info("request",
				"request_id", id,
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.status,
				"duration", time.Since(start),
			)
		}()
		next.ServeHTTP(rw, r)
	})
}
```

The request ID is created here on purpose. We have no way yet to pass it into the handler. That is Lesson 36 (`context`). For now, generate it where you log it.

Use `crypto/rand` for the ID. It fills a byte slice with unpredictable bytes. `encoding/hex` turns those bytes into text that is safe in logs.

```go
func newRequestID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	return hex.EncodeToString(b[:])
}
```

The time fallback is only if `rand.Read` fails. Do not use time as the normal ID. Two requests can share a nanosecond.

---

## What must not go in the log

A request log is stored. It gets copied. It gets searched. Treat it as public.

Do not log:

- `Authorization` or `Cookie` headers
- passwords, API keys, session tokens
- personal data (email, name, address) as fields
- the request body or the response body

Do log: method, path, status, duration, request ID.

Query strings are a judgement call. `/users?token=` will leak a secret if you log `r.URL.RawQuery`.

Bodies are the common leak. One debug log of the body ends up as a password in the log system.

---

## Putting it together in `main`

```go
slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

mux := http.NewServeMux()
// register routes

handler := chain(mux,
	recoverMiddleware,
	loggingMiddleware,
	versionMiddleware,
)

http.ListenAndServe(":8080", handler)
```

---

## Exercise

Take the product API from Lesson 33. Add a `chain` helper. Register at least three middleware with it. List outermost first in `main`.

The chain must:

1. Recover from a panic. If the response is not committed, the client gets a 500. If it is committed, the client keeps the status already sent, and you still log the panic.
2. Log each request with `log/slog`. Fields: request ID, method, path, status, duration. Use the Lesson 33 wrapper so implicit 200 and first-status-wins show up in the log.
3. Keep or add one other middleware you already have, so order is a real choice.

Add a route that panics before writing a response. Add a route that panics after `WriteHeader` or after writing a body. Compare status codes and logs.

Do not put the request ID on `context`. Do not log bodies.
