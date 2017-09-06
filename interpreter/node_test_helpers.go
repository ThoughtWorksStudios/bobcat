package interpreter

import (
	ast "github.com/ThoughtWorksStudios/bobcat/common"
	"log"
	"strings"
	"time"
)

func Field(name string, kind *ast.Node, args ...*ast.Node) *ast.Node {
	ident := ast.IdNode(nil, name)
	if strings.HasPrefix(kind.Kind, "literal-") {
		return ast.StaticFieldNode(nil, ident, kind, nil)
	}

	ns := append(make(ast.NodeSet, 0, len(args)), args...)
	return ast.DynamicFieldNode(nil, ident, kind, ns, nil, false)
}

func Builtin(value string) *ast.Node {
	return ast.BuiltinNode(nil, value)
}

func StringVal(val string) *ast.Node {
	return ast.StrLiteralNode(nil, val)
}

func IntVal(val int64) *ast.Node {
	return ast.IntLiteralNode(nil, val)
}

func FloatVal(val float64) *ast.Node {
	return ast.FloatLiteralNode(nil, val)
}

func DateVal(val string) *ast.Node {
	parsed, err := time.Parse("2006-01-02", val)

	if err != nil {
		log.Fatalf("could not parse %v against YYYY-mm-dd. Error: %v", val, err)
	}

	return ast.DateLiteralNode(nil, parsed)
}

func StringArgs(values ...string) ast.NodeSet {
	args := make(ast.NodeSet, len(values))

	for i, val := range values {
		args[i] = StringVal(val)
	}

	return args
}

func IntArgs(values ...int64) ast.NodeSet {
	args := make(ast.NodeSet, len(values))

	for i, val := range values {
		args[i] = IntVal(val)
	}

	return args
}

func FloatArgs(values ...float64) ast.NodeSet {
	args := make(ast.NodeSet, len(values))

	for i, val := range values {
		args[i] = FloatVal(val)
	}

	return args
}

func StringCollection(vals ...string) *ast.Node {
	value := make(ast.NodeSet, len(vals))
	for idx, str := range vals {
		value[idx] = StringVal(str)
	}
	return ast.CollectionLiteralNode(nil, value)
}

func DateArgs(values ...string) ast.NodeSet {
	args := make(ast.NodeSet, len(values))

	for i, val := range values {
		args[i] = DateVal(val)
	}

	return args
}

func Root(nodes ...*ast.Node) *ast.Node {
	ns := append(make(ast.NodeSet, 0, len(nodes)), nodes...)

	return ast.RootNode(nil, ns)
}

func Generation(count int64, entity *ast.Node) *ast.Node {
	return ast.GenNode(nil, ast.NodeSet{IntVal(count), entity})
}

func Entity(name string, fields ast.NodeSet) *ast.Node {
	var body *ast.Node

	if len(fields) > 0 {
		body = ast.EntityBodyNode(nil, nil, ast.FieldSetNode(nil, fields))
	} else {
		body = ast.EntityBodyNode(nil, nil, nil)
	}

	return ast.EntityNode(nil, Id(name), nil, body)
}

func EntityExtension(name, extends string, fields ast.NodeSet) *ast.Node {
	var body *ast.Node

	if len(fields) > 0 {
		body = ast.EntityBodyNode(nil, nil, ast.FieldSetNode(nil, fields))
	} else {
		body = ast.EntityBodyNode(nil, nil, nil)
	}

	return ast.EntityNode(nil, Id(name), Id(extends), body)
}

func Id(name string) *ast.Node {
	return ast.IdNode(nil, name)
}
