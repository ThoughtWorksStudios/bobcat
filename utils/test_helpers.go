package utils

import (
	"github.com/ThoughtWorksStudios/datagen/dsl"
	"github.com/ThoughtWorksStudios/datagen/generator"
	"testing"
)

func AssertShouldHaveField(t *testing.T, entity *generator.Generator, field dsl.Node) {
	if entity.GetField(field.Name) == nil {
		t.Errorf("Expected entity to have field %s, but it did not", field.Name)
	}
}

func AssertExpectedEqsActual(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Errorf("expected %v, but was %v", expected, actual)
	}
}
