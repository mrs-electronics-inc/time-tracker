---
status: draft
author: Addison Emig
creation_date: 2026-01-06
---

# Web Interface

Add a `--mode web` option to the `serve` command that provides an HTTP-based web interface for time tracking.

## Problem

The current TUI requires a terminal. A web interface would enable:
- Access from any device with a browser
- Easier mobile use
- Sharing/viewing across machines

## Solution

Extend `time-tracker serve` with `--mode web` that starts an HTTP server serving a web-based interface.

```bash
time-tracker serve --mode web                              # Start web server (default port, localhost only)
time-tracker serve --mode web --host 0.0.0.0 --port 8080   # Listen on all interfaces, custom port
```

## Design Decisions

### Default Port

TBD - Need to decide on a default port for the web server.

## Task List

### TBD

- [ ] TBD
