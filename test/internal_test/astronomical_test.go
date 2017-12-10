package internal_test

import (
	"testing"
	"github.com/fahrinh/adhan-go/internal"
	"github.com/fahrinh/adhan-go/doubleutil"
	"github.com/fahrinh/adhan-go/test"
)

func TestSolarCoordinates(t *testing.T) {
	jd := internal.JulianDay(1992, 10, 13, 0)
	solar := internal.NewSolarCoordinates(jd)

	T := internal.JulianCentury(jd)
	L0 := internal.MeanSolarLongitude(T)
	ε0 := internal.MeanObliquityOfTheEcliptic(T)
	εapp := internal.ApparentObliquityOfTheEcliptic(T, ε0)
	M := internal.MeanSolarAnomaly(T)
	C := internal.SolarEquationOfTheCenter(T, M)
	λ := internal.ApparentSolarLongitude(T, L0)
	δ := solar.Declination
	α := doubleutil.UnwindAngle(solar.RightAscension)

	test.IsWithin(t, T, 0.00000000001, -0.072183436)
	test.IsWithin(t, L0, 0.00001, 201.80720)
	test.IsWithin(t, ε0, 0.00001, 23.44023)
	test.IsWithin(t, εapp, 0.00001, 23.43999)
	test.IsWithin(t, M, 0.00001, 278.99397)
	test.IsWithin(t, C, 0.00001, -1.89732)
	test.IsWithin(t, λ, 0.00002, 199.90895)
	test.IsWithin(t, δ, 0.00001, -7.78507)
	test.IsWithin(t, α, 0.00001, 198.38083)

	jd = internal.JulianDay(1987, 4, 10, 0.0)
	solar = internal.NewSolarCoordinates(jd)
	T = internal.JulianCentury(jd)

	θ0 := internal.MeanSiderealTime(T)
	θapp := solar.ApparentSiderealTime
	Ω := internal.AscendingLunarNodeLongitude(T)
	ε0 = internal.MeanObliquityOfTheEcliptic(T)
	L0 = internal.MeanSolarLongitude(T)
	Lp := internal.MeanLunarLongitude(T)
	ΔΨ := internal.NutationInLongitude(T, L0, Lp, Ω)
	Δε := internal.NutationInObliquity(T, L0, Lp, Ω)
	ε := ε0 + Δε

	test.IsWithin(t, θ0, 0.000001, 197.693195)
	test.IsWithin(t, θapp, 0.0001, 197.6922295833)

	test.IsWithin(t, Ω, 0.0001, 11.2531)
	test.IsWithin(t, ΔΨ, 0.0001, -0.0010522)
	test.IsWithin(t, Δε, 0.00001, 0.0026230556)
	test.IsWithin(t, ε0, 0.000001, 23.4409463889)
	test.IsWithin(t, ε, 0.00001, 23.4435694444)

}


