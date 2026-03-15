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

type gzipMode int

const (
	gzipCompress gzipMode = iota
	gzipDecompress
)

// GZipView is the TUI view for gzip compression/decompression.
type GZipView struct {
	input  textarea.Model
	output viewport.Model
	mode   gzipMode
	width  int
	height int
	err    string
}

// NewGZipView creates a new GZip Compress / Decompress tool view.
func NewGZipView() *GZipView {
	ti := textarea.New()
	ti.Placeholder = "Enter text to compress..."
	ti.Focus()

	vp := viewport.New()

	return &GZipView{
		input:  ti,
		output: vp,
		mode:   gzipCompress,
	}
}

func (v *GZipView) Init() tea.Cmd {
	return textarea.Blink
}

func (v *GZipView) Update(msg tea.Msg) (ToolView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))):
			v.mode = gzipCompress
			v.input.Placeholder = "Enter text to compress..."
			v.process()
			return v, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+d"))):
			v.mode = gzipDecompress
			v.input.Placeholder = "Enter base64 gzip data to decompress..."
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

func (v *GZipView) process() {
	input := v.input.Value()
	if input == "" {
		v.output.SetContent("")
		v.err = ""
		return
	}

	var result tools.Result
	switch v.mode {
	case gzipCompress:
		result = tools.GZipCompress(input)
	case gzipDecompress:
		result = tools.GZipDecompress(input)
	}

	if result.Error != "" {
		v.err = result.Error
		v.output.SetContent("")
	} else {
		v.err = ""
		v.output.SetContent(result.Output)
	}
}

func (v *GZipView) View() string {
	title := styles.TitleStyle.Render("GZip Compress / Decompress")

	var modeStr string
	if v.mode == gzipCompress {
		modeStr = styles.ModeActivePill.Render("Compress") + " " + styles.ModeInactivePill.Render("Decompress")
	} else {
		modeStr = styles.ModeInactivePill.Render("Compress") + " " + styles.ModeActivePill.Render("Decompress")
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

func (v *GZipView) KeyHints() string {
	return hint("ctrl+c", "compress") + "  " + hint("ctrl+d", "decompress")
}

func (v *GZipView) SetSize(width, height int) {
	v.width = width
	v.height = height

	inputHeight := max((height-10)/2, 3)
	outputHeight := max((height-10)/2, 3)

	v.input.SetWidth(width - 4)
	v.input.SetHeight(inputHeight)
	v.output.SetWidth(width - 4)
	v.output.SetHeight(outputHeight)
}
