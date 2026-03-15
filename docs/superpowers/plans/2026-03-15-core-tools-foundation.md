# Core Tools Foundation Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the pure business logic layer — 6 Tier-1 tools, Tool interface, Result types, and Registry — as the foundation for all Forge UI surfaces.

**Architecture:** Flat `core/tools/` package with one file per tool containing pure functions. `core/registry/` provides discovery, search, and detection. All functions are stateless, return structs with `Output`/`Error` fields, and have zero I/O.

**Tech Stack:** Go 1.25, `encoding/base64`, `encoding/json`, `crypto/*`, `net/url` (stdlib), `github.com/google/uuid`, `github.com/stretchr/testify` (test)

**Spec:** `docs/superpowers/specs/2026-03-15-core-tools-foundation-design.md`

---

## File Structure

| File | Responsibility |
|------|---------------|
| `go.mod` | Module definition, dependencies |
| `core/tools/tool.go` | `Tool` interface, `Result` struct, category constants |
| `core/tools/base64.go` | `Base64Tool` struct, `Base64Encode`, `Base64Decode` |
| `core/tools/base64_test.go` | Tests for base64 encode/decode/detection |
| `core/tools/jwt.go` | `JWTTool` struct, `JWTDecode`, `JWTValidate`, `JWTDecodeResult` |
| `core/tools/jwt_test.go` | Tests for JWT decode/validate/detection |
| `core/tools/json.go` | `JSONTool` struct, `JSONFormat`, `JSONMinify`, `JSONValidate` |
| `core/tools/json_test.go` | Tests for JSON format/minify/validate/detection |
| `core/tools/hash.go` | `HashTool` struct, `Hash` |
| `core/tools/hash_test.go` | Tests for hash algorithms/detection |
| `core/tools/url.go` | `URLTool` struct, `URLEncode`, `URLDecode`, `URLParse`, `URLParseResult` |
| `core/tools/url_test.go` | Tests for URL encode/decode/parse/detection |
| `core/tools/uuid.go` | `UUIDTool` struct, `UUIDGenerate`, `UUIDValidate`, `UUIDParse`, `UUIDParseResult` |
| `core/tools/uuid_test.go` | Tests for UUID generate/validate/parse/detection |
| `core/registry/registry.go` | `Registry` struct with `Register`, `All`, `ByID`, `ByCategory`, `Search`, `Detect` |
| `core/registry/registry_test.go` | Tests for all registry operations |
| `core/registry/defaults.go` | `Default()` factory with all 6 tools pre-registered |

---

## Chunk 1: Project Bootstrap + Tool Interface

### Task 1: Initialize Go module

**Files:**
- Create: `go.mod`

- [ ] **Step 1: Initialize module and add dependencies**

```bash
cd /Users/stanislawmarek/Desktop/coding/forge
go mod init github.com/StanMarek/forge
go get github.com/google/uuid
go get github.com/stretchr/testify
```

- [ ] **Step 2: Verify go.mod**

Run: `cat go.mod`
Expected: module path is `github.com/StanMarek/forge`, requires `google/uuid` and `testify`

- [ ] **Step 3: Commit**

```bash
git add go.mod go.sum
git commit -m "Initialize Go module with uuid and testify dependencies"
```

### Task 2: Create Tool interface and Result types

**Files:**
- Create: `core/tools/tool.go`

- [ ] **Step 1: Write tool.go**

```go
package tools

// Result is the standard return type for tool operations.
// Success is indicated by Error == "".
type Result struct {
	Output string
	Error  string
}

// Tool defines the metadata interface for tool discovery and routing.
// Tool logic lives in standalone functions, NOT on this interface.
type Tool interface {
	Name() string
	ID() string
	Description() string
	Category() string
	Keywords() []string
	DetectFromClipboard(s string) bool
}
```

- [ ] **Step 2: Verify it compiles**

Run: `go build ./core/tools/`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add core/tools/tool.go
git commit -m "Add Tool interface and Result type"
```

---

## Chunk 2: Base64 Tool

### Task 3: Base64 — tests

**Files:**
- Create: `core/tools/base64_test.go`

- [ ] **Step 1: Write base64 tests**

```go
package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64Encode(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		urlSafe   bool
		noPadding bool
		expected  string
	}{
		{"simple string", "Hello, World!", false, false, "SGVsbG8sIFdvcmxkIQ=="},
		{"empty string", "", false, false, ""},
		{"url-safe encoding", "https://example.com?foo=bar", true, false, "aHR0cHM6Ly9leGFtcGxlLmNvbT9mb289YmFy"},
		{"no padding", "Hello, World!", false, true, "SGVsbG8sIFdvcmxkIQ"},
		{"url-safe no padding", "subjects?_d", true, true, "c3ViamVjdHM_X2Q"},
		{"unicode", "こんにちは", false, false, "44GT44KT44Gr44Gh44Gv"},
		{"binary-like chars", "\x00\x01\x02\x03", false, false, "AAECAw=="},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Base64Encode(tt.input, tt.urlSafe, tt.noPadding)
			assert.Empty(t, result.Error)
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
	}{
		{"simple string", "SGVsbG8sIFdvcmxkIQ==", false, "Hello, World!"},
		{"empty string", "", false, ""},
		{"url-safe decoding", "aHR0cHM6Ly9leGFtcGxlLmNvbT9mb289YmFy", true, "https://example.com?foo=bar"},
		{"no padding accepted", "SGVsbG8sIFdvcmxkIQ", false, "Hello, World!"},
		{"unicode", "44GT44KT44Gr44Gh44Gv", false, "こんにちは"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Base64Decode(tt.input, tt.urlSafe)
			assert.Empty(t, result.Error)
			assert.Equal(t, tt.expected, result.Output)
		})
	}
}

func TestBase64DecodeInvalid(t *testing.T) {
	result := Base64Decode("not-valid-base64!!!", false)
	assert.NotEmpty(t, result.Error)
	assert.Empty(t, result.Output)
}

func TestBase64ToolDetection(t *testing.T) {
	tool := Base64Tool{}
	assert.True(t, tool.DetectFromClipboard("SGVsbG8sIFdvcmxkIQ=="))
	assert.True(t, tool.DetectFromClipboard("AAAA"))
	assert.False(t, tool.DetectFromClipboard("abc"))          // too short
	assert.False(t, tool.DetectFromClipboard("hello world"))  // spaces
	assert.False(t, tool.DetectFromClipboard("AAAAA"))        // not mod 4
	assert.False(t, tool.DetectFromClipboard(""))             // empty
}

