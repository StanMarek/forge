package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------- Tool metadata ----------

func TestRegexTool_Metadata(t *testing.T) {
	tool := RegexTool{}

	assert.Equal(t, "Regex Tester", tool.Name())
	assert.Equal(t, "regex", tool.ID())
	assert.Equal(t, "Testers", tool.Category())
	assert.NotEmpty(t, tool.Description())

	keywords := tool.Keywords()
	for _, kw := range []string{"regex", "regexp", "pattern", "match", "test"} {
		assert.Contains(t, keywords, kw)
	}
}

func TestRegexTool_DetectFromClipboard(t *testing.T) {
	tool := RegexTool{}
	assert.False(t, tool.DetectFromClipboard("anything"))
	assert.False(t, tool.DetectFromClipboard(""))
	assert.False(t, tool.DetectFromClipboard(`\d+`))
}

// ---------- Simple digit match ----------

func TestRegex_SimpleDigitMatch(t *testing.T) {
	r := RegexTest(`\d+`, "Order 42 has 3 items", false)
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, `"42"`)
	assert.Contains(t, r.Output, "pos 6-8")
}

// ---------- Global matches ----------

func TestRegex_GlobalMatches(t *testing.T) {
	r := RegexTest(`\d+`, "Order 42 has 3 items", true)
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "Matches:")
	assert.Contains(t, r.Output, `[0] "42" (pos 6-8)`)
	assert.Contains(t, r.Output, `[1] "3" (pos 13-14)`)
}

// ---------- No match ----------

func TestRegex_NoMatch(t *testing.T) {
	r := RegexTest(`\d+`, "no digits here", false)
	require.Empty(t, r.Error)
	assert.Equal(t, "No matches found", r.Output)
}

func TestRegex_NoMatch_Global(t *testing.T) {
	r := RegexTest(`\d+`, "no digits here", true)
	require.Empty(t, r.Error)
	assert.Equal(t, "No matches found", r.Output)
}

// ---------- Invalid pattern ----------

func TestRegex_InvalidPattern(t *testing.T) {
	r := RegexTest(`[invalid`, "test", false)
	assert.NotEmpty(t, r.Error)
	assert.Contains(t, r.Error, "invalid regex")
	assert.Empty(t, r.Output)
}

// ---------- Empty input ----------

func TestRegex_EmptyInput(t *testing.T) {
	r := RegexTest(`.*`, "", false)
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, `Match: ""`)
}

func TestRegex_EmptyInput_NoMatchPattern(t *testing.T) {
	r := RegexTest(`\d+`, "", false)
	require.Empty(t, r.Error)
	assert.Equal(t, "No matches found", r.Output)
}

// ---------- Groups ----------

func TestRegex_Groups(t *testing.T) {
	r := RegexTest(`(\w+)@(\w+)\.(\w+)`, "user@example.com", false)
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, `"user@example.com"`)
	assert.Contains(t, r.Output, "pos 0-16")
}

// ---------- Output format ----------

func TestRegex_OutputContainsPatternAndInput(t *testing.T) {
	r := RegexTest(`hello`, "say hello world", false)
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, `Pattern: hello`)
	assert.Contains(t, r.Output, `Input: "say hello world"`)
}

// ---------- Tool interface compliance ----------

func TestRegexTool_ImplementsToolInterface(t *testing.T) {
	var _ Tool = RegexTool{}
}
