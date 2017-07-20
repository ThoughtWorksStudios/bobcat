package dsl

import (
	. "github.com/ThoughtWorksStudios/datagen/test_helpers"
	"testing"
)

var genBird Node = Node{
	Kind:     "generation",
	Name:     "Bird",
	Ref:      NewLocation("", 1, 1, 0),
	Args:     NodeSet{},
	Children: NodeSet{},
}
var birdField Node = Node{
	Kind: "field",
	Name: "name",
	Args: NodeSet{},
	Ref:  NewLocation("", 1, 1, 0),
}

var bird Node = Node{
	Kind:     "definition",
	Ref:      NewLocation("", 1, 1, 0),
	Name:     "Bird",
	Children: NodeSet{},
}

var root Node = Node{
	Kind: "root",
	Ref:  NewLocation("", 1, 1, 0),
}

func RequiresDefOrGenerateStatements(t *testing.T) {
	_, err := Parse("", []byte("eek"))
	expectedErrorMsg := "1:1 (0): no match found, expected: \"def\", \"generate\", [ \t\r\n] or EOF"
	ExpectsError(t, expectedErrorMsg, err)
}

func TestParsesBasicEntity(t *testing.T) {
	root.Children = NodeSet{bird}
	actual, err := Parse("", []byte("def Bird {  }"))
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, root.String(), actual.(Node).String())
}

func TestParsesBasicGenerationStatement(t *testing.T) {
	genBird.Args = NodeSet{Node{Kind: "literal-int", Value: 1, Ref: NewLocation("", 1, 15, 14)}}
	root.Children = NodeSet{genBird}
	actual, err := Parse("", []byte("generate Bird(1)"))
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, root.String(), actual.(Node).String())
}

func TestParsedBothBasicEntityAndGenerationStatement(t *testing.T) {
	genBird.Args = NodeSet{Node{Kind: "literal-int", Value: 1, Ref: NewLocation("", 2, 15, 26)}}
	genBird.Ref = NewLocation("", 2, 1, 12)
	root.Children = NodeSet{bird, genBird}
	actual, err := Parse("", []byte("def Bird {}\ngenerate Bird(1)"))
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, root.String(), actual.(Node).String())
}

func TestParseEntityWithField(t *testing.T) {
	field := Node{
		Kind:  "field",
		Name:  "name",
		Value: Node{Kind: "builtin", Ref: NewLocation("", 1, 17, 16), Value: "string"},
		Args:  NodeSet{},
		Ref:   NewLocation("", 1, 12, 11),
	}
	bird.Children = NodeSet{field}
	root.Children = NodeSet{bird}
	actual, err := Parse("", []byte("def Bird { name string }"))
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, root.String(), actual.(Node).String())
}

func TestParseEntityWithFieldWithArgs(t *testing.T) {
	birdField.Value = Node{Kind: "builtin", Ref: NewLocation("", 1, 17, 16), Value: "string"}
	birdField.Args = NodeSet{Node{Kind: "literal-int", Value: 1, Ref: NewLocation("", 1, 24, 23)}}
	birdField.Ref = NewLocation("", 1, 12, 11)
	bird.Children = NodeSet{birdField}
	root.Children = NodeSet{bird}
	actual, err := Parse("", []byte("def Bird { name string(1) }"))
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, root.String(), actual.(Node).String())
}

func TestParseEntitywithFieldWithMultipleArgs(t *testing.T) {
	birdField.Value = Node{Kind: "builtin", Ref: NewLocation("", 1, 17, 16), Value: "integer"}
	arg1 := Node{Kind: "literal-int", Value: 1, Ref: NewLocation("", 1, 25, 24)}
	arg2 := Node{Kind: "literal-int", Value: 5, Ref: NewLocation("", 1, 28, 27)}
	birdField.Args = NodeSet{arg1, arg2}
	birdField.Ref = NewLocation("", 1, 12, 11)
	bird.Children = NodeSet{birdField}
	root.Children = NodeSet{bird}
	actual, err := Parse("", []byte("def Bird { name integer(1, 5) }"))
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, root.String(), actual.(Node).String())
}
