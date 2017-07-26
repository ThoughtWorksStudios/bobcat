package generator

import (
	"bufio"
	"encoding/json"
)

type GeneratedContent map[string][]map[string]interface{}

func NewGeneratedContent() GeneratedContent {
	return GeneratedContent{}
}

func (gc GeneratedContent) Append(existingData GeneratedContent) GeneratedContent {
	for k, v := range existingData {
		if _, ok := gc[k]; !ok {
			gc[k] = make([]map[string]interface{}, 0)
		}
		for _, entity := range v {
			gc[k] = append(gc[k], entity)
		}
	}
	return gc
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
