# Architecture Decision Records

**Project:** Forge — Developer Utility Toolkit
**Date:** 2026-03-15

---

## ADR-001: Go as the Implementation Language

### Context

Forge needs a single language that can produce fast native binaries for CLI/TUI, compile to a self-contained web server, and integrate with desktop frameworks. The primary candidates are Go, Rust, Kotlin (JVM or Native), and TypeScript (Node.js + Electron).

The author is a senior Java/Spring Boot engineer learning Go for this project, which introduces a learning curve consideration.

### Decision

**Go.**

Reasons in order of importance:

1. **Single binary deployment.** `go build` produces one statically-linked binary. No JVM, no Node.js runtime, no shared libraries. This is critical for `forge` as a CLI tool that users install via `curl | tar`.

2. **Cross-platform compilation.** `GOOS=darwin GOARCH=arm64 go build` cross-compiles natively. No Docker, no cross-compiler toolchain (except for CGO, which Forge minimizes). GoReleaser automates this.

3. **Bubbletea.** The best TUI framework for Forge's architecture exists only in Go. The Charm ecosystem (bubbletea + lipgloss + bubbles) is unmatched in any other language for building polished terminal UIs.

4. **Fast startup.** Go binaries start in under 100ms. This matters for CLI tools that are invoked hundreds of times per day.

