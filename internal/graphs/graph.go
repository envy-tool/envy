package graphs // import "envy.pw/cil/internal/graphs"

import (
	"sync"

	"envy.pw/cli/internal/addrs"
)

// Graph represents a directed graph.
//
// Operations on graphs are concurrency-safe.
type Graph struct {
	nodes    NodeSet
	edgesOut map[Node]NodeSet
	edgesIn  map[Node]NodeSet
	l        sync.RWMutex
}

// Node represents a single node in a graph.
type Node interface {
	isGraphNode()
}

type graphNodeImpl struct{} // embed this in graph node types

func (n graphNodeImpl) isGraphNode() {}

// NodeSet is a mutable set of nodes.
type NodeSet map[Node]struct{}

func NewGraph() *Graph {
	return &Graph{
		nodes:    make(NodeSet),
		edgesIn:  make(map[Node]NodeSet),
		edgesOut: make(map[Node]NodeSet),
	}
}

// HasNode returns true if and only if the graph has the given node.
func (g *Graph) HasNode(node Node) bool {
	g.l.RLock()
	ok := g.nodes.Has(node)
	g.l.RUnlock()
	return ok
}

// HasEdge returns true if and only if the graph has an edge from the first
// given node to the second given node.
func (g *Graph) HasEdge(from, to Node) bool {
	g.l.RLock()
	_, ok := g.edgesOut[from][to]
	g.l.RUnlock()
	return ok
}

// Nodes returns a set of all of the nodes in the graph.
func (g *Graph) Nodes() NodeSet {
	g.l.RLock()
	defer g.l.RUnlock()
	return g.nodes.Copy()
}

// Referents returns a set of the nodes that the given node refers to. That is,
// the nodes at the "to" end of edges starting from the given node.
func (g *Graph) Referents(node Node) NodeSet {
	g.l.RLock()
	defer g.l.RUnlock()
	return g.edgesOut[node].Copy()
}

// Referrers returns a set of the nodes that refer to the given node. That is,
// the nodes at the "from" end of edges ending at the given node.
func (g *Graph) Referrers(node Node) NodeSet {
	g.l.RLock()
	defer g.l.RUnlock()
	return g.edgesIn[node].Copy()
}

// Roots returns a set of the nodes that are not referred to by any other nodes.
func (g *Graph) Roots() NodeSet {
	ret := make(NodeSet)
	g.l.RLock()
	for n := range g.nodes {
		if len(g.edgesIn[n]) == 0 {
			ret.Add(n)
		}
	}
	g.l.RUnlock()
	return ret
}

// Leaves returns a set of the nodes that do not refer to any other nodes.
func (g *Graph) Leaves() NodeSet {
	ret := make(NodeSet)
	g.l.RLock()
	for n := range g.nodes {
		if len(g.edgesOut[n]) == 0 {
			ret.Add(n)
		}
	}
	g.l.RUnlock()
	return ret
}

// IsLeaf returns true if and only if the given node does not refer to any
// other nodes.
func (g *Graph) IsLeaf(node Node) bool {
	g.l.RLock()
	defer g.l.RUnlock()
	return len(g.edgesOut[node]) == 0
}

// IsRoot returns true if and only if the given node is not referred to by
// any other nodes.
func (g *Graph) IsRoot(node Node) bool {
	g.l.RLock()
	defer g.l.RUnlock()
	return len(g.edgesIn[node]) == 0
}

// AddNode inserts a new node into the graph, with no incoming or outgoing
// edges.
//
// If the node is already present then this is a no-op.
func (g *Graph) AddNode(node Node) {
	g.l.Lock()
	g.nodes.Add(node)
	g.l.Unlock()
}

// Connect creates a new edge from the first given node to the second given
// node. If such an edge already exists, this is a no-op.
//
// If either of the given nodes is not already in the graph, it is first
// implicitly added.
func (g *Graph) Connect(from, to Node) {
	g.l.Lock()
	g.connect(from, to)
	g.l.Unlock()
}

// Disconnect removes the edge from the first given node to the second given
// node, or does nothing if no such edge exists.
func (g *Graph) Disconnect(from, to Node) {
	g.l.Lock()
	g.connect(from, to)
	g.l.Unlock()
}

func (g *Graph) connect(from, to Node) {
	if _, ok := g.edgesOut[from]; !ok {
		g.edgesOut[from] = make(NodeSet)
	}
	if _, ok := g.edgesIn[to]; !ok {
		g.edgesIn[to] = make(NodeSet)
	}
	g.nodes.Add(from)
	g.nodes.Add(to)
	g.edgesOut[from].Add(to)
	g.edgesIn[to].Add(from)
}

func (g *Graph) disconnect(from, to Node) {
	ns, ok := g.edgesOut[from]
	if !ok {
		return // nothing to do
	}
	delete(ns, to)
	delete(g.edgesIn[to], from) // we assume edgesIn and edgesOut will always be consistent
}

// NodesForAddr searches the graph for nodes that represent the given referencable
// address, returning a set of them.
func (g *Graph) NodesForAddr(addr addrs.Referenceable) NodeSet {
	g.l.RLock()
	ret := make(NodeSet)
	for n := range g.nodes {
		rn, ok := n.(ReferenceableNode)
		if !ok {
			continue
		}
		if rn.ReferenceableAddr() == addr {
			ret.Add(n)
		}
	}
	g.l.RUnlock()
	return ret
}

// Has returns true if and only if the given edge is in the set.
func (ns NodeSet) Has(node Node) bool {
	_, ok := ns[node]
	return ok
}

// Add inserts the given node into the set. If the node is already present
// then this is a no-op.
func (ns NodeSet) Add(node Node) {
	ns[node] = struct{}{}
}

// Remove removes the given node from the set. If the node wasn't already
// present then this is a no-op.
func (ns NodeSet) Remove(node Node) {
	delete(ns, node)
}

// Copy allocates a new NodeSet and shallow-copies the notes from the receiver
// into it.
func (ns NodeSet) Copy() NodeSet {
	ret := make(NodeSet, len(ns))
	for k, v := range ns {
		ret[k] = v
	}
	return ret
}
