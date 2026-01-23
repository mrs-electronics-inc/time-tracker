---
description: Develop the changes for a new pull request
agent: build
---

IMPORTANT: This command must run the full non-interactive flow for creating a PR. That means it MUST run the test suite(s), commit any changes, push the branch, and create the GitHub pull request — all without asking the user for additional input.

The user gave the input: "$ARGUMENTS"

Use the user input as the spec number.

If the user input is empty or invalid, use the previously entered spec number from `/start-pr` (but if `/start-pr` was not previously ran, prompt the user for the spec number).

Required behavior (non-interactive flow)

1. Read the spec for the given number in `specs/` and determine the next incomplete section from the Task List.
2. For each task in the incomplete section:
   - Implement the task.
   - Run the relevant automated tests immediately after implementing the change. Tests must be run and pass before committing.
     - If a change only affects unit tests, run the narrower set of packages to save time.
   - If tests fail, refine the code until tests pass. Do not proceed to committing that TODO item until its tests pass.
   - Once tests pass, update the spec (check off corresponding item) and commit the change locally using a descriptive conventional commit message (example `feat(7): add backup script`).
     - Use: `git add -A && git commit -m "<scope>: <short description>"`
3. After all task items for the current section are completed and committed locally:
   - Push the branch to the remote:
     - `git push -u origin "$(git rev-parse --abbrev-ref HEAD)"`
   - Create the pull request non-interactively using `gh` (GitHub CLI). Provide a clear title and a PR body via a HEREDOC to avoid shell quoting issues. Example:
     - `gh pr create --title "<PR title>" --body "$(cat <<'EOF'\n<PR body>\nEOF\n)"`
   - Ask the user for code review feedback on the new pull request.
   - DO NOT suggest starting on the next section of the task list.

Error handling and constraints

- The command must NOT prompt the user for extra confirmation during the flow. If an operation would normally require input (for example, `gh pr create` in interactive mode), invoke the non-interactive flags and provide the input programmatically (HEREDOC or CLI flags).
- If network push or GH CLI operations fail, surface the error and abort; do not attempt destructive recovery automatically.
