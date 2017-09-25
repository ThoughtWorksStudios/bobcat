package generator

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/emitter"
	"math/rand"
	"time"
)

type Domain struct {
	intervals []FieldType
}

//possible TODO: Refactor the Generator to be an interface{}, and create two implementing types
// 1st: StdGenerator (or something) which is pretty much what Generator currently is
// 2nd: DistributionGenerator (aka DistGen)
//
// By adding DistGen we'd possibly/might gain the following:
//   * The interpreter could be refactored to pass the DistGen to Visit, withXField which would reduce  code duplication
//   * The interpreter interface would be a little more standardized
type Distribution interface {
	One(domain Domain, parentId interface{}, emitter Emitter, scope *Scope) interface{}
	OneFromMultipleIntervals(intervals []FieldType, parentId interface{}, emitter Emitter, scope *Scope) interface{}
	OneFromSingleInterval(interval FieldType, parentId interface{}, emitter Emitter, scope *Scope) interface{}
	isCompatibleDomain(domain string) bool
	supportsMultipleIntervals() bool
	Type() string
}

type WeightedDistribution struct {
	weights []float64
}

func (dist *WeightedDistribution) One(domain Domain, parentId interface{}, emitter Emitter, scope *Scope) interface{} {
	if len(domain.intervals) == 1 {
		return dist.OneFromSingleInterval(domain.intervals[0], parentId, emitter, scope)
	} else {
		return dist.OneFromMultipleIntervals(domain.intervals, parentId, emitter, scope)
	}
}

func (dist *WeightedDistribution) sumOfWeights() float64 {
	var result float64
	for i := 0; i < len(dist.weights); i++ {
		result += dist.weights[i]
	}
	return result
}

func (dist *WeightedDistribution) OneFromMultipleIntervals(intervals []FieldType, parentId interface{}, emitter Emitter, scope *Scope) interface{} {
	rand.Seed(time.Now().UnixNano())
	n := (&FloatType{min: 0.0, max: dist.sumOfWeights()}).One(parentId, emitter, nil, scope).(float64)
	for i := 0; i < len(intervals); i++ {
		if n < dist.weights[i] {
			return dist.OneFromSingleInterval(intervals[i], parentId, emitter, scope)
		}
		n -= dist.weights[i]
	}
	return nil
}

func (dist *WeightedDistribution) OneFromSingleInterval(interval FieldType, parentId interface{}, emitter Emitter, scope *Scope) interface{} {
	return interval.One(parentId, emitter, nil, scope)
}

func (dist *WeightedDistribution) isCompatibleDomain(domain string) bool {
	return true
}

func (dist *WeightedDistribution) supportsMultipleIntervals() bool {
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

func (dist *PercentageDistribution) One(domain Domain, parentId interface{}, emitter Emitter, scope *Scope) interface{} {
	if len(domain.intervals) == 1 {
		return dist.OneFromSingleInterval(domain.intervals[0], parentId, emitter, scope)
	} else {
		return dist.OneFromMultipleIntervals(domain.intervals, parentId, emitter, scope)
	}
}

func (dist *PercentageDistribution) OneFromMultipleIntervals(intervals []FieldType, parentId interface{}, emitter Emitter, scope *Scope) interface{} {
	for i := 0; i < len(intervals); i++ {
		if dist.bins[i] == 0 || dist.weights[i] >= (float64(dist.bins[i])/float64(dist.total)*100.0) {
			dist.bins[i] = dist.bins[i] + 1
			dist.total++
			return dist.OneFromSingleInterval(intervals[i], parentId, emitter, scope)
		}
	}
	return nil
}

func (dist *PercentageDistribution) OneFromSingleInterval(interval FieldType, parentId interface{}, emitter Emitter, scope *Scope) interface{} {
	return interval.One(parentId, emitter, nil, scope)
}

func (dist *PercentageDistribution) isCompatibleDomain(domain string) bool {
	return true
}

func (dist *PercentageDistribution) supportsMultipleIntervals() bool {
	return true
}

func (dist *PercentageDistribution) Type() string {
	return "percentage"
}

type NormalDistribution struct{}

func (dist *NormalDistribution) One(domain Domain, parentId interface{}, emitter Emitter, scope *Scope) interface{} {
	return dist.OneFromSingleInterval(domain.intervals[0], parentId, emitter, scope)
}

func (dist *NormalDistribution) calcMean(min, max float64) float64 {
	return (max + min) / 2.0
}

func (dist *NormalDistribution) OneFromSingleInterval(interval FieldType, parentId interface{}, emitter Emitter, scope *Scope) interface{} {
	floatInterval := interval.(*FloatType)
	min, max := floatInterval.min, floatInterval.max
	rand.Seed(time.Now().UnixNano())
	mean := dist.calcMean(min, max)
	stdDev := dist.calcMean(mean, max)

	result := rand.NormFloat64()*stdDev + mean

	//Need this check because it's possible the result will be
	// 0.9999999999999 smaller/bigger than the min/max
	if result < min || result > max {
		return dist.OneFromSingleInterval(interval, parentId, emitter, scope)
	} else {
		return result
	}
}

func (dist *NormalDistribution) supportsMultipleIntervals() bool {
	return false
}

func (dist *NormalDistribution) isCompatibleDomain(domain string) bool {
	return domain == FLOAT_TYPE
}

func (dist *NormalDistribution) Type() string {
	return "normal"
}

func (dist *NormalDistribution) OneFromMultipleIntervals(intervals []FieldType, parentId interface{}, emitter Emitter, scope *Scope) interface{} {
	return nil
}

type UniformDistribution struct{}

func (dist *UniformDistribution) OneFromMultipleIntervals(intervals []FieldType, parentId interface{}, emitter Emitter, scope *Scope) interface{} {
	return nil
}

func (dist *UniformDistribution) OneFromSingleInterval(interval FieldType, parentId interface{}, emitter Emitter, scope *Scope) interface{} {
	return interval.One(parentId, emitter, nil, scope)
}

func (dist *UniformDistribution) One(domain Domain, parentId interface{}, emitter Emitter, scope *Scope) interface{} {
	return dist.OneFromSingleInterval(domain.intervals[0], parentId, emitter, scope)
}

func (dist *UniformDistribution) isCompatibleDomain(domain string) bool {
	return domain == INT_TYPE || domain == FLOAT_TYPE
}

func (dist *UniformDistribution) Type() string {
	return "normal"
}

func (dist *UniformDistribution) supportsMultipleIntervals() bool {
	return false
}
