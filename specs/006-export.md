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

### Daily-Projects Export

**Prerequisite**: Aggregation function from [spec #5](./005-stats-mode.md) (ProjectDateEntry and AggregateByProjectDate)

- [ ] Test: ExportDailyProjects formats aggregated data as TSV with correct columns and headers
- [ ] Test: ExportDailyProjects escapes tabs and newlines in task descriptions correctly
- [ ] Test: ExportDailyProjects converts durations to minutes
- [ ] Implement: ExportDailyProjects function in utils/ that takes ProjectDateEntry slice and returns TSV string
- [ ] Test: Export command writes TSV to stdout and to file
- [ ] Implement: Add `export` CLI command with `--format daily-projects` (default), output to stdout or `--output` file
- [ ] Test: End-to-end: load sample data, export, verify TSV contents match aggregated stats display

### Raw Export

- [ ] Test: ExportRaw formats raw time entries as TSV with correct columns (Project, Task, Start, End, Duration) and headers
- [ ] Test: ExportRaw filters out blank entries
- [ ] Test: ExportRaw escapes tabs and newlines in project/task names correctly
- [ ] Test: ExportRaw converts durations to minutes
- [ ] Implement: ExportRaw function in utils/ that takes TimeEntry slice and returns TSV string
- [ ] Implement: Add `--format raw` option to export command
- [ ] Test: End-to-end: load sample data, export raw format, verify TSV contains all non-blank entries
