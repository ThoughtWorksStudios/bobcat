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
	fn := IdNode(nil, name)
	ft := StrLiteralNode(nil, value)
	n := StaticFieldNode(nil, fn, ft, nil)
	return n
}

func TestRootNodeReturnsExpectedNode(t *testing.T) {
	node1 := &Node{}
	node2 := &Node{}
	kids := NodeSet{node1, node2}
	expected := &Node{Kind: "root", Children: kids, Ref: location}
	var statements interface{} = []interface{}{node1, node2}
	actual := RootNode(cnt, statements)

	AssertEqual(t, expected.String(), actual.String())
}

func TestEntityNodeReturnsExpectedNode(t *testing.T) {
	field1 := staticStringField("first", "beth")
	field2 := staticStringField("last", "morty")
	fields := NodeSet{field1, field2}
	ent := EntityNode(nil, nil, fields)

	ident := &Node{Kind: "identifier", Value: "Rick"}

	expected := &Node{Kind: "entity", Name: "Rick", Children: fields}
	actual := NamedEntityNode(nil, ident, ent)

	AssertEqual(t, expected.String(), actual.String())
}

func TestEntityNodeHandleExtension(t *testing.T) {
	field1 := staticStringField("first", "beth")
	field2 := staticStringField("last", "morty")
	fields := NodeSet{field1, field2}
	parent := &Node{Kind: "identifier", Value: "RickestRick"}
	ent := EntityNode(nil, parent, fields)

	ident := &Node{Kind: "identifier", Value: "Rick"}

	expected := &Node{Kind: "entity", Name: "Rick", Related: parent, Children: fields}
	actual := NamedEntityNode(nil, ident, ent)

	AssertEqual(t, expected.String(), actual.String())
}

func TestGenNodeReturnsExpectedNodeWithArgs(t *testing.T) {
	field1 := staticStringField("first", "beth")
	field2 := staticStringField("last", "morty")
	fields := NodeSet{field1, field2}
	ent := EntityNode(nil, nil, fields)

	ident := &Node{Kind: "identifier", Value: "Rick"}

	entity := NamedEntityNode(nil, ident, ent)

	count := IntLiteralNode(nil, int64(5))
	args := NodeSet{count}

	expected := &Node{Kind: "generation", Value: entity, Args: args}
	actual := GenNode(nil, entity, args)

	AssertEqual(t, expected.String(), actual.String())
}

func TestGenNodeReturnsExpectedNodeWithoutArgs(t *testing.T) {
	field1 := staticStringField("first", "beth")
	field2 := staticStringField("last", "morty")
	fields := NodeSet{field1, field2}
	ent := EntityNode(nil, nil, fields)
	ident := &Node{Kind: "identifier", Value: "Rick"}
	entity := NamedEntityNode(nil, ident, ent)

	expected := &Node{Kind: "generation", Value: entity, Args: NodeSet{}}
	actual := GenNode(nil, entity, nil)

	AssertEqual(t, expected.String(), actual.String())
}

func TestStaticFieldNode(t *testing.T) {
	morty := &Node{Kind: "builtin", Name: "grandson", Value: "morty"}
	expected := &Node{Kind: "field", Ref: location, Name: "Rick", Value: morty}
	actual := StaticFieldNode(cnt, &Node{Value: "Rick"}, morty, nil)

	AssertEqual(t, expected.String(), actual.String())
}

func TestDynamicNodeWithoutArgsAndBound(t *testing.T) {
	morty := &Node{Kind: "builtin", Name: "grandson", Value: "morty"}
	expected := &Node{Kind: "field", Ref: location, Name: "Rick", Value: morty, Args: NodeSet{}}
	actual := DynamicFieldNode(cnt, &Node{Value: "Rick"}, morty, nil, nil)

	AssertEqual(t, expected.String(), actual.String())
}

func TestDynamicNodeWithArgsAndBound(t *testing.T) {
	morty := &Node{Kind: "builtin", Name: "grandson", Value: "morty"}
	args := NodeSet{&Node{}}
	r := &Node{}
	expected := &Node{Kind: "field", Ref: location, Name: "Rick", Value: morty, Args: args, CountRange: r}
	actual := DynamicFieldNode(cnt, &Node{Value: "Rick"}, morty, args, r)

	AssertEqual(t, expected.String(), actual.String())
}

func TestIDNode(t *testing.T) {
	expected := &Node{Kind: "identifier", Ref: location, Value: "whatever"}
	actual := IdNode(cnt, "whatever")

	AssertEqual(t, expected.String(), actual.String())
}

func TestBuiltinNode(t *testing.T) {
	expected := &Node{Kind: "builtin", Ref: location, Value: "string"}
	actual := BuiltinNode(cnt, "string")

	AssertEqual(t, expected.String(), actual.String())
}

func TestDateLiteralNode(t *testing.T) {
	fullDate, _ := time.Parse("2006-01-02", "2017-07-19")
	expected := &Node{Kind: "literal-date", Ref: location, Value: fullDate}
	actual := DateLiteralNode(cnt, fullDate)

	AssertEqual(t, expected.String(), actual.String())
}

func TestAssembleTimeReturnsError(t *testing.T) {
	_, err := assembleTime("2017-07-19", []string{"13:00:00-0700"})

	ExpectsError(t, "Not a parsable timestamp", err)
}

func TestIntLiteralNode(t *testing.T) {
	expected := &Node{Kind: "literal-int", Ref: location, Value: int64(5)}
	actual := IntLiteralNode(cnt, 5)

	AssertEqual(t, expected.String(), actual.String())
}

func TestFloatLiteralNode(t *testing.T) {
	expected := &Node{Kind: "literal-float", Ref: location, Value: float64(5)}
	actual := FloatLiteralNode(cnt, float64(5))

	AssertEqual(t, expected.String(), actual.String())
}

func TestNullLiteralNode(t *testing.T) {
	expected := &Node{Kind: "literal-null", Ref: location}
	actual := NullLiteralNode(cnt)
	AssertEqual(t, expected.String(), actual.String())
}

func TestBoolLiteralNode(t *testing.T) {
	expected := &Node{Kind: "literal-bool", Ref: location, Value: true}
	actual := BoolLiteralNode(nil, true)
	AssertEqual(t, expected.String(), actual.String())
}

func TestStrLiteralNode(t *testing.T) {
	expected := &Node{Kind: "literal-string", Ref: location, Value: "v"}
	actual := StrLiteralNode(nil, "v")

	AssertEqual(t, expected.String(), actual.String())
}
