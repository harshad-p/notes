# JavaScript — 1.3 How JavaScript Code Is Executed

You can write:

```js
const x = 10;

function calculate() {
    return x * 2;
}

const result = calculate();
```

As a developer, it's easy to think:

> "The runtime reads this file from top to bottom and executes each line."

That's useful as a rough mental model, but technically incomplete.

Before JavaScript can execute the code, the engine has to **parse it**, establish the necessary execution state, and then execute it within an **execution context**.

Let's build that model properly.

---

## 1. Source code

Source code is simply the JavaScript program you write:

```js
const x = 10;

function calculate() {
    return x * 2;
}

const result = calculate();
```

The engine receives this source text.

At this point, it's just text.

The engine needs to determine:

- Is this valid JavaScript?
- What constructs are present?
- Where are the declarations?
- What expressions need to be evaluated?
- What functions exist?
- How are the pieces of code structured?

That's the job of **parsing**.

---

# 2. Parsing

**Parsing** is the process of analyzing JavaScript source code according to the JavaScript language grammar and turning it into an internal representation that the engine can work with.

You don't directly interact with the parser, but it's a fundamental part of execution.

Conceptually:

```text
Source code
    ↓
  Parser
    ↓
Internal representation
    ↓
Execution
```

A common representation is an **Abstract Syntax Tree (AST)**.

You don't need to learn ASTs deeply yet, but you should understand what they represent.

For example:

```js
const x = 10;
```

contains several pieces of structure:

- a variable declaration
- an identifier `x`
- an initializer
- a numeric literal `10`

The parser recognizes that structure.

It isn't merely looking for text like `"const"` and `"10"`.

It understands the **syntactic structure of the program**.

---

## 3. Syntax errors can be detected during parsing

Consider:

```js
const x = ;
```

This doesn't form a valid JavaScript statement.

The parser can detect this before normal execution begins.

That's why something like:

```js
const x = ;
console.log("hello");
```

doesn't normally result in:

> "The first statement failed, but let's continue to the second one."

The source code has a syntax problem that prevents it from being parsed as valid JavaScript.

This is different from a runtime error such as:

```js
const user = null;

user.name;
```

Here, the code is syntactically valid.

The engine can parse it.

The problem occurs when the code is **executed**.

We'll cover errors in detail much later in Chapter 17.

For now, remember the distinction:

```text
Syntax error
    ↓
problem understanding the program structure

Runtime error
    ↓
valid code encounters a problem during execution
```

---

# 4. Execution

Once the engine has understood the program, it can execute it.

For example:

```js
const x = 10;

function calculate() {
    return x * 2;
}

const result = calculate();
```

Conceptually, execution involves:

```text
Create execution environment
        ↓
Execute declarations / statements
        ↓
Evaluate expressions
        ↓
Call functions when reached
        ↓
Continue execution
```

But here's where an important concept enters.

**Where does this execution happen?**

JavaScript uses something called an **execution context**.

---

# 5. What is an execution context?

An **execution context** is the environment in which a piece of JavaScript code is evaluated and executed.

It contains the execution-related state needed by the engine for that code.

You can think of it as:

> **The engine's working environment for executing a particular piece of JavaScript code.**

There are different kinds of execution contexts.

For the scope of this lesson, the important ones are:

1. **Global execution context**
2. **Function execution context**

There are other details, including how modern JavaScript handles modules, but we'll introduce those when modules become relevant.

---

# 6. Global execution context

When JavaScript begins executing a script, it needs an initial execution context.

That's the **global execution context**.

For:

```js
const x = 10;

console.log(x);
```

there is initially a global context in which this code executes.

Conceptually:

```text
Global Execution Context
│
├── information about global code
├── bindings such as x
├── global environment
└── other execution state
```

The global execution context represents the execution of the program at the global level.

It exists before any function is called.

---

# 7. Function execution context

Now consider:

```js
function calculate(x) {
    return x * 2;
}

const result = calculate(10);
```

