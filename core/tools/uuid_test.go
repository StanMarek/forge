package tools

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Standard UUID regex for validation.
var testUUIDRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
var testUUIDNoHyphensRegex = regexp.MustCompile(`^[0-9a-f]{32}$`)
var testUUIDUpperRegex = regexp.MustCompile(`^[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12}$`)
var testUUIDUpperNoHyphensRegex = regexp.MustCompile(`^[0-9A-F]{32}$`)

// --- UUIDTool metadata tests ---

func TestUUIDToolMetadata(t *testing.T) {
	tool := UUIDTool{}
	assert.Equal(t, "UUID Generate / Validate / Parse", tool.Name())
	assert.Equal(t, "uuid", tool.ID())
	assert.Equal(t, "Generate, validate, and parse UUIDs", tool.Description())
	assert.Equal(t, "Generators", tool.Category())
	assert.Contains(t, tool.Keywords(), "uuid")
	assert.Contains(t, tool.Keywords(), "guid")
	assert.Contains(t, tool.Keywords(), "v4")
	assert.Contains(t, tool.Keywords(), "v7")
}

// --- DetectFromClipboard tests ---

func TestUUIDDetectFromClipboard(t *testing.T) {
	tool := UUIDTool{}

	tests := []struct {
		name   string
		input  string
		expect bool
	}{
		{"valid UUID", "550e8400-e29b-41d4-a716-446655440000", true},
		{"valid UUID with spaces", "  550e8400-e29b-41d4-a716-446655440000  ", true},
		{"uppercase UUID", "550E8400-E29B-41D4-A716-446655440000", true},
		{"no hyphens", "550e8400e29b41d4a716446655440000", false},
		{"too short", "550e8400-e29b", false},
		{"random string", "hello world", false},
		{"empty string", "", false},
		{"almost UUID missing char", "550e8400-e29b-41d4-a716-44665544000", false},
		{"UUID with extra char", "550e8400-e29b-41d4-a716-4466554400000", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tool.DetectFromClipboard(tt.input))
		})
	}
}

// --- UUIDGenerate tests ---

func TestUUIDGenerateV4(t *testing.T) {
	result := UUIDGenerate(4, false, false)
	require.Empty(t, result.Error)
	assert.True(t, testUUIDRegex.MatchString(result.Output), "expected valid lowercase UUID with hyphens, got: %s", result.Output)
	// Version nibble should be '4'.
	assert.Equal(t, "4", string(result.Output[14]))
}

func TestUUIDGenerateV7(t *testing.T) {
	result := UUIDGenerate(7, false, false)
	require.Empty(t, result.Error)
	assert.True(t, testUUIDRegex.MatchString(result.Output), "expected valid lowercase UUID with hyphens, got: %s", result.Output)
	// Version nibble should be '7'.
	assert.Equal(t, "7", string(result.Output[14]))
}

func TestUUIDGenerateUppercase(t *testing.T) {
	result := UUIDGenerate(4, true, false)
	require.Empty(t, result.Error)
	assert.True(t, testUUIDUpperRegex.MatchString(result.Output), "expected uppercase UUID, got: %s", result.Output)
}

func TestUUIDGenerateNoHyphens(t *testing.T) {
	result := UUIDGenerate(4, false, true)
	require.Empty(t, result.Error)
	assert.True(t, testUUIDNoHyphensRegex.MatchString(result.Output), "expected UUID without hyphens, got: %s", result.Output)
	assert.Len(t, result.Output, 32)
}

func TestUUIDGenerateUppercaseNoHyphens(t *testing.T) {
	result := UUIDGenerate(4, true, true)
	require.Empty(t, result.Error)
	assert.True(t, testUUIDUpperNoHyphensRegex.MatchString(result.Output), "expected uppercase UUID without hyphens, got: %s", result.Output)
	assert.Len(t, result.Output, 32)
}

func TestUUIDGenerateUnsupportedVersion(t *testing.T) {
	tests := []int{0, 1, 2, 3, 5, 6, 8, 99}
	for _, v := range tests {
		t.Run("version_"+strings.Repeat("x", v), func(t *testing.T) {
			result := UUIDGenerate(v, false, false)
			assert.NotEmpty(t, result.Error)
			assert.Contains(t, result.Error, "unsupported UUID version")
			assert.Contains(t, result.Error, "supported: 4, 7")
			assert.Empty(t, result.Output)
		})
	}
}

