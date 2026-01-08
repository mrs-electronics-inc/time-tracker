---
status: draft
author: Addison Emig
creation_date: 2026-01-06
---

# Headless Subcommand

Add a `headless` subcommand that allows AI agents and automated tests to interact with the TUI programmatically. This enables automated testing and verification of TUI behavior without requiring a real terminal.

## Problem

Currently, testing the TUI requires either:

- Manual interaction with a real terminal
- Unit tests that call model methods directly (which don't verify rendered output)

This makes E2E testing difficult, and AI agents cannot easily verify what the user actually sees because:

- PTY-based approaches require complex terminal emulation
- ANSI escape sequences are difficult to parse and verify
- Structured text output loses styling information (colors indicate state)

## Solution

Run `time-tracker headless` to start an HTTP server that accepts commands and serves rendered screenshots:

```bash
time-tracker headless              # Start on default port 8080
time-tracker headless --port 9000  # Start on custom port
```

This enables:

- **E2E testing**: Automated tests can verify the actual rendered output
- **AI agent interaction**: Agents can use vision capabilities to verify the TUI
- **Visual regression testing**: Compare screenshots across versions
- **Easy debugging**: View renders directly in a browser

## HTTP API

### POST /input - Send Commands

Send a command to the TUI and receive the updated state.

**Request:**
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

**Response:**
```json
{
  "render_url": "/render/2026-01-06T10-45-32.123.png",
  "ansi": "\u001b[1;92mStart             End..."
}
```

### GET /render/latest - Redirect to Current Screen

Redirects (HTTP 302) to the most recent timestamped render. Convenient for quick viewing - just refresh to see the latest state.

### GET /render/{timestamp}.png - Specific Render

Returns a specific render by timestamp. The `render_url` in POST responses points here.

### GET /state - Current State

Returns the current TUI state including ANSI output and link to latest render.

```json
{
  "width": 120,
  "height": 30,
  "mode": "list",
  "render_url": "/render/2026-01-06T10-45-32.123.png",
  "ansi": "\u001b[1;92mStart             End..."
}
```

## Default Configuration

- **Port**: 8080
- **Terminal size**: 120 columns × 30 rows (larger default for better visibility)
- **Render cleanup**: Renders are kept for the lifetime of the server

## Usage

### Starting the Headless Server

```bash
# Via just recipe (recommended for development)
just run-docker headless

# Direct invocation (without Docker)
time-tracker headless --port 8080
```

### Interacting with the Server

```bash
# Send a key command
curl -X POST http://localhost:8080/input \
  -H "Content-Type: application/json" \
  -d '{"cmd": "key", "key": "j"}'

# View latest render in browser (redirects to timestamped URL)
open http://localhost:8080/render/latest

# Resize terminal
curl -X POST http://localhost:8080/input \
  -d '{"cmd": "resize", "rows": 40, "cols": 160}'
```

### AI Agent Workflow

1. Start headless server: `just run-docker headless`
2. Send commands via POST /input
3. View renders via browser at /render/latest
4. Use ANSI output from response for text-based assertions

## Design Decisions

### HTTP Server vs stdin/stdout

We considered two approaches:

| Approach | Pros | Cons |
|----------|------|------|
| stdin/stdout JSON | Simple, no network | Requires volume mounts for images, buffering issues |
| **HTTP server** | Direct image access, browser viewable, stateless | Requires port allocation |

**Decision**: HTTP server. Benefits:
- AI agents can directly navigate to render URLs in browser
- No need for filesystem volume mounts in Docker
- Easy manual debugging via browser
- curl/httpie for scripting
- ANSI + PNG in single response

### Response Includes Both ANSI and PNG URL

The POST /input response includes both:
- `render_url`: For visual verification via vision capabilities
- `ansi`: For text-based assertions and searching

This allows agents to choose the best approach for each verification.

### Larger Default Terminal Size

Default size is 120×30 (vs typical 80×24) because:
- Modern displays can show more content
- Status bars and wide tables render better
- AI vision works better with larger, clearer images

### Render Retention

Renders are kept for the server's lifetime (no cleanup). This enables:
- Debugging by reviewing history
- Visual regression comparisons
- No risk of deleting renders still being viewed

### Font Choice: Fira Code

Embed Fira Code (OFL licensed) because:
- Excellent Unicode coverage including powerline symbols
- Clear rendering at various sizes
- Popular, well-maintained open source font
- Includes bold variant for proper bold rendering

## Task List

### HTTP Server Foundation

- [ ] Add `headless` subcommand with HTTP server
- [ ] Add `--port` flag (default: 8080)
- [ ] Implement POST /input endpoint
- [ ] Implement GET /render/latest redirect endpoint
- [ ] Implement GET /render/{timestamp}.png endpoint
- [ ] Add tests for HTTP endpoints

### Input Handling

- [ ] Convert `key` commands to `tea.KeyMsg`
- [ ] Convert `type` commands to sequence of `tea.KeyMsg`
- [ ] Handle `resize` commands via `tea.WindowSizeMsg`
- [ ] Add tests for command conversion

### Rendering

- [ ] Embed Fira Code Regular and Bold fonts
- [ ] Implement ANSI sequence parser
- [ ] Implement PNG renderer with Ghostty color palette
- [ ] Store renders in memory with timestamp keys
- [ ] Add tests for rendering

### Response Format

- [ ] Include `render_url` in POST response
- [ ] Include `ansi` (raw ANSI string) in POST response
- [ ] Implement GET /state endpoint with render_url and ansi

### Integration

- [ ] Wire up TUI model to HTTP handlers
- [ ] Set default terminal size to 120×30
- [ ] Send initial render on first /render/latest.png request
- [ ] Update `run-docker` recipe to bind port 8080 for headless subcommand

### Documentation

- [ ] Document headless subcommand in README
- [ ] Add example usage for AI agents
- [ ] Update AGENTS.md with headless server workflow
