package generator

import (
	"bufio"
	"encoding/json"
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

func (gc GeneratedContent) WriteToFile(dest string) error {
	out, err := createWriterFor(dest)
	if err != nil {
		return err
	}

	if closeable, doClose := isClosable(out); doClose {
		defer closeable.Close()
	}

	writer := bufio.NewWriter(out)
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "\t")

	if err = encoder.Encode(gc); err != nil {
		return err
	}

	return writer.Flush()
}
