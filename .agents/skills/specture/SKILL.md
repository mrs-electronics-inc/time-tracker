---
name: specture
description: Follow the Specture System for spec-driven development. Use when creating, implementing, or managing specs.
---

# Specture System

Specture is a spec-driven development system. Specs are design documents in the `specs/` tree that live in `SPEC.md` files and may nest to any number of levels. Each spec contains metadata, goals, and design decisions.

Spec numbers are stored in the YAML frontmatter `number` field. The number is a local identifier within its parent scope, and the full dotted reference is derived from the directory tree.

## Design Workflow

When designing a spec (writing or refining a `draft` spec):

1. Create a new branch for the spec design work (e.g., `spec/my-feature`)
2. Write or refine the spec content — description, goals, and design decisions
3. Commit and push the spec changes
4. Open a PR for review and discussion

## Implementation Workflow

When starting implementation of an `approved` spec:

1. **Update the spec status to `in-progress`** by editing the frontmatter `status` field
2. **Create a new branch** for the implementation work (e.g., `impl/my-feature`)
3. Commit the status change as the first commit on the branch

Then follow this loop:

1. Read the spec
2. Analyze the codebase
3. Implement one small chunk of the spec
4. Update the spec docs and tests as needed
5. Commit the implementation changes
6. Push the changes
7. Repeat from step 1

**Critical rules:**

- Use plain-language markdown headings in specs. Do **not** number section headers (`## 1. Design Decisions`, `### 2.1 Foundation`, etc.).
- Any cross-spec mention MUST use an inline markdown link to the other spec file with the correct repo-root-relative path (for example, `[Status command](specs/002-status-command/SPEC.md)`).
- Do NOT edit spec design decisions or descriptions without explicit user permission.
- When all work is complete, update the frontmatter `status` to `completed`.

## CLI Commands

Always use non-interactive flags. Interactive mode will hang waiting for input.
To find specs, use `specture ls`/`specture list` — do **not** scan with `grep`, `find`, or manual filename searching.

### specture list

Use `list` to see all specs at a glance. To inspect a specific spec in detail, read its `SPEC.md` file directly.

**`specture list`** — overview of all specs (ref, name, status, path).

```bash
specture list                            # All specs
specture list --status in-progress       # Filter by status
specture list --status draft,approved    # Multiple statuses
specture list -f json                    # JSON output with ref, name, status, and path
```

Aliases: `list`, `ls`

### specture new

Create a new spec directory with automatic numbering and branch.

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

Checks: valid frontmatter (number and status), no duplicate numbers, description present. Warns on number/filename mismatch.

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

Rename a spec directory and update all markdown links in the specs tree.

```bash
# Rename spec 3 to status-command
specture rename --spec 3 status-command

# Preview changes
specture rename --spec 3 status-command --dry-run
```

## Spec Status Workflow

Specs move through these statuses:

1. **draft** — Being written and refined
2. **approved** — Ready for implementation
3. **in-progress** — Implementation underway
4. **completed** — All planned work is done and goals are achieved
5. **rejected** — Reviewed and rejected

If a spec has no explicit `status` in frontmatter, it is treated as `draft` until the status is set explicitly.

## Commit Messages

Use conventional commits:

- `feat:` — new features
- `fix:` — bug fixes
- `refactor:` — code restructuring
- `docs:` — documentation
- `test:` — test changes

## Precedence

Completed specs are historical records — do not retroactively update them (except to fix typos or factual errors).

## Spec Format Reference

For detailed spec file format (frontmatter fields, sections, naming conventions), see [references/spec-format.md](references/spec-format.md).
