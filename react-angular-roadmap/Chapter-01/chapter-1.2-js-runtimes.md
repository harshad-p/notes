# JavaScript — Chapter 1.2: JavaScript Runtimes

This chapter is about an important distinction:

> **JavaScript is the language. A JavaScript runtime is the environment that makes it possible to run JavaScript code.**

Understanding this distinction becomes important later when you work with **React, Node.js, browsers, APIs, and TypeScript**.

---

## 1. What is a JavaScript runtime?

A **runtime** is the environment in which a program executes.

For JavaScript, a runtime typically provides:

1. A **JavaScript engine** that executes JavaScript.
2. APIs and services that JavaScript code can use.
3. Mechanisms for interacting with the outside world.
4. The overall environment needed to run a JavaScript application.

A simplified picture is:

```text
                JavaScript Runtime
                       |
             +---------+---------+
             |                   |
      JavaScript Engine       Runtime APIs
             |                   |
          executes        provide capabilities
         JavaScript       outside the language
```

For example:

```js
console.log("Hello");
```

`console.log()` is a good example of why the distinction matters.

The **JavaScript language itself** does not define `console.log()` as part of its core language syntax.

The runtime provides `console`.

---

# 2. JavaScript engine vs JavaScript runtime

These terms are often used interchangeably in casual conversation, but they aren't the same thing.

### JavaScript engine

The engine's primary job is:

> **Take JavaScript code and execute it.**

Examples:

- **V8** — Chrome, Node.js, Deno, etc.
- **SpiderMonkey** — Firefox
- **JavaScriptCore** — Safari

The engine handles things such as:

- parsing JavaScript
- understanding the language
- executing code
- managing memory
- garbage collection
- optimizing frequently executed code

We'll go deeper into execution in the next chapter.

### Runtime

The runtime is broader.

It contains the engine **plus additional capabilities provided by the environment**.

For example:

```text
Browser Runtime
│
├── JavaScript Engine
│      └── V8 / SpiderMonkey / JavaScriptCore
│
├── DOM APIs
├── Web APIs
├── Timers
├── Fetch / networking
├── Storage APIs
└── Browser-specific functionality
```

Whereas:

```text
Node.js Runtime
│
├── V8 JavaScript Engine
├── File system APIs
├── HTTP APIs
├── Network APIs
├── Timers
├── Streams
├── Process APIs
└── Other Node.js functionality
```

So:

> **Engine = executes JavaScript.**  
> **Runtime = engine + environment-provided capabilities.**

That's a distinction worth remembering for interviews.

---

# 3. Browser runtimes

A browser provides a JavaScript runtime.

For example, Chrome provides a runtime built around **V8**.

Firefox uses **SpiderMonkey**, while Safari uses **JavaScriptCore**.

The browser runtime gives JavaScript access to browser functionality.

For example:

```js
document.getElementById("title");
```

`document` is not a JavaScript language feature.

It is provided by the **browser environment** and represents the web page's DOM.

Similarly:

```js
fetch("/api/users");
```

`fetch()` is provided by the web platform/runtime.

And:

```js
setTimeout(() => {
    console.log("Done");
}, 1000);
```

The timer capability is provided by the host environment.

This leads to a useful mental model:

```text
JavaScript language
        +
Browser APIs
        ↓
Browser JavaScript environment
```

---

# 4. Node.js

Node.js is another JavaScript runtime.

The important point is:

> **Node.js allows JavaScript to run outside a browser.**

Node.js uses **V8** as its JavaScript engine, but Node.js itself is much more than V8.

For example, Node.js provides APIs for:

```js
import fs from "fs";

const data = fs.readFileSync("data.txt", "utf8");
```

File-system access isn't something JavaScript inherently provides.

Node.js provides the `fs` API.

Similarly, Node.js can provide:

- HTTP servers
- file-system access
- networking
- streams
- processes
- environment variables
- command-line interaction

This is why you can build a backend using Node.js.

For example:

```text
                Node.js
                   |
          +--------+--------+
          |                 |
         V8            Node.js APIs
          |                 |
    JavaScript         filesystem
     execution          HTTP
                        networking
                        streams
                        etc.
```

---

# 5. Why can the same JavaScript run in different environments?

Because the **language** and the **environment** are separate.

Consider:

```js
const x = 10;

if (x > 5) {
    console.log("Large");
}
```

The language defines things such as:

- `const`
- numbers
- `if`
- comparison
- blocks
- function syntax
- objects
- arrays

Those are JavaScript language features.

But:

```js
document.title = "Hello";
```

depends on a browser.

And:

```js
fs.readFileSync("data.txt");
```

depends on Node.js.

So the same JavaScript language can be hosted by different environments.

```text
                 JavaScript
                     |
        +------------+------------+
        |                         |
   Browser Runtime          Node.js Runtime
        |                         |
       V8*                       V8
        |                         |
     DOM, Web APIs          Node.js APIs
     Browser APIs           filesystem, HTTP...
```

