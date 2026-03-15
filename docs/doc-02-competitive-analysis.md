# Competitive Analysis

**Project:** Forge — Developer Utility Toolkit
**Date:** 2026-03-15

---

## Projects Analysed

### 1. DevToys (DevToys-app/DevToys)

**What it is:** Cross-platform developer Swiss Army knife. C#/.NET, Blazor Hybrid UI. Windows/macOS/Linux. 28k+ GitHub stars. The direct inspiration for Forge.

**What it does well:**

- Comprehensive tool set (30+ built-in tools, 44 extensions). Covers most daily developer needs.
- Smart Detection is the killer feature. Clipboard content automatically routes to the right tool. This is the single most important UX innovation to carry into Forge.
- Clean, native-feeling UI with consistent design language.
- Extension system allows community tools without bloating the core.
- CLI companion (`DevToys.CLI`) for scripting use cases.

**What it does poorly:**

- Built on Blazor Hybrid, which bundles a WebView runtime. The macOS and Linux experience is noticeably worse than Windows — WebKitGTK on Linux has rendering quirks, and the app feels sluggish.
- No terminal-native experience. Developers who live in the terminal must context-switch to a GUI.
- The extension API is C#-specific, limiting the contributor base to .NET developers.
- Startup time is noticeable (~2–3 seconds on macOS) due to the .NET runtime.
- No self-hostable web version. You must install the desktop app.

**What Forge should steal:** Smart Detection, tool categorization model, the core tool selection (Base64, JWT, JSON, Hash, URL, UUID).

**What Forge should avoid:** Blazor Hybrid architecture, heavy runtime dependencies, GUI-only approach.

---

### 2. DevUtils (devutils.com)

**What it is:** macOS-native developer toolkit. Paid ($9), closed source. Native Swift/AppKit.

**What it does well:**

- Truly native macOS experience — fast, feels like a system utility.
- Menubar integration: quick access without launching a full app.
- Smart Detection works from system clipboard without the app being in focus.
- Offline-first, privacy-focused.
- Beautiful, polished UI.

**What it does poorly:**

- macOS only. No Windows, Linux, or web.
- Closed source — cannot extend or contribute.
- Paid software in a space with strong free alternatives.
- No CLI mode — cannot integrate into scripts.
- No terminal UI — same context-switching problem as DevToys.

**What Forge should steal:** The menubar/system-tray integration concept (future, not v1). The "always running, always watching clipboard" background mode.

**What Forge should avoid:** Platform lock-in, closed source, paid model.

---

### 3. CyberChef (gchq/CyberChef)

**What it is:** Web-based "Cyber Swiss Army Knife." JavaScript, runs entirely in the browser. Open source, maintained by GCHQ (UK intelligence).

**What it does well:**

- Extremely powerful chaining model: you build "recipes" by stacking operations. Input flows through multiple transforms sequentially.
- 300+ operations covering encoding, encryption, compression, data analysis.
- Runs entirely client-side — no server needed. Privacy by architecture.
- Shareable recipes via URL encoding. You can send someone a link with a pre-built transformation pipeline.
- Excellent for CTF challenges and security analysis.

**What it does poorly:**

- The UI is overwhelming for simple tasks. Opening CyberChef to Base64-decode a string is like using a chainsaw to cut bread.
- No native desktop or terminal experience. Browser-only.
- Performance degrades with large inputs (JavaScript single-thread limitation).
- No smart detection — you must manually select operations.
- The learning curve is steep for the chaining/recipe model.

**What Forge should steal:** Nothing directly for v1 — CyberChef serves a different audience (security analysts). However, the recipe/chaining concept could inspire a "pipeline" feature in v2+.

**What Forge should avoid:** Feature overload, complex UI for simple tasks.

---

### 4. Charm's Soft Serve (charmbracelet/soft-serve)

**What it is:** A self-hostable Git server with a built-in TUI. Go, bubbletea. Not a developer toolkit, but the best publicly available reference for a complex bubbletea application.

**What it does well:**

- Multi-panel bubbletea layout: sidebar navigation + main content area. This is the exact layout Forge needs.
- Proper focus management between panels.
- Clean code architecture: models are composable, styles are centralized.
- Good example of handling async operations (git operations) within the bubbletea event loop.
- Demonstrates how to build a polished, production-grade bubbletea TUI.

**What it does poorly:**

- It is a Git server, not a toolkit — so the tool-per-panel model doesn't apply directly.
- The codebase is large and tightly coupled to Git concepts, making it harder to extract reusable patterns.

