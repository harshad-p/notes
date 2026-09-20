Perfect. Before we get into collections like slices and maps, you need to know how Go controls the flow of execution.

The good news: if you know C#, **90% of this will feel familiar**.

# Lesson 8 — If/Else and Loops

## 1. The `if` statement

Let's start with something simple.

```go
age := 20

if age >= 18 {
	fmt.Println("Adult")
}
```

Output:

```text
Adult
```

Looks almost identical to C#.

---

## 2. `if...else`

```go
age := 16

if age >= 18 {
	fmt.Println("Adult")
} else {
	fmt.Println("Minor")
}
```

Output:

```text
Minor
```

Again, very similar to C#.

---

## ⭐ Difference #1 — No parentheses

In C#, you'd write:

```csharp
if (age >= 18)
{
}
```

In Go:

```go
if age >= 18 {
}
```

Notice:

- ❌ No `()`
- ✅ Curly braces are still required.

This is one of the first syntax differences you'll notice.

---

## 3. `for` loop

Go has **only one looping keyword**:

```text
for
```

There is **no**:

- `while`
- `do...while`
- `foreach`

Everything is done with `for`.

A classic loop:

```go
for i := 1; i <= 5; i++ {
	fmt.Println(i)
}
```

Output:

```text
1
2
3
4
5
```

This should look very familiar.

---

## ⭐ Difference #2 — No `while`

In C#:

```csharp
while (count < 5)
{
    count++;
}
```

In Go:

```go
count := 0

for count < 5 {
	count++
}
```

The `for` loop acts like a `while` loop when you omit the initialization and increment sections.

---

## 4. Infinite loop

Need a loop that never ends?

```go
for {
	fmt.Println("Running...")
}
```

Equivalent to:

```csharp
while (true)
{
}
```

---

## Exercise

Write a program that prints the numbers **1 to 10**.

Then modify it to print **only even numbers**.

**Hint:** Use the remainder operator `%`.

Example:

```go
if i%2 == 0 {
	// even number
}
```

Expected output:

```text
2
4
6
8
10
```

---

## 🧠 What you've learned today

- `if` and `else` work almost like C#.
- Go **doesn't use parentheses** around conditions.
- Go has **only one loop keyword**: `for`.
- `for` can replace `for`, `while`, and `do...while` from C#.

Go's philosophy is to have **fewer language features** while still being expressive. Instead of giving you three different loop constructs, it gives you one flexible one.

---

### Next lesson → **Lesson 9 — Slices: Go's most important collection type**

This is an important lesson because slices are used everywhere in Go, much like `List<T>` is in C#.