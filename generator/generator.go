package generator

import (
	"fmt"
	. "github.com/ThoughtWorksStudios/bobcat/common"
	"github.com/ThoughtWorksStudios/bobcat/logging"
	"os"
	"sort"
	"strings"
	"time"
)

func debug(f string, t ...interface{}) {
	fmt.Fprintf(os.Stderr, f+"\n", t...)
}

type Generator struct {
	name   string
	base   string
	fields FieldSet
	log    logging.ILogger
}

func ExtendGenerator(name string, parent *Generator) *Generator {
	gen := NewGenerator(name, parent.log)
	gen.base = parent.Type()
	gen.fields["$extends"] = &LiteralField{value: gen.base}
	gen.fields["$type"] = &LiteralField{value: gen.Type()}

	for key, _ := range parent.fields {
		if _, hasField := gen.fields[key]; !hasField || !strings.HasPrefix(key, "$") {
			gen.fields[key] = &ReferenceField{referred: parent, fieldName: key}
		}
	}

	return gen
}

func NewGenerator(name string, logger logging.ILogger) *Generator {
	if logger == nil {
		logger = &logging.DefaultLogger{}
	}

	if name == "" {
		name = "$"
	}

	g := &Generator{name: name, fields: make(map[string]Field), log: logger}

	g.fields["$id"] = &UuidField{}

	g.fields["$type"] = &LiteralField{value: g.name}
	g.fields["$species"] = &LiteralField{value: g.name}

	return g
}

func (g *Generator) WithStaticField(fieldName string, fieldValue interface{}) error {
	if f, ok := g.fields[fieldName]; ok && f.Type() != "reference" {
		g.log.Warn("Field %s.%s is already defined; overriding to %v", g.name, fieldName, fieldValue)
	}

	g.fields[fieldName] = &LiteralField{value: fieldValue}
	return nil
}

func (g *Generator) WithEntityField(fieldName string, entityGenerator *Generator, fieldArgs interface{}, fieldBound Bound) error {
	if f, ok := g.fields[fieldName]; ok && f.Type() != "reference" {
		g.log.Warn("Field %s.%s is already defined; overriding.", g.name, fieldName)
	}

	g.fields[fieldName] = &EntityField{entityGenerator: entityGenerator, minBound: fieldBound.Min, maxBound: fieldBound.Max}
	return nil
}

func (g *Generator) WithField(fieldName, fieldType string, fieldArgs interface{}, fieldBound Bound) error {
	if fieldArgs == nil {
		return fmt.Errorf("FieldArgs are nil for field '%s', this should never happen!", fieldName)
	}

	if f, ok := g.fields[fieldName]; ok && f.Type() != "reference" {
		g.log.Warn("Field %s.%s is already defined; overriding to %s(%v)", g.name, fieldName, fieldType, fieldArgs)
	}

	switch fieldType {
	case "string":
		if ln, ok := fieldArgs.(int); ok {
			g.fields[fieldName] = &StringField{length: ln, minBound: fieldBound.Min, maxBound: fieldBound.Max}
		} else {
			return fmt.Errorf("expected field args to be of type 'int' for field %s (%s), but got %v",
				fieldName, fieldType, fieldArgs)
		}
	case "integer":
		if bounds, ok := fieldArgs.([2]int); ok {
			min, max := bounds[0], bounds[1]
			if max < min {
				return fmt.Errorf("max %v cannot be less than min %v", max, min)
			}

			g.fields[fieldName] = &IntegerField{min: min, max: max, minBound: fieldBound.Min, maxBound: fieldBound.Max}
		} else {
			return fmt.Errorf("expected field args to be of type '(min:int, max:int)' for field %s (%s), but got %v", fieldName, fieldType, fieldArgs)
		}
	case "decimal":
		if bounds, ok := fieldArgs.([2]float64); ok {
			min, max := bounds[0], bounds[1]
			if max < min {
				return fmt.Errorf("max %v cannot be less than min %v", max, min)
			}
			g.fields[fieldName] = &FloatField{min: min, max: max, minBound: fieldBound.Min, maxBound: fieldBound.Max}
		} else {
			return fmt.Errorf("expected field args to be of type '(min:float64, max:float64)' for field %s (%s), but got %v", fieldName, fieldType, fieldArgs)
		}
	case "date":
		if bounds, ok := fieldArgs.([2]time.Time); ok {
			min, max := bounds[0], bounds[1]
			field := &DateField{min: min, max: max, minBound: fieldBound.Min, maxBound: fieldBound.Max}
			if !field.ValidBounds() {
				return fmt.Errorf("max %v cannot be before min %v", max, min)
			}
			g.fields[fieldName] = field
		} else {
			return fmt.Errorf("expected field args to be of type 'time.Time' for field %s (%s), but got %v", fieldName, fieldType, fieldArgs)
		}
	case "uuid":
		g.fields[fieldName] = &UuidField{}
	case "dict":
		if dict, ok := fieldArgs.(string); ok {
			g.fields[fieldName] = &DictField{category: dict, minBound: fieldBound.Min, maxBound: fieldBound.Max}
		} else {
			return fmt.Errorf("expected field args to be of type 'string' for field %s (%s), but got %v", fieldName, fieldType, fieldArgs)
		}
	default:
		return fmt.Errorf("Invalid field type '%v'", fieldType)
	}

	return nil
}

func (g *Generator) Type() string {
	if (strings.HasPrefix(g.name, "$") || g.name == "") && g.base != "" {
		return g.base
	}
	return g.name
}

func (g *Generator) Generate(count int64) GeneratedEntities {
	entities := NewGeneratedEntities(count)
	for i := int64(0); i < count; i++ {
		entity := GeneratedFields{}
		for _, name := range sortKeys(g.fields) { // need $name fields generated first
			field := g.fields[name]
			if field.Type() == "entity" { // add reference to parent entity
				field.(*EntityField).entityGenerator.fields["$parent"] = &LiteralField{value: entity["$id"]}
			}
			values := []interface{}{}
			amount := field.Amount()
			if amount == 1 {
				entity[name] = field.GenerateValue()
			} else {
				for i := 0; i < amount; i++ {
					values = append(values, field.GenerateValue())
				}
				entity[name] = values
			}
		}
		entities[i] = entity
	}
	return entities
}

func (g *Generator) String() string {
	return fmt.Sprintf("%s{}", g.name)
}

func sortKeys(fields FieldSet) []string {
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
