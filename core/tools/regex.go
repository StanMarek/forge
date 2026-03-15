package tools

import (
	"fmt"
	"regexp"
	"strings"
)

// RegexTool provides metadata for the Regex Tester tool.
type RegexTool struct{}

func (r RegexTool) Name() string        { return "Regex Tester" }
func (r RegexTool) ID() string          { return "regex" }
func (r RegexTool) Description() string { return "Test regular expressions against input strings" }
func (r RegexTool) Category() string    { return "Testers" }
func (r RegexTool) Keywords() []string {
	return []string{"regex", "regexp", "pattern", "match", "test"}
}

// DetectFromClipboard always returns false for the regex tool.
// Regex patterns are too ambiguous to detect reliably.
func (r RegexTool) DetectFromClipboard(_ string) bool {
	return false
}

// RegexTest compiles pattern and tests it against input.
// If global is true, all matches are returned via FindAllString.
// If global is false, only the first match is returned via FindString.
func RegexTest(pattern string, input string, global bool) Result {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return Result{Error: fmt.Sprintf("invalid regex: %s", err.Error())}
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Pattern: %s\n", pattern)
	fmt.Fprintf(&b, "Input: %q\n", input)

	if global {
		matches := re.FindAllStringIndex(input, -1)
		if len(matches) == 0 {
			return Result{Output: "No matches found"}
		}
		fmt.Fprintln(&b, "Matches:")
		for i, loc := range matches {
			matched := input[loc[0]:loc[1]]
			fmt.Fprintf(&b, "  [%d] %q (pos %d-%d)", i, matched, loc[0], loc[1])
			if i < len(matches)-1 {
				fmt.Fprintln(&b)
			}
		}
	} else {
		loc := re.FindStringIndex(input)
		if loc == nil {
			return Result{Output: "No matches found"}
		}
		matched := input[loc[0]:loc[1]]
		fmt.Fprintf(&b, "Match: %q (pos %d-%d)", matched, loc[0], loc[1])
	}

	return Result{Output: b.String()}
}
