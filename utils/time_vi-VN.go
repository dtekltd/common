package utils

import (
	"fmt"
	"time"
)

func FormatFullVnDate(date time.Time) string {
	// Get the day of the week in Vietnamese
	weekdays := []string{"Chủ nhật", "Thứ 2", "Thứ 3", "Thứ 4", "Thứ 5", "Thứ 6", "Thứ 7"}
	weekday := weekdays[date.Weekday()]

	// Format the date in the desired Vietnamese format
	return fmt.Sprintf("%s, ngày %d tháng %d năm %d", weekday, date.Day(), int(date.Month()), date.Year())
}
