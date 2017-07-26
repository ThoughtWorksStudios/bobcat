package generator

import (
	"fmt"
	"github.com/ThoughtWorksStudios/datagen/logging"
	"os"
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
	g.fields["$id"] = &UuidField{}
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

	g.fields[fieldName] = &LiteralField{value: fieldValue}
	return nil
}

func (g *Generator) WithField(fieldName, fieldType string, fieldOpts interface{}) error {
	if fieldOpts == nil {
		return fmt.Errorf("FieldOpts are nil for field '%s', this should never happen!", fieldName)
	}

	if f, ok := g.fields[fieldName]; ok && f.Type() != "reference" {
		g.log.Warn("Field %s.%s is already defined; overriding to %s(%v)", g.Name, fieldName, fieldType, fieldOpts)
	}

	switch fieldType {
	case "string":
		if ln, ok := fieldOpts.(int); ok {
			g.fields[fieldName] = &StringField{length: ln}
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

			g.fields[fieldName] = &IntegerField{min: min, max: max}
		} else {
			return fmt.Errorf("expected field options to be of type '(min:int, max:int)' for field %s (%s), but got %v", fieldName, fieldType, fieldOpts)
		}
	case "decimal":
		if bounds, ok := fieldOpts.([2]float64); ok {
			min, max := bounds[0], bounds[1]
			if max < min {
				return fmt.Errorf("max %v cannot be less than min %v", max, min)
			}
			g.fields[fieldName] = &FloatField{min: min, max: max}
		} else {
			return fmt.Errorf("expected field options to be of type '(min:float64, max:float64)' for field %s (%s), but got %v", fieldName, fieldType, fieldOpts)
		}
	case "date":
		if bounds, ok := fieldOpts.([2]time.Time); ok {
			min, max := bounds[0], bounds[1]
			field := &DateField{min: min, max: max}
			if !field.ValidBounds() {
				return fmt.Errorf("max %v cannot be before min %v", max, min)
			}
			g.fields[fieldName] = field
		} else {
			return fmt.Errorf("expected field options to be of type 'time.Time' for field %s (%s), but got %v", fieldName, fieldType, fieldOpts)
		}
	case "uuid":
		g.fields[fieldName] = &UuidField{}
	case "dict":
		if dict, ok := fieldOpts.(string); ok {
			g.fields[fieldName] = &DictField{category: dict}
		} else {
			return fmt.Errorf("expected field options to be of type 'string' for field %s (%s), but got %v", fieldName, fieldType, fieldOpts)
		}
	default:
		return fmt.Errorf("Invalid field type '%v'", fieldType)
	}

	return nil
}

func (g *Generator) Generate(count int64) GeneratedContent {
	result := NewGeneratedContent()
	entities := make([]map[string]interface{}, count)
	for i := int64(0); i < count; i++ {

		obj := make(map[string]interface{})
		for name, field := range g.fields {
			obj[name] = field.GenerateValue()
		}
		entities[i] = obj
	}
	result[g.Name] = entities

	return result
}
