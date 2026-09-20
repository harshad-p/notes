Great. Now we understand the basic Go project setup. Let's take the next small step.

# Lesson 2: Structure of a Go Program

Your `main.go`:

```go
package main

import "fmt"

func main() {
	fmt.Println("Hello Go")
}
```

Let's understand each part.

---

## 1. Package

```go
package main
```

Every Go file starts with a package declaration.

A package is a way to organize code.

For now, remember:

- `package main` → this program can be executed
- other packages → reusable code libraries

Think:

C#:

```csharp
namespace MyApp
```

Go:

```go
package myapp
```

(not exactly the same, but similar purpose)

---

## 2. Import

```go
import "fmt"
```

Imports another package.

`fmt` is Go's built-in formatting package.

It gives us functions like:

```go
fmt.Println()
fmt.Printf()
```

Similar to:

C#:

```csharp
using System;
```

---

## 3. Function

```go
func main()
```

`func` means we are creating a function.

`main` is special.

When you run:

```bash
go run .
```

Go looks for:

```go
func main()
```

and starts execution there.

Similar to C#:

```csharp
static void Main()
```

---

## 4. Curly braces

Go uses `{ }` like C# and Java.

Example:

```go
func main() {
    // code here
}
```

One important difference:

Go is strict about formatting.

This works:

```go
func main() {
}
```

This does not:

```go
func main()
{
}
```

Go automatically inserts semicolons, so formatting matters.

---

# Small practical change

Modify your program:

```go
package main

import "fmt"

func main() {
	fmt.Println("My name is Harshad")
	fmt.Println("I am learning Go")
}
```

Run:

```bash
go run .
```

You should see:

```
My name is Harshad
I am learning Go
```

---

## Exercise

Change the program to print:

```
I am a backend developer
I am learning Go for building APIs
```

No new concepts yet. Just get comfortable running Go programs.

# 001 — Hello Go

The first Go program. A breakdown of every line.

---

## The Code

```go
package main

import "fmt"

func main() {
    fmt.Println("golang or go home!")
}
```

---

## Line by Line

### `package main`
Every Go file must declare which **package** it belongs to. The `main` package is special — it's the entry point of any executable program. Without it, Go won't know where to start.

---

### `import "fmt"`
This imports the **`fmt` package** from Go's standard library. `fmt` stands for *format* — it handles printing, string formatting, and input. It's one of the most commonly used packages in Go.

---

### `func main() {`
This declares the **`main` function**. When you run a Go program, this is the first (and in our case, only) function that executes. `func` is the keyword for declaring a function.

---

### `fmt.Println("golang or go home!")`
This calls the **`Println`** function from the `fmt` package. It prints the text to the terminal followed by a new line. The dot (`.`) is how you access something that belongs to a package — `fmt.Println` means "the `Println` function inside `fmt`".

---

## The `go.mod` File

When you ran `go mod init golang-or-go-home`, Go created a `go.mod` file. It looks like this:

```
module golang-or-go-home

go 1.24
```

- **`module`** — declares the name of your module (your project)
- **`go`** — the minimum Go version required

As you add external packages later, this file will also track your dependencies — similar to `package.json` in Node or `requirements.txt` in Python.

---

## Running the Program

```bash
go run main.go
```

- **`go run`** — compiles and runs the file in one step, without producing a binary file
- **`main.go`** — the file to run

To produce an actual executable binary:

```bash
go build main.go
```

That creates a file called `main` (or `main.exe` on Windows) that you can run directly.

Next lesson:
**Variables in Go (`var` vs `:=`) — one of the first places Go feels different from C#**.