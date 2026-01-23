---
status: draft
author: Addison Emig
creation_date: 2026-01-22
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

### Project Matching and Rename Behavior

- **Decision**: Time entries store project name. Export looks up code by exact name match. Project names are case-sensitive and must be unique. When renaming a project, update the project entry name AND rewrite all matching time entries to the new name.
- **Rationale**: Simple, deterministic matching. Single source of truth per project name. Rename is non-breaking because all entries stay consistent.

### Backward Compatibility

- **Decision**: Projects without metadata entries continue to work; they just won't have a code in exports
- **Rationale**: Non-breaking change; users can adopt project metadata incrementally

## Task List

### Project Storage

- [ ] Define `Project` struct with `Name` and `Code` fields
- [ ] Add `projects` key to `data.json` schema
- [ ] Update storage to load/save projects from `data.json`
- [ ] Test: Load projects from `data.json`
- [ ] Test: Save projects to `data.json`
- [ ] Test: Handle missing `projects` key gracefully (empty project list)

### Project Management Commands

- [ ] Implement `project list` command to display all projects
- [ ] Implement `project add <name> [--code <code>]` command
- [ ] Implement `project edit <name> [--code <code>]` command
- [ ] Implement `project remove <name>` command
- [ ] Test: Add project with name only
- [ ] Test: Add project with name and code
- [ ] Test: Edit project code
- [ ] Test: Remove project

### Export Updates

- [ ] Update `ExportDailyProjects` to output `ProjectName`, `ProjectCode`, and `ProjectCategory` columns
- [ ] Implement `--category` filter for export command to show only entries from specified category
- [ ] Test: Export includes `ProjectCode` and `ProjectCategory` when project has metadata
- [ ] Test: Export has empty `ProjectCode` and `ProjectCategory` when project has no metadata
- [ ] Test: `--category` filter works correctly
- [ ] Test: Backward compatible with entries using undefined projects

### TUI Updates

- [ ] Add `projects` view alongside `list` and `stats` views
- [ ] In `projects` view: scroll through all projects
- [ ] In `projects` view: add new projects (name, code, category)
- [ ] In `projects` view: edit existing projects (name, code, category)
- [ ] Consider project autocomplete from defined projects when creating entries
