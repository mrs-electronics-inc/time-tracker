---
status: draft
author: Addison Emig
creation_date: 2026-01-08
---

# Headless Mode

Add a `headless` subcommand that runs the TUI as an HTTP server, enabling AI agents and automated tests to interact programmatically.

## Solution

```bash
time-tracker headless                  # Start on localhost:8484
time-tracker headless --port 9000      # Custom port
time-tracker headless --bind 0.0.0.0   # Expose to network (use with caution)
```

## HTTP API

### `POST /input`

Send an action, receive updated state.

```json
// Request
{"action": "key", "key": "j"}
{"action": "key", "key": "enter"}
{"action": "key", "key": "ctrl+c"}
{"action": "type", "text": "hello world"}
{"action": "resize", "rows": 24, "cols": 80}

// Response
{
  "width": 160,
  "height": 40,
  "mode": "list",
  "render_url": "/render/2026-01-08T10-45-32-123.png",
  "ansi": "\u001b[1;92mStart..."
}
```

**Key format:** Use `tea.KeyMsg.String()` format: `enter`, `esc`, `tab`, `up`, `down`, `shift+tab`, `ctrl+c`, etc.

**Mode values:** From `CurrentMode.Name`: `list`, `start`, `help`, `stats`, etc.

**Timestamp format:** `2026-01-08T10-45-32-123.png` (ISO 8601 with dashes for URL safety, millisecond precision)

The `ansi` field contains raw output from `View()` with all escape sequences.

### `GET /render/latest`

Redirects (302) to most recent render.

### `GET /render/{timestamp}.png`

Returns specific render PNG.

### `GET /state`

Returns current state (same format as `POST /input` response):

```json
{
  "width": 160,
  "height": 40,
  "mode": "list",
  "render_url": "/render/2026-01-08T10-45-32.123.png",
  "ansi": "\u001b[1;92mStart..."
}
```

## Configuration

| Flag | Default | Description |
|------|---------|-------------|
| `--bind` | 127.0.0.1 | Bind address (localhost only by default) |
| `--port` | 8484 | Port number |
| `--max-renders` | 100 | Max renders to keep in memory (FIFO eviction) |

Default terminal size: 160×40

## Security

**This server is for local development only.**

- Binds to localhost by default
- No authentication
- No CORS restrictions
- May expose sensitive time tracking data

Never expose to untrusted networks.

## Usage

```bash
# Start server
just run-docker headless

# Send input via just recipe
just input key j
just input type "hello world"
just input resize 40 160

# Or use curl directly
curl -X POST localhost:8484/input -d '{"action": "key", "key": "j"}'

# View in browser
open http://localhost:8484/render/latest
```

## Design Decisions

- **HTTP vs stdin/stdout**: HTTP allows direct browser access to renders, no volume mounts needed
- **ANSI + PNG**: Response includes both for flexibility (vision vs text assertions)
- **160×40 default**: Large enough for content, small enough to catch layout issues
- **FiraCode Nerd Font**: Required for powerline symbols (regular Fira Code lacks these glyphs)
- **Color palette**: ANSI 16-color palette (see below), default/background color is pure black (`#000000`)

## Color Palette

Default/background color: `#000000` (pure black)

ANSI 16-color palette:

| Index | Name | Hex |
|-------|------|-----|
| 0 | black | `#1D1F21` |
| 1 | red | `#CC6666` |
| 2 | green | `#B5BD68` |
| 3 | yellow | `#F0C674` |
| 4 | blue | `#81A2BE` |
| 5 | magenta | `#B294BB` |
| 6 | cyan | `#8ABEB7` |
| 7 | white | `#C5C8C6` |
| 8 | bright black | `#666666` |
| 9 | bright red | `#D54E53` |
| 10 | bright green | `#B9CA4A` |
| 11 | bright yellow | `#E7C547` |
| 12 | bright blue | `#7AA6DA` |
| 13 | bright magenta | `#C397D8` |
| 14 | bright cyan | `#70C0B1` |
| 15 | bright white | `#EAEAEA` |

## Task List

### HTTP Server Foundation

- [ ] Add `headless` subcommand with HTTP server
- [ ] Add `--port` flag (default: 8484)
- [ ] Add `--bind` flag (default: 127.0.0.1)
- [ ] Add `--max-renders` flag (default: 100)
- [ ] Implement `POST /input` endpoint
- [ ] Implement `GET /render/latest` redirect endpoint
- [ ] Implement `GET /render/{timestamp}.png` endpoint
- [ ] Return JSON error responses for invalid actions
- [ ] Add tests for HTTP endpoints

### Input Handling

- [ ] Convert `key` actions to `tea.KeyMsg`
- [ ] Convert `type` actions to sequence of `tea.KeyMsg`
- [ ] Handle `resize` actions via `tea.WindowSizeMsg`
- [ ] Add tests for action conversion

### Rendering

- [ ] Embed FiraCode Nerd Font (Regular and Bold) for powerline symbol support
- [ ] Implement ANSI sequence parser
- [ ] Implement PNG renderer with color palette from spec
- [ ] Store renders in memory with timestamp keys
- [ ] Implement render eviction when max limit reached (FIFO)
- [ ] Add tests for rendering

### Response Format

- [ ] Include `render_url` in POST response
- [ ] Include `ansi` (raw ANSI string) in POST response
- [ ] Implement `GET /state` endpoint with render_url and ansi

### Integration

- [ ] Configure bubbletea to render ANSI codes to non-tty output
- [ ] Wire up TUI model to HTTP handlers
- [ ] Set default terminal size to 160×40
- [ ] Create initial render on server startup (so /render/latest works immediately)
- [ ] Update `run-docker` recipe to bind port 8484 for headless subcommand
- [ ] Add `input` recipe that wraps curl for sending actions (key, type, resize)

### Documentation

- [ ] Document headless mode in README
- [ ] Document security considerations (localhost binding, no auth, local dev only)
- [ ] Add example usage for AI agents
- [ ] Update AGENTS.md with headless server workflow and new recipes (`just run-docker headless`, `just input`)
