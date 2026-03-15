package tools

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// JWTTool provides metadata for the JWT Decoder tool.
type JWTTool struct{}

func (j JWTTool) Name() string        { return "JWT Decoder" }
func (j JWTTool) ID() string          { return "jwt" }
func (j JWTTool) Description() string { return "Decode and inspect JSON Web Tokens" }
func (j JWTTool) Category() string    { return "Encoders" }
func (j JWTTool) Keywords() []string  { return []string{"jwt", "token", "decode", "json web token"} }

// DetectFromClipboard returns true if the string looks like a JWT (3 dot-separated non-empty segments).
func (j JWTTool) DetectFromClipboard(s string) bool {
	parts := strings.Split(strings.TrimSpace(s), ".")
	if len(parts) != 3 {
		return false
	}
	for _, p := range parts {
		if p == "" {
			return false
		}
	}
	return true
}

// JWTDecodeResult holds the decoded components of a JWT.
type JWTDecodeResult struct {
	Header    string
	Payload   string
	Signature string
	Output    string
	Error     string
}

// JWTDecode splits a JWT into its three parts, base64url-decodes the header
// and payload, and pretty-prints the JSON.
func JWTDecode(token string) JWTDecodeResult {
	token = strings.TrimSpace(token)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return JWTDecodeResult{Error: "invalid JWT: expected 3 dot-separated segments"}
	}
	for i, p := range parts {
		if p == "" {
			return JWTDecodeResult{Error: fmt.Sprintf("invalid JWT: segment %d is empty", i+1)}
		}
	}

	header, err := decodeJWTSegment(parts[0])
	if err != nil {
		return JWTDecodeResult{Error: fmt.Sprintf("invalid JWT header: %v", err)}
	}

	payload, err := decodeJWTSegment(parts[1])
	if err != nil {
		return JWTDecodeResult{Error: fmt.Sprintf("invalid JWT payload: %v", err)}
	}

	sig := parts[2]

	output := fmt.Sprintf("--- Header ---\n%s\n--- Payload ---\n%s\n--- Signature ---\n%s", header, payload, sig)

	return JWTDecodeResult{
		Header:    header,
		Payload:   payload,
		Signature: sig,
		Output:    output,
	}
}

// JWTValidate checks whether a token string is a structurally valid JWT.
// It verifies 3 dot-separated segments, valid base64url encoding, and valid
// JSON in the header and payload. Returns Output:"valid" on success or an Error.
func JWTValidate(token string) Result {
	token = strings.TrimSpace(token)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Result{Error: "invalid JWT: expected 3 dot-separated segments"}
	}
	for i, p := range parts {
		if p == "" {
			return Result{Error: fmt.Sprintf("invalid JWT: segment %d is empty", i+1)}
		}
	}

	// Validate base64url encoding for all three segments.
	for i, p := range parts {
		if _, err := base64.RawURLEncoding.DecodeString(p); err != nil {
			return Result{Error: fmt.Sprintf("invalid JWT: segment %d is not valid base64url: %v", i+1, err)}
		}
	}

	// Validate that header and payload are valid JSON.
	headerBytes, _ := base64.RawURLEncoding.DecodeString(parts[0])
	if !json.Valid(headerBytes) {
		return Result{Error: "invalid JWT: header is not valid JSON"}
	}

	payloadBytes, _ := base64.RawURLEncoding.DecodeString(parts[1])
	if !json.Valid(payloadBytes) {
		return Result{Error: "invalid JWT: payload is not valid JSON"}
	}

	return Result{Output: "valid"}
}

// decodeJWTSegment base64url-decodes a JWT segment and pretty-prints it as JSON.
func decodeJWTSegment(segment string) (string, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(segment)
	if err != nil {
		return "", fmt.Errorf("base64url decode failed: %w", err)
	}

	pretty, err := formatJSON(string(decoded))
	if err != nil {
		return "", fmt.Errorf("JSON formatting failed: %w", err)
	}

	return pretty, nil
}

// formatJSON pretty-prints a JSON string with 2-space indentation.
func formatJSON(s string) (string, error) {
	var obj interface{}
	if err := json.Unmarshal([]byte(s), &obj); err != nil {
		return "", err
	}
	pretty, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return "", err
	}
	return string(pretty), nil
}
