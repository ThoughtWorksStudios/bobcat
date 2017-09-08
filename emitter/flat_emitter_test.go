package emitter

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"strings"
	"testing"
)

func TestFlatEmitter_Lifecycle(t *testing.T) {
	testWriter := &StringWriter{}
	emitter := NewFlatEmitter(testWriter)
	emitter.Init()

	emitter.Emit(EntityResult{"foo": 1}, "testType")
	emitter.Emit(EntityResult{"bar": 2}, "testType")

	emitter.Finalize()

	expected := []string{START, "{\n  \"foo\": 1\n}", DELIMITER, "{\n  \"bar\": 2\n}", END}
	AssertEqual(t, strings.Join(expected, ""), testWriter.String())
}

func TestFlatEmitter_Emit(t *testing.T) {
	testWriter := &StringWriter{}
	emitter := NewFlatEmitter(testWriter)

	testWriter.Reset()
	emitter.Emit(EntityResult{"foo": 1}, "testType")
	AssertEqual(t, "{\n  \"foo\": 1\n}", testWriter.String(), "Emit() produces JSON")

	testWriter.Reset()
	emitter.Emit(EntityResult{"bar": 2}, "testType")
	AssertEqual(t, DELIMITER+"{\n  \"bar\": 2\n}", testWriter.String(), "Emit() only inserts a DELIMITER from the second Emit() onward")
}

func TestFlatEmitter_NextEmitter(t *testing.T) {
	emitter := &FlatEmitter{}
	AssertEqual(t, emitter, emitter.NextEmitter(nil, "", false), "NextEmitter() should yield self")
}

func TestFlatEmitter_Receiver(t *testing.T) {
	emitter := &FlatEmitter{}
	AssertNil(t, emitter.Receiver(), "Receiver() should return nil")
}

func TestFlatEmitter_Finalize(t *testing.T) {
	writer := &StringWriter{}
	err := (&FlatEmitter{writer: writer}).Finalize()
	AssertNil(t, err, "Finalize() should not have thrown error, but did: %v", err)
	AssertEqual(t, END, writer.String(), "Finalize() should have written END token (i.e. %q)", END)
}

func TestFlatEmitter_Init(t *testing.T) {
	writer := &StringWriter{}
	err := (&FlatEmitter{writer: writer}).Init()
	AssertNil(t, err, "Init() should not have thrown error, but did: %v", err)
	AssertEqual(t, START, writer.String(), "Init() should have written START token (i.e. %q)", START)
}

func TestCreateFlatEmitter(t *testing.T) {
	emitter, err := FlatEmitterForFile(".")
	ExpectsError(t, "is a directory", err)
	AssertNil(t, emitter, "Should not have constructed an emitter on error")
}
