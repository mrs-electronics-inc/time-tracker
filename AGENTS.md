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

- **Current version**: Check `models/storage.go` for `CurrentVersion` constant (currently v3)
- **JSON structure**: Use `"time-entries"` key (not `"entries"`), with structure:
  ```json
  {
    "version": 3,
    "time-entries": [
      {
        "start": "2025-12-23T09:00:00Z",
        "project": "ProjectName",
        "title": "Task description"
      }
    ]
  }
  ```
- **Blank entries**: Use empty `project` and `title` strings to mark gaps/end of day. Each workday must end with a blank entry at closing time to prevent overnight durations
- **Duration calculation**: End times are derived from the next entry's start time
- **Date selection**: Use recent dates (within last 14 days) for CLI testing since `stats` command shows the last 14 days by default
- **Version handling**: Older versions (v0-v2) are automatically migrated to v3 on load. Always create new test data with current version only

## GitHub

- **Getting issue descriptions**: Use the GitHub CLI: `gh issue view <number>`
- **Creating PRs**: Use the GitHub CLI: `gh pr create --title "..." --body "..."`
- **PR titles MUST follow conventional commit format** (e.g., `feat:`, `fix:`, `refactor:`, `docs:`, etc.). Since PRs are squashed on merge to main, the PR title becomes the commit message.

## Specture System

This project uses the [Specture System](https://github.com/specture-system/specture) for managing specs. See the `.agents/skills/specture/` skill for the full workflow, or run `specture help` for CLI usage.

## Headless Server

For programmatic TUI interaction, use the headless server:

```bash
# Start the headless server in Docker using just
just run-docker headless --bind 0.0.0.0

# Send keyboard input using just
just input key j
just input key enter
just input type hello
just input resize 40 160

# Get current state
curl localhost:8484/state

# Get latest render as PNG
curl -L localhost:8484/render/latest -o screenshot.png
```

The `/state` and `/input` endpoints return JSON with `width`, `height`, `mode`, `render_url`, and `ansi` fields.

### Using Browser Tools (Preferred)

**When browser/screenshot tools are available, agents SHOULD use them to view renders directly:**

```
# Navigate to the render URL
browser_navigate("http://localhost:8484/render/latest")

# Take a screenshot to see the TUI
browser_take_screenshot()
```

This is preferred over downloading PNGs via curl because:
1. Screenshots are immediately visible in the conversation
2. Browser tools handle redirects automatically
3. No need to manage temporary files

### Key format

Use `tea.KeyMsg.String()` format: `enter`, `esc`, `tab`, `up`, `down`, `shift+tab`, `ctrl+c`, etc.
