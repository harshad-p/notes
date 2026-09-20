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

- Lessons that have been written name the repository folder holding them and their code, and describe what was *actually taught* there — including language features introduced in the middle of backend work. Later lessons describe what they will cover.
- Some folders hold several lessons. The lesson sequence, not the folder name, is the authoritative order.
- Each phase continues from the state the service is actually in. A topic already taught is never rescheduled; where a later lesson revisits a subject, it names the deeper concern it is there to solve.
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
| PostgreSQL, `database/sql` | Illustrative future example in Lessons 18 and 30 | Lessons 38–39 |
| Test fakes | Motivated the interface in Lesson 18 | Lesson 37 |
| Optional `ResponseWriter` interfaces | Named as a caveat in Lesson 33 | Lesson 35 |
| Runes, UTF-8 internals | Named in Lesson 21 without explanation; the minimal bridge is owed there | Lesson 49 |
| Go's full numeric types | `float64` naming deferred in Lesson 4 | Lesson 46 |
| `GOROOT` / `GOPATH` | Shown in Lesson 0 | Covered enough by Lesson 17 |
| `context.Context` | Drafted, then pulled back out as premature | Lesson 36 |

## How the course is taught

- **One service, carried forward.** The product API started in Lesson 25 is the same codebase that gains routing, validation, response contracts, and middleware. Every future phase extends it rather than starting a new toy. Lesson 60 is that service reaching production shape, not a new project.
- **Backend need first, language feature second.** A Go concept is introduced when the service needs it, then deepened later only when the deeper treatment is genuinely different.
- **No architecture by default.** Service layers, repositories, interfaces, factories, and folder conventions appear only when a concrete problem — a dependency to swap, a boundary to defend, a lifecycle to own, a test seam to create — justifies them, and the trade-off is stated. Lesson 18 introduced an interface because a test fake needed one; that is the standard every later abstraction must meet.
- **Depth without narration.** Lessons explain mechanisms, edge cases, and design consequences for an experienced developer. They do not walk through obvious code line by line.
- **One exercise per lesson**, realistic, with no solution or step-by-step hints in the statement.
- **Lesson shape.** What it is → why it exists → the Go mechanics → applied to our service → how it is used in production → alternatives and trade-offs → the mistakes that actually cause incidents → exercise.
- **Interview checkpoints.** Each phase ends with the learner answering that phase's interview questions before any answer is given. The per-lesson *Interview focus* notes below are the source of those questions.
- **Completion gates.** A lesson is marked done in `progress.md` when the concept is understood at a level usable at work — not when it has been read.

---

# Phase 1 — Go Foundations

*Go Language* · Lessons 0–16

### `000-install-go` — Lesson 0: Install Go

The toolchain, `go version`, platform install paths, and a first look at `go env`. `GOROOT`/`GOPATH` are shown and explicitly deferred to the modules lesson.

### `001-understanding-a-go-project` — Lesson 1: Understanding a Go project

A folder is a module: no solution or project file. `go mod init`, what `go.mod` tracks, its rough equivalence to `.csproj`/`package.json`, and the habit of running `go run .` rather than `go run main.go`.

### `002-structure-of-a-go-program` — Lesson 2: Structure of a Go program

`package main` and what makes it executable, `import`, `func main` as the entry point, and Go's brace placement rule as a consequence of automatic semicolon insertion. Also covers `go build` versus `go run` and the `go.mod` contents.

### `003-go-basic-types` — Lessons 3–5: Variables, basic types, constants

**Lesson 3 — Variables:** `var` with an explicit type, `:=` short declaration, and why `:=` cannot redeclare an existing variable.
**Lesson 4 — Basic types:** static typing, `string`, `int`, `float64`, `bool`, and `Println` inserting spaces between arguments. Why the type is named `float64` is deferred; `fmt.Print`/`Printf` are deferred.
**Lesson 5 — Constants:** `const`, the requirement to initialize immediately, and when a constant is preferable to a variable.

### `004-go-functions` — Lessons 6–7: Functions and multiple returns

