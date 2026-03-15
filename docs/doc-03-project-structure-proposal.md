# Project Structure Proposal

**Project:** Forge — Developer Utility Toolkit
**Date:** 2026-03-15

---

## Annotated Directory Tree

```
forge/
├── go.mod                          # Go module definition: module github.com/StanMarek/forge
├── go.sum                          # Dependency checksums (auto-generated)
├── main.go                         # Entrypoint: calls cmd.Execute(). 5 lines, no logic.
├── Makefile                        # Build, test, lint, run targets
├── README.md                       # Project documentation, install instructions, usage
├── LICENSE                         # MIT license text
├── .gitignore                      # Go, build artifacts, IDE, templ generated files
├── .goreleaser.yaml                # Cross-platform release configuration
│
├── cmd/                            # CLI command definitions (cobra)
│   ├── root.go                     # Root cobra command. Registers all subcommands. Default action: launch TUI.
│   ├── base64.go                   # `forge base64 encode|decode` subcommands
│   ├── jwt.go                      # `forge jwt decode|validate` subcommands
│   ├── json.go                     # `forge json format|minify|validate` subcommands
│   ├── hash.go                     # `forge hash <algorithm>` subcommand
│   ├── url.go                      # `forge url encode|decode|parse` subcommands
│   ├── uuid.go                     # `forge uuid generate|validate|parse` subcommands
│   ├── yaml.go                     # `forge yaml to-json|to-yaml` subcommands
│   ├── timestamp.go                # `forge timestamp from-unix|to-unix|now` subcommands
│   ├── number_base.go              # `forge number-base` subcommand
│   ├── regex.go                    # `forge regex` subcommand
│   ├── html_entity.go              # `forge html-entity encode|decode` subcommands
│   ├── password.go                 # `forge password` subcommand
│   ├── lorem.go                    # `forge lorem` subcommand
│   ├── tui.go                      # `forge tui` — launches bubbletea program
│   ├── web.go                      # `forge web` — launches chi HTTP server
│   └── version.go                  # `forge version` — prints build info
│
├── core/                           # Pure business logic. NEVER imports from ui/ or cmd/.
│   ├── tools/                      # Tool interface + all tool implementations
│   │   ├── tool.go                 # Tool interface definition, Result types
│   │   ├── base64.go               # Base64Encode(), Base64Decode() pure functions
│   │   ├── base64_test.go          # Unit tests for base64
│   │   ├── jwt.go                  # JWTDecode(), JWTValidate() pure functions
│   │   ├── jwt_test.go             # Unit tests for jwt
│   │   ├── json.go                 # JSONFormat(), JSONMinify(), JSONValidate()
│   │   ├── json_test.go            # Unit tests for json
│   │   ├── hash.go                 # Hash() with algorithm parameter
│   │   ├── hash_test.go            # Unit tests for hash
│   │   ├── url.go                  # URLEncode(), URLDecode(), URLParse()
│   │   ├── url_test.go             # Unit tests for url
│   │   ├── uuid.go                 # UUIDGenerate(), UUIDValidate(), UUIDParse()
│   │   ├── uuid_test.go            # Unit tests for uuid
│   │   ├── yaml.go                 # YAMLToJSON(), JSONToYAML()
│   │   ├── yaml_test.go            # Unit tests for yaml
│   │   ├── timestamp.go            # TimestampFromUnix(), TimestampToUnix(), TimestampNow()
│   │   ├── timestamp_test.go       # Unit tests for timestamp
│   │   ├── number_base.go          # NumberBaseConvert()
│   │   ├── number_base_test.go     # Unit tests for number base
│   │   ├── regex.go                # RegexTest()
│   │   ├── regex_test.go           # Unit tests for regex
│   │   ├── html_entity.go          # HTMLEntityEncode(), HTMLEntityDecode()
│   │   ├── html_entity_test.go     # Unit tests for html entity
│   │   ├── password.go             # PasswordGenerate()
│   │   ├── password_test.go        # Unit tests for password
│   │   ├── lorem.go                # LoremGenerate()
│   │   └── lorem_test.go           # Unit tests for lorem
│   │
│   ├── registry/                   # Tool registry — registration, lookup, search
│   │   ├── registry.go             # Registry struct: Register(), All(), ByID(), ByCategory(), Search(), Detect()
│   │   ├── registry_test.go        # Unit tests for registry
│   │   └── defaults.go             # DefaultRegistry() — creates registry with all tools pre-registered
│   │
│   └── detection/                  # Smart clipboard detection engine
│       ├── detection.go            # Detector struct, polls clipboard, emits DetectionResult
│       └── detection_test.go       # Unit tests for detection logic (not clipboard I/O)
│
├── ui/                             # All UI surfaces. Imports from core/, never imported by core/.
│   ├── tui/                        # Terminal UI (bubbletea)
│   │   ├── app.go                  # Root AppModel: sidebar + tool panel + detection banner
│   │   ├── sidebar.go              # Sidebar model: tool list, categories, search, navigation
│   │   ├── banner.go               # Detection banner model: show/dismiss/accept
│   │   ├── styles/
│   │   │   └── styles.go           # Lipgloss style definitions, color palette, adaptive colors
│   │   ├── keys/
│   │   │   └── keys.go             # Keybinding definitions (bubbles/key)
│   │   └── views/                  # Per-tool TUI models (each implements tea.Model)
│   │       ├── base64.go           # Base64 TUI: encode/decode radio, input/output textareas
│   │       ├── jwt.go              # JWT TUI: token input, decoded header/payload display
│   │       ├── json.go             # JSON TUI: format/minify/validate modes, indent option
│   │       ├── hash.go             # Hash TUI: algorithm selector, input, hash output
│   │       ├── url.go              # URL TUI: encode/decode/parse modes
│   │       ├── uuid.go             # UUID TUI: generate/validate/parse, version selector
│   │       ├── yaml.go             # YAML TUI: to-json/to-yaml toggle
│   │       ├── timestamp.go        # Timestamp TUI: direction toggle, format selector
│   │       ├── number_base.go      # Number base TUI: input, multi-base output display
│   │       ├── regex.go            # Regex TUI: pattern input, test string, highlighted matches
│   │       ├── html_entity.go      # HTML entity TUI: encode/decode toggle
│   │       ├── password.go         # Password TUI: length slider, option checkboxes, generate button
│   │       └── lorem.go            # Lorem TUI: unit selector, count input
│   │
│   ├── web/                        # Web UI (chi + templ + HTMX)
│   │   ├── server.go               # Chi router setup, middleware, static files, server start
│   │   ├── handlers/               # HTTP handlers — one per tool
│   │   │   ├── index.go            # GET / — renders tool list
│   │   │   ├── base64.go           # GET/POST /tools/base64
│   │   │   ├── jwt.go              # GET/POST /tools/jwt
│   │   │   ├── json.go             # GET/POST /tools/json
│   │   │   ├── hash.go             # GET/POST /tools/hash
│   │   │   ├── url.go              # GET/POST /tools/url
│   │   │   ├── uuid.go             # GET/POST /tools/uuid
│   │   │   ├── yaml.go             # GET/POST /tools/yaml
│   │   │   ├── timestamp.go        # GET/POST /tools/timestamp
│   │   │   ├── number_base.go      # GET/POST /tools/number-base
│   │   │   ├── regex.go            # GET/POST /tools/regex
│   │   │   ├── html_entity.go      # GET/POST /tools/html-entity
│   │   │   ├── password.go         # GET/POST /tools/password
│   │   │   └── lorem.go            # GET/POST /tools/lorem
│   │   ├── templates/              # Templ templates
│   │   │   ├── layout.templ        # Base HTML layout: head, header, main, footer
│   │   │   ├── index.templ         # Tool list page template
│   │   │   └── tools/              # Per-tool templates
│   │   │       ├── base64.templ    # Base64 tool page
│   │   │       ├── jwt.templ       # JWT tool page
│   │   │       ├── json.templ      # JSON tool page
│   │   │       ├── hash.templ      # Hash tool page
│   │   │       ├── url.templ       # URL tool page
│   │   │       ├── uuid.templ      # UUID tool page
│   │   │       └── ...             # (one per tool)
│   │   └── static/                 # Embedded static assets
│   │       ├── style.css           # Custom stylesheet (~200 lines, dark theme)
│   │       ├── htmx.min.js         # Vendored HTMX 2.x
│   │       └── forge.js            # Clipboard helpers, textarea auto-resize (~20 lines)
│   │
│   └── desktop/                    # Desktop app (Wails v2) — deferred to post-v1
│       ├── main.go                 # Wails entry point (stub)
│       ├── app.go                  # Wails app bindings (stub)
│       └── frontend/               # Wails frontend assets
│           └── .gitkeep            # Placeholder — will reuse web templates
│
├── internal/                       # Internal packages — not importable by external code
│   ├── clipboard/                  # Cross-platform clipboard abstraction
│   │   ├── clipboard.go            # Clipboard interface + platform implementation (atotto/clipboard)
│   │   └── clipboard_test.go       # Integration tests (skipped in CI without display)
│   ├── version/                    # Build version info
│   │   └── version.go              # Version, Commit, Date variables (set via ldflags)
│   └── stdin/                      # Stdin reading utilities
│       └── stdin.go                # ReadStdin() — reads from stdin with timeout, detects pipe vs terminal
│
└── .github/                        # GitHub configuration
    └── workflows/
        ├── ci.yml                  # CI: test, lint, build on push/PR
        └── release.yml             # Release: goreleaser on tag push
```

