package logging

import (
	"fmt"
	"log"
)

type ILogger interface {
	Warn(msg string, tokens ...interface{})
	Die(location, msg string, tokens ...interface{})
}

type DefaultLogger struct {
	ILogger
}

func (l *DefaultLogger) Die(location, msg string, tokens ...interface{}) {
	log.Fatalf("FATAL ERROR (%v): %v\n", location, fmt.Sprintf(msg, tokens...))
}

func (l *DefaultLogger) Warn(msg string, tokens ...interface{}) {
	log.Println("WARNING:", fmt.Sprintf(msg, tokens...))
}
