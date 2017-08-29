package common

import (
	"bufio"
	"os"
)

type EntityResult map[string]interface{}

type Emitter interface {
	Receiver() EntityResult
	Emit(entity EntityResult, entityType string) error
	NextEmitter(current EntityResult, key string, isMultiValue bool) Emitter
	Finalize() error
}

func createWriterFor(filename string) (*os.File, *bufio.Writer, error) {
	os_writer, err := os.Create(filename)
	if err != nil {
		return nil, nil, err
	}

	writer := bufio.NewWriter(os_writer)
	return os_writer, writer, nil
}
