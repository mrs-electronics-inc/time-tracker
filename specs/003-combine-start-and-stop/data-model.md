# Data Model: combine-start-and-stop

## Overview
This feature does not introduce new data entities or modify existing ones. The existing Task and TimeEntry models remain unchanged.

## Existing Entities

### Task
- **Fields**: ID (string), Name (string), Project (string)
- **Relationships**: One-to-many with TimeEntry
- **Validation**: Name and Project required

### TimeEntry
- **Fields**: ID (string), TaskID (string), StartTime (time.Time), EndTime (*time.Time), Duration (time.Duration)
- **Relationships**: Many-to-one with Task
- **Validation**: StartTime required, EndTime optional, Duration calculated

## State Transitions
- Task: Created → Active (when started) → Completed (when stopped)
- TimeEntry: Created (start) → Completed (stop)

No changes required for this feature.