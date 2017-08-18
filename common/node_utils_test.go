package common

import (
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"testing"
)

func TestSearchNodesWhenGivenSliceOfNodes(t *testing.T) {
	node1 := &Node{Kind: "integer", Name: "one", Value: 1}
	node2 := &Node{Kind: "integer", Name: "two", Value: 2}
	expected := NodeSet{node1, node2}
	actual := searchNodes([]interface{}{node1, node2})
	AssertEqual(t, expected.String(), actual.String())
}

func TestSearchNodesReturnsEmptyNodeSetWhenReceivesNil(t *testing.T) {
	expected := NodeSet{}
	actual := searchNodes(nil)
	AssertEqual(t, expected.String(), actual.String())
}

func TestSearchNodesWhenGivenListOfNonNodes(t *testing.T) {
	node1 := &Node{Kind: "string", Name: "thing", Value: "blah"}
	node2 := &Node{Kind: "integer", Name: "value", Value: 42}
	node3 := &Node{Kind: "dict", Name: "city", Value: "city"}
	expected := NodeSet{node1, node2, node3}
	weirdArgs := []interface{}{[]interface{}{node1, node2, node3}}
	actual := searchNodes(weirdArgs)
	AssertEqual(t, expected.String(), actual.String())
}

func TestSearchNodesWhenGivenListOfNodesAndValues(t *testing.T) {
	topNode := &Node{Kind: "string", Name: "thing", Value: "blah"}
	node1 := &Node{Kind: "integer", Name: "value", Value: 42}
	node2 := &Node{Kind: "dict", Name: "city", Value: "city"}
	expected := NodeSet{topNode, node1, node2}
	weirdArgs := []interface{}{topNode, []interface{}{node1, node2}}
	actual := searchNodes(weirdArgs)
	AssertEqual(t, expected.String(), actual.String())
}

func TestDelimitedNodeSliceWhereFirstAndRestAreNodes(t *testing.T) {
	first := &Node{Kind: "string", Name: "thing", Value: "blah"}
	n := &Node{Kind: "integer", Name: "value", Value: 42}
	var rest interface{} = []interface{}{n}
	expected := NodeSet{first, n}
	actual := DelimitedNodeSlice(first, rest)
	AssertEqual(t, expected.String(), actual.String())
}

func TestDelimitedNodeSliceWhereRestIsSliceOfNodes(t *testing.T) {
	first := &Node{Kind: "string", Name: "thing", Value: "blah"}
	node1 := &Node{Kind: "integer", Name: "value", Value: 42}
	node2 := &Node{Kind: "dict", Name: "city", Value: "city"}
	expected := NodeSet{first, node1, node2}
	rest := []interface{}{[]interface{}{node1, node2}}
	actual := DelimitedNodeSlice(first, rest)
	AssertEqual(t, expected.String(), actual.String())
}
func TestDelimitedNodeSliceWhereRestIsComplex(t *testing.T) {
	first := &Node{Kind: "string", Name: "thing", Value: "blah"}
	node2 := &Node{Kind: "integer", Name: "value", Value: 42}
	node3 := &Node{Kind: "dict", Name: "city", Value: "city"}
	node4 := &Node{Kind: "decimal", Name: "age"}
	expected := NodeSet{first, node2, node3, node4}
	rest := []interface{}{node2, []interface{}{node3, node4}}
	actual := DelimitedNodeSlice(first, rest)
	AssertEqual(t, expected.String(), actual.String())
}

func TestDefaultToEmptySlice(t *testing.T) {
	expected := NodeSet{}
	actual := DefaultToEmptySlice(nil)
	AssertEqual(t, expected.String(), actual.String())

	node1 := &Node{Kind: "integer", Name: "one", Value: 1}
	node2 := &Node{Kind: "integer", Name: "two", Value: 2}
	expected = NodeSet{node1, node2}
	actual = DefaultToEmptySlice(expected)
	AssertEqual(t, expected.String(), actual.String())
}
