# Go Backend Engineering — Master Curriculum

## Goal

Reach professional Go backend proficiency: build, test, secure, operate, and defend a production Go service — and explain every design decision behind it.

The course uses backend development as its spine. Go language features are introduced at the point where backend work makes them necessary, then deepened later in a new context rather than retaught. The progression is:

```text
Go fundamentals
      ↓
Build useful backend functionality
      ↓
Encounter new Go concepts naturally
      ↓
Build increasingly realistic APIs
      ↓
Deepen Go knowledge where necessary
      ↓
Databases / concurrency / messaging
      ↓
Production engineering
      ↓
Professional Go backend proficiency
```

## The three proficiency areas

Every phase advances at least one of these; the tag under each phase heading shows which.

| Area | What it means |
| --- | --- |
| **Go Language** | Read and write idiomatic, maintainable Go: type system, structs, methods, pointers, interfaces, errors, generics, standard library, package design. |
| **Backend / API Engineering** | Build real systems: HTTP, routing, JSON, API design, middleware, context, external APIs, auth, PostgreSQL, transactions, testing, concurrency, background work, queues, Kafka, event-driven architecture. |
| **Production / Engineering** | Run services for real: configuration, logging, observability, metrics, tracing, reliability, graceful shutdown, health checks, security, performance, profiling, tooling, CI/CD, architecture, Docker, deployment. |

## How to read this

- Each lesson lists the topics it teaches. That list is the source of truth for what the lesson covers.
- Lessons that have been written also name the repository folder holding them and their code — under the phase directory. The lesson sequence, not the folder name, is the authoritative order. Some folders hold several lessons.
- A topic already taught is never rescheduled; where a later lesson revisits a subject, it names the deeper concern it is there to solve.
- `progress.md` records what is done and what comes next.

### What "covered" means here

A topic counts as learned only when a lesson taught it deliberately. Four states are kept distinct, because several topics have been *seen* without being learned:

| State | Meaning |
| --- | --- |
| **Taught** | A written lesson explained it deliberately. |
| **Mentioned** | Named to motivate or foreshadow something else, with the explanation explicitly postponed. |
| **Incidental** | Appeared inside example code without being the subject of the lesson. |
| **Not covered** | Has not appeared at all. |

These have appeared but are **not** learned, and each has a scheduled lesson:

| Seen | Where it appeared | Taught in |
| --- | --- | --- |
| `error` wrapping, custom error types | Lesson 16 promises them later | Lesson 40 |
| Channels, `select` | Not yet seen; `ctx.Done()` will surface one | Lesson 51 |
| PostgreSQL, `database/sql` | Named as a possible future implementation in Lesson 18; mentioned again in Lesson 30 | Lessons 38–39 |
| Test fakes | Motivated the interface in Lesson 18 | Lesson 37 |
| Optional `ResponseWriter` interfaces | Named as a caveat in Lesson 33 | Lesson 35 |
| Runes, UTF-8 internals | Named in Lesson 21 without explanation; the minimal bridge is owed there | Lesson 49 |
| Go's full numeric types | `float64` naming deferred in Lesson 4 | Lesson 46 |
| `GOROOT` / `GOPATH` | Shown in `go env` output in Lesson 0, explicitly postponed | Lesson 17 (top-up owed) |
| `context.Context` | Drafted, then pulled back out as premature | Lesson 36 |

## How the course is taught

- **One service, carried forward.** The product API started in Lesson 25 is the same codebase that gains routing, validation, response contracts, and middleware. Every future phase extends it rather than starting a new toy. Lesson 60 is that service reaching production shape, not a new project.
- **Backend need first, language feature second.** A Go concept is introduced when the service needs it, then deepened later only when the deeper treatment is genuinely different.
- **No architecture by default.** Service layers, repositories, interfaces, factories, and folder conventions appear only when a concrete problem — a dependency to swap, a boundary to defend, a lifecycle to own, a test seam to create — justifies them, and the trade-off is stated. Lesson 18 introduced an interface because a test fake needed one; that is the standard every later abstraction must meet.
- **Depth without narration.** Lessons explain mechanisms, edge cases, and design consequences for an experienced developer. They do not walk through obvious code line by line.
- **One exercise per lesson**, realistic, with no solution or step-by-step hints in the statement.
- **Lesson shape.** What it is → why it exists → the Go mechanics → applied to our service → how it is used in production → alternatives and trade-offs → the mistakes that actually cause incidents → exercise.
- **Interview checkpoints.** Each phase ends with the learner answering that phase's interview questions before any answer is given. The per-lesson *Interview* bullets below are the source of those questions.
- **Completion gates.** A lesson is marked done in `progress.md` when the concept is understood at a level usable at work — not when it has been read.

---

# Phase 1 — Go Foundations

*Go Language* · Lessons 0–16

## Lesson 0 — Install Go

`phase-01-go-foundations/000-install-go`

