package graphs

import (
	"envy.pw/cli/internal/addrs"
)

// CommandNode is a Node representing a Command.
//
// Other packages may embed this type in another struct to create a specialized
// CommandNode.
type CommandNode struct {
	addr addrs.Command
	isGraphNode
}

// Addr returns the address of the command the node represents.
func (n *CommandNode) Addr() addrs.Command {
	return n.addr
}

// ReferenceableAddr is the implementation of ReferenceableNode.
func (n *CommandNode) ReferenceableAddr() addrs.Referenceable {
	return n.Addr()
}

// HelperNode is a Node representing a Helper.
type HelperNode struct {
	addr addrs.Helper
	isGraphNode
}

// Addr returns the address of the helper the node represents.
func (n *HelperNode) Addr() addrs.Helper {
	return n.addr
}

// ReferenceableAddr is the implementation of ReferenceableNode.
func (n *HelperNode) ReferenceableAddr() addrs.Referenceable {
	return n.Addr()
}
