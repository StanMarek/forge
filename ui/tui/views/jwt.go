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

type jwtMode int

const (
	jwtModeFull    jwtMode = iota
	jwtModeHeader
	jwtModePayload
)

// JWTView is the TUI view for JWT decoding.
type JWTView struct {
	input  textarea.Model
	output viewport.Model
	mode   jwtMode
	width  int
	height int
	err    string
}

// NewJWTView creates a new JWT tool view.
func NewJWTView() *JWTView {
	ti := textarea.New()
	ti.Placeholder = "Paste JWT token..."
	ti.Focus()

	vp := viewport.New()

	return &JWTView{
		input:  ti,
		output: vp,
		mode:   jwtModeFull,
	}
}

func (v *JWTView) Init() tea.Cmd {
	return textarea.Blink
}

func (v *JWTView) Update(msg tea.Msg) (ToolView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+f"))):
			v.mode = jwtModeFull
			v.process()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+h"))):
			v.mode = jwtModeHeader
			v.process()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+p"))):
			v.mode = jwtModePayload
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

func (v *JWTView) process() {
	input := v.input.Value()
	if input == "" {
		v.output.SetContent("")
		v.err = ""
		return
	}

	result := tools.JWTDecode(input)

	if result.Error != "" {
		v.err = result.Error
		v.output.SetContent("")
	} else {
		v.err = ""
		switch v.mode {
		case jwtModeFull:
			v.output.SetContent(result.Output)
		case jwtModeHeader:
			v.output.SetContent(result.Header)
		case jwtModePayload:
			v.output.SetContent(result.Payload)
		}
	}
}

func (v *JWTView) View() string {
	title := styles.TitleStyle.Render("JWT Decoder")

	var modeStr string
	switch v.mode {
	case jwtModeFull:
		modeStr = styles.ModeActivePill.Render("Full") + " " +
			styles.ModeInactivePill.Render("Header") + " " +
			styles.ModeInactivePill.Render("Payload")
	case jwtModeHeader:
		modeStr = styles.ModeInactivePill.Render("Full") + " " +
			styles.ModeActivePill.Render("Header") + " " +
			styles.ModeInactivePill.Render("Payload")
	case jwtModePayload:
		modeStr = styles.ModeInactivePill.Render("Full") + " " +
			styles.ModeInactivePill.Render("Header") + " " +
			styles.ModeActivePill.Render("Payload")
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

func (v *JWTView) KeyHints() string {
	return hint("ctrl+f", "full") + "  " + hint("ctrl+h", "header") + "  " + hint("ctrl+p", "payload")
}

func (v *JWTView) SetSize(width, height int) {
	v.width = width
	v.height = height

	inputHeight := max((height-10)/2, 3)
	outputHeight := max((height-10)/2, 3)

	v.input.SetWidth(width - 4)
	v.input.SetHeight(inputHeight)
	v.output.SetWidth(width - 4)
	v.output.SetHeight(outputHeight)
}
