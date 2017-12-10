package internal

import "time"

func JulianDayFromDate(date time.Time) float64  {
	return JulianDay(date.Year(), int(date.Month()), date.Day(), float64(date.Hour() + date.Minute() / 60.0))
}

func JulianDay(year, month, day int, hours float64) float64 {
	var Y int
	if month > 2 {
		Y = year
	} else {
		Y = year - 1
	}

	var M int
	if month > 2 {
		M = month
	} else {
		M = month + 12
	}

	D := float64(day) + (hours / 24)

	A := Y/100
	B := 2 - A + (A/4)

	i0 := int(float64(365.25) * (float64(Y) + 4716))
	i1 := int(float64(30.6001) * (float64(M) + 1))
	return float64(i0) + float64(i1) + float64(D) + float64(B) - 1524.5
}

func JulianCentury(JD float64) float64 {
	return (JD - 2451545.0) / 36525
}