func TestBase64ToolMetadata(t *testing.T) {
	tool := Base64Tool{}
	assert.Equal(t, "base64", tool.ID())
	assert.Equal(t, "Encoders", tool.Category())
	assert.NotEmpty(t, tool.Name())
	assert.NotEmpty(t, tool.Description())
	assert.NotEmpty(t, tool.Keywords())
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./core/tools/ -run TestBase64 -v`
Expected: FAIL — `Base64Encode` undefined

### Task 4: Base64 — implementation

**Files:**
- Create: `core/tools/base64.go`

- [ ] **Step 1: Write base64 implementation**

```go
package tools

import (
	"encoding/base64"
	"regexp"
	"strings"
)

// Base64Tool implements the Tool interface for Base64 encoding/decoding.
type Base64Tool struct{}

func (Base64Tool) Name() string        { return "Base64 Encode / Decode" }
func (Base64Tool) ID() string          { return "base64" }
func (Base64Tool) Description() string { return "Encode and decode Base64 strings" }
func (Base64Tool) Category() string    { return "Encoders" }
func (Base64Tool) Keywords() []string  { return []string{"base64", "encode", "decode", "b64"} }

var base64Regex = regexp.MustCompile(`^[A-Za-z0-9+/=]+$`)

func (Base64Tool) DetectFromClipboard(s string) bool {
	if len(s) < 4 || len(s)%4 != 0 {
		return false
	}
	return base64Regex.MatchString(s)
}

// Base64Encode encodes input as Base64.
func Base64Encode(input string, urlSafe bool, noPadding bool) Result {
	if input == "" {
		return Result{Output: ""}
	}
	var enc *base64.Encoding
	if urlSafe {
		enc = base64.URLEncoding
	} else {
		enc = base64.StdEncoding
	}
	if noPadding {
		enc = enc.WithPadding(base64.NoPadding)
	}
	return Result{Output: enc.EncodeToString([]byte(input))}
}

// Base64Decode decodes a Base64-encoded string.
func Base64Decode(input string, urlSafe bool) Result {
	if input == "" {
		return Result{Output: ""}
	}
	var enc *base64.Encoding
	if urlSafe {
		enc = base64.URLEncoding
	} else {
		enc = base64.StdEncoding
	}
	// Try with padding first, then without
	decoded, err := enc.DecodeString(input)
	if err != nil {
		enc = enc.WithPadding(base64.NoPadding)
		decoded, err = enc.DecodeString(strings.TrimRight(input, "="))
		if err != nil {
			return Result{Error: "invalid base64: " + err.Error()}
		}
	}
	return Result{Output: string(decoded)}
}
```

- [ ] **Step 2: Run tests to verify they pass**

Run: `go test ./core/tools/ -run TestBase64 -v`
Expected: all PASS

- [ ] **Step 3: Commit**

```bash
git add core/tools/base64.go core/tools/base64_test.go
git commit -m "feat: add base64 encode/decode tool with tests"
```

---

## Chunk 3: JWT Tool

### Task 5: JWT — tests

**Files:**
- Create: `core/tools/jwt_test.go`

- [ ] **Step 1: Write JWT tests**

Use this well-known test token: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c`

```go
package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

func TestJWTDecode(t *testing.T) {
	result := JWTDecode(testJWT)
	assert.Empty(t, result.Error)
	assert.Contains(t, result.Header, `"alg"`)
	assert.Contains(t, result.Header, `"HS256"`)
	assert.Contains(t, result.Payload, `"sub"`)
	assert.Contains(t, result.Payload, `"1234567890"`)
	assert.Contains(t, result.Payload, `"John Doe"`)
	assert.Equal(t, "SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c", result.Signature)
	assert.NotEmpty(t, result.Output)
}

func TestJWTDecodeEmpty(t *testing.T) {
	result := JWTDecode("")
	assert.NotEmpty(t, result.Error)
}

func TestJWTDecodeInvalid(t *testing.T) {
	result := JWTDecode("not.a.jwt")
	assert.NotEmpty(t, result.Error)
}

func TestJWTDecodeTwoParts(t *testing.T) {
	result := JWTDecode("only.twoparts")
	assert.NotEmpty(t, result.Error)
}

func TestJWTValidate(t *testing.T) {
	result := JWTValidate(testJWT)
	assert.Empty(t, result.Error)
	assert.Equal(t, "valid", result.Output)
}

func TestJWTValidateInvalid(t *testing.T) {
	result := JWTValidate("garbage")
	assert.NotEmpty(t, result.Error)
}

func TestJWTValidateEmpty(t *testing.T) {
	result := JWTValidate("")
	assert.NotEmpty(t, result.Error)
}

func TestJWTToolDetection(t *testing.T) {
	tool := JWTTool{}
	assert.True(t, tool.DetectFromClipboard(testJWT))
	assert.False(t, tool.DetectFromClipboard("hello"))
	assert.False(t, tool.DetectFromClipboard("two.parts"))
	assert.False(t, tool.DetectFromClipboard(""))
}

func TestJWTToolMetadata(t *testing.T) {
	tool := JWTTool{}
	assert.Equal(t, "jwt", tool.ID())
	assert.Equal(t, "Encoders", tool.Category())
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./core/tools/ -run TestJWT -v`
Expected: FAIL — `JWTDecode` undefined

### Task 6: JWT — implementation

**Files:**
- Create: `core/tools/jwt.go`

- [ ] **Step 1: Write JWT implementation**

```go
package tools

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// JWTDecodeResult holds the decoded components of a JWT token.
type JWTDecodeResult struct {
	Header    string
	Payload   string
	Signature string
	Output    string
	Error     string
}

// JWTTool implements the Tool interface for JWT decoding.
type JWTTool struct{}

func (JWTTool) Name() string        { return "JWT Decoder" }
func (JWTTool) ID() string          { return "jwt" }
func (JWTTool) Description() string { return "Decode and inspect JWT tokens" }
func (JWTTool) Category() string    { return "Encoders" }
func (JWTTool) Keywords() []string  { return []string{"jwt", "token", "decode", "json web token"} }

func (JWTTool) DetectFromClipboard(s string) bool {
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return false
	}
	for _, p := range parts {
		if len(p) == 0 {
			return false
		}
	}
	return true
}

// JWTDecode decodes a JWT token without validating the signature.
func JWTDecode(token string) JWTDecodeResult {
	if token == "" {
		return JWTDecodeResult{Error: "empty token"}
	}
	parts := strings.Split(strings.TrimSpace(token), ".")
	if len(parts) != 3 {
		return JWTDecodeResult{Error: fmt.Sprintf("invalid JWT: expected 3 segments, got %d", len(parts))}
	}

	header, err := decodeJWTSegment(parts[0])
	if err != nil {
		return JWTDecodeResult{Error: "invalid JWT header: " + err.Error()}
	}
	payload, err := decodeJWTSegment(parts[1])
	if err != nil {
		return JWTDecodeResult{Error: "invalid JWT payload: " + err.Error()}
	}

	headerFormatted, _ := formatJSON(header)
	payloadFormatted, _ := formatJSON(payload)

	output := fmt.Sprintf("--- Header ---\n%s\n--- Payload ---\n%s\n--- Signature ---\n%s",
		headerFormatted, payloadFormatted, parts[2])

	return JWTDecodeResult{
		Header:    headerFormatted,
		Payload:   payloadFormatted,
		Signature: parts[2],
		Output:    output,
	}
}

// JWTValidate checks if a string is a structurally valid JWT.
func JWTValidate(token string) Result {
	if token == "" {
		return Result{Error: "empty token"}
	}
	parts := strings.Split(strings.TrimSpace(token), ".")
	if len(parts) != 3 {
		return Result{Error: fmt.Sprintf("invalid JWT: expected 3 segments, got %d", len(parts))}
	}
	for i, name := range []string{"header", "payload"} {
		decoded, err := decodeJWTSegment(parts[i])
		if err != nil {
			return Result{Error: fmt.Sprintf("invalid JWT %s: %s", name, err.Error())}
		}
		if !json.Valid([]byte(decoded)) {
			return Result{Error: fmt.Sprintf("invalid JWT %s: not valid JSON", name)}
		}
	}
	return Result{Output: "valid"}
}

func decodeJWTSegment(segment string) (string, error) {
	// JWT uses base64url encoding without padding
	decoded, err := base64.RawURLEncoding.DecodeString(segment)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func formatJSON(s string) (string, error) {
	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return s, err
	}
	formatted, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return s, err
	}
	return string(formatted), nil
}
```

- [ ] **Step 2: Run tests to verify they pass**

Run: `go test ./core/tools/ -run TestJWT -v`
Expected: all PASS

- [ ] **Step 3: Commit**

```bash
git add core/tools/jwt.go core/tools/jwt_test.go
git commit -m "feat: add JWT decode/validate tool with tests"
```

---

## Chunk 4: JSON Tool

### Task 7: JSON — tests

**Files:**
- Create: `core/tools/json_test.go`

- [ ] **Step 1: Write JSON tests**

```go
package tools

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONFormat(t *testing.T) {
	input := `{"name":"forge","version":1}`
	result := JSONFormat(input, 2, false, false)
	assert.Empty(t, result.Error)
	assert.Contains(t, result.Output, "  \"name\"")
	assert.Contains(t, result.Output, "forge")
}

func TestJSONFormatWithTabs(t *testing.T) {
	input := `{"a":1}`
	result := JSONFormat(input, 0, false, true)
	assert.Empty(t, result.Error)
	assert.Contains(t, result.Output, "\t\"a\"")
}

func TestJSONFormatSortKeys(t *testing.T) {
	input := `{"z":1,"a":2,"m":3}`
	result := JSONFormat(input, 2, true, false)
	assert.Empty(t, result.Error)
	// Keys should appear in alphabetical order: a, m, z
	aIdx := strings.Index(result.Output, `"a"`)
	mIdx := strings.Index(result.Output, `"m"`)
	zIdx := strings.Index(result.Output, `"z"`)
	assert.Less(t, aIdx, mIdx)
	assert.Less(t, mIdx, zIdx)
}

func TestJSONFormatSortKeysNested(t *testing.T) {
	input := `{"z":{"c":1,"a":2},"a":0}`
	result := JSONFormat(input, 2, true, false)
	assert.Empty(t, result.Error)
	// Outer keys sorted: a before z
	aIdx := strings.Index(result.Output, `"a"`)
	zIdx := strings.Index(result.Output, `"z"`)
	assert.Less(t, aIdx, zIdx)
}

func TestJSONFormatIndent4(t *testing.T) {
	input := `{"a":1}`
	result := JSONFormat(input, 4, false, false)
	assert.Empty(t, result.Error)
	assert.Contains(t, result.Output, "    \"a\"")
}

func TestJSONFormatInvalid(t *testing.T) {
	result := JSONFormat("not json", 2, false, false)
	assert.NotEmpty(t, result.Error)
}

func TestJSONFormatEmpty(t *testing.T) {
	result := JSONFormat("", 2, false, false)
	assert.NotEmpty(t, result.Error)
}

func TestJSONMinify(t *testing.T) {
	input := "{\n  \"name\": \"forge\",\n  \"version\": 1\n}"
	result := JSONMinify(input)
	assert.Empty(t, result.Error)
	assert.Equal(t, `{"name":"forge","version":1}`, result.Output)
}

func TestJSONMinifyInvalid(t *testing.T) {
	result := JSONMinify("not json")
	assert.NotEmpty(t, result.Error)
}

func TestJSONValidate(t *testing.T) {
	result := JSONValidate(`{"valid":true}`)
	assert.Empty(t, result.Error)
	assert.Equal(t, "valid", result.Output)
}

func TestJSONValidateInvalid(t *testing.T) {
	result := JSONValidate(`{"missing": }`)
	assert.NotEmpty(t, result.Error)
}

func TestJSONValidateEmpty(t *testing.T) {
	result := JSONValidate("")
	assert.NotEmpty(t, result.Error)
}

func TestJSONToolDetection(t *testing.T) {
	tool := JSONTool{}
	assert.True(t, tool.DetectFromClipboard(`{"key":"value"}`))
	assert.True(t, tool.DetectFromClipboard(`[1,2,3]`))
	assert.False(t, tool.DetectFromClipboard("not json"))
	assert.False(t, tool.DetectFromClipboard(""))
}

func TestJSONToolMetadata(t *testing.T) {
	tool := JSONTool{}
	assert.Equal(t, "json", tool.ID())
	assert.Equal(t, "Formatters", tool.Category())
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./core/tools/ -run TestJSON -v`
Expected: FAIL — `JSONFormat` undefined

### Task 8: JSON — implementation

**Files:**
- Create: `core/tools/json.go`

- [ ] **Step 1: Write JSON implementation**

```go
package tools

import (
	"bytes"
	"encoding/json"
	"strings"
)

// JSONTool implements the Tool interface for JSON formatting.
type JSONTool struct{}

func (JSONTool) Name() string        { return "JSON Formatter" }
func (JSONTool) ID() string          { return "json" }
func (JSONTool) Description() string { return "Format, minify, or validate JSON" }
func (JSONTool) Category() string    { return "Formatters" }
func (JSONTool) Keywords() []string  { return []string{"json", "format", "minify", "validate", "pretty"} }

func (JSONTool) DetectFromClipboard(s string) bool {
	s = strings.TrimSpace(s)
	return len(s) > 0 && json.Valid([]byte(s))
}

// JSONFormat pretty-prints JSON with configurable indentation.
func JSONFormat(input string, indent int, sortKeys bool, useTabs bool) Result {
	input = strings.TrimSpace(input)
	if input == "" {
		return Result{Error: "empty input"}
	}
	if !json.Valid([]byte(input)) {
		return Result{Error: "invalid JSON: " + findJSONError(input)}
	}

	indentStr := strings.Repeat(" ", indent)
	if useTabs {
		indentStr = "\t"
	}

	if sortKeys {
		var v interface{}
		if err := json.Unmarshal([]byte(input), &v); err != nil {
			return Result{Error: "invalid JSON: " + err.Error()}
		}
		formatted, err := json.MarshalIndent(v, "", indentStr)
		if err != nil {
			return Result{Error: "format error: " + err.Error()}
		}
		return Result{Output: string(formatted)}
	}

	var buf bytes.Buffer
	if err := json.Indent(&buf, []byte(input), "", indentStr); err != nil {
		return Result{Error: "format error: " + err.Error()}
	}
	return Result{Output: buf.String()}
}

// JSONMinify removes all unnecessary whitespace from JSON.
func JSONMinify(input string) Result {
	input = strings.TrimSpace(input)
	if input == "" {
		return Result{Error: "empty input"}
	}
	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(input)); err != nil {
		return Result{Error: "invalid JSON: " + err.Error()}
	}
	return Result{Output: buf.String()}
}

// JSONValidate checks if input is well-formed JSON.
func JSONValidate(input string) Result {
	input = strings.TrimSpace(input)
	if input == "" {
		return Result{Error: "empty input"}
	}
	if json.Valid([]byte(input)) {
		return Result{Output: "valid"}
	}
	return Result{Error: "invalid JSON: " + findJSONError(input)}
}

func findJSONError(input string) string {
	var v interface{}
	err := json.Unmarshal([]byte(input), &v)
	if err != nil {
		return err.Error()
	}
	return "unknown error"
}
```

- [ ] **Step 2: Run tests to verify they pass**

Run: `go test ./core/tools/ -run TestJSON -v`
Expected: all PASS

- [ ] **Step 3: Commit**

```bash
git add core/tools/json.go core/tools/json_test.go
git commit -m "feat: add JSON format/minify/validate tool with tests"
```

---

## Chunk 5: Hash Tool

### Task 9: Hash — tests

**Files:**
- Create: `core/tools/hash_test.go`

- [ ] **Step 1: Write hash tests**

Known test vectors — `echo -n "hello world" | sha256sum` = `b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9`

```go
package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashMD5(t *testing.T) {
	result := Hash("hello world", "md5", false)
	assert.Empty(t, result.Error)
	assert.Equal(t, "5eb63bbbe01eeed093cb22bb8f5acdc3", result.Output)
}

func TestHashSHA1(t *testing.T) {
	result := Hash("hello world", "sha1", false)
	assert.Empty(t, result.Error)
	assert.Equal(t, "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed", result.Output)
}

func TestHashSHA256(t *testing.T) {
	result := Hash("hello world", "sha256", false)
	assert.Empty(t, result.Error)
	assert.Equal(t, "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9", result.Output)
}

func TestHashSHA512(t *testing.T) {
	result := Hash("hello world", "sha512", false)
	assert.Empty(t, result.Error)
	assert.Equal(t, "309ecc489c12d6eb4cc40f50c902f2b4d0ed77ee511a7c7a9bcd3ca86d4cd86f989dd35bc5ff499670da34255b45b0cfd830e81f605dcf7dc5542e93ae9cd76f", result.Output)
}

func TestHashUppercase(t *testing.T) {
	result := Hash("hello world", "md5", true)
	assert.Empty(t, result.Error)
	assert.Equal(t, "5EB63BBBE01EEED093CB22BB8F5ACDC3", result.Output)
}

func TestHashUnknownAlgorithm(t *testing.T) {
	result := Hash("hello", "sha999", false)
	assert.NotEmpty(t, result.Error)
	assert.Contains(t, result.Error, "unsupported")
}

func TestHashEmpty(t *testing.T) {
	result := Hash("", "sha256", false)
	assert.Empty(t, result.Error)
	assert.Equal(t, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", result.Output)
}

func TestHashToolDetection(t *testing.T) {
	tool := HashTool{}
	assert.False(t, tool.DetectFromClipboard("anything"))
	assert.False(t, tool.DetectFromClipboard(""))
}

func TestHashToolMetadata(t *testing.T) {
	tool := HashTool{}
	assert.Equal(t, "hash", tool.ID())
	assert.Equal(t, "Generators", tool.Category())
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./core/tools/ -run TestHash -v`
Expected: FAIL — `Hash` undefined

### Task 10: Hash — implementation

**Files:**
- Create: `core/tools/hash.go`

- [ ] **Step 1: Write hash implementation**

```go
package tools

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"strings"
)

// HashTool implements the Tool interface for hash generation.
type HashTool struct{}

func (HashTool) Name() string        { return "Hash Generator" }
func (HashTool) ID() string          { return "hash" }
func (HashTool) Description() string { return "Generate hash digests (MD5, SHA-1, SHA-256, SHA-512)" }
func (HashTool) Category() string    { return "Generators" }
func (HashTool) Keywords() []string  { return []string{"hash", "md5", "sha1", "sha256", "sha512", "digest", "checksum"} }

func (HashTool) DetectFromClipboard(_ string) bool { return false }

// Hash computes a hash digest of the input string.
func Hash(input string, algorithm string, uppercase bool) Result {
	var h hash.Hash
	switch strings.ToLower(algorithm) {
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	case "sha256":
		h = sha256.New()
	case "sha512":
		h = sha512.New()
	default:
		return Result{Error: fmt.Sprintf("unsupported algorithm: %s (supported: md5, sha1, sha256, sha512)", algorithm)}
	}
	h.Write([]byte(input))
	digest := fmt.Sprintf("%x", h.Sum(nil))
	if uppercase {
		digest = strings.ToUpper(digest)
	}
	return Result{Output: digest}
}
```

- [ ] **Step 2: Run tests to verify they pass**

Run: `go test ./core/tools/ -run TestHash -v`
Expected: all PASS

- [ ] **Step 3: Commit**

```bash
git add core/tools/hash.go core/tools/hash_test.go
git commit -m "feat: add hash generator tool with tests"
```

---

## Chunk 6: URL Tool

### Task 11: URL — tests

**Files:**
- Create: `core/tools/url_test.go`

- [ ] **Step 1: Write URL tests**

```go
package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLEncodeComponent(t *testing.T) {
	result := URLEncode("key=value&foo=bar", true)
	assert.Empty(t, result.Error)
	assert.Equal(t, "key%3Dvalue%26foo%3Dbar", result.Output)
}

func TestURLEncodePathSafe(t *testing.T) {
	result := URLEncode("hello world", false)
	assert.Empty(t, result.Error)
	assert.Equal(t, "hello%20world", result.Output)
}

func TestURLEncodeEmpty(t *testing.T) {
	result := URLEncode("", false)
	assert.Empty(t, result.Error)
	assert.Equal(t, "", result.Output)
}

func TestURLDecode(t *testing.T) {
	result := URLDecode("hello%20world")
	assert.Empty(t, result.Error)
	assert.Equal(t, "hello world", result.Output)
}

func TestURLDecodeComponentPlus(t *testing.T) {
	result := URLDecode("hello+world")
	assert.Empty(t, result.Error)
	assert.Equal(t, "hello world", result.Output)
}

func TestURLDecodeInvalid(t *testing.T) {
	result := URLDecode("%zz")
	assert.NotEmpty(t, result.Error)
}

func TestURLDecodeEmpty(t *testing.T) {
	result := URLDecode("")
	assert.Empty(t, result.Error)
	assert.Equal(t, "", result.Output)
}

func TestURLParse(t *testing.T) {
	result := URLParse("https://example.com:8080/path?q=hello&page=1#section")
	assert.Empty(t, result.Error)
	assert.Equal(t, "https", result.Scheme)
	assert.Equal(t, "example.com", result.Host)
	assert.Equal(t, "8080", result.Port)
	assert.Equal(t, "/path", result.Path)
	assert.Equal(t, "q=hello&page=1", result.Query)
	assert.Equal(t, "section", result.Fragment)
	assert.Equal(t, []string{"hello"}, result.Params["q"])
	assert.Equal(t, []string{"1"}, result.Params["page"])
}

func TestURLParseSimple(t *testing.T) {
	result := URLParse("https://example.com")
	assert.Empty(t, result.Error)
	assert.Equal(t, "https", result.Scheme)
	assert.Equal(t, "example.com", result.Host)
}

func TestURLParseInvalid(t *testing.T) {
	result := URLParse("://bad")
	assert.NotEmpty(t, result.Error)
}

func TestURLParseEmpty(t *testing.T) {
	result := URLParse("")
	assert.NotEmpty(t, result.Error)
}

func TestURLToolDetection(t *testing.T) {
	tool := URLTool{}
	assert.True(t, tool.DetectFromClipboard("https://example.com"))
	assert.True(t, tool.DetectFromClipboard("http://localhost:8080"))
	assert.False(t, tool.DetectFromClipboard("ftp://files.example.com"))
	assert.False(t, tool.DetectFromClipboard("not a url"))
	assert.False(t, tool.DetectFromClipboard(""))
}

func TestURLToolMetadata(t *testing.T) {
	tool := URLTool{}
	assert.Equal(t, "url", tool.ID())
	assert.Equal(t, "Encoders", tool.Category())
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./core/tools/ -run TestURL -v`
Expected: FAIL — `URLEncode` undefined

### Task 12: URL — implementation

**Files:**
- Create: `core/tools/url.go`

- [ ] **Step 1: Write URL implementation**

```go
package tools

import (
	"fmt"
	"net/url"
	"strings"
)

// URLParseResult holds the parsed components of a URL.
type URLParseResult struct {
	Scheme   string
	Host     string
	Port     string
	Path     string
	Query    string
	Fragment string
	Params   map[string][]string
	Output   string
	Error    string
}

// URLTool implements the Tool interface for URL encoding/decoding.
type URLTool struct{}

func (URLTool) Name() string        { return "URL Encode / Decode / Parse" }
func (URLTool) ID() string          { return "url" }
func (URLTool) Description() string { return "Encode, decode, or parse URLs" }
func (URLTool) Category() string    { return "Encoders" }
func (URLTool) Keywords() []string  { return []string{"url", "encode", "decode", "parse", "percent", "uri"} }

func (URLTool) DetectFromClipboard(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

// URLEncode percent-encodes a string.
// component=true uses QueryEscape (spaces → +, encodes /, ?, &, =).
// component=false uses PathEscape (spaces → %20, preserves structure).
func URLEncode(input string, component bool) Result {
	if input == "" {
		return Result{Output: ""}
	}
	if component {
		return Result{Output: url.QueryEscape(input)}
	}
	return Result{Output: url.PathEscape(input)}
}

// URLDecode decodes a percent-encoded string.
func URLDecode(input string) Result {
	if input == "" {
		return Result{Output: ""}
	}
	decoded, err := url.QueryUnescape(input)
	if err != nil {
		return Result{Error: "invalid URL encoding: " + err.Error()}
	}
	return Result{Output: decoded}
}

// URLParse breaks a URL into its components.
func URLParse(input string) URLParseResult {
	if input == "" {
		return URLParseResult{Error: "empty input"}
	}
	parsed, err := url.Parse(input)
	if err != nil {
		return URLParseResult{Error: "invalid URL: " + err.Error()}
	}
	if parsed.Scheme == "" {
		return URLParseResult{Error: "invalid URL: missing scheme"}
	}

	host := parsed.Hostname()
	port := parsed.Port()

	var sb strings.Builder
	fmt.Fprintf(&sb, "Scheme:    %s\n", parsed.Scheme)
	fmt.Fprintf(&sb, "Host:      %s\n", host)
	if port != "" {
		fmt.Fprintf(&sb, "Port:      %s\n", port)
	}
	if parsed.Path != "" {
		fmt.Fprintf(&sb, "Path:      %s\n", parsed.Path)
	}
	if parsed.RawQuery != "" {
		fmt.Fprintf(&sb, "Query:     %s\n", parsed.RawQuery)
	}
	if parsed.Fragment != "" {
		fmt.Fprintf(&sb, "Fragment:  %s\n", parsed.Fragment)
	}

	params := parsed.Query()
	if len(params) > 0 {
		sb.WriteString("Params:\n")
		for k, v := range params {
			fmt.Fprintf(&sb, "  %s = %s\n", k, strings.Join(v, ", "))
		}
	}

	return URLParseResult{
		Scheme:   parsed.Scheme,
		Host:     host,
		Port:     port,
		Path:     parsed.Path,
		Query:    parsed.RawQuery,
		Fragment: parsed.Fragment,
		Params:   params,
		Output:   strings.TrimSpace(sb.String()),
	}
}
```

- [ ] **Step 2: Run tests to verify they pass**

Run: `go test ./core/tools/ -run TestURL -v`
Expected: all PASS

- [ ] **Step 3: Commit**

```bash
git add core/tools/url.go core/tools/url_test.go
git commit -m "feat: add URL encode/decode/parse tool with tests"
```

---

## Chunk 7: UUID Tool

### Task 13: UUID — tests

**Files:**
- Create: `core/tools/uuid_test.go`

- [ ] **Step 1: Write UUID tests**

```go
package tools

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func TestUUIDGenerateV4(t *testing.T) {
	result := UUIDGenerate(4, false, false)
	assert.Empty(t, result.Error)
	assert.Regexp(t, uuidRegex, result.Output)
}

func TestUUIDGenerateV7(t *testing.T) {
	result := UUIDGenerate(7, false, false)
	assert.Empty(t, result.Error)
	assert.Regexp(t, uuidRegex, result.Output)
}

func TestUUIDGenerateUppercase(t *testing.T) {
	result := UUIDGenerate(4, true, false)
	assert.Empty(t, result.Error)
	assert.Regexp(t, `^[0-9A-F\-]+$`, result.Output)
}

func TestUUIDGenerateNoHyphens(t *testing.T) {
	result := UUIDGenerate(4, false, true)
	assert.Empty(t, result.Error)
	assert.NotContains(t, result.Output, "-")
	assert.Len(t, result.Output, 32)
}

func TestUUIDGenerateUnsupportedVersion(t *testing.T) {
	result := UUIDGenerate(1, false, false)
	assert.NotEmpty(t, result.Error)
	assert.Contains(t, result.Error, "unsupported UUID version")
}

func TestUUIDValidate(t *testing.T) {
	result := UUIDValidate("550e8400-e29b-41d4-a716-446655440000")
	assert.Empty(t, result.Error)
	assert.Contains(t, result.Output, "valid")
	assert.Contains(t, result.Output, "version 4")
}

func TestUUIDValidateInvalid(t *testing.T) {
	result := UUIDValidate("not-a-uuid")
	assert.NotEmpty(t, result.Error)
}

func TestUUIDValidateEmpty(t *testing.T) {
	result := UUIDValidate("")
	assert.NotEmpty(t, result.Error)
}

func TestUUIDParse(t *testing.T) {
	result := UUIDParse("550e8400-e29b-41d4-a716-446655440000")
	assert.Empty(t, result.Error)
	assert.Equal(t, 4, result.Version)
	assert.Equal(t, "RFC 4122", result.Variant)
	assert.NotEmpty(t, result.Output)
}

func TestUUIDParseInvalid(t *testing.T) {
	result := UUIDParse("garbage")
	assert.NotEmpty(t, result.Error)
}

func TestUUIDParseEmpty(t *testing.T) {
	result := UUIDParse("")
	assert.NotEmpty(t, result.Error)
}

func TestUUIDToolDetection(t *testing.T) {
	tool := UUIDTool{}
	assert.True(t, tool.DetectFromClipboard("550e8400-e29b-41d4-a716-446655440000"))
	assert.True(t, tool.DetectFromClipboard("550E8400-E29B-41D4-A716-446655440000"))
	assert.False(t, tool.DetectFromClipboard("not-a-uuid"))
	assert.False(t, tool.DetectFromClipboard(""))
}

func TestUUIDToolMetadata(t *testing.T) {
	tool := UUIDTool{}
	assert.Equal(t, "uuid", tool.ID())
	assert.Equal(t, "Generators", tool.Category())
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./core/tools/ -run TestUUID -v`
Expected: FAIL — `UUIDGenerate` undefined

### Task 14: UUID — implementation

**Files:**
- Create: `core/tools/uuid.go`

- [ ] **Step 1: Write UUID implementation**

```go
package tools

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// UUIDParseResult holds the parsed components of a UUID.
type UUIDParseResult struct {
	UUID      string
	Version   int
	Variant   string
	Timestamp string
	Output    string
	Error     string
}

// UUIDTool implements the Tool interface for UUID operations.
type UUIDTool struct{}

func (UUIDTool) Name() string        { return "UUID Generate / Validate / Parse" }
func (UUIDTool) ID() string          { return "uuid" }
func (UUIDTool) Description() string { return "Generate, validate, or parse UUIDs" }
func (UUIDTool) Category() string    { return "Generators" }
func (UUIDTool) Keywords() []string  { return []string{"uuid", "guid", "generate", "v4", "v7", "unique"} }

var uuidDetectRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func (UUIDTool) DetectFromClipboard(s string) bool {
	return uuidDetectRegex.MatchString(strings.TrimSpace(s))
}

// UUIDGenerate creates a new UUID of the specified version.
func UUIDGenerate(version int, uppercase bool, noHyphens bool) Result {
	var u uuid.UUID
	var err error
	switch version {
	case 4:
		u, err = uuid.NewRandom()
	case 7:
		u, err = uuid.NewV7()
	default:
		return Result{Error: fmt.Sprintf("unsupported UUID version: %d (supported: 4, 7)", version)}
	}
	if err != nil {
		return Result{Error: "generation error: " + err.Error()}
	}
	output := u.String()
	if noHyphens {
		output = strings.ReplaceAll(output, "-", "")
	}
	if uppercase {
		output = strings.ToUpper(output)
	}
	return Result{Output: output}
}

// UUIDValidate checks if a string is a valid UUID.
func UUIDValidate(input string) Result {
	if input == "" {
		return Result{Error: "empty input"}
	}
	parsed, err := uuid.Parse(strings.TrimSpace(input))
	if err != nil {
		return Result{Error: "invalid UUID: " + err.Error()}
	}
	return Result{Output: fmt.Sprintf("valid (version %d)", parsed.Version())}
}

// UUIDParse parses a UUID and returns its components.
func UUIDParse(input string) UUIDParseResult {
	if input == "" {
		return UUIDParseResult{Error: "empty input"}
	}
	parsed, err := uuid.Parse(strings.TrimSpace(input))
	if err != nil {
		return UUIDParseResult{Error: "invalid UUID: " + err.Error()}
	}

	version := int(parsed.Version())
	variant := variantString(parsed.Variant())

	var sb strings.Builder
	fmt.Fprintf(&sb, "UUID:      %s\n", parsed.String())
	fmt.Fprintf(&sb, "Version:   %d\n", version)
	fmt.Fprintf(&sb, "Variant:   %s", variant)

	var timestamp string
	if version == 7 {
		// UUID v7 embeds a Unix timestamp in the first 48 bits
		sec, _ := parsed.Time().UnixTime()
		timestamp = fmt.Sprintf("%d", sec)
		fmt.Fprintf(&sb, "\nTimestamp: %s", timestamp)
	}

	return UUIDParseResult{
		UUID:      parsed.String(),
		Version:   version,
		Variant:   variant,
		Timestamp: timestamp,
		Output:    sb.String(),
	}
}

func variantString(v uuid.Variant) string {
	switch v {
	case uuid.RFC4122:
		return "RFC 4122"
	case uuid.Reserved:
		return "Reserved"
	case uuid.Microsoft:
		return "Microsoft"
	case uuid.Future:
		return "Future"
	default:
		return "Unknown"
	}
}
```

- [ ] **Step 2: Run tests to verify they pass**

Run: `go test ./core/tools/ -run TestUUID -v`
Expected: all PASS

- [ ] **Step 3: Commit**

```bash
git add core/tools/uuid.go core/tools/uuid_test.go
git commit -m "feat: add UUID generate/validate/parse tool with tests"
```

---

## Chunk 8: Registry

### Task 15: Registry — tests

**Files:**
- Create: `core/registry/registry_test.go`

- [ ] **Step 1: Write registry tests**

```go
package registry

import (
	"testing"

	"github.com/StanMarek/forge/core/tools"
	"github.com/stretchr/testify/assert"
)

func TestRegistryRegisterAndAll(t *testing.T) {
	r := New()
	r.Register(tools.Base64Tool{})
	r.Register(tools.JWTTool{})

	all := r.All()
	assert.Len(t, all, 2)
}

func TestRegistryByID(t *testing.T) {
	r := New()
	r.Register(tools.Base64Tool{})

	tool, ok := r.ByID("base64")
	assert.True(t, ok)
	assert.Equal(t, "base64", tool.ID())

	_, ok = r.ByID("nonexistent")
	assert.False(t, ok)
}

func TestRegistryByCategory(t *testing.T) {
	r := New()
	r.Register(tools.Base64Tool{})
	r.Register(tools.JWTTool{})
	r.Register(tools.HashTool{})

	encoders := r.ByCategory("Encoders")
	assert.Len(t, encoders, 2)

	generators := r.ByCategory("Generators")
	assert.Len(t, generators, 1)

	empty := r.ByCategory("Nonexistent")
	assert.Len(t, empty, 0)
}

func TestRegistrySearch(t *testing.T) {
	r := New()
	r.Register(tools.Base64Tool{})
	r.Register(tools.JWTTool{})
	r.Register(tools.JSONTool{})

	results := r.Search("minify")
	assert.Len(t, results, 1)
	assert.Equal(t, "json", results[0].ID())

	results = r.Search("encode")
	assert.GreaterOrEqual(t, len(results), 1)

	results = r.Search("zzzzz")
	assert.Len(t, results, 0)
}

func TestRegistrySearchCaseInsensitive(t *testing.T) {
	r := New()
	r.Register(tools.Base64Tool{})

	results := r.Search("BASE64")
	assert.Len(t, results, 1)
}

func TestRegistryDetect(t *testing.T) {
	r := Default()

	// JWT should be detected first (highest priority)
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"
	matches := r.Detect(jwt)
	assert.NotEmpty(t, matches)
	assert.Equal(t, "jwt", matches[0].ID())

	// UUID should be detected
	matches = r.Detect("550e8400-e29b-41d4-a716-446655440000")
	assert.NotEmpty(t, matches)
	assert.Equal(t, "uuid", matches[0].ID())

	// URL should be detected
	matches = r.Detect("https://example.com")
	assert.NotEmpty(t, matches)
	assert.Equal(t, "url", matches[0].ID())

	// JSON should be detected
	matches = r.Detect(`{"key":"value"}`)
	assert.NotEmpty(t, matches)
	assert.Equal(t, "json", matches[0].ID())

	// No match
	matches = r.Detect("just plain text")
	assert.Empty(t, matches)
}

func TestRegistryDetectPriority(t *testing.T) {
	r := Default()

	// A valid JSON string that is also valid base64 — JSON should win
	matches := r.Detect(`{"a":1}`)
	if len(matches) > 1 {
		assert.Equal(t, "json", matches[0].ID())
	}
}

func TestDefaultRegistry(t *testing.T) {
	r := Default()
	all := r.All()
	assert.Len(t, all, 6)

	// Verify all 6 tools are registered
	ids := make(map[string]bool)
	for _, tool := range all {
		ids[tool.ID()] = true
	}
	assert.True(t, ids["base64"])
	assert.True(t, ids["jwt"])
	assert.True(t, ids["json"])
	assert.True(t, ids["hash"])
	assert.True(t, ids["url"])
	assert.True(t, ids["uuid"])
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./core/registry/ -v`
Expected: FAIL — package doesn't exist yet

### Task 16: Registry — implementation

**Files:**
- Create: `core/registry/registry.go`
- Create: `core/registry/defaults.go`

- [ ] **Step 1: Write registry.go**

```go
package registry

import (
	"sort"
	"strings"

	"github.com/StanMarek/forge/core/tools"
)

// detectionPriority defines the order for smart clipboard detection.
// Lower index = higher priority.
var detectionPriority = map[string]int{
	"jwt":         0,
	"uuid":        1,
	"url":         2,
	"json":        3,
	"base64":      4,
	"timestamp":   5,
	"html-entity": 6,
	"number-base": 7,
}

// Registry holds registered tools and provides lookup, search, and detection.
type Registry struct {
	tools  map[string]tools.Tool
	order  []string // insertion order
}

// New creates an empty Registry.
func New() *Registry {
	return &Registry{
		tools: make(map[string]tools.Tool),
	}
}

// Register adds a tool to the registry.
func (r *Registry) Register(tool tools.Tool) {
	id := tool.ID()
	if _, exists := r.tools[id]; !exists {
		r.order = append(r.order, id)
	}
	r.tools[id] = tool
}

// All returns all registered tools in registration order.
func (r *Registry) All() []tools.Tool {
	result := make([]tools.Tool, 0, len(r.order))
	for _, id := range r.order {
		result = append(result, r.tools[id])
	}
	return result
}

// ByID returns a tool by its ID.
func (r *Registry) ByID(id string) (tools.Tool, bool) {
	tool, ok := r.tools[id]
	return tool, ok
}

// ByCategory returns all tools in a given category.
func (r *Registry) ByCategory(category string) []tools.Tool {
	var result []tools.Tool
	for _, id := range r.order {
		if r.tools[id].Category() == category {
			result = append(result, r.tools[id])
		}
	}
	return result
}

// Search returns tools matching a query against name, ID, and keywords.
func (r *Registry) Search(query string) []tools.Tool {
	q := strings.ToLower(query)
	var result []tools.Tool
	for _, id := range r.order {
		tool := r.tools[id]
		if matchesTool(tool, q) {
			result = append(result, tool)
		}
	}
	return result
}

// Detect returns tools whose DetectFromClipboard returns true,
// sorted by detection priority (JWT > UUID > URL > JSON > Base64 > ...).
func (r *Registry) Detect(clipboard string) []tools.Tool {
	var matches []tools.Tool
	for _, id := range r.order {
		tool := r.tools[id]
		if tool.DetectFromClipboard(clipboard) {
			matches = append(matches, tool)
		}
	}
	sort.SliceStable(matches, func(i, j int) bool {
		pi := priorityOf(matches[i].ID())
		pj := priorityOf(matches[j].ID())
		return pi < pj
	})
	return matches
}

func matchesTool(tool tools.Tool, query string) bool {
	if strings.Contains(strings.ToLower(tool.Name()), query) {
		return true
	}
	if strings.Contains(strings.ToLower(tool.ID()), query) {
		return true
	}
	for _, kw := range tool.Keywords() {
		if strings.Contains(strings.ToLower(kw), query) {
			return true
		}
	}
	return false
}

func priorityOf(id string) int {
	if p, ok := detectionPriority[id]; ok {
		return p
	}
	return 999
}
```

- [ ] **Step 2: Write defaults.go**

```go
package registry

import "github.com/StanMarek/forge/core/tools"

// Default creates a registry pre-loaded with all Tier-1 tools.
func Default() *Registry {
	r := New()
	r.Register(tools.Base64Tool{})
	r.Register(tools.JWTTool{})
	r.Register(tools.JSONTool{})
	r.Register(tools.HashTool{})
	r.Register(tools.URLTool{})
	r.Register(tools.UUIDTool{})
	return r
}
```

- [ ] **Step 3: Run tests to verify they pass**

Run: `go test ./core/... -v`
Expected: ALL tests pass across tools and registry

- [ ] **Step 4: Commit**

```bash
git add core/registry/
git commit -m "feat: add tool registry with search, detection, and defaults"
```

---

## Chunk 9: Final Verification

### Task 17: Full test suite + cleanup

- [ ] **Step 1: Run complete test suite**

Run: `go test -v -count=1 ./core/...`
Expected: all tests pass, zero failures

- [ ] **Step 2: Run go vet**

Run: `go vet ./core/...`
Expected: no issues

- [ ] **Step 3: Verify no forbidden imports**

Run: `grep -r '"github.com/StanMarek/forge/ui\|"github.com/StanMarek/forge/cmd' core/`
Expected: no matches (core must never import ui or cmd)

- [ ] **Step 4: Commit any cleanup**

Only if changes were needed:
```bash
git add -A
git commit -m "chore: final cleanup and verification for core tools foundation"
```

---

## Task Dependency Summary

Tasks 1-2 are sequential (must complete before anything else).
Tasks 3-14 (the 6 tools) are **independent** — they can be built in parallel by subagents.
Tasks 15-16 (registry) depend on all 6 tools being complete.
Task 17 (verification) depends on everything.

```
Task 1 (go.mod) → Task 2 (tool.go)
                    ├── Tasks 3-4   (base64)  ─┐
                    ├── Tasks 5-6   (jwt)     ─┤
                    ├── Tasks 7-8   (json)    ─┤
                    ├── Tasks 9-10  (hash)    ─├── Tasks 15-16 (registry) → Task 17 (verify)
                    ├── Tasks 11-12 (url)     ─┤
                    └── Tasks 13-14 (uuid)    ─┘
```
