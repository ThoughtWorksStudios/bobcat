package generator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

type GeneratedContent map[string]GeneratedEntities

type GeneratedEntities []GeneratedEntity

type GeneratedEntity map[string]interface{}

func NewGeneratedEntities(count int64) GeneratedEntities {
	return make([]GeneratedEntity, count)
}

func NewGeneratedContent() GeneratedContent {
	return GeneratedContent{}
}

func (gc GeneratedContent) Append(data GeneratedContent) {
	for k, v := range data {
		if _, ok := gc[k]; !ok {
			gc[k] = GeneratedEntities{}
		}
		for _, entity := range v {
			gc[k] = append(gc[k], entity)
		}
	}
}

func (gc GeneratedContent) WriteFilePerKey() error {
	for k, v := range gc {
		out, err := createWriterFor(fmt.Sprintf("%s.json", k))
		if err != nil {
			return err
		}
		d := NewGeneratedContent()
		d[k] = v
		if err = d.writeToFile(out); err != nil {
			return err
		}
	}
	return nil
}

func (gc GeneratedContent) WriteContentToFile(dest string) error {
	out, err := createWriterFor(dest)
	if err != nil {
		return err
	}

	return gc.writeToFile(out)
}

func (gc GeneratedContent) writeToFile(out io.Writer) error {
	if closeable, doClose := isClosable(out); doClose {
		defer closeable.Close()
	}

	writer := bufio.NewWriter(out)
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "\t")

	if err := encoder.Encode(gc); err != nil {
		return err
	}

	return writer.Flush()
}