- The Go toolchain
- `go version`
- Platform install paths
- `go env`
- `GOROOT` / `GOPATH` shown in `go env` output, not explained
- *Deferred:* what those values mean → Lesson 17

## Lesson 1 — Understanding a Go project

`phase-01-go-foundations/001-understanding-a-go-project`

- A folder is a module: no solution or project file
- `go mod init`
- What `go.mod` tracks
- Rough equivalence to `.csproj` / `package.json`
- `go run .` rather than `go run main.go`

## Lesson 2 — Structure of a Go program

`phase-01-go-foundations/002-structure-of-a-go-program`

- `package main` and what makes a program executable
- `import`
- `func main` as the entry point
- Brace placement as a consequence of automatic semicolon insertion
- `go build` versus `go run`
- `go.mod` contents

## Lesson 3 — Variables

`phase-01-go-foundations/003-go-basic-types`

- `var` with an explicit type
- `:=` short declaration
- Why `:=` cannot redeclare an existing variable

## Lesson 4 — Basic types

`phase-01-go-foundations/003-go-basic-types`

- Static typing
- `string`, `int`, `float64`, `bool`
- `Println` inserting spaces between arguments
- *Deferred:* why `float64` includes its width → Lesson 46
- *Deferred:* `fmt.Print` / `Printf` → Lesson 21

## Lesson 5 — Constants

`phase-01-go-foundations/003-go-basic-types`

- `const`
- Must initialize immediately
- When a constant is preferable to a variable

## Lesson 6 — Functions

`phase-01-go-foundations/004-go-functions`

- `func`
- `name type` parameter order versus C#'s `type name`
- Return types
- Grouping parameters of the same type

## Lesson 7 — Multiple return values

`phase-01-go-foundations/004-go-functions`

- Returning `(string, int)`
- Destructuring at the call site
- `(value, error)` previewed as Go's replacement for exceptions
- `error` and `nil` used here, taught in Lesson 16

## Lesson 8 — If/else and loops

`phase-01-go-foundations/005-if-else-and-loops`

- `if` / `else` without parentheses, mandatory braces
- `for` as the only loop keyword
- Counted `for`
- Condition-only `for` (`while` form)
- Infinite `for`
- `%`
- *Deferred:* `switch` → Lesson 26

## Lesson 9 — Slices

`phase-01-go-foundations/006-arrays-and-slices`

- Slice literals
- Indexing
- `append` and why its result must be reassigned
- `len`
- `range` with index and value
- The blank identifier `_`
- Unused variables do not compile

## Lesson 10 — Arrays versus slices

`phase-01-go-foundations/006-arrays-and-slices`

- Fixed-size `[3]int` versus `[]int`
- When an array is the right choice
- Slice header: pointer, length, capacity
- What happens to the backing array on growth
- Empty-slice-then-append for building API responses

## Lesson 11 — Maps

`phase-01-go-foundations/007-maps`

- Map literals
- Reading, add, and update by assignment
- `delete`
- Zero values: a missing key returns the value type's zero value
- Comma-`ok` to distinguish absence from a zero value

## Lesson 12 — Structs

`phase-01-go-foundations/008-structs`

- `type ... struct`
- Literals with field names
- Field access and mutation
- Capitalization as visibility, taught fully in Lesson 17

## Lesson 13 — Methods and receivers

`phase-01-go-foundations/008-structs`

- Methods declared outside the type
- The receiver
- A value receiver operates on a copy

## Lesson 14 — Pointers

`phase-01-go-foundations/008-structs`

- Arguments passed by copy
- `&` for address-of
- `*T` as a pointer type
- Pointer receivers for mutation
- Go taking the address automatically on a method call

## Lesson 15 — Interfaces

`phase-01-go-foundations/008-structs`

- Implicit satisfaction: no `implements` keyword
- Depending on behavior rather than a concrete type
- Two types satisfying one interface

## Lesson 16 — Error handling

`phase-01-go-foundations/009-errors`

- `(value, error)` return convention
- `nil` as "no error"
- `if err != nil`
- `error` as an ordinary interface with `Error() string`
- Errors are values
- `errors.New` versus `fmt.Errorf`
- Early-return error path
- *Deferred:* wrapping, sentinel errors, custom error types → Lesson 40

---

# Phase 2 — Organizing Go Code and Connecting Language Features

*Go Language · Backend / API Engineering* · Lessons 17–24

## Lesson 17 — Packages and modules

`phase-02-organizing-go-code/010-packages-and-modules`

- Several files in one package
- A second package in a subdirectory
- Capitalization as the export rule
- The module path as the root of import paths
- What `go mod init` establishes
- *Top-up owed:* `GOROOT` versus `GOPATH` versus modules
  - `GOROOT` is the toolchain and standard library; you almost never set it
  - `GOPATH` was the pre-modules workspace; putting projects there is obsolete
  - What `GOPATH` still is: module cache (`pkg/mod`) and default `go install` bin directory (`GOBIN` / `GOPATH/bin`)
  - Why a module (`go.mod`) is not "a project inside GOPATH"
  - `go env GOROOT`, `go env GOPATH`, `go env GOMODCACHE`

