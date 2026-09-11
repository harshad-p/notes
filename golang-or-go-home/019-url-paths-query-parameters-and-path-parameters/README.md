# Lesson 27 — URL Paths, Query Parameters & Path Parameters

This is the next piece of HTTP fundamentals we need before building a useful REST API.

In the previous lessons, we had routes such as:

```text
GET /users
POST /users
```

But real APIs need to identify **which resource** we're talking about and often need additional information from the client.

For example:

```text
GET /users/42
```

might mean:

> Give me the user whose ID is 42.

Or:

```text
GET /users?role=admin
```

might mean:

> Give me users filtered by the role `admin`.

These are two different mechanisms:

- **Path parameters** identify something in the URL path.
- **Query parameters** provide additional options or filters.

We'll learn both carefully.

---

## 1. The URL has structure

Consider:

```text
https://example.com/users/42?active=true&sort=name
```

We can conceptually break this into:

```text
https://example.com
        │
        └── host

/users/42
    │
    └── path

?active=true&sort=name
    │
    └── query string
```

For our Go API, we're primarily interested in:

```text
/users/42?active=true
```

The **path** is:

```text
/users/42
```

The **query string** is:

```text
active=true
```

---

# 2. Reading the path from `http.Request`

We already learned that:

```go
r *http.Request
```

represents the incoming HTTP request.

It has a `URL` field.

```go
r.URL
```

represents the parsed URL associated with the request.

One of its fields is:

```go
r.URL.Path
```

which gives us the path portion.

For example, if the client requests:

```text
GET /users/42
```

then:

```go
fmt.Println(r.URL.Path)
```

prints:

```text
/users/42
```

So we can write:

```go
func userHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, r.URL.Path)
}
```

A request to:

```text
http://localhost:8080/users/42
```

would result in:

```text
/users/42
```

being written into the response.

---

# 3. What is a path parameter?

Suppose we want an endpoint for an individual user.

We could have:

```text
GET /users/42
```

Here, `42` is the user's ID.

We call it a **path parameter** because the value is contained within the URL path.

Conceptually, we define a route like:

```text
/users/{id}
```

where `{id}` represents a variable value.

So:

```text
/users/42
/users/17
/users/123
```

all represent the same general endpoint pattern:

```text
/users/{id}
```

but with different IDs.

---

# 4. Important: our current router doesn't automatically extract `{id}`

This is an important distinction.

The simple routing mechanism we've been using:

```go
http.HandleFunc("/users", userHandler)
```

doesn't give us a convenient `{id}` variable like some modern web frameworks do.

We're currently using Go's standard `net/http` routing capabilities, so we need to understand what the request actually contains before introducing more sophisticated routing.

For now, we can inspect the path ourselves.

For:

```text
/users/42
```

we know:

```go
r.URL.Path
```

contains:

```text
/users/42
```

We could then extract `42`.

But before doing that, we need to learn how to manipulate strings safely.

---

# 5. Splitting a path

Go's `strings` package provides:

```go
strings.Split()
```

We already encountered the `strings` package in our earlier strings lesson, but we haven't used `Split`, so let's introduce it properly.

`strings.Split` takes:

```go
strings.Split(value, separator)
```

and divides a string wherever the separator occurs.

For example:

```go
parts := strings.Split("users/42", "/")
```

produces a slice of strings conceptually equivalent to:

```text
["users", "42"]
```

With:

```text
/users/42
```

there is a leading `/`, so:

```go
parts := strings.Split("/users/42", "/")
```

produces:

```text
["", "users", "42"]
```

The first element is empty because the string begins with `/`.

Therefore:

```go
parts[0] // ""
parts[1] // "users"
parts[2] // "42"
```

This is useful for understanding how a path is structured, although manually parsing URLs this way is not how we'd ideally build a production router.

---

# 6. Why manually splitting paths isn't ideal

You might now think:

> Why not just split every URL and extract whatever we need?

You could, but things become complicated surprisingly quickly.

For example:

```text
/users/42
/users/42/orders
/products/15
/products/15/reviews
```

You would have to write increasingly complicated logic to determine which path you're dealing with.

You'd also need to consider:

- missing segments
- unexpected segments
- trailing slashes
- URL encoding
- HTTP methods
- conflicting routes

That's why routing libraries and newer versions of Go's standard `net/http` routing support provide better mechanisms.

We'll get there.

For now, manually inspecting the path is useful because it teaches you what's actually happening underneath the router.

---

# 7. Query parameters

Now let's look at the other mechanism.

