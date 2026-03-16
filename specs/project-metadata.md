---
number: 9
status: in-progress
author: Addison Emig
creation_date: 2026-01-22
approved_by: Addison Emig
approval_date: 2026-03-16
---

# Project Metadata

Add project metadata support to enable structured project definitions with human-readable names and external system codes. This enables integration with external systems (e.g., TWE, issue trackers) while keeping time-tracker generic and user-friendly.

Currently, projects are just free-form strings. This makes it difficult to:

- Maintain consistent project naming
- Export data for import into external systems that require specific project identifiers
- Track time against projects with both a display name and an external code

## Data Format (v4)

Projects are stored in a top-level `projects` array in `data.json`:

```json
{
  "version": 4,
  "time-entries": [
    {
      "start": "2025-12-23T09:00:00Z",
      "project": "Auth Refactor",
      "title": "Fixed authentication bug"
    }
  ],
  "projects": [
    {
      "name": "Auth Refactor",
      "code": "12572",
      "category": "Infrastructure"
    },
    {
      "name": "API Updates",
      "code": "12573",
      "category": "Backend"
    }
  ]
}
```

## Project Metadata Fields

| Field      | Type   | Required | Description                                         |
| ---------- | ------ | -------- | --------------------------------------------------- |
| `name`     | string | yes      | Human-readable display name                         |
| `code`     | string | no       | External system identifier (e.g., project number)   |
| `category` | string | no       | Free-form category (e.g., "Engineering", "Support") |

## Export Format Changes

Update the `daily-projects` export format to include project metadata. Format is **TSV** (tab-delimited).

**Changes:** Replace `Project` with `ProjectName`, add `ProjectCode`, and add `ProjectCategory`

Example output:

```
ProjectName	ProjectCode	ProjectCategory	Date	Duration	Description
Auth Refactor	12572	Infrastructure	2026-01-20	120	Fixed authentication bug, Code review
API Updates	12573	Backend	2026-01-20	60	Updated documentation
```

Empty cells: If a project has no code or category defined, those columns will be empty for that row.

## Design Decisions

### Storage Location

- **Decision**: Store project metadata under a `projects` key inside `data.json`
- **Rationale**: Single data file to manage; projects and time entries stay together

### Data and Persistence Contract

- **Decision**: Bump storage format to version 4.
- **Decision**: Keep existing time-entry methods and extend storage with project-specific methods:
  - `Load() ([]TimeEntry, error)`
  - `Save([]TimeEntry) error`
  - `LoadProjects() ([]Project, error)`
  - `SaveProjects([]Project) error`
- **Decision**: `Project` has `name`, `code`, and `category` string fields; missing optional values are stored as explicit empty strings (`""`), not omitted.
- **Decision**: Continue current migration behavior: migrate older versions in memory on load, persist v4 on next save.
- **Decision**: Operations that mutate both entries and projects (rename/merge) must persist atomically in one file write.
- **Rationale**: Aligns with existing code patterns while adding project support safely and without broad interface churn.

### Naming, Validation, and Matching

- **Decision**: Project name uniqueness is case-insensitive, while preserving original casing for display and storage.
- **Decision**: Add/edit/rename trim leading/trailing whitespace before validation.
- **Decision**: Empty name after trim is invalid; empty `code` and `category` are allowed.
- **Decision**: Time entries continue to store project name strings. Matching from entries to metadata uses exact, case-sensitive name equality.
- **Rationale**: Prevents confusing near-duplicates while keeping deterministic matching behavior.

### Rename and Merge Behavior

- **Decision**: Renaming to a non-existing name rewrites all exact old-name entry references to the new name.
- **Decision**: Renaming to an existing name is a merge:
  - Rewrite all exact old-name entry references to the target name.
  - Remove the source project metadata record.
  - Keep target metadata canonical; on conflict, target `code`/`category` win.
