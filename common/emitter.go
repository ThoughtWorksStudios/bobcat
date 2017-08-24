package common

import (
	"bufio"
	"fmt"
	j "github.com/json-iterator/go"
	"io"
	"os"
)

type GeneratedEntities []EntityResult
type EntityResult map[string]interface{}

type Emitter interface {
	Receiver() EntityResult
	Emit(entity EntityResult) error
	NextEmitter(current EntityResult, key string, isMultiValue bool) Emitter
	Finalize() error
}

type FlatEmitter struct {
	os_writer *os.File
	writer    io.Writer
	encoder   *j.Encoder
}

var delimeter = []byte(",\n")

func createWriterFor(filename string) (*os.File, io.Writer, error) {
	os_writer, err := os.Create(filename)
	if err != nil {
		return nil, nil, err
	}

	writer := bufio.NewWriter(os_writer)
	return os_writer, writer, nil
}

func (f *FlatEmitter) Emit(entity EntityResult) error {
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

type NestedEmitter struct {
	cursor    *Cursor
	os_writer *os.File
	writer    io.Writer
	encoder   *j.Encoder
}

func NewNestedEmitter(filename string) (Emitter, error) {
	emitter := &NestedEmitter{}
	var err error

	if filename != "" { // tests don't need to write to file
		if emitter.os_writer, emitter.writer, err = createWriterFor(filename); err != nil {
			return nil, err
		}
		emitter.encoder = j.ConfigFastest.NewEncoder(emitter.writer)
		emitter.encoder.SetIndent("", "  ")
	}

	emitter.cursor = &Cursor{current: make(EntityResult)}
	return emitter, nil
}

type Cursor struct {
	current      EntityResult
	key          string
	isMultiValue bool
}

func (c *Cursor) Insert(value interface{}) error {
	if c.isMultiValue {
		var result GeneratedEntities
		var ok bool

		if original := c.current[c.key]; nil == original {
			result = make(GeneratedEntities, 0)
		} else {
			if result, ok = original.(GeneratedEntities); !ok {
				return fmt.Errorf("Expected an entity set")
			}
		}

		c.current[c.key] = append(result, value.(EntityResult))
	} else {
		c.current[c.key] = value
	}

	return nil
}

func (n *NestedEmitter) Emit(entity EntityResult) error {
	return n.cursor.Insert(entity)
}

func (n *NestedEmitter) Finalize() error {
	if err := n.encoder.Encode(n.Receiver()); err != nil {
		return err
	}

	if err := n.writer.(*bufio.Writer).Flush(); err != nil {
		return err
	}

	return n.os_writer.Close()
}

func (n *NestedEmitter) Receiver() EntityResult {
	return n.cursor.current
}

func (n *NestedEmitter) NextEmitter(current EntityResult, key string, isMultiValue bool) Emitter {
	return &NestedEmitter{cursor: &Cursor{current: current, key: key, isMultiValue: isMultiValue}}
}
