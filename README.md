# Task-timer

Task-timer is a simple CLI tool to manage your tasks and track their time.
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
git clone https://github.com/LeanMendez/time-tracker.git
```
2. Navigate to the repository directory:
```bash
cd time-tracker
```
From here, you have two options to compile and install the application:

#### Option A: Install with Go
Use the go install command to compile and install the binary automatically in your *$GOPATH/bin* directory:
```bash
go install
```

#### Option B: Build with Task
Use [Task](https://taskfile.dev/) to compile the binary. Make sure you have Task installed before proceeding.
```bash
task build

```
This command will generate the binary at *./bin/time-tracker*.

If you want the binary to be accessible globally, you’ll need to add its path to your system’s PATH environment variable.

## Usage

### Init command
To start using the CLI you need to initialize the application by running 
```bash
timer-cli init [path]
```
This will create a *tasks.json* file in the specified path where all your tasks will be stored.

### Create command
Then you can create tasks by running 
```bash
timer-cli create [task name]
```
Creating a task doens't mean it is started.
You can also start a task by running the `start` command or when creating it by using the `--start` or `-s` flag.
```bash
timer-cli create [task name] --start
```

### Start command
To start a task use 
```bash
timer-cli start [task name]
```

or you can use the task ID
```bash
timer-cli start [task ID]
```


### List command
To list your tasks use 
```bash
timer-cli list
```
It will display all your tasks in a table format.
To list a specific task use 
```bash
timer-cli list [task name]
```

or you can use the task ID
```bash
timer-cli list [task ID]
```

### Remove command
To remove a task use 
```bash
timer-cli remove [task name]
```

or you can use the task ID
```bash
timer-cli remove [task ID]
```

### Stop command
To stop a task use 
```bash
timer-cli stop [task name]
```
or you can use the task ID
```bash
timer-cli stop [task ID]
```
Stopping a task marks it as completed.


## Tech Stack

- [Go](https://go.dev/)
- [Cobra](https://github.com/spf13/cobra)

## License

All the code is under [MIT](https://github.com/LeanMendez/time-tracker/blob/main/LICENSE)
