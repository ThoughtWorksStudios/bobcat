package generator

import (
	"math/rand"
	"time"
)

type Domain struct {
	intervals []Interval
}

type Interval interface {
	Type() string
}

type IntegerInterval struct {
	min int64
	max int64
}

func (i IntegerInterval) Type() string {
	return "integer"
}

type FloatInterval struct {
	min float64
	max float64
}

func (i FloatInterval) Type() string {
	return "integer"
}

type Distribution interface {
	One(domain Domain) interface{}
	OneFromMultipleIntervals(intervals []Interval) interface{}
	OneFromSingleInterval(interval Interval) interface{}
	isCompatibleDomain(domain string) bool
	supportsMultipleDomains() bool
	Type() string
}

type WeightedDistribution struct {
	weights []float64
}

func (dist *WeightedDistribution) One(domain Domain) interface{} {
	if len(domain.intervals) == 1 {
		return (&UniformDistribution{}).One(domain)
	} else {
		dist.OneFromMultipleIntervals(domain.intervals)
	}
	return nil
}

func (dist *WeightedDistribution) OneFromMultipleIntervals(intervals []Interval) interface{} {
	return nil
}

func (dist *WeightedDistribution) OneFromSingleInterval(interval Interval) interface{} {
	return nil
}

func (dist *WeightedDistribution) isCompatibleDomain(domain string) bool {
	return true
}

func (dist *WeightedDistribution) supportsMultipleDomains() bool {
	return true
}

func (dist *WeightedDistribution) Type() string {
	return "weighted"
}

type NormalDistribution struct{}

func (dist *NormalDistribution) One(domain Domain) interface{} {
	return dist.OneFromSingleInterval(domain.intervals[0])
}

func (dist *NormalDistribution) calcMean(min, max float64) float64 {
	return (max + min) / 2.0
}

func (dist *NormalDistribution) OneFromSingleInterval(interval Interval) interface{} {
	floatInterval := interval.(FloatInterval)
	min, max := floatInterval.min, floatInterval.max
	rand.Seed(time.Now().UnixNano())
	mean := dist.calcMean(min, max)
	stdDev := dist.calcMean(mean, max)

	result := rand.NormFloat64()*stdDev + mean

	//Need this check because it's possible the result will be
	// 0.9999999999999 smaller/bigger than the min/max
	if result < min || result > max {
		return dist.OneFromSingleInterval(interval)
	} else {
		return result
	}
}

func (dist *NormalDistribution) supportsMultipleDomains() bool {
	return false
}

func (dist *NormalDistribution) isCompatibleDomain(domain string) bool {
	return domain == "float"
}

func (dist *NormalDistribution) Type() string {
	return "normal"
}

func (dist *NormalDistribution) OneFromMultipleIntervals(intervals []Interval) interface{} {
	return nil
}

type UniformDistribution struct{}

func (dist *UniformDistribution) OneFromMultipleIntervals(intervals []Interval) interface{} {
	return nil
}

func (dist *UniformDistribution) OneFromSingleInterval(interval Interval) interface{} {
	switch interval.Type() {
	case "integer":
		intInterval := interval.(IntegerInterval)
		return intInterval.min + rand.Int63n(intInterval.max-intInterval.min+1)
	case "float":
		floatInterval := interval.(FloatInterval)
		return rand.Float64()*(floatInterval.max-floatInterval.min) + floatInterval.min
	default:
		return nil
	}
}

func (dist *UniformDistribution) One(domain Domain) interface{} {
	return dist.OneFromSingleInterval(domain.intervals[0])
}

func (dist *UniformDistribution) isCompatibleDomain(domain string) bool {
	switch domain {
	case "integer":
		return true
	case "float":
		return true
	default:
		return false
	}
}

func (dist *UniformDistribution) Type() string {
	return "normal"
}

func (dist *UniformDistribution) supportsMultipleDomains() bool {
	return false
}
