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

type htmlEntityMode int

const (
	htmlEntityEncode htmlEntityMode = iota
	htmlEntityDecode
)

// HTMLEntityView is the TUI view for HTML entity encoding/decoding.
type HTMLEntityView struct {
	input  textarea.Model
	output viewport.Model
	mode   htmlEntityMode
	width  int
	height int
	err    string
}

// NewHTMLEntityView creates a new HTML Entity tool view.
func NewHTMLEntityView() *HTMLEntityView {
	ti := textarea.New()
	ti.Placeholder = "Enter text to encode HTML entities..."
	ti.Focus()

	vp := viewport.New()

	return &HTMLEntityView{
		input:  ti,
		output: vp,
		mode:   htmlEntityEncode,
	}
}

func (v *HTMLEntityView) Init() tea.Cmd {
	return textarea.Blink
}

func (v *HTMLEntityView) Update(msg tea.Msg) (ToolView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+e"))):
			v.mode = htmlEntityEncode
			v.input.Placeholder = "Enter text to encode HTML entities..."
			v.process()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+d"))):
			v.mode = htmlEntityDecode
			v.input.Placeholder = "Enter HTML entities to decode..."
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

func (v *HTMLEntityView) process() {
	input := v.input.Value()
	if input == "" {
		v.output.SetContent("")
		v.err = ""
		return
	}

	var result tools.Result
	switch v.mode {
	case htmlEntityEncode:
		result = tools.HTMLEntityEncode(input)
	case htmlEntityDecode:
		result = tools.HTMLEntityDecode(input)
	}

	if result.Error != "" {
		v.err = result.Error
		v.output.SetContent("")
	} else {
		v.err = ""
		v.output.SetContent(result.Output)
	}
}

func (v *HTMLEntityView) View() string {
	title := styles.TitleStyle.Render("HTML Entity Encoder")

	var modeStr string
	if v.mode == htmlEntityEncode {
		modeStr = styles.ModeActivePill.Render("Encode") + " " + styles.ModeInactivePill.Render("Decode")
	} else {
		modeStr = styles.ModeInactivePill.Render("Encode") + " " + styles.ModeActivePill.Render("Decode")
	}
	options := fmt.Sprintf("Mode: %s", modeStr)

	inputLabel := styles.LabelStyle.Render("Input:")
	inputView := v.input.View()

	var outputSection string
	if v.err != "" {
		outputSection = styles.LabelStyle.Render("Error:") + "\n" + styles.ErrorTextStyle.Render(v.err)
	} else {
		outputSection = styles.LabelStyle.Render("Output:") + "\n" + v.output.View()
	}

	return fmt.Sprintf("%s\n\n%s\n\n%s\n%s\n\n%s",
		title, options, inputLabel, inputView, outputSection)
}

func (v *HTMLEntityView) KeyHints() string {
	return hint("ctrl+e", "encode") + "  " + hint("ctrl+d", "decode")
}

func (v *HTMLEntityView) SetSize(width, height int) {
	v.width = width
	v.height = height

	inputHeight := max((height-10)/2, 3)
	outputHeight := max((height-10)/2, 3)

	v.input.SetWidth(width - 4)
	v.input.SetHeight(inputHeight)
	v.output.SetWidth(width - 4)
	v.output.SetHeight(outputHeight)
}
