# Improve TUI Shortcuts

The current start/stop shortcut in the TUI is confusing. It does different things based on the context. We should improve the list mode shortcuts.

## Shortcuts

These shortcuts apply to the list mode

- `n` - new
  - Open the **start entry form** with empty project and task
- `s` - stop
  - Keep as-is
  - Does nothing on a blank entry
- `r` - resume (replaces start)
  - Duplicate selected project and task into new task
  - Use the **start entry form**, so the user can select the start time
  - If used on a blank entry, it acts just the same as the **new** shortcut
- `e` - edit
  - Go to **edit entry form** and allow user to edit any of the fields
  - The **edit entry form** will look the exact same as the **start entry form**, just with a different title (ideally it will be the exact same form behind the scenes)
- `d` - delete
  - Make the selected entry a blank entry
  - There should be a confirmation dialog to protect against accidental deletions

## Task List

TBD
