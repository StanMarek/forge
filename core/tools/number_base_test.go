package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumberBaseTool_Metadata(t *testing.T) {
	tool := NumberBaseTool{}
	assert.Equal(t, "Number Base Converter", tool.Name())
	assert.Equal(t, "number-base", tool.ID())
	assert.Equal(t, "Converters", tool.Category())
	assert.Equal(t, "Convert numbers between decimal, hex, octal, and binary", tool.Description())
	assert.Equal(t, []string{"number", "base", "hex", "decimal", "binary", "octal", "convert"}, tool.Keywords())
}

func TestNumberBaseTool_ImplementsInterface(t *testing.T) {
	var _ Tool = NumberBaseTool{}
}

func TestNumberBaseConvert(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
		errMsg   string
	}{
		{
			name:     "empty input",
			input:    "",
			hasError: true,
			errMsg:   "empty input",
		},
		{
			name:  "decimal 255",
			input: "255",
			expected: "Decimal:  255\n" +
				"Hex:      ff\n" +
				"Octal:    377\n" +
				"Binary:   11111111",
		},
		{
			name:  "hex 0xff",
			input: "0xff",
			expected: "Decimal:  255\n" +
				"Hex:      ff\n" +
				"Octal:    377\n" +
				"Binary:   11111111",
		},
		{
			name:  "hex uppercase 0xFF",
			input: "0xFF",
			expected: "Decimal:  255\n" +
				"Hex:      ff\n" +
				"Octal:    377\n" +
				"Binary:   11111111",
		},
		{
			name:  "binary 0b11111111",
			input: "0b11111111",
			expected: "Decimal:  255\n" +
				"Hex:      ff\n" +
				"Octal:    377\n" +
				"Binary:   11111111",
		},
		{
			name:  "octal 0o377",
			input: "0o377",
			expected: "Decimal:  255\n" +
				"Hex:      ff\n" +
				"Octal:    377\n" +
				"Binary:   11111111",
		},
		{
			name:  "zero",
			input: "0",
			expected: "Decimal:  0\n" +
				"Hex:      0\n" +
				"Octal:    0\n" +
				"Binary:   0",
		},
		{
			name:  "decimal 42",
			input: "42",
			expected: "Decimal:  42\n" +
				"Hex:      2a\n" +
				"Octal:    52\n" +
				"Binary:   101010",
		},
		{
			name:  "hex 0x0",
			input: "0x0",
			expected: "Decimal:  0\n" +
				"Hex:      0\n" +
				"Octal:    0\n" +
				"Binary:   0",
		},
		{
			name:     "invalid decimal",
			input:    "abc",
			hasError: true,
			errMsg:   "invalid number",
		},
		{
			name:     "invalid hex",
			input:    "0xZZZ",
			hasError: true,
			errMsg:   "invalid number",
		},
		{
			name:     "invalid binary",
			input:    "0b1234",
			hasError: true,
			errMsg:   "invalid number",
		},
		{
			name:     "invalid octal",
			input:    "0o999",
			hasError: true,
			errMsg:   "invalid number",
		},
		{
			name:  "large number",
			input: "65535",
			expected: "Decimal:  65535\n" +
				"Hex:      ffff\n" +
				"Octal:    177777\n" +
				"Binary:   1111111111111111",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NumberBaseConvert(tt.input)
			if tt.hasError {
				assert.NotEmpty(t, result.Error)
				assert.Contains(t, result.Error, tt.errMsg)
			} else {
				assert.Equal(t, "", result.Error)
				assert.Equal(t, tt.expected, result.Output)
			}
		})
	}
}

func TestNumberBaseTool_DetectFromClipboard(t *testing.T) {
	tool := NumberBaseTool{}

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"hex prefix", "0xff", true},
		{"hex uppercase", "0XFF", true},
		{"binary prefix", "0b1010", true},
		{"binary uppercase", "0B1010", true},
		{"octal prefix", "0o77", true},
		{"octal uppercase", "0O77", true},
		{"decimal no prefix", "255", false},
		{"empty", "", false},
		{"plain text", "hello", false},
		{"with whitespace", "  0xff  ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tool.DetectFromClipboard(tt.input))
		})
	}
}
