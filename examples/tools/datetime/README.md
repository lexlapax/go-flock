# DateTime Tools Example

This example demonstrates all datetime tools available in go-flock and how to use them with go-llms agents.

## Overview

The datetime_tools package provides comprehensive datetime operations that can be used by LLM agents. All tools implement the go-llms Tool interface and follow the `verb_object` naming convention.

## Available Tools

### 1. `get_current_datetime`
Gets the current date and time in various formats.

**Parameters:**
- `format` (string, required): The output format
- `timezone` (string, optional): Timezone name (e.g., "UTC", "America/New_York")

**Supported Formats:**
- `date` - ISO date (2024-01-01)
- `time` - 24-hour time (15:04:05)
- `datetime_iso` - ISO datetime (2024-01-01T15:04:05)
- `filename` - Filename-safe (2024-01-01_15-04-05)
- `log` - Log format ([2024-01-01 15:04:05])
- `relative` - Relative time ("2 hours ago")
- And many more...

### 2. `calculate_duration`
Calculates the time difference between two dates.

**Parameters:**
- `start` (string, required): Start date
- `end` (string, required): End date

**Returns:** Seconds, minutes, hours, days, and human-readable duration

### 3. `get_weekday`
Gets the day of the week for a given date.

**Parameters:**
- `date` (string, optional): Date to check (defaults to today)

**Returns:** Weekday name, weekday number, and weekend status

### 4. `add_subtract_time`
Adds or subtracts time units from a date.

**Parameters:**
- `date` (string, required): Base date
- `amount` (int, required): Amount to add/subtract
- `unit` (string, required): Time unit (year, month, day, hour, minute, second)
- `operation` (string, required): "add" or "subtract"

### 5. `convert_timezone`
Converts datetime between different timezones.

**Parameters:**
- `datetime` (string, required): DateTime to convert
- `from_timezone` (string, required): Source timezone
- `to_timezone` (string, required): Target timezone

### 6. `parse_datetime`
Parses datetime from various string formats with auto-detection.

**Parameters:**
- `datetime_string` (string, required): String to parse
- `format` (string, optional): Specific format to use
- `timezone` (string, optional): Timezone for parsing

**Auto-detects:** ISO formats, Unix timestamps, common date formats

### 7. `format_datetime`
Formats datetime to specific patterns.

**Parameters:**
- `datetime` (string, required): DateTime to format
- `output_format` (string, required): Target format
- `timezone` (string, optional): Timezone for output

**Predefined Formats:** `rfc3339`, `rfc1123`, `human`, `short`, `long`, etc.

### 8. `compare_datetimes`
Compares two datetimes.

**Parameters:**
- `datetime1` (string, required): First datetime
- `datetime2` (string, required): Second datetime

**Returns:** Before/after/equal status and differences in various units

### 9. `get_business_days`
Calculates business days between two dates.

**Parameters:**
- `start_date` (string, required): Start date
- `end_date` (string, required): End date
- `exclude_holidays` ([]string, optional): List of holiday dates

**Returns:** Business days, weekend days, holidays, and total days

### 10. `convert_unix_timestamp`
Converts between Unix timestamps and datetime.

**Parameters:**
- `value` (string, required): Value to convert
- `direction` (string, required): "to_unix" or "from_unix"
- `unit` (string, optional): "seconds", "milliseconds", or "nanoseconds"

## Running the Example

```bash
cd examples/datetime_tools
go run main.go
```

## Integration with go-llms Agents

These tools can be added to any go-llms agent:

```go
import (
    "github.com/lexlapax/go-llms/pkg/agent/workflow"
    "github.com/lexlapax/go-flock/pkg/tools"
)

// Create an agent
agent := workflow.NewAgent(provider)

// Add datetime tools
agent.AddTool(tools.NewGetCurrentDatetimeTool())
agent.AddTool(tools.NewCalculateDurationTool())
agent.AddTool(tools.NewConvertTimezoneTool())
// ... add more tools as needed

// Now the agent can use these tools when processing requests
result, err := agent.Run(ctx, "What day of the week is Christmas 2024?")
```

## Use Cases

1. **Scheduling**: Calculate business days for project timelines
2. **Data Processing**: Parse timestamps from various sources
3. **Internationalization**: Convert times between timezones
4. **Logging**: Format timestamps for different logging systems
5. **Analytics**: Calculate durations and compare time periods
6. **Automation**: Add/subtract time for scheduled tasks

## Notes

- All date parsing attempts multiple common formats automatically
- Timezone names follow IANA timezone database (e.g., "America/New_York")
- Unix timestamps support seconds, milliseconds, and nanoseconds
- Business days calculation excludes weekends and custom holidays
- All tools handle errors gracefully with descriptive messages