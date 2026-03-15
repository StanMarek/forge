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

func (m AppModel) View() tea.View {
	sidebarStyle := styles.SidebarStyle
	if m.focus == focusSidebar {
		sidebarStyle = styles.SidebarFocusedStyle
	}
	sidebarView := sidebarStyle.Height(m.height - 2).Render(m.sidebar.View())

	toolView := styles.ToolPanelStyle.Render(m.activeView.View())

	content := lipgloss.JoinHorizontal(lipgloss.Top, sidebarView, toolView)

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

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
