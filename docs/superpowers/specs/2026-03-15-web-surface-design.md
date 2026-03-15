# Web Surface — Design Spec

**Date:** 2026-03-15
**Status:** Approved
**Scope:** Chi router, templ templates, HTMX live processing, sidebar layout, 6 Tier-1 tools, `forge web` command

---

## Goal

Implement a self-hostable web UI for Forge. Launched via `forge web`, it serves a dark-themed web app with sidebar navigation and live tool processing via HTMX. All processing happens server-side (localhost) — nothing sent to external services.

## Dependencies

- `github.com/go-chi/chi/v5` — HTTP router
- `github.com/a-h/templ` — Type-safe HTML templates (compile-time)
- HTMX 2.x — Vendored JS file (~14KB), no npm
- `github.com/StanMarek/forge/core/tools` — Business logic
- `github.com/StanMarek/forge/core/registry` — Tool list for sidebar

## Architecture

```
forge web --port 8080
  → chi.Router
    GET  /              → index handler → layout + index template
    GET  /tools/{id}    → tool page handler → layout + sidebar + tool template
    POST /tools/{id}    → HTMX handler → tool output fragment only
    GET  /static/*      → http.FileServer (embedded via go:embed)
    GET  /health        → 200 OK JSON
```

### Request Flow (HTMX)

1. User types in input textarea
2. `hx-trigger="input changed delay:200ms"` fires
3. HTMX POSTs form data to `/tools/{id}`
4. Handler parses form, calls `core/tools.Function()`
5. Handler renders templ output fragment (just the `<div id="output">` content)
6. HTMX swaps `#output` with the response
7. No full page reload — only the output area updates

### Static Assets

Embedded in the Go binary via `go:embed`. No external CDN, no npm, no build step for frontend.

```go
//go:embed static/*
var staticFS embed.FS
```

## Visual Style — Modern Dark

```css
:root {
    --bg:           #0a0a0a;
    --surface:      #171717;
    --surface-hover: #1f1f1f;
    --border:       #262626;
    --border-light: #303030;
    --text:         #fafafa;
    --text-secondary: #a1a1aa;
    --text-muted:   #525252;
    --accent:       #f97316;
    --accent-hover: #fb923c;
    --accent-dim:   rgba(249, 115, 22, 0.1);
    --green:        #86efac;
    --red:          #f87171;
    --cyan:         #67e8f9;
    --yellow:       #fde68a;
    --font-mono:    'JetBrains Mono', 'Fira Code', 'Cascadia Code', ui-monospace, monospace;
    --font-sans:    'Inter', system-ui, -apple-system, sans-serif;
    --radius:       8px;
    --radius-sm:    6px;
    --sidebar-w:    240px;
}
```

### Typography

- Body: `var(--font-sans)`, 14px, color `var(--text)`
- Code/textareas: `var(--font-mono)`, 13px
- Headers: `var(--font-sans)`, bold
- Google Fonts: Inter + JetBrains Mono (loaded in `<head>`)

## Layout

### Header

- Fixed top bar, `var(--bg)` background, bottom border
- Left: "⚒ Forge" logo (anvil emoji + monospace text, accent color on anvil)
- Right: version info in muted text
- Height: ~48px

### Sidebar (Desktop)

- Fixed left, 240px wide, full height below header
- Background: `var(--surface)`
- Right border: `var(--border)`
- Tool list grouped by category (ENCODERS, FORMATTERS, GENERATORS)
- Category headers: uppercase, muted, small font
- Tool items: padding, hover state (`var(--surface-hover)`)
- Active tool: accent left border + accent text + `var(--accent-dim)` background
- Scrollable if tools overflow

### Sidebar (Mobile <768px)

- Hidden by default
- Hamburger button in header toggles visibility
- Opens as full-width overlay
- CSS-only toggle (checkbox hack or `:target`), no JS needed

### Main Content Area

