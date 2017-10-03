package generator

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"testing"
)

// TODO: Flaky??
func TestWeightDistributionOneEnum(t *testing.T) {
	weights := []float64{80.0, 20.0}
	intervals := []FieldType{
		NewLiteralType("one"),
		NewLiteralType("two"),
	}

	dist, err := NewDistribution(WEIGHT_DIST, weights, intervals)
	AssertNil(t, err, "Should not receive error during distribution construction")

	count := 10000

	resultIntervalOne := []string{}
	resultIntervalTwo := []string{}

	for i := 0; i < count; i++ {
		v, err := dist.One(nil, nil, nil)
		AssertNil(t, err, "Should not receive error")
		value := v.(string)

		switch value {
		case "one":
			resultIntervalOne = append(resultIntervalOne, value)
		case "two":
			resultIntervalTwo = append(resultIntervalTwo, value)
		default:
			t.Errorf("Should not have generated a value outside of the domain!")
		}
	}

	AssertEqual(t, 8, int(RoundFloat(float64(len(resultIntervalOne))/1000.0, 1)), "Interval 1 should be approximately 80%")
	AssertEqual(t, 2, int(RoundFloat(float64(len(resultIntervalTwo))/1000.0, 1)), "Interval 2 should be approximately 80%")
}

func TestWeightDistributionOne(t *testing.T) {
	weights := []float64{50.0, 50.0}
	intervalOne := NewLiteralType(int64(10))
	intervalTwo := NewLiteralType(int64(20))

	dist, err := NewDistribution(WEIGHT_DIST, weights, []FieldType{intervalOne, intervalTwo})
	AssertNil(t, err, "Should not receive error during distribution construction")

	count := 10

	resultIntervalOne := []int64{}
	resultIntervalTwo := []int64{}

	for i := 0; i < count; i++ {
		v, err := dist.One(nil, nil, nil)
		AssertNil(t, err, "Should not receive error")
		value := v.(int64)

		if value == int64(10) {
			resultIntervalOne = append(resultIntervalOne, value)
		} else if value == int64(20) {
			resultIntervalTwo = append(resultIntervalTwo, value)
		} else {
			t.Errorf("Should not have generated a value outside of the domain!")
		}
	}

	Assert(t, len(resultIntervalOne) > 0, "expected to generate at least one")
	Assert(t, len(resultIntervalTwo) > 0, "expected to generate at least one")
}

func TestWeightedType(t *testing.T) {
	w := &WeightDistribution{}
	AssertEqual(t, WEIGHT_DIST, w.Type())
}

func TestNormalType(t *testing.T) {
	w := &NormalDistribution{}
	AssertEqual(t, NORMAL_DIST, w.Type())
}
