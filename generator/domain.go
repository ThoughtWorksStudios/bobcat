package generator

import "math/rand"

type Domain struct {
	intervals []Interval
}

type Interval interface {
	Type() string
	contains(value interface{}) bool
	One() interface{}
}

type IntegerInterval struct {
	min int64
	max int64
}

func (i IntegerInterval) One() interface{} {
	return i.min + rand.Int63n(i.max-i.min+1)
}

func (i IntegerInterval) contains(value interface{}) bool {
	if v, ok := value.(int64); ok {
		if v >= i.min && v <= i.max {
			return true
		}
	}
	return false
}

func (i IntegerInterval) Type() string {
	return "integer"
}

type FloatInterval struct {
	min float64
	max float64
}

func (i FloatInterval) One() interface{} {
	return rand.Float64()*(i.max-i.min) + i.min
}

func (i FloatInterval) contains(value interface{}) bool {
	if v, ok := value.(float64); ok {
		if v >= i.min && v <= i.max {
			return true
		}
	}
	return false
}

func (i FloatInterval) Type() string {
	return "integer"
}
