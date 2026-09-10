# Lesson 26 — HTTP Routing, Methods & Status Codes

Our previous server accepted requests, but we haven't yet distinguished **HTTP methods** or deliberately returned **status codes**.

### 1. Checking the HTTP method

A handler receives every method sent to its route unless you restrict it:

```go
func userHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // handle GET
}
```

Go provides constants such as:

```go
http.MethodGet
http.MethodPost
http.MethodPut
http.MethodDelete
http.MethodPatch
```

Prefer these over writing `"GET"` or `"POST"` yourself.

---

### 2. Status codes

You can explicitly set a response status:

```go
w.WriteHeader(http.StatusCreated)
```

Common ones:

| Status | Meaning |
|---|---|
| `200 OK` | Successful request |
| `201 Created` | Resource successfully created |
| `204 No Content` | Successful, but no response body |
| `400 Bad Request` | Invalid request |
| `401 Unauthorized` | Authentication required/failed |
| `403 Forbidden` | Not allowed |
| `404 Not Found` | Resource doesn't exist |
| `405 Method Not Allowed` | HTTP method isn't supported |
| `500 Internal Server Error` | Server-side failure |

For example:

```go
func createUserHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // create user...

    w.WriteHeader(http.StatusCreated)
}
```

### 3. One subtle point

`WriteHeader()` should generally happen **before writing the body**:

```go
w.WriteHeader(http.StatusCreated)
fmt.Fprintln(w, "User created")
```

Once you start writing the response body, Go may already have sent the headers, making it too late to change the status.

---

### Exercise

Create two handlers:

```text
GET  /users
POST /users
```

For `/users`:

- GET → return `"Getting users"` with `200`
- POST → return `"Creating user"` with `201`
- anything else → return `"Method not allowed"` with `405`

You can use separate handlers or one handler that switches on `r.Method`. Try the latter first.

**Next lesson → Lesson 27: URL Paths, Query Parameters & Path Parameters**