# Agent Guidelines

## Local Development

You MUST use the `just` tool for all local development tasks:

```bash
# Run all Go tests
just test

# Build the binary
just build

# Run time-tracker in Docker sandbox with any subcommand and flags
just run-docker start "project-name" "task-name"
just run-docker stop
just run-docker list
just run-docker list --all
just run-docker edit
just run-docker stats
just run-docker stats --weekly
just run-docker stats --rows 7

# Build the Docker image
just build-docker

# View the dev data file from the volume (for debugging)
just inspect-data

# Import JSON data from stdin into the volume (OVERWRITES existing data)
# Always use the latest data version from models/migration_types.go
just import-data < data.json
```

See `justfile` in the repo root for all available recipes.

**IMPORTANT**: Agents must always run time-tracker using `just run-docker`. This runs the CLI in a Docker sandbox with isolated test data. Never run the binary directly on the host system.

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

This project uses the [Specture System](https://github.com/specture-system/specture) for managing specifications and design documents. Each spec file (specs/NNN-name.md) contains:

- Design rationale and decisions (why)
- Task lists for implementation (what)
- Requirements and acceptance criteria (how)

### When to refer to specs

When the user asks about planned features, architectural decisions, or implementation details, refer to the specs/ directory. Specs provide the complete context for understanding the project's design.

### Implementing specs

When implementing a spec, follow this workflow for each task:

1. Complete a single task from the task list
2. Update the spec file by changing `- [ ]` to `- [x]` for that task
3. Commit both the implementation and spec update together with a conventional commit message (e.g., `feat: implement feature X`)
4. Push the changes

This keeps the spec file as a living document that tracks implementation progress, with each task corresponding to one commit.

### Spec editing safety

- Rule: Spec files under `specs/` are long-term design documents. Do NOT record ephemeral or per-session choices (e.g., "user chose 1B") directly inside `specs/` files.

- Rule: Before editing any `specs/` file the agent MUST ask for confirmation. State the exact file path and the change summary. Example prompt:
  - I plan to update `specs/001-new-data-format` to change the 'Blank entries representation' line to 'decision pending'. Reply 'yes' to apply.

- Rule: After receiving approval to edit a `specs/` file, present the staged files and a one-line commit message for explicit confirmation BEFORE committing. Do not proceed without this second confirmation.

- Rule: When in doubt about whether something is a transient implementation choice or a long-term spec decision, ask the user.

See specs/README.md for complete guidelines on the Specture System.
