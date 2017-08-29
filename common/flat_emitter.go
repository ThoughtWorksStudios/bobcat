package common

import (
	"bufio"
	j "github.com/json-iterator/go"
	"os"
	"io"
	"errors"
)

type FlatEmitter struct {
	os_writer *os.File
	writer    io.Writer
	encoder   Encoder
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
	var err error

	if !f.first {
		if _, err = f.writer.Write(delimeter); err != nil {
			return err
		}
	} else {
		f.first = false
	}

	if err = f.encoder.Encode(entity); err != nil {
		return err
	}

	return nil
}

func (f *FlatEmitter) Finalize() error {
	var err error
	if _, err = f.writer.Write(end); err != nil {
		return err
	}

	var bufioWriter *bufio.Writer
	var ok bool
	if bufioWriter, ok = f.writer.(*bufio.Writer); !ok {
		return errors.New("Expected bufioWriter but did not receive one")
	}

	if err = bufioWriter.Flush(); err != nil {
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
	if encoder, ok := emitter.encoder.(*j.Encoder); ok {
		encoder.SetIndent("", "  ")
	}
	return emitter, nil
}
