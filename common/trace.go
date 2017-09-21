package common

import (
	"fmt"
	"log"
	"os"
)

var _TRACE bool
var indent string
var tag string

func init() {
	_TRACE = os.Getenv("TRACE") == "true"
}

// Print arbitrary messages to STDERR if _TRACE is enabled; useful for devel debug output
func Msg(f string, t ...interface{}) {
	if _TRACE {
		printStdErr(f, t...)
	}
}

func Warn(f string, t ...interface{}) {
	tag = "* WARN * "
	printStdErr(f, t...)
	tag = ""
}

// Forcefully exit with message
func Die(f string, t ...interface{}) {
	log.Fatalf(f+"\n", t...)
}

// increase indent of Msg()
func Bump() {
	indent += "|   "
}

// decrease indent of Msg()
func Dunk() {
	if len(indent) > 0 {
		indent = indent[4:]
	}
}

func ResetIndent() {
	indent = ""
}

// Temporarily enables output for the duration of the lambda
func WithTrace(lambda func()) {
	orig := _TRACE
	_TRACE = true
	defer func() { _TRACE = orig }()
	lambda()
}

func printStdErr(f string, t ...interface{}) {
	fmt.Fprintf(os.Stderr, tag+indent+f+"\n", t...)
}
