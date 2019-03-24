package lib

import (
	"fmt"
	"log"

	"github.com/umarcor/run/dep"
	"github.com/umarcor/run/dot"
	"gonum.org/v1/gonum/graph"
)

func PrintGraph(g *dep.DependencyGraph) {
	for _, n := range graph.NodesOf(g.Nodes()) {
		fmt.Printf("%+v\n", n)
	}
	for _, e := range graph.EdgesOf(g.Edges()) {
		fmt.Printf("%+v\n", e)
	}
}

func GetSubgraphs(f string, rv bool) (map[string]*dep.DependencyGraph, error) {
	b, err := dot.ReadFile(f)
	if err != nil {
		return nil, err
	}
	g := dot.Unmarshal(b)
	if g == nil {
		return nil, fmt.Errorf("failed to parse input DOT file")
	}
	d := dep.NewDependencyGraph(g)
	var n map[int64]graph.Node
	if rv {
		n = d.Roots()
	} else {
		n = d.Leafs()
	}
	s := d.Induce(n)
	o := make(map[string]*dep.DependencyGraph)
	for k, n := range s {
		x := dot.DotGraph{n.DirectedGraph}.Node(k).(*dot.DotNode).DOTID()
		o[x] = n
	}
	return o, nil
}

func GetTaskListAll(d map[string]*dep.DependencyGraph) [][]string {
	for k, j := range d {
		fmt.Printf("\nSUBGRAPH [for leaf vertex %d]: %+v\n", k, j)
		str := GetTaskList(j, "")
		if str != nil {
			fmt.Println(str)
		}
	}
	return nil
}

type Task struct {
	ID          int64
	DOTID       string
	Description string
	Cmds        [][]string
	Env         map[string]string
	Sources     map[string]string
	Artifacts   map[string]string
	Results     map[string]string
}

func taskSubgraph(g *dep.DependencyGraph, t string) *dep.DependencyGraph {
	if t == "" {
		return g
	}
	fmt.Println("Filtering subgraph fot task", t)
	r := dot.DotGraph{g.DirectedGraph}.GetNodeByDOTID(t)
	if r == nil {
		log.Fatal("node %s not found in graph with leafs %s", t, g.Leafs())
	}
	// FIXME Check if the task generates some output. If so, let the output node be the target leaf.
	i := g.InduceDir(map[int64]graph.Node{r.ID(): r}, false, true)
	s, ok := i[r.ID()]
	if !ok {
		log.Fatal("subgraph for node %s not found in graph with leafs %s", t, g.Leafs())
	}
	return s
}

func GetTaskList(d *dep.DependencyGraph, t string) []string {

	g := taskSubgraph(d, t)

	//PrintGraph(g)

	s, err := g.Sort()
	if err != nil {
		log.Fatal(err)
	}
	k := 0
	for _, n := range s {
		d := n.(*dot.DotNode)
		if _, e := d.Attribute("shape"); e == nil {
			//if a == "box" {
			fmt.Printf("%d. %s\n", k, d.DOTID())
			//fmt.Println("SHAPE:", a)
			k++
			//}
		}
	}
	//printGraph(j)
	/*
		if b := dot.Marshal(g); b != nil {
			fmt.Println(string(b))
		}
	*/
	return nil
}
