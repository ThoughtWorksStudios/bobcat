package generator

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
)

type TestEmitter struct {
	result EntityResult
}

func (te *TestEmitter) Receiver() EntityResult {
	return te.result
}
func (te *TestEmitter) Emit(entity EntityResult) error {
	return nil
}
func (te *TestEmitter) NextEmitter(current EntityResult, key string, isMultiValue bool) Emitter {
	return nil
}
func (te *TestEmitter) Finalize() error {
	return nil
}

func testEmitter() Emitter {
	return &TestEmitter{}
}
