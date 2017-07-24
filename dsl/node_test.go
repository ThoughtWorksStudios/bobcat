package dsl

import (
	"fmt"
	. "github.com/ThoughtWorksStudios/datagen/test_helpers"
	"testing"
)

func TestNodeToString(t *testing.T) {
	location := NewLocation("eek", 2, 2, 2)
	nodeSet := NodeSet{Node{Kind: "integer", Name: "blah"}}
	node := Node{
		Kind:     "string",
		Name:     "blah",
		Value:    2,
		Ref:      location,
		Args:     nodeSet,
		Children: nodeSet,
	}

	actual := node.String()
	expected := fmt.Sprintf("{ Kind: \"%s\", Ref: \"%s\", Name: \"%s\", Value: %v, Args: %v, Children: %v }", "string", location.String(), "blah", 2, nodeSet, nodeSet)
	AssertEqual(t, expected, actual)
}

func TestNodeWithPositionReturnsValidNodeWithLocation(t *testing.T) {
	c := &current{
		pos:         position{line: 4, col: 3, offset: 2},
		globalStore: map[string]interface{}{"filename": "whatever.spec"},
	}

	node := Node{Name: "blah"}

	expected := NewLocation("whatever.spec", 4, 3, 2).String()
	actual := node.withPos(c).Ref.String()
	AssertEqual(t, expected, actual)
}

func TestNewLocationReturnsValidLocation(t *testing.T) {
	AssertEqual(t, "whatever.spec:4:8 [byte 42]", NewLocation("whatever.spec", 4, 8, 42).String())
}

func TestHasParent(t *testing.T) {
	node := Node{Parent: "eek"}
	AssertEqual(t, true, node.HasParent())
}
