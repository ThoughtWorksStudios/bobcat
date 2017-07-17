package dsl

import (
	"fmt"
	"strings"
)

type Node struct {
	Kind     string
	Name     string
	Value    interface{}
	Args     NodeSet
	Children NodeSet
}

func (n Node) String() string {
	attrs := make([]string, 1)

	attrs[0] = fmt.Sprintf("Kind: %v", n.Kind)
	if n.Name != "" {
		attrs = append(attrs, fmt.Sprintf("Name: %v", n.Name))
	}
	if n.Value != nil {
		attrs = append(attrs, fmt.Sprintf("Value: %v", n.Value))
	}
	if n.Args != nil {
		attrs = append(attrs, fmt.Sprintf("Args: %v", n.Args))
	}
	if n.Children != nil {
		attrs = append(attrs, fmt.Sprintf("Children: %v", n.Children))
	}

	return fmt.Sprintf("{ %s }", strings.Join(attrs, ", "))
}

type NodeSet []Node // bless this with functional shims

type Iterator func(index int, node Node)
type Collector func(index int, node Node) interface{}

func (nodes NodeSet) Each(f Iterator) NodeSet {
	for i, size := 0, len(nodes); i < size; i++ {
		f(i, nodes[i])
	}
	return nodes
}

func (nodes NodeSet) Map(f Collector) []interface{} {
	size := len(nodes)
	result := make([]interface{}, size)
	nodes.Each(func(index int, node Node) {
		result[index] = f(index, node)
	})
	return result
}
