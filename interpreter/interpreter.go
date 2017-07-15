package interpreter

import (
	"fmt"
	"github.com/ThoughtWorksStudios/datagen/dsl"
	"github.com/ThoughtWorksStudios/datagen/generator"
	"github.com/ThoughtWorksStudios/datagen/logging"
	"time"
)

type Interpreter struct {
	entities map[string]*generator.Generator // TODO: should probably be a more generic symbol table or possibly the parent scope
	l        logging.ILogger
}

type NodeSet []dsl.Node
type Iterator func(index int, node dsl.Node)
type Collector func(index int, node dsl.Node) interface{}

func (nodes NodeSet) Each(f Iterator) NodeSet {
	for i, size := 0, len(nodes); i < size; i++ {
		f(i, nodes[i])
	}
	return nodes
}

func (nodes NodeSet) Map(f Collector) []interface{} {
	size := len(nodes)
	result := make([]interface{}, size)
	nodes.Each(func(index int, node dsl.Node) {
		result[index] = f(index, node)
	})
	return result
}

func New(logger logging.ILogger) *Interpreter {
	if logger == nil {
		logger = &logging.DefaultLogger{}
	}

	return &Interpreter{l: logger, entities: make(map[string]*generator.Generator)}
}

func (i *Interpreter) Visit(node dsl.Node) error {
	switch node.Kind {
	case "root":
		NodeSet(node.Children).Each(func(_ int, node dsl.Node) {
			i.Visit(node)
		})
		return nil
	case "definition":
		i.EntityFromNode(node)
		return nil
	case "generation":
		return i.GenerateFromNode(node)
	}

	return nil
}

func (i *Interpreter) defaultArgumentFor(fieldType string) interface{} {
	switch fieldType {
	case "string":
		return 5
	case "integer":
		return [2]int{1, 10}
	case "decimal":
		return [2]float64{1, 10}
	case "date":
		t1, _ := time.Parse("2006-01-02", "1945-01-01")
		t2, _ := time.Parse("2006-01-02", "2017-01-01")
		return [2]time.Time{t1, t2}
	default:
		i.l.Die("Field of type `%s` requires arguments", fieldType)
	}

	return nil
}

func (i *Interpreter) EntityFromNode(node dsl.Node) *generator.Generator {
	entity, fields := generator.NewGenerator(node.Name), node.Children

	for _, field := range fields {
		if field.Kind != "field" {
			i.l.Die("Expected a field declaration, but instead got %v", field)
		}

		declType := field.Value.(dsl.Node).Kind

		if declType == "builtin" {
			i.withDynamicField(entity, field)
		} else {
			i.withStaticField(entity, field)
		}
	}

	i.entities[node.Name] = entity
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
			i.l.Die("Field type `dict` requires exactly 1 argument")
		}
	case "date":
		if numArgs == 2 {
			entity.WithField(field.Name, fieldType, [2]time.Time{valTime(field.Args[0]), valTime(field.Args[1])})
		}
	}
}

func (i *Interpreter) GenerateFromNode(node dsl.Node) error {
	count, ok := node.Args[0].Value.(int64)
	entity, exists := i.entities[node.Name]

	if !ok {
		return err("generate %s takes an integer count", node.Name)
	}

	if count <= int64(1) {
		return err("Must generate at least 1 `%s` entity", node.Name)
	}

	if !exists {
		return err("Unknown symbol `%s` -- expected an entity. Did you mean to define an entity named `%s`?", node.Name, node.Name)
	}

	entity.Generate(count)
	return nil
}
