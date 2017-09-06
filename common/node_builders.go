package common

import (
	"time"
)

func identStr(ident interface{}) string {
	return ident.(*Node).ValStr()
}

func RootNode(l *Location, statements interface{}) *Node {
	node := &Node{
		Kind:     "root",
		Children: searchNodes(statements),
	}
	return node.withPos(l)
}

func ImportNode(l *Location, path string) *Node {
	node := &Node{
		Kind:  "import",
		Value: path,
	}

	return node.withPos(l)
}

func PkNode(l *Location, name, keyType interface{}) *Node {
	node := &Node{
		Kind:    "primary-key-def",
		Value:   name.(*Node),
		Related: keyType.(*Node),
	}

	return node.withPos(l)
}

func EntityNode(l *Location, name, extends, body interface{}) *Node {
	node := &Node{
		Kind:     "entity",
		Children: DefaultToEmptySlice(body),
	}

	if nil != name {
		node.Name = name.(*Node).ValStr()
	}

	if nil != extends {
		node.Related = extends.(*Node)
	}

	return node.withPos(l)
}

func GenNode(l *Location, args interface{}) *Node {
	node := &Node{
		Kind: "generation",
		Args: DefaultToEmptySlice(args),
	}
	return node.withPos(l)
}

func StaticFieldNode(l *Location, ident, fieldValue interface{}, countRange *Node) *Node {
	node := &Node{
		Kind:       "field",
		Name:       identStr(ident),
		Value:      fieldValue.(*Node),
		CountRange: countRange,
	}
	return node.withPos(l)
}

func DynamicFieldNode(l *Location, ident, fieldType, args interface{}, countRange *Node, unique bool) *Node {
	node := &Node{
		Kind:       "field",
		Name:       identStr(ident),
		Value:      fieldType.(*Node),
		Args:       DefaultToEmptySlice(args),
		CountRange: countRange,
		Unique:     unique,
	}
	return node.withPos(l)
}

func RangeNode(l *Location, lower, upper *Node) *Node {
	node := &Node{
		Kind:     "range",
		Children: NodeSet{lower, upper},
	}
	return node.withPos(l)
}

func VariableNode(l *Location, ident, init interface{}) *Node {
	node := &Node{
		Kind:  "variable",
		Name:  identStr(ident),
		Value: init,
	}
	return node.withPos(l)
}

func SequentialNode(l *Location, expressions interface{}) *Node {
	node := &Node{
		Kind:     "sequential",
		Children: expressions.(NodeSet),
	}
	return node.withPos(l)
}

func AssignNode(l *Location, left, right interface{}) *Node {
	identNode, _ := left.(*Node)
	valueNode, _ := right.(*Node)

	node := &Node{
		Kind:     "assignment",
		Children: NodeSet{identNode, valueNode},
	}
	return node.withPos(l)
}

func IdNode(l *Location, value string) *Node {
	node := &Node{
		Kind:  "identifier",
		Value: value,
	}
	return node.withPos(l)
}

func BuiltinNode(l *Location, value string) *Node {
	node := &Node{
		Kind:  "builtin",
		Value: value,
	}
	return node.withPos(l)
}

func DateLiteralNode(l *Location, dateTime time.Time) *Node {
	node := &Node{
		Kind:  "literal-date",
		Value: dateTime,
	}

	return node.withPos(l)
}

func IntLiteralNode(l *Location, val int64) *Node {
	node := &Node{
		Kind:  "literal-int",
		Value: val,
	}

	return node.withPos(l)
}

func FloatLiteralNode(l *Location, val float64) *Node {
	node := &Node{
		Kind:  "literal-float",
		Value: val,
	}

	return node.withPos(l)
}

func NullLiteralNode(l *Location) *Node {
	node := &Node{
		Kind: "literal-null",
	}
	return node.withPos(l)
}

func BoolLiteralNode(l *Location, val bool) *Node {
	node := &Node{
		Kind:  "literal-bool",
		Value: val,
	}

	return node.withPos(l)
}

func StrLiteralNode(l *Location, val string) *Node {
	node := &Node{
		Kind:  "literal-string",
		Value: val,
	}
	return node.withPos(l)
}

func CollectionLiteralNode(l *Location, elements interface{}) *Node {
	node := &Node{
		Kind:     "literal-collection",
		Children: DefaultToEmptySlice(elements),
	}
	return node.withPos(l)
}
