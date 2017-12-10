package adhan

type HighLatitudeRule int

const (
	MIDDLE_OF_THE_NIGHT HighLatitudeRule = iota + 1
	SEVENTH_OF_THE_NIGHT
	TWILIGHT_ANGLE
)
