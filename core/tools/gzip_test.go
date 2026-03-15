package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------- Tool metadata ----------

func TestGZipTool_Metadata(t *testing.T) {
	tool := GZipTool{}

	assert.Equal(t, "GZip Compress / Decompress", tool.Name())
	assert.Equal(t, "gzip", tool.ID())
	assert.Equal(t, "Encoders", tool.Category())
	assert.NotEmpty(t, tool.Description())

	keywords := tool.Keywords()
	for _, kw := range []string{"gzip", "compress", "decompress", "zip"} {
		assert.Contains(t, keywords, kw)
	}
}

func TestGZipTool_DetectFromClipboard(t *testing.T) {
	tool := GZipTool{}
	assert.False(t, tool.DetectFromClipboard("anything"))
	assert.False(t, tool.DetectFromClipboard(""))
}

// ---------- GZipCompress ----------

func TestGZipCompress_HelloWorld(t *testing.T) {
	r := GZipCompress("hello world")
	require.Empty(t, r.Error)
	assert.NotEmpty(t, r.Output)
	// Output should be valid base64
	assert.NotContains(t, r.Output, " ")
}

func TestGZipCompress_EmptyInput(t *testing.T) {
	r := GZipCompress("")
	require.Empty(t, r.Error)
	assert.Equal(t, "", r.Output)
}

// ---------- GZipDecompress ----------

func TestGZipDecompress_Roundtrip(t *testing.T) {
	original := "hello world"
	compressed := GZipCompress(original)
	require.Empty(t, compressed.Error)

	decompressed := GZipDecompress(compressed.Output)
	require.Empty(t, decompressed.Error)
	assert.Equal(t, original, decompressed.Output)
}

func TestGZipDecompress_LargerText(t *testing.T) {
	original := "The quick brown fox jumps over the lazy dog. " +
		"Pack my box with five dozen liquor jugs. " +
		"How vexingly quick daft zebras jump!"
	compressed := GZipCompress(original)
	require.Empty(t, compressed.Error)

	decompressed := GZipDecompress(compressed.Output)
	require.Empty(t, decompressed.Error)
	assert.Equal(t, original, decompressed.Output)
}

func TestGZipDecompress_InvalidBase64(t *testing.T) {
	r := GZipDecompress("not!valid!base64!!!")
	assert.NotEmpty(t, r.Error)
	assert.Contains(t, r.Error, "invalid base64")
	assert.Empty(t, r.Output)
}

func TestGZipDecompress_InvalidGzipData(t *testing.T) {
	// Valid base64 but not valid gzip data
	r := GZipDecompress("aGVsbG8gd29ybGQ=")
	assert.NotEmpty(t, r.Error)
	assert.Contains(t, r.Error, "invalid gzip data")
	assert.Empty(t, r.Output)
}

func TestGZipDecompress_EmptyInput(t *testing.T) {
	r := GZipDecompress("")
	require.Empty(t, r.Error)
	assert.Equal(t, "", r.Output)
}

// ---------- Tool interface compliance ----------

func TestGZipTool_ImplementsToolInterface(t *testing.T) {
	var _ Tool = GZipTool{}
}
