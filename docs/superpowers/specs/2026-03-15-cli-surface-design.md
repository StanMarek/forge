# CLI Surface — Design Spec

**Date:** 2026-03-15
**Status:** Approved
**Scope:** `main.go`, `cmd/`, `internal/version/`, `internal/stdin/`

---

## Goal

Implement the Cobra CLI surface for Forge's 6 Tier-1 tools. Each tool is invoked as `forge <tool> <subcommand> [input] [flags]`. Input comes from args or stdin. Output goes to stdout, errors to stderr.

## Dependencies

- `github.com/spf13/cobra` — CLI framework
- `github.com/StanMarek/forge/core/tools` — business logic (already implemented)
- `github.com/StanMarek/forge/internal/version` — build info
- `github.com/StanMarek/forge/internal/stdin` — stdin reading

## Package: `internal/stdin/stdin.go`

```go
func Read() (string, error)
```

- Uses `os.Stdin.Stat()` to check `ModeCharDevice` — if stdin is a terminal (no pipe), returns an error.
- If piped, reads all of stdin via `io.ReadAll(os.Stdin)`, trims trailing newline.
- Returns the input string or error.

## Package: `internal/version/version.go`

```go
var (
    Version = "dev"
    Commit  = "none"
    Date    = "unknown"
)
```

Set via `-ldflags` at build time. Defaults allow `go run .` to work without ldflags.

## Package: `cmd/`

### Input Resolution Pattern

Every tool command follows this pattern for resolving input:

```go
func resolveInput(args []string) (string, error)
```

1. If `args[0]` exists and is not `"-"`, use it as input.
2. Otherwise, call `stdin.Read()`.
3. If both fail, return error.

This is a shared helper in `cmd/root.go` or a small `cmd/helpers.go`.

### root.go

```go
var rootCmd = &cobra.Command{
    Use:   "forge",
    Short: "A developer's workbench for the terminal, browser, and desktop",
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(2)
    }
}
```

Registers all subcommands in `init()`. No default action — prints help when invoked without subcommand.

### version.go

```go
forge version
```

Output format:
```
forge <version>
commit:  <hash>
built:   <date>
go:      <go version>
```

Uses `runtime.Version()` for Go version.

### base64.go

```
forge base64 encode [input] [--url-safe] [--no-padding]
forge base64 decode [input] [--url-safe]
```

- Parent command `base64` with two subcommands: `encode`, `decode`.
- Flags: `--url-safe` (bool), `--no-padding` (bool, encode only).
- Calls `tools.Base64Encode()` / `tools.Base64Decode()`.
- Prints `result.Output` to stdout. On error: prints to stderr, exits 1.

### jwt.go

```
forge jwt decode [token] [--header-only] [--payload-only] [--compact]
forge jwt validate [token]
```

- Parent command `jwt` with subcommands: `decode`, `validate`.
- `decode` flags: `--header-only`, `--payload-only`, `--compact` (all bool).
- `--header-only`: prints only `result.Header`.
- `--payload-only`: prints only `result.Payload`.
- `--compact`: prints JSON on single lines (re-minify the header/payload).
- `validate`: prints `result.Output` ("valid") or error with exit 1.

### json.go

```
forge json format [input] [--indent N] [--tabs] [--sort-keys]
forge json minify [input]
forge json validate [input]
```

- Parent command `json` with subcommands: `format`, `minify`, `validate`.
- `format` flags: `--indent int` (default 2), `--tabs` (bool), `--sort-keys` (bool).
- `validate`: prints "valid" or error with exit 1.

### hash.go

```
forge hash <algorithm> [input] [--uppercase] [--file PATH]
```

- Single command, not subcommands. Algorithm is the first positional arg.
- Supported algorithms: md5, sha1, sha256, sha512.
- `--uppercase` (bool): output in uppercase hex.
- `--file` (string): hash file contents instead of string input. Reads file bytes, passes to `tools.Hash()`.
- Input resolution: if `--file` is set, read file. Otherwise, resolve from args[1] or stdin.

### url.go

```
forge url encode [input] [--component]
forge url decode [input]
forge url parse [input] [--json]
```

- Parent command `url` with subcommands: `encode`, `decode`, `parse`.
- `encode` flag: `--component` (bool).
- `parse` flag: `--json` (bool) — output as JSON instead of human-readable.
- For `--json` on parse: marshal `URLParseResult` to JSON (excluding Output and Error fields).

### uuid.go

```
forge uuid generate [--version N] [--count N] [--uppercase] [--no-hyphens]
forge uuid validate [input]
forge uuid parse [input] [--json]
```

- Parent command `uuid` with subcommands: `generate`, `validate`, `parse`.
- `generate` flags: `--version int` (default 4), `--count int` (default 1), `--uppercase`, `--no-hyphens`.
- `--count`: loop calling `tools.UUIDGenerate()` N times, one UUID per line.
- `parse` flag: `--json` (bool).
- For `--json` on parse: marshal `UUIDParseResult` to JSON (excluding Output and Error fields).

## Error Handling

- Tool error (`result.Error != ""`): print to stderr, exit 1.
- Usage error (wrong flags, missing subcommand): Cobra handles automatically, exit 2.
- Stdin error (no pipe, no arg): print "no input provided" to stderr, exit 1.

## Testing Strategy

- Integration tests via `exec.Command` are overkill for this milestone.
- The core logic is already tested (200 tests). CLI layer is thin glue.
- Manual smoke testing: `go run . base64 encode "hello"`, `echo "hello" | go run . base64 encode`, etc.
- `internal/stdin/` gets a unit test with a mock reader.

## File Manifest

```
main.go
cmd/root.go
cmd/helpers.go
cmd/base64.go
cmd/jwt.go
cmd/json.go
cmd/hash.go
cmd/url.go
cmd/uuid.go
cmd/version.go
internal/version/version.go
internal/stdin/stdin.go
internal/stdin/stdin_test.go
```

## Out of Scope

- `forge tui` / `forge web` commands — separate milestones.
- Shell completions (`forge completion`) — can add later, Cobra generates them for free.
- `--color` flag / isatty detection — not needed until we add colored output.
- Tier-2 tool commands (yaml, timestamp, etc.) — not yet implemented in core.
