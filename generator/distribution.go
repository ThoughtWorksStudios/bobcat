package generator

import (
	"math/rand"
	"time"
)

type Domain struct {
	intervals []FieldType
}

type Distribution interface {
	One(domain Domain) interface{}
	OneFromMultipleIntervals(intervals []FieldType) interface{}
	OneFromSingleInterval(interval FieldType) interface{}
	isCompatibleDomain(domain string) bool
	supportsMultipleDomains() bool
	Type() string
}

type WeightedDistribution struct {
	weights []float64
}

func (dist *WeightedDistribution) One(domain Domain) interface{} {
	if len(domain.intervals) == 1 {
		return dist.OneFromSingleInterval(domain.intervals[0])
	} else {
		return dist.OneFromMultipleIntervals(domain.intervals)
	}
}

func (dist *WeightedDistribution) sumOfWeights() float64 {
	var result float64
	for i := 0; i < len(dist.weights); i++ {
		result += dist.weights[i]
	}
	return result
}

func (dist *WeightedDistribution) OneFromMultipleIntervals(intervals []FieldType) interface{} {
	rand.Seed(time.Now().UnixNano())
	n := (&FloatType{min: 0.0, max: dist.sumOfWeights()}).One(nil, nil, nil, nil).(float64)
	for i := 0; i < len(intervals); i++ {
		if n < dist.weights[i] {
			return dist.OneFromSingleInterval(intervals[i])
		}
		n -= dist.weights[i]
	}
	return nil
}

func (dist *WeightedDistribution) OneFromSingleInterval(interval FieldType) interface{} {
	return interval.One(nil, nil, nil, nil)
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

type PercentageDistribution struct {
	weights []float64 //the percent associated with intervals[i]
	bins    []int64   //The number of generated values for interval[i]
	total   int64     // the totaly number of values generated
}

func (dist *PercentageDistribution) One(domain Domain) interface{} {
	if len(domain.intervals) == 1 {
		return dist.OneFromSingleInterval(domain.intervals[0])
	} else {
		return dist.OneFromMultipleIntervals(domain.intervals)
	}
}

func (dist *PercentageDistribution) OneFromMultipleIntervals(intervals []FieldType) interface{} {
	for i := 0; i < len(intervals); i++ {
		if dist.bins[i] == 0 || dist.weights[i] >= (float64(dist.bins[i])/float64(dist.total)*100.0) {
			dist.bins[i] = dist.bins[i] + 1
			dist.total++
			return dist.OneFromSingleInterval(intervals[i])
		}
	}
	return nil
}

func (dist *PercentageDistribution) OneFromSingleInterval(interval FieldType) interface{} {
	return interval.One(nil, nil, nil, nil)
}

func (dist *PercentageDistribution) isCompatibleDomain(domain string) bool {
	return true
}

func (dist *PercentageDistribution) supportsMultipleDomains() bool {
	return true
}

func (dist *PercentageDistribution) Type() string {
	return "percentage"
}

type NormalDistribution struct{}

func (dist *NormalDistribution) One(domain Domain) interface{} {
	return dist.OneFromSingleInterval(domain.intervals[0])
}

func (dist *NormalDistribution) calcMean(min, max float64) float64 {
	return (max + min) / 2.0
}

func (dist *NormalDistribution) OneFromSingleInterval(interval FieldType) interface{} {
	floatInterval := interval.(*FloatType)
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

func (dist *NormalDistribution) OneFromMultipleIntervals(intervals []FieldType) interface{} {
	return nil
}

type UniformDistribution struct{}

func (dist *UniformDistribution) OneFromMultipleIntervals(intervals []FieldType) interface{} {
	return nil
}

func (dist *UniformDistribution) OneFromSingleInterval(interval FieldType) interface{} {
	return interval.One(nil, nil, nil, nil)
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
