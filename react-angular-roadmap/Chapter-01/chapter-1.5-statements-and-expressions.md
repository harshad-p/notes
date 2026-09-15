Your planned topics are sufficient. I would add **one small sub-topic** because it makes the distinction much clearer:

- **Expression statements** — already on your list, but we'll explicitly distinguish the *expression* from the *statement containing it*.
- **Expressions as values** — not a separate major topic, but an important property we'll use to understand why expressions can appear in so many places.

I would **not** add operator precedence, coercion, or control flow yet. Those belong later in the curriculum.

# 1.5 Statements and Expressions

This distinction is fundamental to JavaScript syntax, and it becomes particularly useful when you start reading unfamiliar code, understanding arrow functions, and reasoning about what can appear in a particular position.

The core idea is:

> **An expression produces a value. A statement performs an action or controls execution.**

That's a useful starting point, but there's more nuance to it.

---

## 1. What is an expression?

An **expression** is a piece of JavaScript code that can be evaluated to produce a value.

For example:

```js
42
```

is an expression. Its value is `42`.

So is:

```js
10 + 20
```

Its value is `30`.

And:

```js
user.name
```

is an expression whose value is whatever `user.name` evaluates to.

Even a function call is an expression:

```js
calculatePrice()
```

The expression's value is the value returned by `calculatePrice()`.

### The important property

An expression **produces a value**.

For example:

```js
10 + 20
```

→ produces `30`

```js
user.name
```

→ produces the user's name

```js
calculatePrice()
```

→ produces the function's return value

This is why expressions can be used in places where JavaScript expects a value.

---

# 2. Expressions can be composed

One reason expressions are so powerful is that expressions can contain other expressions.

Consider:

```js
price * quantity + shippingCost
```

There are several expressions involved:

```text
price
quantity
price * quantity
shippingCost
price * quantity + shippingCost
```

The larger expression is built from smaller expressions.

This composability is a major characteristic of expressions.

For example:

```js
const total = calculatePrice(quantity) + shippingCost;
```

The right-hand side is an expression.

It contains:

- `calculatePrice(quantity)` — expression
- `shippingCost` — expression
- `calculatePrice(quantity) + shippingCost` — larger expression

The entire expression evaluates to a value that can be assigned to `total`.

---

# 3. Common kinds of expressions

You already know most of these syntactically; what's new here is recognizing them as **expressions**.

### Literals

```js
100
"Berlin"
true
null
```

These represent values, so they are expressions.

### Variable references

```js
price
```

This evaluates to the current value associated with `price`.

### Arithmetic

```js
price * quantity
```

### Comparisons

```js
price > 100
```

This produces a boolean value.

### Function calls

```js
calculateTotal()
```

### Object access

```js
user.name
```

### Object creation

```js
{ name: "Alice" }
```

### Array creation

```js
[1, 2, 3]
```

The exact categories aren't important to memorize right now. The useful question is:

> **Can this piece of code be evaluated to a value?**

If yes, you're probably looking at an expression.

---

# 4. What is a statement?

A **statement** is a unit of code that represents an instruction or action to be performed by the program.

For example:

```js
const price = 100;
```

This is a **variable declaration statement**.

It tells JavaScript to create a binding called `price` and initialize it with a value.

Another example:

```js
if (price > 100) {
    console.log("Expensive");
}
```

The `if` construct is a statement. It controls whether a particular block of code executes.

A loop such as:

```js
while (condition) {
    processItem();
}
```

is also a statement.

So statements are concerned more with **program execution and structure** than simply producing a value.

---

# 5. Statements don't necessarily produce values

This is the most useful contrast:

```text
Expression → evaluates to a value

Statement   → performs an action / controls execution
```

For example:

```js
const price = 100;
```

The purpose of this statement is to establish a variable binding.

Whereas:

```js
price * quantity
```

evaluates to a value.

This distinction is not merely theoretical. It determines **where JavaScript allows particular constructs to appear**.

