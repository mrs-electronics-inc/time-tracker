# Basic Stats Command

## Overview
Implement a new `stats` command for the Time Tracker CLI to display basic time tracking statistics. The command provides daily totals (default), weekly totals (--weekly), and project-based totals (--projects), with output in table format using tablewriter.

## User Requirements

### Primary User Story
As a time tracker user, I want to view basic statistics about my tracked time so that I can understand my productivity patterns and project contributions.

### Acceptance Scenarios
1. **Given** the user has tracked time entries in the past week, **When** they run the stats command, **Then** the system displays daily totals for the past 7 days.
2. **Given** the user has tracked time entries in the past month, **When** they run the stats command, **Then** the system displays weekly totals for the past 4 weeks.
3. **Given** the user has tracked time across multiple projects, **When** they run the stats command, **Then** the system displays totals by project for the past week.

### Functional Requirements
- **FR-001**: System MUST display daily time totals for the past 7 days using YYYY-MM-DD format for dates and HH:MM format for times
- **FR-002**: System MUST display weekly time totals for the past 4 weeks
- **FR-003**: System MUST display time totals grouped by project for the past week (sorted by total time descending)
- **FR-004**: System MUST handle cases where no data exists by displaying "No data available"
- **FR-005**: System MUST calculate totals based on tracked time entries, including currently running entries
- **FR-006**: System MUST prevent combining --weekly and --projects flags

## Clarifications
- Date format: YYYY-MM-DD
- Time format: HH:MM
- No data message: "No data available"
- Weekly starts on Monday
- Project totals sorted by time descending

## Technical Implementation

### Architecture
- **Language**: Go
- **Framework**: Cobra CLI
- **Storage**: JSON file format
- **Output**: Table format using github.com/olekukonko/tablewriter
- **Testing**: Go test with _test.go suffix

### Key Components
- `cmd/stats.go`: Command implementation with flag handling
- `utils/stats_calculations.go`: Calculation functions for daily, weekly, project totals
- `utils/file_storage.go`: Data loading from JSON

### Data Model
- **TimeEntry**: ID, Start, End, Project, Title
- **Project**: Name (string)

### Command Flags
- `--weekly`: Show weekly totals instead of daily
- `--projects`: Group by project instead of time periods

### Error Handling
- Combining --weekly and --projects: Error "cannot combine --weekly and --projects flags"
- No data: Display "No data available"

## Key Decisions
- Use Go time package for date calculations (standard library)
- Load all entries into memory for aggregation (simple for JSON storage)
- Table output with borders for consistency with list command
- Include running entries in calculations
- Sort project totals by time descending
- Remove weekly project combination to simplify

## Validation
- Unit tests for calculation functions
- Contract test for command output format
- Integration tests for command execution
- All tests pass with docker compose