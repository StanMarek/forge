# Forge

A developer utility toolkit for the terminal, browser, and desktop. Encoding, decoding, formatting, hashing, and generation tools — all in one place.

## Install

```bash
brew install StanMarek/tap/forge
```

Or build from source:

```bash
go install github.com/StanMarek/forge@latest
```

## Usage

Launch the interactive TUI:

```bash
forge
```

Use tools directly from the CLI:

```bash
forge base64 encode "hello world"
forge jwt decode "eyJhbGciOi..."
forge hash sha256 "password"
forge uuid generate
forge json format '{"a":1}'
forge url encode "hello world&foo=bar"
```

Start the web UI:

```bash
forge web
```

## Tools

| Category | Tools |
|----------|-------|
| Encoders | Base64, JWT, URL, HTML Entity, Text Escape, Gzip |
| Formatters | JSON, YAML, XML, CSV, Diff |
| Generators | UUID, Password, Lorem Ipsum, Color |
| Converters | Timestamp, Number Base, Cron |
| Analyzers | Hash, Regex, Text Stats |

## Surfaces

- **CLI** — pipe-friendly commands with stdin support
- **TUI** — interactive terminal UI (default when no subcommand given)
- **Web** — browser-based UI via `forge web`
- **Desktop** — native-looking window via `forge desktop`

## License

[MIT](LICENSE)
