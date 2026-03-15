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

type jsonMode int

const (
	jsonModeFormat jsonMode = iota
	jsonModeMinify
	jsonModeValidate
)

// JSONView is the TUI view for JSON formatting, minifying, and validating.
type JSONView struct {
	input    textarea.Model
	output   viewport.Model
	mode     jsonMode
	sortKeys bool
	width    int
	height   int
	err      string
}

// NewJSONView creates a new JSON tool view.
func NewJSONView() *JSONView {
	ti := textarea.New()
	ti.Placeholder = "Paste JSON..."
	ti.Focus()

	vp := viewport.New()

	return &JSONView{
		input:  ti,
		output: vp,
		mode:   jsonModeFormat,
	}
}

func (v *JSONView) Init() tea.Cmd {
	return textarea.Blink
}

func (v *JSONView) Update(msg tea.Msg) (ToolView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+f"))):
			v.mode = jsonModeFormat
			v.process()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+m"))):
			v.mode = jsonModeMinify
			v.process()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+v"))):
			v.mode = jsonModeValidate
			v.process()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+s"))):
			v.sortKeys = !v.sortKeys
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

func (v *JSONView) process() {
	input := v.input.Value()
	if input == "" {
		v.output.SetContent("")
		v.err = ""
		return
	}

	var result tools.Result
	switch v.mode {
	case jsonModeFormat:
		result = tools.JSONFormat(input, 2, v.sortKeys, false)
	case jsonModeMinify:
		result = tools.JSONMinify(input)
	case jsonModeValidate:
		result = tools.JSONValidate(input)
	}

	if result.Error != "" {
		v.err = result.Error
		v.output.SetContent("")
	} else {
		v.err = ""
		v.output.SetContent(result.Output)
	}
}

func (v *JSONView) View() string {
	title := styles.TitleStyle.Render("JSON Formatter")

	var modeStr string
	switch v.mode {
	case jsonModeFormat:
		modeStr = styles.ModeActiveStyle.Render("(●) Format") + "  " +
			styles.ModeInactiveStyle.Render("( ) Minify") + "  " +
			styles.ModeInactiveStyle.Render("( ) Validate")
	case jsonModeMinify:
		modeStr = styles.ModeInactiveStyle.Render("( ) Format") + "  " +
			styles.ModeActiveStyle.Render("(●) Minify") + "  " +
			styles.ModeInactiveStyle.Render("( ) Validate")
	case jsonModeValidate:
		modeStr = styles.ModeInactiveStyle.Render("( ) Format") + "  " +
			styles.ModeInactiveStyle.Render("( ) Minify") + "  " +
			styles.ModeActiveStyle.Render("(●) Validate")
	}

	sortKeysStr := "☐ Sort keys"
	if v.sortKeys {
		sortKeysStr = "☑ Sort keys"
	}
	options := fmt.Sprintf("Mode: %s    %s", modeStr, styles.LabelStyle.Render(sortKeysStr))

	inputLabel := styles.LabelStyle.Render("Input:")
	inputView := v.input.View()

	var outputSection string
	if v.err != "" {
		outputSection = styles.LabelStyle.Render("Error:") + "\n" + styles.ErrorStyle.Render(v.err)
	} else {
		outputSection = styles.LabelStyle.Render("Output:") + "\n" + v.output.View()
	}

	status := styles.StatusBarStyle.Render("ctrl+f: format  ctrl+m: minify  ctrl+v: validate  ctrl+s: sort keys  tab: switch panel")

	return fmt.Sprintf("%s\n\n%s\n\n%s\n%s\n\n%s\n\n%s",
		title, options, inputLabel, inputView, outputSection, status)
}

func (v *JSONView) SetSize(width, height int) {
	v.width = width
	v.height = height

	inputHeight := max((height - 10) / 2, 3)
	outputHeight := max((height - 10) / 2, 3)

	v.input.SetWidth(width - 4)
	v.input.SetHeight(inputHeight)
	v.output.SetWidth(width - 4)
	v.output.SetHeight(outputHeight)
}
