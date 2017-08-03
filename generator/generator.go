package generator

import (
	"fmt"
	"github.com/ThoughtWorksStudios/datagen/logging"
	"os"
	"sort"
	"strings"
	"time"
)

func debug(f string, t ...interface{}) {
	fmt.Fprintf(os.Stderr, f+"\n", t...)
}

type Generator struct {
	Name   string
	Base   string
	parent *Generator
	fields FieldSet
	log    logging.ILogger
}

func ExtendGenerator(name string, parent *Generator) *Generator {
	gen := NewGenerator(name, parent.log)
	gen.parent = parent

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

	g := &Generator{Name: name, fields: make(map[string]Field), log: logger}
	g.fields["$id"] = &UuidField{count: Range{1, 1}}
	return g
}

// For testing purposes
func (g *Generator) GetField(name string) Field {
	return g.fields[name]
}

func (g *Generator) WithStaticField(fieldName string, fieldValue interface{}) error {
	if f, ok := g.fields[fieldName]; ok && f.Type() != "reference" {
		g.log.Warn("Field %s.%s is already defined; overriding to %v", g.Name, fieldName, fieldValue)
	}

	g.fields[fieldName] = &LiteralField{value: fieldValue, count: Range{1, 1}}
	return nil
}

func (g *Generator) WithEntityField(fieldName string, entityGenerator *Generator, fieldValueCount Range) error {
	if f, ok := g.fields[fieldName]; ok && f.Type() != "reference" {
		g.log.Warn("Field %s.%s is already defined; overriding.", g.Name, fieldName)
	}

	g.fields[fieldName] = &EntityField{entityGenerator: entityGenerator, count: fieldValueCount}
	return nil
}

func (g *Generator) WithField(fieldName, fieldType string, fieldOpts interface{}, fieldValueCount Range) error {
	if fieldOpts == nil {
		return fmt.Errorf("FieldOpts are nil for field '%s', this should never happen!", fieldName)
	}

	if f, ok := g.fields[fieldName]; ok && f.Type() != "reference" {
		g.log.Warn("Field %s.%s is already defined; overriding to %s(%v)", g.Name, fieldName, fieldType, fieldOpts)
	}

	switch fieldType {
	case "string":
		if ln, ok := fieldOpts.(int); ok {
			g.fields[fieldName] = &StringField{length: ln, count: fieldValueCount}
		} else {
			return fmt.Errorf("expected field options to be of type 'int' for field %s (%s), but got %v",
				fieldName, fieldType, fieldOpts)
		}
	case "integer":
		if bounds, ok := fieldOpts.([2]int); ok {
			min, max := bounds[0], bounds[1]
			if max < min {
				return fmt.Errorf("max %v cannot be less than min %v", max, min)
			}

			g.fields[fieldName] = &IntegerField{min: min, max: max, count: fieldValueCount}
		} else {
			return fmt.Errorf("expected field options to be of type '(min:int, max:int)' for field %s (%s), but got %v", fieldName, fieldType, fieldOpts)
		}
	case "decimal":
		if bounds, ok := fieldOpts.([2]float64); ok {
			min, max := bounds[0], bounds[1]
			if max < min {
				return fmt.Errorf("max %v cannot be less than min %v", max, min)
			}
			g.fields[fieldName] = &FloatField{min: min, max: max, count: fieldValueCount}
		} else {
			return fmt.Errorf("expected field options to be of type '(min:float64, max:float64)' for field %s (%s), but got %v", fieldName, fieldType, fieldOpts)
		}
	case "date":
		if bounds, ok := fieldOpts.([2]time.Time); ok {
			min, max := bounds[0], bounds[1]
			field := &DateField{min: min, max: max, count: fieldValueCount}
			if !field.ValidBounds() {
				return fmt.Errorf("max %v cannot be before min %v", max, min)
			}
			g.fields[fieldName] = field
		} else {
			return fmt.Errorf("expected field options to be of type 'time.Time' for field %s (%s), but got %v", fieldName, fieldType, fieldOpts)
		}
	case "uuid":
		g.fields[fieldName] = &UuidField{count: Range{1, 1}}
	case "dict":
		if dict, ok := fieldOpts.(string); ok {
			g.fields[fieldName] = &DictField{category: dict, count: fieldValueCount}
		} else {
			return fmt.Errorf("expected field options to be of type 'string' for field %s (%s), but got %v", fieldName, fieldType, fieldOpts)
		}
	default:
		return fmt.Errorf("Invalid field type '%v'", fieldType)
	}

	return nil
}

func (g *Generator) Generate(count int64) GeneratedEntities {
	entities := NewGeneratedEntities(count)
	for i := int64(0); i < count; i++ {
		entity := GeneratedFields{}
		for _, name := range sortKeys(g.fields) { // need $name fields generated first
			field := g.fields[name]
			if field.Type() == "entity" { // add reference to parent entity
				field.(*EntityField).entityGenerator.fields["$parent"] = &LiteralField{value: entity["$id"], count: Range{1,1}}
			}

			fieldCount := field.Count()
			if fieldCount > 1 {
				fieldValue := []interface{}{}
				for j := 0; j < fieldCount; j++ {
					fieldValue = append(fieldValue, field.GenerateValue())
				}
				entity[name] = fieldValue
			} else {
				entity[name] = field.GenerateValue()
			}
			if field.Type() == "entity" {
				labeledEntity := make(map[string]interface{})
				entityName := field.(*EntityField).entityGenerator.Name
				labeledEntity[entityName] = entity[name]
				entity[name] = labeledEntity
			}
		}
		entities[i] = entity
	}
	return entities
}

func sortKeys(fields FieldSet) []string {
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
