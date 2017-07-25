package generator

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

func appendContent(newData interface{}, data []byte) []map[string]interface{} {
	var x []map[string]interface{}
	entityList, _ := newData.([]map[string]interface{})
	json.Unmarshal(data, &x)
	for _, entity := range x {
		entityList = append(entityList, entity)
	}
	return entityList
}

func createWriterFor(filename string) (io.Writer, []byte, error) {
	var f interface{}
	var existingOutput []byte
	var err error
	if _, exists := os.Stat(filename); !os.IsNotExist(exists) {
		existingOutput, _ = ioutil.ReadFile(filename)
	}

	f, err = os.Create(filename)
	if err != nil {
		return nil, nil, err
	}
	out, _ := f.(io.Writer)

	return out, existingOutput, nil
}

func isClosable(v interface{}) (io.Closer, bool) {
	closeable, doClose := v.(io.Closer)
	return closeable, doClose
}