- **Decision**: Case-only renames are valid and rewrite exact old-name matches.
- **Decision**: Successful rename operations (including merge) report rewritten entry count.
- **Decision**: Metadata-only edits report a simple success message without rewrite count.
- **Rationale**: Keeps data consistent and gives users clear, predictable outcomes.

### Remove Behavior

- **Decision**: Removing a project is blocked when any time entries reference it.
- **Decision**: Blocked remove returns a consistent error in CLI and TUI that includes reference count.
- **Rationale**: Avoids silently orphaning metadata-backed workflows.

### Export Behavior

### Backward Compatibility

- **Decision**: Projects without metadata entries continue to work; they just won't have a code in exports
- **Rationale**: Non-breaking change; users can adopt project metadata incrementally
- **Decision**: `daily-projects` header becomes: `ProjectName`, `ProjectCode`, `ProjectCategory`, `Date`, `Duration`, `Description`.
- **Decision**: Add `export --category <value>` filter:
  - Input is trimmed.
  - Whitespace-only value is invalid.
  - Match is exact string and case-insensitive (`strings.EqualFold` semantics).
  - Without `--category`, include all rows.
  - With `--category`, include only rows whose project metadata exists and matches; undefined projects are excluded.
- **Rationale**: Keeps exports backward compatible by default while enabling reliable category-focused subsets.

### TUI Navigation and Scope

- **Decision**: Spec 9 includes a full `projects` TUI view.
- **Decision**: Mode switching uses `Tab` cycle: `list -> stats -> projects -> list`.
- **Decision**: `projects` mode uses `j/k` (and arrows) for navigation and `n/e/d` for add/edit/delete.
- **Decision**: Add/edit project actions use a dedicated project form with fields `Name` (required), `Code` (optional), and `Category` (optional).
- **Decision**: Time-entry forms remain free-form in spec 9; autocomplete/strict project selection is deferred to spec 12.
- **Rationale**: Preserves current interaction patterns while delivering project management in the TUI.

## Task List

### Storage Foundation

- [x] Define `Project` struct with `Name`, `Code`, and `Category` fields.
- [x] Add `projects` key to `data.json` schema.
- [x] Extend storage interfaces and implementations with `LoadProjects`/`SaveProjects`.
- [x] Handle missing `projects` key gracefully (empty project list).

### Versioning and Persistence

- [x] Bump `CurrentVersion` to 4 and extend migration/load-save paths for v4 data.
- [x] Implement atomic persistence for operations that update entries and projects together.

### Domain Logic

- [x] Implement project mutation logic in `TaskManager` and keep command handlers thin.

### Project Management Commands

- [x] Implement `project list` command with columns `Name`, `Code`, `Category`, sorted case-insensitively by name.
- [ ] Implement `project add <name> [--code <code>] [--category <category>]` with trim + validation and case-insensitive uniqueness enforcement.
- [ ] Implement `project edit <name> [--name <new-name>] [--code <code>] [--category <category>]` with rename/merge semantics and rewrite count reporting for rename operations.
- [ ] Implement `project remove <name>` with reference blocking and reference-count error.

### Export Columns

- [ ] Update `ExportDailyProjects` to output `ProjectName`, `ProjectCode`, and `ProjectCategory` columns.
- [ ] Ensure backward compatibility with entries using undefined projects.

### Category Filter

- [ ] Implement `--category` filter for export command to show only entries from specified category.

### TUI Updates

- [ ] Add `projects` view alongside `list` and `stats` views.
- [ ] Implement `Tab` cycle across `list`, `stats`, and `projects` modes.
- [ ] In `projects` view: scroll through all projects, sorted case-insensitively by name.
- [ ] In `projects` view: add and edit projects via project form (`Name`, `Code`, `Category`).
- [ ] In `projects` view: delete projects with the same reference-blocking behavior as CLI.
- [ ] Keep time-entry project fields free-form in this spec; project autocomplete remains in spec 12.
