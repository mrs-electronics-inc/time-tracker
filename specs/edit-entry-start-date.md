---
number: 11
status: approved
author: Addison Emig
creation_date: 2026-03-13
approved_by: Addison Emig
approval_date: 2026-03-18
---

# Edit Entry Start Date

Extend the edit form to allow users to edit the full start date (YYYY-MM-DD) of an entry, not just the time (HH:MM).

Currently, edit mode only allows changing the hour and minute, while preserving the original entry's date. This spec adds the ability to change the date as well.

## Design

### User Interaction

- All form modes (new, edit, resume) now show three separate input fields for the date (year, month, day) in addition to the existing hour and minute fields
- Tab/Shift+Tab navigation cycles through all fields: Project → Title → Year → Month → Day → Hour → Minute
- Date fields default to today's date for new/resume modes, and the entry's existing date for edit mode
- The "assume yesterday" auto-adjustment logic is removed — the user-entered date is always used as-is
- Full date validation: month must be 1-12, day must be valid for the given month/year (including leap year handling)
- Submit saves the entry with the full datetime constructed from all input fields

### Form Layout

- Project
- Title
- Date: `YYYY - MM - DD` (three side-by-side inputs with separators, matching the time layout)
- Time: `HH : MM`

## Task List

### Add date fields to all forms

- [ ] Extract input index constants (`InputProject`, `InputTitle`, `InputHour`, `InputMinute`) in `types.go` and replace hardcoded input indices across all affected TUI code and tests
- [ ] Add year, month, day text inputs to `NewModel` in `model.go`, update index constants to include `InputYear`, `InputMonth`, `InputDay`, and update all form open functions (`openNewMode`, `openEditMode`, `openResumeMode`, `openStartMode`, `openStartModeBlank`) to set date defaults
- [ ] Update `renderFormContent` in `form.go` and `StartMode.RenderContent`/`renderStartContent` in `start.go` to render date as `YYYY - MM - DD` above time
- [ ] Update `parseFormTime` in `form.go` and the inline parsing in `StartMode.HandleKeyMsg` in `start.go` to read date from input fields, add full date validation, and remove the "assume yesterday" logic
