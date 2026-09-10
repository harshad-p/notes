# JavaScript — Lesson 2: Values, Variables, and References

This is a deceptively important topic. A lot of confusing JavaScript behavior comes from not distinguishing **a value** from **a reference to an object**.

## 1. JavaScript has values

Some basic values:

```javascript
42
"hello"
true
false
null
undefined
```

And objects:

```javascript
{ name: "Alice" }
[1, 2, 3]
```

JavaScript broadly divides values into:

### Primitive values

- `string`
- `number`
- `boolean`
- `null`
- `undefined`
- `bigint`
- `symbol`

### Objects

Everything else that's an object, including:

- ordinary objects
- arrays
- functions
- dates
- maps
- sets

For now, the key distinction is **primitive vs object**.

---

# 2. Variables hold values

```javascript
let age = 37;
```

You can think of `age` as referring to the value `37`.

```javascript
let name = "Harshad";
```

`name` refers to the string `"Harshad"`.

With `const`:

```javascript
const age = 37;
```

you cannot make `age` refer to another value:

```javascript
age = 38; // Error
```

But there's an important subtlety.

---

# 3. `const` doesn't make objects immutable

Consider:

```javascript
const user = {
    name: "Alice",
    age: 30
};
```

This is allowed:

```javascript
user.age = 31;
```

But this isn't:

```javascript
user = {};
```

Why?

Because `const` prevents **reassignment of the variable**.

It doesn't prevent mutation of the object that the variable refers to.

Conceptually:

```text
user
 │
 ↓
┌──────────────────┐
│ name: "Alice"    │
│ age: 30          │
└──────────────────┘
```

You can't make `user` point somewhere else:

```text
user ─────X────→ another object
```

But you can modify the existing object.

This distinction becomes **extremely important in React**, because React relies heavily on creating new objects/arrays rather than mutating existing ones.

---

# 4. Primitive assignment

Look at this:

```javascript
let a = 10;
let b = a;

b = 20;
```

What is `a`?

```javascript
console.log(a); // 10
```

`b` received the value `10`.

Changing `b` doesn't affect `a`.

Conceptually:

```text
a → 10
b → 10
```

After:

```javascript
b = 20;
```

we have:

```text
a → 10
b → 20
```

They're independent.

---

# 5. Objects behave differently

Now:

```javascript
const a = {
    value: 10
};

const b = a;

b.value = 20;
```

What is:

```javascript
a.value
```

?

It's:

```text
20
```

Why?

Because `a` and `b` refer to the **same object**.

Conceptually:

```text
a ──────┐
        ↓
     ┌───────────┐
     │ value: 10 │
     └───────────┘
        ↑
b ──────┘
```

After:

```javascript
b.value = 20;
```

there is still only one object:

```text
a ──────┐
        ↓
     ┌───────────┐
     │ value: 20 │
     └───────────┘
        ↑
b ──────┘
```

This is why we say objects have **reference semantics**.

---

# 6. This also applies to arrays

```javascript
const a = [1, 2, 3];
const b = a;

b.push(4);
```

Now:

```javascript
console.log(a);
```

gives:

```text
[1, 2, 3, 4]
```

Again, there's one array and two variables referring to it.

---

# 7. A common mistake

People sometimes think this creates a copy:

```javascript
const b = a;
```

It doesn't.

For an object/array, it copies the **reference**, not the object.

To make a shallow copy:

```javascript
const b = { ...a };
```

or:

```javascript
const b = [...a];
```

We'll go much deeper into this later when we discuss **spread, shallow copies, deep copies, and immutability**.

---

# 8. Why you should care for React

Suppose React has some state:

```javascript
const [user, setUser] = useState({
    name: "Alice",
    age: 30
});
```

You generally don't want to do:

```javascript
user.age = 31;
```

Instead:

```javascript
setUser({
    ...user,
    age: 31
});
```

Why?

Because the second version creates a **new object**.

That's connected directly to what we learned here about references.

We'll eventually get to the deeper question:

> **How does React know that something changed?**

And the answer involves reference identity.

So this seemingly basic JavaScript concept is actually foundational to React.

---

## One thing to remember

There are two ideas you should keep separate:

**Primitive:**

```javascript
let b = a;
```

copies the value.

**Object/array:**

```javascript
let b = a;
```

copies the reference to the same object.

And:

> `const` prevents reassignment, not mutation.

---

### Quick exercise

What does this print?

```javascript
const user1 = {
    name: "Alice",
    address: {
        city: "Berlin"
    }
};

const user2 = { ...user1 };

user2.name = "Bob";
user2.address.city = "Munich";

console.log(user1.name);
console.log(user1.address.city);
```

Think carefully about **both** properties. Don't run it.

Say your answer, then **next**.