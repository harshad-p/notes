# Lesson 17 — Packages and Modules

We've been putting everything into `main.go`. That's fine for learning, but a real backend will quickly have dozens of files.

Go gives us **packages** to organize that code.

### Multiple files in the same package

Suppose we have:

```text
my-first-go-app/
├── go.mod
├── main.go
└── math.go
```

Both files can belong to:

```go
package main
```

`math.go`:

```go
package main

func add(a, b int) int {
	return a + b
}
```

`main.go`:

```go
package main

import "fmt"

func main() {
	result := add(10, 20)
	fmt.Println(result)
}
```

You don't need to import `math.go`. Both files belong to the same package, so their functions are available to each other.

Run:

```bash
go run .
```

and you'll get:

```text
30
```

### Packages are more than files

Now imagine we want a separate `math` package:

```text
my-first-go-app/
├── go.mod
├── main.go
└── math/
    └── math.go
```

`math/math.go`:

```go
package math

func Add(a, b int) int {
	return a + b
}
```

`main.go`:

```go
package main

import (
	"fmt"
	"my-first-go-app/math"
)

func main() {
	result := math.Add(10, 20)
	fmt.Println(result)
}
```

Notice `Add` is capitalized.

That's significant in Go:

```go
func Add()
```

is **exported** and can be used by another package.

```go
func add()
```

is **unexported** and cannot.

This is how Go controls visibility instead of using `public`/`private` keywords.

### What does `go.mod` do?

You created it earlier with:

```bash
go mod init my-first-go-app
```

The module name becomes the root of your package import paths:

```go
"my-first-go-app/math"
```

For a real project, you'd normally use something like:

```bash
go mod init github.com/yourname/my-api
```

Then packages within it can be imported using:

```go
"github.com/yourname/my-api/math"
```

We don't need to change yours yet.

---

### Exercise

Create a `greeting` package:

```text
my-first-go-app/
├── go.mod
├── main.go
└── greeting/
    └── greeting.go
```

Have it expose:

```go
func Hello(name string) string
```

which returns:

```text
Hello Harshad
```

Import the package from `main.go` and print the result.

This is your first step toward structuring an actual Go application.

### Next lesson → **Lesson 18 — Pointers, Structs & Interfaces: Putting the pieces together**

We'll use the things you've learned so far in a small practical program rather than introducing another isolated language feature.