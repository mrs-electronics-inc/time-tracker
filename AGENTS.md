# Agent Guidelines

## Local Development

You MUST use the `just` tool for all local development tasks:

```bash
# Run all Go tests
just test

# Build the Docker image
just build

# Run dev time-tracker with any subcommand and flags
just run start "project-name" "task-name"
just run stop
just run list
just run list --all
just run stats
just run stats --weekly
just run stats --rows 7

# View the dev data file from the volume (for debugging)
just inspect-data

# Edit the dev data file from the volume with your EDITOR
just edit-data

# Import JSON data from stdin into the volume (OVERWRITES existing data)
just import-data < data.json
```

See `justfile` in the repo root for all available recipes.

**IMPORTANT**: Never run the binary directly on the host system. Always use `just run` for CLI testing.

## GitHub

- **Getting issue descriptions**: Use the GitHub CLI: `gh issue view <number>`
- **Creating PRs**: Use the GitHub CLI: `gh pr create --title "..." --body "..."`
- **PR titles MUST follow conventional commit format** (e.g., `feat:`, `fix:`, `refactor:`, `docs:`, etc.). Since PRs are squashed on merge to main, the PR title becomes the commit message.

## Spec Editing Safety

- Rule: Spec files under `specs/` are long-term design documents. Do NOT record ephemeral or per-session choices (e.g., "user chose 1B") directly inside `specs/` files.

- Rule: Before editing any `specs/` file the agent MUST ask for confirmation. The prompt should state the exact file path and the change summary. Example prompt:
  - I plan to update `specs/001-new-data-format` to change the 'Blank entries representation' line to 'decision pending'. Reply 'yes' to apply.

- Rule: After receiving approval to edit a `specs/` file, the agent MUST present the staged files and a one-line commit message for explicit confirmation BEFORE committing. Do not proceed to commit without this second confirmation.

- Rule: The agent MUST NOT commit changes to `specs/` files without explicit user approval. If a commit is requested, the agent should present the staged files and a one-line commit message for confirmation.

- Rule: When in doubt about whether something is a transient implementation choice or a long-term spec decision, ask the user.
