package generator

import (
	"fmt"
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/emitter"
	"strings"
)

var DEFAULT_PK_CONFIG *PrimaryKey = &PrimaryKey{name: "$id", kind: UID_TYPE}

type Generator struct {
	name            string
	extends         string
	declaredType    string
	fields          *FieldSet
	disableMetadata bool
	pkey            *PrimaryKey
}

func ExtendGenerator(name string, parent *Generator, pkey *PrimaryKey, disableMetadata bool) *Generator {
	var gen *Generator

	if pkey == nil {
		gen = NewGenerator(name, parent.pkey, disableMetadata)
		parent.pkey.Inherit(gen, parent) // must reference parent's primary key; this is important for serial fields to continue the sequence
	} else {
		gen = NewGenerator(name, pkey, disableMetadata)
	}

	gen.extends = parent.Type()

	gen.recalculateType()

	if !disableMetadata {
		gen.fields.AddField("$extends", NewField(&LiteralType{value: gen.extends}, nil))
		gen.fields.AddField("$type", NewField(&LiteralType{value: gen.Type()}, nil))
	}

	for _, fieldEntry := range parent.fields.AllFields() {
		key := fieldEntry.Name
		f := fieldEntry.Field
		if !gen.fields.HasField(key) && !strings.HasPrefix(key, "$") && key != parent.PrimaryKeyName() {
			gen.fields.AddField(key, NewField(&ReferenceType{referred: parent.fields, fieldName: key}, f.count))
		}
	}

	return gen
}

func NewGenerator(name string, pkey *PrimaryKey, disableMetadata bool) *Generator {
	if name == "" {
		name = "$"
	}

	g := &Generator{name: name, fields: NewFieldSet(), disableMetadata: disableMetadata}

	g.recalculateType()

	if pkey == nil {
		pkey = DEFAULT_PK_CONFIG
	}

	pkey.Attach(g)

	if !disableMetadata {
		g.fields.AddField("$type", NewField(&LiteralType{value: g.Type()}, nil))
	}

	return g
}

func (g *Generator) HasField(name string) bool {
	return g.fields.HasField(name)
}

func (g *Generator) GetField(name string) *Field {
	return g.fields.GetField(name)
}

func (g *Generator) PrimaryKeyName() string {
	return g.pkey.name
}

func (g *Generator) WithDeferredField(fieldName string, fieldValue DeferredResolver) error {
	g.fields.AddField(fieldName, NewDeferredField(fieldValue))
	return nil
}

func (g *Generator) WithLiteralField(fieldName string, fieldValue interface{}) error {
	g.fields.AddField(fieldName, NewLiteralField(fieldValue))
	return nil
}

func (g *Generator) WithEntityField(fieldName string, entityGenerator *Generator, countRange *CountRange) error {
	g.fields.AddField(fieldName, NewField(&EntityType{entityGenerator: entityGenerator}, countRange))
	return nil
}

func (g *Generator) WithField(fieldName string, fieldType FieldType, count *CountRange) {
	g.fields.AddField(fieldName, NewField(fieldType, count))
}

func (g *Generator) Type() string {
	return g.declaredType
}

func (g *Generator) recalculateType() {
	if (strings.HasPrefix(g.name, "$") || g.name == "") && g.extends != "" {
		g.declaredType = g.extends
	} else {
		g.declaredType = g.name
	}
}

func (g *Generator) Generate(count int64, emitter Emitter, scope *Scope) ([]interface{}, error) {
	ids := make([]interface{}, count)
	idKey := g.PrimaryKeyName()

	for i := int64(0); i < count; i++ {
		if r, err := g.One(nil, emitter, scope); err == nil {
			ids[i] = r[idKey]
		} else {
			return nil, err
		}
	}
	return ids, nil
}

func (g *Generator) One(parentId interface{}, emitter Emitter, scope *Scope) (EntityResult, error) {
	entity := EntityResult{}
	// Need to extend TransientScope once more as a protective layer so that any symbols declared
	// by expressions are NOT set as fields in the final entity result
	childScope := ExtendScope(TransientScope(scope, SymbolTable(entity)))

	idKey := g.PrimaryKeyName()
	id, err := g.fields.GetField(idKey).GenerateValue(nil, emitter, childScope)

	if err != nil {
		return nil, err
	}

	entity[idKey] = id // create this first because we may use it as the parentId when generating child entities

	if parentId != nil {
		entity["$parent"] = parentId
	}

	for _, entry := range g.fields.AllFields() {
		name, field := entry.Name, entry.Field
		if name != idKey { // don't GenerateValue() more than once for id -- it throws off the sequence for serial types
			// GenerateValue MAY populate entity[name] for entity fields
			value, err := field.GenerateValue(id, emitter.NextEmitter(entity, name, field.MultiValue()), childScope)
			if err != nil {
				return nil, err
			}
			if _, isAlreadySet := entity[name]; !isAlreadySet {
				entity[name] = value
			}
		}
	}

	if err = emitter.Emit(entity, g.Type()); err != nil {
		return nil, err
	}

	return entity, nil
}

func (g *Generator) String() string {
	return fmt.Sprintf("%s{}", g.name)
}
