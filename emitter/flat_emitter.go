package emitter

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
	"io"
)

type FlatEmitter struct {
	writer  io.WriteCloser
	encoder Encoder
	first   bool
}

func (f *FlatEmitter) Emit(entity EntityResult, entityType string) error {
	var err error

	if !f.first {
		if _, err = f.writer.Write(DelimeterSeq); err != nil {
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

func (f *FlatEmitter) Init() error {
	if _, err := f.writer.Write(StartSeq); err != nil {
		return err
	}

	return nil
}

func (f *FlatEmitter) Finalize() error {
	if _, err := f.writer.Write(EndSeq); err != nil {
		return err
	}

	return f.writer.Close()
}

func (f *FlatEmitter) NextEmitter(current EntityStore, key string, isMultiValue bool) Emitter {
	return f
}

func (f *FlatEmitter) Receiver() EntityStore {
	return nil
}

/**
 * Creates a FlatEmitter with a generic io.WriterCloser
 */
func NewFlatEmitter(writer io.WriteCloser) Emitter {
	return &FlatEmitter{first: true, writer: writer, encoder: DefaultEncoder(writer)}
}

/**
 * Creates a FlatEmitter with a FileWriter for the given filename
 */
func FlatEmitterForFile(filename string) (Emitter, error) {
	if writer, err := NewFileWriter(filename); err != nil {
		return nil, err
	} else {
		return NewFlatEmitter(writer), nil
	}
}
