# Lesson 19 — `defer`

`defer` is used when you want something to happen **when the current function is about to finish**.

A simple example:

```go
func main() {
	defer fmt.Println("Goodbye")

	fmt.Println("Hello")
}
```

Output:

```text
Hello
Goodbye
```

Even though `defer` appears first, the deferred call happens at the **end of `main()`**.

---

## Why is this useful?

Its most common purpose is **cleanup**.

For example, when working with a file:

```go
file, err := os.Open("data.txt")
if err != nil {
	return
}

defer file.Close()
```

Now you don't have to remember to call:

```go
file.Close()
```

at every possible exit point of the function.

Whether the function finishes normally or returns early, the deferred `Close()` runs.

You'll see the same pattern with things like:

- files
- database connections
- HTTP response bodies
- mutexes

---

## Multiple `defer`s

This is an interesting detail:

```go
func main() {
	defer fmt.Println("First")
	defer fmt.Println("Second")
	defer fmt.Println("Third")
}
```

Output:

```text
Third
Second
First
```

Deferred calls execute in **LIFO order** — *last in, first out*.

Think of it as a stack.

---

## One important detail

The arguments to a deferred function are evaluated **when `defer` is executed**, not when the function eventually runs.

For example:

```go
func main() {
	x := 10

	defer fmt.Println(x)

	x = 20
}
```

Output:

```text
10
```

`x` was evaluated when this happened:

```go
defer fmt.Println(x)
```

This is a subtle detail, but you'll encounter it when debugging deferred code.

---

## Exercise

Write a function:

```go
func test() {
	defer fmt.Println("Function finished")

	fmt.Println("Doing some work")
}
```

Call `test()` from `main()`.

Then add a second `defer` and predict which message will appear first.

---

### Next lesson → **Lesson 20 — Go Modules & External Dependencies**

We'll use a real external package and see what `go get`, `go.mod`, and `go.sum` actually do.