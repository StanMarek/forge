# CLI Surface Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the Cobra CLI surface so all 6 Tier-1 tools are usable from the command line with stdin/pipe support.

**Architecture:** `main.go` calls `cmd.Execute()`. `cmd/root.go` registers subcommands. Each tool gets one file with a parent command and subcommands. `internal/stdin/` handles pipe detection. `internal/version/` holds build metadata.

**Tech Stack:** Go 1.25, `github.com/spf13/cobra`, existing `core/tools/` functions

**Spec:** `docs/superpowers/specs/2026-03-15-cli-surface-design.md`

---

## File Structure

| File | Responsibility |
|------|---------------|
| `internal/version/version.go` | Build version variables (set via ldflags) |
| `internal/stdin/stdin.go` | Pipe detection and stdin reading |
| `internal/stdin/stdin_test.go` | Tests for stdin reading |
| `cmd/helpers.go` | Shared `resolveInput` and `exitWithError` helpers |
| `cmd/root.go` | Root cobra command, subcommand registration, Execute() |
| `cmd/version.go` | `forge version` command |
| `cmd/base64.go` | `forge base64 encode/decode` commands |
| `cmd/jwt.go` | `forge jwt decode/validate` commands |
| `cmd/json.go` | `forge json format/minify/validate` commands |
| `cmd/hash.go` | `forge hash <algo>` command |
| `cmd/url.go` | `forge url encode/decode/parse` commands |
| `cmd/uuid.go` | `forge uuid generate/validate/parse` commands |
| `main.go` | Entrypoint — calls cmd.Execute() |

---

## Chunk 1: Bootstrap — internal packages + main + root

### Task 1: internal/version

**Files:**
- Create: `internal/version/version.go`

- [ ] **Step 1: Write version.go**

```go
package version

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)
```

- [ ] **Step 2: Verify it compiles**

Run: `go build ./internal/version/`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add internal/version/version.go
git commit -m "Add internal/version package with build metadata variables"
```

### Task 2: internal/stdin

**Files:**
- Create: `internal/stdin/stdin.go`
- Create: `internal/stdin/stdin_test.go`

- [ ] **Step 1: Write stdin_test.go**

```go
package stdin

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFromPipe(t *testing.T) {
	// Create a pipe to simulate stdin
	r, w, err := os.Pipe()
	assert.NoError(t, err)

	_, err = w.WriteString("hello world\n")
	assert.NoError(t, err)
	w.Close()

	result, err := readFrom(r)
	assert.NoError(t, err)
	assert.Equal(t, "hello world", result)
}

func TestReadFromPipeMultiline(t *testing.T) {
	r, w, err := os.Pipe()
	assert.NoError(t, err)

	_, err = w.WriteString("line1\nline2\nline3\n")
	assert.NoError(t, err)
	w.Close()

	result, err := readFrom(r)
	assert.NoError(t, err)
	assert.Equal(t, "line1\nline2\nline3", result)
}

func TestReadFromEmpty(t *testing.T) {
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	w.Close()

	result, err := readFrom(r)
	assert.NoError(t, err)
	assert.Equal(t, "", result)
}

func TestReadFromPreservesInternalWhitespace(t *testing.T) {
	r, w, err := os.Pipe()
	assert.NoError(t, err)

	_, err = w.WriteString("  hello  world  \n")
	assert.NoError(t, err)
	w.Close()

	result, err := readFrom(r)
	assert.NoError(t, err)
	assert.Equal(t, "  hello  world  ", result)
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/stdin/ -v`
Expected: FAIL — `readFrom` undefined

- [ ] **Step 3: Write stdin.go**

```go
package stdin

import (
	"errors"
	"io"
	"os"
	"strings"
)

// Read reads all input from stdin if it is piped.
// Returns an error if stdin is a terminal (interactive).
func Read() (string, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return "", errors.New("no input provided")
	}
	if info.Mode()&os.ModeCharDevice != 0 {
		return "", errors.New("no input provided")
	}
	return readFrom(os.Stdin)
}

// readFrom reads all content from a reader, trimming the trailing newline.
// Extracted for testability.
func readFrom(r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(data), "\n"), nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/stdin/ -v`
Expected: all PASS

- [ ] **Step 5: Commit**

```bash
git add internal/stdin/
git commit -m "Add internal/stdin package with pipe detection and reading"
```

### Task 3: cmd/helpers.go + cmd/root.go + cmd/version.go + main.go

**Files:**
- Create: `cmd/helpers.go`
- Create: `cmd/root.go`
- Create: `cmd/version.go`
- Create: `main.go`

- [ ] **Step 1: Write cmd/helpers.go**

```go
package cmd

