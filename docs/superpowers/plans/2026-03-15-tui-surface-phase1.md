# TUI Surface Phase 1 Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the foundational Bubbletea v2 TUI with a two-panel layout (sidebar + tool panel) and a working Base64 tool view.

**Architecture:** `AppModel` (root tea.Model) owns a custom `Sidebar` and an active `ToolView` interface. Focus toggles between panels via tab. Only Base64View is functional; others show a placeholder. Uses Bubbletea v2 APIs: `tea.View` return type, `tea.KeyPressMsg`, no `tea.WithAltScreen`.

**Tech Stack:** `charm.land/bubbletea/v2`, `charm.land/lipgloss/v2`, `charm.land/bubbles/v2` (textarea, viewport, key), existing `core/tools/`

**Spec:** `docs/superpowers/specs/2026-03-15-tui-surface-phase1-design.md`

**IMPORTANT Bubbletea v2 changes from v1:**
- `View()` returns `tea.View`, NOT `string`. Use `tea.NewView(content)`.
- Alt screen: set `view.AltScreen = true` in `View()`, NOT `tea.WithAltScreen()`.
- Key events: use `tea.KeyPressMsg`, NOT `tea.KeyMsg`.
- `tea.NewProgram(model{})` — program options like WithAltScreen are REMOVED.
- Viewport: `viewport.New()` with option funcs, `vp.SetWidth(n)`, `vp.SetHeight(n)`, `vp.SetContent(s)`.
- Textarea: `textarea.New()`, sizing via setters not direct field access.

---

## File Structure

| File | Responsibility |
|------|---------------|
| `ui/tui/styles/styles.go` | Color palette + lipgloss style definitions |
| `ui/tui/keys/keys.go` | Keybinding definitions |
| `ui/tui/views/view.go` | ToolView interface |
| `ui/tui/views/placeholder.go` | "Coming soon" placeholder view |
| `ui/tui/views/base64.go` | Base64 encode/decode tool view |
| `ui/tui/sidebar.go` | Custom sidebar with categories + tool list |
| `ui/tui/app.go` | Root AppModel + New() constructor |
| `cmd/tui.go` | `forge tui` cobra command |

---

## Chunk 1: Foundation — styles, keys, view interface, placeholder

### Task 1: Add Bubbletea v2 dependencies

**Files:**
- Modify: `go.mod`

- [ ] **Step 1: Add charm dependencies**

```bash
cd /Users/stanislawmarek/Desktop/coding/forge
go get charm.land/bubbletea/v2
go get charm.land/lipgloss/v2
go get charm.land/bubbles/v2
go mod tidy
```

- [ ] **Step 2: Verify**

Run: `go build ./...`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add go.mod go.sum
git commit -m "Add bubbletea v2, lipgloss v2, bubbles v2 dependencies"
```

### Task 2: Styles package

**Files:**
- Create: `ui/tui/styles/styles.go`

- [ ] **Step 1: Write styles.go**

```go
package styles

import "charm.land/lipgloss/v2"

// Material-Darker color palette
var (
	Background    = lipgloss.Color("#212121")
	Surface       = lipgloss.Color("#292929")
	Contrast      = lipgloss.Color("#1A1A1A")
	TextPrimary   = lipgloss.Color("#EEFFFF")
	TextSecondary = lipgloss.Color("#B0BEC5")
	TextMuted     = lipgloss.Color("#616161")
	Accent        = lipgloss.Color("#FF9800")
	Green         = lipgloss.Color("#C3E88D")
	Cyan          = lipgloss.Color("#89DDFF")
	Red           = lipgloss.Color("#FF5370")
)

