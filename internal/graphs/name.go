package graphs

import (
	"fmt"
)

// NodeName returns a name for a node that is mainly intended to help with
// debugging.
func NodeName(node Node) string {
	var baseName string
	switch tn := node.(type) {
	case ReferenceableNode:
		baseName = tn.ReferenceableAddr().String()
	case interface { NodeName() string }:
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
