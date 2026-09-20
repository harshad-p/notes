# Lesson 30 — JSON API Design: Request/Response Models & Validation

So far, we've learned how to receive JSON and turn it into a Go struct, and how to serialize a Go struct back to JSON.

The next step is an important API-design question:

> **Should the same struct represent the HTTP request, the HTTP response, and the database entity?**

Often, **no**.

This lesson is about designing the types at the HTTP boundary deliberately, and validating incoming data before it reaches application logic.

---

## 1. The problem with using one struct for everything

Suppose we have:

```go
type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}
```

It is tempting to use this everywhere:

```text
HTTP request
    ↓
User
    ↓
database
    ↓
User
    ↓
HTTP response
```

But this creates problems.

For example, a client creating a user should probably **not provide the database ID**:

```json
{
    "id": 42,
    "name": "Alice",
    "email": "alice@example.com"
}
```

The server/database should normally determine the ID.

Even more importantly, we should never accidentally return:

```json
{
    "id": 42,
    "name": "Alice",
    "email": "alice@example.com",
    "password": "secret"
}
```

The password is internal data and shouldn't be part of the API response.

This is where separate request and response models become useful.

---

# 2. Request models

A request model represents the data the **client is allowed or expected to send**.

For creating a user:

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

The JSON request might therefore be:

```json
{
    "name": "Alice",
    "email": "alice@example.com"
}
```

Notice that there is no `ID`.

That's intentional.

The API contract says:

> "To create a user, provide a name and email."

The request model expresses that contract directly.

---

# 3. Response models

A response model represents the data the **API chooses to expose**.

```go
type UserResponse struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

A response might be:

```json
{
    "id": 42,
    "name": "Alice",
    "email": "alice@example.com"
}
```

There is no password.

This gives us a deliberate boundary:

```text
Client
   |
   | CreateUserRequest
   v
Handler / Application
   |
   | internal model
   v
Database
   |
   | internal model
   v
UserResponse
   |
   v
Client
```

The exact number of types depends on the application. The important idea is that **API contracts don't have to be identical to internal data structures**.

---

# 4. Why this matters beyond passwords

Separating API models from internal models gives you control over your public contract.

Imagine your database model eventually becomes:

```go
type User struct {
    ID           int
    Name         string
    Email        string
    PasswordHash string
    CreatedAt    time.Time
    InternalFlag bool
}
```

You probably don't want all of those fields exposed through the API.

Your response can remain:

```go
type UserResponse struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

The database can evolve without necessarily changing your HTTP contract.

This is one reason DTOs—**Data Transfer Objects**—are common in backend applications.

In Go, though, don't create DTO types mechanically just because "architecture says DTOs." Create separate types when the boundaries actually need different representations.

---

# 5. JSON decoding does not validate your business rules

This is an important distinction.

Suppose we have:

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

and receive:

```json
{
    "name": "",
    "email": "not-an-email"
}
```

This is perfectly valid **JSON**.

So:

```go
json.NewDecoder(r.Body).Decode(&request)
```

can succeed.

But the data may still be invalid according to our application's rules.

There are two separate questions:

### Is the JSON syntactically valid?

For example:

```json
{
    "name": "Alice"
}
```

versus malformed JSON:

```json
{
    "name": "Alice"
```

That's the decoder's job.

### Is the data acceptable?

For example:

```text
name must not be empty
email must be valid
```

That's **validation**.

So:

```text
JSON decoding
     ↓
Can I understand this JSON?
     ↓
Validation
     ↓
Is the supplied data acceptable?
```

These are different responsibilities.

---

# 6. Basic validation in Go

Go doesn't have a built-in validation framework equivalent to something like ASP.NET's model-validation ecosystem.

For simple validation, ordinary Go code is often enough.

For example:

```go
if request.Name == "" {
    http.Error(w, "Name is required", http.StatusBadRequest)
    return
}
```

This isn't sophisticated, but that's precisely the point: for a small API, introducing a validation framework just to check one or two fields may add unnecessary complexity.

For more substantial validation rules, a validation library can make sense. We can cover that later when we actually have a reason to introduce one.

---

# 7. Validation belongs after decoding

A typical create-user handler might therefore look like:

```go
func createUserHandler(w http.ResponseWriter, r *http.Request) {
    var request CreateUserRequest

    err := json.NewDecoder(r.Body).Decode(&request)
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    if request.Name == "" {
        http.Error(w, "Name is required", http.StatusBadRequest)
        return
    }

    if request.Email == "" {
        http.Error(w, "Email is required", http.StatusBadRequest)
        return
    }

    // Application logic comes here.
}
```