// Sidebar styles
var (
	SidebarStyle = lipgloss.NewStyle().
			Width(20).
			Background(Surface).
			BorderRight(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(TextMuted)

	SidebarFocusedStyle = lipgloss.NewStyle().
				Width(20).
				Background(Surface).
				BorderRight(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(Accent)

	CategoryStyle = lipgloss.NewStyle().
			Foreground(TextMuted).
			Bold(true).
			PaddingLeft(1).
			MarginTop(1)

	ActiveItemStyle = lipgloss.NewStyle().
			Foreground(Accent).
			Bold(true).
			PaddingLeft(2)

	NormalItemStyle = lipgloss.NewStyle().
			Foreground(TextSecondary).
			PaddingLeft(2)
)

// Tool panel styles
var (
	ToolPanelStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			PaddingRight(2).
			PaddingTop(1)

	TitleStyle = lipgloss.NewStyle().
			Foreground(TextPrimary).
			Bold(true).
			MarginBottom(1)

	LabelStyle = lipgloss.NewStyle().
			Foreground(TextSecondary)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Red).
			Bold(true)

	ModeActiveStyle = lipgloss.NewStyle().
			Foreground(Accent).
			Bold(true)

	ModeInactiveStyle = lipgloss.NewStyle().
				Foreground(TextMuted)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(TextMuted).
			MarginTop(1)
)
```

- [ ] **Step 2: Verify it compiles**

Run: `go build ./ui/tui/styles/`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add ui/tui/styles/styles.go
git commit -m "Add TUI styles package with Material-Darker palette"
```

### Task 3: Keys package

**Files:**
- Create: `ui/tui/keys/keys.go`

- [ ] **Step 1: Write keys.go**

```go
package keys

import "charm.land/bubbles/v2/key"

// KeyMap defines the application keybindings.
type KeyMap struct {
	SwitchPanel key.Binding
	Quit        key.Binding
	ForceQuit   key.Binding
	Select      key.Binding
	Up          key.Binding
	Down        key.Binding
	Help        key.Binding
}

// Keys is the default set of keybindings.
var Keys = KeyMap{
	SwitchPanel: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch panel"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
	ForceQuit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "force quit"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
}
```

- [ ] **Step 2: Verify it compiles**

Run: `go build ./ui/tui/keys/`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add ui/tui/keys/keys.go
git commit -m "Add TUI keybindings package"
```

### Task 4: ToolView interface + PlaceholderView

**Files:**
- Create: `ui/tui/views/view.go`
- Create: `ui/tui/views/placeholder.go`

- [ ] **Step 1: Write view.go**

```go
package views

import tea "charm.land/bubbletea/v2"

// ToolView is the interface that all tool views implement.
// This is NOT tea.Model — only AppModel satisfies tea.Model.
// ToolView.Update returns (ToolView, tea.Cmd), not (tea.Model, tea.Cmd).
type ToolView interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (ToolView, tea.Cmd)
	View() string
	SetSize(width, height int)
}
```

- [ ] **Step 2: Write placeholder.go**

```go
package views

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// PlaceholderView shows "Coming soon" for unimplemented tools.
type PlaceholderView struct {
	name   string
	width  int
	height int
}

// NewPlaceholder creates a placeholder view for the given tool name.
func NewPlaceholder(name string) *PlaceholderView {
	return &PlaceholderView{name: name}
}

func (p *PlaceholderView) Init() tea.Cmd { return nil }

func (p *PlaceholderView) Update(msg tea.Msg) (ToolView, tea.Cmd) {
	return p, nil
}

func (p *PlaceholderView) View() string {
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#EEFFFF")).Render(p.name)
	msg := lipgloss.NewStyle().Foreground(lipgloss.Color("#616161")).Render("Coming soon.")
	return fmt.Sprintf("%s\n\n%s", title, msg)
}

