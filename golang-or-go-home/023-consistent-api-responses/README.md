# Lesson 31 — HTTP Responses in JSON: Status Codes, Headers & Consistent API Responses

We've already used JSON encoding and HTTP status codes separately. Now we'll put them together into a more realistic API response design.

A production API shouldn't just return:

```text
Product created
```

for one endpoint and raw JSON for another. It should have predictable response behavior.

We'll focus on three things:

1. Returning JSON consistently
2. Choosing the correct HTTP status
3. Understanding the interaction between headers, status, and body

---

## 1. An HTTP response has three important parts

An HTTP response consists conceptually of:

```text
Status
Headers
Body
```

For example:

```text
HTTP/1.1 201 Created
Content-Type: application/json

{
    "id": 42,
    "name": "Keyboard",
    "price": 99.99
}
```

Here:

- `201 Created` is the **status**
- `Content-Type: application/json` is a **header**
- the JSON object is the **body**

We've already worked with each of these individually. The important thing now is understanding how they work together.

---

# 2. `Content-Type` tells the client what the body contains

When we return JSON, we should tell the client that the response body contains JSON:

```go
w.Header().Set("Content-Type", "application/json")
```

`w.Header()` gives us the response's HTTP headers.

`Set` assigns a value to a header name.

So:

```go
w.Header().Set("Content-Type", "application/json")
```

means:

> Set the response's `Content-Type` header to `application/json`.

This is important because the body itself is just bytes. The HTTP protocol needs metadata telling the client how those bytes should be interpreted.

For JSON APIs, this is normally:

```text
Content-Type: application/json
```

---

# 3. Set headers before writing the response

There is an important ordering rule.

HTTP headers are sent before the response body.

Therefore, this:

```go
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(product)
```

is correct.

But if you have already started writing the response, changing headers may be too late.

This is the same reason we discussed in Lesson 26 that status handling has an ordering aspect.

A useful mental model is:

```text
Set headers
    ↓
Set status (if needed)
    ↓
Write body
```

Once the response has started going to the client, you can't freely modify what was already sent.

---

# 4. Returning a JSON response

Suppose we have:

```go
type ProductResponse struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}
```

We can return it as JSON:

```go
func getProductHandler(w http.ResponseWriter, r *http.Request) {
    product := ProductResponse{
        ID:    42,
        Name:  "Keyboard",
        Price: 99.99,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(product)
}
```

The client receives approximately:

```json
{
    "id": 42,
    "name": "Keyboard",
    "price": 99.99
}
```

The encoder writes the JSON directly to the `ResponseWriter`.

Because we haven't explicitly set a status, writing the body causes the response to use the default successful status:

```text
200 OK
```

So this is effectively:

```text
200 OK
Content-Type: application/json

{
    "id": 42,
    ...
}
```

---

# 5. When should we explicitly set the status?

If the status is `200 OK`, you generally don't need:

```go
w.WriteHeader(http.StatusOK)
```

because `200` is the default when the handler successfully writes a response without selecting another status.

For a different status, however, we need to explicitly set it.

For example, when creating a resource:

```go
w.WriteHeader(http.StatusCreated)
```

This produces:

```text
201 Created
```

So a create endpoint could do:

```go
func createProductHandler(w http.ResponseWriter, r *http.Request) {
    product := ProductResponse{
        ID:    42,
        Name:  "Keyboard",
        Price: 99.99,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(product)
}
```

The response is:

```text
201 Created
Content-Type: application/json

{
    "id": 42,
    "name": "Keyboard",
    "price": 99.99
}
```

Notice the order:

```go
w.Header().Set(...)
w.WriteHeader(...)
json.NewEncoder(w).Encode(...)
```

That's intentional.

---

# 6. Status must be written before the body

Consider:

```go
json.NewEncoder(w).Encode(product)
w.WriteHeader(http.StatusCreated)
```

This is wrong.

The JSON encoding writes to `w`, which starts the HTTP response. By that point, Go has already committed the response status, normally as `200 OK`.

Trying to change it afterward is too late.

So:

```text
Headers
   ↓
Status
   ↓
Body
```

is the safe order.

---

# 7. Different operations normally use different statuses

For the API we've been building, a reasonable starting point is:

| Operation | Typical success status |
|---|---:|
| Get a resource | `200 OK` |
| Create a resource | `201 Created` |
| Update a resource | `200 OK` or `204 No Content` |
| Delete a resource | `204 No Content` |

We've already encountered `204 No Content`, but let's clarify its purpose.

A `204` response means:

> The request succeeded, but there is no response body.

For example, after successfully deleting a product:

```go
w.WriteHeader(http.StatusNoContent)
```

The client gets:

```text
204 No Content
```

and no JSON body.

This can be preferable to returning:

```json
{
    "message": "Product deleted"
}
```

when the client doesn't actually need any data.

Again, this isn't an absolute rule. API design is a contract between the client and server.

---

# 8. JSON error responses

Our current errors have been things like:

```go
http.Error(w, "Invalid product ID", http.StatusBadRequest)
```

This produces a plain-text error response.

That's fine for a tiny example, but if our API is otherwise JSON-based, we might want errors to also be JSON.

