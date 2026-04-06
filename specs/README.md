# Specs

Specs are design documents that describe planned changes to this project. Specs live under `specs/` as directories that contain `SPEC.md` files, and those directories may nest to model larger features. Spec numbers are stored in the YAML frontmatter `number` field. The directory tree is the source of truth for how specs are organized.

When writing specs:

- Do not number section headers (`##`/`###`).
- Use inline markdown links relative to this README for any cross-spec mentions (for example, `[Status command](002-status-command/SPEC.md)`).
- Keep specs focused on goals, rationale, and design decisions.

For full documentation on the spec system, workflow, and file format, see the [Specture System](https://github.com/specture-system/specture) repository. Run `specture help` for CLI usage.

This project uses pull requests for collaboration on specs.