func (p *PlaceholderView) SetSize(width, height int) {
	p.width = width
	p.height = height
}
```

- [ ] **Step 3: Verify it compiles**

Run: `go build ./ui/tui/views/`
Expected: no errors

- [ ] **Step 4: Commit**

```bash
git add ui/tui/views/
git commit -m "Add ToolView interface and placeholder view"
```

---

## Chunk 2: Sidebar + AppModel + Base64View

### Task 5: Sidebar

**Files:**
- Create: `ui/tui/sidebar.go`

- [ ] **Step 1: Write sidebar.go**

Custom sidebar (NOT bubbles/list — to support non-selectable category headers).

```go
package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/key"
	"github.com/StanMarek/forge/core/registry"
	"github.com/StanMarek/forge/core/tools"
	"github.com/StanMarek/forge/ui/tui/keys"
	"github.com/StanMarek/forge/ui/tui/styles"
)

type entryKind int

const (
	entryCategory entryKind = iota
	entryTool
)

type sidebarEntry struct {
	kind     entryKind
	label    string
	toolID   string
	toolInfo tools.Tool
}

// Sidebar is a custom sidebar model with category headers and tool items.
type Sidebar struct {
	entries  []sidebarEntry
	cursor   int
	width    int
	height   int
}

// NewSidebar creates a sidebar from the registry, grouped by category.
func NewSidebar(reg *registry.Registry) Sidebar {
	var entries []sidebarEntry
	categories := []string{"Encoders", "Formatters", "Generators"}

	for _, cat := range categories {
		toolsInCat := reg.ByCategory(cat)
		if len(toolsInCat) == 0 {
			continue
		}
		entries = append(entries, sidebarEntry{kind: entryCategory, label: cat})
		for _, t := range toolsInCat {
			entries = append(entries, sidebarEntry{
				kind:     entryTool,
				label:    t.Name(),
				toolID:   t.ID(),
				toolInfo: t,
			})
		}
	}

	s := Sidebar{entries: entries}
	// Set cursor to the first tool item
	for i, e := range entries {
		if e.kind == entryTool {
			s.cursor = i
			break
		}
	}
	return s
}

// SelectedID returns the tool ID of the currently selected item.
func (s *Sidebar) SelectedID() string {
	if s.cursor >= 0 && s.cursor < len(s.entries) && s.entries[s.cursor].kind == entryTool {
		return s.entries[s.cursor].toolID
	}
	return ""
}

// SelectedName returns the display name of the currently selected tool.
func (s *Sidebar) SelectedName() string {
	if s.cursor >= 0 && s.cursor < len(s.entries) && s.entries[s.cursor].kind == entryTool {
		return s.entries[s.cursor].label
	}
	return ""
}

// Update handles key events for the sidebar.
func (s *Sidebar) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, keys.Keys.Up):
			s.moveUp()
		case key.Matches(msg, keys.Keys.Down):
			s.moveDown()
		}
	}
	return nil
}

func (s *Sidebar) moveUp() {
	for i := s.cursor - 1; i >= 0; i-- {
		if s.entries[i].kind == entryTool {
			s.cursor = i
			return
		}
	}
}

func (s *Sidebar) moveDown() {
	for i := s.cursor + 1; i < len(s.entries); i++ {
		if s.entries[i].kind == entryTool {
			s.cursor = i
			return
		}
	}
}

// SetSize sets the sidebar dimensions.
func (s *Sidebar) SetSize(width, height int) {
	s.width = width
	s.height = height
}

