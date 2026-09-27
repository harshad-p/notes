# AGENTS.md

## Rules
1. When the user asks to go next, look for the next chapter in line as per plan.md
1. Create the appropraite file/folder as the rules below

## Chapters

For every chapter:

1. Create a numbered part folder when it does not already exist (`part-00-interview-foundations/`, `part-01-c#-core/`, and so on). Take the folder name from `plan.md`: lowercase, spaces as hyphens, with the part number. 
1. Inside that part folder, create an empty file per chapter. Take the file name from the chapter name in `plan.md`: lowercase, spaces as hyphens, prefixed with the chapter number. Example: `chapter-0.1-how-to-approach-technical-questions.md`.
1. Keep numbering consistent with `plan.md`.

Example:

```text
part-00-interview-foundations/
    chapter-0.1-how-to-approach-technical-questions.md
part-01-c#-core/
    chapter-1.1-value-types-vs-reference-types.md
```

## Progress

Maintain `progress.md`.

progress.md must use Markdown checkboxes to track progress.

Mark it as complete when the user asks. 