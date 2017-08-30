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

func createWriterFor(filename string) (*os.File, io.Writer, error) {
	os_writer, err := os.Create(filename)
	if err != nil {
		return nil, nil, err
	}

	writer := bufio.NewWriter(os_writer)
	return os_writer, writer, nil
}
