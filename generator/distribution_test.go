package generator

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"math"
	"testing"
)

// NOTE: these tests inherently have the potential for flakiness, as they rely on
// randomly generated values to approach a certain shape. We use a reasonably large
// sample size (that, at the same time, isn't too slow) and make judicious use of
// rounding to mitigate this. It appears to do a decent job, but time will tell.

func TestWeightDistribution(t *testing.T) {
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

func TestNormalDistribution(t *testing.T) {
	bounds := []float64{0, 100}

	dist, err := NewDistribution(NORMAL_DIST, bounds, nil)
	AssertNil(t, err, "Should not receive error during distribution construction")

	count := 10000
	values := make([]float64, count)

	sum := 0.0

	for i := 0; i < count; i++ {
		v, err := dist.One(nil, nil, nil)
		AssertNil(t, err, "Should not receive error")
		n := v.(float64)
		sum += n
		values[i] = n
	}

	mean := sum / float64(count)
	sd := stdDev(mean, values)

	expectedMean := 50.0

	// actually, perfectly uniform distribution yields closer to 29.16, so either math/rand or the distro
	// algorithm may be slightly off. or perhaps we'll approach 29.16 with a larger count, or maybe not.
	// close enough though.
	expectedStandDev := 28.0

	Assert(t, withinTolerance(expectedMean, RoundFloat(mean, 1), 1), "Mean should be roughly %f", expectedMean)
	Assert(t, withinTolerance(expectedStandDev, RoundFloat(sd, 1), 1), "Standard deviation should be roughly %f", expectedStandDev)
}

func TestWeightedType(t *testing.T) {
	w := &WeightDistribution{}
	AssertEqual(t, WEIGHT_DIST, w.Type())
}

func TestNormalType(t *testing.T) {
	w := &NormalDistribution{}
	AssertEqual(t, NORMAL_DIST, w.Type())
}

func withinTolerance(expected, actual, tolerance float64) bool {
	return expected == actual || (actual <= (expected+tolerance) && actual >= (expected-tolerance))
}

func stdDev(mean float64, values []float64) float64 {
	count, variance := len(values), 0.0

	for _, val := range values {
		variance += math.Pow(val-mean, 2)
	}

	variance /= float64(count)

	return math.Sqrt(variance)
}
