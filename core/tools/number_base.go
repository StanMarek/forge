package tools

import (
	"fmt"
	"strconv"
	"strings"
)

// NumberBaseTool provides metadata for the Number Base converter tool.
type NumberBaseTool struct{}

func (n NumberBaseTool) Name() string        { return "Number Base Converter" }
func (n NumberBaseTool) ID() string          { return "number-base" }
func (n NumberBaseTool) Description() string { return "Convert numbers between decimal, hex, octal, and binary" }
func (n NumberBaseTool) Category() string    { return "Converters" }
func (n NumberBaseTool) Keywords() []string {
	return []string{"number", "base", "hex", "decimal", "binary", "octal", "convert"}
}

// DetectFromClipboard returns true if s has a 0x, 0b, or 0o prefix.
func (n NumberBaseTool) DetectFromClipboard(s string) bool {
	s = strings.TrimSpace(s)
	lower := strings.ToLower(s)
	return strings.HasPrefix(lower, "0x") ||
		strings.HasPrefix(lower, "0b") ||
		strings.HasPrefix(lower, "0o")
}

// NumberBaseConvert auto-detects the number base from the input prefix and
// converts it to all four representations: decimal, hex, octal, binary.
// Prefixes: 0x = hex, 0b = binary, 0o = octal, else decimal.
func NumberBaseConvert(input string) Result {
	input = strings.TrimSpace(input)
	if input == "" {
		return Result{Error: "empty input"}
	}

	var value int64
	var err error
	lower := strings.ToLower(input)

	switch {
	case strings.HasPrefix(lower, "0x"):
		value, err = strconv.ParseInt(input[2:], 16, 64)
	case strings.HasPrefix(lower, "0b"):
		value, err = strconv.ParseInt(input[2:], 2, 64)
	case strings.HasPrefix(lower, "0o"):
		value, err = strconv.ParseInt(input[2:], 8, 64)
	default:
		value, err = strconv.ParseInt(input, 10, 64)
	}

	if err != nil {
		return Result{Error: fmt.Sprintf("invalid number: %s", err.Error())}
	}

	output := fmt.Sprintf("Decimal:  %d\nHex:      %x\nOctal:    %o\nBinary:   %b",
		value, value, value, value)

	return Result{Output: output}
}
