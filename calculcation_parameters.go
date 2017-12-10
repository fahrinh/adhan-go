package adhan

type CalculationParameters struct {
	Method           CalculationMethod
	FajrAngle        float64
	IshaAngle        float64
	IshaInterval     int
	Madhab           Madhab
	HighLatitudeRule HighLatitudeRule
	Adjustments      *PrayerAdjustments
}

type NightPortions struct {
	fajr float64
	isha float64
}

func NewCalculationParameters() *CalculationParameters {
	c := &CalculationParameters{}
	c.Method = OTHER
	c.Madhab = SHAFI
	c.HighLatitudeRule = MIDDLE_OF_THE_NIGHT
	c.Adjustments = &PrayerAdjustments{}
	return c
}

func (c *CalculationParameters) NightPortions() *NightPortions {
	var n *NightPortions;

	switch c.HighLatitudeRule {
	case MIDDLE_OF_THE_NIGHT:
		n =  &NightPortions{fajr:1.0 / 2.0, isha:1.0 / 2.0}
	case SEVENTH_OF_THE_NIGHT:
		n = &NightPortions{fajr:1.0 / 7.0, isha:1.0 / 7.0}
	case TWILIGHT_ANGLE:
		n =  &NightPortions{fajr:c.FajrAngle / 60.0, isha:c.IshaAngle / 60.0}
	}

	return n
}