// View renders the sidebar.
func (s *Sidebar) View() string {
	var b strings.Builder
	for i, entry := range s.entries {
		switch entry.kind {
		case entryCategory:
			b.WriteString(styles.CategoryStyle.Render(strings.ToUpper(entry.label)))
		case entryTool:
			if i == s.cursor {
				b.WriteString(styles.ActiveItemStyle.Render("▸ " + entry.label))
			} else {
				b.WriteString(styles.NormalItemStyle.Render("  " + entry.label))
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}
```

- [ ] **Step 2: Verify it compiles**

Run: `go build ./ui/tui/`
Expected: may fail because `app.go` doesn't exist yet — that's OK. Check for syntax errors only: `go vet ./ui/tui/sidebar.go` won't work on its own due to package dependencies. We'll verify after Task 7.

- [ ] **Step 3: Commit**

```bash
git add ui/tui/sidebar.go
git commit -m "Add custom sidebar with category headers and cursor navigation"
```

### Task 6: Base64View

**Files:**
- Create: `ui/tui/views/base64.go`

- [ ] **Step 1: Write base64.go**

```go
package views

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	"github.com/StanMarek/forge/core/tools"
	"github.com/StanMarek/forge/ui/tui/styles"
)

type base64Mode int

const (
	modeEncode base64Mode = iota
	modeDecode
)

// Base64View is the TUI view for Base64 encoding/decoding.
type Base64View struct {
	input     textarea.Model
	output    viewport.Model
	mode      base64Mode
	urlSafe   bool
	noPadding bool
	width     int
	height    int
	err       string
}

// NewBase64View creates a new Base64 tool view.
func NewBase64View() *Base64View {
	ti := textarea.New()
	ti.SetPlaceholder("Enter text to encode...")
	ti.Focus()

	vp := viewport.New()

	v := &Base64View{
		input:  ti,
		output: vp,
		mode:   modeEncode,
	}
	return v
}

func (v *Base64View) Init() tea.Cmd {
	return textarea.Blink
}

func (v *Base64View) Update(msg tea.Msg) (ToolView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+e"))):
			v.mode = modeEncode
			v.input.SetPlaceholder("Enter text to encode...")
			v.process()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+d"))):
			v.mode = modeDecode
			v.input.SetPlaceholder("Enter Base64 to decode...")
			v.process()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+u"))):
			v.urlSafe = !v.urlSafe
			v.process()
			return v, nil
		}
	}

	// Forward to input textarea
	var cmd tea.Cmd
	v.input, cmd = v.input.Update(msg)
	cmds = append(cmds, cmd)

	// Process on every input change
	v.process()

	return v, tea.Batch(cmds...)
}

func (v *Base64View) process() {
	input := v.input.Value()
	if input == "" {
		v.output.SetContent("")
		v.err = ""
		return
	}

	var result tools.Result
	switch v.mode {
	case modeEncode:
		result = tools.Base64Encode(input, v.urlSafe, v.noPadding)
	case modeDecode:
		result = tools.Base64Decode(input, v.urlSafe)
	}

	if result.Error != "" {
		v.err = result.Error
		v.output.SetContent("")
	} else {
		v.err = ""
		v.output.SetContent(result.Output)
	}
}

func (v *Base64View) View() string {
	// Title
	title := styles.TitleStyle.Render("Base64 Encode / Decode")

	// Mode indicator
	var modeStr string
	if v.mode == modeEncode {
		modeStr = styles.ModeActiveStyle.Render("(●) Encode") + "  " + styles.ModeInactiveStyle.Render("( ) Decode")
	} else {
		modeStr = styles.ModeInactiveStyle.Render("( ) Encode") + "  " + styles.ModeActiveStyle.Render("(●) Decode")
	}

	// URL-safe indicator
	urlSafeStr := "☐ URL-safe"
	if v.urlSafe {
		urlSafeStr = "☑ URL-safe"
	}
	options := fmt.Sprintf("Mode: %s    %s", modeStr, styles.LabelStyle.Render(urlSafeStr))

	// Input
	inputLabel := styles.LabelStyle.Render("Input:")
	inputView := v.input.View()

	// Output or error
	var outputSection string
	if v.err != "" {
		outputSection = styles.LabelStyle.Render("Error:") + "\n" + styles.ErrorStyle.Render(v.err)
	} else {
		outputSection = styles.LabelStyle.Render("Output:") + "\n" + v.output.View()
	}

	// Status bar
	status := styles.StatusBarStyle.Render("ctrl+e: encode  ctrl+d: decode  ctrl+u: url-safe  tab: switch panel")

	return fmt.Sprintf("%s\n\n%s\n\n%s\n%s\n\n%s\n\n%s",
		title, options, inputLabel, inputView, outputSection, status)
}

