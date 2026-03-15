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
	entries []sidebarEntry
	cursor  int
	width   int
	height  int
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
