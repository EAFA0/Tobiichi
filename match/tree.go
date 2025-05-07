// 使用字典树构建的匹配算法
package tree

// Node defines the interface for a trie node used in conditional matching algorithms.
// Q represents the query type, and D represents the data storage type.
type Node[Q, D any] interface {
	// Next retrieves the next set of nodes based on the query.
	Next(query Q) []Node[Q, D]
	// Leaf checks if the node is a leaf and retrieves its data.
	Leaf() ([]D, bool)
}

// EmptyNode represents a terminal node in the trie with no data.
// It serves as a marker for the end of a matching process.
type EmptyNode[Q, D any] struct{}

func (n EmptyNode[Q, D]) Leaf() ([]D, bool) {
	return nil, false // Empty nodes have no data.
}

// NodeBuilder defines the interface for constructing trie nodes.
// Implementations manage the process of building child nodes and distributing data.
type NodeBuilder[Q, D any] interface {
	// Load initializes the builder with data.
	Load(data []D) NodeBuilder[Q, D]
	// Push adds a child node builder and returns the updated list of builders.
	Push(next NodeBuilder[Q, D]) []NodeBuilder[Q, D]
	// Node finalizes and returns the constructed node.
	Node() Node[Q, D]
}

// DataNode represents a terminal node that stores data.
// Q is the query type, and D is the data type.
type DataNode[Q, D any] struct {
	Data []D
}

func (n DataNode[Q, D]) Next(query Q) []Node[Q, D] {
	return nil // Data nodes have no child nodes.
}

func (n DataNode[Q, D]) Leaf() ([]D, bool) {
	return n.Data, len(n.Data) != 0
}

// DataNodeBuilder is a builder for creating DataNode instances.
type DataNodeBuilder[Q, D any] struct {
	Data []D
}

func (n DataNodeBuilder[Q, D]) Load(list []D) NodeBuilder[Q, D] {
	return DataNodeBuilder[Q, D]{Data: list}
}

func (n DataNodeBuilder[Q, D]) Push(next NodeBuilder[Q, D]) []NodeBuilder[Q, D] {
	return nil
}

func (n DataNodeBuilder[Q, D]) Node() Node[Q, D] {
	return DataNode[Q, D](n)
}

// Build constructs a chain of nodes based on the provided builders.
// The builders are linked in reverse order, starting from the last one.
//
// Parameters:
// - list: Initial dataset.
// - builder: List of node builders in reverse order.
//
// Returns the root node, which serves as the entry point for the matching process.
func Build[Q, D any](list []D, builder []NodeBuilder[Q, D]) Node[Q, D] {
	// Append a DataNodeBuilder to the builder list.
	copyBuilder := make([]NodeBuilder[Q, D], len(builder))
	copy(copyBuilder, builder)
	builder = append(copyBuilder, DataNodeBuilder[Q, D]{})

	// Initialize the root builder with the dataset.
	rootBuilder := builder[0].Load(list)
	// Build the node tree layer by layer.
	current := []NodeBuilder[Q, D]{rootBuilder}
	for _, nextBuilder := range builder[1:] {
		next := make([]NodeBuilder[Q, D], 0, len(current))
		for _, item := range current {
			addition := item.Push(nextBuilder)
			next = append(next, addition...)
		}
		current = next
	}
	// Return the root node.
	return rootBuilder.Node()
}

// Search performs a breadth-first search (BFS) to traverse the trie and match nodes.
//
// Parameters:
// - root: Starting node for the search.
// - query: Query object used for matching.
//
// Returns a collection of data stored in the matching nodes.
func Search[Query, Data any](root Node[Query, Data], query Query) (res []Data) {
	queue := []Node[Query, Data]{root}
	for ; len(queue) > 0; queue = queue[1:] {
		if val, ok := queue[0].Leaf(); ok {
			res = append(res, val...)
		}
		queue = append(queue, queue[0].Next(query)...)
	}

	return
}
