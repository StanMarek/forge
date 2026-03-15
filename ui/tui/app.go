package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/key"
	"charm.land/lipgloss/v2"
	"github.com/StanMarek/forge/core/registry"
	"github.com/StanMarek/forge/internal/version"
	"github.com/StanMarek/forge/ui/tui/keys"
	"github.com/StanMarek/forge/ui/tui/styles"
	"github.com/StanMarek/forge/ui/tui/views"
)

type focus int

const (
	focusSidebar focus = iota
	focusTool
)

const sidebarWidth = 24

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
		m.recalcSizes()
		return m, nil

	case tea.KeyPressMsg:
		if key.Matches(msg, keys.Keys.ForceQuit) {
			return m, tea.Quit
		}

		if key.Matches(msg, keys.Keys.SwitchPanel) {
			if m.focus == focusSidebar {
				m.focus = focusTool
			} else {
				m.focus = focusSidebar
			}
			return m, nil
		}

		if m.focus == focusSidebar {
			if key.Matches(msg, keys.Keys.Quit) {
				return m, tea.Quit
			}

			cmd := m.sidebar.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}

			if key.Matches(msg, keys.Keys.Select) {
				newID := m.sidebar.SelectedID()
				if newID != "" {
					m.activeView = createView(newID, m.sidebar.SelectedName(), 0, 0)
					m.recalcSizes()
					m.focus = focusTool
					initCmd := m.activeView.Init()
					if initCmd != nil {
						cmds = append(cmds, initCmd)
					}
				}
			}
			return m, tea.Batch(cmds...)
		}

		var cmd tea.Cmd
		m.activeView, cmd = m.activeView.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	}

	var cmd tea.Cmd
	m.activeView, cmd = m.activeView.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *AppModel) recalcSizes() {
	panelHeight := m.height - 3 // -3 for status bar + border
	toolWidth := m.width - sidebarWidth - 6
	toolHeight := panelHeight - 2 // -2 for border top/bottom

	m.sidebar.SetSize(sidebarWidth, panelHeight-2)
	if toolWidth > 0 && toolHeight > 0 {
		m.activeView.SetSize(toolWidth, toolHeight)
	}
}

func (m AppModel) View() tea.View {
	panelHeight := m.height - 3

	// Sidebar panel
	sidebarBorder := styles.UnfocusedBorderStyle
	if m.focus == focusSidebar {
		sidebarBorder = styles.FocusedBorderStyle
	}
	sidebarBox := sidebarBorder.
		Width(sidebarWidth).
		Height(panelHeight).
		Render(m.sidebar.View())

	// Tool panel
	toolBorder := styles.UnfocusedBorderStyle
	if m.focus == focusTool {
		toolBorder = styles.FocusedBorderStyle
	}
	toolWidth := m.width - sidebarWidth - 6
	toolBox := toolBorder.
		Width(toolWidth).
		Height(panelHeight).
		Render(styles.ToolPanelStyle.Render(m.activeView.View()))

	panels := lipgloss.JoinHorizontal(lipgloss.Top, sidebarBox, toolBox)

	// Status bar (frameless)
	statusBar := m.renderStatusBar()

	content := lipgloss.JoinVertical(lipgloss.Left, panels, statusBar)

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func (m AppModel) renderStatusBar() string {
	// Tool-specific hints
	toolHints := m.activeView.KeyHints()

	// Global hints
	globalHints := hint("tab", "switch") + "  " +
		hint("↑↓", "navigate") + "  " +
		hint("enter", "select") + "  " +
		hint("q", "quit")

	left := toolHints
	if left != "" {
		left += "  " + styles.StatusBarStyle.Render("│") + "  "
	}
	left += globalHints

	right := styles.StatusValueStyle.Render("forge " + version.Version)

	// Fill gap
	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right) - 2
	if gap < 1 {
		gap = 1
	}

	return " " + left + strings.Repeat(" ", gap) + right + " "
}

func hint(k, desc string) string {
	return styles.StatusKeyStyle.Render(k) + styles.StatusBarStyle.Render(":"+desc)
}

func createView(toolID, toolName string, width, height int) views.ToolView {
	var view views.ToolView
	switch toolID {
	case "base64":
		view = views.NewBase64View()
	case "jwt":
		view = views.NewJWTView()
	case "json":
		view = views.NewJSONView()
	case "hash":
		view = views.NewHashView()
	case "url":
		view = views.NewURLView()
	case "uuid":
		view = views.NewUUIDView()
	default:
		view = views.NewPlaceholder(toolName)
	}
	view.SetSize(width, height)
	return view
}
