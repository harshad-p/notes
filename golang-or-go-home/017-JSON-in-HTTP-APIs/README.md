# Lesson 25 — JSON in HTTP APIs: Request & Response Bodies

We've already learned two separate things:

1. Go structs can be converted to JSON with `json.Marshal()`.
2. Go structs can be created from JSON with `json.Unmarshal()`.

Now we're going to connect those concepts to HTTP.

The goal of this lesson is to understand what actually happens when a client communicates with a Go server.

---

## 1. First: what is an HTTP request?

When your browser, Postman, or another application calls:

```text
GET http://localhost:8080/user
```

it sends an **HTTP request** to your server.

An HTTP request contains several important pieces:

```text
Request
├── Method
├── URL
├── Headers
└── Body
```

For example, a request might conceptually look like:

```text
POST /users HTTP/1.1
Content-Type: application/json

{
    "name": "Harshad",
    "age": 36
}
```

Here:

- `POST` is the **method**
- `/users` is the **path**
- `Content-Type: application/json` is a **header**
- the JSON is the **body**

---

# 2. Starting an HTTP server

Go's standard library contains the `net/http` package.

```go
import "net/http"
```

The package provides functionality for creating HTTP servers and clients.

Our simplest server is:

```go
func main() {
    http.ListenAndServe(":8080", nil)
}
```

### What does `ListenAndServe` mean?

```go
http.ListenAndServe(":8080", nil)
```

tells Go:

> Start an HTTP server and listen for incoming connections on port `8080`.

So when you visit:

```text
http://localhost:8080
```

your request reaches this server.

The first argument:

```go
":8080"
```

specifies the address/port to listen on.

The second argument:

```go
nil
```

is related to how requests are routed. We're not going to rely on that detail yet; we'll introduce routing properly as we need it.

---

# 3. A server needs a handler

A server needs something that knows what to do when a request arrives.

For example:

```go
func helloHandler(w http.ResponseWriter, r *http.Request) {
    // handle request
}
```

This function has two parameters:

```go
w http.ResponseWriter
r *http.Request
```

These are important enough that we should understand them properly.

---

# 4. What is `*http.Request`?

`http.Request` is a struct provided by Go's `net/http` package.

It represents an incoming HTTP request.

Conceptually, you can think of it as Go taking the HTTP request:

```text
POST /users HTTP/1.1
Content-Type: application/json

{"name":"Harshad","age":36}
```

and making its information available through a Go value:

```go
r
```

That value contains information such as:

```go
r.Method
r.URL
r.Header
r.Body
```

For example:

```go
fmt.Println(r.Method)
```

might print:

```text
POST
```

And:

```go
fmt.Println(r.URL.Path)
```

might print:

```text
/users
```

We will explore these properties in more detail in later lessons.

For this lesson, the important thing is:

> `r` represents the incoming HTTP request.

---

# 5. What is `http.ResponseWriter`?

Now the other parameter:

```go
w http.ResponseWriter
```

`http.ResponseWriter` is an interface provided by `net/http`.

It represents the mechanism your handler uses to construct the **HTTP response** that will be sent back to the client.

For example:

```go
fmt.Fprintln(w, "Hello!")
```

writes `"Hello!"` into the HTTP response body.

### What is `fmt.Fprintln`?

This is from the `fmt` package.

You may already know:

```go
fmt.Println("Hello")
```

prints text to the program's standard output, normally your terminal.

`Fprintln` is similar, but instead of always writing to the terminal, it takes a destination as its first argument:

```go
fmt.Fprintln(destination, value)
```

For example:

```go
fmt.Fprintln(os.Stdout, "Hello")
```

writes to standard output.

And:

```go
fmt.Fprintln(w, "Hello")
```

writes to the destination represented by `w`.

In an HTTP handler, `w` represents the HTTP response, so the text becomes part of the response sent to the client.

This distinction is important:

```go
fmt.Println("Hello")
```

→ terminal

```go
fmt.Fprintln(w, "Hello")
```

→ HTTP response

---

# 6. Let's build our first real handler

Put this together:

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

There are a few things here we haven't formally studied yet, particularly:

```go
http.HandleFunc(...)
```

We'll cover routing properly in Lesson 27/26 territory, so for now just understand its purpose:

> It associates `/hello` with `helloHandler`.

When you visit:

```text
http://localhost:8080/hello
```

Go calls:

```go
helloHandler(w, r)
```

and:

```go
fmt.Fprintln(w, "Hello from Go!")
```

