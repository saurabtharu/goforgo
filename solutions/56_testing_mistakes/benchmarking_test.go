// benchmarking_test.go
// Solution for Mistake #89
//
// Fixed:
// 1. Results stored in package-level variables to prevent compiler optimization
// 2. b.ResetTimer called after expensive setup
// 3. Observer effects eliminated by moving setup outside benchmark loops
// 4. Sub-benchmark names include size for meaningful comparison

package main

import (
	"fmt"
	"testing"
)

// Package-level variables to prevent compiler from optimizing away results.
// The compiler can detect unused return values and skip function calls entirely.
var (
	benchResultSort   []int
	benchResultSum    int
	benchResultString string
)

// FIXED: Result stored in package-level variable to prevent compiler optimization.
func BenchmarkSortSlice(b *testing.B) {
	input := GenerateSlice(1_000)

	var r []int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = SortSlice(input)
	}
	benchResultSort = r
}

// FIXED: b.ResetTimer() called after expensive setup.
func BenchmarkSumSlice(b *testing.B) {
	// Expensive setup - excluded from measurement
	data := GenerateSlice(100_000)

	// FIXED: Reset timer so setup time is excluded
	b.ResetTimer()

	var r int
	for i := 0; i < b.N; i++ {
		r = SumSlice(data)
	}
	benchResultSum = r
}

// FIXED: Input created outside the loop to avoid observer effect.
// Only BuildString is being measured now.
func BenchmarkBuildString(b *testing.B) {
	parts := make([]string, 100)
	for j := range parts {
		parts[j] = "word"
	}

	var r string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = BuildString(parts)
	}
	benchResultString = r
}

// FIXED: Input created outside loop AND result captured.
func BenchmarkConcatStrings(b *testing.B) {
	parts := make([]string, 100)
	for j := range parts {
		parts[j] = "word"
	}

	var r string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = ConcatStrings(parts)
	}
	benchResultString = r
}

// FIXED: Sub-benchmarks with proper names, setup outside loops,
// and result capture.
func BenchmarkStringMethods(b *testing.B) {
	sizes := []int{10, 100, 1_000}

	for _, size := range sizes {
		// FIXED: Sub-benchmark name includes size for meaningful comparison
		b.Run(fmt.Sprintf("concat/size=%d", size), func(b *testing.B) {
			// FIXED: Setup outside the loop
			parts := make([]string, size)
			for j := range parts {
				parts[j] = "x"
			}

			var r string
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				r = ConcatStrings(parts)
			}
			benchResultString = r
		})

		b.Run(fmt.Sprintf("builder/size=%d", size), func(b *testing.B) {
			// FIXED: Setup outside the loop
			parts := make([]string, size)
			for j := range parts {
				parts[j] = "x"
			}

			var r string
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				r = BuildString(parts)
			}
			benchResultString = r
		})
	}
}
