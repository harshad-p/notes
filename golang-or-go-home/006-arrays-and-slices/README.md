Great. Now we move into one of the most important Go concepts.

# Lesson 9 — Slices (Go's Most Used Collection Type)

If you come from C#, think of a **slice** as being closest to:

```csharp
List<T>
```

A slice is a **dynamic collection** that can grow or shrink.

---

## 1. Creating a slice

Example:

```go
numbers := []int{10, 20, 30}
```

This creates a slice of integers.

Think:

```csharp
var numbers = new List<int> { 10, 20, 30 };
```

---

## 2. Accessing elements

Like arrays and lists:

```go
fmt.Println(numbers[0])
```

Output:

```text
10
```

Indexing starts at **0**.

So:

```text
Index:    0    1    2
Value:   10   20   30
```

---

## 3. Adding elements

In C#:

```csharp
numbers.Add(40);
```

In Go:

```go
numbers = append(numbers, 40)
```

Example:

```go
numbers := []int{10, 20, 30}

numbers = append(numbers, 40)

fmt.Println(numbers)
```

Output:

```text
[10 20 30 40]
```

---

## Important: `append` returns a new slice

This is a small but important Go concept.

You don't do:

```go
append(numbers, 40)
```

You do:

```go
numbers = append(numbers, 40)
```

Because `append()` may create a new underlying storage area.

---

## 4. Looping through a slice

You can use a normal `for`:

```go
numbers := []int{10, 20, 30}

for i := 0; i < len(numbers); i++ {
	fmt.Println(numbers[i])
}
```

`len()` gives the size.

Output:

```text
10
20
30
```

---

## 5. The Go way: `range`

Most Go developers use:

```go
for index, value := range numbers {
	fmt.Println(index, value)
}
```

Output:

```text
0 10
1 20
2 30
```

This is similar to:

```csharp
foreach(var number in numbers)
{
}
```

but Go also gives you the index.

---

## Ignoring a value with `_`

What if you don't need the index?

This gives an error:

```go
for index, value := range numbers {
	fmt.Println(value)
}
```

because `index` is unused.

Go does not allow unused variables.

So use `_`:

```go
for _, value := range numbers {
	fmt.Println(value)
}
```

Meaning:

> "I know this value exists, but I don't care about it."

You will see `_` everywhere in Go code.

---

# Exercise

Create a slice:

```go
languages := []string{"C#", "Go", "JavaScript"}
```

Then:

1. Add `"Python"` using `append`
2. Loop through it using `range`
3. Print each language

Expected output:

```text
C#
Go
JavaScript
Python
```

---

## Today's takeaway

| C# | Go |
|-|-|
| `List<T>` | slice |
| `.Add()` | `append()` |
| `.Count` | `len()` |
| `foreach` | `for range` |

Slices are **everywhere** in Go:
- API responses
- database results
- JSON data
- lists of users/products/orders

Mastering slices is essential before we build a backend.

---

### Next lesson → **Lesson 10 — Arrays vs Slices (Why Go has both, and when to use each)**
Great. Now let's clear up a common confusion: **arrays vs slices**.

If you come from C#, this lesson is important because you might think:

> "Why does Go need both? Doesn't `List<T>` already solve this?"

Go has a different design.

# Lesson 10 — Arrays vs Slices

## 1. Arrays

An array has a **fixed size**.

Example:

```go
numbers := [3]int{10, 20, 30}
```

The `[3]` means:

> This array can hold exactly 3 integers.

You cannot add another element.

This would fail:

```go
numbers = append(numbers, 40)
```

because the size is fixed.

---

## 2. Slices

A slice is dynamic.

Example:

```go
numbers := []int{10, 20, 30}
```

Notice the difference:

Array:

```go
[3]int
```

Slice:

```go
[]int
```

The slice has no size specified.

You can do:

```go
numbers = append(numbers, 40)
```

Result:

```text
[10 20 30 40]
```

---

## C# comparison

Think of it like:

| Go | C# equivalent |
|-|-|
| Array | `int[]` |
| Slice | `List<int>` |

Example:

C# array:

```csharp
int[] numbers = {10,20,30};
```

Fixed size.

C# list:

```csharp
var numbers = new List<int>{10,20,30};
```

Dynamic size.

---

# Why does Go have arrays?

Good question.

Most application code uses slices.

Arrays exist because they are useful when:

- size is known beforehand
- performance matters
- working with low-level data

Example:

```go
var days [7]string
```

A week always has 7 days.

---

# Slice internals (important idea)

A slice is actually a small structure containing:

```
Slice
 |
 +-- pointer → actual data
 |
 +-- length
 |
 +-- capacity
```

Example:

```go
numbers := []int{10,20,30}
```

Conceptually:

```
numbers
   |
   v
[10][20][30]
 length = 3
```

When you append:

```go
numbers = append(numbers, 40)
```

Go may:

1. use existing space
2. allocate a bigger array
3. copy values

You don't manage this manually.

---

# Creating an empty slice

Very common in backend code:

```go
users := []string{}
```

Then:

```go
users = append(users, "Harshad")
users = append(users, "Alex")
```

Result:

```text
[Harshad Alex]
```

You will see this pattern frequently when building API responses.

---

# Exercise

Create an empty slice:

```go
tasks := []string{}
```

Add:

```
"Learn Go"
"Build API"
"Deploy App"
```

Then print the slice.

Expected:

```text
[Learn Go Build API Deploy App]
```

---

## Today's takeaway

Remember only this:

- **Array** → fixed size → `[3]int`
- **Slice** → dynamic size → `[]int`
- In normal Go backend development, you will use **slices much more often**

---

### Next lesson → **Lesson 11 — Maps in Go (similar to C# Dictionary<TKey,TValue>)**