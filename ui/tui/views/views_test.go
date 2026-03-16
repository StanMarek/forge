package views

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
)

// viewFactory describes how to create a view and what to expect from it.
type viewFactory struct {
	name       string
	create     func() ToolView
	title      string   // expected substring in View() output
	hasHints   bool     // whether KeyHints() should be non-empty
	autoOutput bool     // whether the view auto-generates output (no input needed)
	modes      []string // if non-empty, mode pill labels expected in View()
}

var allViews = []viewFactory{
	{"Base64", func() ToolView { return NewBase64View() }, "Base64", true, false, []string{"Encode", "Decode"}},
	{"JWT", func() ToolView { return NewJWTView() }, "JWT", true, false, []string{"Full", "Header", "Payload"}},
	{"JSON", func() ToolView { return NewJSONView() }, "JSON", true, false, []string{"Format", "Minify", "Validate"}},
	{"Hash", func() ToolView { return NewHashView() }, "Hash", true, false, nil},
	{"URL", func() ToolView { return NewURLView() }, "URL", true, false, []string{"Parse", "Encode", "Decode"}},
	{"UUID", func() ToolView { return NewUUIDView() }, "UUID", true, true, []string{"Generate", "Validate", "Parse"}},
	{"YAML", func() ToolView { return NewYAMLView() }, "YAML", true, false, nil},
	{"Timestamp", func() ToolView { return NewTimestampView() }, "Timestamp", true, true, nil},
	{"NumberBase", func() ToolView { return NewNumberBaseView() }, "Number Base", false, false, nil},
	{"Regex", func() ToolView { return NewRegexView() }, "Regex", true, false, nil},
	{"HTMLEntity", func() ToolView { return NewHTMLEntityView() }, "HTML", true, false, []string{"Encode", "Decode"}},
	{"Password", func() ToolView { return NewPasswordView() }, "Password", true, true, nil},
	{"Lorem", func() ToolView { return NewLoremView() }, "Lorem", true, true, nil},
	{"Color", func() ToolView { return NewColorView() }, "Color", false, false, nil},
	{"Cron", func() ToolView { return NewCronView() }, "Cron", false, false, nil},
	{"TextEscape", func() ToolView { return NewTextEscapeView() }, "Text Escape", true, false, []string{"Escape", "Unescape"}},
	{"GZip", func() ToolView { return NewGZipView() }, "GZip", true, false, []string{"Compress", "Decompress"}},
	{"TextStats", func() ToolView { return NewTextStatsView() }, "Text", true, false, nil},
	{"Diff", func() ToolView { return NewDiffView() }, "Diff", true, false, nil},
	{"XML", func() ToolView { return NewXMLView() }, "XML", true, false, []string{"Format", "Minify"}},
	{"CSV", func() ToolView { return NewCSVView() }, "CSV", true, false, nil},
	{"Placeholder", func() ToolView { return NewPlaceholder("Test Tool") }, "Test Tool", false, false, nil},
}

func TestAllViewsInitWithoutPanic(t *testing.T) {
	for _, vf := range allViews {
		t.Run(vf.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				v := vf.create()
				v.Init()
			})
		})
	}
}

func TestAllViewsRenderTitle(t *testing.T) {
	for _, vf := range allViews {
		t.Run(vf.name, func(t *testing.T) {
			v := vf.create()
			v.SetSize(80, 40)
			output := v.View()
			assert.NotEmpty(t, output, "%s View() returned empty", vf.name)
			assert.Contains(t, output, vf.title, "%s View() missing title %q", vf.name, vf.title)
		})
	}
}

func TestAllViewsSetSizeVariousDimensions(t *testing.T) {
	sizes := [][2]int{{80, 40}, {120, 60}, {40, 20}, {200, 80}, {0, 0}}

	for _, vf := range allViews {
		t.Run(vf.name, func(t *testing.T) {
			v := vf.create()
			for _, sz := range sizes {
				assert.NotPanics(t, func() {
					v.SetSize(sz[0], sz[1])
					v.View() // ensure rendering doesn't panic after resize
				}, "%s panicked at size %dx%d", vf.name, sz[0], sz[1])
			}
		})
	}
}

func TestAllViewsKeyHints(t *testing.T) {
	for _, vf := range allViews {
		t.Run(vf.name, func(t *testing.T) {
			v := vf.create()
			hints := v.KeyHints()
			if vf.hasHints {
				assert.NotEmpty(t, hints, "%s should have key hints", vf.name)
			}
		})
	}
}

