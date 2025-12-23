---
status: draft
author: Addison Emig
creation_date: 2025-12-23
---

# Stats Mode

We have a basic stats command for the CLI, which allows the user to view aggregated stats by day or by week.

We should add a stats mode to the TUI. This mode will be available via the `Tab` key in list mode.

The stats mode will aggregate the data a bit differently than the existing stats command. This new stats output will align more closely with the `daily-projects` export format from [#6](./006-export.md).

We will display a list of rows, where each row has the following columns:

- Project
- Date (ISO 8601 format)
- Duration (minutes)
- Description

The description column should contain a bullet-point list of all the tasks completed for the given combination of project and date.

The stats mode will be a read-only view of the data.

We also need to insert a different-colored row at the end of each week that aggregates the time spent across all projects for that week.

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
