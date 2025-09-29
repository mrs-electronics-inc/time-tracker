# Time Tracker ⏲️👣

Time Tracker is a simple CLI tool to track the time you spend on different projects and tasks throughout the day. It stores time entries in a JSON file and ensures only one entry can be active at a time.

> [!NOTE]
> This project is based on [LeanMendez/time-tracker](https://github.com/LeanMendez/time-tracker). The codebase has been reworked to use a time entry data model instead of tasks.

> [!WARNING]
> This project is a work in progress. Things may break.

## Install

### Requirements

To install this application, you will need **Go 1.23.0 or higher**.
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

2. Navigate to the src directory:

```bash
cd time-tracker/src
```

3. Use the go install command to compile and install the binary automatically in your _$GOPATH/bin_ directory:

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

Output: Started tracking time for "Working on feature" in project "my-project"

### Stop Tracking

To stop the currently running time entry:

```bash
time-tracker stop
# or
time-tracker s
```

This stops the active entry and shows the duration.

Example output: Stopped tracking time for "Working on feature" in project "my-project" (duration: 1h 30m)

If no entry is running, it will show an error.

Note: The `stop` command does not accept arguments. Use `start` to begin tracking.

### List

To list all time entries:

```bash
time-tracker list
```

Displays all entries in chronological order (newest first) with ID, start time, end time (or "running"), project, title, and duration.

## Tech Stack

- [Go](https://go.dev/)
- [Cobra](https://github.com/spf13/cobra)

## Development

### Testing

To run the automated tests from the repository root:

```bash
go test ./src/tests/...
```

This will run contract tests, integration tests, and unit tests.

### AI-driven Workflow

- Install [spec-kit](https://github.com/github/spec-kit) - `uv tool install specify-cli --from git+https://github.com/github/spec-kit.git`
- Install [opencode.ai](https://opencode.ai)
- Login with openrouter - `opencode auth login`
- Select model in `opencode` with the `/models` command (currently recommended: [Grok Code Fast 1](https://openrouter.ai/x-ai/grok-code-fast-1)).
- Use the `/specify` command to describe what you want to build. ([docs](https://github.com/github/spec-kit?tab=readme-ov-file#3-create-the-spec))
- Take a look at the output. Refine as needed.
- Commit
- Use the `/clarify` command to clarify the design.
- Take a look at the output. Refine as needed.
- Commit
- Use the `/plan` command to describe any architecture choices. ([docs](https://github.com/github/spec-kit?tab=readme-ov-file#3-create-the-spec))
- Take a look at the output. Refine as needed.
- Commit
- Use the `/tasks` command to create the task list.
- Take a look at the output. Refine as needed.
- Commit
- Use the `/implement` command to execute the task list.
- Commit
- Push
- Create PR
- Review all changes yourself
- Refine as needed
- Use the `/compact-spec <subdirectory>` command to compact the spec files in `specs/<subdirectory>` into a single `spec.md`, removing boilerplate.
- Commit
- Push
- Request human review
- Refine based on code review feedback

## License

All the code is under the [MIT license](/LICENSE). Contributions are welcome!