For example:

```json
{
    "error": "Invalid product ID"
}
```

We can define a response type:

```go
type ErrorResponse struct {
    Error string `json:"error"`
}
```

Then:

```go
func writeError(w http.ResponseWriter, status int, message string) {
    response := ErrorResponse{
        Error: message,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(response)
}
```

Now a handler can do:

```go
writeError(w, http.StatusBadRequest, "Invalid product ID")
```

and the client receives:

```text
400 Bad Request
Content-Type: application/json

{
    "error": "Invalid product ID"
}
```

---

# 9. Why create a helper?

`writeError` is a small example of a helper function that centralizes repeated HTTP response behavior.

Without it, every handler might contain:

```go
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusBadRequest)
json.NewEncoder(w).Encode(...)
```

repeated over and over.

The helper gives us one place to establish the behavior.

It also means that if we later decide our error format should be:

```json
{
    "error": {
        "message": "Invalid product ID"
    }
}
```

we have one implementation to change rather than dozens.

But don't take this as a rule that every repeated three-line sequence deserves a helper. The abstraction becomes worthwhile when it represents a meaningful, reusable operation—in this case, **writing an API error response**.

---

# 10. A JSON response helper

We can apply the same idea to successful responses:

```go
func writeJSON(w http.ResponseWriter, status int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}
```

There is one new thing here:

```go
any
```

`any` is an alias for Go's empty interface:

```go
interface{}
```

It means the parameter can contain a value of any type.

We haven't needed `any` much so far because we've generally known the concrete type we're working with. Here, however, the helper is deliberately generic: it should be able to encode a `ProductResponse`, `UserResponse`, `ErrorResponse`, or another response type.

We could therefore write:

```go
writeJSON(w, http.StatusOK, product)
```

or:

```go
writeJSON(w, http.StatusCreated, product)
```

The helper doesn't care what concrete value it's encoding. `json.NewEncoder` can serialize the supplied value.

---

# 11. One subtle problem with a generic JSON helper

Our helper currently always calls:

```go
w.WriteHeader(status)
```

That's perfectly reasonable, but remember that `200` is already the default.

So this:

```go
writeJSON(w, http.StatusOK, product)
```

explicitly writes `200`, while directly encoding the product would implicitly produce `200`.

There's nothing wrong with explicitly setting it in a helper. In fact, it can make the helper's behavior predictable:

> Whatever status the caller gives me, I will use it.

That's different from unnecessarily writing `WriteHeader(http.StatusOK)` throughout individual handlers.

---

# 12. What about `204 No Content`?

This is where the helper demonstrates why abstraction should reflect actual behavior.

We should **not** do:

```go
writeJSON(w, http.StatusNoContent, nil)
```

and assume that's equivalent to a proper `204`.

A `204 No Content` response should not contain a response body.

So a delete handler can simply do:

```go
w.WriteHeader(http.StatusNoContent)
```

No JSON encoding is necessary.

This is a good example of why we shouldn't create one giant helper that tries to handle every possible HTTP response.

Different HTTP semantics sometimes deserve different code.

---

# 13. A practical response design

For our API, we can now establish a simple convention:

### Successful JSON response

```text
Content-Type: application/json
```

with an appropriate `2xx` status.

### JSON error response

```json
{
    "error": "..."
}
```

with an appropriate `4xx` or `5xx` status.

### No-content response

Use:

```text
204 No Content
```

when there is deliberately no response body.

This gives clients predictable behavior.

---

# 14. Don't put HTTP semantics into your domain models

One architectural point worth making now.

Avoid doing something like:

```go
type ProductResponse struct {
    StatusCode int     `json:"statusCode"`
    ID         int     `json:"id"`
    Name       string  `json:"name"`
}
```

The HTTP status isn't part of the product.

It's part of the **HTTP response**.

Likewise, your domain model shouldn't need to know that the HTTP API decided to return `201`.

Keep the concerns separate:

```text
Product
    ↓
data


HTTP response
    ↓
status + headers + serialized data
```

The handler or HTTP layer combines these things.

---

# 15. What we've built so far

Our API layer now has a fairly clear responsibility:

```text
HTTP request
     |
     v
Extract path/query values
     |
     v
Parse values
     |
     v
Decode JSON
     |
     v
Validate request
     |
     v
Application logic
     |
     v
Create result
     |
     v
Choose HTTP status
     |
     v
Set response headers
     |
     v
Encode JSON
```

That's already enough structure to build a reasonably clean small HTTP API without introducing a framework.

---

# Exercise

Take the product API from the previous lessons and make its responses consistently JSON-based.

Implement:

```text
GET    /products
GET    /products/{id}
POST   /products
DELETE /products/{id}
```

Use sensible HTTP status codes for successful operations.

For errors, return JSON in the form:

```json
{
    "error": "..."
}
```

Create a small reusable helper for JSON responses and another for JSON errors **if you think the abstraction is justified**.

For `DELETE`, think carefully about whether a response body is actually necessary.

Also consider the order in which you set:

- headers
- status
- body

Don't add a database yet. Keep the data in memory.

### Next: Lesson 32 — Middleware: Intercepting and Wrapping HTTP Requests