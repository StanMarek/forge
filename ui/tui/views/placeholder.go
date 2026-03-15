package views

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/StanMarek/forge/ui/tui/styles"
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

func (p *PlaceholderView) Update(_ tea.Msg) (ToolView, tea.Cmd) {
	return p, nil
}

func (p *PlaceholderView) View() string {
	title := styles.TitleStyle.Render(p.name)
	msg := styles.LabelStyle.Render("Coming soon.")
	return fmt.Sprintf("%s\n\n%s", title, msg)
}

func (p *PlaceholderView) SetSize(width, height int) {
	p.width = width
	p.height = height
}
