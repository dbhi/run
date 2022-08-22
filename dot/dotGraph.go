package dot

import (
	"fmt"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
)

// Graph ensures dotEdge and Node types are created when needed.
type Graph struct {
	*simple.DirectedGraph
}

// NewNode returns a new Node to be added to g. The Node contains the attributes.
func (g Graph) NewNode() graph.Node {
	return &Node{Node: g.DirectedGraph.NewNode(), attrs: make(map[string]string)}
}

// NewNode returns a new Edge to be added to g. The Edge contains the attributes.
func (g Graph) NewEdge(from, to graph.Node) graph.Edge {
	return &dotEdge{Edge: g.DirectedGraph.NewEdge(from, to), attrs: make(map[string]string)}
}

// GetNodeByDOTID gets a Node by it's DOT ID.
func (g Graph) GetNodeByDOTID(DOTID string) graph.Node {
	for _, n := range graph.NodesOf(g.Nodes()) {
		if n.(*Node).DOTID() == DOTID {
			return n
		}
	}
	return nil
}

// Node handles basic DOT serialisation and deserialisation
type Node struct {
	graph.Node
	dotID string
	attrs map[string]string
}

// SetAttribute sets a DOT attribute.
func (n *Node) SetAttribute(attr encoding.Attribute) error {
	n.attrs[attr.Key] = attr.Value
	return nil
}

// DOTID gets the DOT attribute.
func (n *Node) DOTID() string { return n.dotID }

// SetDOTID sets the DOT ID.
func (n *Node) SetDOTID(id string) { n.dotID = id }

// Attributes gets the slice of attributes defined for the node.
func (n *Node) Attributes() []encoding.Attribute {
	attrs := make([]encoding.Attribute, 0, len(n.attrs))
	for k, v := range n.attrs {
		attrs = append(attrs, encoding.Attribute{Key: k, Value: v})
	}
	return attrs
}

// Attribute gets a DOT attribute.
func (n *Node) Attribute(key string) (string, error) {
	a, ok := n.attrs[key]
	if ok {
		return a, nil
	}
	return "", fmt.Errorf("key %s not found in attrs", key)
}

// dotEdge handles basic DOT serialisation and deserialisation
// and edge reversal.
type dotEdge struct {
	graph.Edge
	attrs map[string]string
}

// SetAttribute sets a DOT attribute.
func (e *dotEdge) SetAttribute(attr encoding.Attribute) error {
	e.attrs[attr.Key] = attr.Value
	return nil
}

// Attributes gets the slice of attributes defined for the edge.
func (e *dotEdge) Attributes() []encoding.Attribute {
	attrs := make([]encoding.Attribute, 0, len(e.attrs))
	for k, v := range e.attrs {
		attrs = append(attrs, encoding.Attribute{Key: k, Value: v})
	}
	return attrs
}

// ReversedEdge returns a new Edge with the end point of the edges in the pair swapped.
func (e *dotEdge) ReversedEdge() graph.Edge {
	return &dotEdge{Edge: e.Edge.ReversedEdge(), attrs: e.attrs}
}

// Unmarshal reads a slice of bytes containing a dot encoded graph and returns a
// *simple.DirectedGraph with edges and nodes containing attributes
func Unmarshal(b []byte) *simple.DirectedGraph {
	g := simple.NewDirectedGraph()
	if err := dot.Unmarshal(b, Graph{g}); err != nil {
		return nil
	}
	return g
}

// Marshal returns the DOT encoding for the graph g
func Marshal(g graph.Graph) []byte {
	b, err := dot.Marshal(g, "", "", "")
	if err != nil {
		return nil
	}
	return b
}