- Left margin: `var(--sidebar-w)` on desktop, 0 on mobile
- Padding: 24px
- Max-width: none (fills available space)

### Footer

- Inside main content area, bottom
- Muted text: "forge {version} · All processing happens locally"

## Templates

### layout.templ

Base HTML layout wrapping all pages:

```
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{title} — Forge</title>
  <link href="/static/style.css" rel="stylesheet">
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
  <script src="/static/htmx.min.js"></script>
</head>
<body>
  <header>...</header>
  <aside class="sidebar">...</aside>
  <main>{children}</main>
  <script src="/static/forge.js"></script>
</body>
</html>
```

Parameters: `title string`, `activeTool string`, `children templ.Component`

The sidebar is part of the layout, rendered on every page. `activeTool` highlights the current tool.

### index.templ

Tool grid on the homepage. Cards grouped by category, linking to `/tools/{id}`.

```
<div class="tool-grid">
  for each category:
    <h2>CATEGORY</h2>
    <div class="grid">
      for each tool:
        <a href="/tools/{id}" class="tool-card">
          <h3>{name}</h3>
          <p>{description}</p>
        </a>
    </div>
</div>
```

### Tool templates (one per tool)

Each tool template has TWO components:

1. **Page component** — the full form (rendered on GET)
2. **Output component** — just the output area (rendered on POST, swapped by HTMX)

Example for base64:

```
// Page component — rendered on GET /tools/base64
templ Base64Page(result tools.Result, mode string, urlSafe bool) {
    <h1>Base64 Encode / Decode</h1>
    <form hx-post="/tools/base64" hx-target="#output" hx-trigger="input changed delay:200ms, change">
        // Mode pills
        <div class="mode-pills">
            <label class={modeClass("encode", mode)}>
                <input type="radio" name="mode" value="encode" checked?={mode=="encode"}> Encode
            </label>
            <label class={modeClass("decode", mode)}>
                <input type="radio" name="mode" value="decode" checked?={mode=="decode"}> Decode
            </label>
        </div>
        // Options
        <label class="checkbox">
            <input type="checkbox" name="url-safe" checked?={urlSafe}> URL-safe
        </label>
        // Input
        <label>Input</label>
        <textarea name="input" rows="6" placeholder="Enter text..."></textarea>
        // Output (HTMX target)
        <div id="output">
            @Base64Output(result)
        </div>
    </form>
}

// Output component — rendered on POST (HTMX partial)
templ Base64Output(result tools.Result) {
    if result.Error != "" {
        <div class="error">{result.Error}</div>
    } else {
        <label>Output</label>
        <textarea readonly rows="6">{result.Output}</textarea>
    }
}
```

## Handlers

### Pattern

Every tool handler follows the same pattern:

```go
// GET /tools/base64 — full page
func HandleBase64Page(w http.ResponseWriter, r *http.Request) {
    result := tools.Result{}
    mode := "encode"
    templates.Base64Page(result, mode, false).Render(r.Context(), w)
}

// POST /tools/base64 — HTMX partial
func HandleBase64Process(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    input := r.FormValue("input")
    mode := r.FormValue("mode")
    urlSafe := r.FormValue("url-safe") == "on"

    var result tools.Result
    if mode == "decode" {
        result = tools.Base64Decode(input, urlSafe)
    } else {
        result = tools.Base64Encode(input, urlSafe, false)
    }
    templates.Base64Output(result).Render(r.Context(), w)
}
```

GET renders the full page (layout + sidebar + tool form + empty output).
POST renders just the output fragment (HTMX swaps it into `#output`).

### server.go

