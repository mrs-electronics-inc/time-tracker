---
description: Start a new pull request by checking specs and creating the appropriate branch
agent: build
---

Automates the process of starting work on a new pull request.

The user gave the input: "$ARGUMENTS"

Use the user input as the spec number.

If the user input is empty or invalid, prompt the user for the spec number.

Required behavior and confirmation flow

1. Read the spec from the `specs/` directory and determine the next incomplete section from the Task List.
2. Branch creation rules (agents MUST NOT ask the user about branch behavior):
   - Compute branch as `{spec-number}-{slug(section-header)}` where `slug()` lowercases the header, replaces any non‑alphanumeric sequence with `-`, collapses duplicate `-`, and trims leading/trailing `-`.
   - If current branch equals OR is very similar to the computed name, do nothing; otherwise create and switch with `git checkout -b "<branch>"`.
   - If creating/switching would overwrite uncommitted work, warn and request confirmation.
3. Research the codebase to gather information about the change.
4. Ask the user clarifying questions.
   - Clearly number the questions.
   - Clearly letter the options for each question.
5. Update the Task List section with any new updates based on your research and the user's answers.
6. Explain the current Task List section to the user.
7. When the Task List section is approved by the user, instruct the user to run `/do-pr` to begin implementing the changes. Do not use the word “proceed” as the final prompt — always reference `/do-pr`.
