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

Refer to plan.md for the course curriculum. 

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
