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
	One(domain Domain, parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error)
	OneFromMultipleIntervals(intervals []FieldType, parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error)
	OneFromSingleInterval(interval FieldType, parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error)
	isCompatibleDomain(domain string) bool
	supportsMultipleIntervals() bool
	Type() string
}

type WeightDistribution struct {
	weights []float64
}

func (dist *WeightDistribution) One(domain Domain, parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	if len(domain.intervals) == 1 {
		return dist.OneFromSingleInterval(domain.intervals[0], parentId, emitter, scope)
	} else {
		return dist.OneFromMultipleIntervals(domain.intervals, parentId, emitter, scope)
	}
}

func (dist *WeightDistribution) sumOfWeights() float64 {
	var result float64
	for i := 0; i < len(dist.weights); i++ {
		result += dist.weights[i]
	}
	return result
}

func (dist *WeightDistribution) OneFromMultipleIntervals(intervals []FieldType, parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	rand.Seed(time.Now().UnixNano())
	if val, err := (&FloatType{min: 0.0, max: dist.sumOfWeights()}).One(parentId, emitter, scope); err == nil {
		n := val.(float64)
		for i := 0; i < len(intervals); i++ {
			if n < dist.weights[i] {
				return dist.OneFromSingleInterval(intervals[i], parentId, emitter, scope)
			}
			n -= dist.weights[i]
		}
		return nil, nil
	} else {
		return nil, err
	}
}

func (dist *WeightDistribution) OneFromSingleInterval(interval FieldType, parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	return interval.One(parentId, emitter, scope)
}

func (dist *WeightDistribution) isCompatibleDomain(domain string) bool { return true }

func (dist *WeightDistribution) supportsMultipleIntervals() bool { return true }

func (dist *WeightDistribution) Type() string { return WEIGHT_DIST }

type NormalDistribution struct{}

func (dist *NormalDistribution) One(domain Domain, parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	return dist.OneFromSingleInterval(domain.intervals[0], parentId, emitter, scope)
}

func (dist *NormalDistribution) calcMean(min, max float64) float64 {
	return (max + min) / 2.0
}

func (dist *NormalDistribution) OneFromSingleInterval(interval FieldType, parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
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
		return result, nil
	}
}

func (dist *NormalDistribution) supportsMultipleIntervals() bool { return false }

func (dist *NormalDistribution) isCompatibleDomain(domain string) bool {
	return domain == FLOAT_TYPE
}

func (dist *NormalDistribution) Type() string { return NORMAL_DIST }

func (dist *NormalDistribution) OneFromMultipleIntervals(intervals []FieldType, parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	return nil, nil
}
