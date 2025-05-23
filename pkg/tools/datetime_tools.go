// ABOUTME: This file provides datetime-related tools for working with dates, times, and durations.
// ABOUTME: All tools implement the go-llms Tool interface for use with LLM agents.

package tools

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/lexlapax/go-llms/pkg/agent/domain"
	"github.com/lexlapax/go-llms/pkg/agent/tools"
	sdomain "github.com/lexlapax/go-llms/pkg/schema/domain"
)

// DateTimeFormat represents the available datetime formats
type DateTimeFormat string

const (
	// Date formats
	FormatDate        DateTimeFormat = "date"         // 2024-12-10
	FormatDateSlash   DateTimeFormat = "date_slash"   // 2024/12/12
	FormatDateCompact DateTimeFormat = "date_compact" // 20241210
	FormatDateUS      DateTimeFormat = "date_us"      // 12/10/2024
	FormatDateEU      DateTimeFormat = "date_eu"      // 10.12.2024

	// Time formats
	FormatTime        DateTimeFormat = "time"         // 15:04:05
	FormatTime12      DateTimeFormat = "time_12"      // 3:04:05 PM
	FormatTimeCompact DateTimeFormat = "time_compact" // 150405

	// DateTime formats
	FormatDateTime        DateTimeFormat = "datetime"         // 2024-12-10 15:04:05
	FormatDateTimeISO     DateTimeFormat = "datetime_iso"     // 2024-12-10T15:04:05
	FormatDateTimeRFC     DateTimeFormat = "datetime_rfc"     // Mon, 10 Dec 2024 15:04:05 MST
	FormatDateTimeUnix    DateTimeFormat = "datetime_unix"    // Unix timestamp
	FormatDateTimeCompact DateTimeFormat = "datetime_compact" // 20241210150405

	// Special formats
	FormatFilename    DateTimeFormat = "filename"     // 2024-12-10_15-04-05
	FormatFilenameMD  DateTimeFormat = "filename_md"  // 2024-12-10_15-04-05.md
	FormatFilenameTXT DateTimeFormat = "filename_txt" // 2024-12-10_15-04-05.txt
	FormatFilenameLOG DateTimeFormat = "filename_log" // 2024-12-10_15-04-05.log
	FormatLog         DateTimeFormat = "log"          // [2024-12-10 15:04:05]
	FormatLogCompact  DateTimeFormat = "log_compact"  // [20241210-150405]

	// Relative formats
	FormatRelative DateTimeFormat = "relative" // "2 hours ago", "in 3 days"
	FormatDuration DateTimeFormat = "duration" // Duration since/until a reference time
)

// GetCurrentDateTimeParams defines parameters for the datetime tool
type GetCurrentDateTimeParams struct {
	Format   string `json:"format"`
	Timezone string `json:"timezone,omitempty"`
}

// GetCurrentDateTimeResult defines the result of the datetime tool
type GetCurrentDateTimeResult struct {
	Formatted string `json:"formatted"`
	Format    string `json:"format"`
	Timezone  string `json:"timezone"`
	Unix      int64  `json:"unix"`
}

// GetCurrentDateTimeParamSchema defines the parameter schema
var GetCurrentDateTimeParamSchema = &sdomain.Schema{
	Type: "object",
	Properties: map[string]sdomain.Property{
		"format": {
			Type:        "string",
			Description: "The format for the datetime output",
			Enum: []string{
				"date", "date_slash", "date_compact", "date_us", "date_eu",
				"time", "time_12", "time_compact",
				"datetime", "datetime_iso", "datetime_rfc", "datetime_unix", "datetime_compact",
				"filename", "filename_md", "filename_txt", "filename_log",
				"log", "log_compact",
				"relative", "duration",
			},
		},
		"timezone": {
			Type:        "string",
			Description: "Timezone for the datetime (e.g., 'UTC', 'America/New_York'). Defaults to local timezone",
		},
	},
	Required: []string{"format"},
}

