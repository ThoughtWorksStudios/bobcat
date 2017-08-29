package emitter

import "fmt"

type EmitterMap map[string]Emitter

func (em EmitterMap) FetchOrCreate(key string, provider EmitterProvider) (Emitter, error) {
	if "" == key {
		return nil, fmt.Errorf("Cannot fetch emitter without key")
	}

	if emitter, isPresent := em[key]; isPresent {
		return emitter, nil
	}

	if emitter, err := provider.Get(key); err == nil {
		em[key] = emitter
		return emitter, nil
	} else {
		return nil, err
	}
}

func (em EmitterMap) Finalize() error {
	for _, emitter := range em {
		if err := emitter.Finalize(); err != nil {
			return err
		}
	}

	return nil
}
