package tools

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// YAMLTool provides metadata for the JSON / YAML converter tool.
type YAMLTool struct{}

func (y YAMLTool) Name() string        { return "JSON / YAML Converter" }
func (y YAMLTool) ID() string          { return "yaml" }
func (y YAMLTool) Description() string { return "Convert between JSON and YAML formats" }
func (y YAMLTool) Category() string    { return "Converters" }
func (y YAMLTool) Keywords() []string  { return []string{"yaml", "json", "convert"} }

// DetectFromClipboard returns true if s looks like YAML (but not JSON).
// It tries to unmarshal as YAML, and rejects inputs that start with { or [
// to avoid detecting JSON as YAML.
func (y YAMLTool) DetectFromClipboard(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	if strings.HasPrefix(s, "{") || strings.HasPrefix(s, "[") {
		return false
	}
	var out interface{}
	err := yaml.Unmarshal([]byte(s), &out)
	if err != nil {
		return false
	}
	// Reject plain scalars — a bare word like "hello" unmarshals fine
	// but isn't meaningful YAML for our purposes.
	// We require at least a mapping or sequence.
	switch out.(type) {
	case map[string]interface{}:
		return true
	case []interface{}:
		return true
	}
	return false
}

// YAMLToJSON converts a YAML string to JSON.
// If compact is true, the JSON output has no indentation.
func YAMLToJSON(input string, compact bool) Result {
	if strings.TrimSpace(input) == "" {
		return Result{Error: "empty input"}
	}

	var data interface{}
	if err := yaml.Unmarshal([]byte(input), &data); err != nil {
		return Result{Error: fmt.Sprintf("invalid YAML: %s", err.Error())}
	}

	// yaml.v3 unmarshals map keys as string, which is compatible with
	// encoding/json, but we need to recursively convert map[string]interface{}
	// in case of nested maps.
	data = normalizeYAMLValue(data)

	var jsonBytes []byte
	var err error
	if compact {
		jsonBytes, err = json.Marshal(data)
	} else {
		jsonBytes, err = json.MarshalIndent(data, "", "  ")
	}
	if err != nil {
		return Result{Error: fmt.Sprintf("JSON marshal error: %s", err.Error())}
	}

	return Result{Output: string(jsonBytes)}
}

// JSONToYAML converts a JSON string to YAML.
func JSONToYAML(input string) Result {
	if strings.TrimSpace(input) == "" {
		return Result{Error: "empty input"}
	}

	var data interface{}
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return Result{Error: fmt.Sprintf("invalid JSON: %s", err.Error())}
	}

	yamlBytes, err := yaml.Marshal(data)
	if err != nil {
		return Result{Error: fmt.Sprintf("YAML marshal error: %s", err.Error())}
	}

	return Result{Output: strings.TrimRight(string(yamlBytes), "\n")}
}

// normalizeYAMLValue recursively converts yaml.v3 output types to types
// that encoding/json can handle. In particular, yaml.v3 may produce
// map[string]interface{} which is fine, but nested values need traversal.
func normalizeYAMLValue(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(val))
		for k, v2 := range val {
			out[k] = normalizeYAMLValue(v2)
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(val))
		for i, v2 := range val {
			out[i] = normalizeYAMLValue(v2)
		}
		return out
	default:
		return v
	}
}
