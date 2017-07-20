package dsl

import (
	"strconv"
	"strings"
)

func rootNode(c *current, statements interface{}) (Node, error) {
	node := &Node{
		Kind:     "root",
		Children: searchNodes(statements),
	}
	return node.withPos(c), nil
}

func entityNode(c *current, name, body interface{}) (Node, error) {
	node := &Node{
		Kind:     "definition",
		Name:     name.(Node).Value.(string),
		Children: defaultToEmptySlice(body),
	}
	return node.withPos(c), nil
}

func genNode(c *current, name, body, args interface{}) (Node, error) {
	node := &Node{
		Kind: "generation",
		Name: name.(Node).Value.(string),
		Children: defaultToEmptySlice(body),
		Args: defaultToEmptySlice(args),
	}
	return node.withPos(c), nil
}

func staticFieldNode(c *current, name, fieldValue interface{}) (Node, error) {
	node := &Node{
		Kind:  "field",
		Name:  name.(Node).Value.(string),
		Value: fieldValue.(Node),
	}
	return node.withPos(c), nil
}

func dynamicFieldNode(c *current, name, fieldType, args interface{}) (Node, error) {
	node := &Node{
		Kind:  "field",
		Name:  name.(Node).Value.(string),
		Value: fieldType.(Node),
		Args:  defaultToEmptySlice(args),
	}
	return node.withPos(c), nil
}

func idNode(c *current) (Node, error) {
	node := &Node{
		Kind:  "identifier",
		Value: string(c.text),
	}
	return node.withPos(c), nil
}

func builtinNode(c *current) (Node, error) {
	node := &Node{
		Kind:  "builtin",
		Value: string(c.text),
	}
	return node.withPos(c), nil
}

func dateLiteralNode(c *current, date, localTime interface{}) (Node, error) {
	iso8601Date := charGroupAsString(date)
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

func boolLiteralNode(c *current) (Node, error) {
	val, er := strconv.ParseBool(string(c.text))

	node := &Node{
		Kind:  "literal-bool",
		Value: val,
	}

	return node.withPos(c), er
}

func strLiteralNode(c *current) (Node, error) {
	val, er := strconv.Unquote(string(c.text))

	node := &Node{
		Kind:  "literal-string",
		Value: val,
	}
	return node.withPos(c), er
}
