package doubleutil

import "math"

func NormalizeWithBound(value, max float64) float64 {
	return value - (max * (math.Floor(value/max)))
}

func UnwindAngle(value float64) float64 {
	return NormalizeWithBound(value, 360)
}

func ClosestAngle(angle float64) float64 {
	if angle >= -180 && angle <= 180 {
		return angle
	}
	return angle - (360 * Round(angle / 360))
}

func Round(a float64) float64 {
	if a < 0 {
		return math.Ceil(a - 0.5)
	}
	return math.Floor(a + 0.5)
}

func DegToRad(angdeg float64) float64 {
	return angdeg / 180 * math.Pi
}

func RadToDeg(angrad float64) float64 {
	return angrad * 180 / math.Pi
}