# Lesson 21 — Strings in Go

We've already used strings:

```go
name := "Harshad"
```

Now let's look at a few things that are useful in real Go code.

## 1. Concatenation

You can join strings with `+`:

```go
firstName := "Harshad"
lastName := "Paradkar"

fullName := firstName + " " + lastName

fmt.Println(fullName)
```

Output:

```text
Harshad Paradkar
```

---

## 2. `len()` gives you bytes, not necessarily characters

This is an important Go detail.

```go
name := "Harshad"

fmt.Println(len(name))
```

Output:

```text
7
```

But:

```go
text := "café"

fmt.Println(len(text))
```

may give you:

```text
5
```

even though `"café"` contains 4 visible characters.

That's because Go strings are sequences of **bytes**, and UTF-8 characters can occupy multiple bytes.

You don't need to deal with Unicode internals yet. Just remember:

> `len(string)` gives the number of bytes, not necessarily the number of characters.

---

## 3. Indexing a string

You can access an individual byte:

```go
name := "Harshad"

fmt.Println(name[0])
```

You'll get:

```text
72
```

Not `"H"`.

Why? `name[0]` is a `byte` (`uint8`), and `H` has the numeric UTF-8 value `72`.

If you want to print it as a character:

```go
fmt.Printf("%c\n", name[0])
```

Output:

```text
H
```

This distinction becomes important when we get to `[]byte` and `rune`.

---

## 4. Useful string functions

Go's `strings` package provides common operations:

```go
import "strings"
```

For example:

```go
text := "Hello Go"

fmt.Println(strings.Contains(text, "Go"))
fmt.Println(strings.ToUpper(text))
fmt.Println(strings.ToLower(text))
```

Output:

```text
true
HELLO GO
hello go
```

You'll use functions like these frequently when processing request data.

---

## Small exercise

Write a program that:

1. Creates a string containing `"Go is awesome"`
2. Prints its length
3. Checks whether it contains `"awesome"`
4. Converts it to uppercase

Don't worry about Unicode for the exercise.

### Next lesson → **Lesson 22 — Pointers, Part 2: `*` and `&` in practice**

We'll revisit pointers, but this time from the perspective of actual Go code and pointer semantics rather than repeating Lesson 14.