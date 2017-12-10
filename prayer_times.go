package adhan

import (
	"time"
	"github.com/fahrinh/adhan-go/data"
	"github.com/fahrinh/adhan-go/internal"
	"math"
	"github.com/fahrinh/adhan-go/doubleutil"
)

type PrayerTimes struct {
	Fajr    *time.Time
	Sunrise *time.Time
	Dhuhr   *time.Time
	Asr     *time.Time
	Maghrib *time.Time
	Isha    *time.Time
}

func NewPrayerTimes(coordinates *internal.Coordinates, dc *data.DateComponents, params *CalculationParameters) *PrayerTimes {
	p := &PrayerTimes{}

	date := data.ResolveTime(dc.Year, dc.Month, dc.Day)

	var tempFajr *time.Time
	var tempSunrise *time.Time
	var tempDhuhr *time.Time
	var tempAsr *time.Time
	var tempMaghrib *time.Time
	var tempIsha *time.Time

	year := date.Year()
	dayOfYear := date.YearDay()

	solarTime := internal.NewSolarTime(date, coordinates)

	timeComponents := data.NewTimeComponents(solarTime.Transit)
	var transit *time.Time
	if timeComponents == nil {
		transit = nil
	} else {
		d := timeComponents.DateComponents(date)
		transit = &d
	}

	timeComponents = data.NewTimeComponents(solarTime.Sunrise)
	var sunriseComponents *time.Time
	if timeComponents == nil {
		sunriseComponents = nil
	} else {
		d := timeComponents.DateComponents(date)
		sunriseComponents = &d
	}

	timeComponents = data.NewTimeComponents(solarTime.Sunset)
	var sunsetComponents *time.Time
	if timeComponents == nil {
		sunsetComponents = nil
	} else {
		d := timeComponents.DateComponents(date)
		sunsetComponents = &d
	}

	error := transit == nil || sunriseComponents == nil || sunsetComponents == nil

	if !error {
		tempDhuhr = transit
		tempSunrise = sunriseComponents
		tempMaghrib = sunsetComponents

		s := params.Madhab.GetShadowLength()
		timeComponents = data.NewTimeComponents(
			solarTime.Afternoon(&s))
		if timeComponents != nil {
			d := timeComponents.DateComponents(date)
			tempAsr = &d
		}

		tomorrowSunrise := sunriseComponents.AddDate(0, 0, 1)
		night := tomorrowSunrise.Nanosecond() - sunsetComponents.Nanosecond()

		timeComponents = data.NewTimeComponents(
			solarTime.HourAngle(-params.FajrAngle, false))
		if timeComponents != nil {
			d := timeComponents.DateComponents(date)
			tempFajr = &d
		}

		if params.Method == MOON_SIGHTING_COMMITTEE && coordinates.Latitude >= 55 {
			a := sunriseComponents.Add(time.Second * (-1 * (time.Duration(night) / 7000)))
			tempFajr = &a
		}

		nightPortions := params.NightPortions()

		var safeFajr *time.Time
		if params.Method == MOON_SIGHTING_COMMITTEE {
			s := seasonAdjustedMorningTwilight(coordinates.Latitude, dayOfYear, year, *sunriseComponents)
			safeFajr = &s
		} else {
			portion := nightPortions.fajr
			nightFraction := int(portion) * night / 1000
			s := sunriseComponents.Add(time.Second * -1 * time.Duration(nightFraction))
			safeFajr = &s
		}

		if tempFajr == nil || tempFajr.Before(*safeFajr) {
			tempFajr = safeFajr
		}

		if params.IshaInterval > 0 {
			m := tempMaghrib.Add(time.Second * time.Duration(params.IshaInterval) * 60)
			tempIsha = &m
		} else {
			timeComponents = data.NewTimeComponents(
				solarTime.HourAngle(-params.IshaAngle, true))
			if timeComponents != nil {
				d := timeComponents.DateComponents(date)
				tempIsha = &d
			}

			if params.Method == MOON_SIGHTING_COMMITTEE && coordinates.Latitude >= 55 {
				nightFraction := night / 7000
				s := sunsetComponents.Add(time.Second * time.Duration(nightFraction))
				tempIsha = &s
			}

			var safeIsha *time.Time
			if params.Method == MOON_SIGHTING_COMMITTEE {
				s := seasonAdjustedEveningTwilight(
					coordinates.Latitude, dayOfYear, year, *sunsetComponents)
				safeIsha = &s
			} else {
				portion := nightPortions.isha
				nightFraction := int(portion) * night / 1000
				s := sunsetComponents.Add(time.Second * time.Duration(nightFraction))
				safeIsha = &s
			}

			if (tempIsha == nil || tempIsha.After(*safeIsha)) {
				tempIsha = safeIsha
			}
		}
	}

	var dhuhrOffsetInMinutes int
	if params.Method == MOON_SIGHTING_COMMITTEE {
		dhuhrOffsetInMinutes = 5
	} else if params.Method == UMM_AL_QURA ||
		params.Method == GULF ||
		params.Method == QATAR {
		dhuhrOffsetInMinutes = 0
	} else {
		dhuhrOffsetInMinutes = 1
	}

	var maghribOffsetInMinutes int
	if params.Method == MOON_SIGHTING_COMMITTEE {
		maghribOffsetInMinutes = 3
	} else {
		maghribOffsetInMinutes = 0
	}

	if error || tempAsr == nil {
		p.Fajr = nil
		p.Sunrise = nil
		p.Dhuhr = nil
		p.Asr = nil
		p.Maghrib = nil
		p.Isha = nil
	} else {
		f := data.RoundedMinute(tempFajr.Add(time.Minute * time.Duration(params.Adjustments.Fajr)))
		p.Fajr = &f
		s := data.RoundedMinute(tempSunrise.Add(time.Minute * time.Duration(params.Adjustments.Sunrise)))
		p.Sunrise = &s
		d := data.RoundedMinute(tempDhuhr.Add(time.Minute * time.Duration(params.Adjustments.Dhuhr + dhuhrOffsetInMinutes)))
		p.Dhuhr = &d
		a := data.RoundedMinute(tempAsr.Add(time.Minute * time.Duration(params.Adjustments.Asr)))
		p.Asr = &a
		m := data.RoundedMinute(tempMaghrib.Add(time.Minute * time.Duration(params.Adjustments.Maghrib + maghribOffsetInMinutes)))
		p.Maghrib = &m
		i := data.RoundedMinute(tempIsha.Add(time.Minute * time.Duration(params.Adjustments.Isha)))
		p.Isha = &i
	}

	return p
}

