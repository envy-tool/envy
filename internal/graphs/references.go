package graphs

import (
	"github.com/envy-tool/envy/internal/addrs"
)

// ReferenceableNode is an interface implemented by Nodes that can be
// referenced in expressions.
type ReferenceableNode interface {
	Node
	ReferenceableAddr() addrs.Referenceable
}
