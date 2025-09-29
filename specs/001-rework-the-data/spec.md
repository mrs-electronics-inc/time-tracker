# Feature Specification: Rework the data model

**Feature Branch**: `001-rework-the-data`  
**Created**: Wed Sep 24 2025  
**Status**: Completed  
**Summary**: Rework the data model to store time entries instead of tasks, with each entry having id, start timestamp, end timestamp (nullable), project, and title. Only one time entry can be active at a time. Storage in data.json with {"time-entries": []} format. Update stop command to automatically stop the running entry without arguments.

## Technical Details

### Data Model

- **Time Entry Entity**: id (int, unique), start (timestamp), end (timestamp, nullable), project (string), title (string)
- **Validation Rules**: id unique, start valid timestamp, end null or after start, project/title non-empty, only one entry with end=null
- **State Transitions**: Created (end=null), Running (end=null), Stopped (end=current time)
- **Storage Format**: JSON file data.json with {"time-entries": [{"id":1,"start":"2025-09-24T10:00:00Z","end":"2025-09-24T11:30:00Z","project":"time-tracker","title":"Implement feature"}]}

### Commands

- **list**: Lists all entries newest first, shows ID, start, end, project, title, duration
- **start <project> <title>**: Creates new entry, auto-stops current if running
- **stop**: Stops current running entry, sets end to current time

### Business Rules

- Only one active entry (end=null) at any time
- Starting new entry auto-stops current
- Stopping without args stops current active entry
- Entries identified by numeric ID

## Decisions Made

- **Language/Framework**: Go with Cobra CLI (constitution mandated)
- **Storage**: JSON file format (constitution mandated)
- **Testing**: TDD with Go testing (constitution mandated)
- **Structure**: Single project with src/ and tests/ (CLI tool scope)
- **Performance**: Startup <100ms, low memory usage (constitution mandated)

## Implementation Status

All tasks completed: project structure, Go/Cobra setup, contract tests, integration tests, TimeEntry model, TaskManager service, command implementations, JSON storage, error handling, unit tests, performance tests, README updates.
