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
	Count     NodeSet
	Related  *Node
	Children NodeSet
	Ref      *Location
}

func (n Node) String() string {
	attrs := make([]string, 1)

	attrs[0] = fmt.Sprintf("Kind: %s", strconv.Quote(n.Kind))

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

	if n.Count != nil {
		attrs = append(attrs, fmt.Sprintf("Count: %v", n.Count))
	}

	if n.Related != nil {
		attrs = append(attrs, fmt.Sprintf("Related: %v", n.Related))
	}

	if n.Children != nil {
		attrs = append(attrs, fmt.Sprintf("Children: %v", n.Children))
	}

	return fmt.Sprintf("{ %s }", strings.Join(attrs, ", "))
}

func (n *Node) HasRelation() bool {
	return n.Related != nil
}

func (n *Node) ValNode() Node {
	return n.Value.(Node)
}

func (n *Node) ValStr() string {
	return n.Value.(string)
}

func (n *Node) ValInt() int64 {
	return n.Value.(int64)
}

func (n *Node) ValFloat() float64 {
	return n.Value.(float64)
}

func (n *Node) ValTime() time.Time {
	return n.Value.(time.Time)
}

func (n *Node) withPos(c *current) Node {
	if nil != c {
		filename, _ := c.globalStore["filename"].(string)
		n.Ref = NewLocation(
			filename,
			c.pos.line,
			c.pos.col,
			c.pos.offset,
		)
	}
	return *n
}

func (n *Node) Err(msg string, tokens ...interface{}) error {
	if nil == n.Ref {
		return fmt.Errorf(msg, tokens...)
	} else {
		format := fmt.Sprintf("%v %s", n.Ref, msg)
		return fmt.Errorf(format, tokens...)
	}
}

func (n *Node) WrapErr(inner error) error {
	return n.Err(inner.Error())
}

type NodeSet []Node // bless this with functional shims

func (ns NodeSet) String() string {
	els := make([]string, len(ns))
	for i, l := 0, len(ns); i < l; i++ {
		els[i] = ns[i].String()
	}
	return fmt.Sprintf("[ %s ]", strings.Join(els, ", "))
}

type IterEnv struct {
	Self NodeSet
	Idx  int
	Halt func()
}

func mkEnv(self NodeSet, halt func()) *IterEnv {
	return &IterEnv{Self: self, Idx: 0, Halt: halt}
}

type Iterator func(env *IterEnv, node Node)
type Collector func(env *IterEnv, node Node) interface{}

func (nodes NodeSet) Each(f Iterator) NodeSet {
	abort := false
	env := mkEnv(nodes, func() { abort = true })

	for i, size := 0, len(nodes); i < size; i++ {
		if abort {
			break
		}
		env.Idx = i
		f(env, nodes[i])
	}

	return nodes
}

func (nodes NodeSet) Map(f Collector) []interface{} {
	size := len(nodes)
	result := make([]interface{}, size)
	nodes.Each(func(env *IterEnv, node Node) {
		result[env.Idx] = f(env, node)
	})
	return result
}

type Location struct {
	line, col, offset int
	filename          string
}

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
