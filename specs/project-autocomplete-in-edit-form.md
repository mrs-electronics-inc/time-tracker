---
number: 12
status: draft
author: Addison Emig
creation_date: 2026-03-13
---

# Project Autocomplete in Edit Form

## Overview

Once spec 9 (Project Metadata) ships, the edit form’s `Project` input should stop being a free-form field and only accept values from the managed project list. This spec introduces inline autocomplete powered exclusively by the metadata-backed project list so edited entries stay aligned with the canonical names exports and integrations rely on. There is no fallback to historical entries; the feature depends on the metadata list already existing.

## User Flow

1. Select an entry in the list and press `e` to open edit mode (with the form defaulting to the entry’s project/title/time).
2. While the project input remains focused, typing begins to filter the dropdown directly beneath the input showing matching project names.
3. The autocomplete list performs fuzzy matching (case-insensitive, substring) over the metadata project names and updates in real time, showing only the names (no code/category).
4. Use Up/Down arrow keys to highlight matches; Enter selects the highlighted name and keeps focus on the project input. Tab/Shift+Tab continue to rotate focus between the four form fields exactly as today.
5. Clearing the input or moving focus out closes the dropdown; selecting a project hides it until the user types again.

## Design Decisions

1. **Suggestion source**
   - *Options*: (A) fallback to historical projects, (B) rely solely on the metadata list from spec 9.
   - *Decision*: choose option B so edited entries match the authoritative project definitions that exports use. This enforces consistency immediately after metadata is available.

2. **Matching strategy**
   - *Options*: prefix-only or fuzzy matching.
   - *Decision*: fuzzy matching so users can locate projects using any memorable substring or partial words instead of needing exact prefixes.

3. **Interaction model**
   - *Options*: reuse Tab for suggestions or keep Tab cycling inputs.
   - *Decision*: keep Tab exclusively for field navigation (as it is now) while using arrow keys + Enter for suggestion selection to preserve existing muscle memory.

## Task List

### Data plumbing

- [ ] Provide a helper that returns the ordered list of project names from the spec 9 metadata store (e.g., a public method on the storage layer or task manager).
- [ ] Expose that list through the TUI model so form modes (start/edit) can read it while rendering.
- [ ] Add tests validating the helper’s behavior with populated metadata and when the metadata list is empty.

### Autocomplete behavior

- [ ] Extend the form state to track the current autocomplete matches, the active suggestion index, and whether the dropdown is visible.
- [ ] Implement fuzzy filtering (case-insensitive substring match) keyed to the project input value.
- [ ] Render the matching project names in a dropdown below the project input, highlighting the active suggestion with focused styling consistent with the rest of the form.
- [ ] Handle keyboard navigation: Up/Down move through the list when visible, Enter picks the highlighted suggestion and writes it into the project field, and Tab/Shift+Tab keep cycling form inputs.
- [ ] Close the dropdown when the project field loses focus or the value is cleared.

### Verification

- [ ] Add headless/integration coverage that types into the project field inside edit mode, asserts the suggestion list updates, navigates it, and confirms Enter populates the value while Tab still changes inputs.
- [ ] Unit-test the fuzzy matcher so it respects substring matching, ignores case, and tolerates empty metadata lists without panicking.
- [ ] Document the new behavior in the relevant user-facing guides (e.g., README or TUI help text) once implementation is wired up.
