# Risk Register

**Project:** Forge — Developer Utility Toolkit
**Date:** 2026-03-15

---

## Risk Assessment Matrix

**Likelihood:** Low (< 20%) | Medium (20–60%) | High (> 60%)
**Impact:** Low (workaround exists, minor delay) | Medium (significant rework, weeks of delay) | High (project viability threatened)

---

| # | Risk | Likelihood | Impact | Mitigation |
|---|------|-----------|--------|------------|
| R1 | **Bubbletea v2 complexity exceeds estimates.** Multi-panel focus management, responsive layouts, and async clipboard polling are all non-trivial in bubbletea. The Elm Architecture is elegant but forces all state changes through messages, which can become unwieldy in complex UIs. | **High** | **High** | Build the TUI incrementally. Start with a single-panel app (just one tool, no sidebar). Add the sidebar second. Add focus management third. Add clipboard detection last. Each increment should be a working, testable milestone. Reference Soft Serve's codebase for patterns. Use `bubblelayout` or `bubbletea-nav` libraries if custom layout code becomes unmanageable. Budget 4–6 weeks for TUI development, not the 1–2 weeks a simple CLI would take. |
| R2 | **Go learning curve slows development.** The author is a senior Java/Spring Boot engineer. Go's error handling patterns, interface satisfaction model, goroutines, and package layout conventions differ significantly from Java idioms. | **Medium** | **Medium** | Invest the first week purely in Go fundamentals: the Go tour, Effective Go, and building a small CLI tool (not Forge) to internalize idioms. Follow the Go proverbs. Use `golangci-lint` from day one to catch non-idiomatic code early. Leverage `testify` for familiar assertion-style testing. The Java → Go transition is easier than Java → Rust, but allow 2–3 weeks before expecting full productivity. |
| R3 | **Scope creep beyond the 6 Tier-1 tools.** The temptation to add "just one more tool" before shipping v1 is strong, especially when the tool inventory document lists 20+ candidates. | **High** | **Medium** | Define a hard v1.0 scope: the 6 Tier-1 tools (base64, jwt, json, hash, url, uuid) + TUI + CLI. Web is v1.1. Desktop is v2. Tier-2 tools are v1.2+. Write this scope in the README and reference it in every planning session. If a tool idea emerges during development, add it to the inventory document with a future milestone, not to the current sprint. |
| R4 | **Wails v2 becomes unmaintained as v3 development absorbs all attention.** Wails v3 has been in alpha for an extended period. If the team focuses entirely on v3, v2 bug fixes may slow or stop. | **Medium** | **Low** | The desktop surface is deferred to post-v1. By the time Forge needs it, v3 may be stable. If v2 is abandoned and v3 is not ready, Forge can delay desktop indefinitely — the TUI and web surfaces cover the primary use cases. Worst case: drop Wails entirely and ship the web surface as a PWA for a desktop-like experience. |
| R5 | **Templ breaks API before reaching v1.0.** Templ is at v0.3.x. Pre-v1 libraries may change their API. | **Low** | **Medium** | Templ has been stable in practice with 5,481 known importers. The risk of a breaking API change is low. If it happens, the migration scope is limited to `ui/web/templates/` — the core and TUI are unaffected. Pin the templ version in `go.mod` and upgrade deliberately, not automatically. |
| R6 | **Cross-platform clipboard issues.** Linux clipboard access requires `xclip` or `xsel`. Headless environments (SSH, CI, Docker) have no clipboard. Wayland vs X11 creates fragmentation on Linux. | **High** | **Low** | Wrap clipboard access behind `internal/clipboard/` with a clean interface. On failure (no `xclip`, no display), log a warning and disable smart detection gracefully — the app must still work without clipboard access. Document the `xclip`/`xsel` dependency in the README. In the TUI, show a non-blocking notice: "Clipboard detection unavailable — install xclip for smart detection." |
| R7 | **Terminal rendering inconsistencies across platforms.** macOS Terminal.app, iTerm2, Windows Terminal, Alacritty, Ghostty, and GNOME Terminal all render ANSI codes differently. Box-drawing characters, Unicode, and colors behave inconsistently. | **Medium** | **Medium** | Test on at least 4 terminals during development: macOS Terminal.app, iTerm2 (or Ghostty), Windows Terminal, and one Linux terminal. Use lipgloss's `AdaptiveColor` to handle light/dark detection. Avoid Unicode characters beyond basic box-drawing (stick to `┌─┐│└┘`). Bubbletea v2's Mode 2026 support reduces tearing in modern terminals. Accept that some terminals will render imperfectly and document known issues. |
| R8 | **Bubbletea v2 migration instability.** Bubbletea v2 just shipped. Early adopters may encounter undiscovered bugs or documentation gaps. The import path change (`charm.land/bubbletea/v2`) and View struct model are significant changes from v1. | **Medium** | **Medium** | Start on v2 directly — do not build on v1 and migrate. Use the official upgrade guide. Follow Charm's GitHub issues and discussions for early bug reports. If a blocking bug is found, pin to the last known-good v2 release and report the issue upstream. The Charm team is responsive (active Discord, GitHub). |
| R9 | **HTMX limitations for complex tool interactions.** Tools like Regex Tester (live highlighting) or JSON Formatter (real-time output) may strain HTMX's request/response model, causing sluggish UX with rapid input changes. | **Low** | **Low** | Use `hx-trigger="input changed delay:200ms"` to debounce requests. For truly latency-sensitive tools, add a `<script>` that processes client-side and falls back to the server endpoint. Forge's tools are computationally trivial — even with a 200ms debounce round-trip, the UX should feel responsive on localhost. |
| R10 | **goreleaser CGO cross-compilation for desktop builds.** Wails requires CGO for WebView bindings. CGO cross-compilation is notoriously painful — it requires platform-specific C toolchains, Docker-based build environments, or separate CI runners per OS. | **High** | **Low** | Separate desktop builds from CLI/TUI builds entirely. The CLI/TUI binary is pure Go (no CGO) and cross-compiles trivially via goreleaser. Desktop binaries are built per-platform using GitHub Actions matrix builds (macOS runner for macOS, Windows runner for Windows, Ubuntu runner for Linux). This is more CI config but avoids the CGO cross-compilation nightmare entirely. Since desktop is post-v1, this is a future problem. |

