package dep

import (
	"fmt"

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

// NewDependencyGraph returns a DependencyGraph. If 'g' is nil, a new graph is created
// with 'simple.NewDirectedGraph'.
func NewDependencyGraph(g *simple.DirectedGraph) *DependencyGraph {
	if g != nil {
		return &DependencyGraph{g, nil, nil}
	}
	return &DependencyGraph{simple.NewDirectedGraph(), nil, nil}
}

// Sort returns the topological order computed through 'topo.Sort', which implements reversed
// Tarjan SCC. On failure, the set or sets of nodes that are in directed cycles are provided,
// i.e. circular dependencies in the graph.
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

// IsMid returns true if the given id corresponds to a mid node (i.e. neither a root nor a leaf).
//
// Complexity: O(2) or O(V)
func (d *DependencyGraph) IsMid(id int64) bool {
	return !(d.IsRoot(id) || d.IsLeaf(id))
}

// Induce induces a different subgraph for each of the given nodes of the dependency
// graph, where the node is a root (walk forward) or a leaf (reverse walk).
//
// Complexity: O(len(ns) * (V + E))
func (d *DependencyGraph) Induce(ns map[int64]graph.Node) map[int64]*DependencyGraph {
	return d.induce(ns, func(n graph.Node) bool { return d.IsRoot(n.ID()) }, func(n graph.Node) bool { return d.IsLeaf(n.ID()) })
}

// InduceDir induces a different subgraph for each of the given nodes of the dependency
// graph, where the node can be any vertex (root, leaf or mid). It is possible to
// traverse the graph forward, in reverse or in both directions. Not that the following
// contexts will produce a subgraph with a single node:
//  - a root node with reverse walk only
//  - a leaf node with forward walk only
//  - any node with neither forward nor reverse walk. a warning is shown in this case.
//
// Complexity: O(len(ns) * (V + E))
func (d *DependencyGraph) InduceDir(ns map[int64]graph.Node, fw, rv bool) map[int64]*DependencyGraph {
	return d.induce(ns, func(n graph.Node) bool { return fw }, func(n graph.Node) bool { return rv })
}

func (d *DependencyGraph) induce(ns map[int64]graph.Node, fw, rv func(n graph.Node) bool) map[int64]*DependencyGraph {
	d.rootsAndLeafs()
	o := make(map[int64]*DependencyGraph)
	var w Inducer
	for _, u := range ns {
		w.Induce(d, u, fw(u), rv(u))
		o[u.ID()] = w.Graph
	}
	return o
}

// Inducer is a wrapper around 'traverse.DepthFirst' to induce a subgraph from a dependency
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

// Induce walks a graph starting from a given node and generates a subgraph by copying all the
// traversed edges (and the nodes touched by them). 'fw' and 'rv' allow to select whether a
// forward walk is executed, a reverse walk or both of them.
func (i *Inducer) Induce(d *DependencyGraph, u graph.Node, fw, rv bool) {
	if !(fw || rv) {
		fmt.Println("Induce called with fw=false and rv=false. This might be a mid vertex. Consider using InduceDir instead")
	}
	i.Reset()
	i.Graph = NewDependencyGraph(nil)
	i.Graph.AddNode(u)
	i.Traverse = func(e graph.Edge) bool { i.Graph.SetEdge(e); return true }
	if fw {
		i.Walk(d, u, nil)
		if rv {
			i.Reset()
		}
	}
	if rv {
		i.Walk(reversed{d}, u, nil)
	}
}
