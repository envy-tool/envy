package graphs

import (
	"fmt"
	"sort"
	"strings"
)

// DebugString produces a multi-line string representation of the graph
// that is intended to aid in debugging. It's not a suitable representation
// to show an end-user.
func (g *Graph) DebugString() string {
	var buf strings.Builder

	var nodes []Node
	for n := range g.Nodes() {
		nodes = append(nodes, n)
	}
	sort.SliceStable(nodes, func(i, j int) bool {
		return NodeDebugName(nodes[i]) < NodeDebugName(nodes[j])
	})

	for _, n := range nodes {
		var to []Node
		for other := range g.Referents(n) {
			to = append(to, other)
		}
		sort.SliceStable(to, func(i, j int) bool {
			return NodeDebugName(to[i]) < NodeDebugName(to[j])
		})
		fmt.Fprintf(&buf, "%s\n", NodeDebugName(n))
		for _, other := range to {
			fmt.Fprintf(&buf, "  %s\n", NodeDebugName(other))
		}
	}

	return buf.String()
}

// NodeDebugName returns a name for a node that is mainly intended to help with
// debugging.
func NodeDebugName(node Node) string {
	var baseName string
	switch tn := node.(type) {
	case ReferenceableNode:
		baseName = tn.ReferenceableAddr().String()
	case interface{ NodeName() string }:
		baseName = tn.NodeName()
	default:
		baseName = fmt.Sprintf("%#v", tn)
	}

	type SubtypeName interface {
		NodeSubtypeName() string
	}
	if tn, ok := node.(SubtypeName); ok {
		return fmt.Sprintf("%s (%s)", baseName, tn.NodeSubtypeName())
	}
	return baseName
}
