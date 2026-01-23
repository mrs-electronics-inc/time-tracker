---
status: completed
author: Bennett Moore
creation_date: 2025-12-18
approved_by: Addison Emig
approval_date: 2025-12-18
---

# Command to View the Current Task

It would be great to have a subcommand, available directly from command-line without the TUI, that outputs the currently running task.

The command could be `current`, aliases `curr` and `c`.

## Task List

- [x] Create convenience function to return current running task
- [x] Output something like:

  ```text
  mrs-sdk-qt review #46, duration 18m
  ```

- [x] Show nothing if no task is running