// NewGetCurrentDatetimeTool creates a tool for getting current date/time in various formats
func NewGetCurrentDatetimeTool() domain.Tool {
	return tools.NewTool(
		"get_current_datetime",
		"Gets the current date and time in various formats",
		func(ctx context.Context, params GetCurrentDateTimeParams) (*GetCurrentDateTimeResult, error) {
			// Get current time
			now := time.Now()

			// Handle timezone if specified
			timezone := "Local"
			if params.Timezone != "" {
				loc, err := time.LoadLocation(params.Timezone)
				if err != nil {
					return nil, fmt.Errorf("invalid timezone %s: %w", params.Timezone, err)
				}
				now = now.In(loc)
				timezone = params.Timezone
			}

			// Format based on requested format
			var formatted string
			format := DateTimeFormat(params.Format)

			switch format {
			// Date formats
			case FormatDate:
				formatted = now.Format("2006-01-02")
			case FormatDateSlash:
				formatted = now.Format("2006/01/02")
			case FormatDateCompact:
				formatted = now.Format("20060102")
			case FormatDateUS:
				formatted = now.Format("01/02/2006")
			case FormatDateEU:
				formatted = now.Format("02.01.2006")

			// Time formats
			case FormatTime:
				formatted = now.Format("15:04:05")
			case FormatTime12:
				formatted = now.Format("3:04:05 PM")
			case FormatTimeCompact:
				formatted = now.Format("150405")

			// DateTime formats
			case FormatDateTime:
				formatted = now.Format("2006-01-02 15:04:05")
			case FormatDateTimeISO:
				formatted = now.Format("2006-01-02T15:04:05")
			case FormatDateTimeRFC:
				formatted = now.Format(time.RFC1123)
			case FormatDateTimeUnix:
				formatted = fmt.Sprintf("%d", now.Unix())
			case FormatDateTimeCompact:
				formatted = now.Format("20060102150405")

			// Special formats
			case FormatFilename:
				formatted = now.Format("2006-01-02_15-04-05")
			case FormatFilenameMD:
				formatted = now.Format("2006-01-02_15-04-05") + ".md"
			case FormatFilenameTXT:
				formatted = now.Format("2006-01-02_15-04-05") + ".txt"
			case FormatFilenameLOG:
				formatted = now.Format("2006-01-02_15-04-05") + ".log"
			case FormatLog:
				formatted = now.Format("[2006-01-02 15:04:05]")
			case FormatLogCompact:
				formatted = now.Format("[20060102-150405]")

			// Relative formats
			case FormatRelative:
				formatted = getRelativeTime(now)
			case FormatDuration:
				// For duration, we'll show time since Unix epoch as example
				duration := time.Since(time.Unix(0, 0))
				formatted = duration.String()

			default:
				return nil, fmt.Errorf("unknown format: %s", params.Format)
			}

			return &GetCurrentDateTimeResult{
				Formatted: formatted,
				Format:    params.Format,
				Timezone:  timezone,
				Unix:      now.Unix(),
			}, nil
		},
		GetCurrentDateTimeParamSchema,
	)
}

// getRelativeTime returns a human-readable relative time string
func getRelativeTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < 0 {
		// Future time
		diff = -diff
		switch {
		case diff < time.Minute:
			return fmt.Sprintf("in %d seconds", int(diff.Seconds()))
		case diff < time.Hour:
			return fmt.Sprintf("in %d minutes", int(diff.Minutes()))
		case diff < 24*time.Hour:
			return fmt.Sprintf("in %d hours", int(diff.Hours()))
		default:
			return fmt.Sprintf("in %d days", int(diff.Hours()/24))
		}
	} else {
		// Past time
		switch {
		case diff < time.Minute:
			return fmt.Sprintf("%d seconds ago", int(diff.Seconds()))
		case diff < time.Hour:
			return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
		case diff < 24*time.Hour:
			return fmt.Sprintf("%d hours ago", int(diff.Hours()))
		default:
			return fmt.Sprintf("%d days ago", int(diff.Hours()/24))
		}
	}
}

// Additional helper functions for common datetime operations

// CalculateDurationParams defines parameters for calculating time between dates
type CalculateDurationParams struct {
	Start  string `json:"start"`
	End    string `json:"end"`
	Format string `json:"format,omitempty"`
}

// NewCalculateDurationTool creates a tool for calculating time between two dates
func NewCalculateDurationTool() domain.Tool {
	return tools.NewTool(
		"calculate_duration",
		"Calculates the time difference between two dates",
		func(ctx context.Context, params CalculateDurationParams) (map[string]interface{}, error) {
			// Parse dates (try multiple formats)
			formats := []string{
				"2006-01-02",
				"2006-01-02 15:04:05",
				"2006-01-02T15:04:05",
				time.RFC3339,
			}

			var start, end time.Time
			var err error

			for _, format := range formats {
				start, err = time.Parse(format, params.Start)
				if err == nil {
					break
				}
			}
			if err != nil {
				return nil, fmt.Errorf("could not parse start date: %w", err)
			}

			for _, format := range formats {
				end, err = time.Parse(format, params.End)
				if err == nil {
					break
				}
			}
			if err != nil {
				return nil, fmt.Errorf("could not parse end date: %w", err)
			}

			diff := end.Sub(start)

			return map[string]interface{}{
				"seconds": diff.Seconds(),
				"minutes": diff.Minutes(),
				"hours":   diff.Hours(),
				"days":    diff.Hours() / 24,
				"human":   diff.String(),
			}, nil
		},
		nil,
	)
}

