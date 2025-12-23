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

### Aggregation & Data Model

- [ ] Test: `ProjectDateEntry` struct groups entries by (project, date) with collected task descriptions
- [ ] Implement: `ProjectDateEntry` struct and `AggregateByProjectDate` function in `utils/`
- [ ] Test: Verify aggregation handles blank entries, running entries, and task deduplication
- [ ] Test: Verify weekly grouping for separator calculations

### Stats Mode

- [ ] Test: `StatsMode` renders table with Project | Date | Duration (minutes) | Description columns
- [ ] Test: `StatsMode` renders weekly separator rows with different styling at end of each week
- [ ] Test: `StatsMode` keyboard navigation (k/j, G, ?, q/esc)
- [ ] Implement: Create `stats.go` in `cmd/tui/modes/` with `StatsMode` definition, rendering, and navigation
- [ ] Test: Verify column width calculations and text wrapping for description lists
- [ ] Test: Verify viewport scrolling when content exceeds available height

### Integration

- [ ] Test: Tab key in `ListMode` switches to `StatsMode`
- [ ] Test: Stats mode shows keybinding hints in status bar with Tab to return to list
- [ ] Implement: Add Tab keybinding to `ListMode`, add `StatsMode` to `Model`, add navigation between modes
- [ ] Test: End-to-end: load sample data in TUI list mode, switch to stats, verify aggregation correctness
- [ ] Test: Verify stats mode handles edge cases (no data, single entry, data spanning multiple weeks)
