# Lesson 26 — HTTP Methods, Routing & Status Codes

In Lesson 25, we learned how an HTTP request reaches a handler and how JSON can travel between the request/response bodies and Go structs.

Now we're going to answer a different question:

> **How does our server decide which code should handle a particular HTTP request, and how does it tell the client what happened?**

We'll build this from the ground up.

---

## 1. An HTTP request has a method

Consider these requests:

```text
GET /product
```

and:

```text
POST /product
```

They have the **same path**:

```text
/product
```

but different **methods**.

The method describes what kind of operation the client is requesting.

Some common HTTP methods are:

| Method | Typical purpose |
|---|---|
| GET | Retrieve data |
| POST | Create/submit data |
| PUT | Replace an existing resource |
| PATCH | Partially modify a resource |
| DELETE | Delete a resource |

For example, an API might use:

```text
GET    /users       → retrieve users
POST   /users       → create a user
GET    /users/42    → retrieve user 42
DELETE /users/42    → delete user 42
```

The method and path together help describe what operation is being requested.

---

# 2. Go gives us the method through `http.Request`

We already learned that:

```go
r *http.Request
```

represents the incoming request.

One of its fields is:

```go
r.Method
```

This is a `string`.

So if the client sends:

```text
GET /user
```

then:

```go
fmt.Println(r.Method)
```

prints:

```text
GET
```

If the client sends:

```text
POST /user
```

it prints:

```text
POST
```

We can therefore inspect it ourselves:

```go
func userHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println(r.Method)
}
```

---

# 3. Comparing the method

Because `r.Method` is a string, we could technically write:

```go
if r.Method == "GET" {
    // ...
}
```

That works.

However, Go's `net/http` package provides constants for the standard methods:

```go
http.MethodGet
http.MethodPost
http.MethodPut
http.MethodPatch
http.MethodDelete
```

For example:

```go
if r.Method == http.MethodGet {
    // ...
}
```

`http.MethodGet` is simply a predefined constant whose value is `"GET"`.

So this:

```go
r.Method == http.MethodGet
```

is essentially comparing:

```text
the request's method == "GET"
```

Using the provided constants avoids repeatedly writing method strings yourself and makes the code clearer.

---

# 4. Handling different methods

We can use `if`:

```go
func userHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        fmt.Fprintln(w, "Getting users")
        return
    }

    if r.Method == http.MethodPost {
        fmt.Fprintln(w, "Creating user")
        return
    }
}
```

A request to:

```text
GET /user
```

takes the first branch.

A request to:

```text
POST /user
```

takes the second branch.

The `return` stops the handler at that point.

---

# 5. `switch` is often cleaner

When we're comparing one value against several possibilities, Go's `switch` is convenient.

```go
func userHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        fmt.Fprintln(w, "Getting users")

    case http.MethodPost:
        fmt.Fprintln(w, "Creating user")

    default:
        fmt.Fprintln(w, "Something else")
    }
}
```

Here:

```go
switch r.Method
```

means:

> Examine the value of `r.Method`.

Then each `case` specifies a possible value.

Unlike C#'s traditional `switch`, you don't need `break` after each case. Go automatically stops after the matching case unless you explicitly use special behavior such as `fallthrough`.

We'll use `switch` frequently when handling HTTP methods.

---

# 6. What is routing?

So far, we've used:

```go
http.HandleFunc("/user", userHandler)
```

We need to understand what this actually does.

A **route** associates an incoming request path with a handler.

For example:

```go
http.HandleFunc("/user", userHandler)
```

essentially tells Go:

> When a request comes to the `/user` path, use `userHandler` to handle it.

So:

```text
GET /user
```

can reach:

```go
userHandler(...)
```

and:

```text
POST /user
```

can also reach:

```go
userHandler(...)
```

The route determines **where the request goes**.

The handler can then inspect the method to determine **what operation to perform**.

That's an important distinction:

```text
Route
  ↓
Which handler?

Method
  ↓
Which operation inside the handler?
```

Later we'll learn more sophisticated routing where the router itself can distinguish methods.

For now, we're deliberately doing it manually so you understand the underlying mechanism.

---

# 7. What is an HTTP response status?

So far our handlers have returned text or JSON.

But an HTTP response contains more than a body.

Conceptually:

```text
HTTP response
├── Status
├── Headers
└── Body
```

For example:

```text
HTTP/1.1 200 OK
Content-Type: application/json

{"name":"Harshad","age":36}
```

The status tells the client what happened.

---

# 8. `200 OK`

The most common success status is:

```text
200 OK
```

It means, broadly:

> The request was successfully processed.

Go provides:

```go
http.StatusOK
```

which represents the numeric status code:

```text
200
```

You can explicitly send it with:

```go
w.WriteHeader(http.StatusOK)
```

---

# 9. What does `WriteHeader` actually do?

This is important.

`w` is the `http.ResponseWriter` from Lesson 25.

