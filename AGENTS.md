# Agent Guidelines

## Local Development

- **Go tests**: Run directly with `go test ./...`
- **Binary testing**: Use Docker Compose to test the CLI binary to avoid affecting the host system.
  - Build the image: `docker compose build`
  - Run commands: `docker compose run --remove-orphans time-tracker [args]`
  - Example: `docker compose run --remove-orphans time-tracker start "test-project" "test-task"`

**IMPORTANT**: Never run the binary directly on the host system. Always use Docker Compose for CLI testing.

## Spec Editing Safety

- Rule: Spec files under `specs/` are long-term design documents. Do NOT record ephemeral or per-session choices (e.g., "user chose 1B") directly inside `specs/` files.

- Rule: Before editing any `specs/` file the agent MUST ask for confirmation. The prompt should state the exact file path and the change summary. Example prompt:
  - I plan to update `specs/001-new-data-format` to change the 'Blank entries representation' line to 'decision pending'. Reply 'yes' to apply.

- Rule: The agent MUST NOT commit changes to `specs/` files without explicit user approval. If a commit is requested, the agent should present the staged files and a one-line commit message for confirmation.

- Rule: When in doubt about whether something is a transient implementation choice or a long-term spec decision, ask the user.