Consider:

```text
GET /users?role=admin
```

The path is:

```text
/users
```

and the query string is:

```text
role=admin
```

Unlike a path parameter, `role=admin` isn't identifying a particular resource.

It's providing an **option/filter** for the request.

Another example:

```text
GET /users?role=admin&active=true
```

means the client is supplying two query parameters:

```text
role=admin
active=true
```

---

# 8. Reading query parameters in Go

Go gives us a convenient way to access them through the request URL.

We can use:

```go
r.URL.Query()
```

This returns the query parameters in a type called:

```go
url.Values
```

`url.Values` is essentially a mapping between parameter names and their values.

For example:

```text
/users?role=admin&active=true
```

can be represented conceptually as:

```text
role   → admin
active → true
```

We can retrieve one parameter with:

```go
role := r.URL.Query().Get("role")
```

Then:

```go
fmt.Fprintln(w, role)
```

would output:

```text
admin
```

---

# 9. What does `.Get()` do here?

This is worth being precise about.

```go
r.URL.Query()
```

returns the collection of query parameters.

Then:

```go
.Get("role")
```

asks that collection:

> Give me the value associated with the parameter named `role`.

So:

```go
role := r.URL.Query().Get("role")
```

is essentially:

```text
Request URL
     ↓
extract query parameters
     ↓
find "role"
     ↓
return its value
```

If the request is:

```text
/users?role=admin
```

then:

```go
role == "admin"
```

If the request doesn't contain `role`:

```text
/users
```

then `.Get("role")` returns the empty string.

---

# 10. Multiple query parameters

Suppose we receive:

```text
GET /users?role=admin&active=true
```

We can do:

```go
query := r.URL.Query()

role := query.Get("role")
active := query.Get("active")
```

Now:

```go
role   // "admin"
active // "true"
```

Notice that `active` is still a **string**.

The URL contains text:

```text
active=true
```

Go doesn't automatically know that `"true"` should become the Go boolean:

```go
true
```

If we need a boolean, we'd have to parse it.

We'll cover conversion/parsing properly when we need it rather than introducing another API here.

---

# 11. Query parameters vs path parameters

This distinction is extremely important for API design.

### Path parameter

```text
GET /users/42
```

`42` identifies the specific resource.

You can think:

> **Which user?** → 42

### Query parameter

```text
GET /users?role=admin
```

`role=admin` modifies the request.

You can think:

> **Which users / how should I retrieve them?** → users with role admin

Common examples:

```text
GET /products/123
```

Specific product.

```text
GET /products?category=laptops
```

Filter products.

```text
GET /products?page=2&limit=20
```

Pagination options.

```text
GET /products?sort=price
```

Sorting option.

---

# 12. They can be used together

There's nothing preventing both mechanisms from appearing in one URL.

For example:

```text
GET /users/42/orders?status=paid
```

Here:

```text
/users/42/orders
```

is the path.

The `42` can identify the user.

And:

```text
status=paid
```

is a query parameter.

Conceptually:

```text
/users/{userId}/orders?status={status}
```

means:

> Get orders belonging to user 42, filtered to paid orders.

This pattern is extremely common in REST APIs.

---

# 13. A practical exercise

Let's build an endpoint that demonstrates both concepts.

Create a route:

```text
/users
```

and have its handler inspect the request.

### Request 1

```text
GET /users?role=admin
```

Your handler should return:

```text
Role: admin
```

### Request 2

```text
GET /users?role=customer
```

It should return:

```text
Role: customer
```

### Request 3

```text
GET /users
```

It should detect that no `role` was supplied and return:

```text
Role not specified
```

---

## Then build a second endpoint

Create:

```text
/users/42
```

For now, don't worry about a sophisticated router.

Inside the handler, use:

```go
r.URL.Path
```

and demonstrate that you can obtain:

```text
/users/42
```

Then use the `strings` package to split the path and extract:

```text
42
```

Your output should be:

```text
User ID: 42
```

Also test:

```text
/users/17
/users/100
```

and verify that the extracted value changes.

---

# 14. One thing we're deliberately NOT doing yet

You may notice that manually extracting:

```text
42
```

from:

```text
/users/42
```

is cumbersome.

That's intentional.

Before we introduce a more capable router, I want you to understand the underlying HTTP URL structure first.

In the next part, we'll look at **proper route matching and path-variable extraction**, rather than continuing to build increasingly ugly `strings.Split` logic.

**Next lesson → Lesson 28: Go's HTTP Router — Route Patterns & Path Variables**