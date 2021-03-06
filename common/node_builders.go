package common

import (
	"time"
)

var PRECEDENCE = map[string]int{
	"=":  1,
	"||": 5,
	"&&": 10,
	"<":  15, ">": 15, "<=": 15, ">=": 15, "==": 15, "!=": 15,
	"+": 20, "-": 20,
	"*": 25, "/": 25, "%": 25,
	"**": 30,
}

func identStr(ident interface{}) string {
	return ident.(*Node).ValStr()
}

func RootNode(l *Location, statements interface{}) *Node {
	children := searchNodes(statements)
	node := &Node{
		Kind:     "root",
		Children: children,
	}

	// minor optimization, no need to rewrap
	if 1 == len(children) && children[0].Kind == "sequential" {
		node.Children = children[0].Children
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
		Kind:    "primary-key",
		Value:   name.(*Node),
		Related: keyType.(*Node),
	}

	return node.withPos(l)
}

func LambdaNode(l *Location, name, params, body interface{}) *Node {
	node := &Node{
		Kind:     "lambda",
		Children: searchNodes(body),
		Args:     DefaultToEmptySlice(params),
	}

	if name != nil {
		node.Name = name.(*Node).ValStr()
	}

	return node.withPos(l)
}

func CallNode(l *Location, callable, args interface{}) *Node {
	node := &Node{
		Kind:  "call",
		Value: callable,
		Args:  DefaultToEmptySlice(args),
	}

	return node.withPos(l)
}

func FieldSetNode(l *Location, fields NodeSet) *Node {
	node := &Node{
		Kind:     "field-set",
		Children: fields,
	}

	return node.withPos(l)
}

func EntityBodyNode(l *Location, mod, fieldset interface{}) *Node {
	node := &Node{
		Kind: "entity-body",
	}

	if nil != mod {
		node.Related = mod.(*Node)
	}

	if nil != fieldset {
		node.Value = fieldset.(*Node)
	}

	return node.withPos(l)
}

func EntityNode(l *Location, name, extends, body interface{}) *Node {
	node := &Node{
		Kind:  "entity",
		Value: body.(*Node),
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

func FieldNode(l *Location, ident, fieldValue interface{}, countRange *Node) *Node {
	node := &Node{
		Kind:       "field",
		Name:       identStr(ident),
		Value:      fieldValue.(*Node),
		CountRange: countRange,
	}
	return node.withPos(l)
}

func AssociativeArgumentNode(l *Location, key, value interface{}) *Node {
	node := &Node{
		Kind:    "associative-arg",
		Value:   key,
		Related: value.(*Node),
	}
	return node.withPos(l)
}

func DistributionNode(l *Location, fn string, intervals interface{}) *Node {
	node := &Node{
		Kind: DIST_TYPE,
		Name: fn,
		Args: DefaultToEmptySlice(intervals),
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
		Kind:  "assignment",
		Name:  identNode.ValStr(),
		Value: valueNode,
	}
	return node.withPos(l)
}

func AtomicNode(l *Location, expr interface{}) *Node {
	node := &Node{
		Kind:  "atomic",
		Value: expr.(*Node),
	}
	return node.withPos(l)
}

func BinaryNode(l *Location, head, tail interface{}) *Node {
	rest := tail.([]interface{})
	result := head.(*Node)

	if len(rest) == 0 {
		return result
	}

	priorPrecedence := 0

	if !result.Is("atomic") {
		result = AtomicNode(result.Ref, result)
	}

	lastNode := result

	for _, r := range rest {
		s := r.([]interface{})
		op := string(s[1].([]interface{})[0].([]byte))
		rhs := s[3].(*Node)

		thisPrecedence := PRECEDENCE[op]

		if thisPrecedence >= priorPrecedence && !lastNode.Is("atomic") {
			n := (&Node{
				Kind:    "binary",
				Name:    op,
				Value:   lastNode.Related,
				Related: rhs,
			}).withPos(rhs.Ref)
			lastNode.Related = n
			lastNode = n
		} else {
			result = (&Node{
				Kind:    "binary",
				Name:    op,
				Value:   result,
				Related: rhs,
			}).withPos(rhs.Ref)
			lastNode = result
		}

		priorPrecedence = thisPrecedence
	}
	return result.withPos(l)
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
		Kind: "builtin",
		Name: value,
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
