# CLI Command Contracts

## Command: track (aliases: s, start, stop)

### Invocations

#### Start Tracking
- **Alias**: start
- **Arguments**: <project> <task> (required)
- **Behavior**: Start new time tracking session
- **Output**: Success message with task details
- **Errors**: Missing arguments, invalid project/task

#### Stop Tracking
- **Alias**: stop
- **Arguments**: None (reject if provided)
- **Behavior**: Stop current tracking session
- **Output**: Session summary with duration
- **Errors**: No active session, arguments provided

#### Short Alias with Args (Start)
- **Alias**: s
- **Arguments**: <project> <task> (required)
- **Behavior**: Same as start
- **Output**: Same as start
- **Errors**: Same as start

#### Short Alias without Args (Stop)
- **Alias**: s
- **Arguments**: None (reject if provided)
- **Behavior**: Same as stop
- **Output**: Same as stop
- **Errors**: Same as stop

### Response Format
- Success: "Started tracking: [project] - [task]" or "Stopped tracking: [project] - [task] (duration)"
- Error: "Error: [message]"