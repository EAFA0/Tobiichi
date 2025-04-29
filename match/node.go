// 一些常用节点的定义

// package tree 定义了树结构中常用节点的类型和方法
package tree

import (
	"cmp"
	"sort"
)

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

func (n DataNode[Q, D]) Build(list []D) Node[Q, D] {
	return DataNode[Q, D]{
		Data: list,
	}
}

func (n DataNode[Q, D]) SetNext(next BuildNode[Q, D]) {
	// 数据节点无下一节点
}

// ListItem 定义了列表项的结构，包含一个可排序的键和一个任意类型的值。
type ListItem[K cmp.Ordered, V any] struct {
	Key K
	Val V
}

func findLeftBound[K cmp.Ordered, V any](l []ListItem[K, V], key K) int {
	// 使用 sort.Search 查找第一个 Key >= key 的元素的索引
	return sort.Search(len(l), func(i int) bool {
		return l[i].Key >= key
	})
}

func findRightBound[K cmp.Ordered, V any](l []ListItem[K, V], key K) int {
	// 使用 sort.Search 查找第一个 Key > key 的元素的索引
	idx := sort.Search(len(l), func(i int) bool {
		return l[i].Key > key
	})
	// 返回前一个元素的索引，即最后一个 Key <= key 的元素的索引
	return idx - 1
}

// PickOneWrapper 将单值查找函数包装为多值返回函数。
// 参数:
//
//	obj - 原始单值查找函数
//
// 返回包装后的多值查找函数。
func PickOneWrapper[K cmp.Ordered, V any](obj func(K, []ListItem[K, V]) (res V, ok bool)) func(K, []ListItem[K, V]) []V {
	return func(k K, li []ListItem[K, V]) []V {
		if val, ok := obj(k, li); ok {
			return []V{val}
		}

		return nil
	}
}

// GE 二分查找获取第一个大于等于key的元素。
// 参数:
//
//	key - 查找的键值
//	l   - 已排序的ListItem切片
//
// 返回值:
//
//	res - 找到的元素值
//	ok  - 是否找到有效元素
func GE[K cmp.Ordered, V any](key K, l []ListItem[K, V]) (res V, ok bool) {
	idx := findLeftBound(l, key)
	if idx < len(l) {
		return l[idx].Val, true
	}

	return
}

// LE 二分查找获取最后一个小于等于key的元素。
// 参数说明同GE函数。
func LE[K cmp.Ordered, V any](key K, l []ListItem[K, V]) (res V, ok bool) {
	idx := findRightBound(l, key)
	if idx >= 0 {
		return l[idx].Val, true
	}

	return
}

// GEs 二分查找获取所有大于等于key的元素。
// 参数和返回值同GE函数。
func GEs[K cmp.Ordered, V any](key K, l []ListItem[K, V]) (res []V) {
	idx := findLeftBound(l, key)
	for _, item := range l[idx:] {
		res = append(res, item.Val)
	}

	return
}

// LEs 二分查找获取所有小于等于key的元素。
// 参数和返回值同GE函数。
func LEs[K cmp.Ordered, V any](key K, l []ListItem[K, V]) (res []V) {
	idx := findRightBound(l, key)
	if idx >= 0 {
		for _, item := range l[:idx+1] {
			res = append(res, item.Val)
		}
	}

	return
}

// InRange 范围查找获取[left, right)区间内的所有元素。
// 参数:
//
//	left  - 区间左边界(包含)
//	right - 区间右边界(不包含)
//	l     - 已排序的ListItem切片
func InRange[K cmp.Ordered, V any](left, right K, l []ListItem[K, V]) (res []V) {
	leftIdx := findLeftBound(l, left)
	rightIdx := findLeftBound(l, right)
	for _, item := range l[leftIdx:rightIdx] {
		res = append(res, item.Val)
	}

	return
}

// groupBy 数据分组工具函数。
// 参数:
//
//	l   - 原始数据列表
//	key - 分组键提取函数
//
// 返回按分组键组织的map。
func groupBy[K comparable, V any](l []V, key func(V) K) (res map[K][]V) {
	res = make(map[K][]V)
	for _, item := range l {
		res[key(item)] = append(res[key(item)], item)
	}

	return
}

