# Lesson 29 — Type Conversion & Parsing: Turning URL Strings into Go Types

In Lesson 28, we saw that path parameters and query parameters arrive as **strings**.

For example:

```go
id := r.PathValue("id")
```

For:

```text
GET /users/42
```

`id` contains:

```text
"42"
```

That is fine if we only need to display the ID. But in a real API, we often need to **use the value as a number**.

For example:

```go
userID := 42
```

might be used to query a database:

```sql
SELECT * FROM users WHERE id = 42
```

So we need to turn the string `"42"` into an `int`.

That is the focus of this lesson.

---

## 1. Type conversion vs parsing

There are two related ideas in Go that are worth distinguishing.

### Type conversion

A **type conversion** changes a value from one Go type to another compatible Go type.

For example:

```go
var x int = 42
var y float64 = float64(x)
```

Here, `x` is already a valid numeric value. We are simply converting the `int` value `42` into a `float64` value `42.0`.

There is no text involved.

### Parsing

**Parsing** means taking some representation—often text—and interpreting it as a value of a particular type.

For example:

```text
"42"
```

is text.

We want to interpret that text as the number:

```text
42
```

That's parsing.

This distinction matters for HTTP APIs because URL path and query values are **text**, so we're generally parsing them rather than simply converting them.

---

# 2. The `strconv` package

Go's standard library provides the `strconv` package for converting between strings and basic values.

`strconv` stands for **string conversion**.

We import it like this:

```go
import "strconv"
```

It provides functions for parsing things such as:

- integers
- floating-point numbers
- booleans

and also for converting values back into strings.

For this lesson, we'll primarily use:

```go
strconv.Atoi
```

---

# 3. Parsing a string into an `int`

`strconv.Atoi` means **ASCII to integer**.

It takes a string and attempts to interpret it as an `int`.

For example:

```go
id, err := strconv.Atoi("42")
```

After this:

```text
id  = 42
err = nil
```

The important thing is that `Atoi` returns **two values**:

```go
int, error
```

Why?

Because the string might not actually represent a valid integer.

For example:

```go
id, err := strconv.Atoi("abc")
```

There is no meaningful integer represented by `"abc"`.

So instead of silently producing some strange value, Go gives us an error.

---

# 4. Handling the error

This is especially important in an HTTP handler.

Suppose we have:

```go
mux.HandleFunc("GET /users/{id}", getUserHandler)
```

and:

```go
func getUserHandler(w http.ResponseWriter, r *http.Request) {
    idText := r.PathValue("id")

    id, err := strconv.Atoi(idText)

    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    fmt.Fprintln(w, "User ID:", id)
}
```

Let's go through the new part carefully.

### First:

```go
idText := r.PathValue("id")
```

`idText` is a `string`.

For:

```text
GET /users/42
```

it contains:

```text
"42"
```

### Then:

```go
id, err := strconv.Atoi(idText)
```

`Atoi` tries to parse `"42"` as an integer.

If successful:

```text
id  = 42
err = nil
```

So `id` is now an actual `int`.

### Then:

```go
if err != nil {
```

We check whether parsing failed.

Remember from our earlier error-handling lesson:

```go
nil
```

means there is no error.

Therefore:

```go
err != nil
```

means an error occurred.

If the URL contained something that isn't a valid integer, we send:

```go
http.Error(w, "Invalid user ID", http.StatusBadRequest)
```

and:

```go
return
```

stops the handler immediately.

This is important because we don't want to continue processing with an invalid ID.

---

# 5. Why HTTP APIs need this

Consider an endpoint:

```text
GET /users/{id}
```

We might eventually have:

```go
func getUserHandler(w http.ResponseWriter, r *http.Request) {
    idText := r.PathValue("id")

    id, err := strconv.Atoi(idText)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    user, err := repository.GetUser(id)
    // ...
}
```

Our repository might have:

```go
GetUser(id int)
```

So parsing forms the boundary between the **HTTP world**, where values arrive as strings, and our **application code**, where we want strongly typed values.

Conceptually:

```text
HTTP request
     |
     |  "42"
     v
  Handler
     |
     |  42
     v
Application / Repository
```

This is a very common pattern in backend development.

---

# 6. Invalid input vs valid but unacceptable input

There's another important distinction.

Suppose our endpoint expects a positive user ID.

These are different problems:

```text
/users/abc
```

and:

```text
/users/-5
```

For `"abc"`:

```go
strconv.Atoi("abc")
```

fails.

That's a **parsing error**.

But:

```go
strconv.Atoi("-5")
```

succeeds.

We get:

```go
id = -5
err = nil
```

So parsing was successful.

However, `-5` might not be a valid user ID according to our application's rules.

We therefore need a separate validation step:

```go
if id <= 0 {
    http.Error(w, "User ID must be positive", http.StatusBadRequest)
    return
}
```

This gives us two stages:

```text
"abc"
  |
  | parsing
  X
invalid integer


"-5"
  |
  | parsing
  v
-5
  |
  | validation
  X
invalid user ID
```

This distinction becomes important in real APIs.

---

# 7. Parsing query parameters

The same issue occurs with query parameters.

Suppose we have:

```text
GET /users?limit=20
```

We learned previously that:

```go
r.URL.Query().Get("limit")
```

