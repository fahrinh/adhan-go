package internal

import (
	"math"
	"github.com/fahrinh/adhan-go/doubleutil"
)

type SolarCoordinates struct {
	Declination          float64
	RightAscension       float64
	ApparentSiderealTime float64
}

func NewSolarCoordinates(julianDay float64) *SolarCoordinates {
	s := &SolarCoordinates{}

	T := JulianCentury(julianDay);
	L0 := MeanSolarLongitude(/* julianCentury */ T)
	Lp := MeanLunarLongitude(/* julianCentury */ T)
	Ω := AscendingLunarNodeLongitude(/* julianCentury */ T)
	λ := doubleutil.DegToRad(ApparentSolarLongitude(/* julianCentury*/ T, /* meanLongitude */ L0))

	θ0 := MeanSiderealTime(/* julianCentury */ T)
	ΔΨ := NutationInLongitude(/* julianCentury */ T, /* solarLongitude */ L0,
		/* lunarLongitude */ Lp, /* ascendingNode */ Ω)
	Δε := NutationInObliquity(/* julianCentury */ T, /* solarLongitude */ L0,
		/* lunarLongitude */ Lp, /* ascendingNode */ Ω)

	ε0 := MeanObliquityOfTheEcliptic(/* julianCentury */ T)
	εapp := doubleutil.DegToRad(ApparentObliquityOfTheEcliptic(
		/* julianCentury */ T, /* meanObliquityOfTheEcliptic */ ε0))

	/* Equation from Astronomical Algorithms page 165 */
	s.Declination = doubleutil.RadToDeg(math.Asin(math.Sin(εapp) * math.Sin(λ)))

	/* Equation from Astronomical Algorithms page 165 */
	s.RightAscension = doubleutil.UnwindAngle(
		doubleutil.RadToDeg(math.Atan2(math.Cos(εapp) * math.Sin(λ), math.Cos(λ))))

	/* Equation from Astronomical Algorithms page 88 */
	s.ApparentSiderealTime = θ0 + (((ΔΨ * 3600) * math.Cos(doubleutil.DegToRad(ε0 + Δε))) / 3600)

	return s
}
