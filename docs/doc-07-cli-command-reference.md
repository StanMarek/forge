# CLI Command Reference

**Project:** Forge — Developer Utility Toolkit
**Date:** 2026-03-15

---

## Global Behaviour

All tools follow these conventions:

- If `[input]` is provided as an argument, it is used directly.
- If `[input]` is omitted or is `-`, input is read from stdin.
- Output is written to stdout. Errors are written to stderr.
- Exit code 0 on success, 1 on input error, 2 on usage error.
- All commands support `--help` and `-h` for usage information.
- No color output by default in pipe mode (detected via `os.Stdout` isatty check). Use `--color` to force.

---

## Root Command

```
forge — A developer's workbench for the terminal, browser, and desktop.

Usage:
  forge [command]

Available Commands:
  base64      Encode or decode Base64 strings
  jwt         Decode and inspect JWT tokens
  json        Format, minify, or validate JSON
  hash        Generate hash digests
  url         Encode, decode, or parse URLs
  uuid        Generate, validate, or parse UUIDs
  yaml        Convert between JSON and YAML
  timestamp   Convert between Unix timestamps and datetime
  number-base Convert between number bases
  regex       Test regular expressions
  html-entity Encode or decode HTML entities
  password    Generate random passwords
  lorem       Generate Lorem Ipsum text
  tui         Launch the interactive terminal UI
  web         Launch the web server
  version     Print version information
  completion  Generate shell completion scripts

Flags:
  -h, --help      Show help for forge
  -v, --version   Print version

Use "forge [command] --help" for more information about a command.
```

If no subcommand is given, `forge` launches the TUI by default.

---

## base64

```
forge base64 — Encode or decode Base64 strings

Usage:
  forge base64 <subcommand> [input] [flags]

Subcommands:
  encode    Encode input as Base64
  decode    Decode Base64 input to plaintext

Examples:
  forge base64 encode "hello world"
  forge base64 decode "aGVsbG8gd29ybGQ="
  echo "hello world" | forge base64 encode
  echo "aGVsbG8gd29ybGQ=" | forge base64 decode
  forge base64 encode --url-safe "https://example.com?foo=bar"
  cat binary.dat | forge base64 encode --no-padding
```

### base64 encode

```
forge base64 encode [input] [flags]

Encodes input as Base64.

Arguments:
  input     String to encode. Reads from stdin if omitted or "-".

Flags:
  --url-safe      Use URL-safe Base64 alphabet (RFC 4648 §5).
                  Replaces + with - and / with _
  --no-padding    Omit trailing = padding characters
  -h, --help      Show help
```

### base64 decode

```
forge base64 decode [input] [flags]

Decodes a Base64-encoded string to plaintext.

Arguments:
  input     Base64 string to decode. Reads from stdin if omitted or "-".

Flags:
  --url-safe      Expect URL-safe Base64 alphabet (RFC 4648 §5)
  -h, --help      Show help

Errors:
  Exits with code 1 if input is not valid Base64.
  Error message: "invalid base64: <detail>"
```

---

## jwt

```
forge jwt — Decode and inspect JWT tokens

Usage:
  forge jwt <subcommand> [token] [flags]

Subcommands:
  decode      Decode a JWT token (without signature validation)
  validate    Check if a string is a structurally valid JWT
```

### jwt decode

```
forge jwt decode [token] [flags]

Decodes a JWT token and prints the header and payload as formatted JSON.
Does NOT validate the signature — this is a decoding tool, not a verification tool.

Arguments:
  token     JWT token string. Reads from stdin if omitted or "-".

Flags:
  --header-only    Print only the header
  --payload-only   Print only the payload
  --compact        Print JSON on a single line (no formatting)
  -h, --help       Show help

Output format:
  --- Header ---
  {
    "alg": "HS256",
    "typ": "JWT"
  }
  --- Payload ---
  {
    "sub": "1234567890",
    "name": "John Doe",
    "iat": 1516239022
  }
  --- Signature ---
  SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c

Examples:
  forge jwt decode "eyJhbGciOiJIUzI1NiIs..."
  pbpaste | forge jwt decode
  forge jwt decode --payload-only "eyJ..."
```

### jwt validate

```
forge jwt validate [token] [flags]

Checks if a string is a structurally valid JWT (three dot-separated Base64URL segments).
Prints "valid" or "invalid: <reason>" and sets exit code accordingly.

Arguments:
  token     JWT token string. Reads from stdin if omitted or "-".

Flags:
  -h, --help    Show help

Exit codes:
  0   Valid JWT structure
  1   Invalid JWT structure
```

