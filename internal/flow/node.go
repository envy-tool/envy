package flow

import (
	"envy.pw/cli/internal/states"
	"envy.pw/cli/internal/graphs"
)

// Node is an interface implemented by graph nodes that participate in
// flows.
type Node interface {
	graphs.Node
	Flow(state *states.State, in <-chan Change, out chan<- Change)
}