func (v *Base64View) SetSize(width, height int) {
	v.width = width
	v.height = height

	// Input gets ~40% of available height, output gets the rest
	inputHeight := max((height-10)/2, 3)
	outputHeight := max((height-10)/2, 3)

	v.input.SetWidth(width - 4)
	v.input.SetHeight(inputHeight)
	v.output.SetWidth(width - 4)
	v.output.SetHeight(outputHeight)
}
```

- [ ] **Step 2: Commit**

```bash
git add ui/tui/views/base64.go
git commit -m "Add Base64 TUI view with live encode/decode"
```

### Task 7: AppModel

**Files:**
- Create: `ui/tui/app.go`

- [ ] **Step 1: Write app.go**

```go
package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/key"
	"charm.land/lipgloss/v2"
	"github.com/StanMarek/forge/core/registry"
	"github.com/StanMarek/forge/ui/tui/keys"
	"github.com/StanMarek/forge/ui/tui/styles"
	"github.com/StanMarek/forge/ui/tui/views"
)

type focus int

const (
	focusSidebar focus = iota
	focusTool
)

const sidebarWidth = 22

// AppModel is the root Bubbletea model for the Forge TUI.
type AppModel struct {
	sidebar    Sidebar
	activeView views.ToolView
	focus      focus
	width      int
	height     int
	reg        *registry.Registry
}

// New creates a new AppModel with the default registry.
func New() AppModel {
	reg := registry.Default()
	sidebar := NewSidebar(reg)
	return AppModel{
		sidebar:    sidebar,
		activeView: views.NewBase64View(),
		focus:      focusSidebar,
		reg:        reg,
	}
}

