package common

import (
	"bufio"
	"fmt"
	j "github.com/json-iterator/go"
	"io"
	"os"
)

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

/**
 * A list of entities
 */
type GeneratedEntities []EntityResult

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