---

## Key Structural Decisions

### Why `core/` instead of `pkg/`?

The Go community has moved away from `pkg/` as a convention. The golang-standards/project-layout repo (while not official) recommends `pkg/` for library code intended for external consumption. Forge's tool logic is not a reusable library — it is internal business logic. `core/` communicates intent more clearly: this is the domain logic, the heart of the application. It also avoids the confusion of "should I import `pkg/tools` from my other project?" — no, you shouldn't.

### Why `internal/` for clipboard and version?

Go's `internal/` package has compiler-enforced import restrictions: code in `internal/` can only be imported by code in the parent tree. This is exactly right for `clipboard` (platform-specific I/O that should not leak into `core/`) and `version` (build metadata that only `cmd/` needs). It prevents accidental coupling.

### Why `core/tools/` is flat (no subdirectories per tool)?

Each tool implementation is a single file with ~50–150 lines of pure functions. Subdirectories would add navigation overhead without organizational benefit. A flat package with one file per tool plus one test file per tool is the Go-idiomatic approach for this size. If a tool grows beyond ~300 lines, it should be refactored into its own sub-package.

### Why `ui/tui/views/` instead of `ui/tui/tools/`?

The directory is named `views/` to emphasize that these files contain TUI presentation logic (tea.Model implementations), not business logic. This naming makes the separation crystal clear: `core/tools/` = logic, `ui/tui/views/` = presentation.