**Lesson 6 — Functions:** `func`, the `name type` parameter order versus C#'s `type name`, return types, and grouping parameters of the same type.
**Lesson 7 — Multiple return values:** returning `(string, int)`, destructuring at the call site, and the `(value, error)` convention previewed as Go's replacement for exception-based flow. `error` and `nil` are used here but explicitly deferred to Lesson 16.

### `005-if-else-and-loops` — Lesson 8: If/else and loops

`if`/`else` without parentheses but with mandatory braces; `for` as the only loop keyword, in its counted, condition-only (`while`), and infinite forms; `%` for the even/odd exercise. `switch` is deliberately absent — it arrives in Lesson 26.

### `006-arrays-and-slices` — Lessons 9–10: Slices, then arrays

**Lesson 9 — Slices:** slice literals, indexing, `append` and why its result must be reassigned, `len`, `range` with index and value, the blank identifier `_`, and Go's refusal to compile unused variables.
**Lesson 10 — Arrays versus slices:** fixed-size `[3]int` versus `[]int`, when an array is the right choice, the slice header as pointer/length/capacity, what happens to the backing array on growth, and the empty-slice-then-append pattern used to build API responses.

### `007-maps` — Lesson 11: Maps

Map literals, reading, add/update through plain assignment, `delete`, **zero values introduced here** (a missing key returns the value type's zero value), and the comma-`ok` form as the way to distinguish absence from a zero value.

### `008-structs` — Lessons 12–15: Structs, methods, pointers, interfaces

**Lesson 12 — Structs:** `type ... struct`, literals with field names, field access and mutation, and capitalization as the visibility mechanism (noted, fully explained in Lesson 17).
**Lesson 13 — Methods and receivers:** methods declared outside the type, the receiver, and a value receiver operating on a copy.
**Lesson 14 — Pointers:** arguments passed by copy, `&` for address-of, `*T` as a pointer type, pointer receivers for mutation, and Go taking the address automatically on a method call.
**Lesson 15 — Interfaces:** implicit satisfaction with no `implements` keyword, depending on behavior rather than a concrete type, and two types satisfying one interface.

### `009-errors` — Lesson 16: Error handling

The `(value, error)` return convention, `nil` as "no error" and the `if err != nil` shape, **`error` as an ordinary interface with `Error() string`** so errors are values, `errors.New` versus `fmt.Errorf` with formatted values, and the early-return error path. Wrapping, sentinel errors, and custom error types are explicitly deferred.

---

# Phase 2 — Organizing Go Code and Connecting Language Features

*Go Language · Backend / API Engineering* · Lessons 17–24

### Lesson 17 / `010-packages-and-modules` — Packages and modules

Several files in one package, a second package in a subdirectory, capitalization as the export rule instead of access modifiers, the module path as the root of import paths, and what `go mod init` establishes.

### Lesson 18 / `011-pointers-structs-interfaces` — A small service boundary

The first backend-shaped program: a struct model with behavior, an interface describing only what a consumer needs, an in-memory repository satisfying it implicitly, a service holding the interface as a dependency, explicit wiring in `main` instead of a DI container, a test fake as the motivation for the interface, and composition instead of inheritance. This is interface *use*, not a complete interface-design treatment.

### Lesson 19 / `012-defer` — Cleanup with `defer`

Deferred execution at function exit, cleanup that survives early returns, LIFO ordering, and argument evaluation at `defer` time rather than call time.

### Lesson 20 / `013-external-dependencies` — Dependencies

`go get`, importing a third-party package, the roles of `go.mod` and `go.sum`, and `go mod tidy`.

### Lesson 21 / `014-strings` — Strings in Go

Concatenation; a Go string as a sequence of bytes holding UTF-8-encoded text; `len` counting those bytes rather than human-readable characters; indexing yielding a single `byte`; **`fmt.Printf` and verbs introduced here** via `%c`; and the `strings` package.

**Top-up owed here.** The lesson names `rune` without explaining it, which leaves a gap rather than a deferral. The bridge belongs in this lesson: a `rune` is a Unicode code point, a character outside ASCII occupies several bytes, so indexing or slicing by byte can split one apart, and `range` over a string decodes runes rather than walking bytes. That is enough to know when byte indexing is the wrong tool. Decoding, normalization, and the `unicode/utf8` API stay in Lesson 49.

### Lesson 22 / `015-pointers-part-2` — Pointer syntax in practice

`&`, `*` in a type versus `*` as dereference, mutation through a pointer, and automatic dereferencing for field and method access. Deepens Lesson 14 at the syntax level.

### Lesson 23 / `016-struct-embedding-and-composition` — Composition

Struct embedding, field and method promotion, and why embedding expresses "has-a" composition rather than inheritance.

### Lesson 24 / `017-struct-tags-and-json` — JSON at the model boundary

Struct tags as field metadata, `encoding/json`, `Marshal` producing `[]byte`, `Unmarshal` requiring a pointer so it can populate the target, and `omitempty`. JSON arrives here as a language/stdlib feature; Lesson 25 applies it as API transport.

---

# Phase 3 — Building the HTTP API

*Backend / API Engineering · Go Language* · Lessons 25–31

### Lesson 25 / `018-JSON-in-HTTP-APIs` — HTTP request and response bodies

The shape of an HTTP request (method, URL, headers, body), `net/http`, `ListenAndServe`, the handler signature, what `*http.Request` exposes, `http.ResponseWriter` as the response-construction mechanism, `fmt.Fprintln` writing to a destination rather than stdout, setting response headers with `w.Header().Set`, `r.Body` as a stream, `json.NewEncoder`/`NewDecoder` versus in-memory `Marshal`/`Unmarshal`, `http.Error`, and returning `400` for malformed JSON.

### Lesson 26 / `019-http-routing-methods-and-status-codes` — Methods, routing, status codes

`r.Method` and the `http.MethodGet`-style constants; **`switch` is introduced here**, in the context of method dispatch, and contrasted with C#'s fallthrough rules; manual route registration with `HandleFunc`; the separation of "which handler" (routing) from "which operation" (method); `WriteHeader` and the fact that the status is committed once; `201 Created`; `405 Method Not Allowed`; and the implicit `200` when a handler writes a body without selecting a status.

### Lesson 27 / `020-url-paths-query-parameters-and-path-parameters` — URL input

URL anatomy, `r.URL.Path`, extracting a path segment manually with `strings.Split` and why that does not scale, `r.URL.Query()` and `.Get`, an absent parameter yielding the empty string, multiple filters, and the design distinction between a path parameter identifying a resource and a query parameter modifying retrieval. Automatic path extraction is deliberately withheld until the next lesson.

### Lesson 28 / `021-http-router` — `ServeMux` route patterns

`http.NewServeMux` and passing an explicit router to the server, method-qualified patterns such as `GET /users`, `{id}` wildcards including several in one pattern, `r.PathValue`, why this beats hand-splitting paths, and how the router distinguishes `404 Not Found` from `405 Method Not Allowed`.

### Lesson 29 / `022-type-conversion-and-parsing` — Turning URL strings into Go values

Type conversion between compatible Go types versus parsing text into a value; `strconv.Atoi`, `ParseBool`, `ParseFloat`, and when `ParseInt` is warranted; parsing failure as a `400`; the distinction between a malformed value and a well-formed but unacceptable one; and treating an absent query parameter differently from an invalid one by supplying a default.

### Lesson 30 / `023-json-api-request-response-model` — API contracts and validation

Why request, response, and persistence representations are separate types when their contracts differ, and why not to multiply types when they do not; decoding answering "is this JSON understandable" while validation answers "is this data acceptable"; input validation versus business validation; `400 Bad Request` for both malformed and invalid input; and pointer fields as the way a PATCH request model distinguishes "field absent" from "field set to its zero value".

### Lesson 31 / `024-consistent-api-responses` — Consistent JSON responses

The headers → status → body ordering rule and why writing the body commits the response; choosing success statuses per operation; a JSON error shape and a `writeError` helper; a `writeJSON` helper — which is where **`any`, the alias for the empty interface, is introduced**; why `204 No Content` must not be routed through a body-writing helper; and keeping HTTP status semantics out of domain models.

---

# Phase 4 — Middleware Mechanics

*Backend / API Engineering · Go Language* · Lessons 32–33

### Lesson 32 / `025-middleware` — Middleware mechanics

`http.Handler` as an interface and `http.HandlerFunc` as the function adapter that satisfies it; the `func(next http.Handler) http.Handler` shape; delegating with `next.ServeHTTP` and what happens when you do not; before/after execution around a synchronous downstream call; **`time.Now`/`time.Since` and `log.Printf` are introduced here** for request timing; wrapping the whole router versus wrapping selected handlers; nesting, composition, and why ordering changes behavior; short-circuiting the chain; and the boundary between cross-cutting middleware and endpoint logic. Status capture is explicitly postponed.

### Lesson 33 / `026-middleware-in-depth` — Response observation

`http.ResponseWriter` as an interface with no way to read back the status; decorating it with a struct that **embeds the interface** so unoverridden methods still reach the underlying writer; intercepting `WriteHeader` to record the status and forwarding it so the client still receives it; overriding `Write` to model the implicit `200`; tracking `wroteHeader` so the first status wins; why the wrapper needs a pointer receiver, and that this is about mutation rather than about implementing the interface — with the method-set consequence that `*responseWriter`, not `responseWriter`, is the implementing type; status-and-duration request logging; and the production caveat that a naive wrapper hides `http.Flusher`, `http.Hijacker`, and `io.ReaderFrom`.

---

# Phase 5 — Finish the HTTP Boundary Properly

*Backend / API Engineering · Production / Engineering* · Lessons 34–37

The chain built in Lessons 32–33 is still hand-nested, has no recovery, cannot pass anything to a handler, and has nothing verifying that it behaves. This phase takes middleware from working to professional, introduces the request lifecycle the finished chain needs, and makes the existing API testable before it grows.

#### The middleware track

Middleware is not finished when `func(next http.Handler) http.Handler` works. It is a thread running through the rest of the course, and each concern lands where its prerequisites exist:

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

### Lesson 34 — Middleware composition and ordering

A chain helper replacing nested calls, the reverse-wrapping order that makes a chain read top-to-bottom, building the chain once at startup rather than per request, and a deliberate order for the pieces the API is missing: panic recovery that must cope with a response already committed, and request-scoped structured logging with `log/slog` replacing `log.Printf`. A request ID can be generated and logged here, but not yet handed to the handler — that gap is what Lesson 36 solves. Also what a request log must not contain: credentials, tokens, personal data, and whole request or response bodies. *Interview focus: why ordering changes behavior, and what recovery can still do after the status has been written.*

### Lesson 35 — Middleware at the edge of production

The cross-cutting concerns a public API cannot ship without: CORS and why a permissive policy is a vulnerability, security headers, rate limiting (token bucket, per-key versus global, and what breaks across multiple instances), and request-size and body limits. Returns to `ResponseWriter` for the exact problem Lesson 33 named — implementing a wrapper that preserves `http.Flusher`, `http.Hijacker`, and `io.ReaderFrom` so streaming, upgrades, and efficient copies survive it. Surveys how the standard library, common routers, and full frameworks express middleware, so an unfamiliar production codebase is readable and the cost of adopting one is clear. Finally, the boundary question with a real case rather than a slogan: a rule that depends on the resource being acted upon cannot live in middleware, which is why authorization is only partly a middleware concern when it arrives in Lesson 44. *Interview focus: implementing rate limiting, why wrapping the writer can break a dependency, and what belongs in middleware versus the handler.*

**Phase 5 interview checkpoint (answered before any answers are given):** what middleware is and how a chain is constructed; why order changes behaviour; how request logging and status capture are implemented; what happens when `WriteHeader` is called twice; why wrapping `http.ResponseWriter` can be problematic and how to do it safely; how authentication middleware would be structured; how a request timeout is enforced and what it does not stop.

### Lesson 36 — Request context: cancellation, deadlines, and request-scoped data

With the chain complete, two things it still cannot do motivate this lesson: bound how long a request may run, and hand the request ID from Lesson 34 to the code downstream. `context.Context` as request lifetime plus request-scoped metadata; `r.Context()`; `Done` and `Err`, and `context.Canceled` versus `context.DeadlineExceeded`; `WithTimeout` and `WithDeadline` deriving child contexts; `defer cancel()`; `r.WithContext` inside timeout and request-ID middleware; passing context explicitly as a first parameter rather than storing it on a long-lived service, because a context describes one operation while a service outlives thousands; and cancellation as cooperative — a deadline stops nothing that does not check it.

`WithValue` gets an explicit rule, because misusing it is the most common way Go codebases rot. It is **not** a mechanism for passing application data or dependencies down the call stack. It is for request-scoped metadata that belongs to the request's lifecycle and genuinely has to cross layers that do not otherwise care about it: the correlation ID from Lesson 34, the authenticated principal once Lesson 43 introduces one, and tracing/span information when Lesson 56 adds it. Everything else — a customer, an order, a repository, a logger, a configuration value — is an ordinary parameter or an ordinary dependency, and hiding it in a context only removes it from the function signature where a reader would look for it. The test is whether the value is *about the request* or *for the work*. `Done` returns a channel, used here only as API surface; channels and `select` are deferred to Lesson 51. *Interview focus: what actually happens to work in flight when a client disconnects, and why a context key is not dependency injection.*

### Lesson 37 — Testing the API you already have

`testing`, table-driven tests, subtests, test organization, helpers with `t.Helper`, and meaningful assertions without pulling in an assertion framework; `httptest` for exercising handlers, routes, status codes, headers, JSON contracts, validation, and — the reason this lesson lands here — middleware ordering, status capture, and timeout behavior, none of which can be verified by reading the code. Mocking is treated as a trade-off, not a default: a fake belongs where Lesson 18's interface already provides a seam, and a test that mocks its way to a green result while asserting nothing about behavior is worse than no test. Establishes the unit/integration/end-to-end split the rest of the course follows.

The workflow comes with the tooling, learned by using it rather than as a command reference: `go test` for the package being worked on, `go test ./...` for the whole module as CI will run it, `go test -run` to isolate one failing case or subtest while fixing it, `go test -race` to expose the data races that an HTTP server's concurrent handlers will eventually produce, and `go test -cover` to find untested paths — treated as a diagnostic, never as a target. `testing.B` is introduced as a tool here; profiling comes much later. *Interview focus: what makes an HTTP test useful, stable, and isolated, and when mocking makes a test worse.*

---

# Phase 6 — Make the Data Real: PostgreSQL

*Backend / API Engineering · Go Language* · Lessons 38–42

The API still serves values constructed in memory. This phase replaces that, and the database is what finally forces real error design.

### Lesson 38 — A real schema and migrations

Modelling the API's domain in PostgreSQL: tables, keys, constraints, and the types that matter (`timestamptz`, numeric versus float, `text`); why constraints belong in the database rather than only in validation code; and versioned migrations as the way schema changes ship.

### Lesson 39 — `database/sql` and talking to PostgreSQL from Go

The driver model and why `database/sql` is an abstraction over one, `sql.DB` as a pool rather than a connection, `QueryContext`/`QueryRowContext`/`ExecContext` — the payoff for the context work in Lesson 36 — `Scan` and its type mapping, `NULL` and `sql.Null*`, `defer rows.Close()` and `rows.Err()`, parameterized queries and SQL injection, prepared statements where justified, and pool sizing and connection lifetime. Persistence goes behind an interface only where something must actually vary — the test seam Lesson 18 already established, or a second implementation that genuinely exists. A repository layer whose only job is to mirror the database, or a service layer that forwards one call, is named as the cargo cult it is. Closes with the ecosystem question so the choice is informed rather than inherited: hand-written SQL, `sqlx`-style helpers, query builders, code generation such as `sqlc`, and full ORMs such as GORM — what each buys, what each costs, and why Go leans closer to SQL than C# does to EF. *Interview focus: what `sql.DB` actually holds, and why a query without a context is a liability.*

### Lesson 40 — Errors that cross layers

`sql.ErrNoRows` is the first error whose *identity* matters, which makes this the natural point for the rest of Go's error model: sentinel errors, custom error types carrying data, wrapping with `%w`, `errors.Is` for identity and `errors.As` for extraction, choosing how far to wrap, and mapping a domain outcome to `404`, `409`, or `500` at the transport boundary without leaking driver detail to clients. *Interview focus: error identity versus error text, and where HTTP mapping belongs.*

### Lesson 41 — Transactions and correctness

`BeginTx`, commit and rollback via `defer`, choosing a transaction boundary and why it usually belongs above the repository, isolation levels and the anomalies each permits, row locking and `SELECT ... FOR UPDATE`, deadlocks and retry-safe transactions, statement timeouts, and integration tests running against a real PostgreSQL instance — a disposable container rather than a shared environment or an in-memory substitute that does not share the real engine's semantics. *Interview focus: picking an isolation level, and what a transaction held open across an HTTP call does to a pool.*

### Lesson 42 — Query performance from the application side

Indexes and what they cost on write, reading `EXPLAIN (ANALYZE, BUFFERS)`, the N+1 pattern and how an API's shape causes it, keyset versus offset scanning, connection-pool exhaustion and how it presents as latency, and query timeouts driven by the request context. *Interview focus: diagnosing a query that is slow only in production.*

---

# Phase 7 — Identity and a Grown-Up API

*Backend / API Engineering · Production / Engineering* · Lessons 43–45

### Lesson 43 — Authentication

Password hashing with a memory-hard function and why not SHA-family, the login flow, sessions versus tokens and what each costs, JWT structure, signing, validation pitfalls, and expiry/refresh, secure cookie attributes, and where OAuth2/OIDC fit. Delivered through the middleware and context mechanisms already built: identity into context, credentials never into logs.

### Lesson 44 — Authorization

Authentication answers who; authorization answers whether. Resource ownership checks, role and policy models, why authorization usually cannot live entirely in middleware, failing closed, and avoiding the object-level access flaws that dominate real API breaches. *Interview focus: where an authorization decision belongs, and how it is tested.*

### Lesson 45 — API evolution and usability

Resource-oriented design revisited now that the data is real — what is a resource, what is a sub-resource, what is an action that does not fit REST, and when a non-REST endpoint is the honest answer. Then, since collections now come from a database with indexes behind them: pagination (offset versus keyset, and consistency under concurrent writes), filtering, sorting, and search; PUT versus PATCH using the pointer-presence model from Lesson 30; idempotency keys stored transactionally so a retried write is safe; a stable error contract; versioning and backward compatibility introduced only when a contract genuinely must change; content negotiation where it earns its place rather than by reflex; and documenting the API with OpenAPI. *Interview focus: pagination consistency, and making a payment endpoint safe to retry.*

---

# Phase 8 — Go Depth Where the Service Demands It

*Go Language* · Lessons 46–49

The service is now large enough that its own code raises these questions. Each topic is applied to code already written rather than demonstrated in isolation.

### Lesson 46 — Go's type system in depth

Defined types versus aliases and what each one actually creates; underlying types and the conversion rules that follow from them; why a `UserID` defined as an `int` stops being interchangeable with every other `int` in the service, and what that buys at an API boundary. The numeric types in full: sized integers, signed versus unsigned, overflow behaviour, and why `float64` carries its width in its name — the question left open in Lesson 4. Then the case that matters for a backend: money must not be a `float64`, because binary floating point cannot represent `0.10`, and the alternatives are integer minor units or a decimal type mapped to PostgreSQL `numeric`. Zero values as a deliberate design tool — making the useful state the zero state so a struct is usable before anything is set. *Interview focus: when a defined type earns its keep, and what goes wrong when money is a float.*

### Lesson 47 — Interfaces in depth

A genuine deepening of Lessons 15 and 18, which established implicit satisfaction and using an interface as a dependency. What those lessons could not explain: an interface value is a (type, value) pair, which is why an interface holding a nil `*User` is itself non-nil and `err != nil` fires on an error that was never set — the typed-nil trap, and the pointer-versus-value method set rule that produces it. Method sets stated formally, and why `*responseWriter` rather than `responseWriter` was the implementing type back in Lesson 33. Type assertions and the comma-`ok` form, and type switches — now motivated by `any` from Lesson 31, `errors.As` from Lesson 40, and `Scan` from Lesson 39. Finally, interface placement: defining a small interface where it is consumed rather than beside its implementation, and why an interface with one implementation and no test seam is usually noise. *Interview focus: why a non-nil interface can hold a nil pointer, and which method set a value versus a pointer carries.*

### Lesson 48 — Generics where they earn their place

Type parameters, constraints, type inference, and generic helpers and containers. Generics vary code by type; interfaces vary it by behavior — recognizing which axis a problem sits on, which is a question Lesson 47 has just made precise. Why generics rarely belong in domain logic. *Interview focus: constraints, the absence of method-based specialization, and choosing an interface over a type parameter.*

### Lesson 49 — Packages, standard library, and maintainable boundaries

Package API design, ownership and dependency direction, `internal`, constructors only where invariants or dependencies demand them, and the standard-library packages a backend leans on (`io` and its interfaces, `time` and timezone handling, `net/url`, `os`, `path/filepath`, `slices`, `maps`, `strings.Builder`, and `unicode/utf8` — which finally settles the rune question raised in Lesson 21). Formatting, `go vet`, linting, documentation comments, and reading Go code the way a reviewer does. *Interview focus: package boundaries, and when *not* to introduce an interface.*

---

# Phase 9 — External Systems and Concurrency

*Backend / API Engineering · Production / Engineering* · Lessons 50–54

### Lesson 50 — HTTP clients and failure-aware integration

`http.Client` with explicit timeouts and a configured transport, `NewRequestWithContext`, response-body lifecycle and connection reuse, handling non-2xx responses, retries with exponential backoff and jitter, retryable versus non-retryable failures, client-side idempotency, circuit breakers, and bulkheads. External dependencies are treated as unreliable, and their failure behavior is made observable and testable with `httptest.Server`. The structural question gets the same anti-cargo-cult treatment as the repository did: whether an external call deserves its own client type and interface, or whether the handler calling it directly is the honest design at this size. *Interview focus: retry storms, and why a retry without idempotency is a bug.*

### Lesson 51 — Concurrency foundations

Goroutines and their lifecycle, channel ownership and direction, buffered versus unbuffered semantics, closing and ranging, `select`, timers, `sync.Mutex`/`RWMutex`/`WaitGroup`/`Once`, atomics, and the happens-before model at a practical level. Closes the loop on `ctx.Done()`, which Lesson 36 used without explaining the channel underneath. *Interview focus: when a mutex is simpler and safer than a channel.*

### Lesson 52 — Concurrency patterns and failure modes

Worker pools, bounded concurrency with semaphores, fan-out/fan-in, pipelines, aggregating results and errors with `errgroup`, propagating cancellation, and the failure modes that matter: data races, deadlocks, goroutine leaks, and work that outlives the request that started it. Driven by race-detector-backed tests. *Interview focus: bounding concurrency, and finding a leaking goroutine.*

### Lesson 53 — Background jobs and queues

Separating request-time work from durable asynchronous work; job payloads and schemas, at-least-once execution, visibility timeouts, idempotent handlers, poison messages, dead-letter queues, scheduling, and operational ownership. Deciding when a PostgreSQL-backed worker suffices and when a broker is justified.

### Lesson 54 — Kafka and event-driven workflows

Topics and partitions, ordering guarantees and their scope, consumer groups and rebalancing, offsets and commit strategies, delivery semantics and acknowledgement, retries and DLQs, deduplication, schema evolution, eventual consistency, why distributed transactions across a database and a broker are avoided rather than solved — two-phase commit, sagas, and compensation as the alternatives — and the transactional outbox — which is only implementable because Lesson 41 established transaction boundaries. One workflow is built end to end across API, PostgreSQL, worker, and consumer. *Interview focus: at-least-once delivery, duplicate processing, and consumer lag.*

---

# Phase 10 — Operate It in Production

*Production / Engineering* · Lessons 55–59

### Lesson 55 — Configuration, startup, and graceful shutdown

Typed configuration from environment and secrets, validated at startup so a misconfigured service fails immediately rather than on its first request. Signal handling, `http.Server` timeouts, `Shutdown` and in-flight request draining, and health, readiness, and liveness endpoints that mean genuinely different things.

The other half of the lesson is **lifecycle ownership**, which the service can no longer avoid now that it holds a PostgreSQL connection pool, an `http.Client` and its transport, a Kafka consumer, a worker pool, a metrics exporter, and a tracer provider. For each one: who constructs it, who owns it for the life of the process, and who is responsible for closing or stopping it. The rule is that a dependency is created once in `main`, passed explicitly to whatever uses it, and released by the code that created it — not opened lazily somewhere in a handler, not stashed in a package-level variable, and not closed by a consumer that does not own it.

Ordering follows from the dependency graph. Startup builds from the bottom up — configuration, then the pool and clients, then the workers and consumers that use them, then the HTTP server last, because accepting traffic before its dependencies are ready is what a readiness probe exists to prevent. Shutdown reverses it: stop accepting new work, drain in-flight requests with `Shutdown`, cancel the context the workers and the Kafka consumer are selecting on so they finish their current message and commit their offsets, then close the pool and clients, and flush the metrics exporter and tracer provider last so the shutdown itself is observable. This is where `defer` from Lesson 19, context cancellation from Lesson 36, and the worker and consumer loops from Lessons 52–54 become one coordinated sequence — and where scattered ownership shows its cost, because a dependency closed in the wrong order produces errors on connections that are still in use. *Interview focus: what your service does between receiving SIGTERM and exiting.*

### Lesson 56 — Observability and production debugging

Completes the middleware track: instrumentation middleware emitting request metrics and opening a trace span per request, built on the status wrapper from Lesson 33 and the correlation ID from Lesson 36. Then the system around it — correlation IDs carried across handler, database, client, and worker; metrics and the RED/USE views; SLOs, dashboards, and alerts worth waking up for; distributed tracing with OpenTelemetry and context propagation; log levels, sampling, and cost; and the concepts behind tools such as Datadog. One request is followed across every component. *Interview focus: what you actually look at when latency rises at 3am.*

### Lesson 57 — Reliability and security engineering

Timeout budgets across a call chain, backpressure, overload protection and load shedding, caching and invalidation trade-offs, graceful degradation, TLS, secret handling and rotation, input and output safety, the common API vulnerability classes, dependency and supply-chain risk, and incident-oriented runbooks. *Interview focus: containing a failing dependency instead of amplifying it.*

### Lesson 58 — Performance and Go runtime diagnostics

Benchmarking a representative workload, `pprof` for CPU, heap, and blocking profiles, the execution tracer, escape analysis and allocation reduction, GC behavior and `GOGC`/memory limits, HTTP-level costs such as JSON encoding, body copying, and connection handling, and knowing when a database or network bound makes Go-level optimization pointless. Measure first, change second. *Interview focus: finding the bottleneck before changing code.*

### Lesson 59 — Delivery, deployment, and architecture review

Go tooling in CI, the unit/integration/contract test split, linting and `govulncheck`, multi-stage Docker builds and image hardening, CI/CD pipelines, deployment strategies and rollback, and the Kubernetes concepts a service owner needs: probes, resource requests and limits, configuration and secrets, autoscaling. The service's architecture is then reviewed by dependency direction, operational ownership, and failure behavior — not by folder fashion.

---

# Phase 11 — Professional Go Backend Proficiency

*All three areas* · Lessons 60–61

### Lesson 60 — Production service and design defense

Not a new project: the service that has been growing since Lesson 25, finished and operated. It carries a documented HTTP API, PostgreSQL persistence with migrations, authentication and authorization, meaningful tests, context-aware external integration, bounded background work, an event-driven workflow, structured logs, metrics and traces, health endpoints, a container image, and CI/CD. Its design document must state the API contract, data model, transaction boundaries, concurrency limits, retry and idempotency strategy, observability plan, security assumptions, and failure modes.

### Lesson 61 — Demonstrate professional proficiency

Defend a system design involving an API, a database, a justified cache, external dependencies, workers and queues, authentication, observability, failure handling, scaling, and deployment — and answer the interview material accumulated across the course cold.

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
