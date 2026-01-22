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

## Project Metadata Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | yes | Human-readable display name |
| `code` | string | no | External system identifier (e.g., TWE ProjectNumber) |

## Export Format Changes

Update the `daily-projects` export format to include project metadata:

**Current columns:** `Project`, `Date`, `Duration`, `Description`

**New columns:** `ProjectName`, `ProjectCode`, `Date`, `Duration`, `Description`

Example output:

```
ProjectName	ProjectCode	Date	Duration	Description
Auth Refactor	24-0147	2026-01-20	120	Fixed authentication bug, Code review
API Updates	24-0150	2026-01-20	60	Updated documentation
```

If a project has no code defined, the `ProjectCode` column will be empty for that row.

## Design Decisions

### Storage Location

- **Decision**: Store project metadata under a `projects` key inside `data.json`
- **Rationale**: Single data file to manage; projects and time entries stay together

### Project Matching

- **Decision**: Time entries continue to store only the project name (no code). Export looks up the code from project metadata by matching the entry's project name.
- **Rationale**: No changes to time entry schema; keeps tracking simple. Code resolution happens at export time.

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

- [ ] Update `ExportDailyProjects` to output `ProjectName` and `ProjectCode` columns
- [ ] Test: Export includes `ProjectCode` when project has metadata
- [ ] Test: Export has empty `ProjectCode` when project has no metadata
- [ ] Test: Backward compatible with entries using undefined projects

### TUI Updates

- [ ] Show project code in TUI when available (e.g., "Auth Refactor (24-0147)")
- [ ] Consider project autocomplete from defined projects