import (
	"fmt"
	"os"

	"github.com/StanMarek/forge/internal/stdin"
)

// resolveInput gets input from args or stdin.
// If args[0] exists and is not "-", use it. Otherwise read stdin.
func resolveInput(args []string) (string, error) {
	if len(args) > 0 && args[0] != "-" {
		return args[0], nil
	}
	return stdin.Read()
}

// exitWithError prints an error to stderr and exits with code 1.
func exitWithError(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
```

- [ ] **Step 2: Write cmd/root.go**

```go
package cmd

import (
	"os"

	"github.com/StanMarek/forge/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "forge",
	Short: "A developer's workbench for the terminal, browser, and desktop",
	// TODO: launch TUI when no subcommand is given
}

func init() {
	rootCmd.Version = version.Version
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(base64Cmd)
	rootCmd.AddCommand(jwtCmd)
	rootCmd.AddCommand(jsonCmd)
	rootCmd.AddCommand(hashCmd)
	rootCmd.AddCommand(urlCmd)
	rootCmd.AddCommand(uuidCmd)
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(2)
	}
}
```

Note: This file references commands defined in later tasks. It will NOT compile until all command files exist. That's fine — we'll build them in the next tasks and verify compilation at the end.

- [ ] **Step 3: Write cmd/version.go**

```go
package cmd

import (
	"fmt"
	"runtime"

	"github.com/StanMarek/forge/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("forge %s\n", version.Version)
		fmt.Printf("commit:  %s\n", version.Commit)
		fmt.Printf("built:   %s\n", version.Date)
		fmt.Printf("go:      %s\n", runtime.Version())
	},
}
```

- [ ] **Step 4: Write main.go**

```go
package main

import "github.com/StanMarek/forge/cmd"

func main() {
	cmd.Execute()
}
```

- [ ] **Step 5: Commit (without verifying build — commands not yet defined)**

```bash
git add internal/version/version.go cmd/helpers.go cmd/root.go cmd/version.go main.go
git commit -m "Add CLI bootstrap: main.go, root command, version command, helpers"
```

---

## Chunk 2: Tool commands (all 6, parallelizable)

All 6 tool tasks below are **independent** — they can be built in parallel by subagents. Each creates one file in `cmd/`. After all 6 exist, the project compiles.

### Task 4: cmd/base64.go

**Files:**
- Create: `cmd/base64.go`

- [ ] **Step 1: Write base64.go**

```go
package cmd

import (
	"fmt"

	"github.com/StanMarek/forge/core/tools"
	"github.com/spf13/cobra"
)

var base64Cmd = &cobra.Command{
	Use:   "base64",
	Short: "Encode or decode Base64 strings",
}

