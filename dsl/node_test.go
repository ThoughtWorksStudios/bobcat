package dsl

import (
	"fmt"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
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
		Bound:   nodeSet,
	}

	actual := node.String()
	expected := fmt.Sprintf("{ Kind: \"%s\", Name: \"%s\", Value: %v, Args: %v, Children: %v, Bound: %v }", "string", "blah", 2, nodeSet, nodeSet, nodeSet)
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

func TestHasRelation(t *testing.T) {
	noRelations := &Node{}
	withRelations := &Node{Related: &Node{}}
	Assert(t, !noRelations.HasRelation(), "if node does not have related node, should report false")
	Assert(t, withRelations.HasRelation(), "if node has related node, should report true")
}
