package emitter

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"reflect"
	"testing"
	"os"
)

func TestFlatEmitter_Emit(t *testing.T) {
	emptySlice := []byte{}
	testWriter := &TestWriter{emptySlice}
	testEncoder := &TestEncoder{[]byte{}}
	flat_emitter := FlatEmitter{nil, testWriter, testEncoder, true}

	entityResult := make(EntityResult)
	entityResult["foo"] = "bar"

	flat_emitter.Emit(entityResult, "testType")

	testWritersOutputShouldBeEmpty := reflect.DeepEqual(testWriter.WrittenValue, emptySlice)
	Assert(t, testWritersOutputShouldBeEmpty, "Writer should output an empty value on first emit, but instead had: %v", testWriter.WrittenValue)
	testEncodersOutputShouldEqualEntity := reflect.DeepEqual(entityResult, testEncoder.EncodedValue)
	Assert(t, testEncodersOutputShouldEqualEntity, "Expected encoder to contain '%v', but got '%v'", entityResult, testEncoder.EncodedValue)

	flat_emitter.Emit(entityResult, "testType")

	delimeter := []byte(",\n")
	testWritersOutputShouldEqualDelimiter := reflect.DeepEqual(delimeter, testWriter.WrittenValue)
	Assert(t, testWritersOutputShouldEqualDelimiter, "Writer should output '%v' for subsequent emits, but got: %v", delimeter, testWriter.WrittenValue)
}

func TestFlatEmitter_NextEmitter(t *testing.T) {
	flat_emitter := &FlatEmitter{}

	emitter := flat_emitter.NextEmitter(nil, "", false)

	AssertEqual(t, flat_emitter, emitter, "Next emitter should be the same as original flat emitter")
}

func TestFlatEmitter_Receiver(t *testing.T) {
	flat_emitter := &FlatEmitter{}

	entityResult := flat_emitter.Receiver()

	Assert(t, reflect.DeepEqual(entityResult, EntityResult(nil)), "Values do not match '%v' is not '%v'", entityResult, EntityResult(nil))
}

func TestFlatEmitter_Finalize(t *testing.T) {
	testWriter := &TestWriter{}
	os_writer, _ := os.Create("/dev/null")

	flat_emitter := &FlatEmitter{os_writer, testWriter, nil, true}

	err := flat_emitter.Finalize()

	AssertNil(t, err, "Finalize should not return an error, instead got: %v", err)
	expected := []byte("\n]")
	Assert(t, reflect.DeepEqual(expected, testWriter.WrittenValue), "Expected to Write '%v', but got '%v'", expected, testWriter.WrittenValue)
}

func TestFlatEmitter_Init(t *testing.T) {
	testWriter := &TestWriter{}
	flat_emitter := FlatEmitter{nil, testWriter, nil, true}
	flat_emitter.Init()
	expected := []byte("[\n")
	Assert(t, reflect.DeepEqual(expected, testWriter.WrittenValue), "Expected to Write '%v', but got '%v'", expected, testWriter.WrittenValue)
}

func TestNewFlatEmitter(t *testing.T) {
	flat_emitter, err := NewFlatEmitter("/dev/null")

	AssertNil(t, err, "Creating a flat emitter threw and error: %v", err)
	AssertNotNil(t, flat_emitter, "Flat emitter should not be nil")
}