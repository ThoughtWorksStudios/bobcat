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

func TestMultiValueGenerate(t *testing.T) {
	field := NewField(&IntegerType{1, 10}, &CountRange{3, 3})
	actual := len(field.GenerateValue().([]interface{}))

	AssertEqual(t, 3, actual)
}