The important architectural sequence is:

```text
HTTP request
     ↓
Decode JSON
     ↓
Validate input
     ↓
Application logic
     ↓
Persist / retrieve data
     ↓
Create response
```

The handler shouldn't pass obviously invalid data deeper into the application and hope something else catches it.

---

# 8. `400 Bad Request`

We've already encountered:

```go
http.StatusBadRequest
```

This represents HTTP status **400**.

It's appropriate when the client sends a request that the server cannot reasonably process because the request itself is invalid.

Examples include:

- malformed JSON
- invalid integer in a path parameter
- missing required field
- invalid field value
- invalid query parameter

For example:

```text
POST /users
```

with:

```json
{
    "name": ""
}
```

could result in:

```text
400 Bad Request
```

because the client supplied invalid input.

---

# 9. Validation is more than checking for empty strings

Real validation often involves relationships and constraints.

For example:

```text
age >= 18
quantity > 0
startDate < endDate
currency must be supported
email must have an acceptable format
```

Some validation is purely structural:

```text
name is required
quantity is an integer
```

Other validation involves domain rules:

```text
a booking cannot end before it starts
an order cannot contain zero items
a user cannot change their email to one already belonging to another account
```

That distinction becomes important later.

A useful way to think about it is:

### Input validation

> "Is this request well-formed and acceptable as input?"

### Business validation

> "Does this operation make sense according to the application's domain rules?"

Those don't necessarily belong in the same place.

For example, checking that `quantity > 0` can naturally happen close to the request boundary. Checking whether a particular product can actually be ordered may require database/application logic.

We don't need to solve that architecture yet; just recognize the distinction.

---

# 10. Request models and partial updates

This becomes particularly interesting with `PATCH`.

Suppose we have:

```go
type UpdateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

What does this mean?

If the client sends:

```json
{
    "name": "Bob"
}
```

does that mean:

> Change the name, leave email unchanged

or:

> Set name to Bob and email to its zero value?

With ordinary Go fields, you can't necessarily distinguish:

```text
field wasn't supplied
```

from:

```text
field was supplied with its zero value
```

This is one reason PATCH request models often use pointers:

```go
type UpdateUserRequest struct {
    Name  *string `json:"name"`
    Email *string `json:"email"`
}
```

Now:

```text
nil
```

can represent:

> The client didn't provide this field.

Whereas:

```go
name := ""
```

and a pointer to that value can represent:

> The client explicitly provided an empty string.

This is a useful application of something we learned much earlier about pointers.

We don't need to build PATCH handling yet, but this is an important reason you will see pointer fields in Go API request structs.

---

# 11. Don't confuse API models with database models

Suppose eventually we introduce PostgreSQL.

We might have:

```go
type User struct {
    ID           int
    Name         string
    Email        string
    PasswordHash string
}
```

and:

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

and:

```go
type UserResponse struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

These three types have different purposes.

| Type | Purpose |
|---|---|
| `CreateUserRequest` | What the client may send |
| `User` | Internal/domain/database representation |
| `UserResponse` | What the API exposes |

They might happen to contain overlapping fields, but that's not a reason to merge them.

At the same time, **don't automatically create ten different types for every endpoint**. If two representations genuinely have the same contract and semantics, sharing a type can be perfectly reasonable.

The principle is:

> **Separate representations when their responsibilities or contracts differ—not merely because separation is fashionable.**

---

# 12. A useful mental model

At this point, our API boundary is becoming more structured:

```text
                    HTTP
                     |
              +------+------+
              |             |
           Request       Response
              |             ^
              v             |
        Request Model    Response Model
              |             ^
              v             |
           Validate         |
              |             |
              v             |
       Application Logic ---+
              |
              v
        Domain / Data
```

The HTTP layer is responsible for translating between external representations and our application's typed world.

That includes:

- extracting path/query values
- parsing them into Go types
- decoding JSON
- validating input
- encoding responses
- selecting appropriate HTTP status codes

That boundary will become increasingly important as we build the API.

---

# Exercise

Extend the API we've been building with:

```text
POST /products
```

Define a request model for creating a product with:

- `name`
- `price`

Then:

- Decode the JSON body.
- Reject malformed JSON.
- Validate that `name` isn't empty.
- Validate that `price` is greater than zero.
- Return an appropriate `400` response for invalid input.
- For valid input, return a JSON response representing the product.

**Don't add a database yet.** You can construct the response directly in memory.

One thing to decide yourself: **should the client's create request contain the product ID?** Think about who should be responsible for generating an ID.

### Next: Lesson 31 — HTTP Responses in JSON: Status Codes, Headers & Consistent API Responses