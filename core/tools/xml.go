package tools

import (
	"bytes"
	"encoding/xml"
	"io"
	"strings"
)

// XMLTool provides metadata for the XML Formatter tool.
type XMLTool struct{}

func (XMLTool) Name() string        { return "XML Formatter" }
func (XMLTool) ID() string          { return "xml" }
func (XMLTool) Description() string { return "Format, pretty-print, and minify XML" }
func (XMLTool) Category() string    { return "Formatters" }
func (XMLTool) Keywords() []string  { return []string{"xml", "format", "pretty", "minify"} }

// DetectFromClipboard returns true if the input looks like XML — starts with
// '<' and contains '>'.
func (XMLTool) DetectFromClipboard(s string) bool {
	s = strings.TrimSpace(s)
	return len(s) > 0 && strings.HasPrefix(s, "<") && strings.Contains(s, ">")
}

// XMLFormat pretty-prints XML with 2-space indentation. It decodes all tokens
// from the input and re-encodes them with indentation.
func XMLFormat(input string) Result {
	input = strings.TrimSpace(input)
	if input == "" {
		return Result{Error: "empty input"}
	}

	tokens, err := decodeXMLTokens(input)
	if err != nil {
		return Result{Error: "invalid XML: " + err.Error()}
	}

	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	for _, tok := range tokens {
		if err := enc.EncodeToken(tok); err != nil {
			return Result{Error: "encode error: " + err.Error()}
		}
	}
	if err := enc.Flush(); err != nil {
		return Result{Error: "encode error: " + err.Error()}
	}

	return Result{Output: buf.String()}
}

// XMLMinify removes all unnecessary whitespace from XML. It decodes and
// re-encodes tokens with no indentation.
func XMLMinify(input string) Result {
	input = strings.TrimSpace(input)
	if input == "" {
		return Result{Error: "empty input"}
	}

	tokens, err := decodeXMLTokens(input)
	if err != nil {
		return Result{Error: "invalid XML: " + err.Error()}
	}

	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	for _, tok := range tokens {
		if err := enc.EncodeToken(tok); err != nil {
			return Result{Error: "encode error: " + err.Error()}
		}
	}
	if err := enc.Flush(); err != nil {
		return Result{Error: "encode error: " + err.Error()}
	}

	return Result{Output: buf.String()}
}

// decodeXMLTokens reads all tokens from the input, copying them so they
// survive past the next Decoder.Token() call. Pure whitespace CharData
// tokens are dropped so the encoder can apply its own indentation.
func decodeXMLTokens(input string) ([]xml.Token, error) {
	dec := xml.NewDecoder(strings.NewReader(input))
	var tokens []xml.Token
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		// Skip pure-whitespace text nodes so the encoder controls indentation.
		if cd, ok := tok.(xml.CharData); ok {
			if strings.TrimSpace(string(cd)) == "" {
				continue
			}
		}
		tokens = append(tokens, xml.CopyToken(tok))
	}
	return tokens, nil
}
