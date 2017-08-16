package dsl

import (
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"testing"
	"time"
)

var cnt = &current{
	pos:         position{line: 4, col: 3, offset: 42},
	text:        []byte("wubba lubba dub dub!!!!"),
	globalStore: map[string]interface{}{"filename": "whatever.spec"},
}

var location = NewLocation("whatever.spec", 4, 3, 42)

func staticStringField(name, value string) *Node {
	fn, _ := idNode(nil, name)
	ft, _ := strLiteralNode(nil, value)
	n, _ := staticFieldNode(nil, fn, ft, nil)
	return n
}

func TestRootNodeReturnsExpectedNode(t *testing.T) {
	node1 := &Node{}
	node2 := &Node{}
	kids := NodeSet{node1, node2}
	expected := &Node{Kind: "root", Children: kids, Ref: location}
	var statements interface{} = []interface{}{node1, node2}
	actual, err := rootNode(cnt, statements)

	AssertNil(t, err, "Got an error constructing root node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestEntityNodeReturnsExpectedNode(t *testing.T) {
	field1 := staticStringField("first", "beth")
	field2 := staticStringField("last", "morty")
	fields := NodeSet{field1, field2}
	ent, _ := entityNode(nil, nil, fields)

	ident := &Node{Kind: "identifier", Value: "Rick"}

	expected := &Node{Kind: "entity", Name: "Rick", Children: fields}
	actual, err := namedEntityNode(nil, ident, ent)

	AssertNil(t, err, "Got an error constructing root node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestEntityNodeHandleExtension(t *testing.T) {
	field1 := staticStringField("first", "beth")
	field2 := staticStringField("last", "morty")
	fields := NodeSet{field1, field2}
	parent := &Node{Kind: "identifier", Value: "RickestRick"}
	ent, _ := entityNode(nil, parent, fields)

	ident := &Node{Kind: "identifier", Value: "Rick"}

	expected := &Node{Kind: "entity", Name: "Rick", Related: parent, Children: fields}
	actual, err := namedEntityNode(nil, ident, ent)

	AssertNil(t, err, "Got an error constructing root node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestGenNodeReturnsExpectedNodeWithArgs(t *testing.T) {
	field1 := staticStringField("first", "beth")
	field2 := staticStringField("last", "morty")
	fields := NodeSet{field1, field2}
	ent, _ := entityNode(nil, nil, fields)

	ident := &Node{Kind: "identifier", Value: "Rick"}

	entity, _ := namedEntityNode(nil, ident, ent)

	count, _ := intLiteralNode(nil, "5")
	args := NodeSet{count}

	expected := &Node{Kind: "generation", Value: entity, Args: args}
	actual, err := genNode(nil, entity, args)

	AssertNil(t, err, "Got an error constructing generation node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestGenNodeReturnsExpectedNodeWithoutArgs(t *testing.T) {
	field1 := staticStringField("first", "beth")
	field2 := staticStringField("last", "morty")
	fields := NodeSet{field1, field2}
	ent, _ := entityNode(nil, nil, fields)
	ident := &Node{Kind: "identifier", Value: "Rick"}
	entity, _ := namedEntityNode(nil, ident, ent)

	expected := &Node{Kind: "generation", Value: entity, Args: NodeSet{}}
	actual, err := genNode(nil, entity, nil)

	AssertNil(t, err, "Got an error constructing generation node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestStaticFieldNode(t *testing.T) {
	morty := &Node{Kind: "builtin", Name: "grandson", Value: "morty"}
	expected := &Node{Kind: "field", Ref: location, Name: "Rick", Value: morty}
	actual, err := staticFieldNode(cnt, &Node{Value: "Rick"}, morty, nil)

	AssertNil(t, err, "Got an error constructing root node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestDynamicNodeWithoutArgsAndBound(t *testing.T) {
	morty := &Node{Kind: "builtin", Name: "grandson", Value: "morty"}
	expected := &Node{Kind: "field", Ref: location, Name: "Rick", Value: morty, Args: NodeSet{}}
	actual, err := dynamicFieldNode(cnt, &Node{Value: "Rick"}, morty, nil, nil)

	AssertNil(t, err, "Got an error constructing root node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestDynamicNodeWithArgsAndBound(t *testing.T) {
	morty := &Node{Kind: "builtin", Name: "grandson", Value: "morty"}
	args := NodeSet{&Node{}}
	expected := &Node{Kind: "field", Ref: location, Name: "Rick", Value: morty, Args: args, CountRange: args}
	actual, err := dynamicFieldNode(cnt, &Node{Value: "Rick"}, morty, args, args)

	AssertNil(t, err, "Got an error constructing root node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestIDNode(t *testing.T) {
	expected := &Node{Kind: "identifier", Ref: location, Value: "whatever"}
	actual, err := idNode(cnt, "whatever")

	AssertNil(t, err, "Got an error constructing root node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestBuiltinNode(t *testing.T) {
	expected := &Node{Kind: "builtin", Ref: location, Value: "kidney"}
	actual, err := builtinNode(cnt, "kidney")

	AssertNil(t, err, "Got an error constructing root node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestDateLiteralNode(t *testing.T) {
	fullDate, _ := time.Parse("2006-01-02", "2017-07-19")
	expected := &Node{Kind: "literal-date", Ref: location, Value: fullDate}
	date := "2017-07-19"
	actual, err := dateLiteralNode(cnt, date, []string{})

	AssertNil(t, err, "Got an error constructing root node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestDateLiteralNodeReturnsError(t *testing.T) {
	_, err := dateLiteralNode(cnt, "2017-07-19", []string{"13:00:00-0700"})

	if err == nil {
		t.Errorf("Expected an error, but got none")
	}
}

func TestIntLiteralNode(t *testing.T) {
	expected := &Node{Kind: "literal-int", Ref: location, Value: 5}
	actual, err := intLiteralNode(cnt, "5")

	AssertNil(t, err, "Got an error constructing int literal node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestIntLiteralNodeError(t *testing.T) {
	_, err := intLiteralNode(cnt, string(5))
	ExpectsError(t, `strconv.ParseInt: parsing "\x05": invalid syntax`, err)
}

func TestFloatLiteralNode(t *testing.T) {
	expected := &Node{Kind: "literal-float", Ref: location, Value: float64(5)}
	actual, err := floatLiteralNode(cnt, "5")

	AssertNil(t, err, "Got an error constructing float literal node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestNullLiteralNode(t *testing.T) {
	expected := &Node{Kind: "literal-null", Ref: location}
	actual, err := nullLiteralNode(cnt)
	AssertNil(t, err, "Got an error constructing null literal node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestBoolLiteralNode(t *testing.T) {
	expected := &Node{Kind: "literal-bool", Ref: location, Value: true}
	actual, err := boolLiteralNode(nil, "true")
	AssertNil(t, err, "Got an error constructing bool literal node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestBoolLiteralNodeReturnsError(t *testing.T) {
	_, err := boolLiteralNode(nil, "eek")
	ExpectsError(t, `strconv.ParseBool: parsing "eek": invalid syntax`, err)
}

func TestStrLiteralNode(t *testing.T) {
	expected := &Node{Kind: "literal-string", Ref: location, Value: "v"}
	actual, err := strLiteralNode(nil, `"v"`)

	AssertNil(t, err, "Got an error constructing string literal node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestStrLiteralNodeReturnsError(t *testing.T) {
	_, err := strLiteralNode(nil, "not quoted")
	ExpectsError(t, "invalid syntax", err)
}
