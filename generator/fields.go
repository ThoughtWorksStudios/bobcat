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
	result, err = f.value(parentId, emitter, scope)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (f *Field) value(parentId interface{}, emitter Emitter, scope *Scope) (result interface{}, err error) {
	if !f.count.Multiple() {
		result, err = f.fieldType.One(parentId, emitter, scope)
	} else {
		count := f.count.Count()
		values := make([]interface{}, count)

		for i := int64(0); i < count; i++ {
			if values[i], err = f.fieldType.One(parentId, emitter, scope); err != nil {
				break
			}
		}

		result = values
	}
	return result, err
}

func contains(sl []interface{}, value interface{}) bool {
	for _, a := range sl {
		if a == value {
			return true
		}
	}
	return false

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

func (field *EntityType) Type() string {
	return "entity"
}

func (field *EntityType) One(parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	g := field.entityGenerator
	if result, err := g.One(parentId, emitter, scope); err == nil {
		return result[g.PrimaryKeyName()], nil
	} else {
		return nil, err
	}
}

type LiteralType struct {
	value interface{}
}

func (field *LiteralType) Type() string {
	return "literal"
}

func (field *LiteralType) One(parentId interface{}, emitter Emitter, scope *Scope) (interface{}, error) {
	return field.value, nil
}
