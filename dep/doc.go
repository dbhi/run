/*
Package dep provides dependency graph analysis and manipulation logic.

Following the design style in `gonum/graph/flow`, type `graph.Directed` is wrapped in a new
struct named `DependencyGraph` which includes three unexported fields:
  roots    map of root nodes where the key is the ID
  leafs    map of leaf nodes where the key is the ID
  schedule naive topological sort of the graph as a slice of IDs

Basic features, such as retrieving a map of roots/leafs, inducing a subgraph for a leaf or retrieving a
valid schedule (topological sort) are already implemented. However, it'd be interesting to extend it with
other common operations: `InduceAllIn`, `Reduce`, `Closure` and `Schedule` are on the roadmap.

A topological ordering may be constructed by reversing a postorder numbering of a depth-first search
graph traversal. Hence, it should be possible to construct a solution to the topological order at the same
time that `Induction`, `Reduction` or `Closure` are computed, as long as these are implemented through
`traverse.DepthFirst`. This is explained in https://en.wikipedia.org/wiki/Topological_sorting#Depth-first_search.
Certainly, since Tarjan is mentioned as the first one describing this approach, this algorithm may already
be implemented as `gonum/graph/topo.Sort`. Should ask @kortschak about it.

Overall, the exact solution to the topological sorting provided implicitly by `Induction`, `Reduction` or
Closure is not relevant, as long as it is valid. Finer control of the sorting approaches is provided
through `Schedule`:

  Is it worth implementing Kahn's algorithm for `p==1`?
    - Is it faster than Tarjan?
    - Is it better than the 'naive' implicit solution from `traverse.DepthFirst`?

  Is it worth executing `topo.Sort` for `p==1`?
    - Is the resulting order better that the 'naive' implicit solution from `traverse.DepthFirst`?
    - `topo.Sort` is a general interface for many directed/undirected and (a)cyclic graphs.
      It might be too complex for this specific task.

  How to handle `p>2`? Should we implement Coffman–Graham algorithm? Or some other?

  How to handle `p>3` ? Is Coffman–Graham still suitable? Is there any other? Is there any constraint,
  such as 'p' being required to be a power of two?

ATM, we assume that the graph provided by the user is indeed a DAG, which might or might not be true. We
should provide a `IsDAG` function to check it. Can we use `topo.DirectedCyclesIn(g) == nil`?

For now, all the scheduling algorithms assume a non-weighted graph. In the future, it is possible to take
edge weight into account, as suggested in http://dnaeon.github.io/graphs-and-clojure/.

References:
  - Dependency graph: https://en.wikipedia.org/wiki/Dependency_graph
  - Directed acyclic graph: https://en.wikipedia.org/wiki/Directed_acyclic_graph
  - Induced subgraph: https://en.wikipedia.org/wiki/Induced_subgraph
  - Transitive reduction: https://en.wikipedia.org/wiki/Transitive_reduction
  - Transitive closure: https://en.wikipedia.org/wiki/Transitive_closure
  - Topological sorting: https://en.wikipedia.org/wiki/Topological_sorting
    - Kahn's algoritm: https://en.wikipedia.org/wiki/Topological_sorting#Kahn's_algorithm
    - Coffman–Graham algorithm: https://en.wikipedia.org/wiki/Coffman%E2%80%93Graham_algorithm
*/
package dep
