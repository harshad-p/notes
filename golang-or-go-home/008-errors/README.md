Got it. I'll keep **useful observations and implications**, but remove repetition of concepts you've already understood.

# Lesson 16 — Error Handling

Go does not use exceptions in the usual C# sense for ordinary errors.

Instead, functions commonly return an `error` as one of their return values.

For example:

```go
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide by zero")
	}

	return a / b, nil
}
```

Calling it:

```go
result, err := divide(10, 2)

if err != nil {
	fmt.Println("Error:", err)
	return
}

fmt.Println(result)
```

Output:

```text
5
```

---

## `nil`

You'll see `nil` constantly in Go.

Here:

```go
return a / b, nil
```

`nil` means:

> There is no error.

So this:

```go
if err != nil
```

means:

> An error exists.

And:

```go
if err == nil
```

means:

> Everything is okay.

---

## Errors are values

This is an important Go idea.

`error` is actually an **interface**:

```go
type error interface {
	Error() string
}
```

So an error is simply a value that provides an `Error()` method.

That's why you can do:

```go
fmt.Println(err)
```

or:

```go
message := err.Error()
```

---

## A very common pattern

You'll see this constantly in real Go applications:

```go
user, err := getUser(id)
if err != nil {
	return err
}
```

Then continue using `user`.

This makes the error path explicit and immediately visible.

---

## One useful detail

You don't always have to create an error with `fmt.Errorf`.

For a simple fixed message:

```go
errors.New("user not found")
```

requires:

```go
import "errors"
```

`fmt.Errorf` becomes particularly useful when you want to include values:

```go
fmt.Errorf("user %d not found", id)
```

---

## Exercise

Write:

```go
func divide(a, b float64) (float64, error)
```

It should:

- return the result when `b` isn't zero
- return an error when `b` is zero

Then test both:

```go
divide(10, 2)
divide(10, 0)
```

For now, just :chatgpt-content-reference{index="0"}. We'll later learn how to create and handle custom errors more cleanly.

---

### Next lesson → **Lesson 17 — Packages and Modules**

We'll finally look at how Go projects are split into multiple files and packages—the foundation you'll need before we start building an actual backend.