# DateTime Tools Documentation

The datetime tools package provides comprehensive date and time operations for LLM agents.

## Overview

All datetime tools are located in `pkg/tools/datetime_tools.go` and implement the go-llms `domain.Tool` interface.

## Available Tools

### get_current_datetime

Gets the current date and time in various formats.

**Tool Name:** `get_current_datetime`

**Parameters:**
```go
type GetCurrentDateTimeParams struct {
    Format   string `json:"format"`
    Timezone string `json:"timezone,omitempty"`
}
```

**Supported Formats:**
- `date` - 2024-01-01
- `date_slash` - 2024/01/01
- `date_compact` - 20240101
- `date_us` - 01/01/2024
- `date_eu` - 01.01.2024
- `time` - 15:04:05
- `time_12` - 3:04:05 PM
- `time_compact` - 150405
- `datetime` - 2024-01-01 15:04:05
- `datetime_iso` - 2024-01-01T15:04:05
- `datetime_rfc` - Mon, 01 Jan 2024 15:04:05 MST
- `datetime_unix` - Unix timestamp
- `datetime_compact` - 20240101150405
- `filename` - 2024-01-01_15-04-05
- `filename_md` - 2024-01-01_15-04-05.md
- `filename_txt` - 2024-01-01_15-04-05.txt
- `filename_log` - 2024-01-01_15-04-05.log
- `log` - [2024-01-01 15:04:05]
- `log_compact` - [20240101-150405]
- `relative` - "2 hours ago", "in 3 days"
- `duration` - Duration since Unix epoch

**Returns:**
```go
type GetCurrentDateTimeResult struct {
    Formatted string `json:"formatted"`
    Format    string `json:"format"`
    Timezone  string `json:"timezone"`
    Unix      int64  `json:"unix"`
}
```

### calculate_duration

Calculates the time difference between two dates.

**Tool Name:** `calculate_duration`

**Parameters:**
```go
type CalculateDurationParams struct {
    Start  string `json:"start"`
    End    string `json:"end"`
    Format string `json:"format,omitempty"`
}
```

**Returns:**
```json
{
    "seconds": 777600,
    "minutes": 12960,
    "hours": 216,
    "days": 9,
    "human": "216h0m0s"
}
```

### get_weekday

Gets the day of the week for a given date.

**Tool Name:** `get_weekday`

**Parameters:**
```go
type GetWeekdayParams struct {
    Date string `json:"date,omitempty"` // Defaults to today
}
```

**Returns:**
```json
{
    "weekday": "Monday",
    "weekday_num": 1,      // 0 = Sunday
    "is_weekend": false,
    "date": "2024-01-01"
}
```

### add_subtract_time

Adds or subtracts time units from a date.

**Tool Name:** `add_subtract_time`

**Parameters:**
```go
type AddSubtractTimeParams struct {
    Date      string `json:"date"`
    Amount    int    `json:"amount"`
    Unit      string `json:"unit"`      // year, month, day, hour, minute, second
    Operation string `json:"operation"` // add, subtract
}
```

**Returns:**
```json
{
    "original": "2024-01-01 00:00:00",
    "result": "2024-01-06 00:00:00",
    "result_iso": "2024-01-06T00:00:00Z",
    "operation": "add 5 day(s)",
    "unix_timestamp": 1704499200
}
```

### convert_timezone

Converts datetime between different timezones.

**Tool Name:** `convert_timezone`

**Parameters:**
```go
type ConvertTimezoneParams struct {
    Datetime     string `json:"datetime"`
    FromTimezone string `json:"from_timezone"`
    ToTimezone   string `json:"to_timezone"`
}
```

**Timezone Format:** IANA timezone names (e.g., "America/New_York", "Asia/Tokyo", "UTC")

**Returns:**
```json
{
    "original": "2024-01-01 12:00:00 EST",
    "converted": "2024-01-02 02:00:00 JST",
    "original_offset": "-05:00",
    "converted_offset": "+09:00",
    "unix_timestamp": 1704207600,
    "is_dst_original": false,
    "is_dst_converted": false
}
```

### parse_datetime

Parses datetime from various string formats with auto-detection.

**Tool Name:** `parse_datetime`

**Parameters:**
```go
type ParseDatetimeParams struct {
    DatetimeString string `json:"datetime_string"`
    Format         string `json:"format,omitempty"`     // Optional specific format
    Timezone       string `json:"timezone,omitempty"`
}
```

**Auto-detected Formats:**
- RFC3339, RFC3339Nano
- ISO formats
- Common date formats (MM/DD/YYYY, DD-MM-YYYY)
- Month name formats ("Jan 2, 2024", "January 2, 2024")
- Unix timestamps (seconds and milliseconds)
- RFC1123, RFC822
- Kitchen time

