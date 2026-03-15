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

type uuidMode int

const (
	uuidModeGenerate uuidMode = iota
	uuidModeValidate
	uuidModeParse
)

// UUIDView is the TUI view for UUID generation, validation, and parsing.
type UUIDView struct {
	input     textarea.Model
	output    viewport.Model
	mode      uuidMode
	version   int
	uppercase bool
	noHyphens bool
	width     int
	height    int
	err       string
}

// NewUUIDView creates a new UUID tool view.
func NewUUIDView() *UUIDView {
	ti := textarea.New()
	ti.Placeholder = "Enter UUID..."

	vp := viewport.New()

	v := &UUIDView{
		input:   ti,
		output:  vp,
		mode:    uuidModeGenerate,
		version: 4,
	}

	// Auto-generate one UUID on creation.
	v.process()

	return v
}

func (v *UUIDView) Init() tea.Cmd {
	return nil
}

func (v *UUIDView) Update(msg tea.Msg) (ToolView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+g"))):
			v.mode = uuidModeGenerate
			v.input.Blur()
			v.process()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+4"))):
			if v.mode == uuidModeGenerate {
				v.version = 4
				v.process()
			}
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+7"))):
			if v.mode == uuidModeGenerate {
				v.version = 7
				v.process()
			}
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+u"))):
			v.uppercase = !v.uppercase
			v.process()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+n"))):
			v.noHyphens = !v.noHyphens
			v.process()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+v"))):
			v.mode = uuidModeValidate
			v.input.Placeholder = "Enter UUID to validate..."
			v.input.Focus()
			v.process()
			return v, textarea.Blink
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+p"))):
			v.mode = uuidModeParse
			v.input.Placeholder = "Enter UUID to parse..."
			v.input.Focus()
			v.process()
			return v, textarea.Blink
		}
	}

	// In generate mode, do NOT forward key events to the textarea.
	if v.mode != uuidModeGenerate {
		var cmd tea.Cmd
		v.input, cmd = v.input.Update(msg)
		cmds = append(cmds, cmd)
	}

	v.process()

	return v, tea.Batch(cmds...)
}

func (v *UUIDView) process() {
	switch v.mode {
	case uuidModeGenerate:
		result := tools.UUIDGenerate(v.version, v.uppercase, v.noHyphens)
		if result.Error != "" {
			v.err = result.Error
			v.output.SetContent("")
		} else {
			v.err = ""
			v.output.SetContent(result.Output)
		}
	case uuidModeValidate:
		input := v.input.Value()
		if input == "" {
			v.output.SetContent("")
			v.err = ""
			return
		}
		result := tools.UUIDValidate(input)
		if result.Error != "" {
			v.err = result.Error
			v.output.SetContent("")
		} else {
			v.err = ""
			v.output.SetContent(result.Output)
		}
	case uuidModeParse:
		input := v.input.Value()
		if input == "" {
			v.output.SetContent("")
			v.err = ""
			return
		}
		result := tools.UUIDParse(input)
		if result.Error != "" {
			v.err = result.Error
			v.output.SetContent("")
		} else {
			v.err = ""
			v.output.SetContent(result.Output)
		}
	}
}

func (v *UUIDView) View() string {
	title := styles.TitleStyle.Render("UUID Generate / Validate / Parse")

	var modeStr string
	switch v.mode {
	case uuidModeGenerate:
		modeStr = styles.ModeActivePill.Render("Generate") + "  " +
			styles.ModeInactivePill.Render("Validate") + "  " +
			styles.ModeInactivePill.Render("Parse")
	case uuidModeValidate:
		modeStr = styles.ModeInactivePill.Render("Generate") + "  " +
			styles.ModeActivePill.Render("Validate") + "  " +
			styles.ModeInactivePill.Render("Parse")
	case uuidModeParse:
		modeStr = styles.ModeInactivePill.Render("Generate") + "  " +
			styles.ModeInactivePill.Render("Validate") + "  " +
			styles.ModeActivePill.Render("Parse")
	}

	versionStr := fmt.Sprintf("v%d", v.version)
	if v.mode == uuidModeGenerate {
		if v.version == 4 {
			versionStr = styles.ModeActivePill.Render("v4") + "  " + styles.ModeInactivePill.Render("v7")
		} else {
			versionStr = styles.ModeInactivePill.Render("v4") + "  " + styles.ModeActivePill.Render("v7")
		}
	}

	var uppercaseStr string
	if v.uppercase {
		uppercaseStr = styles.CheckboxOnStyle.Render("● Uppercase")
	} else {
		uppercaseStr = styles.CheckboxOffStyle.Render("○ Uppercase")
	}
	var noHyphensStr string
	if v.noHyphens {
		noHyphensStr = styles.CheckboxOnStyle.Render("● No hyphens")
	} else {
		noHyphensStr = styles.CheckboxOffStyle.Render("○ No hyphens")
	}

	options := fmt.Sprintf("Mode: %s    %s    %s    %s",
		modeStr, versionStr, uppercaseStr, noHyphensStr)

	var inputSection string
	if v.mode == uuidModeGenerate {
		inputSection = styles.LabelStyle.Render("") // hidden in generate mode
	} else {
		inputSection = styles.LabelStyle.Render("Input:") + "\n" + v.input.View()
	}

	var outputSection string
	if v.err != "" {
		outputSection = styles.LabelStyle.Render("Error:") + "\n" + styles.ErrorTextStyle.Render(v.err)
	} else {
		outputSection = styles.LabelStyle.Render("Output:") + "\n" + v.output.View()
	}

	if v.mode == uuidModeGenerate {
		return fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s",
			title, options, outputSection, inputSection)
	}

	return fmt.Sprintf("%s\n\n%s\n\n%s\n%s\n\n%s",
		title, options, inputSection, "", outputSection)
}

func (v *UUIDView) KeyHints() string {
	return hint("ctrl+g", "generate") + "  " + hint("ctrl+v", "validate") + "  " + hint("ctrl+p", "parse") + "  " + hint("ctrl+u", "uppercase")
}

func (v *UUIDView) SetSize(width, height int) {
	v.width = width
	v.height = height

	inputHeight := max((height - 10) / 2, 3)
	outputHeight := max((height - 10) / 2, 3)

	v.input.SetWidth(width - 4)
	v.input.SetHeight(inputHeight)
	v.output.SetWidth(width - 4)
	v.output.SetHeight(outputHeight)
}
