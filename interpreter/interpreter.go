package interpreter

import (
	"fmt"
	"github.com/ThoughtWorksStudios/datagen/dsl"
	"github.com/ThoughtWorksStudios/datagen/generator"
	"log"
	"time"
)

var die = func(msg string, args ...interface{}) {
	log.Fatalln("FATAL:", fmt.Sprintf(msg, args...))
}

func defaultArgumentFor(fieldType string) interface{} {
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
		die("Field of type `%s` requires arguments", fieldType)
	}

	return arg
}

func translateFieldsForEntity(entity *generator.Generator, fields []dsl.Node) {
	for _, field := range fields {
		configureFieldOn(entity, field)
	}
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

func configureFieldOn(entity *generator.Generator, field dsl.Node) {
	declType := field.Value.(dsl.Node).Kind

	if declType == "builtin" {
		withDynamicField(entity, field)
	} else {
		withStaticField(entity, field)
	}
}

func withStaticField(entity *generator.Generator, field dsl.Node) {
	fieldValue := field.Value.(dsl.Node).Value
	entity.WithStaticField(field.Name, fieldValue)
}

func withDynamicField(entity *generator.Generator, field dsl.Node) {
	fieldType, ok := field.Value.(dsl.Node).Value.(string)
	if !ok {
		die("Could not parse field-type for field `%s`. Expected one of the builtin generator types, but instead got: %v", field.Name, field.Value.(dsl.Node).Value)
	}
	numArgs := len(field.Args)

	if 0 == numArgs {
		entity.WithField(field.Name, fieldType, defaultArgumentFor(fieldType))
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
			die("field type `dict` requires exactly 1 argument")
		}
	case "date":
		if numArgs == 2 {
			entity.WithField(field.Name, fieldType, [2]time.Time{valTime(field.Args[0]), valTime(field.Args[1])})
		}
	}
}

func translateEntity(node dsl.Node) *generator.Generator {
	entity := generator.NewGenerator(node.Name)
	translateFieldsForEntity(entity, node.Children)
	return entity
}

func translateEntities(tree dsl.Node) map[string]*generator.Generator {
	entities := make(map[string]*generator.Generator)
	for _, node := range tree.Children {
		if node.Kind == "definition" {
			entities[node.Name] = translateEntity(node)
		}
	}
	return entities
}

func generateEntities(tree dsl.Node, entities map[string]*generator.Generator) error {
	for _, node := range tree.Children {
		if node.Kind == "generation" {
			count, e := node.Args[0].Value.(int64)
			entity, exists := entities[node.Name]

			if e {
				if count <= int64(1) {
					return fmt.Errorf("ERROR: Must generate at least 1 `%s` entity", node.Name)
				} else if !exists {
					return fmt.Errorf("ERROR: %s is undefined; expected entity", node.Name)
				} else {
					entity.Generate(count)
				}
			} else {
				return fmt.Errorf("ERROR: generate %s takes an integer count", node.Name)
			}
		}
	}
	return nil
}

func Translate(tree dsl.Node) error {
	entities := translateEntities(tree)
	return generateEntities(tree, entities)
}
