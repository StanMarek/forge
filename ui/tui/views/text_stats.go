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

type textStatsMode int

const (
	tsStats textStatsMode = iota
	tsLower
	tsUpper
	tsTitle
	tsCamel
	tsSnake
	tsKebab
)

var textStatsModeNames = []string{
	"Stats", "Lower", "Upper", "Title", "Camel", "Snake", "Kebab",
}

// TextStatsView is the TUI view for text analysis and case conversion.
type TextStatsView struct {
	input  textarea.Model
	output viewport.Model
	mode   textStatsMode
	width  int
	height int
	err    string
}

// NewTextStatsView creates a new Text Analyzer tool view.
func NewTextStatsView() *TextStatsView {
	ti := textarea.New()
	ti.Placeholder = "Enter text to analyze..."
	ti.Focus()

	vp := viewport.New()

	return &TextStatsView{
		input:  ti,
		output: vp,
		mode:   tsStats,
	}
}

func (v *TextStatsView) Init() tea.Cmd {
	return textarea.Blink
}

func (v *TextStatsView) Update(msg tea.Msg) (ToolView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+m"))):
			v.mode = (v.mode + 1) % textStatsMode(len(textStatsModeNames))
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

func (v *TextStatsView) process() {
	input := v.input.Value()
	if input == "" {
		v.output.SetContent("")
		v.err = ""
		return
	}

	var result tools.Result
	switch v.mode {
	case tsStats:
		result = tools.TextStats(input)
	case tsLower:
		result = tools.TextCaseConvert(input, "lower")
	case tsUpper:
		result = tools.TextCaseConvert(input, "upper")
	case tsTitle:
		result = tools.TextCaseConvert(input, "title")
	case tsCamel:
		result = tools.TextCaseConvert(input, "camel")
	case tsSnake:
		result = tools.TextCaseConvert(input, "snake")
	case tsKebab:
		result = tools.TextCaseConvert(input, "kebab")
	}

	if result.Error != "" {
		v.err = result.Error
		v.output.SetContent("")
	} else {
		v.err = ""
		v.output.SetContent(result.Output)
	}
}

func (v *TextStatsView) View() string {
	title := styles.TitleStyle.Render("Text Analyzer")

	var modeStr string
	for i, name := range textStatsModeNames {
		if i > 0 {
			modeStr += " "
		}
		if textStatsMode(i) == v.mode {
			modeStr += styles.ModeActivePill.Render(name)
		} else {
			modeStr += styles.ModeInactivePill.Render(name)
		}
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

func (v *TextStatsView) KeyHints() string {
	return hint("ctrl+m", "cycle mode")
}

func (v *TextStatsView) SetSize(width, height int) {
	v.width = width
	v.height = height

	inputHeight := max((height-10)/2, 3)
	outputHeight := max((height-10)/2, 3)

	v.input.SetWidth(width - 4)
	v.input.SetHeight(inputHeight)
	v.output.SetWidth(width - 4)
	v.output.SetHeight(outputHeight)
}
