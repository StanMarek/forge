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

type urlMode int

const (
	urlModeParse  urlMode = iota
	urlModeEncode
	urlModeDecode
)

// URLView is the TUI view for URL encoding/decoding/parsing.
type URLView struct {
	input     textarea.Model
	output    viewport.Model
	mode      urlMode
	component bool
	width     int
	height    int
	err       string
}

// NewURLView creates a new URL tool view.
func NewURLView() *URLView {
	ti := textarea.New()
	ti.Placeholder = "Enter URL..."
	ti.Focus()

	vp := viewport.New()

	return &URLView{
		input:  ti,
		output: vp,
		mode:   urlModeParse,
	}
}

func (v *URLView) Init() tea.Cmd {
	return textarea.Blink
}

func (v *URLView) Update(msg tea.Msg) (ToolView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+p"))):
			v.mode = urlModeParse
			v.input.Placeholder = "Enter URL..."
			v.process()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+e"))):
			v.mode = urlModeEncode
			v.input.Placeholder = "Enter text to URL-encode..."
			v.process()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+d"))):
			v.mode = urlModeDecode
			v.input.Placeholder = "Enter URL-encoded text to decode..."
			v.process()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+o"))):
			v.component = !v.component
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

func (v *URLView) process() {
	input := v.input.Value()
	if input == "" {
		v.output.SetContent("")
		v.err = ""
		return
	}

	switch v.mode {
	case urlModeParse:
		result := tools.URLParse(input)
		if result.Error != "" {
			v.err = result.Error
			v.output.SetContent("")
		} else {
			v.err = ""
			v.output.SetContent(result.Output)
		}
	case urlModeEncode:
		result := tools.URLEncode(input, v.component)
		if result.Error != "" {
			v.err = result.Error
			v.output.SetContent("")
		} else {
			v.err = ""
			v.output.SetContent(result.Output)
		}
	case urlModeDecode:
		result := tools.URLDecode(input)
		if result.Error != "" {
			v.err = result.Error
			v.output.SetContent("")
		} else {
			v.err = ""
			v.output.SetContent(result.Output)
		}
	}
}

func (v *URLView) View() string {
	title := styles.TitleStyle.Render("URL Encode / Decode / Parse")

	var modeStr string
	switch v.mode {
	case urlModeParse:
		modeStr = styles.ModeActiveStyle.Render("(*) Parse") + "  " +
			styles.ModeInactiveStyle.Render("( ) Encode") + "  " +
			styles.ModeInactiveStyle.Render("( ) Decode")
	case urlModeEncode:
		modeStr = styles.ModeInactiveStyle.Render("( ) Parse") + "  " +
			styles.ModeActiveStyle.Render("(*) Encode") + "  " +
			styles.ModeInactiveStyle.Render("( ) Decode")
	case urlModeDecode:
		modeStr = styles.ModeInactiveStyle.Render("( ) Parse") + "  " +
			styles.ModeInactiveStyle.Render("( ) Encode") + "  " +
			styles.ModeActiveStyle.Render("(*) Decode")
	}

	options := fmt.Sprintf("Mode: %s", modeStr)
	if v.mode == urlModeEncode {
		componentStr := "[ ] Component"
		if v.component {
			componentStr = "[x] Component"
		}
		options += "    " + styles.LabelStyle.Render(componentStr)
	}

	inputLabel := styles.LabelStyle.Render("Input:")
	inputView := v.input.View()

	var outputSection string
	if v.err != "" {
		outputSection = styles.LabelStyle.Render("Error:") + "\n" + styles.ErrorStyle.Render(v.err)
	} else {
		outputSection = styles.LabelStyle.Render("Output:") + "\n" + v.output.View()
	}

	status := styles.StatusBarStyle.Render("ctrl+p: parse  ctrl+e: encode  ctrl+d: decode  ctrl+o: component  tab: switch panel")

	return fmt.Sprintf("%s\n\n%s\n\n%s\n%s\n\n%s\n\n%s",
		title, options, inputLabel, inputView, outputSection, status)
}

func (v *URLView) SetSize(width, height int) {
	v.width = width
	v.height = height

	inputHeight := max((height-10)/2, 3)
	outputHeight := max((height-10)/2, 3)

	v.input.SetWidth(width - 4)
	v.input.SetHeight(inputHeight)
	v.output.SetWidth(width - 4)
	v.output.SetHeight(outputHeight)
}