---

# 6. Expression statements

Here's where things become slightly interesting.

An expression can itself be used as a statement.

This is called an **expression statement**.

For example:

```js
calculatePrice();
```

`calculatePrice()` is an expression because it produces a value.

But when you write it by itself:

```js
calculatePrice();
```

it is being used as an **expression statement**.

The expression is:

```js
calculatePrice()
```

The complete statement is:

```js
calculatePrice();
```

That's an important distinction.

---

# 7. Why would you execute an expression and ignore its value?

Consider:

```js
console.log("Processing order");
```

`console.log("Processing order")` is an expression.

It produces a value — specifically, `undefined`.

But you aren't interested in that value.

You're interested in the **side effect** of calling `console.log`.

Similarly:

```js
saveOrder(order);
```

The function call produces some value depending on the function, but perhaps you only care about the action performed by the function.

So JavaScript allows expressions to be used as standalone statements.

This is why:

```js
saveOrder(order);
```

is perfectly valid even if you don't do anything with the returned value.

---

# 8. Another example: assignment

Consider:

```js
price = 100;
```

This is an interesting case because the assignment itself is an **expression**.

The assignment evaluates to the assigned value.

Conceptually:

```js
price = 100
```

produces:

```text
100
```

But when you write:

```js
price = 100;
```

as a standalone instruction, it is an **expression statement**.

This is one reason understanding expressions is useful: things that look like "statements" can sometimes actually contain expressions with their own values.

We'll explore assignment expressions and operators more deeply in Chapter 3.

---

# 9. Statements can contain expressions

This is probably the most important structural relationship to understand.

Statements and expressions aren't two completely separate worlds.

A statement can contain expressions.

For example:

```js
const total = price * quantity;
```

The entire thing is a **declaration statement**.

Inside it:

```js
price * quantity
```

is an expression.

Another example:

```js
if (total > 100) {
    processLargeOrder();
}
```

The `if` is a statement.

Inside it:

```js
total > 100
```

is an expression.

And:

```js
processLargeOrder()
```

is an expression being used as an expression statement inside the block.

So:

```text
Statement
│
├── expressions
├── other syntax
└── execution structure
```

---

# 10. Why the distinction matters

The distinction becomes particularly important when you ask:

> **"What can I put here?"**

Some syntactic positions expect an **expression**.

For example, after an assignment:

```js
const total = ???;
```

the `???` must be something that can serve as an expression.

So these work:

```js
const total = price * quantity;
```

```js
const total = calculateTotal();
```

But a statement such as an `if` cannot simply be inserted there:

```js
const total = if (condition) {
    ...
};
```

That's invalid JavaScript.

Understanding the distinction lets you reason about syntax instead of memorizing isolated rules.

---

# 11. The distinction becomes especially important with arrow functions

We'll properly cover arrow functions in Chapter 6, so I won't teach their syntax in detail yet.

But there is one reason they're worth mentioning here.

These two forms behave differently:

```js
const calculate = () => 42;
```

and:

```js
const calculate = () => {
    return 42;
};
```

The first has an **expression body**.

The second has a **block body containing statements**.

That's one of the reasons arrow functions have the behavior they do.

We'll revisit this properly when we reach functions.

I don't want to pull the rest of arrow-function semantics into this lesson prematurely.

---

# 12. Semicolons

You'll often see statements written with semicolons:

```js
const price = 100;
saveOrder(order);
```

Semicolons are related to JavaScript's rules for terminating statements.

However, JavaScript has a feature called **Automatic Semicolon Insertion (ASI)**, which means semicolons can sometimes be omitted.

For example:

```js
const price = 100
saveOrder(order)
```

is generally valid JavaScript.

However, this does **not** mean:

> "Semicolons are unnecessary."

ASI has specific rules and can produce surprising behavior in certain situations.

We'll discuss ASI when it becomes relevant to JavaScript syntax and coding conventions rather than spending this lesson on it.

