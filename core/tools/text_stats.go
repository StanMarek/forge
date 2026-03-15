package tools

import (
	"fmt"
	"strings"
	"unicode"
)

// TextStatsTool provides metadata for the Text Analyzer tool.
type TextStatsTool struct{}

func (t TextStatsTool) Name() string        { return "Text Analyzer" }
func (t TextStatsTool) ID() string          { return "text-stats" }
func (t TextStatsTool) Description() string { return "Analyze text statistics and convert text case" }
func (t TextStatsTool) Category() string    { return "Text" }
func (t TextStatsTool) Keywords() []string {
	return []string{"text", "stats", "count", "words", "characters", "case", "convert"}
}

// DetectFromClipboard always returns false for the text stats tool.
func (t TextStatsTool) DetectFromClipboard(_ string) bool {
	return false
}

// TextStats computes statistics about the input text: character count,
// word count, line count, sentence count, and byte count.
func TextStats(input string) Result {
	if input == "" {
		output := fmt.Sprintf("Characters:  0\nWords:       0\nLines:       0\nSentences:   0\nBytes:       0")
		return Result{Output: output}
	}

	chars := len([]rune(input))
	bytes := len(input)

	// Word count: split on whitespace, filter empty
	words := 0
	for _, w := range strings.Fields(input) {
		if w != "" {
			words++
		}
	}

	// Line count: count newlines + 1
	lines := strings.Count(input, "\n") + 1

	// Sentence count: count sentence-ending punctuation
	sentences := 0
	for _, r := range input {
		if r == '.' || r == '!' || r == '?' {
			sentences++
		}
	}

	output := fmt.Sprintf("Characters:  %d\nWords:       %d\nLines:       %d\nSentences:   %d\nBytes:       %d",
		chars, words, lines, sentences, bytes)

	return Result{Output: output}
}

// TextCaseConvert converts the input text to the specified case mode.
// Supported modes: lower, upper, title, camel, snake, kebab.
func TextCaseConvert(input string, mode string) Result {
	if input == "" {
		return Result{Output: ""}
	}

	switch strings.ToLower(mode) {
	case "lower":
		return Result{Output: strings.ToLower(input)}

	case "upper":
		return Result{Output: strings.ToUpper(input)}

	case "title":
		return Result{Output: toTitleCase(input)}

	case "camel":
		return Result{Output: toCamelCase(input)}

	case "snake":
		return Result{Output: toSnakeCase(input)}

	case "kebab":
		return Result{Output: toKebabCase(input)}

	default:
		return Result{Error: fmt.Sprintf("unsupported case mode: %s (supported: lower, upper, title, camel, snake, kebab)", mode)}
	}
}

// splitWords splits a string into words by whitespace, underscores, and hyphens.
func splitWords(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		return unicode.IsSpace(r) || r == '_' || r == '-'
	})
}

// toTitleCase capitalizes the first letter of each whitespace-delimited word.
func toTitleCase(s string) string {
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) > 0 {
			runes := []rune(w)
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}
	return strings.Join(words, " ")
}

// toCamelCase converts to camelCase: removes separators, capitalizes each
// word except the first.
func toCamelCase(s string) string {
	words := splitWords(s)
	if len(words) == 0 {
		return ""
	}

	var b strings.Builder
	for i, w := range words {
		if w == "" {
			continue
		}
		runes := []rune(strings.ToLower(w))
		if i == 0 {
			b.WriteString(string(runes))
		} else {
			runes[0] = unicode.ToUpper(runes[0])
			b.WriteString(string(runes))
		}
	}
	return b.String()
}

// toSnakeCase converts to snake_case: lowercase with underscores.
func toSnakeCase(s string) string {
	words := splitWords(s)
	for i, w := range words {
		words[i] = strings.ToLower(w)
	}
	return strings.Join(words, "_")
}

// toKebabCase converts to kebab-case: lowercase with hyphens.
func toKebabCase(s string) string {
	words := splitWords(s)
	for i, w := range words {
		words[i] = strings.ToLower(w)
	}
	return strings.Join(words, "-")
}
