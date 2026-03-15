# TUI Mockups

**Project:** Forge — Developer Utility Toolkit
**Date:** 2026-03-15

---

## Mockup 1: Main Screen — Base64 Tool Active

Standard terminal: 120 columns x 30 rows

```
┌──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│  ⚒ forge                                                                                              v0.1.0       │
├─────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────┤
│                     │                                                                                              │
│  ENCODERS           │  Base64 Encode / Decode                                                                      │
│  ▸ Base64           │  ──────────────────────────────────────────────────────────────────────────────────           │
│    JWT              │                                                                                              │
│    URL              │  Mode:  (●) Encode  ( ) Decode          ☐ URL-safe                                           │
│    HTML Entity      │                                                                                              │
│                     │  Input:                                                                                      │
│  FORMATTERS         │  ┌──────────────────────────────────────────────────────────────────────────────┐             │
│    JSON             │  │ Hello, World!                                                                │             │
│                     │  │                                                                              │             │
│  GENERATORS         │  │                                                                              │             │
│    Hash             │  └──────────────────────────────────────────────────────────────────────────────┘             │
│    UUID             │                                                                                              │
│    Password         │  Output:                                                                                     │
│    Lorem Ipsum      │  ┌──────────────────────────────────────────────────────────────────────────────┐             │
│                     │  │ SGVsbG8sIFdvcmxkIQ==                                                        │             │
│  CONVERTERS         │  │                                                                              │             │
│    YAML             │  │                                                                              │             │
│    Timestamp        │  └──────────────────────────────────────────────────────────────────────────────┘             │
│    Number Base      │                                                                                              │
│                     │  ┌────────┐  ┌────────┐  ┌───────┐                                                           │
│  TESTERS            │  │  Copy  │  │ Paste  │  │ Clear │                                                           │
│    Regex            │  └────────┘  └────────┘  └───────┘                                                           │
│                     │                                                                                              │
│  ╭─────────────────╮│                                                                                              │
│  │ 🔍 search...    ││                                                                                              │
│  ╰─────────────────╯│                                                                                              │
├─────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────┤
│  tab: switch panel   /: search   ↑↓: navigate   enter: select   ctrl+l: clear   ?: help   q: quit                  │
└──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

**Active keybindings:**

- `tab` — switch focus between sidebar and tool panel
- `↑` / `↓` — navigate tool list (when sidebar focused) or scroll content (when tool focused)
- `enter` — select highlighted tool (sidebar) or activate focused element (tool panel)
- `/` — focus search input in sidebar
- `ctrl+l` — clear input and output in active tool
- `ctrl+c` — copy output to clipboard
- `ctrl+v` — paste from clipboard into input
- `?` — show help overlay
- `q` — quit (only when sidebar focused; in tool panel, `q` types into text input)

**Behaviour notes:**

- The sidebar uses the bubbles `list` component with category headers rendered as non-selectable items.
- The currently selected tool (`▸ Base64`) is highlighted with the accent color (#7C3AED).
- Input and output are both `textarea` components. Input is editable; output is read-only.
- Output updates live as the user types (debounced at 100ms for expensive operations like JSON formatting).
- Radio buttons for mode selection (Encode/Decode) and checkboxes (URL-safe) are navigable with arrow keys when the tool panel is focused.

---

## Mockup 2: Sidebar Search Active

```
┌──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│  ⚒ forge                                                                                              v0.1.0       │
├─────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────┤
│                     │                                                                                              │
│  ╭─────────────────╮│  Base64 Encode / Decode                                                                      │
│  │ 🔍 json█        ││  ──────────────────────────────────────────────────────────────────────────────────           │
│  ╰─────────────────╯│                                                                                              │
│                     │  (tool content remains visible but dimmed)                                                    │
│  Matching:          │                                                                                              │
│  ▸ JSON Formatter   │                                                                                              │
│    JSON/YAML Conv.  │                                                                                              │
│    JSONPath Tester  │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
├─────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────┤
│  esc: cancel search   ↑↓: navigate results   enter: select tool                                                     │
└──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

**Active keybindings:**

- `esc` — cancel search, restore full tool list
- `↑` / `↓` — navigate filtered results
- `enter` — select highlighted tool and exit search mode
- Any character — updates search filter in real-time

**Behaviour notes:**

