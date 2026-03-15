# Web UI Mockups

**Project:** Forge — Developer Utility Toolkit
**Date:** 2026-03-15

---

## Design Principles

1. **No JavaScript frameworks.** HTMX + a tiny custom JS file (< 1KB) for clipboard and textarea auto-resize.
2. **No CSS frameworks.** A single custom stylesheet (~200 lines). Clean, minimal, dark-first.
3. **No build step for frontend.** CSS and JS are vendored static files. Only templ requires codegen.
4. **Mobile-friendly.** Sidebar collapses to a hamburger menu on narrow viewports.
5. **Every tool is a page.** URL-addressable: `/tools/base64`, `/tools/jwt`, etc. Bookmarkable, shareable.

---

## Mockup 1: Index Page — Tool List Grouped by Category

### Desktop (> 768px)

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│  ⚒ Forge                                                                       │
│  A developer's workbench for the terminal, browser, and desktop.                │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│  ENCODERS                                                                       │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐              │
│  │ Base64           │  │ JWT Decoder      │  │ URL Encode       │              │
│  │ Encode & decode  │  │ Decode & inspect │  │ Encode, decode   │              │
│  │ Base64 strings   │  │ JWT tokens       │  │ & parse URLs     │              │
│  └──────────────────┘  └──────────────────┘  └──────────────────┘              │
│  ┌──────────────────┐                                                           │
│  │ HTML Entity      │                                                           │
│  │ Encode & decode  │                                                           │
│  │ HTML entities    │                                                           │
│  └──────────────────┘                                                           │
│                                                                                 │
│  FORMATTERS                                                                     │
│  ┌──────────────────┐                                                           │
│  │ JSON Formatter   │                                                           │
│  │ Format, minify   │                                                           │
│  │ & validate JSON  │                                                           │
│  └──────────────────┘                                                           │
│                                                                                 │
│  GENERATORS                                                                     │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐              │
│  │ Hash Generator   │  │ UUID Generator   │  │ Password Gen     │              │
│  │ MD5, SHA-1/256   │  │ Generate v4, v7  │  │ Random secure    │              │
│  │ SHA-512          │  │ validate, parse  │  │ passwords        │              │
│  └──────────────────┘  └──────────────────┘  └──────────────────┘              │
│                                                                                 │
│  CONVERTERS                                                                     │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐              │
│  │ JSON / YAML      │  │ Timestamp        │  │ Number Base      │              │
│  │ Convert between  │  │ Unix ↔ datetime  │  │ Dec, hex, oct    │              │
│  │ JSON and YAML    │  │ converter        │  │ bin converter    │              │
│  └──────────────────┘  └──────────────────┘  └──────────────────┘              │
│                                                                                 │
│  TESTERS                                                                        │
│  ┌──────────────────┐                                                           │
│  │ Regex Tester     │                                                           │
│  │ Test patterns    │                                                           │
│  │ with highlights  │                                                           │
│  └──────────────────┘                                                           │
│                                                                                 │
├─────────────────────────────────────────────────────────────────────────────────┤
│  Forge v0.1.0 · All processing happens locally in your browser.                 │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### HTML structure (templ component sketch):

```
layout.templ:
  <html>
    <head> — stylesheet, HTMX script
    <body>
      <header> — logo, tagline
      <main> — content slot
      <footer> — version, privacy note

index.templ:
  for each category:
    <section>
      <h2>CATEGORY NAME</h2>
      <div class="tool-grid">
        for each tool in category:
          <a href="/tools/{id}" class="tool-card">
            <h3>Tool Name</h3>
            <p>Description</p>
          </a>
```

**Behaviour notes:**

- Tool cards are rendered as `<a>` links — full page navigation, no HTMX needed on the index page.
- Cards use CSS Grid: 3 columns on desktop, 2 on tablet, 1 on mobile.
- Cards have a subtle hover effect (border color change to accent violet).
- No search on the index page in v1 — the tool list is short enough to scan visually. Add search in v1.1 if the tool count exceeds ~20.

---

## Mockup 2: Tool Page — Base64 (Desktop)

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│  ⚒ Forge  ·  ◀ All Tools                                                       │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│  Base64 Encode / Decode                                                         │
│                                                                                 │
│  ┌─ Options ──────────────────────────────────────────────────────────────────┐ │
│  │  Mode:  (●) Encode  ( ) Decode          ☐ URL-safe    ☐ No padding       │ │
│  └────────────────────────────────────────────────────────────────────────────┘ │
│                                                                                 │
│  Input                                                                          │
│  ┌────────────────────────────────────────────────────────────────────────────┐ │
│  │                                                                            │ │
│  │ Hello, World!                                                              │ │
│  │                                                                            │ │
│  │                                                                            │ │
│  └────────────────────────────────────────────────────────────────────────────┘ │
│  [Paste from clipboard]                                                         │
│                                                                                 │
│  Output                                                                         │
│  ┌────────────────────────────────────────────────────────────────────────────┐ │
│  │                                                                            │ │
│  │ SGVsbG8sIFdvcmxkIQ==                                                      │ │
│  │                                                                            │ │
│  │                                                                            │ │
│  └────────────────────────────────────────────────────────────────────────────┘ │
│  [Copy to clipboard]    [Clear]                                                 │
│                                                                                 │
├─────────────────────────────────────────────────────────────────────────────────┤
│  Forge v0.1.0 · All processing happens locally. Nothing is sent to a server.    │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### HTMX interaction flow:

```
1. User types in the Input <textarea>
2. On input change (hx-trigger="input changed delay:200ms"):
     hx-post="/tools/base64"
     hx-target="#output-area"
     hx-swap="innerHTML"
3. Server receives form data: { input, mode, urlSafe, noPadding }
4. Server calls core/tools.Base64Encode() or Base64Decode()
5. Server returns an HTML fragment: just the <textarea> content for output
6. HTMX swaps the output area
```

**Templ component sketch:**

```
tools/base64.templ:
  <form hx-post="/tools/base64" hx-target="#output" hx-trigger="input changed delay:200ms">
    <fieldset> — mode radio buttons, checkboxes
    <label>Input</label>
    <textarea name="input" rows="6">...</textarea>
    <button type="button" onclick="pasteClipboard()">Paste</button>

    <label>Output</label>
    <div id="output">
      <textarea readonly rows="6">{ result }</textarea>
    </div>
    <button type="button" onclick="copyOutput()">Copy</button>
    <button type="button" hx-post="/tools/base64/clear" hx-target="closest form" hx-swap="outerHTML">Clear</button>
  </form>
```

**Behaviour notes:**

- The "◀ All Tools" link in the header navigates back to the index page.
- Input textarea is editable; output textarea is readonly.
- HTMX handles form submission on every input change with a 200ms debounce.
- The "Paste" and "Copy" buttons use a tiny JS helper (~5 lines) for the Clipboard API.
- Error states render inline: if the server returns an error, the output area shows a red-bordered error message instead of the textarea.
- The footer reminds users that processing is local (server-side, but on their self-hosted instance).

---

## Mockup 3: Tool Page — Mobile Layout (< 768px)

```
┌───────────────────────────────┐
│  ☰  ⚒ Forge                   │
├───────────────────────────────┤
│                               │
│  Base64 Encode / Decode       │
│                               │
│  Mode:                        │
│  (●) Encode  ( ) Decode       │
│  ☐ URL-safe                   │
│                               │
│  Input                        │
│  ┌───────────────────────────┐│
│  │ Hello, World!             ││
│  │                           ││
│  └───────────────────────────┘│
│  [Paste]                      │
│                               │
│  Output                       │
│  ┌───────────────────────────┐│
│  │ SGVsbG8sIFdvcmxkIQ==     ││
│  │                           ││
│  └───────────────────────────┘│
│  [Copy]  [Clear]              │
│                               │
├───────────────────────────────┤
│  Forge v0.1.0                 │
└───────────────────────────────┘
```

**Hamburger menu open:**

```
┌───────────────────────────────┐
│  ✕  ⚒ Forge                   │
├───────────────────────────────┤
│                               │
│  ENCODERS                     │
│    Base64                     │
│    JWT Decoder                │
│    URL Encode/Decode          │
│    HTML Entity                │
│                               │
│  FORMATTERS                   │
│    JSON Formatter             │
│                               │
│  GENERATORS                   │
│    Hash Generator             │
│    UUID Generator             │
│    Password Generator         │
│                               │
│  CONVERTERS                   │
│    JSON / YAML                │
│    Timestamp                  │
│    Number Base                │
│                               │
│  TESTERS                      │
│    Regex Tester               │
│                               │
└───────────────────────────────┘
```

**Behaviour notes:**

- On mobile, the header shows a hamburger menu (☰) instead of the "◀ All Tools" breadcrumb.
- Tapping ☰ opens a full-screen navigation overlay listing all tools by category.
- Options (mode, checkboxes) stack vertically instead of inline.
- Button labels shorten: "Paste from clipboard" becomes "Paste".
- Textareas expand to full width with minimal horizontal padding.
- HTMX behavior is identical to desktop — the same endpoints and swap targets.
- The hamburger menu itself does NOT use HTMX — it's a pure CSS toggle (`:target` or checkbox hack) to avoid any JS dependency for navigation.

---

## CSS Architecture

A single stylesheet: `/static/style.css` (~200 lines). Key design tokens:

```css
:root {
    --color-bg:       #0f172a;    /* slate-900 */
    --color-surface:  #1e293b;    /* slate-800 */
    --color-border:   #334155;    /* slate-700 */
    --color-text:     #e2e8f0;    /* slate-200 */
    --color-muted:    #94a3b8;    /* slate-400 */
    --color-accent:   #7c3aed;    /* violet-600 */
    --color-error:    #ef4444;    /* red-500 */
    --color-success:  #22c55e;    /* green-500 */
    --font-mono:      'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
    --font-sans:      system-ui, -apple-system, sans-serif;
    --radius:         6px;
}
```

Light mode is not a v1 priority but the architecture should support it via `@media (prefers-color-scheme: light)` with overridden CSS variables.

---

## Static File Structure

```
ui/web/static/
├── style.css          ← custom stylesheet (~200 lines)
├── htmx.min.js        ← vendored HTMX 2.x (~14KB gzipped)
└── forge.js           ← clipboard helpers, textarea resize (~20 lines)
```

No npm, no bundler, no node_modules. These files are embedded in the Go binary via `go:embed`.

---

## Route Map

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/` | `handleIndex` | Tool list grouped by category |
| GET | `/tools/{id}` | `handleToolPage` | Full tool page with form |
| POST | `/tools/{id}` | `handleToolProcess` | Process input, return HTML fragment |
| GET | `/static/*` | `http.FileServer` | Static CSS/JS assets |
| GET | `/health` | `handleHealth` | Health check endpoint (returns 200 OK) |