// GetWeekdayParams defines parameters for getting weekday
type GetWeekdayParams struct {
	Date string `json:"date,omitempty"`
}

// NewGetWeekdayTool creates a tool for getting the day of the week
func NewGetWeekdayTool() domain.Tool {
	return tools.NewTool(
		"get_weekday",
		"Gets the day of the week for a given date (or today if not specified)",
		func(ctx context.Context, params GetWeekdayParams) (map[string]interface{}, error) {
			var t time.Time
			var err error

			if params.Date == "" {
				t = time.Now()
			} else {
				// Try to parse the date
				formats := []string{
					"2006-01-02",
					"2006-01-02 15:04:05",
					"2006-01-02T15:04:05",
				}

				for _, format := range formats {
					t, err = time.Parse(format, params.Date)
					if err == nil {
						break
					}
				}
				if err != nil {
					return nil, fmt.Errorf("could not parse date: %w", err)
				}
			}

			return map[string]interface{}{
				"weekday":     t.Weekday().String(),
				"weekday_num": int(t.Weekday()), // 0 = Sunday
				"is_weekend":  t.Weekday() == time.Saturday || t.Weekday() == time.Sunday,
				"date":        t.Format("2006-01-02"),
			}, nil
		},
		nil,
	)
}

// AddSubtractTimeParams defines parameters for adding/subtracting time
type AddSubtractTimeParams struct {
	Date      string `json:"date"`
	Amount    int    `json:"amount"`
	Unit      string `json:"unit"`      // day, hour, minute, month, year
	Operation string `json:"operation"` // add, subtract
}

// NewAddSubtractTimeTool creates a tool for adding or subtracting time units
func NewAddSubtractTimeTool() domain.Tool {
	return tools.NewTool(
		"add_subtract_time",
		"Adds or subtracts time units (days, hours, months, years) from a date",
		func(ctx context.Context, params AddSubtractTimeParams) (map[string]interface{}, error) {
			// Parse the date
			var t time.Time
			var err error
			formats := []string{
				"2006-01-02",
				"2006-01-02 15:04:05",
				"2006-01-02T15:04:05",
				time.RFC3339,
			}

			for _, format := range formats {
				t, err = time.Parse(format, params.Date)
				if err == nil {
					break
				}
			}
			if err != nil {
				return nil, fmt.Errorf("could not parse date: %w", err)
			}

			// Apply the operation
			amount := params.Amount
			if params.Operation == "subtract" {
				amount = -amount
			} else if params.Operation != "add" {
				return nil, fmt.Errorf("operation must be 'add' or 'subtract'")
			}

			var result time.Time
			switch params.Unit {
			case "year":
				result = t.AddDate(amount, 0, 0)
			case "month":
				result = t.AddDate(0, amount, 0)
			case "day":
				result = t.AddDate(0, 0, amount)
			case "hour":
				result = t.Add(time.Duration(amount) * time.Hour)
			case "minute":
				result = t.Add(time.Duration(amount) * time.Minute)
			case "second":
				result = t.Add(time.Duration(amount) * time.Second)
			default:
				return nil, fmt.Errorf("invalid unit: %s (must be year, month, day, hour, minute, or second)", params.Unit)
			}

			return map[string]interface{}{
				"original":       t.Format("2006-01-02 15:04:05"),
				"result":         result.Format("2006-01-02 15:04:05"),
				"result_iso":     result.Format(time.RFC3339),
				"operation":      fmt.Sprintf("%s %d %s(s)", params.Operation, params.Amount, params.Unit),
				"unix_timestamp": result.Unix(),
			}, nil
		},
		nil,
	)
}

// ConvertTimezoneParams defines parameters for timezone conversion
type ConvertTimezoneParams struct {
	Datetime     string `json:"datetime"`
	FromTimezone string `json:"from_timezone"`
	ToTimezone   string `json:"to_timezone"`
}

