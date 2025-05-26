// ABOUTME: Example demonstrating all datetime tools available in go-flock
// ABOUTME: Shows comprehensive usage of datetime operations with go-llms agents

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/lexlapax/go-flock/pkg/tools"
)

func main() {
	fmt.Println("go-flock DateTime Tools Comprehensive Example")
	fmt.Println("=============================================")

	ctx := context.Background()

	// Example 1: Get current datetime in various formats
	demonstrateGetCurrentDatetime(ctx)

	// Example 2: Calculate duration between dates
	demonstrateCalculateDuration(ctx)

	// Example 3: Get weekday information
	demonstrateGetWeekday(ctx)

	// Example 4: Add/Subtract time units
	demonstrateAddSubtractTime(ctx)

	// Example 5: Convert between timezones
	demonstrateConvertTimezone(ctx)

	// Example 6: Parse datetime from various formats
	demonstrateParseDatetime(ctx)

	// Example 7: Format datetime to specific patterns
	demonstrateFormatDatetime(ctx)

	// Example 8: Compare datetimes
	demonstrateCompareDatetimes(ctx)

	// Example 9: Calculate business days
	demonstrateBusinessDays(ctx)

	// Example 10: Convert Unix timestamps
	demonstrateUnixTimestamp(ctx)

	fmt.Println("\nThese tools can be integrated with any go-llms agent for comprehensive datetime operations!")
}

func demonstrateGetCurrentDatetime(ctx context.Context) {
	fmt.Println("1. Getting Current DateTime in Different Formats:")
	fmt.Println("------------------------------------------------")

	tool := tools.NewGetCurrentDatetimeTool()

	formats := []struct {
		format string
		desc   string
	}{
		{"date", "ISO date"},
		{"time", "24-hour time"},
		{"datetime_iso", "ISO datetime"},
		{"filename", "Filename-safe"},
		{"log", "Log format"},
		{"relative", "Relative time"},
	}

	for _, f := range formats {
		params := tools.GetCurrentDateTimeParams{Format: f.format}
		result, err := tool.Execute(ctx, params)
		if err != nil {
			log.Printf("Error getting %s: %v", f.format, err)
			continue
		}
		if res, ok := result.(*tools.GetCurrentDateTimeResult); ok {
			fmt.Printf("   %-15s (%s): %s\n", f.format, f.desc, res.Formatted)
		}
	}
	fmt.Println()
}

func demonstrateCalculateDuration(ctx context.Context) {
	fmt.Println("2. Calculating Duration Between Dates:")
	fmt.Println("-------------------------------------")

	tool := tools.NewCalculateDurationTool()

	params := tools.CalculateDurationParams{
		Start: "2024-01-01",
		End:   "2024-12-31",
	}

	result, err := tool.Execute(ctx, params)
	if err != nil {
		log.Fatalf("Error calculating duration: %v", err)
	}

	if res, ok := result.(map[string]any); ok {
		fmt.Printf("   Between %s and %s:\n", params.Start, params.End)
		fmt.Printf("   - Days: %.0f\n", res["days"])
		fmt.Printf("   - Hours: %.0f\n", res["hours"])
		fmt.Printf("   - Human readable: %s\n", res["human"])
	}
	fmt.Println()
}

func demonstrateGetWeekday(ctx context.Context) {
	fmt.Println("3. Getting Weekday Information:")
	fmt.Println("------------------------------")

	tool := tools.NewGetWeekdayTool()

	dates := []string{"2024-12-25", "2024-01-01", "2024-07-04"}

	for _, date := range dates {
		params := tools.GetWeekdayParams{Date: date}
		result, err := tool.Execute(ctx, params)
		if err != nil {
			log.Printf("Error getting weekday for %s: %v", date, err)
			continue
		}
		if res, ok := result.(map[string]any); ok {
			fmt.Printf("   %s: %s (Weekend: %v)\n", date, res["weekday"], res["is_weekend"])
		}
	}
	fmt.Println()
}

func demonstrateAddSubtractTime(ctx context.Context) {
	fmt.Println("4. Adding/Subtracting Time Units:")
	fmt.Println("--------------------------------")

	tool := tools.NewAddSubtractTimeTool()

	operations := []tools.AddSubtractTimeParams{
		{Date: "2024-01-15", Amount: 5, Unit: "day", Operation: "add"},
		{Date: "2024-03-31", Amount: 1, Unit: "month", Operation: "subtract"},
		{Date: "2024-06-15 14:30:00", Amount: 3, Unit: "hour", Operation: "add"},
	}

	for _, params := range operations {
		result, err := tool.Execute(ctx, params)
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}
		if res, ok := result.(map[string]any); ok {
			fmt.Printf("   %s %s %d %s(s) = %s\n",
				params.Date, params.Operation, params.Amount, params.Unit, res["result"])
		}
	}
	fmt.Println()
}

func demonstrateConvertTimezone(ctx context.Context) {
	fmt.Println("5. Converting Between Timezones:")
	fmt.Println("-------------------------------")

	tool := tools.NewConvertTimezoneTool()

	params := tools.ConvertTimezoneParams{
		Datetime:     "2024-01-01 12:00:00",
		FromTimezone: "America/New_York",
		ToTimezone:   "Asia/Tokyo",
	}

	result, err := tool.Execute(ctx, params)
	if err != nil {
		log.Fatalf("Error converting timezone: %v", err)
	}

	if res, ok := result.(map[string]any); ok {
		fmt.Printf("   Original: %s\n", res["original"])
		fmt.Printf("   Converted: %s\n", res["converted"])
		fmt.Printf("   Offset change: %s → %s\n", res["original_offset"], res["converted_offset"])
	}
	fmt.Println()
}

