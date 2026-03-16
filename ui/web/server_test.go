package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthEndpoint(t *testing.T) {
	router := NewRouter()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status":"ok"`)
}

func TestIndexPage(t *testing.T) {
	router := NewRouter()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, "Forge")
	assert.Contains(t, body, "Base64")
}

// toolPageTest defines a GET page test case.
type toolPageTest struct {
	name     string
	path     string
	contains []string // substrings expected in response body
}

// toolProcessTest defines a POST process test case.
type toolProcessTest struct {
	name     string
	path     string
	form     url.Values
	contains []string // substrings expected in response body
}

func TestAllToolPagesRender(t *testing.T) {
	router := NewRouter()

	tests := []toolPageTest{
		{"Base64", "/tools/base64", []string{"Base64"}},
		{"JWT", "/tools/jwt", []string{"JWT"}},
		{"JSON", "/tools/json", []string{"JSON"}},
		{"Hash", "/tools/hash", []string{"Hash"}},
		{"URL", "/tools/url", []string{"URL"}},
		{"UUID", "/tools/uuid", []string{"UUID"}},
		{"YAML", "/tools/yaml", []string{"YAML"}},
		{"Timestamp", "/tools/timestamp", []string{"Timestamp"}},
		{"NumberBase", "/tools/number-base", []string{"Number Base"}},
		{"Regex", "/tools/regex", []string{"Regex"}},
		{"HTMLEntity", "/tools/html-entity", []string{"HTML"}},
		{"Password", "/tools/password", []string{"Password"}},
		{"Lorem", "/tools/lorem", []string{"Lorem"}},
		{"Color", "/tools/color", []string{"Color"}},
		{"Cron", "/tools/cron", []string{"Cron"}},
		{"TextEscape", "/tools/text-escape", []string{"Text Escape"}},
		{"GZip", "/tools/gzip", []string{"GZip"}},
		{"TextStats", "/tools/text-stats", []string{"Text"}},
		{"Diff", "/tools/diff", []string{"Diff"}},
		{"XML", "/tools/xml", []string{"XML"}},
		{"CSV", "/tools/csv", []string{"CSV"}},
	}

	for _, tt := range tests {
		t.Run(tt.name+" GET", func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code, "GET %s returned %d", tt.path, rec.Code)

			body := rec.Body.String()
			assert.NotEmpty(t, body, "GET %s returned empty body", tt.path)
			for _, s := range tt.contains {
				assert.Contains(t, body, s, "GET %s missing %q", tt.path, s)
			}
			// Every tool page should render within the layout
			assert.Contains(t, body, "<html", "GET %s missing <html", tt.path)
			assert.Contains(t, body, "hx-post", "GET %s missing hx-post (HTMX form)", tt.path)
		})
	}
}