func (m AppModel) Init() tea.Cmd {
	return m.activeView.Init()
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		toolWidth := m.width - sidebarWidth - 2
		toolHeight := m.height - 2
		m.sidebar.SetSize(sidebarWidth, m.height-2)
		m.activeView.SetSize(toolWidth, toolHeight)
		return m, nil

	case tea.KeyPressMsg:
		// Force quit from anywhere
		if key.Matches(msg, keys.Keys.ForceQuit) {
			return m, tea.Quit
		}

		// Tab switches focus
		if key.Matches(msg, keys.Keys.SwitchPanel) {
			if m.focus == focusSidebar {
				m.focus = focusTool
			} else {
				m.focus = focusSidebar
			}
			return m, nil
		}

		if m.focus == focusSidebar {
			// Quit only from sidebar
			if key.Matches(msg, keys.Keys.Quit) {
				return m, tea.Quit
			}

			// Check for tool selection
			prevID := m.sidebar.SelectedID()
			cmd := m.sidebar.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}

			// Enter selects a tool
			if key.Matches(msg, keys.Keys.Select) {
				newID := m.sidebar.SelectedID()
				if newID != "" {
					toolWidth := m.width - sidebarWidth - 2
					toolHeight := m.height - 2
					m.activeView = createView(newID, m.sidebar.SelectedName(), toolWidth, toolHeight)
					m.focus = focusTool
					initCmd := m.activeView.Init()
					if initCmd != nil {
						cmds = append(cmds, initCmd)
					}
				}
			}
			_ = prevID
			return m, tea.Batch(cmds...)
		}

		// Tool panel is focused — forward to active view
		var cmd tea.Cmd
		m.activeView, cmd = m.activeView.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	}

	// Forward non-key messages to active view (e.g., blink)
	var cmd tea.Cmd
	m.activeView, cmd = m.activeView.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m AppModel) View() tea.View {
	// Sidebar
	sidebarStyle := styles.SidebarStyle
	if m.focus == focusSidebar {
		sidebarStyle = styles.SidebarFocusedStyle
	}
	sidebarView := sidebarStyle.Height(m.height - 2).Render(m.sidebar.View())

	// Tool panel
	toolView := styles.ToolPanelStyle.Render(m.activeView.View())

	content := lipgloss.JoinHorizontal(lipgloss.Top, sidebarView, toolView)

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

// createView returns the appropriate ToolView for the given tool ID.
func createView(toolID, toolName string, width, height int) views.ToolView {
	var view views.ToolView
	switch toolID {
	case "base64":
		view = views.NewBase64View()
	default:
		view = views.NewPlaceholder(toolName)
	}
	view.SetSize(width, height)
	return view
}
```

- [ ] **Step 2: Verify the whole TUI package compiles**

Run: `go build ./ui/tui/...`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add ui/tui/app.go
git commit -m "Add AppModel with two-panel layout and focus management"
```

---

## Chunk 3: Integration + smoke test

### Task 8: cmd/tui.go + root registration

**Files:**
- Create: `cmd/tui.go`
- Modify: `cmd/root.go`

- [ ] **Step 1: Write cmd/tui.go**

```go
package cmd

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/StanMarek/forge/ui/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the interactive terminal UI",
	Run: func(cmd *cobra.Command, args []string) {
		app := tui.New()
		p := tea.NewProgram(app)
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	},
}
```

- [ ] **Step 2: Add tuiCmd to root.go**

In `cmd/root.go`, add `rootCmd.AddCommand(tuiCmd)` to the `init()` function.

- [ ] **Step 3: Build and verify**

Run: `go build -o bin/forge .`
Expected: builds successfully

- [ ] **Step 4: Commit**

```bash
git add cmd/tui.go cmd/root.go
git commit -m "Add forge tui command to launch terminal UI"
```

### Task 9: Manual smoke test

- [ ] **Step 1: Launch the TUI**

Run: `./bin/forge tui`
Expected: Alt screen, sidebar on left with tool categories, Base64 view on right

- [ ] **Step 2: Test sidebar navigation**

Press `↓`/`↑` — cursor moves between tools, skips category headers.
Press `enter` on Base64 — should already be active, focus switches to tool panel.

- [ ] **Step 3: Test Base64 encoding**

Press `tab` to focus tool panel.
Type "Hello, World!" — output should show `SGVsbG8sIFdvcmxkIQ==` live.
Press `ctrl+d` — switch to decode mode.
Clear and type `SGVsbG8sIFdvcmxkIQ==` — should show `Hello, World!`

- [ ] **Step 4: Test focus toggle**

Press `tab` — focus returns to sidebar (border changes color).
Press `q` — app quits.

- [ ] **Step 5: Test placeholder**

Navigate to JWT in sidebar, press `enter` — should show "Coming soon."
Navigate to Hash, press `enter` — "Coming soon."

- [ ] **Step 6: Test force quit**

Re-launch, press `ctrl+c` from any panel — app quits.

- [ ] **Step 7: Fix any issues found, commit**

If fixes needed:
```bash
git add -A
git commit -m "Fix TUI issues found during smoke testing"
```

### Task 10: Push

- [ ] **Step 1: Run full test suite**

Run: `go test -count=1 ./...`
Expected: all existing tests pass (TUI has no automated tests)

- [ ] **Step 2: Push**

```bash
git push
```

---

## Task Dependency Summary

```
Task 1 (deps)
├── Task 2 (styles)     ─┐
├── Task 3 (keys)       ─┤
├── Task 4 (views)      ─┤
│                         ├── Task 5 (sidebar) ─┐
│                         │                      ├── Task 7 (app) → Task 8 (cmd) → Task 9 (smoke) → Task 10 (push)
│                         ├── Task 6 (base64)  ─┘
```

Tasks 2-4 are independent (parallelizable).
Tasks 5-6 depend on 2-4.
Task 7 depends on everything.
Tasks 8-10 are sequential.
