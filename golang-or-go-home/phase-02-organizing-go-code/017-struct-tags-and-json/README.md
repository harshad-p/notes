# Lesson 24 — Struct Tags & JSON

This is particularly important for backend development because you'll constantly convert Go structs to/from JSON.

## 1. Struct tags

You can attach metadata to struct fields:

```go
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}
```

The part inside backticks is a **struct tag**.

Here it tells Go's JSON package:

> When this field is represented as JSON, call it `name`, not `Name`.

## 2. Struct → JSON

```go
user := User{
    Name: "Harshad",
    Age:  36,
}

data, err := json.Marshal(user)
if err != nil {
    panic(err)
}

fmt.Println(string(data))
```

Output:

```json
{"name":"Harshad","age":36}
```

You need:

```go
import "encoding/json"
```

`json.Marshal()` returns `[]byte`, which is why we use `string(data)` for printing.

## 3. JSON → Struct

The opposite is `json.Unmarshal()`:

```go
data := []byte(`{"name":"Harshad","age":36}`)

var user User

err := json.Unmarshal(data, &user)
if err != nil {
    panic(err)
}

fmt.Println(user.Name)
fmt.Println(user.Age)
```

Notice:

```go
json.Unmarshal(data, &user)
```

We pass a **pointer to `user`** because `Unmarshal` needs to modify it.

This is one of those places where the pointer concepts you've already learned become immediately useful.

## One useful tag option

You can omit empty fields:

```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email,omitempty"`
}
```

If `Email == ""`, the JSON won't contain `email` at all.

(No space between `email,omitempty`)

---

## Exercise

Given:

```go
type Product struct {
    Name  string `json:"name"`
    Price float64 `json:"price"`
}
```

Create a `Product`, convert it to JSON with `json.Marshal()`, and print the resulting JSON.

Then take this JSON:

```json
{"name":"Laptop","price":1299.99}
```

and use `json.Unmarshal()` to populate a `Product`.

**Next lesson → Lesson 25: JSON in HTTP APIs — request and response bodies**