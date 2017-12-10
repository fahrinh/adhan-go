package internal

import (
	"math"
	"github.com/fahrinh/adhan-go/doubleutil"
)

func MeanSolarLongitude(t float64) float64 {
	term1 := 280.4664567
	term2 := 36000.76983 * t
	term3 := 0.0003032 * math.Pow(t, 2)
	l0 := term1 + term2 + term3
	return doubleutil.UnwindAngle(l0)
}

func MeanLunarLongitude(t float64) float64 {
	term1 := 218.3165
	term2 := 481267.8813 * t
	lp := term1 + term2
	return doubleutil.UnwindAngle(lp)
}

func ApparentSolarLongitude(t, l0 float64) float64 {
	longitude := l0 + SolarEquationOfTheCenter(t, MeanSolarAnomaly(t))
	Ω := 125.04 - (1934.136 * t)
	λ := longitude - 0.00569 - (0.00478 * math.Sin(doubleutil.DegToRad(Ω)))
	return doubleutil.UnwindAngle(λ)
}

func AscendingLunarNodeLongitude(t float64) float64 {
	term1 := 125.04452
	term2 := 1934.136261 * t
	term3 := 0.0020708 * math.Pow(t, 2)
	term4 := math.Pow(t, 3) / 450000
	Ω := term1 - term2 + term3 + term4
	return doubleutil.UnwindAngle(Ω)
}

func MeanSolarAnomaly(t float64) float64 {
	term1 := 357.52911
	term2 := 35999.05029 * t
	term3 := 0.0001537 * math.Pow(t, 2)
	m := term1 + term2 + term3
	return doubleutil.UnwindAngle(m)
}

func SolarEquationOfTheCenter(t, m float64) float64 {
	mrad := doubleutil.DegToRad(m)
	term1 := (1.914602 - (0.004817 * t) - (0.000014 * math.Pow(t, 2))) * math.Sin(mrad)
	term2 := (0.019993 - (0.000101 * t)) * math.Sin(2 * mrad)
	term3 := 0.000289 * math.Sin(3 * mrad)
	return term1 + term2 + term3
}

func MeanObliquityOfTheEcliptic(t float64) float64 {
	term1 := 23.439291
	term2 := 0.013004167 * t
	term3 := 0.0000001639 * math.Pow(t, 2)
	term4 := 0.0000005036 * math.Pow(t, 3)
	return term1 - term2 - term3 + term4
}

func ApparentObliquityOfTheEcliptic(t, ε0 float64) float64 {
	O := 125.04 - (1934.136 * t)
	return ε0 + (0.00256 * math.Cos(doubleutil.DegToRad(O)))
}

func MeanSiderealTime(t float64) float64 {
	JD := (t * 36525) + 2451545.0
	term1 := 280.46061837
	term2 := 360.98564736629 * (JD - 2451545)
	term3 := 0.000387933 * math.Pow(t, 2)
	term4 := math.Pow(t, 3) / 38710000
	θ := term1 + term2 + term3 - term4
	return doubleutil.UnwindAngle(θ)
}

func NutationInLongitude(t, l0, lp, Ω float64) float64 {
	term1 := (-17.2/3600) * math.Sin(doubleutil.DegToRad(Ω))
	term2 :=  (1.32/3600) * math.Sin(2 * doubleutil.DegToRad(l0))
	term3 :=  (0.23/3600) * math.Sin(2 * doubleutil.DegToRad(lp))
	term4 :=  (0.21/3600) * math.Sin(2 * doubleutil.DegToRad(Ω))
	return term1 - term2 - term3 + term4
}

func NutationInObliquity(t, l0, lp, Ω float64) float64 {
	term1 :=  (9.2/3600) * math.Cos(doubleutil.DegToRad(Ω))
	term2 := (0.57/3600) * math.Cos(2 * doubleutil.DegToRad(l0))
	term3 := (0.10/3600) * math.Cos(2 * doubleutil.DegToRad(lp))
	term4 := (0.09/3600) * math.Cos(2 * doubleutil.DegToRad(Ω))
	return term1 + term2 + term3 - term4
}

func AltitudeOfCelestialBody(φ, δ, h float64) float64 {
	term1 := math.Sin(doubleutil.DegToRad(φ)) * math.Sin(doubleutil.DegToRad(δ))
	term2 := math.Cos(doubleutil.DegToRad(φ)) * math.Cos(doubleutil.DegToRad(δ)) * math.Cos(doubleutil.DegToRad(h))
	return doubleutil.RadToDeg(math.Asin(term1 + term2))
}

func ApproximateTransit(l, Θ0, α2 float64) float64 {
	lw := l * -1
	return doubleutil.NormalizeWithBound((α2 + lw - Θ0) / 360, 1)
}

func CorrectedTransit(m0, l, Θ0, α2, α1, α3 float64) float64 {
	lw := l * -1;
	θ := doubleutil.UnwindAngle(Θ0 + (360.985647 * m0));
	α := doubleutil.UnwindAngle(InterpolateAngles(
		/* value */ α2, /* previousValue */ α1, /* nextValue */ α3, /* factor */ m0));
	H := doubleutil.ClosestAngle(θ - lw - α);
	Δm := H / -360;
	return (m0 + Δm) * 24;
}

func CorrectedHourAngle(m0, h0 float64, coordinates *Coordinates, afterTransit bool,
Θ0, α2, α1, α3, δ2, δ1, δ3 float64) float64 {
	Lw := coordinates.Longitude * -1
	term1 := math.Sin(doubleutil.DegToRad(h0)) -
		(math.Sin(doubleutil.DegToRad(coordinates.Latitude)) * math.Sin(doubleutil.DegToRad(δ2)))
	term2 := math.Cos(doubleutil.DegToRad(coordinates.Latitude)) * math.Cos(doubleutil.DegToRad(δ2))
	H0 := doubleutil.RadToDeg(math.Acos(term1 / term2))
	var m float64
	if afterTransit {
		m = m0 + (H0 / 360)
	} else {
		m = m0 - (H0 / 360)
	}
	θ := doubleutil.UnwindAngle(Θ0 + (360.985647 * m))
	α := doubleutil.UnwindAngle(InterpolateAngles(
		/* value */ α2, /* previousValue */ α1, /* nextValue */ α3, /* factor */ m))
	δ := Interpolate(/* value */ δ2, /* previousValue */ δ1,
		/* nextValue */ δ3, /* factor */ m)
	H := (θ - Lw - α)
	h := AltitudeOfCelestialBody(/* observerLatitude */ coordinates.Latitude,
		/* Declination */ δ, /* localHourAngle */ H)
	term3 := h - h0
	term4 := 360 * math.Cos(doubleutil.DegToRad(δ)) *
		math.Cos(doubleutil.DegToRad(coordinates.Latitude)) * math.Sin(doubleutil.DegToRad(H))
	Δm := term3 / term4
	return (m + Δm) * 24
}

func Interpolate(y2, y1, y3, n float64) float64 {
	a := y2 - y1
	b := y3 - y2
	c := b - a
	return y2 + ((n/2) * (a + b + (n * c)))
}

func InterpolateAngles(y2, y1, y3, n float64) float64 {
	a := doubleutil.UnwindAngle(y2 - y1)
	b := doubleutil.UnwindAngle(y3 - y2)
	c := b - a
	return y2 + ((n/2) * (a + b + (n * c)))

}