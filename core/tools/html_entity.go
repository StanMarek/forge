package tools

import (
	"html"
	"strings"
)

// HTMLEntityTool provides metadata for the HTML Entity Encoder tool.
type HTMLEntityTool struct{}

func (h HTMLEntityTool) Name() string        { return "HTML Entity Encoder" }
func (h HTMLEntityTool) ID() string          { return "html-entity" }
func (h HTMLEntityTool) Description() string { return "Encode and decode HTML entities" }
func (h HTMLEntityTool) Category() string    { return "Encoders" }
func (h HTMLEntityTool) Keywords() []string {
	return []string{"html", "entity", "encode", "decode", "escape"}
}

// DetectFromClipboard returns true if s contains HTML entity patterns
// such as &amp;, &lt;, &gt;, &quot;, or &#-prefixed sequences.
func (h HTMLEntityTool) DetectFromClipboard(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	return strings.Contains(s, "&amp;") ||
		strings.Contains(s, "&lt;") ||
		strings.Contains(s, "&gt;") ||
		strings.Contains(s, "&quot;") ||
		strings.Contains(s, "&#")
}

// HTMLEntityEncode encodes special HTML characters in input using html.EscapeString.
func HTMLEntityEncode(input string) Result {
	return Result{Output: html.EscapeString(input)}
}

// HTMLEntityDecode decodes HTML entities in input using html.UnescapeString.
func HTMLEntityDecode(input string) Result {
	return Result{Output: html.UnescapeString(input)}
}
