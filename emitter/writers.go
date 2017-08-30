package emitter

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func NewFileWriter(filename string) (io.WriteCloser, error) {
	fw := &FileWriter{}

	if err := fw.Open(filename); err != nil {
		return nil, err
	}

	return fw, nil
}

/**
 * Buffered io.WriteCloser backed by a real file
 */
type FileWriter struct {
	file   *os.File
	writer *bufio.Writer
}

func (f *FileWriter) Open(filename string) error {
	if nil != f.file {
		return fmt.Errorf("Refusing to open %q; this writer is already associated with another open file (%v)", f.file, filename)
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	f.file = file
	f.writer = bufio.NewWriter(file)
	return nil
}

func (f *FileWriter) Write(payload []byte) (int, error) {
	return f.writer.Write(payload)
}

func (f *FileWriter) Close() error {
	if err := f.writer.Flush(); err != nil {
		return err
	}

	return f.file.Close()
}

/**
 * String-backed io.WriterCloser; useful for tests
 */
type StringWriter struct {
	result string
}

func (s *StringWriter) Write(payload []byte) (int, error) {
	s.result += string(payload)
	return len(payload), nil
}

func (s *StringWriter) Close() error {
	return nil
}

func (s *StringWriter) Reset() {
	s.result = ""
}

func (s StringWriter) String() string {
	return s.result
}
