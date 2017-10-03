package generator

import (
	"fmt"
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/emitter"
	"math/rand"
	"time"
)

func NewDistribution(distType string, weights []float64, fields []FieldType) (FieldType, error) {
	switch distType {
	case NORMAL_DIST:
		min, max := weights[0], weights[1]

		if max < min {
			return nil, fmt.Errorf("max cannot be less than min")
		}

		return &NormalDistribution{min: min, max: max}, nil
	case WEIGHT_DIST, PERCENT_DIST:
		total := float64(0)

		for _, w := range weights {
			if w < 0 {
				return nil, fmt.Errorf("weights cannot be negative: %f", w)
			}
			total += w
		}

		if distType == PERCENT_DIST && total != float64(1) {
			return nil, fmt.Errorf("percentage weights do not add to 100%% (i.e. 1.0). total = %f", total)
		}

		return &WeightDistribution{weights: weights, intervals: fields, picker: func() float64 { return rand.Float64() * (total) }}, nil
	default:
		return nil, fmt.Errorf("Unsupported distribution %q", distType)
	}
}

type roulette = func() float64

type WeightDistribution struct {
	weights   []float64
	intervals []FieldType
	picker    roulette
}

func (dist *WeightDistribution) One(parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	if len(dist.intervals) == 1 {
		return dist.intervals[0].One(parentId, emitter, scope)
	}

	rand.Seed(time.Now().UnixNano())

	n := dist.picker()
	for i := 0; i < len(dist.intervals); i++ {
		if n < dist.weights[i] {
			return dist.intervals[i].One(parentId, emitter, scope)
		}
		n -= dist.weights[i]
	}
	return nil, nil
}

func (dist *WeightDistribution) Type() string { return WEIGHT_DIST }

type NormalDistribution struct {
	min, max float64
}

func (dist *NormalDistribution) One(parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	min, max := dist.min, dist.max
	rand.Seed(time.Now().UnixNano())
	mean := calcMean(min, max)
	stdDev := calcMean(mean, max)

	result := rand.NormFloat64()*stdDev + mean

	//Need this check because it's possible the result will be
	// 0.9999999999999 smaller/bigger than the min/max
	if result < min || result > max {
		return dist.One(parentId, emitter, scope)
	} else {
		return result, nil
	}
}

func calcMean(min, max float64) float64 {
	return (max + min) / 2.0
}

func (dist *NormalDistribution) Type() string { return NORMAL_DIST }
