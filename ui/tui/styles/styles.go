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
