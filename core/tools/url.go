package tools

import (
	"fmt"
	"net/url"
	"strings"
)

// URLTool provides metadata for the URL encode/decode/parse tool.
type URLTool struct{}

func (u URLTool) Name() string        { return "URL Encode / Decode / Parse" }
func (u URLTool) ID() string          { return "url" }
func (u URLTool) Description() string { return "Encode, decode, and parse URLs" }
func (u URLTool) Category() string    { return "Encoders" }
func (u URLTool) Keywords() []string {
	return []string{"url", "encode", "decode", "parse", "percent", "uri"}
}

// DetectFromClipboard returns true if s looks like an HTTP or HTTPS URL.
func (u URLTool) DetectFromClipboard(s string) bool {
	s = strings.TrimSpace(s)
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

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

// URLEncode percent-encodes the input string.
// If component is true, url.QueryEscape is used (spaces become +).
// If component is false, url.PathEscape is used (spaces become %20).
func URLEncode(input string, component bool) Result {
	if input == "" {
		return Result{Output: ""}
	}

	var encoded string
	if component {
		encoded = url.QueryEscape(input)
	} else {
		encoded = url.PathEscape(input)
	}

	return Result{Output: encoded}
}

// URLDecode decodes a percent-encoded string.
// It handles both %20 and + as spaces.
func URLDecode(input string) Result {
	if input == "" {
		return Result{Output: ""}
	}

	decoded, err := url.QueryUnescape(input)
	if err != nil {
		return Result{Error: fmt.Sprintf("invalid URL encoding: %s", err.Error())}
	}

	return Result{Output: decoded}
}

// URLParse parses a URL string into its components.
// Returns an error for empty input or URLs without a scheme.
func URLParse(input string) URLParseResult {
	if input == "" {
		return URLParseResult{Error: "empty input"}
	}

	parsed, err := url.Parse(input)
	if err != nil {
		return URLParseResult{Error: fmt.Sprintf("invalid URL: %s", err.Error())}
	}

	if parsed.Scheme == "" {
		return URLParseResult{Error: "missing scheme (e.g. http:// or https://)"}
	}

	host := parsed.Hostname()
	port := parsed.Port()
	params := parsed.Query()

	var b strings.Builder
	b.WriteString(fmt.Sprintf("Scheme:   %s\n", parsed.Scheme))
	b.WriteString(fmt.Sprintf("Host:     %s\n", host))
	if port != "" {
		b.WriteString(fmt.Sprintf("Port:     %s\n", port))
	}
	if parsed.Path != "" {
		b.WriteString(fmt.Sprintf("Path:     %s\n", parsed.Path))
	}
	if parsed.RawQuery != "" {
		b.WriteString(fmt.Sprintf("Query:    %s\n", parsed.RawQuery))
	}
	if parsed.Fragment != "" {
		b.WriteString(fmt.Sprintf("Fragment: %s\n", parsed.Fragment))
	}
	if len(params) > 0 {
		b.WriteString("Params:\n")
		for key, values := range params {
			for _, val := range values {
				b.WriteString(fmt.Sprintf("  %s = %s\n", key, val))
			}
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
		Output:   strings.TrimRight(b.String(), "\n"),
	}
}
