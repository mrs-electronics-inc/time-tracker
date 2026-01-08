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
    #!/usr/bin/env bash
    if [[ "$1" == "headless" ]]; then
        docker run --rm -p 8484:8484 -v time-tracker_config:/root/.config time-tracker headless "${@:2}"
    else
        docker compose run --remove-orphans time-tracker {{ args }}
    fi

# View the dev data file from the volume (for debugging)
inspect-data:
    docker run --rm -v time-tracker_config:/mnt alpine cat /mnt/time-tracker/data.json | jq .

# Import JSON data from stdin into the volume (OVERWRITES existing data)
import-data:
    docker run --rm -v time-tracker_config:/mnt -i alpine tee /mnt/time-tracker/data.json > /dev/null

# Send input to headless server (action: key, type, resize)
input action *args:
    #!/usr/bin/env bash
    case "{{ action }}" in
        key)
            curl -s -X POST localhost:8484/input -d "{\"action\": \"key\", \"key\": \"$1\"}" | jq .
            ;;
        type)
            curl -s -X POST localhost:8484/input -d "{\"action\": \"type\", \"text\": \"$1\"}" | jq .
            ;;
        resize)
            curl -s -X POST localhost:8484/input -d "{\"action\": \"resize\", \"rows\": $1, \"cols\": $2}" | jq .
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
