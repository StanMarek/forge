package tools

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimestampTool_Metadata(t *testing.T) {
	tool := TimestampTool{}
	assert.Equal(t, "Timestamp Converter", tool.Name())
	assert.Equal(t, "timestamp", tool.ID())
	assert.Equal(t, "Converters", tool.Category())
	assert.Equal(t, "Convert between Unix timestamps and human-readable dates", tool.Description())
	assert.Equal(t, []string{"timestamp", "unix", "date", "time", "epoch"}, tool.Keywords())
}

func TestTimestampTool_ImplementsInterface(t *testing.T) {
	var _ Tool = TimestampTool{}
}

func TestTimestampFromUnix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		tz       string
		contains []string
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
			name:     "known timestamp seconds",
			input:    "1700000000",
			tz:       "UTC",
			contains: []string{"2023-11-14T22:13:20Z", "1700000000", "1700000000000"},
		},
		{
			name:     "known timestamp millis",
			input:    "1700000000000",
			tz:       "UTC",
			contains: []string{"2023-11-14T22:13:20Z", "1700000000", "1700000000000"},
		},
		{
			name:     "default timezone is UTC",
			input:    "1700000000",
			tz:       "",
			contains: []string{"2023-11-14T22:13:20Z"},
		},
		{
			name:     "timezone America/New_York",
			input:    "1700000000",
			tz:       "America/New_York",
			contains: []string{"2023-11-14T17:13:20-05:00"},
		},
		{
			name:     "invalid timezone",
			input:    "1700000000",
			tz:       "Invalid/Zone",
			hasError: true,
			errMsg:   "invalid timezone",
		},
		{
			name:     "invalid input",
			input:    "not-a-number",
			hasError: true,
			errMsg:   "invalid unix timestamp",
		},
		{
			name:     "epoch zero",
			input:    "0000000000",
			tz:       "UTC",
			contains: []string{"1970-01-01T00:00:00Z"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TimestampFromUnix(tt.input, tt.tz)
			if tt.hasError {
				assert.NotEmpty(t, result.Error)
				assert.Contains(t, result.Error, tt.errMsg)
			} else {
				assert.Equal(t, "", result.Error)
				for _, substr := range tt.contains {
					assert.Contains(t, result.Output, substr)
				}
			}
		})
	}
}

func TestTimestampToUnix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		millis   bool
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
			name:     "RFC3339",
			input:    "2023-11-14T22:13:20Z",
			expected: "1700000000",
		},
		{
			name:     "RFC3339 with millis flag",
			input:    "2023-11-14T22:13:20Z",
			millis:   true,
			expected: "1700000000000",
		},
		{
			name:     "date only",
			input:    "2023-11-14",
			expected: "1699920000",
		},
		{
			name:     "datetime space separator",
			input:    "2023-11-14 22:13:20",
			expected: "1700000000",
		},
		{
			name:     "invalid datetime",
			input:    "not-a-date",
			hasError: true,
			errMsg:   "unable to parse datetime",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TimestampToUnix(tt.input, tt.millis)
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

func TestTimestampNow(t *testing.T) {
	tests := []struct {
		name     string
		tz       string
		hasError bool
		errMsg   string
	}{
		{
			name: "UTC",
			tz:   "UTC",
		},
		{
			name: "default empty tz",
			tz:   "",
		},
		{
			name: "America/New_York",
			tz:   "America/New_York",
		},
		{
			name:     "invalid timezone",
			tz:       "Invalid/Zone",
			hasError: true,
			errMsg:   "invalid timezone",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TimestampNow(tt.tz)
			if tt.hasError {
				assert.NotEmpty(t, result.Error)
				assert.Contains(t, result.Error, tt.errMsg)
			} else {
				assert.Equal(t, "", result.Error)
				assert.NotEmpty(t, result.Output)
				assert.True(t, strings.Contains(result.Output, "Unix:"))
				assert.True(t, strings.Contains(result.Output, "Unix ms:"))
				assert.True(t, strings.Contains(result.Output, "RFC3339:"))
				assert.True(t, strings.Contains(result.Output, "Human:"))
			}
		})
	}
}

func TestTimestampTool_DetectFromClipboard(t *testing.T) {
	tool := TimestampTool{}

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"10 digit unix seconds", "1700000000", true},
		{"13 digit unix millis", "1700000000000", true},
		{"too short", "170000000", false},
		{"too long", "17000000000000", false},
		{"11 digits", "17000000001", false},
		{"12 digits", "170000000012", false},
		{"empty", "", false},
		{"letters", "abcdefghij", false},
		{"mixed", "170000000a", false},
		{"with whitespace", "  1700000000  ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tool.DetectFromClipboard(tt.input))
		})
	}
}
