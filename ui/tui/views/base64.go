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

type base64Mode int

const (
	modeEncode base64Mode = iota
	modeDecode
)

// Base64View is the TUI view for Base64 encoding/decoding.
type Base64View struct {
	input     textarea.Model
	output    viewport.Model
	mode      base64Mode
	urlSafe   bool
	noPadding bool
	width     int
	height    int
	err       string
}

// NewBase64View creates a new Base64 tool view.
func NewBase64View() *Base64View {
	ti := textarea.New()
	ti.Placeholder = "Enter text to encode..."
	ti.Focus()

	vp := viewport.New()

	return &Base64View{
		input:  ti,
		output: vp,
		mode:   modeEncode,
	}
}

func (v *Base64View) Init() tea.Cmd {
	return textarea.Blink
}

func (v *Base64View) Update(msg tea.Msg) (ToolView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+e"))):
			v.mode = modeEncode
			v.input.Placeholder = "Enter text to encode..."
			v.process()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+d"))):
			v.mode = modeDecode
			v.input.Placeholder = "Enter Base64 to decode..."
			v.process()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+u"))):
			v.urlSafe = !v.urlSafe
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

func (v *Base64View) process() {
	input := v.input.Value()
	if input == "" {
		v.output.SetContent("")
		v.err = ""
		return
	}

	var result tools.Result
	switch v.mode {
	case modeEncode:
		result = tools.Base64Encode(input, v.urlSafe, v.noPadding)
	case modeDecode:
		result = tools.Base64Decode(input, v.urlSafe)
	}

	if result.Error != "" {
		v.err = result.Error
		v.output.SetContent("")
	} else {
		v.err = ""
		v.output.SetContent(result.Output)
	}
}

func (v *Base64View) View() string {
	title := styles.TitleStyle.Render("Base64 Encode / Decode")

	var modeStr string
	if v.mode == modeEncode {
		modeStr = styles.ModeActivePill.Render("Encode") + " " + styles.ModeInactivePill.Render("Decode")
	} else {
		modeStr = styles.ModeInactivePill.Render("Encode") + " " + styles.ModeActivePill.Render("Decode")
	}

	var urlSafeStr string
	if v.urlSafe {
		urlSafeStr = styles.CheckboxOnStyle.Render("● URL-safe")
	} else {
		urlSafeStr = styles.CheckboxOffStyle.Render("○ URL-safe")
	}
	options := fmt.Sprintf("Mode: %s    %s", modeStr, urlSafeStr)

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

func (v *Base64View) KeyHints() string {
	return hint("ctrl+e", "encode") + "  " + hint("ctrl+d", "decode") + "  " + hint("ctrl+u", "url-safe")
}

func (v *Base64View) SetSize(width, height int) {
	v.width = width
	v.height = height

	inputHeight := max((height-10)/2, 3)
	outputHeight := max((height-10)/2, 3)

	v.input.SetWidth(width - 4)
	v.input.SetHeight(inputHeight)
	v.output.SetWidth(width - 4)
	v.output.SetHeight(outputHeight)
}
