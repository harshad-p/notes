# JavaScript — Lesson 1: How JavaScript Actually Runs

Before we get into syntax, I want you to understand the basic execution model. This will pay off later when we get to **closures, async/await, React rendering, and Angular**.

## 1. JavaScript needs a runtime

JavaScript itself is a language. Something has to execute it.

Two common environments are:

- **Browser** → Chrome, Firefox, Safari, etc.
- **Node.js** → JavaScript running outside the browser

For example:

```javascript
const name = "Harshad";
console.log(name);
```

The JavaScript engine reads this code and executes it.

Chrome uses **V8**. Node.js also uses **V8**.

You don't need to know the internals of V8 yet. The important distinction is:

> **JavaScript language ≠ JavaScript runtime.**

The runtime also provides things that aren't part of the core language itself.

For example, browsers provide:

```javascript
fetch()
setTimeout()
document
localStorage
```

Node provides things such as:

```javascript
fs
process
Buffer
```

That's why JavaScript code can behave differently depending on where it runs.

---

# 2. JavaScript executes code in an execution context

Consider:

```javascript
const x = 10;

function add(a, b) {
    return a + b;
}

const result = add(x, 20);
```

When JavaScript starts executing this program, it creates an **execution context**.

You can think of it roughly as:

> "The environment in which this particular piece of JavaScript is being executed."

The initial one is the **global execution context**.

It contains things such as:

- variables
- functions
- information about the current environment
- the value of `this`

Then when we call:

```javascript
add(x, 20);
```

JavaScript creates another execution context for `add`.

So conceptually:

```text
Global Execution Context
        |
        | calls add()
        ↓
add() Execution Context
```

When `add()` finishes, its execution context goes away.

This becomes very important when we later discuss **closures**.

---

# 3. The call stack

JavaScript uses a **call stack** to keep track of which functions are currently executing.

Consider:

```javascript
function first() {
    second();
}

function second() {
    third();
}

function third() {
    console.log("Hello");
}

first();
```

Execution roughly looks like:

```text
first()
   ↓
second()
   ↓
third()
```

The stack becomes:

```text
┌─────────┐
│ third() │
├─────────┤
│ second()│
├─────────┤
│ first() │
├─────────┤
│ global  │
└─────────┘
```

`third()` finishes first, so it is removed.

Then:

```text
second()
first()
global
```

Then `second()` finishes:

```text
first()
global
```

And finally `first()`:

```text
global
```

This is why it's called a **stack**.

It's essentially **last in, first out (LIFO)**.

---

# 4. What happens with recursion?

This also explains recursion.

```javascript
function count(n) {
    if (n === 0) {
        return;
    }

    count(n - 1);
}

count(3);
```

The stack grows:

```text
count(3)
count(2)
count(1)
count(0)
```

Then they return one by one.

If you recurse too deeply, you can eventually get:

```text
RangeError: Maximum call stack size exceeded
```

That's a literal consequence of the call stack becoming too large.

---

# 5. One important thing: JavaScript is synchronous by default

Consider:

```javascript
console.log("A");

console.log("B");

console.log("C");
```

You get:

```text
A
B
C
```

JavaScript executes one piece of code at a time.

Now:

```javascript
console.log("A");

function test() {
    console.log("B");
}

test();

console.log("C");
```

Still:

```text
A
B
C
```

The call stack has to finish the current work before moving on.

This is one reason JavaScript is often described as **single-threaded**.

There is more nuance to that statement once we get into browser runtimes, Web APIs and workers, but for now:

> **One JavaScript execution thread processes synchronous code sequentially.**

---

# 6. Why this matters for React

Later, you'll write things like:

```javascript
function User() {
    console.log("Rendering User");

    return <div>Hello</div>;
}
```

And you'll hear:

> "React renders the component."

What does that actually mean?

At a fundamental level, React is causing JavaScript functions to execute.

Understanding:

- functions
- execution contexts
- the call stack
- synchronous execution

will make React's behavior much less mysterious.

Later we'll build on this with:

**JavaScript execution → asynchronous JavaScript → React rendering → state updates → re-rendering.**

---

# One thing to remember

Don't memorize the diagrams.

Just keep this mental model:

> **JavaScript executes code. Function calls create execution contexts, and the call stack keeps track of those active calls.**

We'll revisit this when we study closures and asynchronous JavaScript.

---

## Small exercise

Without running it, tell me the order in which the numbers are printed:

```javascript
function a() {
    console.log(1);
    b();
    console.log(2);
}

function b() {
    console.log(3);
}

console.log(4);
a();
console.log(5);
```

And, more importantly, **why** that order occurs.

The execution is:

`4 → a() → 1 → b() → 3 → back to a() → 2 → back to global → 5`

The important part is that after `b()` finishes, execution **returns to exactly where it left off in `a()`**.

---