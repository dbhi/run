package dep

import (
	"fmt"
	"os"

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
//
// Complexity: O(V,E)
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

// cleanRootsAndLeafs sets the internal fields 'roots' and 'leafs' empty.
//
// Complexity: O(2)
func (d *DependencyGraph) cleanRootsAndLeafs() {
	d.roots = make(map[int64]graph.Node, 0)
	d.leafs = make(map[int64]graph.Node, 0)
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

// This does not work if defined outside of 'Filter', why?
// func (i *Inducer) Traverse(e graph.Edge) bool { i.Graph.SetEdge(e); return true }
/*
// This should allow to get the topological order from the DFS walk. However, we found
// it not to provide valid result consistently. We are using topo.Sort until we guess
// how to get a naive but valid order from here.
// Also, this does not work if defined outside of 'Filter'
func (i *Inducer) Visit(v graph.Node) {
	i.Graph.schedule = append(i.Graph.schedule, v.ID())              // forward
	i.Graph.schedule = append([]int64{v.ID()}, i.Graph.schedule...)  // reverse
}
*/

// InduceAllIn induces a single subgraph that contain all and only the nodes given
// in the list, and edges that have both ends between any of them.
//
// Complexity: ????
func (d *DependencyGraph) InduceAllIn(ns map[int64]graph.Node) *DependencyGraph {
	fmt.Println("This method is not implemented yet")
	os.Exit(1)
	return nil
}

// Reduce returns the transitive reduction of the dependency graph.
//
// Complexity: ????
func (d *DependencyGraph) Reduce() *DependencyGraph {
	//https://cs.stackexchange.com/questions/7096/transitive-reduction-of-dag
	fmt.Println("This method is not implemented yet")
	os.Exit(1)
	return nil
}

// Closure return the transitive closure of the dependency graph.
//
// Complexity: ????
func (d *DependencyGraph) Closure() *DependencyGraph {
	/*
		//A navie implementation can be:
		for _, from := range graph.NodesOf(d.Nodes()) {
			for _, to := range graph.NodesOf(d.Nodes()) {
				if from != to && topo.PathExistsIn(d, from, to) && !simple.HasEdgeFromTo(from, to) {
					d.SetEdge(d.Edge(from, to))
				}
			}
		}

		// But it's so bad that it is not worth uncommenting it.
		// On the one hand, there is need to analyse N element in each loop
		// On the other hand *PathExistsIn exists as a helper function. If many tests for path existence
		// are being performed, other approaches will be more efficient. *
	*/
	fmt.Println("This method is not implemented yet")
	os.Exit(1)
	return nil
}

// Schedule returns a solution to the topological sorting of the dependency graph for:
// [`p=0`] return unexported field 'schedule' if not nil,
// [`p>1`] 'p' parallel/concurrent execution threads, or
// [`p=-1`] return as many threads as the maximum number of concurrently executable tasks in the graph.
//
// Complexity: ????
func (d *DependencyGraph) Schedule(p int64) [][]int64 {
	fmt.Println("This method is not implemented yet")
	os.Exit(1)
	return nil
}

/*
// This is an alternative implementation of Induce, which does not use `gonum/graph/traverse`.
func subFromWalk(g graph.Directed) map[string]graph.Directed {
	addNode := func(s graph.DirectedGraph, n graph.Node) {
		if s.Node(n.ID()) == nil {
			s.AddNode(n)
		}
	}
	var walk func(g, s graph.DirectedGraph, n graph.Node) []graph.Node
	walk = func(g, s graph.DirectedGraph, n graph.Node) []graph.Node {
		addNode(s, n)
		if l := graph.NodesOf(g.To(n.ID())); len(l) > 0 {
			for _, x := range l {
				addNode(s, x)
				s.SetEdge(g.Edge(x.ID(), n.ID()))
				walk(g, s, x)
			}
		}
		return nil
	}
	o := make(map[string]graph.DirectedGraph)
	for _, n := range getTargets(g) {
		s := simple.NewDirectedGraph()
		walk(g, s, n)
		o[n.(*dot.DotNode).DOTID()] = s
	}
	return o
}
*/
