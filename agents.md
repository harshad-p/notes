# Curriculum Management

You are managing a programming curriculum that contains both theoretical material and practical coding.

## 1. Programming Language

This is a **Python** programming curriculum.

All practical source files must use Python unless explicitly instructed otherwise.

If this file is being reused for another programming language, this section must be changed explicitly before creating practical files.

Do NOT guess the programming language from the chapter topics.

---

## 2. Source of Truth

The curriculum structure is defined by `plan.md`.

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

## 3. Parts / Modules

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

00-python-foundations/
01-python-core/
02-building-agents/

---

## 4. Chapters

Every actual chapter is represented by a folder.

### Named chapter

If `plan.md` explicitly gives the chapter a name, use:

`<chapter-number>-<chapter-name>/`

Example:

01-python-interpreter/  
02-first-python-program/  
03-variables/  

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

Create ONE chapter folder:

0.1-interview-foundations/

or, if no useful name can be derived:

chapter-0.1/

If a meaningful chapter name can be derived from the topics, prefer:

0.1-interview-foundations/

Do not create:

0.1-topic-one/  
0.2-topic-two/  
0.3-topic-three/  

Only create multiple chapters when `plan.md` explicitly defines multiple chapters.

---

## 5. Do Not Pre-create Curriculum Structure

Do not create folders or files for future Parts, Modules, or Chapters.

Only create the filesystem structure required for:

- the current chapter, when explicitly requested, or
- the next chapter when the user says `next`.

The existence of a Part or Chapter in `plan.md` does not mean its folder or files should be created yet.

---

## 6. Chapter Contents

Every chapter folder must contain:

README.md

`README.md` contains the theoretical notes for that chapter.

Create a practical source file only when the chapter requires hands-on coding.

For this Python curriculum:

practice.py

Example:

01-python-interpreter/  
├── README.md  
└── practice.py  

If the chapter is purely theoretical:

01-something/  
└── README.md

Do not create an unnecessary `practice.py`.

Determine whether practical code is appropriate from the chapter topics in `plan.md`.

---

## 7. Go Curriculum

If the Programming Language section is changed to Go, use the following structure instead of `practice.py`.

For example:

0.3-http-server/  
├── README.md  
├── go.mod  
└── main.go  

The chapter folder name must be used as the Go module name.

Example:

module 0.3-http-server

Create `main.go` as an empty source file.

Do NOT create `go.run`.

Do NOT use `practice.go`.

The default Go source file is always:

main.go

---

## 8. Creating the Next Chapter

When the user says `next`:

1. Read `plan.md`.
2. Read `progress.md`.
3. Identify the current chapter.
4. Mark the current chapter as completed in `progress.md`.
5. Identify the next actual chapter according to `plan.md`.
6. Create its Part folder if it does not exist.
7. Create the next chapter folder.
8. Create `README.md`.
9. Create the appropriate practical source file if required.
10. For Go, create `go.mod` and an empty `main.go`.
11. Do not mark the new chapter as completed.

Never skip chapters.

Never turn topic bullet points into chapters.

Never create future chapters beyond the next one.

---

## 9. progress.md

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

Creating a chapter folder does not mean the chapter is completed.

Only mark a chapter `[x]` when the user says `next` or otherwise explicitly indicates that they have completed it.

---

## 10. Existing Files

Never overwrite existing files or their contents.

If a Part folder already exists, reuse it.

If a chapter folder already exists, reuse it.

If `README.md`, `practice.py`, `go.mod`, or `main.go` already exists, leave it unchanged.

Only create missing files.

---

## 11. Content

When preparing a new chapter, create the required files but do not write lesson content into them unless explicitly instructed.

Your responsibility is filesystem and progress management, not teaching.