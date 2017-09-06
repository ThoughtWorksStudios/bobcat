package common

import (
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"testing"
	"time"
)

var ref = NewLocation("whatever.spec", 4, 3, 42)

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
	expected := &Node{Kind: "root", Children: kids, Ref: ref}
	var statements interface{} = []interface{}{node1, node2}
	actual := RootNode(ref, statements)

	AssertEqual(t, expected.String(), actual.String())
}

func TestFieldSetNode(t *testing.T) {
	field1 := staticStringField("first", "beth")
	field2 := staticStringField("last", "morty")
	fields := NodeSet{field1, field2}
	expected := &Node{Kind: "field-set", Children: fields}

	AssertEqual(t, expected.String(), FieldSetNode(nil, fields).String())
}

func TestPrimaryKeyNode(t *testing.T) {
	name := StrLiteralNode(nil, "id")
	val := BuiltinNode(nil, "uid")
	expected := &Node{Kind: "primary-key", Value: name, Related: val}

	AssertEqual(t, expected.String(), PkNode(nil, name, val).String())
}

func TestEntityBodyNode(t *testing.T) {
	field1 := staticStringField("first", "beth")
	field2 := staticStringField("last", "morty")
	fields := FieldSetNode(nil, NodeSet{field1, field2})
	pk := PkNode(nil, StrLiteralNode(nil, "myid"), BuiltinNode(nil, "uid"))

	expected := &Node{Kind: "entity-body", Value: fields, Related: pk}
	AssertEqual(t, expected.String(), EntityBodyNode(nil, pk, fields).String())

	expected = &Node{Kind: "entity-body", Value: fields}
	AssertEqual(t, expected.String(), EntityBodyNode(nil, nil, fields).String())

	expected = &Node{Kind: "entity-body", Related: pk}
	AssertEqual(t, expected.String(), EntityBodyNode(nil, pk, nil).String())

	expected = &Node{Kind: "entity-body"}
	AssertEqual(t, expected.String(), EntityBodyNode(nil, nil, nil).String())
}

func TestEntityNodeReturnsExpectedNode(t *testing.T) {
	field1 := staticStringField("first", "beth")
	field2 := staticStringField("last", "morty")
	fields := FieldSetNode(nil, NodeSet{field1, field2})
	body := EntityBodyNode(nil, nil, fields)

	ident := IdNode(nil, "Rick")

	expected := &Node{Kind: "entity", Name: "Rick", Value: body}
	actual := EntityNode(nil, ident, nil, body)

	AssertEqual(t, expected.String(), actual.String())
}

func TestEntityNodeHandleExtension(t *testing.T) {
	field1 := staticStringField("first", "beth")
	field2 := staticStringField("last", "morty")
	fields := FieldSetNode(nil, NodeSet{field1, field2})
	body := EntityBodyNode(nil, nil, fields)

	ident := IdNode(nil, "Rick")
	parent := IdNode(nil, "RickestRick")

	expected := &Node{Kind: "entity", Name: "Rick", Related: parent, Value: body}
	actual := EntityNode(nil, ident, parent, body)

	AssertEqual(t, expected.String(), actual.String())
}

func TestGenNodeReturnsExpectedNode(t *testing.T) {
	field1 := staticStringField("first", "beth")
	field2 := staticStringField("last", "morty")
	fields := FieldSetNode(nil, NodeSet{field1, field2})
	body := EntityBodyNode(nil, nil, fields)

	ident := IdNode(nil, "Rick")
	entity := EntityNode(nil, ident, nil, body)

	args := NodeSet{IntLiteralNode(nil, int64(5)), entity}

	expected := &Node{Kind: "generation", Args: args}
	actual := GenNode(nil, args)

	AssertEqual(t, expected.String(), actual.String())
}

func TestStaticFieldNode(t *testing.T) {
	morty := &Node{Kind: "builtin", Name: "grandson", Value: "morty"}
	expected := &Node{Kind: "field", Ref: ref, Name: "Rick", Value: morty}
	actual := StaticFieldNode(ref, &Node{Value: "Rick"}, morty, nil)

	AssertEqual(t, expected.String(), actual.String())
}

func TestDynamicNodeWithoutArgsAndBound(t *testing.T) {
	morty := &Node{Kind: "builtin", Name: "grandson", Value: "morty"}
	expected := &Node{Kind: "field", Ref: ref, Name: "Rick", Value: morty, Args: NodeSet{}}
	actual := DynamicFieldNode(ref, &Node{Value: "Rick"}, morty, nil, nil, false)

	AssertEqual(t, expected.String(), actual.String())
}

func TestDynamicNodeWithArgsAndBound(t *testing.T) {
	morty := &Node{Kind: "builtin", Name: "grandson", Value: "morty"}
	args := NodeSet{&Node{}}
	r := &Node{}
	expected := &Node{Kind: "field", Ref: ref, Name: "Rick", Value: morty, Args: args, CountRange: r}
	actual := DynamicFieldNode(ref, &Node{Value: "Rick"}, morty, args, r, false)

	AssertEqual(t, expected.String(), actual.String())
}

func TestIDNode(t *testing.T) {
	expected := &Node{Kind: "identifier", Ref: ref, Value: "whatever"}
	actual := IdNode(ref, "whatever")

	AssertEqual(t, expected.String(), actual.String())
}

func TestBuiltinNode(t *testing.T) {
	expected := &Node{Kind: "builtin", Ref: ref, Value: "string"}
	actual := BuiltinNode(ref, "string")

	AssertEqual(t, expected.String(), actual.String())
}

func TestDateLiteralNode(t *testing.T) {
	fullDate, _ := time.Parse("2006-01-02", "2017-07-19")
	expected := &Node{Kind: "literal-date", Ref: ref, Value: fullDate}
	actual := DateLiteralNode(ref, fullDate)

	AssertEqual(t, expected.String(), actual.String())
}

func TestIntLiteralNode(t *testing.T) {
	expected := &Node{Kind: "literal-int", Ref: ref, Value: int64(5)}
	actual := IntLiteralNode(ref, 5)

	AssertEqual(t, expected.String(), actual.String())
}

func TestFloatLiteralNode(t *testing.T) {
	expected := &Node{Kind: "literal-float", Ref: ref, Value: float64(5)}
	actual := FloatLiteralNode(ref, float64(5))

	AssertEqual(t, expected.String(), actual.String())
}

func TestNullLiteralNode(t *testing.T) {
	expected := &Node{Kind: "literal-null", Ref: ref}
	actual := NullLiteralNode(ref)
	AssertEqual(t, expected.String(), actual.String())
}

func TestBoolLiteralNode(t *testing.T) {
	expected := &Node{Kind: "literal-bool", Ref: ref, Value: true}
	actual := BoolLiteralNode(nil, true)
	AssertEqual(t, expected.String(), actual.String())
}

func TestStrLiteralNode(t *testing.T) {
	expected := &Node{Kind: "literal-string", Ref: ref, Value: "v"}
	actual := StrLiteralNode(nil, "v")

	AssertEqual(t, expected.String(), actual.String())
}
