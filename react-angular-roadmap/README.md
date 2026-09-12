The curriculum below is designed around one rule:

> **No concept will be used as if you already know it before we have introduced and explained it.**

And this is **not a crash course**. We will take the time to understand the language, then TypeScript, then React and Angular properly.

# Roadmap

```text
JAVASCRIPT
    │
    ├── Foundations
    ├── Core language
    ├── Functions & scope
    ├── Objects & prototypes
    ├── Modern JavaScript
    ├── Modules
    ├── Asynchronous JavaScript
    ├── Browser APIs
    └── Practical JavaScript
             │
             ▼
TYPESCRIPT
    │
    ├── TypeScript foundations
    ├── Type system
    ├── Functions & objects
    ├── Generics
    ├── Advanced types
    ├── Classes
    ├── Modules & configuration
    └── Practical TypeScript
             │
       ┌─────┴─────┐
       ▼           ▼
    REACT       ANGULAR
       │           │
       └─────┬─────┘
             ▼
      Comparison & Architecture
```

---

# Part I — JavaScript

## Chapter 1 — JavaScript and the Runtime

Before writing substantial JavaScript, we establish what we're actually working with.

### 1.1 What JavaScript is
- JavaScript as a programming language
- JavaScript vs ECMAScript
- What the language defines
- What the runtime provides

### 1.2 JavaScript runtimes
- What a runtime is
- Browser runtimes
- Node.js
- JavaScript engines
- V8 at a high level

### 1.3 How JavaScript code is executed
- Source code
- Parsing
- Execution
- Global execution context
- Function execution context

### 1.4 The call stack
- What the stack represents
- Function calls
- Returning from functions
- Nested calls
- Recursion and stack overflow

### 1.5 Statements and expressions
- What a statement is
- What an expression is
- Expression statements
- Why the distinction matters

### 1.6 Strict mode
- `"use strict"`
- What strict mode changes
- Why modern JavaScript behaves differently from older JavaScript

---

# Chapter 2 — Values, Types, Variables, and Assignment

This is the **real beginning of JavaScript programming**.

### 2.1 Values
- What a value is
- Literal values
- Numbers
- Strings
- Booleans

### 2.2 Other primitive values
- `undefined`
- `null`
- `bigint`
- `symbol`
- Why these types exist

### 2.3 Variables
- What a variable represents
- Declaring a variable
- `let`
- `const`
- `var`

### 2.4 Assignment
- Assignment with `=`
- Reassignment
- Assignment expressions
- What happens when a variable is assigned a new value

### 2.5 Constants
- What `const` actually guarantees
- Why `const` doesn't mean immutable

### 2.6 Variable naming
- Identifier rules
- Naming conventions
- Reserved words

### 2.7 Inspecting values
- `console.log`
- `typeof`
- Understanding what `typeof` actually tells us
- Its historical quirks

---

# Chapter 3 — Operators and Expressions

Now that we know values and variables, we can meaningfully manipulate them.

### 3.1 Arithmetic operators
- `+`
- `-`
- `*`
- `/`
- `%`
- `**`

### 3.2 Assignment operators
- `=`
- `+=`
- `-=`
- `*=`
- `/=`
- `++`
- `--`

### 3.3 Comparison
- `<`
- `>`
- `<=`
- `>=`
- `===`
- `!==`
- `==`
- `!=`

### 3.4 Logical operators
- `&&`
- `||`
- `!`

### 3.5 Short-circuit evaluation
- How `&&` and `||` actually evaluate
- Why they don't necessarily produce booleans

### 3.6 Conditional expressions
- Ternary operator
- When it is appropriate
- Readability considerations

### 3.7 Operator precedence
- How JavaScript determines evaluation order
- Parentheses

---

# Chapter 4 — Type Conversion and Coercion

Now that operators are understood, we can explain one of JavaScript's most notorious behaviors.

### 4.1 Explicit conversion
- Number conversion
- String conversion
- Boolean conversion

### 4.2 Truthy and falsy values
- What truthiness means
- Falsy values
- How truthiness affects conditions

### 4.3 Implicit coercion
- What coercion means
- When JavaScript converts values automatically

### 4.4 Equality
- Strict equality
- Loose equality
- Why `===` is generally preferred

### 4.5 Common coercion pitfalls

---

# Chapter 5 — Control Flow

Only now do we introduce conditional and repetitive execution.

### 5.1 `if`
### 5.2 `else`
### 5.3 `else if`
### 5.4 Nested conditions
### 5.5 `switch`
### 5.6 `while`
### 5.7 `do...while`
### 5.8 `for`
### 5.9 `break`
### 5.10 `continue`
### 5.11 `for...of`
### 5.12 `for...in`

