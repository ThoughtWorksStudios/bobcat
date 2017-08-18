package interpreter

import (
	"github.com/ThoughtWorksStudios/bobcat/dsl"
	"log"
	"strings"
	"time"
)

func FieldNode(name string, kind *dsl.Node, args ...*dsl.Node) *dsl.Node {
	ident := dsl.IdNode(nil, name)
	if strings.HasPrefix(kind.Kind, "literal-") {
		return dsl.StaticFieldNode(nil, ident, kind, nil)
	}

	ns := append(make(dsl.NodeSet, 0, len(args)), args...)
	return dsl.DynamicFieldNode(nil, ident, kind, ns, nil)
}

func BuiltinNode(value string) *dsl.Node {
	return dsl.BuiltinNode(nil, value)
}

func StringNode(val string) *dsl.Node {
	return dsl.StrLiteralNode(nil, val)
}

func IntNode(val int64) *dsl.Node {
	return dsl.IntLiteralNode(nil, val)
}

func FloatNode(val float64) *dsl.Node {
	return dsl.FloatLiteralNode(nil, val)
}

func DateNode(val string) *dsl.Node {
	parsed, err := time.Parse("2006-01-02", val)

	if err != nil {
		log.Fatalf("could not parse %v against YYYY-mm-dd. Error: %v", val, err)
	}

	return dsl.DateLiteralNode(nil, parsed)
}

func StringArgs(values ...string) dsl.NodeSet {
	args := make(dsl.NodeSet, len(values))

	for i, val := range values {
		args[i] = StringNode(val)
	}

	return args
}

func IntArgs(values ...int64) dsl.NodeSet {
	args := make(dsl.NodeSet, len(values))

	for i, val := range values {
		args[i] = IntNode(val)
	}

	return args
}

func FloatArgs(values ...float64) dsl.NodeSet {
	args := make(dsl.NodeSet, len(values))

	for i, val := range values {
		args[i] = FloatNode(val)
	}

	return args
}

func StringCollectionNode(vals ...string) *dsl.Node {
	value := make(dsl.NodeSet, len(vals))
	for idx, str := range vals {
		value[idx] = StringNode(str)
	}
	return dsl.CollectionLiteralNode(nil, value)
}

func DateArgs(values ...string) dsl.NodeSet {
	args := make(dsl.NodeSet, len(values))

	for i, val := range values {
		args[i] = DateNode(val)
	}

	return args
}

func RootNode(nodes ...*dsl.Node) *dsl.Node {
	ns := append(make(dsl.NodeSet, 0, len(nodes)), nodes...)

	return dsl.RootNode(nil, ns)
}

func GenerationNode(entity *dsl.Node, count int64) *dsl.Node {
	return dsl.GenNode(nil, entity, IntArgs(count))
}

func EntityNode(name string, fields dsl.NodeSet) *dsl.Node {
	return dsl.EntityNode(nil, IdNode(name), nil, fields)
}

func EntityExtensionNode(name, extends string, fields dsl.NodeSet) *dsl.Node {
	return dsl.EntityNode(nil, IdNode(name), IdNode(extends), fields)
}

func IdNode(name string) *dsl.Node {
	return dsl.IdNode(nil, name)
}
