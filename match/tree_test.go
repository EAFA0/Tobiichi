package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 添加针对 tree.go 的单元测试
func TestBuildAndSearch(t *testing.T) {
	data := []map[string]int{
		{"Field1": 1, "Field2": 2},
		{"Field1": 3, "Field2": 4},
	}
	query := map[string]int{"Field1": 1}

	order := []NodeBuilder[map[string]int, map[string]int]{
		UniqueMapNode(
			func(q map[string]int) int { return q["Field1"] },
			func(d map[string]int) int { return d["Field1"] },
		),
	}

	tree := Build(data, order)
	results := Search(tree, query)

	assert.Equal(t, 1, len(results))
	assert.Equal(t, data[0], results[0])
}

// 添加更多针对 tree.go 的单元测试
func TestBuildWithEmptyData(t *testing.T) {
	data := []map[string]int{}
	order := []NodeBuilder[map[string]int, map[string]int]{
		UniqueMapNode(
			func(q map[string]int) int { return q["Field1"] },
			func(d map[string]int) int { return d["Field1"] },
		),
	}
	tree := Build(data, order)
	results := Search(tree, map[string]int{"Field1": 1})
	assert.Equal(t, 0, len(results))
}

func TestSearchWithNoMatch(t *testing.T) {
	data := []map[string]int{
		{"Field1": 1, "Field2": 2},
		{"Field1": 3, "Field2": 4},
	}
	query := map[string]int{"Field1": 5}

	order := []NodeBuilder[map[string]int, map[string]int]{
		UniqueMapNode(
			func(q map[string]int) int { return q["Field1"] },
			func(d map[string]int) int { return d["Field1"] },
		),
	}
	tree := Build(data, order)
	results := Search(tree, query)
	assert.Equal(t, 0, len(results))
}

func TestUniqueSortedNodeBuildAndSearch(t *testing.T) {
	data := []map[string]int{
		{"Field1": 1, "Field2": 2},
		{"Field1": 3, "Field2": 4},
	}
	query := map[string]int{"Field1": 1}

	order := []NodeBuilder[map[string]int, map[string]int]{
		UniqueSortedNode(
			func(q map[string]int) int { return q["Field1"] },
			LE,
			func(d map[string]int) int { return d["Field1"] },
		),
	}

	tree := Build(data, order)
	results := Search(tree, query)

	assert.Equal(t, 1, len(results))
	assert.Equal(t, data[0], results[0])
}

func TestDataNodeBuilder(t *testing.T) {
	// 测试 Load 方法
	builder := DataNodeBuilder[int, string]{}
	loadedBuilder := builder.Load([]string{"A", "B"})
	assert.NotNil(t, loadedBuilder)
	assert.IsType(t, DataNodeBuilder[int, string]{}, loadedBuilder)

	// 测试 Push 方法
	result := builder.Push(nil)
	assert.Nil(t, result)

	// 测试 Node 方法
	node := builder.Node()
	assert.NotNil(t, node)
	assert.IsType(t, DataNode[int, string]{}, node)
}
