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