**What Forge should steal:** The sidebar + content panel layout pattern. The focus management approach. The style centralization pattern.

---

### 5. Glow (charmbracelet/glow)

**What it is:** Terminal-based Markdown reader. Go, bubbletea. By Charm (same team as bubbletea).

**What it does well:**

- Beautiful terminal rendering with lipgloss.
- Clean, focused TUI — does one thing well.
- Good example of viewport scrolling for large content.
- Stash feature (save/organize documents) shows how to add persistence to a bubbletea app.
- Excellent responsive layout handling for different terminal sizes.

**What it does poorly:**

- Single-purpose — no multi-tool architecture to reference.
- No sidebar navigation (list view is separate from document view, not side-by-side).

**What Forge should steal:** Viewport scrolling patterns for tool output. Terminal resize handling. The quality bar for terminal aesthetics.

---

### 6. OpenDev / Webacus

**What they are:** Web-based DevToys alternatives. JavaScript, open source.

**What they do well:**

- Zero-install: just open a URL.
- Clean, modern web UI.
- Cover the basic tool set (Base64, JSON, Hash, UUID, URL).

**What they do poorly:**

- Browser-only — no terminal, no desktop, no CLI.
- Feature set is much smaller than DevToys.
- No smart detection.
- No offline capability (unless manually saved as PWA).

**What Forge should steal:** The web UI simplicity. Forge's web surface should be this clean.

**What Forge should avoid:** Being web-only.

---

## Key Takeaways for Forge

1. **Smart Detection is non-negotiable.** DevToys and DevUtils prove this is the single most valuable UX feature. Forge must ship with it in v1 across all surfaces (TUI clipboard polling, web paste detection via JS).

2. **Terminal-first is the differentiator.** No existing toolkit offers a first-class terminal UI. This is Forge's unique selling point. Every competitor requires launching a GUI or browser.

3. **Start focused, expand later.** CyberChef's 300+ operations are impressive but overwhelming. DevToys' 30 tools are comprehensive but some are niche. Forge should ship 6 must-have tools in Tier 1, expand to 13 in Tier 2, and resist the urge to add tools just because they exist.

4. **Speed matters.** DevToys' 2–3 second startup is noticeable. A Go binary starts in under 100ms. This is a meaningful advantage — lean into it.

5. **Multi-surface is the moat.** No competitor offers TUI + Web + Desktop + CLI from a single codebase. If Forge delivers even TUI + CLI + Web in v1, it occupies a unique position.

6. **Self-hosting is a real use case.** DevToys and DevUtils require installation. CyberChef is web-only. Forge's self-hostable web server (`forge web`) is genuinely useful for teams who want a shared utility server on their internal network.

---

## Features to Prioritise Based on Gaps

| Gap in Existing Tools | Forge's Answer | Priority |
|----------------------|----------------|----------|
| No terminal-native toolkit exists | TUI as primary surface | P0 |
| No tool has TUI + CLI + Web | Multi-surface architecture | P0 |
| Smart Detection only in GUI apps | Clipboard detection in TUI | P0 |
| No self-hostable web toolkit | `forge web` command | P1 |
| DevToys startup is slow | Go binary, <100ms cold start | P1 |
| CyberChef is overwhelming for simple tasks | Clean, focused tool set | P1 |
| No tool offers pipe-friendly CLI | Full stdin/stdout CLI mode | P0 |

---

## Things to Avoid

1. **Feature creep.** Do not try to match DevToys' tool count or CyberChef's operation count in v1. Ship 6 tools that work perfectly.

2. **Heavy runtimes.** Blazor Hybrid (.NET), Electron (Node.js + Chromium). Forge's Go binary should be <15MB and start instantly.

3. **GUI-first thinking.** The TUI is the primary surface. Design the tool interface for terminal constraints first, then adapt up to web and desktop.

4. **Complex chaining/pipelines.** CyberChef's recipe model is powerful but antithetical to Forge's "one tool, one job" philosophy. Unix pipes (`forge base64 decode | forge json format`) provide chaining for free.

5. **Mandatory accounts or cloud features.** Every competitor that works offline wins trust. Forge must never require accounts, internet, or telemetry.

6. **Platform-specific UI frameworks.** DevUtils' AppKit lock-in and DevToys' Blazor quirks demonstrate the cost. Forge's architecture (Go core + platform-agnostic UI layers) avoids this trap.