5. **Mature standard library.** Encoding (base64, JSON, hex), hashing (crypto/*), URL parsing (net/url), and HTTP serving (net/http) are all in the Go stdlib. Most of Forge's core logic requires zero external dependencies.

6. **Learning curve is manageable.** Go is intentionally simple. A Java engineer can be productive in Go within 1–2 weeks. The language has fewer concepts to learn than Rust (ownership, lifetimes, traits) or Kotlin (coroutines, multiplatform).

### Consequences

- The author must invest time learning Go idioms (error handling patterns, interface satisfaction, goroutine patterns).
- Rust would offer better memory safety guarantees, but the learning curve is significantly steeper and the TUI ecosystem is less mature.
- Kotlin/JVM would leverage the author's Java expertise but produces JARs that require a JVM, defeating the single-binary goal.
- TypeScript/Electron would offer the fastest UI development but creates bloated desktop apps and cannot produce native CLI tools.

---

## ADR-002: Rejection of Blazor Hybrid (DevToys Architecture)

### Context

DevToys uses Blazor Hybrid — a .NET framework that renders web content inside a native WebView container. This gives DevToys a single codebase that runs on Windows, macOS, and Linux with a web-like UI.

Forge could adopt a similar approach: build a web UI and embed it in a WebView for the desktop surface.

### Decision

**Reject Blazor Hybrid as an architecture. Forge uses distinct UI layers sharing a pure core.**

Reasons:

1. **WebView quality varies by platform.** DevToys' macOS experience (WKWebView) and Linux experience (WebKitGTK) are noticeably worse than Windows. This is a fundamental limitation of the WebView approach — you inherit every platform's WebView bugs.

2. **No terminal UI path.** Blazor Hybrid produces GUI applications. There is no path from Blazor to a terminal UI. Forge's primary surface is the TUI, which requires a fundamentally different rendering approach.

3. **Runtime dependency.** Blazor Hybrid requires the .NET runtime. Forge's Go binary is self-contained.

4. **Startup time.** DevToys takes 2–3 seconds to launch on macOS. A Go binary starts in <100ms.

**Forge's alternative architecture:**

```
core/tools/ (pure Go functions, zero I/O)
    ↑
    ├── ui/tui/   (bubbletea — terminal)
    ├── ui/web/   (chi + templ + HTMX — browser)
    └── ui/desktop/ (Wails — native webview wrapping web UI)
```

Each UI surface imports from `core/` and adapts the tool logic to its rendering model. The desktop surface reuses the web templates inside a Wails webview, achieving the "write once" benefit of Blazor Hybrid without its downsides.

### Consequences

- Three UI layers must be maintained instead of one. This is more code but cleaner architecture.
- The desktop surface will lag behind TUI and web in feature parity (acceptable — it is lowest priority).
- Each UI surface can be optimized independently. The TUI can use terminal-specific features (alternate screen, mouse events) that a WebView cannot.

---

## ADR-003: Core Layer Must Never Import from UI Layers

### Context

In a multi-surface application, the temptation is to add UI-specific logic to the core layer for convenience: "let's add a `RenderHTML()` method to the Tool interface" or "let's have the core format errors as lipgloss strings."

### Decision

**`core/` must never import any package from `ui/` or `cmd/`. This is a hard, non-negotiable rule.**

Enforcement:

1. Code review: any import of `ui/` or `cmd/` in `core/` is rejected.
2. CI: a linting step checks for disallowed imports (custom Go analysis tool or grep-based check).
3. Architecture: `core/tools/` functions return plain Go types (`string`, `struct`), never formatted terminal output or HTML.

### Consequences

- Tool functions are pure: `func Base64Encode(input string, urlSafe bool) Base64Result`. The result is a struct with `Output string` and `Error string`. No ANSI codes, no HTML, no lipgloss styles.
- Each UI layer is responsible for taking a `Result` and rendering it appropriately: lipgloss styling for TUI, templ template for web, cobra output formatting for CLI.
- Testing is simplified: core tool tests only verify logic, not presentation.
- Adding a new UI surface (e.g., a VS Code extension, a Slack bot) requires zero changes to `core/`.
- The cost is duplication: each UI layer has its own rendering logic for the same tool. This is intentional — rendering is inherently surface-specific.

---

## ADR-004: Wails over Electron for Desktop

### Context

The desktop surface needs a framework to create a native-feeling application. The main candidates are Electron, Tauri, Wails, and Fyne.

### Decision

**Wails v2.**

| Criterion | Electron | Tauri | Wails | Fyne |
|-----------|----------|-------|-------|------|
| Language | JS/TS | Rust | Go | Go |
| Binary size | ~150MB+ | ~5MB | ~10MB | ~10MB |
| Runtime | Chromium + Node.js | System WebView | System WebView | Pure Go |
| Memory usage | ~100MB+ | ~30MB | ~30MB | ~50MB |
| Native feel | No (Chromium) | Yes (WebView) | Yes (WebView) | No (custom rendering) |
| Go integration | Poor | None | Excellent | Excellent |

Reasoning:

1. **Go integration.** Wails is built for Go. Backend functions are directly callable from the frontend via type-safe bindings. No HTTP bridge, no IPC overhead.

2. **Binary size.** ~10MB vs Electron's ~150MB+. For a developer utility tool, this matters.

3. **Frontend reuse.** Wails uses a WebView, so Forge can reuse the same HTML/CSS/HTMX templates from the web surface. This dramatically reduces desktop-specific code.

4. **Fyne was rejected** because its custom rendering engine produces UIs that look "off" — neither native nor web-like. For a tool that developers use daily, aesthetics matter.

5. **Tauri was rejected** because it is Rust-based. Adding Rust to a Go project creates a dual-language codebase with separate build systems, toolchains, and mental models.

### Consequences

- Wails v2 is stable but v3 is in alpha with no release date. If v3 introduces breaking changes, migration will be needed.
- Desktop builds require CGO (for WebView bindings), complicating cross-compilation. Desktop binaries will be built per-platform, not via goreleaser's default cross-compilation.
- The desktop surface is deferred to post-v1, reducing immediate risk.
- On Linux, WebKitGTK must be installed, which is a distribution-dependent requirement.

---

## ADR-005: Templ + Chi over Heavier Web Frameworks

### Context

Go has several web framework options: Gin, Echo, Fiber, chi, and the stdlib `net/http`. For templating: stdlib `html/template`, templ, gomponents.

### Decision

**chi for routing, templ for templating, HTMX for interactivity.**

### Why chi over Gin/Echo/Fiber:

- **Zero dependencies.** Chi has no external dependencies. Gin pulls in encoding libraries. Fiber uses fasthttp (incompatible with `net/http` middleware).
- **stdlib compatible.** Chi's middleware interface is `func(http.Handler) http.Handler` — any stdlib middleware works.
- **Minimal magic.** No struct tag binding, no automatic JSON serialization. Explicit routing, explicit handling.
- **Close to stdlib.** If Go's `http.ServeMux` improves further, migrating from chi is trivial. Migrating from Gin or Echo is not.

### Why templ over html/template:

- **Type safety.** Template errors are caught at compile time. With `html/template`, a typo in a template variable name is a runtime panic.
- **IDE support.** The templ LSP provides autocompletion, go-to-definition, and inline error highlighting.
- **Component model.** Templ components are Go functions with typed parameters. Reuse is natural.
- **HTMX integration.** Templ and HTMX pair naturally — `hx-post`, `hx-target`, `hx-swap` attributes work directly in `.templ` files.

### Why HTMX over React/Vue/Svelte:

- **No build step.** HTMX is a single 14KB JS file. No npm, no webpack, no Vite.
- **Server-rendered HTML.** Forge's web surface renders HTML on the server (via templ). HTMX's model — make HTTP request, receive HTML fragment, swap into DOM — matches this perfectly.
- **Forge is stateless.** Each tool request is independent: input → process → output. There is no client-side state to manage. SPA frameworks solve a problem Forge does not have.

### Consequences

- Templ is pre-v1 (v0.3.x). API changes before 1.0 are possible but unlikely to be disruptive.
- Templ requires a code generation step (`templ generate`). This must be integrated into the Makefile and CI.
- HTMX has a learning curve for developers unfamiliar with the hypermedia approach. The team should read the HTMX documentation thoroughly.
- No client-side routing. Every tool is a full page load. This is intentional — URLs are bookmarkable and shareable.

---

## ADR-006: Tool Interface Design

### Context

Every tool in Forge must work across four surfaces (TUI, Web, CLI, Desktop) and be discoverable via smart detection and search. The interface must be rich enough to support UI rendering but simple enough to keep core logic pure.

### Decision

The `Tool` interface is minimal and metadata-focused:

```go
type Tool interface {
    Name() string                       // Display name: "Base64 Encode / Decode"
    ID() string                         // URL/CLI slug: "base64"
    Description() string                // One-line: "Encode and decode Base64 strings"
    Category() string                   // Grouping: "Encoders"
    Keywords() []string                 // Search: ["base64", "encode", "decode", "b64"]
    DetectFromClipboard(s string) bool  // Smart detection: "does this look like my input?"
}
```

Tool *logic* is not part of the interface. Each tool exposes standalone functions:

```go
func Base64Encode(input string, urlSafe bool) Base64Result
func Base64Decode(input string, urlSafe bool) Base64Result
```

This is deliberate. The interface describes the tool for discovery and routing. The functions implement the tool's behavior.

### Reasoning

1. **UI layers need metadata, not logic.** The sidebar needs `Name()` and `Category()`. The search needs `Keywords()`. The CLI router needs `ID()`. None of them need `Process(input) output` — that coupling would force a single function signature on tools with very different parameters.

2. **Tool parameters vary.** Base64 has `urlSafe bool`. JSON has `indent int` and `sortKeys bool`. Hash has `algorithm string`. A generic `Process(input string) Result` would require stuffing all options into a map or interface{}, losing type safety. Standalone functions with explicit parameters are cleaner.

3. **Pure functions are testable.** `Base64Encode("hello", false)` is trivially testable. A method on an interface behind a registry lookup is not.

4. **DetectFromClipboard is tool-specific.** JWT detection checks for three dot-separated segments. UUID detection checks for a hex pattern with hyphens. JSON detection calls `json.Valid()`. These are simple predicates, not general-purpose parsers.

### Consequences

- Each UI layer must know about individual tool functions. The TUI's `base64.go` calls `tools.Base64Encode()` directly, not via an abstract `tool.Process()`.
- Adding a new tool requires: (1) implement the `Tool` interface, (2) write the logic functions, (3) register in the registry, (4) add UI views in each surface. This is ~4 files of work per tool.
- The registry handles discovery and routing; individual tool functions handle logic. Clean separation.

---

## Summary Table

| ADR | Decision | Key Tradeoff |
|-----|----------|--------------|
| 001 | Go | Learning curve vs. ecosystem fit |
| 002 | No Blazor Hybrid | More code vs. clean architecture |
| 003 | Core never imports UI | Rendering duplication vs. testability |
| 004 | Wails over Electron | Smaller ecosystem vs. smaller binary |
| 005 | Templ + chi + HTMX | Pre-v1 template lib vs. type safety |
| 006 | Minimal Tool interface | Per-tool UI wiring vs. type-safe parameters |
