package tools

import (
	"testing"
)

// ---------------------------------------------------------------------------
// URLTool metadata
// ---------------------------------------------------------------------------

func TestURLToolMetadata(t *testing.T) {
	tool := URLTool{}

	if tool.Name() != "URL Encode / Decode / Parse" {
		t.Errorf("Name() = %q, want %q", tool.Name(), "URL Encode / Decode / Parse")
	}
	if tool.ID() != "url" {
		t.Errorf("ID() = %q, want %q", tool.ID(), "url")
	}
	if tool.Category() != "Encoders" {
		t.Errorf("Category() = %q, want %q", tool.Category(), "Encoders")
	}

	keywords := tool.Keywords()
	expected := []string{"url", "encode", "decode", "parse", "percent", "uri"}
	if len(keywords) != len(expected) {
		t.Fatalf("Keywords() length = %d, want %d", len(keywords), len(expected))
	}
	for i, kw := range expected {
		if keywords[i] != kw {
			t.Errorf("Keywords()[%d] = %q, want %q", i, keywords[i], kw)
		}
	}
}

func TestURLToolImplementsInterface(t *testing.T) {
	var _ Tool = URLTool{}
}

// ---------------------------------------------------------------------------
// DetectFromClipboard
// ---------------------------------------------------------------------------

