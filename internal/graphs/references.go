package graphs

import (
	"envy.pw/cli/internal/addrs"
	"envy.pw/cli/internal/configs"
	"envy.pw/cli/internal/nvdiags"
)

// ReferenceableNode is an interface implemented by Nodes that can be
// referenced in expressions.
type ReferenceableNode interface {
	Node
	ReferenceableAddr() addrs.Referenceable
}

// ReferrerNode is an interface implemented by Nodes that can refer to objects
// using referenceable addresses.
type ReferrerNode interface {
	Node
	References() []configs.Reference
}

// A ReferenceableNodeFactory is a factory function that can produce a graph
// node for the referent of a given reference.
//
// The factory returns nil with no error diagnostics to signal that no node
// or edges are required for the given address.
//
// The factory returns error diagnostics to indicate that the given address
// is somehow invalid.
//
// The factory must make decisions based only on the address of the given
// reference. The source range of the reference is included only for use when
// creating diagnostics.
type ReferenceableNodeFactory func(referrer addrs.Referenceable, ref configs.Reference) (Node, nvdiags.Diagnostics)

// NodeReferences is a helper to get the references for a node that may or
// may not implement ReferrerNode.
func NodeReferences(node Node) []configs.Reference {
	if rn, ok := node.(ReferrerNode); ok {
		return rn.References()
	}
	return nil
}

// NodeReferenceableAddr is a helper to get the referenceable address for
// a node if it has one, or return nil otherwise.
func NodeReferenceableAddr(node Node) addrs.Referenceable {
	if rn, ok := node.(ReferenceableNode); ok {
		return rn.ReferenceableAddr()
	}
	return nil
}

// AddWithReferents adds the given start node along with nodes representing
// objects it refers to directly or indirectly and edges representing the
// dependencies implied by those references.
func (g *Graph) AddWithReferents(start Node, factory ReferenceableNodeFactory) nvdiags.Diagnostics {
	g.l.Lock()
	defer g.l.Unlock()

	g.nodes.Add(start)

	nodes := make(map[addrs.Referenceable]Node)

	// The graph might already contain some referencable nodes.
	for n := range g.nodes {
		addr := NodeReferenceableAddr(n)
		if addr == nil {
			continue
		}
		nodes[addr] = n
	}

	return g.addReferents(start, nodes, factory)
}

func (g *Graph) addReferents(current Node, nodes map[addrs.Referenceable]Node, factory ReferenceableNodeFactory) nvdiags.Diagnostics {
	var diags nvdiags.Diagnostics
	refs := NodeReferences(current)

	for _, ref := range refs {
		target, exists := nodes[ref.Addr]
		if target == nil && !exists {
			newTarget, moreDiags := factory(NodeReferenceableAddr(current), ref)
			diags = diags.Append(moreDiags)
			if moreDiags.HasErrors() {
				continue
			}
			target = newTarget
			nodes[ref.Addr] = newTarget
		}
		if target == nil {
			continue
		}

		g.connect(current, target) // implicitly adds target if it isn't already present
		g.addReferents(target, nodes, factory)
	}

	return diags
}
