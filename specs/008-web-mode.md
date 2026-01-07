---
status: draft
author: Addison Emig
creation_date: 2026-01-06
---

# Web Mode

Add a `--mode web` flag that provides an HTTP-based web interface for time tracking.

## Problem

The current TUI requires a terminal. A web interface would enable:

- Access from any device with a browser
- Easier mobile use
- Sharing/viewing across machines

## Solution

Add `--mode web` to start an HTTP server serving a web-based interface.

```bash
time-tracker --mode web                              # Start web server (default port, localhost only)
time-tracker --mode web --host 0.0.0.0 --port 8080   # Listen on all interfaces, custom port
```

## Design Decisions

### Default Port

TBD - Need to decide on a default port for the web server.

### Frontend Architecture

Consider using [Elm](https://elm-lang.org/) for the web frontend. Bubble Tea was inspired by Elm's architecture (Model-Update-View), which could allow:

- Reusing the existing TUI logic/model on the backend
- Sending user input from the web UI to the same Update functions
- Rendering to HTML instead of terminal output
- Clean separation: Go handles state/logic, Elm handles web presentation

This would provide convenient mobile-friendly inputs while preserving the existing Bubble Tea architecture.

### Render Saving

The `--save-renders` flag (shared across all modes) is disabled by default for web mode. Can be enabled for debugging, saving HTML snapshots of the rendered UI:

```bash
time-tracker --mode web --save-renders                    # saves HTML to /tmp/time-tracker/renders
time-tracker --mode web --save-renders --render-dir /custom/path   # custom directory
```

## Task List

### TBD

- [ ] TBD
