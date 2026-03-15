# Tech Stack Validation Report

**Project:** Forge — Developer Utility Toolkit
**Date:** 2026-03-15
**Author:** Architecture Review

---

## Summary Verdict

| Library | Version | Verdict | Notes |
|---------|---------|---------|-------|
| bubbletea | v2.0.0 | **KEEP** | Best-in-class Go TUI framework, actively maintained, v2 just shipped |
| lipgloss | v2.x | **KEEP** | Companion to bubbletea, no alternative needed |
| bubbles | v2.x | **KEEP** | Essential component library for bubbletea |
| templ | v0.3.x | **KEEP (with caution)** | Pre-v1 but widely adopted, active development |
| chi | v5.x | **KEEP** | Lightweight, idiomatic, zero-dependency router |
| HTMX | 2.x | **KEEP** | Proven pairing with templ in Go ecosystem |
| Wails | v2 stable / v3 alpha | **KEEP v2, WATCH v3** | v2 is production-ready; v3 alpha not ready |
| cobra | v1.x | **KEEP** | De facto standard, used by kubectl/docker/helm |
| goreleaser | v2.x | **KEEP** | Industry standard for Go release automation |
| clipboard | — | **USE atotto/clipboard** | Simpler, sufficient for text-only needs |

---

## bubbletea

**Current version:** v2.0.0 (released 2025, stable as of early 2026)
**Import path:** `charm.land/bubbletea/v2`
**GitHub stars:** 28k+
**Last activity:** Actively maintained, regular releases throughout 2025–2026

**What's new in v2:**

- New "Cursed Renderer" built from scratch on the ncurses rendering algorithm — significantly reduces flicker.
- Mode 2026 support for synchronized output in modern terminals (Ghostty, etc.).
- Declarative View struct replaces imperative commands for terminal features (alt screen, mouse, etc.).
- Up to 30% faster rendering and improved memory management.
- Breaking changes from v1: new import path, message type changes, View struct model.

**Known issues and limitations:**

- Multi-panel focus management is not built-in — you must implement it yourself or use community libraries like `bubbletea-nav` or `bubblelayout`.
- Complex layouts require manual width/height math. There is no CSS-like flexbox; you calculate sizes from `tea.WindowSizeMsg`.
- Async operations (clipboard polling) must be handled via `tea.Cmd` that return messages — cannot mutate state directly.
- v2 migration from v1 is non-trivial (import path change, View struct rewrite, renderer changes).

**Alternatives considered:**

- **tview** (rivo/tview): Widget-based rather than Elm Architecture. Better for form-heavy apps, worse for custom layouts. Less idiomatic Go. Not recommended — bubbletea's composability is a better fit for Forge's multi-surface architecture.
- **tcell**: Too low-level. Would require building everything from scratch.
- **Standard library (no framework)**: Not practical for a multi-panel interactive TUI.

**Verdict:** bubbletea v2 is the right choice. It is the dominant Go TUI framework, actively maintained by Charm, and the Elm Architecture maps cleanly to Forge's tool model. Start with v2 directly — do not build on v1.

---

## lipgloss

**Current version:** v2.x
**Import path:** `charm.land/lipgloss/v2`
**Last activity:** January 2026 (performance optimizations)

Lipgloss is the styling companion to bubbletea. It provides CSS-like declarative styling for terminal output: colors, borders, padding, margins, alignment. It supports adaptive colors (auto-detect light/dark terminals), 256 colors, and true color (24-bit). Clickable hyperlinks in supported terminals.

**Known issues:**

- No layout engine — lipgloss styles individual blocks, but arranging them is manual (JoinHorizontal, JoinVertical, Place).
- Performance can degrade with very large rendered strings due to ANSI escape code processing.

**Alternatives:** None worth considering. Lipgloss is purpose-built for bubbletea and is the only serious terminal styling library in Go.

**Verdict:** Keep. Non-negotiable companion to bubbletea.

---

## bubbles

**Current version:** v2.x (aligned with bubbletea v2)
**Import path:** `charm.land/bubbles/v2`

Bubbles provides pre-built bubbletea components: text input, text area, viewport (scrollable content), spinner, progress bar, paginator, list, table, file picker, help, key bindings, timer, and stopwatch.

**Relevant to Forge:**

- `textinput` — tool input fields
- `textarea` — multi-line input/output
- `viewport` — scrollable output panels
- `list` — sidebar tool list with filtering
- `help` — keybinding help display
- `key` — keybinding definitions

**Known issues:**

- The `list` component includes built-in filtering but can be finicky to customize heavily.
- Components are designed to be embedded in your own models, not used standalone.

**Verdict:** Keep. Provides essential building blocks that would take weeks to build from scratch.

---

## templ

**Current version:** v0.3.x (latest release: February 2026)
**Import path:** `github.com/a-h/templ`
**GitHub stars:** 9k+
**Importers:** 5,481 known packages

