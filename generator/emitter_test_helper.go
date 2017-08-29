package generator

import . "github.com/ThoughtWorksStudios/bobcat/common"

type TestEmitter struct {
	result []EntityResult
}

func (te *TestEmitter) Receiver() EntityResult {
	return te.result[0]
}

func (te *TestEmitter) Emit(entity EntityResult, entityType string) error {
	te.result = append(te.result, entity)
	return nil
}

func (te *TestEmitter) NextEmitter(current EntityResult, key string, isMultiValue bool) Emitter {
	return te
}

func (te *TestEmitter) Finalize() error {
	return nil
}

func (te *TestEmitter) Shift() EntityResult {
	entity := te.result[0]
	te.result = te.result[1:]
	return entity
}

func NewTestEmitter() *TestEmitter {
	return &TestEmitter{}
}
