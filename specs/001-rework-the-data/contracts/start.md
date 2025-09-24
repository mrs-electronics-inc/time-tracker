# Contract: start command

## Command
```
time-tracker start <project> <title>
```

## Description
Starts a new time entry with the specified project and title. If another time entry is currently running, it will be automatically stopped first.

## Arguments
- `project` (string, required): Name of the project
- `title` (string, required): Title/description of the work

## Behavior
1. If a time entry is currently running (end = null), stop it by setting end to current timestamp
2. Create new time entry with:
   - id: next available numeric ID
   - start: current timestamp
   - end: null
   - project: provided argument
   - title: provided argument
3. Save to data.json
4. Output confirmation message

## Error Cases
- Missing project or title: Display usage error
- Unable to write to data.json: Display error message

## Output
```
Started tracking time for "Implement feature" in project "time-tracker"
```