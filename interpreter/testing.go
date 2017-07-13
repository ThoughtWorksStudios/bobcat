package interpreter

import "testing"
import "time"
import "github.com/ThoughtWorksStudios/datagen/dsl"
import "github.com/ThoughtWorksStudios/datagen/generator"

func dictArgs(value string) []dsl.Node {
	return []dsl.Node{dictArg(value)}
}

func timeArgs(min, max string) []dsl.Node {
	minTime, _ := time.Parse("20016-01-02", min)
	maxTime, _ := time.Parse("20016-01-02", max)
	return []dsl.Node{timeArg(minTime), timeArg(maxTime)}
}

func stringArgs(value int64) []dsl.Node {
	return []dsl.Node{stringArg(value)}
}

func intArgs(min, max int64) []dsl.Node {
	return []dsl.Node{stringArg(min), stringArg(max)}
}

func floatArgs(min, max float64) []dsl.Node {
	return []dsl.Node{floatArg(min), floatArg(max)}
}

func stringArg(value int64) dsl.Node {
	return dsl.Node{Kind: "literal-int", Value: value}
}

func intArg(value int64) dsl.Node {
	return dsl.Node{Kind: "literal-int", Value: value}
}

func dictArg(value string) dsl.Node {
	return dsl.Node{Kind: "dict", Value: value}
}

func floatArg(value float64) dsl.Node {
	return dsl.Node{Kind: "decimal", Value: value}
}

func timeArg(value time.Time) dsl.Node {
	return dsl.Node{Kind: "date", Value: value}
}

func rootNode(nodes ...dsl.Node) dsl.Node {
	return dsl.Node{Name: "root", Children: nodes}
}

func generationNode(entityName string, count int64) dsl.Node {
	return dsl.Node{Kind: "generation", Name: entityName, Args: []dsl.Node{intArg(count)}}
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
