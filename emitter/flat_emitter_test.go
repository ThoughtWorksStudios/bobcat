package emitter

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"reflect"
	"testing"
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
