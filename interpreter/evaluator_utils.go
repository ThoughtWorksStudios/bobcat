package interpreter

import (
	"fmt"
	. "github.com/ThoughtWorksStudios/bobcat/builtins"
	. "github.com/ThoughtWorksStudios/bobcat/common"
	"strconv"
	"strings"
	"time"
)

type Lambda struct {
	name     string
	params   []string
	executor DeferredResolver
	scope    *Scope
}

func (l *Lambda) Name() string {
	if l.name == "" {
		return "<anonymous>"
	} else {
		return l.name
	}
}

func (l Lambda) String() string {
	return fmt.Sprintf("lambda %s(%s){ ... }", l.Name(), strings.Join(l.params, ", "))
}

func (l *Lambda) Call(boundArgs ...interface{}) (interface{}, error) {
	if expected, actual := len(l.params), len(boundArgs); expected != actual {
		return nil, fmt.Errorf("%s: mismatched arity; expected %d arguments, but got %d", l.String(), expected, actual)
	}

	syms := make(SymbolTable)

	if len(l.params) > 0 {
		for i, s := range l.params {
			syms[s] = boundArgs[i]
		}
	}

	return l.executor(TransientScope(l.scope, syms))
}

func NewLambda(name string, params []string, body DeferredResolver, scope *Scope) *Lambda {
	return &Lambda{name: name, params: params, executor: body, scope: scope}
}

type ExecQueue struct {
	expr []interface{}
}

/** returns the result of the last statement */
func (eq *ExecQueue) Run(scope *Scope) (interface{}, error) {
	var val interface{}
	var err error

	for _, ex := range eq.expr {
		if res, ok := ex.(DeferredResolver); ok {
			if val, err = res(scope); err != nil {
				return nil, err
			}
		} else {
			val = ex
		}
	}

	return val, nil
}

func NewExecQueue(compiledExpressions []interface{}) *ExecQueue {
	return &ExecQueue{expr: compiledExpressions}
}

func (i *Interpreter) handleDeferredLHS(op string, left DeferredResolver, right interface{}) DeferredResolver {
	return func(scope *Scope) (interface{}, error) {
		if lhs, err := left(scope); err != nil {
			return nil, err
		} else {
			return i.ApplyOperator(op, lhs, right, scope, false)
		}
	}
}

func (i *Interpreter) addToInt(op string, lhs int64, right interface{}, scope *Scope, deferred bool) (interface{}, error) {
	switch right.(type) {
	case int64:
		rhs := right.(int64)
		if "-" == op {
			return lhs - rhs, nil
		}
		return lhs + rhs, nil
	case float64:
		return i.addToFloat(op, float64(lhs), right, scope, deferred)
	case string:
		if "-" == op {
			return nil, incompatible(op)
		}
		return strconv.FormatInt(lhs, 10) + right.(string), nil
	case *TimeWithFormat:
		if "-" == op {
			return nil, refuseTimeAsRHS(op)
		}
		return i.addToTime(op, right.(*TimeWithFormat), lhs, scope, deferred)
	case DeferredResolver:
		closure := func(scope *Scope) (interface{}, error) {
			if rhs, err := right.(DeferredResolver)(scope); err == nil {
				return i.addToInt(op, lhs, rhs, scope, false)
			} else {
				return nil, err
			}
		}
		if deferred {
			return closure, nil
		}
		return closure(scope)
	default:
		return nil, incompatible(op)
	}
}

func (i *Interpreter) addToFloat(op string, lhs float64, right interface{}, scope *Scope, deferred bool) (interface{}, error) {
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
		return strconv.FormatFloat(lhs, 'f', -1, 64) + right.(string), nil
	case *TimeWithFormat:
		if "-" == op {
			return nil, refuseTimeAsRHS(op)
		}
		return i.addToTime(op, right.(*TimeWithFormat), lhs, scope, deferred)
	case DeferredResolver:
		closure := func(scope *Scope) (interface{}, error) {
			if rhs, err := right.(DeferredResolver)(scope); err == nil {
				return i.addToFloat(op, lhs, rhs, scope, false)
			} else {
				return nil, err
			}
		}

		if deferred {
			return closure, nil
		}
		return closure(scope)
	default:
		return nil, incompatible(op)
	}
}

func (i *Interpreter) addToString(op string, lhs string, right interface{}, scope *Scope, deferred bool) (interface{}, error) {
	if "-" == op {
		return nil, incompatible(op)
	}
	switch right.(type) {
	case string:
		return lhs + right.(string), nil
	case int64:
		return lhs + strconv.FormatInt(right.(int64), 10), nil
	case float64:
		return lhs + strconv.FormatFloat(right.(float64), 'f', -1, 64), nil
	case bool:
		return lhs + strconv.FormatBool(right.(bool)), nil
	case *TimeWithFormat:
		return i.addToString(op, lhs, right.(*TimeWithFormat).Formatted(), scope, deferred)
	case DeferredResolver:
		closure := func(scope *Scope) (interface{}, error) {
			if rhs, err := right.(DeferredResolver)(scope); err == nil {
				return i.addToString(op, lhs, rhs, scope, false)
			} else {
				return nil, err
			}
		}

		if deferred {
			return closure, nil
		}
		return closure(scope)
	default:
		return nil, incompatible(op)
	}
}

