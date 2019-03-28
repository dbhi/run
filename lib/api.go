package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/umarcor/run/dep"
	"github.com/umarcor/run/dot"
	"gonum.org/v1/gonum/graph"
)

func ReadFile(f string) ([]byte, error) {
	if len(f) == 0 {
		src := `strict digraph {
// Node definitions.
A [label="yellow"];
B [label="green"];
C [label="red"];
D [label="blue"];
E [label="magenta"];
F [label="purple"];
// Edge definitions.
A -> B;
A -> C -> E -> F;
C -> D;
B -> F;
B -> E;
}`
		fmt.Println("Empty file path! Please provide a DOT file")
		fmt.Println("Using the following content as an example:")
		fmt.Println(src)
		return []byte(src), nil
		//return nil, fmt.Errorf("Empty file path! Please provide a DOT file")
	}
	jf, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer jf.Close()
	return ioutil.ReadAll(jf)
}

func WriteFile(f string, g *dep.DependencyGraph) error {
	if b := dot.Marshal(g); b != nil {
		fmt.Printf("Writing graph to '%s'\n", f)
		return ioutil.WriteFile(f, b, 0644)
	}
	return fmt.Errorf("Marshal failed for file %s", f)
}

func PrintGraph(g *dep.DependencyGraph) {
	for _, n := range graph.NodesOf(g.Nodes()) {
		fmt.Printf("%+v\n", n)
	}
	for _, e := range graph.EdgesOf(g.Edges()) {
		fmt.Printf("%+v\n", e)
	}
}

func GetSubgraphs(f string, rv bool) (map[string]*dep.DependencyGraph, error) {
	b, err := ReadFile(f)
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

	//If you add a `func (n dotNode) String() string { return n.DOTID() }` method, then fmt.Sprint(n) or similar with give you the name in the printout.

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
