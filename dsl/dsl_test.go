package dsl

import (
	. "github.com/ThoughtWorksStudios/datagen/test_helpers"
	"testing"
)

func RequiresDefOrGenerateStatements(t *testing.T) {
	_, err := Parse("", []byte("eek"))
	expectedErrorMsg := "1:1 (0): no match found, expected: \"def\", \"generate\", [ \t\r\n] or EOF"
	ExpectsError(t, expectedErrorMsg, err)
}

func TestParsesBasicEntity(t *testing.T) {
	location := NewLocation("", 1, 1, 0)
	kids := NodeSet{Node{Kind: "definition", Ref: location, Name: "Bird", Children: NodeSet{}}}
	expected := Node{Kind: "root", Ref: location, Children: kids}
	actual, err := Parse("", []byte("def Bird {  }"))
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, expected.String(), actual.(Node).String())
}
