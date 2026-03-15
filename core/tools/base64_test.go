package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64Tool_Metadata(t *testing.T) {
	tool := Base64Tool{}
	assert.Equal(t, "Base64 Encode / Decode", tool.Name())
	assert.Equal(t, "base64", tool.ID())
	assert.Equal(t, "Encoders", tool.Category())
	assert.Equal(t, "Encode and decode Base64 strings", tool.Description())
	assert.Equal(t, []string{"base64", "encode", "decode", "b64"}, tool.Keywords())
}

func TestBase64Tool_ImplementsInterface(t *testing.T) {
	var _ Tool = Base64Tool{}
}

func TestBase64Encode(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		urlSafe   bool
		noPadding bool
		expected  string
	}{
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "hello world",
			input:    "hello",
			expected: "aGVsbG8=",
		},
		{
			name:     "hello world full",
			input:    "hello world",
			expected: "aGVsbG8gd29ybGQ=",
		},
		{
			name:     "no padding needed",
			input:    "abc",
			expected: "YWJj",
		},
		{
			name:      "no padding flag",
			input:     "hello",
			noPadding: true,
			expected:  "aGVsbG8",
		},
		{
			name:     "url safe encoding",
			input:    "subjects?_d",
			urlSafe:  true,
			expected: "c3ViamVjdHM_X2Q=",
		},
		{
			name:     "standard encoding with special chars",
			input:    "subjects?_d",
			urlSafe:  false,
			expected: "c3ViamVjdHM/X2Q=",
		},
		{
			name:      "url safe no padding",
			input:     "subjects?_d",
			urlSafe:   true,
			noPadding: true,
			expected:  "c3ViamVjdHM_X2Q",
		},
		{
			name:     "binary-like content",
			input:    "\x00\x01\x02\x03",
			expected: "AAECAw==",
		},
		{
			name:     "unicode content",
			input:    "Zdrowie!",
			expected: "WmRyb3dpZSE=",
		},
		{
			name:     "long string",
			input:    "The quick brown fox jumps over the lazy dog",
			expected: "VGhlIHF1aWNrIGJyb3duIGZveCBqdW1wcyBvdmVyIHRoZSBsYXp5IGRvZw==",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Base64Encode(tt.input, tt.urlSafe, tt.noPadding)
			assert.Equal(t, "", result.Error)
			assert.Equal(t, tt.expected, result.Output)
		})
	}
}

func TestBase64Decode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		urlSafe  bool
		expected string
		hasError bool
		errMsg   string
	}{
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "hello",
			input:    "aGVsbG8=",
			expected: "hello",
		},
		{
			name:     "hello world",
			input:    "aGVsbG8gd29ybGQ=",
			expected: "hello world",
		},
		{
			name:     "no padding in input",
			input:    "aGVsbG8",
			expected: "hello",
		},
		{
			name:     "double padding",
			input:    "AAECAw==",
			expected: "\x00\x01\x02\x03",
		},
		{
			name:     "url safe decode",
			input:    "c3ViamVjdHM_X2Q=",
			urlSafe:  true,
			expected: "subjects?_d",
		},
		{
			name:     "url safe decode without padding",
			input:    "c3ViamVjdHM_X2Q",
			urlSafe:  true,
			expected: "subjects?_d",
		},
		{
			name:     "standard decode",
			input:    "c3ViamVjdHM/X2Q=",
			urlSafe:  false,
			expected: "subjects?_d",
		},
		{
			name:     "long string",
			input:    "VGhlIHF1aWNrIGJyb3duIGZveCBqdW1wcyBvdmVyIHRoZSBsYXp5IGRvZw==",
			expected: "The quick brown fox jumps over the lazy dog",
		},
		{
			name:     "invalid base64",
			input:    "!!!invalid!!!",
			hasError: true,
			errMsg:   "invalid base64:",
		},
		{
			name:     "invalid characters mixed in",
			input:    "aGVs bG8=",
			hasError: true,
			errMsg:   "invalid base64:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Base64Decode(tt.input, tt.urlSafe)
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

func TestBase64_RoundTrip(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		urlSafe   bool
		noPadding bool
	}{
		{"simple", "hello world", false, false},
		{"url safe", "subjects?_d", true, false},
		{"no padding", "hello", false, true},
		{"url safe no padding", "subjects?_d", true, true},
		{"empty", "", false, false},
		{"unicode", "Zdrowie!", false, false},
		{"binary", "\x00\xff\xfe", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := Base64Encode(tt.input, tt.urlSafe, tt.noPadding)
			assert.Equal(t, "", encoded.Error)

			decoded := Base64Decode(encoded.Output, tt.urlSafe)
			assert.Equal(t, "", decoded.Error)
			assert.Equal(t, tt.input, decoded.Output)
		})
	}
}

func TestBase64Tool_DetectFromClipboard(t *testing.T) {
	tool := Base64Tool{}

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid base64 single pad", "aGVsbG8=", true},
		{"valid base64 double pad", "AAECAw==", true},
		{"valid base64 no pad needed", "YWJj", true},
		{"long valid base64", "VGhlIHF1aWNrIGJyb3duIGZveCBqdW1wcyBvdmVyIHRoZSBsYXp5IGRvZw==", true},
		{"too short", "YQ=", false},
		{"single char", "a", false},
		{"empty", "", false},
		{"not divisible by 4", "aGVsbG8", false},           // len 7, 7%4!=0
		{"contains spaces", "aGVs bG8=", false},
		{"contains invalid chars", "aGV!bG8=", false},
		{"url safe chars not standard", "c3ViamVjdHM_X2Q=", false}, // _ is not in std regex
		{"with whitespace trimmed", "  YWJj  ", true},
		{"only padding", "====", true},                     // matches regex, len 4, div by 4
		{"numbers and letters", "QmFz", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tool.DetectFromClipboard(tt.input))
		})
	}
}