For now:

> Treat semicolons as statement terminators in the style we'll use throughout this course.

---

# 13. A useful way to identify them

When reading JavaScript, ask two different questions.

### Question 1

> **Does this code evaluate to a value?**

If yes, you're looking at an expression.

### Question 2

> **Is this code an instruction/control structure that participates in program execution?**

If yes, you're probably looking at a statement.

Then remember that they can overlap:

```js
saveOrder(order);
```

contains an expression:

```js
saveOrder(order)
```

which is being used as an expression statement.

---

# Interview questions

### 1. What is an expression?

> An expression is JavaScript syntax that can be evaluated to produce a value.

---

### 2. What is a statement?

> A statement is a unit of JavaScript code that represents an action or controls program execution.

---

### 3. What's the difference between an expression and a statement?

A strong concise answer:

> An expression evaluates to a value, whereas a statement represents an instruction or control-flow construct. Statements can contain expressions, and expressions can sometimes be used as expression statements.

---

### 4. What is an expression statement?

> It's a statement consisting of an expression whose resulting value is not necessarily used.

For example:

```js
saveOrder(order);
```

The function call is an expression, and the complete line is an expression statement.

---

### 5. Is assignment an expression or a statement?

This is a nice interview trap.

The assignment operation itself is an **expression**:

```js
price = 100
```

When used on its own:

```js
price = 100;
```

it forms an **expression statement**.

---

### 6. Can a statement produce a value?

Don't make the overly broad claim that statements "never have values."

The more useful distinction is:

> **Expressions have values; statements are syntactic constructs for execution.**

Some language constructs have specification-level completion values, but that's a different concept from an expression's value and isn't something you need to conflate with it at this stage.

---

# The mental model

Keep this:

```text
                  JavaScript code
                        │
             ┌──────────┴──────────┐
             ↓                     ↓
        Expressions            Statements
             │                     │
       produce values       perform/control execution
             │                     │
             └───────┬─────────────┘
                     ↓
             statements can
             contain expressions
                     │
                     ↓
             expression can
             be a statement
             (expression statement)
```

The most important example to remember is:

```js
saveOrder(order);
```

**Expression:**

```js
saveOrder(order)
```

**Statement:**

```js
saveOrder(order);
```

Same piece of code, viewed at two different syntactic levels.

---

## Exercise

For each of the following, determine whether it is an **expression**, a **statement**, or an **expression statement**. If it's a statement containing an expression, identify the expression inside it.

```js
42
```

```js
const total = price * quantity;
```

```js
price * quantity;
```

```js
if (total > 100) {
    processOrder();
}
```

```js
calculateTotal();
```

Don't just classify them — **explain why** for each one.

## Answers
Almost — **3 out of 5**, with one important nuance.

1. `42` → **Expression** 
   It evaluates to the value `42`.

2. `const total = price * quantity;` → **Statement**  
   More specifically, a variable declaration statement.  
   `price * quantity` inside it is an **expression**.

3. `price * quantity;` → **Expression statement**  
   The expression itself is:
   ```js
   price * quantity
   ```
   But because it's being used as a standalone statement, the complete construct is an **expression statement**.

4. `if (total > 100) { processOrder(); }` → **Statement**  
   More specifically, an `if` statement.  
   `total > 100` is an expression, and `processOrder()` is an expression used as an expression statement inside the block.

5. `calculateTotal();` → **Expression statement**  

### The distinction I want you to internalize

Don't classify based on whether something *looks* like an expression.

Ask:

> **What is the syntactic role of this code here?**

For #3:

```js
price * quantity;
```

The part:

```js
price * quantity
```

**is an expression.**

But the complete construct:

```js
price * quantity;
```

**is an expression statement.**

That's exactly the kind of distinction that will help later when we encounter more complicated JavaScript syntax.

You're good on **1.5**. When you're ready, say **next** and we'll do **1.6 Strict Mode**.