# Lesson 1: Understanding a Go Project

Unlike C# where you often start with a `.sln` and `.csproj`, Go uses a simpler concept:

**A folder = a project/module**

Example:

```
my-api/
    main.go
    go.mod
```

- `main.go` → your Go code
- `go.mod` → describes your project and dependencies

---

## Step 1: Create a project folder

Go somewhere you keep your code:

```bash
mkdir my-first-go-app
cd my-first-go-app
```

---

## Step 2: Initialize a Go module

Run:

```bash
go mod init my-first-go-app
```

This creates:

```
my-first-go-app/
    go.mod
```

`go.mod` is similar to:

- C# → `.csproj`
- Node.js → `package.json`

It tracks:
- project name
- Go version
- external libraries later

---

## Step 3: Create your first Go file

Create:

```
main.go
```

Put this inside:

```go
package main

import "fmt"

func main() {
	fmt.Println("Hello Go")
}
```

---

## Step 4: Run it

From the project folder:

```bash
go run .
```

Output:

```
Hello Go
```

Notice:

We did **not** do:

```
go run main.go
```

We used:

```
go run .
```

Meaning:

> Run this entire Go project.

This is the habit we will use for backend projects.

---

## Your exercise

Do only this:

1. Create the folder
2. Run `go mod init`
3. Create `main.go`
4. Run:

```bash
go run .
```

Tell me when you see:

```
Hello Go
```

Then Lesson 2 will be:
**How Go code is structured (`package`, `import`, `func main`)**.

We will go one concept at a time.