**Returns:**
```json
{
    "parsed": "2024-01-01 15:04:05 UTC",
    "iso": "2024-01-01T15:04:05Z",
    "unix_timestamp": 1704121445,
    "detected_format": "2006-01-02T15:04:05Z07:00",
    "year": 2024,
    "month": "January",
    "day": 1,
    "hour": 15,
    "minute": 4,
    "second": 5,
    "weekday": "Monday",
    "timezone": "UTC"
}
```

### format_datetime

Formats datetime to specific patterns.

**Tool Name:** `format_datetime`

**Parameters:**
```go
type FormatDatetimeParams struct {
    Datetime     string `json:"datetime"`
    OutputFormat string `json:"output_format"`
    Timezone     string `json:"timezone,omitempty"`
}
```

**Predefined Output Formats:**
- `rfc3339` - 2006-01-02T15:04:05Z07:00
- `rfc1123` - Mon, 02 Jan 2006 15:04:05 MST
- `kitchen` - 3:04PM
- `stamp` - Jan _2 15:04:05
- `human` - Monday, January 2, 2006 at 3:04 PM MST
- `short` - Jan 2, 2006
- `long` - January 2, 2006 15:04:05 MST
- Custom Go time format strings

**Returns:**
```json
{
    "formatted": "Monday, January 1, 2024 at 2:30 PM UTC",
    "original": "2024-01-01 14:30:00",
    "timezone": "UTC"
}
```

### compare_datetimes

Compares two datetimes.

**Tool Name:** `compare_datetimes`

**Parameters:**
```go
type CompareDatetimesParams struct {
    Datetime1 string `json:"datetime1"`
    Datetime2 string `json:"datetime2"`
}
```

**Returns:**
```json
{
    "datetime1": "2024-01-01 10:00:00",
    "datetime2": "2024-01-05 15:30:00",
    "is_before": true,
    "is_after": false,
    "is_equal": false,
    "difference_seconds": 366600,
    "difference_minutes": 6110,
    "difference_hours": 101.83333,
    "difference_days": 4.243055,
    "difference_human": "101h30m0s"
}
```

### get_business_days

Calculates business days between two dates.

**Tool Name:** `get_business_days`

**Parameters:**
```go
type GetBusinessDaysParams struct {
    StartDate        string   `json:"start_date"`
    EndDate          string   `json:"end_date"`
    ExcludeHolidays  []string `json:"exclude_holidays,omitempty"`
}
```

**Returns:**
```json
{
    "business_days": 10,
    "weekend_days": 4,
    "holidays": 1,
    "total_days": 15,
    "start_date": "2024-01-01",
    "end_date": "2024-01-15"
}
```

### convert_unix_timestamp

Converts between Unix timestamps and datetime.

**Tool Name:** `convert_unix_timestamp`

**Parameters:**
```go
type ConvertUnixTimestampParams struct {
    Value     string `json:"value"`
    Direction string `json:"direction"` // to_unix or from_unix
    Unit      string `json:"unit,omitempty"` // seconds, milliseconds, nanoseconds
}
```

**Returns (to_unix):**
```json
{
    "datetime": "2024-01-01 12:00:00",
    "timestamp": 1704110400,
    "unit": "seconds"
}
```

**Returns (from_unix):**
```json
{
    "timestamp": 1704110400,
    "datetime": "2024-01-01 12:00:00",
    "iso": "2024-01-01T12:00:00Z",
    "unit": "seconds"
}
```

## Usage Examples

### With go-llms Agent

```go
import (
    "github.com/lexlapax/go-llms/pkg/agent/workflow"
    "github.com/lexlapax/go-flock/pkg/tools"
)

// Create agent with datetime tools
agent := workflow.NewAgent(provider)
agent.AddTool(tools.NewGetCurrentDatetimeTool())
agent.AddTool(tools.NewCalculateDurationTool())

// Agent can now use these tools
result, _ := agent.Run(ctx, "What's the current time in Tokyo?")
```

### Direct Tool Usage

```go
// Get current time in different timezone
tool := tools.NewGetCurrentDatetimeTool()
params := tools.GetCurrentDateTimeParams{
    Format: "datetime",
    Timezone: "Asia/Tokyo",
}
result, _ := tool.Execute(ctx, params)
```

## Error Handling

All tools return descriptive errors for:
- Invalid date formats
- Unknown timezones
- Invalid parameters
- Parsing failures

## Best Practices

1. **Date Parsing**: The tools attempt multiple format parsings automatically
2. **Timezones**: Always use IANA timezone names
3. **Unix Timestamps**: Check the unit (seconds vs milliseconds) when converting
4. **Business Days**: Provide holiday lists for accurate calculations
5. **Relative Times**: Use `get_current_datetime` with `relative` format for human-friendly output

## Performance Notes

- All tools are stateless and thread-safe
- Timezone data is cached by Go's time package
- Date parsing tries formats in order of likelihood
- No external dependencies beyond Go standard library