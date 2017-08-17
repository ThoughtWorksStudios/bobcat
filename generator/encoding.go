package generator

import (
	"github.com/json-iterator/go"
	"unsafe"
)

type ValueEncoder struct{}

func (encoder *ValueEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
}

func (encoder *ValueEncoder) EncodeInterface(val interface{}, stream *jsoniter.Stream) {
	if s, ok := val.(GeneratedStringValue); ok {
		jsoniter.WriteToStream(s, stream, StringEncoder{})
	} else if s, ok := val.(GeneratedIntegerValue); ok {
		jsoniter.WriteToStream(s, stream, IntegerEncoder{})
	}
}

func (encoder *ValueEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

type StringEncoder struct{}

func (encoder StringEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	stream.WriteString(*((*string)(ptr)))
}

func (encoder StringEncoder) EncodeInterface(val interface{}, stream *jsoniter.Stream) {
}

func (encoder StringEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return (*((*string)(ptr))) == ""
}

type IntegerEncoder struct{}

func (encoder IntegerEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	stream.WriteInt((*(*int)(ptr)))
}

func (encoder IntegerEncoder) EncodeInterface(val interface{}, stream *jsoniter.Stream) {
}

func (encoder IntegerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *(*int)(ptr) == 0
}

type FloatEncoder struct{}

func (v FloatEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	stream.WriteFloat64((*(*float64)(ptr)))
}

func (v FloatEncoder) EncodeInterface(val interface{}, stream *jsoniter.Stream) {
}

func (v FloatEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

type ListEncoder struct{}

func (v ListEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
}

func (v ListEncoder) EncodeInterface(val interface{}, stream *jsoniter.Stream) {
}

func (v ListEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

type BoolEncoder struct{}

func (v BoolEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	stream.WriteBool(*(*bool)(ptr))
}

func (v BoolEncoder) EncodeInterface(val interface{}, stream *jsoniter.Stream) {
}

func (v BoolEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

type EntityEncoder struct{}

func (v EntityEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
}

func (v EntityEncoder) EncodeInterface(val interface{}, stream *jsoniter.Stream) {
}

func (v EntityEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}
