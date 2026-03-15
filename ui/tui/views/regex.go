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

type regexFocus int

const (
	regexFocusPattern regexFocus = iota
	regexFocusTest
)

// RegexView is the TUI view for regex testing.
type RegexView struct {
	pattern textarea.Model
	test    textarea.Model
	output  viewport.Model
	global  bool
	focus   regexFocus
	width   int
	height  int
	err     string
}

// NewRegexView creates a new Regex Tester tool view.
func NewRegexView() *RegexView {
	pat := textarea.New()
	pat.Placeholder = "Enter regex pattern..."
	pat.Focus()

	test := textarea.New()
	test.Placeholder = "Enter test string..."

	vp := viewport.New()

	return &RegexView{
		pattern: pat,
		test:    test,
		output:  vp,
		global:  true,
		focus:   regexFocusPattern,
	}
}

func (v *RegexView) Init() tea.Cmd {
	return textarea.Blink
}

func (v *RegexView) Update(msg tea.Msg) (ToolView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+g"))):
			v.global = !v.global
			v.process()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+p"))):
			v.focus = regexFocusPattern
			v.pattern.Focus()
			v.test.Blur()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+t"))):
			v.focus = regexFocusTest
			v.test.Focus()
			v.pattern.Blur()
			return v, nil
		}
	}

	var cmd tea.Cmd
	if v.focus == regexFocusPattern {
		v.pattern, cmd = v.pattern.Update(msg)
	} else {
		v.test, cmd = v.test.Update(msg)
	}
	cmds = append(cmds, cmd)

	v.process()

	return v, tea.Batch(cmds...)
}

func (v *RegexView) process() {
	pattern := v.pattern.Value()
	test := v.test.Value()
	if pattern == "" || test == "" {
		v.output.SetContent("")
		v.err = ""
		return
	}

	result := tools.RegexTest(pattern, test, v.global)

	if result.Error != "" {
		v.err = result.Error
		v.output.SetContent("")
	} else {
		v.err = ""
		v.output.SetContent(result.Output)
	}
}

func (v *RegexView) View() string {
	title := styles.TitleStyle.Render("Regex Tester")

	globalStr := styles.CheckboxOffStyle.Render("○ Global")
	if v.global {
		globalStr = styles.CheckboxOnStyle.Render("● Global")
	}
	options := globalStr

	patternLabel := styles.LabelStyle.Render("Pattern:")
	patternView := v.pattern.View()

	testLabel := styles.LabelStyle.Render("Test String:")
	testView := v.test.View()

	var outputSection string
	if v.err != "" {
		outputSection = styles.LabelStyle.Render("Error:") + "\n" + styles.ErrorTextStyle.Render(v.err)
	} else {
		outputSection = styles.LabelStyle.Render("Matches:") + "\n" + v.output.View()
	}

	return fmt.Sprintf("%s\n\n%s\n\n%s\n%s\n\n%s\n%s\n\n%s",
		title, options, patternLabel, patternView, testLabel, testView, outputSection)
}

func (v *RegexView) KeyHints() string {
	return hint("ctrl+g", "global") + "  " + hint("ctrl+p", "pattern") + "  " + hint("ctrl+t", "test")
}

func (v *RegexView) SetSize(width, height int) {
	v.width = width
	v.height = height

	inputHeight := max((height-14)/3, 2)
	outputHeight := max((height-14)/3, 2)

	v.pattern.SetWidth(width - 4)
	v.pattern.SetHeight(inputHeight)
	v.test.SetWidth(width - 4)
	v.test.SetHeight(inputHeight)
	v.output.SetWidth(width - 4)
	v.output.SetHeight(outputHeight)
}