We'll distinguish carefully between these loops rather than treating them as interchangeable syntax.

---

# Chapter 6 — Functions

Now we introduce one of the most important concepts in JavaScript.

### 6.1 Why functions exist
### 6.2 Function declarations
### 6.3 Parameters
### 6.4 Arguments
### 6.5 Return values
### 6.6 Functions without return values
### 6.7 Function expressions
### 6.8 Arrow functions
### 6.9 Default parameters
### 6.10 Rest parameters
### 6.11 Functions as values
### 6.12 Passing functions to other functions

At this point we'll establish the foundation needed for higher-order functions, callbacks, closures, and React.

---

# Chapter 7 — Scope and Closures

Now that functions exist, scope can be properly explained.

### 7.1 Global scope
### 7.2 Function scope
### 7.3 Block scope
### 7.4 Lexical scope
### 7.5 Scope chains
### 7.6 Variable lookup
### 7.7 Shadowing
### 7.8 Closures
### 7.9 Why closures exist
### 7.10 Practical uses of closures

This chapter will be detailed. **Closures are foundational to understanding modern JavaScript and React.**

---

# Chapter 8 — Objects

Now we introduce JavaScript's object model properly.

### 8.1 What an object is
### 8.2 Object literals
### 8.3 Properties
### 8.4 Reading properties
### 8.5 Writing properties
### 8.6 Adding and deleting properties
### 8.7 Dot notation
### 8.8 Bracket notation
### 8.9 Computed property names
### 8.10 Methods
### 8.11 Property shorthand
### 8.12 Object references
### 8.13 Mutation vs reassignment

---

# Chapter 9 — Arrays

### 9.1 What arrays are
### 9.2 Creating arrays
### 9.3 Indexes
### 9.4 Reading and changing elements
### 9.5 Array length
### 9.6 Adding/removing elements
- `push`
- `pop`
- `shift`
- `unshift`
- `splice`

### 9.7 Iterating arrays
### 9.8 Array destructuring
### 9.9 Arrays containing objects
### 9.10 Arrays and references

---

# Chapter 10 — Object and Array Techniques

Only after objects and arrays are understood independently do we introduce these patterns.

### 10.1 Spread syntax
### 10.2 Rest syntax
### 10.3 Shallow copying
### 10.4 Deep copying
### 10.5 Nested objects
### 10.6 Nested arrays
### 10.7 Immutability
### 10.8 Object destructuring
### 10.9 Combining destructuring and defaults

This is where the earlier `user1/user2/address` example belongs.

---

# Chapter 11 — Built-in Objects and Data Structures

### 11.1 `String`
### 11.2 `Number`
### 11.3 `Math`
### 11.4 `Date`
### 11.5 `Map`
### 11.6 `Set`
### 11.7 WeakMap
### 11.8 WeakSet
### 11.9 JSON
### 11.10 When to use each

---

# Chapter 12 — Prototypes and the JavaScript Object Model

This comes **after** objects, because otherwise prototypes make no sense.

### 12.1 What a prototype is
### 12.2 Prototype chains
### 12.3 Property lookup through the chain
### 12.4 `Object.create`
### 12.5 Constructor functions
### 12.6 The `prototype` property
### 12.7 `new`
### 12.8 Classes as syntax over JavaScript's prototype system
### 12.9 `instanceof`

---

# Chapter 13 — `this`

Only after objects, methods, functions, and prototypes are established.

### 13.1 What `this` represents
### 13.2 `this` in method calls
### 13.3 `this` in regular functions
### 13.4 `this` in arrow functions
### 13.5 `call`
### 13.6 `apply`
### 13.7 `bind`
### 13.8 Common `this` mistakes

---

# Chapter 14 — Classes

### 14.1 Class syntax
### 14.2 Constructors
### 14.3 Instance methods
### 14.4 Fields
### 14.5 Static members
### 14.6 Getters and setters
### 14.7 Private fields
### 14.8 Inheritance
### 14.9 `super`
### 14.10 Composition vs inheritance

---

# Chapter 15 — Higher-Order Functions and Functional JavaScript

Now the prerequisites are finally in place.

### 15.1 Functions as values
### 15.2 Higher-order functions
### 15.3 Callbacks
### 15.4 `map`
### 15.5 `filter`
### 15.6 `find`
### 15.7 `some`
### 15.8 `every`
### 15.9 `reduce`
### 15.10 Sorting
### 15.11 Chaining operations
### 15.12 Choosing appropriate array operations

This chapter will be especially relevant to React.

---

# Chapter 16 — Modules

