package emitter

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
	"github.com/json-iterator/go"
	"io"
)

const (
	START     = "[\n"
	DELIMITER = ",\n"
	END       = "\n]"
)

// essentially constants, except that Go doesn't allow slice constants
// due to compile-time restrictions
var StartSeq = []byte(START)
var DelimeterSeq = []byte(DELIMITER)
var EndSeq = []byte(END)

type Encoder interface {
	Encode(val interface{}) error
}

type Emitter interface {
	/** Emit() accepts and processes an entity */
	Emit(entity EntityResult, entityType string) error

	/** Called once when Emitter is created */
	Init() error

	/** Called once after interpreter exits and all generation is complete */
	Finalize() error

	/** NextEmitter() returns a continuation, as an Emitter, to handle subsequent calls to Emit() */
	NextEmitter(current EntityStore, key string, isMultiValue bool) Emitter

	/**
	 * Receiver() returns the EntityStore referenced by the current continuation; certain
	 * Emitters will return nil by design (e.g. streaming Emitters such as FlatEmitter)
	 */
	Receiver() EntityStore
}

func DefaultEncoder(writer io.Writer) Encoder {
	encoder := jsoniter.ConfigFastest.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder
}