## Lesson 18 — A small service boundary

`phase-02-organizing-go-code/011-pointers-structs-interfaces`

- A struct model with behavior
- An interface describing only what a consumer needs
- An in-memory repository satisfying it implicitly
- A service holding the interface as a dependency
- Explicit wiring in `main` instead of a DI container
- A test fake as the motivation for the interface
- Composition instead of inheritance
- Interface *use*, not interface design in depth
- PostgreSQL / `sql.DB` named only as a possible later implementation of the same interface; not taught here → Lessons 38–39

## Lesson 19 — Cleanup with `defer`

`phase-02-organizing-go-code/012-defer`

- Deferred execution at function exit
- Cleanup that survives early returns
- LIFO ordering
- Arguments evaluated at `defer` time, not call time

## Lesson 20 — Dependencies

`phase-02-organizing-go-code/013-external-dependencies`

- `go get`
- Importing a third-party package
- Roles of `go.mod` and `go.sum`
- `go mod tidy`

## Lesson 21 — Strings in Go

`phase-02-organizing-go-code/014-strings`

- Concatenation
- A Go string is a sequence of bytes holding UTF-8-encoded text
- `len` counts bytes, not human-readable characters (`"café"` as the example)
- Indexing yields a `byte` (numeric code, not a character)
- `fmt.Printf` and verbs, introduced via `%c`
- The `strings` package: `Contains`, `ToUpper`, `ToLower`
- `rune` is named, not explained
- *Top-up owed:* a `rune` is a Unicode code point; byte indexing can split a multi-byte character; `range` over a string decodes runes
- *Deferred:* decoding, normalization, `unicode/utf8` → Lesson 49

## Lesson 22 — Pointer syntax in practice

`phase-02-organizing-go-code/015-pointers-part-2`

- `&`
- `*` in a type versus `*` as dereference
- Mutation through a pointer
- Automatic dereferencing for field and method access
- Deepens Lesson 14 at the syntax level

## Lesson 23 — Composition

`phase-02-organizing-go-code/016-struct-embedding-and-composition`

- Struct embedding
- Field and method promotion
- Embedding expresses "has-a" composition, not inheritance

## Lesson 24 — JSON at the model boundary

`phase-02-organizing-go-code/017-struct-tags-and-json`

- Struct tags as field metadata
- `encoding/json`
- `Marshal` producing `[]byte`
- `Unmarshal` requiring a pointer
- `omitempty`
- JSON as a language/stdlib feature; HTTP transport is Lesson 25

---

# Phase 3 — Building the HTTP API

*Backend / API Engineering · Go Language* · Lessons 25–31

## Lesson 25 — HTTP request and response bodies

`phase-03-building-the-http-api/018-JSON-in-HTTP-APIs`

- Shape of an HTTP request: method, URL, headers, body
- `net/http`
- `ListenAndServe`
- The handler signature
- What `*http.Request` exposes
- `http.ResponseWriter` as the response-construction mechanism
- `fmt.Fprintln` writing to a destination rather than stdout
- `w.Header().Set`
- `r.Body` as a stream
- `json.NewEncoder` / `NewDecoder` versus in-memory `Marshal` / `Unmarshal`
- `http.Error`
- Returning `400` for malformed JSON

## Lesson 26 — Methods, routing, status codes

`phase-03-building-the-http-api/019-http-routing-methods-and-status-codes`

- `r.Method` and `http.MethodGet`-style constants
- `switch`, introduced here for method dispatch
- Contrast with C# fallthrough rules
- Manual route registration with `HandleFunc`
- Routing ("which handler") versus method ("which operation")
- `WriteHeader` and the fact that the status is committed once
- `201 Created`
- `405 Method Not Allowed`
- Implicit `200` when a handler writes a body without selecting a status

## Lesson 27 — URL input

`phase-03-building-the-http-api/020-url-paths-query-parameters-and-path-parameters`

- URL anatomy
- `r.URL.Path`
- Extracting a path segment with `strings.Split`, and why that does not scale
- `r.URL.Query()` and `.Get`
- An absent parameter yields the empty string
- Multiple filters
- Path parameter identifies a resource; query parameter modifies retrieval
- *Deferred:* automatic path extraction → Lesson 28

## Lesson 28 — `ServeMux` route patterns

`phase-03-building-the-http-api/021-http-router`

- `http.NewServeMux`
- Passing an explicit router to the server
- Method-qualified patterns such as `GET /users`
- `{id}` wildcards, including several in one pattern
- `r.PathValue`
- Why this beats hand-splitting paths
- How the router distinguishes `404 Not Found` from `405 Method Not Allowed`

## Lesson 29 — Turning URL strings into Go values

`phase-03-building-the-http-api/022-type-conversion-and-parsing`