- Search box moves to the top of the sidebar when active.
- Category headers are hidden during search — only matching tools are shown as a flat list.
- Search matches against tool name, ID, and keywords (fuzzy matching via `list` component's built-in filter).
- The tool panel content is dimmed but remains visible (not replaced with a blank screen).
- Empty search results show: "No tools match your search."

---

## Mockup 3: Smart Detection Banner

```
┌──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│  ⚒ forge                                                                                              v0.1.0       │
├──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│  ┌────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐  │
│  │  📋 Clipboard detected: JWT token  →  Press enter to open JWT Decoder, esc to dismiss                        │  │
│  └────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘  │
├─────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────┤
│                     │                                                                                              │
│  ENCODERS           │  Base64 Encode / Decode                                                                      │
│  ▸ Base64           │  ──────────────────────────────────────────────────────────────────────────────────           │
│    JWT              │                                                                                              │
│    ...              │  (current tool content)                                                                      │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
│                     │                                                                                              │
├─────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────┤
│  enter: open suggested tool   esc: dismiss   tab: switch panel   ?: help   q: quit                                  │
└──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

**Active keybindings:**

- `enter` — switch to the detected tool with clipboard content pre-filled as input
- `esc` — dismiss the detection banner; resume current tool
- All normal keybindings remain active

**Behaviour notes:**

- The detection banner appears between the title bar and the sidebar/content area.
- It uses a distinct background color (dark violet, semi-transparent feel via lipgloss) to draw attention without being jarring.
- The banner auto-dismisses after 10 seconds if no action is taken.
- If multiple tools match, the banner shows the highest-priority match and a count: "JWT token (also matches: Base64) → Press enter..."
- Detection polling runs every 500ms via a background `tea.Cmd` that compares clipboard content against the previous value.
- The banner does not appear if the user is actively typing in an input field (prevents interruption).

---

## Mockup 4: Tool with Error State

```
┌──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│  ⚒ forge                                                                                              v0.1.0       │
├─────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────┤
│                     │                                                                                              │
│  ENCODERS           │  JSON Formatter                                                                              │
│    Base64           │  ──────────────────────────────────────────────────────────────────────────────────           │
│    JWT              │                                                                                              │
│    URL              │  Mode:  (●) Format  ( ) Minify  ( ) Validate     Indent: [2]                                 │
│    HTML Entity      │                                                                                              │
│                     │  Input:                                                                                      │
│  FORMATTERS         │  ┌──────────────────────────────────────────────────────────────────────────────┐             │
│  ▸ JSON             │  │ {"name": "forge", "version": 0.1,                                           │             │
│                     │  │                                                                              │             │
│  GENERATORS         │  │                                                                              │             │
│    Hash             │  └──────────────────────────────────────────────────────────────────────────────┘             │
│    UUID             │                                                                                              │
│    Password         │  ┌ Error ──────────────────────────────────────────────────────────────────────┐              │
│    Lorem Ipsum      │  │  ✗ Invalid JSON: unexpected end of input at position 38                     │              │
│                     │  │    Hint: Missing closing brace '}'. Check line 1.                           │              │
│  CONVERTERS         │  └─────────────────────────────────────────────────────────────────────────────┘              │
│    YAML             │                                                                                              │
│    Timestamp        │  Output:                                                                                     │
│    Number Base      │  ┌──────────────────────────────────────────────────────────────────────────────┐             │
│                     │  │ (no output — fix input errors above)                                        │             │
│  TESTERS            │  │                                                                              │             │
│    Regex            │  │                                                                              │             │
│                     │  └──────────────────────────────────────────────────────────────────────────────┘             │
│  ╭─────────────────╮│                                                                                              │
│  │ 🔍 search...    ││  ┌────────┐  ┌────────┐  ┌───────┐                                                          │
│  ╰─────────────────╯│  │  Copy  │  │ Paste  │  │ Clear │                                                          │
│                     │  └────────┘  └────────┘  └───────┘                                                           │
├─────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────┤
│  tab: switch panel   /: search   ↑↓: navigate   enter: select   ctrl+l: clear   ?: help   q: quit                  │
└──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

**Behaviour notes:**

- The error box appears between input and output areas, replacing the output header temporarily.
- Error border is rendered in red (lipgloss `Color("#EF4444")`), distinct from the normal UI chrome.
- Error messages include: (1) the error itself, (2) a hint when possible (e.g., "Missing closing brace").
- The output area shows a dimmed placeholder message instead of stale output when there's an error.
- Errors update live as the user types — when the input becomes valid, the error box disappears and output renders instantly.
- The `Copy` button is disabled (dimmed) when output is empty/error state.

---

## Mockup 5: Help Overlay

```
┌──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│  ⚒ forge                                                                                              v0.1.0       │
├──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                                                      │
│                          ┌─── Keybindings ──────────────────────────────────┐                                        │
│                          │                                                  │                                        │
│                          │  Navigation                                      │                                        │
│                          │    tab          Switch sidebar / tool panel       │                                        │
│                          │    ↑ / ↓        Navigate tool list or scroll      │                                        │
│                          │    enter        Select tool / activate element    │                                        │
│                          │    /            Focus search                      │                                        │
│                          │    esc          Cancel search / dismiss overlay   │                                        │
│                          │                                                  │                                        │
│                          │  Tool Actions                                    │                                        │
│                          │    ctrl+l       Clear input and output            │                                        │
│                          │    ctrl+c       Copy output to clipboard          │                                        │
│                          │    ctrl+v       Paste clipboard into input        │                                        │
│                          │                                                  │                                        │
│                          │  Application                                     │                                        │
│                          │    ?            Toggle this help overlay          │                                        │
│                          │    q            Quit (sidebar must be focused)    │                                        │
│                          │    ctrl+c       Force quit (from anywhere)        │                                        │
│                          │                                                  │                                        │
│                          │                                                  │                                        │
│                          │  Press ? or esc to close                         │                                        │
│                          └──────────────────────────────────────────────────┘                                        │
│                                                                                                                      │
│                                                                                                                      │
├──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│  ?: close help   esc: close help                                                                                     │
└──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

**Active keybindings:**

- `?` — close help overlay
- `esc` — close help overlay

**Behaviour notes:**

- The help overlay is rendered on top of the main UI (the main content is dimmed behind it).
- The overlay is centered horizontally and vertically in the terminal.
- It uses lipgloss `Place()` to center the help box within the available space.
- The help overlay captures all keystrokes — only `?` and `esc` are processed.
- Keybinding categories are grouped: Navigation, Tool Actions, Application.
- Built using the bubbles `help` component's key map, ensuring bindings stay in sync with actual behavior.

---

## Mockup 6: Narrow Terminal (< 80 columns)

When terminal width drops below 80 columns, the sidebar collapses.

```
┌──────────────────────────────────────────────────────────┐
│  ⚒ forge                                     v0.1.0     │
├──────────────────────────────────────────────────────────┤
│  ◀ Base64 Encode / Decode                     [Tools ▸] │
├──────────────────────────────────────────────────────────┤
│                                                          │
│  Mode: (●) Encode  ( ) Decode     ☐ URL-safe            │
│                                                          │
│  Input:                                                  │
│  ┌──────────────────────────────────────────────────┐    │
│  │ Hello, World!                                    │    │
│  │                                                  │    │
│  └──────────────────────────────────────────────────┘    │
│                                                          │
│  Output:                                                 │
│  ┌──────────────────────────────────────────────────┐    │
│  │ SGVsbG8sIFdvcmxkIQ==                             │    │
│  │                                                  │    │
│  └──────────────────────────────────────────────────┘    │
│                                                          │
│  ┌──────┐  ┌───────┐  ┌───────┐                         │
│  │ Copy │  │ Paste │  │ Clear │                         │
│  └──────┘  └───────┘  └───────┘                         │
│                                                          │
├──────────────────────────────────────────────────────────┤
│  ◀▸: tool list   ctrl+l: clear   ?: help   q: quit      │
└──────────────────────────────────────────────────────────┘
```

**Active keybindings:**

- `◀` / `▸` or `[` / `]` — previous/next tool (replaces sidebar navigation)
- `ctrl+l` — clear input and output
- `?` — help overlay
- `q` — quit
- All tool-specific keybindings remain active

**Behaviour notes:**

- Below 80 columns, the sidebar is hidden entirely.
- A compact header replaces it: `◀ Tool Name [Tools ▸]` with arrow navigation.
- Pressing `[Tools ▸]` (or `/`) opens a full-screen tool picker overlay (similar to search mode).
- The tool panel takes the full width minus a small margin.
- Below 40 columns, Forge prints a warning and suggests widening the terminal. It still renders but may clip content.
- The transition between sidebar and collapsed mode happens dynamically on `tea.WindowSizeMsg`.

---

## Mockup 7: Wide Terminal (> 160 columns)

```
┌──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│  ⚒ forge                                                                                                                                          v0.1.0      │
├───────────────────────┬──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│                       │                                                                                                                                        │
│  ENCODERS             │  Base64 Encode / Decode                                                                                                                │
│  ▸ Base64             │  ────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────                │
│    JWT                │                                                                                                                                        │
│    URL                │  Mode:  (●) Encode  ( ) Decode          ☐ URL-safe    ☐ No padding                                                                    │
│    HTML Entity        │                                                                                                                                        │
│                       │  ┌─ Input ────────────────────────────────────────────────────┐  ┌─ Output ───────────────────────────────────────────────────┐         │
│  FORMATTERS           │  │ Hello, World!                                              │  │ SGVsbG8sIFdvcmxkIQ==                                      │         │
│    JSON               │  │                                                            │  │                                                            │         │
│                       │  │                                                            │  │                                                            │         │
│  GENERATORS           │  │                                                            │  │                                                            │         │
│    Hash               │  │                                                            │  │                                                            │         │
│    UUID               │  │                                                            │  │                                                            │         │
│    Password           │  │                                                            │  │                                                            │         │
│    Lorem Ipsum        │  │                                                            │  │                                                            │         │
│                       │  │                                                            │  │                                                            │         │
│  CONVERTERS           │  │                                                            │  │                                                            │         │
│    YAML               │  └────────────────────────────────────────────────────────────┘  └────────────────────────────────────────────────────────────┘         │
│    Timestamp          │                                                                                                                                        │
│    Number Base        │  ┌────────┐  ┌────────┐  ┌───────┐           Input: 13 bytes   Output: 20 bytes   Encoding: UTF-8                                     │
│                       │  │  Copy  │  │ Paste  │  │ Clear │                                                                                                    │
│  TESTERS              │  └────────┘  └────────┘  └───────┘                                                                                                    │
│    Regex              │                                                                                                                                        │
│                       │                                                                                                                                        │
│  ╭─────────────────╮  │                                                                                                                                        │
│  │ 🔍 search...    │  │                                                                                                                                        │
│  ╰─────────────────╯  │                                                                                                                                        │
├───────────────────────┴──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│  tab: switch panel   /: search   ↑↓: navigate   enter: select   ctrl+l: clear   ?: help   q: quit                                                              │
└──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

**Behaviour notes:**

- Above 160 columns, input and output are displayed side-by-side instead of stacked vertically.
- The sidebar width remains fixed (max 23 characters) — extra width goes to the tool panel.
- A status bar appears below the I/O areas showing byte counts and encoding info.
- The side-by-side layout is particularly useful for encode/decode tools where you want to see input and output simultaneously.
- Tools that don't benefit from side-by-side (like UUID generator) can opt to center their content with generous margins instead.
- Layout breakpoints: < 80 = collapsed sidebar; 80–160 = stacked I/O; > 160 = side-by-side I/O.

---

## Layout Breakpoint Summary

| Terminal Width | Sidebar | I/O Layout | Notes |
|---------------|---------|------------|-------|
| < 40 cols | Hidden + warning | Stacked, clipped | Minimum usable size |
| 40–79 cols | Hidden, header nav | Stacked, full width | Compact mode |
| 80–159 cols | Visible (20 cols) | Stacked (input above output) | Standard mode |
| 160+ cols | Visible (23 cols) | Side-by-side (input left, output right) | Wide mode |

---

## Color Palette

| Element | Color | Hex | Usage |
|---------|-------|-----|-------|
| Accent / Selected | Violet | `#7C3AED` | Selected tool, active borders, buttons |
| Accent subtle | Light violet | `#A78BFA` | Hover states, secondary highlights |
| Error | Red | `#EF4444` | Error borders, error text |
| Success | Green | `#22C55E` | Successful validation indicators |
| Dimmed | Gray | `#6B7280` | Inactive text, disabled buttons, dimmed content |
| Border | Dark gray | `#374151` | Panel borders, separators |
| Background | Terminal default | — | Adapts to user's terminal theme |
| Text | Terminal default | — | Adapts to user's terminal theme |

All colors use lipgloss `AdaptiveColor` to support both light and dark terminals.
