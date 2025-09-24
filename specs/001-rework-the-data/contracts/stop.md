# Contract: stop command

## Command
```
time-tracker stop
```

## Description
Stops the currently running time entry by setting its end timestamp to the current time.

## Arguments
None

## Behavior
1. Find the time entry with end = null
2. If found, set end to current timestamp
3. Save to data.json
4. Output confirmation message with duration

## Error Cases
- No currently running time entry: Display message "No active time entry to stop"
- Unable to write to data.json: Display error message

## Output
```
Stopped tracking time for "Implement feature" in project "time-tracker" (duration: 1h 30m)
```