---

## json

```
forge json — Format, minify, or validate JSON

Usage:
  forge json <subcommand> [input] [flags]

Subcommands:
  format      Pretty-print JSON with indentation
  minify      Remove all whitespace (compact JSON)
  validate    Check if input is valid JSON
```

### json format

```
forge json format [input] [flags]

Pretty-prints JSON with configurable indentation.

Arguments:
  input     JSON string. Reads from stdin if omitted or "-".

Flags:
  --indent int    Number of spaces for indentation (default: 2)
  --tabs          Use tabs instead of spaces for indentation
  --sort-keys     Sort object keys alphabetically
  -h, --help      Show help

Examples:
  forge json format '{"name":"forge","version":1}'
  cat data.json | forge json format --indent 4
  forge json format --tabs < config.json
  curl -s api.example.com/data | forge json format --sort-keys
```

### json minify

```
forge json minify [input] [flags]

Removes all unnecessary whitespace from JSON.

Arguments:
  input     JSON string. Reads from stdin if omitted or "-".

Flags:
  -h, --help    Show help

Examples:
  forge json minify '{ "name" : "forge" }'
  cat pretty.json | forge json minify > compact.json
```

### json validate

```
forge json validate [input] [flags]

Validates that input is well-formed JSON. Prints "valid" or an error message
with the position of the first error.

Arguments:
  input     JSON string. Reads from stdin if omitted or "-".

Flags:
  -h, --help    Show help

Exit codes:
  0   Valid JSON
  1   Invalid JSON (error message printed to stderr)
```

---

## hash

```
forge hash — Generate hash digests

Usage:
  forge hash <algorithm> [input] [flags]

Algorithms:
  md5       MD5 (128-bit) — not cryptographically secure
  sha1      SHA-1 (160-bit) — not cryptographically secure
  sha256    SHA-256 (256-bit)
  sha512    SHA-512 (512-bit)

Arguments:
  input     String to hash. Reads from stdin if omitted or "-".

Flags:
  --uppercase     Output hash in uppercase hex
  --file string   Hash the contents of a file instead of a string
  -h, --help      Show help

Output:
  Prints the hex-encoded hash digest to stdout.

Examples:
  forge hash sha256 "hello world"
  echo -n "hello world" | forge hash sha256
  forge hash md5 --file ./document.pdf
  forge hash sha512 "sensitive data" --uppercase
```

---

## url

```
forge url — Encode, decode, or parse URLs

Usage:
  forge url <subcommand> [input] [flags]

Subcommands:
  encode    URL-encode a string (percent encoding)
  decode    Decode a URL-encoded string
  parse     Break a URL into its components
```

### url encode

```
forge url encode [input] [flags]

URL-encodes a string using percent encoding (RFC 3986).

Arguments:
  input     String to encode. Reads from stdin if omitted or "-".

Flags:
  --component     Encode as a URL component (encodes /, ?, &, = etc.)
                  Default: encodes as a full URL (preserves structure)
  -h, --help      Show help

Examples:
  forge url encode "hello world"
  forge url encode --component "key=value&foo=bar"
```

### url decode

```
forge url decode [input] [flags]

Decodes a URL-encoded (percent-encoded) string.

Arguments:
  input     URL-encoded string. Reads from stdin if omitted or "-".

Flags:
  -h, --help    Show help

Examples:
  forge url decode "hello%20world"
  forge url decode "https%3A%2F%2Fexample.com%2Fpath%3Fq%3Dfoo"
```

### url parse

```
forge url parse [input] [flags]

Parses a URL and displays its components.

Arguments:
  input     URL to parse. Reads from stdin if omitted or "-".

Flags:
  --json          Output as JSON instead of human-readable format
  -h, --help      Show help

Output (default):
  Scheme:    https
  Host:      example.com
  Port:      443
  Path:      /api/v1/users
  Query:     page=1&limit=10
  Fragment:  section1
  Params:
    page    = 1
    limit   = 10

Output (--json):
  {"scheme":"https","host":"example.com","port":"443",...}

Examples:
  forge url parse "https://example.com:8080/path?q=hello#section"
  echo "https://api.example.com/v1/users?page=1" | forge url parse --json
```

---

## uuid

```
forge uuid — Generate, validate, or parse UUIDs

Usage:
  forge uuid <subcommand> [input] [flags]

Subcommands:
  generate    Generate a new UUID
  validate    Check if a string is a valid UUID
  parse       Parse a UUID and show its components
```

### uuid generate

