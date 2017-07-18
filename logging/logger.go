package logging

import (
	"fmt"
	"log"
)

type ILogger interface {
	Warn(msg string, tokens ...interface{})
	Die(prefix, msg string, tokens ...interface{})
}

type DefaultLogger struct {
	ILogger
}

func (l *DefaultLogger) Die(prefix, msg string, tokens ...interface{}) {
	log.SetFlags(0)
	log.Fatalf("ERROR %v: %v\n", prefix, fmt.Sprintf(msg, tokens...))
}

func (l *DefaultLogger) Warn(msg string, tokens ...interface{}) {
	log.SetFlags(0)
	log.Println("WARNING:", fmt.Sprintf(msg, tokens...))
}
