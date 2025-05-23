package tools

import (
	"context"
	"strings"
	"testing"
)

func TestGetCurrentDatetimeTool(t *testing.T) {
	tool := NewGetCurrentDatetimeTool()
	ctx := context.Background()

	tests := []struct {
		name      string
		params    GetCurrentDateTimeParams
		wantErr   bool
		checkFunc func(t *testing.T, result *GetCurrentDateTimeResult)
	}{
		{
			name:   "date format",
			params: GetCurrentDateTimeParams{Format: "date"},
			checkFunc: func(t *testing.T, result *GetCurrentDateTimeResult) {
				// Check format YYYY-MM-DD
				if len(result.Formatted) != 10 {
					t.Errorf("expected date format YYYY-MM-DD, got %s", result.Formatted)
				}
				if !strings.Contains(result.Formatted, "-") {
					t.Errorf("expected date to contain dashes, got %s", result.Formatted)
				}
			},
		},
		{
			name:   "time format",
			params: GetCurrentDateTimeParams{Format: "time"},
			checkFunc: func(t *testing.T, result *GetCurrentDateTimeResult) {
				// Check format HH:MM:SS
				if len(result.Formatted) != 8 {
					t.Errorf("expected time format HH:MM:SS, got %s", result.Formatted)
				}
				if !strings.Contains(result.Formatted, ":") {
					t.Errorf("expected time to contain colons, got %s", result.Formatted)
				}
			},
		},
		{
			name:   "datetime ISO format",
			params: GetCurrentDateTimeParams{Format: "datetime_iso"},
			checkFunc: func(t *testing.T, result *GetCurrentDateTimeResult) {
				if !strings.Contains(result.Formatted, "T") {
					t.Errorf("expected ISO format to contain T separator, got %s", result.Formatted)
				}
			},
		},
		{
			name:   "filename format",
			params: GetCurrentDateTimeParams{Format: "filename"},
			checkFunc: func(t *testing.T, result *GetCurrentDateTimeResult) {
				if strings.Contains(result.Formatted, ":") {
					t.Errorf("filename should not contain colons, got %s", result.Formatted)
				}
				if !strings.Contains(result.Formatted, "_") {
					t.Errorf("filename should contain underscores, got %s", result.Formatted)
				}
			},
		},
		{
			name:   "filename with .md extension",
			params: GetCurrentDateTimeParams{Format: "filename_md"},
			checkFunc: func(t *testing.T, result *GetCurrentDateTimeResult) {
				if !strings.HasSuffix(result.Formatted, ".md") {
					t.Errorf("expected .md extension, got %s", result.Formatted)
				}
			},
		},
		{
			name:   "UTC timezone",
			params: GetCurrentDateTimeParams{Format: "datetime", Timezone: "UTC"},
			checkFunc: func(t *testing.T, result *GetCurrentDateTimeResult) {
				if result.Timezone != "UTC" {
					t.Errorf("expected UTC timezone, got %s", result.Timezone)
				}
			},
		},
		{
			name:   "log format",
			params: GetCurrentDateTimeParams{Format: "log"},
			checkFunc: func(t *testing.T, result *GetCurrentDateTimeResult) {
				if !strings.HasPrefix(result.Formatted, "[") || !strings.HasSuffix(result.Formatted, "]") {
					t.Errorf("log format should be wrapped in brackets, got %s", result.Formatted)
				}
			},
		},
		{
			name:   "unix timestamp",
			params: GetCurrentDateTimeParams{Format: "datetime_unix"},
			checkFunc: func(t *testing.T, result *GetCurrentDateTimeResult) {
				// Should be a number string (all digits)
				for _, char := range result.Formatted {
					if char < '0' || char > '9' {
						t.Errorf("unix timestamp should only contain digits, got %s", result.Formatted)
						break
					}
				}
			},
		},
		{
			name:    "invalid format",
			params:  GetCurrentDateTimeParams{Format: "invalid_format"},
			wantErr: true,
		},
		{
			name:    "invalid timezone",
			params:  GetCurrentDateTimeParams{Format: "date", Timezone: "Invalid/Timezone"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.Execute(ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.checkFunc != nil {
				dateTimeResult, ok := result.(*GetCurrentDateTimeResult)
				if !ok {
					t.Errorf("Execute() returned wrong type: %T", result)
					return
				}
				tt.checkFunc(t, dateTimeResult)

				// Common checks for all successful results
				if dateTimeResult.Format != tt.params.Format {
					t.Errorf("Format mismatch: expected %s, got %s", tt.params.Format, dateTimeResult.Format)
				}
				if dateTimeResult.Unix <= 0 {
					t.Errorf("Unix timestamp should be positive, got %d", dateTimeResult.Unix)
				}
			}
		})
	}
}

func TestCalculateDurationTool(t *testing.T) {
	tool := NewCalculateDurationTool()
	ctx := context.Background()

	tests := []struct {
		name    string
		params  CalculateDurationParams
		wantErr bool
		check   func(t *testing.T, result map[string]interface{})
	}{
		{
			name: "calculate days between dates",
			params: CalculateDurationParams{
				Start: "2024-01-01",
				End:   "2024-01-10",
			},
			check: func(t *testing.T, result map[string]interface{}) {
				days, ok := result["days"].(float64)
				if !ok {
					t.Errorf("days not found in result")
					return
				}
				if days != 9 {
					t.Errorf("expected 9 days, got %f", days)
				}
			},
		},
		{
			name: "calculate hours between datetimes",
			params: CalculateDurationParams{
				Start: "2024-01-01 10:00:00",
				End:   "2024-01-01 15:00:00",
			},
			check: func(t *testing.T, result map[string]interface{}) {
				hours, ok := result["hours"].(float64)
				if !ok {
					t.Errorf("hours not found in result")
					return
				}
				if hours != 5 {
					t.Errorf("expected 5 hours, got %f", hours)
				}
			},
		},
		{
			name: "invalid start date",
			params: CalculateDurationParams{
				Start: "invalid-date",
				End:   "2024-01-01",
			},
			wantErr: true,
		},
		{
			name: "invalid end date",
			params: CalculateDurationParams{
				Start: "2024-01-01",
				End:   "invalid-date",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.Execute(ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.check != nil {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Errorf("Execute() returned wrong type: %T", result)
					return
				}
				tt.check(t, resultMap)
			}
		})
	}
}

func TestGetWeekdayTool(t *testing.T) {
	tool := NewGetWeekdayTool()
	ctx := context.Background()

	tests := []struct {
		name    string
		params  GetWeekdayParams
		wantErr bool
		check   func(t *testing.T, result map[string]interface{})
	}{
		{
			name:   "get weekday for specific date",
			params: GetWeekdayParams{Date: "2024-01-01"}, // Monday
			check: func(t *testing.T, result map[string]interface{}) {
				weekday, ok := result["weekday"].(string)
				if !ok {
					t.Errorf("weekday not found in result")
					return
				}
				if weekday != "Monday" {
					t.Errorf("expected Monday, got %s", weekday)
				}

				isWeekend, ok := result["is_weekend"].(bool)
				if !ok {
					t.Errorf("is_weekend not found in result")
					return
				}
				if isWeekend {
					t.Errorf("Monday should not be weekend")
				}
			},
		},
		{
			name:   "get weekday for weekend",
			params: GetWeekdayParams{Date: "2024-01-06"}, // Saturday
			check: func(t *testing.T, result map[string]interface{}) {
				isWeekend, ok := result["is_weekend"].(bool)
				if !ok {
					t.Errorf("is_weekend not found in result")
					return
				}
				if !isWeekend {
					t.Errorf("Saturday should be weekend")
				}
			},
		},
		{
			name:   "get current weekday (no date specified)",
			params: GetWeekdayParams{},
			check: func(t *testing.T, result map[string]interface{}) {
				// Just check that we got a valid weekday
				weekday, ok := result["weekday"].(string)
				if !ok || weekday == "" {
					t.Errorf("weekday not found or empty in result")
				}
			},
		},
		{
			name:    "invalid date format",
			params:  GetWeekdayParams{Date: "not-a-date"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.Execute(ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.check != nil {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Errorf("Execute() returned wrong type: %T", result)
					return
				}
				tt.check(t, resultMap)
			}
		})
	}
}

func TestGetCurrentDatetimeTool_Metadata(t *testing.T) {
	tool := NewGetCurrentDatetimeTool()

	// Test tool metadata
	if tool.Name() != "get_current_datetime" {
		t.Errorf("expected tool name 'get_current_datetime', got %s", tool.Name())
	}

	if tool.Description() == "" {
		t.Error("tool description should not be empty")
	}

	if tool.ParameterSchema() == nil {
		t.Error("tool parameter schema should not be nil")
	}
}

func TestAddSubtractTimeTool(t *testing.T) {
	tool := NewAddSubtractTimeTool()
	ctx := context.Background()

	tests := []struct {
		name    string
		params  AddSubtractTimeParams
		wantErr bool
		check   func(t *testing.T, result map[string]interface{})
	}{
		{
			name: "add 5 days",
			params: AddSubtractTimeParams{
				Date:      "2024-01-01",
				Amount:    5,
				Unit:      "day",
				Operation: "add",
			},
			check: func(t *testing.T, result map[string]interface{}) {
				expected := "2024-01-06 00:00:00"
				if result["result"] != expected {
					t.Errorf("expected %s, got %s", expected, result["result"])
				}
			},
		},
		{
			name: "subtract 2 months",
			params: AddSubtractTimeParams{
				Date:      "2024-03-15",
				Amount:    2,
				Unit:      "month",
				Operation: "subtract",
			},
			check: func(t *testing.T, result map[string]interface{}) {
				expected := "2024-01-15 00:00:00"
				if result["result"] != expected {
					t.Errorf("expected %s, got %s", expected, result["result"])
				}
			},
		},
		{
			name: "add 3 hours",
			params: AddSubtractTimeParams{
				Date:      "2024-01-01 10:00:00",
				Amount:    3,
				Unit:      "hour",
				Operation: "add",
			},
			check: func(t *testing.T, result map[string]interface{}) {
				expected := "2024-01-01 13:00:00"
				if result["result"] != expected {
					t.Errorf("expected %s, got %s", expected, result["result"])
				}
			},
		},
		{
			name: "invalid unit",
			params: AddSubtractTimeParams{
				Date:      "2024-01-01",
				Amount:    1,
				Unit:      "week",
				Operation: "add",
			},
			wantErr: true,
		},
		{
			name: "invalid operation",
			params: AddSubtractTimeParams{
				Date:      "2024-01-01",
				Amount:    1,
				Unit:      "day",
				Operation: "multiply",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.Execute(ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.check != nil {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Errorf("Execute() returned wrong type: %T", result)
					return
				}
				tt.check(t, resultMap)
			}
		})
	}
}

func TestConvertTimezoneTool(t *testing.T) {
	tool := NewConvertTimezoneTool()
	ctx := context.Background()

	tests := []struct {
		name    string
		params  ConvertTimezoneParams
		wantErr bool
		check   func(t *testing.T, result map[string]interface{})
	}{
		{
			name: "convert NYC to Tokyo",
			params: ConvertTimezoneParams{
				Datetime:     "2024-01-01 10:00:00",
				FromTimezone: "America/New_York",
				ToTimezone:   "Asia/Tokyo",
			},
			check: func(t *testing.T, result map[string]interface{}) {
				// NYC to Tokyo is typically +14 hours
				if !strings.Contains(result["converted"].(string), "00:00:00") {
					t.Errorf("timezone conversion seems incorrect: %s", result["converted"])
				}
			},
		},
		{
			name: "convert UTC to PST",
			params: ConvertTimezoneParams{
				Datetime:     "2024-01-01 10:00:00",
				FromTimezone: "UTC",
				ToTimezone:   "America/Los_Angeles",
			},
			check: func(t *testing.T, result map[string]interface{}) {
				// Should have timezone offset
				if result["original_offset"] == result["converted_offset"] {
					t.Error("offsets should be different for different timezones")
				}
			},
		},
		{
			name: "invalid source timezone",
			params: ConvertTimezoneParams{
				Datetime:     "2024-01-01 10:00:00",
				FromTimezone: "Invalid/Zone",
				ToTimezone:   "UTC",
			},
			wantErr: true,
		},
		{
			name: "invalid target timezone",
			params: ConvertTimezoneParams{
				Datetime:     "2024-01-01 10:00:00",
				FromTimezone: "UTC",
				ToTimezone:   "Invalid/Zone",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.Execute(ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.check != nil {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Errorf("Execute() returned wrong type: %T", result)
					return
				}
				tt.check(t, resultMap)
			}
		})
	}
}

func TestParseDatetimeTool(t *testing.T) {
	tool := NewParseDatetimeTool()
	ctx := context.Background()

	tests := []struct {
		name    string
		params  ParseDatetimeParams
		wantErr bool
		check   func(t *testing.T, result map[string]interface{})
	}{
		{
			name: "parse ISO format",
			params: ParseDatetimeParams{
				DatetimeString: "2024-01-01T10:00:00Z",
			},
			check: func(t *testing.T, result map[string]interface{}) {
				year, _ := result["year"].(int)
				if year != 2024 {
					t.Errorf("expected year 2024, got %v", result["year"])
				}
				if result["detected_format"] != "2006-01-02T15:04:05Z07:00" {
					t.Errorf("expected RFC3339 format detection, got %s", result["detected_format"])
				}
			},
		},
		{
			name: "parse unix timestamp",
			params: ParseDatetimeParams{
				DatetimeString: "1704110400", // 2024-01-01 12:00:00 UTC
			},
			check: func(t *testing.T, result map[string]interface{}) {
				if result["detected_format"] != "unix_seconds" {
					t.Errorf("expected unix_seconds format, got %s", result["detected_format"])
				}
			},
		},
		{
			name: "parse with custom format",
			params: ParseDatetimeParams{
				DatetimeString: "01/15/2024",
				Format:         "01/02/2006",
			},
			check: func(t *testing.T, result map[string]interface{}) {
				if result["month"] != "January" {
					t.Errorf("expected January, got %s", result["month"])
				}
				day, _ := result["day"].(int)
				if day != 15 {
					t.Errorf("expected day 15, got %v", result["day"])
				}
			},
		},
		{
			name: "invalid datetime string",
			params: ParseDatetimeParams{
				DatetimeString: "not a date",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.Execute(ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.check != nil {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Errorf("Execute() returned wrong type: %T", result)
					return
				}
				tt.check(t, resultMap)
			}
		})
	}
}

func TestFormatDatetimeTool(t *testing.T) {
	tool := NewFormatDatetimeTool()
	ctx := context.Background()

	tests := []struct {
		name    string
		params  FormatDatetimeParams
		wantErr bool
		check   func(t *testing.T, result map[string]interface{})
	}{
		{
			name: "format to RFC3339",
			params: FormatDatetimeParams{
				Datetime:     "2024-01-01 10:00:00",
				OutputFormat: "rfc3339",
			},
			check: func(t *testing.T, result map[string]interface{}) {
				formatted := result["formatted"].(string)
				if !strings.Contains(formatted, "T") || !strings.Contains(formatted, "Z") {
					t.Errorf("expected RFC3339 format, got %s", formatted)
				}
			},
		},
		{
			name: "format to human readable",
			params: FormatDatetimeParams{
				Datetime:     "2024-01-01",
				OutputFormat: "human",
			},
			check: func(t *testing.T, result map[string]interface{}) {
				formatted := result["formatted"].(string)
				if !strings.Contains(formatted, "Monday") || !strings.Contains(formatted, "January") {
					t.Errorf("expected human readable format, got %s", formatted)
				}
			},
		},
		{
			name: "custom format",
			params: FormatDatetimeParams{
				Datetime:     "2024-01-01",
				OutputFormat: "Jan 2, 2006",
			},
			check: func(t *testing.T, result map[string]interface{}) {
				if result["formatted"] != "Jan 1, 2024" {
					t.Errorf("expected 'Jan 1, 2024', got %s", result["formatted"])
				}
			},
		},
		{
			name: "invalid datetime",
			params: FormatDatetimeParams{
				Datetime:     "invalid",
				OutputFormat: "rfc3339",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.Execute(ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.check != nil {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Errorf("Execute() returned wrong type: %T", result)
					return
				}
				tt.check(t, resultMap)
			}
		})
	}
}

func TestCompareDatetimesTool(t *testing.T) {
	tool := NewCompareDatetimesTool()
	ctx := context.Background()

	tests := []struct {
		name    string
		params  CompareDatetimesParams
		wantErr bool
		check   func(t *testing.T, result map[string]interface{})
	}{
		{
			name: "compare dates - first before second",
			params: CompareDatetimesParams{
				Datetime1: "2024-01-01",
				Datetime2: "2024-01-10",
			},
			check: func(t *testing.T, result map[string]interface{}) {
				if result["is_before"] != true {
					t.Error("expected datetime1 to be before datetime2")
				}
				if result["is_after"] != false {
					t.Error("expected datetime1 not to be after datetime2")
				}
				if result["difference_days"].(float64) != 9 {
					t.Errorf("expected 9 days difference, got %v", result["difference_days"])
				}
			},
		},
		{
			name: "compare equal dates",
			params: CompareDatetimesParams{
				Datetime1: "2024-01-01 10:00:00",
				Datetime2: "2024-01-01 10:00:00",
			},
			check: func(t *testing.T, result map[string]interface{}) {
				if result["is_equal"] != true {
					t.Error("expected datetimes to be equal")
				}
				if result["difference_seconds"].(float64) != 0 {
					t.Errorf("expected 0 seconds difference, got %v", result["difference_seconds"])
				}
			},
		},
		{
			name: "invalid first datetime",
			params: CompareDatetimesParams{
				Datetime1: "invalid",
				Datetime2: "2024-01-01",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.Execute(ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.check != nil {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Errorf("Execute() returned wrong type: %T", result)
					return
				}
				tt.check(t, resultMap)
			}
		})
	}
}

func TestGetBusinessDaysTool(t *testing.T) {
	tool := NewGetBusinessDaysTool()
	ctx := context.Background()

	tests := []struct {
		name    string
		params  GetBusinessDaysParams
		wantErr bool
		check   func(t *testing.T, result map[string]interface{})
	}{
		{
			name: "calculate business days in a week",
			params: GetBusinessDaysParams{
				StartDate: "2024-01-01", // Monday
				EndDate:   "2024-01-07", // Sunday
			},
			check: func(t *testing.T, result map[string]interface{}) {
				businessDays, _ := result["business_days"].(int)
				if businessDays != 5 {
					t.Errorf("expected 5 business days, got %v", result["business_days"])
				}
				weekendDays, _ := result["weekend_days"].(int)
				if weekendDays != 2 {
					t.Errorf("expected 2 weekend days, got %v", result["weekend_days"])
				}
			},
		},
		{
			name: "calculate with holidays",
			params: GetBusinessDaysParams{
				StartDate:       "2024-01-01",
				EndDate:         "2024-01-07",
				ExcludeHolidays: []string{"2024-01-01"}, // New Year's Day
			},
			check: func(t *testing.T, result map[string]interface{}) {
				businessDays, _ := result["business_days"].(int)
				if businessDays != 4 {
					t.Errorf("expected 4 business days (excluding holiday), got %v", result["business_days"])
				}
				holidays, _ := result["holidays"].(int)
				if holidays != 1 {
					t.Errorf("expected 1 holiday, got %v", result["holidays"])
				}
			},
		},
		{
			name: "invalid date format",
			params: GetBusinessDaysParams{
				StartDate: "01-01-2024",
				EndDate:   "01-07-2024",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.Execute(ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.check != nil {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Errorf("Execute() returned wrong type: %T", result)
					return
				}
				tt.check(t, resultMap)
			}
		})
	}
}

func TestConvertUnixTimestampTool(t *testing.T) {
	tool := NewConvertUnixTimestampTool()
	ctx := context.Background()

	tests := []struct {
		name    string
		params  ConvertUnixTimestampParams
		wantErr bool
		check   func(t *testing.T, result map[string]interface{})
	}{
		{
			name: "datetime to unix seconds",
			params: ConvertUnixTimestampParams{
				Value:     "2024-01-01 00:00:00",
				Direction: "to_unix",
				Unit:      "seconds",
			},
			check: func(t *testing.T, result map[string]interface{}) {
				timestamp := result["timestamp"].(int64)
				if timestamp < 1704067200 || timestamp > 1704153600 { // Allow for timezone differences
					t.Errorf("unexpected timestamp: %d", timestamp)
				}
			},
		},
		{
			name: "unix seconds to datetime",
			params: ConvertUnixTimestampParams{
				Value:     "1704067200", // 2024-01-01 00:00:00 UTC
				Direction: "from_unix",
				Unit:      "seconds",
			},
			check: func(t *testing.T, result map[string]interface{}) {
				datetime := result["datetime"].(string)
				if !strings.Contains(datetime, "2024-01-01") && !strings.Contains(datetime, "2023-12-31") {
					t.Errorf("unexpected datetime: %s", datetime)
				}
			},
		},
		{
			name: "datetime to unix milliseconds",
			params: ConvertUnixTimestampParams{
				Value:     "2024-01-01 00:00:00",
				Direction: "to_unix",
				Unit:      "milliseconds",
			},
			check: func(t *testing.T, result map[string]interface{}) {
				timestamp := result["timestamp"].(int64)
				if timestamp < 1704067200000 {
					t.Errorf("timestamp too small for milliseconds: %d", timestamp)
				}
			},
		},
		{
			name: "invalid direction",
			params: ConvertUnixTimestampParams{
				Value:     "2024-01-01",
				Direction: "sideways",
			},
			wantErr: true,
		},
		{
			name: "invalid unit",
			params: ConvertUnixTimestampParams{
				Value:     "2024-01-01",
				Direction: "to_unix",
				Unit:      "minutes",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tool.Execute(ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.check != nil {
				resultMap, ok := result.(map[string]interface{})
				if !ok {
					t.Errorf("Execute() returned wrong type: %T", result)
					return
				}
				tt.check(t, resultMap)
			}
		})
	}
}
