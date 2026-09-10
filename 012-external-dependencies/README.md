# Lesson 20 — Go Modules & External Dependencies

We've already seen that `go.mod` identifies your Go module. Now let's actually **add a third-party dependency** and see what Go does with it.

We'll use a small, real package rather than writing everything ourselves.

## Step 1 — Add a dependency

From your project directory, run:

```bash
go get github.com/fatih/color
```

Go will download the package and update your `go.mod`.

You should now see something along the lines of:

```go
require github.com/fatih/color v1.x.x
```

The exact version will depend on what Go resolves today.

You'll also get a new file:

```text
go.sum
```

### What are these two files?

**`go.mod`**

Says:

> These are the dependencies my project requires.

**`go.sum`**

Contains checksums that Go uses to verify the downloaded dependency's contents.

For now, think:

```text
go.mod → what dependencies I use
go.sum → verification information for those dependencies
```

---

## Step 2 — Use the package

Change `main.go`:

```go
package main

import (
	"fmt"

	"github.com/fatih/color"
)

func main() {
	fmt.Println("Normal text")
	color.Red("This is red")
}
```

Run:

```bash
go run .
```

You should see the two messages, with the second one displayed using the package's color functionality.

Notice that we didn't manually download anything or copy source code into our project.

The sequence was simply:

```bash
go get github.com/fatih/color
```

then:

```go
import "github.com/fatih/color"
```

and Go takes care of the dependency.

---

## One useful detail

If you delete the dependency from your code and then run:

```bash
go mod tidy
```

Go examines your project and removes dependencies that are no longer needed.

It can also add dependencies that your code imports but that aren't correctly recorded yet.

So a command you'll use frequently is:

```bash
go mod tidy
```

Think of it as:

> "Make my `go.mod` and `go.sum` accurately reflect what this project actually uses."

---

### Exercise

Do the following:

1. Run `go get github.com/fatih/color`
2. Confirm that `go.mod` changed.
3. Confirm that `go.sum` appeared.
4. Import the package and use one of its functions.
5. Run `go mod tidy`.
6. Look at `go.mod` and `go.sum` again.

**Don't worry about memorizing the contents of `go.sum`.** The important thing at this stage is understanding the role of the two files and the dependency workflow.

### Next lesson → **Lesson 21 — Strings in Go**

We'll move back into the language itself: string operations, indexing, and a couple of Go-specific details that become relevant when handling HTTP/JSON data.