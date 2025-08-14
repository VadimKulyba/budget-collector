package datetime

import (
	"log"
	"time"
)

const (
	periodDateFormat       = "01.2006"
	periodBorderDateFormat = "02.01.2006"
)

// GetMonthRangeByPeriod converts a period string to a formatted date range.
//
// Args:
//
//	period: period string in "MM.YYYY" format (e.g., "12.2024")
//
// Returns:
//
//	string: formatted date range string (e.g., "01.12.2024-31.12.2024")
//
// Example:
//
//	range := GetMonthRangeByPeriod("12.2024")
//	// Returns: "01.12.2024-31.12.2024"
func GetMonthRangeByPeriod(period string) string {
	// validate period
	startOfPeriod, err := time.Parse(periodDateFormat, period)
	if err != nil {
		log.Fatal("Invalid period format")
	}

	nextMonth := time.Date(
		startOfPeriod.Year(),
		startOfPeriod.Month()+1,
		1,
		0, 0, 0, 0,
		startOfPeriod.Location(),
	)

	endOfPeriod := nextMonth.Add(-24 * time.Hour)

	return startOfPeriod.Format(periodBorderDateFormat) + "-" + endOfPeriod.Format(periodBorderDateFormat)
}
