<p align="center">
  <img src="./logo.png" width="550"/>
</p>

<p align="center">
  <a title="GoDoc" href="https://godoc.org/github.com/umarcor/run/lib"><img src="https://img.shields.io/badge/godoc-reference-5272B4.svg?longCache=true&style=flat-square&logo=go&logoColor=fff"></a><!--
  -->
  <a title="Releases" href="https://github.com/umarcor/run/releases"><img src="https://img.shields.io/github/commits-since/umarcor/run/latest.svg?longCache=true&style=flat-square"></a>
</p>

---

`run` is a task execution automation package that analyses complex dependency graphs ([multitree](https://en.wikipedia.org/wiki/Multitree) [directed acyclic graph](https://en.wikipedia.org/wiki/Directed_acyclic_graph)), generates filtered subgraphs for each target and provides ordered lists of tasks through [topological sorting](https://en.wikipedia.org/wiki/Topological_sorting). `run` is neither a replacement nor a wrapper for [Make](https://en.wikipedia.org/wiki/Make_(software)), instead it is a complement. The main use case is the combination of multiple build steps, each with a different tool which has it's own build system (be it `make`, `cmake`, `go build`, etc.), and probably involving parameters provided through JSON files, CLI arguments and/or environment variables.

`run/lib` is meant to be used imported to other [golang](https://golang.org/) projects, such as CLI tools or web services. This allows golang developers to make the best of third-party libraries to process data and handle parallel (concurrent) execution seamlessly.

`run/cli` is an example implementation of such a CLI tool which is based on [spf13/cobra](https://github.com/spf13/cobrahttps://github.com/spf13/viper/) and [spf13/viper](https://github.com/spf13/viper/) along with `run/lib`. This is provided as a reference of how to interact with the API of `run/lib`, but it can be used as a standalone tool.

`run` is built on top of [gonum](https://www.gonum.org/). Precisely, types, interfaces and functions defined in [gonum/graph](https://github.com/gonum/gonum/tree/master/graph) ([godoc.org/gonum.org/v1/gonum/graph](https://godoc.org/gonum.org/v1/gonum/graph)) are used to manipulate graphs. Therefore, `run` relies on the list of input formats supported by the package. In the examples, [Graphviz](https://www.graphviz.org/)'s [DOT](https://en.wikipedia.org/wiki/DOT_(graph_description_language)) language is used, which is a widespread output format supported by many tools. For example, [lindenb/makefile2graph](https://github.com/lindenb/makefile2graph) allows to analyse makefiles, and [kisielk/godepgraph](https://github.com/kisielk/godepgraph) generates *a dependency graph of Go packages*.

# Usage

The main input is a large complex graph where developers put the dependencies of their multiple workflows. Some of them are cross-related, some are independent dependency chains. This can be provided as a `graphviz` `dot` file (e.g. [`example/graph.dot`](./example/graph.dot)).

That's enough for basic features, such as reducing the complexity, filtering the nodes/edges, getting topologically ordered lists, etc. In order use execution features, the context of each task/job needs to be defined. This is currently done through either a JSON file (e.g. [`example/config.json`](./example/graph.dot)) or golang sources.

> NOTE: tasks/jobs cannot be defined through golang sources at runtime, unless golang is available. If pre-built binaries are used, new tasks/jobs can only be defined through JSON files.

> NOTE: in the discussion about similar projects below some info is provided about other input formats that we would like to support in the future.

## Induce

``` bash
run induce -g graph.json -o subgraphs leafs
# OR
run induce -c config.json -o subgraphs leafs
# note that '"graph": "graph.dot"' is defined in 'config.json'
```

Generates a DOT subgraph in subdir `subgraphs` for each of the leafs in in `graph.dot`. Each subgraph includes only the dependencies required to build the corresponding leaf.

> WIP:
> - allow to induce the graph of a single leaf.
> - allow to induce the graph of the nodes that depend on a root.
> - allow to induce the graph of a single root.
> - allow to induce the graph of a mid node (either forward, reverse or both).

> HINT: the subgraphs can be shown in a web frontend, so the user can select to visualize all the dependecies or to choose a single target and show the corresponding subgraph.

## List

``` bash
run list -g graph.json NODE
# OR
run list -c config.json NODE
```

Returns an ordered list of tasks/jobs required to execute the given target NODE. The target can be any leaf or mid vertex.

> WIP:
> ``` bash
> run list -c config.json NODE[:FILTER]
> ```
>
> Returns an ordered list of tasks/jobs required to execute the given target NODE. The > target can be any leaf or mid vertex. The optional argument `FILTER` allows to > filter the list to include only a subset of the tasks in the subgraphs corresponding > to the node. It can be either of:
> - `>FNODE` jobs that allow build FNODE.
> - `FNODE>` jobs that depend on FNODE.
> - `>FNODE>` jobs that allow to build FNODE and those that depend on it.
>
> For example:
>
> ``` bash
> # run list -c config.json bin
>
> # run list -c config.json bin:>objB
> # run list -c config.json bin:objA>
> # run list -c config.json bin:>buildB>
> ```

## Exec

``` bash
run exec -g graph.json NODE
# note that the context and logic of the jobs must have been previously defined and built into 'run'
# OR
run exec -c config.json NODE
```

Executes all the tasks until NODE (included), in topological order.

> WIP:
> ``` bash
> run exec -c config.json NODE[:EXCLUDE]
> ```
>

# References

- [gonum](https://www.gonum.org)
  - [godoc.org/gonum.org/v1/gonum](https://godoc.org/gonum.org/v1/gonum)
  - [gonum/gonum#910](https://github.com/gonum/gonum/issues/910)
    - [Preserving labels when marshaling a dot graph](https://groups.google.com/forum/#!topic/gonum-dev/xupu8gEmuIs)
    - [godoc.org/github.com/kortschak/graphprac](https://godoc.org/github.com/kortschak/graphprac)
- [semver.org](https://semver.org)
- [keepachangelog.com](https://keepachangelog.com)
- [dnaeon.github.io/graphs-and-clojure](http://dnaeon.github.io/graphs-and-clojure/)

# Similar projects

Use cases for `run` are similar to those for other tools such as:

- [taskfile.dev](https://taskfile.dev) ([go-task/task](https://github.com/go-task/task)): a task runner/simpler Make alternative written in Go.
  - Users need to write all the configuration details in one or multiple `Taskfile.yml` files. We want to support this. We might accept `Taskfile.yml` files indeed. But this should not be the single source of configuration.
  - Dependencies are described explicitly. A different syntax is used to define dependencies that can be executed concurrently ('dep') and those that need to be executed serially ('task'). A dependency tree is derived, which is a specific type of dependency graph. As a result, some dependencies/tasks are built multiple times, if required by multiple jobs. Instead, we want to process the dependencies as a graph, and be able to reduce and topologically sort the list of tasks.
  - Some golang features are used to handle OS specific tasks. This is something we might want to support.
  - Sources and artifacts are explicitly listed through `sources` and `generates`, respectively. This allows to watch for changes and also to clean the artifacts. We want to support both features. However, we should not constrain the format to specific file paths. Instead we should support either expansion of wildcard or regexps.
  - Golang templating features can be used in the `Taskfile.yml` file. We might want to support this.
  - Wen multiple tasks are executed concurrently, the output can be set to `interleaved`, `group` or `prefixed`. We want to support this feature, although we will probably use a different approach.
  - The tool is distributed as a single static binary. We want to do this too.
- [magefile.org/](https://magefile.org/) ([magefile/mage](https://github.com/magefile/mage)): a Make/rake-like build tool using Go.
  - Users need to write all the configuration details in one or multiple golang sources. This offers great flexibility and is a very powerful approach. We definitely want to support this approach, but this should not be the single source of configuration.
  - Golang is required on the target platform. If pre-built, it is not possible to later define additional tasks ('targets') without golang. We do not want golang to be a required on the target platform in order to add jobs to the configuration which where not defined when the tool was built.
  - Target aliases and Namespaces are supported. This is something we want to do too.
  - `context` is supported. This is something we might want to support.
  - Dependencies are described explicitly and different functions are used (`Deps` or `SerialDeps`). A dependency graph is built so that each dependency is guaranteed to be run exactly one. This is something we want to support too, but it should not be the single input source to the graph. We want to also support filtering some of the tasks at runtime.
  - Two functions are provided to watch files. `target.Path` watches a directory or file not recursively and `target.Dir` watches a directory recursively.
  - It is available as a compile-in library and three helper libraries are provided (`mg`, `sh` and `target`). At some point, we might find it interesting to directly integrate `run` and `mage`. Further analysis is required to do so.
    - `run` can import some features from the `mage` API, transparently to the user.
    - `run` can behave as a frontend/extension to `mage` which:
      - Allows to read add 'target' configurations from other sources at runtime (i.e. DOT graphs or `Taskfile.yml` files).
      - Provide generic functions to execute the targets defined through other sources.

Therefore, it can be said that one of the purposes of `run` is to somehow integrate them.

# ToDo for v1.0.0

- We want to be able to retrieve the topological order of a subset of nodes in one of the subgraphs. We are evaluating three approaches:
  - Retrieve the topological order of the subgraph and then remove the items that are to be ignored.
  - Generate a subsubgraph from the subgraph and retrieve the topological order.
  - We feel that we will need both: first generate a subsubgraph and then optionally remove some items from the topological order.
- Propose `gonum/graph/dep`.
- Support minimal web GUI to show subgraphs, subsubgraphs and task lists.
- Provide basic example implementation of 'Exec'.
- Allow to decide whether a target needs to be regenerated by comparing file modification times.
- Merge graphs from different sources which might share some nodes and edges.
- Dry run mode
- Go's template engine
- ignore certain tasks in a list