# Contributing to time-tracker

## Setup

### With Nix

If you have Nix installed, run:

```bash
nix shell .#default
```

This drops you into a shell with all development dependencies pre-installed, including Go, just, and Docker. After making changes, rebuild the CLI with `just build` to test them.

### Without Nix

Install dependencies and set up pre-commit hooks:

```bash
pip install pre-commit
pre-commit install
```

This ensures code passes all checks (formatting, linting, tests) before committing.

## Development

Use `just` to run development tasks. Run `just --list` to see available recipes.

### Common tasks

```bash
# Run all Go tests
just test

# Build the Docker image
just build

# Run the CLI in Docker sandbox (isolated from your actual data)
just run-dev start "project-name" "task-name"
just run-dev stop
just run-dev list
just run-dev stats
```

### Running against your actual data

After building with `just build`, you can run the compiled binary directly to interact with your actual data file:

```bash
./time-tracker start "project-name" "task-name"
```

**Note:** `just run-dev` runs the CLI in a Docker sandbox with isolated test data, while running the compiled binary directly has direct access to your user's actual data file.

See `justfile` in the repo root for all available recipes.

## GitHub Issues

GitHub issues should be used for bug reports only. Feature requests and refactor requests should be contributed by adding a new file to the `specs/` directory.

## Specs

See [specs/README.md](specs/README.md) for complete guidelines on contributing specs using the Specture System.
