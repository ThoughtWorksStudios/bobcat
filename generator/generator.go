package generator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/ThoughtWorksStudios/datagen/logging"
	"io"
	"os"
	"time"
)

func debug(f string, t ...interface{}) {
	fmt.Fprintf(os.Stderr, f+"\n", t...)
}

type FieldSet map[string]Field

type Generator struct {
	name   string
	parent *Generator
	fields FieldSet
	log    logging.ILogger
}

func ExtendGenerator(name string, parent *Generator) *Generator {
	gen := NewGenerator(name, parent.log)
	gen.parent = parent

	for key, _ := range parent.fields {
		gen.fields[key] = &ReferenceField{referred: parent, fieldName: key}
	}

	return gen
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

// Also for testing purposes
func (g *Generator) GetName() string {
	return g.name
}

func (g *Generator) WithStaticField(fieldName string, fieldValue interface{}) error {
	if f, ok := g.fields[fieldName]; ok && f.Type() != "reference" {
		g.log.Warn("Field %s.%s is already defined; overriding to %v", g.name, fieldName, fieldValue)
	}

	g.fields[fieldName] = &LiteralField{value: fieldValue}
	return nil
}

func (g *Generator) WithField(fieldName, fieldType string, fieldOpts interface{}) error {
	if fieldOpts == nil {
		return fmt.Errorf("FieldOpts are nil for field '%s', this should never happen!", fieldName)
	}

	if f, ok := g.fields[fieldName]; ok && f.Type() != "reference" {
		g.log.Warn("Field %s.%s is already defined; overriding to %s(%v)", g.name, fieldName, fieldType, fieldOpts)
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

func isClosable(v interface{}) (io.Closer, bool) {
	closeable, doClose := v.(io.Closer)
	return closeable, doClose
}

func (g *Generator) writeJsonToStream(v interface{}, out io.Writer) error {
	if out == nil {
		var f interface{}
		var err error
		if f, err = os.Create(fmt.Sprintf("%s.json", g.name)); err != nil {
			return err
		} else {
			out, _ = f.(io.Writer)
		}
	}

	if closeable, doClose := isClosable(out); doClose {
		defer closeable.Close()
	}

	writer := bufio.NewWriter(out)
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "\t")

	if err := encoder.Encode(v); err != nil {
		return err
	}

	return writer.Flush()
}

func (g *Generator) Generate(count int64, out io.Writer) error {
	result := make([]map[string]interface{}, count)
	for i := int64(0); i < count; i++ {

		obj := make(map[string]interface{})
		for name, field := range g.fields {
			obj[name] = field.GenerateValue()
		}
		result[i] = obj
	}

	return g.writeJsonToStream(result, out)
}