// BuildSortedNodeData 构建有序节点数据。
// 流程:
// 1. 按分组键对数据进行分组
// 2. 对分组键进行排序
// 3. 为每个分组构建排序列表项。
func BuildSortedNodeData[K cmp.Ordered, Q, D any](list []D, key func(D) K, next BuildNode[Q, D]) (data []ListItem[K, Node[Q, D]]) {
	group := groupBy(list, key)
	keys := make([]K, 0, len(group))
	for k := range group {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, k := range keys {
		data = append(data, ListItem[K, Node[Q, D]]{
			Key: k, Val: next(group[k]),
		})
	}

	return
}

// BuildMapNodeData 构建哈希映射节点数据。
// 流程:
// 1. 按分组键对数据进行分组
// 2. 为每个分组构建对应的子节点。
func BuildMapNodeData[K comparable, Q, D any](list []D, key func(D) K, next BuildNode[Q, D]) (data map[K]Node[Q, D]) {
	data = make(map[K]Node[Q, D], len(list))
	for k, v := range groupBy(list, key) {
		data[k] = next(v)
	}

	return
}

// SortedNode 排序节点，用于构建有序查询结构。
// K 排序键类型，Q 查询条件类型，D 数据类型。
type SortedNode[K cmp.Ordered, Q, D any] struct {
	NodeBase[Q, D]

	Pick func(Q, []ListItem[K, Node[Q, D]]) []Node[Q, D]
	Data []ListItem[K, Node[Q, D]]

	GroupKey func(D) K
}

func (n SortedNode[K, Q, D]) Next(query Q) []Node[Q, D] {
	return n.Pick(query, n.Data)
}

func (n SortedNode[K, Q, D]) Build(list []D) Node[Q, D] {
	return SortedNode[K, Q, D]{
		Data: BuildSortedNodeData(list, n.GroupKey, n.NextBuilder()),

		NodeBase: n.NodeBase,
		GroupKey: n.GroupKey,
		Pick:     n.Pick,
	}
}

// MapNode 哈希映射节点，用于快速键值查找。
// K 可比较键类型，Q 查询条件类型，D 数据类型。
type MapNode[K comparable, Q, D any] struct {
	NodeBase[Q, D]

	Data map[K]Node[Q, D]
	Pick func(Q, map[K]Node[Q, D]) []Node[Q, D]

	GroupKey func(D) K
}

func (n MapNode[K, Q, D]) Next(query Q) []Node[Q, D] {
	return n.Pick(query, n.Data)
}

func (n MapNode[K, Q, D]) Build(list []D) Node[Q, D] {
	return MapNode[K, Q, D]{
		Data: BuildMapNodeData(list, n.GroupKey, n.NextBuilder()),

		NodeBase: n.NodeBase,
		GroupKey: n.GroupKey,
		Pick:     n.Pick,
	}
}

// UniqueMapNode 创建一个唯一映射节点，用于处理具有唯一键的映射结构。
// 它接受两个函数作为参数：
// - queryKey: 一个函数，用于从查询对象 Q 中提取键 K。
// - groupKey: 一个函数，用于从数据对象 D 中提取键 K。
// 返回一个 MapNode 实例，该实例使用 PickOneInMap 函数来选择节点。
// UniqueMapNode 创建唯一映射节点。
// 参数:
//
//	queryKey  - 从查询条件中提取键的函数
//	groupKey  - 从数据中提取分组键的函数
//
// 返回配置好的MapNode实例。
func UniqueMapNode[K comparable, Q, D any](queryKey func(Q) K, groupKey func(D) K) *MapNode[K, Q, D] {
	return &MapNode[K, Q, D]{
		GroupKey: groupKey,
		Pick: func(q Q, m map[K]Node[Q, D]) []Node[Q, D] {
			if val, ok := m[queryKey(q)]; ok {
				return []Node[Q, D]{val}
			}

			return nil
		},
	}
}

// UniqueSortedNode 创建一个唯一排序节点，用于处理具有唯一键的排序结构。
// 它接受三个函数作为参数：
// - queryKey: 一个函数，用于从查询对象 Q 中提取键 K。
// - queryFunc: 一个函数，用于根据键 K 从排序列表中查找节点，并返回找到的节点和是否找到的布尔值。
// - groupKey: 一个函数，用于从数据对象 D 中提取键 K。
// 返回一个 SortedNode 实例，该实例使用 PickOneInArray 函数来选择节点。
// UniqueSortedNode 创建唯一排序节点。
// 参数:
//
//	queryKey   - 从查询条件中提取键的函数
//	queryFunc  - 在排序列表中查找节点的函数
//	groupKey   - 从数据中提取分组键的函数
//
// 返回配置好的SortedNode实例。
func UniqueSortedNode[K cmp.Ordered, Q, D any](queryKey func(Q) K, queryFunc func(K, []ListItem[K, Node[Q, D]]) (Node[Q, D], bool), groupKey func(D) K) *SortedNode[K, Q, D] {
	return &SortedNode[K, Q, D]{
		GroupKey: groupKey,
		Pick: func(q Q, li []ListItem[K, Node[Q, D]]) []Node[Q, D] {
			return PickOneWrapper(queryFunc)(queryKey(q), li)
		},
	}
}