// NewConvertTimezoneTool creates a tool for converting between timezones
func NewConvertTimezoneTool() domain.Tool {
	return tools.NewTool(
		"convert_timezone",
		"Converts datetime between different timezones",
		func(ctx context.Context, params ConvertTimezoneParams) (map[string]interface{}, error) {
			// Load source timezone
			fromLoc, err := time.LoadLocation(params.FromTimezone)
			if err != nil {
				return nil, fmt.Errorf("invalid source timezone %s: %w", params.FromTimezone, err)
			}

			// Load target timezone
			toLoc, err := time.LoadLocation(params.ToTimezone)
			if err != nil {
				return nil, fmt.Errorf("invalid target timezone %s: %w", params.ToTimezone, err)
			}

			// Parse the datetime in the source timezone
			formats := []string{
				"2006-01-02 15:04:05",
				"2006-01-02T15:04:05",
				time.RFC3339,
			}

			var t time.Time
			for _, format := range formats {
				t, err = time.ParseInLocation(format, params.Datetime, fromLoc)
				if err == nil {
					break
				}
			}
			if err != nil {
				return nil, fmt.Errorf("could not parse datetime: %w", err)
			}

			// Convert to target timezone
			converted := t.In(toLoc)

			return map[string]interface{}{
				"original":         t.Format("2006-01-02 15:04:05 MST"),
				"converted":        converted.Format("2006-01-02 15:04:05 MST"),
				"original_offset":  t.Format("-07:00"),
				"converted_offset": converted.Format("-07:00"),
				"unix_timestamp":   t.Unix(),
				"is_dst_original":  t.IsDST(),
				"is_dst_converted": converted.IsDST(),
			}, nil
		},
		nil,
	)
}

// ParseDatetimeParams defines parameters for parsing datetime strings
type ParseDatetimeParams struct {
	DatetimeString string `json:"datetime_string"`
	Format         string `json:"format,omitempty"`
	Timezone       string `json:"timezone,omitempty"`
}

// NewParseDatetimeTool creates a tool for parsing datetime from various formats
func NewParseDatetimeTool() domain.Tool {
	return tools.NewTool(
		"parse_datetime",
		"Parses datetime from various string formats",
		func(ctx context.Context, params ParseDatetimeParams) (map[string]interface{}, error) {
			var t time.Time
			var err error
			var detectedFormat string

			// Get location if specified
			loc := time.Local
			if params.Timezone != "" {
				loc, err = time.LoadLocation(params.Timezone)
				if err != nil {
					return nil, fmt.Errorf("invalid timezone %s: %w", params.Timezone, err)
				}
			}

			// If format is specified, use it
			if params.Format != "" {
				t, err = time.ParseInLocation(params.Format, params.DatetimeString, loc)
				if err != nil {
					return nil, fmt.Errorf("could not parse with format %s: %w", params.Format, err)
				}
				detectedFormat = params.Format
			} else {
				// Try to auto-detect format
				formats := []string{
					time.RFC3339,
					time.RFC3339Nano,
					"2006-01-02 15:04:05",
					"2006-01-02T15:04:05",
					"2006-01-02",
					"01/02/2006",
					"02-01-2006",
					"Jan 2, 2006",
					"January 2, 2006",
					"2 Jan 2006",
					"Mon, 02 Jan 2006 15:04:05 MST",
					time.RFC1123,
					time.RFC822,
					time.Kitchen,
				}

				// Check if it's a Unix timestamp
				if timestamp, err := strconv.ParseInt(params.DatetimeString, 10, 64); err == nil {
					// Likely a Unix timestamp
					if timestamp > 1e10 { // Milliseconds
						t = time.Unix(0, timestamp*1e6).In(loc)
						detectedFormat = "unix_milliseconds"
					} else { // Seconds
						t = time.Unix(timestamp, 0).In(loc)
						detectedFormat = "unix_seconds"
					}
				} else {
					// Try each format
					for _, format := range formats {
						t, err = time.ParseInLocation(format, params.DatetimeString, loc)
						if err == nil {
							detectedFormat = format
							break
						}
					}
					if err != nil {
						return nil, fmt.Errorf("could not parse datetime string: %s", params.DatetimeString)
					}
				}
			}

			return map[string]interface{}{
				"parsed":          t.Format("2006-01-02 15:04:05 MST"),
				"iso":             t.Format(time.RFC3339),
				"unix_timestamp":  t.Unix(),
				"detected_format": detectedFormat,
				"year":            t.Year(),
				"month":           t.Month().String(),
				"day":             t.Day(),
				"hour":            t.Hour(),
				"minute":          t.Minute(),
				"second":          t.Second(),
				"weekday":         t.Weekday().String(),
				"timezone":        t.Location().String(),
			}, nil
		},
		nil,
	)
}

