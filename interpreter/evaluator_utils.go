package interpreter

import (
	"fmt"
	"strconv"
	"strings"
)

func addToInt(lhs int64, right interface{}, op string) (interface{}, error) {
	switch right.(type) {
	case int64:
		rhs := right.(int64)
		if "-" == op {
			return lhs - rhs, nil
		}
		return lhs + rhs, nil
	case float64:
		return addToFloat(float64(lhs), right, op)
	case string:
		if "-" == op {
			return nil, incompatible(op)
		}
		return (strconv.FormatInt(lhs, 10) + right.(string)), nil
	default:
		return nil, incompatible(op)
	}
}

func addToFloat(lhs float64, right interface{}, op string) (interface{}, error) {
	switch right.(type) {
	case int64:
		rhs := float64(right.(int64))
		if "-" == op {
			return lhs - rhs, nil
		}
		return lhs + rhs, nil
	case float64:
		rhs := right.(float64)
		if "-" == op {
			return lhs - rhs, nil
		}
		return lhs + rhs, nil
	case string:
		if "-" == op {
			return nil, incompatible(op)
		}
		return (strconv.FormatFloat(lhs, 'f', -1, 64) + right.(string)), nil
	default:
		return nil, incompatible(op)
	}
}

func addToString(lhs string, right interface{}, op string) (interface{}, error) {
	if "-" == op {
		return nil, incompatible(op)
	}
	switch right.(type) {
	case string:
		return (lhs + right.(string)), nil
	case int64:
		return (lhs + strconv.FormatInt(right.(int64), 10)), nil
	case float64:
		return (lhs + strconv.FormatFloat(right.(float64), 'f', -1, 64)), nil
	case bool:
		return (lhs + strconv.FormatBool(right.(bool))), nil
	default:
		return nil, incompatible(op)
	}
}

func addToBool(lhs bool, right interface{}, op string) (interface{}, error) {
	if "-" == op {
		return nil, incompatible(op)
	}
	switch right.(type) {
	case string:
		return (strconv.FormatBool(lhs) + right.(string)), nil
	default:
		return nil, incompatible(op)
	}
}

func multByInt(lhs int64, right interface{}, op string) (interface{}, error) {
	switch right.(type) {
	case int64:
		rhs := right.(int64)
		if "/" == op {
			return float64(lhs) / float64(rhs), nil
		}
		return lhs * rhs, nil
	case float64:
		return multByFloat(float64(lhs), right, op)
	case string:
		if "/" == op {
			return nil, incompatible(op)
		}
		if lhs < int64(0) {
			return nil, fmt.Errorf("Cannot multiply string by negative number")
		}
		rhs := right.(string)
		r := make([]string, lhs)
		for i := int64(0); i < lhs; i++ {
			r[i] = rhs
		}
		return strings.Join(r, ""), nil
	default:
		return nil, incompatible(op)
	}
}

func multByFloat(lhs float64, right interface{}, op string) (interface{}, error) {
	switch right.(type) {
	case int64:
		rhs := float64(right.(int64))
		if "/" == op {
			return lhs / rhs, nil
		}
		return lhs * rhs, nil
	case float64:
		rhs := right.(float64)
		if "/" == op {
			return lhs / rhs, nil
		}
		return lhs * rhs, nil
	case string:
		return multByInt(int64(lhs), right, op)
	default:
		return nil, incompatible(op)
	}
}

func multByString(lhs string, right interface{}, op string) (interface{}, error) {
	switch right.(type) {
	case int64:
		rhs := right.(int64)
		return multByInt(rhs, lhs, op)
	case float64:
		rhs := int64(right.(float64))
		return multByInt(rhs, lhs, op)
	default:
		return nil, incompatible(op)
	}
}

func incompatible(op string) error {
	return fmt.Errorf("Incompatible types for operator %q", op)
}
