package flow

import (
	"sync"
	"context"

	"envy.pw/cli/internal/graphs"
	"envy.pw/cli/internal/states"
)

// Change is a named type used to send notifications to and from
// graph nodes participating in a flow.
type Change struct{}

// Run begins a flow for the given graph and blocks until the given context
// expires.
func Run(ctx context.Context, graph *graphs.Graph, state *states.State) {
	type Channels struct {
		In chan<- Change 
		Out <-chan Change
	}
	nodeChans := make(map[graphs.Node]Channels)
	var leaves []chan<- Change

	running := true

	var wg sync.WaitGroup
	nodes := graph.Nodes()
	for n := range nodes {
		fn, ok := n.(Node)
		if !ok {
			continue // this node does not participate in flows
		}

		toNode := make(chan Change)
		fromNode := make(chan Change)

		nodeChans[fn] = Channels{
			In: toNode,
			Out: fromNode,
		}

		wg.Add(1)
		go func (in <-chan Change, out chan<- Change) {
			for {
				fn.Flow(state, in, out)
				// If we've not been cancelled yet then the node has returned
				// early and we must start it up again. If we _have_ been
				// cancelled then it's time to exit.
				if !running {
					break
				}
			}
			close(out)
			wg.Done()
		}(toNode, fromNode)

		if graph.IsLeaf(n) {
			leaves = append(leaves, toNode)
		}
	}

	// Now that each channel has one input and one output channel, we need
	// some additional goroutines/channels to handle the fan-out and fan-in
	// to propagate a notification from a source node to each of its referrers.
	for sn := range nodes {
		sourceChans, ok := nodeChans[sn]
		if !ok {
			continue
		}
		var propagateChans []chan<- Change
		for tn := range graph.Referrers(sn) {
			targetChans, ok := nodeChans[tn]
			if !ok {
				continue
			}

			edgeCh := make(chan Change)
			propagateChans = append(propagateChans, edgeCh)
			go relay(edgeCh, targetChans.In)
		}
		go relayAll(sourceChans.Out, propagateChans)
	}

	// Send an initial "change" to all of the leaves to get the flow started.
	for _, ch := range leaves {
		ch <- Change{}
	}

	// Wait for cancellation
	<-ctx.Done()
	running = false

	// Close all of the leaves to begin graceful shutdown. Closure should
	// then propagate through the flow.
	for _, ch := range leaves {
		close(ch)
	}

	// Wait until all of the flow nodes have exited.
	wg.Wait()
}

func relay(in <-chan Change, out chan<- Change) {
	for change := range in {
		out <- change
	}

	// FIXME: This isn't correct... multiple goroutines can share the same
	// output channel here, so this will panic trying to close the same
	// channel multiple times. Should pass a waitgroup in here instead,
	// so that a different goroutine can be responsible for closing the
	// channel just once when all of these are done.
	close(out)
}

func relayAll(in <-chan Change, out []chan<- Change) {
	for change := range in {
		for _, ch := range out {
			ch <- change
		}
	}
	for _, ch := range out {
		close(ch)
	}
}
