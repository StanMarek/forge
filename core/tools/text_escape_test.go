package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------- Tool metadata ----------

func TestTextEscapeTool_Metadata(t *testing.T) {
	tool := TextEscapeTool{}

	assert.Equal(t, "Text Escape / Unescape", tool.Name())
	assert.Equal(t, "text-escape", tool.ID())
	assert.Equal(t, "Encoders", tool.Category())
	assert.NotEmpty(t, tool.Description())

	keywords := tool.Keywords()
	for _, kw := range []string{"escape", "unescape", "text", "string", "backslash"} {
		assert.Contains(t, keywords, kw)
	}
}

func TestTextEscapeTool_DetectFromClipboard(t *testing.T) {
	tool := TextEscapeTool{}

	assert.True(t, tool.DetectFromClipboard(`hello\nworld`))
	assert.True(t, tool.DetectFromClipboard(`col1\tcol2`))
	assert.True(t, tool.DetectFromClipboard(`line\r`))
	assert.True(t, tool.DetectFromClipboard(`path\\to\\file`))
	assert.True(t, tool.DetectFromClipboard(`say \"hi\"`))
	assert.False(t, tool.DetectFromClipboard("plain text"))
	assert.False(t, tool.DetectFromClipboard(""))
}

// ---------- TextEscape ----------

func TestTextEscape_MultilineString(t *testing.T) {
	input := "hello\nworld\ttab"
	r := TextEscape(input)
	require.Empty(t, r.Error)
	assert.Equal(t, `hello\nworld\ttab`, r.Output)
}

func TestTextEscape_SpecialCharacters(t *testing.T) {
	input := "line1\nline2\r\nend\t\b\f\"\\"
	r := TextEscape(input)
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, `\n`)
	assert.Contains(t, r.Output, `\r`)
	assert.Contains(t, r.Output, `\t`)
	assert.Contains(t, r.Output, `\\`)
	assert.Contains(t, r.Output, `\"`)
}

func TestTextEscape_EmptyInput(t *testing.T) {
	r := TextEscape("")
	require.Empty(t, r.Error)
	assert.Equal(t, "", r.Output)
}

func TestTextEscape_AlreadyPlain(t *testing.T) {
	r := TextEscape("no special chars")
	require.Empty(t, r.Error)
	assert.Equal(t, "no special chars", r.Output)
}

// ---------- TextUnescape ----------

func TestTextUnescape_Basic(t *testing.T) {
	r := TextUnescape(`hello\nworld`)
	require.Empty(t, r.Error)
	assert.Equal(t, "hello\nworld", r.Output)
}

func TestTextUnescape_Tab(t *testing.T) {
	r := TextUnescape(`col1\tcol2`)
	require.Empty(t, r.Error)
	assert.Equal(t, "col1\tcol2", r.Output)
}

func TestTextUnescape_EmptyInput(t *testing.T) {
	r := TextUnescape("")
	require.Empty(t, r.Error)
	assert.Equal(t, "", r.Output)
}

func TestTextUnescape_AlreadyQuoted(t *testing.T) {
	r := TextUnescape(`"hello\nworld"`)
	require.Empty(t, r.Error)
	assert.Equal(t, "hello\nworld", r.Output)
}

// ---------- Roundtrip ----------

func TestTextEscape_Roundtrip(t *testing.T) {
	original := "line1\nline2\ttab\r\nend"
	escaped := TextEscape(original)
	require.Empty(t, escaped.Error)

	unescaped := TextUnescape(escaped.Output)
	require.Empty(t, unescaped.Error)
	assert.Equal(t, original, unescaped.Output)
}

// ---------- Tool interface compliance ----------

func TestTextEscapeTool_ImplementsToolInterface(t *testing.T) {
	var _ Tool = TextEscapeTool{}
}
