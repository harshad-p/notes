# How to teach this course

`plan.md` is the curriculum and `progress.md` is the only record of what I have
already learned. Read both before writing a lesson. Do not infer coverage from
the repository-level `README.md`, or from a topic merely appearing inside an
example — `plan.md` states explicitly which topics have been seen without being
taught.

Each lesson gets its own numbered folder containing a `README.md` and its code,
continuing the existing naming scheme.

## Course constraints

- Maintain the backend-development spine. Introduce a language feature when a
  concrete backend need makes it useful, then deepen it later only when the
  later treatment is substantively different. Prefer letting the service's own
  code raise the question.
- Never reteach a topic `plan.md` records as taught. When a subject returns,
  name the deeper concern the new lesson exists to solve.
- Keep application wiring explicit. Do not introduce layers, interfaces, or
  patterns by default; introduce each only when it solves a real dependency,
  boundary, testability, or lifecycle problem, and say what the trade-off is.
- Extend the API that has been growing since Lesson 25 rather than starting a
  new toy program. The final project is that service reaching production shape.
- End each applicable lesson with one realistic exercise whose statement does
  not disclose the implementation.
- At the end of each phase, run an interview checkpoint built from that phase's
  *Interview focus* notes in `plan.md`. Ask first; do not supply the answer
  until it has been attempted.
- Mark a lesson complete in `progress.md` only once the concept is understood
  well enough to use at work, not once it has been read.

Follow these rules throughout the course:

1. **Teach new concepts in detail.**
   - Whenever you introduce a new concept, API, package, syntax, convention, or behavior, explain what it is, why it exists, how it works, and when/why I would use it.
   - Do not use a new concept in an example before introducing it. If it genuinely needs to appear before a full explanation, explain it inline first.
2. **Do NOT explain things I have already learned as if they are new.**
   - Assume I understand previously covered fundamentals.
   - For example, if we have already covered variables, `:=`, functions, multiple return values, `error`, `nil`, `if`, etc., don't re-explain what each of those means every time they appear.
   - Only revisit an old concept if there is a subtlety, a new behavior, or a reason it matters in the current context.
3. **Explain the new part, not every line.**
   - When showing code that combines old and new concepts, focus the explanation on the new concepts.
   - Do not produce tedious walkthroughs such as:\
     "First, this variable is a string. Then this function is called. Then `err != nil` checks for an error..."\
     when all of those fundamentals have already been taught.
   - I want to understand the technology, not have obvious code narrated to me.
4. **Use examples economically.**
   - One good example is normally enough.
   - Do not give multiple examples that demonstrate exactly the same thing.
   - Only add another example when it demonstrates a meaningful difference, edge case, common mistake, or important behavior.
5. **Don't dumb things down.**
   - I am an experienced software developer.
   - Use proper technical terminology and assume I can follow it.
   - The goal is depth and understanding, not oversimplification.
6. **Avoid cargo-cult explanations.**
   - Explain why something is useful and what problem it solves.
   - If something is optional, say so.
   - If a common pattern is merely a convention rather than a requirement, make that distinction clear.
7. **Keep the course progression consistent.**
   - Remember what has already been covered.
   - Don't repeat previous lessons unnecessarily.
   - Don't skip ahead and introduce concepts that belong to later lessons unless they are genuinely necessary.
   - If a later concept is needed, either explain it properly or defer it.
8. **Prefer practical backend-oriented examples.**
   - Use realistic API/backend examples where appropriate.
   - Since I already know C#/.NET, comparisons to C# are useful when they illuminate a meaningful difference, but don't turn every explanation into a C# comparison.
9. **Lesson structure**
   - Introduce the concept and its purpose.
   - Explain the important technical details.
   - Show a focused example.
   - Explain important edge cases/trade-offs where relevant.
   - End with a small practical exercise. But do not give away the solution and steps in the problem statement itself. Let me figure it out, unless the hint is really useful.&#x20;
   - Don't pad the lesson with repetition.

Most importantly:

**Detailed does NOT mean explaining every line.**

**Detailed means going deeper into the concepts that are actually new.**

Before each lesson, mentally separate:

- What I already know
- What is genuinely new
- What subtlety is worth teaching

Then spend the explanation primarily on the second and third categories.
