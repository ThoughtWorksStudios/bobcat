package interpreter

import (
	"github.com/ThoughtWorksStudios/datagen/dsl"
	"github.com/ThoughtWorksStudios/datagen/generator"
	"log"
	"testing"
	"time"
)

func builtin(value string) dsl.Node {
	return dsl.Node{Kind: "bullitin", Value: value}
}

func stringArgs(values ...string) []dsl.Node {
	i, size := 0, len(values)
	args := make([]dsl.Node, size)

	for _, val := range values {
		args[i] = dsl.Node{Kind: "literal-string", Value: val}
		i = i + 1
	}

	return args
}

func intArgs(values ...int64) []dsl.Node {
	i, size := 0, len(values)
	args := make([]dsl.Node, size)

	for _, val := range values {
		args[i] = dsl.Node{Kind: "literal-int", Value: val}
		i = i + 1
	}

	return args
}

func floatArgs(values ...float64) []dsl.Node {
	i, size := 0, len(values)
	args := make([]dsl.Node, size)

	for _, val := range values {
		args[i] = dsl.Node{Kind: "literal-int", Value: val}
		i = i + 1
	}

	return args
}

func dateArgs(values ...string) []dsl.Node {
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

func rootNode(nodes ...dsl.Node) dsl.Node {
	return dsl.Node{Name: "root", Children: nodes}
}

func generationNode(entityName string, count int64) dsl.Node {
	return dsl.Node{Kind: "generation", Name: entityName, Args: intArgs(count)}
}

func newEntity(name string, fields []dsl.Node) dsl.Node {
	return dsl.Node{Name: name, Kind: "definition", Children: fields}
}

func assertShouldHaveField(t *testing.T, entity *generator.Generator, field dsl.Node) {
	if entity.GetField(field.Name) == nil {
		t.Errorf("Expected entity to have field %s, but it did not", field.Name)
	}
}

func assertExpectedEqsActual(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Errorf("expected %v, but was %v", expected, actual)
	}
}
