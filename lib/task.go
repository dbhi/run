package lib

import (
	"fmt"
	"log"
	"strings"

	"github.com/dbhi/run/dep"
	"github.com/dbhi/run/dot"
)

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

/*
func taskSubgraph(g *dep.DependencyGraph, t string) *dep.DependencyGraph {
	if t == "" {
		return g
	}
	fmt.Println("Filtering subgraph fot task", t)
	r := dot.Graph{g.DirectedGraph}.GetNodeByDOTID(t)
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
*/
/*
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
*/
func GetTaskList(d *dep.DependencyGraph) []string {
	s, err := d.Sort()
	checkErr(err)
	o := make([]string, 0)
	for _, n := range s {
		x := n.(*dot.Node)
		if a, e := x.Attribute("shape"); e == nil && strings.ToLower(a) == "box" {
			o = append(o, x.DOTID())
			continue
		}
		if a, e := x.Attribute("type"); e == nil && strings.ToLower(a) == "job" {
			o = append(o, x.DOTID())
		}
	}
	return o
}

/*
	err := ioutil.WriteGraphToFile("testdata/hello", message, 0644)
	if err != nil {
		log.Fatal(err)
	}*/
/*
	d, ok := s["happ"]
		if ok {
			l := lib.GetTaskList(d, "ghdl -a [UUT]")
			fmt.Println(l)
		}
		//t := lib.GetTaskListAll(s)
*/

func List(f string, args []string) {
	l, r := InduceSubGraphsFromFile(f)
	if len(l) == 0 || len(r) == 0 {
		log.Fatal("Something went wrong. Empty subgraph map!")
	}
	if len(args) != 0 {
		for _, a := range args {
			s, n := GetSubGraph(l, r, a)
			if s == nil {
				log.Fatal("Something went wrong. Empty subgraph!")
			}
			fmt.Printf("[%s]\n", n)
			for _, i := range GetTaskList(s) {
				fmt.Println("  ", i)
			}
		}
		return
	}
	for a, d := range l {
		fmt.Printf("[%s]\n", a)
		for _, i := range GetTaskList(d) {
			fmt.Println("  ", i)
		}
		/* Print list
		if e := WriteGraphToFile(path.Join(o, a+".dot"), d); e != nil {
			log.Fatal(e)
		}
		*/
	}
}

/*
func ListFromFile(f string, args []string) {
	l, r := InduceSubGraphsFromFile(f)
}
*/
