// benchmarking.go
// Solution for Mistake #89
//
// The source code is the same - the fixes are in the benchmark file.

package main

import (
	"fmt"
	"sort"
	"strings"
)

// SortSlice sorts a copy of the input slice and returns it.
func SortSlice(input []int) []int {
	result := make([]int, len(input))
	copy(result, input)
	sort.Ints(result)
	return result
}

// ConcatStrings joins strings using + operator.
func ConcatStrings(parts []string) string {
	result := ""
	for _, p := range parts {
		result += p
	}
	return result
}

// BuildString joins strings using strings.Builder.
func BuildString(parts []string) string {
	var b strings.Builder
	for _, p := range parts {
		b.WriteString(p)
	}
	return b.String()
}

// SumSlice returns the sum of all elements in a slice.
func SumSlice(nums []int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// GenerateSlice creates a slice of n integers for benchmarking setup.
func GenerateSlice(n int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = n - i
	}
	return s
}

func main() {
	fmt.Println("Run 'go test -bench=. -v' to execute the benchmarks.")
}
