package common

import (
  "fmt"
  "os"
)

var TRACE bool;

func init() {
  TRACE = os.Getenv("TRACE") == "true"
}

// print arbitrary messages to STDERR; useful when making debug statements
// for development
func Debug(f string, t ...interface{}) {
  if TRACE {
    fmt.Fprintf(os.Stderr, f+"\n", t...)
  }
}

func WithTrace(f func()) {
  orig := TRACE
  TRACE = true
  defer func() { TRACE = orig }()
  f()
}