# Lesson 3 — Variables

Go gives you **two common ways** to declare variables.

### 1. Explicitly with `var`

```go
var name string = "Harshad"
var age int = 36
```

The type is explicitly specified.

### 2. Using `:=`

```go
name := "Harshad"
age := 36
```

Go figures out the types automatically.

This is called **short variable declaration** and you'll see it *everywhere* in Go code.

For example:

```go
name := "Harshad"   // string
age := 36           // int
developer := true   // bool
```

Think of `:=` roughly as:

> "Create this variable and figure out its type."

### Important difference

`:=` can only be used **when creating a new variable**.

```go
name := "Harshad"

name = "John"       // ✅ assignment
name := "John"      // ❌ already declared
```

This distinction is worth remembering.

---

## Your exercise

Replace your current `main()` with:

```go
func main() {
	name := "Harshad"
	age := 36

	fmt.Println(name)
	fmt.Println(age)
}
```

Run:

```bash
go run .
```

You should get:

```text
Harshad
36
```

Then change `age` and run it again.

**Don't worry about types beyond this yet.** We'll explore them next.

### Next lesson → **Lesson 4 — Basic Types (`string`, `int`, `float64`, `bool`)**

Awesome. Before we learn more syntax, let's answer **why Go has types**.

# Lesson 4 — Basic Types

Go is a **statically typed** language.

That means once a variable has a type, it cannot change.

Example:

```go
name := "Harshad"
```

Go knows:

```text
name → string
```

So this is **not allowed**:

```go
name = 123
```

The compiler will complain because `name` is a `string`, not an `int`.

---

## The four basic types you'll use most

### 1. `string`

Text.

```go
name := "Harshad"
city := "Berlin"
```

---

### 2. `int`

Whole numbers.

```go
age := 36
count := 100
```

---

### 3. `float64`

Numbers with decimals.

```go
price := 19.99
pi := 3.14159
```

You might wonder why it's called `float64` instead of just `float`.

We'll answer that in a later lesson when we discuss number types. For now, just know:

> Decimal numbers are usually `float64`.

---

### 4. `bool`

True or false.

```go
isAdmin := true
isLoggedIn := false
```

---

## Printing multiple values

Instead of:

```go
fmt.Println(name)
fmt.Println(age)
```

You can do:

```go
fmt.Println(name, age)
```

Output:

```text
Harshad 36
```

`Println` automatically puts a space between arguments.

---

## Exercise

Replace your `main()` with:

```go
func main() {
	name := "Harshad"
	age := 36
	height := 5.9
	isDeveloper := true

	fmt.Println(name)
	fmt.Println(age)
	fmt.Println(height)
	fmt.Println(isDeveloper)
}
```

Then try this:

```go
fmt.Println(name, age, height, isDeveloper)
```

Observe the output.

---

## 💡 Tiny challenge

Without running it first, predict the output of:

```go
fmt.Println("Age:", age)
```

Will it print:

```text
Age:36
```

or

```text
Age: 36
```

Make your guess, run it, and tell me what happened.

---

### Next lesson → **Lesson 5 — Constants (`const`) and why Go has both variables and constants**
Correct. 👍

Go's `fmt.Println()` automatically separates multiple arguments with a space.

So:

```go
fmt.Println("Age:", age)
```

is treated as two values:

1. `"Age:"`
2. `age`

and prints:

```text
Age: 36
```

If you wanted **no space**, you would need a different approach (we'll see `fmt.Print` and `fmt.Printf` later).

---

One small thing before moving on: notice the pattern you are learning:

```go
name := "Harshad"
age := 36
```

Go is doing type inference:

```text
name → string
age  → int
```

but the compiler still knows the types.

This gives Go a nice balance:
- less typing than Java/C#
- still catches type mistakes before running

---

# Lesson 5 — Constants

A variable can change:

```go
age := 36

age = 37
```

No problem.

A **constant** cannot change:

```go
const pi = 3.14159
```

After this:

```go
pi = 4
```

Go will reject it.

---

## Why use constants?

For values that should never change.

Examples:

```go
const companyName = "Atain"
const maxRetries = 3
const taxRate = 0.19
```

Imagine a backend API:

```go
const maxLoginAttempts = 5
```

You don't want someone accidentally changing it while the program is running.

---

## Variables vs Constants

| | Variable | Constant |
|-|-|-|
| Keyword | `var` / `:=` | `const` |
| Can change? | ✅ Yes | ❌ No |
| Example | user age | tax rate |

---

## Important Go rule

Constants must have a value immediately:

Works:

```go
const appName = "My API"
```

Does not work:

```go
const appName string
```

because Go needs to know the value.

---

## Exercise

Create:

```go
const appName = "My First Go API"

version := 1
```

Print:

```
My First Go API
1
```

Then try:

```go
appName = "Another API"
```

and see the compiler error.

Next lesson → Lesson 6 — Functions in Go (parameters, return values, and how they differ from C# methods)