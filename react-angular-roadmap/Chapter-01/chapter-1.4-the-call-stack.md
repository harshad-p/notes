# JavaScript — 1.4 The Call Stack

The call stack is one of the most important pieces of the JavaScript execution model.

It explains:

- how JavaScript keeps track of which function is currently executing,
- how it knows where to return after a function finishes,
- how nested function calls work,
- why recursion can eventually produce `Maximum call stack size exceeded`.

---

## 1. What is a stack?

A **stack** is a data structure with **LIFO** semantics:

> **Last In, First Out**

Think of a stack of plates:

```text
      ┌─────────┐
      │ Plate C │ ← last added, first removed
      ├─────────┤
      │ Plate B │
      ├─────────┤
      │ Plate A │
      └─────────┘
```

If you add another plate, it goes on top.

If you remove one, you remove the top one first.

The two fundamental operations are usually described as:

- **push** — add something to the top
- **pop** — remove something from the top

JavaScript's **call stack** follows this LIFO model.

---

# 2. What does the call stack represent?

The **call stack keeps track of the currently active function calls**.

This is important because when one function calls another function, JavaScript needs to remember:

> "I'm currently inside this function, and when the called function finishes, I need to resume here."

Each active function invocation has an associated **stack frame**.

So conceptually:

```text
Call Stack
┌──────────────────────┐
│ current function     │ ← top
├──────────────────────┤
│ function that called │
├──────────────────────┤
│ caller               │
├──────────────────────┤
│ global execution     │
└──────────────────────┘
```

The exact implementation is engine-specific, so don't think of this as a literal JavaScript object containing these things.

It's a **mental model of how execution state is managed**.

---

# 3. Connecting this to execution contexts

This is where the previous lesson becomes useful.

We learned:

> A function invocation gets a function execution context.

Now we can connect that to the call stack.

When a function starts executing, its execution context is associated with an active **stack frame**.

Conceptually:

```text
Function invocation
       ↓
Function execution context
       ↓
Stack frame on call stack
```

So when you see:

> "Function A calls function B"

you should mentally think:

```text
A is executing
   ↓
B is called
   ↓
B becomes the active call
   ↓
B finishes
   ↓
A resumes
```

The call stack gives us the mechanism for tracking that nesting.

---

# 4. A simple function call

Consider:

```js
function greet() {
    console.log("Hello");
}

greet();
```

When the program reaches `greet()`, the function is invoked.

Conceptually:

```text
Before call:

Call Stack
┌───────────────┐
│ Global        │
└───────────────┘
```

Then `greet()` is called:

```text
Call Stack
┌───────────────┐
│ greet()       │ ← currently executing
├───────────────┤
│ Global        │
└───────────────┘
```

`greet()` executes.

When it returns:

```text
Call Stack
┌───────────────┐
│ Global        │
└───────────────┘
```

The `greet()` frame is popped.

Execution continues in the caller.

---

# 5. Returning from a function

Consider:

```js
function calculate(x) {
    return x * 2;
}

const result = calculate(10);
```

When `calculate(10)` begins:

```text
Call Stack

┌──────────────────┐
│ calculate(10)    │ ← top
├──────────────────┤
│ Global           │
└──────────────────┘
```

The function executes:

```js
return x * 2;
```

It produces `20`.

The function finishes.

Its stack frame is removed:

```text
Call Stack

┌──────────────────┐
│ Global           │
└──────────────────┘
```

The caller then receives the return value:

```js
const result = 20;
```

The important idea is:

> **Returning from a function removes that function's active call from the stack and resumes the caller.**

---

# 6. Nested function calls

This is where the stack becomes particularly useful.

Consider:

```js
function getPrice() {
    return calculatePrice();
}

function calculatePrice() {
    return applyDiscount();
}

function applyDiscount() {
    return 100;
}

const price = getPrice();
```

When `getPrice()` is called:

```text
Call Stack

┌─────────────────┐
│ getPrice()      │ ← top
├─────────────────┤
│ Global          │
└─────────────────┘
```

`getPrice()` calls `calculatePrice()`:

