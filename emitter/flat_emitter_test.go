package emitter

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"reflect"
	"testing"
	"os"
)

func TestFlatEmitter_FirstEmitEncodesEntity(t *testing.T) {
	delimeter := []byte{}
	testWriter := &TestWriter{delimeter}
	testEncoder := &TestEncoder{[]byte{}}
	flat_emitter := FlatEmitter{nil, testWriter, testEncoder, true}

	entityResult := make(EntityResult)
	entityResult["foo"] = "bar"

	flat_emitter.Emit(entityResult, "testType")
	Assert(t, reflect.DeepEqual(testWriter.WrittenValue, delimeter), "Did not expect Emit to Write any values")
	Assert(t, reflect.DeepEqual(entityResult, testEncoder.EncodedValue), "Expected Emit to Encode '%v'", entityResult)

	flat_emitter.Emit(entityResult, "testType")
	expected := []byte(",\n")
	Assert(t, reflect.DeepEqual(expected, testWriter.WrittenValue), "Expected Emit to Write '%v'", expected)
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