- Type conversion between compatible Go types versus parsing text into a value
- `strconv.Atoi`
- `ParseBool`
- `ParseFloat`
- When `ParseInt` is warranted
- Parsing failure as a `400`
- Malformed value versus well-formed but unacceptable
- Absent query parameter versus invalid: supplying a default

## Lesson 30 — API contracts and validation

`phase-03-building-the-http-api/023-json-api-request-response-model`

- Separate request, response, and persistence types when their contracts differ
- Not multiplying types when they do not
- Decoding: "is this JSON understandable"
- Validation: "is this data acceptable"
- Input validation versus business validation
- `400 Bad Request` for both malformed and invalid input
- Pointer fields so PATCH can distinguish "field absent" from "field set to zero"

## Lesson 31 — Consistent JSON responses

`phase-03-building-the-http-api/024-consistent-api-responses`

- Headers → status → body, and why writing the body commits the response
- Choosing success statuses per operation
- A JSON error shape and a `writeError` helper
- A `writeJSON` helper
- `any`, the alias for the empty interface, introduced here
- Why `204 No Content` must not go through a body-writing helper
- Keeping HTTP status semantics out of domain models

---

# Phase 4 — Middleware Mechanics

*Backend / API Engineering · Go Language* · Lessons 32–33

## Lesson 32 — Middleware mechanics

`phase-04-middleware-mechanics/025-middleware`

- `http.Handler` as an interface
- `http.HandlerFunc` as the function adapter that satisfies it
- The `func(next http.Handler) http.Handler` shape
- Delegating with `next.ServeHTTP`, and what happens when you do not
- Before/after execution around a synchronous downstream call
- `time.Now` / `time.Since` for request timing
- `log.Printf`
- Wrapping the whole router versus wrapping selected handlers
- Nesting, composition, and why ordering changes behavior
- Short-circuiting the chain
- Boundary between cross-cutting middleware and endpoint logic
- *Deferred:* status capture → Lesson 33

## Lesson 33 — Response observation

`phase-04-middleware-mechanics/026-middleware-in-depth`

- `http.ResponseWriter` as an interface with no way to read back the status
- Decorating it with a struct that embeds the interface
- Unoverridden methods still reaching the underlying writer
- Intercepting `WriteHeader` to record the status and forwarding it
- Overriding `Write` to model the implicit `200`
- `wroteHeader` so the first status wins
- Why the wrapper needs a pointer receiver (mutation, not "to implement the interface")
- Method-set consequence: `*responseWriter`, not `responseWriter`, is the implementing type
- Status-and-duration request logging
- Naive wrappers hide `http.Flusher`, `http.Hijacker`, and `io.ReaderFrom`
- *Deferred:* preserving those optional interfaces → Lesson 35

---

# Phase 5 — Finish the HTTP Boundary Properly

*Backend / API Engineering · Production / Engineering* · Lessons 34–37

The chain from Lessons 32–33 is still hand-nested, has no recovery, cannot pass anything to a handler, and has nothing verifying that it behaves.

### The middleware track

Middleware is a thread through the rest of the course. Each concern lands where its prerequisites exist:

| Concern | Lesson |
| --- | --- |
| `http.Handler`, `HandlerFunc`, wrapping, `next.ServeHTTP`, before/after, short-circuiting, middleware versus handler | 32 |
| `ResponseWriter` wrapping, status capture, implicit `200`, first-write-wins | 33 |
| Reusable chains, nesting and execution order, ordering bugs, panic recovery, structured logging | 34 |
| Rate limiting, CORS, security headers, body limits, `Flusher`/`Hijacker`/`ReaderFrom`, router and framework approaches, middleware versus service logic | 35 |
| Timeout middleware, request/correlation IDs reaching handlers | 36 |
| Testing middleware order and behaviour | 37 |
| Authentication and authorization middleware, and what must not live in it | 43–44 |
| Metrics and tracing middleware, correlation across components | 56 |

## Lesson 34 — Middleware composition and ordering

- A chain helper replacing nested calls
- Reverse-wrapping order so a chain reads top-to-bottom
- Building the chain once at startup, not per request
- Panic recovery that must cope with a response already committed
- Request-scoped structured logging with `log/slog` replacing `log.Printf`
- Generating a request ID and logging it, without yet handing it to the handler
- What a request log must not contain: credentials, tokens, personal data, whole bodies
- **Interview:** why ordering changes behavior; what recovery can still do after the status has been written
- *Deferred:* handing the request ID to downstream code → Lesson 36

## Lesson 35 — Middleware at the edge of production

- CORS, and why a permissive policy is a vulnerability
- Security headers
- Rate limiting: token bucket, per-key versus global, what breaks across multiple instances
- Request-size and body limits
- A `ResponseWriter` wrapper that preserves `http.Flusher`, `http.Hijacker`, and `io.ReaderFrom`
- How the standard library, common routers, and full frameworks express middleware
- What belongs in middleware versus the handler: a rule that depends on the resource cannot live in middleware
- Why authorization is only partly a middleware concern when it arrives in Lesson 44
- **Interview:** implementing rate limiting; why wrapping the writer can break a dependency; middleware versus handler

