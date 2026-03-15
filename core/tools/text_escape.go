package tools

import (
	"strconv"
	"strings"
)

// TextEscapeTool provides metadata for the Text Escape / Unescape tool.
type TextEscapeTool struct{}

func (t TextEscapeTool) Name() string        { return "Text Escape / Unescape" }
func (t TextEscapeTool) ID() string          { return "text-escape" }
func (t TextEscapeTool) Description() string { return "Escape and unescape special characters in text" }
func (t TextEscapeTool) Category() string    { return "Encoders" }
func (t TextEscapeTool) Keywords() []string {
	return []string{"escape", "unescape", "text", "string", "backslash"}
}

// DetectFromClipboard returns true if s contains literal escape sequences
// such as \n, \t, \r, \\, or \".
func (t TextEscapeTool) DetectFromClipboard(s string) bool {
	return strings.Contains(s, `\n`) ||
		strings.Contains(s, `\t`) ||
		strings.Contains(s, `\r`) ||
		strings.Contains(s, `\\`) ||
		strings.Contains(s, `\"`)
}

// TextEscape escapes special characters in the input string.
// It uses strconv.Quote and strips the surrounding quotes.
func TextEscape(input string) Result {
	if input == "" {
		return Result{Output: ""}
	}

	quoted := strconv.Quote(input)
	// strconv.Quote wraps the result in double quotes; strip them.
	escaped := quoted[1 : len(quoted)-1]

	return Result{Output: escaped}
}

// TextUnescape unescapes escape sequences in the input string.
// It replaces literal \n, \t, etc. with the actual characters.
// Uses strconv.Unquote, wrapping the input in quotes if needed.
func TextUnescape(input string) Result {
	if input == "" {
		return Result{Output: ""}
	}

	// If the input is not already wrapped in double quotes, wrap it
	// so strconv.Unquote can process it.
	toUnquote := input
	if !strings.HasPrefix(input, `"`) || !strings.HasSuffix(input, `"`) {
		toUnquote = `"` + input + `"`
	}

	unescaped, err := strconv.Unquote(toUnquote)
	if err != nil {
		return Result{Error: "invalid escape sequence: " + err.Error()}
	}

	return Result{Output: unescaped}
}