creates the response body.

---

# 7. Now let's introduce JSON

We already learned:

```go
json.Marshal(user)
```

and:

```go
json.Unmarshal(data, &user)
```

Those are useful when you have JSON data in memory.

But HTTP APIs deal with **streams of data**.

For example:

```text
Client
   |
   | HTTP request
   | JSON body
   ↓
Go server
```

and:

```text
Go server
   |
   | HTTP response
   | JSON body
   ↓
Client
```

Go therefore provides another pair of tools:

```go
json.NewDecoder(...)
json.NewEncoder(...)
```

Before using them, let's understand what a decoder and encoder are.

---

# 8. Encoder vs Decoder

An **encoder** takes a Go value and converts it into another representation.

For JSON:

```text
Go value
   ↓
JSON encoder
   ↓
JSON
```

A **decoder** does the opposite:

```text
JSON
   ↓
JSON decoder
   ↓
Go value
```

So:

```go
json.NewEncoder(...)
```

is useful when we're sending JSON **out**.

And:

```go
json.NewDecoder(...)
```

is useful when we're receiving JSON **in**.

---

# 9. `json.NewEncoder`

Suppose we have:

```go
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}
```

We could create:

```go
user := User{
    Name: "Harshad",
    Age:  36,
}
```

Now we want to send this user as an HTTP response.

We could use `json.Marshal`:

```go
data, err := json.Marshal(user)
```

which produces JSON data in memory.

But we don't actually need to keep the JSON in memory first.

We already have a destination:

```go
w
```

So we can create an encoder associated with that destination:

```go
encoder := json.NewEncoder(w)
```

Then:

```go
encoder.Encode(user)
```

means:

> Convert `user` into JSON and write the resulting JSON to `w`.

Since `w` is the HTTP response writer, the JSON becomes the HTTP response body.

We can therefore write:

```go
func userHandler(w http.ResponseWriter, r *http.Request) {
    user := User{
        Name: "Harshad",
        Age:  36,
    }

    encoder := json.NewEncoder(w)
    encoder.Encode(user)
}
```

The client receives approximately:

```json
{"name":"Harshad","age":36}
```

`Encode` also writes a newline after the JSON value, which is normally harmless for an HTTP JSON response.

---

# 10. Tell the client that this is JSON

There's another important piece.

The HTTP response has **headers**.

One of those headers tells the client what kind of data we're returning.

For JSON, we normally use:

```text
Content-Type: application/json
```

In Go:

```go
w.Header().Set("Content-Type", "application/json")
```

Let's break that down.

### `w.Header()`

`Header()` gives us the response headers.

### `.Set(...)`

`Set` assigns a value to a header.

So:

```go
w.Header().Set("Content-Type", "application/json")
```

means:

> Set the response's `Content-Type` header to `application/json`.

Now our handler becomes:

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

The response is conceptually:

```text
HTTP response

Content-Type: application/json

{"name":"Harshad","age":36}
```

That's a proper JSON HTTP response.

---

# 11. Now the other direction: receiving JSON

Suppose a client sends:

```text
POST /users

Content-Type: application/json

{
    "name": "Harshad",
    "age": 36
}
```

The JSON exists inside the request body.

Go exposes that body through:

```go
r.Body
```

`r.Body` is an `io.ReadCloser`.

That's a type from Go's standard library representing something that can be **read from** and then **closed**.

You don't need to memorize `io.ReadCloser` yet.

For now, the important mental model is:

> `r.Body` is the incoming data stream containing the HTTP request body.

---

# 12. Using `json.NewDecoder` with `r.Body`

We want to turn:

```text
JSON request body
```

into:

```go
User
```

So we use a JSON decoder.

```go
decoder := json.NewDecoder(r.Body)
```

This creates a decoder that reads JSON from the request body.

Then:

```go
err := decoder.Decode(&user)
```

reads JSON from the body and populates `user`.

Complete example:

```go
func createUserHandler(w http.ResponseWriter, r *http.Request) {
    var user User

    decoder := json.NewDecoder(r.Body)

    err := decoder.Decode(&user)
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    fmt.Println(user.Name)
    fmt.Println(user.Age)
}
```

If the client sends:

```json
{
    "name": "Harshad",
    "age": 36
}
```

then after:

```go
decoder.Decode(&user)
```

we have:

```go
user.Name // "Harshad"
user.Age  // 36
```

---

# 13. Why `&user`?

This is directly connected to our pointer lessons.

We have:

