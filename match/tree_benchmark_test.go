package tree

import (
	"fmt"
	"math/rand"
	"testing"

)

// --- Helper Functions for Benchmark Data Generation ---

// generateBenchmarkData generates a slice of maps representing the data for benchmarking.
// numFields specifies the number of fields in each data map.
// dataSize specifies the number of data items to generate.
func generateBenchmarkData(numFields, dataSize int) []map[string]int {
	data := make([]map[string]int, dataSize)
	for i := 0; i < dataSize; i++ {
		item := make(map[string]int)
		for j := 0; j < numFields; j++ {
			fieldName := fmt.Sprintf("Field%d", j)
			// Assign random values for simplicity, adjust range as needed
			item[fieldName] = rand.Intn(100)
		}
		data[i] = item
	}
	return data
}

// generateBenchmarkQuery generates a query map for benchmarking.
// numFields specifies the number of fields in the query map.
func generateBenchmarkQuery(numFields int) map[string]int {
	query := make(map[string]int)
	for j := 0; j < numFields; j++ {
		fieldName := fmt.Sprintf("Field%d", j)
		// Assign random values for simplicity, ensure it might match some data
		query[fieldName] = rand.Intn(100)
	}
	return query
}

// buildBenchmarkTree builds the matching tree for benchmark data.
// numFields specifies the number of fields to use for building the tree nodes.
func buildBenchmarkTree(data []map[string]int, numFields int) Node[map[string]int, map[string]int] {
	order := make([]NodeBuilder[map[string]int, map[string]int], 0, numFields+1)
	for i := 0; i < numFields; i++ {
		fieldName := fmt.Sprintf("Field%d", i)
		// Use MapNode for exact match on each field
		node := UniqueMapNode(
			func(q map[string]int) int { return q[fieldName] },
			func(d map[string]int) int { return d[fieldName] },
		)
		order = append(order, node)
	}
	order = append(order, DataNodeBuilder[map[string]int, map[string]int]{})
	return Build(data, order)
}

// --- Benchmark Functions ---

// runBenchmark performs the benchmark for a given number of fields and data size.
func runBenchmark(b *testing.B, numFields, dataSize int) {
	data := generateBenchmarkData(numFields, dataSize)
	root := buildBenchmarkTree(data, numFields)
	query := generateBenchmarkQuery(numFields)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Search(root, query)
	}
}

// runBenchmarkWithBuild performs the benchmark including the tree build time.
func runBenchmarkWithBuild(b *testing.B, numFields, dataSize int) {
	data := generateBenchmarkData(numFields, dataSize)
	query := generateBenchmarkQuery(numFields)

	b.ResetTimer() // Reset timer *before* the loop
	for i := 0; i < b.N; i++ {
		// Build the tree inside the loop
		root := buildBenchmarkTree(data, numFields)
		_ = Search(root, query)
	}
}

// --- Benchmarks for Varying Field Counts (DataSize = 1000) ---

func BenchmarkTreeSearch_Fields5_Data1000(b *testing.B) {
	runBenchmark(b, 5, 1000)
}

func BenchmarkLinearSearch_Fields5_Data1000(b *testing.B) {
	runLinearBenchmark(b, 5, 1000)
}
func BenchmarkTreeSearch_Fields10_Data1000(b *testing.B) {
	runBenchmark(b, 10, 1000)
}

func BenchmarkLinearSearch_Fields10_Data1000(b *testing.B) {
	runLinearBenchmark(b, 10, 1000)
}
func BenchmarkTreeSearch_Fields20_Data1000(b *testing.B) {
	runBenchmark(b, 20, 1000)
}

func BenchmarkLinearSearch_Fields20_Data1000(b *testing.B) {
	runLinearBenchmark(b, 20, 1000)
}

// --- Benchmarks for Varying Data Sizes (NumFields = 10) ---

func BenchmarkTreeSearch_Fields10_Data10(b *testing.B) {
	runBenchmark(b, 10, 10)
}

func BenchmarkLinearSearch_Fields10_Data10(b *testing.B) {
	runLinearBenchmark(b, 10, 10)
}

func BenchmarkTreeSearch_Fields10_Data100(b *testing.B) {
	runBenchmark(b, 10, 100)
}

func BenchmarkLinearSearch_Fields10_Data100(b *testing.B) {
	runLinearBenchmark(b, 10, 100)
}

// Optional: Benchmark with larger data size
func BenchmarkTreeSearch_Fields10_Data10000(b *testing.B) {
	runBenchmark(b, 10, 10000)
}

func BenchmarkLinearSearch_Fields10_Data10000(b *testing.B) {
	runLinearBenchmark(b, 10, 10000)
}

// --- Linear Search Benchmark for Comparison ---

// linearSearch performs a simple linear scan for matching data.
func linearSearch(data []map[string]int, query map[string]int) []map[string]int {
	results := make([]map[string]int, 0)
	for _, item := range data {
		match := true
		for key, queryValue := range query {
			if itemValue, ok := item[key]; !ok || itemValue != queryValue {
				match = false
				break
			}
		}
		if match {
			results = append(results, item)
		}
	}
	return results
}

// runLinearBenchmark performs the benchmark for linear search.
func runLinearBenchmark(b *testing.B, numFields, dataSize int) {
	data := generateBenchmarkData(numFields, dataSize)
	query := generateBenchmarkQuery(numFields)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = linearSearch(data, query)
	}
}

// --- Linear Benchmarks for Varying Field Counts (DataSize = 1000) ---

// --- Benchmarks Including Tree Build Time (NumFields = 10) ---

func BenchmarkTreeBuildAndSearch_Fields10_Data100(b *testing.B) {
	runBenchmarkWithBuild(b, 10, 100)
}

func BenchmarkTreeBuildAndSearch_Fields10_Data1000(b *testing.B) {
	runBenchmarkWithBuild(b, 10, 1000)
}

// --- Linear Search Benchmark for Comparison ---