It is how the handler constructs the response.

`WriteHeader` tells Go:

> Set the HTTP response's status code to this value.

For example:

```go
w.WriteHeader(http.StatusCreated)
```

means:

```text
HTTP status = 201
```

It does **not** write a body.

So:

```go
w.WriteHeader(http.StatusCreated)
fmt.Fprintln(w, "User created")
```

produces a response conceptually like:

```text
HTTP/1.1 201 Created

User created
```

The status and body are separate parts of the response.

---

# 10. `201 Created`

When a `POST` request successfully creates a new resource, `201 Created` is commonly appropriate.

Go provides:

```go
http.StatusCreated
```

which represents:

```text
201
```

Example:

```go
func createUserHandler(w http.ResponseWriter, r *http.Request) {
    // create user...

    w.WriteHeader(http.StatusCreated)
    fmt.Fprintln(w, "User created")
}
```

The client receives a `201` rather than the default success response.

---

# 11. `405 Method Not Allowed`

Suppose our endpoint supports:

```text
GET /user
POST /user
```

but someone sends:

```text
DELETE /user
```

The request reached the correct path, but our API doesn't support that method for this endpoint.

The appropriate status is commonly:

```text
405 Method Not Allowed
```

Go provides:

```go
http.StatusMethodNotAllowed
```

which represents:

```text
405
```

We can use:

```go
http.Error(
    w,
    "Method not allowed",
    http.StatusMethodNotAllowed,
)
```

We've already learned `http.Error` in Lesson 25: it sends an error response with the supplied message and status code.

---

# 12. Putting methods and status codes together

Now we can build a more meaningful handler:

```go
func userHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        fmt.Fprintln(w, "Getting users")

    case http.MethodPost:
        w.WriteHeader(http.StatusCreated)
        fmt.Fprintln(w, "Creating user")

    default:
        http.Error(
            w,
            "Method not allowed",
            http.StatusMethodNotAllowed,
        )
    }
}
```

The behavior is now:

```text
GET /user
    ↓
200 OK
Getting users
```

and:

```text
POST /user
    ↓
201 Created
Creating user
```

while:

```text
DELETE /user
    ↓
405 Method Not Allowed
Method not allowed
```

---

# 13. A subtle but important detail about `WriteHeader`

You might notice that we didn't explicitly write:

```go
w.WriteHeader(http.StatusOK)
```

for the GET case.

That's because Go's HTTP server will generally send a successful status automatically if you write a response without explicitly selecting another status.

So:

```go
fmt.Fprintln(w, "Getting users")
```

will normally result in:

```text
200 OK
```

You **can** explicitly write:

```go
w.WriteHeader(http.StatusOK)
```

but there's usually no need when `200` is the desired result.

This is different from `201`, `400`, `404`, etc., where you need to deliberately choose the status.

---

# 14. Why status codes matter

Imagine your API always returned:

```text
200 OK
```

even when something went wrong.

A client would have to inspect the response body to figure out whether the operation succeeded.

Instead, HTTP already gives us a standardized way of communicating the result:

```text
200 → successful retrieval/operation
201 → successfully created
400 → client sent an invalid request
401 → authentication required/failed
403 → authenticated but not allowed
404 → resource not found
405 → method not supported
500 → server encountered an error
```

This becomes particularly important when we build REST APIs.

For example:

```text
GET /users/123
```

If user `123` exists:

```text
200 OK
```

If user `123` doesn't exist:

```text
404 Not Found
```

The body can provide additional information, but the status code communicates the broad result consistently.

---

# 15. Your exercise

Let's make this a proper exercise rather than a tiny syntax test.

Create a new Go program containing a `/users` endpoint.

### Requirements

Register:

```text
/users
```

with a handler.

The handler must support exactly these methods:

### GET

Request:

```text
GET /users
```

Response:

```text
200 OK
```

Body:

```text
Getting users
```

### POST

Request:

```text
POST /users
```

Response:

```text
201 Created
```

Body:

```text
Creating user
```

### DELETE

Request:

```text
DELETE /users
```

Response:

```text
405 Method Not Allowed
```

Body should indicate that the method isn't supported.

### Anything else

For example:

```text
PUT /users
```

should also receive:

```text
405 Method Not Allowed
```

---

## Then test it

Use Postman, `curl`, or another HTTP client and make all four requests:

```text
GET
POST
DELETE
PUT
```

Don't just look at the response body.

**Look specifically at the HTTP status code.**

You should observe:

```text
GET     → 200
POST    → 201
DELETE  → 405
PUT     → 405
```

---

## One final distinction to retain

At this point, keep these three concepts separate:

```text
URL path
    ↓
Routing decides which handler receives the request

HTTP method
    ↓
Handler decides what operation is being requested

HTTP status
    ↓
Server tells the client what happened
```

We'll build on all three in the next lesson.

**Next lesson → Lesson 27: URL Paths, Query Parameters & Path Parameters**