**Phase 5 interview checkpoint (after Lesson 35, before answers):** what middleware is and how a chain is constructed; why order changes behaviour; how request logging and status capture are implemented; what happens when `WriteHeader` is called twice; why wrapping `http.ResponseWriter` can be problematic and how to do it safely; how authentication middleware would be structured; how a request timeout is enforced and what it does not stop.

## Lesson 36 — Request context: cancellation, deadlines, and request-scoped data

- Why this lesson exists now: bound how long a request may run, and hand Lesson 34's request ID downstream
- `context.Context` as request lifetime plus request-scoped metadata
- `r.Context()`
- `Done` and `Err`
- `context.Canceled` versus `context.DeadlineExceeded`
- `WithTimeout` and `WithDeadline` deriving child contexts
- `defer cancel()`
- `r.WithContext` inside timeout and request-ID middleware
- Pass context as a first parameter; do not store it on a long-lived service
- Cancellation is cooperative: a deadline stops nothing that does not check it
- `WithValue` is **not** for application data or dependencies
- `WithValue` is for request-scoped metadata that must cross layers: correlation ID, authenticated principal, tracing/span info
- Ordinary data — a customer, an order, a repository, a logger, configuration — is an ordinary parameter or dependency
- Test: is the value *about the request* or *for the work*?
- `Done` returns a channel, used here only as API surface
- **Interview:** what happens to work in flight when a client disconnects; why a context key is not dependency injection
- *Deferred:* channels and `select` → Lesson 51

## Lesson 37 — Testing the API you already have

- `testing`
- Table-driven tests
- Subtests
- Test organization
- Helpers with `t.Helper`
- Meaningful assertions without an assertion framework
- `httptest` for handlers, routes, status codes, headers, JSON contracts, validation
- Testing middleware ordering, status capture, and timeout behavior
- Mocking as a trade-off, not a default
- A fake belongs where Lesson 18's interface already provides a seam
- Unit versus integration versus end-to-end
- `go test` for the package being worked on
- `go test ./...` for the whole module, as CI will run it
- `go test -run` to isolate a failing case or subtest
- `go test -race` for data races from concurrent handlers
- `go test -cover` as a diagnostic, never as a target
- `testing.B` introduced as a tool; profiling comes later
- **Interview:** what makes an HTTP test useful, stable, and isolated; when mocking makes a test worse

---

# Phase 6 — Make the Data Real: PostgreSQL

*Backend / API Engineering · Go Language* · Lessons 38–42

The API still serves values constructed in memory. This phase replaces that, and the database is what finally forces real error design.

## Lesson 38 — A real schema and migrations

- Modelling the API's domain in PostgreSQL
- Tables, keys, constraints
- Types that matter: `timestamptz`, numeric versus float, `text`
- Why constraints belong in the database, not only in validation code
- Versioned migrations as the way schema changes ship

## Lesson 39 — `database/sql` and talking to PostgreSQL from Go

- The driver model, and why `database/sql` is an abstraction over one
- `sql.DB` as a pool, not a connection
- `QueryContext` / `QueryRowContext` / `ExecContext`
- `Scan` and its type mapping
- `NULL` and `sql.Null*`
- `defer rows.Close()` and `rows.Err()`
- Parameterized queries and SQL injection
- Prepared statements where justified
- Pool sizing and connection lifetime
- Persistence behind an interface only where something must actually vary
- A repository that only mirrors the database, or a service that forwards one call, is cargo-cult
- Hand-written SQL, `sqlx`-style helpers, query builders, `sqlc`, GORM: what each buys and costs
- Why Go leans closer to SQL than C# does to EF
- **Interview:** what `sql.DB` actually holds; why a query without a context is a liability

## Lesson 40 — Errors that cross layers

- `sql.ErrNoRows` as the first error whose identity matters
- Sentinel errors
- Custom error types carrying data
- Wrapping with `%w`
- `errors.Is` for identity
- `errors.As` for extraction
- How far to wrap
- Mapping a domain outcome to `404`, `409`, or `500` at the transport boundary
- Not leaking driver detail to clients
- **Interview:** error identity versus error text; where HTTP mapping belongs

## Lesson 41 — Transactions and correctness

- `BeginTx`
- Commit and rollback via `defer`
- Choosing a transaction boundary, usually above the repository
- Isolation levels and the anomalies each permits
- Row locking and `SELECT ... FOR UPDATE`
- Deadlocks and retry-safe transactions
- Statement timeouts
- Integration tests against a real PostgreSQL instance
- A disposable container, not a shared environment or an in-memory substitute
- **Interview:** picking an isolation level; what a transaction held open across an HTTP call does to a pool

## Lesson 42 — Query performance from the application side

- Indexes and what they cost on write
- Reading `EXPLAIN (ANALYZE, BUFFERS)`
- The N+1 pattern, and how an API's shape causes it
- Keyset versus offset scanning
- Connection-pool exhaustion presenting as latency
- Query timeouts driven by the request context
- **Interview:** diagnosing a query that is slow only in production

