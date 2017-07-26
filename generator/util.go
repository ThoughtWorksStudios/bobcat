package generator

import (
	"io"
	"os"
)

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