When the engine reaches:

```js
calculate(10);
```

it needs an execution environment specifically for that invocation of `calculate`.

It creates a **function execution context**.

Conceptually:

```text
Global Execution Context
        │
        │ calls calculate(10)
        ↓
Function Execution Context
for calculate
        │
        ├── x = 10
        └── execute function body
```

The function's parameters and local variables belong to that function's execution environment.

This is one reason execution contexts are important.

---

# 8. Every function invocation gets its own function execution context

This is an important detail.

Consider:

```js
function double(x) {
    return x * 2;
}

const a = double(10);
const b = double(20);
```

There are **two invocations** of `double`.

Conceptually:

```text
Global Context
   │
   ├── double(10)
   │      ↓
   │   Function Context
   │   x = 10
   │
   └── double(20)
          ↓
       Function Context
       x = 20
```

The two invocations don't share the same parameter binding.

Each invocation gets its own execution context.

This becomes particularly important with recursion and nested function calls.

---

# 9. Execution context contains more than just variables

Don't reduce an execution context to:

> "A box containing local variables."

That's too simplistic.

An execution context involves execution state such as:

- the code being executed
- the current lexical environment
- variable/function bindings
- the relationship to outer lexical environments
- the value of `this` where applicable
- other internal state used by the language/runtime

Some of these concepts we'll study properly later.

For example, **lexical environments and scope chains** belong primarily to Chapter 7.

So you don't need to memorize their internal specification details right now.

The important point is:

> An execution context is the broader execution environment, not merely a local-variable container.

---

# 10. Global context → function context

Let's put everything together with one example:

```js
const multiplier = 2;

function calculate(value) {
    const result = value * multiplier;
    return result;
}

const answer = calculate(10);
```

At a high level:

### Step 1 — Source code

The engine receives the JavaScript source.

### Step 2 — Parsing

The engine parses the source and understands its structure.

### Step 3 — Global execution begins

A global execution context is established.

### Step 4 — Global code executes

The global code reaches:

```js
calculate(10);
```

### Step 5 — Function invocation

A new function execution context is created for this invocation.

That context contains the function's execution state, including:

```text
value = 10
```

The function executes:

```js
const result = value * multiplier;
```

### Step 6 — Function returns

The function produces:

```text
20
```

and its execution finishes.

### Step 7 — Global execution continues

The returned value is assigned to:

```js
answer
```

This gives us:

```text
                Global Execution Context
                         │
                         │ calculate(10)
                         ↓
                Function Execution Context
                         │
                    value = 10
                         │
                    result = 20
                         │
                       return
                         ↓
                Global execution resumes
```

---

# 11. But what happens to the function context?

Once the function finishes executing, its execution context is no longer actively executing.

For example:

```js
function calculate(value) {
    return value * 2;
}

const result = calculate(10);
```

The function context exists for the invocation.

After the function returns, execution proceeds back to the caller.

Exactly **how execution contexts are managed** is where the **call stack** comes in.

That's our next lesson.

For now, don't combine the concepts prematurely:

> **Execution context** describes the environment/state in which code executes.

> **Call stack** describes the mechanism used to keep track of active execution contexts.

We'll go into the call stack properly next.

---

# 12. One subtle point: declarations and execution aren't quite the same thing

Consider:

```js
console.log(message);

var message = "Hello";
```

If you only think:

> "JavaScript executes every line strictly from top to bottom."

you'll have trouble explaining JavaScript behavior like this.

The language's handling of declarations and bindings is more nuanced than simple line-by-line interpretation.

Likewise, function declarations have behavior that differs from ordinary executable statements.

For example:

```js
calculate();

function calculate() {
    console.log("Hello");
}
```

This works.

The function declaration is available when the code executes the call.

This is related to how the execution context is established and how declarations are instantiated.

**We will cover the exact behavior of `var`, `let`, `const`, and function declarations when we get to variables and scope.**

