# New Data Format

We need to improve the format of our data files to make it easier to implement future features.

The overall goal is to remove the `end` field and only have a `start` field on each time entry.

We will have `null` for `project and `title` for empty entries between other entries.

## Task List

### Basic Migration Logic

- [x] Add `version` field to the `data.json` files.
- [x] If the `version` field is missing, assume a version of `0`.

### Add Blank Time Entries

- [x] Add migration logic for version 0 to 1 when loading data.
- [x] Version 1 has the following differences from version 0:
  - [x] Add blank time entries between entries that have a space between the `end` of one entry and the `start` of the next entry.
  - [x] Blank entries representation choice: serialize `project` and `title` as empty strings (`""`)
  - [x] Assign IDs to inserted blank entries using sequential IDs continuing from the current max ID
  - [x] Migration behavior: insert blank entries in-memory only when loading; do not rewrite the on-disk file automatically

### Remove End Field

- [ ] Add migration logic for version 1 to 2 when loading data.
- [ ] Version 2 has the following differences from version 1:
  - [ ] No `end` field when saving time entries (the end of each entry is the start of the next entry).

### Remove Id Field

- [ ] Add migration logic for version 2 to 3 when loading data.
- [ ] Version 3 has the following differences from version 2:
  - [ ] No `id` field for the time entries.
- [ ] The ID column in the output of the list command should be automatically generated
