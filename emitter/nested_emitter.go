package emitter

import (
	"fmt"
	. "github.com/ThoughtWorksStudios/bobcat/common"
	"io"
)

type NestedEmitter struct {
	cursor  *Cursor
	writer  io.WriteCloser
	encoder Encoder
}

func (n *NestedEmitter) Init() error {
	n.cursor = &Cursor{current: make(EntityResult)}
	return nil
}

func (n *NestedEmitter) Emit(entity EntityResult, entityType string) error {
	return n.cursor.Insert(entity)
}

func (n *NestedEmitter) Finalize() error {
	if err := n.encoder.Encode(n.Receiver()); err != nil {
		return err
	}

	return n.writer.Close()
}

func (n *NestedEmitter) NextEmitter(current EntityResult, key string, isMultiValue bool) Emitter {
	return &NestedEmitter{cursor: &Cursor{current: current, key: key, isMultiValue: isMultiValue}}
}

func (n *NestedEmitter) Receiver() EntityResult {
	return n.cursor.current
}

/**
 * Creates a NestedEmitter with a generic io.WriterCloser
 */
func InitNestedEmitter(writer io.WriteCloser) (Emitter, error) {
	emitter := &NestedEmitter{writer: writer, encoder: DefaultEncoder(writer)}

	if err := emitter.Init(); err != nil {
		return nil, err
	}

	return emitter, nil
}

/**
 * Creates a NestedEmitter with a FileWriter for the given filename
 */
func NewNestedEmitter(filename string) (Emitter, error) {
	if writer, err := NewFileWriter(filename); err != nil {
		return nil, err
	} else {
		return InitNestedEmitter(writer)
	}
}

/**
 * Wraps a location to insert an emitted entity.
 *
 * Essentially just holds a reference to a target EntityResult (i.e. the parent),
 * and a field key. Abstracts multi-value awareness from the NestedEmitter.
 */
type Cursor struct {
	current      EntityResult
	key          string
	isMultiValue bool
}

func (c *Cursor) Insert(value interface{}) error {
	if c.isMultiValue {
		var result []EntityResult
		var ok bool

		if original := c.current[c.key]; nil == original {
			result = make([]EntityResult, 0)
		} else {
			if result, ok = original.([]EntityResult); !ok {
				return fmt.Errorf("Expected an entity set")
			}
		}

		c.current[c.key] = append(result, value.(EntityResult))
	} else {
		c.current[c.key] = value
	}

	return nil
}