```
forge uuid generate [flags]

Generates a new UUID.

Flags:
  --version int    UUID version to generate: 4 (random) or 7 (time-ordered)
                   (default: 4)
  --count int      Number of UUIDs to generate (default: 1)
  --uppercase      Output in uppercase
  --no-hyphens     Output without hyphens
  -h, --help       Show help

Examples:
  forge uuid generate
  forge uuid generate --version 7
  forge uuid generate --count 10
  forge uuid generate --version 4 --uppercase --no-hyphens
```

### uuid validate

```
forge uuid validate [input] [flags]

Validates whether a string is a properly formatted UUID.

Arguments:
  input     UUID string. Reads from stdin if omitted or "-".

Flags:
  -h, --help    Show help

Exit codes:
  0   Valid UUID
  1   Invalid UUID

Output:
  "valid (version 4)" or "invalid: <reason>"
```

### uuid parse

```
forge uuid parse [input] [flags]

Parses a UUID and displays its components.

Arguments:
  input     UUID string. Reads from stdin if omitted or "-".

Flags:
  --json          Output as JSON
  -h, --help      Show help

Output:
  UUID:       550e8400-e29b-41d4-a716-446655440000
  Version:    4
  Variant:    RFC 4122
  (v7 only)
  Timestamp:  2024-01-15T10:30:00.000Z

Examples:
  forge uuid parse "550e8400-e29b-41d4-a716-446655440000"
  forge uuid parse --json "01945b12-3c4d-7000-8000-000000000001"
```

---

## yaml

```
forge yaml — Convert between JSON and YAML

Usage:
  forge yaml <subcommand> [input] [flags]

Subcommands:
  to-json     Convert YAML to JSON
  to-yaml     Convert JSON to YAML (alias: from-json)
```

### yaml to-json

```
forge yaml to-json [input] [flags]

Converts YAML input to JSON.

Arguments:
  input     YAML string. Reads from stdin if omitted or "-".

Flags:
  --compact       Output compact JSON (no indentation)
  --indent int    Indentation spaces for JSON output (default: 2)
  -h, --help      Show help

Examples:
  forge yaml to-json "name: forge\nversion: 1"
  cat config.yaml | forge yaml to-json
  forge yaml to-json < docker-compose.yml --compact
```

### yaml to-yaml

```
forge yaml to-yaml [input] [flags]

Converts JSON input to YAML. Alias: forge yaml from-json

Arguments:
  input     JSON string. Reads from stdin if omitted or "-".

Flags:
  -h, --help    Show help

Examples:
  forge yaml to-yaml '{"name":"forge","version":1}'
  cat data.json | forge yaml to-yaml
```

---

## timestamp

```
forge timestamp — Convert between Unix timestamps and human-readable datetime

Usage:
  forge timestamp <subcommand> [input] [flags]

Subcommands:
  to-unix       Convert datetime string to Unix timestamp
  from-unix     Convert Unix timestamp to datetime string
  now           Print the current time in multiple formats
```

### timestamp from-unix

```
forge timestamp from-unix [timestamp] [flags]

Converts a Unix timestamp to a human-readable datetime.

Arguments:
  timestamp    Unix timestamp (seconds or milliseconds).

Flags:
  --format string    Output format: "rfc3339", "iso8601", "human" (default: "rfc3339")
  --tz string        Timezone for output (default: "UTC"). Example: "America/New_York"
  -h, --help         Show help

Examples:
  forge timestamp from-unix 1700000000
  forge timestamp from-unix 1700000000000
  forge timestamp from-unix 1700000000 --tz "Europe/Warsaw"
  forge timestamp from-unix 1700000000 --format human
```

### timestamp to-unix

```
forge timestamp to-unix [datetime] [flags]

Converts a datetime string to a Unix timestamp.

Arguments:
  datetime    Datetime string (supports RFC3339, ISO8601, common formats).

Flags:
  --millis          Output in milliseconds instead of seconds
  -h, --help        Show help

Examples:
  forge timestamp to-unix "2024-01-15T10:30:00Z"
  forge timestamp to-unix "2024-01-15 10:30:00" --millis
```

### timestamp now

```
forge timestamp now [flags]

Prints the current time in multiple formats.

Flags:
  --tz string    Timezone (default: local)
  -h, --help     Show help

Output:
  Unix (s):     1700000000
  Unix (ms):    1700000000000
  RFC 3339:     2024-11-14T22:13:20Z
  ISO 8601:     2024-11-14T22:13:20+00:00
  Human:        Thursday, November 14, 2024 10:13 PM UTC
```

---

## number-base