### 16.1 Why modules exist
### 16.2 `export`
### 16.3 `import`
### 16.4 Named exports
### 16.5 Default exports
### 16.6 Module boundaries
### 16.7 File organization
### 16.8 Circular dependencies
### 16.9 How modules are used in real applications

---

# Chapter 17 — Error Handling

### 17.1 Errors
### 17.2 `throw`
### 17.3 `try`
### 17.4 `catch`
### 17.5 `finally`
### 17.6 Built-in error types
### 17.7 Creating errors
### 17.8 Error propagation
### 17.9 Designing useful error handling

---

# Chapter 18 — Asynchronous JavaScript

Only now do we introduce asynchronous execution, because we already understand the synchronous execution model.

### 18.1 Why asynchronous programming exists
### 18.2 Blocking vs non-blocking work
### 18.3 Timers
### 18.4 The event loop
### 18.5 Task queue
### 18.6 Microtask queue
### 18.7 How the runtime interacts with JavaScript

---

# Chapter 19 — Promises

### 19.1 What a Promise represents
### 19.2 Pending/fulfilled/rejected
### 19.3 Creating Promises
### 19.4 Consuming Promises
### 19.5 `.then`
### 19.6 `.catch`
### 19.7 `.finally`
### 19.8 Promise chaining
### 19.9 Error propagation
### 19.10 Promise concurrency

---

# Chapter 20 — `async` / `await`

### 20.1 `async` functions
### 20.2 `await`
### 20.3 Error handling with `try/catch`
### 20.4 Sequential asynchronous operations
### 20.5 Parallel operations
### 20.6 `Promise.all`
### 20.7 `Promise.allSettled`
### 20.8 Common async mistakes

---

# Chapter 21 — Browser JavaScript

Now we move from the language itself into the browser environment.

### 21.1 The DOM
### 21.2 Selecting elements
### 21.3 Creating/changing elements
### 21.4 Events
### 21.5 Event listeners
### 21.6 Event bubbling
### 21.7 Event capturing
### 21.8 Event delegation
### 21.9 Forms
### 21.10 Browser storage

---

# Chapter 22 — HTTP and APIs from JavaScript

### 22.1 HTTP fundamentals
### 22.2 Requests and responses
### 22.3 `fetch`
### 22.4 HTTP methods
### 22.5 Headers
### 22.6 JSON
### 22.7 Status codes
### 22.8 Handling API errors
### 22.9 Authentication basics
### 22.10 Consuming a REST API

---

# Chapter 23 — JavaScript Tooling

### 23.1 npm
### 23.2 `package.json`
### 23.3 Dependencies
### 23.4 Development dependencies
### 23.5 npm scripts
### 23.6 Bundlers/build tools
### 23.7 Vite
### 23.8 Environment variables
### 23.9 Linting
### 23.10 Formatting

---

# Chapter 24 — JavaScript Testing

### 24.1 Why testing exists
### 24.2 Unit tests
### 24.3 Assertions
### 24.4 Test structure
### 24.5 Mocking
### 24.6 Testing asynchronous code
### 24.7 Practical JavaScript testing

---

# Chapter 25 — JavaScript Consolidation

A small application tying together:

- modules
- objects
- arrays
- functions
- asynchronous code
- API calls
- error handling
- browser interaction
- testing

This gives us a proper JavaScript foundation before moving to TypeScript.

---

# Part II — TypeScript

TypeScript comes **after JavaScript**, not alongside it.

## Chapter 26 — Why TypeScript?

- Problems TypeScript addresses
- Static vs dynamic typing
- Type checking
- Compilation/transpilation
- TypeScript vs JavaScript
- What TypeScript does and does not do at runtime

---

# Chapter 27 — TypeScript Fundamentals

### 27.1 Type annotations
### 27.2 Primitive types
### 27.3 Arrays
### 27.4 Tuples
### 27.5 Objects
### 27.6 Function parameter types
### 27.7 Return types
### 27.8 Optional parameters
### 27.9 Optional properties
### 27.10 Type inference

---

# Chapter 28 — Type System Fundamentals

### 28.1 Union types
### 28.2 Literal types
### 28.3 Type aliases
### 28.4 Interfaces
### 28.5 `null` and `undefined`
### 28.6 `any`
### 28.7 `unknown`
### 28.8 `never`
### 28.9 `void`

---

# Chapter 29 — Narrowing

### 29.1 Why narrowing is necessary
### 29.2 `typeof`
### 29.3 Equality narrowing
### 29.4 Truthiness narrowing
### 29.5 Property checks
### 29.6 Type predicates
### 29.7 Discriminated unions

---

# Chapter 30 — TypeScript Functions

