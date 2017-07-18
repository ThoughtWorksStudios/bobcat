package dsl

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Node struct {
	Kind     string
	Name     string
	Value    interface{}
	Args     NodeSet
	Children NodeSet
	Ref      *Location
}

type Location struct {
	line, col, offset int
	filename          string
}

// NOTE: For testing purposes only
func NewLocation(filename string, line, col, offset int) *Location {
	return &Location{
		filename: filename,
		line:     line,
		col:      col,
		offset:   offset,
	}
}

func (l *Location) String() string {
	return fmt.Sprintf("%s:%d:%d [byte %d]", l.filename, l.line, l.col, l.offset)
}

func (n Node) String() string {
	attrs := make([]string, 1)

	attrs[0] = fmt.Sprintf("Kind: %s", strconv.Quote(n.Kind))

	if n.Ref != nil {
		attrs = append(attrs, fmt.Sprintf("Ref: %s", strconv.Quote(n.Ref.String())))
	}

	if n.Name != "" {
		attrs = append(attrs, fmt.Sprintf("Name: %s", strconv.Quote(n.Name)))
	}

	if n.Value != nil {
		switch n.Value.(type) {
		case time.Time:
			attrs = append(attrs, fmt.Sprintf("Value: %s", strconv.Quote(n.Value.(time.Time).String())))
		case string:
			attrs = append(attrs, fmt.Sprintf("Value: %s", strconv.Quote(n.Value.(string))))
		default:
			attrs = append(attrs, fmt.Sprintf("Value: %v", n.Value))
		}
	}

	if n.Args != nil {
		attrs = append(attrs, fmt.Sprintf("Args: %v", n.Args))
	}

	if n.Children != nil {
		attrs = append(attrs, fmt.Sprintf("Children: %v", n.Children))
	}

	return fmt.Sprintf("{ %s }", strings.Join(attrs, ", "))
}

func (n *Node) withPos(c *current) Node {
	n.Ref = &Location{
		filename: c.globalStore["filename"].(string),
		line:     c.pos.line,
		col:      c.pos.col,
		offset:   c.pos.offset,
	}
	return *n
}

type NodeSet []Node // bless this with functional shims

func (ns NodeSet) String() string {
	els := make([]string, len(ns))
	for i, l := 0, len(ns); i < l; i++ {
		els[i] = ns[i].String()
	}
	return fmt.Sprintf("[ %s ]", strings.Join(els, ", "))
}

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
