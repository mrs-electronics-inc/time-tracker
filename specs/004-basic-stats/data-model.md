# Data Model: Basic Stats

## Entities

### TimeEntry
Represents a completed period of tracked time.

**Fields**:
- ID (int): Unique identifier
- Start (time.Time): When tracking started
- End (*time.Time): When tracking ended (nil if running)
- Project (string): Project name
- Title (string): Task title

**Validation Rules**:
- Start must be before End (if End exists)
- Project and Title cannot be empty
- Duration must be positive

**Relationships**:
- Belongs to Project (via Project field)

### Project
Represents a grouping of tasks.

**Fields**:
- Name (string): Project identifier

**Validation Rules**:
- Name cannot be empty
- Unique across time entries

**Relationships**:
- Has many TimeEntry

## State Transitions
TimeEntry states: Created → Running → Completed
- Created: End is nil, IsRunning() true
- Completed: End set, Duration() calculated

## Data Volume Assumptions
- Typical: 100-1000 time entries per user
- Peak: 10,000+ entries
- Storage: JSON file, loaded entirely into memory for stats calculation