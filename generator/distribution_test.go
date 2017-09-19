package generator

import (
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"testing"
	"time"
)

func TestPercentageDistributionOneInteger(t *testing.T) {
	weights := []float64{50.0, 50.0}
	intervalOne := &IntegerType{min: 1, max: 10}
	intervalTwo := &IntegerType{min: 20, max: 30}
	domain := Domain{intervals: []FieldType{intervalOne, intervalTwo}}
	dist := &PercentageDistribution{weights: weights, bins: make([]int64, len(weights))}

	count := 10

	resultIntervalOne := []interface{}{}
	resultIntervalTwo := []interface{}{}

	for i := 0; i < count; i++ {
		v := dist.One(domain)

		value := v.(int64)

		if value >= intervalOne.min && value <= intervalOne.max {
			resultIntervalOne = append(resultIntervalOne, v)
		} else if value >= intervalTwo.min && value <= intervalTwo.max {
			resultIntervalTwo = append(resultIntervalTwo, v)
		} else {
			t.Errorf("Should not have generated a value outside of the domain!")
		}
	}

	AssertEqual(t, len(resultIntervalOne), 5)
	AssertEqual(t, len(resultIntervalTwo), 5)
}

func TestPercentageDistributionOneLiteralField(t *testing.T) {
	weights := []float64{50.0, 50.0}
	intervalOne := &LiteralType{value: "blah"}
	intervalTwo := &LiteralType{value: "eek"}
	domain := Domain{intervals: []FieldType{intervalOne, intervalTwo}}
	dist := &PercentageDistribution{weights: weights, bins: make([]int64, len(weights))}

	count := 10

	resultIntervalOne := []interface{}{}
	resultIntervalTwo := []interface{}{}

	for i := 0; i < count; i++ {
		v := dist.One(domain)

		value := v.(string)

		if value == "blah" {
			resultIntervalOne = append(resultIntervalOne, v)
		} else if value == "eek" {
			resultIntervalTwo = append(resultIntervalTwo, v)
		} else {
			t.Errorf("Should not have generated a value outside of the domain!")
		}
	}

	AssertEqual(t, len(resultIntervalOne), 5)
	AssertEqual(t, len(resultIntervalTwo), 5)
}

func TestWeightedDistributionOneEnum(t *testing.T) {
	weights := []float64{60.0, 40.0}
	intervalOne := &EnumType{size: 2, values: []interface{}{"one", "two"}}
	intervalTwo := &EnumType{size: 2, values: []interface{}{"three", "four"}}
	domain := Domain{intervals: []FieldType{intervalOne, intervalTwo}}
	dist := &PercentageDistribution{weights: weights, bins: make([]int64, len(weights))}

	count := 10

	resultIntervalOne := []interface{}{}
	resultIntervalTwo := []interface{}{}

	for i := 0; i < count; i++ {
		v := dist.One(domain)
		value := v.(string)

		if value == "one" || value == "two" {
			resultIntervalOne = append(resultIntervalOne, v)
		} else if value == "three" || value == "four" {
			resultIntervalTwo = append(resultIntervalTwo, v)
		} else {
			t.Errorf("Should not have generated a value outside of the domain!")
		}
	}

	AssertEqual(t, len(resultIntervalOne), 6)
	AssertEqual(t, len(resultIntervalTwo), 4)
}

func TestPercentageDistributionOneDate(t *testing.T) {
	weights := []float64{50.0, 50.0}
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "1945-01-02")
	timeMax2, _ := time.Parse("2006-01-02", "1950-01-02")
	intervalOne := &DateType{min: timeMin, max: timeMax}
	intervalTwo := &DateType{min: timeMax, max: timeMax2}
	domain := Domain{intervals: []FieldType{intervalOne, intervalTwo}}
	dist := &PercentageDistribution{weights: weights, bins: make([]int64, len(weights))}

	count := 10

	resultIntervalOne := []interface{}{}
	resultIntervalTwo := []interface{}{}

	for i := 0; i < count; i++ {
		v := dist.One(domain)

		value := v.(*TimeWithFormat).Time

		if value.After(intervalOne.min) && value.Before(intervalOne.max) {
			resultIntervalOne = append(resultIntervalOne, v)
		} else if value.After(intervalTwo.min) && value.Before(intervalTwo.max) {
			resultIntervalTwo = append(resultIntervalTwo, v)
		} else {
			t.Errorf("Should not have generated a value outside of the domain! %v\n", value)
		}
	}

	AssertEqual(t, len(resultIntervalOne), 5)
	AssertEqual(t, len(resultIntervalTwo), 5)
}

func TestWeightedDistributionOne(t *testing.T) {
	weights := []float64{50.0, 50.0}
	intervalOne := &IntegerType{min: 1, max: 10}
	intervalTwo := &IntegerType{min: 20, max: 30}
	domain := Domain{intervals: []FieldType{intervalOne, intervalTwo}}
	dist := &WeightedDistribution{weights: weights}

	count := 10

	resultIntervalOne := []interface{}{}
	resultIntervalTwo := []interface{}{}

	for i := 0; i < count; i++ {
		v := dist.One(domain)
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
	Assert(t, norm.isCompatibleDomain("decimal"), "decmals should be a compatible domain for normal distributions")
	Assert(t, !norm.isCompatibleDomain("integer"), "integers should not be a compatible domain for normal distributions")
}

func TestUniformCompatibleDomain(t *testing.T) {
	uni := &UniformDistribution{}
	Assert(t, uni.isCompatibleDomain("decimal"), "decimals should be a compatible domain for uniform distributions")
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

func TestPercentageShouldSupportMultipleDomains(t *testing.T) {
	w := &PercentageDistribution{}
	Assert(t, w.supportsMultipleDomains(), "percent distributions should support multiple domains")
}

func TestWeightedShouldSupportMultipleDomains(t *testing.T) {
	w := &WeightedDistribution{}
	Assert(t, w.supportsMultipleDomains(), "weighted distributions should support multiple domains")
}
