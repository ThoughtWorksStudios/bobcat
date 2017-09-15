package generator

import (
	"github.com/leesper/go_rng"
	// "math"
	"math/rand"
	"time"
)

type Distribution interface {
	One(min, max interface{}) interface{}
	isCompatibleDomain(domain string) bool
	Type() string
}

type WeightedDistribution struct {
	domainType string
}

func (dist *WeightedDistribution) One(min, max interface{}) interface{} {
	return nil
}

func (dist *WeightedDistribution) isCompatibleDomain(domain string) bool {
	return true
}

func (dist *WeightedDistribution) Type() string {
	return "weighted"
}

type NormalDistribution struct {
	domainType string
}

func (dist *NormalDistribution) calcMean(min, max float64) float64 {
	return (max + min) / 2.0
}

func (dist *NormalDistribution) One(min, max interface{}) interface{} {
	rand.Seed(time.Now().UnixNano())
	floor := min.(float64)
	ceiling := max.(float64)
	mean := dist.calcMean(floor, ceiling)
	stdDev := dist.calcMean(mean, ceiling)

	result := rand.NormFloat64()*stdDev + mean

	//Need this check because it's possible the result will be
	// 0.9999999999999 smaller/bigger than the min/max
	if result < floor || result > ceiling {
		return dist.One(floor, ceiling)
	} else {
		return result
	}
}

func (dist *NormalDistribution) isCompatibleDomain(domain string) bool {
	return domain == "float"
}

func (dist *NormalDistribution) Type() string {
	return "normal"
}

type UniformDistribution struct {
	domainType string
}

func (dist *UniformDistribution) One(min, max interface{}) interface{} {
	uniform := rng.NewUniformGenerator(time.Now().UnixNano())
	switch dist.domainType {
	case "integer":
		return uniform.Int64Range(min.(int64), max.(int64))
	case "float":
		return uniform.Float64Range(min.(float64), max.(float64))
	default:
		return nil
	}
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
