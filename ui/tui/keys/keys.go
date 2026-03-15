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
