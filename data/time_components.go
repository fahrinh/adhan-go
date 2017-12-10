package data

import (
	"math"
	"time"
)

type TimeComponents struct {
	Hours int
	Minutes int
	Seconds int
}

func NewTimeComponents(value float64) *TimeComponents {
	if math.IsInf(value, 0) || math.IsNaN(value) {
		return nil
	}

	hours := math.Floor(value)
	minutes := math.Floor((value - hours) * 60.0)
	seconds := math.Floor((value - (hours + minutes / 60.0)) * 60 * 60)
	return &TimeComponents{Hours:int(hours), Minutes:int(minutes), Seconds:int(seconds)}
}

func (t *TimeComponents) DateComponents(d time.Time) time.Time {
	c := time.Date(d.Year(), d.Month(), d.Day(), 0, t.Minutes, t.Seconds, d.Nanosecond(), time.UTC)
	c.Add(time.Hour * time.Duration(t.Hours))
	return c
}