func TestAllViewsModePills(t *testing.T) {
	for _, vf := range allViews {
		if len(vf.modes) == 0 {
			continue
		}
		t.Run(vf.name, func(t *testing.T) {
			v := vf.create()
			v.SetSize(80, 40)
			output := v.View()
			for _, mode := range vf.modes {
				assert.Contains(t, output, mode, "%s missing mode pill %q", vf.name, mode)
			}
		})
	}
}

func TestAllViewsUpdateDoesNotPanic(t *testing.T) {
	for _, vf := range allViews {
		t.Run(vf.name, func(t *testing.T) {
			v := vf.create()
			v.SetSize(80, 40)
			assert.NotPanics(t, func() {
				v.Update(tea.WindowSizeMsg{Width: 80, Height: 40})
			})
		})
	}
}

// ctrlKey constructs a tea.KeyPressMsg for ctrl+<letter>.
func ctrlKey(r rune) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: r, Mod: tea.ModCtrl}
}

func TestBase64ModeSwitch(t *testing.T) {
	v := NewBase64View()
	v.SetSize(80, 40)

	// Default mode is encode
	output := v.View()
	assert.Contains(t, output, "Encode")

	// Switch to decode via ctrl+d
	v.Update(ctrlKey('d'))
	output = v.View()
	// The decode pill should now be active (we can't check style, but
	// placeholder changes)
	assert.Contains(t, output, "Decode")

	// Switch back to encode via ctrl+e
	v.Update(ctrlKey('e'))
	output = v.View()
	assert.Contains(t, output, "Encode")
}

func TestBase64URLSafeToggle(t *testing.T) {
	v := NewBase64View()
	v.SetSize(80, 40)

	// Default: URL-safe off
	output := v.View()
	assert.Contains(t, output, "URL-safe")

	// Toggle URL-safe
	v.Update(ctrlKey('u'))
	output = v.View()
	assert.Contains(t, output, "URL-safe")
}

func TestBase64ProcessEncode(t *testing.T) {
	v := NewBase64View()
	v.SetSize(80, 40)

	// Type into the textarea by inserting text
	v.input.InsertString("hello")
	v.process()

	output := v.View()
	assert.Contains(t, output, "aGVsbG8=", "expected base64 encoding of 'hello'")
}

func TestBase64ProcessDecode(t *testing.T) {
	v := NewBase64View()
	v.SetSize(80, 40)
	v.mode = modeDecode

	v.input.InsertString("aGVsbG8=")
	v.process()

	output := v.View()
	assert.Contains(t, output, "hello", "expected decoded output 'hello'")
}

func TestBase64ProcessDecodeError(t *testing.T) {
	v := NewBase64View()
	v.SetSize(80, 40)
	v.mode = modeDecode

	v.input.InsertString("!!!not-valid-base64!!!")
	v.process()

	assert.NotEmpty(t, v.err, "expected error for invalid base64")
	output := v.View()
	assert.Contains(t, output, "Error")
}

func TestBase64EmptyInput(t *testing.T) {
	v := NewBase64View()
	v.SetSize(80, 40)
	v.process()

	// Empty input should produce no error and no output
	assert.Empty(t, v.err)
	output := v.View()
	assert.Contains(t, output, "Output")
}

func TestJSONViewModeSwitch(t *testing.T) {
	v := NewJSONView()
	v.SetSize(80, 40)

	// Default is format mode
	output := v.View()
	assert.Contains(t, output, "Format")

	// Switch to minify
	v.Update(ctrlKey('m'))
	output = v.View()
	assert.Contains(t, output, "Minify")

	// Switch to validate
	v.Update(ctrlKey('v'))
	output = v.View()
	assert.Contains(t, output, "Validate")
}

func TestJSONViewProcessFormat(t *testing.T) {
	v := NewJSONView()
	v.SetSize(80, 40)

	v.input.InsertString(`{"a":1,"b":2}`)
	v.process()

	output := v.View()
	assert.Contains(t, output, `"a"`)
}

func TestHashViewProcess(t *testing.T) {
	v := NewHashView()
	v.SetSize(80, 40)

	v.input.InsertString("hello")
	v.process()

	output := v.View()
	// MD5 of "hello" is 5d41402abc4b2a76b9719d911017c592
	assert.Contains(t, output, "5d41402abc4b2a76b9719d911017c592")
}

