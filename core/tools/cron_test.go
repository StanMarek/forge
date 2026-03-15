package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------- Tool metadata ----------

func TestCronTool_Metadata(t *testing.T) {
	tool := CronTool{}

	assert.Equal(t, "Cron Expression Parser", tool.Name())
	assert.Equal(t, "cron", tool.ID())
	assert.Equal(t, "Converters", tool.Category())
	assert.NotEmpty(t, tool.Description())

	keywords := tool.Keywords()
	for _, kw := range []string{"cron", "schedule", "crontab"} {
		assert.Contains(t, keywords, kw)
	}
}

func TestCronTool_DetectFromClipboard(t *testing.T) {
	tool := CronTool{}
	assert.True(t, tool.DetectFromClipboard("* * * * *"))
	assert.True(t, tool.DetectFromClipboard("*/5 * * * *"))
	assert.True(t, tool.DetectFromClipboard("0 9 * * 1-5"))
	assert.False(t, tool.DetectFromClipboard("hello world"))
	assert.False(t, tool.DetectFromClipboard("* * *"))
	assert.False(t, tool.DetectFromClipboard(""))
}

// ---------- Every minute ----------

func TestCronParse_EveryMinute(t *testing.T) {
	r := CronParse("* * * * *")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "Expression: * * * * *")
	assert.Contains(t, r.Output, "Schedule:   Every minute")
	assert.Contains(t, r.Output, "Minute:     Every minute")
	assert.Contains(t, r.Output, "Hour:       Every hour")
	assert.Contains(t, r.Output, "Day:        Every day")
	assert.Contains(t, r.Output, "Month:      Every month")
	assert.Contains(t, r.Output, "Weekday:    Every day of the week")
}

// ---------- Every 5 minutes ----------

func TestCronParse_Every5Min(t *testing.T) {
	r := CronParse("*/5 * * * *")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "Schedule:   Every 5 minutes")
	assert.Contains(t, r.Output, "Minute:     Every 5 minutes")
}

// ---------- 9am weekdays ----------

func TestCronParse_9amWeekdays(t *testing.T) {
	r := CronParse("0 9 * * 1-5")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "Expression: 0 9 * * 1-5")
	assert.Contains(t, r.Output, "At 09:00")
	assert.Contains(t, r.Output, "Monday through Friday")
	assert.Contains(t, r.Output, "Minute:     At minute 0")
	assert.Contains(t, r.Output, "Hour:       At hour 9")
	assert.Contains(t, r.Output, "Weekday:    Monday through Friday")
}

// ---------- Every 2 hours at minute 30 ----------

func TestCronParse_Every2HoursAt30(t *testing.T) {
	r := CronParse("30 */2 * * *")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "Expression: 30 */2 * * *")
	assert.Contains(t, r.Output, "Minute:     At minute 30")
	assert.Contains(t, r.Output, "Hour:       Every 2 hours")
}

// ---------- Specific time daily ----------

func TestCronParse_DailyAt0930(t *testing.T) {
	r := CronParse("30 9 * * *")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "Schedule:   At 09:30 every day")
}

// ---------- Invalid: too few fields ----------

func TestCronParse_TooFewFields(t *testing.T) {
	r := CronParse("* * *")
	assert.NotEmpty(t, r.Error)
	assert.Contains(t, r.Error, "expected 5 fields, got 3")
	assert.Empty(t, r.Output)
}

// ---------- Invalid: too many fields ----------

func TestCronParse_TooManyFields(t *testing.T) {
	r := CronParse("* * * * * *")
	assert.NotEmpty(t, r.Error)
	assert.Contains(t, r.Error, "expected 5 fields, got 6")
	assert.Empty(t, r.Output)
}

// ---------- Invalid: bad values ----------

func TestCronParse_BadMinuteValue(t *testing.T) {
	r := CronParse("60 * * * *")
	assert.NotEmpty(t, r.Error)
	assert.Contains(t, r.Error, "minute")
	assert.Empty(t, r.Output)
}

func TestCronParse_BadHourValue(t *testing.T) {
	r := CronParse("0 25 * * *")
	assert.NotEmpty(t, r.Error)
	assert.Contains(t, r.Error, "hour")
	assert.Empty(t, r.Output)
}

func TestCronParse_NonNumeric(t *testing.T) {
	r := CronParse("abc * * * *")
	assert.NotEmpty(t, r.Error)
	assert.Contains(t, r.Error, "minute")
	assert.Empty(t, r.Output)
}

// ---------- Empty input ----------

func TestCronParse_EmptyInput(t *testing.T) {
	r := CronParse("")
	assert.NotEmpty(t, r.Error)
	assert.Empty(t, r.Output)
}

// ---------- List values ----------

func TestCronParse_ListWeekday(t *testing.T) {
	r := CronParse("0 9 * * 1,3,5")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "Weekday:    Monday, Wednesday, Friday")
}

// ---------- Tool interface compliance ----------

func TestCronTool_ImplementsToolInterface(t *testing.T) {
	var _ Tool = CronTool{}
}
