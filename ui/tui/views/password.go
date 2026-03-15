package views

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	"github.com/StanMarek/forge/core/tools"
	"github.com/StanMarek/forge/ui/tui/styles"
)

// PasswordView is the TUI view for password generation.
type PasswordView struct {
	output    viewport.Model
	length    int
	lowercase bool
	uppercase bool
	digits    bool
	symbols   bool
	width     int
	height    int
	err       string
}

// NewPasswordView creates a new Password Generator tool view.
func NewPasswordView() *PasswordView {
	vp := viewport.New()

	v := &PasswordView{
		output:    vp,
		length:    16,
		lowercase: true,
		uppercase: true,
		digits:    true,
		symbols:   false,
	}
	v.generate()
	return v
}

func (v *PasswordView) Init() tea.Cmd {
	return nil
}

func (v *PasswordView) Update(msg tea.Msg) (ToolView, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+g"))):
			v.generate()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+l"))):
			v.lowercase = !v.lowercase
			v.generate()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+u"))):
			v.uppercase = !v.uppercase
			v.generate()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+d"))):
			v.digits = !v.digits
			v.generate()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+s"))):
			v.symbols = !v.symbols
			v.generate()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+up"))):
			if v.length < 256 {
				v.length++
				v.generate()
			}
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+down"))):
			if v.length > 1 {
				v.length--
				v.generate()
			}
			return v, nil
		}
	}

	return v, nil
}

func (v *PasswordView) generate() {
	result := tools.PasswordGenerate(v.length, v.uppercase, v.lowercase, v.digits, v.symbols, "")
	if result.Error != "" {
		v.err = result.Error
		v.output.SetContent("")
	} else {
		v.err = ""
		v.output.SetContent(result.Output)
	}
}

func (v *PasswordView) View() string {
	title := styles.TitleStyle.Render("Password Generator")

	lengthStr := styles.LabelStyle.Render(fmt.Sprintf("Length: %d", v.length))

	lowercaseStr := styles.CheckboxOffStyle.Render("○ Lowercase")
	if v.lowercase {
		lowercaseStr = styles.CheckboxOnStyle.Render("● Lowercase")
	}
	uppercaseStr := styles.CheckboxOffStyle.Render("○ Uppercase")
	if v.uppercase {
		uppercaseStr = styles.CheckboxOnStyle.Render("● Uppercase")
	}
	digitsStr := styles.CheckboxOffStyle.Render("○ Digits")
	if v.digits {
		digitsStr = styles.CheckboxOnStyle.Render("● Digits")
	}
	symbolsStr := styles.CheckboxOffStyle.Render("○ Symbols")
	if v.symbols {
		symbolsStr = styles.CheckboxOnStyle.Render("● Symbols")
	}

	options := fmt.Sprintf("%s    %s  %s  %s  %s", lengthStr, lowercaseStr, uppercaseStr, digitsStr, symbolsStr)

	var outputSection string
	if v.err != "" {
		outputSection = styles.LabelStyle.Render("Error:") + "\n" + styles.ErrorTextStyle.Render(v.err)
	} else {
		outputSection = styles.LabelStyle.Render("Password:") + "\n" + v.output.View()
	}

	return fmt.Sprintf("%s\n\n%s\n\n%s", title, options, outputSection)
}

func (v *PasswordView) KeyHints() string {
	return hint("ctrl+g", "generate") + "  " +
		hint("ctrl+l", "lower") + "  " +
		hint("ctrl+u", "upper") + "  " +
		hint("ctrl+d", "digits") + "  " +
		hint("ctrl+s", "symbols") + "  " +
		hint("ctrl+↑↓", "length")
}

func (v *PasswordView) SetSize(width, height int) {
	v.width = width
	v.height = height

	outputHeight := max(height-8, 3)

	v.output.SetWidth(width - 4)
	v.output.SetHeight(outputHeight)
}
