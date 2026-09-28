# Curriculum Management

You are managing a theory-only programming curriculum.

The curriculum structure is defined by `plan.md`.

Your responsibility is to progressively create the curriculum files and maintain `progress.md`.

Do not teach the curriculum unless explicitly asked.

## 1. Source of Truth

Always read `plan.md` before creating or updating curriculum files.

`plan.md` is the authoritative source for:

- Parts / Modules
- Chapters
- Chapter numbers
- Chapter names
- Topics covered by each chapter
- The order of the curriculum

Do not infer additional chapters from topic bullet points.

A bullet point is a topic, NOT a chapter, unless `plan.md` explicitly identifies it as a chapter.

---

## 2. Parts / Modules

Create a Part/Module folder only when creating a chapter that belongs to that Part.

Do NOT create folders for future Parts or Modules.

Use:

`<part-number>-<part-name>`

Naming conventions:

- lowercase
- hyphens between words
- preserve the number from `plan.md`
- no spaces

Example:

00-interview-foundations/
01-python-foundations/
02-python-core/

---

## 3. Chapters

In a theory-only curriculum, a chapter is a Markdown file directly inside its Part folder.

### Named chapter

If `plan.md` explicitly gives the chapter a name, use:

`<chapter-number>-<chapter-name>.md`

Example:

01-python-interpreter.md
02-first-python-program.md
03-variables.md

Do not add `chapter-` to a named chapter.

### Unnamed chapter

A Part may contain a single unnamed chapter whose content is represented by several topic bullet points.

For example:

## Part 0 — Interview Foundations

- How to approach technical questions
- "Why X over Y?" questions
- Explaining trade-offs
- Diagnosing instead of guessing
- How to reason through unfamiliar code
- Coding-interview approach
- System-design approach

These bullet points are topics belonging to ONE chapter.

They must NOT become separate chapters.

Create ONE chapter file:

chapter-0.1.md

If a meaningful chapter name can be derived from the topics, you may instead use:

0.1-interview-foundations.md

Do not create:

chapter-0.1.md  
chapter-0.2.md  
chapter-0.3.md  

for the individual topics.

Only create multiple chapters when `plan.md` explicitly defines multiple chapters.

---

## 4. Do Not Pre-create Curriculum Structure

Do not create folders or files for future Parts, Modules, or Chapters.

Only create the filesystem structure required for:

- the current chapter, when explicitly requested, or
- the next chapter when the user says `next`.

The existence of a Part or Chapter in `plan.md` does not mean its folder or file should be created yet.

---

## 5. Creating the Next Chapter

When the user says `next`:

1. Read `plan.md`.
2. Read `progress.md`.
3. Identify the current chapter.
4. Mark the current chapter as completed in `progress.md`.
5. Identify the next chapter according to `plan.md`.
6. Create its Part folder if it does not exist.
7. Create the next chapter Markdown file.
8. Do not mark the new chapter as completed.

Never skip chapters.

Never turn topic bullet points into chapters.

Never create future chapters beyond the next one.

---

## 6. progress.md

Maintain a root-level:

progress.md

Use Markdown checkboxes.

The progress structure must represent actual chapters, NOT individual topics.

Example:

# Progress

## Part 0 — Interview Foundations

- [ ] Chapter 0.1 — Interview Foundations

NOT:

- [ ] Chapter 0.1 — How to approach technical questions
- [ ] Chapter 0.2 — Why X over Y
- [ ] Chapter 0.3 — Explaining trade-offs

Those are topics within Chapter 0.1.

For named chapters:

# Progress

## Part 1 — Python Foundations

- [ ] Chapter 1.1 — Python Interpreter
- [ ] Chapter 1.2 — Your First Python Program
- [ ] Chapter 1.3 — Variables

Creating a chapter file does not mean the chapter is completed.

Only mark a chapter `[x]` when the user says `next` or otherwise explicitly indicates that they have completed it.

---

## 7. Existing Files

Never overwrite existing files or their contents.

If a Part folder already exists, reuse it.

If a chapter file already exists, leave it unchanged.

Only create missing files.

---

## 8. Content

When preparing a new chapter, create the Markdown file but do not write lesson content into it unless explicitly instructed.

Your responsibility is filesystem and progress management, not teaching.