returns a string.

So:

```go
limitText := r.URL.Query().Get("limit")
```

gives:

```text
"20"
```

If we want an integer:

```go
limit, err := strconv.Atoi(limitText)
```

Now `limit` is an `int`.

We can then validate it:

```go
if limit <= 0 {
    http.Error(w, "Limit must be positive", http.StatusBadRequest)
    return
}
```

This pattern is extremely common:

```text
query/path value
       ↓
     string
       ↓
     parse
       ↓
   Go value
       ↓
   validate
       ↓
 business logic
```

---

# 8. What if the query parameter is missing?

This is another case we need to think about.

Suppose the request is:

```text
GET /users
```

There is no `limit`.

Then:

```go
limitText := r.URL.Query().Get("limit")
```

returns an empty string:

```text
""
```

If we then do:

```go
limit, err := strconv.Atoi(limitText)
```

parsing fails.

Whether that's an error depends on our API design.

For example, we might decide:

> If `limit` isn't supplied, use 20.

Then we'd handle the missing value before parsing:

```go
limitText := r.URL.Query().Get("limit")

limit := 20

if limitText != "" {
    parsedLimit, err := strconv.Atoi(limitText)

    if err != nil {
        http.Error(w, "Invalid limit", http.StatusBadRequest)
        return
    }

    limit = parsedLimit
}
```

Notice something important here.

We are not treating **missing** and **invalid** as the same thing.

```text
?limit=20     → supplied and valid
?limit=abc    → supplied but invalid
(no limit)    → not supplied
```

The API can give each case different behavior.

---

# 9. Other useful parsing functions

`strconv` isn't limited to integers.

## `ParseBool`

For a string such as:

```text
"true"
```

we can use:

```go
value, err := strconv.ParseBool("true")
```

It returns:

```text
bool, error
```

This is useful for query parameters such as:

```text
GET /users?includeOrders=true
```

because query parameters are strings, but our application might want:

```go
includeOrders bool
```

---

## `ParseFloat`

For decimal numbers:

```go
price, err := strconv.ParseFloat("19.99", 64)
```

The `64` tells Go that we want a 64-bit floating-point value, so the result is a `float64`.

Again, it returns:

```text
float64, error
```

because `"19.99"` may be valid, while `"hello"` isn't.

---

# 10. `Atoi` vs `ParseInt`

You will also encounter:

```go
strconv.ParseInt
```

For example:

```go
id, err := strconv.ParseInt("42", 10, 64)
```

This is more configurable than `Atoi`.

The arguments mean:

```text
"42"   → string being parsed
10     → base 10
64     → result should fit in 64 bits
```

The result is an `int64`, not an `int`.

For ordinary HTTP IDs where we simply want an `int`, `Atoi` is usually the simpler choice:

```go
id, err := strconv.Atoi(idText)
```

You don't need `ParseInt` just because it exists.

Use the simpler function when its behavior matches what you need.

---

# 11. Type conversion doesn't replace parsing

It's worth seeing why this doesn't work:

```go
id := int(idText)
```

If `idText` is:

```go
string
```

you cannot use `int(...)` to interpret the characters `"42"` as the number `42`.

Type conversion isn't a general-purpose "turn this text into a value" mechanism.

The string contains characters:

```text
'4' '2'
```

Parsing is what interprets those characters as the numeric value:

```text
42
```

So for URL data:

```go
id, err := strconv.Atoi(idText)
```

is the appropriate operation.

---

# 12. A complete example

Let's put the concepts together.

Suppose we have:

```text
GET /products/{id}?limit=10
```

The handler could look like:

```go
func getProductHandler(w http.ResponseWriter, r *http.Request) {
    idText := r.PathValue("id")

    id, err := strconv.Atoi(idText)
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }

    if id <= 0 {
        http.Error(w, "Product ID must be positive", http.StatusBadRequest)
        return
    }

    limitText := r.URL.Query().Get("limit")

    limit := 20

    if limitText != "" {
        parsedLimit, err := strconv.Atoi(limitText)

        if err != nil {
            http.Error(w, "Invalid limit", http.StatusBadRequest)
            return
        }

        limit = parsedLimit
    }

    fmt.Fprintln(w, "Product ID:", id)
    fmt.Fprintln(w, "Limit:", limit)
}
```

The important parts are not the number of lines, but the sequence:

```text
Path/query
   ↓
string
   ↓
parse
   ↓
typed value
   ↓
validate
   ↓
use value
```

This pattern will appear repeatedly as we build APIs.

---

# Exercise

Create a small API with this route:

```text
GET /products/{id}
```

Use the modern `ServeMux` routing from Lesson 28.

Your handler should:

1. Get `id` using `r.PathValue("id")`.
2. Parse it into an `int` using `strconv.Atoi`.
3. Return `400 Bad Request` if it isn't a valid integer.
4. Return `400 Bad Request` if the ID is `0` or negative.
5. Otherwise return something like:

```text
Product ID: 42
```

Test at least:

```text
GET /products/42
GET /products/abc
GET /products/-5
```

The key thing I want you to practice is the distinction between:

```text
parsing
```

and:

```text
validation
```

Don't add a database yet. We haven't introduced that part of the course.

### Next: Lesson 30 — JSON API Design: Request/Response Models & Validation