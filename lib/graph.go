package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

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

func WriteGraphToFile(f string, g *dep.DependencyGraph) error {
	if b := dot.Marshal(g); b != nil {
		fmt.Printf("Writing graph to '%s'\n", f)
		return ioutil.WriteFile(f, b, 0644)
	}
	return fmt.Errorf("Marshal failed for file %s", f)
}

func ReadGraphFromFile(f string) (*dep.DependencyGraph, error) {
	b, err := ReadFile(f)
	if err != nil {
		return nil, err
	}
	g := dot.Unmarshal(b)
	if g == nil {
		return nil, fmt.Errorf("failed to parse input DOT file")
	}
	return dep.NewDependencyGraph(g), nil
}

func PrintGraph(g *dep.DependencyGraph) {
	for _, n := range graph.NodesOf(g.Nodes()) {
		fmt.Printf("%+v\n", n)
	}
	for _, e := range graph.EdgesOf(g.Edges()) {
		fmt.Printf("%+v\n", e)
	}
}

func InduceSubGraphsFromFile(f string) (map[string]*dep.DependencyGraph, map[string]*dep.DependencyGraph) {
	if len(f) == 0 {
		_, err := os.Stat("graph.dot")
		if err == nil {
			f = "graph.dot"
		}
	}
	d, err := ReadGraphFromFile(f)
	checkErr(err)
	return InduceSubGraphs(d)
}

func InduceSubGraphs(d *dep.DependencyGraph) (map[string]*dep.DependencyGraph, map[string]*dep.DependencyGraph) {
	induce := func(d *dep.DependencyGraph, m map[int64]graph.Node) map[string]*dep.DependencyGraph {
		o := make(map[string]*dep.DependencyGraph)
		for k, n := range d.Induce(m) {
			x := dot.DotGraph{n.DirectedGraph}.Node(k).(*dot.DotNode).DOTID()
			o[x] = n
		}
		return o
	}
	return induce(d, d.Leafs()), induce(d, d.Roots())
}

func Induce(f, o string, args []string) {
	l, r := InduceSubGraphsFromFile(f)
	if len(l) == 0 || len(r) == 0 {
		log.Fatal("Something went wrong. Empty subgraph map!")
	}
	if o != "" {
		checkErr(os.MkdirAll(o, 0644))
	}
	if len(args) != 0 {
		for _, a := range args {
			s, n := GetSubGraph(l, r, a)
			if s == nil {
				log.Fatal("Something went wrong. Empty subgraph!")
			}
			if e := WriteGraphToFile(path.Join(o, n+".dot"), s); e != nil {
				log.Fatal(e)
			}
		}
		return
	}
	for a, d := range l {
		if e := WriteGraphToFile(path.Join(o, a+".dot"), d); e != nil {
			log.Fatal(e)
		}
	}
}

/*
parsearg parses arguments with format 'LEAF[|TASK]'. The optional parameter 'TASK' allows
to filter the list to include only a subset of the tasks in the subgraphs corresponding
to the 'LEAF'. These can be either of:
- '>DOTID' tasks that allow build DOTID.
- 'DOTID>' tasks that depend on DOTID.
- '>DOTID>' tasks that allow to build DOTID and those that depend on it.

it return a key for LEAF, a key for TASK, rv and fw
*/
func parsearg(a string) (string, string, bool, bool) {
	s := strings.Split(a, "|")
	fw, rv, t, l := false, false, "", s[0]
	f := func(t string) (string, bool, bool) {
		if t[0] == '>' {
			t, rv = t[1:], true
		}
		if t[len(t)-1] == '>' {
			t, fw = t[:len(t)-1], true
		}
		if !(fw || rv) {
			rv = true
		}
		return t, rv, fw
	}
	if (len(s) > 1) && (len(s[1]) > 0) {
		t = s[1]
		t, rv, fw = f(t)
		return l, t, rv, fw
	}
	l, rv, fw = f(l)
	return l, "", rv, fw
}

func GetSubGraph(l, r map[string]*dep.DependencyGraph, a string) (*dep.DependencyGraph, string) {
	k, t, rv, fw := parsearg(a)

	if len(t) == 0 {
		if rv && !fw {
			d, ok := l[k]
			if !ok {
				// TODO Check if it is a mid node
				fmt.Printf("subgraph rv,!fw for node '%s' not found\n", k)
				return nil, ""
			}
			return d, k + ".rv"
		}
		if !rv && fw {
			d, ok := r[k]
			if !ok {
				// TODO Check if it is a mid node
				fmt.Printf("subgraph !rv,fw for node '%s' not found\n", k)
				return nil, ""
			}
			return d, k + ".fw"
		}
		if rv && fw {
			fmt.Printf("mode rv,fw not implementing yet. skipping subgraph for node '%s'\n", k)
			return nil, ""
		}
	}

	d, ok := l[k]
	if !ok {
		// TODO Check if it is a mid node
		fmt.Printf("subgraph rv,!fw for node '%s' not found\n", k)
		return nil, ""
	}
	n := make(map[int64]graph.Node)
	x := dot.DotGraph{d.DirectedGraph}.GetNodeByDOTID(t)
	if x == nil {
		fmt.Printf("node '%s' not found in subgraph for '%s'!\n", t, k)
		return nil, ""
	}
	n[x.ID()] = x
	s := d.InduceDir(n, fw, rv)
	return s[x.ID()], k + "." + t
}
