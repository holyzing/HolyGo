package GoArithmetic

type T interface{}

type Edge struct {
	from   Node
	to     Node
	weight float64
}

type Node struct {
	data    T
	index   int
	visited bool
	edges   []Edge
}

func (node *Node) addEdge(edge Edge) {
	node.edges = append(node.edges, edge)
}
