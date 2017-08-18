package interpreter

import (
	"bufio"
	"fmt"
	g "github.com/ThoughtWorksStudios/bobcat/generator"
	"github.com/json-iterator/go"
	"io"
	"os"
)

type GenerationOutput interface {
	writeToFile(dest string) error
	write(out io.Writer) error
	addAndAppend(entityName string, entities g.GeneratedEntities) GenerationOutput
	writeFilePerKey() error
}
type FlatOutput g.GeneratedEntities
type NestedOutput map[string]g.GeneratedEntities

func (output FlatOutput) Concat(newEntities g.GeneratedEntities) {
	for _, entity := range newEntities {
		output = append(output, entity)
	}
}

func (output FlatOutput) addAndAppend(entityName string, entities g.GeneratedEntities) GenerationOutput {
	return append(output, entities...)
}

func (output FlatOutput) writeToFile(dest string) error {
	out, err := createWriterFor(dest)
	if err != nil {
		return err
	}

	return output.write(out)
}

func (output FlatOutput) write(out io.Writer) error {
	if closeable, doClose := isClosable(out); doClose {
		defer closeable.Close()
	}

	writer := bufio.NewWriter(out)
	encoder := jsoniter.ConfigFastest.NewEncoder(writer)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(output); err != nil {
		return err
	}

	return writer.Flush()
}

func (output FlatOutput) writeFilePerKey() error {
	return fmt.Errorf("Cannot split output for flat output")
}

func (output NestedOutput) addAndAppend(entityName string, entities g.GeneratedEntities) GenerationOutput {
	if _, ok := output[entityName]; ok {
		output[entityName] = output[entityName].Concat(entities)
	} else {
		output[entityName] = entities
	}
	return output
}

func (output NestedOutput) writeFilePerKey() error {
	for k, v := range output {
		out, err := createWriterFor(fmt.Sprintf("%s.json", k))
		if err != nil {
			return err
		}
		d := NestedOutput{}
		d[k] = v
		if err = d.write(out); err != nil {
			return err
		}
	}
	return nil
}

func (output NestedOutput) writeToFile(dest string) error {
	out, err := createWriterFor(dest)
	if err != nil {
		return err
	}

	return output.write(out)
}

func (output NestedOutput) write(out io.Writer) error {
	if closeable, doClose := isClosable(out); doClose {
		defer closeable.Close()
	}

	writer := bufio.NewWriter(out)
	encoder := jsoniter.ConfigFastest.NewEncoder(writer)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(output); err != nil {
		return err
	}

	return writer.Flush()
}

func isClosable(v interface{}) (io.Closer, bool) {
	closeable, doClose := v.(io.Closer)
	return closeable, doClose
}

func createWriterFor(filename string) (io.Writer, error) {
	var f interface{}
	var err error
	f, err = os.Create(filename)
	if err != nil {
		return nil, err
	}
	out, _ := f.(io.Writer)

	return out, nil
}
