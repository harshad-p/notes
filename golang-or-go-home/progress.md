# Go Backend Engineering — Progress

`plan.md` is the curriculum. This file is the only record of completion state.

A lesson counts as complete when its concepts are understood **and** its exercise has been attempted. Seeing a topic mentioned in passing does not make it complete.

**Next up:** Lesson 34 — Middleware composition and ordering.

---

# Phase 1 — Go Foundations

### Lesson 0 — Install Go

`phase-01-go-foundations/000-install-go`

- [x] Lesson complete
- The Go toolchain
- `go version`
- Platform install paths
- `go env`
- `GOROOT` / `GOPATH` shown, taught properly in Lesson 17

### Lesson 1 — Understanding a Go project

`phase-01-go-foundations/001-understanding-a-go-project`

- [x] Lesson complete
- A folder is a module: no solution or project file
- `go mod init`
- What `go.mod` tracks
- Rough equivalence to `.csproj` / `package.json`
- `go run .` rather than `go run main.go`

### Lesson 2 — Structure of a Go program

`phase-01-go-foundations/002-structure-of-a-go-program`

- [x] Lesson complete
- `package main` and what makes a program executable
- `import`
- `func main` as the entry point
- Brace placement / automatic semicolon insertion
- `go build` versus `go run`
- `go.mod` contents

### Lesson 3 — Variables

`phase-01-go-foundations/003-go-basic-types`

- [x] Lesson complete
- `var` with an explicit type
- `:=` short declaration
- Why `:=` cannot redeclare an existing variable

### Lesson 4 — Basic types

`phase-01-go-foundations/003-go-basic-types`

- [x] Lesson complete
- Static typing
- `string`, `int`, `float64`, `bool`
- `Println` inserting spaces between arguments

### Lesson 5 — Constants

`phase-01-go-foundations/003-go-basic-types`

- [x] Lesson complete
- `const`
- Must initialize immediately
- When a constant is preferable to a variable

### Lesson 6 — Functions

`phase-01-go-foundations/004-go-functions`

- [x] Lesson complete
- `func`
- `name type` parameter order versus C#'s `type name`
- Return types
- Grouping parameters of the same type

### Lesson 7 — Multiple return values

`phase-01-go-foundations/004-go-functions`

- [x] Lesson complete
- Returning `(string, int)`
- Destructuring at the call site
- `(value, error)` previewed
- `error` and `nil` used here, taught in Lesson 16

### Lesson 8 — If/else and loops

`phase-01-go-foundations/005-if-else-and-loops`

- [x] Lesson complete
- `if` / `else` without parentheses, mandatory braces
- `for` as the only loop keyword
- Counted, condition-only, and infinite forms
- `%`

### Lesson 9 — Slices

`phase-01-go-foundations/006-arrays-and-slices`

- [x] Lesson complete
- Slice literals and indexing
- `append` and why its result must be reassigned
- `len`
- `range` with index and value
- The blank identifier `_`
- Unused variables do not compile

### Lesson 10 — Arrays versus slices

`phase-01-go-foundations/006-arrays-and-slices`

- [x] Lesson complete
- `[3]int` versus `[]int`
- When an array is the right choice
- Slice header: pointer, length, capacity
- Backing array on growth
- Empty-slice-then-append for API responses

### Lesson 11 — Maps

`phase-01-go-foundations/007-maps`

- [x] Lesson complete
- Map literals
- Reading, add, and update by assignment
- `delete`
- Zero values: a missing key returns the zero value
- Comma-`ok` to distinguish absence from a zero value

### Lesson 12 — Structs

`phase-01-go-foundations/008-structs`

- [x] Lesson complete
- `type ... struct`
- Literals with field names
- Field access and mutation
- Capitalization as visibility, taught fully in Lesson 17

### Lesson 13 — Methods and receivers

`phase-01-go-foundations/008-structs`

- [x] Lesson complete
- Methods declared outside the type
- The receiver
- A value receiver operates on a copy

### Lesson 14 — Pointers

`phase-01-go-foundations/008-structs`

- [x] Lesson complete
- Arguments passed by copy
- `&` for address-of
- `*T` as a pointer type
- Pointer receivers for mutation
- Automatic address-taking on a method call

### Lesson 15 — Interfaces

`phase-01-go-foundations/008-structs`

- [x] Lesson complete
- Implicit satisfaction: no `implements` keyword
- Depending on behavior rather than a concrete type
- Two types satisfying one interface

### Lesson 16 — Error handling

`phase-01-go-foundations/009-errors`

- [x] Lesson complete
- `(value, error)` return convention
- `nil` as "no error"
- `if err != nil`
- `error` as an interface with `Error() string`
- Errors are values
- `errors.New` versus `fmt.Errorf`
- Early-return error path

