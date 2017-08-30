package emitter

import (
	"bufio"
	. "github.com/ThoughtWorksStudios/bobcat/common"
	j "github.com/json-iterator/go"
	"io"
	"os"
)

type Encoder interface {
	Encode(val interface{}) error
}

type Emitter interface {
	Receiver() EntityResult
	Emit(entity EntityResult, entityType string) error
	NextEmitter(current EntityResult, key string, isMultiValue bool) Emitter
	Init() error
	Finalize() error
}

func NewEncoder(writer io.Writer) Encoder {
	encoder := j.ConfigFastest.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder
}
