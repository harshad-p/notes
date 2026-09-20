# Structs

# Interfaces

- Create an interface type
- Add a method that accepts a type of this interface
- Create your struct type
- Attach a method to your struct type with the same name as that in the interface

And things should work. 

Excellent. This is one of the **most important lessons in Go**. Once you understand structs, you'll be able to model real data like `User`, `Product`, `Order`, etc.

# Lesson 12 — Structs (Go's Answer to Classes)

If you're thinking:

> "Where are the classes?"

The answer is:

> **Go doesn't have classes. It has structs.**

---

## 1. Creating a struct

Let's model a user.

```go
type User struct {
	Name string
	Age  int
}
```

Let's break it down:

```text
type   → Define a new type
User   → Type name
struct → A collection of fields
```

Think of it like this C# class:

```csharp
public class User
{
    public string Name { get; set; }
    public int Age { get; set; }
}
```

---

## 2. Creating an object (struct instance)

Now create a user:

```go
user := User{
	Name: "Harshad",
	Age:  36,
}
```

Print it:

```go
fmt.Println(user)
```

Output:

```text
{Harshad 36}
```

---

## 3. Accessing fields

Exactly what you'd expect:

```go
fmt.Println(user.Name)
fmt.Println(user.Age)
```

Output:

```text
Harshad
36
```

This feels very similar to C#:

```csharp
Console.WriteLine(user.Name);
```

---

## 4. Modifying fields

Structs are mutable.

```go
user.Age = 37

fmt.Println(user.Age)
```

Output:

```text
37
```

---

## Complete example

```go
package main

import "fmt"

type User struct {
	Name string
	Age  int
}

func main() {
	user := User{
		Name: "Harshad",
		Age:  36,
	}

	fmt.Println(user.Name)
	fmt.Println(user.Age)

	user.Age = 37

	fmt.Println(user.Age)
}
```

---

## A small Go style note

Notice this:

```go
type User struct
```

The type name starts with a **capital letter**.

For now, don't worry why.

Later you'll learn:

- `User` → visible outside the package
- `user` → only visible inside the package

This is how Go handles **public vs private** access. There are no `public` or `private` keywords.

---

## Exercise

Create this struct:

```go
type Book struct {
	Title  string
	Author string
	Pages  int
}
```

Then in `main()`:

1. Create a book.
2. Print its title.
3. Change the number of pages.
4. Print the updated page count.

---

## 🧠 Today's takeaway

You now know how Go models data.

| C# | Go |
|----|----|
| `class User` | `type User struct` |
| `new User()` | `User{}` |
| `user.Name` | `user.Name` |

### One important difference

At first glance, structs look like "classes without methods."

That's only **partly true**.

Go **does have methods**, but they are **not declared inside the struct** like C#.

Instead, they're attached **outside** the struct.

For example:

```go
func (u User) Greet() {
	fmt.Println("Hello", u.Name)
}
```

This is one of Go's most distinctive features, and understanding it will make the language "click."

---

### Next lesson → **Lesson 13 — Methods on Structs (why methods live outside the struct and what receivers are)**

This is where you'll see how Go achieves object-oriented programming without classes and inheritance.
# Lesson 13 — Methods and Receivers

We just saw that Go doesn't put methods **inside** a struct like C# does.

Instead, you define the method separately and connect it to the struct using a **receiver**.

Start with our `User`:

```go
type User struct {
	Name string
	Age  int
}
```

Now give `User` a method:

```go
func (u User) Greet() {
	fmt.Println("Hello", u.Name)
}
```

The interesting part is:

```go
(u User)
```

That's the **receiver**.

It means:

> `Greet()` is a method belonging to the `User` type.

So we can do:

```go
user := User{
	Name: "Harshad",
	Age:  36,
}

user.Greet()
```

Output:

```text
Hello Harshad
```

### Compare with C#

C#:

```csharp
public class User
{
    public string Name { get; set; }

    public void Greet()
    {
        Console.WriteLine($"Hello {Name}");
    }
}
```

Go separates the two:

```go
type User struct {
	Name string
}

func (u User) Greet() {
	fmt.Println("Hello", u.Name)
}
```

That's a major Go design choice.

---

## What is `u`?

Here:

```go
func (u User) Greet()
```

`u` is basically the current `User` instance.

So:

```go
fmt.Println(u.Name)
```

means:

> Get the `Name` from the User that called this method.

If:

```go
user.Name == "Harshad"
```

then:

```go
user.Greet()
```

prints:

```text
Hello Harshad
```

---

## One important thing: value receiver

This:

```go
func (u User) Greet()
```

receives a **copy** of the `User`.

We'll soon encounter:

```go
func (u *User) ChangeName(name string)
```

where `*User` means we're working with a pointer to the original struct.

**Don't worry about pointers yet.** We'll take that separately.

---

## Exercise

Using:

```go
type Book struct {
	Title  string
	Author string
}
```