```text
┌─────────────────────┐
│ calculatePrice()    │ ← top
├─────────────────────┤
│ getPrice()          │
├─────────────────────┤
│ Global              │
└─────────────────────┘
```

Then `calculatePrice()` calls `applyDiscount()`:

```text
┌─────────────────────┐
│ applyDiscount()     │ ← top
├─────────────────────┤
│ calculatePrice()    │
├─────────────────────┤
│ getPrice()          │
├─────────────────────┤
│ Global              │
└─────────────────────┘
```

Now `applyDiscount()` returns.

It is popped:

```text
┌─────────────────────┐
│ calculatePrice()    │ ← top
├─────────────────────┤
│ getPrice()          │
├─────────────────────┤
│ Global              │
└─────────────────────┘
```

Then `calculatePrice()` returns.

```text
┌─────────────────────┐
│ getPrice()          │ ← top
├─────────────────────┤
│ Global              │
└─────────────────────┘
```

Then `getPrice()` returns.

```text
┌─────────────────────┐
│ Global              │
└─────────────────────┘
```

Finally the global code continues.

This is exactly the **LIFO** behavior of a stack.

---

# 7. Why does JavaScript need a call stack?

Because function calls can be nested arbitrarily deeply.

Suppose:

```text
A → B → C → D
```

When `D` finishes, JavaScript needs to know:

> "Where do I return?"

It needs to return to `C`.

Then when `C` finishes:

> "Where do I return?"

Back to `B`.

Then `B` → `A`.

The stack naturally solves this.

The most recently called function is always at the top.

So:

```text
A
 ↓
B
 ↓
C
 ↓
D
```

becomes:

```text
D   ← return here first
C
B
A
```

That's why a stack is an appropriate data structure for managing nested function execution.

---

# 8. Recursion

Now we reach one of the most interesting consequences.

**Recursion** occurs when a function calls itself, directly or indirectly.

For example:

```js
function countDown(n) {
    if (n === 0) {
        return;
    }

    countDown(n - 1);
}

countDown(3);
```

The calls build up like this:

```text
countDown(3)
    ↓
countDown(2)
    ↓
countDown(1)
    ↓
countDown(0)
```

While those calls are active, the stack conceptually looks like:

```text
┌──────────────────┐
│ countDown(0)     │
├──────────────────┤
│ countDown(1)     │
├──────────────────┤
│ countDown(2)     │
├──────────────────┤
│ countDown(3)     │
├──────────────────┤
│ Global           │
└──────────────────┘
```

Once `countDown(0)` returns, the calls unwind:

```text
countDown(0) returns
        ↓
countDown(1) returns
        ↓
countDown(2) returns
        ↓
countDown(3) returns
        ↓
Global continues
```

This is called **stack unwinding**.

---

# 9. Stack overflow

The call stack has a finite capacity.

If you keep adding function calls without allowing them to return, eventually the stack becomes too large.

For example:

```js
function forever() {
    forever();
}

forever();
```

There is no terminating condition.

Conceptually:

```text
forever()
forever()
forever()
forever()
forever()
...
```

The stack keeps growing:

```text
┌──────────────┐
│ forever()    │
├──────────────┤
│ forever()    │
├──────────────┤
│ forever()    │
├──────────────┤
│ forever()    │
├──────────────┤
│ ...          │
└──────────────┘
```

Eventually the JavaScript engine cannot accommodate another call.

You typically get an error such as:

```text
RangeError: Maximum call stack size exceeded
```

The exact error wording can vary by engine.

This is **stack overflow**.

---

# 10. Stack overflow isn't only caused by recursion

Recursion is the classic example, but the underlying problem is:

> **Too many active function calls.**

For example, mutually recursive functions can also cause it:

```js
function a() {
    b();
}

function b() {
    a();
}

a();
```

The call sequence becomes:

```text
a()
 ↓
b()
 ↓
a()
 ↓
b()
 ↓
a()
 ↓
...
```

Neither function gets an opportunity to return.

The stack grows until it reaches its limit.

---

# 11. Stack frames

A **stack frame** represents an active function invocation on the call stack.

