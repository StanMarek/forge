package views

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	"github.com/StanMarek/forge/core/tools"
	"github.com/StanMarek/forge/ui/tui/styles"
)

// ColorView is the TUI view for color conversion.
type ColorView struct {
	input  textarea.Model
	output viewport.Model
	width  int
	height int
	err    string
}

// NewColorView creates a new Color Converter tool view.
func NewColorView() *ColorView {
	ti := textarea.New()
	ti.Placeholder = "Enter a color (#hex, rgb(), or hsl())..."
	ti.Focus()

	vp := viewport.New()

	return &ColorView{
		input:  ti,
		output: vp,
	}
}

func (v *ColorView) Init() tea.Cmd {
	return textarea.Blink
}

func (v *ColorView) Update(msg tea.Msg) (ToolView, tea.Cmd) {
	var cmds []tea.Cmd

	var cmd tea.Cmd
	v.input, cmd = v.input.Update(msg)
	cmds = append(cmds, cmd)

	v.process()

	return v, tea.Batch(cmds...)
}

func (v *ColorView) process() {
	input := v.input.Value()
	if input == "" {
		v.output.SetContent("")
		v.err = ""
		return
	}

	result := tools.ColorConvert(input)

	if result.Error != "" {
		v.err = result.Error
		v.output.SetContent("")
	} else {
		v.err = ""
		v.output.SetContent(result.Output)
	}
}

func (v *ColorView) View() string {
	title := styles.TitleStyle.Render("Color Converter")

	inputLabel := styles.LabelStyle.Render("Input:")
	inputView := v.input.View()

	var outputSection string
	if v.err != "" {
		outputSection = styles.LabelStyle.Render("Error:") + "\n" + styles.ErrorTextStyle.Render(v.err)
	} else {
		outputSection = styles.LabelStyle.Render("Output:") + "\n" + v.output.View()
	}

	return fmt.Sprintf("%s\n\n%s\n%s\n\n%s",
		title, inputLabel, inputView, outputSection)
}

func (v *ColorView) KeyHints() string {
	return ""
}

func (v *ColorView) SetSize(width, height int) {
	v.width = width
	v.height = height

	inputHeight := max((height-10)/2, 3)
	outputHeight := max((height-10)/2, 3)

	v.input.SetWidth(width - 4)
	v.input.SetHeight(inputHeight)
	v.output.SetWidth(width - 4)
	v.output.SetHeight(outputHeight)
}