func TestAllToolProcessEndpoints(t *testing.T) {
	router := NewRouter()

	// Note: templ HTML-escapes content in textareas, so:
	//   " → &#34;    & → &amp;    < → &lt;    > → &gt;
	// Assertions must account for this encoding.

	tests := []toolProcessTest{
		{
			"Base64 encode", "/tools/base64",
			url.Values{"input": {"hello world"}, "mode": {"encode"}},
			[]string{"aGVsbG8gd29ybGQ="},
		},
		{
			"Base64 decode", "/tools/base64",
			url.Values{"input": {"aGVsbG8="}, "mode": {"decode"}},
			[]string{"hello"},
		},
		{
			"JWT decode", "/tools/jwt",
			url.Values{"input": {"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"}},
			[]string{"John Doe"},
		},
		{
			"JSON format", "/tools/json",
			url.Values{"input": {`{"a":1}`}, "mode": {"format"}},
			[]string{"&#34;a&#34;"}, // templ encodes " as &#34;
		},
		{
			"JSON minify", "/tools/json",
			url.Values{"input": {"{\n  \"a\": 1\n}"}, "mode": {"minify"}},
			[]string{"{&#34;a&#34;:1}"}, // templ encodes " as &#34;
		},
		{
			"JSON validate valid", "/tools/json",
			url.Values{"input": {`{"a":1}`}, "mode": {"validate"}},
			[]string{"valid"},
		},
		{
			"Hash", "/tools/hash",
			url.Values{"input": {"hello"}},
			[]string{"5d41402abc4b2a76b9719d911017c592"}, // MD5 of "hello"
		},
		{
			"Hash uppercase", "/tools/hash",
			url.Values{"input": {"hello"}, "uppercase": {"on"}},
			[]string{"5D41402ABC4B2A76B9719D911017C592"},
		},
		{
			"URL encode", "/tools/url",
			url.Values{"input": {"hello world"}, "mode": {"encode"}},
			[]string{"hello%20world"},
		},
		{
			"URL decode", "/tools/url",
			url.Values{"input": {"hello%20world"}, "mode": {"decode"}},
			[]string{"hello world"},
		},
		{
			"URL parse", "/tools/url",
			url.Values{"input": {"https://example.com/path?q=1"}, "mode": {"parse"}},
			[]string{"example.com"},
		},
		{
			"UUID generate", "/tools/uuid",
			url.Values{"mode": {"generate"}, "version": {"4"}},
			nil, // just check 200 status — output is random
		},
		{
			"UUID validate valid", "/tools/uuid",
			url.Values{"input": {"550e8400-e29b-41d4-a716-446655440000"}, "mode": {"validate"}},
			[]string{"valid"},
		},
		{
			"YAML to JSON", "/tools/yaml",
			url.Values{"input": {"name: test\nvalue: 42"}, "mode": {"yaml-to-json"}},
			[]string{"&#34;name&#34;", "&#34;test&#34;", "42"}, // templ encodes " as &#34;
		},
		{
			"JSON to YAML", "/tools/yaml",
			url.Values{"input": {`{"name":"test"}`}, "mode": {"json-to-yaml"}},
			[]string{"name:", "test"},
		},
		{
			"Timestamp now", "/tools/timestamp",
			url.Values{"mode": {"now"}},
			nil, // dynamic output, just check 200
		},
		{
			"Timestamp from-unix", "/tools/timestamp",
			url.Values{"input": {"0"}, "mode": {"from-unix"}},
			[]string{"1970"},
		},
		{
			"Number base", "/tools/number-base",
			url.Values{"input": {"42"}},
			[]string{"2a", "101010", "52"},
		},
		{
			"Regex test", "/tools/regex",
			url.Values{"pattern": {"\\d+"}, "input": {"abc123def"}, "global": {"on"}},
			[]string{"123"},
		},
		{
			"HTML entity encode", "/tools/html-entity",
			url.Values{"input": {"<div>"}, "mode": {"encode"}},
			[]string{"&amp;lt;div&amp;gt;"}, // templ double-encodes: &lt; → &amp;lt;
		},
		{
			"HTML entity decode", "/tools/html-entity",
			url.Values{"input": {"&lt;div&gt;"}, "mode": {"decode"}},
			nil, // the output contains literal <div> which is HTML — just check 200
		},
		{
			"Password generate", "/tools/password",
			url.Values{"length": {"32"}, "uppercase": {"on"}, "lowercase": {"on"}, "digits": {"on"}, "symbols": {"on"}},
			nil, // random output, just check 200
		},
		{
			"Lorem paragraphs", "/tools/lorem",
			url.Values{"mode": {"paragraphs"}, "count": {"2"}},
			nil, // lorem text is randomized, just check 200
		},
		{
			"Color convert hex", "/tools/color",
			url.Values{"input": {"#ff9800"}},
			[]string{"rgb", "hsl", "RGB", "HSL"},
		},
		{
			"Cron parse", "/tools/cron",
			url.Values{"input": {"*/5 * * * *"}},
			[]string{"5", "minute"},
		},
		{
			"Text escape", "/tools/text-escape",
			url.Values{"input": {"hello\tworld"}, "mode": {"escape"}},
			nil, // check 200
		},
		{
			"GZip compress", "/tools/gzip",
			url.Values{"input": {"hello world"}, "mode": {"compress"}},
			nil, // base64 encoded gzip — just check 200
		},
		{
			"Text stats", "/tools/text-stats",
			url.Values{"input": {"hello world"}, "mode": {"stats"}},
			[]string{"2", "11"}, // 2 words, 11 chars
		},
		{
			"Text case upper", "/tools/text-stats",
			url.Values{"input": {"hello"}, "mode": {"upper"}},
			[]string{"HELLO"},
		},
		{
			"Diff", "/tools/diff",
			url.Values{"text-a": {"hello"}, "text-b": {"world"}},
			nil, // just check 200 — diff output varies
		},
		{
			"XML format", "/tools/xml",
			url.Values{"input": {"<root><child/></root>"}, "mode": {"format"}},
			[]string{"root", "child"},
		},
		{
			"XML minify", "/tools/xml",
			url.Values{"input": {"<root>\n  <child/>\n</root>"}, "mode": {"minify"}},
			[]string{"root", "child"},
		},
		{
			"CSV json-to-csv", "/tools/csv",
			url.Values{"input": {`[{"name":"alice","age":"30"}]`}, "mode": {"json-to-csv"}},
			[]string{"alice", "30"},
		},
		{
			"CSV csv-to-json", "/tools/csv",
			url.Values{"input": {"name,age\nalice,30"}, "mode": {"csv-to-json"}},
			[]string{"alice", "30"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			body := tt.form.Encode()
			req := httptest.NewRequest(http.MethodPost, tt.path, strings.NewReader(body))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code, "POST %s returned %d", tt.path, rec.Code)
			assert.NotEmpty(t, rec.Body.String(), "POST %s returned empty body", tt.path)

			respBody := rec.Body.String()
			for _, s := range tt.contains {
				assert.Contains(t, respBody, s, "POST %s missing %q in response", tt.path, s)
			}
		})
	}
}

func TestNonexistentRouteReturns404(t *testing.T) {
	router := NewRouter()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/tools/nonexistent", nil)
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestStaticFilesServed(t *testing.T) {
	router := NewRouter()

	paths := []string{"/static/style.css", "/static/forge.js"}
	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, path, nil)
			router.ServeHTTP(rec, req)
			require.Equal(t, http.StatusOK, rec.Code, "static file %s not served", path)
			assert.NotEmpty(t, rec.Body.String())
		})
	}
}

func TestPostWithEmptyInputDoesNotPanic(t *testing.T) {
	router := NewRouter()

	paths := []string{
		"/tools/base64", "/tools/jwt", "/tools/json", "/tools/hash",
		"/tools/url", "/tools/uuid", "/tools/yaml", "/tools/timestamp",
		"/tools/number-base", "/tools/regex", "/tools/html-entity",
		"/tools/password", "/tools/lorem", "/tools/color", "/tools/cron",
		"/tools/text-escape", "/tools/gzip", "/tools/text-stats",
		"/tools/diff", "/tools/xml", "/tools/csv",
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(""))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			router.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code, "POST %s with empty input returned %d", path, rec.Code)
		})
	}
}
