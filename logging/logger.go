package logging

import (
	"fmt"
	"log"
)

type ILogger interface {
	Die(msg string, tokens ...interface{})
	Warn(msg string, tokens ...interface{})
}

type DefaultLogger struct {
	ILogger
}

func (l *DefaultLogger) Die(msg string, tokens ...interface{}) {
	log.Fatalln("FATAL:", fmt.Sprintf(msg, tokens...))
}

func (l *DefaultLogger) Warn(msg string, tokens ...interface{}) {
	log.Println("WARNING:", fmt.Sprintf(msg, tokens...))
}
