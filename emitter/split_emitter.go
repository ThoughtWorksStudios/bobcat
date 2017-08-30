package emitter

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
)

type SplitEmitter struct {
	emitters EmitterMap
	provider EmitterProvider
}

func NewSplitEmitter(filenameTemplate string) (Emitter, error) {
	if provider, err := NewPerTypeEmitterProvider(filenameTemplate); err != nil {
		return nil, err
	} else {
		return &SplitEmitter{provider: provider, emitters: make(EmitterMap)}, nil
	}
}

func (se *SplitEmitter) Emit(entity EntityResult, entityType string) error {
	if emitter, err := se.emitters.FetchOrCreate(entityType, se.provider); err == nil {
		return emitter.Emit(entity, entityType)
	} else {
		return err
	}
}

func (se *SplitEmitter) Finalize() error {
	return se.emitters.Finalize()
}

func (se *SplitEmitter) NextEmitter(current EntityResult, key string, isMultiValue bool) Emitter {
	return se
}

func (se *SplitEmitter) Receiver() EntityResult {
	return nil
}

func (se *SplitEmitter) Init() error {
	return nil
}
