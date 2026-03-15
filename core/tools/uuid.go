package tools

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// uuidRegex matches a standard UUID with hyphens.
var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// UUIDTool provides metadata for the UUID Generate / Validate / Parse tool.
type UUIDTool struct{}

func (u UUIDTool) Name() string        { return "UUID Generate / Validate / Parse" }
func (u UUIDTool) ID() string          { return "uuid" }
func (u UUIDTool) Description() string { return "Generate, validate, and parse UUIDs" }
func (u UUIDTool) Category() string    { return "Generators" }
func (u UUIDTool) Keywords() []string {
	return []string{"uuid", "guid", "generate", "v4", "v7", "unique"}
}

// DetectFromClipboard returns true if the string looks like a UUID.
func (u UUIDTool) DetectFromClipboard(s string) bool {
	return uuidRegex.MatchString(strings.TrimSpace(s))
}

// UUIDParseResult holds the parsed components of a UUID.
type UUIDParseResult struct {
	UUID      string
	Version   int
	Variant   string
	Timestamp string
	Output    string
	Error     string
}

// UUIDGenerate creates a new UUID of the specified version.
// Supported versions: 4 (random), 7 (time-ordered).
// If uppercase is true, the output is uppercased.
// If noHyphens is true, hyphens are removed from the output.
func UUIDGenerate(version int, uppercase bool, noHyphens bool) Result {
	var generated uuid.UUID
	var err error

	switch version {
	case 4:
		generated, err = uuid.NewRandom()
		if err != nil {
			return Result{Error: fmt.Sprintf("failed to generate UUID v4: %v", err)}
		}
	case 7:
		generated, err = uuid.NewV7()
		if err != nil {
			return Result{Error: fmt.Sprintf("failed to generate UUID v7: %v", err)}
		}
	default:
		return Result{Error: fmt.Sprintf("unsupported UUID version: %d (supported: 4, 7)", version)}
	}

	output := generated.String()

	if noHyphens {
		output = strings.ReplaceAll(output, "-", "")
	}

	if uppercase {
		output = strings.ToUpper(output)
	}

	return Result{Output: output}
}

// UUIDValidate checks whether the input is a valid UUID.
// Returns "valid (version N)" on success or an Error.
func UUIDValidate(input string) Result {
	input = strings.TrimSpace(input)
	parsed, err := uuid.Parse(input)
	if err != nil {
		return Result{Error: fmt.Sprintf("invalid UUID: %v", err)}
	}
	return Result{Output: fmt.Sprintf("valid (version %d)", parsed.Version())}
}

// UUIDParse parses a UUID string and returns its components.
func UUIDParse(input string) UUIDParseResult {
	input = strings.TrimSpace(input)
	parsed, err := uuid.Parse(input)
	if err != nil {
		return UUIDParseResult{Error: fmt.Sprintf("invalid UUID: %v", err)}
	}

	ver := int(parsed.Version())
	variant := variantString(parsed.Variant())

	var timestamp string
	if ver == 7 {
		sec, nsec := parsed.Time().UnixTime()
		t := time.Unix(sec, int64(nsec))
		timestamp = t.UTC().Format(time.RFC3339)
	}

	output := fmt.Sprintf("UUID:    %s\nVersion: %d\nVariant: %s", parsed.String(), ver, variant)
	if timestamp != "" {
		output += fmt.Sprintf("\nTime:    %s", timestamp)
	}

	return UUIDParseResult{
		UUID:      parsed.String(),
		Version:   ver,
		Variant:   variant,
		Timestamp: timestamp,
		Output:    output,
	}
}

// variantString maps a uuid.Variant to a human-readable string.
func variantString(v uuid.Variant) string {
	switch v {
	case uuid.RFC4122:
		return "RFC 4122"
	case uuid.Reserved:
		return "Reserved"
	case uuid.Microsoft:
		return "Microsoft"
	case uuid.Future:
		return "Future"
	default:
		return "Unknown"
	}
}