---

# Phase 7 — Identity and a Grown-Up API

*Backend / API Engineering · Production / Engineering* · Lessons 43–45

## Lesson 43 — Authentication

- Password hashing with a memory-hard function, and why not SHA-family
- The login flow
- Sessions versus tokens, and what each costs
- JWT structure, signing, validation pitfalls, expiry and refresh
- Secure cookie attributes
- Where OAuth2 / OIDC fit
- Identity into context via the middleware already built
- Credentials never into logs

## Lesson 44 — Authorization

- Authentication answers who; authorization answers whether
- Resource ownership checks
- Role and policy models
- Why authorization usually cannot live entirely in middleware
- Failing closed
- Avoiding object-level access flaws
- **Interview:** where an authorization decision belongs, and how it is tested

## Lesson 45 — API evolution and usability

- Resource-oriented design with real data: resource, sub-resource, action that does not fit REST
- When a non-REST endpoint is the honest answer
- Pagination: offset versus keyset, consistency under concurrent writes
- Filtering, sorting, and search
- PUT versus PATCH using the pointer-presence model from Lesson 30
- Idempotency keys stored transactionally
- A stable error contract
- Versioning and backward compatibility only when a contract must change
- Content negotiation where it earns its place
- Documenting the API with OpenAPI
- **Interview:** pagination consistency; making a payment endpoint safe to retry

---

# Phase 8 — Go Depth Where the Service Demands It

*Go Language* · Lessons 46–49

Applied to code already written, not demonstrated in isolation.

## Lesson 46 — Go's type system in depth

- Defined types versus aliases, and what each one actually creates
- Underlying types and conversion rules
- Why a `UserID` defined as an `int` is not interchangeable with every other `int`
- Sized integers, signed versus unsigned, overflow behaviour
- Why `float64` carries its width in its name (from Lesson 4)
- Why money must not be a `float64`
- Alternatives: integer minor units, or a decimal type mapped to PostgreSQL `numeric`
- Zero values as a deliberate design tool
- **Interview:** when a defined type earns its keep; what goes wrong when money is a float

## Lesson 47 — Interfaces in depth

- Deepens Lessons 15 and 18; does not re-teach implicit satisfaction
- An interface value is a `(type, value)` pair
- Typed nil: a non-nil interface holding a nil `*User`
- Pointer versus value method sets
- Why Lesson 33's wrapper had to be `*responseWriter`
- Type assertions and the comma-`ok` form
- Type switches, motivated by `any`, `errors.As`, and `Scan`
- Define a small interface at the consumer, not beside the implementation
- An interface with one implementation and no test seam is usually noise
- **Interview:** why a non-nil interface can hold a nil pointer; which method set a value versus a pointer carries

## Lesson 48 — Generics where they earn their place

- Type parameters
- Constraints
- Type inference
- Generic helpers and containers
- Generics vary code by type; interfaces vary it by behavior
- Why generics rarely belong in domain logic
- **Interview:** constraints; the absence of method-based specialization; choosing an interface over a type parameter

## Lesson 49 — Packages, standard library, and maintainable boundaries

- Package API design
- Ownership and dependency direction
- `internal`
- Constructors only where invariants or dependencies demand them
- `io` and its interfaces
- `time` and timezone handling
- `net/url`
- `os`
- `path/filepath`
- `slices`, `maps`
- `strings.Builder`
- `unicode/utf8`: decoding, runes, the deeper treatment owed from Lesson 21
- Formatting, `go vet`, linting
- Documentation comments
- Reading Go code the way a reviewer does
- **Interview:** package boundaries; when *not* to introduce an interface

---

# Phase 9 — External Systems and Concurrency

*Backend / API Engineering · Production / Engineering* · Lessons 50–54

## Lesson 50 — HTTP clients and failure-aware integration

- `http.Client` with explicit timeouts and a configured transport
- `NewRequestWithContext`
- Response-body lifecycle and connection reuse
- Handling non-2xx responses
- Retries with exponential backoff and jitter
- Retryable versus non-retryable failures
- Client-side idempotency
- Circuit breakers
- Bulkheads
- External dependencies as unreliable
- Failure behavior made observable and testable with `httptest.Server`
- Whether an external call deserves its own client type and interface
- **Interview:** retry storms; why a retry without idempotency is a bug

## Lesson 51 — Concurrency foundations

- Goroutines and their lifecycle
- Channel ownership and direction
- Buffered versus unbuffered semantics
- Closing and ranging
- `select`
- Timers
- `sync.Mutex` / `RWMutex` / `WaitGroup` / `Once`
- Atomics
- Happens-before at a practical level
- The channel underneath `ctx.Done()` from Lesson 36
- **Interview:** when a mutex is simpler and safer than a channel

## Lesson 52 — Concurrency patterns and failure modes

