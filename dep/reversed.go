package dep

import (
	"gonum.org/v1/gonum/graph"
)

// reversed provides enough of the traverse.Graph interface to get what we want here.
type reversed struct {
	graph.Directed
}

func (g reversed) From(id int64) graph.Nodes     { return g.Directed.To(id) }
func (g reversed) Edge(u, v int64) graph.Edge    { return g.Directed.Edge(v, u) }
func (g reversed) HasEdgeFromTo(u, v int64) bool { return g.Directed.HasEdgeFromTo(v, u) }
func (g reversed) To(id int64) graph.Nodes       { return g.Directed.From(id) }
