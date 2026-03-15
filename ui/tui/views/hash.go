package views

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	"github.com/StanMarek/forge/core/tools"
	"github.com/StanMarek/forge/ui/tui/styles"
)

// HashView is the TUI view for hash generation.
type HashView struct {
	input     textarea.Model
	output    viewport.Model
	uppercase bool
	width     int
	height    int
}

// NewHashView creates a new Hash Generator tool view.
func NewHashView() *HashView {
	ti := textarea.New()
	ti.Placeholder = "Enter text to hash..."
	ti.Focus()

	vp := viewport.New()

	return &HashView{
		input:  ti,
		output: vp,
	}
}

func (v *HashView) Init() tea.Cmd {
	return textarea.Blink
}

func (v *HashView) Update(msg tea.Msg) (ToolView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+u"))):
			v.uppercase = !v.uppercase
			v.process()
			return v, nil
		}
	}

	var cmd tea.Cmd
	v.input, cmd = v.input.Update(msg)
	cmds = append(cmds, cmd)

	v.process()

	return v, tea.Batch(cmds...)
}

func (v *HashView) process() {
	input := v.input.Value()
	if input == "" {
		v.output.SetContent("")
		return
	}

	md5Result := tools.Hash(input, "md5", v.uppercase)
	sha1Result := tools.Hash(input, "sha1", v.uppercase)
	sha256Result := tools.Hash(input, "sha256", v.uppercase)
	sha512Result := tools.Hash(input, "sha512", v.uppercase)

	out := fmt.Sprintf("MD5:    %s\nSHA1:   %s\nSHA256: %s\nSHA512: %s",
		md5Result.Output, sha1Result.Output, sha256Result.Output, sha512Result.Output)

	v.output.SetContent(out)
}

func (v *HashView) View() string {
	title := styles.TitleStyle.Render("Hash Generator")

	uppercaseStr := "\u2610 Uppercase"
	if v.uppercase {
		uppercaseStr = "\u2611 Uppercase"
	}
	options := styles.LabelStyle.Render(uppercaseStr)

	inputLabel := styles.LabelStyle.Render("Input:")
	inputView := v.input.View()

	outputSection := styles.LabelStyle.Render("Output:") + "\n" + v.output.View()

	status := styles.StatusBarStyle.Render("ctrl+u: uppercase  tab: switch panel")

	return fmt.Sprintf("%s\n\n%s\n\n%s\n%s\n\n%s\n\n%s",
		title, options, inputLabel, inputView, outputSection, status)
}

func (v *HashView) SetSize(width, height int) {
	v.width = width
	v.height = height

	inputHeight := max((height-10)/2, 3)
	outputHeight := max((height-10)/2, 3)

	v.input.SetWidth(width - 4)
	v.input.SetHeight(inputHeight)
	v.output.SetWidth(width - 4)
	v.output.SetHeight(outputHeight)
}