I don't want to prematurely teach hoisting here because we'd then be mixing several concepts before their proper lessons.

The important lesson for now is:

> JavaScript execution is more than simply reading source text one line at a time.

---

# 13. A useful mental model

At this stage, your mental model should be:

```text
JavaScript source
       ↓
    Parsing
       ↓
Program understood
       ↓
Global execution context
       ↓
Execute global code
       ↓
Function call
       ↓
Function execution context
       ↓
Execute function
       ↓
Return
       ↓
Continue caller's execution
```

And remember the distinction:

```text
Execution Context
    = environment/state for executing code

Call Stack
    = structure that tracks active execution contexts
```

We'll connect those two concepts in the next lesson.

---

# Interview questions

### 1. What is an execution context?

A solid answer:

> An execution context is the environment in which JavaScript code is evaluated and executed. It contains the execution state and bindings required for that code.

Don't say simply:

> "It's where variables are stored."

That's an incomplete explanation.

---

### 2. What is the global execution context?

> It is the execution context established for executing code at the global level of a JavaScript program.

---

### 3. What happens when a function is called?

At this level:

> JavaScript creates a new function execution context for that invocation and executes the function within that context.

The **call stack** part comes next.

---

### 4. Does every function call get its own execution context?

Yes.

Each **function invocation** gets its own function execution context.

---

### 5. What is parsing?

> Parsing is the process of analyzing JavaScript source code according to the language grammar and producing an internal representation that the engine can use to execute the program.

---

### 6. What's the difference between a syntax error and a runtime error?

> A syntax error means the source doesn't conform to JavaScript's grammar and cannot be properly parsed. A runtime error occurs while otherwise valid code is being executed.

---

### 7. Is JavaScript literally executed one line at a time?

A good interview answer is:

> That's a useful simplification, but not an accurate complete model. JavaScript source is first parsed, execution occurs within execution contexts, and declarations and other language semantics affect how execution proceeds.

That's a much stronger answer than simply saying "yes."

---

## Small exercise

Without running the code, reason about the **execution contexts** involved here:

```js
function calculate(x) {
    const doubled = x * 2;
    return doubled;
}

const first = calculate(5);
const second = calculate(10);
```

Answer these:

1. What execution context exists when the program starts?
2. How many **function execution contexts** are created?
3. Are the two calls to `calculate` using the same function execution context or different ones?
4. What value does `x` have in each invocation?
5. After each invocation returns, what happens to execution?

Don't worry about the call stack yet — we'll formally introduce that next.

- global execution context
- 2
- different
- 5, 10
- the function execution context is popped

Let's go through them precisely.

1. **What execution context exists when the program starts?**  
   ✅ **Global execution context.**

2. **How many function execution contexts are created?**  
   ✅ **2.**  
   `calculate(5)` creates one, and `calculate(10)` creates another.

3. **Same function execution context or different ones?**  
   ✅ **Different ones.**  
   Each **invocation** of a function gets its own function execution context. The function definition is the same, but the execution contexts are separate.

4. **What value does `x` have in each invocation?**  
   ✅ **5 and 10.**
   - First invocation → `x = 5`
   - Second invocation → `x = 10`

   Consequently, each invocation also has its own `doubled` variable.

5. **After each invocation returns, what happens to execution?**  
   ✅ **The function's execution context is removed from the call stack (popped), and execution resumes in the caller's context.**

So your mental model is now:

```text
Program starts
    ↓
Global execution context
    ↓
calculate(5)
    ↓
Function execution context #1
x = 5
doubled = 10
    ↓
return
    ↓
Function context #1 removed
    ↓
calculate(10)
    ↓
Function execution context #2
x = 10
doubled = 20
    ↓
return
    ↓
Function context #2 removed
    ↓
Global execution continues
```

One particularly important point: **the function execution context belongs to the invocation, not to the function definition.**

That's why this:

```js
calculate(5);
calculate(10);
```

creates **two** function execution contexts, even though there is only **one** function named `calculate`.