Templ is a type-safe HTML templating language for Go. It compiles `.templ` files into Go code, giving you compile-time type checking and IDE support (LSP). It integrates natively with Go's `http.Handler` interface.

**Why templ over html/template:**

- Type-safe: template errors are caught at compile time, not runtime.
- IDE support: autocompletion, go-to-definition, error highlighting.
- Component-based: `.templ` files define reusable components with typed parameters.
- No runtime reflection overhead.

**Known issues:**

- Pre-v1: API may change before 1.0, though the project has been stable in practice.
- Requires a code generation step (`templ generate`) before build.
- Developer tooling (LSP) can occasionally lag behind the latest release.
- Learning curve: `.templ` syntax is Go-like but distinct from standard Go templates.

**HTMX integration:** Templ and HTMX pair cleanly. Multiple production template repositories exist (go-templ-htmx, go-htmx-template). The templ docs include an official HTMX integration guide. HTMX attributes (`hx-post`, `hx-target`, `hx-swap`) work naturally in `.templ` files.

**Alternatives considered:**

- **html/template (stdlib)**: Works but painful for anything beyond trivial templates. No type safety, no IDE support, runtime panics on errors. Not recommended for a project of Forge's scope.
- **gomponents**: Pure Go HTML generation. Interesting but less ergonomic for larger templates.

**Verdict:** Keep. Pre-v1 status is a minor concern, but 5k+ importers and active development make it production-viable. The type safety and HTMX integration are significant advantages.

---

## chi

**Current version:** v5.x
**Import path:** `github.com/go-chi/chi/v5`
**GitHub stars:** 18k+

Chi is a lightweight, idiomatic HTTP router for Go with zero external dependencies. It is fully compatible with `net/http` and supports middleware chaining, URL parameters, and route grouping.

**Why chi over alternatives:**

- Zero dependencies (unlike Gin, Echo, Fiber).
- Fully compatible with `net/http` — middleware from the standard ecosystem works directly.
- Clean middleware stack (built-in logging, recovery, CORS, etc.).
- No magic — explicit routing, no struct tags or reflection.

**Why not stdlib `http.ServeMux`:** Go 1.22 improved the stdlib mux significantly (method-based routing, path parameters). For Forge's web surface, the stdlib mux would actually be sufficient. However, chi adds route grouping and a cleaner middleware API for negligible cost. The decision is close.

**Recommendation:** Keep chi, but note that migrating to stdlib `http.ServeMux` would be trivial if desired later. Chi's middleware ecosystem is the tiebreaker.

**Security note:** A past open redirect vulnerability was found in chi's `RedirectSlashes` middleware. Ensure you use the latest v5.x release.

---

## HTMX

**Current version:** 2.x
**Distribution:** Single JS file via CDN or vendored

HTMX allows HTML elements to make HTTP requests and swap content without full page reloads. It is the natural fit for Forge's web surface: each tool submits input via `hx-post`, receives HTML fragments, and swaps the output area.

**Why HTMX for Forge:**

- No build step, no npm, no bundler. Just a `<script>` tag.
- Server-rendered HTML (via templ) is the response format — no JSON API layer needed.
- Minimal JS footprint (~14KB gzipped).
- Perfect for stateless request/response tools.

**Known issues:**

- HTMX debugging can be opaque (network tab + `hx-trigger` timing issues).
- Complex multi-step interactions (e.g., chained transformations) may need creative `hx-trigger` usage.
- Not suitable if you later want a rich SPA — but Forge's web surface explicitly does not need that.

**Verdict:** Keep. Ideal match for Forge's stateless, tool-per-page architecture.

---

## Wails

**Current version:** v2.9.x (stable), v3 (alpha)
**GitHub stars:** 25k+

