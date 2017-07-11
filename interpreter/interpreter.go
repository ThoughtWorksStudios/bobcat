package interpreter

import "github.com/ThoughtWorksStudios/datagen/dsl"
import "github.com/ThoughtWorksStudios/datagen/generator"
import "fmt"
import "strconv"

func translateFieldsForEntity(entity *generator.Generator, fields []dsl.Node) {
	for _, field := range fields {
		valueType := field.Value.(string)
		args := field.Args[0]
		var parsedArgs interface{}
		switch args.Kind {
		case "string":
			parsedArgs = args.Value.(string)
		case "numeric":
			i, _ := strconv.Atoi(args.Value.(string))
			parsedArgs = i
		case "range":
			rng := args.Args
			min := rng[0].Value.(string)
			max := rng[1].Value.(string)
			if valueType == "decimal" {
				minValue, _ := strconv.ParseFloat(min, 10)
				maxValue, _ := strconv.ParseFloat(max, 10)
				parsedArgs = [2]float64{minValue, maxValue}
			} else if valueType == "integer" {
				minValue, _ := strconv.Atoi(min)
				maxValue, _ := strconv.Atoi(max)
				parsedArgs = [2]int{minValue, maxValue}
			} else if valueType == "date" {
				parsedArgs = [2]string{min, max}
			}
		}
		entity.WithField(field.Name, valueType, parsedArgs)
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
