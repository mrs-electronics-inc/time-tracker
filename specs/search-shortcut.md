---
number: 10
status: draft
author: Addison Emig
creation_date: 2026-03-13
---

# Search Shortcut

Add a search/filter feature to list mode that allows users to quickly filter entries by project and title fields.

## Design

### User Interaction

- Press `/` to enter search mode and show the search input bar
- Type the search query and press Enter to apply the filter
- The list displays only entries matching the query (project or title contains the search term, case-insensitive)
- Navigate filtered results using existing keys: `j`/`k` for up/down, `G` for end
- Press `/` again to edit the search query
- Clear the search term (empty string) and press Enter to show all entries and exit search mode

### UI Layout

The search input bar appears between the table rows and the status bar:
- Header row
- Table rows (filtered if search is active)
- Search input bar (when search is active)
- Status bar

## Task List

### Search State & Logic

- [ ] Test: SearchTerm field tracks current search query
- [ ] Test: FilteredEntries returns entries where project or title contains SearchTerm (case-insensitive)
- [ ] Implement: Add SearchTerm field to Model, add FilteredEntries helper function

### List Mode Integration

- [ ] Test: `/` key in list mode enters search mode and shows input bar
- [ ] Test: Search input accepts text and updates SearchTerm on Enter
- [ ] Test: Empty SearchTerm + Enter clears filter and exits search mode
- [ ] Test: `/` from filtered view allows editing the search query
- [ ] Test: Existing navigation (j/k/G) works on filtered results
- [ ] Implement: Add "/" key handler to ListMode that toggles search mode
- [ ] Implement: Add search input rendering in ListMode.RenderContent

### UI & Rendering

- [ ] Test: Search input bar displays with prompt indicator (e.g., "Search:") and current query
- [ ] Test: When filter is active, list shows only matching entries
- [ ] Test: Column widths adjust correctly with filtered results
- [ ] Implement: renderSearchBar function renders the input field
- [ ] Implement: Update ListMode.RenderContent to include search bar when active

### Navigation & Selection

- [ ] Test: SelectedIdx is reset to 0 when filter is applied
- [ ] Test: Navigation in filtered list doesn't exceed filtered entry count
- [ ] Test: Viewport scrolling works correctly with filtered entries