```
forge number-base — Convert between number bases

Usage:
  forge number-base [input] [flags]

Arguments:
  input     Number to convert. Prefix with 0x (hex), 0b (binary), 0o (octal),
            or no prefix for decimal. Reads from stdin if omitted or "-".

Flags:
  --from int      Source base (2-36). Auto-detected from prefix if omitted.
  --to int        Target base (2-36). If omitted, shows all common bases.
  -h, --help      Show help

Output (default — all bases):
  Decimal:     255
  Hex:         ff
  Octal:       377
  Binary:      11111111

Examples:
  forge number-base 255
  forge number-base 0xff
  forge number-base 0b11111111
  forge number-base --from 16 --to 2 "ff"
  echo "42" | forge number-base
```

---

## regex

```
forge regex — Test regular expressions

Usage:
  forge regex <pattern> [test-string] [flags]

Arguments:
  pattern       Regular expression pattern
  test-string   String to test against. Reads from stdin if omitted or "-".

Flags:
  --global        Find all matches (not just the first)
  --case-insensitive   Case-insensitive matching (re2: (?i) prefix)
  --json          Output matches as JSON
  -h, --help      Show help

Output (default):
  Pattern:  \d+
  Input:    "Order 42 has 3 items"
  Matches:
    [0] "42"  (pos 6-8)
    [1] "3"   (pos 13-14)

Examples:
  forge regex '\d+' "Order 42 has 3 items" --global
  echo "hello@example.com" | forge regex '[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}'
  forge regex --json '(\w+)=(\w+)' "key=value"
```

---

## html-entity

```
forge html-entity — Encode or decode HTML entities

Usage:
  forge html-entity <subcommand> [input] [flags]

Subcommands:
  encode    Encode special characters as HTML entities
  decode    Decode HTML entities to characters
```

### html-entity encode / decode

```
forge html-entity encode [input]
forge html-entity decode [input]

Arguments:
  input     Text to encode/decode. Reads from stdin if omitted or "-".

Flags:
  -h, --help    Show help

Examples:
  forge html-entity encode '<script>alert("xss")</script>'
  forge html-entity decode '&lt;div class=&quot;main&quot;&gt;Hello&lt;/div&gt;'
```

---

## password

```
forge password — Generate random passwords

Usage:
  forge password [flags]

Flags:
  --length int         Password length (default: 16)
  --count int          Number of passwords to generate (default: 1)
  --no-uppercase       Exclude uppercase letters
  --no-lowercase       Exclude lowercase letters
  --no-digits          Exclude digits
  --no-symbols         Exclude symbols
  --symbols string     Custom symbol set (default: "!@#$%^&*()-_=+")
  -h, --help           Show help

Examples:
  forge password
  forge password --length 32
  forge password --length 20 --no-symbols
  forge password --count 5 --length 24
```

---

## lorem

```
forge lorem — Generate Lorem Ipsum text

Usage:
  forge lorem [flags]

Flags:
  --words int         Number of words (default: 0, use --paragraphs)
  --sentences int     Number of sentences (default: 0)
  --paragraphs int    Number of paragraphs (default: 1)
  -h, --help          Show help

Exactly one of --words, --sentences, or --paragraphs must be specified.

Examples:
  forge lorem --paragraphs 3
  forge lorem --words 50
  forge lorem --sentences 5
```

---

## tui

```
forge tui — Launch the interactive terminal UI

Usage:
  forge tui [flags]

Flags:
  --tool string    Open directly to a specific tool (e.g., "base64", "jwt")
  -h, --help       Show help

Examples:
  forge tui
  forge tui --tool jwt
```

---

## web

```
forge web — Launch the self-hosted web server

Usage:
  forge web [flags]

Flags:
  --port int        Port to listen on (default: 8080)
  --host string     Host to bind to (default: "localhost")
  --open            Open browser automatically after starting
  -h, --help        Show help

Examples:
  forge web
  forge web --port 3000
  forge web --host 0.0.0.0 --port 8080
```

---

## version

```
forge version

Prints version, commit hash, build date, and Go version.

Output:
  forge v0.1.0
  commit:  abc1234
  built:   2024-01-15T10:30:00Z
  go:      go1.22.5
```

---

## Shell Completions

```
forge completion <shell>

Generates shell completion scripts.

Shells:
  bash        Bash completion script
  zsh         Zsh completion script
  fish        Fish completion script
  powershell  PowerShell completion script

Examples:
  forge completion bash > /etc/bash_completion.d/forge
  forge completion zsh > "${fpath[1]}/_forge"
  forge completion fish > ~/.config/fish/completions/forge.fish
```
