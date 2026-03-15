package tools

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

// --- JWTTool metadata tests ---

func TestJWTToolMetadata(t *testing.T) {
	tool := JWTTool{}

	assert.Equal(t, "JWT Decoder", tool.Name())
	assert.Equal(t, "jwt", tool.ID())
	assert.Equal(t, "Encoders", tool.Category())
	assert.NotEmpty(t, tool.Description())
	assert.Contains(t, tool.Keywords(), "jwt")
	assert.Contains(t, tool.Keywords(), "token")
	assert.Contains(t, tool.Keywords(), "decode")
	assert.Contains(t, tool.Keywords(), "json web token")
}

// --- DetectFromClipboard tests ---

func TestJWTDetectFromClipboard(t *testing.T) {
	tool := JWTTool{}

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid JWT", testJWT, true},
		{"three segments", "aaa.bbb.ccc", true},
		{"two segments", "aaa.bbb", false},
		{"four segments", "aaa.bbb.ccc.ddd", false},
		{"empty string", "", false},
		{"single segment", "abcdef", false},
		{"empty first segment", ".bbb.ccc", false},
		{"empty middle segment", "aaa..ccc", false},
		{"empty last segment", "aaa.bbb.", false},
		{"with leading/trailing whitespace", "  " + testJWT + "  ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tool.DetectFromClipboard(tt.input))
		})
	}
}

// --- JWTDecode tests ---

func TestJWTDecode_ValidToken(t *testing.T) {
	result := JWTDecode(testJWT)

	require.Empty(t, result.Error, "expected no error")
	assert.NotEmpty(t, result.Header)
	assert.NotEmpty(t, result.Payload)
	assert.NotEmpty(t, result.Signature)

	// Header should contain alg and typ
	assert.Contains(t, result.Header, `"alg"`)
	assert.Contains(t, result.Header, `"HS256"`)
	assert.Contains(t, result.Header, `"typ"`)
	assert.Contains(t, result.Header, `"JWT"`)

	// Payload should contain sub, name, iat
	assert.Contains(t, result.Payload, `"sub"`)
	assert.Contains(t, result.Payload, `"1234567890"`)
	assert.Contains(t, result.Payload, `"name"`)
	assert.Contains(t, result.Payload, `"John Doe"`)
	assert.Contains(t, result.Payload, `"iat"`)

	// Signature should be the raw third segment
	assert.Equal(t, "SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c", result.Signature)
}

func TestJWTDecode_OutputFormat(t *testing.T) {
	result := JWTDecode(testJWT)

	require.Empty(t, result.Error)
	assert.True(t, strings.HasPrefix(result.Output, "--- Header ---\n"))
	assert.Contains(t, result.Output, "\n--- Payload ---\n")
	assert.Contains(t, result.Output, "\n--- Signature ---\n")
}

func TestJWTDecode_PrettyPrintsJSON(t *testing.T) {
	result := JWTDecode(testJWT)

	require.Empty(t, result.Error)
	// Pretty-printed JSON should have 2-space indentation
	assert.Contains(t, result.Header, "  ")
	assert.Contains(t, result.Payload, "  ")
}

func TestJWTDecode_WhitespaceHandling(t *testing.T) {
	result := JWTDecode("  " + testJWT + "  ")

	require.Empty(t, result.Error)
	assert.Contains(t, result.Header, `"alg"`)
}

func TestJWTDecode_TooFewSegments(t *testing.T) {
	result := JWTDecode("only.two")

	assert.NotEmpty(t, result.Error)
	assert.Contains(t, result.Error, "3 dot-separated segments")
}

func TestJWTDecode_TooManySegments(t *testing.T) {
	result := JWTDecode("one.two.three.four")

	assert.NotEmpty(t, result.Error)
	assert.Contains(t, result.Error, "3 dot-separated segments")
}

func TestJWTDecode_EmptySegment(t *testing.T) {
	result := JWTDecode("aaa..ccc")

	assert.NotEmpty(t, result.Error)
	assert.Contains(t, result.Error, "empty")
}

func TestJWTDecode_InvalidBase64Header(t *testing.T) {
	result := JWTDecode("!!!.eyJzdWIiOiIxIn0.sig")

	assert.NotEmpty(t, result.Error)
	assert.Contains(t, result.Error, "header")
}