// FormatDatetimeParams defines parameters for formatting datetime
type FormatDatetimeParams struct {
	Datetime     string `json:"datetime"`
	OutputFormat string `json:"output_format"`
	Timezone     string `json:"timezone,omitempty"`
}

// NewFormatDatetimeTool creates a tool for formatting datetime to specific patterns
func NewFormatDatetimeTool() domain.Tool {
	return tools.NewTool(
		"format_datetime",
		"Formats datetime to specific patterns",
		func(ctx context.Context, params FormatDatetimeParams) (map[string]interface{}, error) {
			// Parse the datetime
			var t time.Time
			var err error
			formats := []string{
				"2006-01-02",
				"2006-01-02 15:04:05",
				"2006-01-02T15:04:05",
				time.RFC3339,
			}

			for _, format := range formats {
				t, err = time.Parse(format, params.Datetime)
				if err == nil {
					break
				}
			}
			if err != nil {
				return nil, fmt.Errorf("could not parse datetime: %w", err)
			}

			// Apply timezone if specified
			if params.Timezone != "" {
				loc, err := time.LoadLocation(params.Timezone)
				if err != nil {
					return nil, fmt.Errorf("invalid timezone %s: %w", params.Timezone, err)
				}
				t = t.In(loc)
			}

			// Format based on output format
			var formatted string
			switch params.OutputFormat {
			case "rfc3339":
				formatted = t.Format(time.RFC3339)
			case "rfc1123":
				formatted = t.Format(time.RFC1123)
			case "kitchen":
				formatted = t.Format(time.Kitchen)
			case "stamp":
				formatted = t.Format(time.Stamp)
			case "human":
				formatted = t.Format("Monday, January 2, 2006 at 3:04 PM MST")
			case "short":
				formatted = t.Format("Jan 2, 2006")
			case "long":
				formatted = t.Format("January 2, 2006 15:04:05 MST")
			default:
				// Use custom format
				formatted = t.Format(params.OutputFormat)
			}

			return map[string]interface{}{
				"formatted": formatted,
				"original":  params.Datetime,
				"timezone":  t.Location().String(),
			}, nil
		},
		nil,
	)
}

// CompareDatetimesParams defines parameters for comparing datetimes
type CompareDatetimesParams struct {
	Datetime1 string `json:"datetime1"`
	Datetime2 string `json:"datetime2"`
}

// NewCompareDatetimesTool creates a tool for comparing two datetimes
func NewCompareDatetimesTool() domain.Tool {
	return tools.NewTool(
		"compare_datetimes",
		"Compares two datetimes",
		func(ctx context.Context, params CompareDatetimesParams) (map[string]interface{}, error) {
			// Parse both datetimes
			formats := []string{
				"2006-01-02",
				"2006-01-02 15:04:05",
				"2006-01-02T15:04:05",
				time.RFC3339,
			}

			var t1, t2 time.Time
			var err error

			for _, format := range formats {
				t1, err = time.Parse(format, params.Datetime1)
				if err == nil {
					break
				}
			}
			if err != nil {
				return nil, fmt.Errorf("could not parse datetime1: %w", err)
			}

			for _, format := range formats {
				t2, err = time.Parse(format, params.Datetime2)
				if err == nil {
					break
				}
			}
			if err != nil {
				return nil, fmt.Errorf("could not parse datetime2: %w", err)
			}

			diff := t2.Sub(t1)

			return map[string]interface{}{
				"datetime1":          t1.Format("2006-01-02 15:04:05"),
				"datetime2":          t2.Format("2006-01-02 15:04:05"),
				"is_before":          t1.Before(t2),
				"is_after":           t1.After(t2),
				"is_equal":           t1.Equal(t2),
				"difference_seconds": diff.Seconds(),
				"difference_minutes": diff.Minutes(),
				"difference_hours":   diff.Hours(),
				"difference_days":    diff.Hours() / 24,
				"difference_human":   diff.String(),
			}, nil
		},
		nil,
	)
}

// GetBusinessDaysParams defines parameters for calculating business days
type GetBusinessDaysParams struct {
	StartDate       string   `json:"start_date"`
	EndDate         string   `json:"end_date"`
	ExcludeHolidays []string `json:"exclude_holidays,omitempty"`
}