---

## Risk Heatmap

```
              Low Impact    Medium Impact    High Impact
           ┌─────────────┬────────────────┬──────────────┐
  High     │  R6          │  R3            │  R1          │
Likelihood │              │                │              │
           ├─────────────┼────────────────┼──────────────┤
  Medium   │  R4          │  R2, R7, R8    │              │
Likelihood │              │                │              │
           ├─────────────┼────────────────┼──────────────┤
  Low      │  R9, R10     │  R5            │              │
Likelihood │              │                │              │
           └─────────────┴────────────────┴──────────────┘
```

---

## Top 3 Risks by Severity (Likelihood x Impact)

1. **R1 — Bubbletea complexity.** The TUI is the primary surface and the hardest to build. Mitigate by building incrementally and budgeting adequate time.

2. **R3 — Scope creep.** The tool inventory is deliberately large to inspire future work. v1 must be ruthlessly scoped to 6 tools + TUI + CLI.

3. **R2 — Go learning curve.** Manageable but real. A dedicated onboarding week before writing Forge code pays for itself many times over.

---

## Risk Review Cadence

Review this register at each milestone:

- **M1: Core tools + tests complete** — Re-evaluate R2 (Go learning curve should be resolved)
- **M2: TUI functional** — Re-evaluate R1, R7, R8 (bubbletea risks should be characterized)
- **M3: Web server functional** — Re-evaluate R5, R9 (templ and HTMX risks should be characterized)
- **M4: v1.0 release** — Re-evaluate R3 (scope creep), add new risks from user feedback
- **M5: Desktop surface** — Re-evaluate R4, R10 (Wails and CGO risks become active)
