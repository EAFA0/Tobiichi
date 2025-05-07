package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 添加针对 node.go 的单元测试
func TestSortedNode(t *testing.T) {
	data := []ListItem[int, Node[int, string]]{
		{Key: 1, Val: DataNode[int, string]{Data: []string{"A"}}},
		{Key: 2, Val: DataNode[int, string]{Data: []string{"B"}}},
		{Key: 3, Val: DataNode[int, string]{Data: []string{"C"}}},
	}

	node := SortedNode[int, int, string]{
		Data: data,
		Pick: func(query int, items []ListItem[int, Node[int, string]]) []Node[int, string] {
			for _, item := range items {
				if item.Key == query {
					return []Node[int, string]{item.Val}
				}
			}
			return nil
		},
	}

	result := node.Next(2)
	assert.Equal(t, 1, len(result))
}

func TestMapNode(t *testing.T) {
	data := map[int]Node[int, string]{
		1: DataNode[int, string]{Data: []string{"A"}},
		2: DataNode[int, string]{Data: []string{"B"}},
	}

	node := MapNode[int, int, string]{
		Data: data,
		Pick: func(query int, items map[int]Node[int, string]) []Node[int, string] {
			if val, ok := items[query]; ok {
				return []Node[int, string]{val}
			}
			return nil
		},
	}

	result := node.Next(1)
	assert.Equal(t, 1, len(result))
}

// 添加更多针对 node.go 的单元测试
func TestPickOneWrapper(t *testing.T) {
	data := []ListItem[int, string]{
		{Key: 1, Val: "A"},
		{Key: 2, Val: "B"},
	}
	wrapped := PickOneWrapper(GE[int, string])
	result := wrapped(1, data)
	assert.Equal(t, []string{"A"}, result)

	result = wrapped(3, data)
	assert.Empty(t, result)
}

func TestGE(t *testing.T) {
	data := []ListItem[int, string]{
		{Key: 1, Val: "A"},
		{Key: 2, Val: "B"},
	}
	val, ok := GE(1, data)
	assert.True(t, ok)
	assert.Equal(t, "A", val)

	_, ok = GE(3, data)
	assert.False(t, ok)
}

func TestLE(t *testing.T) {
	data := []ListItem[int, string]{
		{Key: 1, Val: "A"},
		{Key: 2, Val: "B"},
	}
	val, ok := LE(2, data)
	assert.True(t, ok)
	assert.Equal(t, "B", val)

	_, ok = LE(0, data)
	assert.False(t, ok)
}

func TestInRange(t *testing.T) {
	data := []ListItem[int, string]{
		{Key: 1, Val: "A"},
		{Key: 2, Val: "B"},
		{Key: 3, Val: "C"},
	}
	result := InRange(1, 3, data)
	assert.Equal(t, []string{"A", "B"}, result)
}

func TestGEs(t *testing.T) {
	data := []ListItem[int, string]{
		{Key: 1, Val: "A"},
		{Key: 2, Val: "B"},
		{Key: 3, Val: "C"},
	}
	result := GEs(2, data)
	assert.Equal(t, []string{"B", "C"}, result)
}

func TestLEs(t *testing.T) {
	data := []ListItem[int, string]{
		{Key: 1, Val: "A"},
		{Key: 2, Val: "B"},
		{Key: 3, Val: "C"},
	}
	result := LEs(2, data)
	assert.Equal(t, []string{"A", "B"}, result)
}

func TestSortedNodeBuilder(t *testing.T) {
	builder := SortedNodeBuilder[int, int, string]{
		GroupKey: func(d string) int { return len(d) },
	}
	builder = builder.Load([]string{"A", "BB", "CCC"}).(SortedNodeBuilder[int, int, string])
	assert.Equal(t, 3, len(builder.data))
}

func TestMapNodeBuilder(t *testing.T) {
	builder := MapNodeBuilder[int, int, string]{
		GroupKey: func(d string) int { return len(d) },
	}
	builder = builder.Load([]string{"A", "BB", "CCC"}).(MapNodeBuilder[int, int, string])
	assert.Equal(t, 3, len(builder.data))
}
