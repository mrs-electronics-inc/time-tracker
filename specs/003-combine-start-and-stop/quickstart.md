# Quickstart: combine-start-and-stop

## Installation
Ensure Docker Compose is set up as per project guidelines.

## Usage Examples

### Start tracking
```bash
docker compose run time-tracker start "My Project" "Task 1"
# Output: Started tracking: My Project - Task 1

docker compose run time-tracker s "My Project" "Task 1"
# Output: Started tracking: My Project - Task 1
```

### Stop tracking
```bash
docker compose run time-tracker stop
# Output: Stopped tracking: My Project - Task 1 (1h 30m)

docker compose run time-tracker s
# Output: Stopped tracking: My Project - Task 1 (1h 30m)
```

### Error cases
```bash
docker compose run time-tracker start
# Error: Missing required arguments: project and task

docker compose run time-tracker stop "arg"
# Error: stop command does not accept arguments

docker compose run time-tracker s "only one"
# Error: Missing task argument
```

## Testing
Run the integration tests through Docker Compose to verify functionality.