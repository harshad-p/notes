# Lesson 25 — JSON in HTTP APIs: Request & Response Bodies

Now we connect the JSON knowledge from Lesson 24 to something you'll actually build: an HTTP API.

Go's standard library gives us `net/http`.

### 1. A basic HTTP server

```go
package main

import (
    "fmt"
    "net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello from Go!")
}

func main() {
    http.HandleFunc("/hello", helloHandler)

    http.ListenAndServe(":8080", nil)
}
```

Run it and visit:

```text
http://localhost:8080/hello
```

The important pieces:

```go
func helloHandler(w http.ResponseWriter, r *http.Request)
```

- `r` contains the incoming HTTP request.
- `w` is how you write the HTTP response.

So:

```go
fmt.Fprintln(w, "Hello from Go!")
```

writes directly into the HTTP response body.

---

### 2. Returning JSON

Suppose we have:

```go
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}
```

We can return it:

```go
func userHandler(w http.ResponseWriter, r *http.Request) {
    user := User{
        Name: "Harshad",
        Age:  36,
    }

    w.Header().Set("Content-Type", "application/json")

    json.NewEncoder(w).Encode(user)
}
```

There is an interesting difference from what we did previously.

Previously:

```go
data, err := json.Marshal(user)
```

created JSON as `[]byte`.

Here:

```go
json.NewEncoder(w).Encode(user)
```

encodes the object **directly into the response writer**.

That's very convenient for HTTP responses.

---

### 3. Reading JSON from a request

For a `POST` request:

```go
func createUserHandler(w http.ResponseWriter, r *http.Request) {
    var user User

    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    fmt.Println(user.Name)
}
```

The flow is:

```text
HTTP request
     ↓
r.Body
     ↓
json.Decoder
     ↓
Go struct
```

And the response goes the opposite direction:

```text
Go struct
     ↓
json.Encoder
     ↓
HTTP response
```

This request/response JSON flow is the foundation we'll use when we eventually build a proper REST API.

### Exercise

Create an endpoint:

```text
GET /user
```

that returns:

```json
{
  "name": "Harshad",
  "age": 36
}
```

Use:

- `http.HandleFunc`
- a `User` struct
- `json.NewEncoder`
- `Content-Type: application/json`

Don't worry about POST yet.

**Next lesson → Lesson 26: HTTP Routing, Methods & Status Codes**