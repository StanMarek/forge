# TUI Surface Phase 1 — Design Spec

**Date:** 2026-03-15
**Status:** Approved
**Scope:** App shell, sidebar, base64 tool view, styles, keybindings, `forge tui` command

---

## Goal

Build the foundational TUI using Bubbletea v2: a two-panel layout with a sidebar listing tools and a tool panel showing the active tool. Only the Base64 view is functional — other tools show a placeholder. No clipboard detection, no search, no responsive breakpoints.

## Dependencies

- `charm.land/bubbletea/v2` — TUI framework
- `charm.land/lipgloss/v2` — Terminal styling
- `charm.land/bubbles/v2` — Pre-built components (list, textarea, key)
- `github.com/StanMarek/forge/core/tools` — Business logic
- `github.com/StanMarek/forge/core/registry` — Tool registry for sidebar

## Architecture

```
AppModel (ui/tui/app.go)
├── Sidebar (ui/tui/sidebar.go)  — bubbles list.Model
└── ToolView (interface)
    ├── Base64View (ui/tui/views/base64.go)
    └── PlaceholderView (ui/tui/views/placeholder.go)
```

All state flows through Bubbletea's Elm Architecture: `Init() → Update(msg) → View()`. No global mutable state.

## Package: `ui/tui/styles/styles.go`

Material-Darker color palette:

```go
var (
    Background = lipgloss.Color("#212121")
    Surface    = lipgloss.Color("#292929")
    Contrast   = lipgloss.Color("#1A1A1A")
    TextPrimary   = lipgloss.Color("#EEFFFF")
    TextSecondary = lipgloss.Color("#B0BEC5")
    TextMuted     = lipgloss.Color("#616161")
    Accent     = lipgloss.Color("#FF9800")
    Green      = lipgloss.Color("#C3E88D")
    Cyan       = lipgloss.Color("#89DDFF")
    Red        = lipgloss.Color("#FF5370")
)
```

Pre-built style definitions:
- `SidebarStyle` — fixed width, Surface background, right border
- `ActiveItemStyle` — Accent foreground, bold
- `NormalItemStyle` — TextSecondary foreground
- `CategoryStyle` — TextMuted foreground, uppercase
- `ToolPanelStyle` — padding, Background
- `TitleStyle` — TextPrimary, bold, underline
- `InputStyle` / `OutputStyle` — bordered boxes
- `ErrorStyle` — Red foreground, red border
- `StatusBarStyle` — Contrast background, TextMuted foreground

All colors should use `lipgloss.AdaptiveColor` where appropriate for light/dark terminal support, but default to the dark palette.

## Package: `ui/tui/keys/keys.go`

```go
type KeyMap struct {
    SwitchPanel key.Binding  // tab
    Quit        key.Binding  // q (sidebar only)
    ForceQuit   key.Binding  // ctrl+c
    Select      key.Binding  // enter
    Help        key.Binding  // ? (reserved, not functional yet)
}
```

Global `Keys` variable with default bindings. Uses `bubbles/key` package.

## Package: `ui/tui/views/view.go`

```go
type ToolView interface {
    Init() tea.Cmd
    Update(msg tea.Msg) (ToolView, tea.Cmd)
    View() string
    SetSize(width, height int)
}
```

Every tool view implements this interface. `AppModel` holds a `ToolView` and delegates `Update`/`View` to it.

## Model: `ui/tui/app.go` — AppModel

```go
type focus int

const (
    focusSidebar focus = iota
    focusTool
)

type AppModel struct {
    sidebar    Sidebar
    activeView views.ToolView
    focus      focus
    width      int
    height     int
    registry   *registry.Registry
}
```

### Behavior

- **Init:** Creates sidebar from `registry.Default()`, sets Base64View as initial active view.
- **WindowSizeMsg:** Updates `width`/`height`, recalculates sidebar and tool panel dimensions, calls `SetSize` on both.
- **Tab key:** Toggles `focus` between `focusSidebar` and `focusTool`.
- **When sidebar focused:** Key events go to sidebar. `enter` on a tool → calls `createView(toolID)` → swaps `activeView`. `q` → quits.
- **When tool focused:** Key events go to `activeView.Update()`. `q` types into text fields (not quit).
- **View:** `lipgloss.JoinHorizontal(lipgloss.Top, sidebarView, toolView)` with the focused panel getting a brighter border.

### Tool switching message

