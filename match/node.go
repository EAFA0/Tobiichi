// 一些常用节点的定义

// package tree 定义了树结构中常用节点的类型和方法
package tree

import (
	"cmp"
	"sort"
)

// ListItem defines a structure for list items, containing a sortable key and a value of any type.
// K is the type of the key, which must be ordered, and V is the type of the value.
type ListItem[K cmp.Ordered, V any] struct {
	Key K
	Val V
}

// findLeftBound performs a binary search to find the index of the first element in the list
// where the key is greater than or equal to the specified key.
func findLeftBound[K cmp.Ordered, V any](l []ListItem[K, V], key K) int {
	return sort.Search(len(l), func(i int) bool {
		return l[i].Key >= key
	})
}

// findRightBound performs a binary search to find the index of the last element in the list
// where the key is less than or equal to the specified key.
func findRightBound[K cmp.Ordered, V any](l []ListItem[K, V], key K) int {
	idx := sort.Search(len(l), func(i int) bool {
		return l[i].Key > key
	})
	return idx - 1
}

// PickOneWrapper wraps a single-value lookup function into a multi-value return function.
// Parameters:
// - obj: The original single-value lookup function.
// Returns a function that returns a slice of values.
func PickOneWrapper[K cmp.Ordered, V any](obj func(K, []ListItem[K, V]) (res V, ok bool)) func(K, []ListItem[K, V]) []V {
	return func(k K, li []ListItem[K, V]) []V {
		if val, ok := obj(k, li); ok {
			return []V{val}
		}

		return nil
	}
}

// GE performs a binary search to find the first element in the list
// where the key is greater than or equal to the specified key.
func GE[K cmp.Ordered, V any](key K, l []ListItem[K, V]) (res V, ok bool) {
	idx := findLeftBound(l, key)
	if idx < len(l) {
		return l[idx].Val, true
	}

	return
}

// LE performs a binary search to find the last element in the list
// where the key is less than or equal to the specified key.
func LE[K cmp.Ordered, V any](key K, l []ListItem[K, V]) (res V, ok bool) {
	idx := findRightBound(l, key)
	if idx >= 0 {
		return l[idx].Val, true
	}

	return
}

// GEs retrieves all elements in the list where the key is greater than or equal to the specified key.
func GEs[K cmp.Ordered, V any](key K, l []ListItem[K, V]) (res []V) {
	idx := findLeftBound(l, key)
	for _, item := range l[idx:] {
		res = append(res, item.Val)
	}

	return
}

// LEs retrieves all elements in the list where the key is less than or equal to the specified key.
func LEs[K cmp.Ordered, V any](key K, l []ListItem[K, V]) (res []V) {
	idx := findRightBound(l, key)
	if idx >= 0 {
		for _, item := range l[:idx+1] {
			res = append(res, item.Val)
		}
	}

	return
}

// InRange retrieves all elements in the list where the key is within the range [left, right).
func InRange[K cmp.Ordered, V any](left, right K, l []ListItem[K, V]) (res []V) {
	leftIdx := findLeftBound(l, left)
	rightIdx := findLeftBound(l, right)
	for _, item := range l[leftIdx:rightIdx] {
		res = append(res, item.Val)
	}

	return
}

// groupBy groups a list of items by a specified key extraction function.
// Parameters:
// - l: The original list of items.
// - key: A function to extract the grouping key from each item.
// Returns a map where the keys are the extracted keys and the values are slices of items.
func groupBy[K comparable, V any](l []V, key func(V) K) (res map[K][]V) {
	res = make(map[K][]V)
	for _, item := range l {
		res[key(item)] = append(res[key(item)], item)
	}

	return
}

// SortedNodeBuilder is a builder for creating SortedNode instances.
// It organizes data into a sorted structure and builds child nodes.
type SortedNodeBuilder[K cmp.Ordered, Q, D any] struct {
	Pick     func(Q, []ListItem[K, Node[Q, D]]) []Node[Q, D]
	GroupKey func(D) K
	data     []ListItem[K, []D]
	next     []NodeBuilder[Q, D]
}

