# Lesson 18 — Putting Structs, Methods & Interfaces Together

Rather than adding another isolated feature, let's combine the pieces you've learned into something closer to real application code.

Imagine we're building a tiny user service.

# 1. Define the model

```go
type User struct {
	Name string
	Age  int
}
```

# 2. Give it behavior

```go
func (u User) IsAdult() bool {
	return u.Age >= 18
}
```

Now:

```go
user := User{
	Name: "Harshad",
	Age: 36,
}

fmt.Println(user.IsAdult())
```

prints:

```text
true
```

# 3. Define an interface

Suppose some part of our application only cares about retrieving users:

```go
type UserRepository interface {
	GetUser(id int) (User, error)
}
```

We don't care **how** the user is retrieved.

It could come from:

- PostgreSQL
- an API
- an in-memory map
- a test fake

Anything providing:

```go
GetUser(id int) (User, error)
```

satisfies the interface automatically.

# 4. Put it together

## 1. The interface

We said:

```go
type UserRepository interface {
	GetUser(id int) (User, error)
}
```

This says only:

> Anything that wants to be a `UserRepository` must provide `GetUser`.

It does **not** contain the implementation.

---

## 2. The concrete implementation

Let's say we're temporarily storing users in memory:

```go
type InMemoryUserRepository struct {
	users map[int]User
}
```

Now give it the required method:

```go
func (r InMemoryUserRepository) GetUser(id int) (User, error) {
	user, exists := r.users[id]

	if !exists {
		return User{}, errors.New("user not found")
	}

	return user, nil
}
```

Because `InMemoryUserRepository` has the required method, **it automatically satisfies `UserRepository`**.

No `implements` keyword.

---

## 3. The service

Now our service:

```go
type UserService struct {
	repository UserRepository
}

func (s UserService) GetUser(id int) (User, error) {
	return s.repository.GetUser(id)
}
```

Notice that `UserService` doesn't know that we're using `InMemoryUserRepository`.

It only knows:

```go
repository UserRepository
```

That's the whole point of the interface.

---

## 4. Connecting them

We create the concrete repository:

```go
repository := InMemoryUserRepository{
	users: map[int]User{
		1: {
			Name: "Harshad",
			Age:  36,
		},
	},
}
```

Then inject it into the service:

```go
service := UserService{
	repository: repository,
}
```

Now:

```go
user, err := service.GetUser(1)

if err != nil {
	fmt.Println(err)
	return
}

fmt.Println(user.Name)
```

prints:

```text
Harshad
```

The actual call chain is:

```text
service.GetUser(1)
        ↓
s.repository.GetUser(1)
        ↓
InMemoryUserRepository.GetUser(1)
```

The important part is that **the interface variable contains the concrete implementation**.

---

# Where would this reside in a real application?

You might eventually structure it something like:

```text
my-api/
├── main.go
├── user/
│   ├── user.go
│   ├── service.go
│   ├── repository.go
│   └── postgres_repository.go
└── go.mod
```

For example:

### `repository.go`

```go
type UserRepository interface {
	GetUser(id int) (User, error)
}
```

### `postgres_repository.go`

```go
type PostgresUserRepository struct {
	db *sql.DB
}

func (r PostgresUserRepository) GetUser(id int) (User, error) {
	// Query PostgreSQL here
}
```

### `service.go`

```go
type UserService struct {
	repository UserRepository
}

func (s UserService) GetUser(id int) (User, error) {
	return s.repository.GetUser(id)
}
```

### `main.go`

This is where you wire everything together:

```go
repository := PostgresUserRepository{
	db: db,
}

service := UserService{
	repository: repository,
}
```

This last part is commonly called **dependency injection / application wiring**.

And there's a nice distinction from the .NET world:

> In Go, you don't need a DI container just to accomplish this.

You can construct the dependencies yourself.

---

### One subtle but important point

You might be wondering:

> "Why have `UserRepository` at all? Why not just put `PostgresUserRepository` directly into `UserService`?"

You absolutely **could**.

The interface becomes useful when the service should depend on the **behavior** rather than the specific implementation.

For example, in a test:

```go
type FakeUserRepository struct {
	// test data
}
```

It can implement the same interface, and you can give that to `UserService` instead of connecting to PostgreSQL.

That's where this design starts becoming particularly useful.

This is an important Go pattern: **structs hold dependencies, interfaces describe the behavior they need**.

You will see this heavily in Go backend applications.

---

## One useful observation

Notice that we didn't need inheritance anywhere.

We have:

```text
User
UserRepository
UserService
```

and they are composed together through fields and interfaces.

This is typical Go design: **composition over inheritance**.

---

## Exercise

Don't build the whole repository yet.

Just create:

```go
type Product struct {
	Name  string
	Price float64
}
```

Add:

```go
func (p Product) IsExpensive() bool
```

Return `true` when the price is greater than `100`.

Then create two products and test the method.

---

### Next lesson → **Lesson 19 — `defer`: Go's Simple Cleanup Mechanism**

We'll learn `defer`, which you'll encounter constantly when working with files, database connections, HTTP responses, and other resources.