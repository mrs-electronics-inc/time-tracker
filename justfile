# Just file for common development tasks

# Run all Go tests
test:
    cd src && go test ./...

# Build the Docker image
build:
    docker compose build

# Start a new time entry
start project task:
    docker compose run --remove-orphans time-tracker start "{{ project }}" "{{ task }}"

# Stop the current time entry
stop:
    docker compose run --remove-orphans time-tracker stop

# List time entries for today
list:
    docker compose run --remove-orphans time-tracker list

# List all time entries
list-all:
    docker compose run --remove-orphans time-tracker list --all

# Show daily stats (default)
stats:
    docker compose run --remove-orphans time-tracker stats

# Show weekly stats
stats-weekly:
    docker compose run --remove-orphans time-tracker stats --weekly

# Show stats for N rows
stats-rows rows:
    docker compose run --remove-orphans time-tracker stats --rows {{ rows }}

# Complete workflow: build, start, stop, and list
demo:
    @just build
    @just start "demo-project" "demo-task"
    @echo "Waiting 2 seconds..."
    @sleep 2
    @just stop
    @just list
