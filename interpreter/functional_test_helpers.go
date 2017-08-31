package interpreter

import (
	. "github.com/ThoughtWorksStudios/bobcat/emitter"
	"github.com/json-iterator/go"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

type Deferrable func()
type EmitterFactory func(filename string) (Emitter, error)
type JsonArrayIterator func(element jsoniter.Any, index int)
type JsonYielder func(data jsoniter.Any, raw string)

func withParsed(t *testing.T, filename string, performAssertions JsonYielder) {
	data := readFile(t, filename)
	performAssertions(jsoniter.Get(data), string(data))
}

func eachElem(array jsoniter.Any, iterator JsonArrayIterator) {
	size := array.Size()
	for i := 0; i < size; i++ {
		iterator(array.Get(i), i)
	}
}

func readFile(t *testing.T, filename string) []byte {
	if b, e := ioutil.ReadFile(filename); e == nil {
		return b
	} else {
		t.Fatalf("Failed to read file %q: %v", filename, e)
		return []byte("")
	}
}

func withBasename(runnable func(basename string)) {
	rand.Seed(time.Now().UnixNano())
	runnable(strconv.Itoa(rand.Int()))
}

func cleanup(filename string) Deferrable {
	return func() {
		os.Remove(filename)
	}
}
