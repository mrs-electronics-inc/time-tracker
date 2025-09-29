---
description: Compact the LLM-generated specs.
---

The user input can be provided directly by the agent or as a command argument - you **MUST** consider it before proceeding with the prompt (if not empty).

User input:

$ARGUMENTS

The text the user typed after `/compact-spec` in the triggering message **is** the subdirectory of `specs/` that needs compacted. If the user does not specify a valid subdirectory of `specs/`, please exit early and display an error message and suggest running `/compact-spec` again.

Compact the files in the given subdirectory of `specs/`.

The files were useful for helping an LLM agent implement the given change, but now we need to remove all the boilerplate.

Compact all the files into a single `spec.md` file that summarizes the change and includes any important technical details or decisions made.

Delete the other files in the `specs/` subdirectory and replace with the new `spec.md` file.
