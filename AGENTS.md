# Agent Guidelines

## Local Development

All testing must be done through Docker Compose ONLY to ensure safe execution without affecting the host system.

- Build the image: `docker compose build`
- Run commands: `docker compose run --remove-orphans time-tracker [args]`
- Example: `docker compose run --remove-orphans time-tracker start "test-project" "test-task"`

**IMPORTANT**: Never run the binary directly, use `go run`, or execute the project in any way that affects the host system. Always use Docker Compose for testing.
