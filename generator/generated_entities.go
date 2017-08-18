package generator

import "github.com/json-iterator/go"

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
	Encoder() jsoniter.ValEncoder
}

type GeneratedStringValue string

func (s GeneratedStringValue) Encoder() jsoniter.ValEncoder {
	return &StringEncoder{}
}

type GeneratedIntegerValue int

func (s GeneratedIntegerValue) Encoder() jsoniter.ValEncoder {
	return &IntegerEncoder{}
}

type GeneratedFloatValue float64

func (s GeneratedFloatValue) Encoder() jsoniter.ValEncoder {
	return &FloatEncoder{}
}

type GeneratedListValue []GeneratedValue

func (s GeneratedListValue) Encoder() jsoniter.ValEncoder {
	return &ListEncoder{}
}

type GeneratedBoolValue bool

func (s GeneratedBoolValue) Encoder() jsoniter.ValEncoder {
	return &BoolEncoder{}
}

type GeneratedEntityValue EntityResult

func (s GeneratedEntityValue) Encoder() jsoniter.ValEncoder {
	return &EntityEncoder{}
}
