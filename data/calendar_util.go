package data

import (
	"time"
	"github.com/fahrinh/adhan-go/doubleutil"
)

func IsLeapYear(year int) bool {
	return year % 4 == 0 && !(year % 100 == 0 && year % 400 != 0)
}

func RoundedMinute(when time.Time) time.Time {
	minute := when.Minute()
	second := when.Second()

	t := time.Date(when.Year(),when.Month(), when.Day(), when.Hour(),
		(minute + int(doubleutil.Round(float64(second / 60)))), 0, when.Nanosecond(),
		when.Location())
	return t
}

func ResolveTime(year, month, day int) time.Time {
	return time.Date(year,time.Month(month),day,0,0,0,0,time.UTC)
}

