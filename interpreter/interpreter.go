package interpreter

import "github.com/ThoughtWorksStudios/datagen/dsl"
import "github.com/ThoughtWorksStudios/datagen/generator"
import "fmt"
import "strconv"

func parseRangeArg(args dsl.Node, fieldValueType string) interface{} {
	rng := args.Args
	min := rng[0].Value.(string)
	max := rng[1].Value.(string)
	var arg interface{}
	switch fieldValueType {
	case "decimal":
		minValue, _ := strconv.ParseFloat(min, 10)
		maxValue, _ := strconv.ParseFloat(max, 10)
		arg = [2]float64{minValue, maxValue}
	case "integer":
		minValue, _ := strconv.Atoi(min)
		maxValue, _ := strconv.Atoi(max)
		arg = [2]int{minValue, maxValue}
	case "date":
		arg = [2]string{min, max}
	}
	return arg
}

func parseNumericArg(args dsl.Node, fieldType string) interface{} {
	var i interface{}
	switch fieldType {
	case "integer":
		i, _ = strconv.Atoi(args.Value.(string))
	case "decimal":
		i, _ = strconv.ParseFloat(args.Value.(string), 10)
	case "string":
		i, _ = strconv.Atoi(args.Value.(string))
	}
	return i
}

func parseArguments(args dsl.Node, fieldType string) interface{} {
	var arg interface{}
	switch args.Kind {
	case "string":
		arg = args.Value.(string)
	case "numeric":
		arg = parseNumericArg(args, fieldType)
	case "range":
		arg = parseRangeArg(args, fieldType)
	}
	return arg
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
		arg = [2]string{"1945-01-01", "2017-01-01"}
	case "dict":
		arg = "silly_name"
	}
	return arg
}

func translateFieldsForEntity(entity *generator.Generator, fields []dsl.Node) {
	for _, field := range fields {
		fieldType := field.Value.(string)
		var parsedArgs interface{}
		if len(field.Args) > 0 {
			parsedArgs = parseArguments(field.Args[0], fieldType)
		} else {
			parsedArgs = defaultArgumentFor(fieldType)
		}
		entity.WithField(field.Name, fieldType, parsedArgs)
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

func Translate(tree dsl.Node) {
	entities := translateEntities(tree)
	for _, node := range tree.Children {
		if node.Kind == "generation" {
			count, e := strconv.Atoi(node.Args[0].Value.(string))
			entity, exists := entities[node.Name]
			if e == nil && exists {
				entity.Generate(count)
			} else {
				fmt.Println(e)
			}
		}
	}
}