---

# Phase 2 — Organizing Go Code and Connecting Language Features

### Lesson 17 — Packages and modules

`phase-02-organizing-go-code/010-packages-and-modules`

- [x] Lesson complete
- Several files in one package
- A second package in a subdirectory
- Capitalization as the export rule
- Module path as the root of import paths
- What `go mod init` establishes

### Lesson 18 — A small service boundary

`phase-02-organizing-go-code/011-pointers-structs-interfaces`

- [x] Lesson complete
- A struct model with behavior
- An interface describing only what a consumer needs
- An in-memory repository satisfying it implicitly
- A service holding the interface as a dependency
- Explicit wiring in `main`
- A test fake as the motivation for the interface
- Composition instead of inheritance

### Lesson 19 — Cleanup with `defer`

`phase-02-organizing-go-code/012-defer`

- [x] Lesson complete
- Deferred execution at function exit
- Cleanup that survives early returns
- LIFO ordering
- Arguments evaluated at `defer` time, not call time

### Lesson 20 — Dependencies

`phase-02-organizing-go-code/013-external-dependencies`

- [x] Lesson complete
- `go get`
- Importing a third-party package
- `go.mod` and `go.sum`
- `go mod tidy`

### Lesson 21 — Strings in Go

`phase-02-organizing-go-code/014-strings`

- [x] Lesson complete
- Concatenation
- A string is a sequence of bytes holding UTF-8 text
- `len` counts bytes, not characters
- Indexing yields a `byte`
- `fmt.Printf` and verbs
- The `strings` package
- `rune` as a Unicode code point (bridge only; decoding in Lesson 49)

### Lesson 22 — Pointer syntax in practice

`phase-02-organizing-go-code/015-pointers-part-2`

- [x] Lesson complete
- `&`
- `*` in a type versus `*` as dereference
- Mutation through a pointer
- Automatic dereferencing for field and method access

### Lesson 23 — Composition

`phase-02-organizing-go-code/016-struct-embedding-and-composition`

- [x] Lesson complete
- Struct embedding
- Field and method promotion
- Embedding is "has-a" composition, not inheritance

### Lesson 24 — JSON at the model boundary

`phase-02-organizing-go-code/017-struct-tags-and-json`

- [x] Lesson complete
- Struct tags as field metadata
- `encoding/json`
- `Marshal` producing `[]byte`
- `Unmarshal` requiring a pointer
- `omitempty`

---

# Phase 3 — Building the HTTP API

### Lesson 25 — HTTP request and response bodies

`phase-03-building-the-http-api/018-JSON-in-HTTP-APIs`

- [x] Lesson complete
- Shape of an HTTP request
- `net/http` and `ListenAndServe`
- The handler signature
- `*http.Request` and `http.ResponseWriter`
- `fmt.Fprintln`, `w.Header().Set`, `r.Body`
- `json.NewEncoder` / `NewDecoder`
- `http.Error` and `400` for malformed JSON

### Lesson 26 — Methods, routing, status codes

`phase-03-building-the-http-api/019-http-routing-methods-and-status-codes`

- [x] Lesson complete
- `r.Method` and method constants
- `switch` for method dispatch
- `HandleFunc` routing
- Routing versus method
- `WriteHeader`, `201`, `405`, implicit `200`

### Lesson 27 — URL input

`phase-03-building-the-http-api/020-url-paths-query-parameters-and-path-parameters`

- [x] Lesson complete
- URL anatomy
- `r.URL.Path` and `strings.Split`
- `r.URL.Query()` and `.Get`
- Path parameter versus query parameter

### Lesson 28 — `ServeMux` route patterns

`phase-03-building-the-http-api/021-http-router`

- [x] Lesson complete
- `http.NewServeMux`
- Method-qualified patterns
- `{id}` wildcards and `r.PathValue`
- `404` versus `405`

### Lesson 29 — Turning URL strings into Go values

`phase-03-building-the-http-api/022-type-conversion-and-parsing`

- [x] Lesson complete
- Conversion versus parsing
- `strconv.Atoi`, `ParseBool`, `ParseFloat`, `ParseInt`
- Parsing failure as `400`
- Absent versus invalid query parameters

### Lesson 30 — API contracts and validation

`phase-03-building-the-http-api/023-json-api-request-response-model`

- [x] Lesson complete
- Separate request / response / persistence types when contracts differ
- Decoding versus validation
- Input validation versus business validation
- Pointer fields for PATCH presence

### Lesson 31 — Consistent JSON responses

`phase-03-building-the-http-api/024-consistent-api-responses`

- [x] Lesson complete
- Headers → status → body
- JSON error shape, `writeError`, `writeJSON`
- `any`
- `204 No Content` must not write a body

---

# Phase 4 — Middleware Mechanics

