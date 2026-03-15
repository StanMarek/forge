package tools

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// CSVTool provides metadata for the JSON/CSV converter tool.
type CSVTool struct{}

func (CSVTool) Name() string        { return "JSON to CSV" }
func (CSVTool) ID() string          { return "csv" }
func (CSVTool) Description() string { return "Convert between JSON arrays and CSV" }
func (CSVTool) Category() string    { return "Converters" }
func (CSVTool) Keywords() []string  { return []string{"csv", "json", "table", "tsv", "convert"} }

// DetectFromClipboard always returns false — conversion requires explicit action.
func (CSVTool) DetectFromClipboard(_ string) bool {
	return false
}

// JSONToCSV converts a JSON array of objects to CSV text. The delimiter
// parameter controls the field separator (default ","). Keys are sorted
// alphabetically to produce deterministic output. Nested values are
// stringified as JSON.
func JSONToCSV(input string, delimiter string) Result {
	input = strings.TrimSpace(input)
	if input == "" {
		return Result{Error: "empty input"}
	}

	if delimiter == "" {
		delimiter = ","
	}
	delim := rune(delimiter[0])

	var rows []map[string]interface{}
	if err := json.Unmarshal([]byte(input), &rows); err != nil {
		return Result{Error: "input must be a JSON array of objects: " + err.Error()}
	}

	if len(rows) == 0 {
		return Result{Output: ""}
	}

	// Collect all unique keys, sorted alphabetically.
	keySet := make(map[string]struct{})
	for _, row := range rows {
		for k := range row {
			keySet[k] = struct{}{}
		}
	}
	headers := make([]string, 0, len(keySet))
	for k := range keySet {
		headers = append(headers, k)
	}
	sort.Strings(headers)

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	w.Comma = delim

	// Write header row.
	if err := w.Write(headers); err != nil {
		return Result{Error: "csv write error: " + err.Error()}
	}

	// Write data rows.
	for _, row := range rows {
		record := make([]string, len(headers))
		for i, h := range headers {
			record[i] = stringify(row[h])
		}
		if err := w.Write(record); err != nil {
			return Result{Error: "csv write error: " + err.Error()}
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return Result{Error: "csv write error: " + err.Error()}
	}

	return Result{Output: strings.TrimRight(buf.String(), "\n")}
}

// CSVToJSON converts CSV text to a JSON array of objects. The first row is
// treated as headers. The delimiter parameter controls the field separator
// (default ",").
func CSVToJSON(input string, delimiter string) Result {
	input = strings.TrimSpace(input)
	if input == "" {
		return Result{Error: "empty input"}
	}

	if delimiter == "" {
		delimiter = ","
	}
	delim := rune(delimiter[0])

	r := csv.NewReader(strings.NewReader(input))
	r.Comma = delim

	records, err := r.ReadAll()
	if err != nil {
		return Result{Error: "csv parse error: " + err.Error()}
	}

	if len(records) < 2 {
		// Only headers or empty — return empty array.
		return Result{Output: "[]"}
	}

	headers := records[0]
	var result []map[string]interface{}

	for _, row := range records[1:] {
		obj := make(map[string]interface{}, len(headers))
		for i, h := range headers {
			if i < len(row) {
				obj[h] = row[i]
			} else {
				obj[h] = ""
			}
		}
		result = append(result, obj)
	}

	out, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return Result{Error: "json marshal error: " + err.Error()}
	}

	return Result{Output: string(out)}
}

// stringify converts an interface{} value to a string. Nested objects and
// arrays are serialised as compact JSON. Nil becomes an empty string.
func stringify(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		// Avoid trailing zeros for integers.
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%g", val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		// Nested object or array — serialise as JSON.
		b, err := json.Marshal(val)
		if err != nil {
			return fmt.Sprintf("%v", val)
		}
		return string(b)
	}
}