It conceptually contains information needed to resume/continue execution, such as:

- which function is executing
- its execution state
- local/parameter bindings
- where execution should continue after the function returns
- other engine-specific execution information

Don't interpret this as a strict specification saying:

> "Every stack frame must contain exactly these fields."

The JavaScript specification describes execution using higher-level concepts such as execution contexts; **the actual call-stack implementation is an engine concern**.

This distinction is useful:

```text
JavaScript specification
    ↓
Defines language semantics

JavaScript engine
    ↓
Implements those semantics
    ↓
May use a physical call stack / optimized representations
```

So when we say:

> "JavaScript has a call stack"

we're using a very useful model of how engines manage synchronous execution.

---

# 12. The call stack is synchronous

One particularly important property:

> **A synchronous function call stays on the call stack until it returns.**

For example:

```js
function a() {
    b();
}

function b() {
    c();
}

function c() {
    console.log("done");
}

a();
```

During `c()`:

```text
c()
b()
a()
Global
```

All of those calls are still active.

JavaScript doesn't execute `c()` independently and somehow forget about `b()` and `a()`.

It has to preserve their execution state so it can return through them.

This becomes extremely important when we later study **asynchronous JavaScript and the event loop**.

But we won't jump into that yet.

---

# 13. A practical debugging connection

The call stack isn't merely an academic concept.

You've probably seen it indirectly when debugging applications.

Suppose an error occurs deep inside a call chain:

```text
API handler
   ↓
Service
   ↓
Repository
   ↓
Database helper
   ↓
Error
```

A stack trace can show something conceptually like:

```text
DatabaseHelper()
Repository()
Service()
ApiHandler()
```

This tells you the chain of function calls that led to the failure.

In browser DevTools, Node.js, and many logging systems, **stack traces are an important debugging tool**.

The stack trace and the live call stack aren't exactly the same thing, but they are closely related: the trace records the relevant chain of execution leading to an error.

---

# Interview questions

### 1. What is the call stack?

> The call stack is a LIFO structure used by the JavaScript engine to keep track of active function calls and their execution state.

---

### 2. What happens when a function is called?

> A new stack frame is created for that invocation and pushed onto the call stack. The function executes, and when it returns, that frame is removed.

---

### 3. What happens when a function returns?

> Its active stack frame is popped from the call stack, and execution resumes in the caller.

---

### 4. Why is the call stack LIFO?

Because nested calls must return in reverse order.

If:

```text
A → B → C
```

then:

```text
C returns
B returns
A returns
```

The last call added is therefore the first one that finishes.

---

### 5. What causes a stack overflow?

> Stack overflow occurs when the number/depth of active calls exceeds the available call-stack capacity, commonly due to unbounded recursion or mutual recursion.

---

### 6. What is a stack frame?

> A stack frame represents an active function invocation and contains execution state needed by the engine to manage that invocation and eventually resume its caller.

---

### 7. Is the call stack part of the JavaScript language specification?

Be careful here.

A strong answer is:

> The ECMAScript specification defines execution semantics using concepts such as execution contexts, but the concrete call-stack implementation is an engine/runtime implementation detail.

That's a good senior-level distinction.

---

# One important connection

You now have the first three pieces of the execution model:

```text
Source code
     ↓
Parsing
     ↓
Execution context
     ↓
Function call
     ↓
Stack frame
     ↓
Call stack
```

And:

```text
Function calls
     ↓
push frames
     ↓
nested execution
     ↓
return
     ↓
pop frames
```

This gives us the foundation for the next lesson.

---

## Exercise

Without running the code, determine the **maximum call-stack depth** reached during the execution of this program, counting `global` as the initial frame:

```js
function first() {
    second();
}

function second() {
    third();
}

function third() {
    return;
}

first();
```

Then answer:

1. What is the stack immediately after `first()` is called?
2. What is it when `third()` is executing?
3. In what order are the frames removed?
4. What does the stack look like after `first()` returns?

Then we'll move on to **1.5 — Statements vs Expressions**, where we'll shift from the execution model into the actual structure of JavaScript code.