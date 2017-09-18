package generator

type Domain struct {
	intervals []Interval
}

type Interval interface {
	Type() string
	contains(value interface{}) bool
}

type IntegerInterval struct {
	min int64
	max int64
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
