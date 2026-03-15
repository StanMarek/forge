package tools

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
)

// GZipTool provides metadata for the GZip Compress / Decompress tool.
type GZipTool struct{}

func (g GZipTool) Name() string        { return "GZip Compress / Decompress" }
func (g GZipTool) ID() string          { return "gzip" }
func (g GZipTool) Description() string { return "Compress and decompress data using gzip" }
func (g GZipTool) Category() string    { return "Encoders" }
func (g GZipTool) Keywords() []string {
	return []string{"gzip", "compress", "decompress", "zip"}
}

// DetectFromClipboard always returns false because binary gzip detection
// from clipboard text is unreliable.
func (g GZipTool) DetectFromClipboard(_ string) bool {
	return false
}

// GZipCompress compresses the input string using gzip and returns the
// result as a base64-encoded string (since gzip output is binary).
func GZipCompress(input string) Result {
	if input == "" {
		return Result{Output: ""}
	}

	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	_, err := writer.Write([]byte(input))
	if err != nil {
		return Result{Error: fmt.Sprintf("gzip compress: %s", err.Error())}
	}

	if err := writer.Close(); err != nil {
		return Result{Error: fmt.Sprintf("gzip compress: %s", err.Error())}
	}

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	return Result{Output: encoded}
}

// GZipDecompress base64-decodes the input, then decompresses it with gzip.
// Returns the decompressed plaintext string.
func GZipDecompress(input string) Result {
	if input == "" {
		return Result{Output: ""}
	}

	compressed, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return Result{Error: fmt.Sprintf("invalid base64: %s", err.Error())}
	}

	reader, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return Result{Error: fmt.Sprintf("invalid gzip data: %s", err.Error())}
	}
	defer reader.Close()

	decompressed, err := io.ReadAll(reader)
	if err != nil {
		return Result{Error: fmt.Sprintf("gzip decompress: %s", err.Error())}
	}

	return Result{Output: string(decompressed)}
}
