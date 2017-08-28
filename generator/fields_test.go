package generator

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"testing"
)

func TestGenerateEntity(t *testing.T) {
	g := NewGenerator("testEntity", false)
	fieldType := &EntityType{g}
	emitter := NewTestEmitter()
	subId := fieldType.One("", emitter)

	e := emitter.Shift()

	if nil == e {
		t.Errorf("Expected to generate an entity but got %T %v", e, e)
	}

	AssertEqual(t, "testEntity", e["$type"], "Should have generated an entity of type \"testEntity\"")
	AssertEqual(t, subId, e["$id"])
}

func TestGenerateFloat(t *testing.T) {
	min, max := 4.25, 4.3
	FieldType := &FloatType{min, max}
	actual := FieldType.One("", NewTestEmitter()).(float64)

	if actual < min || actual > max {
		t.Errorf("Generated value '%v' is outside of expected range min: '%v', max: '%v'", actual, min, max)
	}
}

func TestGenerateEnum(t *testing.T) {
	args := []interface{}{"one", "two", "three"}
	FieldType := &EnumType{values: args, size: int64(len(args))}
	actual := FieldType.One("", NewTestEmitter()).(string)

	if actual != "one" && actual != "two" && actual != "three" {
		t.Errorf("Generated value '%v' enum value list: %v", actual, args)
	}
}

func TestMultiValueGenerate(t *testing.T) {
	field := NewField(&IntegerType{1, 10}, &CountRange{3, 3})
	actual := len(field.GenerateValue("", NewTestEmitter()).([]interface{}))

	AssertEqual(t, 3, actual)
}
