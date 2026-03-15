package tools

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
)

// base64StdRegex matches strings that look like standard Base64.
var base64StdRegex = regexp.MustCompile(`^[A-Za-z0-9+/=]+$`)

// Base64Tool provides metadata for the Base64 encode/decode tool.
type Base64Tool struct{}

func (b Base64Tool) Name() string        { return "Base64 Encode / Decode" }
func (b Base64Tool) ID() string          { return "base64" }
func (b Base64Tool) Description() string { return "Encode and decode Base64 strings" }
func (b Base64Tool) Category() string    { return "Encoders" }
func (b Base64Tool) Keywords() []string  { return []string{"base64", "encode", "decode", "b64"} }

// DetectFromClipboard returns true if s looks like a valid Base64-encoded string.
// Requires standard alphabet only, minimum 4 characters, and length divisible by 4.
func (b Base64Tool) DetectFromClipboard(s string) bool {
	s = strings.TrimSpace(s)
	if len(s) < 4 {
		return false
	}
	if len(s)%4 != 0 {
		return false
	}
	return base64StdRegex.MatchString(s)
}

// Base64Encode encodes the input string to Base64.
// If urlSafe is true, RFC 4648 section 5 URL-safe encoding is used.
// If noPadding is true, padding characters (=) are omitted.
func Base64Encode(input string, urlSafe bool, noPadding bool) Result {
	if input == "" {
		return Result{Output: ""}
	}

	var encoding *base64.Encoding
	if urlSafe {
		encoding = base64.URLEncoding
	} else {
		encoding = base64.StdEncoding
	}

	if noPadding {
		encoding = encoding.WithPadding(base64.NoPadding)
	}

	encoded := encoding.EncodeToString([]byte(input))
	return Result{Output: encoded}
}

// Base64Decode decodes a Base64-encoded string.
// If urlSafe is true, RFC 4648 section 5 URL-safe decoding is used.
// It first tries decoding with padding, then without (stripping trailing =).
// On failure, it returns a Result with Error set.
func Base64Decode(input string, urlSafe bool) Result {
	if input == "" {
		return Result{Output: ""}
	}

	var encoding *base64.Encoding
	if urlSafe {
		encoding = base64.URLEncoding
	} else {
		encoding = base64.StdEncoding
	}

	// Try with padding first.
	decoded, err := encoding.DecodeString(input)
	if err == nil {
		return Result{Output: string(decoded)}
	}

	// Try without padding (strip trailing =).
	stripped := strings.TrimRight(input, "=")
	noPadEncoding := encoding.WithPadding(base64.NoPadding)
	decoded, err = noPadEncoding.DecodeString(stripped)
	if err == nil {
		return Result{Output: string(decoded)}
	}

	return Result{Error: fmt.Sprintf("invalid base64: %s", err.Error())}
}