var base64EncodeCmd = &cobra.Command{
	Use:   "encode [input]",
	Short: "Encode input as Base64",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		urlSafe, _ := cmd.Flags().GetBool("url-safe")
		noPadding, _ := cmd.Flags().GetBool("no-padding")
		result := tools.Base64Encode(input, urlSafe, noPadding)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

var base64DecodeCmd = &cobra.Command{
	Use:   "decode [input]",
	Short: "Decode a Base64-encoded string",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		urlSafe, _ := cmd.Flags().GetBool("url-safe")
		result := tools.Base64Decode(input, urlSafe)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

func init() {
	base64EncodeCmd.Flags().Bool("url-safe", false, "Use URL-safe Base64 alphabet (RFC 4648 §5)")
	base64EncodeCmd.Flags().Bool("no-padding", false, "Omit trailing = padding characters")
	base64DecodeCmd.Flags().Bool("url-safe", false, "Expect URL-safe Base64 alphabet")
	base64Cmd.AddCommand(base64EncodeCmd)
	base64Cmd.AddCommand(base64DecodeCmd)
}
```

- [ ] **Step 2: Commit**

```bash
git add cmd/base64.go
git commit -m "Add base64 CLI command with encode/decode subcommands"
```

### Task 5: cmd/jwt.go

**Files:**
- Create: `cmd/jwt.go`

- [ ] **Step 1: Write jwt.go**

```go
package cmd

import (
	"fmt"

	"github.com/StanMarek/forge/core/tools"
	"github.com/spf13/cobra"
)

var jwtCmd = &cobra.Command{
	Use:   "jwt",
	Short: "Decode and inspect JWT tokens",
}

var jwtDecodeCmd = &cobra.Command{
	Use:   "decode [token]",
	Short: "Decode a JWT token",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.JWTDecode(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}

		headerOnly, _ := cmd.Flags().GetBool("header-only")
		payloadOnly, _ := cmd.Flags().GetBool("payload-only")
		compact, _ := cmd.Flags().GetBool("compact")

		switch {
		case headerOnly:
			header := result.Header
			if compact {
				if m := tools.JSONMinify(header); m.Error == "" {
					header = m.Output
				}
			}
			fmt.Println(header)
		case payloadOnly:
			payload := result.Payload
			if compact {
				if m := tools.JSONMinify(payload); m.Error == "" {
					payload = m.Output
				}
			}
			fmt.Println(payload)
		case compact:
			h := result.Header
			if m := tools.JSONMinify(h); m.Error == "" {
				h = m.Output
			}
			p := result.Payload
			if m := tools.JSONMinify(p); m.Error == "" {
				p = m.Output
			}
			fmt.Printf("--- Header ---\n%s\n--- Payload ---\n%s\n--- Signature ---\n%s\n", h, p, result.Signature)
		default:
			fmt.Println(result.Output)
		}
	},
}

var jwtValidateCmd = &cobra.Command{
	Use:   "validate [token]",
	Short: "Check if a string is a structurally valid JWT",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.JWTValidate(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

func init() {
	jwtDecodeCmd.Flags().Bool("header-only", false, "Print only the header")
	jwtDecodeCmd.Flags().Bool("payload-only", false, "Print only the payload")
	jwtDecodeCmd.Flags().Bool("compact", false, "Print JSON on a single line")
	jwtCmd.AddCommand(jwtDecodeCmd)
	jwtCmd.AddCommand(jwtValidateCmd)
}
```

- [ ] **Step 2: Commit**

```bash
git add cmd/jwt.go
git commit -m "Add JWT CLI command with decode/validate subcommands"
```

### Task 6: cmd/json.go

**Files:**
- Create: `cmd/json.go`

- [ ] **Step 1: Write json.go**

```go
package cmd

import (
	"fmt"

	"github.com/StanMarek/forge/core/tools"
	"github.com/spf13/cobra"
)

var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "Format, minify, or validate JSON",
}

var jsonFormatCmd = &cobra.Command{
	Use:   "format [input]",
	Short: "Pretty-print JSON with indentation",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		indent, _ := cmd.Flags().GetInt("indent")
		tabs, _ := cmd.Flags().GetBool("tabs")
		sortKeys, _ := cmd.Flags().GetBool("sort-keys")
		result := tools.JSONFormat(input, indent, sortKeys, tabs)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

var jsonMinifyCmd = &cobra.Command{
	Use:   "minify [input]",
	Short: "Remove all whitespace from JSON",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.JSONMinify(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

var jsonValidateCmd = &cobra.Command{
	Use:   "validate [input]",
	Short: "Check if input is valid JSON",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.JSONValidate(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

func init() {
	jsonFormatCmd.Flags().Int("indent", 2, "Number of spaces for indentation")
	jsonFormatCmd.Flags().Bool("tabs", false, "Use tabs instead of spaces")
	jsonFormatCmd.Flags().Bool("sort-keys", false, "Sort object keys alphabetically")
	jsonCmd.AddCommand(jsonFormatCmd)
	jsonCmd.AddCommand(jsonMinifyCmd)
	jsonCmd.AddCommand(jsonValidateCmd)
}
```

- [ ] **Step 2: Commit**

```bash
git add cmd/json.go
git commit -m "Add JSON CLI command with format/minify/validate subcommands"
```

### Task 7: cmd/hash.go

**Files:**
- Create: `cmd/hash.go`

- [ ] **Step 1: Write hash.go**

```go
package cmd

import (
	"fmt"
	"os"

	"github.com/StanMarek/forge/core/tools"
	"github.com/StanMarek/forge/internal/stdin"
	"github.com/spf13/cobra"
)

var hashCmd = &cobra.Command{
	Use:   "hash <algorithm> [input]",
	Short: "Generate hash digests",
	Long:  "Generate hash digests. Supported algorithms: md5, sha1, sha256, sha512.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		algorithm := args[0]
		uppercase, _ := cmd.Flags().GetBool("uppercase")
		filePath, _ := cmd.Flags().GetString("file")

		var input string
		var err error

		if filePath != "" {
			data, readErr := os.ReadFile(filePath)
			if readErr != nil {
				exitWithError(fmt.Sprintf("cannot read file: %s", readErr.Error()))
			}
			input = string(data)
		} else if len(args) > 1 && args[1] != "-" {
			input = args[1]
		} else {
			input, err = stdin.Read()
			if err != nil {
				exitWithError(err.Error())
			}
		}

		result := tools.Hash(input, algorithm, uppercase)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

func init() {
	hashCmd.Flags().Bool("uppercase", false, "Output hash in uppercase hex")
	hashCmd.Flags().String("file", "", "Hash the contents of a file instead of a string")
}
```

- [ ] **Step 2: Commit**

```bash
git add cmd/hash.go
git commit -m "Add hash CLI command with algorithm arg and file flag"
```

### Task 8: cmd/url.go

**Files:**
- Create: `cmd/url.go`

- [ ] **Step 1: Write url.go**

```go
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/StanMarek/forge/core/tools"
	"github.com/spf13/cobra"
)

var urlCmd = &cobra.Command{
	Use:   "url",
	Short: "Encode, decode, or parse URLs",
}

var urlEncodeCmd = &cobra.Command{
	Use:   "encode [input]",
	Short: "URL-encode a string",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		component, _ := cmd.Flags().GetBool("component")
		result := tools.URLEncode(input, component)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

var urlDecodeCmd = &cobra.Command{
	Use:   "decode [input]",
	Short: "Decode a URL-encoded string",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.URLDecode(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

var urlParseCmd = &cobra.Command{
	Use:   "parse [input]",
	Short: "Parse a URL into components",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.URLParse(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		asJSON, _ := cmd.Flags().GetBool("json")
		if asJSON {
			out := struct {
				Scheme   string              `json:"scheme"`
				Host     string              `json:"host"`
				Port     string              `json:"port,omitempty"`
				Path     string              `json:"path,omitempty"`
				Query    string              `json:"query,omitempty"`
				Fragment string              `json:"fragment,omitempty"`
				Params   map[string][]string `json:"params,omitempty"`
			}{
				Scheme:   result.Scheme,
				Host:     result.Host,
				Port:     result.Port,
				Path:     result.Path,
				Query:    result.Query,
				Fragment: result.Fragment,
				Params:   result.Params,
			}
			data, _ := json.MarshalIndent(out, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Println(result.Output)
		}
	},
}

func init() {
	urlEncodeCmd.Flags().Bool("component", false, "Encode as URL component (encodes /, ?, &, =)")
	urlParseCmd.Flags().Bool("json", false, "Output as JSON")
	urlCmd.AddCommand(urlEncodeCmd)
	urlCmd.AddCommand(urlDecodeCmd)
	urlCmd.AddCommand(urlParseCmd)
}
```

- [ ] **Step 2: Commit**

```bash
git add cmd/url.go
git commit -m "Add URL CLI command with encode/decode/parse subcommands"
```

### Task 9: cmd/uuid.go

**Files:**
- Create: `cmd/uuid.go`

- [ ] **Step 1: Write uuid.go**

```go
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/StanMarek/forge/core/tools"
	"github.com/spf13/cobra"
)

var uuidCmd = &cobra.Command{
	Use:   "uuid",
	Short: "Generate, validate, or parse UUIDs",
}

var uuidGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new UUID",
	Run: func(cmd *cobra.Command, args []string) {
		ver, _ := cmd.Flags().GetInt("version")
		count, _ := cmd.Flags().GetInt("count")
		uppercase, _ := cmd.Flags().GetBool("uppercase")
		noHyphens, _ := cmd.Flags().GetBool("no-hyphens")

		for i := 0; i < count; i++ {
			result := tools.UUIDGenerate(ver, uppercase, noHyphens)
			if result.Error != "" {
				exitWithError(result.Error)
			}
			fmt.Println(result.Output)
		}
	},
}

var uuidValidateCmd = &cobra.Command{
	Use:   "validate [input]",
	Short: "Check if a string is a valid UUID",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.UUIDValidate(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

var uuidParseCmd = &cobra.Command{
	Use:   "parse [input]",
	Short: "Parse a UUID and show its components",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.UUIDParse(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		asJSON, _ := cmd.Flags().GetBool("json")
		if asJSON {
			out := struct {
				UUID      string `json:"uuid"`
				Version   int    `json:"version"`
				Variant   string `json:"variant"`
				Timestamp string `json:"timestamp,omitempty"`
			}{
				UUID:      result.UUID,
				Version:   result.Version,
				Variant:   result.Variant,
				Timestamp: result.Timestamp,
			}
			data, _ := json.MarshalIndent(out, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Println(result.Output)
		}
	},
}

func init() {
	uuidGenerateCmd.Flags().Int("version", 4, "UUID version: 4 (random) or 7 (time-ordered)")
	uuidGenerateCmd.Flags().Int("count", 1, "Number of UUIDs to generate")
	uuidGenerateCmd.Flags().Bool("uppercase", false, "Output in uppercase")
	uuidGenerateCmd.Flags().Bool("no-hyphens", false, "Output without hyphens")
	uuidParseCmd.Flags().Bool("json", false, "Output as JSON")
	uuidCmd.AddCommand(uuidGenerateCmd)
	uuidCmd.AddCommand(uuidValidateCmd)
	uuidCmd.AddCommand(uuidParseCmd)
}
```

- [ ] **Step 2: Commit**

```bash
git add cmd/uuid.go
git commit -m "Add UUID CLI command with generate/validate/parse subcommands"
```

---

## Chunk 3: Build verification + smoke tests

### Task 10: Verify build and add cobra dependency

- [ ] **Step 1: Add cobra dependency and tidy**

```bash
go get github.com/spf13/cobra
go mod tidy
```

- [ ] **Step 2: Build the binary**

Run: `go build -o bin/forge .`
Expected: no errors, binary at `bin/forge`

- [ ] **Step 3: Commit dependency changes**

```bash
git add go.mod go.sum
git commit -m "Add cobra dependency"
```

### Task 11: Smoke tests

- [ ] **Step 1: Test base64**

Run: `./bin/forge base64 encode "Hello, World!"`
Expected: `SGVsbG8sIFdvcmxkIQ==`

Run: `echo "Hello, World!" | ./bin/forge base64 encode`
Expected: `SGVsbG8sIFdvcmxkIQ==`

Run: `./bin/forge base64 decode "SGVsbG8sIFdvcmxkIQ=="`
Expected: `Hello, World!`

- [ ] **Step 2: Test JWT**

Run: `./bin/forge jwt decode "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
Expected: Header, Payload, Signature sections

Run: `./bin/forge jwt validate "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
Expected: `valid`

- [ ] **Step 3: Test JSON**

Run: `./bin/forge json format '{"name":"forge","version":1}'`
Expected: formatted JSON with indentation

Run: `./bin/forge json minify '{ "name" : "forge" }'`
Expected: `{"name":"forge"}`

Run: `./bin/forge json validate '{"valid":true}'`
Expected: `valid`

- [ ] **Step 4: Test hash**

Run: `./bin/forge hash sha256 "hello world"`
Expected: `b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9`

Run: `echo -n "hello world" | ./bin/forge hash sha256`
Expected: same hash

- [ ] **Step 5: Test URL**

Run: `./bin/forge url encode "hello world"`
Expected: `hello%20world`

Run: `./bin/forge url parse "https://example.com:8080/path?q=hello"`
Expected: parsed components

- [ ] **Step 6: Test UUID**

Run: `./bin/forge uuid generate`
Expected: a UUID v4 string

Run: `./bin/forge uuid validate "550e8400-e29b-41d4-a716-446655440000"`
Expected: `valid (version 4)`

Run: `./bin/forge uuid generate --count 3`
Expected: 3 UUIDs, one per line

- [ ] **Step 7: Test version**

Run: `./bin/forge version`
Expected: forge dev / commit none / built unknown / go version

Run: `./bin/forge --version`
Expected: `forge version dev`

- [ ] **Step 8: Test error cases**

Run: `./bin/forge base64 decode "!!!invalid!!!" 2>&1; echo "exit: $?"`
Expected: error to stderr, exit code 1

Run: `./bin/forge hash sha999 "hello" 2>&1; echo "exit: $?"`
Expected: unsupported algorithm error, exit code 1

---

## Task Dependency Summary

```
Task 1 (version) ─┐
Task 2 (stdin)  ──┤
Task 3 (root+helpers+main+version cmd) ──┤
                   │
                   ├── Task 4  (base64) ─┐
                   ├── Task 5  (jwt)    ─┤
                   ├── Task 6  (json)   ─┤
                   ├── Task 7  (hash)   ─├── Task 10 (build) → Task 11 (smoke)
                   ├── Task 8  (url)    ─┤
                   └── Task 9  (uuid)   ─┘
```

Tasks 1-3 are sequential (bootstrap).
Tasks 4-9 (tool commands) are **independent** — parallelizable.
Tasks 10-11 require all files to exist.
