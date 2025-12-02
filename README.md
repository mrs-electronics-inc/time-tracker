# Time Tracker ⏲️👣

Time Tracker is a simple CLI tool to track the time you spend on different projects and tasks throughout the day. It stores time entries in a JSON file and ensures only one entry can be active at a time.

> [!NOTE]
> This project is based on [LeanMendez/time-tracker](https://github.com/LeanMendez/time-tracker). The codebase has been reworked to use a time entry data model instead of tasks.

> [!WARNING]
> This project is a work in progress. Things may break.

## Install

### Requirements

To install this application, you will need **Go 1.24.0 or higher**.
You can download it from the [official Golang website](https://go.dev/dl/).
To verify your Go installation and version, run the following command in your terminal:

```bash
go version
```

### Steps to Install

1. Clone the repository:

```bash
git clone https://github.com/mrs-electronics-inc/time-tracker.git
```

2. Use the go install command to compile and install the binary automatically in your _$GOPATH/bin_ directory:

```bash
go install
```

**Note** - if _$GOPATH_ is not defined on your system, look in the following locations:

- Unix systems: `$HOME/go`
- Windows: `%USERPROFILE%\go`

(It is a good idea to add `GOPATH` to your `PATH`)

> [!NOTE]
> To make things easier, you can create an alias in your shell for the `time-tracker` command.
> We like to use `t`.

## Usage

The time tracker uses a combined command that can be invoked as `start`, `stop`, or `s` (short alias).

### Start Tracking

To start tracking time for a project and task:

```bash
time-tracker start <project> <title>
# or
time-tracker s <project> <title>
```

To resume tracking by entry ID:

```bash
time-tracker start <ID>
# or
time-tracker s <ID>
```

This will start a new time entry or resume an existing one. If another entry is currently running, it will be automatically stopped first.

Examples:

```bash
time-tracker start "my-project" "Working on feature"
time-tracker start 5  # Resume entry with ID 5
```

### Stop Tracking

To stop the currently running time entry:

```bash
time-tracker stop
# or
time-tracker s
```

This stops the active entry and shows the duration.

If no entry is running, it will show an error.

Note: The `stop` command does not accept arguments. Use `start` to begin tracking a new entry.

### List

To list all time entries:

```bash
time-tracker list
```

Displays all entries in chronological order (newest first) with ID, start time, end time (or "running"), project, title, and duration.

### Stats

To view time tracking statistics:

```bash
time-tracker stats
```

Displays daily totals for the past week in table format, including breakdowns by project.

Options:

- `--weekly`: Show weekly totals for the past month

Examples:

```bash
time-tracker stats  # Daily totals
time-tracker stats --weekly  # Weekly totals
```

## Tech Stack

- [Go](https://go.dev/)
- [Cobra](https://github.com/spf13/cobra)

## Development

### Requirements

To develop this project, you need:
- **Go 1.24.0 or higher**
- **Docker and Docker Compose** (for running the app and tests)
- **just** (optional, but recommended for running common tasks)

### Running Commands

All development commands should be run from the repository root using `just`:

```bash
# Run all tests
just test

# Build the Docker image
just build

# Run the CLI with any subcommand
just run list
just run start "project" "task"
just run stop

# View the dev data file (for debugging)
just inspect-data

# Import JSON data into the dev volume
just import-data < data.json
```

### Testing

To run the automated tests:

```bash
just test
```

This will run contract tests, integration tests, and unit tests.

## License

All the code is under the [MIT license](/LICENSE). Contributions are welcome!
