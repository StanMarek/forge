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
		modeStr = styles.ModeActiveStyle.Render("(●) Generate") + "  " +
			styles.ModeInactiveStyle.Render("( ) Validate") + "  " +
			styles.ModeInactiveStyle.Render("( ) Parse")
	case uuidModeValidate:
		modeStr = styles.ModeInactiveStyle.Render("( ) Generate") + "  " +
			styles.ModeActiveStyle.Render("(●) Validate") + "  " +
			styles.ModeInactiveStyle.Render("( ) Parse")
	case uuidModeParse:
		modeStr = styles.ModeInactiveStyle.Render("( ) Generate") + "  " +
			styles.ModeInactiveStyle.Render("( ) Validate") + "  " +
			styles.ModeActiveStyle.Render("(●) Parse")
	}

	versionStr := fmt.Sprintf("v%d", v.version)
	if v.mode == uuidModeGenerate {
		if v.version == 4 {
			versionStr = styles.ModeActiveStyle.Render("(●) v4") + "  " + styles.ModeInactiveStyle.Render("( ) v7")
		} else {
			versionStr = styles.ModeInactiveStyle.Render("( ) v4") + "  " + styles.ModeActiveStyle.Render("(●) v7")
		}
	}

	uppercaseStr := "☐ Uppercase"
	if v.uppercase {
		uppercaseStr = "☑ Uppercase"
	}
	noHyphensStr := "☐ No hyphens"
	if v.noHyphens {
		noHyphensStr = "☑ No hyphens"
	}

	options := fmt.Sprintf("Mode: %s    %s    %s    %s",
		modeStr, versionStr, styles.LabelStyle.Render(uppercaseStr), styles.LabelStyle.Render(noHyphensStr))

	var inputSection string
	if v.mode == uuidModeGenerate {
		inputSection = styles.LabelStyle.Render("") // hidden in generate mode
	} else {
		inputSection = styles.LabelStyle.Render("Input:") + "\n" + v.input.View()
	}

	var outputSection string
	if v.err != "" {
		outputSection = styles.LabelStyle.Render("Error:") + "\n" + styles.ErrorStyle.Render(v.err)
	} else {
		outputSection = styles.LabelStyle.Render("Output:") + "\n" + v.output.View()
	}

	status := styles.StatusBarStyle.Render("ctrl+g: generate  ctrl+v: validate  ctrl+p: parse  ctrl+4/7: version  ctrl+u: uppercase  ctrl+n: no-hyphens")

	if v.mode == uuidModeGenerate {
		return fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s\n\n%s",
			title, options, outputSection, inputSection, status)
	}

	return fmt.Sprintf("%s\n\n%s\n\n%s\n%s\n\n%s\n\n%s",
		title, options, inputSection, "", outputSection, status)
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
