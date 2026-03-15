package tools

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// timestampDigitsRegex matches strings that are exactly 10 or 13 digits.
var timestampDigitsRegex = regexp.MustCompile(`^\d{10}(\d{3})?$`)

// TimestampTool provides metadata for the Timestamp converter tool.
type TimestampTool struct{}

func (ts TimestampTool) Name() string        { return "Timestamp Converter" }
func (ts TimestampTool) ID() string          { return "timestamp" }
func (ts TimestampTool) Description() string { return "Convert between Unix timestamps and human-readable dates" }
func (ts TimestampTool) Category() string    { return "Converters" }
func (ts TimestampTool) Keywords() []string {
	return []string{"timestamp", "unix", "date", "time", "epoch"}
}

// DetectFromClipboard returns true if s is exactly 10 or 13 digits (unix seconds or millis).
func (ts TimestampTool) DetectFromClipboard(s string) bool {
	s = strings.TrimSpace(s)
	return timestampDigitsRegex.MatchString(s)
}

// commonDateFormats lists datetime formats to try when parsing human-readable dates.
var commonDateFormats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02",
	time.RFC1123,
	time.RFC1123Z,
	time.RFC850,
	time.ANSIC,
	time.UnixDate,
}

// TimestampFromUnix converts a unix timestamp (seconds or milliseconds) to a
// human-readable datetime in the specified timezone.
// If tz is empty, UTC is used.
// Input of 10 digits is treated as seconds, 13 digits as milliseconds.
func TimestampFromUnix(input string, tz string) Result {
	input = strings.TrimSpace(input)
	if input == "" {
		return Result{Error: "empty input"}
	}

	if tz == "" {
		tz = "UTC"
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return Result{Error: fmt.Sprintf("invalid timezone: %s", err.Error())}
	}

	n, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return Result{Error: fmt.Sprintf("invalid unix timestamp: %s", err.Error())}
	}

	var t time.Time
	switch len(input) {
	case 13:
		// Milliseconds
		t = time.UnixMilli(n).In(loc)
	default:
		// Seconds
		t = time.Unix(n, 0).In(loc)
	}

	output := fmt.Sprintf("RFC3339:  %s\nUTC:      %s\nLocal:    %s\nUnix:     %d\nUnix ms:  %d",
		t.Format(time.RFC3339),
		t.UTC().Format("2006-01-02 15:04:05 MST"),
		t.Format("2006-01-02 15:04:05 MST"),
		t.Unix(),
		t.UnixMilli(),
	)

	return Result{Output: output}
}

// TimestampToUnix parses a datetime string in RFC3339 or common formats and
// returns the unix timestamp. If millis is true, returns milliseconds.
func TimestampToUnix(input string, millis bool) Result {
	input = strings.TrimSpace(input)
	if input == "" {
		return Result{Error: "empty input"}
	}

	var t time.Time
	var parseErr error
	for _, format := range commonDateFormats {
		t, parseErr = time.Parse(format, input)
		if parseErr == nil {
			break
		}
	}

	if parseErr != nil {
		return Result{Error: fmt.Sprintf("unable to parse datetime: %s", parseErr.Error())}
	}

	if millis {
		return Result{Output: strconv.FormatInt(t.UnixMilli(), 10)}
	}
	return Result{Output: strconv.FormatInt(t.Unix(), 10)}
}

// TimestampNow returns the current time in multiple formats.
// If tz is empty, UTC is used.
func TimestampNow(tz string) Result {
	if tz == "" {
		tz = "UTC"
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return Result{Error: fmt.Sprintf("invalid timezone: %s", err.Error())}
	}

	now := time.Now().In(loc)

	output := fmt.Sprintf("Unix:     %d\nUnix ms:  %d\nRFC3339:  %s\nHuman:    %s",
		now.Unix(),
		now.UnixMilli(),
		now.Format(time.RFC3339),
		now.Format("Mon, 02 Jan 2006 15:04:05 MST"),
	)

	return Result{Output: output}
}