func (p *PrayerTimes) CurrentPrayer() *Prayer {
	return p.CurrentPrayerFromTime(time.Now())
}

func (p *PrayerTimes) CurrentPrayerFromTime(time time.Time) *Prayer {
	when := time.Nanosecond()

	var r Prayer
	if p.Isha.Nanosecond() - when <= 0 {
		r = ISHA
	} else if p.Maghrib.Nanosecond() - when <= 0 {
		r = MAGHRIB
	} else if p.Asr.Nanosecond() - when <= 0 {
		r = ASR
	} else if p.Dhuhr.Nanosecond() - when <= 0 {
		r = DHUHR
	} else if p.Sunrise.Nanosecond() - when <= 0 {
		r = SUNRISE
	} else if p.Fajr.Nanosecond() - when <= 0 {
		r = FAJR
	} else {
		r = NONE
	}

	return &r
}

func (p *PrayerTimes) NextPrayer() *Prayer {
	return p.NextPrayerFromTime(time.Now())
}

func (p *PrayerTimes) NextPrayerFromTime(time time.Time) *Prayer {
	when := time.Nanosecond()

	var r Prayer
	if p.Isha.Nanosecond() - when <= 0 {
		r = NONE
	} else if p.Maghrib.Nanosecond() - when <= 0 {
		r = ISHA
	} else if p.Asr.Nanosecond() - when <= 0 {
		r = MAGHRIB
	} else if p.Dhuhr.Nanosecond() - when <= 0 {
		r = ASR
	} else if p.Sunrise.Nanosecond() - when <= 0 {
		r = DHUHR
	} else if p.Fajr.Nanosecond() - when <= 0 {
		r = SUNRISE
	} else {
		r = FAJR
	}

	return &r
}

