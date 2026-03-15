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

type textEscapeMode int

const (
	textEscapeModeEscape textEscapeMode = iota
	textEscapeModeUnescape
)

// TextEscapeView is the TUI view for text escaping/unescaping.
type TextEscapeView struct {
	input  textarea.Model
	output viewport.Model
	mode   textEscapeMode
	width  int
	height int
	err    string
}

// NewTextEscapeView creates a new Text Escape / Unescape tool view.
func NewTextEscapeView() *TextEscapeView {
	ti := textarea.New()
	ti.Placeholder = "Enter text to escape..."
	ti.Focus()

	vp := viewport.New()

	return &TextEscapeView{
		input:  ti,
		output: vp,
		mode:   textEscapeModeEscape,
	}
}

func (v *TextEscapeView) Init() tea.Cmd {
	return textarea.Blink
}

func (v *TextEscapeView) Update(msg tea.Msg) (ToolView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+e"))):
			v.mode = textEscapeModeEscape
			v.input.Placeholder = "Enter text to escape..."
			v.process()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+u"))):
			v.mode = textEscapeModeUnescape
			v.input.Placeholder = "Enter text to unescape..."
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

func (v *TextEscapeView) process() {
	input := v.input.Value()
	if input == "" {
		v.output.SetContent("")
		v.err = ""
		return
	}

	var result tools.Result
	switch v.mode {
	case textEscapeModeEscape:
		result = tools.TextEscape(input)
	case textEscapeModeUnescape:
		result = tools.TextUnescape(input)
	}

	if result.Error != "" {
		v.err = result.Error
		v.output.SetContent("")
	} else {
		v.err = ""
		v.output.SetContent(result.Output)
	}
}

func (v *TextEscapeView) View() string {
	title := styles.TitleStyle.Render("Text Escape / Unescape")

	var modeStr string
	if v.mode == textEscapeModeEscape {
		modeStr = styles.ModeActivePill.Render("Escape") + " " + styles.ModeInactivePill.Render("Unescape")
	} else {
		modeStr = styles.ModeInactivePill.Render("Escape") + " " + styles.ModeActivePill.Render("Unescape")
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

func (v *TextEscapeView) KeyHints() string {
	return hint("ctrl+e", "escape") + "  " + hint("ctrl+u", "unescape")
}

func (v *TextEscapeView) SetSize(width, height int) {
	v.width = width
	v.height = height

	inputHeight := max((height-10)/2, 3)
	outputHeight := max((height-10)/2, 3)

	v.input.SetWidth(width - 4)
	v.input.SetHeight(inputHeight)
	v.output.SetWidth(width - 4)
	v.output.SetHeight(outputHeight)
}
