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

type yamlMode int

const (
	yamlToJSON yamlMode = iota
	yamlToYAML
)

// YAMLView is the TUI view for JSON/YAML conversion.
type YAMLView struct {
	input  textarea.Model
	output viewport.Model
	mode   yamlMode
	width  int
	height int
	err    string
}

// NewYAMLView creates a new YAML converter tool view.
func NewYAMLView() *YAMLView {
	ti := textarea.New()
	ti.Placeholder = "Enter YAML to convert to JSON..."
	ti.Focus()

	vp := viewport.New()

	return &YAMLView{
		input:  ti,
		output: vp,
		mode:   yamlToJSON,
	}
}

func (v *YAMLView) Init() tea.Cmd {
	return textarea.Blink
}

func (v *YAMLView) Update(msg tea.Msg) (ToolView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+j"))):
			v.mode = yamlToJSON
			v.input.Placeholder = "Enter YAML to convert to JSON..."
			v.process()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+y"))):
			v.mode = yamlToYAML
			v.input.Placeholder = "Enter JSON to convert to YAML..."
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

func (v *YAMLView) process() {
	input := v.input.Value()
	if input == "" {
		v.output.SetContent("")
		v.err = ""
		return
	}

	var result tools.Result
	switch v.mode {
	case yamlToJSON:
		result = tools.YAMLToJSON(input, false)
	case yamlToYAML:
		result = tools.JSONToYAML(input)
	}

	if result.Error != "" {
		v.err = result.Error
		v.output.SetContent("")
	} else {
		v.err = ""
		v.output.SetContent(result.Output)
	}
}

func (v *YAMLView) View() string {
	title := styles.TitleStyle.Render("JSON / YAML Converter")

	var modeStr string
	if v.mode == yamlToJSON {
		modeStr = styles.ModeActivePill.Render("To JSON") + " " + styles.ModeInactivePill.Render("To YAML")
	} else {
		modeStr = styles.ModeInactivePill.Render("To JSON") + " " + styles.ModeActivePill.Render("To YAML")
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

func (v *YAMLView) KeyHints() string {
	return hint("ctrl+j", "to-json") + "  " + hint("ctrl+y", "to-yaml")
}

func (v *YAMLView) SetSize(width, height int) {
	v.width = width
	v.height = height

	inputHeight := max((height-10)/2, 3)
	outputHeight := max((height-10)/2, 3)

	v.input.SetWidth(width - 4)
	v.input.SetHeight(inputHeight)
	v.output.SetWidth(width - 4)
	v.output.SetHeight(outputHeight)
}
