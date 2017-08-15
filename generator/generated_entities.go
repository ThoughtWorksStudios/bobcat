package generator

import (
	"github.com/pquerna/ffjson/ffjson"
	"strconv"
)

type GeneratedEntities []EntityResult

type EntityResult map[string]GeneratedValue

func NewGeneratedEntities(count int64) GeneratedEntities {
	return make([]EntityResult, count)
}

func (ge GeneratedEntities) Concat(newEntities GeneratedEntities) GeneratedEntities {
	for _, entity := range newEntities {
		ge = append(ge, entity)
	}
	return ge
}

type GeneratedValue interface {
	MarshalJSON() ([]byte, error)
}

type GenericGeneratedValue struct {
	Value interface{}
}

func (v GenericGeneratedValue) MarshalJSON() ([]byte, error) {
	return ffjson.Marshal(v.Value)
}

type GeneratedStringValue struct {
	Value string
}

func (v GeneratedStringValue) MarshalJSON() ([]byte, error) {
	return []byte("\"" + v.Value + "\""), nil
}

type GeneratedIntegerValue struct {
	Value int64
}

func (v GeneratedIntegerValue) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(v.Value, 10)), nil
}

type GeneratedEntityValue struct {
	Value EntityResult
}

func (v GeneratedEntityValue) MarshalJSON() ([]byte, error) {
	return ffjson.Marshal(v.Value)
}

type GeneratedLiteralValue struct {
	Value interface{}
}

func (v GeneratedLiteralValue) MarshalJSON() ([]byte, error) {
	if s, ok := v.Value.(string); ok {
		return GeneratedStringValue{Value: s}.MarshalJSON()
	} else {
		return ffjson.Marshal(v.Value)
	}
}
