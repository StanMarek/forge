package tools

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// JSONTool metadata
// ---------------------------------------------------------------------------

func TestJSONToolMetadata(t *testing.T) {
	tool := JSONTool{}
	assert.Equal(t, "JSON Formatter", tool.Name())
	assert.Equal(t, "json", tool.ID())
	assert.Equal(t, "Formatters", tool.Category())
	assert.Contains(t, tool.Keywords(), "json")
	assert.Contains(t, tool.Keywords(), "format")
	assert.Contains(t, tool.Keywords(), "minify")
	assert.Contains(t, tool.Keywords(), "validate")
	assert.Contains(t, tool.Keywords(), "pretty")
}

// ---------------------------------------------------------------------------
// DetectFromClipboard
// ---------------------------------------------------------------------------

func TestJSONDetectFromClipboard(t *testing.T) {
	tool := JSONTool{}

	assert.True(t, tool.DetectFromClipboard(`{"key":"value"}`))
	assert.True(t, tool.DetectFromClipboard(`[1,2,3]`))
	assert.True(t, tool.DetectFromClipboard(`  {"a":1}  `)) // with whitespace
	assert.True(t, tool.DetectFromClipboard(`"hello"`))      // valid JSON string
	assert.True(t, tool.DetectFromClipboard(`42`))            // valid JSON number

	assert.False(t, tool.DetectFromClipboard(""))
	assert.False(t, tool.DetectFromClipboard("   "))
	assert.False(t, tool.DetectFromClipboard("{invalid}"))
	assert.False(t, tool.DetectFromClipboard("not json at all"))
}

// ---------------------------------------------------------------------------
// JSONFormat
// ---------------------------------------------------------------------------

func TestJSONFormat_Basic(t *testing.T) {
	input := `{"name":"Alice","age":30}`
	r := JSONFormat(input, 2, false, false)
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "\"name\": \"Alice\"")
	assert.Contains(t, r.Output, "\"age\": 30")
}

func TestJSONFormat_PreservesKeyOrder(t *testing.T) {
	// Without sortKeys, json.Indent preserves original key order.
	input := `{"zebra":1,"apple":2,"mango":3}`
	r := JSONFormat(input, 2, false, false)
	require.Empty(t, r.Error)

	zebraIdx := strings.Index(r.Output, "\"zebra\"")
	appleIdx := strings.Index(r.Output, "\"apple\"")
	mangoIdx := strings.Index(r.Output, "\"mango\"")
	assert.Less(t, zebraIdx, appleIdx, "zebra should appear before apple (original order)")
	assert.Less(t, appleIdx, mangoIdx, "apple should appear before mango (original order)")
}

func TestJSONFormat_SortKeys(t *testing.T) {
	input := `{"zebra":1,"apple":2,"mango":3}`
	r := JSONFormat(input, 2, true, false)
	require.Empty(t, r.Error)

	appleIdx := strings.Index(r.Output, "\"apple\"")
	mangoIdx := strings.Index(r.Output, "\"mango\"")
	zebraIdx := strings.Index(r.Output, "\"zebra\"")
	assert.Less(t, appleIdx, mangoIdx, "apple should appear before mango (sorted)")
	assert.Less(t, mangoIdx, zebraIdx, "mango should appear before zebra (sorted)")
}

func TestJSONFormat_SortKeysRecursive(t *testing.T) {
	input := `{"z":{"beta":1,"alpha":2},"a":{"delta":3,"charlie":4}}`
	r := JSONFormat(input, 2, true, false)
	require.Empty(t, r.Error)

	// Top-level: "a" before "z"
	aIdx := strings.Index(r.Output, "\"a\"")
	zIdx := strings.Index(r.Output, "\"z\"")
	assert.Less(t, aIdx, zIdx, "top-level 'a' should come before 'z'")

	// Nested under "a": "charlie" before "delta"
	charlieIdx := strings.Index(r.Output, "\"charlie\"")
	deltaIdx := strings.Index(r.Output, "\"delta\"")
	assert.Less(t, charlieIdx, deltaIdx, "nested 'charlie' should come before 'delta'")

	// Nested under "z": "alpha" before "beta"
	alphaIdx := strings.Index(r.Output, "\"alpha\"")
	betaIdx := strings.Index(r.Output, "\"beta\"")
	assert.Less(t, alphaIdx, betaIdx, "nested 'alpha' should come before 'beta'")
}

func TestJSONFormat_UseTabs(t *testing.T) {
	input := `{"key":"value"}`
	r := JSONFormat(input, 2, false, true)
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "\t\"key\"")
}

func TestJSONFormat_CustomIndent(t *testing.T) {
	input := `{"key":"value"}`
	r := JSONFormat(input, 4, false, false)
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "    \"key\"")
}