### Why separate `cmd/` files per tool?

Each cobra command file is ~30–60 lines: flag definitions, argument parsing, calling `core/tools/`, formatting output. Keeping them in separate files makes it easy to add or remove CLI commands without touching unrelated code. The root.go file registers all subcommands.

### Why `ui/web/handlers/` is separate from `ui/web/templates/`?

Handlers contain Go HTTP logic (parsing form data, calling core functions, selecting templates). Templates contain HTML structure. Mixing them would violate separation of concerns. A handler imports a template; a template never imports a handler.

### Why `ui/desktop/` is a stub?

The desktop surface is explicitly deferred to post-v1. The directory exists to establish the pattern and prevent future restructuring, but contains only placeholder files. The plan is to reuse the web templates inside a Wails webview.

### Why `.github/workflows/` is included?

CI/CD is not optional infrastructure — it is part of the project. The release workflow (goreleaser on tag push) and CI workflow (test + lint on every push/PR) should exist from day one.

---

## Dependency Flow

```
cmd/ ──────────────────┐
  │                    │
  ├── imports ──► core/tools/
  ├── imports ──► core/registry/
  ├── imports ──► ui/tui/
  ├── imports ──► ui/web/
  ├── imports ──► internal/version/
  │
  │
ui/tui/ ───────────────┤
  ├── imports ──► core/tools/
  ├── imports ──► core/registry/
  ├── imports ──► core/detection/
  ├── imports ──► internal/clipboard/
  │
  │
ui/web/ ───────────────┤
  ├── imports ──► core/tools/
  ├── imports ──► core/registry/
  │
  │
core/ ─────────────────┤
  ├── core/registry/ imports ──► core/tools/
  ├── core/detection/ imports ──► core/registry/
  │                               core/tools/
  │                               internal/clipboard/
  │
  ╳ core/ NEVER imports ui/ or cmd/
  ╳ internal/ NEVER imports ui/ or cmd/
```

---

## File Count Summary

| Directory | Files | Purpose |
|-----------|-------|---------|
| Root | 7 | Config, entrypoint, docs |
| cmd/ | 17 | CLI commands |
| core/tools/ | 28 | Tool logic + tests (14 tools x 2 files) |
| core/registry/ | 3 | Registry + defaults + tests |
| core/detection/ | 2 | Detection engine + tests |
| ui/tui/ | 18 | TUI app + views + styles + keys |
| ui/web/ | ~30 | Server + handlers + templates + static |
| ui/desktop/ | 3 | Stubs |
| internal/ | 5 | Clipboard, version, stdin |
| .github/ | 2 | CI/CD workflows |
| **Total** | **~115** | |

This is a medium-sized Go project. The flat structure within packages keeps navigation simple while the top-level separation (core / ui / cmd / internal) enforces clean architecture boundaries.
