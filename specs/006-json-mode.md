---
status: draft
author: Addison Emig
creation_date: 2026-01-06
---

# JSON Mode for TUI Testing

Add a `--mode json` flag that allows AI agents to interact with the TUI programmatically. This enables automated testing and verification of TUI behavior without requiring a real terminal.

## Problem

Currently, testing the TUI requires either:
- Manual interaction with a real terminal
- Unit tests that call model methods directly (which don't verify rendered output)

AI agents cannot easily verify what the user actually sees because:
- PTY-based approaches require complex terminal emulation
- ANSI escape sequences are difficult to parse and verify
- Structured text output loses styling information (colors indicate state)

## Solution

Run `time-tracker --mode json` to start JSON mode, which:
1. Accepts commands via JSON on stdin
2. Renders the TUI to a PNG image after each command
3. Returns the image (base64 encoded) via JSON on stdout

This allows AI agents to use vision capabilities to verify the actual rendered TUI, including colors, layout, and styling.

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
{"image_path": "/tmp/time-tracker/screens/2026-01-06T10-45-32.123.png"}
```

Images are saved to `/tmp/time-tracker/screens` by default. Use `--save-images /custom/path` to specify a different directory.

Files are timestamped for easy sorting and debugging.

When the serve command exits (or is killed), it cleans up all images it created. Use `--keep-images` to persist images for debugging and tracing.

### Initial State

```bash
time-tracker --mode json                              # JSON mode with default image path
time-tracker --mode json --save-images /custom/path   # custom image directory
time-tracker --mode json --keep-images                # persist images after exit
```

On startup, JSON mode will:
1. Initialize with default terminal size (80x24)
2. Load existing data (same as normal TUI mode)
3. Send an initial response with the rendered screen

## Design Decisions

### Output Format: PNG Image

We considered several output formats:

| Format | Pros | Cons |
|--------|------|------|
| Plain text (ANSI stripped) | Simple | Loses color/styling information |
| Raw ANSI codes | Lossless | Difficult to parse and verify |
| Structured cell grid | Precise | Verbose, hard for AI to process |
| Base64-encoded PNG | Self-contained | Requires decoding, verbose in JSON |
| **PNG file with path** | AI agents can read directly | Requires filesystem |

**Decision**: PNG files with paths returned in JSON. AI agents can use the file path directly with image reading tools. Simpler than base64 decoding, and files are useful for debugging.

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

- [ ] Add `--mode` flag to root command (default: tui)
- [ ] Add `--save-images` flag for custom image path
- [ ] Add `--keep-images` flag to disable cleanup
- [ ] Add tests for JSON protocol parsing
- [ ] Implement JSON protocol parsing (stdin reader)
- [ ] Add tests for JSON response writing
- [ ] Implement JSON response writing (stdout)

### Input Handling

- [ ] Add tests for key command conversion
- [ ] Convert `key` commands to `tea.KeyMsg`
- [ ] Add tests for type command conversion
- [ ] Convert `type` commands to sequence of `tea.KeyMsg`
- [ ] Add tests for resize command handling
- [ ] Handle `resize` commands via `tea.WindowSizeMsg`

### Rendering

- [ ] Embed a monospace font (e.g., JetBrains Mono, Source Code Pro)
- [ ] Add tests for ANSI sequence parsing
- [ ] Implement ANSI sequence parser to extract text and styles
- [ ] Add tests for image rendering
- [ ] Implement image renderer (text + colors to PNG)
- [ ] Add tests for image file output with timestamped filenames
- [ ] Implement image file output (default: /tmp/time-tracker/screens)
- [ ] Support custom path via --save-images flag

### Integration

- [ ] Add integration tests for serve command
- [ ] Wire up TUI model to JSON mode loop
- [ ] Send initial rendered state on startup
- [ ] Add error handling for invalid commands
- [ ] Implement cleanup of temp images on exit (signal handling)
- [ ] Implement `--keep-images` flag to disable cleanup

### Documentation

- [ ] Document serve command in README
- [ ] Add example usage for AI agents

## Future Work

The `serve` command is designed to support multiple modes:

- `--mode json` - JSON over stdin/stdout for AI agents (this spec)
- `--mode web` - HTTP server with web interface (see [spec 007](./007-web-mode.md))
