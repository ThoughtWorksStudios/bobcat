package common

import (
	"bufio"
	j "github.com/json-iterator/go"
	"io"
	"os"
)

type FlatEmitter struct {
	os_writer *os.File
	writer    io.Writer
	encoder   *j.Encoder
}

var delimeter = []byte(",\n")

func (f *FlatEmitter) Emit(entity EntityResult, entityType string) error {
	if err := f.encoder.Encode(entity); err != nil {
		return err
	}

	f.writer.Write(delimeter)
	return nil
}

func (f *FlatEmitter) Finalize() error {
	if err := f.writer.(*bufio.Writer).Flush(); err != nil {
		return err
	}
	stat, err := f.os_writer.Stat()
	if err != nil {
		return err
	}
	size := stat.Size()
	new_offset := size - int64(len(delimeter))
	if new_offset > 0 {
		f.os_writer.Truncate(new_offset)
		f.os_writer.Seek(new_offset, 0)
	}

	f.os_writer.Write([]byte("\n]"))
	return f.os_writer.Close()
}

func (f *FlatEmitter) NextEmitter(current EntityResult, key string, isMultiValue bool) Emitter {
	return f
}

func (f *FlatEmitter) Receiver() EntityResult {
	return nil
}

func NewFlatEmitter(filename string) (Emitter, error) {
	emitter := &FlatEmitter{}
	var err error
	if emitter.os_writer, emitter.writer, err = createWriterFor(filename); err != nil {
		return nil, err
	}

	emitter.writer.Write([]byte("[\n"))
	emitter.encoder = j.ConfigFastest.NewEncoder(emitter.writer)
	emitter.encoder.SetIndent("", "  ")
	return emitter, nil
}
