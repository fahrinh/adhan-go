package internal

import (
	"time"
	"math"
	"github.com/fahrinh/adhan-go/doubleutil"
)

type SolarTime struct {
	Transit float64
	Sunrise float64
	Sunset  float64

	observer           *Coordinates
	solar              *SolarCoordinates
	prevSolar          *SolarCoordinates
	nextSolar          *SolarCoordinates
	approximateTransit float64
}

func NewSolarTime(today time.Time, coordinates *Coordinates) *SolarTime {
	s := &SolarTime{}

	tomorrow := today.AddDate(0, 0, 1)
	yesterday := today.AddDate(0, 0, -1)

	s.prevSolar = NewSolarCoordinates(JulianDayFromDate(yesterday))
	s.solar = NewSolarCoordinates(JulianDayFromDate(today))
	s.nextSolar = NewSolarCoordinates(JulianDayFromDate(tomorrow))

	s.approximateTransit = ApproximateTransit(coordinates.Longitude,
		s.solar.ApparentSiderealTime, s.solar.RightAscension)
	solarAltitude := -50.0 / 60.0

	s.observer = coordinates
	s.Transit = CorrectedTransit(s.approximateTransit, coordinates.Longitude,
		s.solar.ApparentSiderealTime, s.solar.RightAscension, s.prevSolar.RightAscension,
		s.nextSolar.RightAscension)
	s.Sunrise = CorrectedHourAngle(s.approximateTransit, solarAltitude,
		coordinates, false, s.solar.ApparentSiderealTime, s.solar.RightAscension,
		s.prevSolar.RightAscension, s.nextSolar.RightAscension, s.solar.Declination,
		s.prevSolar.Declination, s.nextSolar.Declination)
	s.Sunset = CorrectedHourAngle(s.approximateTransit, solarAltitude,
		coordinates, true, s.solar.ApparentSiderealTime, s.solar.RightAscension,
		s.prevSolar.RightAscension, s.nextSolar.RightAscension, s.solar.Declination,
		s.prevSolar.Declination, s.nextSolar.Declination)

	return s
}

func (s *SolarTime) HourAngle(angle float64, afterTransit bool) float64 {
	return CorrectedHourAngle(s.approximateTransit, angle,
		s.observer, afterTransit, s.solar.ApparentSiderealTime, s.solar.RightAscension,
		s.prevSolar.RightAscension, s.nextSolar.RightAscension, s.solar.Declination,
		s.prevSolar.Declination, s.nextSolar.Declination)

}

func (s *SolarTime) Afternoon(shadowLength *ShadowLength) float64 {
	tangent := math.Abs(s.observer.Latitude - s.solar.Declination)
	inverse := float64(*shadowLength) + math.Tan(doubleutil.DegToRad(tangent))
	angle := doubleutil.RadToDeg(math.Atan(1.0 / inverse))

	return s.HourAngle(angle, true)
}
