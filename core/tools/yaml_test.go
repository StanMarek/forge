package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYAMLTool_Metadata(t *testing.T) {
	tool := YAMLTool{}
	assert.Equal(t, "JSON / YAML Converter", tool.Name())
	assert.Equal(t, "yaml", tool.ID())
	assert.Equal(t, "Converters", tool.Category())
	assert.Equal(t, "Convert between JSON and YAML formats", tool.Description())
	assert.Equal(t, []string{"yaml", "json", "convert"}, tool.Keywords())
}

func TestYAMLTool_ImplementsInterface(t *testing.T) {
	var _ Tool = YAMLTool{}
}

func TestYAMLToJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		compact  bool
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
			name:     "whitespace only",
			input:    "   ",
			hasError: true,
			errMsg:   "empty input",
		},
		{
			name:     "simple key value",
			input:    "name: forge",
			compact:  true,
			expected: `{"name":"forge"}`,
		},
		{
			name:  "simple key value pretty",
			input: "name: forge",
			expected: `{
  "name": "forge"
}`,
		},
		{
			name: "nested object",
			input: `server:
  host: localhost
  port: 8080`,
			compact:  true,
			expected: `{"server":{"host":"localhost","port":8080}}`,
		},
		{
			name: "array",
			input: `fruits:
  - apple
  - banana
  - cherry`,
			compact:  true,
			expected: `{"fruits":["apple","banana","cherry"]}`,
		},
		{
			name: "nested with arrays",
			input: `database:
  hosts:
    - db1.example.com
    - db2.example.com
  port: 5432
  name: mydb`,
			compact:  true,
			expected: `{"database":{"hosts":["db1.example.com","db2.example.com"],"name":"mydb","port":5432}}`,
		},
		{
			name:     "invalid yaml",
			input:    ":\n  :\n    - :\n  invalid:: yaml::",
			hasError: true,
			errMsg:   "invalid YAML",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := YAMLToJSON(tt.input, tt.compact)
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

func TestJSONToYAML(t *testing.T) {
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
			name:     "whitespace only",
			input:    "   ",
			hasError: true,
			errMsg:   "empty input",
		},
		{
			name:     "simple object",
			input:    `{"name":"forge"}`,
			expected: "name: forge",
		},
		{
			name:     "nested object",
			input:    `{"server":{"host":"localhost","port":8080}}`,
			expected: "server:\n    host: localhost\n    port: 8080",
		},
		{
			name:     "array",
			input:    `{"fruits":["apple","banana","cherry"]}`,
			expected: "fruits:\n    - apple\n    - banana\n    - cherry",
		},
		{
			name:     "invalid json",
			input:    `{invalid json}`,
			hasError: true,
			errMsg:   "invalid JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JSONToYAML(tt.input)
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

func TestYAML_RoundTrip(t *testing.T) {
	tests := []struct {
		name string
		yaml string
	}{
		{
			name: "simple",
			yaml: "name: forge",
		},
		{
			name: "nested",
			yaml: "server:\n    host: localhost\n    port: 8080",
		},
		{
			name: "with array",
			yaml: "items:\n    - one\n    - two\n    - three",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonResult := YAMLToJSON(tt.yaml, false)
			assert.Equal(t, "", jsonResult.Error)

			yamlResult := JSONToYAML(jsonResult.Output)
			assert.Equal(t, "", yamlResult.Error)
			assert.Equal(t, tt.yaml, yamlResult.Output)
		})
	}
}

func TestYAMLTool_DetectFromClipboard(t *testing.T) {
	tool := YAMLTool{}

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"empty", "", false},
		{"json object", `{"key": "value"}`, false},
		{"json array", `["a", "b"]`, false},
		{"valid yaml mapping", "name: forge\nversion: 1.0", true},
		{"valid yaml list", "- item1\n- item2", true},
		{"plain scalar", "hello", false},
		{"number", "42", false},
		{"yaml with nested", "server:\n  host: localhost", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tool.DetectFromClipboard(tt.input))
		})
	}
}
