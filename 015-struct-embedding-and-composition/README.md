# Lesson 23 — Struct Embedding & Composition

Go doesn't have classical inheritance. Instead, it uses **composition**, and struct embedding is the convenient syntax for it.

### 1. Embedding a struct

```go
type Address struct {
    City    string
    Country string
}

type User struct {
    Name string
    Address
}
```

Now:

```go
user := User{
    Name: "Harshad",
    Address: Address{
        City:    "Berlin",
        Country: "Germany",
    },
}

fmt.Println(user.Name)
fmt.Println(user.City)
```

Notice `user.City`.

You could write:

```go
user.Address.City
```

but because `Address` is **embedded**, Go promotes its fields, allowing:

```go
user.City
```

### 2. Methods are promoted too

```go
func (a Address) FullAddress() string {
    return a.City + ", " + a.Country
}
```

You can then do:

```go
fmt.Println(user.FullAddress())
```

The method belongs to `Address`, but `User` gets access to it through embedding.

### 3. This isn't inheritance

There is no "User **is an** Address" relationship here.

It's:

> User **has an** Address.

Embedding is simply a convenient way of expressing that composition.

This becomes particularly useful when building larger structs from reusable pieces.

### Exercise

Create:

```go
type Logger struct {
    Prefix string
}

func (l Logger) Log(message string) {
    fmt.Println(l.Prefix + message)
}
```

Then create a `Service` struct that embeds `Logger`.

Instantiate it with prefix `"SERVICE: "` and call:

```go
service.Log("starting...")
```

You should be able to call `Log()` directly on `service`.

**Next lesson → Lesson 24: Struct Tags & JSON**