```go
type ToolSelectedMsg struct {
    ToolID string
}
```

Sidebar emits this when `enter` is pressed. AppModel receives it and swaps the view.

### createView factory

```go
func createView(toolID string, width, height int) views.ToolView
```

Returns `Base64View` for "base64", `PlaceholderView` for everything else.

## Model: `ui/tui/sidebar.go` — Sidebar

Wraps `bubbles/list.Model` configured with:
- Category headers rendered as non-selectable styled items (TextMuted, uppercase)
- Tool items rendered with ActiveItemStyle when selected, NormalItemStyle otherwise
- Fixed width: 20 characters
- No filtering (search is Phase 2)
- Accent color highlight on selected item

### Item types

```go
type toolItem struct {
    tool tools.Tool
}

type categoryItem struct {
    name string
}
```

Both implement `list.Item`. Category items are rendered differently and cannot be selected (cursor skips them).

### Tool list order

Built from `registry.Default()` grouped by category:
1. **Encoders:** Base64, JWT, URL, HTML Entity (only first 3 exist)
2. **Formatters:** JSON
3. **Generators:** Hash, UUID

## Model: `ui/tui/views/base64.go` — Base64View

### State

```go
type Base64View struct {
    input      textarea.Model
    output     textarea.Model
    mode       int  // 0 = encode, 1 = decode
    urlSafe    bool
    width      int
    height     int
    err        string
}
```

### Behavior

- **Input textarea:** Editable, receives key events when tool panel is focused.
- **Output textarea:** Read-only, displays result.
- **Mode toggle:** `e` for encode, `d` for decode (when not typing in textarea). Or use `tab` within the tool to cycle between input controls — but simplest approach: mode is toggled via a key binding shown in the view.
- **Live processing:** Every time input changes, calls `tools.Base64Encode()` or `tools.Base64Decode()` and updates output.
- **URL-safe toggle:** `u` key toggles the `urlSafe` bool.
- **Error display:** If decode fails, `err` is set and output area shows the error styled with `ErrorStyle`.
- **View layout:** Title → mode indicator → input textarea → output textarea (or error), stacked vertically.

### Key handling within Base64View

- All printable keys → input textarea
- `ctrl+e` — switch to encode mode
- `ctrl+d` — switch to decode mode
- `ctrl+u` — toggle URL-safe

These keybindings avoid conflicting with text input. Regular `e`/`d`/`u` would be swallowed by the textarea.

## Model: `ui/tui/views/placeholder.go` — PlaceholderView

Simple view that shows:
```
<Tool Name>

Coming soon.
```

Implements `ToolView` interface. `Update` is a no-op. `SetSize` stores dimensions.

## Integration: `cmd/tui.go`

```go
var tuiCmd = &cobra.Command{
    Use:   "tui",
    Short: "Launch the interactive terminal UI",
    Run: func(cmd *cobra.Command, args []string) {
        app := tui.New()
        p := tea.NewProgram(app, tea.WithAltScreen())
        if _, err := p.Run(); err != nil {
            fmt.Fprintf(os.Stderr, "error: %v\n", err)
            os.Exit(1)
        }
    },
}
```

Registered in `cmd/root.go` via `rootCmd.AddCommand(tuiCmd)`.

Uses `tea.WithAltScreen()` to take over the terminal.

## File Manifest

```
ui/tui/app.go
ui/tui/sidebar.go
ui/tui/styles/styles.go
ui/tui/keys/keys.go
ui/tui/views/view.go
ui/tui/views/base64.go
ui/tui/views/placeholder.go
cmd/tui.go
```

Modified: `cmd/root.go` (add `tuiCmd` registration)

## Testing Strategy

- TUI code is notoriously hard to unit test. Bubbletea's `teatest` package exists but is overkill for Phase 1.
- **Manual testing:** `go run . tui`, navigate sidebar, select tools, type in base64 input, verify output updates.
- **Core logic is already tested** (200 tests). The TUI is a thin presentation layer calling tested functions.
- **Styles** are visual — verify by running the app.

## Out of Scope (Phase 2+)

- Remaining 5 tool views (jwt, json, hash, url, uuid)
- Sidebar search / filtering
- Clipboard detection banner
- Help overlay (`?` key)
- Responsive breakpoints (narrow <80, wide >160)
- Copy/Paste/Clear buttons
- `forge` (no subcommand) → auto-launch TUI
