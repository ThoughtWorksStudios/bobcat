package emitter

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"strings"
	"testing"
)

func TestNestedEmitter_Lifecycle(t *testing.T) {
	testWriter := &StringWriter{}
	emitter, _ := InitNestedEmitter(testWriter)

	emitter.Emit(EntityResult{"bar": EntityResult{"baz": 1}}, "testType")
	nextEmitter := emitter.NextEmitter(emitter.Receiver(), "second", false)
	nextEmitter.Emit(EntityResult{"bars": EntityResult{"bazs": 1}}, "testType")

	emitter.Finalize()

	expected := []string{"{\n  \"\": {\n    \"bar\": {\n      \"baz\": 1\n    }\n  },\n  \"second\": {\n    \"bars\": {\n      \"bazs\": 1\n    }\n  }\n}"}
	AssertEqual(t, strings.Join(expected, ""), testWriter.String())
}

func TestNestedEmitter_Emit(t *testing.T) {
	emitter := &NestedEmitter{}
	emitter.Init()

	emitter.Emit(EntityResult{}, "testType")
	AssertDeepEqual(t, EntityResult{"": EntityResult{}}, emitter.Receiver(), "")

	nextEmitter := emitter.NextEmitter(EntityResult{}, "new_emitter", true)
	nextEmitter.Emit(EntityResult{"foo": "bar"}, "testType")
	nextEmitter.Emit(EntityResult{"bar": "foo"}, "testType")

	AssertDeepEqual(t, EntityResult{"new_emitter": []EntityResult{EntityResult{"foo": "bar"}, EntityResult{"bar":"foo"}}}, nextEmitter.Receiver())
}

func TestNestedEmitter_NextEmitter(t *testing.T) {
	emitter := &NestedEmitter{}
	actual := emitter.NextEmitter(EntityResult{"foo": "bar"}, "new_emitter", true)

	expected := &NestedEmitter{&Cursor{current: EntityResult{"foo":"bar"}, key: "new_emitter", isMultiValue: true}, nil, nil}
	AssertDeepEqual(t, expected, actual, "NextEmitter() should return a new emitter")
}

func TestNestedEmitter_Receiver(t *testing.T) {
	initialEmitter := &NestedEmitter{}
	initialEmitter.Init()
	currentEntity := EntityResult{"foo": "bar"}
	currentEmitter := initialEmitter.NextEmitter(currentEntity, "new_emitter", true)

	AssertDeepEqual(t, EntityResult{}, initialEmitter.Receiver())
	AssertDeepEqual(t, currentEntity, currentEmitter.Receiver(), "Receiver() should return current entity")
}

func TestNestedEmitter_Finalize(t *testing.T) {
	testWriter := &StringWriter{}
	emitter, _ := InitNestedEmitter(testWriter)
	testWriter.Reset()
	emitter.Emit(EntityResult{"foo": 1}, "testType")

	AssertEqual(t, "", testWriter.String())
	emitter.Finalize()

	AssertEqual(t, "{\n  \"\": {\n    \"foo\": 1\n  }\n}", testWriter.String(), "Finalize() produces JSON")
}

func TestNestedEmitter_Init(t *testing.T) {
	emitter := &NestedEmitter{}
	err := emitter.Init()
	AssertNil(t, err, "Init() should not have thrown error, but did: %v", err)
	AssertDeepEqual(t, EntityResult{}, emitter.Receiver(), "Init() should initialize a cursor")
}

func TestNewNestedEmitter(t *testing.T) {
	emitter, err := NewNestedEmitter(".")
	ExpectsError(t, "is a directory", err)
	AssertNil(t, emitter, "Should not have constructed an emitter on error")

	emitter, err = NewNestedEmitter("/dev/null")
	AssertNil(t, err, "NewNestedEmitter should not have thrown error, but did: %v", err)
}
