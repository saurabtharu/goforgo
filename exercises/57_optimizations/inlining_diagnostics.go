// inlining_diagnostics.go
// Understand function inlining and use runtime diagnostics in Go.
//
// This exercise covers two optimization concepts:
//
// 1. Function inlining: The Go compiler can inline small functions directly
//    into their callers, eliminating function call overhead. Functions that
//    are too complex (too many AST nodes) won't be inlined. Keeping hot-path
//    functions small and simple enables inlining.
//
// 2. Runtime diagnostics: Go provides runtime tools to inspect goroutine
//    counts, memory statistics, and CPU profiling. Understanding these
//    tools helps identify performance bottlenecks.
//
// Fix the code so the inlineable function is actually simple enough to
// inline, and add proper diagnostic output using runtime facilities.

package main

import (
	"fmt"
	"runtime"
)

// addSimple is a function that should be easily inlineable.
// Small leaf functions with no complex control flow get inlined by the compiler.
func addSimple(a, b int) int {
	return a + b
}

// processComplex is too complex to be inlined. It has too many operations
// and control flow branches for the compiler's inlining budget.
func processComplex(values []int) int {
	total := 0
	for _, v := range values {
		if v > 0 {
			total += v * 2
		} else if v < -10 {
			total -= v
		} else {
			total += v + 1
		}
	}
	if total > 1000 {
		total = 1000
	}
	return total
}

// multiplyInlineable should be a simple, inlineable function.
// TODO: This function is currently too complex to be inlined.
// Simplify it to just perform: a * b (remove all the unnecessary checks).
func multiplyInlineable(a, b int) int {
	result := 0
	if a == 0 || b == 0 {
		return 0
	}
	if a < 0 && b < 0 {
		result = (-a) * (-b)
	} else if a < 0 {
		result = -((-a) * b)
	} else if b < 0 {
		result = -(a * (-b))
	} else {
		result = a * b
	}
	if result > 1_000_000 {
		result = 1_000_000
	}
	return result
}

// benchmarkInlining runs inlineable vs non-inlineable functions in a hot loop.
func benchmarkInlining() {
	const iterations = 10_000_000

	// This should be fast due to inlining
	sum := 0
	for i := 0; i < iterations; i++ {
		sum += addSimple(i, i+1)
	}
	fmt.Printf("  addSimple result: %d\n", sum)

	// This should also be fast if multiplyInlineable is simplified
	product := 0
	for i := 1; i < iterations; i++ {
		product += multiplyInlineable(i, 2)
	}
	fmt.Printf("  multiplyInlineable result: %d\n", product)
}

// printDiagnostics should display runtime diagnostic information.
// TODO: Complete this function to show:
//   - Number of goroutines (use runtime.NumGoroutine())
//   - Number of CPUs available (use runtime.NumCPU())
//   - GOMAXPROCS value (use runtime.GOMAXPROCS(0) to read without changing)
//   - Memory statistics (use runtime.ReadMemStats)
//     Show: Alloc, TotalAlloc, Sys, and NumGC
func printDiagnostics() {
	fmt.Println("  Go version: ", runtime.Version())
	fmt.Println("  Compiler:   ", runtime.Compiler)
	fmt.Println("  GOOS/GOARCH:", runtime.GOOS+"/"+runtime.GOARCH)

	// FIX: Add goroutine count
	// fmt.Printf("  Goroutines:  %d\n", ???)

	// FIX: Add CPU count
	// fmt.Printf("  NumCPU:      %d\n", ???)

	// FIX: Add GOMAXPROCS (read current value without changing it)
	// fmt.Printf("  GOMAXPROCS:  %d\n", ???)

	// FIX: Add memory statistics
	// Use runtime.ReadMemStats(&memStats) and print:
	//   - Alloc: currently allocated heap bytes
	//   - TotalAlloc: cumulative bytes allocated
	//   - Sys: total bytes obtained from system
	//   - NumGC: number of completed GC cycles
}

// formatBytes converts bytes to a human-readable string.
// This is a helper for displaying memory stats.
func formatBytes(bytes uint64) string {
	const (
		kb = 1_024
		mb = 1_024 * kb
		gb = 1_024 * mb
	)

	switch {
	case bytes >= gb:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(gb))
	case bytes >= mb:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(mb))
	case bytes >= kb:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(kb))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func main() {
	fmt.Println("=== Function Inlining and Diagnostics ===")

	// Part 1: Inlining
	fmt.Println("\n--- Part 1: Function Inlining ---")
	fmt.Println("Functions that are small and simple get inlined by the compiler.")
	fmt.Println("You can check with: go build -gcflags='-m' to see inlining decisions.")
	fmt.Println()

	benchmarkInlining()

	// Demonstrate inlining concept
	fmt.Println()
	fmt.Println("  addSimple: small leaf function -> likely inlined")
	fmt.Println("  processComplex: complex control flow -> not inlined")
	fmt.Println("  multiplyInlineable: check if simplified -> should be inlined")

	// Verify multiplyInlineable works correctly
	if multiplyInlineable(3, 4) == 12 && multiplyInlineable(-2, 5) == -10 {
		fmt.Println("  multiplyInlineable: produces correct results")
	}

	// Part 2: Runtime diagnostics
	fmt.Println("\n--- Part 2: Runtime Diagnostics ---")
	printDiagnostics()

	// Force a GC to see updated stats
	runtime.GC()

	fmt.Println("\n  After GC:")
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	fmt.Printf("  HeapObjects: %d\n", memStats.HeapObjects)
	fmt.Printf("  GC cycles:   %d\n", memStats.NumGC)

	fmt.Println("\nDiagnostics complete!")
}
