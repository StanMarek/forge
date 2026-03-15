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

type diffFocus int

const (
	diffFocusA diffFocus = iota
	diffFocusB
)

// DiffView is the TUI view for text comparison.
type DiffView struct {
	textA  textarea.Model
	textB  textarea.Model
	output viewport.Model
	focus  diffFocus
	width  int
	height int
	err    string
}

// NewDiffView creates a new Text Diff tool view.
func NewDiffView() *DiffView {
	ta := textarea.New()
	ta.Placeholder = "Enter text A..."
	ta.Focus()

	tb := textarea.New()
	tb.Placeholder = "Enter text B..."

	vp := viewport.New()

	return &DiffView{
		textA:  ta,
		textB:  tb,
		output: vp,
		focus:  diffFocusA,
	}
}

func (v *DiffView) Init() tea.Cmd {
	return textarea.Blink
}

func (v *DiffView) Update(msg tea.Msg) (ToolView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+a"))):
			v.focus = diffFocusA
			v.textA.Focus()
			v.textB.Blur()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+b"))):
			v.focus = diffFocusB
			v.textB.Focus()
			v.textA.Blur()
			return v, nil
		}
	}

	var cmd tea.Cmd
	if v.focus == diffFocusA {
		v.textA, cmd = v.textA.Update(msg)
	} else {
		v.textB, cmd = v.textB.Update(msg)
	}
	cmds = append(cmds, cmd)

	v.process()

	return v, tea.Batch(cmds...)
}

func (v *DiffView) process() {
	a := v.textA.Value()
	b := v.textB.Value()
	if a == "" && b == "" {
		v.output.SetContent("")
		v.err = ""
		return
	}

	result := tools.DiffText(a, b)

	if result.Error != "" {
		v.err = result.Error
		v.output.SetContent("")
	} else {
		v.err = ""
		v.output.SetContent(result.Output)
	}
}

func (v *DiffView) View() string {
	title := styles.TitleStyle.Render("Text Diff")

	textALabel := styles.LabelStyle.Render("Text A:")
	textAView := v.textA.View()

	textBLabel := styles.LabelStyle.Render("Text B:")
	textBView := v.textB.View()

	var outputSection string
	if v.err != "" {
		outputSection = styles.LabelStyle.Render("Error:") + "\n" + styles.ErrorTextStyle.Render(v.err)
	} else {
		outputSection = styles.LabelStyle.Render("Diff:") + "\n" + v.output.View()
	}

	return fmt.Sprintf("%s\n\n%s\n%s\n\n%s\n%s\n\n%s",
		title, textALabel, textAView, textBLabel, textBView, outputSection)
}

func (v *DiffView) KeyHints() string {
	return hint("ctrl+a", "text A") + "  " + hint("ctrl+b", "text B")
}

func (v *DiffView) SetSize(width, height int) {
	v.width = width
	v.height = height

	inputHeight := max((height-14)/3, 2)
	outputHeight := max((height-14)/3, 2)

	v.textA.SetWidth(width - 4)
	v.textA.SetHeight(inputHeight)
	v.textB.SetWidth(width - 4)
	v.textB.SetHeight(inputHeight)
	v.output.SetWidth(width - 4)
	v.output.SetHeight(outputHeight)
}