func TestURLDetectFromClipboard(t *testing.T) {
	tool := URLTool{}

	tests := []struct {
		input string
		want  bool
	}{
		{"http://example.com", true},
		{"https://example.com", true},
		{"https://example.com/path?q=1", true},
		{"  https://example.com  ", true}, // leading/trailing whitespace
		{"ftp://example.com", false},
		{"example.com", false},
		{"not a url", false},
		{"", false},
	}

	for _, tc := range tests {
		got := tool.DetectFromClipboard(tc.input)
		if got != tc.want {
			t.Errorf("DetectFromClipboard(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

// ---------------------------------------------------------------------------
// URLEncode
// ---------------------------------------------------------------------------

func TestURLEncodeEmpty(t *testing.T) {
	r := URLEncode("", true)
	if r.Output != "" {
		t.Errorf("URLEncode empty: Output = %q, want %q", r.Output, "")
	}
	if r.Error != "" {
		t.Errorf("URLEncode empty: Error = %q, want empty", r.Error)
	}
}

func TestURLEncodeComponent(t *testing.T) {
	// component=true uses QueryEscape: space → +
	r := URLEncode("hello world", true)
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if r.Output != "hello+world" {
		t.Errorf("URLEncode component: Output = %q, want %q", r.Output, "hello+world")
	}
}

func TestURLEncodeComponentSpecialChars(t *testing.T) {
	r := URLEncode("a=1&b=2", true)
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if r.Output != "a%3D1%26b%3D2" {
		t.Errorf("URLEncode component special: Output = %q, want %q", r.Output, "a%3D1%26b%3D2")
	}
}

func TestURLEncodeFullURL(t *testing.T) {
	// component=false uses PathEscape: space → %20
	r := URLEncode("hello world", false)
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if r.Output != "hello%20world" {
		t.Errorf("URLEncode path: Output = %q, want %q", r.Output, "hello%20world")
	}
}

func TestURLEncodePathPreservesSlash(t *testing.T) {
	// PathEscape does NOT encode forward slashes in path segments,
	// but it does encode them within a single segment.
	r := URLEncode("/foo/bar", false)
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	// url.PathEscape encodes / as %2F
	if r.Output != "%2Ffoo%2Fbar" {
		t.Errorf("URLEncode path slash: Output = %q, want %q", r.Output, "%2Ffoo%2Fbar")
	}
}

func TestURLEncodeUnicode(t *testing.T) {
	r := URLEncode("café", true)
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	expected := "caf%C3%A9"
	if r.Output != expected {
		t.Errorf("URLEncode unicode: Output = %q, want %q", r.Output, expected)
	}
}

// ---------------------------------------------------------------------------
// URLDecode
// ---------------------------------------------------------------------------

func TestURLDecodeEmpty(t *testing.T) {
	r := URLDecode("")
	if r.Output != "" {
		t.Errorf("URLDecode empty: Output = %q, want %q", r.Output, "")
	}
	if r.Error != "" {
		t.Errorf("URLDecode empty: Error = %q, want empty", r.Error)
	}
}

func TestURLDecodePercent20(t *testing.T) {
	r := URLDecode("hello%20world")
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if r.Output != "hello world" {
		t.Errorf("URLDecode %%20: Output = %q, want %q", r.Output, "hello world")
	}
}

func TestURLDecodePlus(t *testing.T) {
	r := URLDecode("hello+world")
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if r.Output != "hello world" {
		t.Errorf("URLDecode +: Output = %q, want %q", r.Output, "hello world")
	}
}

func TestURLDecodeSpecialChars(t *testing.T) {
	r := URLDecode("a%3D1%26b%3D2")
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if r.Output != "a=1&b=2" {
		t.Errorf("URLDecode special: Output = %q, want %q", r.Output, "a=1&b=2")
	}
}

func TestURLDecodeUnicode(t *testing.T) {
	r := URLDecode("caf%C3%A9")
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if r.Output != "café" {
		t.Errorf("URLDecode unicode: Output = %q, want %q", r.Output, "café")
	}
}

func TestURLDecodeInvalid(t *testing.T) {
	r := URLDecode("%ZZ")
	if r.Error == "" {
		t.Errorf("URLDecode invalid: expected error, got Output = %q", r.Output)
	}
}

func TestURLDecodeAlreadyDecoded(t *testing.T) {
	r := URLDecode("hello")
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if r.Output != "hello" {
		t.Errorf("URLDecode plain: Output = %q, want %q", r.Output, "hello")
	}
}

// ---------------------------------------------------------------------------
// URLParse
// ---------------------------------------------------------------------------

func TestURLParseEmpty(t *testing.T) {
	r := URLParse("")
	if r.Error == "" {
		t.Errorf("URLParse empty: expected error")
	}
}

func TestURLParseMissingScheme(t *testing.T) {
	r := URLParse("example.com/path")
	if r.Error == "" {
		t.Errorf("URLParse no scheme: expected error")
	}
}

func TestURLParseSimple(t *testing.T) {
	r := URLParse("https://example.com")
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if r.Scheme != "https" {
		t.Errorf("Scheme = %q, want %q", r.Scheme, "https")
	}
	if r.Host != "example.com" {
		t.Errorf("Host = %q, want %q", r.Host, "example.com")
	}
	if r.Port != "" {
		t.Errorf("Port = %q, want empty", r.Port)
	}
}

func TestURLParseWithPort(t *testing.T) {
	r := URLParse("http://localhost:8080/api")
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if r.Scheme != "http" {
		t.Errorf("Scheme = %q, want %q", r.Scheme, "http")
	}
	if r.Host != "localhost" {
		t.Errorf("Host = %q, want %q", r.Host, "localhost")
	}
	if r.Port != "8080" {
		t.Errorf("Port = %q, want %q", r.Port, "8080")
	}
	if r.Path != "/api" {
		t.Errorf("Path = %q, want %q", r.Path, "/api")
	}
}

func TestURLParseFull(t *testing.T) {
	r := URLParse("https://example.com:443/path/to/resource?key=val&foo=bar#section")
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if r.Scheme != "https" {
		t.Errorf("Scheme = %q", r.Scheme)
	}
	if r.Host != "example.com" {
		t.Errorf("Host = %q", r.Host)
	}
	if r.Port != "443" {
		t.Errorf("Port = %q", r.Port)
	}
	if r.Path != "/path/to/resource" {
		t.Errorf("Path = %q", r.Path)
	}
	if r.Query != "key=val&foo=bar" {
		t.Errorf("Query = %q", r.Query)
	}
	if r.Fragment != "section" {
		t.Errorf("Fragment = %q", r.Fragment)
	}
	if r.Params == nil {
		t.Fatal("Params is nil")
	}
	if vals, ok := r.Params["key"]; !ok || len(vals) != 1 || vals[0] != "val" {
		t.Errorf("Params[key] = %v", r.Params["key"])
	}
	if vals, ok := r.Params["foo"]; !ok || len(vals) != 1 || vals[0] != "bar" {
		t.Errorf("Params[foo] = %v", r.Params["foo"])
	}
}

func TestURLParseMultiValueParam(t *testing.T) {
	r := URLParse("https://example.com/search?tag=go&tag=rust")
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if vals, ok := r.Params["tag"]; !ok || len(vals) != 2 {
		t.Errorf("Params[tag] = %v, want [go rust]", r.Params["tag"])
	}
}

func TestURLParseOutputContainsScheme(t *testing.T) {
	r := URLParse("https://example.com/path")
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if r.Output == "" {
		t.Fatal("Output is empty")
	}
	// Output should contain the scheme line
	if !contains(r.Output, "Scheme:") {
		t.Errorf("Output missing Scheme line: %s", r.Output)
	}
	if !contains(r.Output, "Host:") {
		t.Errorf("Output missing Host line: %s", r.Output)
	}
	if !contains(r.Output, "Path:") {
		t.Errorf("Output missing Path line: %s", r.Output)
	}
}

func TestURLParseNoPath(t *testing.T) {
	r := URLParse("https://example.com")
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	// Path should be empty for bare domain
	if r.Path != "" {
		t.Errorf("Path = %q, want empty", r.Path)
	}
}

func TestURLParseWithFragment(t *testing.T) {
	r := URLParse("https://example.com/page#top")
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if r.Fragment != "top" {
		t.Errorf("Fragment = %q, want %q", r.Fragment, "top")
	}
}

// ---------------------------------------------------------------------------
// URLEncode / URLDecode roundtrip
// ---------------------------------------------------------------------------

func TestURLEncodeDecodeRoundtrip(t *testing.T) {
	inputs := []string{
		"hello world",
		"a=1&b=2",
		"café",
		"foo bar/baz?qux=1",
		"https://example.com/path?q=hello world",
	}

	for _, input := range inputs {
		encoded := URLEncode(input, true)
		if encoded.Error != "" {
			t.Fatalf("encode error for %q: %s", input, encoded.Error)
		}
		decoded := URLDecode(encoded.Output)
		if decoded.Error != "" {
			t.Fatalf("decode error for %q: %s", encoded.Output, decoded.Error)
		}
		if decoded.Output != input {
			t.Errorf("roundtrip failed for %q: got %q", input, decoded.Output)
		}
	}
}

// contains is a small helper for substring checks.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
