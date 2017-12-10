package adhan

import "github.com/fahrinh/adhan-go/internal"

type Madhab int

const (
	SHAFI = iota + 1
	HANAFI
)

func (m Madhab) GetShadowLength() internal.ShadowLength {
	var s internal.ShadowLength

	switch m {
	case SHAFI:
		s =  internal.SINGLE
	case HANAFI:
		s =  internal.DOUBLE
	}

	return s
}