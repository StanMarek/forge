# Tool Inventory

**Project:** Forge — Developer Utility Toolkit
**Date:** 2026-03-15

---

## Design Principles for Tool Selection

1. **Daily use**: Only include tools developers reach for regularly, not once-a-year utilities.
2. **Stateless**: Every tool takes input, produces output, no database or filesystem state.
3. **Detectable**: Prioritize tools where clipboard smart detection adds real value.
4. **CLI-friendly**: Every tool must work in non-interactive pipe mode.
5. **No network**: All tools run offline. No API calls, no DNS lookups.

---

## v1.0 Tool Inventory

### Tier 1 — Must Ship (Core)

| Tool ID | Category | Name | Input | Output | DetectFromClipboard | Priority |
|---------|----------|------|-------|--------|---------------------|----------|
| `base64` | Encoders | Base64 Encode/Decode | Text or Base64 string | Encoded/decoded text | `true` — regex match `^[A-Za-z0-9+/=]+$` (min 4 chars, mod 4 length) | **Must** |
| `jwt` | Encoders | JWT Decoder | JWT token string | Decoded header, payload, signature | `true` — matches `xxxxx.yyyyy.zzzzz` pattern (3 dot-separated base64url segments) | **Must** |
| `json` | Formatters | JSON Formatter | JSON string | Pretty-printed or minified JSON | `true` — valid JSON (`json.Valid()`) | **Must** |
| `hash` | Generators | Hash Generator | Text string | Hash digest (hex) | `false` — hashing is one-way, detection is meaningless | **Must** |
| `url` | Encoders | URL Encode/Decode/Parse | URL or text | Encoded/decoded text, or parsed components | `true` — starts with `http://` or `https://` | **Must** |
| `uuid` | Generators | UUID Generate/Validate | UUID string or none | Generated UUID or parsed info | `true` — matches UUID regex `[0-9a-f]{8}-[0-9a-f]{4}-...` | **Must** |

### Tier 2 — Should Ship (High Value)

| Tool ID | Category | Name | Input | Output | DetectFromClipboard | Priority |
|---------|----------|------|-------|--------|---------------------|----------|
| `yaml` | Converters | JSON/YAML Converter | JSON or YAML string | Converted to the other format | `true` — valid YAML that isn't also valid JSON | **Should** |
| `timestamp` | Converters | Unix Timestamp Converter | Unix timestamp or datetime string | Converted datetime / Unix epoch | `true` — matches 10 or 13 digit number | **Should** |
| `number-base` | Converters | Number Base Converter | Number in any base | Converted to decimal, hex, octal, binary | `true` — matches `0x`, `0b`, `0o` prefixed numbers | **Should** |
| `regex` | Testers | Regex Tester | Pattern + test string | Match results, groups, highlights | `false` — regex patterns are too ambiguous to detect | **Should** |
| `html-entity` | Encoders | HTML Entity Encode/Decode | Text with HTML entities or raw text | Encoded/decoded HTML | `true` — contains `&amp;`, `&lt;`, `&#...;` patterns | **Should** |
| `password` | Generators | Password Generator | Length + options | Random password | `false` — generative tool, no clipboard input | **Should** |
| `lorem` | Generators | Lorem Ipsum Generator | Word/sentence/paragraph count | Lorem ipsum text | `false` — generative tool, no clipboard input | **Should** |

### Tier 3 — Nice to Have (v1.1+)

| Tool ID | Category | Name | Input | Output | DetectFromClipboard | Priority |
|---------|----------|------|-------|--------|---------------------|----------|
| `xml` | Formatters | XML Formatter | XML string | Pretty-printed XML | `true` — starts with `<?xml` or `<` with matching close tag | **Nice** |
| `sql` | Formatters | SQL Formatter | SQL query string | Formatted SQL | `true` — starts with common SQL keywords (SELECT, INSERT, etc.) | **Nice** |
| `cron` | Converters | Cron Expression Parser | Cron expression | Human-readable schedule | `true` — matches 5 or 6 space-separated fields | **Nice** |
| `color` | Converters | Color Converter | Hex, RGB, HSL color | Converted color formats + preview | `true` — matches `#[0-9a-f]{3,8}` or `rgb(...)` | **Nice** |
| `diff` | Text | Text Diff/Compare | Two text inputs | Side-by-side diff | `false` — requires two inputs | **Nice** |
| `markdown` | Text | Markdown Preview | Markdown text | Rendered preview (web only) | `true` — contains common markdown syntax (`#`, `**`, `- `) | **Nice** |
| `gzip` | Encoders | GZip Compress/Decompress | Text or binary | Compressed/decompressed output | `false` — binary detection is unreliable | **Nice** |
| `qrcode` | Generators | QR Code Generator | Text or URL | QR code image (web/desktop) or ASCII (TUI) | `false` — generative tool | **Nice** |
| `escape` | Encoders | String Escape/Unescape | Text with escape sequences | Escaped/unescaped text | `true` — contains `\n`, `\t`, `\"` etc. | **Nice** |
| `jsonpath` | Testers | JSONPath Tester | JSON + JSONPath expression | Query results | `false` — requires two inputs | **Nice** |

---

## Tools Explicitly Excluded from v1

These are tools DevToys includes that Forge should **not** ship in v1, with reasoning:

| Tool | Reason for Exclusion |
|------|---------------------|
| Image Compressor (PNG/JPEG) | Requires image processing libraries, adds CGO dependency, niche use case for a dev toolkit |
| Color Blindness Simulator | Too niche, requires image rendering capabilities |
| Certificate Parser | Complex X.509 parsing, niche audience, better served by `openssl` CLI |
| File Splitter | OS-level operation, better done with `split` command |
| Language-specific converters (JSON to C#, JSON to PHP) | Scope creep — these are codegen tools, not utility tools |
| Text Comparer (advanced) | Full diff engines are complex; start with basic diff in Tier 3, expand if demanded |
| RESX Translator | .NET specific, not relevant to Forge's audience |

---

## Category Summary

| Category | Tier 1 Tools | Tier 2 Tools | Tier 3 Tools |
|----------|-------------|-------------|-------------|
| Encoders | base64, jwt, url | html-entity | gzip, escape |
| Formatters | json | — | xml, sql |
| Generators | hash, uuid | password, lorem | qrcode |
| Converters | — | yaml, timestamp, number-base | cron, color |
| Testers | — | regex | jsonpath |
| Text | — | — | diff, markdown |

---

## Smart Detection Priority

When clipboard content matches multiple tools, use this priority order (most specific first):

1. **JWT** — three dot-separated base64url segments is highly specific
2. **UUID** — UUID format is unmistakable
3. **URL** — `http://` or `https://` prefix is definitive
4. **JSON** — valid JSON parsing is deterministic
5. **Base64** — most ambiguous detector; only trigger if no higher-priority match
6. **Timestamp** — 10/13 digit numbers could be many things; low confidence
7. **HTML entity** — `&amp;` patterns are fairly specific
8. **Number base** — `0x`/`0b`/`0o` prefixes are specific but rare in clipboard

When multiple tools match, present all matches in the detection banner, with the highest-priority match pre-selected.

---

## Implementation Notes

- Every tool in Tier 1 and Tier 2 must have complete unit tests before the tool is considered done.
- Tool functions must be pure: `func DoThing(input string, options Options) Result`. No I/O, no global state.
- Each `Result` struct contains `Output string` and `Error string`. Success is `Error == ""`.
- Tier 3 tools should be designed into the architecture (registered in the registry) but can be implemented as stubs that return "Coming soon" in the v1.0 release.
