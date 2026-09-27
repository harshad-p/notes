# AGENTS.md

## Rules
1. When the user asks to go next, look for the next chapter in line as per plan.md
1. Create the appropriate file/folder as the rules below

## Chapters

A named chapter is a numbered heading in `plan.md`, such as `### 1.1 Types and Type System`. Bullets under a part are topics inside a chapter, not separate chapters. A part with no numbered headings has one nameless chapter. Its file is the chapter number only.

For every chapter:

1. Create a numbered part folder when it does not already exist (`part-00-interview-foundations/`, `part-01-c#-core/`, and so on). Take the folder name from `plan.md`: lowercase, spaces as hyphens, with the part number.
2. Inside that part folder, create an empty file per chapter. Take the file name from the chapter name in `plan.md`: lowercase, spaces as hyphens, prefixed with the chapter number. Example: `chapter-1.1-types-and-type-system.md`. A nameless chapter has no name segment.
3. Keep numbering consistent with `plan.md`.

Example:

```text
part-00-interview-foundations/
    chapter-0.1.md
part-01-c#-core/
    chapter-1.1-types-and-type-system.md
```

## Progress

Maintain `progress.md`.

progress.md must use Markdown checkboxes to track progress.

Mark a chapter complete when the user asks.
