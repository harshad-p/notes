Excellent. This is another collection you'll use constantly in backend development.

# Lesson 11 — Maps

A **map** stores data as **key → value** pairs.

If you know C#, think:

```csharp
Dictionary<TKey, TValue>
```

---

## 1. Creating a map

```go
ages := map[string]int{
	"Harshad": 36,
	"Alice":   30,
	"Bob":     25,
}
```

Let's read it:

```text
map[
    key type   -> string
    value type -> int
]
```

So:

- Key = `string`
- Value = `int`

---

## 2. Reading a value

Use the key:

```go
fmt.Println(ages["Harshad"])
```

Output:

```text
36
```

Just like a dictionary lookup in C#:

```csharp
Console.WriteLine(ages["Harshad"]);
```

---

## 3. Adding or updating

To add:

```go
ages["Charlie"] = 28
```

To update:

```go
ages["Harshad"] = 37
```

There's no separate `Add()` method like in C#. Assigning a value creates the key if it doesn't exist or updates it if it does.

---

## 4. Deleting

Use the built-in `delete()` function.

```go
delete(ages, "Bob")
```

Now `"Bob"` is gone.

---

## ⚠️ A Go gotcha

What happens here?

```go
fmt.Println(ages["David"])
```

There is no `"David"` key.

Will it crash?

**No.**

Output:

```text
0
```

Why?

Because `0` is the **zero value** for `int`.

Similarly:

| Type | Zero value |
|------|------------|
| `int` | `0` |
| `string` | `""` (empty string) |
| `bool` | `false` |

This behavior surprises many people coming from C#.

---

## 5. Checking if a key exists

Suppose `"David"` could actually have age `0` (unlikely, but imagine another example).

How do you know if the key exists?

Go returns **two values**:

```go
age, exists := ages["David"]
```

Then:

```go
if exists {
	fmt.Println(age)
} else {
	fmt.Println("User not found")
}
```

This is another example of Go's multiple return values in action.

---

## Exercise

Create this map:

```go
favoriteLanguages := map[string]string{
	"Harshad": "Go",
	"Alice":   "C#",
}
```

Then:

1. Print Harshad's favorite language.
2. Add `"Bob": "Python"`.
3. Change Alice's language to `"Rust"`.
4. Delete Bob.
5. Check whether `"Charlie"` exists.

---

## 🧠 Today's takeaway

| Operation | Go |
|-----------|----|
| Create | `map[string]int{}` |
| Read | `ages["Harshad"]` |
| Add/Update | `ages["Bob"] = 25` |
| Delete | `delete(ages, "Bob")` |
| Check existence | `value, exists := ages["Bob"]` |

Maps are used everywhere in backend applications:

- HTTP headers
- Configuration values
- Caching
- Looking up users by ID
- JSON objects

So getting comfortable with them will pay off.

---

### Next lesson → **Lesson 12 — Structs (Go's answer to classes)**

This is one of the biggest conceptual shifts from C#. We'll answer questions like:

- If Go doesn't have traditional classes, how do you model a `User` or `Product`?
- Where do methods belong?
- How does Go achieve object-oriented programming without inheritance?

This is where Go starts to feel fundamentally different from C#.