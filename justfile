# Just file for common development tasks

# Run all Go tests
test:
    cd src && go test ./...

# Build the Docker image
build:
    docker compose build

# Run dev time-tracker with any subcommand and flags
run *args:
    docker compose run --remove-orphans time-tracker {{ args }}

# View the dev data file from the volume (for debugging)
inspect-data:
    docker run --rm -v time-tracker_config:/mnt alpine cat /mnt/time-tracker/data.json | jq .

# Edit the dev data file from the volume with your EDITOR
edit-data:
    #!/usr/bin/env bash
    tmpfile=$(mktemp)
    docker run --rm -v time-tracker_config:/mnt alpine cat /mnt/time-tracker/data.json > "$tmpfile"
    editor=${EDITOR:-nano}
    $editor "$tmpfile"
    docker run --rm -v time-tracker_config:/mnt -i alpine tee /mnt/time-tracker/data.json < "$tmpfile" > /dev/null
    rm "$tmpfile"

# Import JSON data from stdin into the volume (OVERWRITES existing data)
import-data:
    docker run --rm -v time-tracker_config:/mnt -i alpine tee /mnt/time-tracker/data.json > /dev/null

