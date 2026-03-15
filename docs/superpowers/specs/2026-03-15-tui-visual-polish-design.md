# TUI Visual Polish — Design Spec

**Date:** 2026-03-15
**Status:** Approved
**Scope:** Visual overhaul of all TUI components — styles, borders, layout, status bar

---

## Goal

Transform the TUI from a flat, 2000s-looking terminal app into a modern, lazygit-quality interface. Key changes: rounded borders, two-tier focus coloring, bordered panels with titles, pill-style mode indicators, frameless status bar.

## Reference

Lazygit's design patterns:
- Rounded borders (`╭╮╰╯`) universally
- Bold+color active border as sole focus indicator
- Muted inactive borders
- Frameless bottom status bar
- Semantic color use (accent for active, muted for inactive)
- Full-width selected row background in active panel

## Changes

### 1. styles/styles.go — Complete Rewrite

**Colors** — keep Material-Darker palette, add new surface colors:

```go
var (
    Background    = lipgloss.Color("#212121")
    Surface       = lipgloss.Color("#292929")
    SurfaceLight  = lipgloss.Color("#2C2C2C") // lifted surface for selected items
    Contrast      = lipgloss.Color("#1A1A1A")
    TextPrimary   = lipgloss.Color("#EEFFFF")
    TextSecondary = lipgloss.Color("#B0BEC5")
    TextMuted     = lipgloss.Color("#616161")
    Accent        = lipgloss.Color("#FF9800")
    Green         = lipgloss.Color("#C3E88D")
    Cyan          = lipgloss.Color("#89DDFF")
    Red           = lipgloss.Color("#FF5370")
    Yellow        = lipgloss.Color("#FFCB6B")
)
```

**Border styles** — all rounded:

```go
var Border = lipgloss.RoundedBorder()
```

**Panel styles:**

```go
// Focused panel: accent border + bold
FocusedBorderStyle = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(Accent)

// Unfocused panel: muted border
UnfocusedBorderStyle = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(TextMuted)
```

**Sidebar styles:**

```go
// Selected item: accent text + subtle background fill
ActiveItemStyle = lipgloss.NewStyle().
    Foreground(Accent).
    Bold(true).
    Background(SurfaceLight).
    PaddingLeft(1).
    PaddingRight(1)

// Normal item: secondary text, no background
NormalItemStyle = lipgloss.NewStyle().
    Foreground(TextSecondary).
    PaddingLeft(2)

// Category header: muted, uppercase
CategoryStyle = lipgloss.NewStyle().
    Foreground(TextMuted).
    Bold(true).
    PaddingLeft(1).
    MarginTop(1)
```

**Mode pill styles** (for encode/decode/format toggles):

```go
// Active mode: inverted pill — accent bg, dark text
ModeActivePill = lipgloss.NewStyle().
    Foreground(Background).
    Background(Accent).
    Bold(true).
    Padding(0, 1)

// Inactive mode: muted text, no background
ModeInactivePill = lipgloss.NewStyle().
    Foreground(TextMuted).
    Padding(0, 1)
```

**Input/Output box styles:**

```go
// Labeled box for input/output areas
InputBoxStyle = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(TextMuted)

// Error box
ErrorBoxStyle = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(Red)
```

**Status bar** (frameless):

```go
StatusKeyStyle = lipgloss.NewStyle().
    Foreground(Cyan)

StatusValueStyle = lipgloss.NewStyle().
    Foreground(Green)

StatusBarStyle = lipgloss.NewStyle().
    Foreground(TextMuted)
```

**Other:**

```go
TitleStyle = lipgloss.NewStyle().
    Foreground(TextPrimary).
    Bold(true)

LabelStyle = lipgloss.NewStyle().
    Foreground(TextSecondary)

ErrorTextStyle = lipgloss.NewStyle().
    Foreground(Red).
    Bold(true)

CheckboxOnStyle = lipgloss.NewStyle().
    Foreground(Accent)

CheckboxOffStyle = lipgloss.NewStyle().
    Foreground(TextMuted)
```

### 2. app.go — Layout Changes

**Panel rendering:**

