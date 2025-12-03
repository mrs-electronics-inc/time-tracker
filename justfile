# Just file for common development tasks

# Run all Go tests
test:
    go test ./...

# Build the Docker image
build:
    docker compose build

# Run dev time-tracker with any subcommand and flags
run *args:
    docker compose run --remove-orphans time-tracker {{ args }}

# View the dev data file from the volume (for debugging)
inspect-data:
    docker run --rm -v time-tracker_config:/mnt alpine cat /mnt/time-tracker/data.json | jq .

# Import JSON data from stdin into the volume (OVERWRITES existing data)
import-data:
    docker run --rm -v time-tracker_config:/mnt -i alpine tee /mnt/time-tracker/data.json > /dev/null

