# Just file for common development tasks

# List recipes
default:
    @just --list

# Run all Go tests
test:
    go test ./...

# Format Go code
fmt:
    gofmt -w .

# Lint Go code
lint:
    go vet ./...

# Build the binary for direct use
build:
    go build -o time-tracker

# Build the Docker image
build-docker:
    docker compose build

# Run time-tracker in Docker sandbox with any subcommand and flags
run-docker *args:
    docker compose run --remove-orphans --rm -p 8484:8484 time-tracker {{ args }}

# View the dev data file from the volume (for debugging)
inspect-data:
    docker run --rm -v time-tracker_config:/mnt alpine cat /mnt/time-tracker/data.json | jq .

# Import JSON data from stdin into the volume (OVERWRITES existing data)
import-data:
    docker run --rm -v time-tracker_config:/mnt -i alpine tee /mnt/time-tracker/data.json > /dev/null

# Send input to headless server (action: key, type, resize)
input action first_arg="" second_arg="":
    #!/usr/bin/env bash
    case "{{ action }}" in
        key)
            curl -s -X POST localhost:8484/input -d "{\"action\": \"key\", \"key\": \"{{ first_arg }}\"}" | jq .
            ;;
        type)
            curl -s -X POST localhost:8484/input -d "{\"action\": \"type\", \"text\": \"{{ first_arg }}\"}" | jq .
            ;;
        resize)
            curl -s -X POST localhost:8484/input -d "{\"action\": \"resize\", \"rows\": {{ first_arg }}, \"cols\": {{ second_arg }}}" | jq .
            ;;
        *)
            echo "Usage: just input <key|type|resize> <args>"
            echo "  just input key j"
            echo "  just input key enter"
            echo "  just input type 'hello world'"
            echo "  just input resize 40 160"
            exit 1
            ;;
    esac