create a method:

```go
func (b Book) Describe()
```

that prints something like:

```text
The Go Programming Language by Alan Donovan
```

Then create a `Book` and call:

```go
book.Describe()
```

---

### Next lesson → **Lesson 14 — Pointers in Go**

This is an important one. We'll slow down here because pointers are much easier to understand when you already know structs and methods.
# Lesson 14 — Pointers in Go

Don't worry — we're not going deep into memory management. We just need the practical idea.

You already know that when you pass a normal value to a function, Go passes a **copy**.

```go
func changeAge(user User) {
	user.Age = 40
}
```

If we do:

```go
user := User{Name: "Harshad", Age: 36}

changeAge(user)

fmt.Println(user.Age)
```

we still get:

```text
36
```

Why?

`changeAge()` received a **copy** of `user`.

---

## A pointer lets us work with the original

A pointer stores the **address of a value** rather than a copy of the value.

We use `*`:

```go
func changeAge(user *User) {
	user.Age = 40
}
```

Then:

```go
user := User{Name: "Harshad", Age: 36}

changeAge(&user)

fmt.Println(user.Age)
```

Output:

```text
40
```

Two symbols to understand:

```go
&user
```

means:

> Give me the address of `user`.

And:

```go
*User
```

means:

> A pointer to a `User`.

---

## Why did we need this?

This is particularly important with **methods**.

Remember:

```go
func (u User) Greet()
```

`u` is a copy.

But if a method needs to modify the original struct:

```go
func (u *User) ChangeAge(age int) {
	u.Age = age
}
```

Now:

```go
user := User{Name: "Harshad", Age: 36}

user.ChangeAge(40)

fmt.Println(user.Age)
```

prints:

```text
40
```

Notice something interesting:

We called:

```go
user.ChangeAge(40)
```

—not:

```go
(&user).ChangeAge(40)
```

Go automatically handles that for us in this situation.

---

## 🧠 The practical rule

For now, remember just this:

**Value receiver:**

```go
func (u User) SomeMethod()
```

→ method works with a copy.

**Pointer receiver:**

```go
func (u *User) SomeMethod()
```

→ method can modify the original.

You'll encounter pointer receivers constantly in real Go code.

---

### Exercise

Create:

```go
type User struct {
	Name string
	Age  int
}
```

Then create:

```go
func (u *User) Birthday() {
	u.Age++
}
```

Create a user who is 36, call:

```go
user.Birthday()
```
You're right. I was padding the lessons with repetition instead of using the fact that you've already understood the previous concept.

I'll change the approach:

- **No re-explaining concepts from the previous section.**
- **No "remember" sections that repeat what was just taught.**
- **No unnecessary C# comparisons unless there's a meaningful difference.**
- **Short lesson → example → exercise.**
- I'll assume you've understood something once you've moved on.
- If a concept needs reinforcement, I'll do it through the **next practical example**, not by repeating the explanation.

Let's continue.

# Lesson 15 — Interfaces

Go interfaces are quite different from C# interfaces.

In C#, a type explicitly declares that it implements an interface:

```csharp
class Dog : IAnimal
```

In Go, **there is no declaration saying "I implement this interface."**

Suppose we define:

```go
type Speaker interface {
	Speak()
}
```

Now define a `Dog`:

```go
type Dog struct {
	Name string
}

func (d Dog) Speak() {
	fmt.Println("Woof!")
}
```

`Dog` automatically satisfies `Speaker` because it has the required `Speak()` method.

We can therefore do:

```go
func makeSound(s Speaker) {
	s.Speak()
}
```

And:

```go
dog := Dog{Name: "Buddy"}

makeSound(dog)
```

Output:

```text
Woof!
```

There's no:

```go
type Dog implements Speaker
```

Go figures it out automatically.

---

## Why is this useful?

It allows code to depend on **behavior**, rather than a specific type.

For example:

```go
type Speaker interface {
	Speak()
}
```

Anything that has `Speak()` can be passed to:

```go
makeSound()
```

So this also works:

```go
type Cat struct{}

func (c Cat) Speak() {
	fmt.Println("Meow!")
}
```

Now both `Dog` and `Cat` satisfy `Speaker`.

---

## Exercise

Create:

```go
type Animal interface {
	Speak()
}
```

Then create:

```go
type Dog struct{}
type Cat struct{}
```

Give both types a `Speak()` method.

Finally create:

```go
func makeSound(a Animal)
```

and use it with both animals.

---

### Next lesson → **Lesson 16 — Error Handling in Go**

We'll look at the `error` type, `nil`, and the famous:

```go
if err != nil
```

pattern.
and print the age.

You should get:

```text
37
```

Don't worry about understanding every detail of pointers yet. **The next few lessons will reinforce this naturally.**

---

### Next lesson → **Lesson 15 — Interfaces in Go**

This is another *very* important difference from C#: Go interfaces are **implicitly implemented**.
