package generator

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"testing"
)

// TODO: Flaky??
func TestWeightDistributionOneEnum(t *testing.T) {
	weights := []float64{80.0, 20.0}
	intervalOne := &EnumType{size: 2, values: []interface{}{"one", "two"}}
	intervalTwo := &EnumType{size: 2, values: []interface{}{"three", "four"}}
	domain := Domain{intervals: []FieldType{intervalOne, intervalTwo}}
	dist := &WeightDistribution{weights: weights}

	count := 10000

	resultIntervalOne := []interface{}{}
	resultIntervalTwo := []interface{}{}

	for i := 0; i < count; i++ {
		v, err := dist.One(domain, nil, nil, nil)
		AssertNil(t, err, "Should not receive error")
		value := v.(string)

		if value == "one" || value == "two" {
			resultIntervalOne = append(resultIntervalOne, v)
		} else if value == "three" || value == "four" {
			resultIntervalTwo = append(resultIntervalTwo, v)
		} else {
			t.Errorf("Should not have generated a value outside of the domain!")
		}
	}

	AssertEqual(t, 8, int(RoundFloat(float64(len(resultIntervalOne))/1000.0, 1)), "Interval 1 should be approximately 80%")
	AssertEqual(t, 2, int(RoundFloat(float64(len(resultIntervalTwo))/1000.0, 1)), "Interval 2 should be approximately 80%")
}

func TestWeightDistributionOne(t *testing.T) {
	weights := []float64{50.0, 50.0}
	intervalOne := &IntegerType{min: 1, max: 10}
	intervalTwo := &IntegerType{min: 20, max: 30}
	domain := Domain{intervals: []FieldType{intervalOne, intervalTwo}}
	dist := &WeightDistribution{weights: weights}

	count := 10

	resultIntervalOne := []interface{}{}
	resultIntervalTwo := []interface{}{}

	for i := 0; i < count; i++ {
		v, err := dist.One(domain, nil, nil, nil)
		AssertNil(t, err, "Should not receive error")
		value := v.(int64)

		if value >= intervalOne.min && value <= intervalOne.max {
			resultIntervalOne = append(resultIntervalOne, v)
		} else if value >= intervalTwo.min && value <= intervalTwo.max {
			resultIntervalTwo = append(resultIntervalTwo, v)
		} else {
			t.Errorf("Should not have generated a value outside of the domain!")
		}
	}

	Assert(t, len(resultIntervalOne) > 0, "expected to generate at least one")
	Assert(t, len(resultIntervalTwo) > 0, "expected to generate at least one")
}

func TestNormalCompatibleDomain(t *testing.T) {
	norm := &NormalDistribution{}
	Assert(t, norm.isCompatibleDomain(FLOAT_TYPE), "floats should be a compatible domain for normal distributions")
	Assert(t, !norm.isCompatibleDomain(INT_TYPE), "ints should not be a compatible domain for normal distributions")
}

func TestUniformCompatibleDomain(t *testing.T) {
	uni := &UniformDistribution{}
	Assert(t, uni.isCompatibleDomain(FLOAT_TYPE), "floats should be a compatible domain for uniform distributions")
	Assert(t, uni.isCompatibleDomain(INT_TYPE), "ints should be a compatible domain for uniform distributions")
	Assert(t, !uni.isCompatibleDomain(STRING_TYPE), "strings should not be a compatible domain for uniform distributions")
}

func TestNormalShouldntSupportMultipleIntervals(t *testing.T) {
	norm := &NormalDistribution{}
	Assert(t, !norm.supportsMultipleIntervals(), "normal distributions don't support multiple domains")
}

func TestUniformShouldntSupportMultipleIntervals(t *testing.T) {
	uni := &UniformDistribution{}
	Assert(t, !uni.supportsMultipleIntervals(), "uniform distributions don't support multiple domains")
}

func TestWeightedShouldSupportMultipleIntervals(t *testing.T) {
	w := &WeightDistribution{}
	Assert(t, w.supportsMultipleIntervals(), "weight distributions should support multiple domains")
}

func TestWeightedType(t *testing.T) {
	w := &WeightDistribution{}
	AssertEqual(t, WEIGHT_DIST, w.Type())
}

func TestNormalType(t *testing.T) {
	w := &NormalDistribution{}
	AssertEqual(t, NORMAL_DIST, w.Type())
}

func TestUniformType(t *testing.T) {
	w := &UniformDistribution{}
	AssertEqual(t, UNIFORM_DIST, w.Type())
}
