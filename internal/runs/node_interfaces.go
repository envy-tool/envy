package runs

import (
	"envy.pw/cli/internal/addrs"
	"envy.pw/cli/internal/graphs"
)

// referrerNode is an interface implemented by graph nodes that can refer
// to referenceable addresses.
type referrerNode interface {
	graphs.Node
	ReferenceAddrs() []addrs.Referenceable
}
