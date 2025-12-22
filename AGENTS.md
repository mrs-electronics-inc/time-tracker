# Agent Guidelines

## Local Development

You MUST use the `just` tool for all local development tasks:

```bash
# Run all Go tests
just test

# Build the Docker image
just build

# Run dev time-tracker with any subcommand and flags
just run-dev start "project-name" "task-name"
just run-dev stop
just run-dev list
just run-dev list --all
just run-dev edit
just run-dev stats
just run-dev stats --weekly
just run-dev stats --rows 7

# View the dev data file from the volume (for debugging)
just inspect-data

# Import JSON data from stdin into the volume (OVERWRITES existing data)
# Always use the latest data version from models/migration_types.go
just import-data < data.json
```

See `justfile` in the repo root for all available recipes.

**IMPORTANT**: Never run the binary directly on the host system. Always use `just run` for CLI testing.

Vendor directory is gitignored; dependencies are fetched from the network during builds.

**Updating vendorHash**: When `go.mod` changes:
1. Set `vendorHash = "";` in `flake.nix` (empty string)
2. Run `nix build 2>&1 | grep -E "(specified|got):"`
3. Nix will show the correct hash:
   ```
   specified: sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=
      got:    sha256-ZknVM8bMM0kLIbuV4Bv4XsbgtyhlKyP7p2AVOE1k0GA=
   ```
4. Copy the `got:` hash and update `vendorHash` in `flake.nix`
5. Run `nix build` again to verify it works
6. **DO NOT run `go mod vendor`** - the vendor directory should remain empty/deleted

## Data Format

- **Always use latest version**: Check `models/migration_types.go` for the current data format
- **Blank entries**: Use empty `project` and `title` strings to mark gaps/end of day. Each workday must end with a blank entry at closing time to prevent overnight durations
- **Duration calculation**: End times are derived from the next entry's start time

## GitHub

- **Getting issue descriptions**: Use the GitHub CLI: `gh issue view <number>`
- **Creating PRs**: Use the GitHub CLI: `gh pr create --title "..." --body "..."`
- **PR titles MUST follow conventional commit format** (e.g., `feat:`, `fix:`, `refactor:`, `docs:`, etc.). Since PRs are squashed on merge to main, the PR title becomes the commit message.

## Specture System

This project uses the Specture System for managing specifications and design documents. When the user asks about planned features, architectural decisions, or implementation details, refer to the specs/ directory in the repository. Each spec file (specs/NNN-name.md) contains:

- Design rationale and decisions
- Task lists for implementation
- Requirements and acceptance criteria

The specs/ directory also contains README.md with complete guidelines on how the spec system works.

Be sure to prompt the user for explicit permission before editing the design in any spec file.

When implementing a spec, check off each item in the task list as you go.

## Spec Editing Safety

- Rule: Spec files under `specs/` are long-term design documents. Do NOT record ephemeral or per-session choices (e.g., "user chose 1B") directly inside `specs/` files.

- Rule: Before editing any `specs/` file the agent MUST ask for confirmation. The prompt should state the exact file path and the change summary. Example prompt:
  - I plan to update `specs/001-new-data-format` to change the 'Blank entries representation' line to 'decision pending'. Reply 'yes' to apply.

- Rule: After receiving approval to edit a `specs/` file, the agent MUST present the staged files and a one-line commit message for explicit confirmation BEFORE committing. Do not proceed to commit without this second confirmation.

- Rule: The agent MUST NOT commit changes to `specs/` files without explicit user approval. If a commit is requested, the agent should present the staged files and a one-line commit message for confirmation.

- Rule: When in doubt about whether something is a transient implementation choice or a long-term spec decision, ask the user.
