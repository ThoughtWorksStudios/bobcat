package generator

import (
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"testing"
)

func TestPercentageDistributionOne(t *testing.T) {
	weights := []float64{50.0, 50.0}
	intervalOne := IntegerInterval{min: 1, max: 10}
	intervalTwo := IntegerInterval{min: 20, max: 30}
	domain := Domain{intervals: []Interval{intervalOne, intervalTwo}}
	dist := &PercentageDistribution{weights: weights, bins: make([]int64, len(weights))}

	count := 10

	resultIntervalOne := []interface{}{}
	resultIntervalTwo := []interface{}{}

	for i := 0; i < count; i++ {
		v := dist.One(domain)

		if intervalOne.contains(v) {
			resultIntervalOne = append(resultIntervalOne, v)
		} else if intervalTwo.contains(v) {
			resultIntervalTwo = append(resultIntervalTwo, v)
		} else {
			t.Errorf("Should not have generated a value outside of the domain!")
		}
	}

	AssertEqual(t, len(resultIntervalOne), 5)
	AssertEqual(t, len(resultIntervalTwo), 5)
}

func TestWeightedDistributionOne(t *testing.T) {
	weights := []float64{50.0, 50.0}
	intervalOne := IntegerInterval{min: 1, max: 10}
	intervalTwo := IntegerInterval{min: 20, max: 30}
	domain := Domain{intervals: []Interval{intervalOne, intervalTwo}}
	dist := &WeightedDistribution{weights: weights}

	count := 10

	resultIntervalOne := []interface{}{}
	resultIntervalTwo := []interface{}{}

	for i := 0; i < count; i++ {
		v := dist.One(domain)

		if intervalOne.contains(v) {
			resultIntervalOne = append(resultIntervalOne, v)
		} else if intervalTwo.contains(v) {
			resultIntervalTwo = append(resultIntervalTwo, v)
		} else {
			t.Errorf("Should not have generated a value outside of the domain!")
		}
	}

	AssertEqual(t, len(resultIntervalOne), 5)
	AssertEqual(t, len(resultIntervalTwo), 5)
}

func TestNormalCompatibleDomain(t *testing.T) {
	norm := &NormalDistribution{}
	Assert(t, norm.isCompatibleDomain("float"), "floats should be a compatible domain for normal distributions")
	Assert(t, !norm.isCompatibleDomain("integer"), "integers should not be a compatible domain for normal distributions")
}

func TestUniformCompatibleDomain(t *testing.T) {
	uni := &UniformDistribution{}
	Assert(t, uni.isCompatibleDomain("float"), "floats should be a compatible domain for uniform distributions")
	Assert(t, uni.isCompatibleDomain("integer"), "integers should be a compatible domain for uniform distributions")
	Assert(t, !uni.isCompatibleDomain("string"), "strings should not be a compatible domain for uniform distributions")
}

func TestNormalShouldntSupportMultipleDomains(t *testing.T) {
	norm := &NormalDistribution{}
	Assert(t, !norm.supportsMultipleDomains(), "normal distributions don't support multiple domains")
}

func TestUniformShouldntSupportMultipleDomains(t *testing.T) {
	uni := &UniformDistribution{}
	Assert(t, !uni.supportsMultipleDomains(), "uniform distributions don't support multiple domains")
}

func TestWeightedShouldSupportMultipleDomains(t *testing.T) {
	w := &WeightedDistribution{}
	Assert(t, w.supportsMultipleDomains(), "weighted distributions should support multiple domains")
}
