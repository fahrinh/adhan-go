package test

import (
	_ "github.com/smartystreets/assertions"
	"testing"
	"math"
)

func IsWithin(t *testing.T, actual, tolerance, expected float64) {

	lowerBound := expected-tolerance
	upperBound := expected+tolerance

	if !(math.Abs(actual - expected) <= math.Abs(tolerance)) {
		t.Error("actual: ", actual, "lowerBound: ", lowerBound, "upperBound: ", upperBound)
	}
	//assertions.ShouldBeBetween(actual, lowerBound, upperBound)
}
