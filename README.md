# Time Tracker ⏲️👣

Time Tracker is a simple CLI tool to track the time you spend on different projects and tasks throughout the day. It stores time entries in a JSON file and ensures only one entry can be active at a time.

> [!NOTE]
> This project is based on [LeanMendez/time-tracker](https://github.com/LeanMendez/time-tracker). The codebase has been reworked to use a time entry data model instead of tasks.

> [!WARNING]
> This project is a work in progress. Things may break.

## Install

### Requirements

To install this application, you will need **Go 1.25.0 or higher**.
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

### Start

To start tracking time for a project and task:

```bash
time-tracker start <project> <title>
```

This will start a new time entry. If another entry is currently running, it will be automatically stopped first.

Example:
```bash
time-tracker start "my-project" "Working on feature"
```

Output: Started tracking time for "Working on feature" in project "my-project"

### Stop

To stop the currently running time entry:

```bash
time-tracker stop
```

This stops the active entry and shows the duration.

Example output: Stopped tracking time for "Working on feature" in project "my-project" (duration: 1h 30m)

If no entry is running, it will show an error.

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
- Review all changes yourself before requesting review from a human
- Refine based on code review feedback

## License

All the code is under [MIT](/LICENSE)
