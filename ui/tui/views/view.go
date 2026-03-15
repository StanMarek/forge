package views

import (
	tea "charm.land/bubbletea/v2"
	"github.com/StanMarek/forge/ui/tui/styles"
)

// ToolView is the interface that all tool views implement.
// This is NOT tea.Model — only AppModel satisfies tea.Model.
type ToolView interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (ToolView, tea.Cmd)
	View() string
	SetSize(width, height int)
	KeyHints() string // tool-specific keybinding hints for the status bar
}

// hint formats a keybinding hint for the status bar.
func hint(k, desc string) string {
	return styles.StatusKeyStyle.Render(k) + styles.StatusBarStyle.Render(":"+desc)
}
