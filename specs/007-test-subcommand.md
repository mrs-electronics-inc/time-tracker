---
status: approved
author: Addison Emig
creation_date: 2026-01-06
approved_by: Bennett Moore
approval_date: 2026-01-07
---

# Test Subcommand

Add a `test` subcommand that allows AI agents and automated tests to interact with the TUI programmatically. This enables automated testing and verification of TUI behavior without requiring a real terminal.

## Problem

Currently, testing the TUI requires either:

- Manual interaction with a real terminal
- Unit tests that call model methods directly (which don't verify rendered output)

This makes E2E testing difficult, and AI agents cannot easily verify what the user actually sees because:

- PTY-based approaches require complex terminal emulation
- ANSI escape sequences are difficult to parse and verify
- Structured text output loses styling information (colors indicate state)

## Solution

Run `time-tracker test` to start test mode, which:

1. Accepts commands via JSON on stdin
2. Renders the TUI to a PNG image after each command
3. Returns the image file path via JSON on stdout

This enables:

- **E2E testing**: Automated tests can verify the actual rendered output
- **AI agent interaction**: Agents can use vision capabilities to verify the TUI
- **Visual regression testing**: Compare screenshots across versions

## Protocol

### Input Commands

```json
{"cmd": "key", "key": "j"}
{"cmd": "key", "key": "enter"}
{"cmd": "key", "key": "esc"}
{"cmd": "key", "key": "up"}
{"cmd": "key", "key": "down"}
{"cmd": "key", "key": "tab"}
{"cmd": "type", "text": "hello world"}
{"cmd": "resize", "rows": 24, "cols": 80}
```

### Output Response

```json
{ "render_path": "/tmp/time-tracker/renders/2026-01-06T10-45-32.123.png" }
```

Renders are saved to `/tmp/time-tracker/renders` by default. Use `--render-dir /custom/path` to specify a different directory.

Files are timestamped for easy sorting and debugging.

When time-tracker exits (or is killed), it cleans up all renders it created. Use `--keep-renders` to persist them for debugging and tracing.

### Initial State

```bash
time-tracker test                              # test mode with default render path
time-tracker test --render-dir /custom/path    # custom render directory
time-tracker test --keep-renders               # persist renders after exit
```

On startup, test mode will:

1. Initialize with default terminal size (80x24)
2. Load existing data (same as normal TUI mode)
3. Send an initial response with the rendered screen

## Design Decisions

### Output Format: PNG Image

We considered several output formats:

| Format                     | Pros                        | Cons                               |
| -------------------------- | --------------------------- | ---------------------------------- |
| Plain text (ANSI stripped) | Simple                      | Loses color/styling information    |
| Raw ANSI codes             | Lossless                    | Difficult to parse and verify      |
| Structured cell grid       | Precise                     | Verbose, hard for AI to process    |
| Base64-encoded PNG         | Self-contained              | Requires decoding, verbose in JSON |
| **PNG file with path**     | AI agents can read directly | Requires filesystem                |

**Decision**: PNG files with paths returned in JSON. AI agents can use the file path directly with image reading tools. Simpler than base64 decoding, and files are useful for debugging.

### Render Directory

The `--render-dir` flag specifies where to save renders. Renders are always saved in test mode. The `--keep-renders` flag disables cleanup on exit.

### Rendering Approach: Use Existing Dependencies

We considered:

- External binary (e.g., `textimg`) - adds deployment complexity
- Inline rendering with existing deps - ~200-300 lines of new code

**Decision**: Implement rendering using existing dependencies:

- `charmbracelet/x/ansi` for parsing ANSI sequences (already a dependency)
- `golang.org/x/image/font` for text rendering
- Embed a suitable open-source monospace font (MIT/OFL licensed)
- Standard library `image/png` for encoding

### Protocol Format: JSON Lines

Simple JSON objects over stdin/stdout, one per line. Easy to parse, widely supported.

## Task List

### Foundation

- [x] Add `test` subcommand to root command
- [x] Add `--render-dir` flag for custom render path
- [x] Add `--keep-renders` flag to disable cleanup
- [x] Add tests for JSON protocol parsing
- [x] Implement JSON protocol parsing (stdin reader)
- [x] Add tests for JSON response writing
- [x] Implement JSON response writing (stdout)

### Input Handling

- [x] Add tests for key command conversion
- [x] Convert `key` commands to `tea.KeyMsg`
- [x] Add tests for type command conversion
- [x] Convert `type` commands to sequence of `tea.KeyMsg`
- [x] Add tests for resize command handling
- [x] Handle `resize` commands via `tea.WindowSizeMsg`

### Rendering

- [ ] Embed a monospace font (e.g., JetBrains Mono, Source Code Pro)
- [x] Add tests for ANSI sequence parsing
- [x] Implement ANSI sequence parser to extract text and styles
- [x] Add tests for image rendering
- [x] Implement image renderer (text + colors to PNG)
- [x] Add tests for render file output with timestamped filenames
- [x] Implement render file output (default: /tmp/time-tracker/renders)
- [x] Support custom path via --render-dir flag

### Integration

- [x] Add integration tests for test subcommand
- [x] Wire up TUI model to test mode loop
- [x] Send initial rendered state on startup
- [x] Add error handling for invalid commands
- [x] Implement cleanup of temp renders on exit (signal handling)
- [x] Implement `--keep-renders` flag to disable cleanup

### Documentation

- [ ] Document test subcommand in README
- [ ] Add example usage for AI agents
- [ ] Update AGENTS.md to instruct agents to use `time-tracker test` for E2E testing of new features