### 30.1 Function types
### 30.2 Callback types
### 30.3 Optional parameters
### 30.4 Default parameters
### 30.5 Rest parameters
### 30.6 Overloads
### 30.7 Generic functions

---

# Chapter 31 — TypeScript Objects and Interfaces

### 31.1 Structural typing
### 31.2 Interfaces
### 31.3 Extending interfaces
### 31.4 Readonly properties
### 31.5 Index signatures
### 31.6 Interfaces vs type aliases
### 31.7 Modeling real application data

---

# Chapter 32 — Generics

### 32.1 Why generics exist
### 32.2 Generic functions
### 32.3 Generic types
### 32.4 Generic interfaces
### 32.5 Generic constraints
### 32.6 Multiple type parameters
### 32.7 Practical generic designs

---

# Chapter 33 — Advanced TypeScript

### 33.1 `keyof`
### 33.2 `typeof`
### 33.3 Indexed access types
### 33.4 `in`
### 33.5 Mapped types
### 33.6 Conditional types
### 33.7 Template literal types
### 33.8 Utility types

---

# Chapter 34 — TypeScript Classes

### 34.1 Classes
### 34.2 Access modifiers
### 34.3 `readonly`
### 34.4 Constructors
### 34.5 Inheritance
### 34.6 Abstract classes
### 34.7 Interfaces and classes
### 34.8 Composition

---

# Chapter 35 — TypeScript Configuration and Modules

### 35.1 TypeScript compiler
### 35.2 `tsconfig.json`
### 35.3 Strict mode
### 35.4 Module configuration
### 35.5 Target
### 35.6 Source maps
### 35.7 Type declarations
### 35.8 `.d.ts` files
### 35.9 Using JavaScript libraries from TypeScript

---

# Chapter 36 — Practical TypeScript

### 36.1 Typing API responses
### 36.2 Typing asynchronous code
### 36.3 Error handling
### 36.4 Generic API utilities
### 36.5 Type-safe application models
### 36.6 TypeScript testing
### 36.7 TypeScript project

---

# Part III — React

Only now do we introduce React.

And we'll use **TypeScript with React**, rather than learning React in JavaScript and then having to relearn everything.

## Chapter 37 — React Fundamentals

- What React is
- Why React exists
- React's mental model
- Creating a React project
- Project structure
- Components
- JSX
- JSX vs HTML
- Expressions inside JSX

---

# Chapter 38 — Components and Props

- Component composition
- Props
- Prop types
- Passing data
- Nested components
- `children`
- Component boundaries
- Reusable components

---

# Chapter 39 — Rendering

- How React renders
- Conditional rendering
- Rendering lists
- Keys
- Component identity
- Re-rendering
- Render phase vs DOM updates

---

# Chapter 40 — Events and State

- Event handling
- State
- `useState`
- State updates
- Functional updates
- State and rendering
- Objects in state
- Arrays in state
- Immutable state updates

---

# Chapter 41 — React Forms

- Controlled components
- Input state
- Multiple inputs
- Form submission
- Validation
- Errors
- Reusable form components

---

# Chapter 42 — Effects and Synchronization

- `useEffect`
- Dependencies
- Cleanup
- Effects vs event handlers
- Synchronizing with external systems
- Fetching data
- Common effect mistakes

This chapter will be particularly detailed because misunderstanding `useEffect` causes a huge amount of bad React code.

---

# Chapter 43 — React Hooks

- Rules of hooks
- `useRef`
- `useMemo`
- `useCallback`
- `useReducer`
- Custom hooks
- Designing useful hooks
- When **not** to use a hook

---

# Chapter 44 — Context and State Architecture

- Context
- Providers
- Consuming context
- Context limitations
- Local vs shared state
- Derived state
- State lifting
- Application state architecture

---

# Chapter 45 — React + APIs

- HTTP requests
- Loading state
- Error state
- Successful data
- Request cancellation
- Authentication
- API abstraction
- Typed API responses

---

# Chapter 46 — Routing

- Client-side routing
- Routes
- Route parameters
- Nested routes
- Navigation
- Protected routes
- URL state

---

# Chapter 47 — React Architecture

- Component organization
- Feature-based organization
- Reusable components
- Presentational vs logic concerns
- Custom hooks
- API/service layers
- Avoiding over-abstraction

---

# Chapter 48 — React Performance

- Why components re-render
- Reference equality
- Memoization
- `React.memo`
- `useMemo`
- `useCallback`
- Expensive computations
- Avoiding premature optimization
- React DevTools

---

# Chapter 49 — React Testing

- Component testing
- User interaction testing
- Mocking APIs
- Testing hooks
- Integration testing
- What not to test

