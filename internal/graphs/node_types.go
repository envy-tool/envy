package graphs

import (
	"envy.pw/cli/internal/addrs"
)

// CommandNode is a Node representing a Command.
//
// Other packages may embed this type in another struct to create a specialized
// CommandNode.
type CommandNode struct {
	Addr addrs.Command
	graphNodeImpl
}

var _ Node = (*CommandNode)(nil)

// ReferenceableAddr is the implementation of ReferenceableNode.
func (n *CommandNode) ReferenceableAddr() addrs.Referenceable {
	return n.Addr
}

// HelperNode is a Node representing a Helper.
type HelperNode struct {
	Addr addrs.Helper
	graphNodeImpl
}

var _ Node = (*HelperNode)(nil)

// ReferenceableAddr is the implementation of ReferenceableNode.
func (n *HelperNode) ReferenceableAddr() addrs.Referenceable {
	return n.Addr
}
