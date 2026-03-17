---
number: 10
status: in-progress
author: Addison Emig
creation_date: 2026-03-13
approved_by: Addison Emig
approval_date: 2026-03-17
---

# Search Shortcut

Add a search/filter feature to list mode that allows users to quickly filter entries by project and title fields.

## Design

### User Interaction

- Press `/` to enter search mode and show the search input bar
- Type the search query and press Enter to apply the filter (no live filtering while typing)
- The list displays only entries matching the query (project or title contains the search term, case-insensitive)
- Navigate filtered results using existing keys: `j`/`k` for up/down, `G` for end
- Press `/` again to focus and edit the search query while keeping the currently applied results
- Press `Esc` while editing to clear the search, exit search mode, and show all entries
- Search uses a simple query string only (no field syntax, regex, or fuzzy matching in this version)

### UI Layout

The search input bar appears between the table rows and the status bar, and remains visible while search mode is active:

- Header row
- Table rows (filtered if search is active)
- Search input bar (when search is active)
- Status bar

### Selection and Empty State

- When a filter is applied:
  - Keep the current selection if the selected entry still matches
  - Otherwise move selection to the last filtered result
- When there are zero matches, use the same empty-state behavior as no data (no actionable selection), but with search-specific messaging
- Filter state persists when switching away from list mode and returning

## Task List

### Search State & Logic

- [x] Add search state to TUI model (active flag, query draft, applied query)
- [ ] Implement case-insensitive substring matcher across `project` and `title`
- [ ] Add filtering helper that returns visible entries while preserving source entry index mapping
- [ ] Implement apply-search behavior on `Enter` (update filtered list from applied query)
- [ ] Implement clear-search behavior on `Esc` while editing (clear query, exit search mode, restore full list)

### List Mode Integration

- [ ] Add `/` key handling in list mode to enter/focus search input
- [ ] Keep normal list navigation (`j`/`k`/`G`) for filtered results outside input-edit focus
- [ ] Preserve current selection on apply when still matched; otherwise select last filtered result
- [ ] Ensure selection is non-actionable when filtered result count is zero

### UI & Rendering

- [ ] Render search input bar between rows and status bar when search mode is active
- [ ] Keep search input bar visible after filter apply while search mode remains active
- [ ] Show search-specific empty message when filter has zero matches (distinct from no-data message)
- [ ] Ensure viewport and row rendering use filtered rows consistently

### Navigation & Selection

- [ ] Update selection/viewport helpers to work with filtered row sets and mapped source indices
- [ ] Ensure edit/resume/delete/stop actions operate on the correct underlying entry from filtered selection
- [ ] Persist active filter state across mode switches and list re-entry
