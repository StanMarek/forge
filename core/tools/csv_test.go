package tools

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------- Tool metadata ----------

func TestCSVTool_Metadata(t *testing.T) {
	tool := CSVTool{}

	assert.Equal(t, "JSON to CSV", tool.Name())
	assert.Equal(t, "csv", tool.ID())
	assert.Equal(t, "Converters", tool.Category())
	assert.NotEmpty(t, tool.Description())

	keywords := tool.Keywords()
	for _, kw := range []string{"csv", "json", "table", "tsv", "convert"} {
		assert.Contains(t, keywords, kw)
	}
}

func TestCSVTool_DetectFromClipboard(t *testing.T) {
	tool := CSVTool{}
	assert.False(t, tool.DetectFromClipboard("anything"))
	assert.False(t, tool.DetectFromClipboard(""))
}

// ---------- JSONToCSV ----------

func TestJSONToCSV_Simple(t *testing.T) {
	input := `[{"name":"Alice","age":30},{"name":"Bob","age":25}]`
	r := JSONToCSV(input, ",")
	require.Empty(t, r.Error)

	lines := strings.Split(r.Output, "\n")
	require.Len(t, lines, 3)
	assert.Equal(t, "age,name", lines[0]) // sorted alphabetically
	assert.Contains(t, r.Output, "Alice")
	assert.Contains(t, r.Output, "Bob")
	assert.Contains(t, r.Output, "30")
	assert.Contains(t, r.Output, "25")
}

func TestJSONToCSV_EmptyArray(t *testing.T) {
	r := JSONToCSV("[]", ",")
	require.Empty(t, r.Error)
	assert.Equal(t, "", r.Output)
}

func TestJSONToCSV_NonArray(t *testing.T) {
	r := JSONToCSV(`{"key":"value"}`, ",")
	assert.NotEmpty(t, r.Error)
	assert.Contains(t, r.Error, "JSON array")
}

func TestJSONToCSV_EmptyInput(t *testing.T) {
	r := JSONToCSV("", ",")
	assert.NotEmpty(t, r.Error)
	assert.Equal(t, "empty input", r.Error)
}

func TestJSONToCSV_NestedValues(t *testing.T) {
	input := `[{"name":"Alice","meta":{"role":"admin"}}]`
	r := JSONToCSV(input, ",")
	require.Empty(t, r.Error)
	// Nested value should be stringified as JSON.
	assert.Contains(t, r.Output, `"{""role"":""admin""}"`)
}

func TestJSONToCSV_MissingKeys(t *testing.T) {
	input := `[{"a":"1","b":"2"},{"a":"3"}]`
	r := JSONToCSV(input, ",")
	require.Empty(t, r.Error)

	lines := strings.Split(r.Output, "\n")
	require.Len(t, lines, 3)
	assert.Equal(t, "a,b", lines[0])
	// Second row should have empty value for missing key "b".
	assert.Equal(t, "3,", lines[2])
}

func TestJSONToCSV_TSVDelimiter(t *testing.T) {
	input := `[{"x":"1","y":"2"}]`
	r := JSONToCSV(input, "\t")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "x\ty")
	assert.Contains(t, r.Output, "1\t2")
}

func TestJSONToCSV_BooleanAndNull(t *testing.T) {
	input := `[{"flag":true,"empty":null,"num":42}]`
	r := JSONToCSV(input, ",")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "true")
	assert.Contains(t, r.Output, "42")
}

// ---------- CSVToJSON ----------

func TestCSVToJSON_Simple(t *testing.T) {
	input := "name,age\nAlice,30\nBob,25"
	r := CSVToJSON(input, ",")
	require.Empty(t, r.Error)

	var result []map[string]interface{}
	err := json.Unmarshal([]byte(r.Output), &result)
	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "Alice", result[0]["name"])
	assert.Equal(t, "30", result[0]["age"])
	assert.Equal(t, "Bob", result[1]["name"])
}

func TestCSVToJSON_HeaderOnly(t *testing.T) {
	r := CSVToJSON("name,age", ",")
	require.Empty(t, r.Error)
	assert.Equal(t, "[]", r.Output)
}

func TestCSVToJSON_EmptyInput(t *testing.T) {
	r := CSVToJSON("", ",")
	assert.NotEmpty(t, r.Error)
	assert.Equal(t, "empty input", r.Error)
}

func TestCSVToJSON_TSVDelimiter(t *testing.T) {
	input := "x\ty\n1\t2"
	r := CSVToJSON(input, "\t")
	require.Empty(t, r.Error)

	var result []map[string]interface{}
	err := json.Unmarshal([]byte(r.Output), &result)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "1", result[0]["x"])
	assert.Equal(t, "2", result[0]["y"])
}

// ---------- Roundtrip ----------

func TestCSV_Roundtrip(t *testing.T) {
	original := `[{"city":"New York","country":"US"},{"city":"London","country":"UK"}]`

	// JSON -> CSV
	csvResult := JSONToCSV(original, ",")
	require.Empty(t, csvResult.Error)

	// CSV -> JSON
	jsonResult := CSVToJSON(csvResult.Output, ",")
	require.Empty(t, jsonResult.Error)

	var result []map[string]interface{}
	err := json.Unmarshal([]byte(jsonResult.Output), &result)
	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "New York", result[0]["city"])
	assert.Equal(t, "US", result[0]["country"])
	assert.Equal(t, "London", result[1]["city"])
	assert.Equal(t, "UK", result[1]["country"])
}

// ---------- Tool interface compliance ----------

func TestCSVTool_ImplementsToolInterface(t *testing.T) {
	var _ Tool = CSVTool{}
}
