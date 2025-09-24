# Data Model: Rework the data model

## Entities

### Time Entry
Represents a period of tracked time.

**Fields**:
- `id` (int): Unique numeric identifier for the time entry
- `start` (timestamp): When the time tracking started
- `end` (timestamp, nullable): When the time tracking ended, null if currently running
- `project` (string): Name of the project
- `title` (string): Title/description of the work

**Validation Rules**:
- `id` must be unique across all time entries
- `start` must be a valid timestamp, set by program to current time
- `end` must be null or a timestamp after `start`
- `project` and `title` must be non-empty strings
- Only one time entry can have `end` = null at any time

**State Transitions**:
- Created: `end` = null, `start` = current time
- Running: `end` = null
- Stopped: `end` = current time

**Relationships**:
- No relationships to other entities (single active entry enforced by business logic)

## Storage Format
Time entries are stored in `data.json` with the following structure:
```json
{
  "time-entries": [
    {
      "id": 1,
      "start": "2025-09-24T10:00:00Z",
      "end": "2025-09-24T11:30:00Z",
      "project": "time-tracker",
      "title": "Implement data model"
    },
    {
      "id": 2,
      "start": "2025-09-24T14:00:00Z",
      "end": null,
      "project": "time-tracker",
      "title": "Add CLI commands"
    }
  ]
}
```

## Business Rules
- Only one time entry can be active (end = null) at any given time
- Starting a new entry automatically stops any currently running entry
- Stopping without arguments stops the currently running entry
- Time entries are identified by numeric ID for internal operations