```go
func NewRouter() chi.Router {
    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    // Static files
    r.Handle("/static/*", http.StripPrefix("/static/",
        http.FileServerFS(staticFS)))

    // Health
    r.Get("/health", handleHealth)

    // Index
    r.Get("/", handleIndex)

    // Tools
    r.Get("/tools/base64", handlers.HandleBase64Page)
    r.Post("/tools/base64", handlers.HandleBase64Process)
    // ... same for jwt, json, hash, url, uuid

    return r
}

func Start(host string, port int) error {
    r := NewRouter()
    addr := fmt.Sprintf("%s:%d", host, port)
    fmt.Printf("Forge web server running at http://%s\n", addr)
    return http.ListenAndServe(addr, r)
}
```

### Tool-specific handler notes

**hash:** POST receives `input` string. Returns all 4 hash digests in the output fragment (same as TUI — md5, sha1, sha256, sha512 stacked).

**jwt:** POST receives `token` string. Returns decoded header + payload + signature.

**json:** POST receives `input`, `mode` (format/minify/validate), `sort-keys`, `indent`. Returns formatted/minified output or validation result.

**url:** POST receives `input`, `mode` (parse/encode/decode), `component`. Returns result.

**uuid:** POST for generate doesn't need input — generates and returns. POST for validate/parse receives `input`.

## Static Files

### style.css (~250 lines)

Single stylesheet. CSS variables for theming. No framework, no Tailwind.

Key sections:
- Reset + base styles
- Header
- Sidebar + mobile hamburger
- Main content
- Tool form (mode pills, checkboxes, textareas, output)
- Tool grid (index page cards)
- Error states
- Responsive breakpoints

### htmx.min.js

Vendored HTMX 2.x (~14KB gzipped). Downloaded once and committed.

### forge.js (~20 lines)

Clipboard helpers using the browser Clipboard API:

```js
function copyOutput() {
    const output = document.querySelector('#output textarea');
    if (output) navigator.clipboard.writeText(output.value);
}
```

Uses `document.createElement()` and `textContent` for DOM manipulation (per CLAUDE.md security hook — no `innerHTML`).

## cmd/web.go

```go
var webCmd = &cobra.Command{
    Use:   "web",
    Short: "Launch the web server",
    Run: func(cmd *cobra.Command, args []string) {
        port, _ := cmd.Flags().GetInt("port")
        host, _ := cmd.Flags().GetString("host")
        if err := web.Start(host, port); err != nil {
            fmt.Fprintf(os.Stderr, "error: %v\n", err)
            os.Exit(1)
        }
    },
}

func init() {
    webCmd.Flags().Int("port", 8080, "Port to listen on")
    webCmd.Flags().String("host", "localhost", "Host to bind to")
}
```

Registered in `cmd/root.go`.

## Templ Workflow

1. Write `.templ` files in `ui/web/templates/`
2. Run `templ generate` to produce `_templ.go` files
3. Commit both `.templ` and `_templ.go` to git
4. No templ binary needed at build time — generated Go code is checked in

Install templ: `go install github.com/a-h/templ/cmd/templ@latest`

## File Manifest

```
ui/web/server.go
ui/web/handlers/index.go
ui/web/handlers/base64.go
ui/web/handlers/jwt.go
ui/web/handlers/json.go
ui/web/handlers/hash.go
ui/web/handlers/url.go
ui/web/handlers/uuid.go
ui/web/templates/layout.templ
ui/web/templates/index.templ
ui/web/templates/tools/base64.templ
ui/web/templates/tools/jwt.templ
ui/web/templates/tools/json.templ
ui/web/templates/tools/hash.templ
ui/web/templates/tools/url.templ
ui/web/templates/tools/uuid.templ
ui/web/static/style.css
ui/web/static/htmx.min.js
ui/web/static/forge.js
cmd/web.go
```

Modified: `cmd/root.go` (add webCmd)

## Testing Strategy

- Core logic is tested (200+ tests). Handlers are thin glue.
- Manual testing: `go run . web`, open browser, use each tool.
- Optionally: `httptest` for handler response codes in a later pass.

## Out of Scope

- Light mode / theme switching
- PWA / service worker
- Search on index page
- Websocket-based live updates (HTMX polling is sufficient)
- Authentication / multi-user
