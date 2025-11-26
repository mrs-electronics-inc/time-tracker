# Just file for common development tasks

# Run all Go tests
test:
    cd src && go test ./...

# Run time-tracker with any subcommand and flags
run *args:
    docker compose run --remove-orphans time-tracker {{ args }}

# View the data file from the volume (for debugging)
inspect-data:
    docker run --rm -v time-tracker_config:/mnt alpine cat /mnt/time-tracker/data.json | jq .
