# Core Tools Foundation — Design Spec

**Date:** 2026-03-15
**Status:** Approved
**Scope:** `core/tools/`, `core/registry/`, `go.mod`

---

## Goal

Implement the pure business logic layer for Forge's 6 Tier-1 tools. This is the foundation every UI surface depends on. No I/O, no global state, no UI imports.

## Module

- **Path:** `github.com/StanMarek/forge`
- **Go version:** 1.22+
- **External deps:** `github.com/google/uuid`, `github.com/stretchr/testify` (test only)

## Package: `core/tools/tool.go`

### Tool Interface

```go
type Tool interface {
    Name() string                      // Display name: "Base64 Encode / Decode"
    ID() string                        // URL/CLI slug: "base64"
    Description() string               // One-liner
    Category() string                  // Grouping: "Encoders", "Formatters", etc.
    Keywords() []string                // Search terms
    DetectFromClipboard(s string) bool // Smart detection predicate
}
```

### Result Type

```go
type Result struct {
    Output string
    Error  string
}
```

Success is `Error == ""`. Each tool may define additional result structs (e.g., `JWTDecodeResult`, `URLParseResult`) that embed or extend this pattern, but all must carry `Output` and `Error` fields.

## Tool Implementations

### 1. base64 (`core/tools/base64.go`)

**Category:** Encoders

```go
func Base64Encode(input string, urlSafe bool, noPadding bool) Result
func Base64Decode(input string, urlSafe bool) Result
```

- Uses `encoding/base64` stdlib.
- URL-safe uses RFC 4648 section 5 alphabet.
- `noPadding` strips trailing `=` characters.
- Detection: standard alphabet only — regex `^[A-Za-z0-9+/=]+$`, min 4 chars, length divisible by 4. URL-safe base64 is not detected (too ambiguous without padding).

### 2. jwt (`core/tools/jwt.go`)

**Category:** Encoders

```go
type JWTDecodeResult struct {
    Header    string
    Payload   string
    Signature string
    Output    string // Formatted combined output (header + payload + signature)
    Error     string
}

func JWTDecode(token string) JWTDecodeResult
func JWTValidate(token string) Result
```

- Splits on `.`, base64url-decodes header and payload segments.
- Does NOT validate signatures — decode only.
- `JWTValidate` checks structural validity: 3 dot-separated base64url segments, valid JSON in header and payload.
- Detection: matches `xxxxx.yyyyy.zzzzz` (3 dot-separated non-empty segments).

### 3. json (`core/tools/json.go`)

**Category:** Formatters

```go
func JSONFormat(input string, indent int, sortKeys bool, useTabs bool) Result
func JSONMinify(input string) Result
func JSONValidate(input string) Result
```

- Uses `encoding/json` stdlib.
- `JSONFormat` pretty-prints with configurable indent. Default indent: 2 spaces.
- `sortKeys` unmarshals to `map[string]interface{}` and re-marshals. Keys are sorted recursively at all nesting levels (Go's `json.Marshal` sorts map keys alphabetically).
- `JSONValidate` returns `Output: "valid"` or `Error` with position info.
- Detection: `json.Valid([]byte(input))`.

### 4. hash (`core/tools/hash.go`)

**Category:** Generators

```go
func Hash(input string, algorithm string, uppercase bool) Result
```

- Supported algorithms: `md5`, `sha1`, `sha256`, `sha512`.
- Uses `crypto/md5`, `crypto/sha1`, `crypto/sha256`, `crypto/sha512` stdlib.
- Output is hex-encoded digest.
- Returns `Error` for unknown algorithm.
- Detection: `false` always (hashing is one-way).

### 5. url (`core/tools/url.go`)

**Category:** Encoders

```go
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

func URLEncode(input string, component bool) Result
func URLDecode(input string) Result
func URLParse(input string) URLParseResult
```

- Uses `net/url` stdlib.
- `component` mode uses `url.QueryEscape` — encodes everything including `/`, `?`, `&`, `=`. Spaces become `+`.
- Non-component mode uses `url.PathEscape` — preserves URL structure characters. Spaces become `%20`.
- Detection: starts with `http://` or `https://`.

### 6. uuid (`core/tools/uuid.go`)

**Category:** Generators

```go
type UUIDParseResult struct {
    UUID      string
    Version   int
    Variant   string
    Timestamp string // v7 only, empty for others
    Output    string
    Error     string
}

func UUIDGenerate(version int, uppercase bool, noHyphens bool) Result
func UUIDValidate(input string) Result
func UUIDParse(input string) UUIDParseResult
```

- Uses `github.com/google/uuid`.
- Supports v4 (random) and v7 (time-ordered). Returns `Error: "unsupported UUID version: N"` for any other version.
- `UUIDValidate` returns `Output: "valid (version N)"` or `Error`.
- Generation of multiple UUIDs (`--count`) is handled by the CLI layer via repeated calls to `UUIDGenerate`.
- Detection: matches UUID regex `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`.

## Package: `core/registry/`

### registry.go

```go
type Registry struct { /* map[string]Tool */ }

func New() *Registry
func (r *Registry) Register(tool Tool)
func (r *Registry) All() []Tool
func (r *Registry) ByID(id string) (Tool, bool)
func (r *Registry) ByCategory(category string) []Tool
func (r *Registry) Search(query string) []Tool
func (r *Registry) Detect(clipboard string) []Tool
```

- `Search` matches against `Name()`, `ID()`, `Keywords()` (case-insensitive substring).
- `Detect` calls `DetectFromClipboard` on all registered tools, returns matches sorted by a hardcoded priority map. Tier-1 priority: JWT > UUID > URL > JSON > Base64. Hash always returns `false`. Full priority order (for future Tier-2 tools): JWT > UUID > URL > JSON > Base64 > Timestamp > HTML entity > Number base. Tools not in the priority map sort after all prioritized tools, in registration order.

### defaults.go

```go
func Default() *Registry
```

Returns a registry pre-loaded with all 6 Tier-1 tools.

## Testing Strategy

- One `_test.go` per tool file.
- Use `github.com/stretchr/testify/assert` for assertions.
- Test cases cover: happy path, edge cases, error cases, empty input, unicode input.
- Registry tests: registration, lookup, search, detection priority.
- All tests must pass with `go test ./core/...`.

## File Manifest

```
go.mod
go.sum
core/tools/tool.go
core/tools/base64.go
core/tools/base64_test.go
core/tools/jwt.go
core/tools/jwt_test.go
core/tools/json.go
core/tools/json_test.go
core/tools/hash.go
core/tools/hash_test.go
core/tools/url.go
core/tools/url_test.go
core/tools/uuid.go
core/tools/uuid_test.go
core/registry/registry.go
core/registry/registry_test.go
core/registry/defaults.go
```

## Out of Scope

- **`core/detection/`** — The clipboard detection engine (polling, background detection) is out of scope. It depends on `internal/clipboard/` which involves I/O. The registry's `Detect()` method provides the pure detection logic; the detection engine wraps it with clipboard polling in a later milestone.
- **`cmd/`** — CLI commands are a separate milestone.
- **`ui/`** — All UI surfaces are separate milestones.

## Constraints

- `core/` must NEVER import from `ui/` or `cmd/`.
- All tool functions are pure: no I/O, no global state, no side effects.
- Tool functions return plain Go types, never formatted output.
- External dependencies limited to `google/uuid` (runtime) and `testify` (test).
