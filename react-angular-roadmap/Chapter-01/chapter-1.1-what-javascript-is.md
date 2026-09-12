# JavaScript — Chapter 1.1: What JavaScript Is

Let's start from the actual foundation.

You already know how to write JavaScript, so some of this may feel familiar. That's fine. The goal here is to establish the terminology and mental model we'll build on later.

## 1. JavaScript is a programming language

JavaScript is a **programming language**. It defines things such as:

- values and types
- variables
- operators
- functions
- objects
- classes
- control flow
- modules
- error handling

For example:

```javascript
let age = 37;

if (age >= 18) {
    console.log("Adult");
}
```

The language defines what `let`, `if`, `>=`, and `console.log(...)` mean—or, more precisely, which parts belong to the language and which come from the environment.

That distinction is important.

---

## 2. JavaScript vs ECMAScript

You will sometimes hear **ECMAScript** instead of JavaScript.

ECMAScript is the **standard/specification** that defines the language.

JavaScript is the practical language implemented according to that specification.

You don't need to worry about the distinction in everyday programming, but it explains terminology such as:

> "ES2015 introduced `let`, `const`, classes, arrow functions, modules..."

ES2015 is a version of the ECMAScript standard.

Modern JavaScript generally means JavaScript based on modern ECMAScript features.

---

## 3. JavaScript doesn't execute itself

A `.js` file is just source code.

Something needs to **run** that code.

That something is a **JavaScript runtime**, which includes a JavaScript engine.

For example:

```text
JavaScript source code
        ↓
    JavaScript engine
        ↓
     execution
```

The engine is responsible for understanding and executing JavaScript.

One well-known engine is **V8**, developed by Google. Chrome uses V8, and Node.js uses V8 as well.

There are other engines, such as SpiderMonkey in Firefox and JavaScriptCore in Safari.

For our purposes, the important idea is:

> **Your JavaScript code needs an environment capable of executing it.**

---

# 4. Browser vs Node.js

This distinction will become very important once we start building React applications.

JavaScript can run in different environments.

### In a browser

A browser provides JavaScript with browser-specific capabilities.

For example:

```javascript
document.querySelector("button");
```

`document` is part of the browser environment. It represents the web page's DOM.

Similarly:

```javascript
fetch("/api/users");
```

`fetch` is provided by the environment.

### In Node.js

Node.js allows JavaScript to run outside the browser.

For example, Node provides APIs for things such as:

- files
- networking
- processes
- servers

So:

```text
             JavaScript
                 |
        +--------+--------+
        |                 |
     Browser           Node.js
        |                 |
       DOM             Files
      fetch           Servers
     Storage          Networking
```

The **JavaScript language itself is largely the same**.

The environment around it is different.

This is why code such as:

```javascript
document.querySelector(...)
```

works in a browser but doesn't normally exist in a Node.js program.

---

# 5. Language vs environment

This distinction is worth making very clear.

Suppose you write:

```javascript
const result = Math.max(10, 20);
```

`Math` is part of JavaScript's standard built-in functionality.

But:

```javascript
document.querySelector("#app");
```

depends on the browser environment.

And:

```javascript
import fs from "node:fs";
```

uses functionality provided by Node.js.

So when we learn JavaScript, we'll primarily be learning the **language** first.

Later we'll learn the environments that JavaScript commonly runs in.

---

# 6. What happens when JavaScript runs?

Take this simple program:

```javascript
const x = 10;
const y = 20;

const result = x + y;

console.log(result);
```

The JavaScript engine has to:

1. Read and understand the source code.
2. Determine what the code means.
3. Execute it.
4. Keep track of variables and function calls while doing so.

The execution process is considerably more sophisticated internally than that, but this is the basic model we'll use initially.

Later, we'll examine parts of this process in detail:

- execution contexts
- call stack
- scope
- closures
- event loop
- asynchronous execution

Those aren't separate magical features. They're pieces of the execution model.

---

# 7. A useful distinction: source code vs execution

When you write:

```javascript
const x = 10;
```

there are two different things to distinguish:

**Source code:**

```javascript
const x = 10;
```

and the **runtime state created while executing it**, which includes a variable named `x` associated with the value `10`.

This distinction will become useful when we later discuss:

- scope
- closures
- objects and references
- memory
- React state

---

# 8. Where we're going from here

We're going to build the execution model progressively.

Right now, the important mental model is:

```text
JavaScript language
       ↓
JavaScript runtime
       ↓
JavaScript engine
       ↓
execute source code
```

The browser or Node.js then provides additional APIs around the language.

We haven't yet discussed **how execution is organized internally**. That's the next lesson.

---

### What you should take away

1. **JavaScript is a programming language.**
2. **ECMAScript is the standard defining the language.**
3. **A JavaScript engine executes JavaScript.**
4. **A runtime provides the environment and additional APIs.**
5. **Browser and Node.js provide different environments around JavaScript.**

Next we'll go into **Chapter 1.2 — How JavaScript Code Is Executed**, where we'll introduce **execution contexts** properly.