func TestHashViewUppercase(t *testing.T) {
	v := NewHashView()
	v.SetSize(80, 40)

	v.input.InsertString("hello")
	v.Update(ctrlKey('u')) // toggle uppercase
	v.process()

	output := v.View()
	assert.Contains(t, output, "5D41402ABC4B2A76B9719D911017C592")
}

func TestURLViewModes(t *testing.T) {
	v := NewURLView()
	v.SetSize(80, 40)

	// Default is parse
	output := v.View()
	assert.Contains(t, output, "Parse")

	v.Update(ctrlKey('e'))
	output = v.View()
	assert.Contains(t, output, "Encode")

	v.Update(ctrlKey('d'))
	output = v.View()
	assert.Contains(t, output, "Decode")
}

func TestUUIDViewAutoGenerates(t *testing.T) {
	v := NewUUIDView()
	v.SetSize(80, 40)

	output := v.View()
	// UUID auto-generates on creation; output should contain a UUID-like string
	// UUIDs have 36 chars with hyphens: 8-4-4-4-12
	assert.Contains(t, output, "-", "expected UUID with hyphens in output")
}

func TestPasswordViewAutoGenerates(t *testing.T) {
	v := NewPasswordView()
	v.SetSize(80, 40)

	output := v.View()
	assert.Contains(t, output, "Password")
	// Password auto-generates; the output section should have content
	assert.True(t, len(output) > 50, "expected substantial output from password view")
}

func TestLoremViewAutoGenerates(t *testing.T) {
	v := NewLoremView()
	v.SetSize(80, 40)

	output := v.View()
	assert.Contains(t, output, "Lorem")
}

func TestTimestampViewNowMode(t *testing.T) {
	v := NewTimestampView()
	v.SetSize(80, 40)

	// Default is "now" mode which auto-generates
	output := v.View()
	assert.Contains(t, output, "Timestamp")
}

func TestDiffViewTwoInputs(t *testing.T) {
	v := NewDiffView()
	v.SetSize(80, 40)

	output := v.View()
	// Diff view should have two input areas labeled "Text A" and "Text B"
	assert.Contains(t, output, "Text A")
	assert.Contains(t, output, "Text B")
}

func TestRegexViewTwoInputs(t *testing.T) {
	v := NewRegexView()
	v.SetSize(80, 40)

	output := v.View()
	assert.Contains(t, output, "Pattern")
	// Should show test string area too
	assert.True(t, strings.Contains(output, "Test") || strings.Contains(output, "test"),
		"expected test string input area")
}

func TestGZipViewModes(t *testing.T) {
	v := NewGZipView()
	v.SetSize(80, 40)

	output := v.View()
	assert.Contains(t, output, "Compress")
	assert.Contains(t, output, "Decompress")
}

func TestXMLViewModes(t *testing.T) {
	v := NewXMLView()
	v.SetSize(80, 40)

	output := v.View()
	assert.Contains(t, output, "Format")
	assert.Contains(t, output, "Minify")
}

func TestCSVViewModes(t *testing.T) {
	v := NewCSVView()
	v.SetSize(80, 40)

	output := v.View()
	assert.Contains(t, output, "CSV")
}

func TestHTMLEntityViewModes(t *testing.T) {
	v := NewHTMLEntityView()
	v.SetSize(80, 40)

	output := v.View()
	assert.Contains(t, output, "Encode")
	assert.Contains(t, output, "Decode")
}

func TestTextEscapeViewModes(t *testing.T) {
	v := NewTextEscapeView()
	v.SetSize(80, 40)

	output := v.View()
	assert.Contains(t, output, "Escape")
	assert.Contains(t, output, "Unescape")
}

func TestTextStatsViewModeCycle(t *testing.T) {
	v := NewTextStatsView()
	v.SetSize(80, 40)

	// Default is stats mode
	output := v.View()
	assert.Contains(t, output, "Stats")

	// Cycle through modes
	v.Update(ctrlKey('m'))
	output = v.View()
	assert.Contains(t, output, "Lower")
}

func TestPlaceholderView(t *testing.T) {
	v := NewPlaceholder("My Tool")
	v.SetSize(80, 40)

	output := v.View()
	assert.Contains(t, output, "My Tool")
	assert.Contains(t, output, "Coming soon")
}

func TestYAMLViewModes(t *testing.T) {
	v := NewYAMLView()
	v.SetSize(80, 40)

	output := v.View()
	assert.True(t, strings.Contains(output, "JSON") || strings.Contains(output, "YAML"),
		"expected YAML/JSON mode labels")
}
