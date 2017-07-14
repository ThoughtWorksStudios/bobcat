package test_helpers

import (
	"github.com/ThoughtWorksStudios/datagen/dsl"
	"log"
	"time"
)

func Builtin(value string) dsl.Node {
	return dsl.Node{Kind: "builtin", Value: value}
}

func StaticNode(value interface{}) dsl.Node {
	return dsl.Node{Kind: "static", Value: value}
}

func StringArgs(values ...string) []dsl.Node {
	i, size := 0, len(values)
	args := make([]dsl.Node, size)

	for _, val := range values {
		args[i] = dsl.Node{Kind: "literal-string", Value: val}
		i = i + 1
	}

	return args
}

func IntArgs(values ...int64) []dsl.Node {
	i, size := 0, len(values)
	args := make([]dsl.Node, size)

	for _, val := range values {
		args[i] = dsl.Node{Kind: "literal-int", Value: val}
		i = i + 1
	}

	return args
}

func FloatArgs(values ...float64) []dsl.Node {
	i, size := 0, len(values)
	args := make([]dsl.Node, size)

	for _, val := range values {
		args[i] = dsl.Node{Kind: "literal-int", Value: val}
		i = i + 1
	}

	return args
}

func DateArgs(values ...string) []dsl.Node {
	i, size := 0, len(values)
	args := make([]dsl.Node, size)

	for _, val := range values {
		parsed, err := time.Parse("2006-01-02", val)
		if err != nil {
			log.Fatalf("could not parse %v against YYYY-mm-dd. Error: %v", val, err)
		}

		args[i] = dsl.Node{Kind: "literal-int", Value: parsed}
		i = i + 1
	}

	return args
}

func RootNode(nodes ...dsl.Node) dsl.Node {
	return dsl.Node{Name: "root", Children: nodes}
}

func GenerationNode(entityName string, count int64) dsl.Node {
	return dsl.Node{Kind: "generation", Name: entityName, Args: IntArgs(count)}
}

func EntityNode(name string, fields []dsl.Node) dsl.Node {
	return dsl.Node{Name: name, Kind: "definition", Children: fields}
}
