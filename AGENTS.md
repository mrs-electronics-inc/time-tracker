# Agent Guidelines

## Local Development

All testing must be done through Docker Compose ONLY to ensure safe execution without affecting the host system.

- Build the image: `docker compose build`
- Run commands: `docker compose run --remove-orphans time-tracker [args]`
- Example: `docker compose run --remove-orphans time-tracker start "test-project" "test-task"`

**IMPORTANT**: Never run the binary directly, use `go run`, or execute the project in any way that affects the host system. Always use Docker Compose for testing.

## Spec Editing Safety

- Rule: Spec files under `specs/` are long-term design documents. Do NOT record ephemeral or per-session choices (e.g., "user chose 1B") directly inside `specs/` files.

- Rule: If the agent needs to capture a transient decision, the default behaviour is to NOT record it. Instead the agent should:
  - Prompt the user for an explicit instruction about where to store the decision (if the user wants it stored).
  - If the user requests recording, create a per-spec notes file `specs/<num>-<slug>.notes.md` or append an entry to `.bots/decisions.md` (only when the user asks).

- Rule: If the agent detects transient language in an existing spec (for example, "user chose 1B"), the agent MUST NOT edit the spec to remove or change the text. Instead it should notify the user and wait for instruction. The user may fix the spec manually or ask the agent to sanitize it.

- Rule: Before editing any `specs/` file the agent MUST ask for confirmation. The prompt should state the exact file path and the change summary. Example prompt:
  - `I plan to update `specs/001-new-data-format` to change the 'Blank entries representation' line to 'decision pending'. Reply 'yes' to apply.`

- Rule: The agent MUST NOT commit changes to `specs/` files without explicit user approval. If a commit is requested, the agent should present the staged files and a one-line commit message for confirmation.

- Rule: Branch creation and switching that may affect uncommitted changes must be confirmed. If uncommitted changes are present, the agent should ask:
  - `Detected uncommitted changes. Create and switch to branch '1-add-blank-time-entries' and include these changes in the new branch? (yes/no)`

- Rule: When in doubt about whether something is a transient implementation choice or a long-term spec decision, ask the user.

### Bad vs Good examples

- Bad (do not put this in a spec file):
  - `- [ ] Blank entries representation choice: serialize project and title as empty strings (user chose 1B)`

- Good (spec file as long-term design doc):
  - `- [ ] Blank entries representation: decision pending — see specs/001-new-data-format.notes.md`

- Good (if user explicitly asked to record decision in a notes file):
  - Create `specs/001-new-data-format.notes.md` containing:
    - `2025-11-14 — Transient decision: blank entries will use empty strings (chosen during session).` 

### Mistake handling

- If a spec file is modified incorrectly (contains transient text):
  - Do not push or merge the change.
  - Notify the user immediately and wait for their instruction to revert, sanitize, or leave as-is.
  - If instructed to sanitize, provide the diff and request confirmation before committing.