func (p *PrayerTimes) TimeForPrayer(prayer *Prayer) *time.Time {
	var r *time.Time

	switch *prayer {
	case FAJR:
		r = p.Fajr
	case SUNRISE:
		r = p.Sunrise
	case DHUHR:
		r = p.Dhuhr
	case ASR:
		r = p.Asr
	case MAGHRIB:
		r = p.Maghrib
	case ISHA:
		r = p.Isha
	case NONE:
	default:
		r = nil
	}

	return r
}


func seasonAdjustedMorningTwilight(latitude float64, day, year int, sunrise time.Time) time.Time {
	a := 75 + ((28.65 / 55.0) * math.Abs(latitude))
	b := 75 + ((19.44 / 55.0) * math.Abs(latitude))
	c := 75 + ((32.74 / 55.0) * math.Abs(latitude))
	d := 75 + ((48.10 / 55.0) * math.Abs(latitude))

	var adjustment float64
	dyy := float64(DaysSinceSolstice(day, year, latitude))
	if ( dyy < 91) {
		adjustment = a + ( b-a )/91.0*dyy
	} else if ( dyy < 137) {
		adjustment = b + ( c-b )/46.0*( dyy-91 )
	} else if ( dyy < 183 ) {
		adjustment = c + ( d-c )/46.0*( dyy-137 )
	} else if ( dyy < 229 ) {
		adjustment = d + ( c-d )/46.0*( dyy-183 )
	} else if ( dyy < 275 ) {
		adjustment = c + ( b-c )/46.0*( dyy-229 )
	} else {
		adjustment = b + ( a-b )/91.0*( dyy-275 )
	}

	return sunrise.Add(time.Second * time.Duration(-doubleutil.Round(adjustment * 60.0)))
}

func seasonAdjustedEveningTwilight(latitude float64, day, year int, sunset time.Time) time.Time {
	a := 75 + ((25.60 / 55.0) * math.Abs(latitude))
	b := 75 + ((2.050 / 55.0) * math.Abs(latitude))
	c := 75 - ((9.210 / 55.0) * math.Abs(latitude))
	d := 75 + ((6.140 / 55.0) * math.Abs(latitude))

	var adjustment float64
	dyy := float64(DaysSinceSolstice(day, year, latitude))
	if ( dyy < 91) {
		adjustment = a + ( b-a )/91.0*dyy
	} else if ( dyy < 137) {
		adjustment = b + ( c-b )/46.0*( dyy-91 )
	} else if ( dyy < 183 ) {
		adjustment = c + ( d-c )/46.0*( dyy-137 )
	} else if ( dyy < 229 ) {
		adjustment = d + ( c-d )/46.0*( dyy-183 )
	} else if ( dyy < 275 ) {
		adjustment = c + ( b-c )/46.0*( dyy-229 )
	} else {
		adjustment = b + ( a-b )/91.0*( dyy-275 )
	}

	return sunset.Add(time.Second * time.Duration(-doubleutil.Round(adjustment * 60.0)))
}

func DaysSinceSolstice(dayOfYear, year int, latitude float64) int {
	var daysSinceSolistice int
	northernOffset := 10
	isLeapYear := data.IsLeapYear(year)

	var southernOffset int
	if isLeapYear {
		southernOffset = 173
	} else {
		southernOffset = 172
	}

	var daysInYear int
	if isLeapYear {
		daysInYear = 366
	} else {
		daysInYear = 365
	}

	if latitude >= 0 {
		daysSinceSolistice = dayOfYear + northernOffset
		if daysSinceSolistice >= daysInYear {
			daysSinceSolistice = daysSinceSolistice - daysInYear
		}
	} else {
		daysSinceSolistice = dayOfYear - southernOffset
		if (daysSinceSolistice < 0) {
			daysSinceSolistice = daysSinceSolistice + daysInYear
		}
	}
	return daysSinceSolistice
}
