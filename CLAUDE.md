# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Forge is a developer utility toolkit written in Go. It provides encoding/decoding, formatting, hashing, and generation tools across four surfaces: CLI, TUI (terminal UI), Web, and Desktop. Think "DevToys but in Go, for the terminal first."

**Module path:** `github.com/StanMarek/forge`

## Build & Development Commands

```bash
# Initialize (first time)
go mod init github.com/StanMarek/forge
go mod tidy

# Build
go build -o bin/forge .

# Run
go run .                    # launches TUI by default
go run . base64 encode "hello"  # CLI mode

# Test
go test ./...               # all tests
go test ./core/tools/...    # core tools only
go test -run TestBase64 ./core/tools/  # single test
go test -v -count=1 ./...   # verbose, no cache

# Generate templ templates (when web UI exists)
templ generate

# Lint
golangci-lint run ./...

# Build with version info
go build -ldflags "-X github.com/StanMarek/forge/internal/version.Version=v0.1.0 -X github.com/StanMarek/forge/internal/version.Commit=$(git rev-parse --short HEAD) -X github.com/StanMarek/forge/internal/version.Date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o bin/forge .
```

## Architecture

### Hard Rule: Dependency Direction

```
cmd/ ──► core/     (allowed)
cmd/ ──► ui/       (allowed)
cmd/ ──► internal/ (allowed)
ui/  ──► core/     (allowed)
ui/  ──► internal/clipboard/ (allowed, TUI only)
core/ ──► ui/      (FORBIDDEN — never, ever)
core/ ──► cmd/     (FORBIDDEN — never, ever)
```

`core/` must never import from `ui/` or `cmd/`. This is non-negotiable. Core functions return plain Go types (`string`, `struct`), never formatted terminal output or HTML.

### Package Layout

- **`main.go`** — Entrypoint, calls `cmd.Execute()`. ~5 lines.
- **`cmd/`** — Cobra CLI commands. One file per tool (~30-60 lines each). `root.go` registers all subcommands. Default action (no subcommand) launches TUI.
- **`core/tools/`** — Pure business logic. Flat package, one file + one test file per tool. Tool functions are stateless: `func Base64Encode(input string, urlSafe bool) Base64Result`. No I/O, no global state.
- **`core/tools/tool.go`** — `Tool` interface (metadata only: `Name()`, `ID()`, `Description()`, `Category()`, `Keywords()`, `DetectFromClipboard()`). Tool logic lives in standalone functions, NOT on the interface.
- **`core/registry/`** — Tool registry: `Register()`, `All()`, `ByID()`, `ByCategory()`, `Search()`, `Detect()`. `defaults.go` creates registry with all tools pre-registered.
- **`core/detection/`** — Smart clipboard detection engine.
- **`ui/tui/`** — Bubbletea v2 TUI. `app.go` (root model), `sidebar.go`, `banner.go`, `styles/`, `keys/`, `views/` (one tea.Model per tool).
- **`ui/web/`** — Chi + templ + HTMX web surface. `server.go`, `handlers/` (one per tool), `templates/` (.templ files), `static/` (CSS, vendored HTMX, JS).
- **`ui/desktop/`** — Wails v2 desktop app. Deferred to post-v1, stubs only.
- **`internal/clipboard/`** — Cross-platform clipboard abstraction (atotto/clipboard).
- **`internal/version/`** — Build version info (set via ldflags).
- **`internal/stdin/`** — Stdin reading with timeout, pipe vs terminal detection.

### Tool Result Convention

Every tool function returns a result struct with `Output string` and `Error string`. Success is `Error == ""`. Each UI layer renders results independently — TUI uses lipgloss, web uses templ templates, CLI uses plain text.

### Adding a New Tool

1. Implement `Tool` interface in `core/tools/<toolname>.go`
2. Write pure logic functions in the same file
3. Write tests in `core/tools/<toolname>_test.go`
4. Register in `core/registry/defaults.go`
5. Add cobra command in `cmd/<toolname>.go`
6. Add TUI view in `ui/tui/views/<toolname>.go`
7. Add web handler in `ui/web/handlers/<toolname>.go` + template in `ui/web/templates/tools/<toolname>.templ`

### CLI Conventions

- If `[input]` is provided as an arg, use it directly. If omitted or `-`, read from stdin.
- Output to stdout, errors to stderr.
- Exit codes: 0 success, 1 input error, 2 usage error.
- No color in pipe mode (isatty check). `--color` flag to force.

## Gotchas

- A security hook blocks `innerHTML` usage in HTML/JS files. Use `textContent` for text and `document.createElement()` for DOM construction.

## Tech Stack

| Layer | Library | Import Path |
|-------|---------|-------------|
| TUI | bubbletea v2 | `charm.land/bubbletea/v2` |
| TUI styling | lipgloss v2 | `charm.land/lipgloss/v2` |
| TUI components | bubbles v2 | `charm.land/bubbles/v2` |
| Web templates | templ v0.3.x | `github.com/a-h/templ` |
| Web router | chi v5 | `github.com/go-chi/chi/v5` |
| Web interactivity | HTMX 2.x | Vendored JS file |
| CLI | cobra v1.8.x | `github.com/spf13/cobra` |
| Desktop (deferred) | Wails v2 | `github.com/wailsapp/wails/v2` |
| Clipboard | atotto/clipboard | `github.com/atotto/clipboard` |
| UUID | google/uuid | `github.com/google/uuid` |
| Testing | testify | `github.com/stretchr/testify` |
| Release | goreleaser v2 | `.goreleaser.yaml` |

**Go version:** 1.22+ minimum.

## Tier 1 Tools (must ship first)

base64, jwt, json, hash, url, uuid — implement these in `core/tools/` with full tests before touching any UI code.

## Smart Detection Priority (clipboard)

JWT > UUID > URL > JSON > Base64 > Timestamp > HTML entity > Number base (most specific first).

## Design Theme & Mockups

The UI uses the **Material-Darker** color scheme. Key tokens:
- Background: `#212121`, Surface: `#292929`, Contrast: `#1A1A1A`
- Text: `#EEFFFF` (primary), `#B0BEC5` (secondary), `#616161` (muted)
- Accent: `#FF9800` (orange), Semantic: `#C3E88D` (green), `#89DDFF` (cyan), `#FF5370` (red)

HTML mockups in `docs/mockups/`:
- `mockup-darker-index.html` — Tool grid homepage
- `mockup-darker-tool-base64.html` — Standalone tool page
- `mockup-darker-tool-sidebar.html` — Tool page with left sidebar navigation
- `mockup-darker-desktop.html` — macOS desktop app with translucent sidebar (backdrop-filter vibrancy)

Figma file: `c3SImq2TRDuP7vcYfdVu6Q` (all mockups captured)
To serve mockups locally: `python3 -m http.server 8923` from `docs/mockups/`.

## Design Documents

The `docs/` directory contains the full design specifications:
- `doc-00` — Implementation roadmap and open questions
- `doc-03` — Project structure (annotated directory tree)
- `doc-04` — Tool inventory with detection rules
- `doc-05` / `doc-06` — TUI and web UI mockups
- `doc-07` — CLI command reference (all flags and output formats)
- `doc-08` — Architecture Decision Records (ADRs)
