package dsl

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// recursively selects Node values into a flattened slice
func searchNodes(v interface{}) NodeSet {
	if nil == v {
		return NodeSet{}
	}

	vars := v.([]interface{})
	nodes := make(NodeSet, 0)

	for _, val := range vars {
		n, isNode := val.(Node)

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
func delimitedNodeSlice(first, rest interface{}) NodeSet {
	res := make(NodeSet, 1)
	res[0] = first.(Node)

	if nil != rest {
		res = append(res, searchNodes(rest)...)
	}

	return res
}

// a convenience function to nil-check, returning an empty
// Node slice in place of nil
func defaultToEmptySlice(nodes interface{}) NodeSet {
	if nil == nodes {
		return NodeSet{}
	}

	return nodes.(NodeSet)
}

/**
 * Parses date and date + timestamp in ISO-8601 variations just like
 * JavaScript. Specifically:
 *
 * YYYY-MM-DD
 * YYYY-mm-ddTHH:MM:SS
 * YYYY-mm-ddTHH:MM:SSZ
 * YYYY-mm-ddTHH:MM:SS-0000
 * YYYY-mm-ddTHH:MM:SS-00:00
 */
func ParseDateLikeJS(tstamp string) (time.Time, error) {
	// you'll just have to take my word on this
	re := regexp.MustCompile("^([0-9]{4}-[0-9]{2}-[0-9]{2})(?:(T[0-9]{2}:[0-9]{2}:[0-9]{2})(Z|(?:[+-][0-9]{2}:?[0-9]{2}))?)?$")

	format := "2006-01-02" // default to parsing only the date

	m := re.FindStringSubmatch(tstamp)

	if m == nil {
		return time.Time{}, fmt.Errorf("Not a parsable timestamp: %s", tstamp)
	}

	parts := []string{m[1], m[2], strings.Replace(strings.Replace(m[3], ":", "", -1), "Z", "", -1)}

	if m[2] != "" {
		format = format + "T15:04:05"
	}

	if m[3] != "" && m[3] != "Z" {
		format = format + "-0700"
	}

	return time.Parse(format, strings.Join(parts, ""))
}
