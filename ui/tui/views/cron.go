package views

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	"github.com/StanMarek/forge/core/tools"
	"github.com/StanMarek/forge/ui/tui/styles"
)

// CronView is the TUI view for cron expression parsing.
type CronView struct {
	input  textarea.Model
	output viewport.Model
	width  int
	height int
	err    string
}

// NewCronView creates a new Cron Expression Parser tool view.
func NewCronView() *CronView {
	ti := textarea.New()
	ti.Placeholder = "Enter cron expression (e.g. */5 * * * *)..."
	ti.Focus()

	vp := viewport.New()

	return &CronView{
		input:  ti,
		output: vp,
	}
}

func (v *CronView) Init() tea.Cmd {
	return textarea.Blink
}

func (v *CronView) Update(msg tea.Msg) (ToolView, tea.Cmd) {
	var cmds []tea.Cmd

	var cmd tea.Cmd
	v.input, cmd = v.input.Update(msg)
	cmds = append(cmds, cmd)

	v.process()

	return v, tea.Batch(cmds...)
}

func (v *CronView) process() {
	input := v.input.Value()
	if input == "" {
		v.output.SetContent("")
		v.err = ""
		return
	}

	result := tools.CronParse(input)

	if result.Error != "" {
		v.err = result.Error
		v.output.SetContent("")
	} else {
		v.err = ""
		v.output.SetContent(result.Output)
	}
}

func (v *CronView) View() string {
	title := styles.TitleStyle.Render("Cron Expression Parser")

	inputLabel := styles.LabelStyle.Render("Input:")
	inputView := v.input.View()

	var outputSection string
	if v.err != "" {
		outputSection = styles.LabelStyle.Render("Error:") + "\n" + styles.ErrorTextStyle.Render(v.err)
	} else {
		outputSection = styles.LabelStyle.Render("Output:") + "\n" + v.output.View()
	}

	return fmt.Sprintf("%s\n\n%s\n%s\n\n%s",
		title, inputLabel, inputView, outputSection)
}

func (v *CronView) KeyHints() string {
	return ""
}

func (v *CronView) SetSize(width, height int) {
	v.width = width
	v.height = height

	inputHeight := max((height-10)/2, 3)
	outputHeight := max((height-10)/2, 3)

	v.input.SetWidth(width - 4)
	v.input.SetHeight(inputHeight)
	v.output.SetWidth(width - 4)
	v.output.SetHeight(outputHeight)
}
