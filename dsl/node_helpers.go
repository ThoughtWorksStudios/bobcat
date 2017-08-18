package dsl

import (
	"strings"
	"time"
)

func identStr(ident interface{}) string {
	return ident.(*Node).ValStr()
}

func assembleTime(date, localTime interface{}) (time.Time, error) {
	iso8601Date := date.(string)
	var ts []string

	if localTime != nil {
		ts = localTime.([]string)
	}

	str := strings.Join(append([]string{iso8601Date}, ts...), "")
	return ParseDateLikeJS(str)
}

func RootNode(c *current, statements interface{}) *Node {
	node := &Node{
		Kind:     "root",
		Children: searchNodes(statements),
	}
	return node.withPos(c)
}

func ImportNode(c *current, path string) *Node {
	node := &Node{
		Kind:  "import",
		Value: path,
	}

	return node.withPos(c)
}

func EntityNode(c *current, name, extends *Node, body interface{}) *Node {
	node := &Node{
		Kind:     "entity",
		Children: defaultToEmptySlice(body),
	}

	if nil != name {
		node.Name = name.ValStr()
	}

	if nil != extends {
		node.Related = extends
	}

	return node.withPos(c)
}

func GenNode(c *current, entity, args interface{}) *Node {
	node := &Node{
		Kind:  "generation",
		Value: entity,
		Args:  defaultToEmptySlice(args),
	}
	return node.withPos(c)
}

func StaticFieldNode(c *current, ident, fieldValue interface{}, countRange *Node) *Node {
	node := &Node{
		Kind:       "field",
		Name:       identStr(ident),
		Value:      fieldValue.(*Node),
		CountRange: countRange,
	}
	return node.withPos(c)
}

func DynamicFieldNode(c *current, ident, fieldType, args interface{}, countRange *Node) *Node {
	node := &Node{
		Kind:       "field",
		Name:       identStr(ident),
		Value:      fieldType.(*Node),
		Args:       defaultToEmptySlice(args),
		CountRange: countRange,
	}
	return node.withPos(c)
}

func RangeNode(c *current, lower, upper *Node) *Node {
	node := &Node{
		Kind:     "range",
		Children: NodeSet{lower, upper},
	}
	return node.withPos(c)
}

func AssignNode(c *current, left, right interface{}) *Node {
	identNode, _ := left.(*Node)
	valueNode, _ := right.(*Node)

	if valueNode.Name == "" && valueNode.Kind == "entity" {
		valueNode.Name = identNode.ValStr()
	}

	node := &Node{
		Kind:     "assignment",
		Children: NodeSet{identNode, valueNode},
	}
	return node.withPos(c)
}

func IdNode(c *current, value string) *Node {
	node := &Node{
		Kind:  "identifier",
		Value: value,
	}
	return node.withPos(c)
}

func BuiltinNode(c *current, value string) *Node {
	node := &Node{
		Kind:  "builtin",
		Value: value,
	}
	return node.withPos(c)
}

func DateLiteralNode(c *current, dateTime time.Time) *Node {
	node := &Node{
		Kind:  "literal-date",
		Value: dateTime,
	}

	return node.withPos(c)
}

func IntLiteralNode(c *current, val int64) *Node {
	node := &Node{
		Kind:  "literal-int",
		Value: val,
	}

	return node.withPos(c)
}

func FloatLiteralNode(c *current, val float64) *Node {
	node := &Node{
		Kind:  "literal-float",
		Value: val,
	}

	return node.withPos(c)
}

func NullLiteralNode(c *current) *Node {
	node := &Node{
		Kind: "literal-null",
	}
	return node.withPos(c)
}

func BoolLiteralNode(c *current, val bool) *Node {
	node := &Node{
		Kind:  "literal-bool",
		Value: val,
	}

	return node.withPos(c)
}

func StrLiteralNode(c *current, val string) *Node {
	node := &Node{
		Kind:  "literal-string",
		Value: val,
	}
	return node.withPos(c)
}

func CollectionLiteralNode(c *current, elements interface{}) *Node {
	node := &Node{
		Kind:     "literal-collection",
		Children: defaultToEmptySlice(elements),
	}
	return node.withPos(c)
}
