# Lesson 28 — Go's HTTP Router: Route Patterns & Path Variables

In Lesson 27, we deliberately extracted `/users/42` ourselves. That was useful for understanding the URL, but **we don't want application code manually splitting URL strings every time we need an ID**.

Now we'll learn how Go's HTTP router can do this for us.

One important correction to our previous lesson: **modern Go's standard `net/http` router can match method-specific patterns and named path variables directly.** Since you're using Go 1.26, we'll use that rather than introducing a third-party router prematurely.

---

## 1. What is a router?

A **router** is responsible for deciding which handler should receive an incoming HTTP request.

Suppose our API has:

```text
GET  /users
POST /users
GET  /users/42
GET  /products
```

When a request arrives, something needs to determine:

```text
What path was requested?
What HTTP method was used?
Which handler should process it?
```

That's the router's job.

Conceptually:

```text
HTTP request
     │
     ▼
   Router
     │
     ├── GET /users       → listUsers
     ├── POST /users      → createUser
     ├── GET /users/{id}  → getUser
     └── ...
```

This means your handlers can concentrate on **handling** requests rather than figuring out which route they belong to.

---

# 2. `http.ServeMux`

Go's standard library provides a router called `ServeMux`.

The name comes from **multiplexer**.

A multiplexer is something that receives incoming requests and directs them to one of several possible destinations.

We can create one explicitly:

```go
mux := http.NewServeMux()
```

`http.NewServeMux()` creates a new `ServeMux`.

Think of `mux` as:

> Our application's collection of HTTP routes.

We can then register routes on it.

---

# 3. Registering a route

We previously used:

```go
http.HandleFunc("/users", userHandler)
```

That uses a default/global router.

Now we'll create our own:

```go
mux := http.NewServeMux()

mux.HandleFunc("/users", userHandler)
```

The difference is that we're explicitly creating the router instead of relying on the package-level default.

This becomes useful as an application grows because we can construct and configure the router ourselves and then give it to the HTTP server.

---

# 4. Giving the router to the server

Earlier we had:

```go
http.ListenAndServe(":8080", nil)
```

The second argument is the handler that receives incoming HTTP requests.

A `ServeMux` implements the required handler behavior, so we can give it to the server:

```go
mux := http.NewServeMux()

mux.HandleFunc("/users", userHandler)

http.ListenAndServe(":8080", mux)
```

The flow is now:

```text
HTTP request
     │
     ▼
HTTP server
     │
     ▼
   mux
     │
     ▼
userHandler
```

---

# 5. Method-specific routes

Here's where modern Go becomes much more convenient.

We can register:

```go
mux.HandleFunc("GET /users", getUsersHandler)
mux.HandleFunc("POST /users", createUserHandler)
```

Notice that the method is part of the pattern:

```text
GET /users
POST /users
```

Now the router itself distinguishes the methods.

We no longer need:

```go
if r.Method == http.MethodGet {
    ...
}

if r.Method == http.MethodPost {
    ...
}
```

inside the handler.

Instead:

```go
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Getting users")
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Creating user")
}
```

The router has already decided which handler should run.

That's a cleaner separation of responsibilities:

```text
Router
    ↓
Which handler?

Handler
    ↓
What should happen?
```

---

# 6. Path variables

Now we get to the interesting part.

We want:

```text
GET /users/42
```

to call a handler that can access:

```text
42
```

as the user ID.

We can define the route:

```go
mux.HandleFunc("GET /users/{id}", getUserHandler)
```

The `{id}` is a **path wildcard**.

It means:

> This part of the path can contain a variable value, and that value should be associated with the name `id`.

So all of these can match:

```text
/users/1
/users/42
/users/123
```

The fixed portion is:

```text
/users/
```

and the variable portion is:

```text
{id}
```

---

# 7. Getting the path variable

Once the router has matched:

```text
GET /users/{id}
```

the handler can retrieve the value using:

```go
id := r.PathValue("id")
```

Let's break this down.

`r` is our `*http.Request`.

`PathValue` is a method on `http.Request`.

```go
r.PathValue("id")
```

means:

> Give me the value that was captured for the path variable named `"id"`.

For:

```text
GET /users/42
```

this gives:

```text
"42"
```

For:

```text
GET /users/123
```

it gives:

```text
"123"
```

Notice that the result is a **string**.

We'll deal with converting `"42"` into the integer `42` when we cover parsing/conversion in more detail.

---

# 8. Complete example

Here's a small application using everything we've introduced in this lesson:

```go
package main

import (
    "fmt"
    "net/http"
)

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Getting users")
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Creating user")
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")

    fmt.Fprintln(w, "User ID:", id)
}

func main() {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /users", getUsersHandler)
    mux.HandleFunc("POST /users", createUserHandler)
    mux.HandleFunc("GET /users/{id}", getUserHandler)

    http.ListenAndServe(":8080", mux)
}
```

Now:

```text
GET /users
```

goes to:

```go
getUsersHandler
```

while:

```text
POST /users
```

goes to:

```go
createUserHandler
```

and:

```text
GET /users/42
```

goes to:

```go
getUserHandler
```

where:

```go
r.PathValue("id")
```

is:

```text
"42"
```

---

# 9. Why this is better than `strings.Split`

Previously we had to do something like:

```go
parts := strings.Split(r.URL.Path, "/")
id := parts[2]
```

Now we simply say:

```go
id := r.PathValue("id")
```

The router understands that:

```text
/users/{id}
```

contains a variable and performs the matching/extraction for us.

Your application code no longer needs to know that the ID happened to be the third element after splitting on `/`.

This is a major improvement because your code expresses the **meaning** rather than the mechanics.

---

# 10. Path variables can appear in the middle

They aren't restricted to the end.

For example:

```go
mux.HandleFunc("GET /users/{userID}/orders/{orderID}", handler)
```

matches:

```text
/users/42/orders/817
```

Then:

```go
userID := r.PathValue("userID")
orderID := r.PathValue("orderID")
```

gives:

```text
userID → "42"
orderID → "817"
```

This lets the URL express relationships between resources.

For example:

```text
/users/42/orders/817
```

can naturally mean:

> Order 817 belonging to user 42.

---

# 11. Query parameters still work

Path variables and query parameters are independent.

For example:

```text
GET /users/42?includeOrders=true
```

contains:

**Path variable:**

```text
id = "42"
```

retrieved with:

```go
r.PathValue("id")
```

**Query parameter:**

```text
includeOrders = "true"
```

retrieved with:

```go
r.URL.Query().Get("includeOrders")
```

So you might have:

```go
func getUserHandler(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    includeOrders := r.URL.Query().Get("includeOrders")

    fmt.Fprintln(w, "User ID:", id)
    fmt.Fprintln(w, "Include orders:", includeOrders)
}
```

This gives you a useful mental model:

```text
/users/42?includeOrders=true
     │             │
     │             └── Query parameter
     │
     └── Path variable
```

---

# 12. What happens when a route doesn't match?

Suppose we only register:

```go
mux.HandleFunc("GET /users/{id}", getUserHandler)
```

and somebody requests:

```text
GET /products/10
```

There is no matching route.

The router therefore doesn't call `getUserHandler`.

The HTTP server returns a **404 Not Found** response.

This is another advantage of having routing separated from your application logic: the router handles the question of whether a request corresponds to a known endpoint.

---

# 13. What happens when the path matches but the method doesn't?

Suppose we have:

```go
mux.HandleFunc("GET /users/{id}", getUserHandler)
```

and the client sends:

```text
POST /users/42
```

The path pattern exists, but the registered method is `GET`, not `POST`.

Modern `ServeMux` can distinguish this situation from a completely unknown path and respond with **405 Method Not Allowed** when an appropriate path exists for another method.

This means our router can now handle much of what we manually implemented in Lesson 26.

---

# Exercise

Build this API:

```text
GET    /products
POST   /products
GET    /products/{id}
DELETE /products/{id}
```

Create a separate handler for each route.

### `GET /products`

Return:

```text
Getting products
```

### `POST /products`

Return:

```text
Creating product
```

### `GET /products/{id}`

Extract the ID using:

```go
r.PathValue("id")
```

and return:

```text
Getting product: 42
```

for:

```text
GET /products/42
```

### `DELETE /products/{id}`

Again extract the ID and return:

```text
Deleting product: 42
```

---

### Then test these requests

```text
GET    /products
POST   /products
GET    /products/42
GET    /products/99
DELETE /products/42
PUT    /products/42
GET    /something-else
```

Pay attention to the difference between:

```text
405 Method Not Allowed
```

and:

```text
404 Not Found
```

You should now be able to explain **why** the router produces each one.

---

### One important thing we're postponing

We're currently getting:

```go
id := r.PathValue("id")
```

as a string.

A real API will usually need to turn that into an integer:

```text
"42"
  ↓
42
```

We haven't properly covered string-to-number conversion yet, so **we won't sneak it into this lesson**.

That's next.

**Next lesson → Lesson 29: Type Conversion & Parsing — Turning URL Strings into Go Types**