func TestUUIDGenerateUniqueness(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		result := UUIDGenerate(4, false, false)
		require.Empty(t, result.Error)
		assert.False(t, seen[result.Output], "duplicate UUID generated: %s", result.Output)
		seen[result.Output] = true
	}
}

// --- UUIDValidate tests ---

func TestUUIDValidateValid(t *testing.T) {
	result := UUIDValidate("550e8400-e29b-41d4-a716-446655440000")
	require.Empty(t, result.Error)
	assert.Equal(t, "valid (version 4)", result.Output)
}

func TestUUIDValidateV7(t *testing.T) {
	// Generate a v7 UUID and validate it.
	gen := UUIDGenerate(7, false, false)
	require.Empty(t, gen.Error)
	result := UUIDValidate(gen.Output)
	require.Empty(t, result.Error)
	assert.Equal(t, "valid (version 7)", result.Output)
}

func TestUUIDValidateWithWhitespace(t *testing.T) {
	result := UUIDValidate("  550e8400-e29b-41d4-a716-446655440000  ")
	require.Empty(t, result.Error)
	assert.Equal(t, "valid (version 4)", result.Output)
}

func TestUUIDValidateInvalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"not a UUID", "hello-world"},
		{"too short", "550e8400-e29b-41d4"},
		{"invalid characters", "zzzzzzzz-zzzz-zzzz-zzzz-zzzzzzzzzzzz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UUIDValidate(tt.input)
			assert.NotEmpty(t, result.Error)
			assert.Contains(t, result.Error, "invalid UUID")
			assert.Empty(t, result.Output)
		})
	}
}

// --- UUIDParse tests ---

func TestUUIDParseV4(t *testing.T) {
	result := UUIDParse("550e8400-e29b-41d4-a716-446655440000")
	require.Empty(t, result.Error)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", result.UUID)
	assert.Equal(t, 4, result.Version)
	assert.Equal(t, "RFC 4122", result.Variant)
	assert.Empty(t, result.Timestamp, "v4 UUID should have no timestamp")
	assert.Contains(t, result.Output, "Version: 4")
	assert.Contains(t, result.Output, "Variant: RFC 4122")
	assert.NotContains(t, result.Output, "Time:")
}

func TestUUIDParseV7WithTimestamp(t *testing.T) {
	// Generate a v7 UUID and parse it.
	gen := UUIDGenerate(7, false, false)
	require.Empty(t, gen.Error)

	result := UUIDParse(gen.Output)
	require.Empty(t, result.Error)
	assert.Equal(t, gen.Output, result.UUID)
	assert.Equal(t, 7, result.Version)
	assert.Equal(t, "RFC 4122", result.Variant)
	assert.NotEmpty(t, result.Timestamp, "v7 UUID should have a timestamp")
	assert.Contains(t, result.Output, "Time:")
}

func TestUUIDParseWithWhitespace(t *testing.T) {
	result := UUIDParse("  550e8400-e29b-41d4-a716-446655440000  ")
	require.Empty(t, result.Error)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", result.UUID)
}

func TestUUIDParseInvalid(t *testing.T) {
	result := UUIDParse("not-a-uuid")
	assert.NotEmpty(t, result.Error)
	assert.Contains(t, result.Error, "invalid UUID")
	assert.Empty(t, result.UUID)
}

func TestUUIDParseOutputFormat(t *testing.T) {
	result := UUIDParse("550e8400-e29b-41d4-a716-446655440000")
	require.Empty(t, result.Error)
	lines := strings.Split(result.Output, "\n")
	assert.GreaterOrEqual(t, len(lines), 3)
	assert.True(t, strings.HasPrefix(lines[0], "UUID:"))
	assert.True(t, strings.HasPrefix(lines[1], "Version:"))
	assert.True(t, strings.HasPrefix(lines[2], "Variant:"))
}

// --- variantString tests ---

func TestVariantStringMapping(t *testing.T) {
	// The test UUID 550e8400-e29b-41d4-a716-446655440000 is RFC 4122 variant.
	// We test the helper indirectly through UUIDParse.
	result := UUIDParse("550e8400-e29b-41d4-a716-446655440000")
	require.Empty(t, result.Error)
	assert.Equal(t, "RFC 4122", result.Variant)
}
