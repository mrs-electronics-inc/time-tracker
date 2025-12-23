---
status: draft
author: Addison Emig
creation_date: 2025-12-23
---

# Export

Add a CLI command to export data as TSV (tab-separated values). The goal is that the exported data can easily be imported into other software.

There should be two export formats to start:

- `daily-projects` (the default)
  - Each row represents all the tasks completed for a single project in a single day.
  - The tasks completed for the project in the given day are combined into a "description" column.
  - This should match the display output of the stats mode from spec [#5](./005-stats-mode.md)
  - Columns:
    - Project
    - Date (ISO 8601 format)
    - Duration (minutes)
    - Description
- `raw`
  - Each row represents a single time entry.
  - Filter out all the blank entries and include an End column and a Duration column
  - Columns:
    - Project
    - Task
    - Start (ISO 8601 format)
    - End (ISO 8601 format)
    - Duration (minutes)

## Design Decisions

### File Format

TSV is not the most amazing file format, but it seems the best option for our use case for the following reasons:

- TSV is superior to CSV because you don't have to add special handling for commas in your data.
- TSV can easily be imported into spreadsheet software while JSON can not.

## Task List

### TBD

- [ ] TBD
- [ ] TBD

### TBD

- [ ] TBD
- [ ] TBD

### TBD

- [ ] TBD
- [ ] TBD
