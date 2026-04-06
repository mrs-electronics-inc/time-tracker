# Spec File Format Reference

Detailed format specification for Specture spec files. Loaded on demand when creating or editing specs.

## File Location and Naming

Specs live in the `specs/` tree as directories containing `SPEC.md` files. Nested specs are created as child directories under their parent spec directory.

```
specs/mvp/SPEC.md
specs/mvp/backend/SPEC.md
specs/refactor-database-layer/SPEC.md
```

Spec numbers are stored in the YAML frontmatter `number` field (not in the directory name). The number is a local identifier within the spec's parent scope, and the full dotted reference is derived from the directory tree. `specture new` auto-assigns the next available number within the selected scope.

## Complete Example

```markdown
---
number: 0
status: draft
author: Your Name
creation_date: 2025-12-18
---

# Feature or Change Name

Description as paragraphs and/or bulleted list.

## Design Decisions

### Decision Title

- Chosen: Option A
  - Pro 1
  - Pro 2
- Considered: Option B
  - Pro 1
  - Con 1

```

## Frontmatter

YAML frontmatter between `---` delimiters at the top of the file.

### Required Fields

| Field    | Values                                                         |
| -------- | -------------------------------------------------------------- |
| `number` | Non-negative integer (0, 1, 2, ...). Auto-assigned by `specture new`. |
| `status` | `draft`, `approved`, `in-progress`, `completed`, or `rejected` |

### Optional Fields

| Field           | Format       | Description                          |
| --------------- | ------------ | ------------------------------------ |
| `author`        | Free text    | Person(s) who proposed/wrote the spec |
| `creation_date` | `YYYY-MM-DD` | Date the spec was created            |
| `approved_by`   | Free text    | Person(s) who approved the spec      |
| `approval_date` | `YYYY-MM-DD` | Date the spec was approved           |

### Status Values

- **`draft`** — Being written and refined. May go through multiple iterations.
- **`approved`** — Team has agreed on the design; ready for implementation.
- **`in-progress`** — Implementation underway.
- **`completed`** — All planned work is done and goals are achieved.
- **`rejected`** — Reviewed and rejected. Document **why** if merging a rejected spec.

## Required Sections

Do not number markdown headings in spec files. Use plain titles like `## Design Decisions` and `### Foundation`, not `## 1. Design Decisions` or `### 2.1 Foundation`.

### Title (H1)

A clear, descriptive `# Heading` summarizing what is being proposed. Must be the first heading in the document.

### Description

Overview of the proposed change, immediately after the title. Can be paragraphs or bulleted lists. Consider including:

- What is being proposed
- Why it's needed
- What problem it solves
- High-level approach

For large descriptions, use additional `##` sections (e.g., `## Ideas`, `## Goals`, `## Benefits`).

Keep implementation progress out of the spec content. Use the spec for rationale and design decisions rather than a step-by-step task checklist.

When referencing another spec, always use an inline markdown link with the correct repo-root-relative path to that file (for example, `[Status command](specs/002-status-command/SPEC.md)`).

### Goals

A list of specific goals that this spec aims to achieve. May include sublists.

### Design Decisions

Document major design choices with options considered and their trade-offs:

```markdown
## Design Decisions

### Choice Title

- Chosen: Option A
  - Reason it was selected
  - Another advantage
- Considered: Option B
  - Advantage of this option
  - Why it was not selected
```

Include as many decision points as needed. No obligation for small or trivial specs.

## Scope Guidelines

**Specs are for planned changes**: features, refactors, redesigns, tooling improvements.

**Not for bugs**: If software doesn't match what specs describe, use the issue tracker.

**When to split specs**: Changes are independent, can be deployed separately, solve different problems.

**When to combine**: Changes share design decisions or have task dependencies.

## Validation

Run `specture validate` to check specs.
