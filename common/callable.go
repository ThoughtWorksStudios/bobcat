package common

type Callable interface {
	Name() string
	Call(args ...interface{}) (interface{}, error)
}
