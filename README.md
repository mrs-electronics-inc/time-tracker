# Time-tracker

Time-tracker is a simple CLI tool to manage your tasks and track their time.
It will create a JSON file where all the data is stored in the location you specify in the init command.

It is a project in progress. Any feedback is welcome.

## Install

### Requirements

To install this application, you will need **Go 1.23.2 or higher**.  
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

2. Navigate to the repository directory:

```bash
cd time-tracker
```

3. Use the go install command to compile and install the binary automatically in your _$GOPATH/bin_ directory:

```bash
go install
```

**Note** - if _$GOPATH_ is not defined on your system, look in the following locations:

- Unix systems: `$HOME/go`
- Windows: `%USERPROFILE%\go`

(It is a good idea to add `GOPATH` to your `PATH`)

## Usage

### Init command

To start using the CLI you need to initialize the application by running

```bash
time-tracker init [path]
```

This will create a `tasks.json` file in the specified path where all your tasks will be stored.

### Create command

Then you can create tasks by running

```bash
time-tracker create [task name]
```

Creating a task does not mean it is started.
You can also start a task by running the `start` command or when creating it by using the `--start` or `-s` flag.

```bash
time-tracker create [task name] --start
```

### Start command

To start a task use

```bash
time-tracker start [task name]
```

or you can use the task ID

```bash
time-tracker start [task ID]
```

### List command

To list your tasks use

```bash
time-tracker list
```

It will display all your tasks in a table format.
To list a specific task use

```bash
time-tracker list [task name]
```

or you can use the task ID

```bash
time-tracker list [task ID]
```

### Remove command

To remove a task use

```bash
time-tracker remove [task name]
```

or you can use the task ID

```bash
time-tracker remove [task ID]
```

### Stop command

To stop a task use

```bash
time-tracker stop [task name]
```

or you can use the task ID

```bash
time-tracker stop [task ID]
```

Stopping a task marks it as completed.

## Tech Stack

- [Go](https://go.dev/)
- [Cobra](https://github.com/spf13/cobra)

## Development

### AI-driven Workflow

1. Install [spec-kit](https://github.com/github/spec-kit) - `uv tool install specify-cli --from git+https://github.com/github/spec-kit.git`
1. Install [opencode.ai](https://opencode.ai)
1. Login with openrouter - `opencode auth login`
1. Select model in `opencode` with the `/models` command (currently recommended: [Grok Code Fast 1](https://openrouter.ai/x-ai/grok-code-fast-1)).
1. Use the `/specify` command to describe what you want to build. ([docs](https://github.com/github/spec-kit?tab=readme-ov-file#3-create-the-spec))
1. Take a look at the output. Refine as needed.
1. Commit
1. Use the `/plan` command to describe any architecture choices. ([docs](https://github.com/github/spec-kit?tab=readme-ov-file#3-create-the-spec))
1. Take a look at the output. Refine as needed.
1. Commit
1. Use the `/tasks` command to create the task list.
1. Take a look at the output. Refine as needed.
1. Commit
1. Use the `/implement` command to execute the task list.
1. Commit
1. Push
1. Create PR
1. Review all changes yourself before requesting review from a human
1. Refine based on code review feedback

## License

All the code is under [MIT](/LICENSE)