// NewGetBusinessDaysTool creates a tool for calculating business days between dates
func NewGetBusinessDaysTool() domain.Tool {
	return tools.NewTool(
		"get_business_days",
		"Calculates business days between two dates",
		func(ctx context.Context, params GetBusinessDaysParams) (map[string]interface{}, error) {
			// Parse dates
			startDate, err := time.Parse("2006-01-02", params.StartDate)
			if err != nil {
				return nil, fmt.Errorf("could not parse start date: %w", err)
			}

			endDate, err := time.Parse("2006-01-02", params.EndDate)
			if err != nil {
				return nil, fmt.Errorf("could not parse end date: %w", err)
			}

			// Parse holidays
			holidays := make(map[string]bool)
			for _, holiday := range params.ExcludeHolidays {
				h, err := time.Parse("2006-01-02", holiday)
				if err == nil {
					holidays[h.Format("2006-01-02")] = true
				}
			}

			// Calculate business days
			businessDays := 0
			weekendDays := 0
			holidayCount := 0
			totalDays := 0

			current := startDate
			for !current.After(endDate) {
				totalDays++
				dateStr := current.Format("2006-01-02")

				if current.Weekday() == time.Saturday || current.Weekday() == time.Sunday {
					weekendDays++
				} else if holidays[dateStr] {
					holidayCount++
				} else {
					businessDays++
				}

				current = current.AddDate(0, 0, 1)
			}

			return map[string]interface{}{
				"business_days": businessDays,
				"weekend_days":  weekendDays,
				"holidays":      holidayCount,
				"total_days":    totalDays,
				"start_date":    startDate.Format("2006-01-02"),
				"end_date":      endDate.Format("2006-01-02"),
			}, nil
		},
		nil,
	)
}

// ConvertUnixTimestampParams defines parameters for unix timestamp conversion
type ConvertUnixTimestampParams struct {
	Value     string `json:"value"`
	Direction string `json:"direction"`      // to_unix or from_unix
	Unit      string `json:"unit,omitempty"` // seconds, milliseconds, nanoseconds
}

// NewConvertUnixTimestampTool creates a tool for converting unix timestamps
func NewConvertUnixTimestampTool() domain.Tool {
	return tools.NewTool(
		"convert_unix_timestamp",
		"Converts between Unix timestamps and datetime",
		func(ctx context.Context, params ConvertUnixTimestampParams) (map[string]interface{}, error) {
			unit := params.Unit
			if unit == "" {
				unit = "seconds"
			}

			if params.Direction == "to_unix" {
				// Parse datetime and convert to unix
				formats := []string{
					"2006-01-02",
					"2006-01-02 15:04:05",
					"2006-01-02T15:04:05",
					time.RFC3339,
				}

				var t time.Time
				var err error
				for _, format := range formats {
					t, err = time.Parse(format, params.Value)
					if err == nil {
						break
					}
				}
				if err != nil {
					return nil, fmt.Errorf("could not parse datetime: %w", err)
				}

				var timestamp int64
				switch unit {
				case "seconds":
					timestamp = t.Unix()
				case "milliseconds":
					timestamp = t.UnixMilli()
				case "nanoseconds":
					timestamp = t.UnixNano()
				default:
					return nil, fmt.Errorf("invalid unit: %s (must be seconds, milliseconds, or nanoseconds)", unit)
				}

				return map[string]interface{}{
					"datetime":  params.Value,
					"timestamp": timestamp,
					"unit":      unit,
				}, nil

			} else if params.Direction == "from_unix" {
				// Convert unix to datetime
				timestamp, err := strconv.ParseInt(params.Value, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid unix timestamp: %w", err)
				}

				var t time.Time
				switch unit {
				case "seconds":
					t = time.Unix(timestamp, 0)
				case "milliseconds":
					t = time.UnixMilli(timestamp)
				case "nanoseconds":
					t = time.Unix(0, timestamp)
				default:
					return nil, fmt.Errorf("invalid unit: %s (must be seconds, milliseconds, or nanoseconds)", unit)
				}

				return map[string]interface{}{
					"timestamp": timestamp,
					"datetime":  t.Format("2006-01-02 15:04:05"),
					"iso":       t.Format(time.RFC3339),
					"unit":      unit,
				}, nil

			} else {
				return nil, fmt.Errorf("direction must be 'to_unix' or 'from_unix'")
			}
		},
		nil,
	)
}
