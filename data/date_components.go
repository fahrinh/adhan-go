package data

import "time"

type DateComponents struct {
	Year int
	Month int
	Day int
}

func NewDateComponents(t time.Time) *DateComponents {
	return &DateComponents{Year:t.Year(),Month:int(t.Month()),Day:t.Day()}
}

func (c *DateComponents) toTime() time.Time {
	return time.Date(c.Year,time.Month(c.Month),c.Day,0,0,0,0,time.UTC)
}