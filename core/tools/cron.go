package tools

import (
	"fmt"
	"strconv"
	"strings"
)

// CronTool provides metadata for the Cron Expression Parser tool.
type CronTool struct{}

func (c CronTool) Name() string        { return "Cron Expression Parser" }
func (c CronTool) ID() string          { return "cron" }
func (c CronTool) Description() string { return "Parse cron expressions into human-readable descriptions" }
func (c CronTool) Category() string    { return "Converters" }
func (c CronTool) Keywords() []string {
	return []string{"cron", "schedule", "crontab"}
}

// DetectFromClipboard returns true if s looks like a 5-field cron expression.
func (c CronTool) DetectFromClipboard(s string) bool {
	s = strings.TrimSpace(s)
	fields := strings.Fields(s)
	return len(fields) == 5
}

// CronParse parses a 5-field cron expression and returns a human-readable description.
func CronParse(expression string) Result {
	expr := strings.TrimSpace(expression)
	fields := strings.Fields(expr)

	if len(fields) != 5 {
		return Result{
			Error: fmt.Sprintf("invalid cron expression: expected 5 fields, got %d", len(fields)),
		}
	}

	minuteField := fields[0]
	hourField := fields[1]
	dayField := fields[2]
	monthField := fields[3]
	weekdayField := fields[4]

	// Validate each field.
	if err := validateField(minuteField, 0, 59, "minute"); err != "" {
		return Result{Error: err}
	}
	if err := validateField(hourField, 0, 23, "hour"); err != "" {
		return Result{Error: err}
	}
	if err := validateField(dayField, 1, 31, "day"); err != "" {
		return Result{Error: err}
	}
	if err := validateField(monthField, 1, 12, "month"); err != "" {
		return Result{Error: err}
	}
	if err := validateField(weekdayField, 0, 7, "weekday"); err != "" {
		return Result{Error: err}
	}

	minuteDesc := describeMinute(minuteField)
	hourDesc := describeHour(hourField)
	dayDesc := describeDay(dayField)
	monthDesc := describeMonth(monthField)
	weekdayDesc := describeWeekday(weekdayField)

	schedule := buildScheduleSummary(minuteField, hourField, dayField, monthField, weekdayField)

	output := fmt.Sprintf(
		"Expression: %s\nSchedule:   %s\n\nMinute:     %s\nHour:       %s\nDay:        %s\nMonth:      %s\nWeekday:    %s",
		expr, schedule, minuteDesc, hourDesc, dayDesc, monthDesc, weekdayDesc,
	)

	return Result{Output: output}
}

// validateField validates a single cron field.
func validateField(field string, min, max int, name string) string {
	if field == "*" {
		return ""
	}

	// Handle step values: */N or N/M
	if strings.Contains(field, "/") {
		parts := strings.SplitN(field, "/", 2)
		if parts[0] != "*" {
			if err := validateValue(parts[0], min, max, name); err != "" {
				return err
			}
		}
		step, e := strconv.Atoi(parts[1])
		if e != nil || step < 1 {
			return fmt.Sprintf("invalid step value in %s field: %s", name, parts[1])
		}
		return ""
	}

	// Handle lists: N,M,...
	if strings.Contains(field, ",") {
		for _, part := range strings.Split(field, ",") {
			if err := validateValue(strings.TrimSpace(part), min, max, name); err != "" {
				return err
			}
		}
		return ""
	}

	// Handle ranges: N-M
	if strings.Contains(field, "-") {
		parts := strings.SplitN(field, "-", 2)
		if err := validateValue(parts[0], min, max, name); err != "" {
			return err
		}
		if err := validateValue(parts[1], min, max, name); err != "" {
			return err
		}
		return ""
	}

	// Plain number.
	return validateValue(field, min, max, name)
}

// validateValue validates a single numeric value.
func validateValue(s string, min, max int, name string) string {
	v, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Sprintf("invalid value in %s field: %s", name, s)
	}
	if v < min || v > max {
		return fmt.Sprintf("value out of range in %s field: %d (expected %d-%d)", name, v, min, max)
	}
	return ""
}

// describeMinute describes the minute field.
func describeMinute(field string) string {
	if field == "*" {
		return "Every minute"
	}
	if strings.HasPrefix(field, "*/") {
		step := strings.TrimPrefix(field, "*/")
		return fmt.Sprintf("Every %s minutes", step)
	}
	return fmt.Sprintf("At minute %s", field)
}