func (s SortedNodeBuilder[K, Q, D]) Load(list []D) NodeBuilder[Q, D] {
	group := groupBy(list, s.GroupKey)
	keys := make([]K, 0, len(group))
	for k := range group {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	data := make([]ListItem[K, []D], 0, len(group))
	for _, k := range keys {
		data = append(data, ListItem[K, []D]{
			Key: k, Val: group[k],
		})
	}
	return SortedNodeBuilder[K, Q, D]{
		next:     make([]NodeBuilder[Q, D], len(data)),
		data:     data,
		Pick:     s.Pick,
		GroupKey: s.GroupKey,
	}
}

// Push adds child node builders to the current builder and returns the updated list of builders.
func (s SortedNodeBuilder[K, Q, D]) Push(next NodeBuilder[Q, D]) []NodeBuilder[Q, D] {
	for index, elem := range s.data {
		s.next[index] = next.Load(elem.Val)
	}
	return s.next
}

// Node finalizes and returns the constructed SortedNode.
func (s SortedNodeBuilder[K, Q, D]) Node() Node[Q, D] {
	data := make([]ListItem[K, Node[Q, D]], len(s.data))
	for i, elem := range s.data {
		data[i] = ListItem[K, Node[Q, D]]{
			Key: elem.Key,
			Val: s.next[i].Node(),
		}
	}

	return SortedNode[K, Q, D]{
		Pick: s.Pick,
		Data: data,
	}
}

type SortedNode[K cmp.Ordered, Q, D any] struct {
	EmptyNode[Q, D]
	Data []ListItem[K, Node[Q, D]]
	Pick func(Q, []ListItem[K, Node[Q, D]]) []Node[Q, D]
}

// Next retrieves the next set of nodes based on the query for a SortedNode.
func (n SortedNode[K, Q, D]) Next(query Q) []Node[Q, D] {
	return n.Pick(query, n.Data)
}

// MapNodeBuilder is a builder for creating MapNode instances.
// It organizes data into a map structure and builds child nodes.
type MapNodeBuilder[K comparable, Q, D any] struct {
	Pick     func(Q, map[K]Node[Q, D]) []Node[Q, D]
	GroupKey func(D) K
	data     map[K][]D
	next     map[K]NodeBuilder[Q, D]
}

func (m MapNodeBuilder[K, Q, D]) Load(list []D) NodeBuilder[Q, D] {
	group := groupBy(list, m.GroupKey)
	return MapNodeBuilder[K, Q, D]{
		next:     make(map[K]NodeBuilder[Q, D], len(group)),
		data:     group,
		Pick:     m.Pick,
		GroupKey: m.GroupKey,
	}
}

// Push adds child node builders to the current builder and returns the updated list of builders.
func (m MapNodeBuilder[K, Q, D]) Push(next NodeBuilder[Q, D]) []NodeBuilder[Q, D] {
	for k, elem := range m.data {
		m.next[k] = next.Load(elem)
	}
	nextBuilder := make([]NodeBuilder[Q, D], 0, len(m.data))
	for _, elem := range m.next {
		nextBuilder = append(nextBuilder, elem)
	}
	return nextBuilder
}

// Node finalizes and returns the constructed MapNode.
func (m MapNodeBuilder[K, Q, D]) Node() Node[Q, D] {
	data := make(map[K]Node[Q, D], len(m.data))
	for k := range m.data {
		data[k] = m.next[k].Node()
	}
	return MapNode[K, Q, D]{
		Data: data,
		Pick: m.Pick,
	}
}

type MapNode[K comparable, Q, D any] struct {
	EmptyNode[Q, D]
	Data map[K]Node[Q, D]
	Pick func(Q, map[K]Node[Q, D]) []Node[Q, D]
}

// Next retrieves the next set of nodes based on the query for a MapNode.
func (n MapNode[K, Q, D]) Next(query Q) []Node[Q, D] {
	return n.Pick(query, n.Data)
}

// UniqueMapNode creates a MapNodeBuilder for unique key mappings.
// It uses queryKey to extract keys from queries and groupKey to extract keys from data.
func UniqueMapNode[K comparable, Q, D any](queryKey func(Q) K, groupKey func(D) K) MapNodeBuilder[K, Q, D] {
	return MapNodeBuilder[K, Q, D]{
		GroupKey: groupKey,
		Pick: func(q Q, m map[K]Node[Q, D]) []Node[Q, D] {
			if val, ok := m[queryKey(q)]; ok {
				return []Node[Q, D]{val}
			}

			return nil
		},
	}
}

// UniqueSortedNode creates a SortedNodeBuilder for unique key mappings in a sorted structure.
// It uses queryKey to extract keys from queries, queryFunc to find nodes in the sorted list, and groupKey to extract keys from data.
func UniqueSortedNode[K cmp.Ordered, Q, D any](queryKey func(Q) K, queryFunc func(K, []ListItem[K, Node[Q, D]]) (Node[Q, D], bool), groupKey func(D) K) SortedNodeBuilder[K, Q, D] {
	return SortedNodeBuilder[K, Q, D]{
		GroupKey: groupKey,
		Pick: func(q Q, li []ListItem[K, Node[Q, D]]) []Node[Q, D] {
			return PickOneWrapper(queryFunc)(queryKey(q), li)
		},
	}
}
