package common

import (
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"testing"
)

func TestAmountWithZeroAsBounds(t *testing.T) {
	actual := determineAmount(0, 0)

	AssertEqual(t, 1, actual)
}

func TestAmountWithSameValueAsBounds(t *testing.T) {
	actual := determineAmount(4, 4)

	AssertEqual(t, 4, actual)
}

func TestAmountWithInMinAndMax(t *testing.T) {
	min, max := 4, 7
	actual := determineAmount(min, max)

	if actual < min || actual > max {
		t.Errorf("Generated value '%v' is outside of expected range min: '%v', max: '%v'", actual, min, max)
	}
}
