// 使用字典树构建的匹配算法
package tree

// Node 定义字典树节点接口，用于实现条件匹配算法。
// Q: 查询条件类型 D: 数据存储类型。
type Node[Q, D any] interface {
	Next(query Q) []Node[Q, D] // 继续查找下一条件
	Leaf() ([]D, bool)
}

// EmptyNode 空节点实现，作为叶子节点的终止标识。
// 不包含实际数据，用于终止匹配流程。
type EmptyNode[Q, D any] struct{}

func (n EmptyNode[Q, D]) Leaf() ([]D, bool) {
	return nil, false // 空节点无数据
}

// NodeBuilder 节点构建器接口，负责组织节点构建逻辑。
// 实现类需要管理子节点构建过程和数据分发规则。
type NodeBuilder[Q, D any] interface {
	Load(data []D) NodeBuilder[Q, D]                 // 加载节点数据
	Push(next NodeBuilder[Q, D]) []NodeBuilder[Q, D] // 下推构建节点, 为空时说明构造完成
	Node() Node[Q, D]                                // 返回构建好的节点实例
}

// DataNode 数据存储节点，用于保存最终匹配结果。
// S 查询条件类型，T 数据类型。
type DataNode[Q, D any] struct {
	Data []D
}

func (n DataNode[Q, D]) Next(query Q) []Node[Q, D] {
	return nil // 数据节点无下一节点
}

func (n DataNode[Q, D]) Leaf() ([]D, bool) {
	return n.Data, len(n.Data) != 0
}

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

// Build 构造节点链，根据构建器顺序反向链接节点。
// 参数:
//
//	list - 初始数据集合
//	order - 节点构建器列表(构建顺序从后往前链接)
//
// 返回根节点，启动匹配流程的入口。
func Build[Q, D any](list []D, builder []NodeBuilder[Q, D]) Node[Q, D] {
	// append 一个数据节点
	copyBuilder := make([]NodeBuilder[Q, D], len(builder))
	copy(copyBuilder, builder)
	builder = append(copyBuilder, DataNodeBuilder[Q, D]{})

	// 构建根节点, 只有根节点需要 load 数据
	rootBuilder := builder[0].Load(list)
	// 分层构建节点树
	current := []NodeBuilder[Q, D]{rootBuilder}
	for _, nextBuilder := range builder[1:] {
		next := make([]NodeBuilder[Q, D], 0, len(current))
		for _, item := range current {
			addition := item.Push(nextBuilder)
			next = append(next, addition...)
		}
		current = next
	}
	// 返回顶层节点树
	return rootBuilder.Node()
}

// Search 执行广度优先搜索(BFS)，遍历匹配节点链。
// 参数:
//
//	root - 搜索起始节点
//	query - 查询条件对象
//
// 返回所有匹配节点中存储的数据集合。
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
