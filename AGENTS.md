# time-tracker Development Guidelines

Auto-generated from all feature plans. Last updated: 2025-09-29

## Active Technologies

- Go + Cobra (001-rework-the-data)
- Docker (002-docker-container)
- GitHub Actions (002-docker-container)
- JSON file format (001-rework-the-data)

## Project Structure

```
src/
tests/
```

## Commands

# Add commands for Go

## Code Style

Go: Follow standard conventions

## Recent Changes

- 003-combine-start-and-stop: Combined start and stop commands into single command with aliases
- 002-docker-container: Added Docker support with CI automation for safe LLM agent testing
- 001-rework-the-data: Added Go + Cobra

<!-- MANUAL ADDITIONS START -->

## Testing Instructions

All testing must be done through Docker Compose ONLY to ensure safe execution without affecting the host system.

- Build the image: `docker compose build`
- Run commands: `docker compose run --remove-orphans time-tracker [args]`
- Example: `docker compose run --remove-orphans time-tracker start "test-project" "test-task"`

**IMPORTANT**: Never run the binary directly, use `go run`, or execute the project in any way that affects the host system. Always use Docker Compose for testing.

<!-- MANUAL ADDITIONS END -->