func TestJSONFormat_EmptyInput(t *testing.T) {
	r := JSONFormat("", 2, false, false)
	assert.Equal(t, "empty input", r.Error)

	r = JSONFormat("   ", 2, false, false)
	assert.Equal(t, "empty input", r.Error)
}

func TestJSONFormat_InvalidJSON(t *testing.T) {
	r := JSONFormat("{bad}", 2, false, false)
	assert.NotEmpty(t, r.Error)
	assert.Contains(t, r.Error, "invalid JSON")
}

func TestJSONFormat_WhitespaceTrimed(t *testing.T) {
	input := `  {"key":"value"}  `
	r := JSONFormat(input, 2, false, false)
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "\"key\": \"value\"")
}

func TestJSONFormat_Array(t *testing.T) {
	input := `[1,2,3]`
	r := JSONFormat(input, 2, false, false)
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "1")
}

func TestJSONFormat_SortKeysInvalidJSON(t *testing.T) {
	r := JSONFormat("{bad}", 2, true, false)
	assert.NotEmpty(t, r.Error)
	assert.Contains(t, r.Error, "invalid JSON")
}

// ---------------------------------------------------------------------------
// JSONMinify
// ---------------------------------------------------------------------------

func TestJSONMinify_Basic(t *testing.T) {
	input := `{
  "name": "Alice",
  "age": 30
}`
	r := JSONMinify(input)
	require.Empty(t, r.Error)
	assert.Equal(t, `{"name":"Alice","age":30}`, r.Output)
}

func TestJSONMinify_AlreadyCompact(t *testing.T) {
	input := `{"a":1}`
	r := JSONMinify(input)
	require.Empty(t, r.Error)
	assert.Equal(t, `{"a":1}`, r.Output)
}

func TestJSONMinify_EmptyInput(t *testing.T) {
	r := JSONMinify("")
	assert.Equal(t, "empty input", r.Error)

	r = JSONMinify("   ")
	assert.Equal(t, "empty input", r.Error)
}

func TestJSONMinify_InvalidJSON(t *testing.T) {
	r := JSONMinify("{bad}")
	assert.NotEmpty(t, r.Error)
	assert.Contains(t, r.Error, "invalid JSON")
}

func TestJSONMinify_WhitespaceTrimmed(t *testing.T) {
	input := `  { "key" : "value" }  `
	r := JSONMinify(input)
	require.Empty(t, r.Error)
	assert.Equal(t, `{"key":"value"}`, r.Output)
}

func TestJSONMinify_Array(t *testing.T) {
	input := `[ 1 , 2 , 3 ]`
	r := JSONMinify(input)
	require.Empty(t, r.Error)
	assert.Equal(t, `[1,2,3]`, r.Output)
}

// ---------------------------------------------------------------------------
// JSONValidate
// ---------------------------------------------------------------------------

func TestJSONValidate_ValidObject(t *testing.T) {
	r := JSONValidate(`{"key":"value"}`)
	require.Empty(t, r.Error)
	assert.Equal(t, "valid", r.Output)
}

func TestJSONValidate_ValidArray(t *testing.T) {
	r := JSONValidate(`[1,2,3]`)
	require.Empty(t, r.Error)
	assert.Equal(t, "valid", r.Output)
}

func TestJSONValidate_ValidScalar(t *testing.T) {
	r := JSONValidate(`"hello"`)
	require.Empty(t, r.Error)
	assert.Equal(t, "valid", r.Output)

	r = JSONValidate(`42`)
	require.Empty(t, r.Error)
	assert.Equal(t, "valid", r.Output)

	r = JSONValidate(`true`)
	require.Empty(t, r.Error)
	assert.Equal(t, "valid", r.Output)

	r = JSONValidate(`null`)
	require.Empty(t, r.Error)
	assert.Equal(t, "valid", r.Output)
}

func TestJSONValidate_EmptyInput(t *testing.T) {
	r := JSONValidate("")
	assert.Equal(t, "empty input", r.Error)

	r = JSONValidate("   ")
	assert.Equal(t, "empty input", r.Error)
}

func TestJSONValidate_Invalid(t *testing.T) {
	r := JSONValidate("{bad}")
	assert.NotEmpty(t, r.Error)
	assert.Contains(t, r.Error, "invalid JSON")
}

func TestJSONValidate_InvalidTrailingComma(t *testing.T) {
	r := JSONValidate(`{"a":1,}`)
	assert.NotEmpty(t, r.Error)
	assert.Contains(t, r.Error, "invalid JSON")
}

func TestJSONValidate_WhitespaceTrimmed(t *testing.T) {
	r := JSONValidate(`  {"key":"value"}  `)
	require.Empty(t, r.Error)
	assert.Equal(t, "valid", r.Output)
}

// ---------------------------------------------------------------------------
// Tool interface compliance
// ---------------------------------------------------------------------------

func TestJSONToolImplementsToolInterface(t *testing.T) {
	var _ Tool = JSONTool{}
}
