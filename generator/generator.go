package generator

import (
	"encoding/json"
	"fmt"
	"github.com/ThoughtWorksStudios/datagen/logging"
	"os"
	"time"
)

type Generator struct {
	name   string
	fields map[string]Field
	log    logging.ILogger
}

func NewGenerator(name string, logger logging.ILogger) *Generator {
	if logger == nil {
		logger = &logging.DefaultLogger{}
	}
	return &Generator{name: name, fields: make(map[string]Field), log: logger}
}

// For testing purposes
func (g *Generator) GetField(name string) Field {
	return g.fields[name]
}

func (g *Generator) WithStaticField(fieldName string, fieldValue interface{}) {
	if _, ok := g.fields[fieldName]; ok {
		g.log.Warn("already defined field: %s", fieldName)
	}

	g.fields[fieldName] = &LiteralField{value: fieldValue}
}

func (g *Generator) WithField(fieldName, fieldType string, fieldOpts interface{}) {
	if fieldOpts == nil {
		g.log.Die("FieldOpts are nil for field '%s', this should never happen!", fieldName)
	}

	if _, ok := g.fields[fieldName]; ok {
		g.log.Warn("already defined field: %s", fieldName)
	}

	switch fieldType {
	case "string":
		if ln, ok := fieldOpts.(int); ok {
			g.fields[fieldName] = &StringField{length: ln}
		} else {
			g.log.Die("expected field options to be of type 'int' for field %s (%s), but got %v",
				fieldName, fieldType, fieldOpts)
		}
	case "integer":
		if bounds, ok := fieldOpts.([2]int); ok {
			min, max := bounds[0], bounds[1]
			if max < min {
				g.log.Die("max %d cannot be less than min %d\n", max, min)
			}

			g.fields[fieldName] = &IntegerField{min: min, max: max}
		} else {
			g.log.Die("expected field options to be of type '(min:int, max:int)' for field %s (%s), but got %v", fieldName, fieldType, fieldOpts)
		}
	case "decimal":
		if bounds, ok := fieldOpts.([2]float64); ok {
			min, max := bounds[0], bounds[1]
			if max < min {
				g.log.Die("max %d cannot be less than min %d\n", max, min)
			}
			g.fields[fieldName] = &FloatField{min: min, max: max}
		} else {
			g.log.Die("expected field options to be of type '(min:float64, max:float64)' for field %s (%s), but got %v", fieldName, fieldType, fieldOpts)
		}
	case "date":
		if bounds, ok := fieldOpts.([2]time.Time); ok {
			min, max := bounds[0], bounds[1]
			field := &DateField{min: min, max: max}
			if !field.ValidBounds() {
				g.log.Die("max %s cannot be before min %s\n", max, min)
			}
			g.fields[fieldName] = field
		} else {
			g.log.Die("expected field options to be of type 'time.Time' for field %s (%s), but got %v", fieldName, fieldType, fieldOpts)
		}
	case "dict":
		if dict, ok := fieldOpts.(string); ok {
			g.fields[fieldName] = &DictField{category: dict}
		} else {
			g.log.Die("expected field options to be of type 'string' for field %s (%s), but got %v", fieldName, fieldType, fieldOpts)
		}
	default:
		g.log.Die("Invalid field type '%s'", fieldType)
	}
}

func (g *Generator) Generate(count int64) {
	result := make([]map[string]interface{}, count)
	for i := int64(0); i < count; i++ {

		obj := make(map[string]interface{})
		for name, field := range g.fields {
			obj[name] = field.GenerateValue()
		}
		result[i] = obj
	}

	marsh, _ := json.MarshalIndent(result, "", "\t")
	writeToFile(marsh, fmt.Sprintf("%s.json", g.name))
}

var writeToFile = func(json []byte, filename string) {
	dest, _ := os.Create(filename)
	defer dest.Close()
	dest.Write(json)
}
