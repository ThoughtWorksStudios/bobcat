package common

import (
  "fmt"
  "os"
)

// print arbitrary messages to STDERR; useful when making debug statements
// for development
func Debug(f string, t ...interface{}) {
  fmt.Fprintf(os.Stderr, f+"\n", t...)
}
