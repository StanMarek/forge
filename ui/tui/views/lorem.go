package views

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	"github.com/StanMarek/forge/core/tools"
	"github.com/StanMarek/forge/ui/tui/styles"
)

type loremMode int

const (
	loremWords loremMode = iota
	loremSentences
	loremParagraphs
)

// LoremView is the TUI view for lorem ipsum generation.
type LoremView struct {
	output viewport.Model
	mode   loremMode
	count  int
	width  int
	height int
	err    string
}

// NewLoremView creates a new Lorem Ipsum Generator tool view.
func NewLoremView() *LoremView {
	vp := viewport.New()

	v := &LoremView{
		output: vp,
		mode:   loremWords,
		count:  20,
	}
	v.generate()
	return v
}

func (v *LoremView) Init() tea.Cmd {
	return nil
}

func (v *LoremView) Update(msg tea.Msg) (ToolView, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+w"))):
			v.mode = loremWords
			v.count = 20
			v.generate()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+s"))):
			v.mode = loremSentences
			v.count = 5
			v.generate()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+p"))):
			v.mode = loremParagraphs
			v.count = 3
			v.generate()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+up"))):
			v.count++
			v.generate()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+down"))):
			if v.count > 1 {
				v.count--
				v.generate()
			}
			return v, nil
		}
	}

	return v, nil
}

func (v *LoremView) generate() {
	var result tools.Result
	switch v.mode {
	case loremWords:
		result = tools.LoremGenerate(v.count, 0, 0)
	case loremSentences:
		result = tools.LoremGenerate(0, v.count, 0)
	case loremParagraphs:
		result = tools.LoremGenerate(0, 0, v.count)
	}

	if result.Error != "" {
		v.err = result.Error
		v.output.SetContent("")
	} else {
		v.err = ""
		v.output.SetContent(result.Output)
	}
}

func (v *LoremView) View() string {
	title := styles.TitleStyle.Render("Lorem Ipsum Generator")

	var modeStr string
	switch v.mode {
	case loremWords:
		modeStr = styles.ModeActivePill.Render("Words") + " " + styles.ModeInactivePill.Render("Sentences") + " " + styles.ModeInactivePill.Render("Paragraphs")
	case loremSentences:
		modeStr = styles.ModeInactivePill.Render("Words") + " " + styles.ModeActivePill.Render("Sentences") + " " + styles.ModeInactivePill.Render("Paragraphs")
	case loremParagraphs:
		modeStr = styles.ModeInactivePill.Render("Words") + " " + styles.ModeInactivePill.Render("Sentences") + " " + styles.ModeActivePill.Render("Paragraphs")
	}

	countStr := styles.LabelStyle.Render(fmt.Sprintf("Count: %d", v.count))
	options := fmt.Sprintf("Mode: %s    %s", modeStr, countStr)

	var outputSection string
	if v.err != "" {
		outputSection = styles.LabelStyle.Render("Error:") + "\n" + styles.ErrorTextStyle.Render(v.err)
	} else {
		outputSection = styles.LabelStyle.Render("Output:") + "\n" + v.output.View()
	}

	return fmt.Sprintf("%s\n\n%s\n\n%s", title, options, outputSection)
}

func (v *LoremView) KeyHints() string {
	return hint("ctrl+w", "words") + "  " +
		hint("ctrl+s", "sentences") + "  " +
		hint("ctrl+p", "paragraphs") + "  " +
		hint("ctrl+↑↓", "count")
}

func (v *LoremView) SetSize(width, height int) {
	v.width = width
	v.height = height

	outputHeight := max(height-8, 3)

	v.output.SetWidth(width - 4)
	v.output.SetHeight(outputHeight)
}
