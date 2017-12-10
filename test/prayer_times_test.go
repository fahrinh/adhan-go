package test

import (
	"testing"
	"time"
	"github.com/fahrinh/adhan-go"
)

func TestDaysSinceSolstice(t *testing.T) {
	daysSinceSolsticeTest(t, 11, 2016, 1, 1, 1)
	daysSinceSolsticeTest(t, 10, 2015, 12, 31, 1)
	daysSinceSolsticeTest(t, 10, 2016, 12, 31, 1)
}

func daysSinceSolsticeTest(t *testing.T, value, year, month, day int, latitude float64) {
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	dayOfYear := date.YearDay()
	d := adhan.DaysSinceSolstice(dayOfYear, date.Year(), latitude)
	if d != value {
		t.Error("expected: ", value, "actual: ", d)
	}
}