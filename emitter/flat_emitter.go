package emitter

import (
	"bufio"
	. "github.com/ThoughtWorksStudios/bobcat/common"
	j "github.com/json-iterator/go"
	"io"
	"os"
)

const (
	START     = "[\n"
	DELIMITER = ",\n"
	END       = "\n]"
)

// essentially constants, except that Go doesn't allow slice constants
// due to compile-time restrictions
var start = []byte(START)
var delimeter = []byte(DELIMITER)
var end = []byte(END)

type FlatEmitter struct {
	os_writer *os.File
	writer    io.Writer
	encoder   Encoder
	first     bool
}

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

	if bufioWriter, ok := f.writer.(*bufio.Writer); ok {
		if err = bufioWriter.Flush(); err != nil {
			return err
		}
	}

	return f.os_writer.Close()
}

func (f *FlatEmitter) NextEmitter(current EntityResult, key string, isMultiValue bool) Emitter {
	return f
}

func (f *FlatEmitter) Receiver() EntityResult {
	return nil
}

func (f *FlatEmitter) Init() error {
	if _, err := f.writer.Write(start); err != nil {
		return err
	}
	f.encoder = j.ConfigFastest.NewEncoder(f.writer)
	if encoder, ok := f.encoder.(*j.Encoder); ok {
		encoder.SetIndent("", "  ")
	}
	return nil
}


func NewFlatEmitter(filename string) (Emitter, error) {
	emitter := &FlatEmitter{first: true}
	var err error
	if emitter.os_writer, emitter.writer, err = createWriterFor(filename); err != nil {
		return nil, err
	}

	if err := emitter.Init(); err != nil {
		return nil, err
	}

	return emitter, nil
}