func demonstrateParseDatetime(ctx context.Context) {
	fmt.Println("6. Parsing DateTime from Various Formats:")
	fmt.Println("----------------------------------------")

	tool := tools.NewParseDatetimeTool()

	dateStrings := []string{
		"2024-01-01T15:04:05Z",
		"1704110400", // Unix timestamp
		"Jan 2, 2024",
		"2024-01-02",
	}

	for _, dateStr := range dateStrings {
		params := tools.ParseDatetimeParams{DatetimeString: dateStr}
		result, err := tool.Execute(ctx, params)
		if err != nil {
			log.Printf("Error parsing %s: %v", dateStr, err)
			continue
		}
		if res, ok := result.(map[string]any); ok {
			fmt.Printf("   '%s' → %s (format: %s)\n",
				dateStr, res["parsed"], res["detected_format"])
		}
	}
	fmt.Println()
}

func demonstrateFormatDatetime(ctx context.Context) {
	fmt.Println("7. Formatting DateTime to Specific Patterns:")
	fmt.Println("-------------------------------------------")

	tool := tools.NewFormatDatetimeTool()

	formats := []struct {
		format string
		desc   string
	}{
		{"rfc3339", "RFC3339"},
		{"human", "Human readable"},
		{"short", "Short format"},
		{"2006-01-02 15:04 MST", "Custom"},
	}

	datetime := "2024-01-15 14:30:00"

	for _, f := range formats {
		params := tools.FormatDatetimeParams{
			Datetime:     datetime,
			OutputFormat: f.format,
		}
		result, err := tool.Execute(ctx, params)
		if err != nil {
			log.Printf("Error formatting: %v", err)
			continue
		}
		if res, ok := result.(map[string]any); ok {
			fmt.Printf("   %-15s: %s\n", f.desc, res["formatted"])
		}
	}
	fmt.Println()
}

func demonstrateCompareDatetimes(ctx context.Context) {
	fmt.Println("8. Comparing DateTimes:")
	fmt.Println("----------------------")

	tool := tools.NewCompareDatetimesTool()

	params := tools.CompareDatetimesParams{
		Datetime1: "2024-01-01 10:00:00",
		Datetime2: "2024-01-05 15:30:00",
	}

	result, err := tool.Execute(ctx, params)
	if err != nil {
		log.Fatalf("Error comparing datetimes: %v", err)
	}

	if res, ok := result.(map[string]any); ok {
		fmt.Printf("   Comparing %s and %s:\n", params.Datetime1, params.Datetime2)
		fmt.Printf("   - First is before second: %v\n", res["is_before"])
		fmt.Printf("   - Difference: %.1f days (%.1f hours)\n",
			res["difference_days"], res["difference_hours"])
		fmt.Printf("   - Human readable: %s\n", res["difference_human"])
	}
	fmt.Println()
}

func demonstrateBusinessDays(ctx context.Context) {
	fmt.Println("9. Calculating Business Days:")
	fmt.Println("----------------------------")

	tool := tools.NewGetBusinessDaysTool()

	params := tools.GetBusinessDaysParams{
		StartDate:       "2024-01-01",
		EndDate:         "2024-01-15",
		ExcludeHolidays: []string{"2024-01-01", "2024-01-15"}, // New Year's Day, MLK Day
	}

	result, err := tool.Execute(ctx, params)
	if err != nil {
		log.Fatalf("Error calculating business days: %v", err)
	}

	if res, ok := result.(map[string]any); ok {
		fmt.Printf("   Between %s and %s:\n", params.StartDate, params.EndDate)
		fmt.Printf("   - Total days: %.0f\n", res["total_days"])
		fmt.Printf("   - Business days: %.0f\n", res["business_days"])
		fmt.Printf("   - Weekend days: %.0f\n", res["weekend_days"])
		fmt.Printf("   - Holidays: %.0f\n", res["holidays"])
	}
	fmt.Println()
}

func demonstrateUnixTimestamp(ctx context.Context) {
	fmt.Println("10. Converting Unix Timestamps:")
	fmt.Println("------------------------------")

	tool := tools.NewConvertUnixTimestampTool()

	// Convert datetime to unix
	params1 := tools.ConvertUnixTimestampParams{
		Value:     "2024-01-01 12:00:00",
		Direction: "to_unix",
		Unit:      "seconds",
	}

	result1, err := tool.Execute(ctx, params1)
	if err != nil {
		log.Fatalf("Error converting to unix: %v", err)
	}

	if res, ok := result1.(map[string]any); ok {
		fmt.Printf("   DateTime to Unix: %s → %d\n", params1.Value, res["timestamp"])

		// Convert back from unix
		params2 := tools.ConvertUnixTimestampParams{
			Value:     fmt.Sprintf("%d", res["timestamp"]),
			Direction: "from_unix",
			Unit:      "seconds",
		}

		result2, err := tool.Execute(ctx, params2)
		if err != nil {
			log.Fatalf("Error converting from unix: %v", err)
		}

		if res2, ok := result2.(map[string]any); ok {
			fmt.Printf("   Unix to DateTime: %s → %s\n", params2.Value, res2["datetime"])
		}
	}
}
