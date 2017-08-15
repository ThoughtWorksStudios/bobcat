package dsl

import (
	"strconv"
	"strings"
)

func identStr(ident interface{}) string {
	return ident.(Node).Value.(string)
}

func rootNode(c *current, statements interface{}) (Node, error) {
	node := &Node{
		Kind:     "root",
		Children: searchNodes(statements),
	}
	return node.withPos(c), nil
}

func importNode(c *current, path string) (Node, error) {
	node := &Node{
		Kind:  "import",
		Value: path,
	}

	return node.withPos(c), nil
}

func namedEntityNode(c *current, identifier, entity interface{}) (Node, error) {
	node, _ := entity.(Node)

	if nil != identifier {
		node.Name = identStr(identifier)
	}

	return node.withPos(c), nil
}

func entityNode(c *current, extends, body interface{}) (Node, error) {
	node := &Node{
		Kind:     "entity",
		Children: defaultToEmptySlice(body),
	}

	if nil != extends {
		if parentIdentNode, ok := extends.(Node); ok {
			node.Related = &parentIdentNode
		} else {
			return *node, node.Err("Entity cannot extend %T %v", extends, extends)
		}
	}

	return node.withPos(c), nil
}

func genNode(c *current, entity, args interface{}) (Node, error) {
	node := &Node{
		Kind:  "generation",
		Value: entity,
		Args:  defaultToEmptySlice(args),
	}
	return node.withPos(c), nil
}

func staticFieldNode(c *current, ident, fieldValue interface{}) (Node, error) {
	node := &Node{
		Kind:  "field",
		Name:  identStr(ident),
		Value: fieldValue.(Node),
	}
	return node.withPos(c), nil
}

func dynamicFieldNode(c *current, ident, fieldType, args interface{}, countRange NodeSet) (Node, error) {
	node := &Node{
		Kind:       "field",
		Name:       identStr(ident),
		Value:      fieldType.(Node),
		Args:       defaultToEmptySlice(args),
		CountRange: countRange,
	}
	return node.withPos(c), nil
}

func assignNode(c *current, left, right interface{}) (Node, error) {
	identNode, _ := left.(Node)
	valueNode, _ := right.(Node)

	if valueNode.Name == "" && valueNode.Kind == "entity" {
		valueNode.Name = identNode.ValStr()
	}

	node := &Node{
		Kind:     "assignment",
		Children: NodeSet{identNode, valueNode},
	}
	return node.withPos(c), nil
}

func idNode(c *current, value string) (Node, error) {
	node := &Node{
		Kind:  "identifier",
		Value: value,
	}
	return node.withPos(c), nil
}

func builtinNode(c *current, value string) (Node, error) {
	node := &Node{
		Kind:  "builtin",
		Value: value,
	}
	return node.withPos(c), nil
}

func dateLiteralNode(c *current, date, localTime interface{}) (Node, error) {
	iso8601Date := date.(string)
	var ts []string

	if localTime != nil {
		ts = localTime.([]string)
	}

	str := strings.Join(append([]string{iso8601Date}, ts...), "")
	parsed, er := ParseDateLikeJS(str)

	node := &Node{
		Kind:  "literal-date",
		Value: parsed,
	}

	return node.withPos(c), er
}

func intLiteralNode(c *current, s string) (Node, error) {
	val, er := strconv.ParseInt(s, 10, 64)
	node := &Node{
		Kind:  "literal-int",
		Value: val,
	}

	return node.withPos(c), er
}

func floatLiteralNode(c *current, s string) (Node, error) {
	val, er := strconv.ParseFloat(s, 64)

	node := &Node{
		Kind:  "literal-float",
		Value: val,
	}

	return node.withPos(c), er
}

func nullLiteralNode(c *current) (Node, error) {
	node := &Node{
		Kind: "literal-null",
	}
	return node.withPos(c), nil
}

func boolLiteralNode(c *current, value string) (Node, error) {
	val, er := strconv.ParseBool(value)

	node := &Node{
		Kind:  "literal-bool",
		Value: val,
	}

	return node.withPos(c), er
}

func strLiteralNode(c *current, value string) (Node, error) {
	val, er := strconv.Unquote(value)

	node := &Node{
		Kind:  "literal-string",
		Value: val,
	}
	return node.withPos(c), er
}
