package generator

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

func appendData(entities map[string][]map[string]interface{}, existingData []byte) map[string][]map[string]interface{} {
	var x map[string]interface{}
	json.Unmarshal(existingData, &x)
	for k, v := range x {
		r, _ := v.([]interface{})
		if _, ok := entities[k]; !ok {
			entities[k] = make([]map[string]interface{}, 0)
		}
		for _, ent := range r {
			entity, _ := ent.(map[string]interface{})
			entities[k] = append(entities[k], entity)
		}
	}
	return entities
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
