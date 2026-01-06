---
status: draft
author: Addison Emig
creation_date: 2026-01-06
---

# Serve Command with JSON Mode for TUI Testing

Add a `serve` subcommand with a `--mode json` option that allows AI agents to interact with the TUI programmatically. This enables automated testing and verification of TUI behavior without requiring a real terminal.

## Problem

Currently, testing the TUI requires either:
- Manual interaction with a real terminal
- Unit tests that call model methods directly (which don't verify rendered output)

AI agents cannot easily verify what the user actually sees because:
- PTY-based approaches require complex terminal emulation
- ANSI escape sequences are difficult to parse and verify
- Structured text output loses styling information (colors indicate state)

## Solution

Add a `time-tracker serve --mode json` command that:
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
{"image": "iVBORw0KGgoAAAANSUhEUgAA..."}
```

The `image` field contains a base64-encoded PNG of the rendered TUI.

### Initial State

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
| **PNG image** | AI vision processes it naturally | Larger payload |

**Decision**: PNG image. AI agents are much better at processing image data directly than parsing structured representations. This also tests what the user actually sees.

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

### Infrastructure

- [ ] Add `cmd/serve.go` with cobra command structure and `--mode` flag
- [ ] Implement JSON protocol parsing (stdin reader)
- [ ] Implement JSON response writing (stdout)

### Input Handling

- [ ] Convert `key` commands to `tea.KeyMsg`
- [ ] Convert `type` commands to sequence of `tea.KeyMsg`
- [ ] Handle `resize` commands via `tea.WindowSizeMsg`

### Rendering

- [ ] Embed a monospace font (e.g., JetBrains Mono, Source Code Pro)
- [ ] Implement ANSI sequence parser to extract text and styles
- [ ] Implement image renderer (text + colors to PNG)
- [ ] Base64 encode PNG for JSON response

### Integration

- [ ] Wire up TUI model to serve command loop
- [ ] Send initial rendered state on startup
- [ ] Add error handling for invalid commands

### Documentation

- [ ] Document serve command in README
- [ ] Add example usage for AI agents

## Future Work

The `serve` command is designed to support multiple modes:

- `--mode json` - JSON over stdin/stdout for AI agents (this spec)
- `--mode web` - HTTP server with web interface (future spec)
