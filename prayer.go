package adhan

type Prayer int

const (
	NONE Prayer = iota + 1
	FAJR
	SUNRISE
	DHUHR
	ASR
	MAGHRIB
	ISHA
)