- Worker pools
- Bounded concurrency with semaphores
- Fan-out / fan-in
- Pipelines
- Aggregating results and errors with `errgroup`
- Propagating cancellation
- Data races
- Deadlocks
- Goroutine leaks
- Work that outlives the request that started it
- Race-detector-backed tests
- **Interview:** bounding concurrency; finding a leaking goroutine

## Lesson 53 — Background jobs and queues

- Request-time work versus durable asynchronous work
- Job payloads and schemas
- At-least-once execution
- Visibility timeouts
- Idempotent handlers
- Poison messages
- Dead-letter queues
- Scheduling
- Operational ownership
- When a PostgreSQL-backed worker suffices, and when a broker is justified

## Lesson 54 — Kafka and event-driven workflows

- Topics and partitions
- Ordering guarantees and their scope
- Consumer groups and rebalancing
- Offsets and commit strategies
- Delivery semantics and acknowledgement
- Retries and DLQs
- Deduplication
- Schema evolution
- Eventual consistency
- Why distributed transactions across a database and a broker are avoided
- Two-phase commit, sagas, and compensation as the alternatives
- Transactional outbox, using Lesson 41's transaction boundaries
- One workflow end to end: API, PostgreSQL, worker, consumer
- **Interview:** at-least-once delivery, duplicate processing, consumer lag

---

# Phase 10 — Operate It in Production

*Production / Engineering* · Lessons 55–59

## Lesson 55 — Configuration, startup, and graceful shutdown

- Typed configuration from environment and secrets
- Validation at startup so misconfiguration fails immediately
- Signal handling
- `http.Server` timeouts
- `Shutdown` and in-flight request draining
- Health, readiness, and liveness as genuinely different endpoints
- Lifecycle ownership of long-lived dependencies:
  - PostgreSQL connection pool
  - HTTP client / transport
  - Kafka consumer
  - Worker pool
  - Metrics exporter
  - Tracer provider
- Who creates each, who owns it, who closes or stops it
- Created once in `main`, passed explicitly, released by the creator
- Startup order: configuration → pools and clients → workers and consumers → HTTP server last
- Shutdown order: stop accepting work → drain requests → cancel worker/consumer contexts → close pools and clients → flush telemetry last
- Ties together `defer` (Lesson 19), context cancellation (Lesson 36), and worker/consumer loops (Lessons 52–54)
- **Interview:** what the service does between SIGTERM and exit

## Lesson 56 — Observability and production debugging

- Instrumentation middleware: request metrics and a trace span per request
- Built on the status wrapper from Lesson 33 and the correlation ID from Lesson 36
- Correlation IDs across handler, database, client, and worker
- Metrics and the RED / USE views
- SLOs, dashboards, and alerts worth waking up for
- Distributed tracing with OpenTelemetry and context propagation
- Log levels, sampling, and cost
- Concepts behind tools such as Datadog
- Following one request across every component
- **Interview:** what you actually look at when latency rises at 3am

## Lesson 57 — Reliability and security engineering

- Timeout budgets across a call chain
- Backpressure
- Overload protection and load shedding
- Caching and invalidation trade-offs
- Graceful degradation
- TLS
- Secret handling and rotation
- Input and output safety
- Common API vulnerability classes
- Dependency and supply-chain risk
- Incident-oriented runbooks
- **Interview:** containing a failing dependency instead of amplifying it

## Lesson 58 — Performance and Go runtime diagnostics

- Benchmarking a representative workload
- `pprof` for CPU, heap, and blocking profiles
- The execution tracer
- Escape analysis and allocation reduction
- GC behavior and `GOGC` / memory limits
- HTTP-level costs: JSON encoding, body copying, connection handling
- When a database or network bound makes Go-level optimization pointless
- Measure first, change second
- **Interview:** finding the bottleneck before changing code

## Lesson 59 — Delivery, deployment, and architecture review

- Go tooling in CI
- Unit / integration / contract test split
- Linting and `govulncheck`
- Multi-stage Docker builds and image hardening
- CI/CD pipelines
- Deployment strategies and rollback
- Kubernetes for a service owner: probes, resource requests and limits, configuration and secrets, autoscaling
- Architecture review by dependency direction, operational ownership, and failure behavior — not folder fashion

---

# Phase 11 — Professional Go Backend Proficiency

*All three areas* · Lessons 60–61

## Lesson 60 — Production service and design defense

- Not a new project: the service growing since Lesson 25, finished and operated
- Documented HTTP API
- PostgreSQL persistence with migrations
- Authentication and authorization
- Meaningful tests
- Context-aware external integration
- Bounded background work
- An event-driven workflow
- Structured logs, metrics, and traces
- Health endpoints
- A container image and CI/CD
- Design document covering: API contract, data model, transaction boundaries, concurrency limits, retry and idempotency strategy, observability plan, security assumptions, failure modes

## Lesson 61 — Demonstrate professional proficiency

- Defend a system design involving an API, a database, a justified cache, external dependencies, workers and queues, authentication, observability, failure handling, scaling, and deployment
- Answer the interview material accumulated across the course cold

