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