### Lesson 32 — Middleware mechanics

`phase-04-middleware-mechanics/025-middleware`

- [x] Lesson complete
- `http.Handler` and `http.HandlerFunc`
- `func(next http.Handler) http.Handler`
- `next.ServeHTTP` and short-circuiting
- Before/after execution
- `time.Now` / `time.Since` and `log.Printf`
- Wrapping the router versus selected handlers
- Nesting, composition, and ordering

### Lesson 33 — Response observation

`phase-04-middleware-mechanics/026-middleware-in-depth`

- [x] Lesson complete
- Decorating `http.ResponseWriter` by embedding the interface
- Intercepting `WriteHeader` and `Write`
- Implicit `200` and first-status-wins
- Pointer receiver and method-set consequence
- Status-and-duration logging
- Naive wrappers hide `Flusher` / `Hijacker` / `ReaderFrom`

---

# Phase 5 — Finish the HTTP Boundary Properly

### Lesson 34 — Middleware composition and ordering

- [ ] Lesson complete
- Chain helper replacing nested calls
- Reverse-wrapping order
- Build the chain once at startup
- Panic recovery after a committed response
- Request-scoped `log/slog`
- Generate and log a request ID, without handing it downstream yet
- What must not appear in request logs

### Lesson 35 — Middleware at the edge of production

- [ ] Lesson complete
- CORS and why a permissive policy is a vulnerability
- Security headers
- Rate limiting (token bucket, per-key versus global, multi-instance)
- Request-size and body limits
- Preserving `Flusher` / `Hijacker` / `ReaderFrom`
- How routers and frameworks express middleware
- What belongs in middleware versus the handler

### Lesson 36 — Request context: cancellation, deadlines, and request-scoped data

- [ ] Lesson complete
- `context.Context`, `r.Context()`, `Done`, `Err`
- `Canceled` versus `DeadlineExceeded`
- `WithTimeout`, `WithDeadline`, `defer cancel()`, `r.WithContext`
- Do not store context on a long-lived service
- Cancellation is cooperative
- `WithValue` for request-scoped metadata only, not dependencies
- Hand the request ID from Lesson 34 to downstream code

### Lesson 37 — Testing the API you already have

- [ ] Lesson complete
- `testing`, table-driven tests, subtests, `t.Helper`
- `httptest` for handlers, contracts, and middleware
- Unit versus integration versus end-to-end
- Mocking as a trade-off
- `go test`, `go test ./...`, `-run`, `-race`, `-cover`
- `testing.B`

---

# Phase 6 — Make the Data Real: PostgreSQL

### Lesson 38 — A real schema and migrations

- [ ] Lesson complete
- Tables, keys, constraints
- `timestamptz`, numeric versus float, `text`
- Constraints in the database
- Versioned migrations

### Lesson 39 — `database/sql` and talking to PostgreSQL from Go

- [ ] Lesson complete
- Drivers and `sql.DB` as a pool
- `QueryContext` / `QueryRowContext` / `ExecContext`
- `Scan`, `NULL`, `sql.Null*`
- `rows.Close()`, parameterized queries, prepared statements
- Pool sizing
- When a repository/interface is justified
- Hand-written SQL versus `sqlx` / `sqlc` / GORM

### Lesson 40 — Errors that cross layers

- [ ] Lesson complete
- `sql.ErrNoRows`
- Sentinel errors and custom error types
- Wrapping with `%w`
- `errors.Is` and `errors.As`
- Mapping domain outcomes to HTTP status without leaking driver detail

### Lesson 41 — Transactions and correctness

- [ ] Lesson complete
- `BeginTx`, commit, rollback
- Transaction boundaries
- Isolation levels
- `SELECT ... FOR UPDATE`, deadlocks, retries
- Statement timeouts
- Integration tests against real PostgreSQL

### Lesson 42 — Query performance from the application side

- [ ] Lesson complete
- Indexes and write cost
- `EXPLAIN (ANALYZE, BUFFERS)`
- N+1
- Keyset versus offset
- Pool exhaustion presenting as latency
- Query timeouts from request context

---

# Phase 7 — Identity and a Grown-Up API

### Lesson 43 — Authentication

- [ ] Lesson complete
- Memory-hard password hashing
- Login flow
- Sessions versus tokens
- JWT structure, signing, validation, expiry/refresh
- Secure cookies
- OAuth2 / OIDC fit
- Identity into context; credentials never into logs

### Lesson 44 — Authorization

- [ ] Lesson complete
- Who versus whether
- Resource ownership, roles, policies
- Why authorization cannot live entirely in middleware
- Failing closed
- Object-level access flaws

### Lesson 45 — API evolution and usability

