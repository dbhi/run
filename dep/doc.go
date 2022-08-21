/*
Package dep provides dependency graph analysis and manipulation logic.

Following the design style in `gonum/graph/flow`, type `graph.Directed` is wrapped in a new
struct named `DependencyGraph` which includes two unexported fields:

- roots: map of root nodes where the key is the ID
- leafs: map of leaf nodes where the key is the ID

Basic features, such as retrieving a map of roots/leafs, inducing a subgraph for a leaf or retrieving a
valid schedule (topological sort) are already implemented. However, it'd be interesting to extend it with
other common operations; `InduceAllIn`, `Reduce`, `Closure` and `Schedule` are on the roadmap.

References:
  - Dependency graph: https://en.wikipedia.org/wiki/Dependency_graph
  - Directed acyclic graph: https://en.wikipedia.org/wiki/Directed_acyclic_graph
  - Induced subgraph: https://en.wikipedia.org/wiki/Induced_subgraph
  - Transitive reduction: https://en.wikipedia.org/wiki/Transitive_reduction
  - Transitive closure: https://en.wikipedia.org/wiki/Transitive_closure
  - Topological sorting: https://en.wikipedia.org/wiki/Topological_sorting
    -- Kahn's algoritm: https://en.wikipedia.org/wiki/Topological_sorting#Kahn's_algorithm
    -- Coffman–Graham algorithm: https://en.wikipedia.org/wiki/Coffman%E2%80%93Graham_algorithm
*/
package dep
