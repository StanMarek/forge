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

type csvMode int

const (
	csvToCSV csvMode = iota
	csvToJSON
)

// CSVView is the TUI view for JSON/CSV conversion.
type CSVView struct {
	input  textarea.Model
	output viewport.Model
	mode   csvMode
	width  int
	height int
	err    string
}

// NewCSVView creates a new JSON to CSV tool view.
func NewCSVView() *CSVView {
	ti := textarea.New()
	ti.Placeholder = "Enter JSON array to convert to CSV..."
	ti.Focus()

	vp := viewport.New()

	return &CSVView{
		input:  ti,
		output: vp,
		mode:   csvToCSV,
	}
}

func (v *CSVView) Init() tea.Cmd {
	return textarea.Blink
}

func (v *CSVView) Update(msg tea.Msg) (ToolView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))):
			v.mode = csvToCSV
			v.input.Placeholder = "Enter JSON array to convert to CSV..."
			v.process()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+j"))):
			v.mode = csvToJSON
			v.input.Placeholder = "Enter CSV to convert to JSON..."
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

func (v *CSVView) process() {
	input := v.input.Value()
	if input == "" {
		v.output.SetContent("")
		v.err = ""
		return
	}

	var result tools.Result
	switch v.mode {
	case csvToCSV:
		result = tools.JSONToCSV(input, ",")
	case csvToJSON:
		result = tools.CSVToJSON(input, ",")
	}

	if result.Error != "" {
		v.err = result.Error
		v.output.SetContent("")
	} else {
		v.err = ""
		v.output.SetContent(result.Output)
	}
}

func (v *CSVView) View() string {
	title := styles.TitleStyle.Render("JSON to CSV")

	var modeStr string
	if v.mode == csvToCSV {
		modeStr = styles.ModeActivePill.Render("To CSV") + " " + styles.ModeInactivePill.Render("To JSON")
	} else {
		modeStr = styles.ModeInactivePill.Render("To CSV") + " " + styles.ModeActivePill.Render("To JSON")
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

func (v *CSVView) KeyHints() string {
	return hint("ctrl+c", "to-csv") + "  " + hint("ctrl+j", "to-json")
}

func (v *CSVView) SetSize(width, height int) {
	v.width = width
	v.height = height

	inputHeight := max((height-10)/2, 3)
	outputHeight := max((height-10)/2, 3)

	v.input.SetWidth(width - 4)
	v.input.SetHeight(inputHeight)
	v.output.SetWidth(width - 4)
	v.output.SetHeight(outputHeight)
}
