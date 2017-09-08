package common

import "fmt"

type EntityStore interface {
	Set(key string, entity EntityResult) error
	AppendTo(key string, entity EntityResult) error
}

type EntityResult map[string]interface{}

func (er EntityResult) Set(key string, entity EntityResult) error {
	er[key] = entity
	return nil
}

func (er EntityResult) AppendTo(key string, entity EntityResult) error {
	var result []EntityResult
	var ok bool

	if original, isPresent := er[key]; !isPresent {
		result = make([]EntityResult, 0)
	} else {
		if result, ok = original.([]EntityResult); !ok {
			return fmt.Errorf("Expected an entity set")
		}
	}

	er[key] = append(result, entity)

	return nil
}
