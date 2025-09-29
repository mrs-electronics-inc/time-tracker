# Feature: Combine Start and Stop Commands

## Overview
Combine the existing start and stop commands into a single command accessible via aliases `s`, `start`, or `stop`. When called as `stop`, reject any arguments and prompt to use `start`. When called as `start`, require arguments or error.

## Clarifications
- **'s' alias behavior**: With arguments: start tracking; without arguments: stop tracking

## User Scenarios
### Primary User Story
Users need a unified command to manage time tracking sessions, allowing them to start or stop tracking with simple aliases, ensuring consistent behavior based on how the command is invoked.

### Acceptance Scenarios
1. **Given** no active time tracking session, **When** user runs `start "project-name" "task-name"`, **Then** a new tracking session starts for the specified project and task.
2. **Given** an active time tracking session, **When** user runs `stop`, **Then** the current session stops and time is recorded.
3. **Given** an active time tracking session, **When** user runs `stop "any-argument"`, **Then** system errors and prompts user to use `start` instead.
4. **Given** no active session, **When** user runs `start` without arguments, **Then** system errors due to missing required arguments.
5. **Given** no active time tracking session, **When** user runs `s "project-name" "task-name"`, **Then** a new tracking session starts for the specified project and task.
6. **Given** an active time tracking session, **When** user runs `s`, **Then** the current session stops and time is recorded.

### Edge Cases
- When user runs `s` with arguments, it starts tracking with the provided project and task.
- When user runs `s` without arguments, it stops the current tracking session.
- The system distinguishes invocations by the alias used: 'start' and 's' with args behave as start, 'stop' and 's' without args behave as stop.

## Requirements
### Functional Requirements
- **FR-001**: System MUST provide a single command accessible via aliases 's', 'start', and 'stop'.
- **FR-002**: When the command is invoked as 'stop', system MUST reject any provided arguments and prompt user to use 'start'.
- **FR-003**: When the command is invoked as 'start', system MUST require project and task arguments or error.
- **FR-004**: When the command is invoked as 's', system MUST start tracking if arguments are provided, or stop tracking if no arguments are provided.

## Implementation Summary
Combined the existing start and stop commands into a single command accessible via aliases 's', 'start', and 'stop'. The command behavior depends on the alias used and arguments provided: 'start' and 's' with args start tracking, 'stop' and 's' without args stop tracking. Implemented using Cobra's command aliases and CalledAs() to determine invocation mode.

## Technical Decisions
- **Command Aliases**: Use Cobra's native alias support with CalledAs() detection for different behaviors.
- **Argument Validation**: Strict validation based on invocation: start requires 2 args or 1 ID, stop rejects args, s adapts based on arg count.
- **Data Models**: No changes - existing Task and TimeEntry models unchanged.
- **Testing**: Follow Go conventions with _test.go suffix.

## Technical Context
- **Language**: Go
- **Framework**: Cobra (CLI)
- **Storage**: JSON file format
- **Platform**: Linux
- **Performance**: <100ms startup

## Usage Examples
```bash
# Start tracking
docker compose run time-tracker start "My Project" "Task 1"
docker compose run time-tracker s "My Project" "Task 1"

# Resume by ID
docker compose run time-tracker start 5

# Stop tracking
docker compose run time-tracker stop
docker compose run time-tracker s
```

## Files Changed
- `src/cmd/track.go` (new combined command)
- `src/cmd/start.go` (removed)
- `src/cmd/stop.go` (removed)
- `README.md` (updated usage)
- `Dockerfile` (removed USER, added dirs)
- `docker-compose.yml` (added volumes)
- `AGENTS.md` (updated)