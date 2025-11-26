# Just file for common development tasks

# Run all Go tests
test:
    cd src && go test ./...

# Run time-tracker with any subcommand and flags
run *args:
    docker compose run --remove-orphans time-tracker {{ args }}
