package common

import (
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"testing"
)

func TestCountWithZeroAsBounds(t *testing.T) {
	actual := determineCount(0, 0)

	AssertEqual(t, 1, actual)
}

func TestCountWithSameValueAsBounds(t *testing.T) {
	actual := determineCount(4, 4)

	AssertEqual(t, 4, actual)
}

func TestCountWithInMinAndMax(t *testing.T) {
	min, max := 4, 7
	actual := determineCount(min, max)

	if actual < min || actual > max {
		t.Errorf("Generated value '%v' is outside of expected range min: '%v', max: '%v'", actual, min, max)
	}
}
