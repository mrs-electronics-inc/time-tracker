---
status: draft
author: Addison Emig
creation_date: 2026-01-06
---

# Headless Mode

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

Run `time-tracker headless` to start an HTTP server that accepts input and serves rendered screenshots:

```bash
time-tracker headless              # Start on default port 8484
time-tracker headless --port 9000  # Start on custom port
```

This enables:

- **E2E testing**: Automated tests can verify the actual rendered output
- **AI agent interaction**: Agents can use vision capabilities to verify the TUI
- **Visual regression testing**: Compare screenshots across versions
- **Easy debugging**: View renders directly in a browser

## HTTP API

### POST /input - Send Actions

Send an action to the TUI and receive the updated state.

**Request:**
```json
{"action": "key", "key": "j"}
{"action": "key", "key": "enter"}
{"action": "key", "key": "esc"}
{"action": "key", "key": "up"}
{"action": "key", "key": "down"}
{"action": "key", "key": "tab"}
{"action": "type", "text": "hello world"}
{"action": "resize", "rows": 24, "cols": 80}
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

- **Port**: 8484
- **Terminal size**: 160 columns × 40 rows (moderate size to catch layout issues on smaller terminals)
- **Render cleanup**: Renders are kept for the lifetime of the server

## Usage

### Starting the Headless Server

```bash
# Via just recipe (recommended for development)
just run-docker headless

# Direct invocation (without Docker)
time-tracker headless
```

### Sending Input

```bash
# Send a key
just input key j
just input key enter
just input key tab

# Type text
just input type "hello world"

# Resize terminal
just input resize 40 160
```

### Interacting with the Server

```bash
# Send a key action
curl -X POST http://localhost:8484/input \
  -H "Content-Type: application/json" \
  -d '{"action": "key", "key": "j"}'

# View latest render in browser (redirects to timestamped URL)
open http://localhost:8484/render/latest

# Resize terminal
curl -X POST http://localhost:8484/input \
  -d '{"action": "resize", "rows": 40, "cols": 160}'
```

### AI Agent Workflow

1. Start headless server: `just run-docker headless`
2. Send input: `just input key j` or `just input type "text"`
3. View renders via browser at http://localhost:8484/render/latest
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
- [ ] Add `--port` flag (default: 8484)
- [ ] Implement POST /input endpoint
- [ ] Implement GET /render/latest redirect endpoint
- [ ] Implement GET /render/{timestamp}.png endpoint
- [ ] Add tests for HTTP endpoints

### Input Handling

- [ ] Convert `key` actions to `tea.KeyMsg`
- [ ] Convert `type` actions to sequence of `tea.KeyMsg`
- [ ] Handle `resize` actions via `tea.WindowSizeMsg`
- [ ] Add tests for action conversion

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
- [ ] Set default terminal size to 160×40
- [ ] Send initial render on first /render/latest.png request
- [ ] Update `run-docker` recipe to bind port 8484 for headless subcommand
- [ ] Add `input` recipe for sending actions (key, type, resize)

### Documentation

- [ ] Document headless mode in README
- [ ] Add example usage for AI agents
- [ ] Update AGENTS.md with headless server workflow and new recipes (`just run-docker headless`, `just input`)
