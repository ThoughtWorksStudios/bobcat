package generator

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"testing"
)

func TestGenerateEntity(t *testing.T) {
	g := NewGenerator("testEntity", GetLogger(t))
	fieldType := &EntityType{g}
	e := fieldType.GenerateSingle()

	if kind, ok := e.(EntityResult); !ok {
		t.Errorf("Expected to generate an entity but got %v", kind)
	}
}

func TestGenerateFloat(t *testing.T) {
	min, max := 4.25, 4.3
	FieldType := &FloatType{min, max}
	actual := FieldType.GenerateSingle().(float64)

	if actual < min || actual > max {
		t.Errorf("Generated value '%v' is outside of expected range min: '%v', max: '%v'", actual, min, max)
	}
}

func TestGenerateEnum(t *testing.T) {
	args := []interface{}{"one", "two", "three"}
	FieldType := &EnumType{values: args}
	actual := FieldType.GenerateSingle().(string)

	if actual != "one" && actual != "two" && actual != "three" {
		t.Errorf("Generated value '%v' enum value list: %v", actual, args)
	}
}

func TestMultiValueGenerate(t *testing.T) {
	field := NewField(&IntegerType{1, 10}, &CountRange{3, 3})
	actual := len(field.GenerateValue().([]interface{}))

	AssertEqual(t, 3, actual)
}
