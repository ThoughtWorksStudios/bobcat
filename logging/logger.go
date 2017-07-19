package logging

import (
	"fmt"
	"log"
)

func init() {
	log.SetFlags(0) // clear flags, no timestamps
}

type ILogger interface {
	Die(err error)
	Warn(msg string, tokens ...interface{})
}

type DefaultLogger struct {
	ILogger
}

func (l *DefaultLogger) Die(err error) {
	log.Fatalf("[ERROR] %s", err.Error())
}

func (l *DefaultLogger) Warn(msg string, tokens ...interface{}) {
	log.Println("[WARN] " /* trailing space is intentional; matches same width as [ERROR] */, fmt.Sprintf(msg, tokens...))
}
