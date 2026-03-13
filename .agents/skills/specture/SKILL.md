---
name: specture
description: Follow the Specture System for spec-driven development. Use when creating, implementing, or managing specs.
---

# Specture System

Specture is a spec-driven development system. Specs are design documents in the `specs/` directory that describe planned changes — features, refactors, redesigns, tooling improvements. Each spec contains the design rationale, decisions, and an implementation task list.

Spec numbers are stored in the YAML frontmatter `number` field. New specs use slug-only filenames (e.g., `my-feature.md`). Older specs may retain `NNN-slug.md` filenames — both naming patterns are valid.

## Design Workflow

When designing a spec (writing or refining a `draft` spec):

1. Create a new branch for the spec design work (e.g., `spec/my-feature`)
2. Write or refine the spec content — description, design decisions, task list
3. Commit and push the spec changes
4. Open a PR for review and discussion

## Implementation Workflow

When starting implementation of an `approved` spec:

1. **Update the spec status to `in-progress`** by editing the frontmatter `status` field
2. **Create a new branch** for the implementation work (e.g., `impl/my-feature`)
3. Commit the status change as the first commit on the branch

Then follow this loop:

1. Run `specture status` to see the current spec and next task
2. Complete one or more tasks from the task list
3. Edit the spec file: change `- [ ]` to `- [x]` for every task completed in this commit
4. Stage both the implementation files **and** the spec file update
5. Commit together with a conventional commit message (e.g., `feat: implement feature X`)
6. Push the changes
7. Repeat from step 1

**Critical rules:**

- Every commit that completes a task MUST include the spec file checkbox update alongside the implementation changes. Never commit implementation without the corresponding `- [x]` update. This is the most important rule.
- If a single commit completes multiple tasks, check off all of them in that same commit. Do NOT make separate empty commits just to check off tasks that were already implemented.
- Do NOT edit spec design decisions or descriptions without explicit user permission. You may only mark tasks complete and add/remove tasks during implementation.
- When editing a spec, keep the design decisions section and task list in sync. If a description is updated, update all corresponding task descriptions to match, and vice versa.
- When all tasks are checked off, update the frontmatter `status` to `completed`.

## CLI Commands

Always use non-interactive flags. Interactive mode will hang waiting for input.

### specture list and specture status

Use `list` to see all specs at a glance, then `status` to drill into a specific one.

**`specture list`** — overview of all specs (number, status, progress, name).

```bash
specture list                            # All specs
specture list --status in-progress       # Filter by status
specture list --status draft,approved    # Multiple statuses
specture list -f json                    # JSON output with full metadata
```

Aliases: `list`, `ls`

**`specture status`** — detailed view of one spec, including tasks and current task.

```bash
specture status                          # Current in-progress spec
specture status --spec 3                 # Specific spec by number
specture status -f json                  # JSON output
```

Typical workflow: run `specture list` to find the spec you need, then `specture status --spec N` to see its tasks and progress.

### specture new

Create a new spec file with automatic numbering and branch.

```bash
# Non-interactive: provide title via flag (required for agents)
specture new --title "Feature name"

# Pipe full body content
cat spec-body.md | specture new --title "Feature name"

# Skip branch creation
specture new --title "Feature name" --no-branch

# Skip opening editor
specture new --title "Feature name" --no-editor

# Preview without creating anything
specture new --title "Feature name" --dry-run
```

Aliases: `new`, `n`, `add`, `a`

### specture validate

Validate that specs follow the Specture System format.

```bash
# Validate all specs
specture validate

# Validate a specific spec by number
specture validate --spec 3
specture validate -s 42
```

Checks: valid frontmatter (number and status), no duplicate numbers, description present, task list present. Warns on number/filename mismatch.

Aliases: `validate`, `v`

### specture setup

Initialize or update the Specture System in a repository.

```bash
# Non-interactive setup
specture setup --yes

# Preview without changes
specture setup --dry-run

# Force AGENTS.md update prompt
specture setup --update-agents --yes
```

Aliases: `setup`, `update`, `u`

### specture rename

Rename a spec file and update all markdown links in the specs directory.

```bash
# Rename spec 3 to status-command.md
specture rename --spec 3 status-command

# Preview changes
specture rename --spec 3 status-command --dry-run
```

## Spec Status Workflow

Specs move through these statuses:

1. **draft** — Being written and refined
2. **approved** — Ready for implementation
3. **in-progress** — Implementation underway (tasks being checked off)
4. **completed** — All tasks done
5. **rejected** — Reviewed and rejected

If a spec has no explicit `status` in frontmatter, it is inferred from tasks:
- No task list or no complete tasks → `draft`
- Mix of complete and incomplete → `in-progress`
- All tasks complete → `completed`

## Commit Messages

Use conventional commits:

- `feat:` — new features
- `fix:` — bug fixes
- `refactor:` — code restructuring
- `docs:` — documentation
- `test:` — test changes

## Precedence

Higher-numbered specs take precedence over lower-numbered ones when they conflict. Completed specs are historical records — do not retroactively update them (except to fix typos or factual errors).

## Spec Format Reference

For detailed spec file format (frontmatter fields, sections, naming conventions), see [references/spec-format.md](references/spec-format.md).
