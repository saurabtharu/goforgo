// inlining_diagnostics.go - SOLUTION
// Functions simplified for inlining, full runtime diagnostics added.

package main

import (
	"fmt"
	"runtime"
)

func addSimple(a, b int) int {
	return a + b
}

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

// multiplyInlineable simplified to a single expression.
// Fixed: Go's * operator handles all sign combinations correctly.
func multiplyInlineable(a, b int) int {
	return a * b
}

func benchmarkInlining() {
	const iterations = 10_000_000

	sum := 0
	for i := 0; i < iterations; i++ {
		sum += addSimple(i, i+1)
	}
	fmt.Printf("  addSimple result: %d\n", sum)

	product := 0
	for i := 1; i < iterations; i++ {
		product += multiplyInlineable(i, 2)
	}
	fmt.Printf("  multiplyInlineable result: %d\n", product)
}

// printDiagnostics displays comprehensive runtime information.
// Fixed: all diagnostic outputs are populated.
func printDiagnostics() {
	fmt.Println("  Go version: ", runtime.Version())
	fmt.Println("  Compiler:   ", runtime.Compiler)
	fmt.Println("  GOOS/GOARCH:", runtime.GOOS+"/"+runtime.GOARCH)

	// Fixed: goroutine count
	fmt.Printf("  Goroutines:  %d\n", runtime.NumGoroutine())

	// Fixed: CPU count
	fmt.Printf("  NumCPU:      %d\n", runtime.NumCPU())

	// Fixed: GOMAXPROCS (0 reads without changing)
	fmt.Printf("  GOMAXPROCS:  %d\n", runtime.GOMAXPROCS(0))

	// Fixed: memory statistics
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	fmt.Printf("  Alloc:       %s\n", formatBytes(memStats.Alloc))
	fmt.Printf("  TotalAlloc:  %s\n", formatBytes(memStats.TotalAlloc))
	fmt.Printf("  Sys:         %s\n", formatBytes(memStats.Sys))
	fmt.Printf("  NumGC:       %d\n", memStats.NumGC)
}

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

	runtime.GC()

	fmt.Println("\n  After GC:")
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	fmt.Printf("  HeapObjects: %d\n", memStats.HeapObjects)
	fmt.Printf("  GC cycles:   %d\n", memStats.NumGC)

	fmt.Println("\nDiagnostics complete!")
}
