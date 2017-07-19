package dsl

import (
	"fmt"
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
	if actual != expected {
		t.Errorf("Didn't get expected value\nexpected: %v \ngot %v", actual, expected)
	}
}

func TestNodeWithPositionReturnsValidNodeWithLocation(t *testing.T) {
	expected := NewLocation("whatever.spec", 4, 3, 2)
	c := &current{
		pos:         position{line: 4, col: 3, offset: 2},
		globalStore: map[string]interface{}{"filename": "whatever.spec"},
	}
	node := Node{Name: "blah"}
	actual := node.withPos(c).Ref
	if actual.String() != expected.String() {
		t.Errorf("Didn't get expected value\nexpected: %v \ngot %v", actual, expected)
	}

}

func TestNewLocationReturnsValidLocation(t *testing.T) {
	expected := "whatever.spec:4:8 [byte 42]"
	actual := NewLocation("whatever.spec", 4, 8, 42).String()
	if expected != actual {
		t.Errorf("Didn't get expected value\nexpected: %v \ngot %v", actual, expected)
	}
}
