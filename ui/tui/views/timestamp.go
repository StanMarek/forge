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

type timestampMode int

const (
	tsFromUnix timestampMode = iota
	tsToUnix
	tsNow
)

// TimestampView is the TUI view for timestamp conversion.
type TimestampView struct {
	input  textarea.Model
	output viewport.Model
	mode   timestampMode
	width  int
	height int
	err    string
}

// NewTimestampView creates a new Timestamp converter tool view.
func NewTimestampView() *TimestampView {
	ti := textarea.New()
	ti.Placeholder = "Enter unix timestamp..."
	ti.Focus()

	vp := viewport.New()

	return &TimestampView{
		input:  ti,
		output: vp,
		mode:   tsFromUnix,
	}
}

func (v *TimestampView) Init() tea.Cmd {
	return textarea.Blink
}

func (v *TimestampView) Update(msg tea.Msg) (ToolView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+f"))):
			v.mode = tsFromUnix
			v.input.Placeholder = "Enter unix timestamp..."
			v.process()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+t"))):
			v.mode = tsToUnix
			v.input.Placeholder = "Enter datetime (RFC3339, etc.)..."
			v.process()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+n"))):
			v.mode = tsNow
			v.process()
			return v, nil
		}
	}

	if v.mode != tsNow {
		var cmd tea.Cmd
		v.input, cmd = v.input.Update(msg)
		cmds = append(cmds, cmd)
	}

	v.process()

	return v, tea.Batch(cmds...)
}

func (v *TimestampView) process() {
	if v.mode == tsNow {
		result := tools.TimestampNow("")
		if result.Error != "" {
			v.err = result.Error
			v.output.SetContent("")
		} else {
			v.err = ""
			v.output.SetContent(result.Output)
		}
		return
	}

	input := v.input.Value()
	if input == "" {
		v.output.SetContent("")
		v.err = ""
		return
	}

	var result tools.Result
	switch v.mode {
	case tsFromUnix:
		result = tools.TimestampFromUnix(input, "")
	case tsToUnix:
		result = tools.TimestampToUnix(input, false)
	}

	if result.Error != "" {
		v.err = result.Error
		v.output.SetContent("")
	} else {
		v.err = ""
		v.output.SetContent(result.Output)
	}
}

func (v *TimestampView) View() string {
	title := styles.TitleStyle.Render("Timestamp Converter")

	var modeStr string
	switch v.mode {
	case tsFromUnix:
		modeStr = styles.ModeActivePill.Render("From Unix") + " " + styles.ModeInactivePill.Render("To Unix") + " " + styles.ModeInactivePill.Render("Now")
	case tsToUnix:
		modeStr = styles.ModeInactivePill.Render("From Unix") + " " + styles.ModeActivePill.Render("To Unix") + " " + styles.ModeInactivePill.Render("Now")
	case tsNow:
		modeStr = styles.ModeInactivePill.Render("From Unix") + " " + styles.ModeInactivePill.Render("To Unix") + " " + styles.ModeActivePill.Render("Now")
	}
	options := fmt.Sprintf("Mode: %s", modeStr)

	var body string
	if v.mode == tsNow {
		var outputSection string
		if v.err != "" {
			outputSection = styles.LabelStyle.Render("Error:") + "\n" + styles.ErrorTextStyle.Render(v.err)
		} else {
			outputSection = styles.LabelStyle.Render("Output:") + "\n" + v.output.View()
		}
		body = fmt.Sprintf("%s\n\n%s\n\n%s", title, options, outputSection)
	} else {
		inputLabel := styles.LabelStyle.Render("Input:")
		inputView := v.input.View()

		var outputSection string
		if v.err != "" {
			outputSection = styles.LabelStyle.Render("Error:") + "\n" + styles.ErrorTextStyle.Render(v.err)
		} else {
			outputSection = styles.LabelStyle.Render("Output:") + "\n" + v.output.View()
		}
		body = fmt.Sprintf("%s\n\n%s\n\n%s\n%s\n\n%s",
			title, options, inputLabel, inputView, outputSection)
	}

	return body
}

func (v *TimestampView) KeyHints() string {
	return hint("ctrl+f", "from-unix") + "  " + hint("ctrl+t", "to-unix") + "  " + hint("ctrl+n", "now")
}

func (v *TimestampView) SetSize(width, height int) {
	v.width = width
	v.height = height

	inputHeight := max((height-10)/2, 3)
	outputHeight := max((height-10)/2, 3)

	v.input.SetWidth(width - 4)
	v.input.SetHeight(inputHeight)
	v.output.SetWidth(width - 4)
	v.output.SetHeight(outputHeight)
}