```go
// Sidebar wrapped in rounded border box
sidebarBorder := styles.UnfocusedBorderStyle
if m.focus == focusSidebar {
    sidebarBorder = styles.FocusedBorderStyle
}
sidebarBox := sidebarBorder.
    Width(sidebarWidth).
    Height(m.height - 3). // -3 for status bar
    Render(m.sidebar.View())

// Tool panel wrapped in rounded border box with tool name as title
toolBorder := styles.UnfocusedBorderStyle
if m.focus == focusTool {
    toolBorder = styles.FocusedBorderStyle
}
toolBox := toolBorder.
    Width(m.width - sidebarWidth - 4).
    Height(m.height - 3).
    Render(m.activeView.View())

panels := lipgloss.JoinHorizontal(lipgloss.Top, sidebarBox, toolBox)
```

**Status bar** — rendered separately below panels:

```go
// Frameless status bar at the bottom
leftHints := styles.StatusKeyStyle.Render("tab") + styles.StatusBarStyle.Render(": switch  ") +
    styles.StatusKeyStyle.Render("↑↓") + styles.StatusBarStyle.Render(": navigate  ") +
    styles.StatusKeyStyle.Render("enter") + styles.StatusBarStyle.Render(": select  ") +
    styles.StatusKeyStyle.Render("q") + styles.StatusBarStyle.Render(": quit")
rightInfo := styles.StatusValueStyle.Render("forge " + version.Version)

// Pad between left and right to fill width
statusBar := lipgloss.JoinHorizontal(lipgloss.Top, leftHints, gap, rightInfo)

content := lipgloss.JoinVertical(lipgloss.Left, panels, statusBar)
```

### 3. sidebar.go — Visual Improvements

- Selected item gets `ActiveItemStyle` with `SurfaceLight` background (full-width fill)
- Category headers followed by no separator (the spacing from `MarginTop(1)` is enough with the new borders)
- Prefix: `▸` for selected, space for others (already exists, keep it)

### 4. All 6 views — Consistent Patterns

Each view's `View()` method changes:

**Mode indicators** — replace text radio buttons with pills:

```go
// Before:
styles.ModeActiveStyle.Render("(●) Encode") + "  " + styles.ModeInactiveStyle.Render("( ) Decode")

// After:
styles.ModeActivePill.Render("Encode") + " " + styles.ModeInactivePill.Render("Decode")
```

**Checkbox toggles** — styled:

```go
// Before:
"☐ URL-safe" / "☑ URL-safe"

// After:
styles.CheckboxOnStyle.Render("● URL-safe") / styles.CheckboxOffStyle.Render("○ URL-safe")
```

**Status bar in views** — REMOVED from individual views. The app-level status bar handles keybinding hints now. Each view only renders its content; the app wraps it in a bordered box.

**Tool-specific keybinding hints** — passed up to app via a method on ToolView:

Add to `ToolView` interface:
```go
type ToolView interface {
    Init() tea.Cmd
    Update(msg tea.Msg) (ToolView, tea.Cmd)
    View() string
    SetSize(width, height int)
    KeyHints() string  // e.g. "ctrl+e: encode  ctrl+d: decode"
}
```

The app-level status bar combines the tool's `KeyHints()` with global hints.

### 5. Error display

Errors wrapped in `ErrorBoxStyle` (rounded red border):

```go
if v.err != "" {
    errorBox := styles.ErrorBoxStyle.
        Width(v.width - 4).
        Render(styles.ErrorTextStyle.Render(v.err))
    // render errorBox instead of output
}
```

## File Manifest

```
Modify: ui/tui/styles/styles.go    — complete rewrite
Modify: ui/tui/app.go              — bordered panels, status bar, KeyHints
Modify: ui/tui/sidebar.go          — selected item background
Modify: ui/tui/views/view.go       — add KeyHints() to interface
Modify: ui/tui/views/base64.go     — pills, remove status bar, add KeyHints
Modify: ui/tui/views/jwt.go        — same
Modify: ui/tui/views/json.go       — same
Modify: ui/tui/views/hash.go       — same
Modify: ui/tui/views/url.go        — same
Modify: ui/tui/views/uuid.go       — same
Modify: ui/tui/views/placeholder.go — add KeyHints
```

## Testing

Manual only — run `go run . tui` and visually verify:
- Rounded borders on both panels
- Orange border on focused panel, muted on unfocused
- Pill-style mode indicators
- Selected sidebar item has subtle background
- Frameless status bar at bottom
- Error states show red-bordered box
