---
status: draft
author: Addison Emig
creation_date: 2025-12-23
---

# Export

Add a CLI command to export time tracking data as TSV (tab-separated values) for easy import into other software.

## Export Formats

### Daily-Projects Format (default)

Each row aggregates all tasks completed for a single project in a single day. The task titles are combined into a comma-separated "description" column. This format matches the display output of stats mode from [spec #5](./005-stats-mode.md).

**Columns:**
- Project
- Date (ISO 8601 format)
- Duration (minutes)
- Description (comma-separated task titles)

### Raw Format

Each row represents a single time entry (blank entries are filtered out).

**Columns:**
- Project
- Task
- Start (ISO 8601 format)
- End (ISO 8601 format)
- Duration (minutes)

## Design Decisions

### File Format: TSV

- **Decision**: Use tab-separated values (TSV) instead of CSV or JSON
- **Rationale**:
  - TSV handles commas in data without special escaping
  - Spreadsheet software can import TSV directly, unlike JSON
  - Simpler than CSV for this use case

### TSV Writing and Escaping

- **Decision**: Use Go's `encoding/csv` package with `Comma: '\t'` for all TSV output
- **Rationale**: Standard library correctly handles quoting and escaping, ensuring data integrity and round-trip safety (export → import → export produces identical output)
- **Implementation Requirements**:
  - All TSV output must include a header row as the first line
  - Use `csv.Writer` with `Comma` set to `'\t'`
  - Call `Flush()` after writing all records to ensure proper output
- **Testing Requirements**: All escaping tests must verify that fields with tabs, newlines, and quotes are properly quoted/escaped and can be read back correctly

### Handling Running and Incomplete Entries

- **Decision**: Exclude running entries (entries without End time) from both daily-projects and raw exports
- **Rationale**: Exports should contain only completed, finalized data. Running entries are transient state and will be included once they are completed. This avoids ambiguity about duration calculations and keeps exports deterministic.
- **Implementation**: Filter out entries where `End == nil` before processing for export

## Task List

### Daily-Projects Export

**Prerequisite**: Aggregation function from [spec #5](./005-stats-mode.md) (`ProjectDateEntry` and `AggregateByProjectDate`)

- [ ] Test: `ExportDailyProjects` writes TSV with header row and correct columns (Project, Date, Duration, Description)
- [ ] Test: `ExportDailyProjects` uses `encoding/csv` with `Comma: '\t'` for proper escaping
- [ ] Test: Description column joins task titles with comma-space separator (`, `) for single-line TSV format
- [ ] Test: Fields with tabs, newlines, and quotes are properly quoted and can round-trip (read back correctly)
- [ ] Test: `ExportDailyProjects` excludes running entries (entries without End time)
- [ ] Test: `ExportDailyProjects` converts durations to minutes
- [ ] Implement: `ExportDailyProjects` function in `utils/` that takes `ProjectDateEntry` slice and returns TSV string
- [ ] Test: Export command writes TSV to stdout and to file
- [ ] Implement: Add `export` CLI command with `--format daily-projects` (default), output to stdout or `--output` file
- [ ] Test: End-to-end: load sample data, export, verify TSV contents match aggregated data

### Raw Export

- [ ] Test: `ExportRaw` writes TSV with header row and correct columns (Project, Task, Start, End, Duration)
- [ ] Test: `ExportRaw` uses `encoding/csv` with `Comma: '\t'` for proper escaping
- [ ] Test: `ExportRaw` filters out blank entries
- [ ] Test: `ExportRaw` excludes running entries (entries without End time)
- [ ] Test: Fields with tabs, newlines, and quotes are properly quoted and can round-trip (read back correctly)
- [ ] Test: `ExportRaw` converts durations to minutes
- [ ] Implement: `ExportRaw` function in `utils/` that takes `TimeEntry` slice and returns TSV string
- [ ] Implement: Add `--format raw` option to export command
- [ ] Test: End-to-end: load sample data, export raw format, verify TSV contains all completed entries (excluding blank and running)
