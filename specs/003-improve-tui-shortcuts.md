---
status: in-progress
author: Addison Emig, Bennett Moore, Jason Luke
creation_date: 2025-12-18
---

# Improve TUI Shortcuts

The current start/stop shortcut in the TUI is confusing. It does different things based on the context. We should improve the list mode shortcuts.

## Shortcuts

These shortcuts apply to the list mode

- `n` - new
  - Open the unified form in "new" mode with empty project and task
  - User selects start time
- `s` - stop
  - Stop the currently running entry
  - Does nothing on blank entries or non-running entries
- `r` - resume
  - Open the unified form in "resume" mode with project and task pre-filled from selected entry
  - User selects start time
  - Disabled on blank entries (does nothing)
- `e` - edit
  - Open the unified form in "edit" mode with all fields pre-filled from selected entry
  - User can edit project, task, and start time (not end time)
  - Note: End time is not editable because our duration model derives it from the next entry's start time. Editing end time would require complex logic (inserting blanks or adjusting next entry). Keep it simple for now.
- `d` - delete
  - Show a confirmation modal dialog
  - On confirmation, make the selected entry a blank entry (gap)

## Task List

### Unified Form Infrastructure

- [ ] Create shared form helpers (in start.go or new file)
  - `renderFormContent(m *Model, title string)` - renders form with given title
  - `handleFormSubmit(m *Model, mode string)` - handles submit for new/edit/resume
  - `openNewMode(m *Model)` - setup and open new entry form
  - `openEditMode(m *Model, entry TimeEntry)` - setup and open edit form
  - `openResumeMode(m *Model, entry TimeEntry)` - setup and open resume form
- [ ] Create `NewMode`, `EditMode`, `ResumeMode` using shared helpers
  - Keep `StartMode` for backward compatibility or remove if not used

### Implement List Shortcuts

- [ ] Implement `n` shortcut in `list.go` - calls `openNewMode()`
- [ ] Implement `r` shortcut in `list.go` - calls `openResumeMode()`, disabled on blank entries
- [ ] Implement `e` shortcut in `list.go` - calls `openEditMode()`
- [ ] Refactor `s` shortcut - only stop running entries, does nothing on blank/non-running

### Delete Confirmation Modal

- [ ] Implement `d` shortcut with delete confirmation modal dialog
  - Create `ConfirmMode` for the modal
  - Shows entry details and Yes/No buttons

### Polish & Testing

- [ ] Update KeyBindings in list.go to show all shortcuts (n, s, r, e, d)
- [ ] Add tests for new shortcuts
