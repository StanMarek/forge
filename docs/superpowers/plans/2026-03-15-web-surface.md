# Web Surface Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement a self-hostable web UI with sidebar navigation, HTMX live processing, and all 6 Tier-1 tools.

**Architecture:** Chi router serves templ-rendered HTML pages. Each tool has a GET handler (full page) and POST handler (HTMX partial for output). Static assets (CSS, HTMX, JS) embedded via go:embed. Modern dark theme.

**Tech Stack:** `github.com/go-chi/chi/v5`, `github.com/a-h/templ`, HTMX 2.x (vendored), existing `core/tools/`

**Spec:** `docs/superpowers/specs/2026-03-15-web-surface-design.md`

**Templ binary:** `~/go/bin/templ` (v0.3.1001). Run `~/go/bin/templ generate` after writing .templ files.

---

## Task Dependency Summary

```
Task 1 (deps + static assets) ─┐
Task 2 (CSS)                   ─┤
Task 3 (templ layout+index)    ─┤── Task 9 (server.go) → Task 10 (cmd/web.go) → Task 11 (smoke test)
Tasks 4-8 (tool templates      ─┤
  + handlers, parallelizable)  ─┘
```

Tasks 4-8 (one per tool + index) are independent and parallelizable.

---

## Chunk 1: Foundation — deps, static assets, CSS, layout template

### Task 1: Add chi dependency + create static assets

- [ ] **Step 1: Add chi**

```bash
go get github.com/go-chi/chi/v5
go get github.com/a-h/templ
go mod tidy
```

- [ ] **Step 2: Download HTMX 2.x**

```bash
mkdir -p ui/web/static
curl -L https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js -o ui/web/static/htmx.min.js
```

- [ ] **Step 3: Create forge.js**

Create `ui/web/static/forge.js`:

```js
function copyOutput() {
    var output = document.querySelector('#output textarea');
    if (output) {
        navigator.clipboard.writeText(output.value);
        var btn = document.querySelector('[data-copy-btn]');
        if (btn) {
            var original = btn.textContent;
            btn.textContent = 'Copied!';
            setTimeout(function() { btn.textContent = original; }, 1500);
        }
    }
}
```

- [ ] **Step 4: Commit**

```bash
git add go.mod go.sum ui/web/static/
git commit -m "Add chi, templ deps and static assets (HTMX, forge.js)"
```

### Task 2: CSS stylesheet

Create `ui/web/static/style.css` — the complete modern dark theme.

This is the full CSS file (~300 lines) covering: reset, variables, header, sidebar, sidebar mobile, main content, tool forms, mode pills, checkboxes, textareas, output areas, error states, tool grid cards, footer, utility classes.

- [ ] **Step 1: Write style.css**

- [ ] **Step 2: Commit**

```bash
git add ui/web/static/style.css
git commit -m "Add web UI stylesheet with modern dark theme"
```

### Task 3: Layout + index templates

- [ ] **Step 1: Create directory structure**

```bash
mkdir -p ui/web/templates/tools ui/web/handlers
```

- [ ] **Step 2: Write layout.templ**

`ui/web/templates/layout.templ` — base HTML with header, sidebar, main slot, footer. Parameters: `title string`, `activeTool string`, `contents templ.Component`.

The sidebar is built from a hardcoded tool list (not registry — templates don't import Go packages at runtime, and the list is static for v1).

- [ ] **Step 3: Write index.templ**

`ui/web/templates/index.templ` — tool grid homepage with cards linking to `/tools/{id}`.

- [ ] **Step 4: Generate templ**

```bash
~/go/bin/templ generate
```

- [ ] **Step 5: Verify compilation**

```bash
go build ./ui/web/...
```

- [ ] **Step 6: Commit**

```bash
git add ui/web/templates/
git commit -m "Add layout and index templ templates"
```

---

## Chunk 2: Tool templates + handlers (parallelizable)

Each tool needs TWO files:
1. `ui/web/templates/tools/{tool}.templ` — page component (GET) + output fragment (POST)
2. `ui/web/handlers/{tool}.go` — GET and POST handlers

Plus the index handler.

### Task 4: Index handler

Create `ui/web/handlers/index.go`:

```go
package handlers

import (
    "net/http"
    "github.com/StanMarek/forge/ui/web/templates"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
    templates.IndexPage().Render(r.Context(), w)
}
```

### Task 5: Base64 (template + handler)

**Template** `ui/web/templates/tools/base64.templ`:
- `Base64Page(result, mode, urlSafe, noPadding)` — full form with HTMX
- `Base64Output(result)` — output fragment for HTMX swap

**Handler** `ui/web/handlers/base64.go`:
- GET: render page with empty result
- POST: parse form (input, mode, url-safe, no-padding), call core, render output fragment

### Task 6: JWT + JSON (templates + handlers)

Same pattern. JWT handler uses `tools.JWTDecodeResult`. JSON handler parses indent with `strconv.Atoi`.

### Task 7: Hash + URL (templates + handlers)

Hash handler calls `tools.Hash()` 4 times. URL handler returns `URLParseResult` for parse mode.

### Task 8: UUID (template + handler)

UUID handler: generate mode creates UUID without input, validate/parse modes use input. Returns `UUIDParseResult` for parse.

After all tools: run `~/go/bin/templ generate` and `go build ./ui/web/...`

---

## Chunk 3: Server + command + smoke test

### Task 9: server.go

`ui/web/server.go` — Chi router setup, go:embed for static files, route registration, `Start(host, port)` function.

### Task 10: cmd/web.go

Cobra command `forge web --port 8080 --host localhost`. Register in root.go.

### Task 11: Build + smoke test

Build binary, launch `forge web`, test all 6 tools in browser.
