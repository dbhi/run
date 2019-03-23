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

func GetSubgraphs(f string, rv bool) (map[int64]*dep.DependencyGraph, error) {
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
	return s, nil
}

func GetTaskListAll(d map[int64]*dep.DependencyGraph) [][]string {
	for k, j := range d {
		fmt.Printf("\nSUBGRAPH [for leaf vertex %d]: %+v\n", k, j)
		str := GetTaskList(j)
		if str != nil {
			fmt.Println(str)
		}
	}
	return nil
}

func GetTaskList(g *dep.DependencyGraph) []string {
	PrintGraph(g)
	s, err := g.Sort()
	if err != nil {
		log.Fatal(err)
	}
	for _, n := range s {
		d := n.(*dot.DotNode)
		if a, e := d.Attribute("shape"); e == nil {
			if a == "box" {
				fmt.Println(d.DOTID())
				//fmt.Println("SHAPE:", a)
			}
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