---

# Interview and system-design coverage

Interview readiness is built lesson by lesson, not bolted on at the end. Each topic is answerable once its lesson is complete.

| Area | Topics | Covered by |
| --- | --- | --- |
| **Go** | pointers, receivers, zero values, implicit interface satisfaction | Lessons 11, 13–15, 18, 22 |
| | error values and the `err != nil` convention | Lesson 16 |
| | error wrapping, `errors.Is`/`As`, sentinel and typed errors | Lesson 40 |
| | defined types, conversion rules, numeric types, money | Lesson 46 |
| | method sets, interface values, typed nil, type switches | Lesson 47 |
| | generics and constraints | Lesson 48 |
| **HTTP** | request lifecycle, routing, status codes, JSON contracts | Lessons 25–31 |
| | middleware mechanics, ordering, status capture | Lessons 32–33 |
| | chain composition, recovery, rate limiting, writer compatibility | Lessons 34–35 |
| | request lifecycle, cancellation, deadlines, context | Lesson 36 |
| | testing handlers and middleware | Lesson 37 |
| | pagination, idempotency, versioning | Lesson 45 |
| **Databases** | indexes, query plans, slow-query diagnosis | Lesson 42 |
| | transactions, isolation, locking, boundaries | Lesson 41 |
| | connection pools and exhaustion | Lessons 39, 42 |
| **Concurrency** | goroutines, channels, `select`, mutexes | Lesson 51 |
| | worker pools, races, deadlocks, leaks, cancellation | Lesson 52 |
| **Distributed systems** | retries, duplicate processing, idempotency | Lessons 45, 50, 53 |
| | queues, Kafka, consumer groups, DLQs | Lessons 53–54 |
| | eventual consistency, outbox pattern | Lesson 54 |
| **Security** | password handling, tokens, authorization failures | Lessons 43–44, 57 |
| **Production** | observability, reliability, performance, deployment | Lessons 55–59 |
| **Architecture** | when an interface, repository, or service layer is justified | Lessons 18, 39, 49 |
| **System design** | full service design and defense | Lessons 60–61 |

---

# What proficiency means

```text
                    GO BACKEND PROFICIENCY
                             │
          ┌──────────────────┼──────────────────┐
          │                  │                  │
          ▼                  ▼                  ▼
   GO LANGUAGE        BACKEND / API       PRODUCTION /
                      ENGINEERING          ENGINEERING
          │                  │                  │
          └──────────────────┼──────────────────┘
                             │
                             ▼
                  PROFESSIONAL GO BACKEND
                             │
                 ┌───────────┴───────────┐
                 │                       │
                 ▼                       ▼
        Production Project       Interview / Design
```

## Gates along the way

Each gate is a checkpoint, not a new lesson: it is passed by demonstrating the capability on the service being built.

| Gate | Reached after | You can independently |
| --- | --- | --- |
| **1 — Go developer** | Phase 2 | Read and write idiomatic Go; use structs, methods, pointers, interfaces, slices, maps, and packages; handle errors the Go way. |
| **2 — Go API developer** | Phase 5 | Build an HTTP API end to end: routing, JSON contracts, validation, error responses, a professional middleware chain, and context propagation. |
| **3 — Professional Go backend developer** | Phase 10 | Add PostgreSQL and transactions, meaningful tests, concurrency, external integrations with retries and timeouts, messaging, observability, security, and graceful shutdown. |
| **4 — Senior-level backend competency** | Phase 11 | Reason about service boundaries, distributed systems, consistency, reliability, scalability, asynchronous architecture, and operational trade-offs — and justify the approach you chose. |

Proficiency is not "has seen all of these topics." The course is finished when the following are true without assistance:

- **Go Language** — Write and review idiomatic Go, and reason correctly about value versus pointer semantics, receivers and method sets, interface values and nil, zero values, error identity and wrapping, generics, package boundaries, and when an abstraction is not warranted.
- **Backend / API Engineering** — Design, build, and test an HTTP API with deliberate contracts; use middleware and context correctly; integrate unreliable external services; persist data in PostgreSQL with correct transaction boundaries; run bounded concurrent and background work; and build an event-driven workflow that tolerates duplicates.
- **Production / Engineering** — Configure, secure, observe, debug, profile, deploy, and evolve a service, with reliability and operational trade-offs stated explicitly rather than assumed.

Concretely, the end state is the ability to:

1. Build a real Go backend and understand every line of it.
2. Make architectural decisions and justify why an abstraction does or does not belong.
3. Reason about concurrency, including what breaks under load.
4. Work confidently with databases and external services.
5. Design reliable, evolvable APIs.
6. Handle failure deliberately — timeouts, retries, idempotency, degradation.
7. Test the system at the right level.
8. Observe and troubleshoot it in production.
9. Secure it.
10. Deploy and operate it.
11. Explain the design and its trade-offs in a senior-level interview.
