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

/**
 * Parses a JSON file and binds the resultant jsoniter.Any value as a
 * local variable to the closure. Intended to make parsed value readily
 * available to assertion statements executed within the closure.
 */
func withParsed(t *testing.T, filename string, performAssertions JsonYielder) {
	data := readFile(t, filename)
	performAssertions(jsoniter.Get(data), string(data))
}

/**
 * Convenience func to iterate over a jsoniter.Any when it is expected
 * to be an array. Binds each element and its index as local variables
 * to the iterator closure. Useful for performing a set of assertions
 * on each element.
 */
func eachElem(array jsoniter.Any, iterator JsonArrayIterator) {
	size := array.Size()
	for i := 0; i < size; i++ {
		iterator(array.Get(i), i)
	}
}

/** convenience func to read a file to string */
func readFile(t *testing.T, filename string) []byte {
	if b, e := ioutil.ReadFile(filename); e == nil {
		return b
	} else {
		t.Fatalf("Failed to read file %q: %v", filename, e)
		return []byte("")
	}
}

/**
 * Generates a unique file basename (i.e. sans `.json` extension) for
 * each invocation, and executes a runnable, binding the basename as a
 * local variable to the closure. Useful for functional tests that need
 * to create files. Uniqueness allows for concurrent test execution.
 */
func withBasename(runnable func(basename string)) {
	rand.Seed(time.Now().UnixNano())
	runnable(strconv.Itoa(rand.Int()))
}

/**
 * Returns a deferrable function to remove a file. Used for teardown,
 * and was designed to complement withBasename()
 */
func cleanup(filename string) Deferrable {
	return func() {
		os.Remove(filename)
	}
}
