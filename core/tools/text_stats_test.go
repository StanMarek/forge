package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------- Tool metadata ----------

func TestTextStatsTool_Metadata(t *testing.T) {
	tool := TextStatsTool{}

	assert.Equal(t, "Text Analyzer", tool.Name())
	assert.Equal(t, "text-stats", tool.ID())
	assert.Equal(t, "Text", tool.Category())
	assert.NotEmpty(t, tool.Description())

	keywords := tool.Keywords()
	for _, kw := range []string{"text", "stats", "count", "words", "characters", "case", "convert"} {
		assert.Contains(t, keywords, kw)
	}
}

func TestTextStatsTool_DetectFromClipboard(t *testing.T) {
	tool := TextStatsTool{}
	assert.False(t, tool.DetectFromClipboard("anything"))
	assert.False(t, tool.DetectFromClipboard(""))
}

// ---------- TextStats ----------

func TestTextStats_HelloWorld(t *testing.T) {
	r := TextStats("Hello, World! How are you?")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "Characters:  26")
	assert.Contains(t, r.Output, "Words:       5")
	assert.Contains(t, r.Output, "Lines:       1")
	assert.Contains(t, r.Output, "Sentences:   2") // '!' and '?'
	assert.Contains(t, r.Output, "Bytes:       26")
}

func TestTextStats_EmptyString(t *testing.T) {
	r := TextStats("")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "Characters:  0")
	assert.Contains(t, r.Output, "Words:       0")
	assert.Contains(t, r.Output, "Lines:       0")
	assert.Contains(t, r.Output, "Sentences:   0")
	assert.Contains(t, r.Output, "Bytes:       0")
}

func TestTextStats_Multiline(t *testing.T) {
	input := "Line one.\nLine two.\nLine three."
	r := TextStats(input)
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "Lines:       3")
	assert.Contains(t, r.Output, "Words:       6")
	assert.Contains(t, r.Output, "Sentences:   3")
}

func TestTextStats_MultipleSentences(t *testing.T) {
	input := "Hello! How are you? I'm fine."
	r := TextStats(input)
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "Sentences:   3")
}

func TestTextStats_UnicodeCharacters(t *testing.T) {
	input := "caf\u00e9" // "café"
	r := TextStats(input)
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "Characters:  4")
	assert.Contains(t, r.Output, "Bytes:       5") // é is 2 bytes in UTF-8
}

// ---------- TextCaseConvert ----------

func TestTextCaseConvert_Lower(t *testing.T) {
	r := TextCaseConvert("Hello World", "lower")
	require.Empty(t, r.Error)
	assert.Equal(t, "hello world", r.Output)
}

func TestTextCaseConvert_Upper(t *testing.T) {
	r := TextCaseConvert("Hello World", "upper")
	require.Empty(t, r.Error)
	assert.Equal(t, "HELLO WORLD", r.Output)
}

func TestTextCaseConvert_Title(t *testing.T) {
	r := TextCaseConvert("hello world", "title")
	require.Empty(t, r.Error)
	assert.Equal(t, "Hello World", r.Output)
}

func TestTextCaseConvert_Camel(t *testing.T) {
	r := TextCaseConvert("hello world", "camel")
	require.Empty(t, r.Error)
	assert.Equal(t, "helloWorld", r.Output)
}

func TestTextCaseConvert_CamelFromSnake(t *testing.T) {
	r := TextCaseConvert("hello_world_test", "camel")
	require.Empty(t, r.Error)
	assert.Equal(t, "helloWorldTest", r.Output)
}

func TestTextCaseConvert_CamelFromKebab(t *testing.T) {
	r := TextCaseConvert("hello-world-test", "camel")
	require.Empty(t, r.Error)
	assert.Equal(t, "helloWorldTest", r.Output)
}

func TestTextCaseConvert_Snake(t *testing.T) {
	r := TextCaseConvert("Hello World", "snake")
	require.Empty(t, r.Error)
	assert.Equal(t, "hello_world", r.Output)
}

func TestTextCaseConvert_SnakeFromKebab(t *testing.T) {
	r := TextCaseConvert("hello-world-test", "snake")
	require.Empty(t, r.Error)
	assert.Equal(t, "hello_world_test", r.Output)
}

func TestTextCaseConvert_Kebab(t *testing.T) {
	r := TextCaseConvert("Hello World", "kebab")
	require.Empty(t, r.Error)
	assert.Equal(t, "hello-world", r.Output)
}

func TestTextCaseConvert_KebabFromSnake(t *testing.T) {
	r := TextCaseConvert("hello_world_test", "kebab")
	require.Empty(t, r.Error)
	assert.Equal(t, "hello-world-test", r.Output)
}

func TestTextCaseConvert_UnsupportedMode(t *testing.T) {
	r := TextCaseConvert("hello", "pascal")
	assert.NotEmpty(t, r.Error)
	assert.Contains(t, r.Error, "unsupported case mode: pascal")
	assert.Contains(t, r.Error, "supported: lower, upper, title, camel, snake, kebab")
	assert.Empty(t, r.Output)
}

func TestTextCaseConvert_EmptyInput(t *testing.T) {
	r := TextCaseConvert("", "lower")
	require.Empty(t, r.Error)
	assert.Equal(t, "", r.Output)
}

func TestTextCaseConvert_CaseInsensitiveMode(t *testing.T) {
	r := TextCaseConvert("Hello World", "UPPER")
	require.Empty(t, r.Error)
	assert.Equal(t, "HELLO WORLD", r.Output)
}

// ---------- Tool interface compliance ----------

func TestTextStatsTool_ImplementsToolInterface(t *testing.T) {
	var _ Tool = TextStatsTool{}
}
