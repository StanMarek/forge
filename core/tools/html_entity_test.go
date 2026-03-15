package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------- Tool metadata ----------

func TestHTMLEntityTool_Metadata(t *testing.T) {
	tool := HTMLEntityTool{}

	assert.Equal(t, "HTML Entity Encoder", tool.Name())
	assert.Equal(t, "html-entity", tool.ID())
	assert.Equal(t, "Encoders", tool.Category())
	assert.NotEmpty(t, tool.Description())

	keywords := tool.Keywords()
	for _, kw := range []string{"html", "entity", "encode", "decode", "escape"} {
		assert.Contains(t, keywords, kw)
	}
}

func TestHTMLEntityTool_DetectFromClipboard(t *testing.T) {
	tool := HTMLEntityTool{}

	// Should detect HTML entities
	assert.True(t, tool.DetectFromClipboard("&amp; test"))
	assert.True(t, tool.DetectFromClipboard("&lt;div&gt;"))
	assert.True(t, tool.DetectFromClipboard("&quot;hello&quot;"))
	assert.True(t, tool.DetectFromClipboard("&#60;"))
	assert.True(t, tool.DetectFromClipboard("&#x3C;"))

	// Should not detect plain text
	assert.False(t, tool.DetectFromClipboard("hello world"))
	assert.False(t, tool.DetectFromClipboard(""))
	assert.False(t, tool.DetectFromClipboard("<div>not encoded</div>"))
}

// ---------- Encode ----------

func TestHTMLEntityEncode_Script(t *testing.T) {
	r := HTMLEntityEncode(`<script>alert("xss")</script>`)
	require.Empty(t, r.Error)
	assert.Equal(t, "&lt;script&gt;alert(&#34;xss&#34;)&lt;/script&gt;", r.Output)
}

func TestHTMLEntityEncode_AllSpecialChars(t *testing.T) {
	r := HTMLEntityEncode(`<>&"'`)
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "&lt;")
	assert.Contains(t, r.Output, "&gt;")
	assert.Contains(t, r.Output, "&amp;")
	assert.Contains(t, r.Output, "&#34;")
}

func TestHTMLEntityEncode_PlainText(t *testing.T) {
	r := HTMLEntityEncode("hello world")
	require.Empty(t, r.Error)
	assert.Equal(t, "hello world", r.Output)
}

func TestHTMLEntityEncode_Empty(t *testing.T) {
	r := HTMLEntityEncode("")
	require.Empty(t, r.Error)
	assert.Equal(t, "", r.Output)
}

// ---------- Decode ----------

func TestHTMLEntityDecode_Div(t *testing.T) {
	r := HTMLEntityDecode("&lt;div&gt;")
	require.Empty(t, r.Error)
	assert.Equal(t, "<div>", r.Output)
}

func TestHTMLEntityDecode_NumericEntities(t *testing.T) {
	r := HTMLEntityDecode("&#60;p&#62;")
	require.Empty(t, r.Error)
	assert.Equal(t, "<p>", r.Output)
}

func TestHTMLEntityDecode_PlainText(t *testing.T) {
	r := HTMLEntityDecode("hello world")
	require.Empty(t, r.Error)
	assert.Equal(t, "hello world", r.Output)
}

func TestHTMLEntityDecode_Empty(t *testing.T) {
	r := HTMLEntityDecode("")
	require.Empty(t, r.Error)
	assert.Equal(t, "", r.Output)
}

// ---------- Roundtrip ----------

func TestHTMLEntity_Roundtrip(t *testing.T) {
	original := `<div class="test">&copy; 2024</div>`
	encoded := HTMLEntityEncode(original)
	require.Empty(t, encoded.Error)

	decoded := HTMLEntityDecode(encoded.Output)
	require.Empty(t, decoded.Error)
	assert.Equal(t, original, decoded.Output)
}

// ---------- Tool interface compliance ----------

func TestHTMLEntityTool_ImplementsToolInterface(t *testing.T) {
	var _ Tool = HTMLEntityTool{}
}
