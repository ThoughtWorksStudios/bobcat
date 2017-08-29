package common

import (
	"bufio"
	j "github.com/json-iterator/go"
	"os"
)

type FlatEmitter struct {
	os_writer *os.File
	writer    *bufio.Writer
	encoder   *j.Encoder
	first     bool
}

const (
	START     = "[\n"
	DELIMITER = ",\n"
	END       = "\n]"
)

var start = []byte(START)
var delimeter = []byte(DELIMITER)
var end = []byte(END)

func (f *FlatEmitter) Emit(entity EntityResult, entityType string) error {
	if !f.first {
		f.writer.Write(delimeter)
	} else {
		f.first = false
	}

	if err := f.encoder.Encode(entity); err != nil {
		return err
	}

	return nil
}

func (f *FlatEmitter) Finalize() error {
	f.writer.Write(end)

	if err := f.writer.Flush(); err != nil {
		return err
	}

	return f.os_writer.Close()
}

func (f *FlatEmitter) NextEmitter(current EntityResult, key string, isMultiValue bool) Emitter {
	return f
}

func (f *FlatEmitter) Receiver() EntityResult {
	return nil
}

func NewFlatEmitter(filename string) (Emitter, error) {
	emitter := &FlatEmitter{first: true}
	var err error
	if emitter.os_writer, emitter.writer, err = createWriterFor(filename); err != nil {
		return nil, err
	}

	emitter.writer.Write(start)
	emitter.encoder = j.ConfigFastest.NewEncoder(emitter.writer)
	emitter.encoder.SetIndent("", "  ")
	return emitter, nil
}
