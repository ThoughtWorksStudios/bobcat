package generator

import (
	. "github.com/ThoughtWorksStudios/bobcat/builtins"
	. "github.com/ThoughtWorksStudios/bobcat/common"
)

type PrimaryKey struct {
	name string
	kind string
}

func (pk *PrimaryKey) Field() *Field {
	var builtin Callable

	switch pk.kind {
	case SERIAL_TYPE:
		builtin, _ = NewBuiltin(SERIAL_TYPE)
	default:
		builtin, _ = NewBuiltin(UID_TYPE)
	}
	return NewDeferredField(func(_ *Scope) (interface{}, error) { return builtin.Call() })
}

func (pk *PrimaryKey) Inherit(target *Generator, source *Generator) {
	target.fields.AddField(pk.name, &Field{fieldType: &ReferenceType{referred: source.fields, fieldName: source.PrimaryKeyName()}})
	target.pkey = pk
}

func (pk *PrimaryKey) Attach(target *Generator) {
	target.fields.AddField(pk.name, pk.Field())
	target.pkey = pk
}

func NewPrimaryKeyConfig(name, kind string) *PrimaryKey {
	return &PrimaryKey{name: name, kind: kind}
}
