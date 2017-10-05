package generator

import (
	"fmt"
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/emitter"
)

type Field struct {
	fieldType FieldType
	count     *CountRange
}

func (f *Field) MultiValue() bool {
	return f.count.Multiple()
}

func (f *Field) Type() string {
	return f.fieldType.Type()
}

func (f *Field) underlyingType() string {
	ft := f.fieldType

	for ft.Type() == "reference" {
		rt := ft.(*ReferenceType)
		ft = rt.referred.GetField(rt.fieldName).fieldType
	}

	return ft.Type()
}

func (f Field) String() string {
	return fmt.Sprintf(`{ type: %q, underlying: %q, multiVal: %t`, f.Type(), f.underlyingType(), f.MultiValue())
}

func (f *Field) GenerateValue(parentId interface{}, emitter Emitter, scope *Scope) (result interface{}, err error) {
	if !f.count.Multiple() {
		result, err = f.One(parentId, emitter, scope)
	} else {
		count := f.count.Count()
		values := make([]interface{}, count)

		for i := int64(0); i < count; i++ {
			if values[i], err = f.One(parentId, emitter, scope); err != nil {
				result = nil
				break
			}
		}
		result = values
	}

	return
}

func (f *Field) One(parentId interface{}, emitter Emitter, scope *Scope) (result interface{}, err error) {
	if result, err = f.fieldType.One(parentId, emitter, scope); err != nil {
		return nil, err
	} else {
		if entity, ok := result.(*Generator); ok {
			var val EntityResult
			if val, err = entity.One(parentId, emitter, scope); err != nil {
				return nil, err
			}
			result = val[entity.PrimaryKeyName()]
		}
	}

	return
}

type FieldType interface {
	Type() string
	One(parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error)
}

func NewDeferredField(closure DeferredResolver) *Field {
	return NewField(NewDeferredType(closure), nil)
}

func NewLiteralField(fieldValue interface{}) *Field {
	return NewField(NewLiteralType(fieldValue), nil)
}

func NewLiteralType(fieldValue interface{}) *LiteralType {
	return &LiteralType{value: fieldValue}
}
func NewDeferredType(closure DeferredResolver) *DeferredType {
	return &DeferredType{closure: closure}
}

func NewField(fieldType FieldType, count *CountRange) *Field {
	return &Field{fieldType: fieldType, count: count}
}

type DeferredType struct {
	closure DeferredResolver
}

func (f *DeferredType) Type() string {
	return "deferred"
}

func (f *DeferredType) One(parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	return f.closure(scope)
}

type ReferenceType struct {
	referred  *FieldSet
	fieldName string
}

func (f *ReferenceType) Type() string {
	return "reference"
}

func (f *ReferenceType) One(parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	ref := f.referred.GetField(f.fieldName).fieldType
	return ref.One(parentId, emitter, scope)
}

type EntityType struct {
	entityGenerator *Generator
}

func (f *EntityType) Type() string {
	return "entity"
}

func (f *EntityType) One(parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	entity := f.entityGenerator

	if result, err := entity.One(parentId, emitter, scope); err == nil {
		return result[entity.PrimaryKeyName()], nil
	} else {
		return nil, err
	}
}

type LiteralType struct {
	value interface{}
}

func (f *LiteralType) Type() string {
	return "literal"
}

func (f *LiteralType) One(parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	return f.value, nil
}
