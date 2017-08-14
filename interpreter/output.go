package interpreter

import (
	"bufio"
	"fmt"
	g "github.com/ThoughtWorksStudios/bobcat/generator"
	"github.com/json-iterator/go"
	"io"
	"os"
)

type GenerationOutput map[string]g.GeneratedEntities

func (output GenerationOutput) addAndAppend(entityName string, entities g.GeneratedEntities) {
	if _, ok := output[entityName]; ok {
		output[entityName] = output[entityName].Concat(entities)
	} else {
		output[entityName] = entities
	}
}

func (output GenerationOutput) writeFilePerKey() error {
	for k, v := range output {
		out, err := createWriterFor(fmt.Sprintf("%s.json", k))
		if err != nil {
			return err
		}
		d := GenerationOutput{}
		d[k] = v
		if err = d.write(out); err != nil {
			return err
		}
	}
	return nil
}

func (output GenerationOutput) writeToFile(dest string) error {
	out, err := createWriterFor(dest)
	if err != nil {
		return err
	}

	return output.write(out)
}

func (output GenerationOutput) write(out io.Writer) error {
	if closeable, doClose := isClosable(out); doClose {
		defer closeable.Close()
	}

	writer := bufio.NewWriter(out)
	encoder := jsoniter.NewEncoder(writer)
	encoder.SetIndent("", "\t")

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
