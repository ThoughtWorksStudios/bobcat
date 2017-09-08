package emitter

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"github.com/json-iterator/go"
	"strings"
	"testing"
)

func TestNestedEmitter_Lifecycle(t *testing.T) {
	testWriter := &StringWriter{}
	emitter := NewNestedEmitter(testWriter)
	emitter.Init()

	emitter.Emit(EntityResult{"bar": EntityResult{"baz": 1}}, "testType")
	nextEmitter := emitter.NextEmitter(emitter.Receiver(), "second", false)
	nextEmitter.Emit(EntityResult{"bars": EntityResult{"bazs": 1}}, "testType")

	emitter.Finalize()

	expected := []string{"{\n  \"\": {\n    \"bar\": {\n      \"baz\": 1\n    }\n  },\n  \"second\": {\n    \"bars\": {\n      \"bazs\": 1\n    }\n  }\n}"}

	expectedTree := jsoniter.UnmarshalFromString(strings.Join(expected, ""), make(EntityResult))
	actualTree := jsoniter.UnmarshalFromString(testWriter.String(), make(EntityResult))
	AssertDeepEqual(t, expectedTree, actualTree)
}

func TestNestedEmitter_Emit(t *testing.T) {
	testWriter := &StringWriter{}
	emitter := NewNestedEmitter(testWriter)
	emitter.Init()

	testWriter.Reset()
	emitter.Emit(EntityResult{"foo": 1}, "testType")
	AssertEqual(t, "{\n  \"foo\": 1\n}", testWriter.String(), "Emit() produces JSON")

	testWriter.Reset()
	emitter.Emit(EntityResult{"bar": 2}, "testType")
	AssertEqual(t, DELIMITER+"{\n  \"bar\": 2\n}", testWriter.String(), "Emit() only inserts a DELIMITER from the second Emit() onward")

	testWriter.Reset()
	emitter.Emit(EntityResult{"baz": EntityResult{"quu": 3}}, "testType")
	AssertEqual(t, DELIMITER+"{\n  \"baz\": {\n    \"quu\": 3\n  }\n}", testWriter.String(), "Emit() maintains nesting")
}

func TestNestedEmitter_NextEmitter(t *testing.T) {
	testWriter := &StringWriter{}
	emitter := NewNestedEmitter(testWriter)

	actual := emitter.NextEmitter(EntityResult{"foo": "bar"}, "new_emitter", true)

	expected := &NestedEmitter{&Cursor{current: EntityResult{"foo": "bar"}, key: "new_emitter", isMultiValue: true}, nil}
	AssertDeepEqual(t, expected, actual, "NextEmitter() should return a new emitter")
}

func TestNestedEmitter_Receiver(t *testing.T) {
	initialEmitter := NewNestedEmitter(&StringWriter{})
	initialEmitter.Init()
	currentEntity := EntityResult{"foo": "bar"}
	currentEmitter := initialEmitter.NextEmitter(currentEntity, "new_emitter", true)

	AssertDeepEqual(t, currentEntity, currentEmitter.Receiver(), "Receiver() should return current entity")
}

func TestNestedEmitter_Finalize(t *testing.T) {
	testWriter := &StringWriter{}
	emitter := NewNestedEmitter(testWriter)
	err := emitter.Finalize()
	AssertNil(t, err, "Finalize() should not have thrown error, but did: %v", err)
	AssertEqual(t, END, testWriter.String(), "Finalize() should have written END token (i.e. %q)", END)
}

func TestNestedEmitter_Init(t *testing.T) {
	testWriter := &StringWriter{}
	emitter := NewNestedEmitter(testWriter)

	err := emitter.Init()
	AssertNil(t, err, "Init() should not have thrown error, but did: %v", err)

	_, ok := emitter.Receiver().(*EntityOutputter)
	Assert(t, ok, "Top-level Receiver() yields an EntityOutputter from its Cursor")
}

func TestNewNestedEmitter(t *testing.T) {
	emitter, err := NestedEmitterForFile(".")
	ExpectsError(t, "is a directory", err)
	AssertNil(t, emitter, "Should not have constructed an emitter on error")

	emitter, err = NestedEmitterForFile("/dev/null")
	AssertNil(t, err, "NewNestedEmitter should not have thrown error, but did: %v", err)
}
