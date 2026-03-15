package styles

import "charm.land/lipgloss/v2"

// Material-Darker color palette
var (
	Background   = lipgloss.Color("#212121")
	Surface      = lipgloss.Color("#292929")
	SurfaceLight = lipgloss.Color("#2C2C2C")
	Contrast     = lipgloss.Color("#1A1A1A")

	TextPrimary   = lipgloss.Color("#EEFFFF")
	TextSecondary = lipgloss.Color("#B0BEC5")
	TextMuted     = lipgloss.Color("#616161")

	Accent = lipgloss.Color("#FF9800")
	Green  = lipgloss.Color("#C3E88D")
	Cyan   = lipgloss.Color("#89DDFF")
	Red    = lipgloss.Color("#FF5370")
	Yellow = lipgloss.Color("#FFCB6B")
)

// Border — rounded everywhere
var Border = lipgloss.RoundedBorder()

// Panel border styles
var (
	FocusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(Accent)

	UnfocusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(TextMuted)
)

// Sidebar styles
var (
	CategoryStyle = lipgloss.NewStyle().
			Foreground(TextMuted).
			Bold(true).
			PaddingLeft(1).
			MarginTop(1)

	ActiveItemStyle = lipgloss.NewStyle().
			Foreground(Accent).
			Bold(true).
			Background(SurfaceLight).
			PaddingLeft(1).
			PaddingRight(1)

	NormalItemStyle = lipgloss.NewStyle().
			Foreground(TextSecondary).
			PaddingLeft(2)
)

// Tool panel styles
var (
	ToolPanelStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1)

	TitleStyle = lipgloss.NewStyle().
			Foreground(TextPrimary).
			Bold(true)

	LabelStyle = lipgloss.NewStyle().
			Foreground(TextSecondary)
)

// Mode pill styles
var (
	ModeActivePill = lipgloss.NewStyle().
			Foreground(Background).
			Background(Accent).
			Bold(true).
			Padding(0, 1)

	ModeInactivePill = lipgloss.NewStyle().
				Foreground(TextMuted).
				Padding(0, 1)
)

// Checkbox styles
var (
	CheckboxOnStyle = lipgloss.NewStyle().
			Foreground(Accent)

	CheckboxOffStyle = lipgloss.NewStyle().
				Foreground(TextMuted)
)

// Input/Output box styles
var (
	InputBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(TextMuted)

	ErrorBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Red)

	ErrorTextStyle = lipgloss.NewStyle().
			Foreground(Red).
			Bold(true)
)

// Status bar styles (frameless bottom bar)
var (
	StatusKeyStyle = lipgloss.NewStyle().
			Foreground(Cyan).
			Bold(true)

	StatusValueStyle = lipgloss.NewStyle().
				Foreground(Green)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(TextMuted)
)
