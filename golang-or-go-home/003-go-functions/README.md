# Lesson 6 — Functions

This is the first lesson where Go starts to feel like a real programming language rather than just syntax.

A function is simply a reusable block of code.

You've already used one:

```go
func main() {

}
```

`main()` is just a special function that Go runs first.

---

## Creating your own function

Let's create one that greets a user.

```go
package main

import "fmt"

func greet() {
	fmt.Println("Hello!")
}

func main() {
	greet()
}
```

Output:

```text
Hello!
```

Notice something?

We didn't write:

```csharp
public static void Greet()
```

We simply wrote:

```go
func greet()
```

Go intentionally keeps function declarations simple.

---

## Functions with parameters

Functions become useful when they can accept data.

```go
func greet(name string) {
	fmt.Println("Hello", name)
}
```

Call it like this:

```go
func main() {
	greet("Harshad")
}
```

Output:

```text
Hello Harshad
```

### Compare with C#

C#:

```csharp
void Greet(string name)
```

Go:

```go
func greet(name string)
```

**Notice the difference.**

In C#, we write:

```text
type name
```

In Go, we write:

```text
name type
```

This is one of the biggest syntax differences you'll encounter.

---

## Returning a value

Functions can also return values.

```go
func square(number int) int {
	return number * number
}
```

Let's use it:

```go
func main() {
	result := square(5)

	fmt.Println(result)
}
```

Output:

```text
25
```

### Reading this declaration

```go
func square(number int) int
```

Break it into pieces:

```
func
│
├── function name → square
│
├── parameter → number int
│
└── returns → int
```

The last `int` is the return type.

---

## Exercise

Write this program:

```go
package main

import "fmt"

func multiply(a int, b int) int {
	return a * b
}

func main() {
	result := multiply(6, 7)

	fmt.Println(result)
}
```

It should print:

```text
42
```

---

## 💡 Small Go improvement

Go lets you shorten parameters that have the same type.

Instead of:

```go
func multiply(a int, b int) int
```

you can write:

```go
func multiply(a, b int) int
```

Exactly the same meaning.

Most Go developers prefer this style.

---

## Today's takeaway

You learned that a Go function has four parts:

```go
func multiply(a, b int) int {
	return a * b
}
```

- `func` → declares a function
- `multiply` → function name
- `(a, b int)` → parameters
- `int` → return type

Once you're comfortable with this, you're already reading a large percentage of everyday Go functions.

---

### Mini challenge (no Googling 😊)

Can you write a function named `add` that:

- accepts **two integers**
- returns their sum

Then call it from `main()` with `10` and `20`, and print the result.

Don't worry if you make mistakes—that's exactly how we'll learn.
Great. Let's continue.

# Lesson 7 — Multiple Return Values (A Very Go Feature)

This is one of the places where Go is quite different from C#.

In C#, if a method needs to return multiple values, you might use:

- a class/object
- a tuple
- `out` parameters

Example in C#:

```csharp
(bool success, string message) Login()
```

Go has built-in support for this.

---

## A function can return multiple values

Example:

```go
package main

import "fmt"

func getUser() (string, int) {
	return "Harshad", 36
}

func main() {
	name, age := getUser()

	fmt.Println(name)
	fmt.Println(age)
}
```

Output:

```text
Harshad
36
```

---

## How does this work?

Look at the function:

```go
func getUser() (string, int)
```

The return section is:

```go
(string, int)
```

Meaning:

> This function returns two values:
> 1. a string
> 2. an integer

Then:

```go
name, age := getUser()
```

receives both values.

---

# Why is this useful?

The biggest use case in Go is **error handling**.

You will see this everywhere:

```go
user, err := findUser(id)
```

A function often returns:

1. The actual result
2. An error

Example:

```go
func divide(a, b int) (int, error) {

	if b == 0 {
		return 0, fmt.Errorf("cannot divide by zero")
	}

	return a / b, nil
}
```

Usage:

```go
result, err := divide(10, 2)

if err != nil {
	fmt.Println(err)
	return
}

fmt.Println(result)
```

Output:

```text
5
```

Don't worry about `error` and `nil` yet. We will have a full lesson on error handling.

For now remember:

> In Go, functions commonly return `(value, error)`.

This pattern replaces many exception-based flows you may be used to in C#.

---

# Important difference from C#

C#:

```csharp
try
{
    var user = GetUser();
}
catch(Exception ex)
{
}
```

Go:

```go
user, err := GetUser()

if err != nil {
    // handle error
}
```

Go prefers:
- explicit checking
- fewer hidden control flows

---

# Exercise

Write a function:

```go
func calculate(a, b int) (int, int)
```

It should return:

1. Sum
2. Difference

Example:

```go
sum, diff := calculate(10, 3)

fmt.Println(sum)
fmt.Println(diff)
```

Expected output:

```text
13
7
```

Next lesson → Lesson 8 — If/Else and Loops in Go (what is different from C#)