# Spec Guidelines

> This project uses the [Specture System](https://github.com/specture-system/specture), and this document outlines how it works. As the Specture System is improved and updated, this file will also be updated.

## Overview

Specs are design documents that describe planned changes to the system. They serve as a blueprint for implementation and a historical record of design decisions. They are continually improved during the design and implementation of a change, but left static after the change is complete.

This system is inspired by the [PEP](https://peps.python.org/) (Python Enhancement Proposal) and [BIP](https://github.com/bitcoin/bips) (Bitcoin Improvement Proposal) systems, adapting their formal proposal processes for general software development.

## Spec Scope

**Specs are for planned changes**: new features, major refactors, redesigns, and tooling improvements. Use the issue tracker for bugs: cases where the software doesn't match what's described in the specs.

In traditional project management, this scope would be split between epics (grouping work) and issues (individual tasks). Specture consolidates this: a spec includes the design rationale (why), the decisions (what and how), and the implementation task list all in one document. This keeps related information together and makes it easier for the reader to understand the full context. An individual spec can vary wildly in size, from as large as a traditional epic to as small as a tiny UI improvement.

### When to split or combine specs

A spec should be a coherent unit of work. Split specs when changes are independent: they can be implemented and deployed separately, have different design rationales, or solve different problems. Keep specs together when they share design decisions or task dependencies.

Examples:

- One spec: "Add authentication system" (related design choices, shared infrastructure)
- Two specs: "Add authentication" and "Redesign dashboard UI" (independent changes)
- One minimal spec: "Rename demo app consistently" (small but still a planned change)

## Spec File Structure

Each spec file should be a Markdown document with a numeric prefix in the `specs/` directory, for example `specs/000-mvp.md`.

**Example format:**

```markdown
---
status: draft
author: Your Name
creation_date: 2025-12-18
---

# Feature or Change Name

Description as paragraphs and/or bulleted list.

## Task List

### Foundation

- [ ] Task 1
- [ ] Task 2
- [ ] Task 3

### Core Implementation

- [ ] Task 1
- [ ] Task 2

### Polish and Documentation

- [ ] Task 1
- [ ] Task 2
```

### Frontmatter

Each spec should begin with YAML frontmatter containing metadata:

#### Required Fields

- `status`
  - `draft` - Spec is in the process of being written and improved upon
  - `approved` - Spec has been approved and is awaiting implementation
  - `in-progress` - Implementation is underway
  - `completed` - All tasks have been completed
  - `rejected` - Spec was reviewed but rejected

#### Optional Fields

- `author` - Person(s) who proposed or wrote the spec
- `creation_date` - Date the spec was created (YYYY-MM-DD format)
- `approved_by` - Person(s) who approved the spec
- `approval_date` - Date the spec was approved (YYYY-MM-DD format)

### Title

The spec should start with a clear, descriptive H1 heading that summarizes what is being proposed.

### Description

An overview of the proposed change. It might be a couple sentences for a small change or dozens of paragraphs for a large change.

A few things to consider including:

- What is being proposed
- Why it's needed
- What problem it solves
- High-level approach

Feel free to use paragraph form or bulleted list form, whichever better matches the requirements of clearly describing the proposal.

For large descriptions, please add separate sections with their own headers, for example: `## Ideas`, `## Goals`, `## Benefits`.

### Design Decisions

An optional section to document the design process. For each major decision made during the design process, include the options considered along with the pros and cons of each. This creates a valuable historical record of why certain choices were made.

The number of decisions mentioned will most likely be proportional to the size of the change being discussed. Include as many decision points as needed. There is no obligation to include them for small specs or changes with trivial design choices.

### Task List

A detailed breakdown of implementation tasks using markdown checklists. Split into logical sections if needed.

#### Task List Best Practices

- Make tasks specific and actionable
- Order tasks logically (dependencies first)
- Group related tasks into sections
- Include testing, documentation, and deployment tasks where applicable
- Keep individual tasks reasonably sized and self-contained
- Avoid getting into nitty-gritty implementation details

## File Naming

Use descriptive, kebab-case filenames with a numeric prefix:

- `000-mvp.md`
- `001-add-authentication-system.md`
- `013-refactor-database-layer.md`
- `314-redesign-api-endpoints.md`

## Workflow

### Draft

1. Create a new spec file with frontmatter status set to `draft`
2. Write the spec with required sections: Title, Description, and Task List. Include Design Decisions if there are meaningful choices to document
3. Keep the spec updated as you refine the proposal
4. Open pull requests to share and refine the spec with the team
5. The spec may go through multiple iterations while still in draft status as it evolves

#### Tips

- **Be clear, not clever**: Write for future readers who may not have context
- **Document alternatives**: Even rejected options are valuable to record
- **Link to discussions**: Reference pull request comments, issues, or other specs
- **Focus on "why"**: The code shows "how", the spec should explain "why"

### Approved

It is up to the project maintainers to determine when they are ready to merge a spec with the status set to `approved`. The requirements for what defines an approved spec will vary by project.

1. Once the team agrees the spec is ready for implementation, update status to `approved` in a pull request
2. Merge the pull request with the updated status
3. Implementation can now begin

### In-Progress

1. Implementation begins once a spec is `approved`
2. For each task, follow this workflow:
   - Complete a single task from the task list
   - Update the spec file by changing `- [ ]` to `- [x]` for that task
   - Commit both the implementation and spec update together with a conventional commit message (e.g., `feat: implement feature X`)
   - Push the changes
3. Keep the spec updated as implementation reveals new details
   - Add tasks or task sections as needed
   - Remove or update existing tasks as needed
4. When all tasks are completed, move to `completed` status

This workflow keeps the spec file as a living document that tracks implementation progress, with each task corresponding to one commit.

### Completed

1. Mark the spec status as `completed` when all tasks in its task list are checked off
2. If there is still remaining work to be done for a change, be sure that the task list reflects this reality and keep the status as `in-progress`

### Rejected

For most rejected specs, the project maintainers will not merge the spec. There's no reason to merge a proposal you don't intend to implement.

However, some rejected specs may be useful to merge with the status set to `rejected`. These can act as a historical record of what was considered and rejected. The important thing in these cases is to clearly document **why** a proposal was rejected.

## Precedence System

The requirements listed in one spec may become outdated with time. Once a spec has status `completed`, it should be treated as a historical record and not retroactively updated. It is a bad idea to try to go back and update completed specs. This would be tedious and error-prone. Inevitably, something would be missed.

The exception is to fix overlooked mistakes: typos, errors in documentation, or factual inaccuracies should be corrected.

Instead, we rely on a simple precedence system. The numeric prefix at the beginning of each spec defines its precedence. The higher the number is, the higher the precedence. If two specs contradict each other on any particular point, the higher numbered spec takes precedence.

In general, the numbers should be incremented over time with each new spec added to the project.

Some tricky situations might arise where it becomes necessary to number specs non-incrementally, especially when a team is working on drafting multiple specs at once. A project's spec number assignment system should be optimized for the needs of that project. Overall, no matter what scheme you determine for assigning numbers to each new spec, stick to the rule that higher number means higher precedence.
