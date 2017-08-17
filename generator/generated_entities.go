package generator

import (
	"bytes"
	"github.com/json-iterator/go"
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

var (
	IntegerEncoder GeneratedIntegerValue
	StringEncoder  GeneratedStringValue
	BoolEncoder    GeneratedBoolValue
	FloatEncoder   GeneratedFloatValue
	ListEncoder    GeneratedListValue
	EntityEncoder  GeneratedEntityValue
)

type GeneratedValue interface {
	Encode(ptr unsafe.Pointer, stream *jsoniter.Stream)
	EncodeInterface(val interface{}, stream *jsoniter.Stream)
	IsEmpty(ptr unsafe.Pointer) bool
}

type GeneratedStringValue string

var StringBuffer bytes.Buffer

func (v GeneratedStringValue) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	StringBuffer.Reset()
	StringBuffer.WriteByte('"')
	// StringBuffer.WriteString(string(v))
	StringBuffer.WriteString((*(*string)(ptr)))
	StringBuffer.WriteByte('"')

	stream.Write(StringBuffer.Bytes())
}

func (v GeneratedStringValue) EncodeInterface(val interface{}, stream *jsoniter.Stream) {
	// v, _ = val.(GeneratedStringValue)
	jsoniter.WriteToStream(val, stream, v)
}

func (v GeneratedStringValue) IsEmpty(ptr unsafe.Pointer) bool {
	return v == ""
}

type GeneratedIntegerValue int

func (v GeneratedIntegerValue) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	stream.WriteInt((*(*int)(ptr)))
}

func (v GeneratedIntegerValue) EncodeInterface(val interface{}, stream *jsoniter.Stream) {
	jsoniter.WriteToStream(val, stream, v)
}

func (v GeneratedIntegerValue) IsEmpty(ptr unsafe.Pointer) bool {
	return *(*int)(ptr) == 0
}

type GeneratedFloatValue float64

func (v GeneratedFloatValue) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	stream.WriteFloat64((*(*float64)(ptr)))
}

func (v GeneratedFloatValue) EncodeInterface(val interface{}, stream *jsoniter.Stream) {
	jsoniter.WriteToStream(val, stream, v)
}

func (v GeneratedFloatValue) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

type GeneratedListValue []GeneratedValue

func (v GeneratedListValue) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
}

func (v GeneratedListValue) EncodeInterface(val interface{}, stream *jsoniter.Stream) {
	v, _ = val.(GeneratedListValue)
	jsoniter.WriteToStream(v, stream, v)
}

func (v GeneratedListValue) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

type GeneratedBoolValue bool

func (v GeneratedBoolValue) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	stream.WriteBool(bool(v))
}

func (v GeneratedBoolValue) EncodeInterface(val interface{}, stream *jsoniter.Stream) {
	v, _ = val.(GeneratedBoolValue)
	jsoniter.WriteToStream(v, stream, v)
}

func (v GeneratedBoolValue) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

type GeneratedEntityValue EntityResult

func (v GeneratedEntityValue) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
}

func (v GeneratedEntityValue) EncodeInterface(val interface{}, stream *jsoniter.Stream) {
	v, _ = val.(GeneratedEntityValue)
	jsoniter.WriteToStream(v, stream, v)
}

func (v GeneratedEntityValue) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}
