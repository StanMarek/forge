package tools

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------- Tool metadata ----------

func TestXMLTool_Metadata(t *testing.T) {
	tool := XMLTool{}

	assert.Equal(t, "XML Formatter", tool.Name())
	assert.Equal(t, "xml", tool.ID())
	assert.Equal(t, "Formatters", tool.Category())
	assert.NotEmpty(t, tool.Description())

	keywords := tool.Keywords()
	for _, kw := range []string{"xml", "format", "pretty", "minify"} {
		assert.Contains(t, keywords, kw)
	}
}

func TestXMLTool_DetectFromClipboard(t *testing.T) {
	tool := XMLTool{}
	assert.True(t, tool.DetectFromClipboard("<root/>"))
	assert.True(t, tool.DetectFromClipboard("  <html>content</html>  "))
	assert.False(t, tool.DetectFromClipboard("just text"))
	assert.False(t, tool.DetectFromClipboard(""))
	assert.False(t, tool.DetectFromClipboard("< no closing"))
}

// ---------- XMLFormat ----------

func TestXMLFormat_Simple(t *testing.T) {
	input := `<root><child>value</child></root>`
	r := XMLFormat(input)
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "<root>")
	assert.Contains(t, r.Output, "  <child>value</child>")
	assert.Contains(t, r.Output, "</root>")
}

func TestXMLFormat_Nested(t *testing.T) {
	input := `<a><b><c>deep</c></b></a>`
	r := XMLFormat(input)
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "    <c>deep</c>")
}

func TestXMLFormat_AlreadyFormatted(t *testing.T) {
	input := "<root>\n  <child>value</child>\n</root>"
	r := XMLFormat(input)
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "<child>value</child>")
}

func TestXMLFormat_WithAttributes(t *testing.T) {
	input := `<item id="1" name="test">content</item>`
	r := XMLFormat(input)
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "content")
}

func TestXMLFormat_InvalidXML(t *testing.T) {
	r := XMLFormat("<unclosed>")
	assert.NotEmpty(t, r.Error)
	assert.Contains(t, r.Error, "invalid XML")
}

func TestXMLFormat_Empty(t *testing.T) {
	r := XMLFormat("")
	assert.NotEmpty(t, r.Error)
	assert.Equal(t, "empty input", r.Error)
}

// ---------- XMLMinify ----------

func TestXMLMinify_Basic(t *testing.T) {
	input := "<root>\n  <child>value</child>\n</root>"
	r := XMLMinify(input)
	require.Empty(t, r.Error)
	assert.False(t, strings.Contains(r.Output, "\n"), "minified output should not contain newlines")
	assert.Contains(t, r.Output, "<root><child>value</child></root>")
}

func TestXMLMinify_AlreadyMinified(t *testing.T) {
	input := `<root><child>value</child></root>`
	r := XMLMinify(input)
	require.Empty(t, r.Error)
	assert.Equal(t, `<root><child>value</child></root>`, r.Output)
}

func TestXMLMinify_InvalidXML(t *testing.T) {
	r := XMLMinify("<bad")
	assert.NotEmpty(t, r.Error)
	assert.Contains(t, r.Error, "invalid XML")
}

func TestXMLMinify_Empty(t *testing.T) {
	r := XMLMinify("")
	assert.NotEmpty(t, r.Error)
	assert.Equal(t, "empty input", r.Error)
}

// ---------- Tool interface compliance ----------

func TestXMLTool_ImplementsToolInterface(t *testing.T) {
	var _ Tool = XMLTool{}
}