---

# Chapter 50 — React Production Application

Build a reasonably substantial application using:

- React
- TypeScript
- routing
- forms
- API calls
- authentication
- reusable components
- state management
- testing
- proper project structure

---

# Part IV — Angular

Then we learn Angular **from its own architectural model**, rather than trying to translate React concepts directly into Angular.

## Chapter 51 — Angular Fundamentals

- What Angular is
- Angular's philosophy
- Angular CLI
- Project structure
- Standalone Angular
- Components
- Templates

---

# Chapter 52 — Angular Templates

- Interpolation
- Property binding
- Event binding
- Attribute binding
- Two-way binding
- Template expressions
- Conditional rendering
- List rendering
- Modern Angular control flow

---

# Chapter 53 — Components

- Component metadata
- Component lifecycle
- Inputs
- Outputs
- Component communication
- Content projection
- Component composition

---

# Chapter 54 — Services and Dependency Injection

- Why services exist
- Dependency injection
- Providers
- Injectors
- Service lifetime
- Dependency graphs
- Practical application architecture

---

# Chapter 55 — Angular Forms

### Template-driven forms
- Controls
- Validation
- Submission

### Reactive forms
- `FormControl`
- `FormGroup`
- `FormArray`
- Validators
- Custom validators
- Dynamic forms

---

# Chapter 56 — HTTP

- `HttpClient`
- GET/POST/PUT/PATCH/DELETE
- Typed responses
- Request configuration
- Error handling
- Interceptors
- Authentication

---

# Chapter 57 — RxJS

We will learn the RxJS concepts Angular actually requires.

- Observables
- Subscriptions
- Observable lifecycle
- Operators
- `map`
- `filter`
- `tap`
- `switchMap`
- `mergeMap`
- `concatMap`
- `catchError`
- `combineLatest`
- Subjects
- BehaviorSubject
- Unsubscription

We won't attempt to memorize the entire RxJS ecosystem.

---

# Chapter 58 — Angular Routing

- Router
- Routes
- Navigation
- Parameters
- Query parameters
- Child routes
- Lazy loading
- Guards
- Resolvers

---

# Chapter 59 — Angular Signals

- Why signals exist
- Signals
- Reading/writing signals
- Computed signals
- Effects
- Signals vs RxJS
- Signals and component state
- Modern Angular patterns

---

# Chapter 60 — Angular Change Detection

- How Angular detects changes
- Component rendering
- Change detection strategy
- `OnPush`
- Signals and change detection
- Performance considerations

---

# Chapter 61 — Angular Architecture

- Feature organization
- Shared functionality
- Services
- Components
- State management
- Dependency injection
- Smart/dumb component considerations
- Avoiding over-engineering

---

# Chapter 62 — Angular Testing

- Unit testing
- Component testing
- Service testing
- HTTP testing
- Dependency injection in tests
- Integration testing

---

# Chapter 63 — Angular Production Application

Build a substantial Angular application using:

- TypeScript
- standalone components
- routing
- services
- dependency injection
- forms
- HTTP
- RxJS
- signals
- authentication
- testing
- production-oriented architecture

---

# Part V — React and Angular Together

Finally, we'll deliberately compare them.

## Chapter 64 — React vs Angular

### Components
- React components
- Angular components

### State
- React state
- Angular signals
- Angular services/state

### Data flow
- Props
- Inputs/Outputs

### Dependency injection
- Angular DI
- React alternatives

### Side effects
- React effects
- Angular effects
- RxJS

### Forms
- React forms
- Angular forms

### Routing
- React ecosystem
- Angular Router

### HTTP
- React ecosystem
- Angular HttpClient

### Architecture
- Library/ecosystem approach
- Full-framework approach

---

# Chapter 65 — Choosing Between React and Angular

We'll finish by looking at the practical engineering question:

- When React is a better fit
- When Angular is a better fit
- Team considerations
- Project size
- Architecture
- Ecosystem
- Hiring/skills
- Maintainability
- Enterprise applications
- Startup applications
- Performance
- Developer experience

And we'll build one small feature **once in React and once in Angular** so the differences become concrete.

---

I'm going to treat this as a **dependency graph**, not just a list of topics.

For example, we will **not** do:

> "Here's `useState`. Here's a state example."

before you've learned:

> variables → assignment → objects → references → functions → closures → immutability

Likewise, we won't introduce:

> `switchMap`

before you've learned:

> asynchronous JavaScript → Promises → Observables → subscriptions → operators.

And we won't throw TypeScript generics at you before you've understood:

> types → functions → objects → interfaces → unions → narrowing.
