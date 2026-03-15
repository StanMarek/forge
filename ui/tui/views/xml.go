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

type xmlMode int

const (
	xmlFormat xmlMode = iota
	xmlMinify
)

// XMLView is the TUI view for XML formatting/minifying.
type XMLView struct {
	input  textarea.Model
	output viewport.Model
	mode   xmlMode
	width  int
	height int
	err    string
}

// NewXMLView creates a new XML Formatter tool view.
func NewXMLView() *XMLView {
	ti := textarea.New()
	ti.Placeholder = "Enter XML to format..."
	ti.Focus()

	vp := viewport.New()

	return &XMLView{
		input:  ti,
		output: vp,
		mode:   xmlFormat,
	}
}

func (v *XMLView) Init() tea.Cmd {
	return textarea.Blink
}

func (v *XMLView) Update(msg tea.Msg) (ToolView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+f"))):
			v.mode = xmlFormat
			v.input.Placeholder = "Enter XML to format..."
			v.process()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+m"))):
			v.mode = xmlMinify
			v.input.Placeholder = "Enter XML to minify..."
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

func (v *XMLView) process() {
	input := v.input.Value()
	if input == "" {
		v.output.SetContent("")
		v.err = ""
		return
	}

	var result tools.Result
	switch v.mode {
	case xmlFormat:
		result = tools.XMLFormat(input)
	case xmlMinify:
		result = tools.XMLMinify(input)
	}

	if result.Error != "" {
		v.err = result.Error
		v.output.SetContent("")
	} else {
		v.err = ""
		v.output.SetContent(result.Output)
	}
}

func (v *XMLView) View() string {
	title := styles.TitleStyle.Render("XML Formatter")

	var modeStr string
	if v.mode == xmlFormat {
		modeStr = styles.ModeActivePill.Render("Format") + " " + styles.ModeInactivePill.Render("Minify")
	} else {
		modeStr = styles.ModeInactivePill.Render("Format") + " " + styles.ModeActivePill.Render("Minify")
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

func (v *XMLView) KeyHints() string {
	return hint("ctrl+f", "format") + "  " + hint("ctrl+m", "minify")
}

func (v *XMLView) SetSize(width, height int) {
	v.width = width
	v.height = height

	inputHeight := max((height-10)/2, 3)
	outputHeight := max((height-10)/2, 3)

	v.input.SetWidth(width - 4)
	v.input.SetHeight(inputHeight)
	v.output.SetWidth(width - 4)
	v.output.SetHeight(outputHeight)
}
