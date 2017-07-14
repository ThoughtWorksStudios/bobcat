package interpreter

import (
	"fmt"
	"github.com/ThoughtWorksStudios/datagen/dsl"
	"github.com/ThoughtWorksStudios/datagen/generator"
	"github.com/ThoughtWorksStudios/datagen/logging"
	"log"
	"time"
)

type Interpreter struct {
	l logging.ILogger
}

func New(logger logging.ILogger) *Interpreter {
	if logger == nil {
		logger = &logging.DefaultLogger{}
	}

	return &Interpreter{l: logger}
}

func (i *Interpreter) defaultArgumentFor(fieldType string) interface{} {
	var arg interface{}

	switch fieldType {
	case "string":
		arg = 5
	case "integer":
		arg = [2]int{1, 10}
	case "decimal":
		arg = [2]float64{1, 10}
	case "date":
		t1, _ := time.Parse("2006-01-02", "1945-01-01")
		t2, _ := time.Parse("2006-01-02", "2017-01-01")
		arg = [2]time.Time{t1, t2}
	default:
		i.l.Die("Field of type `%s` requires arguments", fieldType)
	}

	return arg
}

func (i *Interpreter) translateFieldsForEntity(entity *generator.Generator, fields []dsl.Node) *generator.Generator {
	for _, field := range fields {
		declType := field.Value.(dsl.Node).Kind

		if declType == "builtin" {
			i.withDynamicField(entity, field)
		} else {
			i.withStaticField(entity, field)
		}
	}

	return entity
}

func valStr(n dsl.Node) string {
	return n.Value.(string)
}

func valInt(n dsl.Node) int {
	return int(n.Value.(int64))
}

func valFloat(n dsl.Node) float64 {
	return n.Value.(float64)
}

func valTime(n dsl.Node) time.Time {
	return n.Value.(time.Time)
}

func err(msg string, tokens ...interface{}) error {
	return fmt.Errorf("ERROR: "+msg, tokens...)
}

func (i *Interpreter) withStaticField(entity *generator.Generator, field dsl.Node) {
	fieldValue := field.Value.(dsl.Node).Value
	entity.WithStaticField(field.Name, fieldValue)
}

func (i *Interpreter) withDynamicField(entity *generator.Generator, field dsl.Node) {
	fieldType, ok := field.Value.(dsl.Node).Value.(string)
	if !ok {
		i.l.Die("Could not parse field-type for field `%s`. Expected one of the builtin generator types, but instead got: %v", field.Name, field.Value.(dsl.Node).Value)
	}
	numArgs := len(field.Args)

	if 0 == numArgs {
		entity.WithField(field.Name, fieldType, i.defaultArgumentFor(fieldType))
		return
	}

	switch fieldType {
	case "integer":
		if numArgs == 2 {
			entity.WithField(field.Name, fieldType, [2]int{valInt(field.Args[0]), valInt(field.Args[1])})
		}
	case "decimal":
		if numArgs == 2 {
			entity.WithField(field.Name, fieldType, [2]float64{valFloat(field.Args[0]), valFloat(field.Args[1])})
		}
	case "string":
		if numArgs == 1 {
			entity.WithField(field.Name, fieldType, valInt(field.Args[0]))
		}
	case "dict":
		if numArgs == 1 {
			entity.WithField(field.Name, fieldType, valStr(field.Args[0]))
		} else {
			i.l.Die("field type `dict` requires exactly 1 argument")
		}
	case "date":
		if numArgs == 2 {
			entity.WithField(field.Name, fieldType, [2]time.Time{valTime(field.Args[0]), valTime(field.Args[1])})
		}
	}
}

func (i *Interpreter) translateEntities(tree dsl.Node) map[string]*generator.Generator {
	entities := make(map[string]*generator.Generator)
	for _, node := range tree.Children {
		if node.Kind == "definition" {
			entities[node.Name] = i.translateFieldsForEntity(generator.NewGenerator(node.Name), node.Children)
		}
	}
	return entities
}

func (i *Interpreter) generateEntities(tree dsl.Node, entities map[string]*generator.Generator) error {
	for _, node := range tree.Children {
		if node.Kind == "generation" {
			count, e := node.Args[0].Value.(int64)
			entity, exists := entities[node.Name]

			if e {
				if count <= int64(1) {
					return err("Must generate at least 1 `%s` entity", node.Name)
				} else if !exists {
					return err("%s is undefined; expected entity", node.Name)
				} else {
					entity.Generate(count)
				}
			} else {
				return err("generate %s takes an integer count", node.Name)
			}
		}
	}
	return nil
}

func (i *Interpreter) Consume(tree dsl.Node) error {
	return i.generateEntities(tree, i.translateEntities(tree))
}