\* Chrome uses V8; other browsers use different engines.

---

# 6. V8 at a high level

Since you'll encounter V8 constantly in modern JavaScript development, it's worth understanding what it actually does.

**V8 is Google's open-source JavaScript engine.**

It is written primarily in C++ and is used by:

- Google Chrome
- Node.js
- several other JavaScript environments

At a very high level:

```text
JavaScript source code
        ↓
      V8
        ↓
 parse / understand code
        ↓
 execute code
        ↓
 optimize frequently executed code
```

V8 doesn't simply translate the entire JavaScript program into machine code once and then stop.

Modern engines use sophisticated techniques to make JavaScript fast, including:

- parsing
- bytecode generation
- interpretation
- runtime profiling
- JIT compilation/optimization
- deoptimization when assumptions become invalid
- garbage collection

We don't need to dive deeply into compiler internals yet. The important concept is:

> **V8 is responsible for understanding and executing JavaScript; Node.js or the browser supplies the surrounding environment.**

---

# 7. A subtle but important point: JavaScript doesn't "need a browser"

You'll sometimes hear:

> "JavaScript is a browser language."

Historically, JavaScript became famous through browsers, but this is no longer accurate.

JavaScript is a **general-purpose programming language**.

It can run in:

- browsers
- servers
- command-line programs
- desktop applications
- mobile environments
- embedded environments
- various other runtimes

The runtime determines what capabilities are available.

For example, this is valid JavaScript:

```js
const numbers = [1, 2, 3];

const doubled = numbers.map(n => n * 2);
```

It doesn't require a browser.

But this:

```js
document.querySelector("button");
```

requires an environment that provides `document`, such as a browser.

---

# 8. Host environment

You will also encounter the term **host environment**.

It essentially refers to the environment that hosts/runs the JavaScript program and provides capabilities beyond the core language.

For example:

```text
                 JavaScript
                     |
        +------------+------------+
        |                         |
    Browser                    Node.js
     host                       host
        |                         |
   DOM / Web APIs          Node.js APIs
        |                         |
        +-----------+-------------+
                    |
              JavaScript engine
```

The terminology can vary depending on the context, but the distinction is useful:

**JavaScript language** → defines the language.

**Engine** → executes the language.

**Host/runtime** → provides the surrounding capabilities.

---

# 9. One distinction we'll need later: runtime APIs vs language features

Consider:

```js
const users = [1, 2, 3];

users.map(x => x * 2);
```

`const`, arrays, arrow functions, and `map()` belong to JavaScript's language/built-in ecosystem.

Now consider:

```js
fetch("/api/users");
```

`fetch` is an API provided by the host environment.

And:

```js
setTimeout(() => {}, 1000);
```

The timer mechanism involves the host environment.

This distinction will become **very important when we study asynchronous JavaScript and the event loop**.

We'll return to it there rather than trying to explain the event loop prematurely.

---

# Interview Questions

These are the kinds of questions you should be able to answer after this chapter.

### 1. What is a JavaScript runtime?

A good interview answer:

> A JavaScript runtime is the environment that executes JavaScript. It includes a JavaScript engine plus APIs and other capabilities provided by the host environment.

---

### 2. What is the difference between a JavaScript engine and a runtime?

> The engine executes JavaScript code, while the runtime includes the engine plus the APIs and services provided by the environment in which JavaScript runs.

---

### 3. What is V8?

> V8 is Google's JavaScript engine. It executes JavaScript and is used by Chrome and Node.js, among other environments.

---

### 4. Is Node.js a JavaScript engine?

**No.**

Node.js is a **JavaScript runtime** that uses V8 as its JavaScript engine.

This is a very common interview question.

---

### 5. Is JavaScript a browser-only language?

**No.**

JavaScript is a programming language that can run in many environments. Browsers are one type of host environment, while Node.js is another.

---

### 6. Why does `document` work in a browser but not normally in Node.js?

Because `document` is part of the browser's DOM environment. Node.js does not provide the browser DOM by default.

---

### 7. Why can Node.js access files while browser JavaScript normally can't?

Because Node.js provides filesystem APIs. Browsers intentionally provide a different, restricted environment for security reasons.

---

## The mental model to keep

You don't need to memorize all the implementation details. Keep this hierarchy in your head:

```text
JavaScript
   │
   │ language
   ↓
JavaScript Engine
   │
   │ executes JavaScript
   ↓
Runtime / Host Environment
   │
   ├── Browser → DOM, Web APIs, etc.
   │
   └── Node.js → filesystem, server/network APIs, etc.
```

And the key example:

> **Chrome:** V8 + browser environment  
> **Node.js:** V8 + Node.js environment

That distinction will pay off repeatedly once we get into **async JavaScript, the event loop, Node.js, React, and APIs**.