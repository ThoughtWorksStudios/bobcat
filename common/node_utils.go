package common

// recursively selects Node values into a flattened slice
func searchNodes(v interface{}) NodeSet {
	if nil == v {
		return NodeSet{}
	}

	if ns, ok := v.(NodeSet); ok {
		return ns
	}

	vars := v.([]interface{})
	nodes := make(NodeSet, 0)

	for _, val := range vars {
		n, isNode := val.(*Node)

		if isNode {
			nodes = append(nodes, n)
		} else {
			more, isSlice := val.([]interface{})
			if isSlice {
				nodes = append(nodes, searchNodes(more)...)
			}
		}

	}

	return nodes
}

// convenience function to join a single Node with a
// Node slice representing 0 or more Node values; Often
// used to handle arguments, filtering out whitespace and
// delimiter matches
func DelimitedNodeSlice(first, rest interface{}) NodeSet {
	res := make(NodeSet, 1)
	res[0] = first.(*Node)

	if nil != rest {
		res = append(res, searchNodes(rest)...)
	}

	return res
}

// a convenience function to nil-check, returning an empty
// Node slice in place of nil
func DefaultToEmptySlice(nodes interface{}) NodeSet {
	if nil == nodes {
		return NodeSet{}
	}

	return nodes.(NodeSet)
}