func (i *Interpreter) addToBool(op string, lhs bool, right interface{}, scope *Scope, deferred bool) (interface{}, error) {
	if "-" == op {
		return nil, incompatible(op)
	}
	switch right.(type) {
	case string:
		return strconv.FormatBool(lhs) + right.(string), nil
	case DeferredResolver:
		closure := func(scope *Scope) (interface{}, error) {
			if rhs, err := right.(DeferredResolver)(scope); err == nil {
				return i.addToBool(op, lhs, rhs, scope, false)
			} else {
				return nil, err
			}
		}

		if deferred {
			return closure, nil
		}
		return closure(scope)
	default:
		return nil, incompatible(op)
	}
}

func (i *Interpreter) addToTime(op string, lhs *TimeWithFormat, right interface{}, scope *Scope, deferred bool) (interface{}, error) {

	switch right.(type) {
	case int64:
		nanoPerMs := int64(time.Millisecond / time.Nanosecond)
		rhs := time.Duration(right.(int64) * nanoPerMs)

		var result time.Time
		if "-" == op {
			result = lhs.Time.Add(-rhs)
		} else {
			result = lhs.Time.Add(rhs)
		}
		return NewTimeWithFormat(result, lhs.Format), nil
	case float64:
		return i.addToTime(op, lhs, int64(right.(float64)), scope, deferred)
	case string:
		return i.addToString(op, lhs.Formatted(), right, scope, deferred)
	case *TimeWithFormat:
		if "+" == op {
			return nil, fmt.Errorf("Cannot add Time to another Time")
		}

		return int64(lhs.Time.Sub(right.(*TimeWithFormat).Time) / time.Millisecond), nil // gets duration in milliseconds
	case DeferredResolver:
		closure := func(scope *Scope) (interface{}, error) {
			if rhs, err := right.(DeferredResolver)(scope); err == nil {
				return i.addToTime(op, lhs, rhs, scope, false)
			} else {
				return nil, err
			}
		}
		if deferred {
			return closure, nil
		}
		return closure(scope)
	default:
		return nil, incompatible(op)
	}
}

func (i *Interpreter) multByInt(op string, lhs int64, right interface{}, scope *Scope, deferred bool) (interface{}, error) {
	switch right.(type) {
	case int64:
		rhs := right.(int64)
		if "/" == op {
			return float64(lhs) / float64(rhs), nil
		}
		return lhs * rhs, nil
	case float64:
		return i.multByFloat(op, float64(lhs), right, scope, deferred)
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
	case DeferredResolver:
		closure := func(scope *Scope) (interface{}, error) {
			if rhs, err := right.(DeferredResolver)(scope); err == nil {
				return i.multByInt(op, lhs, rhs, scope, false)
			} else {
				return nil, err
			}
		}

		if deferred {
			return closure, nil
		}
		return closure(scope)
	default:
		return nil, incompatible(op)
	}
}

func (i *Interpreter) multByFloat(op string, lhs float64, right interface{}, scope *Scope, deferred bool) (interface{}, error) {
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
		return i.multByInt(op, int64(lhs), right, scope, deferred)
	case DeferredResolver:
		closure := func(scope *Scope) (interface{}, error) {
			if rhs, err := right.(DeferredResolver)(scope); err == nil {
				return i.multByFloat(op, lhs, rhs, scope, false)
			} else {
				return nil, err
			}
		}

		if deferred {
			return closure, nil
		}
		return closure(scope)
	default:
		return nil, incompatible(op)
	}
}

func (i *Interpreter) multByString(op string, lhs string, right interface{}, scope *Scope, deferred bool) (interface{}, error) {
	switch right.(type) {
	case int64:
		rhs := right.(int64)
		return i.multByInt(op, rhs, lhs, scope, deferred)
	case float64:
		rhs := int64(right.(float64))
		return i.multByInt(op, rhs, lhs, scope, deferred)
	case DeferredResolver:
		closure := func(scope *Scope) (interface{}, error) {
			if rhs, err := right.(DeferredResolver)(scope); err == nil {
				return i.multByString(op, lhs, rhs, scope, false)
			} else {
				return nil, err
			}
		}

		if deferred {
			return closure, nil
		}
		return closure(scope)
	default:
		return nil, incompatible(op)
	}
}

func refuseTimeAsRHS(op string) error {
	return fmt.Errorf("Refusing to coerce Time to numeric as right-hand side for `%s` operator", op)
}

func incompatible(op string) error {
	return fmt.Errorf("Incompatible types for operator %q", op)
}
