---
number: 12
status: in-progress
author: Addison Emig
creation_date: 2026-03-13
approved_by: Addison Emig
approval_date: 2026-03-18
---

# Project Autocomplete in Entry Forms

## Overview

Once [Project Metadata](./project-metadata.md) ships, the `Project` input in all form modes (new, edit, resume) should offer inline autocomplete suggestions from the project list. The input remains free-text — users can type any project name — but suggestions encourage using canonical names that exports and integrations rely on.

## User Flow

1. Select an entry in the list and press `e` to open edit mode (with the form defaulting to the entry’s project/title/time).
2. While the project input is focused, typing shows a ghost-text completion inline (the first matching project name, greyed out after the cursor).
3. Up/Down arrow keys cycle through matching suggestions.
4. Tab accepts the current suggestion. Tab/Shift+Tab continue to navigate between form fields when there is no active suggestion.
5. The input remains free-text — the user can ignore suggestions and type any project name.

## Design Decisions

### Use Built-in Bubbles Suggestions

- **Decision**: Use the `textinput.SetSuggestions` / `ShowSuggestions` API from `charmbracelet/bubbles` rather than building a custom dropdown. This provides inline ghost-text completion with prefix matching out of the box.
- **Decision**: Before calling `SetSuggestions`, normalize project names by trimming whitespace, dropping empty names, deduping case-insensitively while preserving the first-seen display casing, and sorting case-insensitively (with raw string tie-breaker) for deterministic ordering.

### Keybinding Changes

- **Decision**: Remove Up/Down as field navigation aliases (currently redundant with Tab/Shift+Tab). When the project input has matched suggestions, Up/Down cycle suggestions and Tab accepts. When there are no suggestions, Tab/Shift+Tab navigate fields as before.

## Task List

### Keybinding Changes

- [x] Remove Up/Down as field navigation aliases from `handleFormKeyMsg`
- [x] Update form help text to show only Tab/Shift+Tab for field navigation

### Project Autocomplete

- [x] Enable `ShowSuggestions` on the project text input
- [x] Load project names via `LoadProjects`, normalize (trim, drop empty, case-insensitive dedupe preserving first-seen display casing, case-insensitive sort with raw string tie-breaker), and call `SetSuggestions` when opening form modes (new, edit, resume)
- [ ] When project input is focused and has matched suggestions, let Tab pass through to accept the suggestion instead of navigating fields
