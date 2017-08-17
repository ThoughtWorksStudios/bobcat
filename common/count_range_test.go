package common

import (
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"testing"
)

func TestCountWithZeroAsBounds(t *testing.T) {
	actual := determineCount(0, 0)

	AssertEqual(t, int64(0), actual)
}

func TestCountWithSameValueAsBounds(t *testing.T) {
	actual := determineCount(4, 4)

	AssertEqual(t, int64(4), actual)
}

func TestCountWithInMinAndMax(t *testing.T) {
	min, max := int64(4), int64(7)
	actual := determineCount(min, max)

	if actual < min || actual > max {
		t.Errorf("Generated value '%v' is outside of expected range min: '%v', max: '%v'", actual, min, max)
	}
}

func TestValidateWithNegativeBounds(t *testing.T) {
	c := &CountRange{Min: -1, Max: 1}
	ExpectsError(t, "Count range bounds must not be negative", c.Validate())
}

func TestValidateWithMaxLessThanMin(t *testing.T) {
	c := &CountRange{Min: 10, Max: 5}
	ExpectsError(t, "Count range max cannot be less than min", c.Validate())
}
