# Quickstart: Rework the data model

## Setup
1. Ensure Go is installed
2. Clone the repository
3. Build the time-tracker CLI: `go build -o time-tracker`

## Basic Usage

### Start tracking time
```
./time-tracker start "my-project" "Working on feature"
```
Output: Started tracking time for "Working on feature" in project "my-project"

### Stop tracking time
```
./time-tracker stop
```
Output: Stopped tracking time for "Working on feature" in project "my-project" (duration: 2h 15m)

### List time entries
```
./time-tracker list
```
Output: Shows all time entries with ID, timestamps, project, title, and duration

## Test Scenarios

### Scenario 1: Start and stop a time entry
1. Run `start "test-project" "Test task"`
2. Wait a few seconds
3. Run `stop`
4. Run `list` - should show completed entry with duration

### Scenario 2: Auto-stop when starting new entry
1. Run `start "project1" "Task 1"`
2. Run `start "project2" "Task 2"`
3. Run `list` - should show Task 1 completed and Task 2 running

### Scenario 3: Stop when no active entry
1. Ensure no active entry (run `stop` if needed)
2. Run `stop`
3. Should show "No active time entry to stop"

## Data Storage
Time entries are stored in `data.json` in the current directory with format:
```json
{
  "time-entries": [
    {
      "id": 1,
      "start": "2025-09-24T10:00:00Z",
      "end": "2025-09-24T12:00:00Z",
      "project": "my-project",
      "title": "Working on feature"
    }
  ]
}
```