// describeHour describes the hour field.
func describeHour(field string) string {
	if field == "*" {
		return "Every hour"
	}
	if strings.HasPrefix(field, "*/") {
		step := strings.TrimPrefix(field, "*/")
		return fmt.Sprintf("Every %s hours", step)
	}
	return fmt.Sprintf("At hour %s", field)
}

// describeDay describes the day-of-month field.
func describeDay(field string) string {
	if field == "*" {
		return "Every day"
	}
	return fmt.Sprintf("On day %s", field)
}

// describeMonth describes the month field.
func describeMonth(field string) string {
	if field == "*" {
		return "Every month"
	}
	monthNames := map[int]string{
		1: "January", 2: "February", 3: "March", 4: "April",
		5: "May", 6: "June", 7: "July", 8: "August",
		9: "September", 10: "October", 11: "November", 12: "December",
	}
	v, err := strconv.Atoi(field)
	if err == nil {
		if name, ok := monthNames[v]; ok {
			return fmt.Sprintf("In %s", name)
		}
	}
	return fmt.Sprintf("In month %s", field)
}

// describeWeekday describes the day-of-week field.
func describeWeekday(field string) string {
	if field == "*" {
		return "Every day of the week"
	}
	dayNames := map[int]string{
		0: "Sunday", 1: "Monday", 2: "Tuesday", 3: "Wednesday",
		4: "Thursday", 5: "Friday", 6: "Saturday", 7: "Sunday",
	}
	// Handle range like 1-5
	if strings.Contains(field, "-") {
		parts := strings.SplitN(field, "-", 2)
		start, err1 := strconv.Atoi(parts[0])
		end, err2 := strconv.Atoi(parts[1])
		if err1 == nil && err2 == nil {
			startName := dayNames[start]
			endName := dayNames[end]
			if startName != "" && endName != "" {
				return fmt.Sprintf("%s through %s", startName, endName)
			}
		}
	}
	// Handle list like 1,3,5
	if strings.Contains(field, ",") {
		names := []string{}
		for _, part := range strings.Split(field, ",") {
			v, err := strconv.Atoi(strings.TrimSpace(part))
			if err == nil {
				if name, ok := dayNames[v]; ok {
					names = append(names, name)
					continue
				}
			}
			names = append(names, part)
		}
		return strings.Join(names, ", ")
	}
	v, err := strconv.Atoi(field)
	if err == nil {
		if name, ok := dayNames[v]; ok {
			return fmt.Sprintf("On %s", name)
		}
	}
	return fmt.Sprintf("On weekday %s", field)
}

// buildScheduleSummary generates a human-readable summary for common patterns.
func buildScheduleSummary(minute, hour, day, month, weekday string) string {
	// * * * * * — every minute
	if minute == "*" && hour == "*" && day == "*" && month == "*" && weekday == "*" {
		return "Every minute"
	}

	// */N * * * * — every N minutes
	if strings.HasPrefix(minute, "*/") && hour == "*" && day == "*" && month == "*" && weekday == "*" {
		step := strings.TrimPrefix(minute, "*/")
		return fmt.Sprintf("Every %s minutes", step)
	}

	// N * * * * — at minute N of every hour
	if isNumber(minute) && hour == "*" && day == "*" && month == "*" && weekday == "*" {
		return fmt.Sprintf("At minute %s of every hour", minute)
	}

	// N N * * * — at specific time daily
	if isNumber(minute) && isNumber(hour) && day == "*" && month == "*" && weekday == "*" {
		return fmt.Sprintf("At %s:%s every day", padTime(hour), padTime(minute))
	}

	// N N * * weekdays — at specific time on specific days
	if isNumber(minute) && isNumber(hour) && day == "*" && month == "*" && weekday != "*" {
		return fmt.Sprintf("At %s:%s, %s", padTime(hour), padTime(minute), describeWeekday(weekday))
	}

	// N */N * * * — at minute N every N hours
	if isNumber(minute) && strings.HasPrefix(hour, "*/") && day == "*" && month == "*" && weekday == "*" {
		step := strings.TrimPrefix(hour, "*/")
		return fmt.Sprintf("At minute %s every %s hours", minute, step)
	}

	// Fall back to field-by-field description.
	return fmt.Sprintf("%s, %s, %s, %s, %s",
		describeMinute(minute), describeHour(hour),
		describeDay(day), describeMonth(month), describeWeekday(weekday))
}

// isNumber checks if a string is a plain integer.
func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// padTime pads a time component to 2 digits.
func padTime(s string) string {
	v, err := strconv.Atoi(s)
	if err != nil {
		return s
	}
	return fmt.Sprintf("%02d", v)
}
