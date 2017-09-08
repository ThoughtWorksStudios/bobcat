package emitter

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
	"io"
)

type NestedEmitter struct {
	cursor *Cursor
	writer io.WriteCloser
}

func (n *NestedEmitter) Init() error {
	if _, err := n.writer.Write(StartSeq); err != nil {
		return err
	}

	n.cursor = &Cursor{current: &EntityOutputter{writer: n.writer, encoder: DefaultEncoder(n.writer), first: true}}
	return nil
}

func (n *NestedEmitter) Emit(entity EntityResult, entityType string) error {
	return n.cursor.Insert(entity)
}

func (n *NestedEmitter) Finalize() error {
	if _, err := n.writer.Write(EndSeq); err != nil {
		return err
	}

	return n.writer.Close()
}

func (n *NestedEmitter) NextEmitter(current EntityStore, key string, isMultiValue bool) Emitter {
	return &NestedEmitter{cursor: &Cursor{current: current, key: key, isMultiValue: isMultiValue}}
}

func (n *NestedEmitter) Receiver() EntityStore {
	return n.cursor.current
}

/**
 * Creates a NestedEmitter with a generic io.WriterCloser
 */
func NewNestedEmitter(writer io.WriteCloser) Emitter {
	return &NestedEmitter{writer: writer}
}

/**
 * Creates a NestedEmitter with a FileWriter for the given filename
 */
func NestedEmitterForFile(filename string) (Emitter, error) {
	if writer, err := NewFileWriter(filename); err != nil {
		return nil, err
	} else {
		return NewNestedEmitter(writer), nil
	}
}

/**
 * Wraps a location to insert an emitted entity.
 *
 * Essentially just holds a reference to a target EntityResult (i.e. the parent),
 * and a field key. Abstracts multi-value awareness from the NestedEmitter.
 */
type Cursor struct {
	current      EntityStore
	key          string
	isMultiValue bool
}

func (c *Cursor) Insert(entity EntityResult) error {
	if c.isMultiValue {
		return c.current.AppendTo(c.key, entity)
	}

	return c.current.Set(c.key, entity)
}

type EntityOutputter struct {
	writer  io.Writer
	encoder Encoder
	first   bool
}

func (eo *EntityOutputter) Set(key string, entity EntityResult) error {
	return eo.outputJSON(entity)
}

func (eo *EntityOutputter) AppendTo(key string, entity EntityResult) error {
	return eo.outputJSON(entity)
}

func (eo *EntityOutputter) outputJSON(entity EntityResult) error {
	if !eo.first {
		if _, err := eo.writer.Write(DelimeterSeq); err != nil {
			return err
		}
	} else {
		eo.first = false
	}

	return eo.encoder.Encode(entity)
}
