// benchmarking_test.go
// Mistake #89: Fix these inaccurate benchmarks!
//
// Problems to fix:
// 1. Results not stored - compiler may optimize the function call away
// 2. b.ResetTimer not called after expensive setup
// 3. Observer effect - benchmark measures more than intended
//
// Fix each benchmark to produce accurate, reliable results.
// After fixing all benchmarks, delete TestFixBenchmarks at the bottom.

package main

import "testing"

// BUG #1: Compiler optimization - the result of SortSlice is not stored.
// The Go compiler can detect that the return value is never used and may
// optimize away the entire function call. Store the result in a package-level
// variable to prevent this.
//
// Pattern to fix:
//   var result []int  // package-level variable
//   func BenchmarkX(b *testing.B) {
//       var r []int
//       for i := 0; i < b.N; i++ { r = SortSlice(input) }
//       result = r  // prevent compiler optimization
//   }

func BenchmarkSortSlice(b *testing.B) {
	input := GenerateSlice(1_000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// BUG: Return value not captured - compiler may skip the call entirely
		SortSlice(input)
	}
}

// BUG #2: No b.ResetTimer after expensive setup.
// GenerateSlice(100_000) takes non-trivial time, and that setup time
// is included in the benchmark measurement. Call b.ResetTimer()
// after setup to exclude it.

func BenchmarkSumSlice(b *testing.B) {
	// Expensive setup - this time should NOT be measured
	data := GenerateSlice(100_000)

	// BUG: Missing b.ResetTimer() here!
	// The setup time above is incorrectly included in the benchmark.

	for i := 0; i < b.N; i++ {
		SumSlice(data)
	}
}

// BUG #3: Observer effect - the benchmark is measuring both string
// building AND slice generation inside the loop. The slice generation
// is not what we want to benchmark. Move setup outside the loop.

func BenchmarkBuildString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// BUG: Creating the input slice on every iteration.
		// This measures both slice creation AND string building.
		parts := make([]string, 100)
		for j := range parts {
			parts[j] = "word"
		}
		BuildString(parts)
	}
}

// BUG #4: Same observer effect AND missing result capture.

func BenchmarkConcatStrings(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// BUG: Creating input inside the loop measures the wrong thing
		parts := make([]string, 100)
		for j := range parts {
			parts[j] = "word"
		}
		// BUG: Result not captured
		ConcatStrings(parts)
	}
}

// BUG #5: Benchmark with sub-benchmarks has ALL the issues above.

func BenchmarkStringMethods(b *testing.B) {
	sizes := []int{10, 100, 1_000}

	for _, size := range sizes {
		// BUG: Sub-benchmark name doesn't include size
		b.Run("concat", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// BUG: Setup inside loop (observer effect)
				parts := make([]string, size)
				for j := range parts {
					parts[j] = "x"
				}
				ConcatStrings(parts)
			}
		})

		b.Run("builder", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// BUG: Setup inside loop (observer effect)
				parts := make([]string, size)
				for j := range parts {
					parts[j] = "x"
				}
				BuildString(parts)
			}
		})
	}
}

// This test forces the exercise to fail until you fix all benchmarks.
// Delete this function after:
// 1. Storing benchmark results in package-level variables
// 2. Adding b.ResetTimer() after expensive setup
// 3. Moving input creation outside benchmark loops
// 4. Including size in sub-benchmark names
func TestFixBenchmarks(t *testing.T) {
	t.Fatal("EXERCISE: Fix all benchmarks above for accuracy. " +
		"Store results to prevent compiler optimization, use b.ResetTimer() " +
		"after setup, and move input creation outside b.N loops. " +
		"Then delete this function.")
}