```go
var user User
```

At this point, `user` is an empty `User`.

The decoder needs to **modify** it.

So we give it the address:

```go
&user
```

Conceptually:

```text
user
┌──────────────────┐
│ Name: ""         │
│ Age: 0           │
└──────────────────┘
       ↑
       │
      &user
```

The decoder can then populate that existing `User`.

After decoding:

```text
user
┌──────────────────┐
│ Name: "Harshad"  │
│ Age: 36          │
└──────────────────┘
```

This is why the pointer knowledge from earlier matters in real Go code.

---

# 14. Handling invalid JSON

`Decode()` returns an error:

```go
err := decoder.Decode(&user)
```

So we must handle it:

```go
if err != nil {
    http.Error(w, "Invalid JSON", http.StatusBadRequest)
    return
}
```

We've already learned ordinary Go error handling:

```go
value, err := something()
if err != nil {
    // handle error
}
```

Same principle here.

---

# 15. What is `http.Error`?

We're using another function, so let's actually define it.

```go
http.Error(w, "Invalid JSON", http.StatusBadRequest)
```

`http.Error` is a convenience function from `net/http` for sending an HTTP error response.

The arguments are:

```go
http.Error(
    responseWriter,
    errorMessage,
    statusCode,
)
```

So:

```go
http.Error(w, "Invalid JSON", http.StatusBadRequest)
```

sends an HTTP `400 Bad Request` response with the supplied error message.

The:

```go
return
```

is important because after sending the error, we don't want the handler to continue processing the invalid request.

---

# 16. Putting request and response together

Now we can build a small API that does both directions.

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
    user := User{
        Name: "Harshad",
        Age:  36,
    }

    w.Header().Set("Content-Type", "application/json")

    err := json.NewEncoder(w).Encode(user)
    if err != nil {
        return
    }
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
    var user User

    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    fmt.Println("Name:", user.Name)
    fmt.Println("Age:", user.Age)

    w.Header().Set("Content-Type", "application/json")

    json.NewEncoder(w).Encode(user)
}

func main() {
    http.HandleFunc("/user", getUserHandler)
    http.HandleFunc("/create-user", createUserHandler)

    http.ListenAndServe(":8080", nil)
}
```

Don't worry about every detail of this program yet. The important part of this lesson is understanding the **data flow**.

### GET

```text
Client
   │
   │ GET /user
   ↓
getUserHandler
   │
   │ Go User
   ↓
json.NewEncoder(w)
   │
   │ JSON
   ↓
HTTP response
```

### POST

```text
Client
   │
   │ POST /create-user
   │ JSON body
   ↓
r.Body
   │
   ↓
json.NewDecoder(r.Body)
   │
   ↓
Go User
```

That's the fundamental JSON mechanism behind a huge number of Go REST APIs.

---

# 17. Your exercise — build it yourself

Don't copy the complete example above.

Create a fresh program.

### Step 1 — Create the struct

Create:

```go
type Product struct {
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}
```

### Step 2 — Create a GET endpoint

Create:

```text
GET /product
```

It should return:

```json
{
    "name": "Laptop",
    "price": 1299.99
}
```

Requirements:

- Create a `Product`
- Set `Content-Type`
- Use `json.NewEncoder`
- Encode the product into the response

### Step 3 — Create a POST endpoint

Create:

```text
POST /product
```

The client will send:

```json
{
    "name": "Keyboard",
    "price": 99.99
}
```

Your handler should:

1. Create an empty `Product`.
2. Read the request body.
3. Decode the JSON into the `Product`.
4. Handle a decoding error with `400 Bad Request`.
5. Print the product's name and price.
6. Return the decoded product as JSON.

### Step 4 — Test it

Use Postman, `curl`, or another HTTP client.

For the POST request, send:

```json
{
    "name": "Keyboard",
    "price": 99.99
}
```

Then deliberately send malformed JSON such as:

```json
{
    "name": "Keyboard",
    "price":
}
```

and verify that your server returns a `400`.

### The concepts you should now be able to explain

Before we move on, you should understand the purpose of each of these:

```go
http.Request
http.ResponseWriter
r.Body
w.Header()
fmt.Fprintln()
json.NewEncoder()
json.NewDecoder()
Encode()
Decode()
http.Error()
```

We're **not** going to move on simply because you've seen the syntax. The point of this lesson is that you understand what happens to JSON as it travels between a client and your Go program.

**Next lesson → Lesson 26: HTTP Routing, Methods & Status Codes**