- [ ] Lesson complete
- Resource, sub-resource, and non-REST actions
- Pagination: offset versus keyset
- Filtering, sorting, search
- PUT versus PATCH
- Idempotency keys
- Stable error contract
- Versioning and backward compatibility
- Content negotiation
- OpenAPI

---

# Phase 8 — Go Depth Where the Service Demands It

### Lesson 46 — Go's type system in depth

- [ ] Lesson complete
- Defined types versus aliases
- Underlying types and conversion rules
- Numeric types and why `float64` includes its width
- Why money is not a `float64`
- Zero values as a design tool

### Lesson 47 — Interfaces in depth

- [ ] Lesson complete
- Interface values as `(type, value)`
- Typed nil
- Method sets; pointer versus value
- Type assertions and type switches
- Small interfaces defined at the consumer

### Lesson 48 — Generics where they earn their place

- [ ] Lesson complete
- Type parameters, constraints, type inference
- Generic helpers and containers
- Generics versus interfaces
- Why generics rarely belong in domain logic

### Lesson 49 — Packages, standard library, and maintainable boundaries

- [ ] Lesson complete
- Package API design, ownership, `internal`
- Constructors only when needed
- `io`, `time`, `net/url`, `os`, `path/filepath`, `slices`, `maps`, `strings.Builder`
- `unicode/utf8` (deeper rune/UTF-8 treatment)
- Formatting, `go vet`, linting, documentation comments

---

# Phase 9 — External Systems and Concurrency

### Lesson 50 — HTTP clients and failure-aware integration

- [ ] Lesson complete
- `http.Client` timeouts and transport
- `NewRequestWithContext`
- Response-body lifecycle
- Retries, backoff, jitter, idempotency
- Circuit breakers and bulkheads
- `httptest.Server`

### Lesson 51 — Concurrency foundations

- [ ] Lesson complete
- Goroutines
- Channels: ownership, direction, buffered versus unbuffered, close, range
- `select`, timers
- `Mutex` / `RWMutex` / `WaitGroup` / `Once`, atomics
- Happens-before
- The channel under `ctx.Done()`

### Lesson 52 — Concurrency patterns and failure modes

- [ ] Lesson complete
- Worker pools, bounded concurrency, fan-out/fan-in, pipelines
- `errgroup`
- Data races, deadlocks, goroutine leaks
- Work that outlives the request

### Lesson 53 — Background jobs and queues

- [ ] Lesson complete
- Request-time versus durable async work
- Payloads, at-least-once, visibility timeouts
- Idempotent handlers, poison messages, DLQs
- Scheduling and operational ownership
- PostgreSQL worker versus a broker

### Lesson 54 — Kafka and event-driven workflows

- [ ] Lesson complete
- Topics, partitions, ordering
- Consumer groups, offsets, delivery semantics
- Retries, DLQs, deduplication, schema evolution
- Eventual consistency
- Why not distributed transactions; sagas and compensation
- Transactional outbox
- One workflow: API, PostgreSQL, worker, consumer

---

# Phase 10 — Operate It in Production

### Lesson 55 — Configuration, startup, and graceful shutdown

- [ ] Lesson complete
- Typed config and startup validation
- Signals, `http.Server` timeouts, `Shutdown`
- Health / readiness / liveness
- Lifecycle ownership of pool, HTTP client, Kafka consumer, workers, metrics, tracer
- Startup order and shutdown order

### Lesson 56 — Observability and production debugging

- [ ] Lesson complete
- Metrics and tracing middleware
- Correlation IDs across components
- RED / USE, SLOs, dashboards, alerts
- OpenTelemetry and context propagation
- Log levels, sampling, cost

### Lesson 57 — Reliability and security engineering

- [ ] Lesson complete
- Timeout budgets, backpressure, load shedding
- Caching and graceful degradation
- TLS, secrets, input/output safety
- API vulnerability classes
- Supply-chain risk and runbooks

### Lesson 58 — Performance and Go runtime diagnostics

- [ ] Lesson complete
- Benchmarking a representative workload
- `pprof`, execution tracer, escape analysis
- GC and memory limits
- HTTP-level costs
- Measure first, change second

### Lesson 59 — Delivery, deployment, and architecture review

- [ ] Lesson complete
- CI tooling, test split, linting, `govulncheck`
- Docker and CI/CD
- Kubernetes probes, resources, config, autoscaling
- Architecture review by dependency direction and failure behavior

---

# Phase 11 — Professional Go Backend Proficiency

### Lesson 60 — Production service and design defense

- [ ] Lesson complete
- Finish and operate the service started in Lesson 25
- Design document: contract, data model, transactions, concurrency, retries, observability, security, failure modes

### Lesson 61 — Demonstrate professional proficiency

- [ ] Lesson complete
- Defend a full system design
- Answer the interview material accumulated across the course
