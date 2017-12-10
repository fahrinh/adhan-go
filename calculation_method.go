package adhan

type CalculationMethod int

const (
	MUSLIM_WORLD_LEAGUE     CalculationMethod = iota + 1
	EGYPTIAN
	KARACHI
	UMM_AL_QURA
	GULF
	MOON_SIGHTING_COMMITTEE
	NORTH_AMERICA
	KUWAIT
	QATAR
	OTHER
)

func (c CalculationMethod) GetParameters() *CalculationParameters {
	p := &CalculationParameters{}
	p.Method = c

	switch c {
	case MUSLIM_WORLD_LEAGUE:
		p.FajrAngle = 18.0
		p.IshaAngle = 17.0
	case EGYPTIAN:
		p.FajrAngle = 20.0
		p.IshaAngle = 18.0
	case KARACHI:
		p.FajrAngle = 18.0
		p.IshaAngle = 18.0
	case UMM_AL_QURA:
		p.FajrAngle = 18.5
		p.IshaInterval = 90
	case GULF:
		p.FajrAngle = 19.5
		p.IshaInterval = 90
	case MOON_SIGHTING_COMMITTEE:
		p.FajrAngle = 18.0
		p.IshaAngle = 18.0
	case NORTH_AMERICA:
		p.FajrAngle = 15.0
		p.IshaAngle = 15.0
	case KUWAIT:
		p.FajrAngle = 18.0
		p.IshaAngle = 17.5
	case QATAR:
		p.FajrAngle = 18.0
		p.IshaInterval = 90
	case OTHER:
		p.FajrAngle = 0.0
		p.IshaAngle = 0.0
	}

	return p
}
