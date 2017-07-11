package interpreter

import "github.com/ThoughtWorksStudios/datagen/dsl"
import "github.com/ThoughtWorksStudios/datagen/generator"
import "fmt"
import "strconv"

func assignFields(entity *generator.Generator, fields []dsl.Node) {
	for _, field := range fields {
		valueType := field.Value.(string)
		args := field.Args[0]
		var parsedArgs interface{}
		switch args.Kind {
		case "string":
			parsedArgs = args.Value.(string)
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
		case "numeric":
			i, _ := strconv.Atoi(args.Value.(string))
			parsedArgs = i
		}
		entity.WithField(field.Name, valueType, parsedArgs)
	}
}

func Generate(tree dsl.Node) {
	entities := make(map[string]*generator.Generator)
	for _, node := range tree.Children {
		if node.Kind == "definition" {
			entity := generator.NewGenerator(node.Name)
			assignFields(entity, node.Children)
			entities[node.Name] = entity
		} else if node.Kind == "generation" {
			fmt.Println(entities)
			fmt.Println(node.Name)
			count, e := strconv.Atoi(node.Args[0].Value.(string))
			if e == nil {
				entities[node.Name].Generate(count)
			} else {
				fmt.Println(e)
			}
		}
	}
}
