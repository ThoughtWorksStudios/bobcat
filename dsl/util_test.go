package dsl

import "testing"
import "time"

func TestSearchNodesWhenGivenSliceOfNodes(t *testing.T) {
	node1 := Node{Kind: "integer", Name: "one", Value: 1}
	node2 := Node{Kind: "integer", Name: "two", Value: 2}
	expected := NodeSet{node1, node2}
	actual := searchNodes([]interface{}{node1, node2})
	if actual.String() != expected.String() {
		t.Errorf("Didn't get expected value\nexpected: %v \ngot       %v", expected, actual)
	}
}

func TestSearchNodesReturnsEmptyNodeSetWhenReceivesNil(t *testing.T) {
	expected := NodeSet{}
	actual := searchNodes(nil)
	if actual.String() != expected.String() {
		t.Errorf("expected searchNodes(nil) to return empty NodeSet, but got %v", actual)
	}
}

func TestSearchNodesWhenGivenListOfNonNodes(t *testing.T) {
	node1 := Node{Kind: "string", Name: "thing", Value: "blah"}
	node2 := Node{Kind: "integer", Name: "value", Value: 42}
	node3 := Node{Kind: "dict", Name: "city", Value: "city"}
	expected := NodeSet{node1, node2, node3}
	weirdArgs := []interface{}{[]interface{}{node1, node2, node3}}
	actual := searchNodes(weirdArgs)
	if actual.String() != expected.String() {
		t.Errorf("Didn't get expected value\nexpected: %v \ngot       %v", expected, actual)
	}
}

func TestSearchNodesWhenGivenListOfNodesAndValues(t *testing.T) {
	topNode := Node{Kind: "string", Name: "thing", Value: "blah"}
	node1 := Node{Kind: "integer", Name: "value", Value: 42}
	node2 := Node{Kind: "dict", Name: "city", Value: "city"}
	expected := NodeSet{topNode, node1, node2}
	weirdArgs := []interface{}{topNode, []interface{}{node1, node2}}
	actual := searchNodes(weirdArgs)
	if actual.String() != expected.String() {
		t.Errorf("Didn't get expected value\nexpected: %v \ngot       %v", expected, actual)
	}
}

func TestDelimitedNodeSliceWhereFirstAndRestAreNodes(t *testing.T) {
	first := Node{Kind: "string", Name: "thing", Value: "blah"}
	n := Node{Kind: "integer", Name: "value", Value: 42}
	var rest interface{} = []interface{}{n}
	expected := NodeSet{first, n}
	actual := delimitedNodeSlice(first, rest)
	if actual.String() != expected.String() {
		t.Errorf("Didn't get expected value\nexpected: %v \ngot       %v", expected, actual)
	}
}

func TestDelimitedNodeSliceWhereRestIsSliceOfNodes(t *testing.T) {
	first := Node{Kind: "string", Name: "thing", Value: "blah"}
	node1 := Node{Kind: "integer", Name: "value", Value: 42}
	node2 := Node{Kind: "dict", Name: "city", Value: "city"}
	expected := NodeSet{first, node1, node2}
	rest := []interface{}{[]interface{}{node1, node2}}
	actual := delimitedNodeSlice(first, rest)
	if actual.String() != expected.String() {
		t.Errorf("Didn't get expected value\nexpected: %v \ngot       %v", expected, actual)
	}
}
func TestDelimitedNodeSliceWhereRestIsComplex(t *testing.T) {
	first := Node{Kind: "string", Name: "thing", Value: "blah"}
	node2 := Node{Kind: "integer", Name: "value", Value: 42}
	node3 := Node{Kind: "dict", Name: "city", Value: "city"}
	node4 := Node{Kind: "decimal", Name: "age"}
	expected := NodeSet{first, node2, node3, node4}
	rest := []interface{}{node2, []interface{}{node3, node4}}
	actual := delimitedNodeSlice(first, rest)
	if actual.String() != expected.String() {
		t.Errorf("Didn't get expected value\nexpected: %v \ngot       %v", expected, actual)
	}
}

func TestCharGroupAsString(t *testing.T) {
	expected := "1:3"
	var input interface{} = []interface{}{[]uint8{'1'}, []uint8{':'}, []uint8{'3'}}
	actual := charGroupAsString(input)
	if actual != expected {
		t.Errorf("Didn't get expected value\nexpected: %v \ngot       %v", expected, actual)
	}
}

func TestParseDateLikeJSWithTimeZone(t *testing.T) {
	input := "2017-07-19T13:00:00-07:00"
	expected, _ := time.Parse("2006-01-02 15:04:00 (MST)", "2017-07-19 13:00:00 -0700 PDT")
	actual, err := ParseDateLikeJS(input)
	if err != nil {
		t.Errorf("Got an error while parsing date: %v", err)
	} else if actual.Equal(expected) {
		t.Errorf("Didn't get expected value\nexpected: %v \ngot       %v", expected, actual)
	}
}

func TestParseDateLikeJSUTC(t *testing.T) {
	input := "2017-07-19T13:00:00Z"
	expected, _ := time.Parse("2006-01-02 15:04:00 (MST)", "2017-07-19 13:00:00 +0000 UTC")

	actual, err := ParseDateLikeJS(input)
	if err != nil {
		t.Errorf("Got an error while parsing date: %v", err)
	} else if actual.Equal(expected) {
		t.Errorf("Didn't get expected value\nexpected: %v \ngot       %v", expected, actual)
	}
}

func TestParseDateLikeJSReturnsError(t *testing.T) {
	input := "2017-07-19T13:00:00Z-700"
	expected := "Not a parsable timestamp: 2017-07-19T13:00:00Z-700"
	_, err := ParseDateLikeJS(input)
	if err == nil || err.Error() != expected {
		t.Errorf("Didn't get the expected error\nexpected: %v \ngot       %v", expected, err)
	}
}

func TestDefaultToEmptySlice(t *testing.T) {
	expected := NodeSet{}
	actual := defaultToEmptySlice(nil)
	if actual.String() != expected.String() {
		t.Errorf("expected defaultToEmptySlice(nil) to return an empty NodeSet, but got %v", actual)
	}

	node1 := Node{Kind: "integer", Name: "one", Value: 1}
	node2 := Node{Kind: "integer", Name: "two", Value: 2}
	expected = NodeSet{node1, node2}
	actual = defaultToEmptySlice(expected)
	if actual.String() != expected.String() {
		t.Errorf("Didn't get expected value\nexpected: %v \ngot       %v", expected, actual)
	}

}
