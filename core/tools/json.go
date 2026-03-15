package tools

import (
	"bytes"
	"encoding/json"
	"sort"
	"strings"
)

// JSONTool provides JSON formatting, minification, and validation.
type JSONTool struct{}

func (JSONTool) Name() string        { return "JSON Formatter" }
func (JSONTool) ID() string          { return "json" }
func (JSONTool) Description() string { return "Format, minify, and validate JSON" }
func (JSONTool) Category() string    { return "Formatters" }
func (JSONTool) Keywords() []string  { return []string{"json", "format", "minify", "validate", "pretty"} }

func (JSONTool) DetectFromClipboard(s string) bool {
	s = strings.TrimSpace(s)
	return len(s) > 0 && json.Valid([]byte(s))
}

// JSONFormat pretty-prints JSON. When sortKeys is true, keys are sorted
// recursively. When useTabs is true, indentation uses tabs instead of spaces.
// Without sortKeys, json.Indent is used to preserve original key order.
func JSONFormat(input string, indent int, sortKeys bool, useTabs bool) Result {
	input = strings.TrimSpace(input)
	if input == "" {
		return Result{Error: "empty input"}
	}

	indentStr := strings.Repeat(" ", indent)
	if useTabs {
		indentStr = "\t"
	}

	if sortKeys {
		var raw interface{}
		if err := json.Unmarshal([]byte(input), &raw); err != nil {
			return Result{Error: "invalid JSON: " + err.Error()}
		}
		raw = sortKeysRecursive(raw)
		out, err := json.MarshalIndent(raw, "", indentStr)
		if err != nil {
			return Result{Error: "marshal error: " + err.Error()}
		}
		return Result{Output: string(out)}
	}

	var buf bytes.Buffer
	if err := json.Indent(&buf, []byte(input), "", indentStr); err != nil {
		return Result{Error: "invalid JSON: " + err.Error()}
	}
	return Result{Output: buf.String()}
}

// JSONMinify removes all whitespace from JSON using json.Compact.
func JSONMinify(input string) Result {
	input = strings.TrimSpace(input)
	if input == "" {
		return Result{Error: "empty input"}
	}

	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(input)); err != nil {
		return Result{Error: "invalid JSON: " + err.Error()}
	}
	return Result{Output: buf.String()}
}

// JSONValidate checks whether the input is valid JSON.
// Returns Output "valid" on success, or Error with details on failure.
func JSONValidate(input string) Result {
	input = strings.TrimSpace(input)
	if input == "" {
		return Result{Error: "empty input"}
	}

	if json.Valid([]byte(input)) {
		return Result{Output: "valid"}
	}

	detail := findJSONError(input)
	return Result{Error: "invalid JSON: " + detail}
}

// findJSONError returns the error message from json.Unmarshal.
func findJSONError(input string) string {
	var raw interface{}
	if err := json.Unmarshal([]byte(input), &raw); err != nil {
		return err.Error()
	}
	return "unknown error"
}

// sortKeysRecursive walks a decoded JSON value and replaces every
// map[string]interface{} with a sortedMap that marshals keys in order.
func sortKeysRecursive(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		sm := &sortedMap{keys: keys, values: make(map[string]interface{}, len(val))}
		for _, k := range keys {
			sm.values[k] = sortKeysRecursive(val[k])
		}
		return sm
	case []interface{}:
		for i, item := range val {
			val[i] = sortKeysRecursive(item)
		}
		return val
	default:
		return v
	}
}

// sortedMap is a JSON-marshalable ordered map that emits keys in a
// predetermined order.
type sortedMap struct {
	keys   []string
	values map[string]interface{}
}

func (sm *sortedMap) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, k := range sm.keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		keyBytes, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}
		buf.Write(keyBytes)
		buf.WriteByte(':')
		valBytes, err := json.Marshal(sm.values[k])
		if err != nil {
			return nil, err
		}
		buf.Write(valBytes)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}
