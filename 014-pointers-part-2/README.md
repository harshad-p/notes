# Lesson 22 — Pointers, Part 2: `*` and `&` in practice

We covered pointers conceptually in Lesson 14. Now let's make sure the syntax is completely comfortable, because you'll encounter it constantly in Go APIs and structs.

## 1. `&` — get the address

```go
user := User{Name: "Harshad", Age: 36}

p := &user
```

- `user` is the actual `User`
- `p` is a `*User` — a pointer to that `User`
- `&user` means **"give me the address of user"**

## 2. `*` — dereference

```go
fmt.Println(p.Age)
```

Go automatically dereferences `p` when accessing fields.

You can also explicitly dereference:

```go
fmt.Println((*p).Age)
```

Both work.

More importantly:

```go
p.Age = 40
```

changes the original `user`:

```go
fmt.Println(user.Age) // 40
```

## 3. `*` also appears in the type

This:

```go
var p *User
```

means:

> `p` is a pointer to a `User`.

Whereas this:

```go
var user User
```

means:

> `user` is a `User`.

So the same `*` has two closely related uses:

```go
*User      // pointer type
*p         // dereference pointer
```

## One subtle Go detail

You normally don't need to manually dereference struct pointers:

```go
p.Age
```

instead of:

```go
(*p).Age
```

Go handles this automatically for field/method access.

## Exercise

Given:

```go
type Counter struct {
    Value int
}

func main() {
    c := Counter{Value: 10}

    p := &c

    // change c.Value to 25 using p
    // print c.Value
}
```

Write the two missing lines.

**Next lesson → Lesson 23: Struct Embedding & Composition**