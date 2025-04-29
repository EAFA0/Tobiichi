// 使用字典树构建的匹配算法
package tree

// Node 定义字典树节点接口，用于实现条件匹配算法。
// Q: 查询条件类型 D: 数据存储类型。
type Node[Q, D any] interface {
	Next(query Q) []Node[Q, D] // 继续查找下一条件
	Leaf() ([]D, bool)
}

// NodeBuilder 节点构建器接口，负责组织节点构建逻辑。
// 实现类需要管理子节点构建过程和数据分发规则。
type NodeBuilder[Q, D any] interface {
	Build(data []D) Node[Q, D]
	SetNext(next BuildNode[Q, D])
}

// BuildNode 节点构建函数类型，用于链式构建节点层级。
// 接收数据集合，返回构建完成的节点实例。
type BuildNode[Q, D any] func(data []D) Node[Q, D]

// EmptyNode 空节点实现，作为叶子节点的终止标识。
// 不包含实际数据，用于终止匹配流程。
type EmptyNode[Q, D any] struct{}

func (n EmptyNode[Q, D]) Leaf() ([]D, bool) {
	return nil, false // 空节点无数据
}

// NodeLinker 节点链接器，维护节点间的链式关系。
// 通过next字段连接下一个条件匹配节点。
type NodeLinker[Q, D any] struct {
	next BuildNode[Q, D] // 下一个条件
}

func (b *NodeLinker[Q, D]) SetNext(next BuildNode[Q, D]) {
	b.next = next
}

func (b *NodeLinker[Q, D]) NextBuilder() BuildNode[Q, D] {
	return b.next
}

// Build 构造节点链，根据构建器顺序反向链接节点。
// 参数:
//
//	list - 初始数据集合
//	order - 节点构建器列表(构建顺序从后往前链接)
//
// 返回根节点，启动匹配流程的入口。
func Build[Q, D any](list []D, order []NodeBuilder[Q, D]) Node[Q, D] {
	for index := len(order) - 1; index > 0; index-- {
		order[index-1].SetNext(order[index].Build)
	}

	return order[0].Build(list)
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

// NodeBase 基础节点结构体，组合常用节点功能。
// 包含节点链接能力和空节点终止逻辑。
type NodeBase[Q, D any] struct {
	NodeLinker[Q, D]
	EmptyNode[Q, D]
}
