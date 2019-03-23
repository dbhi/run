package dep

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
	"gonum.org/v1/gonum/graph/traverse"
)

// DependencyGraph is a dependency graph; a kind of directed acyclic graph (DAG).
type DependencyGraph struct {
	*simple.DirectedGraph
	roots map[int64]graph.Node
	leafs map[int64]graph.Node
}

// NewDependencyGraph returns a DependencyGraph.
func NewDependencyGraph(g *simple.DirectedGraph) *DependencyGraph {
	if g != nil {
		return &DependencyGraph{g, nil, nil}
	}
	return &DependencyGraph{simple.NewDirectedGraph(), nil, nil}
}

// Sort returns the topological order computed through topo.Sort
func (d *DependencyGraph) Sort() ([]graph.Node, error) {
	return topo.Sort(d)
}

// rootsAndLeafs analyses all the nodes to identify the roots and leafs.
//
// Complexity: O(V)
func (d *DependencyGraph) rootsAndLeafs() {
	if (d.roots == nil) || (d.leafs == nil) {
		d.roots = make(map[int64]graph.Node, 0)
		d.leafs = make(map[int64]graph.Node, 0)
		for _, n := range graph.NodesOf(d.Nodes()) {
			id := n.ID()
			if len(graph.NodesOf(d.From(id))) == 0 {
				d.leafs[id] = n
				continue
			}
			if len(graph.NodesOf(d.To(id))) == 0 {
				d.roots[id] = n
			}
		}
	}
}

// Roots returns a map containing the nodes without incoming edges (i.e. root nodes).
//
// Complexity: O(1) or O(V)
func (d *DependencyGraph) Roots() map[int64]graph.Node {
	if d.roots == nil {
		d.rootsAndLeafs()
	}
	return d.roots
}

// Leafs returns a map containing the nodes without outgoing edges (i.e. leaf nodes).
//
// Complexity: O(1) or O(V)
func (d *DependencyGraph) Leafs() map[int64]graph.Node {
	if d.leafs == nil {
		d.rootsAndLeafs()
	}
	return d.leafs
}

// IsRoot returns true if the given id corresponds to a node which is a root.
//
// Complexity: O(1) or O(V)
func (d *DependencyGraph) IsRoot(id int64) bool {
	if d.roots == nil {
		return len(graph.NodesOf(d.From(id))) == 0
	}
	_, ok := d.roots[id]
	return ok
}

// IsLeaf returns true if the given id corresponds to a node which is a leaf.
//
// Complexity: O(1) or O(V)
func (d *DependencyGraph) IsLeaf(id int64) bool {
	if d.leafs == nil {
		return len(graph.NodesOf(d.From(id))) == 0
	}
	_, ok := d.leafs[id]
	return ok
}

// Induce induces a different subgraph for each of the given nodes of the dependency
// graph, where the node is any vertex:
//  [root] walk forward
//  [leaf] walk reverse
//  [mid]  walk first forward and then walk reverse
//
// Complexity: O(len(ns) * (V + E))
//
// TODO For mid nodes (!leafs && !roots), allow the user to select the direction of the
// walk (forward, reverse or both).
func (d *DependencyGraph) Induce(ns map[int64]graph.Node) map[int64]*DependencyGraph {
	o := make(map[int64]*DependencyGraph)
	var w Inducer
	for _, u := range ns {
		w.Reset()
		w.Graph = NewDependencyGraph(nil)
		w.Graph.AddNode(u)
		if !d.IsLeaf(u.ID()) {
			w.Filter(d, u, false)
			w.Reset()
		}
		if !d.IsRoot(u.ID()) {
			w.Filter(d, u, true)
		}
		o[u.ID()] = w.Graph
	}
	return o
}

// Inducer is a wrapper around traverse.DepthFirst to induce a subgraph from a dependency
// graph; a kind of directed acyclic graph (DAG).
type Inducer struct {
	traverse.DepthFirst
	Graph *DependencyGraph
}

// NewInducer returns a Inducer for a dependency graph.
func NewInducer(d *DependencyGraph) *Inducer {
	var t traverse.DepthFirst
	return &Inducer{DepthFirst: t, Graph: d}
}

// Filter walks a graph starting from a given node and copies all the traversed edges
// and nodes touched by them, to a induced subgraph.
func (i *Inducer) Filter(d *DependencyGraph, u graph.Node, rv bool) {
	i.EdgeFilter = func(e graph.Edge) bool { i.Graph.SetEdge(e); return true }
	if rv { // Walk the graph backwards.
		i.Walk(reversed{d}, u, nil)
	} else {
		i.Walk(d, u, nil)
	}
}
