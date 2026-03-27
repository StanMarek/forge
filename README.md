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

<img width="1309" height="677" alt="image" src="https://github.com/user-attachments/assets/1d8752c2-cf7a-4993-bf31-2d452a5816f1" />


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

<img width="1207" height="885" alt="image" src="https://github.com/user-attachments/assets/eb0583af-e7b0-4312-88f1-dc66bb7e19a6" />
<img width="1200" height="847" alt="image" src="https://github.com/user-attachments/assets/bc73589e-ae6f-4d7c-b3af-c28aa4e5966c" />


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