Wails lets you build desktop applications using Go for the backend and web technologies for the frontend. It embeds a native webview (not Chromium/Electron), resulting in small binaries (~10MB vs Electron's ~150MB+).

**v2 status:** Production-ready. Stable API. Supports Windows, macOS, Linux.

**v3 status:** Alpha. API is "reasonably stable" but documentation and tooling are incomplete. No release timeline announced. Some apps run v3 in production, but the team is taking a measured approach to stabilization.

**Recommendation:** Build on Wails v2. The desktop surface is explicitly lower priority than TUI and Web. Wails v2 is stable and sufficient. When v3 reaches stable, migration can be planned as a separate effort. Do not build on alpha software for a surface that isn't the primary focus.

**Known issues:**

- macOS webview (WKWebView) has quirks around file dialogs and system integration.
- Linux support depends on WebKitGTK, which varies across distributions.
- Hot reload in development is less smooth than Electron's dev experience.
- Build times are slower than pure Go due to CGO requirements.

**Alternatives considered:**

- **Electron**: Rejected. Enormous binary size, requires Node.js runtime, defeats the purpose of a Go project.
- **Tauri**: Rust-based, not Go. Would require a separate frontend build pipeline.
- **Fyne**: Pure Go GUI toolkit, but creates non-native-looking UIs. Not suitable for a tool that should feel modern.

**Verdict:** Keep Wails v2. Watch v3 for future migration.

---

## cobra

**Current version:** v1.8.x
**Import path:** `github.com/spf13/cobra`
**GitHub stars:** 38k+
**Used by:** 173,000+ projects (kubectl, Docker, Hugo, Helm, etc.)

Cobra is the de facto standard for building Go CLI applications. It provides command/subcommand trees, flag parsing, auto-generated help, shell completions (bash, zsh, fish, PowerShell), and man page generation.

**Why cobra for Forge:**

- Forge's CLI surface maps perfectly to cobra's command tree: `forge base64 encode`, `forge jwt decode`, etc.
- Auto-generated `--help` and shell completions for free.
- Pairs with Viper for configuration management if needed later.
- Every Go developer knows cobra — low learning curve for contributors.

**Alternatives considered:**

- **urfave/cli**: Simpler, better for flat CLIs. Forge has nested subcommands, making cobra a better fit.
- **kong**: Struct-tag-based. Interesting but less mainstream. Wouldn't gain enough to justify the unfamiliarity.
- **stdlib `flag` + custom routing**: Too much boilerplate for 6+ tools with subcommands.

**Verdict:** Keep. No reason to deviate from the industry standard.

---

## goreleaser

**Current version:** v2.x (open source), v2.13.x (Pro)
**Configuration:** `.goreleaser.yaml` with `version: 2` schema

GoReleaser automates Go binary builds, checksums, archives, changelogs, and publishing to GitHub Releases, Homebrew, Scoop, Docker registries, and more.

**Forge needs:**

- Cross-platform builds: linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64.
- Archive formats: `tar.gz` (Linux/macOS), `zip` (Windows).
- GitHub Releases publishing.
- Homebrew tap (future).
- Checksum file generation.

**GitHub Actions integration:** Straightforward. GoReleaser provides an official GitHub Action (`goreleaser/goreleaser-action`). Triggered on tag push, builds and publishes in ~2 minutes.

**Known issues:**

- Free version lacks some features (monorepo support, includes). Not relevant for Forge.
- CGO cross-compilation (needed for Wails desktop builds) requires additional setup (Docker, cross-compilers). Recommendation: build desktop binaries separately from CLI/TUI releases.

**Verdict:** Keep. Industry standard, well-documented, works out of the box for Forge's needs.

---

## Clipboard Library Recommendation

**Two main options:**

| Library | Text | Image | Platforms | CGO Required | Dependencies |
|---------|------|-------|-----------|-------------|--------------|
| `atotto/clipboard` | Yes | No | macOS, Windows, Linux (xclip/xsel) | No | Minimal |
| `golang.design/x/clipboard` | Yes | Yes (PNG) | macOS, Windows, Linux, Android, iOS | Yes | Heavier |

**Recommendation: `atotto/clipboard`**

Forge only needs text clipboard access (detecting Base64 strings, JWT tokens, URLs, JSON, UUIDs). Image clipboard support is unnecessary. `atotto/clipboard` is simpler, does not require CGO, and has fewer platform-specific complications.

**Caveat:** On Linux, `atotto/clipboard` requires `xclip` or `xsel` to be installed. Document this as a dependency. On headless Linux (servers, CI), clipboard operations should gracefully degrade (log a warning, disable smart detection).

**Abstraction:** Wrap clipboard access behind a `internal/clipboard` interface so the implementation can be swapped later if image support becomes needed.

---

## Final Recommended Stack

| Layer | Library | Version | Import Path |
|-------|---------|---------|-------------|
| TUI framework | bubbletea | v2.0.0 | `charm.land/bubbletea/v2` |
| TUI styling | lipgloss | v2.x | `charm.land/lipgloss/v2` |
| TUI components | bubbles | v2.x | `charm.land/bubbles/v2` |
| Web templating | templ | v0.3.x | `github.com/a-h/templ` |
| Web router | chi | v5.x | `github.com/go-chi/chi/v5` |
| Web interactivity | HTMX | 2.x | Vendored JS file |
| Desktop | Wails | v2.9.x | `github.com/wailsapp/wails/v2` |
| CLI | cobra | v1.8.x | `github.com/spf13/cobra` |
| Build/release | goreleaser | v2.x | `.goreleaser.yaml` config |
| Clipboard | atotto/clipboard | latest | `github.com/atotto/clipboard` |
| Testing | testify | v1.9.x | `github.com/stretchr/testify` |
| UUID generation | google/uuid | v1.6.x | `github.com/google/uuid` |

**Go version:** 1.22+ minimum. Use range-over-int and enhanced stdlib routing where appropriate. Target 1.23 if available at development start.

This stack is cohesive, well-maintained, and avoids unnecessary dependencies. Every library has a clear purpose, active maintenance, and a path forward.
