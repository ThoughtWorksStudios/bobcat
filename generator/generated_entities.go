package generator

import (
	"github.com/json-iterator/go"
	"strconv"
	"unsafe"
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
	Encode(ptr unsafe.Pointer, stream *jsoniter.Stream)
	EncodeInterface(val interface{}, stream *jsoniter.Stream)
	IsEmpty(ptr unsafe.Pointer) bool
}

type GenericGeneratedValue struct {
	Value interface{}
}

// func (v *GenericGeneratedValue) MarshalJSON() ([]byte, error) {
// 	return ffjson.Marshal(v.Value)
// }

func (v *GenericGeneratedValue) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

func (v *GenericGeneratedValue) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
}

func (v *GenericGeneratedValue) EncodeInterface(val interface{}, stream *jsoniter.Stream) {
	jsoniter.WriteToStream(val, stream, v)
}

type GeneratedStringValue struct {
	Value string
}

func (v *GeneratedStringValue) MarshalJSON() ([]byte, error) {
	return []byte("\"" + v.Value + "\""), nil
}

func (v *GeneratedStringValue) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	stream.WriteString((*(*string)(ptr)))
}

func (v *GeneratedStringValue) EncodeToInterface(val interface{}, stream *jsoniter.Stream) {
	jsoniter.WriteToStream(val, stream, v)
}

func (v *GeneratedStringValue) EncodeInterface(val interface{}, stream *jsoniter.Stream) {
	jsoniter.WriteToStream(val, stream, v)
}

func (v *GeneratedStringValue) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

type GeneratedIntegerValue struct {
	Value int64
}

func (v *GeneratedIntegerValue) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(v.Value, 10)), nil
}

func (v *GeneratedIntegerValue) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	stream.WriteInt64(*(*int64)(ptr))
}

func (v *GeneratedIntegerValue) EncodeInterface(val interface{}, stream *jsoniter.Stream) {
	jsoniter.WriteToStream(val, stream, v)
}

func (v *GeneratedIntegerValue) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

type GeneratedLiteralValue struct {
	Value interface{}
}

func (v *GeneratedLiteralValue) EncodeInterface(val interface{}, stream *jsoniter.Stream) {
	jsoniter.WriteToStream(val, stream, v)
}

func (v *GeneratedLiteralValue) MarshalJSON() ([]byte, error) {
	if s, ok := v.Value.(string); ok {
		return []byte("\"" + s + "\""), nil
	} else {
		return jsoniter.Marshal(v.Value)
	}
}

func (v *GeneratedLiteralValue) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	stream.WriteString((*(*string)(ptr)))
}

func (v *GeneratedLiteralValue) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}
