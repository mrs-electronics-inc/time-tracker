---
number: 11
status: draft
author: Addison Emig
creation_date: 2026-03-13
---

# Edit Entry Start Date

Extend the edit form to allow users to edit the full start date (YYYY-MM-DD) of an entry, not just the time (HH:MM).

Currently, edit mode only allows changing the hour and minute, while preserving the original entry's date. This spec adds the ability to change the date as well.

## Design

### User Interaction

- When editing an entry, the form now shows three input fields for the start datetime:
  - Date (YYYY-MM-DD format)
  - Hour (HH)
  - Minute (MM)
- Tab/Shift+Tab navigation cycles through all fields (project, title, date, hour, minute)
- Validation ensures date is in valid ISO 8601 format
- Submit saves the entry with the new full datetime

### Form Layout

- Project
- Title
- Date (YYYY-MM-DD)
- Time (HH:MM)

## Task List

### Form Input Updates

- [ ] Test: Add date input field to Model.Inputs array
- [ ] Test: Date field displays in YYYY-MM-DD format
- [ ] Test: Tab navigation includes the new date field
- [ ] Implement: Add fourth input field for date before the time fields

### Edit Mode Integration

- [ ] Test: openEditMode pre-fills date field with entry's start date in YYYY-MM-DD format
- [ ] Test: Pre-fill hour and minute as before
- [ ] Implement: Update openEditMode to extract and set the date field

### Form Rendering

- [ ] Test: renderFormContent displays date label and input
- [ ] Test: Form layout shows Date field between Title and Time
- [ ] Implement: Update renderFormContent to include date field rendering

### Date Parsing & Validation

- [ ] Test: parseFormTime reads date from inputs and combines with time
- [ ] Test: Invalid date formats are rejected with clear error message
- [ ] Test: Valid dates like "2025-12-23" are parsed correctly
- [ ] Implement: Update parseFormTime to parse and validate date field
- [ ] Implement: Add date validation helper function

### Navigation & Focus

- [ ] Test: FocusIndex correctly increments/decrements across all 5 fields (project, title, date, hour, minute)
- [ ] Test: Shift+Tab from first field wraps to last field
- [ ] Implement: Adjust form field navigation logic for new field count
