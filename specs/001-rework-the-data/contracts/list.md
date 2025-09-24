# Contract: list command

## Command
```
time-tracker list
```

## Description
Lists all time entries from data.json.

## Arguments
None

## Behavior
1. Read data.json
2. Display all time entries in chronological order (newest first)
3. For each entry show: ID, start time, end time (or "running"), project, title, duration

## Error Cases
- data.json not found or unreadable: Display error message
- No time entries: Display "No time entries found"

## Output
```
ID  Start               End                 Project     Title               Duration
1   2025-09-24 10:00    2025-09-24 11:30    time-tracker Implement feature    1h 30m
2   2025-09-24 14:00    running             time-tracker Add CLI commands     30m
```