func TestJWTDecode_InvalidBase64Payload(t *testing.T) {
	result := JWTDecode("eyJhbGciOiJIUzI1NiJ9.!!!.sig")

	assert.NotEmpty(t, result.Error)
	assert.Contains(t, result.Error, "payload")
}

func TestJWTDecode_InvalidJSONInHeader(t *testing.T) {
	// "bm90anNvbg" is base64url for "notjson"
	result := JWTDecode("bm90anNvbg.eyJzdWIiOiIxIn0.sig")

	assert.NotEmpty(t, result.Error)
	assert.Contains(t, result.Error, "header")
}

func TestJWTDecode_EmptyString(t *testing.T) {
	result := JWTDecode("")

	assert.NotEmpty(t, result.Error)
}

// --- JWTValidate tests ---

func TestJWTValidate_ValidToken(t *testing.T) {
	result := JWTValidate(testJWT)

	assert.Empty(t, result.Error)
	assert.Equal(t, "valid", result.Output)
}

func TestJWTValidate_WhitespaceHandling(t *testing.T) {
	result := JWTValidate("  " + testJWT + "\n")

	assert.Empty(t, result.Error)
	assert.Equal(t, "valid", result.Output)
}

func TestJWTValidate_TooFewSegments(t *testing.T) {
	result := JWTValidate("only.two")

	assert.NotEmpty(t, result.Error)
	assert.Contains(t, result.Error, "3 dot-separated segments")
}

func TestJWTValidate_EmptySegment(t *testing.T) {
	result := JWTValidate("aaa..ccc")

	assert.NotEmpty(t, result.Error)
	assert.Contains(t, result.Error, "empty")
}

func TestJWTValidate_InvalidBase64(t *testing.T) {
	result := JWTValidate("!!!.eyJzdWIiOiIxIn0.sig")

	assert.NotEmpty(t, result.Error)
	assert.Contains(t, result.Error, "base64url")
}

func TestJWTValidate_InvalidJSONHeader(t *testing.T) {
	// "bm90anNvbg" is base64url for "notjson"
	result := JWTValidate("bm90anNvbg.eyJzdWIiOiIxIn0.c2ln")

	assert.NotEmpty(t, result.Error)
	assert.Contains(t, result.Error, "header")
	assert.Contains(t, result.Error, "JSON")
}

func TestJWTValidate_InvalidJSONPayload(t *testing.T) {
	// header is valid JSON, payload "bm90anNvbg" is "notjson"
	result := JWTValidate("eyJhbGciOiJIUzI1NiJ9.bm90anNvbg.c2ln")

	assert.NotEmpty(t, result.Error)
	assert.Contains(t, result.Error, "payload")
	assert.Contains(t, result.Error, "JSON")
}

func TestJWTValidate_EmptyString(t *testing.T) {
	result := JWTValidate("")

	assert.NotEmpty(t, result.Error)
}

// --- Helper function tests ---

func TestDecodeJWTSegment_Valid(t *testing.T) {
	// "eyJhbGciOiJIUzI1NiJ9" is base64url for {"alg":"HS256"}
	decoded, err := decodeJWTSegment("eyJhbGciOiJIUzI1NiJ9")

	require.NoError(t, err)
	assert.Contains(t, decoded, `"alg"`)
	assert.Contains(t, decoded, `"HS256"`)
}

func TestDecodeJWTSegment_InvalidBase64(t *testing.T) {
	_, err := decodeJWTSegment("!!!")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "base64url")
}

func TestDecodeJWTSegment_InvalidJSON(t *testing.T) {
	// "bm90anNvbg" is base64url for "notjson"
	_, err := decodeJWTSegment("bm90anNvbg")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JSON")
}

func TestFormatJSON_Valid(t *testing.T) {
	result, err := formatJSON(`{"a":1,"b":"hello"}`)

	require.NoError(t, err)
	assert.Contains(t, result, "  ")
	assert.Contains(t, result, `"a"`)
	assert.Contains(t, result, `"b"`)
}

func TestFormatJSON_Invalid(t *testing.T) {
	_, err := formatJSON("not json")

	assert.Error(t, err)
}
