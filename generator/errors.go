package generator

import "fmt"

type FatalError struct {
	message string
}

func (e FatalError) Error() string {
	return e.message
}

func NewFatalError(msg string, tokens ...interface{}) *FatalError {
	return &FatalError{message: fmt.Sprintf(msg, tokens...)}
}

type WarningError struct {
	message string
}

func (e WarningError) Error() string {
	return e.message
}

func NewWarningError(msg string, tokens ...interface{}) *WarningError {
	return &WarningError{message: fmt.Sprintf(msg, tokens...)}
}
