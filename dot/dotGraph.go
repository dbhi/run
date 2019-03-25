package dot

import (
	"fmt"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
)

// DotGraph ensures dotEdge and DotNode types are created when needed.
type DotGraph struct {
	*simple.DirectedGraph
}

// NewNode returns a new unique Node to be added to g. The Node contains the attributes.
func (g DotGraph) NewNode() graph.Node {
	return &DotNode{Node: g.DirectedGraph.NewNode(), attrs: make(map[string]string)}
}

func (g DotGraph) NewEdge(from, to graph.Node) graph.Edge {
	return &dotEdge{Edge: g.DirectedGraph.NewEdge(from, to), attrs: make(map[string]string)}
}

func (g DotGraph) GetNodeByDOTID(DOTID string) graph.Node {
	for _, n := range graph.NodesOf(g.Nodes()) {
		if n.(*DotNode).DOTID() == DOTID {
			return n
		}
	}
	return nil
}

// DotNode handles basic DOT serialisation and deserialisation
type DotNode struct {
	graph.Node
	dotID string
	attrs map[string]string
}

// SetAttribute sets a DOT attribute.
func (n *DotNode) SetAttribute(attr encoding.Attribute) error {
	n.attrs[attr.Key] = attr.Value
	return nil
}

func (n *DotNode) DOTID() string      { return n.dotID }
func (n *DotNode) SetDOTID(id string) { n.dotID = id }
func (n *DotNode) Attributes() []encoding.Attribute {
	attrs := make([]encoding.Attribute, 0, len(n.attrs))
	for k, v := range n.attrs {
		attrs = append(attrs, encoding.Attribute{Key: k, Value: v})
	}
	return attrs
}

// Attribute gets a DOT attribute.
func (n *DotNode) Attribute(key string) (string, error) {
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

func (e *dotEdge) Attributes() []encoding.Attribute {
	attrs := make([]encoding.Attribute, 0, len(e.attrs))
	for k, v := range e.attrs {
		attrs = append(attrs, encoding.Attribute{Key: k, Value: v})
	}
	return attrs
}

func (e *dotEdge) ReversedEdge() graph.Edge {
	return &dotEdge{Edge: e.Edge.ReversedEdge(), attrs: e.attrs}
}

// Unmarshal reads a slice of bytes containing a dot encoded graph and returns a
// *simple.DirectedGraph with edges and nodes containing attributes
func Unmarshal(b []byte) *simple.DirectedGraph {
	g := simple.NewDirectedGraph()
	if err := dot.Unmarshal(b, DotGraph{g}); err != nil {
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
